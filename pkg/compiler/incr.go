package compiler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-interpreter/gi/pkg/ast"
	"github.com/go-interpreter/gi/pkg/constant"
	"github.com/go-interpreter/gi/pkg/token"
	"github.com/go-interpreter/gi/pkg/types"
	"sort"
	"strings"

	"github.com/go-interpreter/gi/pkg/compiler/analysis"
	"github.com/neelance/astrewrite"
	"golang.org/x/tools/go/gcimporter15"
)

func IncrementallyCompile(a *Archive, importPath string, files []*ast.File, fileSet *token.FileSet, importContext *ImportContext, minify bool) (*Archive, error) {

	pp("jea debug, top of incrementallyCompile(): here is what files has:")
	j := 0
	for _, file := range files {
		for _, decl := range file.Nodes {
			pp("decl[%v] = '%#v'", j, decl)
		}
		j++
	}

	var newCodeText [][]byte

	var typesInfo *types.Info
	if a == nil {
		typesInfo = &types.Info{
			Types:      make(map[ast.Expr]types.TypeAndValue),
			Defs:       make(map[*ast.Ident]types.Object),
			Uses:       make(map[*ast.Ident]types.Object),
			Implicits:  make(map[ast.Node]types.Object),
			Selections: make(map[*ast.SelectorExpr]*types.Selection),
			Scopes:     make(map[ast.Node]*types.Scope),
		}
	} else {
		typesInfo = a.typesInfo
	}

	var importError error
	var errList ErrorList
	var previousErr error
	var config *types.Config
	if a != nil {
		config = a.config
	} else {
		config = &types.Config{
			DisableUnusedImportCheck: true, // jea add
			Importer: packageImporter{
				importContext: importContext,
				importError:   &importError,
			},
			Sizes: sizes32,
			Error: func(err error) {
				panic(fmt.Sprintf("where error? err = '%v'", err))
				if previousErr != nil && previousErr.Error() == err.Error() {
					return
				}
				errList = append(errList, err)
				previousErr = err
			},
		}
	}
	pp("about to call config.Check")
	var pkg *types.Package
	var check *types.Checker
	if a != nil {
		pkg = a.pkg
		check = a.check
	}
	var err error
	pkg, check, err = config.Check(pkg, check, importPath, fileSet, files, typesInfo)
	if importError != nil {
		//pp("config.Check: importError")
		return nil, importError
	}
	if errList != nil {
		//pp("config.Check: errList is not nil")

		if len(errList) > 10 {
			pos := token.NoPos
			if last, ok := errList[9].(types.Error); ok {
				pos = last.Pos
			}
			errList = append(errList[:10], types.Error{Fset: fileSet, Pos: pos, Msg: "too many errors"})
		}
		return nil, errList
	}
	if err != nil {
		return nil, err
		//pp("config.Check err = '%v'", err)
	}
	//pp("got past config.Check")

	importContext.Packages[importPath] = pkg

	exportData := gcimporter.BExportData(nil, pkg)
	encodedFileSet := bytes.NewBuffer(nil)
	if err := fileSet.Write(json.NewEncoder(encodedFileSet).Encode); err != nil {
		return nil, err
	}

	simplifiedFiles := make([]*ast.File, len(files))
	for i, file := range files {
		simplifiedFiles[i] = astrewrite.Simplify(file, typesInfo, false)
	}

	isBlocking := func(f *types.Func) bool {
		archive, err := importContext.Import(f.Pkg().Path())
		if err != nil {
			panic(err)
		}
		fullName := f.FullName()
		for _, d := range archive.Declarations {
			if string(d.FullName) == fullName {
				return d.Blocking
			}
		}
		panic(fullName)
	}
	//pp("about to call AnalyzePkg")
	pkgInfo := analysis.AnalyzePkg(simplifiedFiles, fileSet, typesInfo, pkg, isBlocking)
	c := &funcContext{
		FuncInfo: pkgInfo.InitFuncInfo,
		p: &pkgContext{
			Info:                 pkgInfo,
			additionalSelections: make(map[*ast.SelectorExpr]selection),

			pkgVars:      make(map[string]string),
			objectNames:  make(map[types.Object]string),
			varPtrNames:  make(map[*types.Var]string),
			escapingVars: make(map[*types.Var]bool),
			indentation:  1,
			dependencies: make(map[types.Object]bool),
			minify:       minify,
			fileSet:      fileSet,
		},
		allVars:     make(map[string]int),
		flowDatas:   map[*types.Label]*flowData{nil: {}},
		caseCounter: 1,
		labelCases:  make(map[*types.Label]int),
	}
	for name := range reservedKeywords {
		c.allVars[name] = 1
	}
	pp("got past AnalyzePkg")

	// ==============================
	// actual incrmental compilation
	// modifications start here
	// ==============================

	// imports
	var importDecls []*Decl
	var importedPaths []string
	for _, importedPkg := range pkg.Imports() {
		if importedPkg == types.Unsafe {
			// Prior to Go 1.9, unsafe import was excluded by Imports() method,
			// but now we do it here to maintain previous behavior.
			continue
		}
		c.p.pkgVars[importedPkg.Path()] = c.newVariableWithLevel(importedPkg.Name(), true)
		importedPaths = append(importedPaths, importedPkg.Path())
	}
	sort.Strings(importedPaths)
	for _, impPath := range importedPaths {
		id := c.newIdent(fmt.Sprintf(`%s.$init`, c.p.pkgVars[impPath]), types.NewSignature(nil, nil, nil, false))
		call := &ast.CallExpr{Fun: id}
		c.Blocking[call] = true
		c.Flattened[call] = true
		importDecls = append(importDecls, &Decl{
			Vars:     []string{c.p.pkgVars[impPath]},
			DeclCode: []byte(fmt.Sprintf("\t%s = $packages[\"%s\"];\n", c.p.pkgVars[impPath], impPath)),
			InitCode: c.CatchOutput(1, func() { c.translateStmt(&ast.ExprStmt{X: call}, nil) }),
		})
	}

	collectDependencies := func(f func()) []string {
		c.p.dependencies = make(map[types.Object]bool)
		f()
		var deps []string
		for o := range c.p.dependencies {
			qualifiedName := o.Pkg().Path() + "." + o.Name()
			if f, ok := o.(*types.Func); ok && f.Type().(*types.Signature).Recv() != nil {
				deps = append(deps, qualifiedName+"~")
				continue
			}
			deps = append(deps, qualifiedName)
		}
		sort.Strings(deps)
		return deps
	}

	varsWithInit := make(map[*types.Var]bool)
	for _, init := range c.p.InitOrder {
		for _, o := range init.Lhs {
			varsWithInit[o] = true
		}
	}

	var functions []*ast.FuncDecl
	var vars []*types.Var
	var varDecls []*Decl
	var funcDecls []*Decl
	var mainFunc *types.Func

	for _, file := range simplifiedFiles {
		pp("file.Nodes has %v elements", len(file.Nodes))
		for _, decl := range file.Nodes {

			// fill out vars and functions

			switch d := decl.(type) {
			case *ast.FuncDecl:
				pp("next decl from file.Nodes is a funcDecl: '%#v'", d)
				sig := c.p.Defs[d.Name].(*types.Func).Type().(*types.Signature)
				var recvType types.Type
				if sig.Recv() != nil {
					recvType = sig.Recv().Type()
					if ptr, isPtr := recvType.(*types.Pointer); isPtr {
						recvType = ptr.Elem()
					}
				}
				if sig.Recv() == nil {
					c.objectName(c.p.Defs[d.Name].(*types.Func)) // register toplevel name
				}
				if !isBlank(d.Name) {
					fun := d
					functions = append(functions, fun)

					// jea: codegen here and now, in order.
					o := c.p.Defs[fun.Name].(*types.Func)
					funcInfo := c.p.FuncDeclInfos[o]
					de := Decl{
						FullName: o.FullName(),
						Blocking: len(funcInfo.Blocking) != 0,
					}
					if fun.Recv == nil {
						de.Vars = []string{c.objectName(o)}
						de.DceObjectFilter = o.Name()
						switch o.Name() {
						case "main":
							mainFunc = o
							de.DceObjectFilter = ""
						case "init":
							de.InitCode = c.CatchOutput(1, func() {
								id := c.newIdent("", types.NewSignature(nil, nil, nil, false))
								c.p.Uses[id] = o
								call := &ast.CallExpr{Fun: id}
								if len(c.p.FuncDeclInfos[o].Blocking) != 0 {
									c.Blocking[call] = true
								}
								c.translateStmt(&ast.ExprStmt{X: call}, nil)
							})
							de.DceObjectFilter = ""
						}
					}
					if fun.Recv != nil {
						recvType := o.Type().(*types.Signature).Recv().Type()
						ptr, isPointer := recvType.(*types.Pointer)
						namedRecvType, _ := recvType.(*types.Named)
						if isPointer {
							namedRecvType = ptr.Elem().(*types.Named)
						}
						de.DceObjectFilter = namedRecvType.Obj().Name()
						if !fun.Name.IsExported() {
							de.DceMethodFilter = o.Name() + "~"
						}
					}

					de.DceDeps = collectDependencies(func() {
						de.DeclCode = c.translateToplevelFunction(fun, funcInfo)
					})
					funcDecls = append(funcDecls, &de)
					pp("place3, appending to newCodeText: de.DeclCode='%s'", string(de.DeclCode))
					newCodeText = append(newCodeText, de.DeclCode)

					// end of function codegen now
				}
			case *ast.GenDecl:
				pp("next decl from file.Nodes is a GenDecl: '%#v'", d.Specs[0].(*ast.ValueSpec).Names[0])
				switch d.Tok {
				case token.TYPE:
					for _, spec := range d.Specs {
						o := c.p.Defs[spec.(*ast.TypeSpec).Name].(*types.TypeName)
						c.p.typeNames = append(c.p.typeNames, o)
						c.objectName(o) // register toplevel name
					}
				case token.VAR:
					for _, spec := range d.Specs {
						for _, name := range spec.(*ast.ValueSpec).Names {
							if !isBlank(name) {
								o := c.p.Defs[name].(*types.Var)
								vars = append(vars, o)
								c.objectName(o) // register toplevel name

								// jea: codegen here and now, in order.

								var de Decl
								if !o.Exported() {
									de.Vars = []string{c.objectName(o)}
								}
								if c.p.HasPointer[o] && !o.Exported() {
									de.Vars = append(de.Vars, c.varPtrName(o))
								}
								if _, ok := varsWithInit[o]; !ok {
									de.DceDeps = collectDependencies(func() {
										de.InitCode = []byte(fmt.Sprintf("\t\t%s = %s;\n", c.objectName(o), c.translateExpr(c.zeroValue(o.Type())).String()))
									})
								}
								de.DceObjectFilter = o.Name()
								varDecls = append(varDecls, &de)
								pp("place 1, appending to newCodeText: de.InitCode='%s'", string(de.InitCode))
								newCodeText = append(newCodeText, de.InitCode)

								// end codegen here and now for vars
							}
						}
					}
				case token.CONST:
					// skip, constants are inlined
				}
			default:
				pp("next decl from file.Nodes is an unknown/default type: '%#v'", decl)
				c.output = nil
				switch s := decl.(type) {
				case ast.Stmt:
					c.translateStmt(s, nil)
					pp("in codegen, %T/val='%#v'", s, s)
				default:
					pp("in codegen, unknown type %T", s)
					continue
				}

				_, isExprStmt := decl.(*ast.ExprStmt)
				if !isExprStmt {
					newCodeText = append(newCodeText, c.output)
				} else {

					n := len(c.output)
					var ele string
					if bytes.HasSuffix(c.output, []byte(";\n")) {
						ele = string(bytes.TrimLeft(c.output[:n-2], " \t"))
					} else {
						ele = string(c.output)
					}
					tmp := fmt.Sprintf(`if "table" == type(%[1]s) then print(%[1]s); print("key", "value"); `+
						`for i, v in pairs(%[1]s) do print(i, v); end else print(%[1]s) end `,
						ele)
					newCodeText = append(newCodeText, []byte(tmp))
				}
				pp("place5, appending to newCodeText: c.output='%s'", string(c.output))
				c.output = nil
			}
		}
	}

	// ===========================
	// variables
	// ===========================

	pp("jea, at variables, in package.go:336. vars='%#v'", vars)
	//var varDecls []*Decl
	// moved up above, to stay in sequence entered order.
	//	for _, o := range vars {
	//	}
	pp("jea, package.go:362. c.p.InitOrder='%#v'", c.p.InitOrder)
	for _, init := range c.p.InitOrder {
		lhs := make([]ast.Expr, len(init.Lhs))
		for i, o := range init.Lhs {
			ident := ast.NewIdent(o.Name())
			c.p.Defs[ident] = o
			lhs[i] = c.setType(ident, o.Type())
			varsWithInit[o] = true
		}
		var d Decl
		d.DceDeps = collectDependencies(func() {
			c.localVars = nil
			d.InitCode = c.CatchOutput(1, func() {
				c.translateStmt(&ast.AssignStmt{
					Lhs: lhs,
					Tok: token.DEFINE,
					Rhs: []ast.Expr{init.Rhs},
				}, nil)
			})
			d.Vars = append(d.Vars, c.localVars...)
		})
		if len(init.Lhs) == 1 {
			if !analysis.HasSideEffect(init.Rhs, c.p.Info.Info) {
				d.DceObjectFilter = init.Lhs[0].Name()
			}
		}
		varDecls = append(varDecls, &d)
		pp("place2, appending to newCodeText: d.InitCode='%s'", string(d.InitCode))
		newCodeText = append(newCodeText, d.InitCode)
	}

	pp("jea, functions in package.go:393")

	// ===========================
	// functions
	// ===========================
	// moved up to stay in order.
	//	for _, fun := range functions {
	//	}
	if pkg.Name() == "main" {
		if mainFunc == nil {
			return nil, fmt.Errorf("missing main function")
		}
		id := c.newIdent("", types.NewSignature(nil, nil, nil, false))
		c.p.Uses[id] = mainFunc
		call := &ast.CallExpr{Fun: id}
		ifStmt := &ast.IfStmt{
			Cond: c.newIdent("$pkg === $mainPkg", types.Typ[types.Bool]),
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{X: call},
					&ast.AssignStmt{
						Lhs: []ast.Expr{c.newIdent("$mainFinished", types.Typ[types.Bool])},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{c.newConst(types.Typ[types.Bool], constant.MakeBool(true))},
					},
				},
			},
		}
		if len(c.p.FuncDeclInfos[mainFunc].Blocking) != 0 {
			c.Blocking[call] = true
			c.Flattened[ifStmt] = true
		}
		funcDecls = append(funcDecls, &Decl{
			InitCode: c.CatchOutput(1, func() {
				c.translateStmt(ifStmt, nil)
			}),
		})
	}

	// named types
	var typeDecls []*Decl
	for _, o := range c.p.typeNames {
		if o.IsAlias() {
			continue
		}
		typeName := c.objectName(o)
		d := Decl{
			Vars:            []string{typeName},
			DceObjectFilter: o.Name(),
		}
		d.DceDeps = collectDependencies(func() {
			d.DeclCode = c.CatchOutput(0, func() {
				typeName := c.objectName(o)
				lhs := typeName
				if isPkgLevel(o) {
					lhs += " = $pkg." + encodeIdent(o.Name())
				}
				size := int64(0)
				constructor := "null"
				switch t := o.Type().Underlying().(type) {
				case *types.Struct:
					params := make([]string, t.NumFields())
					for i := 0; i < t.NumFields(); i++ {
						params[i] = fieldName(t, i) + "_"
					}
					constructor = fmt.Sprintf("function(%s) {\n\t\tthis.$val = this;\n\t\tif (arguments.length === 0) {\n", strings.Join(params, ", "))
					for i := 0; i < t.NumFields(); i++ {
						constructor += fmt.Sprintf("\t\t\tthis.%s = %s;\n", fieldName(t, i), c.translateExpr(c.zeroValue(t.Field(i).Type())).String())
					}
					constructor += "\t\t\treturn;\n\t\t}\n"
					for i := 0; i < t.NumFields(); i++ {
						constructor += fmt.Sprintf("\t\tthis.%[1]s = %[1]s_;\n", fieldName(t, i))
					}
					constructor += "\t}"
				case *types.Basic, *types.Array, *types.Slice, *types.Chan, *types.Signature, *types.Interface, *types.Pointer, *types.Map:
					size = sizes32.Sizeof(t)
				}
				c.Printf(`%s = $newType(%d, %s, "%s.%s", %t, "%s", %t, %s);`, lhs, size, typeKind(o.Type()), o.Pkg().Name(), o.Name(), o.Name() != "", o.Pkg().Path(), o.Exported(), constructor)
			})
			d.MethodListCode = c.CatchOutput(0, func() {
				named := o.Type().(*types.Named)
				if _, ok := named.Underlying().(*types.Interface); ok {
					return
				}
				var methods []string
				var ptrMethods []string
				for i := 0; i < named.NumMethods(); i++ {
					method := named.Method(i)
					name := method.Name()
					if reservedKeywords[name] {
						name += "$"
					}
					pkgPath := ""
					if !method.Exported() {
						pkgPath = method.Pkg().Path()
					}
					t := method.Type().(*types.Signature)
					entry := fmt.Sprintf(`{prop: "%s", name: "%s", pkg: "%s", typ: $funcType(%s)}`, name, method.Name(), pkgPath, c.initArgs(t))
					if _, isPtr := t.Recv().Type().(*types.Pointer); isPtr {
						ptrMethods = append(ptrMethods, entry)
						continue
					}
					methods = append(methods, entry)
				}
				if len(methods) > 0 {
					c.Printf("%s.methods = [%s];", c.typeName(named), strings.Join(methods, ", "))
				}
				if len(ptrMethods) > 0 {
					c.Printf("%s.methods = [%s];", c.typeName(types.NewPointer(named)), strings.Join(ptrMethods, ", "))
				}
			})
			switch t := o.Type().Underlying().(type) {
			case *types.Array, *types.Chan, *types.Interface, *types.Map, *types.Pointer, *types.Slice, *types.Signature, *types.Struct:
				d.TypeInitCode = c.CatchOutput(0, func() {
					c.Printf("%s.init(%s);", c.objectName(o), c.initArgs(t))
				})
			}
		})
		typeDecls = append(typeDecls, &d)
	}

	// anonymous types
	for _, t := range c.p.anonTypes {
		d := Decl{
			Vars:            []string{t.Name()},
			DceObjectFilter: t.Name(),
		}
		d.DceDeps = collectDependencies(func() {
			d.DeclCode = []byte(fmt.Sprintf("\t%s = $%sType(%s);\n", t.Name(), strings.ToLower(typeKind(t.Type())[5:]), c.initArgs(t.Type())))
		})
		typeDecls = append(typeDecls, &d)
	}

	var allDecls []*Decl
	for _, d := range append(append(append(importDecls, typeDecls...), varDecls...), funcDecls...) {
		d.DeclCode = removeWhitespace(d.DeclCode, minify)
		d.MethodListCode = removeWhitespace(d.MethodListCode, minify)
		d.TypeInitCode = removeWhitespace(d.TypeInitCode, minify)
		d.InitCode = removeWhitespace(d.InitCode, minify)
		allDecls = append(allDecls, d)
	}

	// raw top level statements! now integrated above in order

	if len(c.p.errList) != 0 {
		return nil, c.p.errList
	}

	if a == nil {
		return &Archive{
			ImportPath:   importPath,
			Name:         pkg.Name(),
			Imports:      importedPaths,
			ExportData:   exportData,
			Declarations: allDecls,
			FileSet:      encodedFileSet.Bytes(),
			Minified:     minify,
			NewCodeText:  newCodeText,
			typesInfo:    typesInfo,
			config:       config,
			pkg:          pkg,
			check:        check,
		}, nil
	} else {
		a.pkg = pkg
		a.check = check
		a.NewCodeText = newCodeText
	}
	return a, nil
}

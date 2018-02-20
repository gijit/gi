package compiler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gijit/gi/pkg/ast"
	"github.com/gijit/gi/pkg/constant"
	"github.com/gijit/gi/pkg/token"
	"github.com/gijit/gi/pkg/types"
	"sort"
	"strings"

	"github.com/gijit/gi/pkg/compiler/analysis"
	"github.com/neelance/astrewrite"
	"golang.org/x/tools/go/gcimporter15"
)

// start with GopherJS Compile again, and do
// as minimal an adaption to LuaJIT as possible,
// in order to preserve the ability to compile
// entire packages. For incremental interactivity,
// see the incr.go:34 IncrementallyCompile() function.
//
func FullPackageCompile(importPath string, files []*ast.File, fileSet *token.FileSet, importContext *ImportContext, minify bool) (*Archive, error) {
	typesInfo := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}

	var importError error
	var errList ErrorList
	var previousErr error
	config := &types.Config{
		Importer: packageImporter{
			importContext: importContext,
			importError:   &importError,
		},
		Sizes: sizes64,
		Error: func(err error) {
			if previousErr != nil && previousErr.Error() == err.Error() {
				return
			}
			errList = append(errList, err)
			previousErr = err
		},
	}

	typesPkg, _, err := config.Check(nil, nil, importPath, fileSet, files, typesInfo, addPreludeToNewPkg)
	if importError != nil {
		return nil, importError
	}
	if errList != nil {
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
	}
	importContext.Packages[importPath] = typesPkg

	exportData := gcimporter.BExportData(nil, typesPkg)
	encodedFileSet := bytes.NewBuffer(nil)
	if err := fileSet.Write(json.NewEncoder(encodedFileSet).Encode); err != nil {
		return nil, err
	}

	simplifiedFiles := make([]*ast.File, len(files))
	for i, file := range files {
		simplifiedFiles[i] = astrewrite.Simplify(file, typesInfo, false)
	}

	isBlocking := func(f *types.Func) bool {
		return false

		// jea: avoid this, it imports shadow/math/rand by full path
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
	pkgInfo := analysis.AnalyzePkg(simplifiedFiles, fileSet, typesInfo, typesPkg, isBlocking)
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
			files:        files,
		},
		allVars:     make(map[string]int),
		flowDatas:   map[*types.Label]*flowData{nil: {}},
		caseCounter: 1,
		labelCases:  make(map[*types.Label]int),
	}
	for name := range reservedKeywords {
		c.allVars[name] = 1
	}

	// imports
	var importDecls []*Decl
	var importedPaths []string
	for _, importedPkg := range typesPkg.Imports() {
		if importedPkg == types.Unsafe {
			// Prior to Go 1.9, unsafe import was excluded by Imports() method,
			// but now we do it here to maintain previous behavior.
			continue
		}
		c.p.pkgVars[importedPkg.Path()] = c.newVariableWithLevel(importedPkg.Name(), true)
		fmt.Printf("importedPkg.Path() = '%s'; importedPkg='%#v'\n", importedPkg.Path(), importedPkg)
		importedPaths = append(importedPaths, importedPkg.Path())
	}
	sort.Strings(importedPaths)
	for _, impPath := range importedPaths {
		id := c.newIdent(fmt.Sprintf(`%s._init`, c.p.pkgVars[impPath]), types.NewSignature(nil, nil, nil, false))
		call := &ast.CallExpr{Fun: id}
		c.Blocking[call] = true
		c.Flattened[call] = true
		importDecls = append(importDecls, &Decl{
			Vars:     []string{c.p.pkgVars[impPath]},
			DeclCode: []byte(fmt.Sprintf("\t%s = __packages[\"%s\"];\n\t__go_import(\"%s\");\n", c.p.pkgVars[impPath], impPath, impPath)),
			//DeclCode: []byte(fmt.Sprintf("\t%s = __packages[\"%s\"];\n", c.p.pkgVars[impPath], impPath)),
			//DeclCode: []byte(fmt.Sprintf("\t__go_import(\"%s\");\n", impPath)),
			InitCode: c.CatchOutput(1, func() { c.translateStmt(&ast.ExprStmt{X: call}, nil) }),
		})
	}

	var functions []*ast.FuncDecl
	var vars []*types.Var
	for _, file := range simplifiedFiles {
		for _, decl := range file.Nodes {
			switch d := decl.(type) {
			case *ast.FuncDecl:
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
					functions = append(functions, d)
				}
			case *ast.GenDecl:
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
							}
						}
					}
				case token.CONST:
					// skip, constants are inlined
				}
			}
		}
	}

	collectDependencies := func(f func()) []string {
		c.p.dependencies = make(map[types.Object]bool)
		f()
		var deps []string
		for o := range c.p.dependencies {
			qualifiedName := o.Pkg().Path() + "." + o.Name()
			if f, ok := o.(*types.Func); ok && f.Type().(*types.Signature).Recv() != nil {
				//deps = append(deps, qualifiedName+"~")
				deps = append(deps, qualifiedName+"___tilde_")
				continue
			}
			deps = append(deps, qualifiedName)
		}
		sort.Strings(deps)
		return deps
	}

	// variables
	var varDecls []*Decl
	varsWithInit := make(map[*types.Var]bool)
	for _, init := range c.p.InitOrder {
		for _, o := range init.Lhs {
			varsWithInit[o] = true
		}
	}
	for _, o := range vars {
		var d Decl
		if !o.Exported() {
			d.Vars = []string{c.objectName(o)}
		}
		if c.p.HasPointer[o] && !o.Exported() {
			d.Vars = append(d.Vars, c.varPtrName(o))
		}
		if _, ok := varsWithInit[o]; !ok {
			d.DceDeps = collectDependencies(func() {
				d.InitCode = []byte(fmt.Sprintf("\t\t%s = %s;\n", c.objectName(o), c.translateExpr(c.zeroValue(o.Type()), nil).String()))
			})
		}
		d.DceObjectFilter = o.Name()
		varDecls = append(varDecls, &d)
	}
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
	}

	// functions
	var funcDecls []*Decl
	var mainFunc *types.Func
	for _, fun := range functions {
		o := c.p.Defs[fun.Name].(*types.Func)
		funcInfo := c.p.FuncDeclInfos[o]
		d := Decl{
			FullName: o.FullName(),
			Blocking: len(funcInfo.Blocking) != 0,
		}
		if fun.Recv == nil {
			d.Vars = []string{c.objectName(o)}
			d.DceObjectFilter = o.Name()
			switch o.Name() {
			case "main":
				mainFunc = o
				d.DceObjectFilter = ""
			case "init":
				d.InitCode = c.CatchOutput(1, func() {
					id := c.newIdent("", types.NewSignature(nil, nil, nil, false))
					c.p.Uses[id] = o
					call := &ast.CallExpr{Fun: id}
					if len(c.p.FuncDeclInfos[o].Blocking) != 0 {
						c.Blocking[call] = true
					}
					c.translateStmt(&ast.ExprStmt{X: call}, nil)
				})
				d.DceObjectFilter = ""
			}
		}
		if fun.Recv != nil {
			recvType := o.Type().(*types.Signature).Recv().Type()
			ptr, isPointer := recvType.(*types.Pointer)
			namedRecvType, _ := recvType.(*types.Named)
			if isPointer {
				namedRecvType = ptr.Elem().(*types.Named)
			}
			d.DceObjectFilter = namedRecvType.Obj().Name()
			if !fun.Name.IsExported() {
				d.DceMethodFilter = o.Name() + "___tilde_" //jea: was: "~"
			}
		}

		d.DceDeps = collectDependencies(func() {
			d.DeclCode = c.translateToplevelFunction(fun, funcInfo)
		})
		funcDecls = append(funcDecls, &d)
	}
	if typesPkg.Name() == "main" {
		if mainFunc == nil {
			return nil, fmt.Errorf("missing main function")
		}
		id := c.newIdent("", types.NewSignature(nil, nil, nil, false))
		c.p.Uses[id] = mainFunc
		call := &ast.CallExpr{Fun: id}
		ifStmt := &ast.IfStmt{
			Cond: c.newIdent("_pkg == _mainPkg", types.Typ[types.Bool]),
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{X: call},
					&ast.AssignStmt{
						Lhs: []ast.Expr{c.newIdent("_mainFinished", types.Typ[types.Bool])},
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
					lhs += " = _pkg." + encodeIdent(o.Name())
				}
				size := int64(0)
				constructor := "null"
				switch t := o.Type().Underlying().(type) {
				case *types.Struct:
					params := make([]string, t.NumFields())
					for i := 0; i < t.NumFields(); i++ {
						params[i] = fieldName(t, i) + "_"
					}

					if t.NumFields() == 0 {
						constructor = fmt.Sprintf("function(self) return self; end")
					} else {
						constructor = fmt.Sprintf("function(self, ...)\n\t\t if self == nil then self = {}; end\n")
						constructor += fmt.Sprintf("\t\t\t local %s = ... ;\n", strings.Join(params, ", "))
						for i := 0; i < t.NumFields(); i++ {
							constructor += fmt.Sprintf("\t\t\t self.%[1]s = %[1]s_ or %[2]s;\n", fieldName(t, i), c.translateExpr(c.zeroValue(t.Field(i).Type()), nil).String())
						}
						constructor += "\t\t return self; \n\t end;\n"
						//constructor += fmt.Sprintf("\n\t %s.__constructor = %s;\n", typeName, constructor)
					}
				case *types.Basic, *types.Array, *types.Slice, *types.Chan, *types.Signature, *types.Interface, *types.Pointer, *types.Map:
					size = sizes64.Sizeof(t)
				}
				c.Printf(`%s = __newType(%d, %s, "%s.%s", %t, "%s", %t, %s);`, lhs, size, typeKind(o.Type()), o.Pkg().Name(), o.Name(), o.Name() != "", o.Pkg().Path(), o.Exported(), constructor)
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
						name += "_"
					}
					pkgPath := ""
					if !method.Exported() {
						pkgPath = method.Pkg().Path()
					}
					t := method.Type().(*types.Signature)
					entry := fmt.Sprintf(`{__prop= "%s", __name= "%s", __pkg= "%s", typ= __funcType(%s)}`, name, method.Name(), pkgPath, c.initArgs(t))
					if _, isPtr := t.Recv().Type().(*types.Pointer); isPtr {
						ptrMethods = append(ptrMethods, entry)
						continue
					}
					methods = append(methods, entry)
				}
				if len(methods) > 0 {
					c.Printf("%s.__methods_desc = {%s};", c.typeName(0, named), strings.Join(methods, ", "))
				}
				if len(ptrMethods) > 0 {
					c.Printf("%s.__methods_desc = {%s};", c.typeName(0, types.NewPointer(named)), strings.Join(ptrMethods, ", "))
				}
			})
			switch t := o.Type().Underlying().(type) {
			case *types.Array, *types.Chan, *types.Interface, *types.Map, *types.Pointer, *types.Slice, *types.Signature, *types.Struct:
				d.TypeInitCode = c.CatchOutput(0, func() {
					c.Printf("%s.init(%s); -- fullpkg.go:418", "__type__"+c.objectName(o), c.initArgs(t))
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
			d.DeclCode = []byte(fmt.Sprintf("\t%s = __%sType(%s); -- fullpkg.go:432\n", t.Name(), strings.ToLower(typeKind(t.Type())[5:]), c.initArgs(t.Type())))
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

	if len(c.p.errList) != 0 {
		return nil, c.p.errList
	}

	return &Archive{
		ImportPath:   importPath,
		Name:         typesPkg.Name(),
		Imports:      importedPaths,
		ExportData:   exportData,
		Declarations: allDecls,
		FileSet:      encodedFileSet.Bytes(),
		Minified:     minify,
	}, nil
}

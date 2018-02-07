package compiler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gijit/gi/pkg/ast"
	//"github.com/gijit/gi/pkg/constant"
	"github.com/gijit/gi/pkg/printer"
	"github.com/gijit/gi/pkg/token"
	"github.com/gijit/gi/pkg/types"
	"github.com/gijit/gi/pkg/verb"
	"os"
	"sort"
	"strings"

	"github.com/gijit/gi/pkg/compiler/analysis"
	"github.com/neelance/astrewrite"
	"golang.org/x/tools/go/gcimporter15"
	//luajit "github.com/glycerine/golua/lua"
	//"github.com/glycerine/luar"
)

func addPreludeToNewPkg(pkg *types.Package) {
	//
	// allow static type checking of the __gijit_printQuoted
	// REPL utility function. It wraps strings
	// in backticks for better printing.
	//
	scope := pkg.Scope()
	scope.Insert(getFunForGijitPrintQuoted(pkg))
}

func IncrementallyCompile(a *Archive, importPath string, files []*ast.File, fileSet *token.FileSet, importContext *ImportContext, minify bool) (*Archive, error) {

	pp("jea debug, top of incrementallyCompile()."+
		" importPath='%s' here is what files has:", importPath)
	j := 0
	for _, file := range files {
		for _, decl := range file.Nodes {
			pp("decl[%v] = '%#v'", j, decl)
		}
		j++
	}

	var newCodeText [][]byte
	var funcSrcCache map[string]string

	var typesInfo *types.Info
	if a == nil {
		typesInfo = &types.Info{
			Types:      make(map[ast.Expr]types.TypeAndValue),
			Defs:       make(map[*ast.Ident]types.Object),
			Uses:       make(map[*ast.Ident]types.Object),
			Implicits:  make(map[ast.Node]types.Object), // imports, but those without renames?
			Selections: make(map[*ast.SelectorExpr]*types.Selection),
			Scopes:     make(map[ast.Node]*types.Scope),
		}
		funcSrcCache = make(map[string]string)
	} else {
		typesInfo = a.TypesInfo
		funcSrcCache = a.FuncSrcCache
		//pp("typesInfo.Types = '%#v'", typesInfo.Types)
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
			//Sizes: sizes32,
			Sizes: sizes64,
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
	pkg, check, err = config.Check(pkg, check, importPath, fileSet, files, typesInfo, addPreludeToNewPkg)
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

	pp("got past config.Check")
	obj := pkg.Scope().Lookup("fmt.Sprintf")
	if verb.VerboseVerbose {
		pp("Sprintf obj is: '%#v'\n", obj)
		//goon.Dump(obj)
	}

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
		return false
		pp("incr.go: isBlocking: f.Pkg().Path() = '%s'", f.Pkg().Path())
		// hardcode "fmt" for now
		if f.Pkg().Path() == "fmt" {
			return false
		}

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

	// jea
	// at the repl, we need to just
	// ignore all re-ordering and take
	// everything in order as the user gives it to us.
	// so eventually we'd rather eliminate this,
	// but can't for now.
	varsWithInit := make(map[*types.Var]bool)
	for _, init := range c.p.InitOrder {
		for _, o := range init.Lhs {
			varsWithInit[o] = true
		}
	}
	// jea: I had commented the above. then reverted, which
	// seemed to fix the missing struct initializing value
	// in repl_test.go tests 027 and 028. But
	// now statements are in the wrong order. Mrrrrummph!
	//    Minor rant:
	// We need to treat the top level like in a repl like its a
	// function where sequence of statements matters!

	var typeDecls []*Decl
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

				pp("next is an *ast.FuncDecl...:'%#v'. with source:", d)
				if verb.Verbose {
					err := printer.Fprint(os.Stdout, fileSet, d)
					panicOn(err)
				}
				// cache the source for checking at the repl
				var by bytes.Buffer
				err = printer.Fprint(&by, fileSet, d)
				panicOn(err)
				funcSrcCache[d.Name.Name] = by.String()
				pp("stored in c.p.funcSrcCache['%s'] the value '%s'", d.Name.Name, funcSrcCache[d.Name.Name])

				pp("with AST:")
				if verb.Verbose {
					ast.Print(fileSet, d)
				}
				pp("done showing AST for *ast.FuncDecl...:'%#v'", d)

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
							_ = mainFunc // keep compiler happy
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
				// jea: could also be *ast.TypeSpec here when declaring a struct!
				pp("next is an *ast.GenDecl...:'%#v'. with source:", d)
				if verb.Verbose {
					err := printer.Fprint(os.Stdout, fileSet, d)
					panicOn(err)
				}
				pp("with AST:")
				if verb.Verbose {
					ast.Print(fileSet, d)
				}
				pp("done showing AST for *ast.GenDecl...:'%#v'", d)

				switch ds := d.Specs[0].(type) {
				case *ast.TypeSpec:
					pp("next decl from file.Nodes is a *ast.TypeSpec: '%#v'", ds)
				case *ast.ValueSpec:
					pp("next decl from file.Nodes is a *ast.VaueSpec: '%#v'", ds.Names[0])
				default:
					pp("next decl from file.Nodes is an unrecognized *ast.GenDecl: '%#v'", ds)
				}
				switch d.Tok {
				case token.TYPE:
					pp("we're in the token.TYPE!")
					for _, spec := range d.Specs {
						o := c.p.Defs[spec.(*ast.TypeSpec).Name].(*types.TypeName)
						c.p.typeNames = append(c.p.typeNames, o)
						c.objectName(o) // register toplevel name

						// jea: codegen here and now, in order.

						// interface Dog codegen here
						decl, by := c.oneNamedType(collectDependencies, o)
						newCodeText = append(newCodeText, by)
						typeDecls = append(typeDecls, decl)
						pp("named type codegen for '%s' generated: '%s'", o, string(by))
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

								//jea: can't filter here if we want the
								// correct order!
								if _, ok := varsWithInit[o]; !ok {

									de.DceDeps = collectDependencies(func() {
										// this is producing the $ifaceNil for the next line, and the &Beagle{word:"hiya"} part is getting lost.
										// the var snoopy Dog = &Beagle{word:"hiya"}
										// versus at place2, the Beagle does get printed!
										//

										// jea, this is getting our "var x [3]int" decl,
										// which needs to end up in code.
										de.InitCode = []byte(fmt.Sprintf("\t\t%s = %s;\n", c.objectName(o), c.translateExpr(c.zeroValue(o.Type()), nil).String()))

										pp("placeN+1, appending to newCodeText: d.InitCode='%s'", string(de.InitCode))
										newCodeText = append(newCodeText, de.InitCode)
									})

								} else {

									// jea: move in from place2 to sequential order.
									nm := c.objectName(o)
									pp("looking up var with init: '%s'", nm)

									info, ok := check.ObjMap[o]
									if !ok {
										panic(fmt.Sprintf("huh? where is the variable '%s'?", nm))
									}

									// from types/initorder, so we re-create the proper structure.
									// BEGIN INITORDER COPY

									// n:1 variable declarations such as: a, b = f()
									// introduce a node for each lhs variable (here: a, b);
									// but they all have the same initializer - emit only
									// one, for the first variable seen

									infoLhs := info.Lhs // possibly nil (see declInfo.lhs field comment)
									if infoLhs == nil {
										infoLhs = []*types.Var{o}
									}
									init := &types.Initializer{infoLhs, info.Init}
									// END INITORDER COPY

									// from place2, variables below, moved here
									// for proper sequencing of user's orders..
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

									// end codegen here and now for vars
								}
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

					// *ast.ExprStmt
					wrapWithPrint := true
					switch y := d.(type) {
					case *ast.ExprStmt:
						switch z := y.X.(type) {
						case *ast.CallExpr:
							switch id := z.Fun.(type) {
							case *ast.Ident:
								//fmt.Printf("z.Fun is Ident: %#v\n", id)
								if id.Name != "len" {
									// so "delete" doesn't print
									wrapWithPrint = false
								}
							default:
								wrapWithPrint = false
							}
						default:
						}
					default:
					}

					n := len(c.output)
					var ele string
					if bytes.HasSuffix(c.output, []byte(";\n")) {
						ele = string(bytes.TrimLeft(c.output[:n-2], " \t"))
					} else {
						ele = string(c.output)
					}
					var tmp string
					if !wrapWithPrint || strings.HasPrefix(ele, "print") {
						tmp = ele + ";"
					} else {
						pp("wrapping last line of '%s' in print at the repl", ele)
						key := fmt.Sprintf("%s", ele)
						fsrc, haveSrc := funcSrcCache[key]
						if haveSrc {
							tmp = fmt.Sprintf(`print([===[%s]===]);`, fsrc)
							pp("cache hit for '%s' -> '%s'. tmp is '%s'", key, fsrc, tmp)
						} else {
							pp("no cache hit for '%s'", key)

							// only wrap the last of the lines in print,
							// so that any helper/pre-amble anon types etc can be
							// defined without messing with the final print.
							splt := strings.Split(ele, "\n")
							nsplit := len(splt)
							if splt[nsplit-1] == "" {
								splt = splt[:nsplit-1]
								nsplit = len(splt)
							}
							if nsplit <= 1 {
								tmp = fmt.Sprintf(`print(%s);`, ele)
							} else {
								tmp = fmt.Sprintf("%s;\nprint(%s);", strings.Join(splt[:nsplit-1], "\n"), splt[nsplit-1])
							}
						}
					}
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

	// jea: we don't do c.p.InitOrder, that would confuse the repl
	// experience.
	// was comment start
	/*
		pp("jea, at variables, in package.go:392. vars='%#v'", vars)
		//pp("jea, package.go:393. c.p.InitOrder='%#v'", c.p.InitOrder)
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
	*/
	//jea: was comment end
	pp("jea, functions in package.go:393")

	// ===========================
	// functions
	// ===========================
	// moved up to stay in order.
	//	for _, fun := range functions {
	//	}

	// jea: don't treat main as special/don't require a main func.
	// At least for now.
	/*
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
	*/

	// moved up above to preserve sequence of entry.
	// 	var typeDecls []*Decl
	//typeDecls, _ = c.namedTypes(typeDecls, collectDependencies)
	typeDecls, _ = c.anonymousTypes(typeDecls, collectDependencies)

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
			TypesInfo:    typesInfo,
			config:       config,
			pkg:          pkg,
			check:        check,
			FuncSrcCache: funcSrcCache,
		}, nil
	} else {
		a.pkg = pkg
		a.check = check
		a.NewCodeText = newCodeText
		a.FuncSrcCache = funcSrcCache
	}
	return a, nil
}

// range over c.p.typeNames
func (c *funcContext) namedTypes(typeDecls []*Decl, collectDependencies func(f func()) []string) ([]*Decl, []byte) {
	var allby []byte
	for _, o := range c.p.typeNames {
		if o.IsAlias() {
			continue
		}
		one, by := c.oneNamedType(collectDependencies, o)
		typeDecls = append(typeDecls, one)
		allby = append(allby, by...)
	}
	return typeDecls, allby
}

func (c *funcContext) oneNamedType(collectDependencies func(f func()) []string, o *types.TypeName) (*Decl, []byte) {

	var allby []byte

	typeName := c.objectName(o)
	pp("on namedTypes, typeName='%s'", typeName)
	d := Decl{
		Vars:            []string{typeName},
		DceObjectFilter: o.Name(),
	}
	set_constructor := ""
	constructor := ""
	d.DceDeps = collectDependencies(func() {
		// interface Dog getting codegen here
		d.DeclCode = c.CatchOutput(0, func() {
			typeName := c.objectName(o)
			lhs := typeName
			if isPkgLevel(o) {
				// jea: might need to attend to package names
				//  eventually, or not.
				//lhs += " = $pkg." + encodeIdent(o.Name())
			}
			size := int64(0)

			switch t := o.Type().Underlying().(type) {
			case *types.Struct:
				pp("incr.go:580, in a Struct")
				params := make([]string, t.NumFields())
				for i := 0; i < t.NumFields(); i++ {
					params[i] = fieldName(t, i) + "_"
				}

				// have pointer types printed after the type they point to.
				prev := c.TypeNameSetting
				defer func() {
					c.TypeNameSetting = prev
				}()
				c.TypeNameSetting = DELAYED

				// add debug code?
				diag := ""
				//if verb.Verbose || verb.VerboseVerbose {
				//  diag = fmt.Sprintf("\n\t\t print(\"top of ctor for type '%s'\")", typeName)
				//}
				if t.NumFields() == 0 {
					//constructor = fmt.Sprintf("function(self) %s\n\t\t self.__gi_val=self; return self; end", diag)
					constructor = fmt.Sprintf("function(self) %s\n\t\t return self; end", diag)
				} else {
					constructor = fmt.Sprintf("function(self, ...) %s\n\t\t if self == nil then self = {}; end\n\t\t local args={...};\n\t\t if #args == 0 then\n", diag)
					//constructor = fmt.Sprintf("function(self, ...) %s\n\t\t self.__gi_val=self;\n\t\t local args={...};\n\t\t if #args == 0 then\n", diag)

					for i := 0; i < t.NumFields(); i++ {

						// the translateExpr call here is what
						// eventually calls c.typeName() and thus
						// generates the deferred 'anon_ptrType' and sibling type
						// variables for any pointers in the members.
						//
						constructor += fmt.Sprintf("\t\t\t self.%s = %s;\n", fieldName(t, i), c.translateExpr(c.zeroValue(t.Field(i).Type()), nil).String())
					}
					constructor += fmt.Sprintf("\t\t else \n\t\t\t local %s = ... ;\n", strings.Join(params, ", "))
					for i := 0; i < t.NumFields(); i++ {
						constructor += fmt.Sprintf("\t\t\t self.%[1]s = %[1]s_;\n", fieldName(t, i))
					}
					constructor += "\t\t end\n\t\t return self; \n\t end "
				}
				set_constructor = fmt.Sprintf("\n\t __type__%s.__constructor = %s;\n", typeName, constructor)
			case *types.Basic, *types.Array, *types.Slice, *types.Chan, *types.Signature, *types.Interface, *types.Pointer, *types.Map:
				//size = sizes32.Sizeof(t)
				size = sizes64.Sizeof(t)
				_ = size
			}
			c.Printf(`__type__%s = __gi_NewType(%d, %s, "%s", "%s", "%s.%s", %t, "%s", %t, nil);`, lhs, size, typeKind(o.Type()), o.Pkg().Name(), o.Name(), o.Pkg().Name(), o.Name(), o.Name() != "", o.Pkg().Path(), o.Exported())
			//c.Printf(`%s = $newType(%d, %s, "%s.%s", %t, "%s", %t, %s);`, lhs, size, typeKind(o.Type()), o.Pkg().Name(), o.Name(), o.Name() != "", o.Pkg().Path(), o.Exported(), constructor)

			// jea: GopherJS can defer init which adds the methods
			// to the interface, but at the REPL we cannot.
			switch o.Type().Underlying().(type) {
			case *types.Interface:
				pp("just printed __gi_NewType for an interface, now, where are the methods?")
				// are they just below here...?
			}

		})
		allby = append(allby, d.DeclCode...)
		d.MethodListCode = c.CatchOutput(0, func() {
			named := o.Type().(*types.Named)
			if _, ok := named.Underlying().(*types.Interface); ok {
				pp("is interface, not! skipping... ???")
				return
			}
			var methods []string
			var ptrMethods []string

			pp("named.NumMethods() = %v", named.NumMethods())
			for i := 0; i < named.NumMethods(); i++ {
				pp("on method i=%v, '%v'\n", i, named.Method(i).Name())
				method := named.Method(i)
				name := method.Name()
				if reservedKeywords[name] {
					name += "_"
				}
				// interface methods don't depend on any
				// particular package, so this is the empty string.
				pkgPath := ""
				if !method.Exported() {
					pkgPath = method.Pkg().Path()
				}
				t := method.Type().(*types.Signature)
				entry := fmt.Sprintf(`{__prop= "%s", __name= "%s", __pkg="%s", __typ= __gi_funcType(%s)}`, name, method.Name(), pkgPath, c.initArgs(t))

				// https://golang.org/ref/spec#Method_sets
				//
				//   The method set of the corresponding
				//   pointer type *T is the set of all
				//   methods declared with receiver *T or T
				//   (that is, it also contains the method set of T).
				//
				// jea: true seems to be needed to have 102 face_test green.
				if true {
					ptrMethods = append(ptrMethods, entry)
				} else {
					if _, isPtr := t.Recv().Type().(*types.Pointer); isPtr {
						ptrMethods = append(ptrMethods, entry)
						continue
					}
				}
				methods = append(methods, entry)
			}
			if len(methods) > 0 {
				// jea: the call to c.typeName() will add to anonType if named is anonymous, which obviously is unlikely since we're in the named type function.

				tnn := c.typeName(named)
				pp("tnn = '%s'", tnn)
				c.Printf("__type__%s.__methods_desc = {%s}; -- incr.go:817 for methods\n", tnn, strings.Join(methods, ", "))

			}
			if len(ptrMethods) > 0 {
				c.TypeNameSetting = SKIP_ANON
				pn := c.typeName(types.NewPointer(named)) // "kind_ptrType"
				pp("newPtrTypeName='%s'", pn)
				pp("c.objectName(o)='%s'", c.objectName(o))
				pn = c.objectName(o)
				// so these are the methods for B (test 102 face_test), but
				// we'll need to get them to __type__B.__ptr and not to ptrType.
				c.Printf("__type__%s.__ptr.__methods_desc = {%s}; -- incr.go:827 for ptr_methods\n", pn, strings.Join(ptrMethods, ", "))
			}
		})
		allby = append(allby, d.MethodListCode...)

		switch t := o.Type().Underlying().(type) {
		case *types.Array, *types.Chan, *types.Interface, *types.Map, *types.Pointer, *types.Slice, *types.Signature, *types.Struct:
			d.TypeInitCode = c.CatchOutput(0, func() {
				// jea: we need to initialize our interfaces with
				// their methods.
				c.Printf("__type__%s.__init(%s);", c.objectName(o), c.initArgs(t))
				_ = t // jea add
				// after methods init, then constructor
				if set_constructor != "" {
					c.Printf(set_constructor)
				}
			})
			// example of what is generated:
			// Dog.init([{prop: "Write", name: "Write", pkg: "", typ: $funcType([String], [String], false)}]);
			allby = append(allby, d.TypeInitCode...)
		}
	})
	return &d, allby
}

// range over c.p.anonTypes
func (c *funcContext) anonymousTypes(typeDecls []*Decl, collectDependencies func(f func()) []string) ([]*Decl, []byte) {

	// anonymous types
	var allby []byte
	for _, t := range c.p.anonTypes {
		one, by := c.oneAnonType(t, collectDependencies)
		typeDecls = append(typeDecls, one)
		allby = append(allby, by...)
	}
	return typeDecls, allby
}

func (c *funcContext) oneAnonType(t *types.TypeName, collectDependencies func(f func()) []string) (*Decl, []byte) {

	d := Decl{
		Vars:            []string{t.Name()},
		DceObjectFilter: t.Name(),
	}
	d.DceDeps = collectDependencies(func() {
		d.DeclCode = []byte(fmt.Sprintf("\t%s = $%sType(%s);\n", t.Name(), strings.ToLower(typeKind(t.Type())[5:]), c.initArgs(t.Type())))
	})

	return &d, d.DeclCode
}

package compiler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/glycerine/gi/pkg/ast"
	"github.com/glycerine/gi/pkg/constant"
	"github.com/glycerine/gi/pkg/printer"
	"github.com/glycerine/gi/pkg/token"
	"github.com/glycerine/gi/pkg/types"
	"runtime/debug"
	"sort"
	"strings"

	"github.com/glycerine/gi/pkg/compiler/analysis"
	"github.com/neelance/astrewrite"
	"golang.org/x/tools/go/gcimporter15"
)

// start with GopherJS Compile again, and do
// as minimal an adaption to LuaJIT as possible,
// in order to preserve the ability to compile
// entire packages. For incremental interactivity,
// see the incr.go:34 IncrementallyCompile() function.
//
func FullPackageCompile(importPath string, files []*ast.File, fileSet *token.FileSet, importContext *ImportContext, minify bool, depth int) (arch *Archive, err error) {
	vv("FullPackageCompile() top. importPath='%s'", importPath)
	defer func() {
		if arch != nil {
			pp("end of FullPackageCompile, returning arch='%#v'\n and arch.Pkg='%#v', err='%v'", arch, arch.Pkg, err)
		}
		pp("and stack is '%s'", string(debug.Stack()))
	}()

	// In incr.go, newCodeText collects Lua source from de.DeclCode for FuncDecl,
	// type definitions, variables de.InitCode, etc.
	// It is probably redundant with Archive.Declarations now.
	//
	var newCodeText [][]byte

	funcSrcCache := make(map[string]string)
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
		FullPackage:                   true,
		AllowOverShadowedNakedReturns: true,  // jea add
		DisableUnusedImportCheck:      false, // jea, differs from incremental check.
		AllowUnusedVar:                true,
		Importer: packageImporter{
			importContext: importContext,
			importError:   &importError,
		},
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

	vv("config.Check on importPath='%s'\n", importPath)
	prelude := addPreludeToNewPkg
	if importPath != "main" {
		prelude = nil
	}
	typesPkg, chk, err := config.Check(nil, nil, importPath, fileSet, files, typesInfo, prelude, depth+1)
	vv("back from config.Check on importPath='%s', err='%v', typesPkg='%#v'\n", importPath, err, typesPkg)
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

	vv("about to call gcimporter.BExportData, with typesPkg='%#v'", typesPkg)
	exportData := gcimporter.BExportData(nil, typesPkg)

	vv("back from gcimporter.BExportData")
	encodedFileSet := bytes.NewBuffer(nil)
	if err := fileSet.Write(json.NewEncoder(encodedFileSet).Encode); err != nil {
		pp("got err back from fileSet.Write: err='%v'", err)
		return nil, err
	}
	pp("fileSet.Write gave nil err='%v'", err)

	simplifiedFiles := make([]*ast.File, len(files))
	for i, file := range files {
		vv("simplifying file from pkg '%s' by calling astrewrite", file.Name.Name)
		simplifiedFiles[i] = astrewrite.Simplify(file, typesInfo, false)
	}

	isBlocking := func(f *types.Func) bool {
		return false

		// jea: avoid this, it imports shadow/math/rand by full path
		archive, err := importContext.Import(f.Pkg().Path(), "", 0)
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
	vv("about to call AnalyzePkg, on importPath='%s'", importPath)
	pkgInfo := analysis.AnalyzePkg(simplifiedFiles, fileSet, typesInfo, typesPkg, isBlocking)
	vv("back from AnalyzePkg on importPath='%s', pkgInfo = '%p'", importPath, pkgInfo)
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
		c.p.pkgVars[importedPkg.Path()] = c.newVariableWithLevel(importedPkg.Name(), true, false)
		pp("importedPkg.Path() = '%s'; importedPkg='%#v'\n", importedPkg.Path(), importedPkg)
		importedPaths = append(importedPaths, importedPkg.Path())
	}
	sort.Strings(importedPaths)
	for _, impPath := range importedPaths {
		id := c.newIdent(fmt.Sprintf(`%s.__init`, c.p.pkgVars[impPath]), types.NewSignature(nil, nil, nil, false))
		call := &ast.CallExpr{Fun: id}
		c.Blocking[call] = true
		c.Flattened[call] = true
		importDecls = append(importDecls, &Decl{
			Vars: []string{c.p.pkgVars[impPath]},

			// jea: we suspect these two import methods are colliding,
			// example:
			//   fmt = __packages["github.com/glycerine/gi/pkg/compiler/shadow/fmt"];
			//   __go_import("fmt");
			//
			// confirm by going back to just one, the direct assignment:
			//
			DeclCode: []byte(fmt.Sprintf("\t%s = __packages[\"%s\"];\n", c.p.pkgVars[impPath], impPath)),
			//DeclCode: []byte(fmt.Sprintf("\t%s = __packages[\"%s\"];\n\t__go_import(\"%s\");\n", c.p.pkgVars[impPath], impPath, omitAnyShadowPathPrefix(impPath))),

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

				// cache the source for checking at the repl
				var by bytes.Buffer
				err = printer.Fprint(&by, fileSet, d)
				panicOn(err)
				funcSrcCache[d.Name.Name] = by.String()
				pp("stored in funcSrcCache['%s'] the value '%s'", d.Name.Name, funcSrcCache[d.Name.Name])

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
			vv("varsWithInit true for o='%#v'", o)
			varsWithInit[o] = true
		}
	}
	for _, o := range vars {
		vv("ranging over vars, o = '%#v'", o)
		var d Decl
		if !o.Exported() {
			vv("o not exported")
			d.Vars = []string{c.objectName(o)}
		} else {
			vv("o exported")
		}
		if c.p.HasPointer[o] && !o.Exported() {
			d.Vars = append(d.Vars, c.varPtrName(o))
		}
		if _, ok := varsWithInit[o]; !ok {
			d.DceDeps = collectDependencies(func() {
				d.InitCode = []byte(fmt.Sprintf("\t\t%s = %s; --fullpkg.go:277\n",
					// c.objectName(o),
					c.objectNameWithPackagePrefix(o),
					c.translateExpr(c.zeroValue(o.Type()), nil).String()))
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
			//test 1002: Verbose = false is written here.
			// jea add:
			d.InitCode = append(d.InitCode, []byte(fmt.Sprintf(
				"\t\t\t--[[ fullpkg.go:321 --]] print(\"99999999 jea debug! package initialization code called, for '%s'!\")\n", importPath))...)
			d.Vars = append(d.Vars, c.localVars...)
			// jea add:
			d.initializePackageVars(lhs)
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
		pp("doing codegen for function '%s'", d.FullName)
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
				d.InitCode = append(d.InitCode, []byte(" --[[ fullpkg.go:343 --]]")...)
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
				pp("func: adding ___tilde_ for d.DecMethodFilter")
				d.DceMethodFilter = o.Name() + "___tilde_" //jea: was: "~"
			}
		}

		d.DceDeps = collectDependencies(func() {
			d.DeclCode = c.translateToplevelFunction(fun, funcInfo)
			pp("translateToplevelFunction returned '%s'", d.DeclCode)
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
				lhsPre := fmt.Sprintf("__type__.%s", typesPkg.Name())
				lhs := fmt.Sprintf("%s.%s", lhsPre, typeName)
				// jea comment out for now... b/c getting stuff like:
				//
				// __type__.GONZAGA = _pkg.GONZAGA --[[ fullpkg.go:395 --]]  = __newType(16, __kindInterface, "spkg_tst.GONZAGA", true, "github.com/glycerine/gi/pkg/compiler/spkg_tst", true, nil);
				//
				/*if isPkgLevel(o) {
					lhs += " = _pkg." + encodeIdent(o.Name()) + " --[[ fullpkg.go:395 --]] "
				}
				*/
				size := int64(0)
				constructor := "nil"
				switch t := o.Type().Underlying().(type) {
				case *types.Struct:
					params := make([]string, t.NumFields())
					for i := 0; i < t.NumFields(); i++ {
						params[i] = fieldName(t, i) + "_"
					}

					if t.NumFields() == 0 {
						constructor = fmt.Sprintf("function(self) return self; end")
					} else {
						constructor = fmt.Sprintf("function(self, ...) if self == nil then self = {}; end; ")
						constructor += fmt.Sprintf("local %s = ... ; ", strings.Join(params, ", "))
						for i := 0; i < t.NumFields(); i++ {
							constructor += fmt.Sprintf(" self.%[1]s = %[1]s_ or %[2]s; ", fieldName(t, i), c.translateExpr(c.zeroValue(t.Field(i).Type()), nil).String())
						}
						constructor += " return self; end " // jea: can't have a semicolon aftre 'end;' will mess up LuaJIT's parse.
						//constructor += fmt.Sprintf("\n\t %s.__constructor = %s;\n", typeName, constructor)
					}
				case *types.Basic, *types.Array, *types.Slice, *types.Chan, *types.Signature, *types.Interface, *types.Pointer, *types.Map:
					size = sizes64.Sizeof(t)
				}
				c.Printf("%[1]s = %[1]s or {};\n", lhsPre)
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

					entry := fmt.Sprintf(`{__prop= "%s", __name= %s, __pkg= "%s", typ= __funcType(%s)}`, name, encodeString(method.Name()), pkgPath, c.initArgs(t))
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
					c.Printf("%s.%s.%s.init(%s); -- fullpkg.go:463", "__type__", typesPkg.Name(), c.objectName(o), c.initArgs(t))
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
			d.DeclCode = []byte(fmt.Sprintf("\t%s = __%sType(%s); -- fullpkg.go:479\n", t.Name(), strings.ToLower(typeKind(t.Type())[6:]), c.initArgs(t.Type())))
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

	vv("at end of FullPackageCompile of importPath='%s', here are allDecls:", importPath)
	for k, d := range allDecls {
		vv("allDecls[k=%v] has .DeclCode = '%s'", k, string(d.DeclCode))
	}
	return &Archive{
		SavedArchive: SavedArchive{
			ImportPath:   importPath,
			Name:         typesPkg.Name(),
			Imports:      importedPaths,
			ExportData:   exportData,
			Declarations: allDecls,
			FileSet:      encodedFileSet.Bytes(),
			Minified:     minify,
		},
		Pkg:          typesPkg,
		NewCodeText:  newCodeText,
		TypesInfo:    typesInfo,
		Config:       config,
		Check:        chk,
		FuncSrcCache: funcSrcCache,
	}, nil
}

// jea add:
func (d *Decl) initializePackageVars(lhs []ast.Expr) {
	code := bytes.NewBuffer(nil)
	for _, e := range lhs {
		nm := nameHelper(e)
		d.Vars = append(d.Vars, nm)

		if ast.IsExported(nm) {
			fmt.Fprintf(code, "\t\t__pkg.%[1]s = %[1]s; -- fullpkg.go:556\n", nm)
		}
	}
	d.InitCode = append(d.InitCode, code.Bytes()...)
}

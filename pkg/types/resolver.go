// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"fmt"
	"github.com/gijit/gi/pkg/ast"
	"github.com/gijit/gi/pkg/constant"
	"github.com/gijit/gi/pkg/token"
	"strconv"
	"strings"
	"unicode"
)

// A DeclInfo describes a package-level const, type, var, or func declaration.
type DeclInfo struct {
	File  *Scope        // scope of file containing this declaration
	Lhs   []*Var        // lhs of n:1 variable declarations, or nil
	Typ   ast.Expr      // type, or nil
	Init  ast.Expr      // init/orig expression, or nil
	Fdecl *ast.FuncDecl // func declaration, or nil
	Alias bool          // type alias declaration

	// The deps field tracks initialization expression dependencies.
	// As a special (overloaded) case, it also tracks dependencies of
	// interface types on embedded interfaces (see ordering.go).
	deps objSet // lazily initialized
}

// An objSet is simply a set of objects.
type objSet map[Object]bool

// hasInitializer reports whether the declared object has an initialization
// expression or function body.
func (d *DeclInfo) hasInitializer() bool {
	return d.Init != nil || d.Fdecl != nil && d.Fdecl.Body != nil
}

// addDep adds obj to the set of objects d's init expression depends on.
func (d *DeclInfo) addDep(obj Object) {
	m := d.deps
	if m == nil {
		m = make(objSet)
		d.deps = m
	}
	m[obj] = true
}

// arityMatch checks that the lhs and rhs of a const or var decl
// have the appropriate number of names and init exprs. For const
// decls, init is the value spec providing the init exprs; for
// var decls, init is nil (the init exprs are in s in this case).
func (check *Checker) arityMatch(s, init *ast.ValueSpec) {
	l := len(s.Names)
	r := len(s.Values)
	if init != nil {
		r = len(init.Values)
	}

	switch {
	case init == nil && r == 0:
		// var decl w/o init expr
		if s.Type == nil {
			check.errorf(s.Pos(), "missing type or init expr")
		}
	case l < r:
		if l < len(s.Values) {
			// init exprs from s
			n := s.Values[l]
			check.errorf(n.Pos(), "extra init expr %s", n)
			// TODO(gri) avoid declared but not used error here
		} else {
			// init exprs "inherited"
			check.errorf(s.Pos(), "extra init expr at %s", check.fset.Position(init.Pos()))
			// TODO(gri) avoid declared but not used error here
		}
	case l > r && (init != nil || r != 1):
		n := s.Names[r]
		check.errorf(n.Pos(), "missing init expr for %s", n)
	}
}

func validatedImportPath(path string) (string, error) {
	s, err := strconv.Unquote(path)
	if err != nil {
		return "", err
	}
	if s == "" {
		return "", fmt.Errorf("empty string")
	}
	const illegalChars = `!"#$%&'()*,:;<=>?[\]^{|}` + "`\uFFFD"
	for _, r := range s {
		if !unicode.IsGraphic(r) || unicode.IsSpace(r) || strings.ContainsRune(illegalChars, r) {
			return s, fmt.Errorf("invalid character %#U", r)
		}
	}
	return s, nil
}

// declarePkgObj declares obj in the package scope, records its ident -> obj mapping,
// and updates check.ObjMap. The object must not be a function or method.
func (check *Checker) declarePkgObj(ident *ast.Ident, obj Object, d *DeclInfo) {
	assert(ident.Name == obj.Name())

	// spec: "A package-scope or file-scope identifier with name init
	// may only be declared to be a function with this (func()) signature."
	if ident.Name == "init" {
		check.errorf(ident.Pos(), "cannot declare init - must be func")
		return
	}

	// spec: "The main package must have package name main and declare
	// a function main that takes no arguments and returns no value."
	if ident.Name == "main" && check.pkg.name == "main" {
		check.errorf(ident.Pos(), "cannot declare main - must be func")
		return
	}

	check.declare(check.pkg.scope, ident, obj, token.NoPos)
	check.ObjMap[obj] = d
	obj.setOrder(uint32(len(check.ObjMap)))
}

// filename returns a filename suitable for debugging output.
func (check *Checker) filename(fileNo int) string {
	file := check.files[fileNo]
	if pos := file.Pos(); pos.IsValid() {
		return check.fset.File(pos).Name()
	}
	return fmt.Sprintf("file[%d]", fileNo)
}

func (check *Checker) importPackage(pos token.Pos, path, dir string) *Package {
	// If we already have a package for the given (path, dir)
	// pair, use it instead of doing a full import.
	// Checker.impMap only caches packages that are marked Complete
	// or fake (dummy packages for failed imports). Incomplete but
	// non-fake packages do require an import to complete them.
	key := importKey{path, dir}
	imp := check.impMap[key]
	if imp != nil {
		return imp
	}

	// no package yet => import it
	if path == "C" && check.conf.FakeImportC {
		imp = NewPackage("C", "C")
		imp.fake = true
	} else {
		// ordinary import
		var err error
		if importer := check.conf.Importer; importer == nil {
			err = fmt.Errorf("Config.Importer not installed")
		} else if importerFrom, ok := importer.(ImporterFrom); ok {
			imp, err = importerFrom.ImportFrom(path, dir, 0)
			if imp == nil && err == nil {
				err = fmt.Errorf("Config.Importer.ImportFrom(%s, %s, 0) returned nil but no error", path, dir)
			}
		} else {
			// jea: our import "fmt" is calling here.
			imp, err = importer.Import(path)
			fmt.Printf("\n jea debug: types/resolver.go:164, back from importer.Import(path='%s') imp='%#v', err='%v'\n", path, imp, err)
			if imp == nil && err == nil {
				fmt.Printf("\n jea debug: imp was nil!?! err was nil\n")
				err = fmt.Errorf("Config.Importer.Import(%s) returned nil but no error", path)
			}
		}
		// make sure we have a valid package name
		// (errors here can only happen through manipulation of packages after creation)
		if err == nil && imp != nil && (imp.name == "_" || imp.name == "") {
			fmt.Printf("\n jea debug: imp.name = '%s' so nilling out imp.\n", imp.name)
			err = fmt.Errorf("invalid package name: %q", imp.name)
			imp = nil // create fake package below
		}
		if err != nil {
			check.errorf(pos, "could not import %s (%s)", path, err)
			if imp == nil {
				// create a new fake package
				// come up with a sensible package name (heuristic)
				name := path
				if i := len(name); i > 0 && name[i-1] == '/' {
					name = name[:i-1]
				}
				if i := strings.LastIndex(name, "/"); i >= 0 {
					name = name[i+1:]
				}
				imp = NewPackage(path, name)
			}
			// continue to use the package as best as we can
			imp.fake = true // avoid follow-up lookup failures
		}
	}

	// package should be complete or marked fake, but be cautious
	if imp.complete || imp.fake {
		check.impMap[key] = imp
		return imp
	}

	// jea: import "fmt" is getting here, returning nil.
	fmt.Printf("\n jea debug: types/resolver.go:201: something went wrong, returning nil. imp.complete='%v', imp.fake='%v'\n", imp.complete, imp.fake)
	// something went wrong (importer may have returned incomplete package without error)
	return nil
}

// collectObjects collects all file and package objects and inserts them
// into their respective scopes. It also performs imports and associates
// methods with receiver base type names.
func (check *Checker) collectObjects() {
	pkg := check.pkg

	// pkgImports is the set of packages already imported by any package file seen
	// so far. Used to avoid duplicate entries in pkg.imports. Allocate and populate
	// it (pkg.imports may not be empty if we are checking test files incrementally).
	// Note that pkgImports is keyed by package (and thus package path), not by an
	// importKey value. Two different importKey values may map to the same package
	// which is why we cannot use the check.impMap here.
	var pkgImports = make(map[*Package]bool)
	for _, imp := range pkg.imports {
		pkgImports[imp] = true
	}

	for fileNo, file := range check.files {
		// The package identifier denotes the current package,
		// but there is no corresponding package object.
		check.recordDef(file.Name, nil)

		// Use the actual source file extent rather than *ast.File extent since the
		// latter doesn't include comments which appear at the start or end of the file.
		// Be conservative and use the *ast.File extent if we don't have a *token.File.
		pos, end := file.Pos(), file.End()
		if f := check.fset.File(file.Pos()); f != nil {
			pos, end = token.Pos(f.Base()), token.Pos(f.Base()+f.Size())
		}
		fileScope := NewScope(check.pkg.scope, pos, end, check.filename(fileNo))
		check.recordScope(file, fileScope)

		// determine file directory, necessary to resolve imports
		// FileName may be "" (typically for tests) in which case
		// we get "." as the directory which is what we would want.
		fileDir := dir(check.fset.Position(file.Name.Pos()).Filename)

		for _, decl := range file.Nodes {
			switch d := decl.(type) {
			case *ast.ExprStmt:
			// jea: check the ExprStmt below, handle other Stmt here.
			case ast.Stmt:
				// jea: allow all statements at top level.
				check.simpleStmt(d)
				continue
			}

			switch d := decl.(type) {
			case *ast.BadDecl:
				// ignore

			case *ast.GenDecl:
				var last *ast.ValueSpec // last ValueSpec with type or init exprs seen
				for iota, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.ImportSpec:
						// import package
						path, err := validatedImportPath(s.Path.Value)
						if err != nil {
							check.errorf(s.Path.Pos(), "invalid import path (%s)", err)
							continue
						}

						imp := check.importPackage(s.Path.Pos(), path, fileDir)
						if imp == nil {
							fmt.Printf("\n jea debug types/resolver.go: check.importPackage(path='%s') returned nil.\n", path)
							continue
						}
						fmt.Printf("\n jea debug types/resolver.go:277: good, check.importPackage(path='%s') was not nil. imp='%#v'.  our spec s = '%#v'\n", path, imp, s)

						// add package to list of explicit imports
						// (this functionality is provided as a convenience
						// for clients; it is not needed for type-checking)
						if !pkgImports[imp] {
							pkgImports[imp] = true
							pkg.imports = append(pkg.imports, imp)
						}

						// local name overrides imported package name
						name := imp.name
						if s.Name != nil {
							name = s.Name.Name
							if path == "C" {
								// match cmd/compile (not prescribed by spec)
								check.errorf(s.Name.Pos(), `cannot rename import "C"`)
								continue
							}
							if name == "init" {
								check.errorf(s.Name.Pos(), "cannot declare init - must be func")
								continue
							}
						}

						obj := NewPkgName(s.Pos(), pkg, name, imp)
						if s.Name != nil {
							// in a dot-import, the dot represents the package
							fmt.Printf("\n 88888 jea debug: check.recordDef(s.Name='%s') \n", s.Name)
							check.recordDef(s.Name, obj)
						} else {
							// jea: import "fmt" is going here!
							fmt.Printf("\n 9999 jea debug: check.recordImplicit(s.Name='%s') \n", s.Name)
							check.recordImplicit(s, obj)
						}
						fmt.Printf("\n 00000 jea debug: resolver.go:311 past the recordDef / recordImplicit \n")

						if path == "C" {
							// match cmd/compile (not prescribed by spec)
							obj.used = true
						}

						// add import to file scope
						if name == "." {
							// merge imported scope with file scope
							for _, obj := range imp.scope.elems {
								// A package scope may contain non-exported objects,
								// do not import them!
								if obj.Exported() {
									// TODO(gri) When we import a package, we create
									// a new local package object. We should do the
									// same for each dot-imported object. That way
									// they can have correct position information.
									// (We must not modify their existing position
									// information because the same package - found
									// via Config.Packages - may be dot-imported in
									// another package!)
									check.declare(fileScope, nil, obj, token.NoPos)
								}
							}
							// add position to set of dot-import positions for this file
							// (this is only needed for "imported but not used" errors)
							check.addUnusedDotImport(fileScope, imp, s.Pos())
						} else {
							// declare imported package object in file scope

							// jea: Ah HAH! this is where the imported package
							//  should be getting added to the scopes.
							fmt.Printf("jea debug, fileScope before adding fmt: '%s'\n\n", fileScope)
							check.declare(fileScope, nil, obj, token.NoPos)

							// jea try addint to upper scope...
							fmt.Printf("jea debug, fileScope after adding fmt: '%s'\n\n", fileScope)
						}

					case *ast.ValueSpec:
						switch d.Tok {
						case token.CONST:
							// determine which initialization expressions to use
							switch {
							case s.Type != nil || len(s.Values) > 0:
								last = s
							case last == nil:
								last = new(ast.ValueSpec) // make sure last exists
							}

							// declare all constants
							for i, name := range s.Names {
								obj := NewConst(name.Pos(), pkg, name.Name, nil, constant.MakeInt64(int64(iota)))

								var init ast.Expr
								if i < len(last.Values) {
									init = last.Values[i]
								}

								d := &DeclInfo{File: fileScope, Typ: last.Type, Init: init}
								check.declarePkgObj(name, obj, d)
							}

							check.arityMatch(s, last)

						case token.VAR:
							lhs := make([]*Var, len(s.Names))
							// If there's exactly one rhs initializer, use
							// the same DeclInfo d1 for all lhs variables
							// so that each lhs variable depends on the same
							// rhs initializer (n:1 var declaration).
							var d1 *DeclInfo
							if len(s.Values) == 1 {
								// The lhs elements are only set up after the for loop below,
								// but that's ok because declareVar only collects the DeclInfo
								// for a later phase.
								d1 = &DeclInfo{File: fileScope, Lhs: lhs, Typ: s.Type, Init: s.Values[0]}
							}

							// declare all variables
							for i, name := range s.Names {
								obj := NewVar(name.Pos(), pkg, name.Name, nil)
								lhs[i] = obj

								d := d1
								if d == nil {
									// individual assignments
									var init ast.Expr
									if i < len(s.Values) {
										init = s.Values[i]
									}
									d = &DeclInfo{File: fileScope, Typ: s.Type, Init: init}
								}

								check.declarePkgObj(name, obj, d)
							}

							check.arityMatch(s, nil)

						default:
							check.invalidAST(s.Pos(), "invalid token %s", d.Tok)
						}

					case *ast.TypeSpec:
						obj := NewTypeName(s.Name.Pos(), pkg, s.Name.Name, nil)
						check.declarePkgObj(s.Name, obj, &DeclInfo{File: fileScope, Typ: s.Type, Alias: s.Assign.IsValid()})

					default:
						check.invalidAST(s.Pos(), "unknown ast.Spec node %T", s)
					}
				}

			case *ast.FuncDecl:
				name := d.Name.Name
				obj := NewFunc(d.Name.Pos(), pkg, name, nil)
				if d.Recv == nil {
					// regular function
					if name == "init" {
						// don't declare init functions in the package scope - they are invisible
						obj.parent = pkg.scope
						check.recordDef(d.Name, obj)
						// init functions must have a body
						if d.Body == nil {
							check.softErrorf(obj.pos, "missing function body")
						}
					} else {
						check.declare(pkg.scope, d.Name, obj, token.NoPos)
					}
				} else {
					// method
					check.recordDef(d.Name, obj)
					// Associate method with receiver base type name, if possible.
					// Ignore methods that have an invalid receiver, or a blank _
					// receiver name. They will be type-checked later, with regular
					// functions.
					if list := d.Recv.List; len(list) > 0 {
						typ := list[0].Type
						if ptr, _ := typ.(*ast.StarExpr); ptr != nil {
							typ = ptr.X
						}
						if base, _ := typ.(*ast.Ident); base != nil && base.Name != "_" {
							check.assocMethod(base.Name, obj)
						}
					}
				}
				info := &DeclInfo{File: fileScope, Fdecl: d}
				check.ObjMap[obj] = info
				obj.setOrder(uint32(len(check.ObjMap)))

				// jea: allow assignments at the top level!
			case *ast.AssignStmt:
				pp("we see an AssignStmt, '%#v'/ pos:end=%v:%v", d, d.Pos(), d.End())
				check.simpleStmt(d)

				// jea: allow ++ and -- at toplevel
			case *ast.IncDecStmt:
				check.simpleStmt(d)

				// jea: allow for loops at the top level!
			case *ast.ForStmt:
				pp("we see an *ast.ForStmt, '%#v'/ pos:end=%v:%v", d, d.Pos(), d.End())
				check.simpleStmt(d)

				// jea: allow evaluation of expressions at the top level!
			case *ast.ExprStmt:
				pp("we see an *ast.ExprStmt, '%#v'/ pos:end=%v:%v", d, d.Pos(), d.End())

				var x operand
				kind := check.rawExpr(&x, d.X, nil)
				var msg string
				switch x.mode {
				default:
					if kind == statement {
						//jea
						//return
					}
				case builtin:
					msg = "must be called"
					check.errorf(x.pos(), "%s %s", &x, msg)
				case typexpr:
					msg = "is not an expression"
					check.errorf(x.pos(), "%s %s", &x, msg)
				}
				check.NewCode = append(check.NewCode, &NewStuff{Node: d, IsExpr: true})

			default:
				check.invalidAST(d.Pos(), "unknown ast.Node node %T", d)
			}
		}
	}

	// verify that objects in package and file scopes have different names
	for _, scope := range check.pkg.scope.children /* file scopes */ {
		for _, obj := range scope.elems {
			if alt := pkg.scope.Lookup(obj.Name()); alt != nil {
				if pkg, ok := obj.(*PkgName); ok {
					check.errorf(alt.Pos(), "%s already declared through import of %s", alt.Name(), pkg.Imported())
					check.reportAltDecl(pkg)
				} else {
					check.errorf(alt.Pos(), "%s already declared through dot-import of %s", alt.Name(), obj.Pkg())
					// TODO(gri) dot-imported objects don't have a position; reportAltDecl won't print anything
					check.reportAltDecl(obj)
				}
			}
		}
	}
}

// packageObjects typechecks all package objects in objList, but not function bodies.
func (check *Checker) packageObjects(objList []Object) {
	// add new methods to already type-checked types (from a prior Checker.Files call)
	for _, obj := range objList {
		if obj, _ := obj.(*TypeName); obj != nil && obj.typ != nil {
			check.addMethodDecls(obj)
		}
	}

	// pre-allocate space for type declaration paths so that the underlying array is reused
	typePath := make([]*TypeName, 0, 8)

	for _, obj := range objList {
		check.objDecl(obj, nil, typePath)
	}

	// At this point we may have a non-empty check.methods map; this means that not all
	// entries were deleted at the end of typeDecl because the respective receiver base
	// types were not found. In that case, an error was reported when declaring those
	// methods. We can now safely discard this map.
	check.methods = nil
}

// functionBodies typechecks all function bodies.
func (check *Checker) functionBodies() {
	for _, f := range check.funcs {
		check.funcBody(f.decl, f.name, f.sig, f.body)
	}
}

// unusedImports checks for unused imports.
func (check *Checker) unusedImports() {
	// if function bodies are not checked, packages' uses are likely missing - don't check
	if check.conf.IgnoreFuncBodies {
		return
	}

	// spec: "It is illegal (...) to directly import a package without referring to
	// any of its exported identifiers. To import a package solely for its side-effects
	// (initialization), use the blank identifier as explicit package name."

	// check use of regular imported packages
	for _, scope := range check.pkg.scope.children /* file scopes */ {
		for _, obj := range scope.elems {
			if obj, ok := obj.(*PkgName); ok {
				// Unused "blank imports" are automatically ignored
				// since _ identifiers are not entered into scopes.
				if !obj.used {
					path := obj.imported.path
					base := pkgName(path)
					if obj.name == base {
						check.softErrorf(obj.pos, "%q imported but not used", path)
					} else {
						check.softErrorf(obj.pos, "%q imported but not used as %s", path, obj.name)
					}
				}
			}
		}
	}

	// check use of dot-imported packages
	for _, unusedDotImports := range check.unusedDotImports {
		for pkg, pos := range unusedDotImports {
			check.softErrorf(pos, "%q imported but not used", pkg.path)
		}
	}
}

// pkgName returns the package name (last element) of an import path.
func pkgName(path string) string {
	if i := strings.LastIndex(path, "/"); i >= 0 {
		path = path[i+1:]
	}
	return path
}

// dir makes a good-faith attempt to return the directory
// portion of path. If path is empty, the result is ".".
// (Per the go/build package dependency tests, we cannot import
// path/filepath and simply use filepath.Dir.)
func dir(path string) string {
	if i := strings.LastIndexAny(path, `/\`); i > 0 {
		return path[:i]
	}
	// i <= 0
	return "."
}

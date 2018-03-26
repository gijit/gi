// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"fmt"

	"github.com/gijit/gi/pkg/ast"
	"github.com/gijit/gi/pkg/constant"
	"github.com/gijit/gi/pkg/token"
)

var DisallowMethodRedefinition = false

func (check *Checker) reportAltDecl(obj Object) {
	if pos := obj.Pos(); pos.IsValid() {
		// We use "other" rather than "previous" here because
		// the first declaration seen may not be textually
		// earlier in the source.
		check.errorf(pos, "\tother declaration of %s", obj.Name()) // secondary error, \t indented
	}
}

func (check *Checker) declare(scope *Scope, id *ast.Ident, obj Object, pos token.Pos) {
	pp("declare called for obj.Name()='%s', id='%#v', scope='%#v'", obj.Name(), id, scope)

	// spec: "The blank identifier, represented by the underscore
	// character _, may be used in a declaration like any other
	// identifier but the declaration does not introduce a new
	// binding."

	// jea debug
	//	if obj.Name() == "fmt" {
	//		panic("where declare")
	//	}

	if obj.Name() != "_" {
		var alt Object
		if scope != nil {
			// jea replace:
			// alt = scope.Insert(obj)
			alt = scope.Replace(obj)
			pp("types/decl.go:40 we did scope.Replace obj for obs = '%#v'", obj)
		} else {
			// at top level package scope
			//jea: alt = check.pkg.scope.Insert(obj)
			alt = check.pkg.scope.Replace(obj)
		}
		if alt != nil {
			// re-declaration of variables, such as `a:=1;a:=1` errors out here.
			pp("previous re-declaration errors now never happen at the repl, so we detect for repl purposes: alt.Name()='%s'. obj='%T'/'%#v'", alt.Name(), obj, obj)
			switch alt.(type) {
			case *TypeName:
				pp("we have a TypeName being re-declared: '%s'", alt.Name())
			}
			/*		panic("should never reach")
					check.errorf(obj.Pos(), "%s redeclared in this block", obj.Name())
					check.reportAltDecl(alt)
					return
			*/
		}
		obj.setScopePos(pos)
	}
	pp("declare: just before id != nil check. id='%#v'", id)
	if id != nil {
		//check.recordDef(id, obj)
		check.recordDefAtScope(id, obj, scope, nil)
	}
}

// objDecl type-checks the declaration of obj in its respective (file) context.
// See check.typ for the details on def and path.
func (check *Checker) objDecl(obj Object, def *Named, path []*TypeName) {
	vv("jea debug: types/decl.go:76, check.objDecl running top.")
	if obj.Type() != nil {
		return // already checked - nothing to do
	}

	if trace {
		check.trace(obj.Pos(), "-- declaring %s", obj.Name())
		if obj.Name() == "Error" {
			//panic("here?")
		}
		check.indent++
		defer func() {
			check.indent--
			check.trace(obj.Pos(), "=> %s", obj)
		}()
	}

	d := check.ObjMap[obj]
	if d == nil {
		check.dump("%s: %s should have been declared", obj.Pos(), obj.Name())
		//debug
		/*
			vv("debug the missing '%v', here is check.ObjMap:", obj.Name())
			for k, v := range check.ObjMap {
				vv("k='%v',   v='%#v'\n", k.Name(), v)
			}
		*/
		unreachable()
	}
	pp("obj = '%#v'/'%s'", obj, obj)
	pp("obj.Name() = '%s'", obj.Name())
	pp("d = '%#v'/'%s'", d, d)

	// save/restore current context and setup object context
	defer func(ctxt context) {
		check.context = ctxt
	}(check.context)
	check.context = context{
		scope: d.File,
	}

	// Const and var declarations must not have initialization
	// cycles. We track them by remembering the current declaration
	// in check.decl. Initialization expressions depending on other
	// consts, vars, or functions, add dependencies to the current
	// check.decl.
	switch obj := obj.(type) {
	case *Const:
		check.decl = d // new package-level const decl
		check.constDecl(obj, d.Typ, d.Init)
	case *Var:
		check.decl = d // new package-level var decl
		check.varDecl(obj, d.Lhs, d.Typ, d.Init)
	case *TypeName:
		// invalid recursive types are detected via path
		check.typeDecl(obj, d.Typ, def, path, d.Alias) //  // spkg_tst3.Error here
	case *Func:
		// functions may be recursive - no need to track dependencies
		// jea: new function declarations happen here.
		pp("check.objDecl calling check.funcDecl with obj='%v', d='%#v'", obj.Name(), d)
		check.funcDecl(obj, d)
	default:
		unreachable()
	}
}

func (check *Checker) constDecl(obj *Const, typ, init ast.Expr) {
	assert(obj.typ == nil)

	if obj.visited {
		obj.typ = Typ[Invalid]
		return
	}
	obj.visited = true

	// use the correct value of iota
	assert(check.iota == nil)
	check.iota = obj.val
	defer func() { check.iota = nil }()

	// provide valid constant value under all circumstances
	obj.val = constant.MakeUnknown()

	// determine type, if any
	if typ != nil {
		t := check.typ(typ)
		if !isConstType(t) {
			check.errorf(typ.Pos(), "invalid constant type %s", t)
			obj.typ = Typ[Invalid]
			return
		}
		obj.typ = t
	}

	// check initialization
	var x operand
	if init != nil {
		check.expr(&x, init)
	}
	check.initConst(obj, &x)
}

func (check *Checker) varDecl(obj *Var, lhs []*Var, typ, init ast.Expr) {
	assert(obj.typ == nil)

	if obj.visited {
		obj.typ = Typ[Invalid]
		return
	}
	obj.visited = true

	// var declarations cannot use iota
	assert(check.iota == nil)

	// determine type, if any
	if typ != nil {
		obj.typ = check.typ(typ)
		// We cannot spread the type to all lhs variables if there
		// are more than one since that would mark them as checked
		// (see Checker.objDecl) and the assignment of init exprs,
		// if any, would not be checked.
		//
		// TODO(gri) If we have no init expr, we should distribute
		// a given type otherwise we need to re-evalate the type
		// expr for each lhs variable, leading to duplicate work.
	}

	// check initialization
	if init == nil {
		if typ == nil {
			// error reported before by arityMatch
			obj.typ = Typ[Invalid]
		}
		return
	}

	if lhs == nil || len(lhs) == 1 {
		assert(lhs == nil || lhs[0] == obj)
		var x operand
		check.expr(&x, init)
		check.initVar(obj, &x, "variable declaration")
		return
	}

	if debug {
		// obj must be one of lhs
		found := false
		for _, lhs := range lhs {
			if obj == lhs {
				found = true
				break
			}
		}
		if !found {
			panic("inconsistent lhs")
		}
	}

	// We have multiple variables on the lhs and one init expr.
	// Make sure all variables have been given the same type if
	// one was specified, otherwise they assume the type of the
	// init expression values (was issue #15755).
	if typ != nil {
		for _, lhs := range lhs {
			lhs.typ = obj.typ
		}
	}

	check.initVars(lhs, []ast.Expr{init}, token.NoPos)
}

// underlying returns the underlying type of typ; possibly by following
// forward chains of named types. Such chains only exist while named types
// are incomplete.
func underlying(typ Type) Type {
	for {
		n, _ := typ.(*Named)
		if n == nil {
			break
		}
		typ = n.underlying
	}
	return typ
}

func (n *Named) setUnderlying(typ Type) {
	pp("jea debug, decl.go: Named.setUnderlying for typ='%#v'/'%s'", typ, typ)
	switch x := typ.(type) {
	case *Named:
		if x.methods != nil {
			pp("len of x.methods is %v", len(x.methods))
			for i, me := range x.methods {
				pp("x.methods[i=%v] = '%#v'; me.name='%s', me.typ: '%s'", i, me, me.name, me.typ)
			}
		}
	}
	if n != nil {
		n.underlying = typ
	}
}

func (check *Checker) typeDecl(obj *TypeName, typ ast.Expr, def *Named, path []*TypeName, alias bool) {
	assert(obj.typ == nil)

	// type declarations cannot use iota
	assert(check.iota == nil)

	if alias {
		obj.typ = Typ[Invalid]
		obj.typ = check.typExpr(typ, nil, append(path, obj))

	} else {

		named := &Named{obj: obj}
		def.setUnderlying(named) // 444444 jea: stacktrace of setUnderlying here.
		obj.typ = named          // make sure recursive type declarations terminate

		pp("4444444 jea debug, about to call check.typExpr to determine underlying")
		pp("path='%s', obj='%s', named='%s', typ='%s'", path, obj, named, typ)
		// determine underlying type of named
		check.typExpr(typ, named, append(path, obj))

		// The underlying type of named may be itself a named type that is
		// incomplete:
		//
		//	type (
		//		A B
		//		B *C
		//		C A
		//	)
		//
		// The type of C is the (named) type of A which is incomplete,
		// and which has as its underlying type the named type B.
		// Determine the (final, unnamed) underlying type by resolving
		// any forward chain (they always end in an unnamed type).
		named.underlying = underlying(named.underlying)

	}

	// check and add associated methods
	// TODO(gri) It's easy to create pathological cases where the
	// current approach is incorrect: In general we need to know
	// and add all methods _before_ type-checking the type.
	// See https://play.golang.org/p/WMpE0q2wK8
	check.addMethodDecls(obj) //  // spkg_tst3.Error here
}

func (check *Checker) addMethodDecls(obj *TypeName) {
	// get associated methods
	methods := check.Methods[obj.name]
	if len(methods) == 0 {
		return // no methods
	}
	delete(check.Methods, obj.name)

	// use an objset to check for name conflicts
	var mset objset

	// spec: "If the base type is a struct type, the non-blank method
	// and field names must be distinct."
	base, _ := obj.typ.(*Named) // nil if receiver base type is type alias
	if base != nil {
		if t, _ := base.underlying.(*Struct); t != nil {
			for _, fld := range t.fields {
				if fld.name != "_" {
					assert(mset.insert(fld) == nil)
				}
			}
		}

		// Checker.Files may be called multiple times; additional package files
		// may add methods to already type-checked types. Add pre-existing methods
		// so that we can detect redeclarations.

		// jea update: we want to *allow* redeclarations at the repl
		// but we still do this to catch element and method
		// name collisions.
		for _, m := range base.methods {
			assert(m.name != "_")
			assert(mset.insert(m) == nil)
		}
	}

	// type-check methods
	for _, m := range methods {
		// spec: "For a base type, the non-blank names of methods bound
		// to it must be unique."
		if m.name != "_" {
			if alt := mset.insert(m); alt != nil {
				switch alt.(type) {
				case *Var:
					check.errorf(m.pos, "field and method with the same name %s", m.name)
				case *Func:
					if DisallowMethodRedefinition {
						check.errorf(m.pos,
							"method %s already declared for %s",
							m.name, obj)
					} else {
						// jea: allow function re-definition at the repl.
						pp("mset starting as: '%s'", mset.String())

						pp("doing mset.replace with m = '%s', and alt= '%s'", m, alt)
						prior := mset.replace(m)

						delete(check.ObjMap, prior)

						// jea: do we need to delete prior in check.Defs as well?
						// Hmm... we have an Object, not an *ast.Ident.
						//delete(check.Defs, prior)

						pp("check.Methods has len %v", len(check.Methods))
						for i, slc := range check.Methods {
							fmt.Printf("methods[%v] = \n", i)
							for j, fn := range slc {
								fmt.Printf("   [%v] = '%s'\n", j, fn.String())
							}
						}

						// need to delete the method from the type too.
						for i, curm := range base.methods {
							if curm == prior {
								base.methods = append(base.methods[:i], base.methods[i+1:]...)
								break
							}
						}

						alt = nil
						pp("mset is now, after replace(m): '%s'", mset.String())
						goto proceed
					}
				default:
					unreachable()
				}
				check.reportAltDecl(alt)
				continue
			}
		}

	proceed:
		// type-check
		check.objDecl(m, nil, nil) // spkg_tst3.Error here

		// methods with blank _ names cannot be found - don't keep them
		if base != nil && m.name != "_" {
			base.methods = append(base.methods, m)
		}
	}
}

func (check *Checker) funcDecl(obj *Func, decl *DeclInfo) {
	//pp("top of Checker.funcDecl, obj.Name()='%s', Type='%s', recv type='%#v'", obj.Name(), obj.Type(), decl.Fdecl.Recv.List[0].Type)

	receiverPrefix := ""
	if decl.Fdecl.Recv != nil && len(decl.Fdecl.Recv.List) > 0 {
		switch x := decl.Fdecl.Recv.List[0].Type.(type) {
		case *ast.Ident:
			//pp("receiver ident is '%s'", x.Name)
			receiverPrefix = x.Name + "."
		case *ast.StarExpr:
			switch y := x.X.(type) {
			case *ast.Ident:
				//pp("receiver ident is pointer to '%s'", y.Name)
				receiverPrefix = y.Name + "."
			}
		}
	}
	methodName := receiverPrefix + obj.Name() // for both methods and functions
	pp("methodName='%s'", methodName)

	assert(obj.typ == nil)

	// func declarations cannot use iota
	assert(check.iota == nil)

	sig := new(Signature)
	obj.typ = sig // guard against cycles
	fdecl := decl.Fdecl

	// jea is this a re-declaration?
	prior := check.scope.Lookup(obj.Name())
	if prior != nil {
		pp("prior was not nil for obj='%s', prior='%#v'", obj.Name(), prior)
	}

	// jea this is the call that defines new function declaration signatures!
	check.funcType(sig, fdecl.Recv, fdecl.Type, methodName)

	if sig.recv == nil && obj.name == "init" && (sig.params.Len() > 0 || sig.results.Len() > 0) {
		check.errorf(fdecl.Pos(), "func init must have no arguments and no return values")
		// ok to continue
	}

	// function body must be type-checked after global declarations
	// (functions implemented elsewhere have no body)
	if !check.conf.IgnoreFuncBodies && fdecl.Body != nil {
		check.later(obj.name, decl, sig, fdecl.Body)
	}
}

func (check *Checker) declStmt(decl ast.Decl) {
	pkg := check.pkg

	switch d := decl.(type) {
	case *ast.BadDecl:
		// ignore

	case *ast.GenDecl:
		var last *ast.ValueSpec // last ValueSpec with type or init exprs seen
		for iota, spec := range d.Specs {
			switch s := spec.(type) {
			case *ast.ValueSpec:
				switch d.Tok {
				case token.CONST:
					// determine which init exprs to use
					switch {
					case s.Type != nil || len(s.Values) > 0:
						last = s
					case last == nil:
						last = new(ast.ValueSpec) // make sure last exists
					}

					// declare all constants
					lhs := make([]*Const, len(s.Names))
					for i, name := range s.Names {
						obj := NewConst(name.Pos(), pkg, name.Name, nil, constant.MakeInt64(int64(iota)))
						lhs[i] = obj

						var init ast.Expr
						if i < len(last.Values) {
							init = last.Values[i]
						}

						check.constDecl(obj, last.Type, init)
					}

					check.arityMatch(s, last)

					// spec: "The scope of a constant or variable identifier declared
					// inside a function begins at the end of the ConstSpec or VarSpec
					// (ShortVarDecl for short variable declarations) and ends at the
					// end of the innermost containing block."
					scopePos := s.End()
					for i, name := range s.Names {
						check.declare(check.scope, name, lhs[i], scopePos)
					}

				case token.VAR:
					lhs0 := make([]*Var, len(s.Names))
					for i, name := range s.Names {
						lhs0[i] = NewVar(name.Pos(), pkg, name.Name, nil)
					}

					// initialize all variables
					for i, obj := range lhs0 {
						var lhs []*Var
						var init ast.Expr
						switch len(s.Values) {
						case len(s.Names):
							// lhs and rhs match
							init = s.Values[i]
						case 1:
							// rhs is expected to be a multi-valued expression
							lhs = lhs0
							init = s.Values[0]
						default:
							if i < len(s.Values) {
								init = s.Values[i]
							}
						}
						check.varDecl(obj, lhs, s.Type, init)
						if len(s.Values) == 1 {
							// If we have a single lhs variable we are done either way.
							// If we have a single rhs expression, it must be a multi-
							// valued expression, in which case handling the first lhs
							// variable will cause all lhs variables to have a type
							// assigned, and we are done as well.
							if debug {
								for _, obj := range lhs0 {
									assert(obj.typ != nil)
								}
							}
							break
						}
					}

					check.arityMatch(s, nil)

					// declare all variables
					// (only at this point are the variable scopes (parents) set)
					scopePos := s.End() // see constant declarations
					for i, name := range s.Names {
						// see constant declarations
						check.declare(check.scope, name, lhs0[i], scopePos)
					}

				default:
					check.invalidAST(s.Pos(), "invalid token %s", d.Tok)
				}

			case *ast.TypeSpec:
				obj := NewTypeName(s.Name.Pos(), pkg, s.Name.Name, nil)
				// spec: "The scope of a type identifier declared inside a function
				// begins at the identifier in the TypeSpec and ends at the end of
				// the innermost containing block."
				scopePos := s.Name.Pos()
				check.declare(check.scope, s.Name, obj, scopePos)
				check.typeDecl(obj, s.Type, nil, nil, s.Assign.IsValid())

			default:
				check.invalidAST(s.Pos(), "const, type, or var declaration expected")
			}
		}

	default:
		check.invalidAST(d.Pos(), "unknown ast.Decl node %T", d)
	}
}

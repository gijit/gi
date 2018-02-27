package compiler

import (
	"bytes"
	"fmt"
	"github.com/gijit/gi/pkg/ast"
	"github.com/gijit/gi/pkg/constant"
	"github.com/gijit/gi/pkg/printer"
	"github.com/gijit/gi/pkg/token"
	"github.com/gijit/gi/pkg/types"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"

	"github.com/gijit/gi/pkg/compiler/analysis"
	"github.com/gijit/gi/pkg/compiler/astutil"
	"github.com/gijit/gi/pkg/compiler/typesutil"
	"github.com/gijit/gi/pkg/verb"
)

var _ = debug.Stack

var debugNilCount int
var pp = verb.PP

type expression struct {
	str    string
	parens bool
}

func (e *expression) String() string {
	return e.str
}

func (e *expression) StringWithParens() string {
	if e.parens {
		return "(" + e.str + ")"
	}
	return e.str
}

// jea add
func (c *funcContext) exprToString(expr ast.Expr) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, c.p.fileSet, expr)
	return buf.String()
}

// desiredType can be nil. When present, for example, it guides
// the proper signed vs. unsigned translation of int,int64 types.
func (c *funcContext) translateExpr(expr ast.Expr, desiredType types.Type) (xprn *expression) {

	exprType := c.p.TypeOf(expr)
	desiredStr := "<nil>"
	if desiredType != nil {
		desiredStr = desiredType.String()
	}
	if expr == nil || c == nil {
		pp("000000 TOP OF gi TRANSLATE EXPR: jea debug, translateExpr(expr='<nil>')")
	} else {
		var exprTypeStr string
		if exprType == nil {
			exprTypeStr = "<nil>"
		} else {
			exprTypeStr = exprType.String()
		}
		pp("000000 TOP OF gi TRANSLATE EXPR: jea debug, translateExpr(expr='%s', exprType='%v', desiredType='%v'). c.p.Types[expr].Value='%#v', stack=\n%s\n", c.exprToString(expr), exprTypeStr, desiredStr, c.p.Types[expr].Value, "") //, string(debug.Stack()))
	}
	defer func() {
		//pp("returning from translateExpr.\n trace:\n%s\n", string(debug.Stack()))
		if xprn == nil {
			pp("444444 returning from gi translateExpr(): '<nil>'")
		} else {
			pp("444444 returning from gi translateExpr(): '%s', and c.output = '%#v'", xprn.str, string(c.output))
		}
	}()
	if value := c.p.Types[expr].Value; value != nil {
		basic := exprType.Underlying().(*types.Basic)
		switch {
		case isBoolean(basic):
			return c.formatExpr("%s", strconv.FormatBool(constant.BoolVal(value)))
		case isInteger(basic):

			// jea: for now, all int types are 64-bit
			// jea: TODO: come back and handle int32, uint32, int16, uint16, int8, uint8
			//if is64Bit(basic) {
			k := basic.Kind()
			if desiredType != nil {
				switch bk := desiredType.(type) {
				case *types.Basic:
					k = bk.Kind()
					pp("k = '%#v'/'%s'", k, k)
				}
			}
			pp("isInteger and k = '%v'", k) // k of 6 is int64, 2 is int.
			if k == types.Int64 || k == types.Int || k == types.Int32 || k == types.Int16 || k == types.Int8 {

				d, ok := constant.Int64Val(constant.ToInt(value))
				if !ok {
					panic("could not get exact int")
				}
				return c.formatExpr("%sLL", strconv.FormatInt(d, 10))
				//return c.formatExpr("new %s(%s, %s)", c.typeName(0, exprType), strconv.FormatInt(d>>32, 10), strconv.FormatUint(uint64(d)&(1<<32-1), 10))
			}

			d, ok := constant.Uint64Val(constant.ToInt(value))
			if !ok {
				panic("could not get exact uint")
			}
			return c.formatExpr("%sULL", strconv.FormatUint(d, 10))

			//return c.formatExpr("new %s(%s, %s)", c.typeName(0, exprType), strconv.FormatUint(d>>32, 10), strconv.FormatUint(d&(1<<32-1), 10))
			//}
			/*
				d, ok := constant.Int64Val(constant.ToInt(value))
				if !ok {
					panic("could not get exact int")
				}

				kd := basic.Kind()
				if kd == types.Int64 || kd == types.Int {
					return c.formatExpr("%sLL", strconv.FormatInt(d, 10))
				}
				if kd == types.Uint64 || kd == types.Uint {
					return c.formatExpr("%sULL", strconv.FormatUint(d, 10))
				}
			*/
		case isFloat(basic):
			f, _ := constant.Float64Val(value)
			return c.formatExpr("%s", strconv.FormatFloat(f, 'g', -1, 64))
		case isComplex(basic):
			r, _ := constant.Float64Val(constant.Real(value))
			i, _ := constant.Float64Val(constant.Imag(value))
			if basic.Kind() == types.UntypedComplex {
				exprType = types.Typ[types.Complex128]
			}
			// jea, luajit has native complex number syntax
			return c.formatExpr("%s+%si", strconv.FormatFloat(r, 'g', -1, 64), strconv.FormatFloat(i, 'g', -1, 64))
			//return c.formatExpr("new %s(%s, %s)", c.typeName(0, exprType), strconv.FormatFloat(r, 'g', -1, 64), strconv.FormatFloat(i, 'g', -1, 64))
		case isString(basic):
			pp("jea, in translateExpr(), value = '%v'", value)
			return c.formatExpr("%s", encodeString(constant.StringVal(value)))
		default:
			panic("Unhandled constant type: " + basic.String())
		}
	}

	pp("expr is type %T", expr)
	var obj types.Object
	switch e := expr.(type) {
	case *ast.SelectorExpr:
		obj = c.p.Uses[e.Sel]
	case *ast.Ident:
		pp("e is '%#v'", e)
		obj = c.p.Defs[e]
		if obj == nil {
			obj = c.p.Uses[e]
		}
	}
	pp("obj is '%#v'", obj)
	if obj == nil {
		debugNilCount++
		if debugNilCount > 2 {
			//panic("where obj nil?")
		}
	}

	if obj != nil && typesutil.IsJsPackage(obj.Pkg()) {
		switch obj.Name() {
		case "Global":
			return c.formatExpr("__global")
		case "Module":
			return c.formatExpr("__module")
		case "Undefined":
			return c.formatExpr("undefined")
		}
	}

	//pp("expressions.go:115, expr is '%#v'/Type=%T", expr, expr)
	switch e := expr.(type) {
	case *ast.CompositeLit:
		// jea: this depends on the initializers!
		pp("expressions.go we have an *ast.CompositeLit: '%#v'", e)
		if ptrType, isPointer := exprType.(*types.Pointer); isPointer {
			pp("isPointer is true")
			exprType = ptrType.Elem()
		}

		collectIndexedElements := func(elementType types.Type) []string {
			var elements []string
			i := 0
			zero := c.translateExpr(c.zeroValue(elementType), nil).String()
			for _, element := range e.Elts {
				if kve, isKve := element.(*ast.KeyValueExpr); isKve {
					key, ok := constant.Int64Val(constant.ToInt(c.p.Types[kve.Key].Value))
					if !ok {
						panic("could not get exact int")
					}
					i = int(key)
					element = kve.Value
				}
				for len(elements) <= i {
					elements = append(elements, zero)
				}
				elements[i] = c.translateImplicitConversionWithCloning(element, elementType).String()
				i++
			}
			return elements
		}

		pp("in expressions.go, at underlying switch, exprType='%T'/'%#v'", exprType, exprType)
		switch t := exprType.Underlying().(type) {
		case *types.Array:
			elements := collectIndexedElements(t.Elem())

			// print anon_arrayType immedaitely
			c.TypeNameSetting = IMMEDIATE
			anonTypeName := c.typeName(0, exprType)
			if len(elements) == 0 {
				pp("expressions.go about to call typeName(t) on t='%#v'", t)
				tn := c.typeName(0, t)
				pp("expressions.go making array of size %v with tn='%v' from t='%#v'", t.Len(), tn, t)
				// gijit below, but try to go back to gophrejs style __zero() call.
				//ret := c.formatExpr("%s.zero()", tn)
				ret := c.formatExpr("%s()", tn) // will do zero for us *and* set proper metatable.
				pp("about to return array zero() for 0 len array, c.output='%s'", string(c.output))
				return ret
				//return c.formatExpr(fmt.Sprintf(`__gi_NewArray({}, "%s", %v, %s)`, typeKind(t.Elem()), t.Len(), zero))
			}
			zero := c.translateExpr(c.zeroValue(t.Elem()), nil).String()
			for len(elements) < int(t.Len()) {
				elements = append(elements, zero)
			}
			return c.formatExpr(fmt.Sprintf(`%s({[0]=%s})`, anonTypeName, strings.Join(elements, ", "))) // , typeKind(t.Elem()), t.Len(), zero))
			//return c.formatExpr(`__toNativeArray(%s, {[0]=%s})`, typeKind(t.Elem()), strings.Join(elements, ", "))

		case *types.Slice:
			//zero := c.translateExpr(c.zeroValue(t.Elem()), nil).String()
			ele := strings.Join(collectIndexedElements(t.Elem()), ", ")
			if len(ele) > 0 {
				// jea: do 0-based indexing of slices, not 1-based.
				ele = "[0]=" + ele
			}
			return c.formatExpr("%s({%s})", c.typeName(0, exprType), ele)
			//return c.formatExpr(fmt.Sprintf(`_gi_NewSlice("%s",{%s}, %s)`, c.typeName(0, t.Elem()), ele, zero))
			//return c.formatExpr("new %s([%s])", c.typeName(0, exprType), strings.Join(collectIndexedElements(t.Elem()), ", "))
		case *types.Map:
			entries := make([]string, len(e.Elts))
			for i, element := range e.Elts {
				kve := element.(*ast.KeyValueExpr)
				entries[i] = fmt.Sprintf(`[%s]=%s`, c.translateImplicitConversionWithCloning(kve.Key, t.Key()), c.translateImplicitConversionWithCloning(kve.Value, t.Elem()))
			}
			joined := strings.Join(entries, ", ")
			pp("joined = '%#v'", joined)
			return c.formatExpr("__makeMap({%s}, %s, %s, %s)", joined, c.typeName(0, t.Key()), c.typeName(0, t.Elem()), c.typeName(0, exprType))
		case *types.Struct:
			pp("in expressions.go, for *types.Struct")
			elements := make([]string, t.NumFields())
			isKeyValue := true
			if len(e.Elts) != 0 {
				_, isKeyValue = e.Elts[0].(*ast.KeyValueExpr)
			}
			if !isKeyValue {
				for i, element := range e.Elts {
					elements[i] = c.translateImplicitConversionWithCloning(element, t.Field(i).Type()).String()
				}
			}
			if isKeyValue {
				for i := range elements {
					elements[i] = c.translateExpr(c.zeroValue(t.Field(i).Type()), nil).String()
				}
				for _, element := range e.Elts {
					kve := element.(*ast.KeyValueExpr)
					for j := range elements {
						if kve.Key.(*ast.Ident).Name == t.Field(j).Name() {
							elements[j] = c.translateImplicitConversionWithCloning(kve.Value, t.Field(j).Type()).String()
							break
						}
					}
				}
			}
			//flds := structFieldTypes(t)
			sele := "nil"
			if len(elements) > 0 {
				sele = strings.Join(elements, ", ")
			}

			return c.formatExpr("%s.ptrToNewlyConstructed(%s)", c.typeName(0, exprType), sele)
			//return c.formatExpr("%s(%s)", c.typeName(0, exprType), sele)
			//return c.formatExpr("%s({}, %s)", c.typeName(0, exprType), sele)

			// first lua attempt:
			//vals := structFieldNameValuesForLua(t, elements)
			//return c.formatExpr(`__reg:NewInstance("%s",{%s})`, c.typeName(0, exprType), strings.Join(vals, ", "))
			// original:
			//return c.formatExpr("new %s.__ptr(%s)", c.typeName(0, exprType), strings.Join(elements, ", "))
		default:
			panic(fmt.Sprintf("Unhandled CompositeLit type: %T\n", t))
		}

	case *ast.FuncLit:
		pp("expressions.go:213 we have a *ast.FuncLit: '%#v'", e)
		_, fun, _ := translateFunction(e.Type, nil, e.Body, c, exprType.(*types.Signature), c.p.FuncLitInfos[e], "", false)
		if len(c.p.escapingVars) != 0 {
			names := make([]string, 0, len(c.p.escapingVars))
			for obj := range c.p.escapingVars {
				names = append(names, c.p.objectNames[obj])
			}
			sort.Strings(names)
			list := strings.Join(names, ", ")
			return c.formatExpr("(function(%s) { return %s; })(%s)", list, fun, list)
		}
		return c.formatExpr("(%s)", fun)

	case *ast.UnaryExpr:
		pp("we have UnaryExpr:  *ast.UnaryExpr: '%#v', with typeName='%s'", e, c.typeName(0, c.p.TypeOf(e)))
		t := c.p.TypeOf(e.X)
		switch e.Op {
		case token.AND:
			pp("we have token.AND, exprType = '%#v", exprType)
			if typesutil.IsJsObject(exprType) {
				return c.formatExpr("%e.object", e.X)
			}

			switch t.Underlying().(type) {
			case *types.Struct, *types.Array:
				pp("underlying is struct or array")
				te := c.translateExpr(e.X, nil)
				pp("after translateExpr on e.X underlying struct or array, te = '%#v'", te)
				pp("after translateExpr on underlying struct or array, type of t = '%#v'", t)
				// jea: gopherjs didn't represent struct values directly?? try
				// wrapping with a newDataPointer or other pointer generating construct...
				//return c.formatExpr("%s(%s)", c.typeName(0, c.p.TypeOf(e)), te)
				// Arg. It turns out:
				// The problem with the above comes  when comparing points to a struct... two
				// points to the same struct have to be equal. Test 121 in struct_test.go.

				return c.translateExpr(e.X, nil)
				// gopherjs:
				// return c.translateExpr(e.X)
			}
			pp("underlying is not struct nor array")

			switch x := astutil.RemoveParens(e.X).(type) {
			case *ast.CompositeLit:
				return c.formatExpr("__newDataPointer(%e, %s)", x, c.typeName(0, c.p.TypeOf(e)))
			case *ast.Ident:
				obj := c.p.Uses[x].(*types.Var)
				if c.p.escapingVars[obj] {
					return c.formatExpr("(%2s(function() return this.__target[0]; end, function(__v) this.__target[0] = __v; end, %1s))", c.p.objectNames[obj], c.typeName(0, exprType))
					// return c.formatExpr("(%1s.__ptr || (%1s.__ptr = new %2s(function() { return this.__target[0]; }, function(__v) { this.__target[0] = __v; }, %1s)))", c.p.objectNames[obj], c.typeName(exprType))
				}

				// basic taking address of value to get pointer. See ptr_test.go, test 099.
				//return c.formatExpr(`__ptrType(function() return %1s; end, function(v) %2s; end, "%s")`, c.objectName(obj), c.translateAssign(x, c.newIdent("v", exprType), false), starToAmp(exprType.String()))

				pp("basic taking address of value to get pointer...")
				pp("c.typeName(0, exprType) = '%v'", c.typeName(0, exprType))
				pp("exprType = '%#v'", exprType)
				return c.formatExpr(`%2s(function() return %3s; end, function(__v) %4s end, %3s)`, c.varPtrName(obj), c.typeName(0, exprType), c.objectName(obj), c.translateAssign(x, c.newIdent("__v", exprType), false))
				//return c.formatExpr(`%2s({}, function() return %3s; end, function(__v) %4s end, %3s)`, c.varPtrName(obj), c.typeName(0, exprType), c.objectName(obj), c.translateAssign(x, c.newIdent("__v", exprType), false))

			case *ast.SelectorExpr:
				sel, ok := c.p.SelectionOf(x)
				if !ok {
					// qualified identifier
					obj := c.p.Uses[x.Sel].(*types.Var)
					return c.formatExpr(`(%2s(function() return %3s; end, function(__v) %4s end))`, c.varPtrName(obj), c.typeName(0, exprType), c.objectName(obj), c.translateAssign(x, c.newIdent("__v", exprType), false))
					//return c.formatExpr(`(%1s || (%1s = %2s(function() { return %3s; }, function(__v) { %4s })))`, c.varPtrName(obj), c.typeName(0, exprType), c.objectName(obj), c.translateAssign(x, c.newIdent("__v", exprType), false))
				}
				newSel := &ast.SelectorExpr{X: c.newIdent("this.__target", c.p.TypeOf(x.X)), Sel: x.Sel}
				c.setType(newSel, exprType)
				c.p.additionalSelections[newSel] = sel
				return c.formatExpr("(%3s(function() return %4e; end, function(__v) %5s end, %1e))", x.X, x.Sel.Name, c.typeName(0, exprType), newSel, c.translateAssign(newSel, c.newIdent("__v", exprType), false))
				//return c.formatExpr("(%1e.__ptr_%2s || (%1e.__ptr_%2s = %3s(function() { return %4e; }, function(__v) { %5s }, %1e)))", x.X, x.Sel.Name, c.typeName(0, exprType), newSel, c.translateAssign(newSel, c.newIdent("__v", exprType), false))
			case *ast.IndexExpr:
				if _, ok := c.p.TypeOf(x.X).Underlying().(*types.Slice); ok {
					return c.formatExpr("__indexPtr(%1e.__array, %1e.__offset + %2e, %3s)", x.X, x.Index, c.typeName(0, exprType))
				}
				return c.formatExpr("__indexPtr(%e, %e, %s)", x.X, x.Index, c.typeName(0, exprType))
			case *ast.StarExpr:
				return c.translateExpr(x.X, nil)
			default:
				panic(fmt.Sprintf("Unhandled: %T\n", x))
			}

		case token.ARROW:
			call := &ast.CallExpr{
				Fun:  c.newIdent("__recv", types.NewSignature(nil, types.NewTuple(types.NewVar(0, nil, "", t)), types.NewTuple(types.NewVar(0, nil, "", exprType), types.NewVar(0, nil, "", types.Typ[types.Bool])), false)),
				Args: []ast.Expr{e.X},
			}
			c.Blocking[call] = true
			if _, isTuple := exprType.(*types.Tuple); isTuple {
				return c.formatExpr("%e", call)
			}
			// jea:  what is this doing???
			// return c.formatExpr("%e[1]", call)
			return c.formatExpr("")
		}

		basic := t.Underlying().(*types.Basic)
		switch e.Op {
		case token.ADD:
			tx := c.translateExpr(e.X, nil)
			pp("e.Op is token.ADD at expressions.go:283, tx='%s'", tx)
			return tx
		case token.SUB:
			switch {
			case is64Bit(basic):
				return c.formatExpr("%1s(-%2h, -%2l)", c.typeName(0, t), e.X)
			case isComplex(basic):
				return c.formatExpr("-(%1e)", e.X)
				//return c.formatExpr("-(%1r+%1i)", e.X)
				//return c.formatExpr("%1s(-%2r, -%2i)", c.typeName(0, t), e.X)
			case isUnsigned(basic):
				return c.fixNumber(c.formatExpr("-%e", e.X), basic)
			default:
				return c.formatExpr("-%e", e.X)
			}
		case token.XOR:
			if is64Bit(basic) {
				return c.formatExpr("%1s(~%2h, ~%2l >>> 0)", c.typeName(0, t), e.X)
			}
			return c.fixNumber(c.formatExpr("~%e", e.X), basic)
		case token.NOT:
			return c.formatExpr(" not %e", e.X)
		default:
			panic(e.Op)
		}

	case *ast.BinaryExpr:
		if e.Op == token.NEQ {
			return c.formatExpr(" not (%s)", c.translateExpr(&ast.BinaryExpr{
				X:  e.X,
				Op: token.EQL,
				Y:  e.Y,
			}, nil))
		}

		t := c.p.TypeOf(e.X)
		t2 := c.p.TypeOf(e.Y)
		_, isInterface := t2.Underlying().(*types.Interface)
		if isInterface || types.Identical(t, types.Typ[types.UntypedNil]) {
			t = t2
		}

		if basic, isBasic := t.Underlying().(*types.Basic); isBasic && isNumeric(basic) {
			if isComplex(basic) {
				switch e.Op {
				case token.EQL:
					return c.formatExpr("((%1e) == (%2e))", e.X, e.Y)
				case token.ADD, token.SUB:
					return c.formatExpr("((%1e) %3t (%2e))", e.X, e.Y, e.Op)
					//return c.formatExpr("%3s(%1r %4t %2r, %1i %4t %2i)", e.X, e.Y, c.typeName(0, t), e.Op)
				case token.MUL:
					return c.formatExpr("%3s((%1e) * (%2e))", e.X, e.Y, c.typeName(0, t))
				case token.QUO:
					return c.formatExpr("((%1e) / (%2e))", e.X, e.Y)
				default:
					panic(e.Op)
				}
			}

			switch e.Op {
			case token.EQL:
				return c.formatParenExpr("((%e) == (%e))", e.X, e.Y)
			case token.LSS, token.LEQ, token.GTR, token.GEQ:
				return c.formatExpr("((%e) %t (%e))", e.X, e.Op, e.Y)
			case token.ADD, token.SUB:
				pp("token.ADD or SUB,calling c.formatExpr, e.X='%s', op='%v', e.Y='%s'", c.exprToString(e.X), e.Op.String(), c.exprToString(e.Y))
				xx := c.formatExpr("((%e) %t (%e))", e.X, e.Op, e.Y)
				pp("token.ADD or SUB,calling c.fixNumber(xx) with xx = '%s'", x2s(xx))
				return c.fixNumber(xx, basic)
			case token.MUL:
				switch basic.Kind() {
				case types.Int64, types.Int,
					types.Uint64, types.Uint, types.Uintptr:
					return c.formatParenExpr("((%e) * (%e))", e.X, e.Y)
				}
				return c.fixNumber(c.formatExpr("((%e) * (%e))", e.X, e.Y), basic)
			case token.QUO:
				if isInteger(basic) {
					// cut off decimals
					// shift := ">>"
					//if isUnsigned(basic) {
					//	shift = ">>>"
					//}

					// jea:
					return c.formatExpr(`__integerByZeroCheck(%1e / %2e)`, e.X, e.Y)
					// return c.formatExpr(`(%1s = %2e / %3e, (%1s == %1s && %1s ~= 1/0 && %1s ~= -1/0) ? %1s %4s 0 : error("integer divide by zero"))`, c.newVariable("_q"), e.X, e.Y, shift)
				}
				if basic.Kind() == types.Float32 {
					return c.fixNumber(c.formatExpr("((%e) / (%e))", e.X, e.Y), basic)
				}
				return c.formatExpr("((%e) / (%e))", e.X, e.Y)
			case token.REM:
				return c.formatExpr(`__integerByZeroCheck(%1e %% %2e)`, e.X, e.Y)
			case token.SHL, token.SHR:
				op := e.Op.String()
				if e.Op == token.SHR && isUnsigned(basic) {
					op = ">>>"
				}
				if v := c.p.Types[e.Y].Value; v != nil {
					i, _ := constant.Uint64Val(constant.ToInt(v))
					if i >= 32 {
						return c.formatExpr("0")
					}
					return c.fixNumber(c.formatExpr("%e %s %s", e.X, op, strconv.FormatUint(i, 10)), basic)
				}
				if e.Op == token.SHR && !isUnsigned(basic) {
					return c.fixNumber(c.formatParenExpr("%e >> __min(%f, 31)", e.X, e.Y), basic)
				}
				y := c.newVariable("y")
				return c.fixNumber(c.formatExpr("(%s = %f, %s < 32 ? (%e %s %s) : 0)", y, e.Y, y, e.X, op, y), basic)
			case token.AND, token.OR:
				if isUnsigned(basic) {
					return c.formatParenExpr("(%e %t %e) >>> 0", e.X, e.Op, e.Y)
				}
				return c.formatParenExpr("%e %t %e", e.X, e.Op, e.Y)
			case token.AND_NOT:
				return c.fixNumber(c.formatParenExpr("%e & ~%e", e.X, e.Y), basic)
			case token.XOR:
				return c.fixNumber(c.formatParenExpr("%e ^ %e", e.X, e.Y), basic)
			default:
				panic(e.Op)
			}
		}

		switch e.Op {
		case token.ADD, token.LSS, token.LEQ, token.GTR, token.GEQ:
			// jea:
			if e.Op == token.ADD {
				if b, isBasic := c.p.TypeOf(e.X).Underlying().(*types.Basic); isBasic && isString(b) {
					// string concat is ".." in Lua, rather than + as in Go.
					return c.formatExpr("%e .. %e", e.X, e.Y)
				}
			}
			return c.formatExpr("%e %t %e", e.X, e.Op, e.Y)
		case token.LAND:
			if c.Blocking[e.Y] {
				skipCase := c.caseCounter
				c.caseCounter++
				resultVar := c.newVariable("_v")
				c.Printf("if (not (%s)) then %s = false; __s = %d; goto ::s::; end", c.translateExpr(e.X, nil), resultVar, skipCase)
				c.Printf("%s = %s; case %d:", resultVar, c.translateExpr(e.Y, nil), skipCase)
				//c.Printf("%s = %s; case %d:", resultVar, c.translateExpr(e.Y, nil), skipCase)
				return c.formatExpr("%s", resultVar)
			}
			return c.formatExpr("%e and %e", e.X, e.Y)
		case token.LOR:
			if c.Blocking[e.Y] {
				skipCase := c.caseCounter
				c.caseCounter++
				resultVar := c.newVariable("_v")
				// was "continue s;"
				c.Printf("if (%s) then %s = true; __s = %d; goto ::s::; end", c.translateExpr(e.X, nil), resultVar, skipCase)
				c.Printf("%s = %s; case %d:", resultVar, c.translateExpr(e.Y, nil), skipCase)
				return c.formatExpr("%s", resultVar)
			}
			return c.formatExpr("%e or %e", e.X, e.Y)
		case token.EQL:
			switch u := t.Underlying().(type) {
			case *types.Array, *types.Struct:
				return c.formatExpr("__equal(%e, %e, %s)", e.X, e.Y, c.typeName(0, t))
			case *types.Interface:
				pp("e.Y='%#v'", e.Y)
				return c.formatExpr("__interfaceIsEqual(%s, %s)", c.translateImplicitConversion(e.X, t), c.translateImplicitConversion(e.Y, t))
			case *types.Pointer:
				if _, ok := u.Elem().Underlying().(*types.Array); ok {
					return c.formatExpr("__equal(%s, %s, %s)", c.translateImplicitConversion(e.X, t), c.translateImplicitConversion(e.Y, t), c.typeName(0, u.Elem()))
				}
			case *types.Basic:
				if isBoolean(u) {
					if b, ok := analysis.BoolValue(e.X, c.p.Info.Info); ok && b {
						return c.translateExpr(e.Y, nil)
					}
					if b, ok := analysis.BoolValue(e.Y, c.p.Info.Info); ok && b {
						return c.translateExpr(e.X, nil)
					}
				}
			}
			x := c.formatExpr("%s == %s", c.translateImplicitConversion(e.X, t), c.translateImplicitConversion(e.Y, t))
			//fmt.Printf("debug: at an == formating... x='%#v'\n", x)
			//panic("where?")
			return x
		default:
			panic(e.Op)
		}

	case *ast.ParenExpr:
		return c.formatParenExpr("%e", e.X)

	case *ast.IndexExpr:
		switch t := c.p.TypeOf(e.X).Underlying().(type) {
		case *types.Array, *types.Pointer:
			pattern := rangeCheck("%1e[%2f]", c.p.Types[e.Index].Value != nil, true)
			if _, ok := t.(*types.Pointer); ok { // check pointer for nix (attribute getter causes a panic)
				pattern = `(%1e.nilCheck, ` + pattern + `)`
			}
			return c.formatExpr(pattern, e.X, e.Index)
		case *types.Slice:
			// jea: we really do need these runtime range checks. Lua will
			//   just give us nils back silently instead of complaining about
			//   out of bounds access.
			pp("expressions.go:484 slice, e.X='%#v', e.Index='%#v'", e.X, e.Index)
			return c.formatExpr(rangeCheck("%1e[%2f]", c.p.Types[e.Index].Value != nil, false), e.X, e.Index)
			// return c.formatExpr(rangeCheck("%1e.__array[%1e.__offset + %2f]", c.p.Types[e.Index].Value != nil, false), e.X, e.Index)
		case *types.Map:
			if typesutil.IsJsObject(c.p.TypeOf(e.Index)) {
				c.p.errList = append(c.p.errList, types.Error{Fset: c.p.fileSet, Pos: e.Index.Pos(), Msg: "cannot use js.Object as map key"})
			}
			key := fmt.Sprintf("%s", c.translateImplicitConversion(e.Index, t.Key()))
			pp("t.Key()='%T'/%#v", t.Key(), t.Key())
			switch bas := t.Key().(type) {
			case *types.Basic:
				if bas.Kind() != types.String {
					key = `"` + key + `"`
				}
			}
			//key := fmt.Sprintf("%s.keyFor(%s)", c.typeName(0, t.Key()), c.translateImplicitConversion(e.Index, t.Key()))
			if _, isTuple := exprType.(*types.Tuple); isTuple {

				return c.formatExpr(` %1e('get', %2s, %3e) `, e.X, key, c.zeroValue(t.Elem()))

				// gi example
				// {m('get',0, zerovalue)}

				// js example:
				/*	_tuple = (_entry = m[0], _entry !== undefined ? [_entry.v, true] : [0, false]);
					a = _tuple[0];
					ok = _tuple[1];
				*/
				// produced by (well, just the 1st line starting with _entry):
				// return c.formatExpr(`(%1s = %2e[%3s], %1s !== undefined ? [%1s.v, true] : [%4e, false])`, c.newVariable("_entry"), e.X, key, c.zeroValue(t.Elem()))
			}
			//return c.formatExpr(fmt.Sprintf(`%s('get', %s)`, e.X, key))

			return c.formatExpr(` %1e('get', %2s, %3e) `, e.X, key, c.zeroValue(t.Elem()))

			//return c.formatExpr(`(%1s = %2e[%3s], %1s !== undefined ? %1s.v : %4e)`, c.newVariable("_entry"), e.X, key, c.zeroValue(t.Elem()))
		case *types.Basic:
			return c.formatExpr("__utf8.sub(%e,%f+1,%f+1)", e.X, e.Index, e.Index)
		default:
			panic(fmt.Sprintf("Unhandled IndexExpr: %T\n", t))
		}

	case *ast.SliceExpr:
		pp("expressions.go:529 we have an *ast.SliceExpr: '%#v'", e)
		if b, isBasic := c.p.TypeOf(e.X).Underlying().(*types.Basic); isBasic && isString(b) {
			switch {
			// e is a slice expression, the slice is from [Low:High).
			case e.Low == nil && e.High == nil:
				return c.translateExpr(e.X, nil)
			case e.Low == nil:
				return c.formatExpr("string.sub(%e, 1, %f)", e.X, e.High)
				//return c.formatExpr("__substring(%e, 0, %f)", e.X, e.High)
			case e.High == nil:
				return c.formatExpr("string.sub(%e, %f+1)", e.X, e.Low)
				//return c.formatExpr("__substring(%e, %f)", e.X, e.Low)
			default:
				return c.formatExpr("string.sub(%e, %f+1, %f)", e.X, e.Low, e.High)
				//return c.formatExpr("__substring(%e, %f, %f)", e.X, e.Low, e.High)
			}
		}
		slice := c.translateConversionToSlice(e.X, exprType)
		switch {
		case e.Low == nil && e.High == nil:
			return c.formatExpr("%s", slice)
		case e.Low == nil:
			if e.Max != nil {
				return c.formatExpr("__subslice(%s, 0, %f, %f)", slice, e.High, e.Max)
			}
			return c.formatExpr("__subslice(%s, 0, %f)", slice, e.High)
		case e.High == nil:
			return c.formatExpr("__subslice(%s, %f)", slice, e.Low)
		default:
			if e.Max != nil {
				return c.formatExpr("__subslice(%s, %f, %f, %f)", slice, e.Low, e.High, e.Max)
			}
			return c.formatExpr("__subslice(%s, %f, %f)", slice, e.Low, e.High)
		}

	case *ast.SelectorExpr:
		sel, ok := c.p.SelectionOf(e)
		if !ok {
			// qualified identifier
			// 99999 fmt-tracking from Sprintf alone.
			pp("expressions.go:586, sel = '%#v', e='%#v'", sel, e) // sel is nil, e is not
			pp("jea expressions.go:589, our ast.SelectorExpr.X='%#v', .Sel='%#v'", e.X, e.Sel)
			// e.X.Name == "fmt", e.Sel.Name == "Sprintf"; both are *ast.Ident
			return c.formatExpr("%s", c.objectName(obj))
		}

		switch sel.Kind() {
		case types.FieldVal:
			fields, jsTag := c.translateSelection(sel, e.Pos())
			if jsTag != "" {
				if _, ok := sel.Type().(*types.Signature); ok {
					return c.formatExpr("__internalize(%1e.%2s.%3s, %4s, %1e.%2s)", e.X, strings.Join(fields, "."), jsTag, c.typeName(0, sel.Type()))
				}
				return c.internalize(c.formatExpr("%e.%s.%s", e.X, strings.Join(fields, "."), jsTag), sel.Type())
			}
			return c.formatExpr("%e.%s", e.X, strings.Join(fields, "."))
		case types.MethodVal:
			sel, _ := c.p.SelectionOf(e)
			recvType := sel.Recv()
			pp("case types.MethodVal: typeName of e.X = '%#v'", c.typeName(0, recvType))
			recvr := c.makeReceiver(e)
			return c.formatExpr(`__methodVal(%s, "%s", "%s")`, recvr, sel.Obj().(*types.Func).Name(), c.typeName(0, recvType))
			//return c.formatExpr(`__methodVal(%s, "%s")`, c.makeReceiver(e), sel.Obj().(*types.Func).Name())
		case types.MethodExpr:
			if !sel.Obj().Exported() {
				c.p.dependencies[sel.Obj()] = true
			}
			if _, ok := sel.Recv().Underlying().(*types.Interface); ok {
				return c.formatExpr(`__ifaceMethodExpr("%s")`, sel.Obj().(*types.Func).Name())
			}
			return c.formatExpr(`__methodExpr(%s, "%s")`, c.typeName(0, sel.Recv()), sel.Obj().(*types.Func).Name())
		default:
			panic(fmt.Sprintf("unexpected sel.Kind(): %T", sel.Kind()))
		}

	case *ast.CallExpr:
		plainFun := astutil.RemoveParens(e.Fun)

		if astutil.IsTypeExpr(plainFun, c.p.Info.Info) {
			return c.formatExpr("(%s)", c.translateConversion(e.Args[0], c.p.TypeOf(plainFun)))
		}

		sig := c.p.TypeOf(plainFun).Underlying().(*types.Signature)

		switch f := plainFun.(type) {
		case *ast.Ident:
			obj := c.p.Uses[f]
			if o, ok := obj.(*types.Builtin); ok {
				return c.translateBuiltin(o.Name(), sig, e.Args, e.Ellipsis.IsValid(), exprType)
			}
			if typesutil.IsJsPackage(obj.Pkg()) && obj.Name() == "InternalObject" {
				return c.translateExpr(e.Args[0], nil)
			}
			return c.translateCall(e, sig, c.translateExpr(f, nil))

		case *ast.SelectorExpr:
			sel, ok := c.p.SelectionOf(f)
			if !ok {
				// qualified identifier
				obj := c.p.Uses[f.Sel]
				if typesutil.IsJsPackage(obj.Pkg()) {
					switch obj.Name() {
					case "Debugger":
						return c.formatExpr("debugger")
					case "InternalObject":
						return c.translateExpr(e.Args[0], nil)
					}
				}
				return c.translateCall(e, sig, c.translateExpr(f, nil))
			}

			externalizeExpr := func(e ast.Expr) string {
				t := c.p.TypeOf(e)
				if types.Identical(t, types.Typ[types.UntypedNil]) {
					return "null"
				}
				return c.externalize(c.translateExpr(e, nil).String(), t)
			}
			externalizeArgs := func(args []ast.Expr) string {
				s := make([]string, len(args))
				for i, arg := range args {
					s[i] = externalizeExpr(arg)
				}
				return strings.Join(s, ", ")
			}

			switch sel.Kind() {
			case types.MethodVal:
				recv := c.makeReceiver(f)
				declaredFuncRecv := sel.Obj().(*types.Func).Type().(*types.Signature).Recv().Type()
				if typesutil.IsJsObject(declaredFuncRecv) {
					globalRef := func(id string) string {
						if recv.String() == "__global" && id[0] == '_' && len(id) > 1 {
							return id
						}
						return recv.String() + "." + id
					}
					switch sel.Obj().Name() {
					case "Get":
						if id, ok := c.identifierConstant(e.Args[0]); ok {
							return c.formatExpr("%s", globalRef(id))
						}
						return c.formatExpr("%s[__externalize(%e, __type__string)]", recv, e.Args[0])
					case "Set":
						if id, ok := c.identifierConstant(e.Args[0]); ok {
							return c.formatExpr("%s = %s", globalRef(id), externalizeExpr(e.Args[1]))
						}
						return c.formatExpr("%s[__externalize(%e, __String)] = %s", recv, e.Args[0], externalizeExpr(e.Args[1]))
					case "Delete":
						return c.formatExpr("delete %s[__externalize(%e, __String)]", recv, e.Args[0])
					case "Length":
						return c.formatExpr("__parseInt(%s.length)", recv)
					case "Index":
						return c.formatExpr("%s[%e]", recv, e.Args[0])
					case "SetIndex":
						return c.formatExpr("%s[%e] = %s", recv, e.Args[0], externalizeExpr(e.Args[1]))
					case "Call":
						if id, ok := c.identifierConstant(e.Args[0]); ok {
							if e.Ellipsis.IsValid() {
								objVar := c.newVariable("obj")
								return c.formatExpr("(%s = %s, %s.%s.apply(%s, %s))", objVar, recv, objVar, id, objVar, externalizeExpr(e.Args[1]))
							}
							return c.formatExpr("%s(%s)", globalRef(id), externalizeArgs(e.Args[1:]))
						}
						if e.Ellipsis.IsValid() {
							objVar := c.newVariable("obj")
							return c.formatExpr("(%s = %s, %s[__externalize(%e, __String)].apply(%s, %s))", objVar, recv, objVar, e.Args[0], objVar, externalizeExpr(e.Args[1]))
						}
						return c.formatExpr("%s[__externalize(%e, __String)](%s)", recv, e.Args[0], externalizeArgs(e.Args[1:]))
					case "Invoke":
						if e.Ellipsis.IsValid() {
							return c.formatExpr("%s.apply(undefined, %s)", recv, externalizeExpr(e.Args[0]))
						}
						return c.formatExpr("%s(%s)", recv, externalizeArgs(e.Args))
					case "New":
						if e.Ellipsis.IsValid() {
							return c.formatExpr("(__global.Function.prototype.bind.apply(%s, [undefined].concat(%s)))", recv, externalizeExpr(e.Args[0]))
						}
						return c.formatExpr("(%s)(%s)", recv, externalizeArgs(e.Args))
					case "Bool":
						return c.internalize(recv, types.Typ[types.Bool])
					case "String":
						return c.internalize(recv, types.Typ[types.String])
					case "Int":
						return c.internalize(recv, types.Typ[types.Int])
					case "Int64":
						return c.internalize(recv, types.Typ[types.Int64])
					case "Uint64":
						return c.internalize(recv, types.Typ[types.Uint64])
					case "Float":
						return c.internalize(recv, types.Typ[types.Float64])
					case "Interface":
						return c.internalize(recv, types.NewInterface(nil, nil))
					case "Unsafe":
						return recv
					default:
						panic("Invalid js package object: " + sel.Obj().Name())
					}
				}

				methodName := sel.Obj().Name()
				if reservedKeywords[methodName] {
					methodName += "_"
				}

				isLuar := typesutil.IsLuarObject(declaredFuncRecv)

				//fmt.Printf("\n isLuar='%v', recv = '%#v', declaredFuncRecv = '%#v'\n",
				//	isLuar, recv, declaredFuncRecv)

				if isLuar {
					// jea: then change back to . for Luar cdata methods. imp_test 087
					return c.translateCall(e, sig, c.formatExpr("%s.%s", recv, methodName))
				}
				// jea: change to object:method call for Lua.
				return c.translateCall(e, sig, c.formatExpr("%s:%s", recv, methodName))
			case types.FieldVal:
				fields, jsTag := c.translateSelection(sel, f.Pos())
				if jsTag != "" {
					call := c.formatExpr("%e.%s.%s(%s)", f.X, strings.Join(fields, "."), jsTag, externalizeArgs(e.Args))
					switch sig.Results().Len() {
					case 0:
						return call
					case 1:
						return c.internalize(call, sig.Results().At(0).Type())
					default:
						c.p.errList = append(c.p.errList, types.Error{Fset: c.p.fileSet, Pos: f.Pos(), Msg: "field with js tag can not have func type with multiple results"})
					}
				}
				return c.translateCall(e, sig, c.formatExpr("%e.%s", f.X, strings.Join(fields, ".")))

			case types.MethodExpr:
				return c.translateCall(e, sig, c.translateExpr(f, nil))

			default:
				panic(fmt.Sprintf("unexpected sel.Kind(): %T", sel.Kind()))
			}
		default:
			return c.translateCall(e, sig, c.translateExpr(plainFun, nil))
		}

	case *ast.StarExpr:
		if typesutil.IsJsObject(c.p.TypeOf(e.X)) {
			return c.formatExpr("__jsObjectPtr(%e)", e.X)
		}
		if c1, isCall := e.X.(*ast.CallExpr); isCall && len(c1.Args) == 1 {
			if c2, isCall := c1.Args[0].(*ast.CallExpr); isCall && len(c2.Args) == 1 && types.Identical(c.p.TypeOf(c2.Fun), types.Typ[types.UnsafePointer]) {
				if unary, isUnary := c2.Args[0].(*ast.UnaryExpr); isUnary && unary.Op == token.AND {
					return c.translateExpr(unary.X, nil) // unsafe conversion
				}
			}
		}
		switch exprType.Underlying().(type) {
		case *types.Struct, *types.Array:
			return c.translateExpr(e.X, nil)
		}
		// jea: pointer dereference
		//return c.formatExpr("%e[0]", e.X) // any key will do.
		return c.formatExpr("%e.__get()", e.X)

	case *ast.TypeAssertExpr:
		if e.Type == nil {
			return c.translateExpr(e.X, nil)
		}
		t := c.p.TypeOf(e.Type)
		if _, isTuple := exprType.(*types.Tuple); isTuple {
			// jea, type assertion place 2; face_test 101 goes here.
			// return both converted-interface-value, and ok.
			return c.formatExpr(`__assertType(%e, %s, 2)`, e.X, c.typeName(0, t))
			//return c.formatExpr("__assertType(%e, %s, true)", e.X, c.typeName(0, t))
		}
		// jea, type assertion place 0: only return value, without the 2nd 'ok' return.
		return c.formatExpr(`__assertType(%e, %s, 0)`, e.X, c.typeName(0, t))

	case *ast.Ident:
		if e.Name == "_" {
			panic("Tried to translate underscore identifier.")
		}
		pp("under *ast.Ident, obj='%#v'/%T", obj, obj)
		switch o := obj.(type) {
		case *types.Var, *types.Const:
			pp("jea debug, line 937 expressions.go")
			return c.formatExpr("%s", c.objectName(o))
		case *types.Func:
			return c.formatExpr("%s", c.objectName(o))
		case *types.TypeName:
			return c.formatExpr("%s", c.typeName(0, o.Type()))
		case *types.Nil:
			if typesutil.IsJsObject(exprType) {
				return c.formatExpr("null")
			}
			switch t := exprType.Underlying().(type) {
			case *types.Basic:
				if t.Kind() != types.UnsafePointer {
					panic("unexpected basic type")
				}
				return c.formatExpr("0")
			case *types.Slice, *types.Pointer:
				return c.formatExpr("%s.__nil", c.typeName(0, exprType))
			case *types.Chan:
				return c.formatExpr("__chanNil")
			case *types.Map:
				return c.formatExpr("false")
			case *types.Interface:
				return c.formatExpr("nil")
			case *types.Signature:
				return c.formatExpr("__throwNilPointerError")
			default:
				panic(fmt.Sprintf("unexpected type: %T", t))
			}
		default:
			pp("e = '%#v', o = '%#v'", e, o)
			panic(fmt.Sprintf("Unhandled object: %T\n", o))
		}

	case nil:
		return c.formatExpr("")
	case *ast.BasicLit:
		pp("expressions.go:818 we have an *ast.BasicLit: '%#v'", e)
		// JEA DEBUG: we added this case. what to do here?

		return &expression{
			str: e.Value, // JEA: Guessing, might not be right.
		}
	default:
		panic(fmt.Sprintf("Unhandled expression: %T\n", e))

	}
}

func (c *funcContext) translateCall(e *ast.CallExpr, sig *types.Signature, fun *expression) *expression {
	pp("top of translateCall, len(e.Args)='%v', e.Args='%#v'", len(e.Args), e.Args)
	for i := range e.Args {
		pp("top of translateCall, e.Args[i=%v]='%#v'", i, e.Args[i])
	}
	args := c.translateArgs(sig, e.Args, e.Ellipsis.IsValid())
	if !c.Blocking[e] {
		joined := strings.Join(args, ", ")
		pp("c.Blocking[e] is false, joined = '%v'", joined)
		return c.formatExpr("%s(%s)", fun, joined)
	}

	pp("c.Blocking[e] is true")
	// jea
	//resumeCase := c.caseCounter
	c.caseCounter++
	returnVar := "_r"
	if sig.Results().Len() != 0 {
		returnVar = c.newVariable("_r")
	}
	// jea
	c.Printf(" %[1]s = %[2]s(%[3]s);", returnVar, fun, strings.Join(args, ", "))
	// hmm... tests fail with this extra scheduler call:
	//c.Printf(" %[1]s = %[2]s(%[3]s); __task.scheduler();", returnVar, fun, strings.Join(args, ", "))
	// jea debug:
	//c.Printf("/*jea expressions.go:873*/ %[1]s = %[2]s(%[3]s);", returnVar, fun, strings.Join(args, ", "))

	//c.Printf("%[1]s = %[2]s(%[3]s); /* */ __s = %[4]d; case %[4]d: if(__c) { __c = false; %[1]s = %[1]s.__blk(); } if (%[1]s && %[1]s.__blk !== undefined) { break s; }", returnVar, fun, strings.Join(args, ", "), resumeCase)
	if sig.Results().Len() != 0 {
		return c.formatExpr("%s", returnVar)
	}
	return c.formatExpr("")

}

func (c *funcContext) makeReceiver(e *ast.SelectorExpr) *expression {
	sel, _ := c.p.SelectionOf(e)
	if !sel.Obj().Exported() {
		c.p.dependencies[sel.Obj()] = true
	}

	x := e.X
	recvType := sel.Recv()
	if len(sel.Index()) > 1 {
		for _, index := range sel.Index()[:len(sel.Index())-1] {
			if ptr, isPtr := recvType.(*types.Pointer); isPtr {
				recvType = ptr.Elem()
			}
			s := recvType.Underlying().(*types.Struct)
			recvType = s.Field(index).Type()
		}

		fakeSel := &ast.SelectorExpr{X: x, Sel: ast.NewIdent("o")}
		c.p.additionalSelections[fakeSel] = &fakeSelection{
			kind:  types.FieldVal,
			recv:  sel.Recv(),
			index: sel.Index()[:len(sel.Index())-1],
			typ:   recvType,
		}
		x = c.setType(fakeSel, recvType)
	}

	_, isPointer := recvType.Underlying().(*types.Pointer)
	methodsRecvType := sel.Obj().Type().(*types.Signature).Recv().Type()
	_, pointerExpected := methodsRecvType.(*types.Pointer)
	if !isPointer && pointerExpected {
		recvType = types.NewPointer(recvType)
		x = c.setType(&ast.UnaryExpr{Op: token.AND, X: x}, recvType)
	}
	if isPointer && !pointerExpected {
		x = c.setType(x, methodsRecvType)
	}

	recv := c.translateImplicitConversionWithCloning(x, methodsRecvType)
	if isWrapped(recvType) {
		recv = c.formatExpr("%s(%s)", c.typeName(0, methodsRecvType), recv)
	}
	return recv
}

func (c *funcContext) translateBuiltin(name string, sig *types.Signature, args []ast.Expr, ellipsis bool, exprType types.Type) *expression {
	switch name {
	case "new":
		t := sig.Results().At(0).Type().(*types.Pointer)
		if c.p.Pkg.Path() == "syscall" && types.Identical(t.Elem().Underlying(), types.Typ[types.Uintptr]) {
			return c.formatExpr("__newByteArray(16)")
		}
		switch t.Elem().Underlying().(type) {
		case *types.Struct, *types.Array:
			return c.formatExpr("%e", c.zeroValue(t.Elem()))
		default:
			return c.formatExpr("__newDataPointer(%e, %s)", c.zeroValue(t.Elem()), c.typeName(0, t))
		}
	case "make":
		switch argType := c.p.TypeOf(args[0]).Underlying().(type) {
		case *types.Slice:
			t := c.typeName(0, c.p.TypeOf(args[0]))
			//zero := c.zeroValue(argType.Elem())
			if len(args) == 3 {
				return c.formatExpr("__makeSlice(%s, %f, %f)",
					t, args[1], args[2])
			}
			return c.formatExpr("__makeSlice(%s, %f)",
				t, args[1])

		case *types.Map:
			//if len(args) == 2 && c.p.Types[args[1]].Value == nil {
			//	return c.formatExpr(`((%1f < 0 || %1f > 2147483647) ? error("makemap: size out of range") : {})`, args[1])
			//}

			// return c.formatExpr("{}") // gopherjs
			// jea: we at least need a metatable for tostring, etc... so this minimal won't do.
			//return c.formatExpr("{}")

			// __makeMap(keyForFunc, entries, keyType, elemType, mapType)
			t := exprType.Underlying().(*types.Map)
			return c.formatExpr("__makeMap({}, %s, %s, %s)", c.typeName(0, t.Key()), c.typeName(0, t.Elem()), c.typeName(0, exprType))
		case *types.Chan:
			length := "0"
			if len(args) == 2 {
				length = c.formatExpr("%f", args[1]).String()
			}
			return c.formatExpr("__Chan(%s, %s)", c.typeName(0, c.p.TypeOf(args[0]).Underlying().(*types.Chan).Elem()), length)
		default:
			panic(fmt.Sprintf("Unhandled make type: %T\n", argType))
		}
	case "len":
		switch argType := c.p.TypeOf(args[0]).Underlying().(type) {
		case *types.Basic:
			//return c.formatExpr("%e.length", args[0])
			return c.formatExpr(" #%e", args[0])
		case *types.Slice:
			return c.formatExpr(" #%e", args[0])
		case *types.Pointer:
			return c.formatExpr("(%e, %d)", args[0], argType.Elem().(*types.Array).Len())
		case *types.Map:
			return c.formatExpr(" #%e", args[0])
		case *types.Chan:
			return c.formatExpr("%e.__buffer.length", args[0])
		// length of array is constant
		default:
			panic(fmt.Sprintf("Unhandled len type: %T\n", argType))
		}
	case "cap":
		switch argType := c.p.TypeOf(args[0]).Underlying().(type) {
		case *types.Slice, *types.Chan:
			return c.formatExpr("%e.__capacity", args[0])
		case *types.Pointer:
			return c.formatExpr("(%e, %d)", args[0], argType.Elem().(*types.Array).Len())
		// capacity of array is constant
		default:
			panic(fmt.Sprintf("Unhandled cap type: %T\n", argType))
		}
	case "panic":
		return c.formatExpr("panic(%s)", c.translateImplicitConversion(args[0], types.NewInterface(nil, nil)))
	case "append":
		if ellipsis || len(args) == 1 {
			argStr := c.translateArgs(sig, args, ellipsis)
			return c.formatExpr("__appendSlice(%s, %s)", argStr[0], argStr[1])
		}
		sliceType := sig.Results().At(0).Type().Underlying().(*types.Slice)
		return c.formatExpr("append(%e, %s)", args[0], strings.Join(c.translateExprSlice(args[1:], sliceType.Elem()), ", "))
	case "delete":
		keyType := c.p.TypeOf(args[0]).Underlying().(*types.Map).Key()
		pp("delete, keyType='%v'", keyType)
		return c.formatExpr(`%e("delete",%s)`, args[0], c.translateImplicitConversion(args[1], keyType))
	case "copy":
		if basic, isBasic := c.p.TypeOf(args[1]).Underlying().(*types.Basic); isBasic && isString(basic) {
			return c.formatExpr("__copyString(%e, %e)", args[0], args[1])
		}
		return c.formatExpr("__copySlice(%e, %e)", args[0], args[1])
	case "print", "println":
		return c.formatExpr("print(%s)", strings.Join(c.translateExprSlice(args, nil), ", "))
	case "complex":
		argStr := c.translateArgs(sig, args, ellipsis)
		return c.formatExpr("%s(%s, %s)", c.typeName(0, sig.Results().At(0).Type()), argStr[0], argStr[1])
	case "real":
		// there is already a LuaJIT real() function
		// available, from complex.lua, can we just cut straight to that?
		return c.formatExpr("real(%e)", args[0])
	case "imag":
		// there is already a LuaJIT imag() function
		// available, from complex.lua, can we just cut straight to that?
		return c.formatExpr("imag(%e)", args[0])
	case "recover":
		return c.formatExpr("recover()")
	case "close":
		return c.formatExpr(`__close(%e)`, args[0])
	default:
		panic(fmt.Sprintf("Unhandled builtin: %s\n", name))
	}
}

func (c *funcContext) identifierConstant(expr ast.Expr) (string, bool) {
	val := c.p.Types[expr].Value
	if val == nil {
		return "", false
	}
	s := constant.StringVal(val)
	if len(s) == 0 {
		return "", false
	}
	for i, c := range s {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (i > 0 && c >= '0' && c <= '9') || c == '_') {
			return "", false
		}
	}
	return s, true
}

func (c *funcContext) translateExprSlice(exprs []ast.Expr, desiredType types.Type) []string {
	parts := make([]string, len(exprs))
	for i, expr := range exprs {
		parts[i] = c.translateImplicitConversion(expr, desiredType).String()
	}
	return parts
}

func (c *funcContext) translateConversion(expr ast.Expr, desiredType types.Type) *expression {

	// jea debug
	var buf bytes.Buffer
	fset := token.NewFileSet()
	printer.Fprint(&buf, fset, expr)
	//pp("debug: translate Conversion sees expr: '%v'.\n", string(buf.Bytes()))

	exprType := c.p.TypeOf(expr)
	if types.Identical(exprType, desiredType) {
		return c.translateExpr(expr, nil)
	}

	if c.p.Pkg.Path() == "reflect" {
		if call, isCall := expr.(*ast.CallExpr); isCall && types.Identical(c.p.TypeOf(call.Fun), types.Typ[types.UnsafePointer]) {
			if ptr, isPtr := desiredType.(*types.Pointer); isPtr {
				if named, isNamed := ptr.Elem().(*types.Named); isNamed {
					switch named.Obj().Name() {
					case "arrayType", "chanType", "funcType", "interfaceType", "mapType", "ptrType", "sliceType", "structType":
						return c.formatExpr("%e.kindType", call.Args[0]) // unsafe conversion
					default:
						return c.translateExpr(expr, nil)
					}
				}
			}
		}
	}

	switch t := desiredType.Underlying().(type) {
	case *types.Basic:
		switch {
		case isInteger(t):
			basicExprType := exprType.Underlying().(*types.Basic)
			switch {
			case is64Bit(t):
				if !is64Bit(basicExprType) {
					if basicExprType.Kind() == types.Uintptr { // this might be an Object returned from reflect.Value.Pointer()
						return c.formatExpr("%1s(%2e)", c.typeName(0, desiredType), expr)
						//return c.formatExpr("%1s(0, %2e.constructor == Number ? %2e : 1)", c.typeName(0, desiredType), expr)
					}
					return c.formatExpr("%s(%e)", c.typeName(0, desiredType), expr)
				}
				return c.formatExpr("%s(%e)", c.typeName(0, desiredType), expr)
			case is64Bit(basicExprType):
				if !isUnsigned(t) && !isUnsigned(basicExprType) {
					return c.formatParenExpr("%e", expr)
				}
				return c.formatExpr("%s", c.translateExpr(expr, nil), t)
			case isFloat(basicExprType):
				// jea
				//return c.formatParenExpr("%e >> 0", expr)
				return c.formatParenExpr("int(%e)", expr)
			case types.Identical(exprType, types.Typ[types.UnsafePointer]):
				return c.translateExpr(expr, nil)
			default:
				// 201 not here
				return c.fixNumber(c.translateExpr(expr, nil), t)
			}
		case isFloat(t):
			if t.Kind() == types.Float32 && exprType.Underlying().(*types.Basic).Kind() == types.Float64 {
				return c.formatExpr("__fround(%e)", expr) // fround returns the nearest 32-bit single precision float representation of a Number.
			}
			return c.formatExpr("tonumber(%f)", expr)
		case isComplex(t):
			return c.formatExpr("%1s(%2e)", c.typeName(0, desiredType), expr)
		case isString(t):
			value := c.translateExpr(expr, nil)
			switch et := exprType.Underlying().(type) {
			case *types.Basic:
				if is64Bit(et) {
					value = c.formatExpr("%s.__low", value)
				}
				if isNumeric(et) {
					return c.formatExpr("__encodeRune(%s)", value)
				}
				return value
			case *types.Slice:
				if types.Identical(et.Elem().Underlying(), types.Typ[types.Rune]) {
					return c.formatExpr("__runesToString(%s)", value)
				}

				return c.formatExpr("__bytesToString(%s)", value)
			default:
				panic(fmt.Sprintf("Unhandled conversion: %v\n", et))
			}
		case t.Kind() == types.UnsafePointer:
			if unary, isUnary := expr.(*ast.UnaryExpr); isUnary && unary.Op == token.AND {
				if indexExpr, isIndexExpr := unary.X.(*ast.IndexExpr); isIndexExpr {
					return c.formatExpr("__sliceToArray(%s)", c.translateConversionToSlice(indexExpr.X, types.NewSlice(types.Typ[types.Uint8])))
				}
				if ident, isIdent := unary.X.(*ast.Ident); isIdent && ident.Name == "_zero" {
					return c.formatExpr("__newByteArray(0)")
				}
			}
			if ptr, isPtr := c.p.TypeOf(expr).(*types.Pointer); c.p.Pkg.Path() == "syscall" && isPtr {
				if s, isStruct := ptr.Elem().Underlying().(*types.Struct); isStruct {
					array := c.newVariable("_array")
					target := c.newVariable("_struct")
					// jea: sizes32 -> sizes64
					c.Printf("%s = __newByteArray(%d);", array, sizes64.Sizeof(s)) // jea
					c.Delayed(func() {
						c.Printf("%s = %s, %s;", target, c.translateExpr(expr, nil), c.loadStruct(array, target, s))
					})
					return c.formatExpr("%s", array)
				}
			}
			if call, ok := expr.(*ast.CallExpr); ok {
				if id, ok := call.Fun.(*ast.Ident); ok && id.Name == "new" {
					// jea: sizes32 -> sizes64
					return c.formatExpr("__newByteArray(%d)", int(sizes64.Sizeof(c.p.TypeOf(call.Args[0]))))
				}
			}
		}

	case *types.Slice:
		switch et := exprType.Underlying().(type) {
		case *types.Basic:
			if isString(et) {
				if types.Identical(t.Elem().Underlying(), types.Typ[types.Rune]) {
					return c.formatExpr("%s(__stringToRunes(%e))", c.typeName(0, desiredType), expr)
				}
				return c.formatExpr("%s(__stringToBytes(%e))", c.typeName(0, desiredType), expr)
			}
		case *types.Array, *types.Pointer:
			return c.formatExpr("%s(%e)", c.typeName(0, desiredType), expr)
		}

	case *types.Pointer:
		switch u := t.Elem().Underlying().(type) {
		case *types.Array:
			return c.translateExpr(expr, nil)
		case *types.Struct:
			if c.p.Pkg.Path() == "syscall" && types.Identical(exprType, types.Typ[types.UnsafePointer]) {
				array := c.newVariable("_array")
				target := c.newVariable("_struct")
				return c.formatExpr("(%s = %e, %s = %e, %s, %s)", array, expr, target, c.zeroValue(t.Elem()), c.loadStruct(array, target, u), target)
			}
			return c.formatExpr("__pointerOfStructConversion(%e, %s)", expr, c.typeName(0, t))
		}

		if !types.Identical(exprType, types.Typ[types.UnsafePointer]) {
			exprTypeElem := exprType.Underlying().(*types.Pointer).Elem()
			ptrVar := c.newVariable("_ptr")
			getterConv := c.translateConversion(c.setType(&ast.StarExpr{X: c.newIdent(ptrVar, exprType)}, exprTypeElem), t.Elem())
			setterConv := c.translateConversion(c.newIdent("__v", t.Elem()), exprTypeElem)
			return c.formatExpr("(%3s(function() return %4s; end, function(__v)  %1s.__set(%5s); end, %1s.__target))", ptrVar, expr, c.typeName(0, desiredType), getterConv, setterConv)
		}

	case *types.Interface:
		if types.Identical(exprType, types.Typ[types.UnsafePointer]) {
			return c.translateExpr(expr, nil)
		}
	}

	return c.translateImplicitConversionWithCloning(expr, desiredType)
}

func (c *funcContext) translateImplicitConversionWithCloning(expr ast.Expr, desiredType types.Type) *expression {
	pp("translateImplicitConversionWithCloning(expr='%#v', desiredType='%s', at: '%s'", expr, desiredType, "") // string(debug.Stack()))
	switch desiredType.Underlying().(type) {
	case *types.Struct, *types.Array:
		// either a struct or an array
		switch expr.(type) {
		case nil, *ast.CompositeLit:
			// nothing
		default:
			// this is the __gi_clone that is called for value receivers on methods.
			// __clone in gohperjs.
			// And for passing an array to a function argument (by-value of course).
			pp("debug __clone arg: c.typeName(0, desiredType)='%s'", c.typeName(0, desiredType))

			typName, isAnon, anonType, createdNm := c.typeNameWithAnonInfo(desiredType)
			pp("debug __clone arg: c.typeName(0, desiredType)='%s'; createdNm='%s'; isAnon='%v', anonType='%#v'", typName, createdNm, isAnon, anonType)
			if isAnon {
				return c.formatExpr(`__clone(%e, %s)`, expr, c.typeName(0, anonType.Type()))
			} else {
				return c.formatExpr(`__clone(%e, %s)`, expr, typName)
			}

		}
	}

	return c.translateImplicitConversion(expr, desiredType)
}

func (c *funcContext) translateImplicitConversion(expr ast.Expr, desiredType types.Type) *expression {
	pp("translateImplicitConversion top: desiredType='%#v', expr='%#v'\n", desiredType, expr)

	if desiredType == nil {
		pp("YYY 1 translateImplicitConversion exiting early on desiredType == nil")
		return c.translateExpr(expr, nil)
	}

	exprType := c.p.TypeOf(expr)
	pp("exprType = '%v'", exprType)
	if types.Identical(exprType, desiredType) {
		pp("YYY 2 translateImplicitConversion exiting early, b/c types are identical, exprType='%#v' and desiredType='%#v'. expr ='%#v'", exprType, desiredType, expr)
		ret := c.translateExpr(expr, nil)
		pp("end of YYY 2, ret = '%#v'/'%s'", ret, ret)
		return ret
	}

	basicExprType, isBasicExpr := exprType.Underlying().(*types.Basic)
	if isBasicExpr && basicExprType.Kind() == types.UntypedNil {
		pp("YYY 3 translateImplicitConversion exiting early")
		return c.formatExpr("%e", c.zeroValue(desiredType))
	}

	switch desiredType.Underlying().(type) {
	case *types.Slice:
		pp("YYY 4 translateImplicitConversion exiting early")
		return c.formatExpr("__subslice(%1s(%2e.__array), %2e.__offset, %2e.__offset + %2e.__length)", c.typeName(0, desiredType), expr)

	case *types.Interface:
		if typesutil.IsJsObject(exprType) {
			pp("YYY 5 translateImplicitConversion exiting early")
			// wrap JS object into js.Object struct when converting to interface
			return c.formatExpr("__jsObjectPtr(%e)", expr)
		}
		if isWrapped(exprType) {
			pp("isWrapped is true for exprType='%#v'", exprType)

			pp("YYY 6 translateImplicitConversion exiting early")
			// jea
			// string literals are converting to new `new String("string")`
			// which we don't need.
			// likewise, function references have junk wrapped around them.
			// Example: fmt.Printf -> new funcType(fmt.Printf)
			//return c.formatExpr("%s(%e)", c.typeName(0, exprType), expr)

			// references to functions arrive here.
			return c.formatExpr("%e", expr)
		}
		pp("!isWrapped for exprType='%#v'", exprType)
		if _, isStruct := exprType.Underlying().(*types.Struct); isStruct {
			pp("YYY 7 translateImplicitConversion exiting early")
			//return c.formatExpr("%1e.__constructor.__elem(%1e)", expr)
			//return c.formatExpr("%1e.__typ.elem(%1e)", expr)
			return c.formatExpr("(%1e)", expr)
		}
	}
	pp("bottom of expressions.go:1250 calling c.translateExpr, for expr='%#v', exprType='%v'", expr, exprType)
	return c.translateExpr(expr, desiredType)
}

func (c *funcContext) translateConversionToSlice(expr ast.Expr, desiredType types.Type) *expression {
	// try reverting back got the GopherJS code:
	switch c.p.TypeOf(expr).Underlying().(type) {
	case *types.Array, *types.Pointer:
		return c.formatExpr("%s(%e)", c.typeName(0, desiredType), expr)
	}
	return c.translateExpr(expr, nil)

	/* version 1, pre tsys.lua fullport of gopherJS type system.
	typeExpr := c.p.TypeOf(expr)
	switch x := typeExpr.Underlying().(type) {
	case *types.Pointer:
		// currently just a copy of the below
		et := x.Elem()
		pp("array to slice conversion, desiredType = '%#v', typeExpr='%#v', x='%#v', et='%s'", desiredType, typeExpr, x, et.String())
		zero := c.translateExpr(c.zeroValue(x.Elem()), nil).String()
		return c.formatExpr(`_gi_NewSlice("%s", %e, %s)`, c.typeName(0, et), expr, zero)
		//return c.formatExpr(`_gi_NewSlice("%s", %e)`, c.typeName(0, desiredType), expr)
	case *types.Array:
		et := x.Elem()
		pp("array to slice conversion, desiredType = '%#v', typeExpr='%#v', x='%#v', et='%s'", desiredType, typeExpr, x, et.String())
		zero := c.translateExpr(c.zeroValue(x.Elem()), nil).String()
		// c.typeName(0, desiredType)
		//return c.formatExpr(`_gi_NewSlice("%s", %e)`, c.typeName(0, et), expr)
		return c.formatExpr(`_gi_NewSlice("%s", %e, %s)`, c.typeName(0, et), expr, zero)
	}
	return c.translateExpr(expr, nil)
	*/
}

func (c *funcContext) loadStruct(array, target string, s *types.Struct) string {
	view := c.newVariable("_view")
	code := fmt.Sprintf("%s = new DataView(%s.buffer, %s.byteOffset)", view, array, array)
	var fields []*types.Var
	var collectFields func(s *types.Struct, path string)
	collectFields = func(s *types.Struct, path string) {
		for i := 0; i < s.NumFields(); i++ {
			field := s.Field(i)
			if fs, isStruct := field.Type().Underlying().(*types.Struct); isStruct {
				collectFields(fs, path+"."+fieldName(s, i))
				continue
			}
			fields = append(fields, types.NewVar(0, nil, path+"."+fieldName(s, i), field.Type()))
		}
	}
	collectFields(s, target)
	// jea: sizes32 -> sizes64
	offsets := sizes64.Offsetsof(fields)
	for i, field := range fields {
		switch t := field.Type().Underlying().(type) {
		case *types.Basic:
			if isNumeric(t) {
				code += fmt.Sprintf(", %s = %s.get%s(%d, true)", field.Name(), view, toJavaScriptType(t), offsets[i])
			}
		case *types.Array:
			code += fmt.Sprintf(", %s = __newAnyArrayValue(%s, __min(%s.byteOffset + %d, #%s.buffer)) -- expressions.go:1498 jea: almost surely wrong, figure out context to fix this.\n", field.Name(), typeKind(t.Elem()), array, offsets[i], array)
			//code += fmt.Sprintf(`, %s = new (__nativeArray(%s))(%s.buffer, __min(%s.byteOffset + %d, %s.buffer.byteLength))`, field.Name(), typeKind(t.Elem()), array, array, offsets[i], array)
		}
	}
	return code
}

func x2s(x *expression) string {
	if x == nil {
		return "<nil>"
	}
	return x.str
}

func (c *funcContext) fixNumber(value *expression, basic *types.Basic) (xprn *expression) {
	pp("top of fixNumber with value='%s'", x2s(value))
	defer func() {
		pp("returning from fixNumber with xprn='%s'", x2s(xprn))
	}()
	switch basic.Kind() {
	case types.Int8:
		// jea
		fallthrough
		//return c.formatParenExpr("%s << 24 >> 24", value)
	case types.Uint8:
		// jea
		fallthrough
		//return c.formatParenExpr("%s << 24 >>> 24", value)
	case types.Int16:
		// jea
		fallthrough
		//return c.formatParenExpr("%s << 16 >> 16", value)
	case types.Uint16:
		// jea
		fallthrough
		//return c.formatParenExpr("%s << 16 >>> 16", value)
	case types.Uint32, types.Uint, types.Uintptr:
		// jea
		fallthrough
		//return c.formatParenExpr("%s >>> 0", value)
	case types.Int32, types.Int, types.Int64, types.UntypedInt:
		// jea
		//return c.formatParenExpr("%s >> 0", value)
		return c.formatParenExpr("%s", value)
	case types.Float32:
		// jea:
		return value
		//return c.formatExpr("__fround(%s)", value) // // fround returns the nearest 32-bit single precision float representation of a Number.
	case types.Float64:
		return value
	default:
		panic(fmt.Sprintf("fixNumber: unhandled basic.Kind(): %s", basic.String()))
	}
}

func (c *funcContext) asFloat64(value *expression, basic *types.Basic) (xprn *expression) {
	pp("top of asFloat64 with value='%s'", x2s(value))
	defer func() {
		pp("returning from asFloat64 with xprn='%s'", x2s(xprn))
	}()
	switch basic.Kind() {
	case types.Int8,
		types.Uint8,
		types.Int16,
		types.Uint16,
		types.Uint32,
		types.Uint,
		types.Uintptr,
		types.Int32,
		types.Int,
		types.Int64,
		types.UntypedInt,
		types.Float32:
		return c.formatParenExpr("tonumber(%s)", value)

		//case types.Float32:
		//		//return value
		//		//
		//		return c.formatParenExpr("tonumber(%s)", value)
		// return c.formatExpr("__fround(%s)", value) // // fround returns the nearest 32-bit single precision float representation of a Number.

	case types.Float64:
		return value

	default:
		panic(fmt.Sprintf("asFloat64: unhandled basic.Kind(): %s", basic.String()))
	}
}

func (c *funcContext) internalize(s *expression, t types.Type) *expression {
	if typesutil.IsJsObject(t) {
		return s
	}
	switch u := t.Underlying().(type) {
	case *types.Basic:
		switch {
		case isBoolean(u):
			return c.formatExpr("not not(%s)", s)
			//return c.formatExpr("!!(%s)", s)
		case isInteger(u) && !is64Bit(u):
			return c.fixNumber(c.formatExpr("__parseInt(%s)", s), u)
		case isFloat(u):
			return c.formatExpr("__parseFloat(%s)", s)
		}
	}
	return c.formatExpr("__internalize(%s, %s)", s, c.typeName(0, t))
}

func (c *funcContext) formatExpr(format string, a ...interface{}) *expression {
	//
	return c.formatExprInternal(format, a, false)
}

func (c *funcContext) formatParenExpr(format string, a ...interface{}) *expression {
	pp("top of formatParenExpr")
	return c.formatExprInternal(format, a, true)
}

func (c *funcContext) formatExprInternal(format string, a []interface{}, parens bool) (xprn *expression) {
	pp("111111 top of formatExprInternal(), format='%s', parens='%v', len(a)=%v, a='%#v'.", format, parens, len(a), a)
	defer func() {
		if xprn == nil {
			pp("222222 returning from formatExprInternal, xrpn='<nil>'")
		} else {
			pp("222222 returning from formatExprInternal, xrpn='%s'", xprn.str)
		}
	}()
	defer func() {
		if xprn != nil {
			pp("formatExprInternal('%s') returning '%s'", format, xprn.str)
			//if xprn.str == `r:Get` {
			//	panic("where?")
			//}
		} else {
			pp("expressions.go:1357, formatExprInternal('%s') returning nil", format)
		}
	}()
	/*
		for i := range a {
			x, isX := a[i].(*expression)
			if isX {
				pp("a[i=%v] = '%v'", i, x)
				//if x.str == "r:Get" {
				//	panic("where?")
				//}
			}
		}
	*/
	processFormat := func(f func(uint8, uint8, int)) {
		n := 0
		for i := 0; i < len(format); i++ {
			b := format[i]
			if b == '%' {
				i++
				k := format[i]
				if k >= '0' && k <= '9' {
					n = int(k - '0' - 1)
					i++
					k = format[i]
				}
				f(0, k, n)
				n++
				continue
			}
			f(b, 0, 0)
		}
	}

	counts := make([]int, len(a))
	processFormat(func(b, k uint8, n int) {
		switch k {
		case 'e', 'f', 'h', 'l', 'r', 'i':
			counts[n]++
		}
	})

	out := bytes.NewBuffer(nil)
	vars := make([]string, len(a))
	hasAssignments := false
	for i, e := range a {
		if counts[i] <= 1 {
			continue
		}
		if _, isIdent := e.(*ast.Ident); isIdent {
			continue
		}
		if val := c.p.Types[e.(ast.Expr)].Value; val != nil {
			continue
		}
		if !hasAssignments {
			hasAssignments = true
			out.WriteByte('(')
			parens = false
		}
		v := c.newVariable("x")
		out.WriteString(v + " = " + c.translateExpr(e.(ast.Expr), nil).String() + ", ")
		vars[i] = v
	}

	processFormat(func(b, k uint8, n int) {
		writeExprWithSuffix := func(suffix string) {
			if vars[n] != "" {
				out.WriteString(vars[n] + suffix)
				return
			}
			stuff := c.translateExpr(a[n].(ast.Expr), nil).StringWithParens() + suffix
			pp("jea debug, stuff='%#v', where suffix was '%s'\n", stuff, suffix)
			out.WriteString(stuff)
		}
		switch k {
		case 0:
			out.WriteByte(b)
		case 's':
			if e, ok := a[n].(*expression); ok {
				out.WriteString(e.StringWithParens())
				return
			}
			out.WriteString(a[n].(string))
		case 'd':
			out.WriteString(strconv.Itoa(a[n].(int)))
		case 't':
			out.WriteString(a[n].(token.Token).String())
		case 'e':
			e := a[n].(ast.Expr)
			if val := c.p.Types[e].Value; val != nil {
				out.WriteString(c.translateExpr(e, nil).String())
				return
			}
			writeExprWithSuffix("")
		case 'f':
			e := a[n].(ast.Expr)
			if val := c.p.Types[e].Value; val != nil {
				d, _ := constant.Int64Val(constant.ToInt(val))
				out.WriteString(strconv.FormatInt(d, 10))
				return
			}
			if is64Bit(c.p.TypeOf(e).Underlying().(*types.Basic)) {
				out.WriteString("__flatten64(")
				writeExprWithSuffix("")
				out.WriteString(")")
				return
			}
			writeExprWithSuffix("")
		case 'h':
			panic("jea: we shouldn't be asking for h format")
			e := a[n].(ast.Expr)
			if val := c.p.Types[e].Value; val != nil {
				d, _ := constant.Uint64Val(constant.ToInt(val))
				if c.p.TypeOf(e).Underlying().(*types.Basic).Kind() == types.Int64 {
					out.WriteString(strconv.FormatInt(int64(d)>>32, 10))
					return
				}
				out.WriteString(strconv.FormatUint(d>>32, 10))
				return
			}
			writeExprWithSuffix(".__high")
		case 'l':
			panic("jea: we shouldn't be asking for l format")
			if val := c.p.Types[a[n].(ast.Expr)].Value; val != nil {
				d, _ := constant.Uint64Val(constant.ToInt(val))
				out.WriteString(strconv.FormatUint(d&(1<<32-1), 10))
				return
			}
			writeExprWithSuffix(".__low")
		case 'r':
			if val := c.p.Types[a[n].(ast.Expr)].Value; val != nil {
				r, _ := constant.Float64Val(constant.Real(val))
				out.WriteString(strconv.FormatFloat(r, 'g', -1, 64))
				return
			}
			out.WriteString("real(")
			writeExprWithSuffix("")
			out.WriteString(")")
			//writeExprWithSuffix(".__real")
		case 'i':
			if val := c.p.Types[a[n].(ast.Expr)].Value; val != nil {
				i, _ := constant.Float64Val(constant.Imag(val))
				out.WriteString(strconv.FormatFloat(i, 'g', -1, 64))
				return
			}
			out.WriteString("imag(")
			writeExprWithSuffix("")
			out.WriteString(")")
			//writeExprWithSuffix(".__imag")
		case '%':
			out.WriteRune('%')
		default:
			panic(fmt.Sprintf("formatExpr: %%%c%d", k, n))
		}
	})

	if hasAssignments {
		out.WriteByte(')')
	}
	// 99999 fmt-tracking from Sprintf alone.
	return &expression{str: out.String(), parens: parens}
}

func structFieldTypes(t *types.Struct) (r []string) {
	n := t.NumFields()
	for i := 0; i < n; i++ {
		fld := t.Field(i)
		r = append(r, fmt.Sprintf(`["%s"]="%s"`, fld.Name(), fld.Type()))
	}
	return
}

func structFieldNameValuesForLua(t *types.Struct, ele []string) (r []string) {
	n := t.NumFields()
	if len(ele) != n {
		panic(fmt.Sprintf("internal error, field count %v == n != len(ele)=%v", n, len(ele)))
	}
	for i := 0; i < n; i++ {
		fld := t.Field(i)
		r = append(r, fmt.Sprintf(`["%s"]=%s`, fld.Name(), ele[i]))
	}
	return
}

func starToAmp(a string) string {
	if len(a) == 0 {
		return a
	}
	if a[0] == '*' {
		return "&" + a[1:]
	}
	return a
}

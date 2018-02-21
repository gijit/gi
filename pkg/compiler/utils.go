package compiler

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gijit/gi/pkg/ast"
	"github.com/gijit/gi/pkg/constant"
	"github.com/gijit/gi/pkg/token"
	"github.com/gijit/gi/pkg/types"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/gijit/gi/pkg/compiler/analysis"
	"github.com/gijit/gi/pkg/compiler/typesutil"
)

func (c *funcContext) Write(b []byte) (int, error) {
	c.writePos()
	c.output = append(c.output, b...)
	pp("func.ContextWrite, c.output is now '%v'", string(c.output))
	// DEBUG:
	if strings.HasPrefix(string(c.output), "	sum1 = adder(5, 5)") {
		//panic("where")
	}
	return len(b), nil
}

func (c *funcContext) Printf(format string, values ...interface{}) {
	c.Write([]byte(strings.Repeat("\t", c.p.indentation)))
	fmt.Fprintf(c, format, values...)
	c.Write([]byte{'\n'})
	c.Write(c.delayedOutput)
	c.delayedOutput = nil
}

func (c *funcContext) PrintCond(cond bool, onTrue, onFalse string) {
	if !cond {
		c.Printf("/* %s */ %s", strings.Replace(onTrue, "*/", "<star>/", -1), onFalse)
		return
	}
	c.Printf("%s", onTrue)
}

func (c *funcContext) SetPos(pos token.Pos) {
	c.posAvailable = true
	c.pos = pos
}

func (c *funcContext) writePos() {
	// jea debug: turn off writePos() for now
	return
	if c.posAvailable {
		c.posAvailable = false
		c.Write([]byte{'\b'})
		binary.Write(c, binary.BigEndian, uint32(c.pos))
	}
}

func (c *funcContext) Indent(f func()) {
	c.p.indentation++
	f()
	c.p.indentation--
}

func (c *funcContext) CatchOutput(indent int, f func()) []byte {
	origoutput := c.output
	c.output = nil
	c.p.indentation += indent
	f()
	c.writePos()
	catched := c.output
	c.output = origoutput
	c.p.indentation -= indent
	return catched
}

func (c *funcContext) Delayed(f func()) {
	c.delayedOutput = c.CatchOutput(0, f)
}

func (c *funcContext) translateArgs(sig *types.Signature, argExprs []ast.Expr, ellipsis bool) []string {
	pp("top of translateArgs, len(argExprs)=%v, ellipsis='%v'", len(argExprs), ellipsis)
	if len(argExprs) == 1 {
		if tuple, isTuple := c.p.TypeOf(argExprs[0]).(*types.Tuple); isTuple {
			pp("translateArgs: we have 1 argExpr that is a tuple; unpacking it into an expanded argExprs")
			c.Printf("\n// utils.go:88 translateArgs: we have 1 argExpr that is a tuple; unpacking it into an expanded argExprs\n")
			tupleVar := c.newVariable("_tuple")
			c.Printf("%s = %s;", tupleVar, c.translateExpr(argExprs[0], nil))
			argExprs = make([]ast.Expr, tuple.Len())
			for i := range argExprs {
				argExprs[i] = c.newIdent(c.formatExpr("%s[%d]", tupleVar, i).String(), tuple.At(i).Type())
			}
		}
	}

	paramsLen := sig.Params().Len()

	pp("sig.Variadic()=%v, ellipsis=%v, paramsLen=%v,  len(argExprs)=%v", sig.Variadic(), ellipsis, paramsLen, len(argExprs))

	var varargType *types.Slice
	numFixedArgs := paramsLen
	if sig.Variadic() && !ellipsis {
		// in the *signature*, not actuals.
		varargType = sig.Params().At(paramsLen - 1).Type().(*types.Slice)
		pp("varargType: '%#v'/elem type='%T'", varargType, varargType.Elem())
		numFixedArgs--
	}

	preserveOrder := false
	for i := 1; i < len(argExprs); i++ {
		preserveOrder = preserveOrder || c.Blocking[argExprs[i]]
	}
	pp("preserveOrder=%v", preserveOrder)
	pp("len(argExprs)=%v", len(argExprs))

	args := make([]string, len(argExprs))
	for i, argExpr := range argExprs {
		var argType types.Type
		switch {
		case varargType != nil && i >= paramsLen-1:
			argType = varargType.Elem()
		default:
			argType = sig.Params().At(i).Type()
		}

		arg := c.translateImplicitConversionWithCloning(argExpr, argType).String()

		if preserveOrder && c.p.Types[argExpr].Value == nil {
			argVar := c.newVariable("_arg")
			c.Printf("%s = %s;", argVar, arg)
			arg = argVar
		}

		args[i] = arg
	}

	pp("jea debug utils.go: argExprs = '%#v'", argExprs)
	for i := range argExprs {
		pp("jea debug utils.go: argExprs[i=%v] = '%#v'", i, argExprs[i])
	}

	pp("jea debug utils.go: len(args)=%v, args = '%#v', which is after c.translateImplicitConversionWithCloning", len(args), args)
	for i := range args {
		pp("jea debug  args[i=%v] = '%#v'", i, args[i])
	}

	pp("jea debug utils.go:151 varargType = '%#v'", varargType)

	// jea add, then comment out. Put this into luar.

	if ellipsis {
		pp("ellipsis true, paramsLen=%v, args='%#v'", paramsLen, args)
		return append(args[:paramsLen-1], fmt.Sprintf(`__lazy_ellipsis(%s)`, strings.Join(args[paramsLen-1:], ", ")))
	}

	// jea debug experiment... what if we turn off variadic for a moment?
	//if varargType == nil {
	if true {

		pp("jea debug utils.go:153 varargType is nil, returning args='%#v'", args)
		return args
	}
	// jea add:
	if len(args) == numFixedArgs {
		pp("jea debug utils.go:158 len(args)==numFixedArgs, returning args='%#v'", args)
		return args
	}

	// INVAR: varargType != nil

	// the 'awesome new' in this expression
	// fmt.Sprintf("hello %v", awesome new sliceType([new Int(3)]));
	//return append(args[:paramsLen-1], fmt.Sprintf("awesome new %s([%s])", c.typeName(varargType), strings.Join(args[paramsLen-1:], ", ")))

	// c.typeName(varargType) : "sliceType" -> "_gi_NewSlice"
	newOper := translateTypeNameToNewOper(c.typeName(0, varargType))

	pp("jea debug, utils.go: paramsLen = %v; newOper=%#v", paramsLen, newOper)
	for i := range args {
		pp("jea debug, utils.go: args[i=%v]='%v'", i, args[i])
	}

	// what is the

	// the ones >= paramsLen-1 are those from the variadic last type.
	return append(args[:paramsLen-1], fmt.Sprintf(`%s("interface{}",{%s})`, newOper, strings.Join(args[paramsLen-1:], ", ")))
}

func (c *funcContext) translateSelection(sel selection, pos token.Pos) ([]string, string) {
	var fields []string
	t := sel.Recv()
	for _, index := range sel.Index() {
		if ptr, isPtr := t.(*types.Pointer); isPtr {
			t = ptr.Elem()
		}
		s := t.Underlying().(*types.Struct)
		if jsTag := getJsTag(s.Tag(index)); jsTag != "" {
			jsFieldName := s.Field(index).Name()
			for {
				fields = append(fields, fieldName(s, 0))
				ft := s.Field(0).Type()
				if typesutil.IsJsObject(ft) {
					return fields, jsTag
				}
				ft = ft.Underlying()
				if ptr, ok := ft.(*types.Pointer); ok {
					ft = ptr.Elem().Underlying()
				}
				var ok bool
				s, ok = ft.(*types.Struct)
				if !ok || s.NumFields() == 0 {
					c.p.errList = append(c.p.errList, types.Error{Fset: c.p.fileSet, Pos: pos, Msg: fmt.Sprintf("could not find field with type *js.Object for 'js' tag of field '%s'", jsFieldName), Soft: true})
					return nil, ""
				}
			}
		}
		fields = append(fields, fieldName(s, index))
		t = s.Field(index).Type()
	}
	return fields, ""
}

var nilObj = types.Universe.Lookup("nil")

func (c *funcContext) zeroValue(ty types.Type) ast.Expr {
	switch t := ty.Underlying().(type) {
	case *types.Basic:
		switch {
		case isBoolean(t):
			return c.newConst(ty, constant.MakeBool(false))
		case isNumeric(t):
			return c.newConst(ty, constant.MakeInt64(0))
		case isString(t):
			return c.newConst(ty, constant.MakeString(""))
		case t.Kind() == types.UnsafePointer:
			// fall through to "nil"
		case t.Kind() == types.UntypedNil:
			panic("Zero value for untyped nil.")
		default:
			panic(fmt.Sprintf("Unhandled basic type: %v\n", t))
		}
	case *types.Array, *types.Struct:
		return c.setType(&ast.CompositeLit{}, ty)
	case *types.Chan, *types.Interface, *types.Map, *types.Signature, *types.Slice, *types.Pointer:
		// fall through to "nil"
	default:
		panic(fmt.Sprintf("Unhandled type: %T\n", t))
	}
	id := c.newIdent("nil", ty)
	c.p.Uses[id] = nilObj
	return id
}

func (c *funcContext) newConst(t types.Type, value constant.Value) ast.Expr {
	id := &ast.Ident{}
	c.p.Types[id] = types.TypeAndValue{Type: t, Value: value}
	return id
}

func (c *funcContext) newVariable(name string) string {
	return c.newVariableWithLevel(name, false)
}

func (c *funcContext) gensym(name string) string {
	next := atomic.AddInt64(&c.genSymCounter, 1)
	return fmt.Sprintf("__gensym_%v_%s", next, name)
}

func (c *funcContext) newVariableWithLevel(name string, pkgLevel bool) string {
	pp("newVariableWithLevel begins, with name='%s'", name)
	if name == "" {
		panic("newVariable: empty name")
	}
	name = encodeIdent(name)
	pp("newVariableWithLevel begins, after encodeIdent, name='%s'", name)

	if c.p.minify {
		i := 0
		for {
			offset := int('a')
			if pkgLevel {
				offset = int('A')
			}
			j := i
			name = ""
			for {
				name = string(offset+(j%26)) + name
				j = j/26 - 1
				if j == -1 {
					break
				}
			}
			if c.allVars[name] == 0 {
				break
			}
			i++
		}
	}
	n := c.allVars[name]
	c.allVars[name] = n + 1
	varName := name
	// jea... what purpose does this serve? we
	// want to re-use the same variable name when
	// we re-declare vars, not add another variable with _1 tacked
	// on at the end.
	// Answer, this was generating different tmp variables with different
	// names. Use a gensym instead.

	if false { // jea add
		if n > 0 {
			//fmt.Printf("repeated name '%s'! funcContext c = '%#v'", name, c)
			//fmt.Printf("repeated name '%s'! funcContext c.parent = '%#v'", name, *c.parent)
			varName = fmt.Sprintf("%s_%d", name, n)
		}
	}
	pp("in newVariableWithLevel(), varName = '%s'", varName)
	if pkgLevel {
		for c2 := c.parent; c2 != nil; c2 = c2.parent {
			c2.allVars[name] = n + 1
		}
		return varName
	}

	c.localVars = append(c.localVars, varName)
	return varName
}

func (c *funcContext) newIdent(name string, t types.Type) *ast.Ident {
	ident := ast.NewIdent(name)
	c.setType(ident, t)
	obj := types.NewVar(0, c.p.Pkg, name, t)
	c.p.Uses[ident] = obj
	c.p.objectNames[obj] = name
	return ident
}

func (c *funcContext) setType(e ast.Expr, t types.Type) ast.Expr {
	c.p.Types[e] = types.TypeAndValue{Type: t}
	return e
}

func (c *funcContext) pkgVar(pkg *types.Package) string {
	if pkg == c.p.Pkg {
		if !c.topLevelRepl {
			return "__pkg"
		}
		return ""
	}

	pkgVar, found := c.p.pkgVars[pkg.Path()]
	if !found {
		pp("not found!")
		pkgVar = fmt.Sprintf(`__packages["%s"]`, pkg.Path())
	}
	return pkgVar
}

func isVarOrConst(o types.Object) bool {
	switch o.(type) {
	case *types.Var, *types.Const:
		return true
	}
	return false
}

func isPkgLevel(o types.Object) bool {
	return o.Parent() != nil && o.Parent().Parent() == types.Universe
}

func (c *funcContext) objectName(o types.Object) (nam string) {
	defer func() {
		pp("objectName called for o='%#v', returning '%s'", o, nam)
	}()

	if isPkgLevel(o) {
		c.p.dependencies[o] = true

		if o.Pkg() != c.p.Pkg || (isVarOrConst(o) && o.Exported()) {
			// jea, foregoing pkg vars with the $Pkg. prefix, for now.
			// return o.Name() // jea was this, until we needed fmt.Sprintf
			// jea debug here for fmt.Sprintf
			pp("o.Pkg() = '%#v', o.Name()='%#v'", o.Pkg(), o.Name())
			pkgPrefix := c.pkgVar(o.Pkg())
			if pkgPrefix == "" {
				return o.Name()
			}
			return pkgPrefix + "." + o.Name()

		}
	}

	name, ok := c.p.objectNames[o]
	pp("utils.go:307, name='%v', ok='%v'", name, ok)
	if !ok {
		name = c.newVariableWithLevel(o.Name(), isPkgLevel(o))
		pp("name='%#v', o.Name()='%v'", name, o.Name())
		c.p.objectNames[o] = name
	}

	if v, ok := o.(*types.Var); ok && c.p.escapingVars[v] {
		return name + "[0]"
	}
	return name
}

func (c *funcContext) varPtrName(o *types.Var) string {
	if isPkgLevel(o) && o.Exported() {
		return c.pkgVar(o.Pkg()) + "." + o.Name() + "_ptr"
	}

	name, ok := c.p.varPtrNames[o]
	if !ok {
		name = c.newVariableWithLevel(o.Name()+"_ptr", isPkgLevel(o))
		c.p.varPtrNames[o] = name
	}
	return name
}

type typeNameSetting int

const (
	IMMEDIATE typeNameSetting = 0 // default
	DELAYED   typeNameSetting = 1
	SKIP_ANON typeNameSetting = 2
)

func (tns typeNameSetting) String() string {
	switch tns {
	case IMMEDIATE:
		return "IMMEDIATE"
	case DELAYED:
		return "DELAYED"
	case SKIP_ANON:
		return "SKIP_ANON"
	}
	panic("unknown tns")
}

// re-port typeName, did not back mat_test 501 package level vars/funcs
func (c *funcContext) typeNameReport(level int, ty types.Type) string {

	switch t := ty.(type) {
	case *types.Basic:
		return "__type__." + toJavaScriptType(t)
	case *types.Named:
		if t.Obj().Name() == "error" {
			return "__type__.error"
		}
		return "__type__." + c.objectName(t.Obj())
	case *types.Interface:
		if t.Empty() {
			return "__type__.emptyInterface"
		}
	}

	anonType, ok := c.p.anonTypeMap.At(ty).(*types.TypeName)
	if !ok {
		c.initArgs(ty) // cause all embedded types to be registered
		varName := c.newVariableWithLevel(strings.ToLower(typeKind(ty)[6:])+"Type", true)
		anonType = types.NewTypeName(token.NoPos, c.p.Pkg, varName, ty) // fake types.TypeName
		c.p.anonTypes = append(c.p.anonTypes, anonType)
		c.p.anonTypeMap.Set(ty, anonType)
	}
	c.p.dependencies[anonType] = true
	return anonType.Name()
}

func (c *funcContext) typeName(level int, ty types.Type) (res string) {

	res, _, _, _ = c.typeNameWithAnonInfo(ty)
	return
}

func (c *funcContext) typeNameWithAnonInfo(
	ty types.Type,
) (
	res string,
	isAnon bool,
	anonType *types.TypeName,
	createdVarName string,
) {

	whenAnonPrint := c.TypeNameSetting

	pp("in typeNameWithAnonInfo, ty='%#v'; whenAnonPrint=%v", ty, whenAnonPrint)
	defer func() {
		// funcContext.typeName returning with res = 'sliceType'
		//fmt.Printf("funcContext.typeName returning with res = '%s'\n", res)
		if res == "__type__.anon_sliceType" {
			//panic("where?")
		}
	}()
	switch t := ty.(type) {
	case *types.Basic:
		jst := toJavaScriptType(t)
		pp("in typeName, basic, calling toJavaScriptType t='%#v', got jst='%s'", t, jst)
		res = "__type__." + jst
		return
	case *types.Named:
		if t.Obj().Name() == "error" {
			res = "__type__.error"
			return
		}
		res = "__type__." + c.objectName(t.Obj())
		return
	case *types.Interface:
		if t.Empty() {
			res = "__type__.emptyInterface"
			return
		}
	}

	anonType, isAnon = c.p.anonTypeMap.At(ty).(*types.TypeName)
	if !isAnon {
		isAnon = true

		// cause all embedded types to be registered
		c.initArgs(ty)

		// [6:] takes prefix "__kind" off.
		low := "__type__.anon_" + strings.ToLower(typeKind(ty)[6:]) + "Type"

		// typeKind(ty)='_kindSlice', low='sliceType'
		pp("typeKind(ty)='%s', low='%s'", typeKind(ty), low)

		varName := c.newVariableWithLevel(low, true)

		anonType = types.NewTypeName(token.NoPos, c.p.Pkg, varName, ty) // fake types.TypeName
		c.p.anonTypes = append(c.p.anonTypes, anonType)
		pp("just added anonType='%s', where whenAnonPrint='%v'", anonType.Name(), whenAnonPrint)
		c.p.anonTypeMap.Set(ty, anonType)
		createdVarName = varName

		anonTypePrint := fmt.Sprintf("\n\t%s = __%sType(%s); -- '%s' anon type printing. utils.go:506\n", varName, strings.ToLower(typeKind(anonType.Type())[6:]), c.initArgs(anonType.Type()), whenAnonPrint.String())
		// gotta generate the type immediately for the REPL.
		// But the pointer  needs to come after the struct it references.

		//fmt.Printf("whenAnonPrint = %v, for anonType: '%s'.\n", whenAnonPrint, anonType.Name())
		switch whenAnonPrint {
		case DELAYED:
			//panic("where?")
			c.Delayed(func() {
				c.Printf(anonTypePrint)
			})
		case IMMEDIATE:
			c.Printf(anonTypePrint)
			pp("done with IMMEDIATE printing of anonTypePrint='%v'", anonTypePrint)
		case SKIP_ANON:
		}
	}
	c.p.dependencies[anonType] = true
	res = anonType.Name()
	return
}

func (c *funcContext) externalize(s string, t types.Type) string {
	if typesutil.IsJsObject(t) {
		return s
	}
	switch u := t.Underlying().(type) {
	case *types.Basic:
		if isNumeric(u) && !is64Bit(u) && !isComplex(u) {
			return s
		}
		if u.Kind() == types.UntypedNil {
			return "null"
		}
	}
	return fmt.Sprintf("__externalize(%s, %s)", s, c.typeName(0, t))
}

func (c *funcContext) handleEscapingVars(n ast.Node) {
	newEscapingVars := make(map[*types.Var]bool)
	for escaping := range c.p.escapingVars {
		newEscapingVars[escaping] = true
	}
	c.p.escapingVars = newEscapingVars

	var names []string
	objs := analysis.EscapingObjects(n, c.p.Info.Info)
	sort.Slice(objs, func(i, j int) bool {
		if objs[i].Name() == objs[j].Name() {
			return objs[i].Pos() < objs[j].Pos()
		}
		return objs[i].Name() < objs[j].Name()
	})
	for _, obj := range objs {
		names = append(names, c.objectName(obj))
		c.p.escapingVars[obj] = true
	}
	sort.Strings(names)
	for _, name := range names {
		c.Printf("%s = [%s];", name, name)
	}
}

func fieldName(t *types.Struct, i int) string {
	name := t.Field(i).Name()
	if name == "_" || reservedKeywords[name] {
		return fmt.Sprintf("%s__%d", name, i)
	}
	return name
}

func typeKind(ty types.Type) (res string) {
	defer func() {
		pp("typeKind called on ty='%#v', returning res='%s'", ty, res)
	}()
	switch t := ty.Underlying().(type) {
	case *types.Basic:
		return "__kind" + toJavaScriptTypeUppercase(t) //toJavaScriptType(t)
	case *types.Array:
		return "__kindArray"
	case *types.Chan:
		return "__kindChan"
	case *types.Interface:
		return "__kindInterface"
	case *types.Map:
		return "__kindMap"
	case *types.Signature:
		return "__kindFunc"
	case *types.Slice:
		return "__kindSlice"
	case *types.Struct:
		return "__kindStruct"
	case *types.Pointer:
		return "__kindPtr"
	default:
		panic(fmt.Sprintf("Unhandled type: %T\n", t))
	}
}

func toJavaScriptType(t *types.Basic) string {
	switch t.Kind() {
	case types.UntypedInt:
		return "int"
	case types.Byte:
		return "uint8"
	case types.Rune:
		return "int32"
	case types.UnsafePointer:
		return "UnsafePointer"
	default:
		name := t.String()
		//jea
		//return strings.ToUpper(name[:1]) + name[1:]
		return name
	}
}

func toJavaScriptTypeUppercase(t *types.Basic) string {
	switch t.Kind() {
	case types.UntypedInt:
		return "Int"
	case types.Byte:
		return "Uint8"
	case types.Rune:
		return "Int32"
	case types.UnsafePointer:
		return "UnsafePointer"
	default:
		name := t.String()
		return strings.ToUpper(name[:1]) + name[1:]
	}
}

func is64Bit(t *types.Basic) bool {
	return t.Kind() == types.Int64 || t.Kind() == types.Uint64
}

func isBoolean(t *types.Basic) bool {
	return t.Info()&types.IsBoolean != 0
}

func isComplex(t *types.Basic) bool {
	return t.Info()&types.IsComplex != 0
}

func isFloat(t *types.Basic) bool {
	return t.Info()&types.IsFloat != 0
}

func isInteger(t *types.Basic) bool {
	return t.Info()&types.IsInteger != 0
}

func isNumeric(t *types.Basic) bool {
	return t.Info()&types.IsNumeric != 0
}

func isString(t *types.Basic) bool {
	return t.Info()&types.IsString != 0
}

func isUnsigned(t *types.Basic) bool {
	return t.Info()&types.IsUnsigned != 0
}

func isBlank(expr ast.Expr) bool {
	if expr == nil {
		return true
	}
	if id, isIdent := expr.(*ast.Ident); isIdent {
		return id.Name == "_"
	}
	return false
}

func nameHelper(expr ast.Expr) string {
	if expr == nil {
		return "_"
	}
	if id, isIdent := expr.(*ast.Ident); isIdent {
		return id.Name
	}
	return "_"
}

func isWrapped(ty types.Type) bool {
	switch t := ty.Underlying().(type) {
	case *types.Basic:
		// jea simplify, we aren't in 32-bit javascript land.
		return false
		/*
						if isString(t) {
							return false
						}
			            // jea add
						kind := t.Kind()
						if kind == types.Int || kind == types.Uint {
							return false
						}

			return !is64Bit(t) && !isComplex(t) && t.Kind() != types.UntypedNil
		*/
	case *types.Array, *types.Chan, *types.Map, *types.Signature:
		return true
	case *types.Pointer:
		_, isArray := t.Elem().Underlying().(*types.Array)
		return isArray
	}
	return false
}

func encodeString(s string) string {
	pp("jea debug: encodeString called with s='%s'", s)
	buffer := bytes.NewBuffer(nil)
	for _, r := range []byte(s) {
		switch r {
		case '\b':
			buffer.WriteString(`\b`)
		case '\f':
			buffer.WriteString(`\f`)
		case '\n':
			buffer.WriteString(`\n`)
		case '\r':
			buffer.WriteString(`\r`)
		case '\t':
			buffer.WriteString(`\t`)
		case '\v':
			buffer.WriteString(`\v`)
		case '"':
			buffer.WriteString(`\"`)
		case '\\':
			buffer.WriteString(`\\`)
		default:
			if r < 0x20 || r > 0x7E {
				fmt.Fprintf(buffer, `\x%02X`, r)
				continue
			}
			buffer.WriteByte(r)
		}
	}
	return `"` + buffer.String() + `"`
}

func getJsTag(tag string) string {
	for tag != "" {
		// skip leading space
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// scan to colon.
		// a space or a quote is a syntax error
		i = 0
		for i < len(tag) && tag[i] != ' ' && tag[i] != ':' && tag[i] != '"' {
			i++
		}
		if i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// scan quoted string to find value
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		if name == "js" {
			value, _ := strconv.Unquote(qvalue)
			return value
		}
	}
	return ""
}

func needsSpace(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '$'
}

func removeWhitespace(b []byte, minify bool) []byte {
	if !minify {
		return b
	}

	var out []byte
	var previous byte
	for len(b) > 0 {
		switch b[0] {
		case '\b':
			out = append(out, b[:5]...)
			b = b[5:]
			continue
		case ' ', '\t', '\n':
			if (!needsSpace(previous) || !needsSpace(b[1])) && !(previous == '-' && b[1] == '-') {
				b = b[1:]
				continue
			}
		case '"':
			out = append(out, '"')
			b = b[1:]
			for {
				i := bytes.IndexAny(b, "\"\\")
				out = append(out, b[:i]...)
				b = b[i:]
				if b[0] == '"' {
					break
				}
				// backslash
				out = append(out, b[:2]...)
				b = b[2:]
			}
		case '/':
			if b[1] == '*' {
				i := bytes.Index(b[2:], []byte("*/"))
				b = b[i+4:]
				continue
			}
		}
		out = append(out, b[0])
		previous = b[0]
		b = b[1:]
	}
	return out
}

func setRangeCheck(constantIndex, array bool) string {
	if constantIndex && array {
		return "%1e[%2f] = %3s"
	}

	return "__gi_SetRangeCheck(%1e, %2f, %3s)"
}

func rangeCheck(pattern string, constantIndex, array bool) string {

	if constantIndex && array {
		return pattern
	}
	//panic("where?")
	return "__gi_GetRangeCheck(%1e, %2f)"
	/*
		lengthProp := "$length"
		if array {
			lengthProp = "length"
		}
		check := "%2f >= %1e." + lengthProp
		if !constantIndex {
			check = "(%2f < 0 || " + check + ")"
		}
		return "(" + check + ` ? ($throwRuntimeError("index out of range"), undefined) : ` + pattern + ")"
	*/
}

func endsWithReturn(stmts []ast.Stmt) bool {
	if len(stmts) > 0 {
		if _, ok := stmts[len(stmts)-1].(*ast.ReturnStmt); ok {
			return true
		}
	}
	return false
}

func encodeIdent(name string) string {
	return strings.Replace(url.QueryEscape(name), "%", "$", -1)
}

func stripOuterParen(s string) (r string) {
	r = strings.TrimSpace(s)
	n := len(r)
	if n < 2 {
		return r
	}
	if r[0] != '(' || r[len(r)-1] != ')' {
		return r
	}
	if n == 2 {
		return ""
	}
	return r[1 : n-1]
}

// function(a) blah -> blah
func stripFirstFunctionAndArg(s string) (head, body string) {
	st := strings.TrimSpace(s)
	r := st
	n := len(r)
	if n < 10 {
		return s, ""
	}
	if !strings.HasPrefix(r, "function") {
		return s, ""
	}
	r = strings.TrimSpace(r[8:])
	if len(r) < 2 {
		return s, ""
	}
	if r[0] != '(' {
		return s, ""
	}
	r = r[1:]
	n = len(r)
	// trim up to and include the first ')'
	pos := -1
	depth := 1
posloop:
	for i := 0; i < n; i++ {
		switch r[i] {
		case ')':
			depth--
			if depth == 0 {
				pos = i
				break posloop
			}
		case '(':
			depth++
		}
	}
	if pos == -1 {
		return s, ""
	}
	body = r[pos+1:]
	lenbod := len(body)
	lenhead := len(st) - lenbod
	head = st[:lenhead]
	return
}

func translateTypeNameToNewOper(typeName string) string {
	switch typeName {
	case "sliceType":
		return "_gi_NewSlice HUH?? --where is this used? utils.go:942"
	}
	panic(fmt.Sprintf("what here? for typeName = '%s'", typeName))
}

package compiler

import (
	"bytes"
	"os"
	//"encoding/json"
	"fmt"
	"github.com/glycerine/gi/pkg/ast"
	//"github.com/glycerine/gi/pkg/constant"
	"github.com/glycerine/gi/pkg/printer"
	"github.com/glycerine/gi/pkg/token"
	"github.com/glycerine/gi/pkg/types"
	"sort"
	"strings"

	"github.com/glycerine/gi/pkg/compiler/analysis"
	//"github.com/neelance/astrewrite"
	//"golang.org/x/tools/go/gcimporter15"
	"golang.org/x/tools/go/types/typeutil"
)

type pkgContext struct {
	*analysis.Info
	additionalSelections map[*ast.SelectorExpr]selection

	typeNames []*types.TypeName

	// jea add
	typeDepend        *dfsState
	typeDefineLuaCode map[types.Object]string
	importedPackages  map[string]*types.Package

	pkgVars      map[string]string
	objectNames  map[types.Object]string
	varPtrNames  map[*types.Var]string
	anonTypes    []*types.TypeName
	anonTypeMap  typeutil.Map
	escapingVars map[*types.Var]bool
	indentation  int
	dependencies map[types.Object]bool
	minify       bool
	fileSet      *token.FileSet
	files        []*ast.File
	errList      ErrorList
}

func (p *pkgContext) SelectionOf(e *ast.SelectorExpr) (selection, bool) {
	if sel, ok := p.Selections[e]; ok {
		return sel, true
	}
	if sel, ok := p.additionalSelections[e]; ok {
		return sel, true
	}
	return nil, false
}

type selection interface {
	Kind() types.SelectionKind
	Recv() types.Type
	Index() []int
	Obj() types.Object
	Type() types.Type
}

type fakeSelection struct {
	kind  types.SelectionKind
	recv  types.Type
	index []int
	obj   types.Object
	typ   types.Type
}

func (sel *fakeSelection) Kind() types.SelectionKind { return sel.kind }
func (sel *fakeSelection) Recv() types.Type          { return sel.recv }
func (sel *fakeSelection) Index() []int              { return sel.index }
func (sel *fakeSelection) Obj() types.Object         { return sel.obj }
func (sel *fakeSelection) Type() types.Type          { return sel.typ }

type funcContext struct {
	*analysis.FuncInfo
	p             *pkgContext
	parent        *funcContext
	sig           *types.Signature
	allVars       map[string]int
	localVars     []string
	resultNames   []ast.Expr
	flowDatas     map[*types.Label]*flowData
	caseCounter   int
	labelCases    map[*types.Label]int
	output        []byte
	delayedOutput []byte
	posAvailable  bool
	pos           token.Pos

	genSymCounter int64

	intType types.Type

	TypeNameSetting typeNameSetting

	topLevelRepl bool

	PkgNameOverride      bool
	PkgNameOverrideValue string
}

type flowData struct {
	postStmt  func()
	beginCase int
	endCase   int
}

type ImportContext struct {
	Packages map[string]*types.Package
	Import   func(path, pkgDir string, depth int) (*Archive, error)
}

// packageImporter implements go/types.Importer interface.
type packageImporter struct {
	importContext *ImportContext
	importError   *error // A pointer to importError in Compile.
}

func (pi packageImporter) Import(path string, depth int) (*types.Package, error) {
	if path == "unsafe" {
		return types.Unsafe, nil
	}

	pp("pi = '%#v', pi.importContext='%#v'", pi, pi.importContext)
	pp("path='%s', pi.importContext.Import='%#v'", path, pi.importContext.Import)
	a, err := pi.importContext.Import(path, "", depth+1)
	pp("jea debug: a *Archive back from pi.importContext.Import('%s') (err='%v') archive is '%#v'", path, err, a)
	if err != nil {
		if *pi.importError == nil {
			// If import failed, show first error of import only (https://github.com/gopherjs/gopherjs/issues/119).
			*pi.importError = err
		}
		return nil, err
	}

	tyPack := a.Pkg
	pi.importContext.Packages[a.ImportPath] = tyPack

	// jea: import "fmt" gives not nil tyPack.
	pp("end of compiler.packageImporter.Import(path='%s'), tyPack is '%#v'.", path, tyPack)
	return tyPack, nil
}

func isPrim(ty types.Type) bool {
	et := elemType(ty)
	if et == nil {
		return false
	}
	// want emptyInterface to be primitive too!
	iface, isIface := et.(*types.Interface)
	if isIface {
		if iface.Empty() {
			return true
		}
	}

	_, ok := et.(*types.Basic)
	pp("for et=%#v, ty='%#v'/%T, isPrim=%v", et, ty, ty, ok)
	return ok
}

// mirror what initArgs will switch on.
func elemType(ty types.Type) types.Type {

	switch t := ty.(type) {
	case *types.Basic:
		return ty
	case *types.Array:
		return t.Elem()
	case *types.Chan:
		return t.Elem()
	case *types.Slice:
		return t.Elem()
	}
	return nil
}

func (c *funcContext) initArgsNoPkgForPrimitives(ty types.Type) string {
	if isPrim(ty) {
		c.PkgNameOverride = true
		c.PkgNameOverrideValue = ""
		s := c.initArgs(ty)
		c.PkgNameOverride = false
		return s
	}
	return c.initArgs(ty)
}
func (c *funcContext) initArgs(ty types.Type) string {

	prev := c.TypeNameSetting
	c.TypeNameSetting = IMMEDIATE
	defer func() {
		c.TypeNameSetting = prev
	}()

	//fmt.Printf("\n initArgs: ty = '%#v'\n", ty)
	// &types.Tuple{vars:[]*types.Var{(*types.Var)(0xc4201800f0), (*types.Var)(0xc420180140)}}'
	switch t := ty.(type) {
	case *types.Array:
		return fmt.Sprintf("%s, %d", c.typeName(t.Elem(), t), t.Len())
	case *types.Chan:
		return fmt.Sprintf("%s, %t, %t", c.typeName(t.Elem(), t), t.Dir()&types.SendOnly != 0, t.Dir()&types.RecvOnly != 0)
	case *types.Interface:
		methods := make([]string, t.NumMethods())
		for i := range methods {
			method := t.Method(i)
			pkgPath := ""
			if !method.Exported() {
				pkgPath = method.Pkg().Path()
			}
			methods[i] = fmt.Sprintf(`{__prop= "%s", __name= %s, __pkg= "%s", __typ= __funcType(%s)}`, method.Name(), encodeString(method.Name()), pkgPath, c.initArgs(method.Type()))
		}
		return fmt.Sprintf("{%s}", strings.Join(methods, ", "))
	case *types.Map:
		return fmt.Sprintf("%s, %s", c.typeName(t.Key(), t), c.typeName(t.Elem(), t))
	case *types.Pointer:
		//vv("t.Elem()='%#v',  t='%#v'", t.Elem().String(), t.String()) // t.Elem="main.S", t="*main.S"
		return fmt.Sprintf("%s", c.typeName(t.Elem(), nil))
	case *types.Slice:
		return fmt.Sprintf("%s", c.typeName(t.Elem(), nil))
	case *types.Signature:
		params := make([]string, t.Params().Len())
		for i := range params {
			params[i] = c.typeName(t.Params().At(i).Type(), t)
		}
		results := make([]string, t.Results().Len())
		for i := range results {
			results[i] = c.typeName(t.Results().At(i).Type(), t)
		}
		return fmt.Sprintf("{%s}, {%s}, %t", strings.Join(params, ", "), strings.Join(results, ", "), t.Variadic())
	case *types.Struct:
		pkgPath := ""
		fields := make([]string, t.NumFields())
		for i := range fields {
			field := t.Field(i)
			if !field.Exported() {
				pkgPath = field.Pkg().Path()
			}
			fields[i] = fmt.Sprintf(`{__prop= "%s", __name= %s, __anonymous= %t, __exported= %t, __typ= %s, __tag= %s}`, fieldName(t, i), encodeString(field.Name()), field.Anonymous(), field.Exported(), c.typeName(field.Type(), t), encodeString(t.Tag(i)))
		}
		return fmt.Sprintf(`"%s", {%s}`, pkgPath, strings.Join(fields, ", "))
	case *types.Tuple:
		// A Tuple represents an ordered list of variables;
		// a nil *Tuple is a valid (empty) tuple.
		// Tuples are used as components of signatures and to
		// represent the type of multiple
		// assignments; they are not first class types of Go.
		// vars []*Var

		results := make([]string, t.Len())
		for i := range results {
			results[i] = c.typeName(t.At(i).Type(), t)
		}

		return fmt.Sprintf("{%s}", strings.Join(results, ", "))
	default:
		panic("invalid type")
	}
}

func (c *funcContext) translateToplevelFunction(fun *ast.FuncDecl, info *analysis.FuncInfo) []byte {
	pp("translateToplevelFunction called! fun.Name.Name='%s'", fun.Name.Name)

	pkgName := c.getPkgName()

	o := c.p.Defs[fun.Name].(*types.Func)
	sig := o.Type().(*types.Signature)
	var recv *ast.Ident
	if fun.Recv != nil && fun.Recv.List[0].Names != nil {
		recv = fun.Recv.List[0].Names[0]
	}

	var joinedParams string

	primaryFunction := func(isMethod bool, funcRef string) []byte {
		if fun.Body == nil {
			return []byte(fmt.Sprintf("\t%s = function() \n\t\t__throwRuntimeError(\"native function not implemented: %s\");\n\t end ;\n", funcRef, o.FullName()))
		}

		params, fun, _ := translateFunction(fun.Type, recv, fun.Body, c, sig, info, funcRef, isMethod)
		pp("funcRef in translateFunction, package.go:698 is '%s'; isMethod='%v'; fun='%#v'; recv='%#v'; fun='%#v'; params='%#v';", funcRef, isMethod, fun, recv, fun, params)
		joinedParams = strings.Join(params, ", ")
		return []byte(fmt.Sprintf("\t%s = %s;\n", funcRef, fun))
	}

	code := bytes.NewBuffer(nil)

	if fun.Recv == nil {
		funcRef := c.objectName(o)
		code.Write(primaryFunction(false, funcRef))
		if fun.Name.IsExported() {
			fmt.Fprintf(code, "\t__pkg.%s = %s;\n", encodeIdent(fun.Name.Name), funcRef)
			//fmt.Fprintf(code, "\t%s = %s;\n", encodeIdent(fun.Name.Name), funcRef)
		}
		return code.Bytes()
	}

	recvType := sig.Recv().Type()
	ptr, isPointer := recvType.(*types.Pointer)
	namedRecvType, _ := recvType.(*types.Named)
	if isPointer {
		namedRecvType = ptr.Elem().(*types.Named)
	}
	typeName := "__type__." + pkgName + c.objectName(namedRecvType.Obj())
	funName := fun.Name.Name
	if reservedKeywords[funName] {
		funName += "_"
	}

	signatureDetail := ""
	if _, isStruct := namedRecvType.Underlying().(*types.Struct); isStruct {
		ptrAddMe := typeName + ".ptr.prototype." + funName
		code.Write(primaryFunction(true, ptrAddMe))
		// get the comma right: either function(this) or function(this, a,b,c)
		jp := ", " + joinedParams
		if joinedParams == "" {
			jp = ""
		}
		fmt.Fprintf(code, "\t%[1]s.prototype.%[2]s = function(this %[3]s)  return %[1]s.ptr.prototype.%[2]s(this.__val %[3]s); end;\n", typeName, funName, jp)

		signatureDetail = c.getMethodDetailsSig(o)
		// add to struct
		fmt.Fprintf(code, "\n %s.__addToMethods(%s); -- package.go:344\n", typeName, signatureDetail)
		// add to pointer to struct
		fmt.Fprintf(code, "\n %s.__addToMethods(%s); -- package.go:346\n", typeName+".ptr", signatureDetail)

		return code.Bytes()
	}

	if isPointer {
		if _, isArray := ptr.Elem().Underlying().(*types.Array); isArray {
			code.Write(primaryFunction(false, typeName+".prototype."+funName))
			fmt.Fprintf(code, "\t__ptrType(%s).prototype.%s = function(this, %s)  return (%s(this.__get())).%s(%s); end;\n", typeName, funName, joinedParams, typeName, funName, joinedParams)
			return code.Bytes()
		}
		return primaryFunction(false, fmt.Sprintf("__ptrType(%s).prototype.%s", typeName, funName))
	}

	value := "this.__get()"
	if isWrapped(recvType) {
		value = fmt.Sprintf("%s(%s)", typeName, value)
	}
	code.Write(primaryFunction(false, typeName+".prototype."+funName))
	fmt.Fprintf(code, "\t__ptrType(%s).prototype.%s = function(%s) return %s.%s(%s); end;\n", typeName, funName, joinedParams, value, funName, joinedParams)
	return code.Bytes()
}

func translateFunction(typ *ast.FuncType, recv *ast.Ident, body *ast.BlockStmt, outerContext *funcContext, sig *types.Signature, info *analysis.FuncInfo, funcRef string, isMethod bool) (params []string, fun string, recvName string) {
	if info == nil {
		panic("nil info")
	}

	if false {
		// debug only
		pp("translateFunction called, body = ")
		printer.Fprint(os.Stdout, outerContext.p.fileSet, body)
	}

	c := &funcContext{
		FuncInfo:    info,
		p:           outerContext.p,
		parent:      outerContext,
		sig:         sig,
		allVars:     make(map[string]int, len(outerContext.allVars)),
		localVars:   []string{},
		flowDatas:   map[*types.Label]*flowData{nil: {}},
		caseCounter: 1,
		labelCases:  make(map[*types.Label]int),
	}
	for k, v := range outerContext.allVars {
		c.allVars[k] = v
	}
	prevEV := c.p.escapingVars
	preComputedNamedNames := []string{}
	preComputedZeroRet := []string{}

	for _, param := range typ.Params.List {
		if len(param.Names) == 0 {
			params = append(params, c.newVariable("param"))
			continue
		}
		for _, ident := range param.Names {
			if isBlank(ident) {
				params = append(params, c.newVariable("param"))
				continue
			}
			params = append(params, c.objectName(c.p.Defs[ident]))
		}
	}

	bodyOutput := string(c.CatchOutput(1, func() {
		if len(c.Blocking) != 0 {
			c.p.Scopes[body] = c.p.Scopes[typ]
			c.handleEscapingVars(body)
		}

		if c.sig != nil && c.sig.Results().Len() != 0 && c.sig.Results().At(0).Name() != "" {
			c.resultNames = make([]ast.Expr, c.sig.Results().Len())
			for i := 0; i < c.sig.Results().Len(); i++ {
				result := c.sig.Results().At(i)
				objName := c.objectName(result)
				zeroV := c.translateExpr(c.zeroValue(result.Type()), nil).String()

				// NB doesn't work with "local %s = %s'", but it's
				// okay: this doesn't collide or
				// pollute the global env because code in defer.lua's
				// __actuallyCall gives the function its own environment.
				if c.HasDefer {
					// will write to the intermediate env
					// we establish with setfenv()
					c.Printf("%s = %s;", objName, zeroV)
				} else {
					c.Printf("local %s = %s;", objName, zeroV)
				}
				preComputedNamedNames = append(preComputedNamedNames, `"`+objName+`"`)
				preComputedZeroRet = append(preComputedZeroRet, zeroV)
				id := ast.NewIdent("")
				c.p.Uses[id] = result
				c.resultNames[i] = c.setType(id, result.Type())
			}
		}

		if recv != nil && !isBlank(recv) {
			recvName = c.translateExpr(recv, nil).String()

			// jea: omit "r = self" now in favor of
			// specifying 'r' as a part of the arguments explicitly
			// from the start, so we get the receiver name 'r' correct,
			// and don't face collisions with vars that happened to
			// be named 'self'.

			if isWrapped(c.p.TypeOf(recv)) {
				c.Printf("%[1]s = %[1]s.__val; -- isWrapped(recv) true at package.go:361\n", recvName)
			}
		}

		c.translateStmtList(body.List)
		if len(c.Flattened) != 0 && !endsWithReturn(body.List) {
			c.translateStmt(&ast.ReturnStmt{}, nil)
		}
	}))

	pp("bodyOutput = '%s'", bodyOutput)

	sort.Strings(c.localVars)

	var prefix, suffix, functionName string
	namedNames := ""
	zeroret := ""

	// jea temp disable with false
	if false {
		if len(c.Flattened) != 0 {
			c.localVars = append(c.localVars, "__s")
			prefix = prefix + " __s = 0;"
		}

		if c.HasDefer {
			c.localVars = append(c.localVars, "__deferred")
			suffix = " }" + suffix
			if len(c.Blocking) != 0 {
				suffix = " }" + suffix
			}
		}

		if len(c.Blocking) != 0 {
			c.localVars = append(c.localVars, "_r")
			if funcRef == "" {
				funcRef = "_b"
				functionName = " _b"
			}
			var stores, loads string
			for _, v := range c.localVars {
				loads += fmt.Sprintf("%s = __f.%s; ", v, v)
				stores += fmt.Sprintf("__f.%s = %s; ", v, v)
			}
			prefix = prefix + " var __f, __c = false; if (this !== undefined && this.__blk !== undefined) { __f = this; __c = true; " + loads + "}"
			suffix = " if (__f === undefined) { __f = { __blk: " + funcRef + " }; } " + stores + "return __f;" + suffix
		}

		if c.HasDefer {
			prefix = prefix + " var __err = null; try {"
			deferSuffix := " } catch(err) { __err = err;"
			if len(c.Blocking) != 0 {
				//deferSuffix += " __s = -1;"
			}
			if c.resultNames == nil && c.sig.Results().Len() > 0 {
				deferSuffix += fmt.Sprintf(" return%s;", c.translateResults(nil))
			}
			deferSuffix += " } finally { __callDeferred(__deferred, __err);"
			if c.resultNames != nil {
				namedNames = strings.Join(preComputedNamedNames, ", ") //fmt.Sprintf("{%s}", c.translateResultsAllQuoted(c.resultNames))
				zeroret = strings.Join(preComputedZeroRet, ", ")
				deferSuffix += fmt.Sprintf(" if (!__curGoroutine.asleep) { return %s; }", c.translateResults(c.resultNames))
			}
			if len(c.Blocking) != 0 {
				deferSuffix += " if(__curGoroutine.asleep) {"
			}
			suffix = deferSuffix + suffix
		}

		if len(c.Flattened) != 0 {
			prefix = prefix + " ::s:: while (true) do\n --switch (__s)\n\t\t if __s == 0 then\n"
			suffix = " end; return; end\n" + suffix
		}

	} else { // end if false, jea temp disable

		// jea: compute namedNames, we need!
		if c.HasDefer {
			if c.resultNames != nil {
				namedNames = strings.Join(preComputedNamedNames, ", ") //fmt.Sprintf("{%s}", c.translateResultsAllQuoted(c.resultNames))
				zeroret = strings.Join(preComputedZeroRet, ", ")
			}
		}
	}

	formals := strings.Join(params, ", ")
	functionWord := "function"
	if isMethod {
		// commenting out makes 029 go green... hmmm.
		// Otherwise we have a missing 'function' word vvv here  in the declaration
		//        __type__.Beagle.ptr.prototype.Write =      (b,with)
		//functionWord = ""
	}

	c.p.escapingVars = prevEV

	if c.HasDefer {
		pp("jea TODO: prefix is '%s'... should we not discard?", prefix)
		//		prefix = prefix + ...
		return params, fmt.Sprintf(`
%s%s(...) 
   local __orig = {...}
   local __defers={}
   local __zeroret = {%s}
   local __namedNames = {%s}
   local __actual=function(%s)
      %s
   end
   return __actuallyCall("%s", __actual, __namedNames, __zeroret, __defers, __orig)
end
`,
			functionWord, functionName, zeroret, namedNames, formals,
			bodyOutput, functionName), recvName

		//prefix = prefix + " __deferred = []; __deferred.index = __curGoroutine.deferStack.length; __curGoroutine.deferStack.push(__deferred);"
	}

	if prefix != "" {
		bodyOutput = strings.Repeat("\t", c.p.indentation+1) + "\n--jea package.go:553 \n" + prefix + "\n" + bodyOutput
	}
	if suffix != "" {
		bodyOutput = bodyOutput + strings.Repeat("\t", c.p.indentation+1) + "\n--jea package.go:556\n" + suffix + "\n"
	}
	if len(c.localVars) != 0 {
		// jea, these are javascript only: at the top of a
		// function body, have vars a,b,c for the formal parameters f(a,b,c).
		//bodyOutput = fmt.Sprintf("%svar %s;\n", strings.Repeat("\t", c.p.indentation+1), strings.Join(c.localVars, ", ")) + bodyOutput
	}

	recvInsert := ""
	if recvName != "" {
		recvInsert = recvName
		if formals != "" {
			recvInsert = recvInsert + ","
		}
	}
	return params, fmt.Sprintf("%s%s(%s%s) \n%s%s end",
			functionWord, functionName, recvInsert, formals,
			bodyOutput, strings.Repeat("\t", c.p.indentation)),
		recvName
}

// Copyright (c) 2010-2016 Steve Donovan
// Licensed under the MIT license found in the LICENSE file.

package luar

import (
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"
	"unsafe"

	"github.com/gijit/gi/pkg/verb"
	"github.com/glycerine/golua/lua"
)

var pp = verb.PP

// ConvError records a conversion error from value 'From' to value 'To'.
type ConvError struct {
	From interface{}
	To   interface{}
}

// ErrTableConv arises when some table entries could not be converted.
// The table conversion result is usable.
// TODO: Work out a more relevant name.
// TODO: Should it be a type instead embedding the actual error?
var ErrTableConv = errors.New("some table elements could not be converted")

func (l ConvError) Error() string {
	return fmt.Sprintf("cannot convert %v to %v", l.From, l.To)
}

// Lua 5.1 'lua_tostring' function only supports string and numbers. Extend it for internal purposes.
// From the Lua 5.3 source code.
func luaToString(L *lua.State, idx int) string {
	switch L.Type(idx) {
	case lua.LUA_TNUMBER:
		L.PushValue(idx)
		defer L.Pop(1)
		return L.ToString(-1)
	case lua.LUA_TSTRING:
		return L.ToString(-1)
	case lua.LUA_TBOOLEAN:
		b := L.ToBoolean(idx)
		if b {
			return "true"
		}
		return "false"
	case lua.LUA_TNIL:
		return "nil"
	}
	return fmt.Sprintf("%s: %p", L.LTypename(idx), L.ToPointer(idx))
}

func luaDesc(L *lua.State, idx int) string {
	return fmt.Sprintf("Lua value '%v' (%v)", luaToString(L, idx), L.LTypename(idx))
}

// NullT is the type of Null.
// Having a dedicated type allows us to make the distinction between zero values and Null.
type NullT int

// Map is an alias for map of strings.
type Map map[string]interface{}

var (
	// Null is the definition of 'luar.null' which is used in place of 'nil' when
	// converting slices and structs.
	Null = NullT(0)
)

var (
	tslice = typeof((*[]interface{})(nil))
	tmap   = typeof((*map[string]interface{})(nil))
	nullv  = reflect.ValueOf(Null)
)

// visitor holds the index to the table in LUA_REGISTRYINDEX with all the tables
// we ran across during a GoToLua conversion.
type visitor struct {
	L     *lua.State
	index int
}

func newVisitor(L *lua.State) visitor {
	var v visitor
	v.L = L
	v.L.NewTable()
	v.index = v.L.Ref(lua.LUA_REGISTRYINDEX)
	return v
}

func (v *visitor) close() {
	v.L.Unref(lua.LUA_REGISTRYINDEX, v.index)
}

// Mark value on top of the stack as visited using the registry index.
func (v *visitor) mark(val reflect.Value) {
	ptr := val.Pointer()
	v.L.RawGeti(lua.LUA_REGISTRYINDEX, v.index)
	// Copy value on top.
	v.L.PushValue(-2)
	// Set value to table.
	v.L.RawSeti(-2, int(ptr))
	v.L.Pop(1)
}

// Push visited value on top of the stack.
// If the value was not visited, return false and push nothing.
func (v *visitor) push(val reflect.Value) bool {
	ptr := val.Pointer()
	v.L.RawGeti(lua.LUA_REGISTRYINDEX, v.index)
	v.L.RawGeti(-1, int(ptr))
	if v.L.IsNil(-1) {
		// Not visited.
		v.L.Pop(2)
		return false
	}
	v.L.Replace(-2)
	return true
}

// Init makes and initializes a new pre-configured Lua state.
//
// It populates the 'luar' table with some helper functions/values:
//
//   method: ProxyMethod
//   unproxify: Unproxify
//
//   chan: MakeChan
//   complex: MakeComplex
//   map: MakeMap
//   slice: MakeSlice
//
//   null: Null
//
// It replaces the 'pairs'/'ipairs' functions with ProxyPairs/ProxyIpairs
// respectively, so that __pairs/__ipairs can be used, Lua 5.2 style. It allows
// for looping over Go composite types and strings.
//
// It also replaces the 'type' function with ProxyType.
//
// It is not required for using the 'GoToLua' and 'LuaToGo' functions.
func Init() *lua.State {
	var L = lua.NewState()
	L.OpenLibs()
	Register(L, "luar", Map{
		// Functions.
		"unproxify": Unproxify,

		"method": ProxyMethod,

		"chan":    MakeChan,
		"complex": Complex,
		"map":     MakeMap,
		"slice":   MakeSlice,

		// Values.
		"null": Null,
	})
	Register(L, "", Map{
		"ipairs": ProxyIpairs,
		"pairs":  ProxyPairs,
		"type":   ProxyType,
	})
	return L
}

func isNil(v reflect.Value) bool {
	nullables := [...]bool{
		reflect.Chan:      true,
		reflect.Func:      true,
		reflect.Interface: true,
		reflect.Map:       true,
		reflect.Ptr:       true,
		reflect.Slice:     true,
	}

	kind := v.Type().Kind()
	if int(kind) >= len(nullables) {
		return false
	}
	return nullables[kind] && v.IsNil()
}

func copyMapToTable(L *lua.State, v reflect.Value, visited visitor) {
	n := v.Len()
	L.CreateTable(0, n)
	visited.mark(v)
	for _, key := range v.MapKeys() {
		val := v.MapIndex(key)
		goToLua(L, key, true, visited)
		if isNil(val) {
			val = nullv
		}
		goToLua(L, val, false, visited)
		L.SetTable(-3)
	}
}

// Also for arrays.
func copySliceToTable(L *lua.State, v reflect.Value, visited visitor) {
	pp("top of copySliceToTable")

	vp := v
	for v.Kind() == reflect.Ptr {
		// For arrays.
		v = v.Elem()
	}

	n := v.Len()
	L.CreateTable(n, 0)
	if v.Kind() == reflect.Slice {
		visited.mark(v)
	} else if vp.Kind() == reflect.Ptr {
		visited.mark(vp)
	}

	for i := 0; i < n; i++ {
		L.PushInteger(int64(i + 1))
		val := v.Index(i)
		if isNil(val) {
			val = nullv
		}
		goToLua(L, val, false, visited)
		L.SetTable(-3)
	}
}

func copyStructToTable(L *lua.State, v reflect.Value, visited visitor) {
	// If 'vstruct' is a pointer to struct, use the pointer to mark as visited.
	vp := v
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	n := v.NumField()
	L.CreateTable(n, 0)
	if vp.Kind() == reflect.Ptr {
		visited.mark(vp)
	}

	for i := 0; i < n; i++ {
		st := v.Type()
		field := st.Field(i)
		key := field.Name
		tag := field.Tag.Get("lua")
		if tag != "" {
			key = tag
		}
		goToLua(L, key, false, visited)
		val := v.Field(i)
		goToLua(L, val, false, visited)
		L.SetTable(-3)
	}
}

func callGoFunction(L *lua.State, v reflect.Value, args []reflect.Value) []reflect.Value {
	defer func() {
		if x := recover(); x != nil {
			// jea debug:
			pp("recovering panic in luar.go, raising error x='%v'", x)
			L.RaiseError(fmt.Sprintf("error %s", x))
		}
	}()
	results := v.Call(args)
	return results
}

func goToLuaFunction(L *lua.State, v reflect.Value) lua.LuaGoFunction {
	switch f := v.Interface().(type) {
	case func(*lua.State) int:
		return f
	}

	t := v.Type()
	argsT := make([]reflect.Type, t.NumIn())
	for i := range argsT {
		argsT[i] = t.In(i)
	}

	return func(L *lua.State) int {
		var lastT reflect.Type
		isVariadic := t.IsVariadic()

		if isVariadic {
			n := len(argsT)
			lastT = argsT[n-1].Elem()
			argsT = argsT[:n-1]
		}

		args := make([]reflect.Value, len(argsT))
		for i, t := range argsT {
			val := reflect.New(t)
			err := LuaToGo(L, i+1, val.Interface())
			if err != nil {
				pp("problem point 1")
				L.RaiseError(fmt.Sprintf("cannot convert Go function argument #%v: %v", i, err))
			}
			args[i] = val.Elem()
		}

		if isVariadic {
			pp("we have a variadic function!. len(argsT)=%v", len(argsT))
			n := L.GetTop()
			for i := len(argsT) + 1; i <= n; i++ {
				// jea: assumes any varargs in the actual call have been
				// pushed onto the stack.

				val := reflect.New(lastT)
				pp("about to call LuaToGo with val from lastT: '%#v'/%T", val.Interface(), val.Interface())
				err := LuaToGo(L, i, val.Interface())
				if err != nil {
					pp("problem point 2")
					L.RaiseError(fmt.Sprintf("cannot convert Go function argument #%v: %v", i, err))
				}
				args = append(args, val.Elem())
			}
			argsT = argsT[:len(argsT)+1]
		}
		results := callGoFunction(L, v, args)
		for _, val := range results {
			if val.Kind() == reflect.Struct {
				// If the function returns a struct (and not a pointer to a struct),
				// calling GoToLua directly will convert it to a table, making the
				// methods inaccessible. We work around that issue by forcibly passing a
				// pointer to a struct.
				valp := reflect.New(val.Type())
				valp.Elem().Set(val)
				val = valp
			}
			GoToLuaProxy(L, val)
		}
		return len(results)
	}
}

// GoToLua pushes a Go value 'val' on the Lua stack.
//
// It unboxes interfaces.
//
// Pointers are followed recursively. Slices, structs and maps are copied over as tables.
func GoToLua(L *lua.State, a interface{}) {
	visited := newVisitor(L)
	goToLua(L, a, false, visited)
	visited.close()
}

// GoToLuaProxy is like GoToLua but pushes a proxy on the Lua stack when it makes sense.
//
// A proxy is a Lua userdata that wraps a Go value.
//
// Pointers are preserved.
//
// Structs and arrays need to be passed as pointers to be proxified, otherwise
// they will be copied as tables.
//
// Predeclared scalar types are never proxified as they have no methods.
func GoToLuaProxy(L *lua.State, a interface{}) {
	visited := newVisitor(L)
	goToLua(L, a, true, visited)
	visited.close()
}

// TODO: Check if we really need multiple pointer levels since pointer methods
// can be called on non-pointers.
func goToLua(L *lua.State, a interface{}, proxify bool, visited visitor) {
	pp("++ goToLua top. a='%#v'/type='%T', proxify='%v', visited='%#v'", a, a, proxify, visited)
	/* jea debug int64:
	switch x := a.(type) {
	case reflect.Value:
		ty := x.Type()
		pp("ty = '%v', kind='%v'", ty.String(), ty.Kind())
		switch ty.Kind() {
		case reflect.Int, reflect.Int64:
			y := x.Int()
			pp("goToLua: we have an int, y = %v", y)
			if y == 2 {
				pp("in luar.go, where called from?")
				fmt.Printf("%s\n", string(debug.Stack()))
			}
		}
	}
	*/

	var v reflect.Value
	v, ok := a.(reflect.Value)
	if !ok {
		v = reflect.ValueOf(a)
	}
	if !v.IsValid() {
		L.PushNil()
		return
	}

	if v.Kind() == reflect.Interface && !v.IsNil() {
		// Unbox interface.
		v = reflect.ValueOf(v.Interface())
	}

	// Follow pointers if not proxifying. We save the original pointer Value in case we proxify.
	vp := v
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if !v.IsValid() {
		L.PushNil()
		return
	}

	// As a special case, we always proxify Null, the empty element for slices and maps.
	if v.CanInterface() && v.Interface() == Null {
		makeValueProxy(L, v, cInterfaceMeta)
		return
	}

	switch v.Kind() {
	case reflect.Float64, reflect.Float32:
		if proxify && isNewType(v.Type()) {
			makeValueProxy(L, vp, cNumberMeta)
		} else {
			L.PushNumber(v.Float())
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		pp("in goToLua at switch v.Kind(), Int types")
		if proxify && isNewType(v.Type()) {
			pp("in goToLua at switch v.Kind(), Int types, calling makeValueProxy")
			makeValueProxy(L, vp, cNumberMeta)
		} else {
			pp("in goToLua at switch v.Kind(), Int types, doing PushInt64")
			L.PushInt64(v.Int())
			pp("in goToLua at switch v.Kind(), Int types, *after* PushInt64")
			if verb.VerboseVerbose {
				DumpLuaStack(L)
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if proxify && isNewType(v.Type()) {
			makeValueProxy(L, vp, cNumberMeta)
		} else {
			L.PushUint64(v.Uint())
		}
	case reflect.String:
		if proxify && isNewType(v.Type()) {
			makeValueProxy(L, vp, cStringMeta)
		} else {
			L.PushString(v.String())
		}
	case reflect.Bool:
		if proxify && isNewType(v.Type()) {
			makeValueProxy(L, vp, cInterfaceMeta)
		} else {
			L.PushBoolean(v.Bool())
		}
	case reflect.Complex128, reflect.Complex64:
		makeValueProxy(L, vp, cComplexMeta)
	case reflect.Array:
		// It needs be a pointer to be a proxy, otherwise values won't be settable.
		if proxify && vp.Kind() == reflect.Ptr {
			makeValueProxy(L, vp, cSliceMeta)
		} else {
			// See the case of struct.
			if vp.Kind() == reflect.Ptr && visited.push(vp) {
				return
			}
			copySliceToTable(L, vp, visited)
		}
	case reflect.Slice:
		if proxify {
			makeValueProxy(L, vp, cSliceMeta)
		} else {
			if visited.push(v) {
				return
			}
			pp("luar.go: reflect.Slice, about to call copySliceToTable")
			copySliceToTable(L, v, visited)
		}
	case reflect.Map:
		if proxify {
			makeValueProxy(L, vp, cMapMeta)
		} else {
			if visited.push(v) {
				return
			}
			copyMapToTable(L, v, visited)
		}
	case reflect.Struct:
		if proxify && vp.Kind() == reflect.Ptr {
			if vp.CanInterface() {
				switch v := vp.Interface().(type) {
				case error:
					L.PushString(v.Error())
				case *LuaObject:
					// TODO: Move out of 'proxify' condition? LuaObject is meant to be
					// manipulated from the Go side, it is not useful in Lua.
					if v.l == L {
						v.Push()
					} else {
						// TODO: What shall we do when LuaObject state is not the current
						// state? Copy across states? Is it always possible?
						L.PushNil()
					}
				default:
					makeValueProxy(L, vp, cStructMeta)
				}
			} else {
				makeValueProxy(L, vp, cStructMeta)
			}
		} else {
			// Use vp instead of v to detect cycles from the very first element, if a pointer.
			if vp.Kind() == reflect.Ptr && visited.push(vp) {
				return
			}
			copyStructToTable(L, vp, visited)
		}
	case reflect.Chan:
		makeValueProxy(L, vp, cChannelMeta)
	case reflect.Func:
		L.PushGoFunction(goToLuaFunction(L, v))
	default:
		if val, ok := v.Interface().(error); ok {
			L.PushString(val.Error())
		} else if v.IsNil() {
			L.PushNil()
		} else {
			makeValueProxy(L, vp, cInterfaceMeta)
		}
	}
}

func luaIsEmpty(L *lua.State, idx int) bool {
	L.PushNil()
	if idx < 0 {
		idx--
	}
	if L.Next(idx) != 0 {
		L.Pop(2)
		return false
	}
	return true
}

func luaMapLen(L *lua.State, idx int) int {
	L.PushNil()
	if idx < 0 {
		idx--
	}
	len := 0
	for L.Next(idx) != 0 {
		len++
		L.Pop(1)
	}
	return len
}

func copyTableToMap(L *lua.State, idx int, v reflect.Value, visited map[uintptr]reflect.Value) (status error) {
	t := v.Type()
	if v.IsNil() {
		v.Set(reflect.MakeMap(t))
	}
	te, tk := t.Elem(), t.Key()

	// See copyTableToSlice.
	ptr := L.ToPointer(idx)
	if !luaIsEmpty(L, idx) {
		visited[ptr] = v
	}

	L.PushNil()
	if idx < 0 {
		idx--
	}
	for L.Next(idx) != 0 {
		// key at -2, value at -1
		key := reflect.New(tk).Elem()
		err := luaToGo(L, -2, key, visited)
		if err != nil {
			// here is where fmt.Sprintf( table) is failing.
			pp("ErrTableConv about to be status, since luaToGo failed for key at -2: '%v'. tk='%s', key='%s'. stack:\n%s\n", err, tk, key,
				string(debug.Stack()))
			status = ErrTableConv
			L.Pop(1)
			continue
		}
		val := reflect.New(te).Elem()
		err = luaToGo(L, -1, val, visited)
		if err != nil {
			pp("ErrTableConv about to be status, since luaToGo failed for key '%s'", key.Interface())
			status = ErrTableConv
			L.Pop(1)
			continue
		}
		v.SetMapIndex(key, val)
		L.Pop(1)
	}

	return
}

// Also for arrays, but isSlice will be false. TODO: Create special function for arrays?
func copyTableToSlice(L *lua.State, idx int, v reflect.Value, visited map[uintptr]reflect.Value, isSlice bool) (status error) {
	pp("top of copyTableToSlice. here is stack:")
	if verb.VerboseVerbose {
		DumpLuaStack(L)
	}

	t := v.Type()
	n := int(L.ObjLen(idx))
	pp("in copyTableToSlice, n='%v', t='%v'. top=%v, idx=%v", n, t, L.GetTop(), idx)

	// detect _gi_Slice and specialize for it.
	if isSlice {
		L.GetGlobal("_giPrivateSliceProps") // stack++
	} else {
		L.GetGlobal("_giPrivateArrayProps") // stack++
	}
	if !L.IsNil(-1) {
		// we are running under `gi`
		// is this a _gi_Slice? it is if the _giPrivateSliceProps key is present.

		// since we increased the stack depth by 1, adjust idx.
		adj := idx
		if idx < 0 && idx > -10000 {
			adj--
		}

		pp("we are running under `gi`, _giPrivateSliceProps key was found in _G. top is now %v, idx=%v, adj=%v", L.GetTop(), idx, adj)

		// get table[key]. replaces key with value,
		// i.e. replace the key _giPrivateSliceProps with
		//  the actual table it represents.
		L.GetTable(adj)
		pp("under `gi`, after GetTable(adj), top is %v, and Top is nil: %v", L.GetTop(), L.IsNil(-1))
		if !L.IsNil(-1) {
			// yes, is _gi_Slice
			// leave the props on the top of the stack, we'll use
			// them immediately.
			return copyGiTableToSlice(L, adj, v, visited, isSlice)
		} else {
			L.Pop(1)
		}
	} else {
		L.Pop(1)
	}

	// Adjust the length of the array/slice.
	if n > v.Len() {
		if t.Kind() == reflect.Array {
			n = v.Len()
		} else {
			// Slice
			v.Set(reflect.MakeSlice(t, n, n))
		}
	} else if n < v.Len() {
		if t.Kind() == reflect.Array {
			// Nullify remaining elements.
			for i := n; i < v.Len(); i++ {
				v.Index(i).Set(reflect.Zero(t.Elem()))
			}
		} else {
			// Slice
			v.SetLen(n)
		}
	}

	// Do not add empty slices to the list of visited elements.
	// The empty Lua table is a single instance object and gets re-used across maps, slices and others.
	// Arrays cannot be cyclic since the interface type will ask for slices.
	if n > 0 && t.Kind() != reflect.Array {
		ptr := L.ToPointer(idx)
		visited[ptr] = v
	}

	te := t.Elem()
	for i := 1; i <= n; i++ {
		L.RawGeti(idx, i)
		val := reflect.New(te).Elem()
		err := luaToGo(L, -1, val, visited)
		if err != nil {
			pp("ErrTableConv about to be status, since luaToGo failed for val '%v'", val.Interface())
			status = ErrTableConv
			L.Pop(1)
			continue
		}
		v.Index(i - 1).Set(val)
		L.Pop(1)
	}

	return
}

func copyTableToStruct(L *lua.State, idx int, v reflect.Value, visited map[uintptr]reflect.Value) (status error) {
	pp("top of copyTableToStruct")
	t := v.Type()

	// See copyTableToSlice.
	ptr := L.ToPointer(idx)
	if !luaIsEmpty(L, idx) {
		visited[ptr] = v.Addr()
	}

	// Associate Lua keys with Go fields: tags have priority over matching field
	// name.
	fields := map[string]string{}
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("lua")
		if tag != "" {
			fields[tag] = field.Name
			continue
		}
		fields[field.Name] = field.Name
		pp("added to fields, field.Name='%v'", field.Name)
	}
	pp("fields is now '%#v'", fields)

	L.PushNil()
	if idx < 0 {
		idx--
	}
	for L.Next(idx) != 0 {
		L.PushValue(-2)
		// Warning: ToString changes the value on stack.
		key := L.ToString(-1)
		L.Pop(1)
		f := v.FieldByName(fields[key])
		// jea: set private fields too.
		//if f.CanSet() {
		val := reflect.New(f.Type()).Elem()
		err := luaToGo(L, -1, val, visited)
		if err != nil {
			pp("ErrTableConv about to be status, since luaToGo failed for val '%v'", val.Interface())
			status = ErrTableConv
			L.Pop(1)
			continue
		}
		//f.Set(val)
		setField(f, val)
		//} // jea
		L.Pop(1)
	}

	return
}

// setField works on private and public fields
func setField(fld, val reflect.Value) {
	fieldPtr := reflect.NewAt(fld.Type(), unsafe.Pointer(fld.UnsafeAddr()))
	fieldPtr.Elem().Set(val)
}

// LuaToGo converts the Lua value at index 'idx' to the Go value.
//
// The Go value must be a non-nil pointer.
//
// Conversions to string and numbers are straightforward.
//
// Lua 'nil' is converted to the zero value of the specified Go value.
//
// If the Lua value is non-nil, pointers are dereferenced (multiple times if
// required) and the pointed value is the one that is set. If 'nil', then the Go
// pointer is set to 'nil'. To set a pointer's value to its zero value, use
// 'luar.null'.
//
// The Go value can be an interface, in which case the type is inferred. When
// converting a table to an interface, the Go value is a []interface{} slice if
// all its elements are indexed consecutively from 1, or a
// map[string]interface{} otherwise.
//
// Existing entries in maps and structs are kept. Arrays and slices are reset.
//
// Nil maps and slices are automatically allocated.
//
// Proxies are unwrapped to the Go value, if convertible.
// Userdata that is not a proxy will be converted to a LuaObject if the Go value
// is an interface or a LuaObject.
func LuaToGo(L *lua.State, idx int, a interface{}) error {
	// jea debug:
	//verb.VerboseVerbose = true

	// LuaToGo should not pop the Lua stack to be consistent with L.ToString(), etc.
	// It is also easier in practice when we want to keep working with the value on stack.

	v := reflect.ValueOf(a)
	// TODO: Test interfaces with methods.
	// TODO: Allow unreferenced map? encoding/json does not do it.
	if v.Kind() != reflect.Ptr {
		return errors.New("not a pointer")
	}
	if v.IsNil() {
		return errors.New("nil pointer")
	}

	v = v.Elem()
	// If the Lua value is 'nil' and the Go value is a pointer, nullify the pointer.
	if v.Kind() == reflect.Ptr && L.IsNil(idx) {
		v.Set(reflect.Zero(v.Type()))
		return nil
	}

	return luaToGo(L, idx, v, map[uintptr]reflect.Value{})
}

/*
from lua.h
** basic types

define LUA_TNONE		(-1)
define LUA_TNIL		        0
define LUA_TBOOLEAN		    1
define LUA_TLIGHTUSERDATA	2
define LUA_TNUMBER		    3
define LUA_TSTRING		4
define LUA_TTABLE		5
define LUA_TFUNCTION	6
define LUA_TUSERDATA	7
define LUA_TTHREAD		8
*/

func luaToGo(L *lua.State, idx int, v reflect.Value, visited map[uintptr]reflect.Value) error {
	pp("-- top of luaToGo, here is stack:")
	if verb.VerboseVerbose {
		DumpLuaStack(L)
	}

	// Derefence 'v' until a non-pointer.
	// This initializes the values, which will be useless effort if the conversion fails.
	// This must be done here so that the copyTable* functions can also call luaToGo on pointers.
	vp := v
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		vp = v
		v = v.Elem()
	}
	kind := v.Kind()

	ltype := L.Type(idx)
	// Typename() is useless and wrong.
	//ltypename := L.Typename(idx)
	pp("ltype = '%v'", ltype)

	switch ltype {
	case lua.LUA_TNIL:
		pp("luar.go, type of idx == LUA_TNIL")
		v.Set(reflect.Zero(v.Type()))
	case lua.LUA_TBOOLEAN:
		pp("luar.go, type of idx == LUA_TBOOLEAN")
		if kind != reflect.Bool && kind != reflect.Interface {
			return ConvError{From: luaDesc(L, idx), To: v.Type()}
		}
		v.Set(reflect.ValueOf(L.ToBoolean(idx)))
	case lua.LUA_TNUMBER:
		pp("luar.go, type of idx == LUA_TNUMBER")
		switch k := unsizedKind(v); k {
		case reflect.Int64, reflect.Uint64, reflect.Float64, reflect.Interface:
			// We do not use ToInteger as it may truncate the value. Let Go truncate
			// instead in Convert().
			f := reflect.ValueOf(L.ToNumber(idx))
			v.Set(f.Convert(v.Type()))
		case reflect.Complex128:
			v.SetComplex(complex(L.ToNumber(idx), 0))
		default:
			return ConvError{From: luaDesc(L, idx), To: v.Type()}
		}
	case lua.LUA_TSTRING:
		pp("luar.go, type of idx == LUA_TSTRING: '%s'", L.ToString(idx))
		if kind != reflect.String && kind != reflect.Interface {
			return ConvError{From: luaDesc(L, idx), To: v.Type()}
		}
		v.Set(reflect.ValueOf(L.ToString(idx)))
	case lua.LUA_TUSERDATA:
		pp("luar.go, type of idx == LUA_TUSERDATA")
		if isValueProxy(L, idx) {
			pp("luar.go, type of idx == LUA_TUSERDATA, isValueProxy is true")
			val, typ := valueOfProxy(L, idx)
			if val.Interface() == Null {
				// Special case for Null.
				v.Set(reflect.Zero(v.Type()))
				return nil
			}

			for !typ.ConvertibleTo(v.Type()) && val.Kind() == reflect.Ptr {
				val = val.Elem()
				typ = typ.Elem()
			}
			if !typ.ConvertibleTo(v.Type()) {
				return ConvError{From: fmt.Sprintf("proxy (%v)", typ), To: v.Type()}
			}
			// We automatically convert between types. This behaviour is consistent
			// with LuaToGo conversions elsewhere.
			v.Set(val.Convert(v.Type()))
			return nil
		} else if kind != reflect.Interface || v.Type() != reflect.TypeOf(LuaObject{}) {
			pp("luar.go, type of idx == LUA_TUSERDATA, ConvError happening!??, from: '%s', to: '%s'", luaDesc(L, idx), v.Type())
			// jea try this, so that we wrap into a lua ref
			// This makes gi: fmt.Printf("%v", fmt.Printf) work.
			v.Set(reflect.ValueOf(NewLuaObject(L, idx)))
			//return ConvError{From: luaDesc(L, idx), To: v.Type()}
		}
		// Wrap the userdata into a LuaObject.
		v.Set(reflect.ValueOf(NewLuaObject(L, idx)))
	case lua.LUA_TTABLE:
		pp("luar.go, Ltype of idx == LUA_TTABLE or ('%v')\n", ltype)

		// TODO: Check what happens if visited is not of the right type.
		ptr := L.ToPointer(idx)
		if val, ok := visited[ptr]; ok {
			pp("visited[ptr] is true")
			if v.Kind() == reflect.Struct {
				vp.Set(val)
			} else {
				v.Set(val)
			}
			return nil
		}
		pp("visited[ptr] was false, kind='%#v'/'%v'", kind, kind)

		switch kind {
		case reflect.Array:
			return copyTableToSlice(L, idx, v, visited, false)
		case reflect.Slice:
			return copyTableToSlice(L, idx, v, visited, true)
		case reflect.Map:
			return copyTableToMap(L, idx, v, visited)
		case reflect.Struct:
			return copyTableToStruct(L, idx, v, visited)
		case reflect.Interface:
			// jea: the original L.ObjLen reults was wrong b/c our _gi_Slice start indexing at 0 not 1.
			//n := int(L.ObjLen(idx)) // does not call __len metamethod. Problem.
			n := getLenByCallingMetamethod(L, idx)
			pp("n back from getLenByCallingMetamethod = %v at idx=%v", n, idx)

			switch v.Elem().Kind() {
			case reflect.Map:
				return copyTableToMap(L, idx, v.Elem(), visited)
			case reflect.Slice:
				// Need to make/resize the slice here since interface values are not adressable.
				v.Set(reflect.MakeSlice(v.Elem().Type(), n, n))
				return copyTableToSlice(L, idx, v.Elem(), visited, true)
				// jea debug: add default: case
			default:
				pp("v.Elem().Kind() = '%#v', v='%#v'/type='%T'", v.Elem().Kind(), v, v) // 0x0, nil interface, reflect.Value
			}

			/* jea, not sure why this map conversion is here, but
			               it messes up imp_test 065 as just one example

						mapLen := luaMapLen(L, idx)
						pp("jea: mapLen = %v, n = %v", mapLen, n)
						if mapLen != n {
							v.Set(reflect.MakeMap(tmap))
							// jea: why are we copying a vararg table to a map???
							return copyTableToMap(L, idx, v.Elem(), visited)
						}
			*/
			v.Set(reflect.MakeSlice(tslice, n, n))
			return copyTableToSlice(L, idx, v.Elem(), visited, true)
		default:
			pp("luar.go ConvError: from '%v' to '%v'\n stack:\n%s\n",
				luaDesc(L, idx), v.Type(),
				string(debug.Stack()))
			return ConvError{From: luaDesc(L, idx), To: v.Type()}
		}
	case 10: // LUA_TCDATA aka cdata
		pp("luaToGo cdata case, L.Type(idx) = '%v'", L.Type(idx))
		ctype := L.LuaJITctypeID(-1)
		pp("luar.go sees ctype = %v", ctype)
		switch ctype {
		case 5: //  int8
		case 6: //  uint8
		case 7: //  int16
		case 8: //  uint16
		case 9: //  int32
		case 10: //  uint32
		case 11: //  int64
			val := L.CdataToInt64(idx)
			f := reflect.ValueOf(val)
			vi := v.Interface()
			pp("luar.go calling L.CdataToInt64, got val=%v/'%T', v=%v/'%T'", val, val, vi, vi)
			//v.Set(f.Convert(v.Type())) // don't do this universally,
			// since it will coerce uints
			// and then we won't get the type mistmatch error that is important.
			// Instead let v.Set(f) panic on wrong type.

			// allow int64 to convert to int
			if v.Kind() == reflect.Int {
				v.Set(f.Convert(v.Type()))
				//setField(f.Convert(v.Type()), v)
			} else {

				if !canAndDidAssign(&f, &v) {
					return ConvError{From: luaDesc(L, idx), To: v.Type()}
				}
				// huh?
				// go test -v -run TestArray
				// panic: reflect.Set: value of type int64 is not assignable to type string
				//v.Set(f)
				//setField(f, v)
			}
			return nil
		case 12: //  uint64
			val := L.CdataToUint64(idx)
			//pp("luar.go calling L.CdataToUint64, got val='%#v'", val)
			f := reflect.ValueOf(val)
			//v.Set(f.Convert(v.Type())) // don't do this, since it will
			// coerce int64, and then we won't get the approprirate type
			// mismatch error. Instead, let v.Set(f) panic on wrong type.

			// allow uint64 to convert to uint
			if v.Kind() == reflect.Uint {
				v.Set(f.Convert(v.Type()))
			} else {
				v.Set(f)
			}
			return nil
		case 13: //  float32
		case 14: //  float64

		case 0: // means it wasn't a ctype
		}

		return ConvError{From: luaDesc(L, idx), To: v.Type()}

	default:
		return ConvError{From: luaDesc(L, idx), To: v.Type()}
	}

	return nil
}

func isNewType(t reflect.Type) bool {
	types := [...]reflect.Type{
		reflect.Invalid:    nil, // Invalid Kind = iota
		reflect.Bool:       typeof((*bool)(nil)),
		reflect.Int:        typeof((*int)(nil)),
		reflect.Int8:       typeof((*int8)(nil)),
		reflect.Int16:      typeof((*int16)(nil)),
		reflect.Int32:      typeof((*int32)(nil)),
		reflect.Int64:      typeof((*int64)(nil)),
		reflect.Uint:       typeof((*uint)(nil)),
		reflect.Uint8:      typeof((*uint8)(nil)),
		reflect.Uint16:     typeof((*uint16)(nil)),
		reflect.Uint32:     typeof((*uint32)(nil)),
		reflect.Uint64:     typeof((*uint64)(nil)),
		reflect.Uintptr:    typeof((*uintptr)(nil)),
		reflect.Float32:    typeof((*float32)(nil)),
		reflect.Float64:    typeof((*float64)(nil)),
		reflect.Complex64:  typeof((*complex64)(nil)),
		reflect.Complex128: typeof((*complex128)(nil)),
		reflect.String:     typeof((*string)(nil)),
	}

	pt := types[int(t.Kind())]
	return pt != t
}

// Register makes a number of Go values available in Lua code.
// 'values' is a map of strings to Go values.
//
// - If table is non-nil, then create or reuse a global table of that name and
// put the values in it.
//
// - If table is '' then put the values in the global table (_G).
//
// - If table is '*' then assume that the table is already on the stack.
func Register(L *lua.State, table string, values Map) {
	pop := true
	if table == "*" {
		pop = false
	} else if len(table) > 0 {
		L.GetGlobal(table)
		if L.IsNil(-1) {
			L.Pop(1)
			L.NewTable()
			L.SetGlobal(table)
			L.GetGlobal(table)
		}
	} else {
		L.GetGlobal("_G")
	}
	for name, val := range values {
		GoToLuaProxy(L, val)
		L.SetField(-2, name)
	}
	if pop {
		L.Pop(1)
	}
}

// Closest we'll get to a typeof operator.
func typeof(a interface{}) reflect.Type {
	return reflect.TypeOf(a).Elem()
}

// jea
func DumpLuaStack(L *lua.State) {
	fmt.Printf("\n%s\n", DumpLuaStackAsString(L))
}

func DumpLuaStackAsString(L *lua.State) (s string) {
	var top int

	top = L.GetTop()
	s += fmt.Sprintf("========== begin DumpLuaStack: top = %v\n", top)
	for i := top; i >= 1; i-- {
		t := L.Type(i)
		s += fmt.Sprintf("DumpLuaStack: i=%v, t= %v\n", i, t)
		switch t {
		case lua.LUA_TSTRING:
			s += fmt.Sprintf(" String : \t%v\n", L.ToString(i))
		case lua.LUA_TBOOLEAN:
			s += fmt.Sprintf(" Bool : \t\t%v\n", L.ToBoolean(i))
		case lua.LUA_TNUMBER:
			s += fmt.Sprintf(" Number : \t%v\n", L.ToNumber(i))
		case lua.LUA_TTABLE:
			s += fmt.Sprintf(" Table : \n%s\n", dumpTableString(L, i))

		case 10: // LUA_TCDATA aka cdata
			//pp("Dump cdata case, L.Type(idx) = '%v'", L.Type(i))
			ctype := L.LuaJITctypeID(i)
			//pp("luar.go Dump sees ctype = %v", ctype)
			switch ctype {
			case 5: //  int8
			case 6: //  uint8
			case 7: //  int16
			case 8: //  uint16
			case 9: //  int32
			case 10: //  uint32
			case 11: //  int64
				val := L.CdataToInt64(i)
				s += fmt.Sprintf(" int64: '%v'\n", val)
			case 12: //  uint64
				val := L.CdataToUint64(i)
				s += fmt.Sprintf(" uint64: '%v'\n", val)
			case 13: //  float32
			case 14: //  float64

			case 0: // means it wasn't a ctype
			}

		case lua.LUA_TUSERDATA:
			s += fmt.Sprintf(" Type(code %v/ LUA_TUSERDATA) : no auto-print available.\n", t)
		case lua.LUA_TFUNCTION:
			s += fmt.Sprintf(" Type(code %v/ LUA_TFUNCTION) : no auto-print available.\n", t)
		default:
			s += fmt.Sprintf(" Type(code %v) : no auto-print available.\n", t)
		}
	}
	s += fmt.Sprintf("========= end of DumpLuaStack\n")
	return
}

func dumpTableString(L *lua.State, index int) (s string) {

	// Push another reference to the table on top of the stack (so we know
	// where it is, and this function can work for negative, positive and
	// pseudo indices
	L.PushValue(index)
	// stack now contains: -1 => table
	L.PushNil()
	// stack now contains: -1 => nil; -2 => table
	for L.Next(-2) != 0 {

		// stack now contains: -1 => value; -2 => key; -3 => table
		// copy the key so that lua_tostring does not modify the original
		L.PushValue(-2)
		// stack now contains: -1 => key; -2 => value; -3 => key; -4 => table
		key := L.ToString(-1)
		value := L.ToString(-2)
		s += fmt.Sprintf("'%s' => '%s'\n", key, value)
		// pop value + copy of key, leaving original key
		L.Pop(2)
		// stack now contains: -1 => key; -2 => table
	}
	// stack now contains: -1 => table (when lua_next returns 0 it pops the key
	// but does not push anything.)
	// Pop table
	L.Pop(1)
	// Stack is now the same as it was on entry to this function
	return
}

func giSliceGetRawHelper(L *lua.State, idx int, v reflect.Value, visited map[uintptr]reflect.Value) (int, reflect.Type) {
	pp("top of giSliceGetRawHelper. idx=%v, here is stack:", idx)
	if verb.VerboseVerbose {
		DumpLuaStack(L)
	}

	t := v.Type()
	getfield(L, -1, "len")
	if L.IsNil(-1) {
		panic("what? should be a `len` member of props for _gi_Slice")
	}
	n := int(L.ToNumber(-1))
	L.Pop(1)
	pp("copyGiTableToSlice after getting n=%v, stack is:", n)
	if verb.VerboseVerbose {
		DumpLuaStack(L)
	}

	/* sample
	       luar.go:1125 copyGiTableToSlice after getting n=3, stack is:
		   ========== begin DumpLuaStack: top = 4
		   DumpLuaStack: i=4, t= 5
		    Table :
		   'len' => '3'
		   'typeKind' => 'int'

		   DumpLuaStack: i=3, t= 5
		    Table :
		   '' => ''
		   '' => ''
		   'Typeof' => '_gi_Slice'

		   DumpLuaStack: i=2, t= 0
		    Type(code 0) : no auto-print available.
		   DumpLuaStack: i=1, t= 5
		    Table :

		   ========= end of DumpLuaStack
	*/

	// get the raw table
	L.GetGlobal("_giPrivateRaw") // stack++
	if L.IsNil(-1) {
		panic(`could not lookup "_giPrivateRaw" in global table`)
	}

	// since we increased the stack depth by 1, adjust idx.
	if idx < 0 && idx > -10000 {
		idx--
	}

	pp("found the global string _giPrivateRaw, here is stack, with adjusted idx=%v:", idx)
	if verb.VerboseVerbose {
		DumpLuaStack(L)
	}

	// get table[key]. replaces key with value,
	// i.e. replace the key _giPrivateRaw with
	//  the actual table it represents.
	L.GetTable(idx)
	pp("under `gi`, after GetTable(idx=%v), top is %v, and Top is nil: %v", idx, L.GetTop(), L.IsNil(-1))
	if L.IsNil(-1) {
		panic("_giPrivateRaw not found in _gi_Slice outer value!")
	}
	pp("in copyGiTableToSlice. after fetching raw table to the top of the stack, here is stack:")
	if verb.VerboseVerbose {
		DumpLuaStack(L)
	}

	// just leave the raw, remove the props table and the outer table.
	L.Replace(-3)
	L.Pop(1)
	pp("after popping the props and outer and leaving just the raw:")
	if verb.VerboseVerbose {
		DumpLuaStack(L)
	}

	return n, t
}

// props is on top of stack. The actual table at idx, which props describes.
func copyGiTableToSlice(L *lua.State, idx int, v reflect.Value, visited map[uintptr]reflect.Value, isSlice bool) (status error) {
	pp("top of copyGiTableToSlice. idx=%v, here is stack:", idx)
	if verb.VerboseVerbose {
		DumpLuaStack(L)
	}

	// extract out the raw underlying table
	n, t := giSliceGetRawHelper(L, idx, v, visited)

	pp("in copyGiTableToSlice, n='%v', t='%v'", n, t)

	// Adjust the length of the array/slice.
	if n > v.Len() {
		if t.Kind() == reflect.Array {
			n = v.Len()
		} else {
			// Slice
			v.Set(reflect.MakeSlice(t, n, n))
		}
	} else if n < v.Len() {
		if t.Kind() == reflect.Array {
			// Nullify remaining elements.
			for i := n; i < v.Len(); i++ {
				v.Index(i).Set(reflect.Zero(t.Elem()))
			}
		} else {
			// Slice
			v.SetLen(n)
		}
	}

	// Do not add empty slices to the list of visited elements.
	// The empty Lua table is a single instance object and gets re-used across maps, slices and others.
	// Arrays cannot be cyclic since the interface type will ask for slices.
	if n > 0 && t.Kind() != reflect.Array {
		ptr := L.ToPointer(idx)
		visited[ptr] = v
	}

	te := t.Elem()
	for i := 0; i < n; i++ {
		L.RawGeti(idx, i)
		val := reflect.New(te).Elem()
		err := luaToGo(L, -1, val, visited)
		if err != nil {
			pp("ErrTableConv about to be status, since luaToGo failed for val '%v'", val.Interface())
			status = ErrTableConv
			L.Pop(1)
			continue
		}
		v.Index(i).Set(val)
		L.Pop(1)
	}

	return
}

// getfield will
// assume that table is on the stack top, and
// returns with the value (that which corresponds to key) on
// the top of the stack.
// If value not present, then a nil is on top of the stack.
// To clean the stack completely, Pop(1).
func getfield(L *lua.State, tableIdx int, key string) {
	// copy up front, so that we work for
	// pseudo indexes, abs, and relative.
	L.PushValue(tableIdx)

	// setup to query.
	L.PushString(key)

	// lua_gettable: It receives the
	// position of the table in the stack,
	// pops the key from the top stack, and
	// pushes the corresponding value.
	//
	// void lua_gettable (lua_State *L, int index);
	// Pushes onto the stack the value t[k],
	// where t is the value at the given valid index
	// and k is the value at the top of the stack.
	//
	// This function pops the key from the stack
	// (putting the resulting value in its place).
	// As in Lua, this function may trigger a
	// metamethod for the "index" event (see §2.8).
	//
	L.GetTable(-2) // get table[key]

	// remove the copy of the table we made up front.
	L.Remove(-2)
}

// jea add
func getLenByCallingMetamethod(L *lua.State, idx int) int {
	pp("top of getLenByCallingMetamethod for idx=%v, here is stack:", idx)
	if verb.VerboseVerbose {
		DumpLuaStack(L)
	}
	pp("trace:\n%s\n", string(debug.Stack()))

	top := L.GetTop()
	//
	// lua_getmetatable: Pushes onto the stack the
	// metatable of the value at the given acceptable
	// index. If the index is not valid, or if the
	// value does not have a metatable, the function
	// returns 0 and pushes nothing on the stack.
	//
	found := L.GetMetaTable(idx)
	if !found {
		return int(L.ObjLen(idx))
	}

	L.PushString("__len")

	// lua_gettable: It receives the
	// position of the table in the stack,
	// pops the key from the top stack, and
	// pushes the corresponding value.
	// lua_rawget: same but no metamethods.
	L.RawGet(-2) // get table[key]
	if L.IsNil(-1) {
		// __len method not found in metatable
		return int(L.ObjLen(idx))
	}
	pp("after RawGet was not nil, top =%v, stack is", top)
	if verb.VerboseVerbose {
		DumpLuaStack(L)
	}

	defer L.SetTop(top)

	// INVAR: __len method is on top of stack.

	// stack: __len method, the metatable, _gi_Slice table

	//	pp("we think __len method is top of stack, followed by metable.")
	//	if verb.VerboseVerbose {
	//		DumpLuaStack(L)
	//	}

	// gotta get rid of the metable first, prior to the call, since
	// __len method expects the actual table to be its self parameter.
	L.Remove(-2)

	pp("after Remove(-2), top =%v, stack is", top)
	if verb.VerboseVerbose {
		DumpLuaStack(L)
	}

	// Setup the call with the table as the argument
	L.PushValue(idx)

	pp("about to call __len, stack is:")
	if verb.VerboseVerbose {
		DumpLuaStack(L)
	}

	err := L.Call(1, 1)
	if err != nil {
		panic(err)
	}
	fLen := L.ToNumber(-1)
	// NOTE: won't work for tables of len > 2^52 or 4 peta items.
	// Since that's way bigger than viable ram, we don't worry about it here.
	pp("getLenByCallingMetamethod returning %v", fLen)
	return int(fLen)
}

func canAndDidAssign(f, v *reflect.Value) (res bool) {
	pp("top of canAndDidAssign, f.Type='%v', v.Type='%T'", f.Interface(), v.Interface()) // 'string'

	res = true
	defer func() {
		if r := recover(); r != nil {
			pp("canAndDidAssign recover caught: '%v'", r)
			res = false
		}
	}()
	v.Set(*f)
	//f.Convert(v.Type())
	return
}

package luar

import (
	"errors"
	"reflect"

	"github.com/aarzilli/golua/lua"
)

// LuaObject encapsulates a Lua object like a table or a function.
//
// We do not make the type distinction since metatables can make tables callable
// and functions indexable.
type LuaObject struct {
	l   *lua.State
	ref int
}

var (
	ErrLuaObjectCallResults   = errors.New("results must be a pointer to pointer/slice/struct")
	ErrLuaObjectCallable      = errors.New("LuaObject must be callable")
	ErrLuaObjectIndexable     = errors.New("not indexable")
	ErrLuaObjectUnsharedState = errors.New("LuaObjects must share the same state")
)

// NewLuaObject creates a new LuaObject from stack index.
func NewLuaObject(L *lua.State, idx int) *LuaObject {
	L.PushValue(idx)
	ref := L.Ref(lua.LUA_REGISTRYINDEX)
	return &LuaObject{l: L, ref: ref}
}

// NewLuaObjectFromName creates a new LuaObject from the object designated by
// the sequence of 'subfields'.
func NewLuaObjectFromName(L *lua.State, subfields ...interface{}) *LuaObject {
	L.GetGlobal("_G")
	defer L.Pop(1)
	err := get(L, subfields...)
	if err != nil {
		return nil
	}
	val := NewLuaObject(L, -1)
	L.Pop(1)
	return val
}

// NewLuaObjectFromValue creates a new LuaObject from a Go value.
// Note that this will convert any slices or maps into Lua tables.
func NewLuaObjectFromValue(L *lua.State, val interface{}) *LuaObject {
	GoToLua(L, val)
	return NewLuaObject(L, -1)
}

// Call calls a Lua function, given the desired results and the arguments.
// 'results' must be a pointer to a pointer/struct/slice.
//
// - If a pointer, then only the first result is stored to that pointer.
//
// - If a struct with 'n' fields, then the first n results are stored in the field.
//
// - If a slice, then all the results are stored in the slice. The slice is re-allocated if necessary.
//
// If the function returns more values than can be stored in the 'results'
// argument, they will be ignored.
//
// If 'results' is nil, results will be discarded.
func (lo *LuaObject) Call(results interface{}, args ...interface{}) error {
	L := lo.l
	// Push the callable value.
	lo.Push()
	if !L.IsFunction(-1) {
		if !L.GetMetaField(-1, "__call") {
			L.Pop(1)
			return ErrLuaObjectCallable
		}
		// We leave the __call metamethod on stack.
		L.Remove(-2)
	}

	// Push the args.
	for _, arg := range args {
		GoToLuaProxy(L, arg)
	}

	// Special case: discard the results.
	if results == nil {
		err := L.Call(len(args), 0)
		if err != nil {
			L.Pop(1)
			return err
		}
		return nil
	}

	resptr := reflect.ValueOf(results)
	if resptr.Kind() != reflect.Ptr {
		return ErrLuaObjectCallResults
	}
	res := resptr.Elem()

	switch res.Kind() {
	case reflect.Ptr:
		err := L.Call(len(args), 1)
		defer L.Pop(1)
		if err != nil {
			return err
		}
		return LuaToGo(L, -1, res.Interface())

	case reflect.Slice:
		residx := L.GetTop() - len(args)
		err := L.Call(len(args), lua.LUA_MULTRET)
		if err != nil {
			L.Pop(1)
			return err
		}

		nresults := L.GetTop() - residx + 1
		defer L.Pop(nresults)
		t := res.Type()

		// Adjust the length of the slice.
		if res.IsNil() || nresults > res.Len() {
			v := reflect.MakeSlice(t, nresults, nresults)
			res.Set(v)
		} else if nresults < res.Len() {
			res.SetLen(nresults)
		}

		for i := 0; i < nresults; i++ {
			err = LuaToGo(L, residx+i, res.Index(i).Addr().Interface())
			if err != nil {
				return err
			}
		}

	case reflect.Struct:
		exportedFields := []reflect.Value{}
		for i := 0; i < res.NumField(); i++ {
			if res.Field(i).CanInterface() {
				exportedFields = append(exportedFields, res.Field(i).Addr())
			}
		}
		nresults := len(exportedFields)
		err := L.Call(len(args), nresults)
		if err != nil {
			L.Pop(1)
			return err
		}
		defer L.Pop(nresults)
		residx := L.GetTop() - nresults + 1

		for i := 0; i < nresults; i++ {
			err = LuaToGo(L, residx+i, exportedFields[i].Interface())
			if err != nil {
				return err
			}
		}

	default:
		return ErrLuaObjectCallResults
	}

	return nil
}

// Close frees the Lua reference of this object.
func (lo *LuaObject) Close() {
	lo.l.Unref(lua.LUA_REGISTRYINDEX, lo.ref)
}

// get pushes the Lua value indexed at the sequence of 'subfields' from the
// indexable value on top of the stack.
//
// It pushes nothing on error.
//
// Numeric indices start from 1: see Set().
func get(L *lua.State, subfields ...interface{}) error {
	// TODO: See if worth exporting.

	// Duplicate iterable since the following loop removes the last table on stack
	// and we don't want to pop it to be consistent with lua.GetField and
	// lua.GetTable.
	L.PushValue(-1)

	for _, field := range subfields {
		if L.IsTable(-1) {
			GoToLua(L, field)
			L.GetTable(-2)
		} else if L.GetMetaField(-1, "__index") {
			L.PushValue(-2)
			GoToLua(L, field)
			err := L.Call(2, 1)
			if err != nil {
				L.Pop(1)
				return err
			}
		} else {
			return ErrLuaObjectIndexable
		}
		// Remove last iterable.
		L.Remove(-2)
	}
	return nil
}

// Get stores in 'a' the Lua value indexed at the sequence of 'subfields'.
// 'a' must be a pointer as in LuaToGo.
func (lo *LuaObject) Get(a interface{}, subfields ...interface{}) error {
	lo.Push()
	defer lo.l.Pop(1)
	err := get(lo.l, subfields...)
	if err != nil {
		return err
	}
	defer lo.l.Pop(1)
	return LuaToGo(lo.l, -1, a)
}

// GetObject returns the LuaObject indexed at the sequence of 'subfields'.
func (lo *LuaObject) GetObject(subfields ...interface{}) (*LuaObject, error) {
	lo.Push()
	defer lo.l.Pop(1)
	err := get(lo.l, subfields...)
	if err != nil {
		return nil, err
	}
	val := NewLuaObject(lo.l, -1)
	lo.l.Pop(1)
	return val, nil
}

// Push pushes this LuaObject on the stack.
func (lo *LuaObject) Push() {
	lo.l.RawGeti(lua.LUA_REGISTRYINDEX, lo.ref)
}

// Set sets the value at the sequence of 'subfields' with the value 'a'.
// Numeric indices start from 1, as in Lua: if we started from zero, access to
// index 0 or negative indices would be shifted awkwardly.
func (lo *LuaObject) Set(a interface{}, subfields ...interface{}) error {
	parentKeys := subfields[:len(subfields)-1]
	parent, err := lo.GetObject(parentKeys...)
	if err != nil {
		return err
	}

	L := parent.l
	parent.Push()
	defer L.Pop(1)

	lastField := subfields[len(subfields)-1]
	if L.IsTable(-1) {
		GoToLuaProxy(L, lastField)
		GoToLuaProxy(L, a)
		L.SetTable(-3)
	} else if L.GetMetaField(-1, "__newindex") {
		L.PushValue(-2)
		GoToLuaProxy(L, lastField)
		GoToLuaProxy(L, a)
		err := L.Call(3, 0)
		if err != nil {
			L.Pop(1)
			return err
		}
	} else {
		return ErrLuaObjectIndexable
	}
	return nil
}

// Setv copies values between two tables in the same Lua state.
// It overwrites existing values.
func (lo *LuaObject) Setv(src *LuaObject, keys ...string) error {
	// TODO: Rename? This function seems to be too specialized, is it worth
	// keeping at all?
	L := lo.l
	if L != src.l {
		return ErrLuaObjectUnsharedState
	}
	lo.Push()
	defer L.Pop(1)
	loIdx := L.GetTop()

	var set func(int, string)
	if L.IsTable(loIdx) {
		set = L.SetField
	} else if L.GetMetaField(loIdx, "__newindex") {
		L.Pop(1)
		set = func(idx int, key string) {
			resultIdx := L.GetTop()
			L.GetMetaField(loIdx, "__newindex")
			L.PushValue(loIdx)
			L.PushString(key)
			L.PushValue(resultIdx)
			L.Remove(resultIdx)
			L.Call(3, 0)
		}
	} else {
		return ErrLuaObjectIndexable
	}

	src.Push()
	defer src.l.Pop(1)
	srcIdx := L.GetTop()
	var get func(int, string)
	if L.IsTable(srcIdx) {
		get = L.GetField
	} else if L.GetMetaField(srcIdx, "__index") {
		L.Pop(1)
		get = func(idx int, key string) {
			L.GetMetaField(srcIdx, "__index")
			L.PushValue(srcIdx)
			L.PushString(key)
			L.Call(2, 1)
		}
	} else {
		return ErrLuaObjectIndexable
	}

	for _, key := range keys {
		get(srcIdx, key)
		set(loIdx, key)
	}

	return nil
}

// LuaTableIter is the Go equivalent of a Lua table iterator.
type LuaTableIter struct {
	lo *LuaObject
	// keyRef is LUA_NOREF before iteration.
	keyRef int
	// Reference to the iterator in case the metamethod gets changed while
	// iterating.
	iterRef int
	// TODO: See if this is an idiomatic implementation of error storage.
	err error
}

// Error returns the error that happened during last iteration, if any.
func (ti *LuaTableIter) Error() error {
	return ti.err
}

// Iter creates a Lua iterator.
func (lo *LuaObject) Iter() (*LuaTableIter, error) {
	L := lo.l
	lo.Push()
	defer L.Pop(1)
	if L.IsTable(-1) {
		return &LuaTableIter{lo: lo, keyRef: lua.LUA_NOREF, iterRef: lua.LUA_NOREF}, nil
	} else if L.GetMetaField(-1, "__pairs") {
		// __pairs(t) = iterator, t, first-key.
		L.PushValue(-2)
		// Only keep iterator on stack, hence '1' result only.
		err := L.Call(1, 1)
		if err != nil {
			L.Pop(1)
			return nil, err
		}
		ref := L.Ref(lua.LUA_REGISTRYINDEX)
		return &LuaTableIter{lo: lo, keyRef: lua.LUA_NOREF, iterRef: ref}, nil
	} else {
		return nil, ErrLuaObjectIndexable
	}
}

// Next gets the next key/value pair from the indexable value.
//
// 'value' must be a valid argument for LuaToGo. As a special case, 'value' can
// be nil to make it possible to loop over keys without caring about associated
// values.
func (ti *LuaTableIter) Next(key, value interface{}) bool {
	if ti.lo == nil {
		ti.err = errors.New("empty iterator")
		return false
	}
	L := ti.lo.l

	if ti.iterRef == lua.LUA_NOREF {
		// Must be a table. This requires the Iter() function to set
		// ref=LUA_NOREF.

		// Push table.
		ti.lo.Push()
		defer L.Pop(1)

		if ti.keyRef == lua.LUA_NOREF {
			L.PushNil()
		} else {
			L.RawGeti(lua.LUA_REGISTRYINDEX, ti.keyRef)
		}

		if L.Next(-2) == 0 {
			L.Unref(lua.LUA_REGISTRYINDEX, ti.keyRef)
			return false
		}

	} else {
		L.RawGeti(lua.LUA_REGISTRYINDEX, ti.iterRef)
		ti.lo.Push()

		if ti.keyRef == lua.LUA_NOREF {
			L.PushNil()
		} else {
			L.RawGeti(lua.LUA_REGISTRYINDEX, ti.keyRef)
		}

		err := L.Call(2, 2)
		if err != nil {
			L.Pop(1)
			ti.err = err
			return false
		}
		if L.IsNil(-2) {
			L.Pop(2)
			L.Unref(lua.LUA_REGISTRYINDEX, ti.iterRef)
			return false
		}
	}

	err := LuaToGo(L, -2, key)
	if err != nil {
		ti.err = err
		return false
	}
	if value != nil {
		err = LuaToGo(L, -1, value)
		if err != nil {
			ti.err = err
			return false
		}
	}

	// Drop value, key is now on top.
	L.Pop(1)

	// Replace former key reference with new key.
	L.Unref(lua.LUA_REGISTRYINDEX, ti.keyRef)
	ti.keyRef = L.Ref(lua.LUA_REGISTRYINDEX)
	return true
}

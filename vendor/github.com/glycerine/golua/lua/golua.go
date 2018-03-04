package lua

/*
#cgo CFLAGS: -I ${SRCDIR}/../../../LuaJIT/LuaJIT/src
#cgo LDFLAGS: -lm
#cgo linux LDFLAGS: -ldl
#cgo darwin LDFLAGS: -ldl

#include <lua.h>
#include <lualib.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"sync"
	"unsafe"
)

// Type of allocation functions to use with NewStateAlloc
type Alloc func(ptr unsafe.Pointer, osize uint, nsize uint) unsafe.Pointer

// This is the type of go function that can be registered as lua functions
type LuaGoFunction func(L *State) int

// Wrapper to keep cgo from complaining about incomplete ptr type
//export State
type State struct {
	// Wrapped lua_State object
	s *C.lua_State

	// index of this object inside the goStates array
	Index uintptr

	Shared *SharedByAllCoroutines

	IsMainCoro bool // if true, then will be registered

	MainCo  *State       // always points to the main coro
	CmainCo *C.lua_State // ditto

	uPos int // position in uniqArray
}

type SharedByAllCoroutines struct {
	// Registry of go object that have been pushed to Lua VM
	registry []interface{}

	// Freelist for funcs indices, to allow for freeing
	freeIndices []uint
}

func newSharedByAllCoroutines() *SharedByAllCoroutines {
	return &SharedByAllCoroutines{
		registry:    make([]interface{}, 0, 8),
		freeIndices: make([]uint, 0, 8),
	}
}

var goStates map[uintptr]*State
var goStatesMutex sync.Mutex

func init() {
	goStates = make(map[uintptr]*State, 16)
}

func registerGoState(L *State) {
	goStatesMutex.Lock()
	defer goStatesMutex.Unlock()
	L.Index = uintptr(unsafe.Pointer(L))
	goStates[L.Index] = L
}

func unregisterGoState(L *State) {
	goStatesMutex.Lock()
	defer goStatesMutex.Unlock()
	delete(goStates, L.Index)
}

func getGoState(gostateindex uintptr) *State {
	goStatesMutex.Lock()
	defer goStatesMutex.Unlock()
	return goStates[gostateindex]
}

/* jea added to debug
func printGoStates() (biggest uintptr) {
	goStatesMutex.Lock()
	defer goStatesMutex.Unlock()
	maxReg := 0
	for gostateindex, v := range goStates {
		fmt.Printf("gostateindex=%v, state='%#v'\n", gostateindex, v)
		cur := len(v.registry)
		if cur > maxReg {
			maxReg = cur
			biggest = gostateindex
		}
	}
	return biggest
}
*/

//export golua_callgofunction
func golua_callgofunction(curThread *C.lua_State, gostateindex uintptr, mainIndex uintptr, fid uint) int {

	fmt.Printf("jea debug: golua_callgofunction, here is stack of curThread at top:\n")
	DumpLuaStack(&State{s: curThread})

	fmt.Printf("jea debug: golua_callgofunction or __call on userdata top: gostateindex='%#v', curThread is '%p'/'%#v'\nstack:\n%s\n", gostateindex, curThread, curThread, string(debug.Stack()))

	var L1 *State
	if gostateindex == 0 {
		// lua side created goroutine, first time seen;
		// and not yet registered on the go-side.
		fmt.Printf("debug: first time this coroutine has been seen on the Go side\n")

		L := getGoState(mainIndex)
		L1 = L.ToThreadHelper(curThread)
	} else {

		// this is the __call() for the MT_GOFUNCTION
		L1 = getGoState(gostateindex)
	}

	fmt.Printf("L1 corresponding to gostateindex '%v' -> '%#v'\n", gostateindex, L1)
	//biggest := printGoStates()
	// jea: to debug, try substituting the biggest...
	//if biggest != gostateindex {
	//	L1 = getGoState(biggest)
	//}
	if fid < 0 {
		panic(&LuaError{0, "Requested execution of an unknown function", L1.StackTrace()})
	}
	f := L1.Shared.registry[fid].(LuaGoFunction)
	fmt.Printf("\n jea debug golua_callgofunction: f back from registry for fid=%#v, is f=%#v\n", fid, f)

	fmt.Printf("\n jea debug: in golua_callgofunction(): L1 stack is:\n")
	DumpLuaStack(L1)

	fmt.Printf("\n jea debug, in golua_callgofunction(): right before final f(L1) call.\n")
	return f(L1)
}

var typeOfBytes = reflect.TypeOf([]byte(nil))

//export golua_interface_newindex_callback
func golua_interface_newindex_callback(gostateindex uintptr, iid uint, field_name_cstr *C.char) int {
	L := getGoState(gostateindex)
	iface := L.Shared.registry[iid]
	ifacevalue := reflect.ValueOf(iface).Elem()

	field_name := C.GoString(field_name_cstr)

	fval := ifacevalue.FieldByName(field_name)

	if fval.Kind() == reflect.Ptr {
		fval = fval.Elem()
	}

	luatype := LuaValType(C.lua_type(L.s, 3))

	switch fval.Kind() {
	case reflect.Bool:
		if luatype == LUA_TBOOLEAN {
			fval.SetBool(int(C.lua_toboolean(L.s, 3)) != 0)
			return 1
		} else {
			L.PushString("Wrong assignment to field " + field_name)
			return -1
		}

	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		if luatype == LUA_TNUMBER {
			fval.SetInt(int64(C.lua_tointeger(L.s, 3)))
			return 1
		} else {
			L.PushString("Wrong assignment to field " + field_name)
			return -1
		}

	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		if luatype == LUA_TNUMBER {
			fval.SetUint(uint64(C.lua_tointeger(L.s, 3)))
			return 1
		} else {
			L.PushString("Wrong assignment to field " + field_name)
			return -1
		}

	case reflect.String:
		if luatype == LUA_TSTRING {
			fval.SetString(C.GoString(C.lua_tolstring(L.s, 3, nil)))
			return 1
		} else {
			L.PushString("Wrong assignment to field " + field_name)
			return -1
		}

	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		if luatype == LUA_TNUMBER {
			fval.SetFloat(float64(C.lua_tonumber(L.s, 3)))
			return 1
		} else {
			L.PushString("Wrong assignment to field " + field_name)
			return -1
		}
	case reflect.Slice:
		if fval.Type() == typeOfBytes {
			if luatype == LUA_TSTRING {
				fval.SetBytes(L.ToBytes(3))
				return 1
			} else {
				L.PushString("Wrong assignment to field " + field_name)
				return -1
			}
		}
	}

	L.PushString("Unsupported type of field " + field_name + ": " + fval.Type().String())
	return -1
}

//export golua_interface_index_callback
func golua_interface_index_callback(gostateindex uintptr, iid uint, field_name *C.char) int {
	L := getGoState(gostateindex)
	iface := L.Shared.registry[iid]
	ifacevalue := reflect.ValueOf(iface).Elem()

	fval := ifacevalue.FieldByName(C.GoString(field_name))

	if fval.Kind() == reflect.Ptr {
		fval = fval.Elem()
	}

	switch fval.Kind() {
	case reflect.Bool:
		L.PushBoolean(fval.Bool())
		return 1

	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		L.PushInteger(fval.Int())
		return 1

	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		L.PushInteger(int64(fval.Uint()))
		return 1

	case reflect.String:
		L.PushString(fval.String())
		return 1

	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		L.PushNumber(fval.Float())
		return 1
	case reflect.Slice:
		if fval.Type() == typeOfBytes {
			L.PushBytes(fval.Bytes())
			return 1
		}
	}

	L.PushString("Unsupported type of field: " + fval.Type().String())
	return -1
}

//export golua_gchook
func golua_gchook(gostateindex uintptr, id uint) int {
	L1 := getGoState(gostateindex)
	L1.unregister(id)
	return 0
}

//export golua_callpanicfunction
func golua_callpanicfunction(gostateindex uintptr, id uint) int {
	L1 := getGoState(gostateindex)
	f := L1.Shared.registry[id].(LuaGoFunction)
	return f(L1)
}

//export golua_idtointerface
func golua_idtointerface(id uint) interface{} {
	return id
}

//export golua_cfunctiontointerface
func golua_cfunctiontointerface(f *uintptr) interface{} {
	return f
}

//export golua_callallocf
func golua_callallocf(fp uintptr, ptr uintptr, osize uint, nsize uint) uintptr {
	return uintptr((*((*Alloc)(unsafe.Pointer(fp))))(unsafe.Pointer(ptr), osize, nsize))
}

//export go_panic_msghandler
func go_panic_msghandler(gostateindex uintptr, z *C.char) {
	L := getGoState(gostateindex)
	s := C.GoString(z)

	panic(&LuaError{LUA_ERRERR, s, L.StackTrace()})
}

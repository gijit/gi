// This package provides access to the excellent lua language interpreter from go code.
//
// Access to most of the functions in lua.h and lauxlib.h is provided as well as additional convenience functions to publish Go objects and functions to lua code.
//
// The documentation of this package is no substitute for the official lua documentation and in many instances methods are described only with the name of their C equivalent
package lua

/*
#cgo CFLAGS: -I ${SRCDIR}/../../../LuaJIT/LuaJIT/src
#cgo LDFLAGS: ${SRCDIR}/../../../LuaJIT/LuaJIT/src/libluajit.a -lm
#cgo (linux OR darwin) LDFLAGS: -ldl

#include <lua.h>
#include <stdlib.h>

#include "golua.h"

*/
import "C"

import (
	"unsafe"

	"fmt"
	"github.com/gijit/gi/pkg/verb"
)

var pp = verb.PP

type LuaStackEntry struct {
	Name        string
	Source      string
	ShortSource string
	CurrentLine int
}

func newState(L *C.lua_State) *State {
	newstate := &State{
		s:          L,
		Shared:     newSharedByAllCoroutines(),
		IsMainCoro: true,
		CmainCo:    L,
		AllCoro:    make(map[int]*State),
	}
	newstate.MainCo = newstate

	// only for main states, not additional coroutines:
	registerGoState(newstate) // sets Index
	newstate.Upos = int(C.clua_setgostate(L, C.size_t(newstate.Index)))

	// assert(Upos == 1)
	if newstate.Upos != 1 {
		panic(fmt.Sprintf("assert violated: we expected newstate.Upos for the main coro to always be at index 1: our code depends on that!"))
	}
	newstate.MainCo.AllCoro[newstate.Upos] = newstate
	C.clua_initstate(L)
	return newstate
}

func (L *State) getFreeIndex() (index uint, ok bool) {
	freelen := len(L.Shared.freeIndices)
	//if there exist entries in the freelist
	if freelen > 0 {
		i := L.Shared.freeIndices[freelen-1] //get index
		//fmt.Printf("Free indices before: %v\n", L.Shared.freeIndices)
		L.Shared.freeIndices = L.Shared.freeIndices[0 : freelen-1] //'pop' index from list
		//fmt.Printf("Free indices after: %v\n", L.freeIndices)
		return i, true
	}
	return 0, false
}

//returns the registered function id
func (L *State) register(f interface{}) uint {
	//fmt.Printf("Registering %v\n")
	index, ok := L.getFreeIndex()
	//fmt.Printf("\tfreeindex: index = %v, ok = %v\n", index, ok)
	if ok {
		// reuse
		L.Shared.registry[index] = f
		return index
	}
	// add a new index
	L.Shared.registry = append(L.Shared.registry, f)
	index = uint(len(L.Shared.registry)) - 1
	return index
}

func (L *State) unregister(fid uint) {
	//fmt.Printf("Unregistering %d (len: %d, value: %v)\n", fid, len(*L.registry), (*L.registry)[fid])
	if (fid < uint(len(L.Shared.registry))) && (L.Shared.registry[fid] != nil) {
		L.Shared.registry[fid] = nil
		L.Shared.freeIndices = append(L.Shared.freeIndices, fid)
	}
}

// Like lua_pushcfunction pushes onto the stack a go function as user data
func (L *State) PushGoFunction(f LuaGoFunction) {
	fid := L.register(f)
	C.clua_pushgofunction(L.s, C.uint(fid))
}

// PushGoClosure pushes a lua.LuaGoFunction to the stack wrapped in a Closure.
// this permits the go function to reflect lua type 'function' when checking with type()
// this implements behaviour akin to lua_pushcfunction() in lua C API.
func (L *State) PushGoClosure(f LuaGoFunction) {
	L.PushGoFunction(f)      // leaves Go function userdata on stack
	C.clua_pushcallback(L.s) // wraps the userdata object with a closure making it into a function
}

// Sets a metamethod to execute a go function
//
// The code:
//
// 	L.LGetMetaTable(tableName)
// 	L.SetMetaMethod(methodName, function)
//
// is the logical equivalent of:
//
// 	L.LGetMetaTable(tableName)
// 	L.PushGoFunction(function)
// 	L.SetField(-2, methodName)
//
// except this wouldn't work because pushing a go function results in user data not a cfunction
func (L *State) SetMetaMethod(methodName string, f LuaGoFunction) {
	L.PushGoFunction(f)      // leaves Go function userdata on stack
	C.clua_pushcallback(L.s) // wraps the userdata object with a closure making it into a function
	L.SetField(-2, methodName)
}

// Pushes a Go struct onto the stack as user data.
//
// The user data will be rigged so that lua code can access and change the public members of simple types directly
func (L *State) PushGoStruct(iface interface{}) {
	iid := L.register(iface)
	C.clua_pushgostruct(L.s, C.uint(iid))
}

// Push a pointer onto the stack as user data.
//
// This function doesn't save a reference to the interface, it is the responsibility of the caller of this function to insure that the interface outlasts the lifetime of the lua object that this function creates.
func (L *State) PushLightUserdata(ud *interface{}) {
	//push
	C.lua_pushlightuserdata(L.s, unsafe.Pointer(ud))
}

// Creates a new user data object of specified size and returns it
func (L *State) NewUserdata(size uintptr) unsafe.Pointer {
	return unsafe.Pointer(C.lua_newuserdata(L.s, C.size_t(size)))
}

// Sets the AtPanic function, returns the old one
//
// BUG(everyone_involved): passing nil causes serious problems
func (L *State) AtPanic(panicf LuaGoFunction) (oldpanicf LuaGoFunction) {
	fid := uint(0)
	if panicf != nil {
		fid = L.register(panicf)
	}
	oldres := interface{}(C.clua_atpanic(L.s, C.uint(fid)))
	switch i := oldres.(type) {
	case C.uint:
		f := L.Shared.registry[uint(i)].(LuaGoFunction)
		//free registry entry
		L.unregister(uint(i))
		return f
	case C.lua_CFunction:
		return func(L1 *State) int {
			return int(C.clua_callluacfunc(L1.s, i))
		}
	}
	//generally we only get here if the panicf got set to something like nil
	//potentially dangerous because we may silently fail
	return nil
}

func (L *State) pcall(nargs, nresults, errfunc int) int {
	return int(C.lua_pcall(L.s, C.int(nargs), C.int(nresults), C.int(errfunc)))
}

func (L *State) callEx(nargs, nresults int, catch bool) (err error) {
	if catch {
		defer func() {
			if err2 := recover(); err2 != nil {
				if _, ok := err2.(error); ok {
					err = err2.(error)
				}
				return
			}
		}()
	}

	// switch to coro here?
	/*
		thread := L.ToThread(1)
		if thread != nil {
			narg := L.GetTop()

			fmt.Printf("callEx: L stack just before XMove:")
			DumpLuaStack(L)
			fmt.Printf("callEx: thread stack just before XMove:")
			DumpLuaStack(thread)

			XMove(L, thread, narg-1)

			fmt.Printf("callEx: L stack just after XMove:")
			DumpLuaStack(L)
			fmt.Printf("callEx: thread stack just after XMove:")
			DumpLuaStack(thread)

			fmt.Printf("callEx(): thread was first arg, replacing L with running coro.\n")
			L = thread
		}
	*/

	L.GetGlobal(C.GOLUA_DEFAULT_MSGHANDLER)
	// We must record where we put the error handler in the stack otherwise it will be impossible to remove after the pcall when nresults == LUA_MULTRET
	erridx := L.GetTop() - nargs - 1
	L.Insert(erridx)
	pp("callEx: golua lua.go, stack just before pcall, L = '%p'/'%#v'\n", L, L)
	//DumpLuaStack(L)
	r := L.pcall(nargs, nresults, erridx)
	pp("callEx: golua lua.go, stack just after pcall:")
	if verb.VerboseVerbose {
		DumpLuaStack(L)
	}
	L.Remove(erridx)
	if r != 0 {
		err = &LuaError{r, L.ToString(-1), L.StackTrace()}
		if !catch {
			panic(err)
		}
	}
	return
}

// lua_call
func (L *State) Call(nargs, nresults int) (err error) {
	return L.callEx(nargs, nresults, true)
}

// Like lua_call but panics on errors
func (L *State) MustCall(nargs, nresults int) {
	L.callEx(nargs, nresults, false)
}

// lua_checkstack
func (L *State) CheckStack(extra int) bool {
	return C.lua_checkstack(L.s, C.int(extra)) != 0
}

// lua_close
func (L *State) Close() {
	C.lua_close(L.s)
	unregisterGoState(L)
}

// lua_concat
func (L *State) Concat(n int) {
	C.lua_concat(L.s, C.int(n))
}

// lua_createtable
func (L *State) CreateTable(narr int, nrec int) {
	C.lua_createtable(L.s, C.int(narr), C.int(nrec))
}

// lua_equal
func (L *State) Equal(index1, index2 int) bool {
	return C.lua_equal(L.s, C.int(index1), C.int(index2)) == 1
}

// lua_gc
func (L *State) GC(what, data int) int { return int(C.lua_gc(L.s, C.int(what), C.int(data))) }

// lua_getfenv
func (L *State) GetfEnv(index int) { C.lua_getfenv(L.s, C.int(index)) }

// lua_getfield
func (L *State) GetField(index int, k string) {
	Ck := C.CString(k)
	defer C.free(unsafe.Pointer(Ck))
	C.lua_getfield(L.s, C.int(index), Ck)
}

// Pushes on the stack the value of a global variable (lua_getglobal)
func (L *State) GetGlobal(name string) { L.GetField(LUA_GLOBALSINDEX, name) }

// lua_getmetatable
func (L *State) GetMetaTable(index int) bool {
	return C.lua_getmetatable(L.s, C.int(index)) != 0
}

// lua_gettable
func (L *State) GetTable(index int) { C.lua_gettable(L.s, C.int(index)) }

// lua_gettop
func (L *State) GetTop() int { return int(C.lua_gettop(L.s)) }

// lua_insert
func (L *State) Insert(index int) { C.lua_insert(L.s, C.int(index)) }

// Returns true if lua_type == LUA_TBOOLEAN
func (L *State) IsBoolean(index int) bool {
	return LuaValType(C.lua_type(L.s, C.int(index))) == LUA_TBOOLEAN
}

// Returns true if the value at index is a LuaGoFunction
func (L *State) IsGoFunction(index int) bool {
	return C.clua_isgofunction(L.s, C.int(index)) != 0
}

// Returns true if the value at index is user data pushed with PushGoStruct
func (L *State) IsGoStruct(index int) bool {
	return C.clua_isgostruct(L.s, C.int(index)) != 0
}

// Returns true if the value at index is user data pushed with PushGoFunction
func (L *State) IsFunction(index int) bool {
	return LuaValType(C.lua_type(L.s, C.int(index))) == LUA_TFUNCTION
}

// Returns true if the value at index is light user data
func (L *State) IsLightUserdata(index int) bool {
	return LuaValType(C.lua_type(L.s, C.int(index))) == LUA_TLIGHTUSERDATA
}

// lua_isnil
func (L *State) IsNil(index int) bool { return LuaValType(C.lua_type(L.s, C.int(index))) == LUA_TNIL }

// lua_isnone
func (L *State) IsNone(index int) bool { return LuaValType(C.lua_type(L.s, C.int(index))) == LUA_TNONE }

// lua_isnoneornil
func (L *State) IsNoneOrNil(index int) bool { return int(C.lua_type(L.s, C.int(index))) <= 0 }

// lua_isnumber
func (L *State) IsNumber(index int) bool { return C.lua_isnumber(L.s, C.int(index)) == 1 }

// lua_isstring
func (L *State) IsString(index int) bool { return C.lua_isstring(L.s, C.int(index)) == 1 }

// lua_istable
func (L *State) IsTable(index int) bool {
	return LuaValType(C.lua_type(L.s, C.int(index))) == LUA_TTABLE
}

// lua_isthread
func (L *State) IsThread(index int) bool {
	return LuaValType(C.lua_type(L.s, C.int(index))) == LUA_TTHREAD
}

// lua_isuserdata
func (L *State) IsUserdata(index int) bool { return C.lua_isuserdata(L.s, C.int(index)) == 1 }

// lua_lessthan
func (L *State) LessThan(index1, index2 int) bool {
	return C.lua_lessthan(L.s, C.int(index1), C.int(index2)) == 1
}

// Creates a new lua interpreter state with the given allocation function
func NewStateAlloc(f Alloc) *State {
	// jea: ../example/alloc.go panics... hmm.
	ls := C.clua_newstate(unsafe.Pointer(&f))
	return newState(ls)
}

// lua_newtable
func (L *State) NewTable() {
	C.lua_createtable(L.s, 0, 0)
}

// lua_newthread
func (L *State) NewThread() *State {

	s := C.lua_newthread(L.s)
	return L.ToThreadHelper(s)
}

// lua_next
func (L *State) Next(index int) int {
	return int(C.lua_next(L.s, C.int(index)))
}

// lua_objlen
//
// Returns the "length" of the value at the
// given acceptable index: for strings, this
// is the string length; for tables, this
// is the result of the length operator ('#');
// for userdata, this is the size of the block
// of memory allocated for the userdata;
// for other values, it is 0.
// jea note: Despite the misleading description
// above, in 5.1 or LuaJit, ObjLen does not call
// the metamethod __len(). In 5.2 it was split
// into len and rawlen, with len calling the __len
// metamethod.
//
func (L *State) ObjLen(index int) uint {
	return uint(C.lua_objlen(L.s, C.int(index)))
}

// lua_pop
func (L *State) Pop(n int) {
	//Why is this implemented this way? I don't get it... maybe it
	// is just inlining manually the actual implementation.
	//C.lua_pop(L.s, C.int(n));
	C.lua_settop(L.s, C.int(-n-1))
}

// lua_pushboolean
func (L *State) PushBoolean(b bool) {
	var bint int
	if b {
		bint = 1
	} else {
		bint = 0
	}
	C.lua_pushboolean(L.s, C.int(bint))
}

// lua_pushstring
func (L *State) PushString(str string) {
	Cstr := C.CString(str)
	defer C.free(unsafe.Pointer(Cstr))
	C.lua_pushlstring(L.s, Cstr, C.size_t(len(str)))
}

func (L *State) PushBytes(b []byte) {
	C.lua_pushlstring(L.s, (*C.char)(unsafe.Pointer(&b[0])), C.size_t(len(b)))
}

// lua_pushinteger
func (L *State) PushInteger(n int64) {
	C.lua_pushinteger(L.s, C.lua_Integer(n))
}

// lua_pushnil
func (L *State) PushNil() {
	C.lua_pushnil(L.s)
}

// lua_pushnumber
func (L *State) PushNumber(n float64) {
	C.lua_pushnumber(L.s, C.lua_Number(n)) // lua_Number is a cast
}

// lua_pushthread
func (L *State) PushThread() (isMain bool) {
	return C.lua_pushthread(L.s) != 0
}

// lua_pushvalue
func (L *State) PushValue(index int) {
	C.lua_pushvalue(L.s, C.int(index))
}

// lua_rawequal
func (L *State) RawEqual(index1 int, index2 int) bool {
	return C.lua_rawequal(L.s, C.int(index1), C.int(index2)) != 0
}

// lua_rawget
func (L *State) RawGet(index int) {
	C.lua_rawget(L.s, C.int(index))
}

// lua_rawgeti
func (L *State) RawGeti(index int, n int) {
	C.lua_rawgeti(L.s, C.int(index), C.int(n))
}

// lua_rawset
func (L *State) RawSet(index int) {
	C.lua_rawset(L.s, C.int(index))
}

// lua_rawseti
func (L *State) RawSeti(index int, n int) {
	C.lua_rawseti(L.s, C.int(index), C.int(n))
}

// Registers a Go function as a global variable
func (L *State) Register(name string, f LuaGoFunction) {
	L.PushGoFunction(f)
	L.SetGlobal(name)
}

// lua_remove
func (L *State) Remove(index int) {
	C.lua_remove(L.s, C.int(index))
}

// lua_replace
func (L *State) Replace(index int) {
	C.lua_replace(L.s, C.int(index))
}

// lua_resume
func (L *State) Resume(narg int) int {
	return int(C.lua_resume(L.s, C.int(narg)))
}

// lua_setallocf
func (L *State) SetAllocf(f Alloc) {
	C.clua_setallocf(L.s, unsafe.Pointer(&f))
}

// lua_setfenv
func (L *State) SetfEnv(index int) {
	C.lua_setfenv(L.s, C.int(index))
}

// lua_setfield
func (L *State) SetField(index int, k string) {
	Ck := C.CString(k)
	defer C.free(unsafe.Pointer(Ck))
	C.lua_setfield(L.s, C.int(index), Ck)
}

// lua_setglobal
func (L *State) SetGlobal(name string) {
	Cname := C.CString(name)
	defer C.free(unsafe.Pointer(Cname))
	C.lua_setfield(L.s, C.int(LUA_GLOBALSINDEX), Cname)
}

// lua_setmetatable
func (L *State) SetMetaTable(index int) {
	C.lua_setmetatable(L.s, C.int(index))
}

// lua_settable
func (L *State) SetTable(index int) {
	C.lua_settable(L.s, C.int(index))
}

// lua_settop
func (L *State) SetTop(index int) {
	C.lua_settop(L.s, C.int(index))
}

// lua_status
func (L *State) Status() int {
	return int(C.lua_status(L.s))
}

// lua_toboolean
func (L *State) ToBoolean(index int) bool {
	return C.lua_toboolean(L.s, C.int(index)) != 0
}

// Returns the value at index as a Go function (it must be something pushed with PushGoFunction)
func (L *State) ToGoFunction(index int) (f LuaGoFunction) {
	if !L.IsGoFunction(index) {
		return nil
	}
	fid := C.clua_togofunction(L.s, C.int(index))
	if fid < 0 {
		return nil
	}
	return L.Shared.registry[fid].(LuaGoFunction)
}

// Returns the value at index as a Go Struct (it must be something pushed with PushGoStruct)
func (L *State) ToGoStruct(index int) (f interface{}) {
	if !L.IsGoStruct(index) {
		return nil
	}
	fid := C.clua_togostruct(L.s, C.int(index))
	if fid < 0 {
		return nil
	}
	pp("querying fid = '%v'\n", fid)
	return L.Shared.registry[fid]
}

// lua_tostring
func (L *State) ToString(index int) string {
	var size C.size_t
	r := C.lua_tolstring(L.s, C.int(index), &size)
	return C.GoStringN(r, C.int(size))
}

func (L *State) ToBytes(index int) []byte {
	var size C.size_t
	b := C.lua_tolstring(L.s, C.int(index), &size)
	return C.GoBytes(unsafe.Pointer(b), C.int(size))
}

// lua_tointeger
func (L *State) ToInteger(index int) int {
	return int(C.lua_tointeger(L.s, C.int(index)))
}

// lua_tonumber
func (L *State) ToNumber(index int) float64 {
	return float64(C.lua_tonumber(L.s, C.int(index)))
}

func (L *State) CdataToInt64(index int) int64 {
	return int64(C.lua_cdata_to_int64(L.s, C.int(index)))
}

func (L *State) CdataToInt32(index int) int32 {
	return int32(C.lua_cdata_to_int32(L.s, C.int(index)))
}

func (L *State) CdataToUint64(index int) uint64 {
	return uint64(C.lua_cdata_to_uint64(L.s, C.int(index)))
}

// lua_topointer
func (L *State) ToPointer(index int) uintptr {
	return uintptr(C.lua_topointer(L.s, C.int(index)))
}

// lua_tothread
func (L *State) ToThread(index int) *State {
	ptr := (*C.lua_State)(unsafe.Pointer(C.lua_tothread(L.s, C.int(index))))
	if ptr == nil {
		return nil
	}
	return L.ToThreadHelper(ptr)
}

func (L *State) ToThreadHelper(ptr *C.lua_State) *State {
	if ptr == nil {
		return nil
	}
	upos := int(C.clua_dedup_coro(ptr))
	already := L.MainCo.AllCoro[upos]
	if already != nil {
		fmt.Printf("ToThreadHelper, already known at upos=%v\n", upos)
		return already
	}
	fmt.Printf("ToThreadHelper, new at upos=%v\n", upos)

	newstate := &State{
		s:       ptr,
		Shared:  L.Shared,
		MainCo:  L.MainCo,
		CmainCo: L.MainCo.s,
		Index:   -1, // not the main state/main thread.
		Upos:    upos,
	}
	// don't register non-main threads in gostates[].

	// asserts that (Upos != 1)
	if newstate.Upos == 1 {
		panic(fmt.Sprintf("assert violated: we expected newstate.Upos to not be 1 for any non-main thread/coroutine stated! our code in c-golua.c depends on that"))

	}
	if newstate.Upos == -1 {
		panic(fmt.Sprintf("assert violated: we expected newstate.Upos to not be -1 for any not known coroutine!"))

	}
	newstate.MainCo.AllCoro[newstate.Upos] = newstate
	return newstate
}

// lua_touserdata
func (L *State) ToUserdata(index int) unsafe.Pointer {
	return unsafe.Pointer(C.lua_touserdata(L.s, C.int(index)))
}

// lua_type
func (L *State) Type(index int) LuaValType {
	return LuaValType(C.lua_type(L.s, C.int(index)))
}

// lua_typename
func (L *State) Typename(tp int) string {
	return C.GoString(C.lua_typename(L.s, C.int(tp)))
}

// lua_xmove
func XMove(from *State, to *State, n int) {
	C.lua_xmove(from.s, to.s, C.int(n))
}

// lua_yield
func (L *State) Yield(nresults int) int {
	return int(C.lua_yield(L.s, C.int(nresults)))
}

// Restricted library opens

// Calls luaopen_base
func (L *State) OpenBase() {
	C.clua_openbase(L.s)
}

// Calls luaopen_io
func (L *State) OpenIO() {
	C.clua_openio(L.s)
}

// Calls luaopen_math
func (L *State) OpenMath() {
	C.clua_openmath(L.s)
}

// Calls luaopen_package
func (L *State) OpenPackage() {
	C.clua_openpackage(L.s)
}

// Calls luaopen_string
func (L *State) OpenString() {
	C.clua_openstring(L.s)
}

// Calls luaopen_table
func (L *State) OpenTable() {
	C.clua_opentable(L.s)
}

// Calls luaopen_os
func (L *State) OpenOS() {
	C.clua_openos(L.s)
}

// Sets the maximum number of operations to execute at instrNumber, after this the execution ends
func (L *State) SetExecutionLimit(instrNumber int) {
	C.clua_setexecutionlimit(L.s, C.int(instrNumber))
}

// Returns the current stack trace
func (L *State) StackTrace() []LuaStackEntry {
	r := []LuaStackEntry{}
	var d C.lua_Debug
	Sln := C.CString("Sln")
	defer C.free(unsafe.Pointer(Sln))

	for depth := 0; C.lua_getstack(L.s, C.int(depth), &d) > 0; depth++ {
		C.lua_getinfo(L.s, Sln, &d)
		ssb := make([]byte, C.LUA_IDSIZE)
		for i := 0; i < C.LUA_IDSIZE; i++ {
			ssb[i] = byte(d.short_src[i])
			if ssb[i] == 0 {
				ssb = ssb[:i]
				break
			}
		}
		ss := string(ssb)

		r = append(r, LuaStackEntry{C.GoString(d.name), C.GoString(d.source), ss, int(d.currentline)})
	}

	return r
}

func (L *State) RaiseError(msg string) {
	st := L.StackTrace()
	prefix := ""
	if len(st) >= 1 {
		prefix = fmt.Sprintf("%s:%d: ", st[1].ShortSource, st[1].CurrentLine)
	}
	panic(&LuaError{0, prefix + msg, st})
}

func (L *State) NewError(msg string) *LuaError {
	return &LuaError{0, msg, L.StackTrace()}
}

// LuaJIT only: return ctype of the cdata at the top of the stack.
func (L *State) LuaJITctypeID(idx int) uint32 {
	res := C.clua_luajit_ctypeid(L.s, C.int(idx))
	return uint32(res)
}

func (L *State) PushInt64(n int64) {
	C.clua_luajit_push_cdata_int64(L.s, C.int64_t(n))
}

func (L *State) PushUint64(u uint64) {
	C.clua_luajit_push_cdata_uint64(L.s, C.uint64_t(u))
}

// jea
func DumpLuaStack(L *State) {
	fmt.Printf("\n%s\n", DumpLuaStackAsString(L))
}

func DumpLuaStackAsString(L *State) (s string) {
	var top int

	top = L.GetTop()
	isMain := L.PushThread()
	thr := L.ToThread(-1)
	L.SetTop(top)
	s += fmt.Sprintf("========== begin DumpLuaStack (of coro %p/lua.State=%p; isMain=%v): top = %v\n", thr, thr.s, isMain, top)
	for i := top; i >= 1; i-- {

		t := L.Type(i)
		s += fmt.Sprintf("DumpLuaStack: i=%v, t= %v\n", i, t)
		s += LuaStackPosToString(L, i)
	}
	s += fmt.Sprintf("========= end of DumpLuaStack\n")
	return
}

func LuaStackPosToString(L *State, i int) string {
	t := L.Type(i)

	switch t {
	case LUA_TNONE: // -1
		return fmt.Sprintf("LUA_TNONE; i=%v was invalid index\n", i)
	case LUA_TNIL:
		return fmt.Sprintf("LUA_TNIL: nil\n")
	case LUA_TSTRING:
		return fmt.Sprintf(" String : \t%v\n", L.ToString(i))
	case LUA_TBOOLEAN:
		return fmt.Sprintf(" Bool : \t\t%v\n", L.ToBoolean(i))
	case LUA_TNUMBER:
		return fmt.Sprintf(" Number : \t%v\n", L.ToNumber(i))
	case LUA_TTABLE:
		return fmt.Sprintf(" Table : \n%s\n", dumpTableString(L, i))

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
			return fmt.Sprintf(" int64: '%v'\n", val)
		case 12: //  uint64
			val := L.CdataToUint64(i)
			return fmt.Sprintf(" uint64: '%v'\n", val)
		case 13: //  float32
		case 14: //  float64

		case 0: // means it wasn't a ctype
		}

	case LUA_TUSERDATA:
		return fmt.Sprintf(" Type(code %v/ LUA_TUSERDATA) : 0x%x with pointer %p\n", t, L.ToPointer(i), L.ToUserdata(i))
	case LUA_TFUNCTION:
		return fmt.Sprintf(" Type(code %v/ LUA_TFUNCTION) : 0x%x\n", t, L.ToPointer(i))
	case LUA_TTHREAD:
		return fmt.Sprintf(" Type(code %v/ LUA_TTHREAD) : 0x%x\n", t, L.ToPointer(i))
	case LUA_TLIGHTUSERDATA:
		return fmt.Sprintf(" Type(code %v/ LUA_TLIGHTUSERDATA) : 0x%x with pointer %p\n", t, L.ToPointer(i), L.ToUserdata(i))
	default:
	}
	return fmt.Sprintf(" Type(code %v) : 0x%x, no auto-print available.\n", t, L.ToPointer(i))
}

func dumpTableString(L *State, index int) (s string) {

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

package lua

import (
	"fmt"
	"testing"
	"unsafe"
)

type TestStruct struct {
	IntField    int
	StringField string
	FloatField  float64
}

func TestGoStruct(t *testing.T) {
	L := NewState()
	L.OpenLibs()
	defer L.Close()

	ts := &TestStruct{10, "test", 2.3}

	L.CheckStack(1)

	L.PushGoStruct(ts)
	L.SetGlobal("t")

	L.GetGlobal("t")
	if !L.IsGoStruct(-1) {
		t.Fatal("Not go struct")
	}

	tsr := L.ToGoStruct(-1).(*TestStruct)
	if tsr != ts {
		t.Fatal("Retrieved something different from what we inserted")
	}

	L.Pop(1)

	L.PushString("This is not a struct")
	if L.ToGoStruct(-1) != nil {
		t.Fatal("Non-GoStruct value attempted to convert into GoStruct should result in nil")
	}

	L.Pop(1)
}

func TestCheckStringSuccess(t *testing.T) {
	L := NewState()
	L.OpenLibs()
	defer L.Close()

	Test := func(L *State) int {
		L.PushString("this is a test")
		L.CheckString(-1)
		return 0
	}

	L.Register("test", Test)
	err := L.DoString("test()")
	if err != nil {
		t.Fatalf("DoString did return an error: %v\n", err.Error())
	}
}

/*
func TestCheckStringFail(t *testing.T) {
	L := NewState()
	L.OpenLibs()
	defer L.Close()

	Test := func(L *State) int {
		L.CheckString(-1)
		return 0
	}

	L.Register("test", Test)
	err := L.DoString("test();")
	if err == nil {
		t.Fatal("DoString did not return an error\n")
	}
}
*/

// jea: works under OpenLibs, because generally
//      we will need pcall/xpcall.
//
// See https://github.com/aarzilli/golua#on-error-handling
//   for why they are hidden. We should probably
//   be hiding these then using unsafe_pcall, and
//   verify that Lua code using them doesn't call
//   back into Go code that panics.
/*
Shortened version of the on-error-handling link:

Lua's exceptions are incompatible with Go.
golua works around this incompatibility by
setting up protected execution environments
in lua.State.DoString, lua.State.DoFile, and
lua.State.Call and turning every exception
into a Go panic.

This means that:

In general you can't do any exception handling
from Lua, pcall and xpcall are renamed to
unsafe_pcall and unsafe_xpcall. They are only
safe to be called from Lua code that never
calls back to Go. Use at your own risk.
*/
func TestPCallHidden(t *testing.T) {
	L := NewState()
	//L.OpenLibs()
	L.OpenBase()
	defer L.Close()

	err := L.DoString("pcall(print, \"ciao\")")
	if err == nil {
		t.Fatal("Can use pcall\n")
	}

	err = L.DoString("unsafe_pcall(print, \"ciao\")")
	if err != nil {
		t.Fatal("Can not use unsafe_pcall\n")
	}
}

func TestCall(t *testing.T) {
	L := NewState()
	L.OpenLibs()
	defer L.Close()

	test := func(L *State) int {
		arg1 := L.ToString(1)
		arg2 := L.ToString(2)
		arg3 := L.ToString(3)

		if arg1 != "Argument1" {
			t.Fatal("Got wrong argument (1)")
		}

		if arg2 != "Argument2" {
			t.Fatal("Got wrong argument (2)")
		}

		if arg3 != "Argument3" {
			t.Fatal("Got wrong argument (3)")
		}

		L.PushString("Return1")
		L.PushString("Return2")

		return 2
	}

	L.Register("test", test)

	L.PushString("Dummy")
	L.GetGlobal("test")
	L.PushString("Argument1")
	L.PushString("Argument2")
	L.PushString("Argument3")
	err := L.Call(3, 2)

	if err != nil {
		t.Fatalf("Error executing call: %v\n", err)
	}

	dummy := L.ToString(1)
	ret1 := L.ToString(2)
	ret2 := L.ToString(3)

	if dummy != "Dummy" {
		t.Fatal("The stack was disturbed")
	}

	if ret1 != "Return1" {
		t.Fatalf("Wrong return value (1) got: <%s>", ret1)
	}

	if ret2 != "Return2" {
		t.Fatalf("Wrong return value (2) got: <%s>", ret2)
	}
}

// equivalent to basic.go
func TestLikeBasic(t *testing.T) {
	L := NewState()
	defer L.Close()
	L.OpenLibs()

	testCalled := 0

	test := func(L *State) int {
		testCalled++
		return 0
	}

	test2Arg := -1
	test2Argfrombottom := -1
	test2 := func(L *State) int {
		test2Arg = L.CheckInteger(-1)
		test2Argfrombottom = L.CheckInteger(1)
		return 0
	}

	L.GetField(LUA_GLOBALSINDEX, "print")
	L.PushString("Hello World!")
	if err := L.Call(1, 0); err != nil {
		t.Fatalf("Call to print returned error")
	}

	L.PushGoFunction(test)
	L.PushGoFunction(test)
	L.PushGoFunction(test)
	L.PushGoFunction(test2)
	L.PushInteger(42)
	if err := L.Call(1, 0); err != nil {
		t.Fatalf("Call to print returned error")
	}
	if (test2Arg != 42) || (test2Argfrombottom != 42) {
		t.Fatalf("Call to test2 didn't work")
	}

	if err := L.Call(0, 0); err != nil {
		t.Fatalf("Call to print returned error")
	}
	if err := L.Call(0, 0); err != nil {
		t.Fatalf("Call to print returned error")
	}
	if err := L.Call(0, 0); err != nil {
		t.Fatalf("Call to print returned error")
	}
	if testCalled != 3 {
		t.Fatalf("Test function not called the correct number of times: %d\n", testCalled)
	}

	// this will fail as we didn't register test2 function
	if err := L.DoString("test2(42)"); err == nil {
		t.Fatal("No error when calling unregistered function")
	}
}

// equivalent to quickstart.go
func TestLikeQuickstart(t *testing.T) {
	adder := func(L *State) int {
		a := L.ToInteger(1)
		b := L.ToInteger(2)
		L.PushInteger(int64(a + b))
		return 1
	}

	L := NewState()
	defer L.Close()
	L.OpenLibs()

	L.Register("adder", adder)

	if err := L.DoString("return adder(2, 2)"); err != nil {
		t.Fatalf("Error during call to adder: %v\n", err)
	}
	if r := L.ToInteger(1); r != 4 {
		t.Fatalf("Wrong return value from adder (was: %d)\n", r)
	}
}

// equivalent to userdata.go
func TestLikeUserdata(t *testing.T) {
	type Userdata struct {
		a, b int
	}

	userDataProper := func(L *State) {
		rawptr := L.NewUserdata(uintptr(unsafe.Sizeof(Userdata{})))
		var ptr *Userdata
		ptr = (*Userdata)(rawptr)
		ptr.a = 2
		ptr.b = 3

		rawptr2 := L.ToUserdata(-1)
		ptr2 := (*Userdata)(rawptr2)

		if ptr != ptr2 {
			t.Fatalf("Failed to create userdata\n")
		}
	}

	testCalled := 0
	test := func(L *State) int {
		testCalled++
		return 0
	}

	goDefinedFunctions := func(L *State) {
		// example_function is registered inside Lua VM
		L.Register("test", test)

		// This code demonstrates checking that a value on the stack is a go function
		L.CheckStack(1)
		L.GetGlobal("test")
		if !L.IsGoFunction(-1) {
			t.Fatalf("IsGoFunction failed to recognize a Go function object")
		}
		L.Pop(1)

		// We call example_function from inside Lua VM
		testCalled = 0
		if err := L.DoString("test()"); err != nil {
			t.Fatalf("Error executing test function: %v\n", err)
		}
		if testCalled != 1 {
			t.Fatalf("It appears the test function wasn't actually called\n")
		}
	}

	type TestObject struct {
		AField int
	}

	goDefinedObjects := func(L *State) {
		z := &TestObject{42}

		L.PushGoStruct(z)
		L.SetGlobal("z")

		// This code demonstrates checking that a value on the stack is a go object
		L.CheckStack(1)
		L.GetGlobal("z")
		if !L.IsGoStruct(-1) {
			t.Fatal("IsGoStruct failed to recognize a Go struct\n")
		}
		L.Pop(1)

		// This code demonstrates access and assignment to a field of a go object
		if err := L.DoString("return z.AField"); err != nil {
			t.Fatal("Couldn't execute code")
		}
		before := L.ToInteger(-1)
		L.Pop(1)
		if before != 42 {
			t.Fatalf("Wrong value of z.AField before change (%d)\n", before)
		}
		if err := L.DoString("z.AField = 10;"); err != nil {
			t.Fatal("Couldn't execute code")
		}
		if err := L.DoString("return z.AField"); err != nil {
			t.Fatal("Couldn't execute code")
		}
		after := L.ToInteger(-1)
		L.Pop(1)
		if after != 10 {
			t.Fatalf("Wrong value of z.AField after change (%d)\n", after)
		}
	}

	L := NewState()
	defer L.Close()
	L.OpenLibs()

	userDataProper(L)
	goDefinedFunctions(L)
	goDefinedObjects(L)
}

func TestStackTrace(t *testing.T) {
	L := NewState()
	defer L.Close()
	L.OpenLibs()

	err := L.DoFile("../example/calls.lua")
	if err == nil {
		t.Fatal("No error returned from the execution of calls.lua")
	}

	le := err.(*LuaError)

	if le.Code() != LUA_ERRERR {
		t.Fatalf("Wrong kind of error encountered running calls.lua: %v (%d %d)\n", le, le.Code(), LUA_ERRERR)
	}

	if len(le.StackTrace()) != 6 {
		t.Fatalf("Wrong size of stack trace (%v)\n", le.StackTrace())
	}
}

func TestConv(t *testing.T) {
	L := NewState()
	defer L.Close()
	L.OpenLibs()

	L.PushString("10")
	n := L.ToNumber(-1)
	if n != 10 {
		t.Fatalf("Wrong conversion (str -> int)")
	}
	if L.Type(-1) != LUA_TSTRING {
		t.Fatalf("Wrong type (str)")
	}

	L.Pop(1)

	L.PushInteger(10)
	s := L.ToString(-1)
	if s != "10" {
		t.Fatalf("Wrong conversion (int -> str)")
	}

	L.Pop(1)

	L.PushString("a\000test")
	s = L.ToString(-1)
	if s != "a\000test" {
		t.Fatalf("Wrong conversion (str -> str): <%s>", s)
	}
}

func Test101CoroutineRunning(t *testing.T) {
	L := NewState()
	L.OpenLibs()
	defer L.Close()

	butterCalled := 0
	butter := func(L *State) int {
		butterCalled++
		fmt.Printf("in butter() callback! here is butterVm's stack:\n")
		DumpLuaStack(L)
		fmt.Printf("above is butterVm's stack\n")
		tot := 0.0
		tot += L.ToNumber(-1)
		tot += L.ToNumber(-2)
		tot += L.ToNumber(-3)
		L.Pop(3)
		L.PushNumber(tot)
		fmt.Printf("tot was %v\n", tot)
		return 1
	}

	L.Register("butter", butter)

	// This code demonstrates checking that a value on the stack is a go function
	L.CheckStack(1)
	L.GetGlobal("butter")
	if !L.IsGoFunction(-1) {
		t.Fatalf("IsGoFunction failed to recognize a Go function object")
	}
	L.Pop(1)

	// We call example_function from inside Lua VM
	butterCalled = 0
	if err := L.DoString("a = {coroutine.resume(coroutine.create(function() return butter(4,5,6); end))}; for k,v in pairs(a) do print('a.key= ',k, '  value:', v); end; b = a[2]"); err != nil {
		t.Fatalf("Error executing butter function: %v\n", err)
	}
	if butterCalled != 1 {
		t.Fatalf("It appears the butter function wasn't actually called\n")
	}
	fmt.Printf("butterCalled = %v\n", butterCalled)

	L.GetGlobal("b")
	top := L.GetTop()
	obsB := L.ToNumber(top)
	if obsB != 15.0 {
		t.Fatalf("butter() summing 4+5+6 should have given 15, but instead got %v\n", obsB)
	}
}

func Test102LuaRegsitryIsPerState(t *testing.T) {
	// We test the  assumption
	// that the Lua Registry is shared by all
	// coroutines within a main C.lua State.
	// However two different C.lua_States
	// are expected to have distinct registries.
	//
	// If validated, we'll use this fact to
	// store a pointer to the main C.lua_State
	// in the registry, and have all
	// coroutines use the key to find their
	// main state.
	//
	// Result: the assumption was confirmed.
	// The registry is distinct per main State,
	// and shared by coroutines within one state.

	L := NewState()
	L.OpenLibs()
	defer L.Close()

	L2 := NewState()
	L2.OpenLibs()
	defer L2.Close()

	key := "lua_test.my_registry_key"
	val := "lua_test.my_value_for_testing"
	L.PushString(key)
	L.PushString(val)
	L.SetTable(LUA_REGISTRYINDEX)
	top := L.GetTop()
	if top != 0 {
		panic("expected empty stack")
	}
	L.PushString(key)
	L.GetTable(LUA_REGISTRYINDEX)
	if L.IsNil(-1) {
		panic("expected value back")
	}
	obsVal := L.ToString(-1)
	if obsVal != val {
		panic("expected obsVal to match val")
	}
	//fmt.Printf("good: retreived val from L registry\n")
	L.Pop(1)

	// now query the L2 registry
	L2.PushString(key)
	L2.GetTable(LUA_REGISTRYINDEX)
	if !L2.IsNil(-1) {
		fmt.Printf("bad, expected nil, got: '%s'\n", LuaStackPosToString(L2, -1))
		panic("expected nil back when querying L2 registry for key")
	}
	//fmt.Printf("good: did not retreived val from L2 registry under key\n")

	// now check that a new coroutine in L sees the same registry.
	L3 := L.NewThread()
	L3.PushString(key)
	L3.GetTable(LUA_REGISTRYINDEX)
	if L3.IsNil(-1) {
		panic("expected value back")
	}
	obsVal3 := L3.ToString(-1)
	if obsVal3 != val {
		panic("expected obsVal3 to match val")
	}
	//fmt.Printf("good: retreived val from L3 registry\n")
}

func assert(t *testing.T, b bool, msg string) {
	if !b {
		t.Fatal(msg)
	}
}

func Test103ToThreadDeduplicatesCoroutines(t *testing.T) {
	// when ToThread encounters the same coroutine
	// again, it should return the prior *State, and
	// not generate a new wrapper for the same coroutine.

	L := NewState()
	L.OpenLibs()
	defer L.Close()

	isMain := L.PushThread()
	if !isMain {
		t.Fatal("should have gotten isMain true!")
	}
	thr := L.ToThread(-1)
	if thr != L {
		t.Fatal("ToThread should dedup")
	}

	L2 := NewState()
	L2.OpenLibs()
	defer L2.Close()

	L2.PushThread()
	thr2 := L2.ToThread(-1)
	if thr2 != L2 {
		t.Fatal("ToThread should dedup L2")
	}

	if thr == thr2 {
		t.Fatal("thr should not equal thr2")
	}

	// make some new coroutines, make sure they are deduped.
	co2b := L2.NewThread()
	if co2b == thr2 || co2b == thr {
		t.Fatal("co2b should not equal thr2 or thr")
	}

	co2b_isMain := co2b.PushThread()
	if co2b.ToThread(-1) != co2b {
		t.Fatal("co2b was not deduped!")
	}
	assert(t, !co2b_isMain, "co2b_isMain should not have been a main thread!")

	// coroutines from Lua first
	err := co2b.DoString("a = 1; return coroutine.create(function() return a; end)")
	if err != nil {
		t.Fatalf("DoString returned an error: %v\n", err)
	}

	fmt.Printf("calling ToThread()")
	thr4 := co2b.ToThread(-1)
	fmt.Printf("thr4 = '%p'/'%#v'\n", thr4, thr4)
	assert(t, thr4.AllCoro == nil, "non-main coroutines should have nil AllCoro maps")
	/*
		for k, v := range L2.AllCoro {
			fmt.Printf("\n L2.AllCoro k = '%#v', v='%p'/'%#v'\n", k, v, v)
		}
	*/
	assert(t, L2.AllCoro[thr4.Upos] == thr4, "thr4 should be found in L2's AllCoro, at Upos")
}

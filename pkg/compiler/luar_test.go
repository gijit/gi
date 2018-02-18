package compiler

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/gijit/gi/pkg/verb"
	cv "github.com/glycerine/goconvey/convey"
	"github.com/glycerine/golua/lua"
	"github.com/glycerine/luar"
)

var _ = math.MinInt64

func getPanicValue(f func()) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			switch x := r.(type) {
			case error:
				err = x
				return
			case string:
				err = fmt.Errorf(x)
				return
			}
			panic(fmt.Sprintf("unknown panic type: '%v'/type=%T", r, r))
		}
	}()
	f()
	return nil
}

func Test053GoToLuarThenLuarToGo(t *testing.T) {

	cv.Convey(`luar.GoToLua followed by luar.LuaToGo should invert cleanly, giving us back a copy like the original`, t, func() {

		//vm, err := NewLuaVmWithPrelude(nil)
		//panicOn(err)

		vm := luar.Init()
		a := []int{6, 7, 8}
		luar.GoToLua(vm, &a)
		DumpLuaStack(vm)
		b := []int{}
		top := vm.GetTop()
		luar.LuaToGo(vm, top, &b)
		pp("a == '%#v'", a)
		pp("b == '%#v'", b)

		cv.So(b, cv.ShouldResemble, a)

	})
}

func Test054GoToLuarThenLuarToGo(t *testing.T) {

	cv.Convey(`luar.GoToLua followed by luar.LuaToGo for []interface{}`, t, func() {

		//vm, err := NewLuaVmWithPrelude(nil)
		//panicOn(err)

		var six int64 = 6

		vm := luar.Init()
		a := []interface{}{six, "hello"}
		luar.GoToLua(vm, &a)

		DumpLuaStack(vm)
		b := []interface{}{}
		top := vm.GetTop()
		luar.LuaToGo(vm, top, &b)
		pp("a == '%#v'", a)
		pp("b == '%#v'", b)

		// ah, luar converts integers to doubles. arg.
		for i := range a {
			pp("deep equal? %v", reflect.DeepEqual(a[i], b[i]))
			if !reflect.DeepEqual(a[i], b[i]) {
				// diff at i=0: '6'/int != '6'/float64
				pp("diff at i=%v: '%#v'/%T != '%#v'/%T", i, a[i], a[i], b[i], b[i])
			}
		}
		cv.So(b, cv.ShouldResemble, a)

	})
}

func Test055_cdata_Int64_GoToLuar_Then_LuarToGo(t *testing.T) {

	cv.Convey(`luar.GoToLua then LuarToGo should preserve int64 cdata values`, t, func() {

		// +10000 so we catch the loss of precision
		// from the conversion to double and back.
		var a int64 = math.MinInt64 + 10000

		//vm := luar.Init()

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		// seems we should do this at all, but rather
		// use luajit to inject the data
		//luar.GoToLua(vm, &a)

		//putOnTopOfStack := fmt.Sprintf(`return int32(%dLL)`, a) // int32
		putOnTopOfStack := fmt.Sprintf(`return %dLL`, a) // int64
		interr := vm.LoadString(putOnTopOfStack)
		if interr != 0 {
			pp("interr %v on vm.LoadString for dofile on '%s'", interr, putOnTopOfStack)
			msg := DumpLuaStackAsString(vm)
			vm.Pop(1)
			panic(fmt.Errorf("error in vm.LoadString of '%s': Details: '%s'", putOnTopOfStack, msg))
		}
		err = vm.Call(0, 1)
		if err != nil {
			msg := DumpLuaStackAsString(vm)
			vm.Pop(1)
			panic(fmt.Errorf("error in Call after vm.LoadString of '%s': '%v'. Details: '%s'", putOnTopOfStack, err, msg))
		}
		ctype := vm.LuaJITctypeID(-1)
		pp("ctype = %v", ctype)

		// ctype ==  0 means it wasn't a ctype
		// ctype ==  5 for int8
		// ctype ==  6 for uint8
		// ctype ==  7 for int16
		// ctype ==  8 for uint16
		// ctype ==  9 for int32
		// ctype == 10 for uint32
		// ctype == 11 for int64
		// ctype == 12 for uint64
		// ctype == 13 for float32
		// ctype == 14 for float64

		DumpLuaStack(vm)
		var b int64
		top := vm.GetTop()
		panicOn(luar.LuaToGo(vm, top, &b))
		pp("a == '%#v'", a)
		pp("b == '%#v'", b)

		cv.So(b, cv.ShouldResemble, a)

	})
}

func Test056_cdata_Uint64_unsigned_int64_GoToLuar_Then_LuarToGo(t *testing.T) {

	cv.Convey(`luar.GoToLua then LuarToGo should preserve uint64 cdata values`, t, func() {

		var a uint64 = uint64(math.MaxInt64 + 10000)

		//vm := luar.Init()

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		// seems we should do this at all, but rather
		// use luajit to inject the data
		//luar.GoToLua(vm, &a)

		putOnTopOfStack := fmt.Sprintf(`return %dULL`, a) // uint64
		interr := vm.LoadString(putOnTopOfStack)
		if interr != 0 {
			pp("interr %v on vm.LoadString for dofile on '%s'", interr, putOnTopOfStack)
			msg := DumpLuaStackAsString(vm)
			vm.Pop(1)
			panic(fmt.Errorf("error in vm.LoadString of '%s': Details: '%s'", putOnTopOfStack, msg))
		}
		err = vm.Call(0, 1)
		if err != nil {
			msg := DumpLuaStackAsString(vm)
			vm.Pop(1)
			panic(fmt.Errorf("error in Call after vm.LoadString of '%s': '%v'. Details: '%s'", putOnTopOfStack, err, msg))
		}
		ctype := vm.LuaJITctypeID(-1)
		pp("ctype = %v", ctype)

		// ctype ==  0 means it wasn't a ctype
		// ctype ==  5 for int8
		// ctype ==  6 for uint8
		// ctype ==  7 for int16
		// ctype ==  8 for uint16
		// ctype ==  9 for int32
		// ctype == 10 for uint32
		// ctype == 11 for int64
		// ctype == 12 for uint64
		// ctype == 13 for float32
		// ctype == 14 for float64

		DumpLuaStack(vm)
		var b uint64
		top := vm.GetTop()
		panicOn(luar.LuaToGo(vm, top, &b))
		pp("a == '%#v'", a)
		pp("b == '%#v'", b)

		cv.So(b, cv.ShouldResemble, a)

	})
}

func Test057_cdata_GoToLuar_Then_LuarToGo_Mistypes_are_flagged(t *testing.T) {

	cv.Convey(`luar.GoToLua then LuarToGo: trying to decode uint64 into int64 should be caught`, t, func() {

		var a uint64 = uint64(math.MaxInt64 + 10000)

		//vm := luar.Init()

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		// seems we should do this at all, but rather
		// use luajit to inject the data
		//luar.GoToLua(vm, &a)

		putOnTopOfStack := fmt.Sprintf(`return %dULL`, a) // uint64
		interr := vm.LoadString(putOnTopOfStack)
		if interr != 0 {
			pp("interr %v on vm.LoadString for dofile on '%s'", interr, putOnTopOfStack)
			msg := DumpLuaStackAsString(vm)
			vm.Pop(1)
			panic(fmt.Errorf("error in vm.LoadString of '%s': Details: '%s'", putOnTopOfStack, msg))
		}
		err = vm.Call(0, 1)
		if err != nil {
			msg := DumpLuaStackAsString(vm)
			vm.Pop(1)
			panic(fmt.Errorf("error in Call after vm.LoadString of '%s': '%v'. Details: '%s'", putOnTopOfStack, err, msg))
		}
		ctype := vm.LuaJITctypeID(-1)
		pp("ctype = %v", ctype)

		// ctype ==  0 means it wasn't a ctype
		// ctype ==  5 for int8
		// ctype ==  6 for uint8
		// ctype ==  7 for int16
		// ctype ==  8 for uint16
		// ctype ==  9 for int32
		// ctype == 10 for uint32
		// ctype == 11 for int64
		// ctype == 12 for uint64
		// ctype == 13 for float32
		// ctype == 14 for float64

		DumpLuaStack(vm)
		var b int64
		top := vm.GetTop()
		err = getPanicValue(func() { luar.LuaToGo(vm, top, &b) })
		// this test is all about getting a type mis-match error back:
		pp("b = %v, a = %v", b, a)
		cv.So(err, cv.ShouldNotBeNil)
		cv.So(err.Error(), cv.ShouldResemble, "reflect.Set: value of type uint64 is not assignable to type int64")
	})
}

func Test060_LuaToGo_handles_slices(t *testing.T) {

	cv.Convey(`if the compiled lua code creates a slice, then passes it to a function in a compiled Go package, then LuaToGo should handle our custom slices-with-metables`, t, func() {

		src := `a := []int{5,6,4};`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		// TODO undo the verb.Verbose setting
		verb.VerboseVerbose = true

		vm.GetGlobal("__gijit_tsys")
		cv.So(vm.IsNil(-1), cv.ShouldBeFalse)
		pp("good: we found __gijit_tsys")
		vm.Pop(1)

		start := vm.GetTop()
		vm.GetGlobal("_MIS_SPELLED")
		cv.So(vm.IsNil(-1), cv.ShouldBeTrue)
		pp("good: we got nil for _MIS_SPELLED")
		cv.So(vm.GetTop(), cv.ShouldEqual, start+1)
		pp(`good: we confirmed that GetGlobal("unknown") still pushes nil onto top of stack`)
		vm.Pop(1)

		pp("stack prior to translation:")
		DumpLuaStack(vm)

		translation := inc.Tr([]byte(src))
		pp("go:'%s'  -->  '%s' in lua\n", src, string(translation))

		cv.So(string(translation), matchesLuaSrc, `
  	__type__anon_sliceType = __sliceType(__type__int); 
  
  	a = __type__anon_sliceType({[0]=5LL, 6LL, 4LL});
    `)

		LoadAndRunTestHelper(t, vm, translation)

		vm.GetGlobal("a")
		pp("stack after GetGlobal on 'a':")
		DumpLuaStack(vm)
		b := []int{}
		top := vm.GetTop()

		//pp("ObjLen should invoke the # operation")
		//olen := vm.ObjLen(top)
		//cv.So(olen, cv.ShouldEqual, 3) // failing here, apparently objlen does not invoke the metatable method!

		// query slice for len

		// this really should get us 3 back
		getfield(vm, -1, "__length")
		aLen := vm.ToNumber(-1)
		cv.So(aLen, cv.ShouldEqual, 3)
		vm.Pop(1)

		pp("good: aLen was %v, stack is now:", aLen)
		DumpLuaStack(vm)

		// Line 286: - cannot convert Lua value 'function: %!p(uintptr=75307960)' (function) to []int
		panicOn(luar.LuaToGo(vm, top, &b))
		cv.So(b, cv.ShouldResemble, []int{5, 6, 4})
	})
}

// getfield will
// assume that table is on the stack top, and
// returns with the value (that which corresponds to key) on
// the top of the stack.
// If value not present, then a nil is on top of the stack,
func getfield(L *lua.State, tableIdx int, key string) {
	L.PushValue(tableIdx)
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

type myGoTestStruct struct {
	priv int
	Pub  int
}

func testOp(m *myGoTestStruct) string {
	return fmt.Sprintf("priv:%v, Pub:%v", m.priv, m.Pub)
}

func Test068_LuaToGo_handles_structs(t *testing.T) {

	cv.Convey(`we want to translate lua structs into go structs, private and public fields`, t, func() {

		src := `
type myGoTestStruct struct {
	priv int
	Pub  int
}

ts := &myGoTestStruct{
	priv: 4,
	Pub:  5,
}

// stub definition, to keep the type checker happy for now.
// We actually want to test calling into compiled go, so
// we'll redefine testOp later.
func testOp(m *myGoTestStruct) string { return "" }
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(src))
		pp("go:'%s'  -->  '%s' in lua\n", src, string(translation))

		//		cv.So(string(translation), matchesLuaSrc,
		//			`a =_gi_NewSlice("int",{[0]=5,6,4});`)

		LoadAndRunTestHelper(t, vm, translation)

		// replace the stub definition of testOp with the real one.
		luar.Register(vm, "", luar.Map{
			"testOp": testOp,
		})

		src = `s := testOp(ts)`
		translation = inc.Tr([]byte(src))
		pp("got to here, go:'%s'  -->  '%s' in lua\n", src, string(translation))

		//		cv.So(string(translation), matchesLuaSrc,
		//			`a =_gi_NewSlice("int",{[0]=5,6,4});`)

		LoadAndRunTestHelper(t, vm, translation)

		LuaMustString(vm, "s", "priv:4, Pub:5")
	})
}

package compiler

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
	"github.com/glycerine/luar"
)

var _ = math.MinInt64

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
		cv.So(err.Error(), cv.ShouldResemble, "reflect.Set: value of type uint64 is not assignable to type int64")
	})
}

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

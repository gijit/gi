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

		vm := luar.Init()

		// seems we should do this at all, but rather
		// use luajit to inject the data
		//luar.GoToLua(vm, &a)

		putOnTopOfStack := fmt.Sprintf(`return %dLL`, a)
		interr := vm.LoadString(putOnTopOfStack)
		if interr != 0 {
			pp("interr %v on vm.LoadString for dofile on '%s'", interr, putOnTopOfStack)
			msg := DumpLuaStackAsString(vm)
			vm.Pop(1)
			panic(fmt.Errorf("error in vm.LoadString of '%s': Details: '%s'", putOnTopOfStack, msg))
		}
		err := vm.Call(0, 1)
		if err != nil {
			msg := DumpLuaStackAsString(vm)
			vm.Pop(1)
			panic(fmt.Errorf("error in Call after vm.LoadString of '%s': '%v'. Details: '%s'", putOnTopOfStack, err, msg))
		}
		ctype := vm.LuaJITctypeID()
		pp("ctype = %v", ctype)
		DumpLuaStack(vm)
		var b int64
		top := vm.GetTop()
		panicOn(luar.LuaToGo(vm, top, &b))
		pp("a == '%#v'", a)
		pp("b == '%#v'", b)

		cv.So(b, cv.ShouldResemble, a)

	})
}

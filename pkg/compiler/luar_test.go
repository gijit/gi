package compiler

import (
	"math"
	"reflect"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
	"github.com/glycerine/luar"
)

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
		luar.GoToLuaProxy(vm, &a)

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

func Test055_Int64_GoToLuar_Then_LuarToGo(t *testing.T) {

	cv.Convey(`luar.GoToLua then LuarToGo should preserve int64 cdata values`, t, func() {

		var a int64 = math.MinInt64

		vm := luar.Init()
		luar.GoToLua(vm, &a)

		DumpLuaStack(vm)
		var b int64
		top := vm.GetTop()
		luar.LuaToGo(vm, top, &b)
		pp("a == '%#v'", a)
		pp("b == '%#v'", b)

		cv.So(b, cv.ShouldResemble, a)

	})
}

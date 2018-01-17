package compiler

import (
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

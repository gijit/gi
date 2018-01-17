package compiler

import (
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test033DefersRunOnPanic(t *testing.T) {

	cv.Convey(`panic invokes only those infers encountered on the path of control, in last-declared-first-to-run order`, t, func() {

		code := `
a := -1
b := 0
d1a := -1
func f() (ret0 int, ret1 int) {
  defer func(a int) {
      d1a = a // d1a should be -1, because defer captures variables at the call point.
      print("first defer running, a=", a, " b=",b)
      r := recover()
      if r != nil {
          b = b + 3
      }
  }(a)
  a = 0
  panic("ouch")
  defer func() {
      println("second defer running, a=", a, " b=",b)
      b = b * 7
  }()
  a = 1
  b = 1
  return b, 58
}
f()
// now b should be set to 3
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		inc := NewIncrState(vm)
		translation := inc.Tr([]byte(code))

		pp("translation='%s'", string(translation))

		// too comple and fragile to verify code. Just  verify that it happens correctly

		LuaRunAndReport(vm, string(translation))
		LuaMustInt(vm, "r1", 0)
		LuaMustInt(vm, "r2", 0)
		LuaMustInt(vm, "b", 3)
		LuaMustInt(vm, "d1a", -1)
	})
}

func Test035DefersRunWithoutPanic(t *testing.T) {

	cv.Convey(`defers run without panic`, t, func() {

		code := `func f () { defer println("say hello, gracie"); }`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		inc := NewIncrState(vm)
		translation := inc.Tr([]byte(code))

		pp("translation='%s'", string(translation))

		// too comple and fragile to verify code. Just  verify that it happens correctly

		LuaRunAndReport(vm, string(translation))

	})
}

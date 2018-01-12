package compiler

import (
	"testing"

	cv "github.com/glycerine/goconvey/convey"
	luajit "github.com/glycerine/golua/lua"
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

		inc := NewIncrState()
		translation := inc.Tr([]byte(code))

		pp("translation='%s'", string(translation))

		// too comple and fragile to verify code. Just  verify that it happens correctly
		vm := luajit.NewState()
		defer vm.Close()
		vm.OpenLibs()
		files, err := FetchPrelude(".")
		panicOn(err)
		LuaDoFiles(vm, files)

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

		inc := NewIncrState()
		translation := inc.Tr([]byte(code))

		pp("translation='%s'", string(translation))

		// too comple and fragile to verify code. Just  verify that it happens correctly
		vm := luajit.NewState()
		defer vm.Close()
		vm.OpenLibs()
		files, err := FetchPrelude(".")
		panicOn(err)
		LuaDoFiles(vm, files)

		LuaRunAndReport(vm, string(translation))

	})
}

/* old stuff, right??


func Test0__DeferOnUnwind(t *testing.T) {

	cv.Convey(`defer runs upon panic, in last-declared-first-to-run order`, t, func() {

		code := `
a := 0
b := 0
func f() {
  defer func() {
      r = recover()
      b = b * r
      if a == 1 {
          b = b + 3
      }
  }()
  defer func() {
      b = b * 7
  }()
  panic(11)
  a = 1
  b = 1
}
f()
// now b should be set to 10.
`

		inc := NewIncrState()
		translation := inc.Tr([]byte(code))

		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace,
			`
a = 0;
b = 0;
function f()
  _defers={}
  _defers[1+#_defers] = function()
        if a == 1 then b = b + 3 end
  end
  _defers[1+#_defers] = function()
        b = b * 7
  end
  a = 1
  b = 1
  for _i = #_defers, 1, -1 do
     _defers[_i]()
  end
end
f()
`)

		// and verify that it happens correctly
		vm := luajit.NewState()
		defer vm.Close()
		vm.OpenLibs()
		files, err := FetchPrelude(".")
		panicOn(err)
		LuaDoFiles(vm, files)

		LuaRunAndReport(vm, string(translation))
		LuaMustInt(vm, "b", 10)
	})
}

func Test0__Defer(t *testing.T) {

	cv.Convey(`defer runs methods at function end, in last-declared-first-to-run order`, t, func() {

		code := `
a := 0
b := 0
func f() {
  defer func() {
      if a == 1 {
          b = b + 3
      }
  }()
  defer func() {
      b = b * 7
  }()
  a = 1
  b = 1
}
f()
// now b should be set to 10.
`

		inc := NewIncrState()
		translation := inc.Tr([]byte(code))

		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace,
			`
a = 0;
b = 0;
function f()
  _defers={}
  _defers[1+#_defers] = function()
        if a == 1 then b = b + 3 end
  end
  _defers[1+#_defers] = function()
        b = b * 7
  end
  a = 1
  b = 1
  for _i = #_defers, 1, -1 do
     _defers[_i]()
  end
end
f()
`)

		// and verify that it happens correctly
		vm := luajit.NewState()
		defer vm.Close()
		vm.OpenLibs()
		files, err := FetchPrelude(".")
		panicOn(err)
		LuaDoFiles(vm, files)

		LuaRunAndReport(vm, string(translation))
		LuaMustInt(vm, "b", 10)
	})
}

*/

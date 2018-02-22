package compiler

import (
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test033DefersRunOnPanic(t *testing.T) {

	cv.Convey(`panic invokes only those defers encountered on the path of control, in last-declared-first-to-run order`, t, func() {

		code := `
a := -1
b := 0
d1a := -1
func f() (ret0 int, ret1 int) {
  defer func(a int) {
      d1a = a // d1a should be -1, because defer captures variables at the call point.
      println("first defer running, a=", a, " b=",b)
      r := recover()
      if r != nil {
          println("rocover was not nil, recovered from a panic")
          b = b + 3
          ret1 = b
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
r0, r1 := f()
// now b should be set to 3
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		inc := NewIncrState(vm, nil)
		translation, err := inc.Tr([]byte(code))
		panicOn(err)

		pp("translation='%s'", string(translation))

		// too comple and fragile to verify code. Just  verify that it happens correctly

		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "r1", 3)
		LuaMustInt64(vm, "d1a", -1)
	})
}

func Test034TwoDefersRunOnPanic(t *testing.T) {

	cv.Convey(`panic invokes only those infers encountered on the path of control, in last-declared-first-to-run order`, t, func() {

		code := `
a := -1
b := 0
d1a := -1
func f() (ret0 int, ret1 int) {
  defer func(a int) {
      d1a = a // d1a should be -1, because defer captures variables at the call point.
      println("first defer running, a=", a, " b=",b)
      r := recover()
      if r != nil {
          b = b + 3
          ret1 = b
      }
      ret0 = ret0 + b
  }(a)
  a = 0
  defer func() {
      println("second defer running, a=", a, " b=",b)
      b = (b+1) * 7
      ret0 = ret0 + b
  }()
  a = 1
  b = 1
  return b, 58
}
r0, r1 := f()
// now b should be set to 3
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		inc := NewIncrState(vm, nil)
		translation, err := inc.Tr([]byte(code))
		panicOn(err)

		pp("translation='%s'", string(translation))

		// too comple and fragile to verify code. Just  verify that it happens correctly

		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "r0", 29)
		LuaMustInt64(vm, "d1a", -1)
	})
}

func Test035DefersRunWithoutPanic(t *testing.T) {

	cv.Convey(`defers run without panic`, t, func() {

		code := `func f () { defer println("say hello, gracie"); }`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		inc := NewIncrState(vm, nil)
		translation, err := inc.Tr([]byte(code))
		panicOn(err)

		pp("translation='%s'", string(translation))

		// too comple and fragile to verify code. Just  verify that it happens correctly

		LuaRunAndReport(vm, string(translation))

	})
}

func Test035bNamedReturnValuesAreReturned(t *testing.T) {

	cv.Convey(`__namedNames was missing named return values, so they weren't being returned`, t, func() {

		code := `func f() (r int) {defer func() { println(" r was ", r); r++; println(" r after++ is ", r) }(); r = 3; return r}; a := f();`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		inc := NewIncrState(vm, nil)
		translation, err := inc.Tr([]byte(code))
		panicOn(err)

		pp("translation='%s'", string(translation))

		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "a", 4)
	})
}

func Test035cNamedReturnValuesDontPolluteGlobalEnv(t *testing.T) {

	cv.Convey(`the named return values of a function should not contaminate the global env`, t, func() {

		code := `glob:=3; func f() (glob, x int) { glob =2; x = 1; return }; a, _ := f()`

		// glob should stay 3
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		inc := NewIncrState(vm, nil)
		translation, err := inc.Tr([]byte(code))
		panicOn(err)

		pp("translation='%s'", string(translation))

		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "glob", 3)
		LuaMustBeInGlobalEnv(vm, "glob")
		LuaMustNotBeInGlobalEnv(vm, "x")
		LuaMustNotBeInGlobalEnv(vm, "y") // control
	})
}

func Test033bDefersWorkOnDirectionFunctionCalls(t *testing.T) {

	cv.Convey(`defers on direct method calls, not function literals, also work`, t, func() {

		code := `
a := 0
func double_a() { a = a * 2 }
func f() {
  defer double_a()
  a++
}
f()
// now 'a' should be 2, but if the defer ran the function right away, then 'a' would be 1
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		inc := NewIncrState(vm, nil)
		translation, err := inc.Tr([]byte(code))
		panicOn(err)

		pp("translation='%s'", string(translation))

		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "a", 2)
		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test033cDefersWorkOnDirectionFunctionCalls(t *testing.T) {

	cv.Convey(`defer on a repeated direct function call`, t, func() {

		code := `
import "fmt"

var result string

func addInt(i int) { result += fmt.Sprint(i) }

func test1helper() {
	for i := 0; i < 10; i++ {
		defer addInt(i)
	}
}
test1helper()
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		inc := NewIncrState(vm, nil)
		translation, err := inc.Tr([]byte(code))
		panicOn(err)

		pp("translation='%s'", string(translation))

		if false {
			by, err := inc.Tr([]byte(code))
			panicOn(err)

			cv.So(string(by), matchesLuaSrc, `
result = ""

addInt = function(i) 
   result = result .. (fmt.Sprint(i));
end;

test1helper =    function(...) 
   local __orig = {...}
   local __defers={}
   local __zeroret = {}
   local __namedNames = {}
   local __actual=function()
      local i = 0LL;
      while (true) do
         if (not (i < 10LL)) then break; end

         local __defer_func = function(i)
            local i = i
            
            __defers[1+#__defers] = function()
               addInt(i)
            end
         end
         __defer_func(i)

         i = i + (1LL);
      end 

   end
   return __actuallyCall("", __actual, __namedNames, __zeroret, __defers, __orig)
end

test1helper();

`)
		}
		LuaRunAndReport(vm, string(translation))
		LuaMustString(vm, "result", "9876543210")
		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test033dDefersWorkOnDirectionFunctionCalls(t *testing.T) {

	cv.Convey(`defer on a repeated direct function call`, t, func() {

		code := `
import "fmt"

var result string

func addDotDotDot(v ...interface{}) { result += fmt.Sprint(v...) }

func test2helper() {
	for i := 0; i < 10; i++ {
		defer addDotDotDot(i)
	}
}
test2helper()
`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		inc := NewIncrState(vm, nil)
		translation, err := inc.Tr([]byte(code))
		panicOn(err)

		pp("translation='%s'", string(translation))
		LuaRunAndReport(vm, string(translation))
		LuaMustString(vm, "result", "9876543210")
		cv.So(true, cv.ShouldBeTrue)
	})
}

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

func Test033dDefersWorkOnDirectionFunctionCalls(t *testing.T) {

	cv.Convey(`defer on a repeated direct function call`, t, func() {

		code := `

var result int

var base int = 1

func addDotDotDot(v int) { result += (base * v); base=base*10; }

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
		LuaMustInt64(vm, "result", 123456789)
		cv.So(true, cv.ShouldBeTrue)
	})
}

// should print:
// go func() { defer fmt.Printf("hi there %#v\n", time.Now()) }()

func Test033mTest2FromRecoverStressTest(t *testing.T) {
	// taken from recover.go in the golang test suite

	cv.Convey(`test2 from the recover.go test in the golang test suite`, t, func() {

		code := `
import "runtime"

got_correct_value := false

func mustRecover(x interface{}) {
    println("top of mustRecover")

    a := doubleRecover()
    println("a is ", a)

    b := recover()
    println("b is ", b)

    c := recover()
    println("c is ", c)

	mustRecoverBody(a, b, c, x)
    println("done with mustRecover()")
}

func die() {
    println("die() invoked, calling runtime.Breakpoint()")
	runtime.Breakpoint() // can't depend on panic
}

func mustRecoverBody(v1, v2, v3, x interface{}) {
	v := v1
	if v != nil {
		println("spurious recover", v)
		die()
	}
	v = v2
	if v == nil {
		println("missing recover ")
		//println("assert x is int: ", x.(int)) // crashing by itself.
		die() // panic is useless here
	}
	if v != x {
		println("wrong value", v, x)
		die()
	}
    got_correct_value = true

	// the value should be gone now regardless
	v = v3
	if v != nil {
		println("recover didn't recover")
		die()
	}
    println("mustRecoverBody reached end without die.")
}


func doubleRecover() interface{} {
	return recover()
}

func test2() {
	// Recover only sees the panic argument
	// if it is called from a deferred call.
	// It does not see the panic when called from a call within 
    // a deferred call (too late)
	// nor does it see the panic when it *is* the deferred call (too early).
	defer mustRecover(2)
	defer recover() // should be no-op
    println("about to panic(2)")
	panic(2)
}
test2()
println("test2() ran")
`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		inc := NewIncrState(vm, nil)
		translation, err := inc.Tr([]byte(code))
		panicOn(err)

		pp("translation='%s'", string(translation))
		LuaRunAndReport(vm, string(translation))
		LuaMustBool(vm, "got_correct_value", true)
		cv.So(true, cv.ShouldBeTrue)
	})

	/* current translation looks okay, but is crashing
			       the scheduler. r2:

				__go_import("runtime")

		        got_correct_value = false;
				mustRecover = function(x)
				   mustRecoverBody(doubleRecover(), recover(), recover(), x);
				   print("done with mustRecover()");
				end;

				die = function()
				   print("die() invoked, calling runtime.Breakpoint()");
				   runtime.Breakpoint();
				end;

				mustRecoverBody = function(v1, v2, v3, x)
				   local v = v1;
				   if ( not (__interfaceIsEqual(v, nil))) then
				      print("spurious recover", v);
							die();
				   end
				   v = v2;
				   if (__interfaceIsEqual(v, nil)) then
				      print("missing recover", __assertType(x, __type__.int, 0));
				      die();
				   end
				   if ( not (__interfaceIsEqual(v, x))) then
				      print("wrong value", v, x);
				      die();
				   end
	               got_correct_value = true
				   v = v3;
				   if ( not (__interfaceIsEqual(v, nil))) then
				      print("recover didn't recover");
				      die();
				   end
				   print("mustRecoverBody reached end without die.");
				end;

				doubleRecover = function()
				   return  recover() ;
				end;

				test2 =
				   function(...)
				   local __orig = {...}
				   local __defers={}
				   local __zeroret = {}
				   local __namedNames = {}
				   local __actual=function()

				      local __defer_func = function(__gensym_1__arg)
				         local __gensym_1__arg = __gensym_1__arg;

				         __defers[1+#__defers] = function()
				            mustRecover(__gensym_1__arg);
				         end
				      end
				      __defer_func(2LL);


				      local __defer_func = function()

				         __defers[1+#__defers] = function()
				            recover();
				         end
				      end
				      __defer_func();

				      print("about to panic(2)");
				      panic(2LL);

				   end
				   return __actuallyCall("", __actual, __namedNames, __zeroret, __defers, __orig)
				   end
				;
				test2();

				print("test2() ran");
	*/
}

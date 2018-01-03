package compiler

import (
	"fmt"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
	"github.com/glycerine/luajit"
)

func Test001LuaTranslation(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("assignment", t, func() {
		cv.So(string(inc.Tr([]byte("a := 10;"))), cv.ShouldMatchModuloWhiteSpace, "a = 10;")
		pp("GOOD: past 1st")

		cv.So(string(inc.Tr([]byte("func adder(a, b int) int { return a + b};  sum1 := adder(5,5)"))), cv.ShouldMatchModuloWhiteSpace,
			`adder = function(a, b) 
				         return a + b;
                     end;
			         sum1 = adder(5,5);`)

		pp("GOOD: past 2nd")

		cv.So(string(inc.Tr([]byte("sum2 := adder(a,a)"))), cv.ShouldMatchModuloWhiteSpace,
			`sum2 = adder(a, a);`)
		pp("GOOD: past 3rd")

		cv.So(string(inc.Tr([]byte("func multiplier(a, b int) int { return a * b};  prod1 := multiplier(5,5)"))), cv.ShouldMatchModuloWhiteSpace,
			`multiplier = function(a, b) 
				         return (a * b);
                     end;
			         prod1 = multiplier(5,5);`)

	})
}

func Test002LuaEvalIncremental(t *testing.T) {

	// and then eval!
	vm := luajit.Newstate()
	defer vm.Close()
	vm.Openlibs()

	inc := NewIncrState()

	srcs := []string{"a := 10;", "func adder(a, b int) int { return a + b}; ", "sum := adder(a,a);"}
	for _, src := range srcs {
		translation := inc.Tr([]byte(src))
		fmt.Printf("go:'%s'  -->  '%s' in lua\n", src, translation)
		//fmt.Printf("go:'%#v'  -->  '%#v' in lua\n", src, translation)

		err := vm.Loadstring(string(translation))
		panicOn(err)
		err = vm.Pcall(0, 0, 0)
		panicOn(err)
		DumpLuaStack(vm)

		//fmt.Printf("v back = '%#v'\n", v)
	}
	vm.Getglobal("sum")
	top := vm.Gettop()
	value_int := vm.Tointeger(top)

	fmt.Printf("value_int=%v", value_int)
	if value_int != 20 {
		panic(fmt.Sprintf("expected 20, got %v", value_int))
	}
}

// func Test003ImportsAtRepl(t *testing.T) {
// 	inc := NewIncrState()

// 	cv.Convey("imports", t, func() {
// 		cv.So(string(inc.Tr([]byte(`import "fmt"; fmt.Printf("hello world!")`))), cv.ShouldMatchModuloWhiteSpace, "")
// 		pp("GOOD: past 1st import")
// 	})
// }

func Test004ExpressionsAtRepl(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("expressions alone at top level", t, func() {
		cv.So(string(inc.Tr([]byte(`a:=10; a`))), cv.ShouldMatchModuloWhiteSpace, "a=10; print(a);")
	})
}

func Test005BacktickStringsToLua(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("In order for Lua to parse Go's literal backtick strings containing ]=] correctly, Go's backtick strings `like ]=] this one` need to be translated from Go into Lua as three concatenated by .. strings: [=[\\nlike ]=] .. ']=]' .. [=[\\n this one]=]. Then Go `backtick` strings without ']=]' can translate to lua: [=[\\nbacktick]=]", t, func() {
		cv.So(string(inc.Tr([]byte("s:=`like ]=] this one`"))), cv.ShouldEqual, "s = [=[\nlike ]=] .. ']=]' .. [=[\n this one]=];")
		cv.So(string(inc.Tr([]byte("s:=`and this one`"))), cv.ShouldEqual, "s = [=[\nand this one]=];")
	})
}

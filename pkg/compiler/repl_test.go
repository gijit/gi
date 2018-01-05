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
		//pp("go:'%s'  -->  '%s' in lua\n", src, translation)
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

	pp("value_int=%v", value_int)
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
		cv.So(string(inc.Tr([]byte(`a:=10;`))), cv.ShouldMatchModuloWhiteSpace, `a=10;`)
	})
}

func Test005BacktickStringsToLua(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("Go backtick delimited strings should translate to Lua", t, func() {
		cv.So(string(inc.Tr([]byte("s:=`\n\n\"hello\"\n\n`"))), cv.ShouldMatchModuloWhiteSpace, `s = "\n\n\"hello\"\n\n";`)
		cv.So(string(inc.Tr([]byte("r:=`\n\n]]\"hello\"\n\n`"))), cv.ShouldMatchModuloWhiteSpace, `r = "\n\n]]\"hello\"\n\n";`)
	})
}

func Test006RedefinitionOfVariablesAllowed(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("At the repl, `a:=1; a:=2;` is allowed. We disable the traditional Go re-definition checks at the REPL", t, func() {
		cv.So(string(inc.Tr([]byte("a:=1; a:=2;"))), cv.ShouldMatchModuloWhiteSpace, `a=1; a=2;`)

		// and in separate calls:
		cv.So(string(inc.Tr([]byte("r:=`\n\n]]\"hello\"\n\n`"))), cv.ShouldMatchModuloWhiteSpace, `r = "\n\n]]\"hello\"\n\n";`)
		// and redefinition of r should work:
		cv.So(string(inc.Tr([]byte("r:=`a new definition`"))), cv.ShouldMatchModuloWhiteSpace, `r = "a new definition";`)
	})
}

func Test007SettingPreviouslyDefinedVariables(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("At the repl, in separate commands`a:=1; a=2;` sets a to 2", t, func() {

		// and in separate calls, where the second call sets the earlier variable.
		cv.So(string(inc.Tr([]byte("a:=1"))), cv.ShouldMatchModuloWhiteSpace, `a=1;`)
		cv.So(string(inc.Tr([]byte("b:=2"))), cv.ShouldMatchModuloWhiteSpace, `b=2;`)
		// and redefinition of a should work:
		pp("starting on a=2;")
		cv.So(string(inc.Tr([]byte("a=2;"))), cv.ShouldMatchModuloWhiteSpace, `a=2;`)
	})
}

func Test008IfThenElseInAFunction(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("if then else within a closure/function should compile into lua", t, func() {

		code := `a:=20; func hmm() { if a > 30 { println("over 30") } else {println("under or at 30")} }`
		// and in separate calls, where the second call sets the earlier variable.
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `a=20; hmm = function() if (a > 30) then print("over 30"); else print("under or at 30"); end end;`)
	})
}

func Test009NumericalForLoop(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("numerical for loops should compile into lua", t, func() {

		code := `a:=5; func hmm() { for i:=0; i < a; i++ { println(i) } }`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `a=5;
  	hmm = function() 
  		i = 0;
  		while (true) do
  			if (not (i < a)) then break; end
  			print(i);
  			i = i + (1);
  		 end
 	 end;
`)
	})
}

func Test010Slice(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("slice literal should compile into lua", t, func() {

		code := `a:=[]int{1,2,3}`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `a=_giSlice{[0]=1,2,3};`)
	})
}

func Test011MapAndRangeForLoop(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("maps and range for loops should compile into lua", t, func() {

		code := `a:=make(map[int]int); a[1]=10; a[2]=20; func hmm() { for k, v := range a { println(k," ",v) } }`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
a = {};
a[1] = 10;
a[2] = 20;
hmm = function() for k, v in pairs(a) do print(k, " ", v);  end end;`)
	})
}

func Test012SliceRangeForLoop(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("range over a slice should compile into lua", t, func() {

		code := `a:=[]int{1,2,3}; func hmm() { for k, v := range a { println(k," ",v) } }`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
a=_giSlice{[0]=1,2,3};
hmm = function() for k, v in pairs(a) do print(k, " ", v);  end end;`)
	})
}

func Test012KeyOnlySliceRangeForLoop(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("key only range over a slice should compile into lua", t, func() {

		code := `a:=[]int{1,2,3}; func hmm() { for i := range a { println(i, a[i]) } }`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
a=_giSlice{[0]=1,2,3};
hmm = function() for i, _ in pairs(a) do print(i, _getRangeCheck(a, i)); end end;`)
	})
}

func Test014SetAStringSliceToEmptyString(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("setting a string slice element should compile into lua", t, func() {

		code := `b := []string{"hi","gophers!"}; b[0]=""`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `b=_giSlice{[0]="hi","gophers!"}; _setRangeCheck(b, 0, "");`)
	})
}

func Test015LenOfSlice(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("len(x) where `x` is a slice should compile", t, func() {

		code := `x := []string{"hi","gophers!"}; bb := len(x)`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `x=_giSlice{[0]="hi","gophers!"}; bb = #x;`)
	})
}

// big L variable name

// print(len(x))

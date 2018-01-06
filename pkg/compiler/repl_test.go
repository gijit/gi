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

		// at top-level
		code := `for i:=0; i < 10; i++ { i=i+2 }`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
  		i = 0;
  		while (true) do
  			if (not (i < 10)) then break; end
            i = i + 2;
  			i = i + (1);
  		 end
`)

		// inside a func
		code = `a:=5; func hmm() { for i:=0; i < a; i++ { println(i) } }`
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
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `a=_gi_NewSlice("Int",{[0]=1,2,3});`)
	})
}

func Test011MapAndRangeForLoop(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("maps and range for loops should compile into lua", t, func() {

		code := `a:=make(map[int]int); a[1]=10; a[2]=20; func hmm() { for k, v := range a { println(k," ",v) } }`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
a = _gi_NewMap("Int", "Int", {});
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
a=_gi_NewSlice("Int",{[0]=1,2,3});
hmm = function() for k, v in pairs(a) do print(k, " ", v);  end end;`)
	})
}

func Test012KeyOnlySliceRangeForLoop(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("key only range over a slice should compile into lua", t, func() {

		code := `a:=[]int{1,2,3}; func hmm() { for i := range a { println(i, a[i]) } }`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
a=_gi_NewSlice("Int",{[0]=1,2,3});
hmm = function() for i, _ in pairs(a) do print(i, _gi_GetRangeCheck(a, i)); end end;`)
	})
}

func Test013SetAStringSliceToEmptyString(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("setting a string slice element should compile into lua", t, func() {

		code := `b := []string{"hi","gophers!"}; b[0]=""`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `b=_gi_NewSlice("String",{[0]="hi","gophers!"}); _setRangeCheck(b, 0, "");`)
	})
}

func Test014LenOfSlice(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("len(x) where `x` is a slice should compile", t, func() {

		code := `x := []string{"hi","gophers!"}; bb := len(x)`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `x=_gi_NewSlice("String",{[0]="hi","gophers!"}); bb = #x;`)
	})
}

func Test015ArrayCreation(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("creating arrays via x := [3]int{1,2,3} where `x` is a slice should compile", t, func() {

		code := `x := [3]int{1,2,3}; bb := len(x)`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `x=_gi_NewArray({[0]=1,2,3}, "_kindInt", 3); bb = 3;`)

		// and empty array with size 3

		code = `var x [3]int`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `x=_gi_NewArray({}, "_kindInt", 3);`)

		// upper case names too
		code = `LX := len(x)`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `LX = 3;`)

		// printing length
		code = `println(len(x))`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `print(3);`)
	})
}

func Test016MapCreation(t *testing.T) {
	inc := NewIncrState()

	cv.Convey(`creating maps via x := map[int]string{3:"hello", 4:"gophers"} should compile`, t, func() {

		// create using make
		code := `y := make(map[int]string)`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `y=_gi_NewMap("Int", "String", {});`)

		// create with literal
		code = `x := map[int]string{3:"hello", 4:"gophers"}`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `x=_gi_NewMap("Int", "String", {[3]="hello", [4]="gophers"});`)

	})
}

func Test017DeleteFromMap(t *testing.T) {
	inc := NewIncrState()

	cv.Convey(`delete from a map, x := map[int]string{3:"hello", 4:"gophers"}, with delete(x, 3) should remove the key 3 with value "hello"`, t, func() {

		code := `x := map[int]string{3:"hello", 4:"gophers"}`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `x=_gi_NewMap("Int", "String", {[3]="hello", [4]="gophers"});`)
		code = `delete(x, 3)`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `x[3]= nil;`)
	})
}

func Test018ReadFromMap(t *testing.T) {

	cv.Convey(`read a map, x := map[int]string{3:"hello", 4:"gophers"}. reading key 3 should provide the value "hello"`, t, func() {

		// and then eval!
		vm := luajit.Newstate()
		defer vm.Close()
		vm.Openlibs()

		files, err := FetchPrelude(".")
		panicOn(err)
		SetupPrelude(vm, files)

		inc := NewIncrState()

		srcs := []string{`x := map[int]string{3:"hello", 4:"gophers"}`, "x3 := x[3]"}
		expect := []string{`x=_gi_NewMap("Int", "String", {[3]="hello", [4]="gophers"});`, `x3 = x[3];`}
		for i, src := range srcs {
			translation := inc.Tr([]byte(src))
			//pp("go:'%s'  -->  '%s' in lua\n", src, translation)
			//fmt.Printf("go:'%#v'  -->  '%#v' in lua\n", src, translation)
			cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, expect[i])

			err := vm.Loadstring(string(translation))
			panicOn(err)
			err = vm.Pcall(0, 0, 0)
			if err != nil {
				fmt.Printf("error: '%v'\n", err)
				DumpLuaStack(vm)
				vm.Pop(1)
				t.Fatal(err)
			} else {
				DumpLuaStack(vm)
			}
			//fmt.Printf("v back = '%#v'\n", v)
		}
		mustLuaString(vm, "x3", "hello")
	})
}

func Test018ReadFromSlice(t *testing.T) {

	cv.Convey(`read a slice, x := []int{3, 4}; reading pos/index 0 should provide the value 3`, t, func() {

		vm := luajit.Newstate()
		defer vm.Close()
		vm.Openlibs()

		files, err := FetchPrelude(".")
		panicOn(err)
		SetupPrelude(vm, files)

		inc := NewIncrState()

		srcs := []string{`x := []int{3, 4}`, "x3 := x[0]"}
		expect := []string{`x=_gi_NewSlice("Int", {[0]=3, 4});`, `x3 = _gi_GetRangeCheck(x,0);`}
		for i, src := range srcs {
			translation := inc.Tr([]byte(src))
			pp("go:'%s'  -->  '%s' in lua\n", src, translation)
			//fmt.Printf("go:'%#v'  -->  '%#v' in lua\n", src, translation)
			cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, expect[i])

			err := vm.Loadstring(string(translation))
			panicOn(err)
			err = vm.Pcall(0, 0, 0)
			if err != nil {
				fmt.Printf("error: '%v'\n", err)
				DumpLuaStack(vm)
				vm.Pop(1)
				t.Fatal(err)
			} else {
				DumpLuaStack(vm)
			}
			//fmt.Printf("v back = '%#v'\n", v)
		}
		vm.Getglobal("x3")
		top := vm.Gettop()
		value_int := vm.Tointeger(top)

		pp("value_int=%v", value_int)
		if value_int != 3 {
			panic(fmt.Sprintf("expected `3`, got %v", value_int))
		}

	})
}

func Test019TopLevelScope(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("top level numerical for loops should be able to refer to other top level variables", t, func() {

		// at top-level
		code := `j:=5; for i:=0; i < 3; i++ { j++ }`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
        j = 5;
  		i = 0;
  		while (true) do
  			if (not (i < 3)) then break; end
            j = j + (1);
  			i = i + (1);
  		 end
`)

	})
}

func Test020StructTypeDeclarations(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("declaring a struct with `type A struct{}` should compile and pass type checking, but produce no lua code (all retained in the type checker, no instantiated value is produced here.", t, func() {

		code := `type A struct{}`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `

`)

	})
}

func Test021StructTypeValues(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("Given `type A struct{}`, when `var a A` is declared, a struct value should be compiled on the lua back end.", t, func() {

		code := `type A struct{}`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
`)
		code = `var a A`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
a=_gi_NewStruct("A",{},{});
`)

	})
}

func Test022StructTypeValues(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("Given `type A struct{ B int }`, when `var a A` is declared, a struct value should be compiled on the lua back end.", t, func() {

		code := `type A struct{ B int}`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
`)
		code = `var a = A{B:43}`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
a=_gi_NewStruct("A",{["B"]="int"},{["B"]=43});
`)

	})
}

/*
func Test023CopyingStructValues(t *testing.T) {
	inc := NewIncrState()

	cv.Convey("Given `type A struct{ B int }`, when `var a =A{B:23}` and then `cp := a; cp.B = 78` then a.B should still be 23 because a full copy/clone should have been made of a during the `cp := a` operation.", t, func() {

		code := `type A struct{ B int}; var a = A{B:23}`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
a=_gi_NewStruct("A",{["B"]="int"},{["B"]=23});
`)
		code = `cp := a`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
cp=_gi_Clone(a);
`)
	})
}
*/

// a, b, c := 1,2,3
func Test024MultipleAssignment(t *testing.T) {

	cv.Convey("Multiple assignment, a, b, c := 1,2,3 should work", t, func() {

		src := `a, b, c := 1,2,3`
		//cv.So(string(inc.Tr([]byte(src))), cv.ShouldMatchModuloWhiteSpace, `a, b, c = 1, 2, 3;`)

		// verify that it happens correctly.
		vm := luajit.Newstate()
		defer vm.Close()
		vm.Openlibs()

		inc := NewIncrState()

		translation := inc.Tr([]byte(src))
		//pp("go:'%s'  -->  '%s' in lua\n", src, translation)
		//fmt.Printf("go:'%#v'  -->  '%#v' in lua\n", src, translation)

		err := vm.Loadstring(string(translation))
		panicOn(err)
		err = vm.Pcall(0, 0, 0)
		panicOn(err)
		DumpLuaStack(vm)

		mustLuaInt(vm, "a", 1)
		mustLuaInt(vm, "b", 2)
		mustLuaInt(vm, "c", 3)
	})
}

func mustLuaInt(vm *luajit.State, varname string, expect int) {

	vm.Getglobal(varname)
	top := vm.Gettop()
	value_int := vm.Tointeger(top)

	pp("value_int=%v", value_int)
	if value_int != expect {
		panic(fmt.Sprintf("expected %v, got %v for '%v'", expect, value_int, varname))
	}
}

func mustLuaString(vm *luajit.State, varname string, expect string) {

	vm.Getglobal(varname)
	top := vm.Gettop()
	value_string := vm.Tostring(top)

	pp("value_string=%v", value_string)
	if value_string != expect {
		panic(fmt.Sprintf("expected %v, got %v for '%v'", expect, value_string, varname))
	}
}

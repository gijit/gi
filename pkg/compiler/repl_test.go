package compiler

import (
	"fmt"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
	luajit "github.com/glycerine/golua/lua"
)

func LoadAndRunTestHelper(t *testing.T, vm *luajit.State, translation []byte) {
	trans := string(translation)
	interr := vm.LoadString(trans)
	if interr != 0 {
		DumpLuaStack(vm)
		t.Fatalf(`could not LoadString("%s")`, trans)
	}
	err := vm.Call(0, 0)
	if err != nil {
		fmt.Printf("error: '%v'\n", err)
		DumpLuaStack(vm)
		vm.Pop(1)
		panic(err)
	}
}

func Test001LuaTranslation(t *testing.T) {

	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("assignment", t, func() {
		cv.So(string(inc.Tr([]byte("a := 10;"))), cv.ShouldMatchModuloWhiteSpace, "a = 10LL;")
		pp("GOOD: past 1st")

		cv.So(string(inc.Tr([]byte("func adder(a, b int) int { return a + b};  sum1 := adder(5,5)"))), cv.ShouldMatchModuloWhiteSpace,
			`adder = function(a, b) 
				         return a + b;
                     end;
			         sum1 = adder(5LL,5LL);`)

		pp("GOOD: past 2nd")

		cv.So(string(inc.Tr([]byte("sum2 := adder(a,a)"))), cv.ShouldMatchModuloWhiteSpace,
			`sum2 = adder(a, a);`)
		pp("GOOD: past 3rd")

		cv.So(string(inc.Tr([]byte("func multiplier(a, b int) int { return a * b};  prod1 := multiplier(5,5)"))), cv.ShouldMatchModuloWhiteSpace,
			`multiplier = function(a, b) 
				         return (a * b);
                     end;
			         prod1 = multiplier(5LL,5LL);`)

	})
}

func Test002LuaEvalIncremental(t *testing.T) {

	// and then eval!
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	srcs := []string{"a := 10;", "func adder(a, b int) int { return a + b}; ", "sum := adder(a,a);"}
	for _, src := range srcs {
		translation := inc.Tr([]byte(src))
		//pp("go:'%s'  -->  '%s' in lua\n", src, translation)
		//fmt.Printf("go:'%#v'  -->  '%#v' in lua\n", src, translation)

		LoadAndRunTestHelper(t, vm, translation)
		//fmt.Printf("v back = '%#v'\n", v)
	}
	LuaMustInt64(vm, "sum", 20)
}

// func Test003ImportsAtRepl(t *testing.T) {
// 	inc := NewIncrState()

// 	cv.Convey("imports", t, func() {
// 		cv.So(string(inc.Tr([]byte(`import "fmt"; fmt.Printf("hello world!")`))), cv.ShouldMatchModuloWhiteSpace, "")
// 		pp("GOOD: past 1st import")
// 	})
// }

func Test004ExpressionsAtRepl(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("expressions alone at top level", t, func() {
		cv.So(string(inc.Tr([]byte(`a:=10;`))), cv.ShouldMatchModuloWhiteSpace, `a=10LL;`)
	})
}

func Test005BacktickStringsToLua(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("Go backtick delimited strings should translate to Lua", t, func() {
		cv.So(string(inc.Tr([]byte("s:=`\n\n\"hello\"\n\n`"))), cv.ShouldMatchModuloWhiteSpace, `s = "\n\n\"hello\"\n\n";`)
		cv.So(string(inc.Tr([]byte("r:=`\n\n]]\"hello\"\n\n`"))), cv.ShouldMatchModuloWhiteSpace, `r = "\n\n]]\"hello\"\n\n";`)
	})
}

func Test006RedefinitionOfVariablesAllowed(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("At the repl, `a:=1; a:=2;` is allowed. We disable the traditional Go re-definition checks at the REPL", t, func() {
		cv.So(string(inc.Tr([]byte("a:=1; a:=2;"))), cv.ShouldMatchModuloWhiteSpace, `a=1LL; a=2LL;`)

		// and in separate calls:
		cv.So(string(inc.Tr([]byte("r:=`\n\n]]\"hello\"\n\n`"))), cv.ShouldMatchModuloWhiteSpace, `r = "\n\n]]\"hello\"\n\n";`)
		// and redefinition of r should work:
		cv.So(string(inc.Tr([]byte("r:=`a new definition`"))), cv.ShouldMatchModuloWhiteSpace, `r = "a new definition";`)
	})
}

func Test007SettingPreviouslyDefinedVariables(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("At the repl, in separate commands`a:=1; a=2;` sets a to 2", t, func() {

		// and in separate calls, where the second call sets the earlier variable.
		cv.So(string(inc.Tr([]byte("a:=1"))), cv.ShouldMatchModuloWhiteSpace, `a=1LL;`)
		cv.So(string(inc.Tr([]byte("b:=2"))), cv.ShouldMatchModuloWhiteSpace, `b=2LL;`)
		// and redefinition of a should work:
		pp("starting on a=2;")
		cv.So(string(inc.Tr([]byte("a=2;"))), cv.ShouldMatchModuloWhiteSpace, `a=2LL;`)
	})
}

func Test008IfThenElseInAFunction(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("if then else within a closure/function should compile into lua", t, func() {

		code := `a:=20; func hmm() { if a > 30 { println("over 30") } else {println("under or at 30")} }`
		// and in separate calls, where the second call sets the earlier variable.
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `a=20LL; hmm = function() if (a > 30LL) then print("over 30"); else print("under or at 30"); end end;`)
	})
}

func Test009NumericalForLoop(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("numerical for loops should compile into lua", t, func() {

		// at top-level
		code := `for i:=0; i < 10; i++ { i=i+2 }`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
  		i = 0LL;
  		while (true) do
  			if (not (i < 10LL)) then break; end
            i = i + 2LL;
  			i = i + (1LL);
  		 end
`)

		// inside a func
		code = `a:=5; func hmm() { for i:=0; i < a; i++ { println(i) } }`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `a=5LL;
  	hmm = function() 
  		local i = 0LL;
  		while (true) do
  			if (not (i < a)) then break; end
  			print(i);
  			i = i + (1LL);
  		 end
 	 end;
`)

	})
}

func Test010Slice(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("slice literal should compile into lua", t, func() {

		code := `a:=[]int{1,2,3}`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
	__type__anon_sliceType = __sliceType(__type__int);  -- 'IMMEDIATE' anon type printing.
a=__type__anon_sliceType({[0]=1LL,2LL,3LL});`)
	})
}

func Test011MapAndRangeForLoop(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("maps and range for loops should compile into lua", t, func() {

		code := `a:=make(map[int]int); a[1]=10; a[2]=20; func hmm() { for k, v := range a { println(k," ",v) } }`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
a = {};
a["1LL"] = 10LL;
a["2LL"] = 20LL;
hmm = function() for k, v in pairs(a) do print(k, " ", v);  end end;`)
	})
}

func Test012SliceRangeForLoop(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("range over a slice should compile into lua", t, func() {

		code := `a:=[]int{1,2,3}; func hmm() { for k, v := range a { println(k," ",v) } }`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
__type__anon_sliceType = __sliceType(__type__int); -- 'IMMEDIATE' anon type printing.  
a = __type__anon_sliceType({[0]=1LL, 2LL, 3LL});
hmm = function() for k, v in pairs(a) do print(k, " ", v);  end end;`)
	})
}

func Test012KeyOnlySliceRangeForLoop(t *testing.T) {

	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("key only range over a slice should compile into lua", t, func() {

		code := `a:=[]int{1,2,3}; func hmm() { for i := range a { println(i, a[i]) } }`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
  	__type__anon_sliceType = __sliceType(__type__int); -- 'IMMEDIATE' anon type printing.
  	a = __type__anon_sliceType({[0]=1LL, 2LL, 3LL});
hmm = function() for i, _ in pairs(a) do print(i, _gi_GetRangeCheck(a, i)); end end;`)
	})
}

func Test013SetAStringSliceToEmptyString(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("setting a string slice element should compile into lua", t, func() {

		code := `b := []string{"hi","gophers!"}; b[0]=""`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
__type__anon_sliceType = __sliceType(__type__string); -- 'IMMEDIATE' anon type printing.  
b = __type__anon_sliceType({[0]="hi", "gophers!"});
_gi_SetRangeCheck(b, 0, "");
`)
	})
}

func Test014LenOfSlice(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("len(x) where `x` is a slice should compile", t, func() {

		code := `x := []string{"hi","gophers!"}; bb := len(x)`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
  	__type__anon_sliceType = __sliceType(__type__string); -- 'IMMEDIATE' anon type printing.
  
  	x = __type__anon_sliceType({[0]="hi", "gophers!"});

    bb = #x;`)
	})
}

func Test015ArrayCreation(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("creating arrays via x := [3]int{1,2,3} where `x` is a slice should compile", t, func() {

		code := `x := [3]int{1,2,3}; bb := len(x)`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
__type__anon_arrayType = __arrayType(__type__int, 3); -- 'IMMEDIATE' anon type printing.
x = __type__anon_arrayType({[0]=1LL, 2LL, 3LL});
bb = 3LL;`)

		// and empty array with size 3

		// type already declared above, so will be reused.
		code = `var x [3]int`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
x = __type__anon_arrayType.zero();
`)

		// upper case names too
		code = `LX := len(x)`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `LX = 3LL;`)

		// printing length
		code = `println(len(x))`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `print(3LL);`)
	})
}

func Test015_5_ArrayCreation(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("creating arrays via x := [3]int{1,2,3} where `x` is a slice should compile", t, func() {

		// and empty array with size 3

		code := `var x [3]int`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
__type__anon_arrayType = __arrayType(__type__int, 3); -- 'IMMEDIATE' anon type printing.
x = __type__anon_arrayType();
`)
	})
}

func Test016MapCreation(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey(`creating maps via x := map[int]string{3:"hello", 4:"gophers"} should compile`, t, func() {

		// create using make
		code := `y := make(map[int]string)`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `y=_gi_NewMap("int", "string", {});`)

		// create with literal
		code = `x := map[int]string{3:"hello", 4:"gophers"}`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `x=_gi_NewMap("int", "string", {["3LL"]="hello", ["4LL"]="gophers"});`)

	})
}

func Test017DeleteFromMap(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey(`delete from a map, x := map[int]string{3:"hello", 4:"gophers"}, with delete(x, 3) should remove the key 3 with value "hello"`, t, func() {

		code := `x := map[int]string{3:"hello", 4:"gophers"}`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `x=_gi_NewMap("int", "string", {["3LL"]="hello", ["4LL"]="gophers"});`)
		code = `delete(x, 3)`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `x("delete",3LL);`)
	})
}

func Test018ReadFromMap(t *testing.T) {

	cv.Convey(`read a map, x := map[int]string{3:"hello", 4:"gophers"}. reading key 3 should provide the value "hello"`, t, func() {

		// and then eval!
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		srcs := []string{`x := map[int]string{3:"hello", 4:"gophers"}`, "x3 := x[3]"}
		expect := []string{`x=_gi_NewMap("int", "string", {["3LL"]="hello", ["4LL"]="gophers"});`, `x3 = x('get',"3LL", "");`}
		for i, src := range srcs {
			translation := inc.Tr([]byte(src))
			//pp("go:'%s'  -->  '%s' in lua\n", src, translation)
			fmt.Printf("go:'%s'  -->  '%s' in lua\n", string(src), string(translation))
			cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, expect[i])

			LoadAndRunTestHelper(t, vm, translation)
			//fmt.Printf("v back = '%#v'\n", v)
		}
		LuaMustString(vm, "x3", "hello")
	})
}

func Test018ReadFromSlice(t *testing.T) {

	cv.Convey(`read a slice, x := []int{3, 4}; reading pos/index 0 should provide the value 3`, t, func() {

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		srcs := []string{`x := []int{3, 4}`, "x3 := x[0]"}
		expect := []string{`x=_gi_NewSlice("int", {[0]=3LL, 4LL},0LL);`, `x3 = _gi_GetRangeCheck(x,0);`}
		for i, src := range srcs {
			translation := inc.Tr([]byte(src))
			pp("go:'%s'  -->  '%s' in lua\n", src, translation)
			//fmt.Printf("go:'%#v'  -->  '%#v' in lua\n", src, translation)
			cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, expect[i])

			LoadAndRunTestHelper(t, vm, translation)
			//fmt.Printf("v back = '%#v'\n", v)
		}
		LuaMustInt64(vm, "x3", 3)
	})
}

func Test019TopLevelScope(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("top level numerical for loops should be able to refer to other top level variables", t, func() {

		// at top-level
		code := `j:=5; for i:=0; i < 3; i++ { j++ }`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
        j = 5LL;
  		i = 0LL;
  		while (true) do
  			if (not (i < 3LL)) then break; end
            j = j + (1LL);
  			i = i + (1LL);
  		 end
`)

	})
}

func Test020StructTypeDeclarations(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("declaring a struct with `type A struct{}` should compile and pass type checking, and register a prototype", t, func() {

		code := `type A struct{}`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldStartWithModuloWhiteSpace, `
       __type__A = __gi_NewType(0, __gi_kind_Struct, "main", "A", "main.A", true, "main",
`)

	})
}

func Test021StructTypeValues(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("Given `type A struct{}`, when `var a A` is declared, a struct value should be compiled on the lua back end.", t, func() {

		code := `type A struct{}`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldStartWithModuloWhiteSpace, `
       __type__A = __gi_NewType(0, __gi_kind_Struct, "main", "A", "main.A", true, "main",
`)
		code = `var a A`
		//cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `a=__reg:NewInstance("A",{});`)
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `a = __type__A.__ptr({}, nil);`)

	})
}

func Test022StructTypeValues(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("Given `type A struct{ B int }`, when `var a A` is declared, a struct value should be compiled on the lua back end.", t, func() {

		code := `type A struct{ B int}`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldStartWithModuloWhiteSpace, `
   __type__A = __gi_NewType(0, __gi_kind_Struct, "main", "A", "main.A", true, "main", true
`)
		code = `var a = A{B:43}`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
a = __type__A.__ptr({}, 43LL);
`)
		// a=__reg:NewInstance("A",{["B"]=43LL});

	})
}

// come back to this.
/*
func Test023CopyingStructValues(t *testing.T) {
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm,nil)


	cv.Convey("Given `type A struct{ B int }`, when `var a =A{B:23}` and then `cp := a; cp.B = 78` then a.B should still be 23 because a full copy/clone should have been made of a during the `cp := a` operation.", t, func() {

		code := `type A struct{ B int}; var a = A{B:23}`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
a=__reg:NewInstance("A",{["B"]=23});
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

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		//cv.So(string(inc.Tr([]byte(src))), cv.ShouldMatchModuloWhiteSpace, `a, b, c = 1, 2, 3;`)

		// verify that it happens correctly.
		translation := inc.Tr([]byte(src))
		//pp("go:'%s'  -->  '%s' in lua\n", src, translation)
		//fmt.Printf("go:'%#v'  -->  '%#v' in lua\n", src, translation)

		LoadAndRunTestHelper(t, vm, translation)

		LuaMustInt64(vm, "a", 1)
		LuaMustInt64(vm, "b", 2)
		LuaMustInt64(vm, "c", 3)
	})
}

func Test025ComplexNumbers(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("a := 6.67428e-11i should compile, since luajit has builtin support for complex number syntax", t, func() {

		code := `a := 6.67428e-11i`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `
	a = 0+6.67428e-11i;`)
	})
}

func Test026LenOfString(t *testing.T) {

	cv.Convey(`a := "hi"; b := len(a); should return b of 2, so the len() on a string works.`, t, func() {

		code := `a := "hi"; b :=len(a);`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))

		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, `
	a = "hi"; b = #a;`)

		// and verify that it happens correctly.

		LoadAndRunTestHelper(t, vm, translation)
		LuaMustInt(vm, "b", 2)
	})
}

func Test029StructMethods(t *testing.T) {

	cv.Convey(`verify that 027 actually executes correctly on the repl: a simple method call through an interface to a struct method should translate`, t, func() {

		code := `
type Dog interface {
    Write(with string) string
}

type Beagle struct{
    word string
}

func (b *Beagle) Write(with string) string {
    return b.word + ":it was a dark and stormy night, " + with
}

var snoopy Dog = &Beagle{word:"hiya"}

book := snoopy.Write("with a pen")`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))

		/*
					cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace,
						`
			        Dog = __reg:RegisterInterface("Dog","main","main");
			        Beagle = __reg:RegisterStruct("Beagle","main","main");

				    function Beagle:Write(with)
			            b = self;
			  		    return b.word .. ":it was a dark and stormy night, " .. with;
			     	end;
			        __reg:AddMethod("struct", "Beagle", "Write", Beagle.Write)
			        snoopy = __reg:NewInstance("Beagle",{["word"]="hiya"});

			  	    _r = snoopy:Write("with a pen");
			  	    book = _r;
			`)
		*/

		// and verify that it happens correctly

		LuaRunAndReport(vm, string(translation))
		LuaMustString(vm, "book", "hiya:it was a dark and stormy night, with a pen")
	})
}

func Test030MethodRedefinitionAllowed(t *testing.T) {

	cv.Convey(`methods can be re-defined at the repl`, t, func() {

		code := `
type Beagle struct{
    word string
}
func (b *Beagle) Write(with string) string {
    return b.word + ":it was a sunny day at the beach, " + with
}
func (b *Beagle) Write(with string) string {
    return b.word + ":it was a dark and stormy night, " + with
}
var snoopy = &Beagle{word:"hiya"}
book := snoopy.Write("with a pen")
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))
		LuaMustString(vm, "book", "hiya:it was a dark and stormy night, with a pen")
	})
}

func Test031ValueOkDoubleReturnFromMapQuery(t *testing.T) {

	cv.Convey(`given m := map[int]int{1:1}; then a, ok := m[0] should provide ok false and a = 0, while a, ok := m[1]; should provide ok true and a = 1.`, t, func() {

		code := `
m := map[int]int{1:1}
a0, ok0 := m[0]
a1, ok1 := m[1]
alone0 := m[0]
alone1 := m[1]
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))
		LuaMustBool(vm, "ok1", true) // fail here
		LuaMustInt64(vm, "a1", 1)
		LuaMustBool(vm, "ok0", false)
		LuaMustInt64(vm, "a0", 0)

		LuaMustInt64(vm, "alone0", 0)
		LuaMustInt64(vm, "alone1", 1)
	})
}

func Test032DeleteOnMapsAndMapsCanStoreNil(t *testing.T) {

	cv.Convey(`delete on maps should work. maps should allow nil as both key and value, and len should still be correct`, t, func() {

		code := `
m := map[int]int{1:1, 2:2}
len1 := len(m)
delete(m, 1)
len2 := len(m)
delete(m, 2)
len3 := len(m)
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))
		LuaMustInt(vm, "len1", 2)
		LuaMustInt(vm, "len2", 1)
		LuaMustInt(vm, "len3", 0)
	})
}

func Test036Println(t *testing.T) {

	cv.Convey(`println should translate to print so basic reporting can be done without fmt`, t, func() {

		code := `
println("hello")
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))

		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace,
			`print("hello");`)

		/*
			// and verify that it happens correctly
			vm := luajit.NewState()
			defer vm.Close()
			vm.OpenLibs()
			files, err := FetchPrelude(".")
			panicOn(err)
			LuaDoFiles(vm, files)

			LuaRunAndReport(vm, string(translation))
			LuaMustInt(vm, "len1", 2)
			LuaMustInt(vm, "len2", 1)
			LuaMustInt(vm, "len3", 0)
		*/
	})
}

func Test037Println(t *testing.T) {

	cv.Convey(`named return values`, t, func() {

		code := `
func f() (a,b,c int, d string) {
  a = 1
  b = 2
  c = 3
  d = "hi"
  return
};
x,y,z,s := f()
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))

		//		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace,
		//			``)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "x", 1)
		LuaMustInt64(vm, "y", 2)
		LuaMustInt64(vm, "z", 3)
		LuaMustString(vm, "s", "hi")

	})
}

func Test038Switch(t *testing.T) {

	cv.Convey(`switch with value at top should compile`, t, func() {

		code := `
a := 7;
b := 2;
c := 0
switch b {
case 1:
  c = a*1
case 2:
  c = a*10
case 3:
  c = a*100
default:
  c = -1
}
`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))

		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace,
			`
	a = 7LL;
  	b = 2LL;
  	c = 0LL;
  	_1 = b;
  	if (_1 == (1LL)) then 
  		c = (a * 1LL);
  	 elseif (_1 == (2LL)) then 
  		c = (a * 10LL);
  	 elseif (_1 == (3LL)) then 
  		c = (a * 100LL);
  	 else  
  		c = -1LL;
  	 end 
`)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "c", 70)

	})
}

func Test039SwitchInFunction(t *testing.T) {

	cv.Convey(`switch statement inside a function should compile`, t, func() {

		code := `func f() int {
a := 7;
b := 2;
c := 0
switch b {
case 1:
  c = a*1
case 2:
  c = a*10
case 3:
  c = a*100
default:
  c = -1
}
return c}
myc := f()
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "myc", 70)

	})
}

func Test040Switch(t *testing.T) {

	cv.Convey(`switch with no value at top should compile`, t, func() {

		code := `
a := 7;
b := 2;
c := 0
switch {
case b == 1:
  c = a*1
case b == 2:
  c = a*10
case b == 3:
  c = a*100
default:
  c = -1
}
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))

		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace,
			`
	a = 7LL;
  	b = 2LL;
  	c = 0LL;
  	if ((b == 1LL)) then 
  		c = (a * 1LL);
  	 elseif ((b == 2LL)) then 
  		c = (a * 10LL);
  	 elseif ((b == 3LL)) then 
  		c = (a * 100LL);
  	 else 
  		c = -1LL;
  	 end 
`)

		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "c", 70)

	})
}

func Test042LenAtRepl(t *testing.T) {

	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey(`a := []int{3}; len(a)' at the repl, len(a) should give us 1, so it should get wrapped in a print()`, t, func() {

		code := `a := []int{3}; len(a)`
		cv.So(string(inc.Tr([]byte(code))), cv.ShouldMatchModuloWhiteSpace, `a = _gi_NewSlice("int",{[0]=3LL}, 0LL); print(#a);`)
	})
}

func Test043IntegerDivideByZero(t *testing.T) {

	cv.Convey(`integers divided by zero or taken modulo zero should produce an error`, t, func() {

		code := `
a := 0;
b := 1/a;
m := 1%a
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))

		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace,
			`
	a = 0LL;
    b = __integerByZeroCheck(1LL / a);
    m = __integerByZeroCheck(1LL % a);
`)

		codeWithCatch := `
c:=0
func f() {
    defer func() {
       // divide by zero should have fired a panic
       if recover() != nil {
           c = 1
       }
    }()
	a := 0;
    b := 1 / a
    _ = b
}
f();
`
		translation = inc.Tr([]byte(codeWithCatch))

		fmt.Printf("\n translation='%s'\n", translation)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		// check for exception
		LuaMustInt64(vm, "c", 1)

	})
}

func Test069MethodRedefinitionAllowed(t *testing.T) {

	cv.Convey(`methods can be re-defined, including changing their signature, at the repl`, t, func() {

		code := `
 type S struct { a int }
 func (s *S) inc(b int) int { return s.a + b}

 // new signature in addition to new body: so we recognize fresh/old
 func (s *S) inc(b, c int) int { s.a++; s.a += b + c; return s.a }
 var s S
 a := s.inc(3, 4)
`

		//   Line 1074: - where error? err = '8:13: too many arguments'

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "a", 8)
	})
}

func Test070VariablesInsideFunctionsAreLocal(t *testing.T) {

	cv.Convey(`given func f() int { a := 1; return a}, the variable "a" should be local and should not introduce a global binding or overwrite a global variable named "a"`, t, func() {

		code := `
a := 2
func f() int { a := 1; return a }
f()
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace,
			`
a = 2LL;
f = function()
  local a = 1LL;
  return a;
end;
f();
`)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "a", 2)
	})
}

func Test028CopyAStruct(t *testing.T) {

	cv.Convey(`copy a struct value`, t, func() {

		code := `
type Beagle struct{ 
    word string 
    a, b, c int
} 
var snoopy = Beagle{word:"hiya", b: 2}
denver := snoopy
d := denver.b
w := denver.word
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))

		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "d", 2)
		LuaMustString(vm, "w", "hiya")
		cv.So(true, cv.ShouldBeTrue)
	})
}

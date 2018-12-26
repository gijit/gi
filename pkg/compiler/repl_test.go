package compiler

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	//"github.com/gijit/gi/pkg/verb"
	cv "github.com/glycerine/goconvey/convey"
	//luajit "github.com/glycerine/golua/lua"
)

func init() {
	defaultTestMode = true
}

func init() {
	runtime.LockOSThread()
}

func TestMain(m *testing.M) {
	if reserveMainThread {
		go func() {
			os.Exit(m.Run())
		}()
		MainCThread()
	} else {
		os.Exit(m.Run())
	}
}

var matchesLuaSrc = cv.ShouldMatchModuloWhiteSpaceAndLuaComments
var startsWithLuaSrc = cv.ShouldStartWithModuloWhiteSpaceAndLuaComments

func LoadAndRunTestHelper(t *testing.T, lvm *LuaVm, translation []byte) {
	trans := string(translation)
	err := LuaRun(lvm, trans, true)
	if err != nil {
		fmt.Printf("error from LuaRun:\n%v\n", err)
		t.Fatalf(`could not LuaRun("%s")`, trans)
	}
}

func Test001LuaTranslation(t *testing.T) {

	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("assignment", t, func() {
		by, err := inc.Tr([]byte("a := 10;"))
		panicOn(err)
		LuaRunAndReport(vm, string(by))
		LuaMustInt64(vm, "a", 10)

		pp("GOOD: past 1st")

		by, err = inc.Tr([]byte("func adder(a, b int) int { return a + b};  sum1 := adder(5,5)"))
		panicOn(err)
		lua := string(by)

		LuaRunAndReport(vm, lua)
		LuaMustInt64(vm, "sum1", 10)

		pp("GOOD: past 2nd")

		by, err = inc.Tr([]byte("sum2 := adder(a,a)"))
		panicOn(err)

		LuaRunAndReport(vm, string(by))
		LuaMustInt64(vm, "sum2", 20)
		pp("GOOD: past 3rd")

		by, err = inc.Tr([]byte("func multiplier(a, b int) int { return a * b};  prod1 := multiplier(5,5)"))
		panicOn(err)
		LuaRunAndReport(vm, string(by))
		LuaMustInt64(vm, "prod1", 25)

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
		translation := inc.trMust([]byte(src))
		//pp("go:'%s'  -->  '%s' in lua\n", src, translation)
		//fmt.Printf("go:'%#v'  -->  '%#v' in lua\n", src, translation)

		LoadAndRunTestHelper(t, vm, translation)
		//fmt.Printf("v back = '%#v'\n", v)
	}
	LuaMustInt64(vm, "sum", 20)
}

func Test004ExpressionsAtRepl(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("expressions alone at top level", t, func() {
		cv.So(string(inc.trMust([]byte(`a:=10;`))), matchesLuaSrc, `a=10LL;`)
	})
}

func Test005BacktickStringsToLua(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("Go backtick delimited strings should translate to Lua", t, func() {
		cv.So(string(inc.trMust([]byte("s:=`\n\n\"hello\"\n\n`"))), matchesLuaSrc, `s = "\n\n\"hello\"\n\n";`)
		cv.So(string(inc.trMust([]byte("r:=`\n\n]]\"hello\"\n\n`"))), matchesLuaSrc, `r = "\n\n]]\"hello\"\n\n";`)
	})
}

// we've had to disable this to get the type system working, for now.
func Test006RedefinitionOfVariablesAllowed(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("At the repl, `a:=1; a:=2;` is allowed. We disable the traditional Go re-definition checks at the REPL", t, func() {
		cv.So(string(inc.trMust([]byte("a:=1; a:=2;"))), matchesLuaSrc, `a=1LL; a=2LL;`)

		// and in separate calls:
		cv.So(string(inc.trMust([]byte("r:=`\n\n]]\"hello\"\n\n`"))), matchesLuaSrc, `r = "\n\n]]\"hello\"\n\n";`)
		// and redefinition of r should work:
		cv.So(string(inc.trMust([]byte("r:=`a new definition`"))), matchesLuaSrc, `r = "a new definition";`)
	})
}

func Test007SettingPreviouslyDefinedVariables(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("At the repl, in separate commands`a:=1; a=2;` sets a to 2", t, func() {

		// and in separate calls, where the second call sets the earlier variable.
		cv.So(string(inc.trMust([]byte("a:=1"))), matchesLuaSrc, `a=1LL;`)
		cv.So(string(inc.trMust([]byte("b:=2"))), matchesLuaSrc, `b=2LL;`)
		// and redefinition of a should work:
		pp("starting on a=2;")
		cv.So(string(inc.trMust([]byte("a=2;"))), matchesLuaSrc, `a=2LL;`)
	})
}

func Test008IfThenElseInAFunction(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("if then else within a closure/function should compile into lua", t, func() {

		code := `a:=20; over:=0; func hmm() { if a > 30 { println("over 30"); over=1; } else {println("under or at 30"); over=-1; } }; hmm();`
		// and in separate calls, where the second call sets the earlier variable.
		LuaRunAndReport(vm, string(inc.trMust([]byte(code))))
		LuaMustInt64(vm, "over", -1)
	})
}

func Test009NumericalForLoop(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("numerical for loops should compile into lua", t, func() {

		// at top-level
		code := `a:= 0; for i:=0; i < 4; i++ { a+=i }`
		lua := string(inc.trMust([]byte(code)))
		LuaRunAndReport(vm, lua)
		LuaMustInt64(vm, "a", 6)

		// inside a func
		code = `a:=5; b:=0; func hmm() { for i:=0; i < a; i++ { println(i); b-=i } }; hmm();`
		lua = string(inc.trMust([]byte(code)))
		LuaRunAndReport(vm, lua)
		LuaMustInt64(vm, "b", -10)
		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test010Slice(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("slice literal should compile into lua", t, func() {

		code := `a:=[]int{1,2,3}`
		cv.So(string(inc.trMust([]byte(code))), matchesLuaSrc, `
	__type__.anon_sliceType = __sliceType(__type__.int);
a=__type__.anon_sliceType({[0]=1LL,2LL,3LL});`)
	})
}

func Test011MapAndRangeForLoop(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("maps and range for loops should compile into lua", t, func() {

		code := `a:=make(map[int]int); a[1]=10; a[2]=20; ktot:=0; vtot:=0; func hmm() { for k, v := range a { ktot+=k; vtot+=v; } }; hmm();`
		lua := string(inc.trMust([]byte(code)))
		pp("lua='%s'", lua)
		LuaRunAndReport(vm, lua)
		LuaMustInt64(vm, "ktot", 3)
		LuaMustInt64(vm, "vtot", 30)

	})
}

func Test012SliceRangeForLoop(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("range over a slice should compile into lua", t, func() {

		code := `a:=[]int{1,2,3}; tot:=0; func hmm() { for k, v := range a { tot += v; _ = k } }; hmm()`
		lua := string(inc.trMust([]byte(code)))
		fmt.Printf("lua='%s'", lua)
		LuaRunAndReport(vm, lua)
		LuaMustInt64(vm, "tot", 6)
	})
}

func Test012KeyOnlySliceRangeForLoop(t *testing.T) {

	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("key only range over a slice should compile into lua", t, func() {

		code := `a:=[]int{1,2,3}; itot:=0;  for i := range a { itot+=i }`
		lua := string(inc.trMust([]byte(code)))
		fmt.Printf("lua='%s'", lua)
		LuaRunAndReport(vm, lua)
		LuaMustInt64(vm, "itot", 3)
	})
}

func Test012AssignNotDefineRangeForLoop(t *testing.T) {

	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("key only range over a slice should compile into lua", t, func() {

		code := `a:=[]int{1,2,3}; itot:=0; i:=0; for i = range a { itot+=i }`
		lua := string(inc.trMust([]byte(code)))
		fmt.Printf("lua='%s'", lua)
		LuaRunAndReport(vm, lua)
		LuaMustInt64(vm, "itot", 3)
	})
}

func Test012DoubleRangeLoop(t *testing.T) {

	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("nested range loops", t, func() {

		code := `a:=[]int{1,2,3}; b:=[]int{4,5,6}; vtot:=0; for i := range a { for j := range b { vtot += a[i]*b[j] } }`
		lua := string(inc.trMust([]byte(code)))
		fmt.Printf("lua='%s'", lua)
		LuaRunAndReport(vm, lua)
		LuaMustInt64(vm, "vtot", 90)
	})
}

func Test013SetAStringSliceToEmptyString(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("setting a string slice element should compile into lua", t, func() {

		code := `b := []string{"hi","gophers!"}; b[0]=""`
		cv.So(string(inc.trMust([]byte(code))), matchesLuaSrc, `
__type__.anon_sliceType = __sliceType(__type__.string); 
b = __type__.anon_sliceType({[0]="hi", "gophers!"});
__gi_SetRangeCheck(b, 0, "");
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
		cv.So(string(inc.trMust([]byte(code))), matchesLuaSrc, `
  	__type__.anon_sliceType = __sliceType(__type__.string); 
  
  	x = __type__.anon_sliceType({[0]="hi", "gophers!"});

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
		cv.So(string(inc.trMust([]byte(code))), matchesLuaSrc, `
__type__.anon_arrayType = __arrayType(__type__.int, 3); 
x = __type__.anon_arrayType({[0]=1LL, 2LL, 3LL});
bb = 3LL;`)

		// and empty array with size 3

		// type already declared above, so will be reused.
		code = `var x [3]int`
		cv.So(string(inc.trMust([]byte(code))), matchesLuaSrc, `
__type__.anon_arrayType = __arrayType(__type__.int, 3); 
x = __type__.anon_arrayType();
`)

		// upper case names too
		code = `LX := len(x)`
		cv.So(string(inc.trMust([]byte(code))), matchesLuaSrc, `LX = 3LL;`)

		// printing length
		code = `println(len(x))`
		cv.So(string(inc.trMustPre([]byte(code), false)), matchesLuaSrc, `print(3LL);`)
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
		cv.So(string(inc.trMust([]byte(code))), matchesLuaSrc, `
__type__.anon_arrayType = __arrayType(__type__.int, 3); 
x = __type__.anon_arrayType();
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
		cv.So(string(inc.trMust([]byte(code))), matchesLuaSrc, `
	__type__.anon_mapType = __mapType(__type__.int, __type__.string);
	y = __makeMap({}, __type__.int, __type__.string, __type__.anon_mapType);
`)

		// create with literal
		code = `x := map[int]string{3:"hello", 4:"gophers"}`
		cv.So(string(inc.trMust([]byte(code))), matchesLuaSrc, `
  	__type__.anon_mapType = __mapType(__type__.int, __type__.string); 
  	x = __makeMap({[3LL]="hello", [4LL]="gophers"}, __type__.int, __type__.string, __type__.anon_mapType);
`)

	})
}

func Test016bMapCreation(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey(`creating maps new named types should compile`, t, func() {

		code := `type Yumo map[int]string; yesso := Yumo{2: "two"}`
		cv.So(string(inc.trMust([]byte(code))), matchesLuaSrc, `
__type__.Yumo = __newType(8, __kindMap, "main.Yumo", true, "main", true, nil);

__type__.Yumo.init(__type__.int, __type__.string);

__type__.anon_mapType = __mapType(__type__.int, __type__.string); 

yesso = __makeMap({[2LL]="two"}, __type__.int, __type__.string, __type__.Yumo);

`)

	})
}

func Test017DeleteFromMap(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey(`delete from a map, x := map[int]string{3:"hello", 4:"gophers"}, with delete(x, 3) should remove the key 3 with value "hello"`, t, func() {

		code := `x := map[int]string{3:"hello", 4:"gophers"}`
		cv.So(string(inc.trMust([]byte(code))), matchesLuaSrc, `
  	__type__.anon_mapType = __mapType(__type__.int, __type__.string); 
  
  	x = __makeMap({[3LL]="hello", [4LL]="gophers"}, __type__.int, __type__.string, __type__.anon_mapType);
`)
		code = `delete(x, 3)`
		cv.So(string(inc.trMustPre([]byte(code), false)), matchesLuaSrc, `x("delete",3LL);`)
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
		// 	expect := []string{`
		// __type__.anon_mapType = __mapType(__type__.int, __type__.string);
		// x = __makeMap(__type__.int.keyFor, {[3LL]="hello", [4LL]="gophers"}, __type__.int, __type__.string, __type__.anon_mapType);`,
		// 		`   x3 =  x('get', "3LL", "");`}

		for _, src := range srcs {
			translation := inc.trMust([]byte(src))
			//pp("go:'%s'  -->  '%s' in lua\n", src, translation)
			fmt.Printf("go:'%s'  -->  '%s' in lua\n", string(src), string(translation))
			//cv.So(string(translation), matchesLuaSrc, expect[i])

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
		expect := []string{`  	__type__.anon_sliceType = __sliceType(__type__.int); 
  
  	x = __type__.anon_sliceType({[0]=3LL, 4LL});`, `x3 = __gi_GetRangeCheck(x,0);`}
		for i, src := range srcs {
			translation := inc.trMust([]byte(src))
			pp("go:'%s'  -->  '%s' in lua\n", src, translation)
			//fmt.Printf("go:'%#v'  -->  '%#v' in lua\n", src, translation)
			cv.So(string(translation), matchesLuaSrc, expect[i])

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
		lua := string(inc.trMust([]byte(code)))
		LuaRunAndReport(vm, lua)
		LuaMustInt64(vm, "j", 8)

	})
}

func Test020StructTypeDeclarations(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("declaring a struct with `type A struct{}` should compile and pass type checking, and register a prototype", t, func() {

		code := `type A struct{}`
		cv.So(string(inc.trMust([]byte(code))), startsWithLuaSrc, `
	__type__.A = __newType(0, __kindStruct, "main.A", true, "main", true, nil);
  	__type__.A.init("", {});  	
  	 __type__.A.__constructor = function() 
  		 return {}; end;`)

	})
}

func Test021StructTypeValues(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("Given `type A struct{}`, when `var a A` is declared, a struct value should be compiled on the lua back end.", t, func() {

		code := `type A struct{}`
		output := string(inc.trMust([]byte(code)))
		LuaRunAndReport(vm, output)

		code = `var a A`
		output = string(inc.trMust([]byte(code)))
		LuaRunAndReport(vm, output)

		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test022StructTypeValues(t *testing.T) {
	vm, err := NewLuaVmWithPrelude(nil)
	panicOn(err)
	defer vm.Close()
	inc := NewIncrState(vm, nil)

	cv.Convey("Given `type A struct{ B int }`, when `var a A` is declared, a struct value should be compiled on the lua back end.", t, func() {

		code := `type A struct{ B int}`
		output := string(inc.trMust([]byte(code)))
		pp("output = '%s'", output)
		LuaRunAndReport(vm, output)

		code = `var a = A{B:43}; ab := a.B;`
		output = string(inc.trMust([]byte(code)))
		pp("output = '%s'", output)
		LuaRunAndReport(vm, output)

		LuaMustInt64(vm, "ab", 43)
		cv.So(true, cv.ShouldBeTrue)
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
		cv.So(string(inc.trMust([]byte(code))), matchesLuaSrc, `
a=__reg:NewInstance("A",{["B"]=23});
`)
		code = `cp := a`
		cv.So(string(inc.trMust([]byte(code))), matchesLuaSrc, `
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

		//cv.So(string(inc.trMust([]byte(src))), matchesLuaSrc, `a, b, c = 1, 2, 3;`)

		// verify that it happens correctly.
		translation := inc.trMust([]byte(src))
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
		cv.So(string(inc.trMust([]byte(code))), matchesLuaSrc, `
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

		translation := inc.trMust([]byte(code))

		cv.So(string(translation), matchesLuaSrc, `
	a = "hi"; b = #a;`)

		// and verify that it happens correctly.

		LoadAndRunTestHelper(t, vm, translation)
		LuaMustInt(vm, "b", 2)
	})
}

func Test029StructMethods(t *testing.T) {

	cv.Convey(`verify that interface + struct + methods actually executes correctly on the repl: a simple method call through an interface to a struct method should translate`, t, func() {

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

		translation := inc.trMust([]byte(code))
		// and verify that it happens correctly

		LuaRunAndReport(vm, string(translation))
		LuaMustString(vm, "book", "hiya:it was a dark and stormy night, with a pen")
		cv.So(true, cv.ShouldBeTrue)
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

		translation := inc.trMust([]byte(code))

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

		translation := inc.trMust([]byte(code))

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

		translation := inc.trMust([]byte(code))

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

		translation := inc.trMustPre([]byte(code), false)

		cv.So(string(translation), matchesLuaSrc,
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

		translation := inc.trMust([]byte(code))

		//		cv.So(string(translation), matchesLuaSrc,
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

		translation := inc.trMust([]byte(code))

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

		translation := inc.trMust([]byte(code))

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

		translation := inc.trMust([]byte(code))

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
		cv.So(string(inc.trMust([]byte(code))), matchesLuaSrc, `
    __type__.anon_sliceType = __sliceType(__type__.int);   
 	a = __type__.anon_sliceType({[0]=3LL});
    print(#a);
`)
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

		translation := inc.trMust([]byte(code))

		cv.So(string(translation), matchesLuaSrc,
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
		translation = inc.trMust([]byte(codeWithCatch))

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

		translation := inc.trMust([]byte(code))

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

		translation := inc.trMust([]byte(code))
		cv.So(string(translation), matchesLuaSrc,
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

		translation := inc.trMust([]byte(code))

		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "d", 2)
		LuaMustString(vm, "w", "hiya")
		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test199CalculatorExpression(t *testing.T) {

	cv.Convey(`a math expression "= math.Exp(2);" ending in a semicolon should still be computed, after removing the semicolon`, t, func() {

		code := `import "math"`
		code2 := `= math.Exp(2);`

		// expect 7.3890560989307
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.trMust([]byte(code))
		LuaRunAndReport(vm, string(translation))

		translation2 := inc.trMust([]byte(code2))
		LuaRunAndReport(vm, string(translation2))
		LuaRunAndReport(vm, "chk = tonumber(__gijit_ans[0])")
		LuaMustFloat64(vm, "chk", 7.3890560989307)

		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test1990complexReturnValue(t *testing.T) {

	cv.Convey(`a complex return value should work, exercising statements.go:337`, t, func() {

		code := `A0 := func() float64 { return 3 }`
		code2 := `A1 := func() float64 { return 1 + A0() }; a1 := A1()`

		// expect A1() to give 4.0
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.trMust([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)
		LuaRunAndReport(vm, string(translation))

		translation2 := inc.trMust([]byte(code2))
		fmt.Printf("\n translation2='%s'\n", translation2)
		LuaRunAndReport(vm, string(translation2))
		LuaMustFloat64(vm, "a1", 4.0)

		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test1991bytesReturnValue(t *testing.T) {

	cv.Convey(`a []byte return value from a native Go function should work, as should string conversion from []byte to string applied to that native return`, t, func() {

		// non-native
		code := `f := func() []byte { return []byte("hello") }; a := f()`
		code2 := `b := string(a)`

		// expect b to give "hello"
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.trMust([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)
		LuaRunAndReport(vm, string(translation))

		translation2 := inc.trMust([]byte(code2))
		fmt.Printf("\n translation2='%s'\n", translation2)
		LuaRunAndReport(vm, string(translation2))
		LuaMustString(vm, "b", "hello")

		// native:

		// this also produces a collision between the
		// Go package name 'debug' and the Lua library
		// 'debug'. For now we'll remove the Lua library
		// from the keyword list since we aren't using it
		// yet.
		code3 := `import "runtime/debug";`
		code4 := `d := string(debug.Stack()); e := d[:9]`

		translation3 := inc.trMust([]byte(code3))
		fmt.Printf("\n translation3='%s'\n", translation3)
		LuaRunAndReport(vm, string(translation3))

		translation4 := inc.trMust([]byte(code4))
		fmt.Printf("\n translation4='%s'\n", translation4)
		LuaRunAndReport(vm, string(translation4))
		LuaMustString(vm, "e", "goroutine")

		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test1992keywordProtected(t *testing.T) {

	cv.Convey(`the 'in' keyword in Lua is legal variable name in Go, see that it is protected`, t, func() {

		code := `
func f(in int) int {
	return in + 1
}
b := f(1)
`
		// expect b to give 2
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.trMust([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)
		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "b", 2)
		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test1993ArrayOfArrayHasTypeAndLenWorking(t *testing.T) {

	cv.Convey(`var a [3][3]float64;  a[2][2] = 3.14; The assignment should succeed and not fail a bounds check.`, t, func() {

		code := `
var a [3][3]float64
a[2][2] = 3.14
b := len(a[2])
`
		// expect b to give 3
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.trMust([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)
		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "b", 3)
		cv.So(true, cv.ShouldBeTrue)
	})
}

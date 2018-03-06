package compiler

import (
	"testing"

	//"github.com/gijit/gi/pkg/verb"
	cv "github.com/glycerine/goconvey/convey"
)

func Test050CallFmtSprintf(t *testing.T) {

	cv.Convey(`call to fmt.Sprintf simplest, no varargs`, t, func() {

		src := `import "fmt"; a := fmt.Sprintf("hello no-args")`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation, err := inc.Tr([]byte(src))
		panicOn(err)
		pp("go:'%s'  -->  '%s' in lua\n", src, string(translation))
		//fmt.Printf("go:'%#v'  -->  '%#v' in lua\n", src, translation)

		cv.So(string(translation), matchesLuaSrc,
			`a = fmt.Sprintf("hello no-args");`)

		*dbg = true
		LoadAndRunTestHelper(t, vm, translation)
		*dbg = false

		LuaMustString(vm, "a", "hello no-args")
	})
}

/* comment out native imports for now

func Test051CallFmtSprintf(t *testing.T) {

	// big Q here: in what format does Luar expect varargs to Sprintf?
	// i.e. this is where we need to match what luar expects...
	//   in order to pass args to Go functions.
	//
	cv.Convey(`call to fmt.Sprintf, single hard coded argument`, t, func() {

		src := `import "fmt"; a := fmt.Sprintf("hello one: %v", 1)`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation, err := inc.Tr([]byte(src))
		panicOn(err)
		pp("go:'%s'  -->  '%s' in lua\n", src, string(translation))
		//fmt.Printf("go:'%#v'  -->  '%#v' in lua\n", src, translation)

		cv.So(string(translation), matchesLuaSrc,
			`a = fmt.Sprintf("hello one: %v", 1LL);`)

		LoadAndRunTestHelper(t, vm, translation)

		LuaMustString(vm, "a", "hello one: 1")
	})
}

func Test052CallFmtSprintf(t *testing.T) {

	cv.Convey(`call to fmt.Sprintf should run, example: a := fmt.Sprintf("hello %v", 3)`, t, func() {

		src := `import "fmt"; a := fmt.Sprintf("hello %v %v", 3, 4)`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation, err := inc.Tr([]byte(src))
		panicOn(err)
		pp("go:'%s'  -->  '%s' in lua\n", src, string(translation))
		//fmt.Printf("go:'%#v'  -->  '%#v' in lua\n", src, translation)

		cv.So(string(translation), matchesLuaSrc,
			`a = fmt.Sprintf("hello %v %v", 3LL, 4LL);`)

		LoadAndRunTestHelper(t, vm, translation)

		LuaMustString(vm, "a", "hello 3 4")
	})
}

func Test058CallFmtIncr(t *testing.T) {

	cv.Convey(`Given a pre-compiled Go function fmt.Incr, we should be able to call it from gi`, t, func() {

		src := `import "gitesting"; a := gitesting.Incr(1);` // then a should be 2

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation, err := inc.Tr([]byte(src))
		panicOn(err)
		pp("go:'%s'  -->  '%s' in lua\n", src, string(translation))

		cv.So(string(translation), matchesLuaSrc,
			`a = gitesting.Incr(1LL);`)

		LoadAndRunTestHelper(t, vm, translation)

		LuaMustInt64(vm, "a", 2) // can't be LuaMustInt, since `a` is int64!!
	})
}

func Test059CallFmtSummer(t *testing.T) {

	cv.Convey(`Given a pre-compiled Go function gitesting.SummerAny(a ...int), we should be able to call it from gi using gitesting.SummerAny(1, 2, 3);`, t, func() {
		cv.So(SummerAny(1, 2, 3), cv.ShouldEqual, 6)
		pp("good: SummerAny(1,2,3) gave us 6 as expected.")

		src := `import "gitesting"; a := gitesting.SummerAny(1, 2, 3);` // then a should be 6

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation, err := inc.Tr([]byte(src))
		panicOn(err)
		pp("go:'%s'  -->  '%s' in lua\n", src, string(translation))

		cv.So(string(translation), matchesLuaSrc,
			`a = gitesting.SummerAny(1LL, 2LL, 3LL);`)

		LoadAndRunTestHelper(t, vm, translation)

		LuaMustInt64(vm, "a", 6)
	})
}

func Test061CallFmtSummerWithDots(t *testing.T) {

	cv.Convey(`Given b := []int{8,9} and a pre-compiled Go function gitesting.SummerAny(a ...int), the call gitesting.SummaryAny(b...) should expand b into the varargs of SummerAny`, t, func() {

		cv.So(SummerAny(1, 2, 3), cv.ShouldEqual, 6)
		pp("good: SummerAny(1,2,3) gave us 6 as expected.")

		src := `import "gitesting"; b := []int{8,9}; a := gitesting.SummerAny(b...);` // then a = 17

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation, err := inc.Tr([]byte(src))
		panicOn(err)
		pp("go:'%s'  -->  '%s' in lua\n", src, string(translation))

		cv.So(string(translation), matchesLuaSrc, `
  	__type__.anon_sliceType = __sliceType(__type__.int); -- 'IMMEDIATE' anon type printing.
  	b = __type__.anon_sliceType({[0]=8LL, 9LL});
  	a = gitesting.SummerAny(__lazy_ellipsis(b));
`)

		LoadAndRunTestHelper(t, vm, translation)

		LuaMustInt64(vm, "a", 17)
	})
}

func Test062SprintfOneSlice(t *testing.T) {

	cv.Convey(`a := fmt.Sprintf("%#v\n", []int{4,5,6}); should make the string version of the int slice, as opposed to just the 4.`, t, func() {

		src := `import "fmt"; a := fmt.Sprintf("yip %#v eee\n", []int{4,5,6});`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		// need the side effect of loading `import "fmt"` package.
		translation, err := inc.Tr([]byte(src))
		panicOn(err)
		pp("go:'%s'  -->  '%s' in lua\n", src, string(translation))

		cv.So(string(translation), matchesLuaSrc, `
__type__.anon_sliceType = __sliceType(__type__.int); -- 'IMMEDIATE' anon type printing.

  	a = fmt.Sprintf("yip %#v eee\n", __type__.anon_sliceType({[0]=4LL, 5LL, 6LL}));`)
		LoadAndRunTestHelper(t, vm, translation)

		LuaMustString(vm, "a", "yip []interface {}{4, 5, 6} eee\n")
	})
}

func Test063SprintfOneSlice(t *testing.T) {

	cv.Convey(`a := fmt.Sprintf("%v %v %v\n", []interface{}{4,5,6}...); should unpack the slice into 3 different args`, t, func() {

		src := `import "fmt"; a := fmt.Sprintf("yee %v %v %v haw\n", []interface{}{4,5,6}...);`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		// need the side effect of loading `import "fmt"` package.
		translation, err := inc.Tr([]byte(src))
		pp("go:'%s'  -->  '%s' in lua\n", src, string(translation))

		cv.So(string(translation), matchesLuaSrc, `
  	__type__.anon_sliceType = __sliceType(__type__.emptyInterface);

     a = fmt.Sprintf("yee %v %v %v haw\n", __lazy_ellipsis(__type__.anon_sliceType({[0]=4LL, 5LL, 6LL})));
			`)

		LoadAndRunTestHelper(t, vm, translation)

		LuaMustString(vm, "a", "yee 4 5 6 haw\n")
	})
}

func Test064SprintfOneSlice(t *testing.T) {

	cv.Convey(`a := fmt.Sprintf("%v %v\n", "hello", []int{4,5,6}); should send the slice as the 3rd arg`, t, func() {

		src := `import "fmt"; a := fmt.Sprintf("%v %v\n", "hello", []int{4,5,6});`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		// need the side effect of loading `import "fmt"` package.
		translation, err := inc.Tr([]byte(src))
		pp("go:'%s'  -->  '%s' in lua\n", src, string(translation))

		cv.So(string(translation), matchesLuaSrc, `
  	     __type__.anon_sliceType = __sliceType(__type__.int);
      	 a = fmt.Sprintf("%v %v\n", "hello", __type__.anon_sliceType({[0]=4LL, 5LL, 6LL}));
        `)
		LoadAndRunTestHelper(t, vm, translation)

		LuaMustString(vm, "a", "hello [4 5 6]\n")
	})
}

// fmt.Printf("heya %#v %v\n", "hello", []int{55,56})
func Test065PrintfItselfAndOneSlice(t *testing.T) {

	cv.Convey(`fmt.Printf("heya %#v %v %v\n", "hello", []int{55,56}, fmt.Printf); should compile and run, printing a reference to itself`, t, func() {

		src := `import "fmt"; fmt.Printf("heya %#v %v\n", "hello", []int{55,56}, fmt.Printf)`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		// need the side effect of loading `import "fmt"` package.
		translation, err := inc.Tr([]byte(src))
		pp("go:'%s'  -->  '%s' in lua\n", src, string(translation))

		cv.So(string(translation), matchesLuaSrc, `
     	__type__.anon_sliceType = __sliceType(__type__.int);

     	fmt.Printf("heya %#v %v\n", "hello", __type__.anon_sliceType({[0]=55LL, 56LL}), fmt.Printf);
        `)
		LoadAndRunTestHelper(t, vm, translation)

	})
}

func Test087ShadowRegexp(t *testing.T) {

	cv.Convey(`import "regexp" and then using it should function via the shadow_regexp import system`, t, func() {

		src := `import "regexp"
    a := regexp.MustCompile("llo")
	loc := a.FindStringIndex("hello")
	lenloc := len(loc)
	a0 := loc[0]
	a1 := loc[1]
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		// need the side effect of loading `import "fmt"` package.
		translation, err := inc.Tr([]byte(src))
		pp("go:'%s'  -->  '%s' in lua\n", src, string(translation))

		// gotta translate to a.FindStringIndex(), not a:FindStringIndex()
		// because FindStringIndex doesn't take a first 'self' parameter.

		// expressions.go:864-869 for ':' versus '.' in method calls.
		cv.So(string(translation), matchesLuaSrc,
			`
	a = regexp.MustCompile("llo");
	loc = a.FindStringIndex("hello");
	lenloc =  #loc;
	a0 = __gi_GetRangeCheck(loc, 0);
	a1 = __gi_GetRangeCheck(loc, 1);
`)
		LoadAndRunTestHelper(t, vm, translation)

		LuaMustInt64(vm, "a0", 2)
		LuaMustInt64(vm, "a1", 5)
	})
}

*/

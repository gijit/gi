package compiler

import (
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test050CallFmtSprintf(t *testing.T) {

	cv.Convey(`call to fmt.Sprintf simplest, no varargs`, t, func() {

		src := `import "fmt"; a := fmt.Sprintf("hello no-args")`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm)

		translation := inc.Tr([]byte(src))
		pp("go:'%s'  -->  '%s' in lua\n", src, string(translation))
		//fmt.Printf("go:'%#v'  -->  '%#v' in lua\n", src, translation)

		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace,
			`a = fmt.Sprintf("hello no-args");`)

		LoadAndRunTestHelper(t, vm, translation)

		LuaMustString(vm, "a", "hello no-args")
	})
}

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
		inc := NewIncrState(vm)

		translation := inc.Tr([]byte(src))
		pp("go:'%s'  -->  '%s' in lua\n", src, string(translation))
		//fmt.Printf("go:'%#v'  -->  '%#v' in lua\n", src, translation)

		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace,
			`a = fmt.Sprintf("hello one: %v", _gi_NewSlice("interface{}",{1}));`)

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
		inc := NewIncrState(vm)

		translation := inc.Tr([]byte(src))
		pp("go:'%s'  -->  '%s' in lua\n", src, string(translation))
		//fmt.Printf("go:'%#v'  -->  '%#v' in lua\n", src, translation)

		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace,
			`a = fmt.Sprintf("hello %v %v", _gi_NewSlice("interface{}",{3, 4}));`)

		LoadAndRunTestHelper(t, vm, translation)

		LuaMustString(vm, "a", "hello 3 4")
	})
}

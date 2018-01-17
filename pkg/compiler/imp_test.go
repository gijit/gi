package compiler

import (
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test050CallFmtSprintf(t *testing.T) {

	cv.Convey(`call to fmt.Sprintf should run, example: a := fmt.Sprintf("hello %v", 3)`, t, func() {

		src := `import "fmt"; a := fmt.Sprintf("hello %v %v", 3, 4)`

		type person struct {
			Name string
			Age  int
		}

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm)

		//user := &person{"Dolly", 46}

		/*
			// globals
			luar.Register(vm, "", luar.Map{
				// Constants can be registered.
				"msg": "foo",
				// And other values as well.
				"user": user,
			})
		*/

		translation := inc.Tr([]byte(src))
		pp("go:'%s'  -->  '%s' in lua\n", src, string(translation))
		//fmt.Printf("go:'%#v'  -->  '%#v' in lua\n", src, translation)

		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace,
			`a = fmt.Sprintf("hello %v %v", _gi_NewSlice("interface{}",{3, 4}));`)

		LoadAndRunTestHelper(t, vm, translation)

		LuaMustString(vm, "a", "hello 3 4")
	})
}

func Test051CallFmtSprintf(t *testing.T) {

	cv.Convey(`call to fmt.Sprintf simpler, no varargs`, t, func() {

		src := `import "fmt"; a := fmt.Sprintf("hello no-args")`

		type person struct {
			Name string
			Age  int
		}

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

package compiler

import (
	"fmt"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
	//luajit "github.com/glycerine/golua/lua"
	"github.com/glycerine/luar"
)

func Test050CallFmtSprintf(t *testing.T) {

	cv.Convey(`call to fmt.Sprintf should run, example: a := fmt.Sprintf("hello %v", 3)`, t, func() {

		src := `import "fmt"; a := fmt.Sprintf("hello %v", 3)`
		inc := NewIncrState()

		type person struct {
			Name string
			Age  int
		}

		vm := NewLuaVmWithPrelude()
		defer vm.Close()

		//user := &person{"Dolly", 46}

		// fmt
		luar.Register(vm, "fmt", luar.Map{
			// Go functions may be registered directly.
			"Sprintf": fmt.Sprintf,
		})

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
			`    a = fmt.Sprintf("hello %v", 3)`)

		LoadAndRunTestHelper(t, vm, translation)

		LuaMustString(vm, "a", "hello 3")
	})
}

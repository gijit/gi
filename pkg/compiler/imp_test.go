package compiler

import (
	"testing"

	cv "github.com/glycerine/goconvey/convey"
	luajit "github.com/glycerine/golua/lua"
)

func Test050ImportFmt(t *testing.T) {

	cv.Convey(`import "fmt"; a := fmt.Sprintf("hello %v", 3)`, t, func() {

		src := `import "fmt"; a := fmt.Sprintf("hello %v", 3)`
		inc := NewIncrState()
		//cv.So(string(inc.Tr([]byte(src))), cv.ShouldMatchModuloWhiteSpace, `a, b, c = 1, 2, 3;`)

		// verify that it happens correctly.
		vm := luajit.NewState()
		defer vm.Close()
		vm.OpenLibs()

		files, err := FetchPrelude(".")
		panicOn(err)
		LuaDoFiles(vm, files)

		translation := inc.Tr([]byte(src))
		pp("go:'%s'  -->  '%s' in lua\n", src, translation)
		//fmt.Printf("go:'%#v'  -->  '%#v' in lua\n", src, translation)

		LoadAndRunTestHelper(t, vm, translation)

		LuaMustString(vm, "a", "hello 3")
	})
}

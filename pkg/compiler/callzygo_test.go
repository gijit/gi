package compiler

import (
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test1502CallZygoFromGijit(t *testing.T) {
	cv.Convey(`within gijit code: a, err := __zygo("3 + 4"); should return nil error`, t, func() {

		src := `
a, err := __zygo("3 + 4");
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation, err := inc.Tr([]byte(src))
		panicOn(err)
		vv("go:'%s'  -->  '%s' in lua\n", src, string(translation))

		LoadAndRunTestHelper(t, vm, translation)

		//LuaMustInt64(vm, "a", 7)
		LuaMustBeNilGolangError(vm, "err")
	})
}

package compiler

import (
	"fmt"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test110ConvertFloat64ToInt64(t *testing.T) {

	cv.Convey(`converting with int() should take a float64 into an int`, t, func() {

		// don't know why the int() is getting lost. a literal copy
		// into lua works fine, because we defined int() as a constructor
		// for the int64 type, in pkg/compiler/int64.lua.

		code := `
 f := float64(3.5)
i := int(f)
`
		// just translate to `i = int(f)` and it will work.
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation, err := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		//cv.So(string(translation), matchesLuaSrc, ``)

		LuaRunAndReport(vm, string(translation))

		LuaMustInt64(vm, "i", 3)

	})
}

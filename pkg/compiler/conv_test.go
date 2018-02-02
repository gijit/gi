package compiler

import (
	"fmt"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test110ConvertFloat64ToInt64(t *testing.T) {

	cv.Convey(`converting with int() should take a float64 into an int`, t, func() {

		code := `
 f := float64(3.5)
i := int(f)
`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		//cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, ``)

		LuaRunAndReport(vm, string(translation))

		LuaMustInt64(vm, "i", 3)

	})
}

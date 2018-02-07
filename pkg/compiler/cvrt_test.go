package compiler

import (
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test201ConvertToFloat64ActuallyDoes(t *testing.T) {

	cv.Convey(`float64(a+b) should do its job, converting from int to float64`, t, func() {

		code := `
a:=1; b:= 2; c := float64(a + b)
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))

		LuaRunAndReport(vm, string(translation))
		LuaMustFloat64(vm, "c", 3.0)
		cv.So(true, cv.ShouldBeTrue)
	})
}

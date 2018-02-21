package compiler

import (
	"fmt"
	"testing"

	//"github.com/gijit/gi/pkg/verb"
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

		by, err := inc.Tr([]byte(code))
		panicOn(err)
		translation := string(by)
		fmt.Printf("translation = '%s'\n", translation)

		by, err = inc.Tr([]byte(code))
		panicOn(err)
		cv.So(string(by), matchesLuaSrc, `
				a = 1LL;
				b = 2LL;
				c = (tonumber((a+b)));
				`)
		LuaRunAndReport(vm, translation)
		LuaMustFloat64(vm, "c", 3.0)
		cv.So(true, cv.ShouldBeTrue)
	})
}

package compiler

import (
	"fmt"
	"testing"

	//"github.com/glycerine/gi/pkg/verb"
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

		LuaRunAndReport(vm, translation)
		LuaMustFloat64(vm, "c", 3.0)
		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test203ConvertBytesToStringAndBack(t *testing.T) {

	cv.Convey(`a:=[]byte("hi"); b:=string(a); should result in "hi" back in b.`, t, func() {

		code := `a:=[]byte("hi"); b:=string(a);`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		by, err := inc.Tr([]byte(code))
		panicOn(err)
		translation := string(by)
		fmt.Printf("translation = '%s'\n", translation)

		LuaRunAndReport(vm, translation)
		LuaMustString(vm, "b", "hi")
		cv.So(true, cv.ShouldBeTrue)
	})
}

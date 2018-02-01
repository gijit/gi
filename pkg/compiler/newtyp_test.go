package compiler

import (
	"fmt"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

var _ = fmt.Printf
var _ = testing.T{}
var _ = cv.So

func Test105NewTypeDeclaration(t *testing.T) {

	cv.Convey(`declare a new named type`, t, func() {

		code := `
type Bean int
b := Bean(99)
`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		//cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, ``)

		LuaRunAndReport(vm, string(translation))

		LuaMustInt64(vm, "b", 99)

	})
}

func Test106NewTypeDeclaration(t *testing.T) {

	cv.Convey(`declare a new named type`, t, func() {

		code := `
type Bean struct{
  counter int
}
b := Bean{counter: 3}
c := b.counter
`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		//cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, ``)

		LuaRunAndReport(vm, string(translation))

		LuaMustInt64(vm, "c", 3)

	})
}

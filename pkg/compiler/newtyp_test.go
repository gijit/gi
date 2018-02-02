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

func Test107NewTypeForFloat64(t *testing.T) {

	cv.Convey(`declare a new named type for a basic float64`, t, func() {

		code := `
type F float64
var f F = 2.5
g := F(3.4)
`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		//cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, ``)

		LuaRunAndReport(vm, string(translation))

		LuaMustFloat64(vm, "f", 2.5)
		LuaMustFloat64(vm, "g", 3.4)

	})
}

/*
// unfinished
func Test108TypesHavePackagePath(t *testing.T) {

	cv.Convey(`types at the repl have the "main" package (or whatever package we are working in) short name AND full path attached to their type, so we can distinguish types with the same name that come from different packages (even vendored) whose ultimate paths differ`, t, func() {

		code := `
type Bean struct{
  counter int
}

`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		//cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, ``)

		LuaRunAndReport(vm, string(translation))

//		LuaMustInt64(vm, "c", 3)

	})
}
*/

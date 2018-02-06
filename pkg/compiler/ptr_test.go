package compiler

import (
	"fmt"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test098Pointers(t *testing.T) {

	cv.Convey(`taking pointers and referencing through them`, t, func() {

		code := `
type S struct{
  val int
}
var s S
a := &s
a.val = 10
b := s.val
`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		//		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, ``)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		LuaMustInt64(vm, "b", 10)

	})
}

func Test099PointerDeference(t *testing.T) {

	cv.Convey(`dereferencing a pointer`, t, func() {

		code := `
a:= 1
b := &a
c := *b
`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		//		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, ``)

		LuaRunAndReport(vm, string(translation))

		LuaMustInt64(vm, "c", 1)

	})
}

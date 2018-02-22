package compiler

import (
	"fmt"
	"testing"

	//"github.com/gijit/gi/pkg/verb"
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
read := a.val
a.val = 10
b := s.val
`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation, err := inc.Tr([]byte(code))
		panicOn(err)
		fmt.Printf("\n translation='%s'\n", translation)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		LuaMustInt64(vm, "b", 10)
		LuaMustInt64(vm, "read", 0)

	})
}

func Test099PointerDeference(t *testing.T) {

	cv.Convey(`dereferencing a pointer`, t, func() {

		code := `
a:= 1
b := &a
c := *b
*b = 3
`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation, err := inc.Tr([]byte(code))
		panicOn(err)
		fmt.Printf("\n translation='%s'\n", translation)

		LuaRunAndReport(vm, string(translation))

		LuaMustInt64(vm, "c", 1)
		LuaMustInt64(vm, "a", 3)

	})
}

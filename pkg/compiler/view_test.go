package compiler

import (
	"fmt"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test510StructsShouldBeViewable(t *testing.T) {

	cv.Convey(`__tostring(str) where str is a struct, should print the struct's fields`, t, func() {

		code := `
type S struct {
   b int
}
var s S;
s.b = 23
chk := __tostring(s)
chk2 := __tostring(&s)
`
		// c0 should be 4, a1 should be 4
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.trMust([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		// By rights, chk should say `main.S{b: 23LL, }`.
		// But `&main.S{b: 23LL, }` is correct because
		// in order for pointer comparison to work
		// (and have two pointers to the same strut
		// have equal values), even a basic struct
		// must need be a pointer. *shrug*. That
		// was the GopherJS design. It works, so
		// we keep it.
		LuaMustString(vm, "chk", `&main.S{b: 23LL, }`)
		LuaMustString(vm, "chk2", `&main.S{b: 23LL, }`)

	})
}

func Test511SlicesShouldBeViewable(t *testing.T) {

	cv.Convey(`tostring([]string{"a","b"}) should stringify the slice`, t, func() {

		code := `
chk := __tostring([]string{"a","b"})
`
		// c0 should be 4, a1 should be 4
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.trMust([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		LuaMustString(vm, "chk", `[]string{[0]= "a", [1]= "b", }`)

	})
}

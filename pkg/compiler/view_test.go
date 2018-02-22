package compiler

import (
	"fmt"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test510StructsShouldBeViewable(t *testing.T) {

	cv.Convey(`tostring(str) where str is a struct, should print the struct's fields`, t, func() {

		code := `
type Str struct {
 s string
 i int
 S string
 I int
}
var str Str;
chk := tostring(str)
chk2 := tostring(&str)
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

		LuaMustString(vm, "chk", `&main.Str{s: "", i: 0LL, S: "", I: 0LL, }`)
		LuaMustString(vm, "chk2", `&main.Str{s: "", i: 0LL, S: "", I: 0LL, }`)

	})
}

func Test511SlicesShouldBeViewable(t *testing.T) {

	cv.Convey(`tostring([]string{"a","b"}) should stringify the slice`, t, func() {

		code := `
chk := tostring([]string{"a","b"})
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

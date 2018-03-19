package compiler

import (
	"fmt"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test1100LabelsAndGotoAtTopLevel(t *testing.T) {

	cv.Convey("goto and labels should work", t, func() {

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		code := `
package main

import "fmt"

var i = 0
func main() {

lab:
	fmt.Printf("hi %v\n", i)
	i++
	if i < 3 {
		goto lab
	}
}
main()
`

		translation := inc.trMust([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)
		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "i", 3)
		cv.So(true, cv.ShouldBeTrue)
	})
}

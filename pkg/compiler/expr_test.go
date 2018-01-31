package compiler

import (
	"fmt"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test094MathExpressionsShouldWorkAtREPL(t *testing.T) {

	cv.Convey(`simple math expressions like 3 + 4 should return results at the REPL`, t, func() {

		code := `
3 + 4
`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, `
println(3 + 4)
`)

	})
}

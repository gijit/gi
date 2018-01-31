package compiler

import (
	"fmt"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test094MathExpressionsShouldWorkAtREPL(t *testing.T) {

	cv.Convey(`simple math expressions like = 3 + 4 should return results at the REPL. Following the Lua convention, in order to view a simple expression, the user adds an equals sign '=' to the start of the line.`, t, func() {

		code := `
= 3 + 4
`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, `
__gijit_ans = 7LL;
__gijit_printQuoted(__gijit_ans);
`)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		// check for exception
		LuaMustInt64(vm, "__gijit_ans", 7)

	})
}

func Test095StringExpressionsShouldWorkAtREPL(t *testing.T) {

	cv.Convey(`simple string expressions like = "hi" + " there" should return results at the REPL. Following the Lua convention, in order to view a simple expression, the user adds an equals sign '=' to the start of the line.`, t, func() {

		code := `
= "hi" + " there"
`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, `
__gijit_ans = "hi there";
__gijit_printQuoted(__gijit_ans);
`)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		// check for exception
		LuaMustString(vm, "__gijit_ans", "hi there")

	})
}

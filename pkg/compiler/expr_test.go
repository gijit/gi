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

		//cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, ``)

		withHelp := append(translation, []byte("\n a = __gijit_ans[0]\n")...)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(withHelp))

		LuaMustInt64(vm, "a", 7)

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

		//cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, ``)

		withHelp := append(translation, []byte("\n a = __gijit_ans[0]\n")...)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(withHelp))

		// check for exception
		LuaMustString(vm, "a", "hi there")

	})
}

func Test096MultipleExpressionsAtOnceAtTheREPL(t *testing.T) {

	cv.Convey(`Multiple expressions after the equals sign: = 1, 2+5, "hi"+" there" should print all three expressions`, t, func() {

		code := `
= 1, 2+5, "hi" + " there"
`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		//		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, ``)

		withHelp := append(translation, []byte("\n a = __gijit_ans[0]\n  b = __gijit_ans[1]\n  c = __gijit_ans[2]\n ")...)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(withHelp))

		LuaMustInt64(vm, "a", 1)
		LuaMustInt64(vm, "b", 7)
		LuaMustString(vm, "c", "hi there")

	})
}

func Test097SingleExpressionAtTheREPL(t *testing.T) {

	cv.Convey(`One expression after the equals sign: = 2+5, should print 7LL`, t, func() {

		code := `
= 2+5
`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		//		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, ``)

		withHelp := append(translation, []byte("\n a = __gijit_ans[0]\n")...)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(withHelp))

		LuaMustInt64(vm, "a", 7)

	})
}

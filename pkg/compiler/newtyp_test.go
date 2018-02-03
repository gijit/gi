package compiler

import (
	"fmt"
	"strings"
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

		code2 := `	var h float64 = g;`
		// expect:
		// ./m.go:27:6: cannot use g (type F) as type float64 in assignment

		gotPanic := false
		var pval string
		f := func() {
			defer func() {
				r := recover()
				if r != nil {
					gotPanic = true
					pval = fmt.Sprintf("%v", r)
				}
			}()
			translation2 := inc.Tr([]byte(code2)) // should panic
			fmt.Printf("\n translation2='%s'\n", translation2)
		}
		f()
		cv.So(gotPanic, cv.ShouldBeTrue)
		cv.So(strings.Contains(pval, "cannot use g (variable of type F) as float64"),
			cv.ShouldBeTrue)
	})
}

func Test108SyntaxErrorDoesNotMessUpTypeSystem(t *testing.T) {

	cv.Convey(`a syntax error was disabling subsequent type checks`, t, func() {

		// this sequence is messing with the type checker's state
		/*
		   type F float64     (1) // a new type, is not equal to float64
		   var f F = 2.5      (2) // an instance of that type.
		   var float64 w      (3) // a syntax error involving the float64 type.
		   var a float64 = f  (4) // *should* be a type assignment error.
		*/
		// Without line (3), the type checker correctly rejects line (4).
		// (change `if true` to `if false` below to verify that).
		// With line (3), after the syntax error, the type checker allows (4).
		//
		code := `
type F float64
var f F = 2.5
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

		gotPanic := false
		var pval string

		if true {
			// this is a deliberate syntax error, but it
			// is/was upsetting subsequent type checks:
			code2 := `var float64 w`
			// expect:
			// ./m.go:27:6: cannot use g (type F) as type float64 in assignment

			f := func() {
				defer func() {
					r := recover()
					if r != nil {
						gotPanic = true
						pval = fmt.Sprintf("%v", r)
					}
				}()
				translation2 := inc.Tr([]byte(code2)) // should panic
				fmt.Printf("\n translation2='%s'\n", translation2)
			}
			f()
			cv.So(gotPanic, cv.ShouldBeTrue)
			cv.So(strings.Contains(pval, "bad identifier: cannot use 'float64' as an identifier in gijit"),
				cv.ShouldBeTrue)
		}

		code3 := `var a float64 = f`

		gotPanic = false
		pval = ""
		f2 := func() {
			defer func() {
				r := recover()
				if r != nil {
					gotPanic = true
					pval = fmt.Sprintf("%v", r)
				}
			}()
			translation3 := inc.Tr([]byte(code3)) // should panic, but was not.
			fmt.Printf("\n translation3='%s'\n", translation3)
		}
		f2()
		fmt.Printf("pval = '%v'\n", pval)
		cv.So(gotPanic, cv.ShouldBeTrue)
		cv.So(strings.Contains(pval, "cannot use f (variable of type F) as float64 value in variable declaration"),
			cv.ShouldBeTrue)

	})
}

func Test109NewTypeSpaceAndVariableSpaceAreSeparate(t *testing.T) {

	cv.Convey(`In Go, the type namespace and the variable namespace are distinct, so that one can have a variable name that is the same as a type name, and there is no compile error. Apparently this allows the introduction of new pre-defined types without breaking old code.`, t, func() {

		code := `
type F struct{
  a int
}
type G struct{
   F F
}
g := G{F:F{a:2}}
two := g.F.a;`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		//cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, ``)

		LuaRunAndReport(vm, string(translation))

		LuaMustInt64(vm, "two", 2)
	})
}

/*
// unfinished
func Test110TypesHavePackagePath(t *testing.T) {

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

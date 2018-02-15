package compiler

/* come back to this
import (
	"fmt"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test400RuntimeReflection(t *testing.T) {

	cv.Convey("runtime type reflection should be able to distinguish between `Inch` and `Meter`, different named types that share the same `int` basic type", t, func() {

		code := `
package main

import (
	"fmt"
	"reflect"
)

type Meter int
type Inch int

type Unit interface {
	PerMile(miles int) float64
}

func (m Meter) PerMile(miles int) float64 {
	return float64(miles) * 1609.34
}

func (n Inch) PerMile(miles int) float64 {
	return float64(miles) * 63359.84251872
}

func Comparable(a Unit, b Unit) bool {
	ta := reflect.TypeOf(a)
	tb := reflect.TypeOf(b)
	return ta == tb
}

// func main() {
	m0 := Meter(0)
	i0 := Inch(0)

	m1 := Meter(1)
	i1 := Inch(1)

	// should be
	m0m1 := Comparable(m0, m1) // true
	fmt.Printf("m0m1 = %v\n", m0m1)
	m1m0 := Comparable(m1, m0) // true
	fmt.Printf("m1m0 = %v\n", m1m0)

	i0i1 := Comparable(i0, i1) // true
	fmt.Printf("i0i1 = %v\n", i0i1)
	i1i0 := Comparable(i1, i0) // true
	fmt.Printf("i1i0 = %v\n", i1i0)

	m0i0 := Comparable(m0, i0) // false
	fmt.Printf("m0i0 = %v\n", m0i0)
	i0m0 := Comparable(i0, m0) // false
	fmt.Printf("i0m0 = %v\n", i0m0)

	i1m1 := Comparable(i1, m1) // false
	fmt.Printf("i1m1 = %v\n", i1m1)
	m1i1 := Comparable(m1, i1) // false
	fmt.Printf("m1i1 = %v\n", m1i1)

	m1i0 := Comparable(m1, i0) // false
	fmt.Printf("m1i0 = %v\n", m1i0)
// }

`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		//cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, ``)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		LuaMustBool(vm, "m0m1", true)
		LuaMustBool(vm, "m1m0", true)

		LuaMustBool(vm, "i0i1", true)
		LuaMustBool(vm, "i1i0", true)

		LuaMustBool(vm, "m0i0", false)
		LuaMustBool(vm, "i0m0", false)

		LuaMustBool(vm, "i1m1", false)
		LuaMustBool(vm, "m1i1", false)

		LuaMustBool(vm, "m1i0", false)

	})
}
*/

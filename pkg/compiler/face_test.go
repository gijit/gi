package compiler

import (
	"fmt"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

var _ = fmt.Printf
var _ = testing.T{}
var _ = cv.So

/*
	       0) declare a var as an interface

		        a) one-value conversion:
		             as := any.(Stringer)

		    b)    two-value conversion check:

		   type Stringer interface {
		      String() string
		   }
		    if v, ok := any.(Stringer); ok {
		        return v.String()
		    }

		    c) type switch:

		   func ToString(any interface{}) string {

		    switch v := any.(type) {
		    case int:
		        return strconv.Itoa(v)
		    case float:
		        return strconv.Ftoa(v, 'g', -1)
		    }
		    return "???"
		    }

		                d) assignment /compile time check:

		       var s Stringer = &MyType{}

*/

func Test100InterfaceDeclaration(t *testing.T) {

	cv.Convey(`declare an interface`, t, func() {

		code := `
type Counter interface {
   Next() int
}
type S struct {
   v int
}
func (s *S) Next() int {
  s.v++
  return s.v
}
var c Counter = &S{}
a := c.Next()
b := c.Next()

r := &S{}
r1 := r.Next()
r2 := r.Next()

type ByTen struct {
   v int
}
func (s *ByTen) Next() int {
   s.v += 10
   return s.v
}
bt := &ByTen{}
c = bt
d := c.Next()
e := c.Next()
	`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		LuaMustInt64(vm, "a", 1)
		LuaMustInt64(vm, "b", 2)

		LuaMustInt64(vm, "d", 10)
		LuaMustInt64(vm, "e", 20)

		LuaMustInt64(vm, "r1", 1)
		LuaMustInt64(vm, "r2", 2)
		cv.So(true, cv.ShouldBeTrue)
	})
}

/* WIP
func Test101InterfaceConversion(t *testing.T) {

	// work in progress

	cv.Convey(`two-value interface conversion check`, t, func() {

		code := `
		package main

		import (
			"fmt"
		)

		type Counter interface {
			Next() int
		}
		type S struct {
			v int
		}

		func (s *S) Next() int {
			s.v++
			return s.v
		}

		type Bad struct {
			v int
		}

		//func main() {

			s := &S{}

			asCounter_s, s_ok := interface{}(s).(Counter)
			sNil := asCounter_s == nil

			a := asCounter_s.Next()
			b := asCounter_s.Next()

			bad := &Bad{}

			asCounter_bad, bad_ok := interface{}(bad).(Counter)
			acbIsNil := asCounter_bad == nil

			fmt.Printf("s_ok=%v, asCounter_s=%v, sNil=%v, a=%v, b=%v, acbIsNil=%v, bad_ok=%v\n", s_ok, asCounter_s, sNil, a, b, acbIsNil, bad_ok)
			// s_ok=true, asCounter_s=&{2}, sNil=false, a=1, b=2, acbIsNil=true, bad_ok=false

		//}
			`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)


		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		LuaMustBool(vm, "sNil", false)
		LuaMustBool(vm, "s_ok", true)

		LuaMustInt64(vm, "a", 1)
		LuaMustInt64(vm, "b", 2)

		LuaMustBool(vm, "abcIsNil", true)
		LuaMustBool(vm, "bad_ok", false)

	})

}
*/

func Test102InterfaceMethodset(t *testing.T) {

	cv.Convey(`the methodsets of interfaces and structs can be compared to check for interface satisfaction.`, t, func() {
		code := `
package main

import (
	"fmt"
)

type Bowser interface {
	Hi()
}

type Possum interface {
	Hi()
    Pebbles()
}

type Unsat interface {
	Hi()
    Pebbles()
    MissMe()
}

type B struct{}

func (b *B) Hi() {
	fmt.Printf("B.Hi called\n")
}
func (b *B) Pebbles() {}

    chk := 0
	var v Bowser = &B{}
	switch v.(type) {
    case Possum:
		fmt.Printf("ooh! it types as a Possum!\n")
        chk = 2
	case Bowser:
		fmt.Printf("yabadadoo! it types as a Bowser!\n")
        chk = 1
	}
    fmt.Printf("chk = '%v'\n", chk)

    // and verify that v implements Bowser too:
    asBowser, isBowser := v.(Bowser)
    asIsNil := (asBowser == nil)

    // negative check, should not convert:
    asUn, isUn := v.(Unsat)
    asUnNil := (asUn == nil)
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "chk", 2)
		LuaMustBool(vm, "isBowser", true)
		LuaMustBool(vm, "asIsNil", false)

		LuaMustBool(vm, "isUn", false)
		LuaMustBool(vm, "asUnNil", true)
	})
}

func Test202InterfaceMethodset(t *testing.T) {

	cv.Convey(`the methodsets of interfaces and structs can be compared to check for interface satisfaction, including the method types not just names`, t, func() {
		code := `
package main

import (
	"fmt"
)

type Bowser interface {
	Hi() int
}

type Possum interface {
	Hi() int
    Pebbles() int
}

type Unsat interface {
	Hi() int
    Pebbles() int
    MissMe() int
}

type B struct{}

func (b *B) Hi() {
	fmt.Printf("B.Hi called\n")
}
func (b *B) Pebbles() {}

chk := 0
var v interface{} = &B{}
switch v.(type) {
    case Possum:
		fmt.Printf("ooh! it types as a Possum!\n")
        chk = 2
	case Bowser:
		fmt.Printf("yabadadoo! it types as a Bowser!\n")
        chk = 1
	}
    fmt.Printf("chk = '%v'\n", chk)

    // and verify that v does not implement Bowser too:
    asBowser, isBowser := v.(Bowser)
    asIsNil := (asBowser == nil)

    // negative check, should not convert:
    asUn, isUn := v.(Unsat)
    asUnNil := (asUn == nil)
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "chk", 0)
		LuaMustBool(vm, "isBowser", false)
		LuaMustBool(vm, "asIsNil", true)

		LuaMustBool(vm, "isUn", false)
		LuaMustBool(vm, "asUnNil", true)
	})
}

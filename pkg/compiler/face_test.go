package compiler

import (
	"fmt"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test092Interfaces(t *testing.T) {

	/*
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

	cv.Convey(``, t, func() {

		code := `

`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, `
`)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		// check for exception
		LuaMustInt64(vm, "c0", 4)
		LuaMustInt64(vm, "a1", 4)

	})
}

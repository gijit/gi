package compiler

import (
	"testing"

	//"github.com/gijit/gi/pkg/verb"
	cv "github.com/glycerine/goconvey/convey"
)

func Test500MatrixDeclOfDoubleSlice(t *testing.T) {

	cv.Convey(`[][]float inside matrix struct`, t, func() {

		src := `
type Matrix struct {
	A    [][]float64
}
m := &Matrix{A:[][]float64{[]float64{1,2}}}
e := m.A[0][1]
`
		// e == 2
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(src))
		pp("go:'%s'  -->  '%s' in lua\n", src, string(translation))
		//fmt.Printf("go:'%#v'  -->  '%#v' in lua\n", src, translation)

		//cv.So(string(translation), matchesLuaSrc, ``)

		LoadAndRunTestHelper(t, vm, translation)

		LuaMustFloat64(vm, "e", 2)
	})
}

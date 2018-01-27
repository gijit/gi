package compiler

import (
	"fmt"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test080Int64ArraysByGoProxyCopyAppend(t *testing.T) {

	cv.Convey(`Proxies for [3]int64 should be allocated from Lua and passable to a Go native function`, t, func() {

		// a := [3]int64{1,2,3}
		// should call, a = __lua2go(_gi_NewArray({1,2,3}, "int64", 3))
		code := `
   import "gitesting"
   a := [3]int64{1,2,3}
   a[0]--
   sum := gitesting.SumArrayInt64(a)
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace,
			// 		a = __lua2go(_gi_NewArray({[0]=1LL,2LL,3LL}, "int64", 3));
			`
		a = _gi_NewArray({[0]=1LL,2LL,3LL}, "int64", 3);
        a[0] = (a[0] - (1LL));
        sum = gitesting.SumArrayInt64(__gi_clone(a, "arrayType"));
`)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		// check for exception
		LuaMustInt64(vm, "sum", 6)

	})
}

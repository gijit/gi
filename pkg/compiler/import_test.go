package compiler

import (
	"fmt"
	"testing"

	//"github.com/gijit/gi/pkg/verb"
	cv "github.com/glycerine/goconvey/convey"
)

func Test1000ImportAGoSourcePackage(t *testing.T) {

	cv.Convey(`import a Go source package`, t, func() {

		code := `
import "github.com/gijit/gi/pkg/compiler/spkg_tst"
caught := spkg_tst.Fish(2)
`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation, err := inc.Tr([]byte(code))
		panicOn(err)
		fmt.Printf("\n translation='%s'\n", translation)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		LuaMustInt64(vm, "caught", 4)
	})
}

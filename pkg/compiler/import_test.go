package compiler

import (
	"fmt"
	"os"
	"testing"

	//"github.com/gijit/gi/pkg/verb"
	cv "github.com/glycerine/goconvey/convey"
)

func Test1000ImportAGoSourcePackage(t *testing.T) {

	cv.Convey(`import a Go source package`, t, func() {

		FishMultipliesBy(2)
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

func Test1001NoCachingOfImportsOfGoSourcePackages(t *testing.T) {

	cv.Convey(`since they may be in flux, importing a Go source package must re-read the source every time, and not use a cached version`, t, func() {

		defer FishMultipliesBy(2)
		for i := 1; i <= 2; i++ {
			FishMultipliesBy(i + 1) // 2, then 3
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

			LuaMustInt64(vm, "caught", int64(2*(i+1)))
		}
	})
}

func FishMultipliesBy(i int) {
	f, err := os.Create("spkg_tst/spkg.go")
	panicOn(err)
	fmt.Fprintf(f, `
package spkg_tst

func Fish(numPole int) (fishCaught int) {
	return numPole * %v
}
`, i)
	f.Close()
}

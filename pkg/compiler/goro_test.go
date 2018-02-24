package compiler

import (
	//"fmt"
	"testing"

	//"github.com/gijit/gi/pkg/token"
	//"github.com/gijit/gi/pkg/types"
	cv "github.com/glycerine/goconvey/convey"
	//"github.com/glycerine/luar"
)

func Test700StartGoroutine(t *testing.T) {

	cv.Convey(`start a new goroutine that gets its own *golua.State`, t, func() {

		r, err := NewGoro(nil)
		panicOn(err)

		r2, err := NewGoro(nil)
		panicOn(err)
		_, _ = r, r2
		/*
			ch := make(chan int, 1)
			ch <- 57

			luar.Register(vm, "", luar.Map{
				"ch": ch,
			})

			code := `b := 3`
			inc := NewIncrState(vm, nil)
			translation, err := inc.Tr([]byte(code))
			panicOn(err)
			pp("translation='%s'", string(translation))
			LuaRunAndReport(vm, string(translation))
			LuaMustInt64(vm, "b", 3)

			// allow ch to type check
			pkg := inc.pkgMap["main"].Arch.Pkg
			scope := pkg.Scope()
			nt64 := types.Typ[types.Int64]
			chVar := types.NewVar(token.NoPos, pkg, "ch", types.NewChan(types.SendRecv, nt64))
			scope.Insert(chVar)

			code = `a := <- ch;`

			translation, err = inc.Tr([]byte(code))
			panicOn(err)

			pp("translation='%s'", string(translation))

			LuaRunAndReport(vm, string(translation))
			LuaMustInt64(vm, "a", 57)
			cv.So(true, cv.ShouldBeTrue)
		*/
	})
}

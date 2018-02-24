package compiler

import (
	//"fmt"
	"testing"

	"github.com/gijit/gi/pkg/token"
	"github.com/gijit/gi/pkg/types"
	cv "github.com/glycerine/goconvey/convey"
	//"github.com/glycerine/luar"
)

func Test700StartGoroutine(t *testing.T) {

	cv.Convey(`start a new goroutine that gets its own *golua.State`, t, func() {

		r, err := NewGoro(nil, nil)
		panicOn(err)

		//r2, err := NewGoro(nil)
		//panicOn(err)

		t0 := r.newTicket()

		ch := make(chan int, 1)
		ch <- 57

		t0.regmap["ch"] = ch

		// first run instantiates the main package so we can add 'ch' to it.
		code := `b := 3`
		inc := NewIncrState(r.vm, nil)
		translation, err := inc.Tr([]byte(code))
		panicOn(err)
		pp("translation='%s'", string(translation))
		LuaRunAndReport(r.vm, string(translation))
		LuaMustInt64(r.vm, "b", 3)

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

		t0.run = translation

		t0.varname["a"] = true

		// execute the `a := <- ch;`
		panicOn(t0.Do())

		ai := t0.varname["a"].(int64)

		cv.So(ai, cv.ShouldEqual, 57)

	})
}

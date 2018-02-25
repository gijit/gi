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

		// key change: an *unbuffer* channel
		ch := make(chan int)
		go func() {
			ch <- 57
		}()

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

/*
func Test701StartTwoGoroutines(t *testing.T) {

	cv.Convey(`start two goroutines that communicate`, t, func() {

		r0, err := NewGoro(nil, nil)
		panicOn(err)

		r1, err := NewGoro(nil, nil)
		panicOn(err)

		t0 := r0.newTicket()
		t1 := r1.newTicket()

		// Big question:
		// how do these two vms learn about their shared channel?

		// the go func itself is a closure, typically grabbing
		// all the variables it sees in scope.

		code0 := `ch := make(chan int)`
		inc := NewIncrState(r.vm, nil)
		translation0, err := inc.Tr([]byte(code))
		panicOn(err)
		pp("translation='%s'", string(translation))
		LuaRunAndReport(r.vm, string(translation))
		LuaMustInt64(r.vm, "b", 3)

		code1 = `a := <- ch;`

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
*/
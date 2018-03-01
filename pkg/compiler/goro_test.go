package compiler

import (
	//"fmt"
	"testing"

	//"github.com/gijit/gi/pkg/token"
	//"github.com/gijit/gi/pkg/types"
	cv "github.com/glycerine/goconvey/convey"
	//"github.com/glycerine/luar"
)

func Test707ReplGoroVsBackendGoro(t *testing.T) {

	cv.Convey(`In order to allow background goroutines to run, the frontend of the repl runs on its own goroutine, and the backend of runs its own goroutine to keep the scheduler alive and running LuaJIT code. Therefore we should see, even when waiting at the REPL and not typing any input, that background goroutines are running.`, t, func() {

		code := `
  a := 1
  aa := []int{}
  ch := make(chan int)
  go func() {
      for i :=0; i < 3; i++ {
         got := <-ch
         a += 1 + got
         aa = append(aa, a)
         println("a is now ", a)
      }
  }()
  for j:=0; j < 3; j++ {
      ch <- j
  }
`
		// 'a' should be 7
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		inc := NewIncrState(vm, nil)
		translation, err := inc.Tr([]byte(code))
		panicOn(err)

		pp("translation='%s'", string(translation))
		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "a", 7)
		LuaMustEvalToInt64(vm, "aa[0]", 1)
		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test708ReplGoroVsBackendGoro(t *testing.T) {

	cv.Convey(`In order to allow background goroutines to run, the frontend of the repl runs on its own goroutine, and the backend of runs its own goroutine to keep the scheduler alive and running LuaJIT code. Therefore we should see, even when waiting at the REPL and not typing any input, that background goroutines are running.  Send and receive should work going from the repl to the background goroutine.`, t, func() {

		code := `
  accumRecv := []int{}
  a0 := 0
  c1 := make(chan int)
  c2 := make(chan int)
  nextSend := 2
  go func() {
      for {
         select {
            case c1 <- nextSend:
               println("background goro sent ", nextSend)
               nextSend++
            case r := <- c2:
               accumRecv = append(accumRecv, r)
               a0 = accumRecv[0]
               println("background goro received ", r)
         }
      }
  }()
`
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		*dbg = true
		inc := NewIncrState(vm, nil)
		translation, err := inc.Tr([]byte(code))
		panicOn(err)

		pp("translation='%s'", string(translation))
		LuaRunAndReport(vm, string(translation))

		// 2nd interaction at the repl
		code2 := ` j2 := <-c1; c2 <- 33`

		translation, err = inc.Tr([]byte(code2))
		panicOn(err)

		pp("2nd translation = '%s'", string(translation))
		LuaRunAndReport(vm, string(translation))

		LuaMustInt64(vm, "j2", 2)
		LuaMustInt64(vm, "nextSend", 3)
		LuaMustInt64(vm, "a0", 33)
		LuaMustEvalToInt64(vm, "accumRecv[0]", 33)
		cv.So(true, cv.ShouldBeTrue)
	})
}

// while we work on lua-only goroutines, comment this out.
/*

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
*/

/* not done:
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

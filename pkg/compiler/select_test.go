package compiler

import (
	"fmt"
	"testing"

	"github.com/gijit/gi/pkg/token"
	"github.com/gijit/gi/pkg/types"
	cv "github.com/glycerine/goconvey/convey"
	"github.com/glycerine/luar"
)

func Test600RecvOnChannel(t *testing.T) {

	cv.Convey(`in Lua, receive an integer on a buffered channel, previously sent by Go`, t, func() {

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		ch := make(chan int, 1)
		ch <- 57

		luar.Register(vm, "", luar.Map{
			"ch": ch,
		})

		// first run instantiates the main package so we can add 'ch' to it.
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
	})
}

func Test601SendOnChannel(t *testing.T) {

	cv.Convey(`in Lua, send to a buffered channel`, t, func() {

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		ch := make(chan int, 1)

		luar.Register(vm, "", luar.Map{
			"ch": ch,
		})

		// first run instantiates the main package so we can add 'ch' to it.
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

		code = `ch <- 6; send:=true;`

		translation, err = inc.Tr([]byte(code))
		panicOn(err)

		pp("translation='%s'", string(translation))

		LuaRunAndReport(vm, string(translation))

		a := <-ch
		fmt.Printf("a received! a = %v\n", a)
		//LuaMustBool(vm, "sent", true)
		cv.So(a, cv.ShouldEqual, 6)
	})
}

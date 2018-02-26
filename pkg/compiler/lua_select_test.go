package compiler

// lua-only verion of select_test

// For the moment, we'll comment this out while
// we get the fully lua system going.

import (
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test800LuaOnlyRecvOnChannel(t *testing.T) {

	cv.Convey(`in all-Lua, receive an integer on a buffered channel, previously sent by Lua`, t, func() {

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		code := `ch:=make(chan int, 1); ch <- 57; a:= <- ch`
		inc := NewIncrState(vm, nil)
		translation, err := inc.Tr([]byte(code))
		panicOn(err)
		//*dbg = true
		pp("translation='%s'", string(translation))
		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "a", 57)

		cv.So(true, cv.ShouldBeTrue)
	})
}

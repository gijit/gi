package compiler

// chan_test.go: tests of the all-lua channels

import (
	//"fmt"
	"testing"
	"time"

	//"github.com/gijit/gi/pkg/token"
	//"github.com/gijit/gi/pkg/types"
	cv "github.com/glycerine/goconvey/convey"
	//"github.com/glycerine/luar"
)

func Test900ForeverBlockingSelect(t *testing.T) {

	cv.Convey("select{} should block the goroutine forever", t, func() {

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		// with default: present we should not block
		// _selection = __task.select({{}});
		code := ` a:= 0; go func() { a = 1; select{ default: }; a= 2; }() // should not block`
		translation, err := inc.Tr([]byte(code))
		//*dbg = true
		pp("translation='%s'", string(translation))

		/*
			LuaRunAndReport(vm, string(translation))
			LuaMustInt64(vm, "a", 2)
		*/

		//  _r = __task.select({});
		code = ` b:= 0; go func() { b = 1; select{}; b= 2; }() // should block goroutine forever`
		translation, err = inc.Tr([]byte(code))
		panicOn(err)
		pp("translation='%s'", string(translation))

		LuaRunAndReport(vm, string(translation))
		select {
		case <-time.After(1 * time.Second):
		}
		LuaMustInt64(vm, "b", 1)
		cv.So(true, cv.ShouldBeTrue)
	})
}

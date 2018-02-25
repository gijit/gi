package compiler

// zgoro_test: tests of the lua-only channels

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

		code := ` a:= 0; go func() { a = 1; select{}; a= 2; }() // should block goroutine forever`
		inc := NewIncrState(vm, nil)
		translation, err := inc.Tr([]byte(code))
		panicOn(err)
		pp("translation='%s'", string(translation))
		LuaRunAndReport(vm, string(translation))
		select {
		case <-time.After(1 * time.Second):
		}
		LuaMustInt64(vm, "a", 1)
		cv.So(true, cv.ShouldBeTrue)
	})
}

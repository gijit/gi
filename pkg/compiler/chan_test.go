package compiler

// chan_test.go: tests of the all-lua channels

import (
	"fmt"
	"testing"
	"time"

	//"github.com/gijit/gi/pkg/token"
	//"github.com/gijit/gi/pkg/types"
	cv "github.com/glycerine/goconvey/convey"
	//"github.com/glycerine/luar"
)

func Test900SendAndRecvAllLu(t *testing.T) {

	cv.Convey("select{} should block the goroutine forever, unless also select{default:}", t, func() {

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		// with default: present we should not block
		// _selection = __task.select({{}});
		code := ` a:= 0; go func() { println("top of go-started func"); a = 1; select{ default: }; a= 2; }() // should not block`
		translation, err := inc.Tr([]byte(code))
		//*dbg = true
		fmt.Printf("translation='%s'\n", string(translation))

		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "a", 2)

		//  _r = __task.select({});
		code = ` b:= 0; go func() { b = 1; select{}; b= 2; }() // should block goroutine forever`
		translation, err = inc.Tr([]byte(code))
		panicOn(err)
		fmt.Printf("translation='%s'\n", string(translation))

		LuaRunAndReport(vm, string(translation))
		select {
		case <-time.After(1 * time.Second):
		}
		LuaMustInt64(vm, "b", 1)
		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test901(t *testing.T) {

	cv.Convey("In the all-lua go/coroutine system, ch := make(chan int, 1); ch <- 56;  b := <-ch; write and read back b of 57", t, func() {

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		// with default: present we should not block
		// _selection = __task.select({{}});
		code := ` ch := make(chan int, 1); ch <- 56;  b := <-ch; `
		translation, err := inc.Tr([]byte(code))
		//*dbg = true
		fmt.Printf("translation='%s'\n", string(translation))

		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "b", 56)
	})
}

func Test902(t *testing.T) {

	cv.Convey("spawn goroutine, send and receive on unbuffered channel, in the all-lua system.", t, func() {

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		// with default: present we should not block
		// _selection = __task.select({{}});
		code := ` ch := make(chan int); go func() {ch <- 56;}(); b := <-ch; `
		translation, err := inc.Tr([]byte(code))
		//*dbg = true
		fmt.Printf("translation='%s'\n", string(translation))

		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "b", 56)
	})
}

package compiler

import (
	//"fmt"
	"testing"

	//"github.com/gijit/gi/pkg/verb"
	cv "github.com/glycerine/goconvey/convey"
)

func Test921TimeImports(t *testing.T) {

	cv.Convey(`import "time"; now := time.Now() should work`+
		` to read the time.`, t, func() {

		code := `import "time"`
		code2 := `now := time.Now()`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.trMust([]byte(code))
		LuaRunAndReport(vm, string(translation))

		translation2 := inc.trMust([]byte(code2))
		LuaRunAndReport(vm, string(translation2))
		LuaRunAndReport(vm, "nowNil = (now == nil)")
		LuaMustBool(vm, "nowNil", false)

		cv.So(true, cv.ShouldBeTrue)

	})
}

/* work in progress
func Test922TimeoutsInSelect(t *testing.T) {

	cv.Convey(`channel timeouts in a select statement`, t, func() {

		code := `import "time"`
		code2 := `
ch := make(chan int)
toCount:= 0
go func() {
    for {
       select {
          case <-time.After(time.Millisecond*100):
            if toCount < 3 {
               toCount++
            }
            println("timeout! toCount is now ", toCount)
       }
    }
}()
time.Sleep(5 * time.Millisecond*100)
`
		// toCount should be 3
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.trMust([]byte(code))
		LuaRunAndReport(vm, string(translation))

		translation2 := inc.trMust([]byte(code2))
		LuaRunAndReport(vm, string(translation2))
		LuaMustInt(vm, "toCount", 3)

		cv.So(true, cv.ShouldBeTrue)

	})
}
*/

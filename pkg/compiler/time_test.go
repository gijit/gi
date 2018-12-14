package compiler

import (
	"bytes"
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

func Test923TimeStruct(t *testing.T) {

	cv.Convey(`import "time"; var tm time.Time; should work to instantiate a time struct value.`, t, func() {

		code := `import "time"`
		code2 := `var tm time.Time`
		code3 := `tm = time.Now()`
		code4 := `type S struct { tm time.Time };`
		code5 := `var s S;`
		code6 := `s.tm = time.Now();`
		code7 := `u := s.tm.UnixNano()`
		// a:=time.Now(); b:=time.Now(); d:= b.Sub(a)

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.trMust([]byte(code))
		LuaRunAndReport(vm, string(translation))

		//vv("------------ starting translation2 -----")
		translation2 := inc.trMust([]byte(code2))
		LuaRunAndReport(vm, string(translation2))
		cv.So(string(translation2), matchesLuaSrc, `tm = __type__.time.Time();`)

		//vv("------------ starting translation3 -----")
		translation3 := inc.trMust([]byte(code3))
		LuaRunAndReport(vm, string(translation3))
		cv.So(string(translation3), matchesLuaSrc, `tm = time.Now();`)

		//vv("------------ starting translation4 -----")
		translation4 := inc.trMust([]byte(code4))
		LuaRunAndReport(vm, string(translation4))

		//vv("------------ starting translation5 -----")
		translation5 := inc.trMust([]byte(code5))
		LuaRunAndReport(vm, string(translation5))
		// get the second statement
		cv.So(string(bytes.Split(translation5, []byte("\n\n"))[1]), matchesLuaSrc, `s = __type__.S.ptrToNewlyConstructed(__type__.time.Time());`)

		//vv("------------ starting translation6 -----")
		translation6 := inc.trMust([]byte(code6))
		LuaRunAndReport(vm, string(translation6))

		//vv("------------ starting translation7 -----")
		translation7 := inc.trMust([]byte(code7))
		LuaRunAndReport(vm, string(translation7))

		obsU := LuaToInt64(vm, "u")
		cv.So(obsU, cv.ShouldBeGreaterThan, 0)

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

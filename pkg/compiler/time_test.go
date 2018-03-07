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

/*
time.Sleep(10 * time.Second)
problem in golua_callgofunction, panic happened: '[string "time.Sleep(10000000000LL);"]:1: cannot convert Go function argument #0: cannot convert Lua value 'cdata: 82734008' (cdata) to time.Duration' at
*/

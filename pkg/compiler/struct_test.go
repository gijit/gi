package compiler

import (
	"fmt"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test093NewMethodsShouldBeRegistered(t *testing.T) {

	cv.Convey(`new methods defined on types should be registered with the __reg for the type and be added to the methodset that __reg holds for that type`, t, func() {

		code := `
type S struct{}
func (s *S) hi() string {
   return "hi called!"
}
`
		// __reg:AddMethod should get called.
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, `
	S = __reg:RegisterStruct("S","main","main");
	function S:hi() 
		s = self;
		return "hi called!";
	end;
    __reg:AddMethod("struct", "S", "hi", S.hi)
`)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		//LuaMustInt64(vm, "c0", 4)

	})
}

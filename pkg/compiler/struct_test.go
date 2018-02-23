package compiler

import (
	"fmt"
	cv "github.com/glycerine/goconvey/convey"
	"testing"
)

func Test093NewMethodsShouldBeRegistered(t *testing.T) {

	cv.Convey(`new methods defined on types should be registered with the __reg for the type and be added to the methodset that __reg holds for that type`, t, func() {

		code := `
type S struct{}
func (s *S) hi() string {
   return "hi called!"
}
var s S
h := s.hi()
`
		// __reg:AddMethod should get called.
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.trMust([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		LuaMustString(vm, "h", "hi called!")
		cv.So(true, cv.ShouldBeTrue)

	})
}

func Test120PointersInsideStructs(t *testing.T) {

	cv.Convey(`pointers inside structs should work`, t, func() {

		code := `

    type Ragdoll struct {
	    Andy *Ragdoll
    }
`
		// same should be true
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.trMust([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		// The mutual dependence between __type__.Ragdol and __anon_ptrType
		//  for its Andy *Ragdoll pointer means we can't define
		//  the __constructor in the call to __gi_NewType. So
		//  we pass nil and the later rawset it.

		// exact details get out of date quickly. So only
		// check that it loads without error now.
		LuaRunAndReport(vm, string(translation))
		cv.So(true, cv.ShouldBeTrue)

	})
}

func Test121PointersInsideStructs(t *testing.T) {

	cv.Convey(`pointers inside structs should work`, t, func() {

		code := `

    type Ragdoll struct {
	    Andy *Ragdoll
    }

	var doll Ragdoll
	doll.Andy = &doll
    same := (doll.Andy == &doll)
`
		// comparison of pointers. &doll

		// same should be true
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.trMust([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		LuaMustBool(vm, "same", true)
		cv.So(true, cv.ShouldBeTrue)

	})
}

func Test122ManyPointersInsideStructs(t *testing.T) {

	cv.Convey(`pointers inside structs should work`, t, func() {

		code := `

    type Bunny struct {
           Velvet string
    }

    type Ragdoll struct {
	    Andy *Ragdoll
        bun1  *Bunny
        bun2  *Bunny
    }

	var doll Ragdoll
    bunny1 := &Bunny{}
    bunny2 := bunny1
	doll.Andy = &doll
    doll.bun1 = bunny1
    doll.bun2 = bunny2
    same := (doll.Andy == &doll)
    same2 := (doll.bun1 == doll.bun2)
`
		// same should be true
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.trMust([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		LuaRunAndReport(vm, string(translation))

		LuaMustBool(vm, "same", true)
		LuaMustBool(vm, "same2", true)
		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test123PointersInsideStructStartsNil(t *testing.T) {

	cv.Convey(`pointers inside structs should begin nil; albiet a special nil value that Go code can recognize`, t, func() {

		code := `
    type B struct {
           V *int
    }
    var b B
    a := b.V
    aIsNil := (a == nil)
    var p *int
    pIsNil := (p == nil)
    p == nil
`

		// a should be nil
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.trMust([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		LuaRunAndReport(vm, string(translation))

		//LuaMustBeNil(vm, "p")
		LuaMustBool(vm, "aIsNil", true)
		LuaMustBool(vm, "pIsNil", true)
		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test124ValueFromStructPointer(t *testing.T) {

	cv.Convey(`a value cloned from a struct pointers should have a copy of the members`, t, func() {

		code := `
    type B struct {
           b int
    }
    ptr := &B{b:5}
    inst := *ptr
    mem := inst.b
`

		// mem should be 5
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.trMust([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		LuaRunAndReport(vm, string(translation))

		LuaMustInt(vm, "mem", 5)
		cv.So(true, cv.ShouldBeTrue)
	})
}

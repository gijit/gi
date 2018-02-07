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
var s S
h := s.hi()
`
		// __reg:AddMethod should get called.
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		//cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, ``)

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

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		// The mutual dependence between __type__Ragdol and __anon_ptrType
		//  for its Andy *Ragdoll pointer means we can't define
		//  the __constructor in the call to __gi_NewType. So
		//  we pass nil and the later rawset it.

		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, `
__type__Ragdoll = __gi_NewType(0, __gi_kind_Struct, "main", "Ragdoll", "main.Ragdoll", true, "main", true, nil);

anon_ptrType = __ptrType(__type__Ragdoll); -- 'DELAYED' anon type printing.

__type__Ragdoll.__init("", {{__prop= "Andy", __name= "Andy", __anonymous= false, __exported= true, __typ= anon_ptrType, __tag= ""}});

__type__Ragdoll.__constructor = function(self, ...) 
		 if self == nil then self = {}; end
		 local args={...};
		 if #args == 0 then
			 self.Andy = anon_ptrType.__nil;
		 else 
			 local Andy_ = ... ;
			 self.Andy = Andy_;
		 end
		 return self; 
end;
`)

		// and verify that it happens correctly
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
		// same should be true
		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, `

__type__Ragdoll = __gi_NewType(0, __gi_kind_Struct, "main", "Ragdoll", "main.Ragdoll", true, "main", true, nil);

anon_ptrType = __ptrType(__type__Ragdoll); -- 'DELAYED' anon type printing.

__type__Ragdoll.__init("", {{__prop= "Andy", __name= "Andy", __anonymous= false, __exported= true, __typ= anon_ptrType, __tag= ""}});

__type__Ragdoll.__constructor = function(self, ...) 
		 if self == nil then self = {}; end
		 local args={...};
		 if #args == 0 then
			 self.Andy = anon_ptrType.__nil;
		 else 
			 local Andy_ = ... ;
			 self.Andy = Andy_;
		 end
		 return self; 
end;

doll = __type__Ragdoll.__ptr({}, anon_ptrType.__nil);

doll.Andy = doll;
same = doll.Andy == doll;
`)

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

		translation := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		//cv.So(string(translation), cv.ShouldMatchModuloWhiteSpace, ``)
		LuaRunAndReport(vm, string(translation))

		LuaMustBool(vm, "same", true)
		LuaMustBool(vm, "same2", true)
		cv.So(true, cv.ShouldBeTrue)
	})
}

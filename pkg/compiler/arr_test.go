package compiler

import (
	"fmt"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test080Int64ArraysByGoProxyCopyAppend(t *testing.T) {

	cv.Convey(`[3]int64 should be allocated from Lua and passable to a Go native function. A test of Luar functionality.`, t, func() {

		// a := [3]int64{1,2,3}
		code := `
   import "gitesting"
   a := [3]int64{1,3,4}
   sum := gitesting.SumArrayInt64(a)
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation, err := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		// too fragile
		//cv.So(string(translation), matchesLuaSrc,``)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		// check for exception
		LuaMustInt64(vm, "sum", 8)

	})
}

func Test081CloneOfInt64Array(t *testing.T) {

	cv.Convey(`__gi_clone([3]int64) should return a clone of the array`, t, func() {

		code := `
   a := [3]int64{1,3,4}
   b := a
   c := b[2]
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation, err := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		// too fragile.
		//cv.So(string(translation), matchesLuaSrc,``)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		LuaMustInt64(vm, "c", 4)

	})
}

func Test082IncrementOnInt64Arrays(t *testing.T) {

	cv.Convey(`a := [3]int64{1,3,4}; a[0]++ should leave a[0] at 2.`, t, func() {
		code := `
   import "gitesting"
   a := [3]int64{1,3,4}
   a[0]++
   a[2]--
   b := a[0]
   c := a[2]
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation, err := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		cv.So(string(translation), matchesLuaSrc,
			`
  	__type__.anon_arrayType = __arrayType(__type__.int64, 3); -- 'IMMEDIATE' anon type printing.
  
  	a = __type__.anon_arrayType({[0]=1LL, 3LL, 4LL});
  	__gi_SetRangeCheck(a, 0, (a[0] + (1LL)));
  	__gi_SetRangeCheck(a, 2, (a[2] - (1LL)));
  	b = a[0];
  	c = a[2];
`)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		LuaMustInt64(vm, "b", 2)
		LuaMustInt64(vm, "c", 3)

	})
}

func Test083Int64ArraysCopyByValue(t *testing.T) {

	cv.Convey(`arrays copy by value, so a:= [3]int64{1,3,4}; b := a; should leave b as an independent copy.`, t, func() {

		//      b[i]++ generated wrong code, TODO: fix.
		code := `
   a := [3]int64{0,1,2}
   b := a
   for i := range b {
     b[i] = b[i]+1
   }
   a0 := a[0]
   a1 := a[1]
   a2 := a[2]

   b0 := b[0]
   b1 := b[1]
   b2 := b[2]
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation, err := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		LuaMustInt64(vm, "b0", 1)
		LuaMustInt64(vm, "b1", 2)
		LuaMustInt64(vm, "b2", 3)

		LuaMustInt64(vm, "a0", 0)
		LuaMustInt64(vm, "a1", 1)
		LuaMustInt64(vm, "a2", 2)

	})
}

// working for slices, but not for arrays
func Test084ForRangeOverArrayAndChangeValue(t *testing.T) {

	cv.Convey(`for i := range a { a[i] = a[i] + 1 } should change the value of a[i]`, t, func() {

		code := `
   b := [1]int{0}
   for i := range b {
     b[i] = b[i]+1
   }
   b0 := b[0]
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation, err := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		// check for exception
		LuaMustInt64(vm, "b0", 1)

	})
}

// works, compare to 084 trace
func Test085ForRangeOverSliceAndChangeValue(t *testing.T) {

	cv.Convey(`for i := range a { a[i] = a[i] + 1 } should change the value of a[i]`, t, func() {

		code := `
   b := []int{0}
   for i := range b {
     b[i] = b[i]+1
   }
   b0 := b[0]
`

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()
		inc := NewIncrState(vm, nil)

		translation, err := inc.Tr([]byte(code))
		fmt.Printf("\n translation='%s'\n", translation)

		// and verify that it happens correctly
		LuaRunAndReport(vm, string(translation))

		// check for exception
		LuaMustInt64(vm, "b0", 1)
		cv.So(true, cv.ShouldBeTrue)
	})
}

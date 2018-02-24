package compiler

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gijit/gi/pkg/token"
	"github.com/gijit/gi/pkg/types"
	cv "github.com/glycerine/goconvey/convey"
	golua "github.com/glycerine/golua/lua"
	"github.com/glycerine/luar"
)

func Test600RecvOnChannel(t *testing.T) {

	cv.Convey(`in Lua, receive an integer on a buffered channel, previously sent by Go`, t, func() {

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		ch := make(chan int, 1)
		ch <- 57

		luar.Register(vm, "", luar.Map{
			"ch": ch,
		})

		// first run instantiates the main package so we can add 'ch' to it.
		code := `b := 3`
		inc := NewIncrState(vm, nil)
		translation, err := inc.Tr([]byte(code))
		panicOn(err)
		pp("translation='%s'", string(translation))
		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "b", 3)

		// allow ch to type check
		pkg := inc.pkgMap["main"].Arch.Pkg
		scope := pkg.Scope()
		nt64 := types.Typ[types.Int64]
		chVar := types.NewVar(token.NoPos, pkg, "ch", types.NewChan(types.SendRecv, nt64))
		scope.Insert(chVar)

		code = `a := <- ch;`

		translation, err = inc.Tr([]byte(code))
		panicOn(err)

		pp("translation='%s'", string(translation))

		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "a", 57)
		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test601SendOnChannel(t *testing.T) {

	cv.Convey(`in Lua, send to a buffered channel`, t, func() {

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		ch := make(chan int, 1)

		luar.Register(vm, "", luar.Map{
			"ch": ch,
		})

		// first run instantiates the main package so we can add 'ch' to it.
		code := `b := 3`
		inc := NewIncrState(vm, nil)
		translation, err := inc.Tr([]byte(code))
		panicOn(err)
		pp("translation='%s'", string(translation))
		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "b", 3)

		// allow ch to type check
		pkg := inc.pkgMap["main"].Arch.Pkg
		scope := pkg.Scope()
		nt64 := types.Typ[types.Int64]
		chVar := types.NewVar(token.NoPos, pkg, "ch", types.NewChan(types.SendRecv, nt64))
		scope.Insert(chVar)

		code = `pre:= true; ch <- 6;`

		translation, err = inc.Tr([]byte(code))
		panicOn(err)

		pp("translation='%s'", string(translation))

		LuaRunAndReport(vm, string(translation))

		a := <-ch
		fmt.Printf("a received! a = %v\n", a)
		cv.So(a, cv.ShouldEqual, 6)
		LuaMustBool(vm, "pre", true)
	})
}

func Test602BlockingSendOnChannel(t *testing.T) {

	cv.Convey(`in Lua, send to an unbuffered channel`, t, func() {

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		ch := make(chan int)

		luar.Register(vm, "", luar.Map{
			"ch": ch,
		})

		// first run instantiates the main package so we can add 'ch' to it.
		code := `b := 3`
		inc := NewIncrState(vm, nil)
		translation, err := inc.Tr([]byte(code))
		panicOn(err)
		pp("translation='%s'", string(translation))
		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "b", 3)

		// allow ch to type check
		pkg := inc.pkgMap["main"].Arch.Pkg
		scope := pkg.Scope()
		nt64 := types.Typ[types.Int64]
		chVar := types.NewVar(token.NoPos, pkg, "ch", types.NewChan(types.SendRecv, nt64))
		scope.Insert(chVar)

		var a int
		done := make(chan bool)
		go func() {
			a = <-ch
			close(done)
		}()

		code = `pre:= true; ch <- 7;`

		translation, err = inc.Tr([]byte(code))
		panicOn(err)

		pp("translation='%s'", string(translation))

		LuaRunAndReport(vm, string(translation))

		<-done
		fmt.Printf("a received! a = %v\n", a)
		cv.So(a, cv.ShouldEqual, 7)
		LuaMustBool(vm, "pre", true)
	})
}

func Test603RecvOnUnbufferedChannel(t *testing.T) {

	cv.Convey(`in Lua, receive an integer on an unbuffered channel, previously sent by Go`, t, func() {

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		ch := make(chan int)

		luar.Register(vm, "", luar.Map{
			"ch": ch,
		})

		// first run instantiates the main package so we can add 'ch' to it.
		code := `b := 3`
		inc := NewIncrState(vm, nil)
		translation, err := inc.Tr([]byte(code))
		panicOn(err)
		pp("translation='%s'", string(translation))
		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "b", 3)

		// allow ch to type check
		pkg := inc.pkgMap["main"].Arch.Pkg
		scope := pkg.Scope()
		nt64 := types.Typ[types.Int64]
		chVar := types.NewVar(token.NoPos, pkg, "ch", types.NewChan(types.SendRecv, nt64))
		scope.Insert(chVar)

		code = `a := <- ch;`

		translation, err = inc.Tr([]byte(code))
		panicOn(err)

		pp("translation='%s'", string(translation))

		go func() {
			ch <- 16
		}()
		LuaRunAndReport(vm, string(translation))

		LuaMustInt64(vm, "a", 16)
		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test604Select(t *testing.T) {

	cv.Convey(`in Lua, select on two channels`, t, func() {

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		chInt := make(chan int)
		chStr := make(chan string)

		luar.Register(vm, "", luar.Map{
			"chInt": chInt,
			"chStr": chStr,
		})

		// first run instantiates the main package so we can add 'ch' to it.
		code := `b := 3`
		inc := NewIncrState(vm, nil)
		translation, err := inc.Tr([]byte(code))
		panicOn(err)
		pp("translation='%s'", string(translation))
		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "b", 3)

		// allow channels to type check
		pkg := inc.pkgMap["main"].Arch.Pkg
		scope := pkg.Scope()

		nt := types.Typ[types.Int]
		chIntVar := types.NewVar(token.NoPos, pkg, "chInt", types.NewChan(types.SendRecv, nt))
		scope.Insert(chIntVar)

		stringType := types.Typ[types.String]
		chStrVar := types.NewVar(token.NoPos, pkg, "chStr", types.NewChan(types.SendRecv, stringType))
		scope.Insert(chStrVar)

		code = `
a := 0
b := ""
for i := 0; i < 2; i++ {
  select {
    case a = <- chInt:
    case b = <- chStr:
  }
}
`
		translation, err = inc.Tr([]byte(code))
		panicOn(err)
		pp("translation='%s'", string(translation))

		go func() {
			chInt <- 43
		}()
		go func() {
			chStr <- "hello select"
		}()
		LuaRunAndReport(vm, string(translation))

		LuaMustString(vm, "b", "hello select")
		LuaMustInt64(vm, "a", 43)
		cv.So(true, cv.ShouldBeTrue)
	})
}

func Test605Select(t *testing.T) {

	cv.Convey(`in Lua, select with send and recv`, t, func() {

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		chInt := make(chan int)
		chStr := make(chan string)

		luar.Register(vm, "", luar.Map{
			"chInt": chInt,
			"chStr": chStr,
		})

		// first run instantiates the main package so we can add 'ch' to it.
		code := `b := 3`
		inc := NewIncrState(vm, nil)
		translation, err := inc.Tr([]byte(code))
		panicOn(err)
		pp("translation='%s'", string(translation))
		LuaRunAndReport(vm, string(translation))
		LuaMustInt64(vm, "b", 3)

		// allow channels to type check
		pkg := inc.pkgMap["main"].Arch.Pkg
		scope := pkg.Scope()

		nt := types.Typ[types.Int]
		chIntVar := types.NewVar(token.NoPos, pkg, "chInt", types.NewChan(types.SendRecv, nt))
		scope.Insert(chIntVar)

		stringType := types.Typ[types.String]
		chStrVar := types.NewVar(token.NoPos, pkg, "chStr", types.NewChan(types.SendRecv, stringType))
		scope.Insert(chStrVar)

		code = `
a := 0
b := 0
c := 0
igot := 0
done := false

for a ==0 || b==0 || c==0 {
  select {
    case igot = <- chInt:
       a=1
    case chStr <- "yumo":
       //send happend
       b=1
    default:
       c=1
  }
}
done = true
`
		translation, err = inc.Tr([]byte(code))
		panicOn(err)
		pp("translation='%s'", string(translation))

		go func() {
			chInt <- 43
		}()
		strReceived := ""
		go func() {
			strReceived = <-chStr
		}()
		LuaRunAndReport(vm, string(translation))

		cv.So(strReceived, cv.ShouldEqual, "yumo")

		LuaMustInt64(vm, "a", 1)
		LuaMustInt64(vm, "b", 1)
		LuaMustInt64(vm, "c", 1)
		LuaMustBool(vm, "done", true)
		LuaMustInt64(vm, "igot", 43)
	})
}

func Test606MakeChannel(t *testing.T) {

	cv.Convey(`in Lua, make a new channel`, t, func() {

		vm, err := NewLuaVmWithPrelude(nil)
		panicOn(err)
		defer vm.Close()

		code := `ch := make(chan int, 1); ch <- 23; a := <- ch`
		inc := NewIncrState(vm, nil)
		translation, err := inc.Tr([]byte(code))
		panicOn(err)
		//*dbg = true
		pp("translation='%s'", string(translation))
		LuaRunAndReport(vm, string(translation))

		LuaMustInt64(vm, "a", 23)

		// get the 'ch' channel out and use it in Go.
		varname := "ch"
		vm.GetGlobal(varname)
		top := vm.GetTop()
		if vm.IsNil(top) {
			panic(fmt.Sprintf("global variable '%s' is nil", varname))
		}
		// is it a table or a cdata. if table, look for t.__native
		// to get the actual Go channel.

		*dbg = true
		pp("before optional unwrappng, stack: '%s'", DumpLuaStackAsString(vm))

		// write method to get the channel out of the vm
		//  and interact with it in Go

		t := vm.Type(top)
		switch t {
		case golua.LUA_TTABLE:
			vm.GetField(top, "__native")
			if vm.IsNil(-1) {
				panic("no __native field, ch was not a table-wrapped channel")
			}
			vm.Remove(-2)
		case golua.LUA_TUSERDATA:
			// okay
		default:
			panic("expected table-enclosed Go channel or direct USERDATA with channel pointer")
		}

		pp("after (unwrapping optionally) stack: '%s'", DumpLuaStackAsString(vm))

		top = vm.GetTop()
		var i interface{}
		//var ch chan int
		//pch := &ch
		//ich := interface{}(pch)
		//ich := reflect.ValueOf(pch)
		_, err = luar.LuaToGo(vm, top, &i)
		panicOn(err)

		// i = '&reflect.Value{typ:(*reflect.rtype)(0x45c8d40), ptr:(unsafe.Pointer)(0xc4200620e0), flag:0x12}'
		fmt.Printf("i = '%#v'\n", i)

		//rf = reflect.NewAt(rf.Type(),
		var _ = reflect.Copy
		fmt.Printf("i.Type() = '%[1]T'/'%[1]#v'\n", i.(*reflect.Value).Type())

		// temp:
		vm.Pop(1)
		// do a positive control first: push what we know is a channel, and get it back out.

		pp("stack prior to MakeChan: '%s'", DumpLuaStackAsString(vm))

		// MakeChan creates a 'chan interface{}' proxy and pushes it on the stack.
		// Optional argument: size (number)
		// Returns: proxy (chan interface{})
		vm.PushNumber(1)

		pp("stack after PushNumber(1): '%s'", DumpLuaStackAsString(vm))

		luar.MakeChan(vm)
		vm.SetGlobal("c2")
		vm.GetGlobal("c2")
		var ii interface{}

		pp("stack prior to LuaToGo: '%s'", DumpLuaStackAsString(vm))

		vm.Remove(-2)
		pp("stack after Remove(-2) to get rid of the channel size: '%s'", DumpLuaStackAsString(vm))

		_, err = luar.LuaToGo(vm, -1, &ii)
		panicOn(err)

		pp("stack after to LuaToGo: '%s'", DumpLuaStackAsString(vm))

		pp("ii = '%#v'", ii) // ii = '(chan interface {})(0xc420066e40)'

		var pci chan interface{} = ii.(chan interface{})

		pp("pci = '%#v'", pci)
		// verify it is a buffered chan
		pci <- 1
		get := <-pci
		cv.So(get, cv.ShouldEqual, 1)
	})
}

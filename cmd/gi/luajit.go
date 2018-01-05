package main

import (
	"bufio"
	"fmt"
	"github.com/glycerine/luajit"
	"github.com/go-interpreter/gi/pkg/compiler"
	"io"
	"os"
	"strings"
)

func (cfg *GIConfig) LuajitMain() {

	vm := luajit.Newstate()
	defer vm.Close()

	vm.Openlibs()

	setupPrelude(vm)

	inc := compiler.NewIncrState()
	_ = inc
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("gi> ")
		src, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			fmt.Printf("[EOF]\n")
			return
		}
		panicOn(err)
		if isPrefix {
			panic("line too long")
		}
		var use string

		if !cfg.RawLua {
			translation, err := translateAndCatchPanic(inc, src)
			if err != nil {
				fmt.Printf("oops: '%v' on input '%s'\n", err, string(src))
				translation = "\n"
				// still write, so we get another prompt
			} else {
				p("got translation of line from Go into lua: '%s'\n", strings.TrimSpace(string(translation)))
			}
			use = translation

		} else {
			use = string(src) + "\n"
		}

		//fmt.Printf("sending use='%v'\n", use)
		err = vm.Loadstring(use)
		if err != nil {
			fmt.Printf("error from Lua vm.LoadString(): '%v'. supplied lua with: '%s'\nlua stack:\n", err, use[:len(use)-1])
			DumpLuaStack(vm)
			vm.Pop(1)
			continue
		}
		err = vm.Pcall(0, 0, 0)
		panicOn(err)
		DumpLuaStack(vm)
		reader.Reset(os.Stdin)
	}
}

func DumpLuaStack(L *luajit.State) {
	var top int

	top = L.Gettop()
	for i := 1; i <= top; i++ {
		t := L.Type(i)
		switch t {
		case luajit.Tstring:
			fmt.Println("String : \t", L.Tostring(i))
		case luajit.Tboolean:
			fmt.Println("Bool : \t\t", L.Toboolean(i))
		case luajit.Tnumber:
			fmt.Println("Number : \t", L.Tonumber(i))
		default:
			fmt.Println("Type : \t\t", L.Typename(i))
		}
	}
	print("\n")
}

func setupPrelude(vm *luajit.State) {
	err := vm.Loadstring(compiler.GiLuaSliceMap)
	panicOn(err)
	err = vm.Pcall(0, 0, 0)
	panicOn(err)
}

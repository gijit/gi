package main

import (
	"bufio"
	"fmt"
	"github.com/glycerine/luajit"
	"github.com/go-interpreter/gi/pkg/incr/compiler"
	"io"
	"os"
)

func LuajitMain() {

	vm := luajit.Newstate()
	defer vm.Close()

	vm.Openlibs()

	inc := compiler.NewIncrState()
	_ = inc
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		src, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			fmt.Printf("[EOF]\n")
			return
		}
		panicOn(err)
		if isPrefix {
			panic("line too long")
		}
		translation := inc.Tr([]byte(src))
		fmt.Printf("go:'%s'  -->  '%s' in lua\n", src, translation)

		err = vm.Loadstring(string(translation))
		//err = vm.Loadstring(string(src))
		panicOn(err)
		err = vm.Pcall(0, 0, 0)
		panicOn(err)
		DumpLuaStack(vm)
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

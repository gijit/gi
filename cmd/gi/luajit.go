package main

import (
	"bufio"
	"fmt"
	"github.com/glycerine/luajit"
	"github.com/go-interpreter/gi/pkg/compiler"
	"github.com/go-interpreter/gi/pkg/verb"
	"io"
	"os"
	"strings"
)

func (cfg *GIConfig) LuajitMain() {

	vm := luajit.Newstate()
	defer vm.Close()

	vm.Openlibs()

	err := cfg.setupPrelude(vm)
	panicOn(err)

	inc := compiler.NewIncrState()
	_ = inc
	reader := bufio.NewReader(os.Stdin)
	goPrompt := "gi> "
	luaPrompt := "raw luajit gi> "
	prompt := goPrompt
	if cfg.RawLua {
		prompt = luaPrompt
	}

	for {
		fmt.Printf(prompt)
		src, isPrefix, err := reader.ReadLine()
		if err == io.EOF {
			fmt.Printf("[EOF]\n")
			return
		}
		panicOn(err)
		if isPrefix {
			panic("line too long")
		}
		use := string(src)
		cmd := strings.TrimSpace(use)
		switch cmd {
		case ":v":
			verb.Verbose = true
			verb.VerboseVerbose = false
			continue
		case ":vv":
			verb.Verbose = true
			verb.VerboseVerbose = true
			continue
		case ":raw":
			cfg.RawLua = true
			prompt = luaPrompt
			fmt.Printf("Raw LuaJIT language mode.\n")
			continue
		case ":go":
			cfg.RawLua = false
			prompt = goPrompt
			fmt.Printf("Go language mode.\n")
			continue
		case ":help":
			fmt.Printf(`
======================
gi: a go interpreter
https://github.com/go-interpreter/gi
command prompt help: 
simply type Go expressions or statements
directly at the prompt, or use one of 
these special commands:
======================
 :v  turns on verbose debug prints
 :vv turns on very verbose prints
 :raw changes to raw-luajit entry mode
 :go  change back from raw mode to Go mode
 ctrl-d to exit
`)
			continue
		}

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
			use += "\n"
		}

		p("sending use='%v'\n", use)
		err = vm.Loadstring(use)
		if err != nil {
			fmt.Printf("error from Lua vm.LoadString(): '%v'. supplied lua with: '%s'\nlua stack:\n", err, use[:len(use)-1])
			DumpLuaStack(vm)
			vm.Pop(1)
			continue
		}
		err = vm.Pcall(0, 0, 0)
		if err != nil {
			fmt.Printf("error from Lua vm.Pcall(0,0,0): '%v'. supplied lua with: '%s'\nlua stack:\n", err, use[:len(use)-1])
			DumpLuaStack(vm)
			vm.Pop(1)
			continue
		}
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

func DumpLuaStackAsString(L *luajit.State) string {
	var top int
	s := ""
	top = L.Gettop()
	for i := 1; i <= top; i++ {
		t := L.Type(i)
		switch t {
		case luajit.Tstring:
			s += fmt.Sprintf("String : \t%v", L.Tostring(i))
		case luajit.Tboolean:
			s += fmt.Sprintf("Bool : \t\t%v", L.Toboolean(i))
		case luajit.Tnumber:
			s += fmt.Sprintf("Number : \t%v", L.Tonumber(i))
		default:
			s += fmt.Sprintf("Type : \t\t%v", L.Typename(i))
		}
	}
	return s
}

func (cfg *GIConfig) setupPrelude(vm *luajit.State) error {
	for _, f := range cfg.preludeFiles {
		err := vm.Loadstring(fmt.Sprintf(`dofile("%s")`, f))
		if err != nil {
			msg := DumpLuaStackAsString(vm)
			vm.Pop(1)
			return fmt.Errorf("error in setupPrelude during LoadString on file '%s': '%v'. Details: '%s'", f, err, msg)
		}
		err = vm.Pcall(0, 0, 0)
		if err != nil {
			msg := DumpLuaStackAsString(vm)
			vm.Pop(1)
			return fmt.Errorf("error in setupPrelude during Pcall on file '%s': '%v'. Details: '%s'", f, err, msg)
		}
	}
	return nil
}

package main

import (
	"bufio"
	"fmt"
	"github.com/glycerine/luajit"
	"github.com/go-interpreter/gi/pkg/compiler"
	"github.com/go-interpreter/gi/pkg/verb"
	"io"
	"os"
	"strconv"
	"strings"
)

func (cfg *GIConfig) LuajitMain() {

	vm := luajit.Newstate()
	defer vm.Close()

	vm.Openlibs()

	err := compiler.LuaDoFiles(vm, cfg.preludeFiles)
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
		low := strings.ToLower(cmd)
		switch low {
		case ":ast":
			inc.PrintAST = true
			continue
		case ":noast":
			inc.PrintAST = false
			continue
		case ":q":
			fmt.Printf("quiet mode\n")
			verb.Verbose = false
			verb.VerboseVerbose = false
			continue
		case ":v":
			fmt.Printf("verbose mode.\n")
			verb.Verbose = true
			verb.VerboseVerbose = false
			continue
		case ":vv":
			fmt.Printf("very verbose mode.\n")
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
		case ":prelude":
			fmt.Printf("Reloading prelude...\n")
			err := compiler.LuaDoFiles(vm, cfg.preludeFiles)
			if err != nil {
				fmt.Printf("error during prelude reload: '%v'", err)
			}
			continue
		case ":help", ":h":
			fmt.Printf(`
======================
gi: a go interpreter
https://github.com/go-interpreter/gi
command prompt help: 
simply type Go expressions or statements
directly at the prompt, or use one of 
these special commands:
======================
 :v          turns on verbose debug printing
 :vv         turns on very verbose printing
 :q          quiets the debug prints (default)
 :raw        changes to raw-luajit entry mode
 :go         change back from raw mode to Go mode
 :ast        print the Go AST prior to translation
 :noast      stop printing the Go AST
 :do <path>  run dofile(path) on a .lua file
 :h          show this help (same as :help)
 ctrl-d to exit
`)
			continue
		}

		if strings.HasPrefix(low, ":do") {
			files := strings.TrimSpace(low[3:])
			splt := strings.Split(files, ",")
			var final, show []string
			for i := range splt {
				tmp := strings.TrimSpace(splt[i])
				home := os.Getenv("HOME")
				if home != "" {
					tmp = strings.Replace(tmp, "~/", home+"/", 1)
				}
				if len(tmp) > 0 {
					final = append(final, tmp)
					show = append(show, strconv.Quote(tmp))
				}
			}
			if len(final) > 0 {
				fmt.Printf("running dofile(%s)\n", strings.Join(show, ","))
				err := compiler.LuaDoFiles(vm, final)
				if err != nil {
					fmt.Printf("error during dofile(): '%v'\n", err)
				}
			} else {
				fmt.Printf("nothing to do.\n")
			}
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

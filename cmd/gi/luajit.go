package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gijit/gi/pkg/compiler"
	"github.com/gijit/gi/pkg/front"
	"github.com/gijit/gi/pkg/verb"
	//luajit "github.com/glycerine/golua/lua"
)

func (cfg *GIConfig) LuajitMain() {

	vmCfg := compiler.NewVmConfig()
	vmCfg.PreludePath = cfg.PreludePath
	vmCfg.Quiet = cfg.Quiet
	vmCfg.NotTestMode = !cfg.IsTestMode
	vm, err := compiler.NewLuaVmWithPrelude(vmCfg)
	panicOn(err)
	defer vm.Close()
	inc := compiler.NewIncrState(vm, vmCfg)

	var t0, t1 time.Time
	var history []string
	home := os.Getenv("HOME")
	var histFn string
	var histFile *os.File
	var sessionStartAfter int
	if home != "" {
		histFn = home + string(os.PathSeparator) + ".gijit.hist"

		// open and close once to read back history
		history, err = readHistory(histFn)
		lh := len(history)
		if lh > 0 {
			sessionStartAfter = lh
		}
		panicOn(err)

		// re-open for append new history
		histFile, err = os.OpenFile(histFn,
			os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC,
			0600)
		panicOn(err)
	}

	_ = inc
	reader := bufio.NewReader(os.Stdin)
	goPrompt := "gi> "
	goMorePrompt := ">>>    "
	luaPrompt := "raw luajit gi> "
	isDo := false
	isSource := false

	var prompter *Prompter
	if !cfg.NoLiner {
		prompter = NewPrompter(goPrompt)
		for i := range history {
			prompter.prompter.AppendHistory(history[i])
		}
	}

	prompt := goPrompt
	if cfg.RawLua {
		prompt = luaPrompt
	}
	prevSrc := ""
	var prompterLine string
	var by []byte

top:
	for {
		if cfg.NoLiner {
			fmt.Printf(prompt)
			by, err = reader.ReadBytes('\n')
		} else {
			prompterLine, err = prompter.Getline(&prompt)
			by = []byte(prompterLine)
		}
		if err == io.EOF {
			if len(by) > 0 {
				fmt.Printf("\n on EOF, but len(by) = %v, by='%s'", len(by), string(by))
				// process bytes first,
				// return next time.
				err = nil
			} else {
				fmt.Printf("[EOF]\n")
				return
			}
		}
		panicOn(err)
		use := string(by)
		src := use
		cmd := bytes.TrimSpace(by)
		low := string(bytes.ToLower(cmd))
		if len(low) > 1 && low[0] == ':' {
			if low[:2] == "::" {
				// likely the start of a lua label for a goto, not a special : command.
				goto notColonCmd
			}
			if low[1] == '-' || (low[1] >= '0' && low[1] <= '9') {
				// replay history, one command, or a range.

				// check for range
				num, err := getHistoryRange(low[1:], history)
				if err != nil {
					fmt.Printf("%s\n", err.Error())
					continue top
				}

				switch len(num) {
				case 1:
					fmt.Printf("replay history %03d:\n", num[0])
					src = history[num[0]-1]
					fmt.Printf("%s\n", src)
				case 2:
					if num[1] < num[0] {
						fmt.Printf("bad history request, end before beginning.\n")
						continue top
					}
					fmt.Printf("replay history %03d - %03d:\n", num[0], num[1])
					src = strings.Join(history[num[0]-1:num[1]], "\n") + "\n"
					fmt.Printf("%s\n", src)
				}
			}
		}
		if len(low) > 3 && low[:3] == ":rm" {
			// remove some commands from history
			var beg, end int
			history, histFile, beg, end, err = removeCommands(history, histFn, histFile, low[3:])
			if err != nil {
				fmt.Printf("%s\n", err.Error())
			}
			if end >= 0 {
				delcount := (end - beg + 1)
				if end < sessionStartAfter {
					// deleted history before our session, adjust marker
					sessionStartAfter -= delcount

				} else if beg < sessionStartAfter {
					// delete history crosses into our session
					sessionStartAfter = beg
				}
			}
			continue top
		}
		switch low {
		case ":ast":
			inc.PrintAST = true
			continue top
		case ":noast":
			inc.PrintAST = false
			continue top
		case ":q":
			fmt.Printf("quiet mode\n")
			verb.Verbose = false
			verb.VerboseVerbose = false
			continue top
		case ":v":
			fmt.Printf("verbose mode.\n")
			verb.Verbose = true
			verb.VerboseVerbose = false
			continue top
		case ":vv":
			fmt.Printf("very verbose mode.\n")
			verb.Verbose = true
			verb.VerboseVerbose = true
			continue top
		case ":clear", ":reset":
			history = history[:0]
			if histFn != "" {
				histFile, err = os.OpenFile(histFn,
					os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_SYNC,
					0600)
				panicOn(err)
			}
			sessionStartAfter = 0
			fmt.Printf("history cleared.\n")
			continue top
		case ":h":
			if len(history) == 0 {
				fmt.Printf("history: empty\n")
				fmt.Printf("----- current session: -----\n")
				continue top
			}
			fmt.Printf("history:\n")
			newline := "\n"
			if sessionStartAfter == 0 {
				fmt.Printf("----- current session: -----\n")
			}
			for i, h := range history {
				lenh := len(h)
				switch {
				case lenh == 0:
					newline = "\n"
				case h[lenh-1] == '\n':
					newline = ""
				default:
					newline = "\n"
				}
				fmt.Printf("%03d: %s%s", i+1, h, newline)
				if i+1 == sessionStartAfter {
					fmt.Printf("----- current session: -----\n")
				}
			}
			fmt.Printf("\n")
			continue top
		case ":raw", ":r":
			cfg.RawLua = true
			prompt = luaPrompt
			fmt.Printf("Raw LuaJIT language mode.\n")
			continue top
		case ":go", ":g", ":":
			cfg.RawLua = false
			prompt = goPrompt
			fmt.Printf("Go language mode.\n")
			continue top
		case ":prelude", ":reload":
			fmt.Printf("Reloading prelude...\n")

			files, err := compiler.FetchPreludeFilenames(cfg.PreludePath, cfg.Quiet)
			if err != nil {
				fmt.Printf("error during prelude reload: '%v'", err)
				continue top
			}
			err = compiler.LuaDoFiles(vm, files)
			if err != nil {
				fmt.Printf("error during prelude reload: '%v'", err)
			}
			continue top
		case ":help", ":?":
			fmt.Printf(`
======================
gijit: a go interpreter, just-in-time
https://github.com/gijit/gi
command prompt help: 
simply type Go expressions or statements
directly at the prompt, or use one of 
these special commands:
======================
 :v          turns on verbose debug printing
 :vv         turns on very verbose printing
 :q          quiets the debug prints (default)
 :r or :raw  change to raw-luajit entry mode
 :g or :go   change back from raw to default Go mode 
 :ast        print the Go AST prior to translation
 :noast      stop printing the Go AST
 :?          show this help (:help does the same)
 :h          show command line history
 :30         replay command number 30 from history
 :1-10       replay commands 1 - 10 inclusive
 :reset      reset and clear history (also :clear)
 :rm 3-4     remove commands 3-4 from history
 :do <path>  run dofile(path) on a .lua file
 :source <path>   re-play/source Go lines from a file
 ctrl-d to exit
`)
			continue top
		}

		isDo = strings.HasPrefix(low, ":do")
		isSource = strings.HasPrefix(low, ":source")
		if isDo || isSource {
			off := 3
			action := "running dofile"
			nm := "dofiles"
			if isSource {
				off = 7
				action = "sourcing Go from"
				nm = "source Go files"
			}

			files := strings.TrimSpace(low[off:])
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
			var err error
			if len(final) > 0 {
				fmt.Printf("%s (%s)\n", nm, strings.Join(show, ","))
				if isDo {
					err = compiler.LuaDoFiles(vm, final)
				} else {
					by, err = sourceGoFiles(final)
					if err != nil {
						fmt.Printf("error during %s: '%v'\n", action, err)
					} else {
						src = string(by)
						goto notColonCmd
					}
				}
				if err != nil {
					fmt.Printf("error during %s: '%v'\n", action, err)
				}
			} else {
				fmt.Printf("nothing to do.\n")
			}
			continue top
		}

	notColonCmd:
		isContinuation := len(prevSrc) > 0
		if !cfg.RawLua {
			if isContinuation {
				src = prevSrc + "\n" + src
			}
			//fmt.Printf("src = '%s'\n", src)
			//fmt.Printf("prevSrc = '%s'\n", prevSrc)

			eof, syntaxErr, empty := front.TopLevelParseGoSource([]byte(src))
			if empty {
				prevSrc = ""
				continue top
			}
			//fmt.Printf("eof = %v, syntaxErr = %v\n", eof, syntaxErr)
			if eof && !syntaxErr {
				prompt = goMorePrompt
				// get another line of input
				prevSrc = src
				continue top
			}
			prevSrc = ""

			prompt = goPrompt
			translation, err := translateAndCatchPanic(inc, []byte(src))
			if err != nil {
				fmt.Printf("oops: '%v' on input '%s'\n", err, strings.TrimSpace(string(src)))
				translation = "\n"
				// still write, so we get another prompt
			} else {
				p("got translation of line from Go into lua: '%s'\n", strings.TrimSpace(string(translation)))
			}
			use = translation

		} else {
			// :r/raw mode
			use = src
		}

		p("sending use='%v'\n", use)
		srcLines := strings.Split(src, "\n")
		//fmt.Printf("appending to history: src='%#v', srcLines='%#v'\n", src, srcLines)
		lensrc := len(srcLines)
		histBeg := len(history)
		if lensrc > 1 && strings.TrimSpace(srcLines[lensrc-1]) == "" {
			history = append(history, srcLines[:lensrc-1]...)
		} else {
			history = append(history, srcLines[:len(srcLines)]...)
		}
		histEnd := len(history)
		if histFile != nil {
			for i := histBeg; i < histEnd; i++ {
				fmt.Fprintf(histFile, "%s\n", history[i])
			}
			histFile.Sync()
		}
		t0 = time.Now()
		// 	loadstring: returns 0 if there are no errors or 1 in case of errors.
		interr := vm.LoadString(use)
		if interr != 0 {
			fmt.Printf("error from Lua vm.LoadString(): supplied lua with: '%s'\nlua stack:\n", use[:len(use)-1])
			compiler.DumpLuaStack(vm)
			vm.Pop(1)
			continue top
		}
		err = vm.Call(0, 0)
		if err != nil {
			fmt.Printf("error from Lua vm.Call(0,0): '%v'. supplied lua with: '%s'\nlua stack:\n", err, use[:len(use)-1])
			compiler.DumpLuaStack(vm)
			vm.Pop(1)
			continue top
		}
		t1 = time.Now()
		// jea debug:
		//compiler.DumpLuaStack(vm)
		fmt.Printf("\n")
		reader.Reset(os.Stdin)
		fmt.Printf("elapsed: '%v'\n", t1.Sub(t0))
	}
}

func sourceGoFiles(files []string) ([]byte, error) {
	var buf bytes.Buffer
	for _, f := range files {
		fd, err := os.Open(f)
		if err != nil {
			return nil, err
		}
		defer fd.Close()
		by, err := ioutil.ReadAll(fd)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(&buf, bytes.NewBuffer(by))
		if err != nil {
			return nil, err
		}
		// newline between files.
		fmt.Fprintf(&buf, "\n")
	}
	bb := buf.Bytes()
	fmt.Printf("sourceGoFiles() returning '%s'\n", string(bb))
	return bb, nil
}

package compiler

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gijit/gi/pkg/front"
	"github.com/gijit/gi/pkg/verb"
	golua "github.com/glycerine/golua/lua"
)

var p = verb.P

func (cfg *GIConfig) LuajitMain() {
	r := NewRepl(cfg)
	defer r.vm.Close()
	r.Loop()
}

type Repl struct {
	inc   *IncrState
	vmCfg *GIConfig
	vm    *golua.State

	t0 time.Time
	t1 time.Time

	history  []string
	home     string
	histFn   string
	histFile *os.File

	sessionStartAfter int

	goPrompt     string
	goMorePrompt string
	luaPrompt    string
	calcPrompt   string
	isDo         bool
	isSource     bool

	prompter *Prompter
	cfg      *GIConfig
	prompt   string

	prevSrc      string
	prompterLine string
	reader       *bufio.Reader
}

func NewRepl(cfg *GIConfig) *Repl {

	vm, err := NewLuaVmWithPrelude(cfg)

	panicOn(err)
	inc := NewIncrState(vm, cfg)

	r := &Repl{cfg: cfg, vm: vm, inc: inc}
	r.home = os.Getenv("HOME")
	if r.home != "" {
		r.histFn = r.home + string(os.PathSeparator) + ".gijit.hist"

		// open and close once to read back history
		r.history, err = readHistory(r.histFn)
		lh := len(r.history)
		if lh > 0 {
			r.sessionStartAfter = lh
		}
		panicOn(err)

		// re-open for append new history
		r.histFile, err = os.OpenFile(r.histFn,
			os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC,
			0600)
		panicOn(err)
	}

	r.reader = bufio.NewReader(os.Stdin)
	r.goPrompt = "gi> "
	r.calcPrompt = "calc mode> "
	//r.goMorePrompt = ">>>    "
	r.luaPrompt = "raw luajit gi> "
	r.isDo = false
	r.isSource = false

	if !r.cfg.NoLiner {
		r.prompter = NewPrompter(r.goPrompt)
		for i := range r.history {
			r.prompter.prompter.AppendHistory(r.history[i])
		}
	}
	r.setPrompt()
	r.prevSrc = ""
	r.prompterLine = ""
	return r
}

func (r *Repl) Loop() {
	for {
		src, err := r.Read()
		if err == io.EOF {
			return
		}

		err = r.Eval(src)
		if err == io.EOF {
			return
		}
	}
}

func (r *Repl) Read() (src string, err error) {

	var by []byte

readtop:
	if r.cfg.NoLiner {
		if r.prompt != "" {
			fmt.Printf(r.prompt)
		}
		by, err = r.reader.ReadBytes('\n')
	} else {
		r.prompterLine, err = r.prompter.Getline(&(r.prompt))
		by = []byte(r.prompterLine)
	}
	if err == io.EOF {
		if len(by) > 0 {
			fmt.Printf("\n on EOF, but len(by) = %v, by='%s'", len(by), string(by))
			// process bytes first,
			// return next time.
			return
		} else {
			fmt.Printf("[EOF]\n")
			return "", err
		}
	}
	panicOn(err)
	use := string(by)
	src = use
	cmd := bytes.TrimSpace(by)
	low := string(bytes.ToLower(cmd))
	if len(low) > 1 && low[0] == ':' {
		if low[:2] == "::" {
			// likely the start of a lua label for a goto, not a special : command.
			return
		}
		if low[1] == '-' || (low[1] >= '0' && low[1] <= '9') {
			// replay history, one command, or a range.

			// check for range
			num, err := getHistoryRange(low[1:], r.history)
			if err != nil {
				fmt.Printf("%s\n", err.Error())
				return "", err
			}

			switch len(num) {
			case 1:
				fmt.Printf("replay history %03d:\n", num[0])
				src = r.history[num[0]-1]
				fmt.Printf("%s\n", src)
			case 2:
				if num[1] < num[0] {
					fmt.Printf("bad history request, end before beginning.\n")
					return "", nil
				}
				fmt.Printf("replay history %03d - %03d:\n", num[0], num[1])
				src = strings.Join(r.history[num[0]-1:num[1]], "\n") + "\n"
				fmt.Printf("%s\n", src)
			}
		}
	}
	if len(low) > 3 && low[:3] == ":rm" {
		// remove some commands from history
		var beg, end int
		r.history, r.histFile, beg, end, err = removeCommands(r.history, r.histFn, r.histFile, low[3:])
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}
		if end >= 0 {
			delcount := (end - beg + 1)
			if end < r.sessionStartAfter {
				// deleted history before our session, adjust marker
				r.sessionStartAfter -= delcount

			} else if beg < r.sessionStartAfter {
				// delete history crosses into our session
				r.sessionStartAfter = beg
			}
		}
		return "", nil
	}
	switch low {
	case ":ast":
		r.inc.PrintAST = true
		return "", nil
	case ":noast":
		r.inc.PrintAST = false
		return "", nil
	case ":q":
		fmt.Printf("quiet mode\n")
		verb.Verbose = false
		verb.VerboseVerbose = false
		return "", nil
	case ":v":
		fmt.Printf("verbose mode.\n")
		verb.Verbose = true
		verb.VerboseVerbose = false
		return "", nil
	case ":vv":
		fmt.Printf("very verbose mode.\n")
		verb.Verbose = true
		verb.VerboseVerbose = true
		return "", nil
	case ":clear", ":reset":
		r.history = r.history[:0]
		if r.histFn != "" {
			r.histFile, err = os.OpenFile(r.histFn,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_SYNC,
				0600)
			panicOn(err)
		}
		r.sessionStartAfter = 0
		fmt.Printf("history cleared.\n")
		return "", nil
	case ":h":
		if len(r.history) == 0 {
			fmt.Printf("history: empty\n")
			fmt.Printf("----- current session: -----\n")
			return "", nil
		}
		fmt.Printf("history:\n")
		newline := "\n"
		if r.sessionStartAfter == 0 {
			fmt.Printf("----- current session: -----\n")
		}
		for i, h := range r.history {
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
			if i+1 == r.sessionStartAfter {
				fmt.Printf("----- current session: -----\n")
			}
		}
		fmt.Printf("\n")
		return "", nil
	case ":ls":
		r.displayCmd(`ls`)
		goto readtop
	case ":gls":
		r.displayCmd(`gls`)
		goto readtop
	case ":glst":
		r.displayCmd(`glst`)
		goto readtop
	case ":lst":
		r.displayCmd(`lst`)
		goto readtop
	case ":r":
		r.cfg.RawLua = true
		r.cfg.CalculatorMode = false
		r.prompt = r.luaPrompt
		fmt.Printf("Raw LuaJIT language mode.\n")
		goto readtop

	case ":go", ":g", ":":
		r.cfg.RawLua = false
		r.cfg.CalculatorMode = false
		r.prompt = r.goPrompt
		fmt.Printf("Go language mode.\n")
		return "", nil

	case "==":
		r.cfg.RawLua = false
		r.cfg.CalculatorMode = true
		fmt.Printf("Calculator mode.\n")
		r.prompt = r.calcPrompt
		return "", nil

	case ":prelude", ":reload":
		fmt.Printf("Reloading prelude...\n")

		files, err := FetchPreludeFilenames(r.cfg.PreludePath, r.cfg.Quiet)
		if err != nil {
			fmt.Printf("error during prelude reload: '%v'", err)
			return "", err
		}
		err = LuaDoFiles(r.vm, files)
		if err != nil {
			fmt.Printf("error during prelude reload: '%v'", err)
		}
		return "", nil
	case ":help", ":?":
		fmt.Printf(`
======================
gijit: a go interpreter, just-in-time
https://github.com/gijit/gi
Type Go expressions or statements
directly after the prompt, or use one of 
these special commands.
======================
 :v              Turn on verbose debug printing.
 :vv             Turn on very verbose printing.
 :q              Quiet the debug prints (default).
 :r              Change to raw-luajit Lua entry mode.
 :g or :go       Change back from raw to default Go mode.
 :ast            Print the Go AST prior to translation.
 :noast          Stop printing the Go AST.
 :?              Show this help (:help does the same).
 :h              Show command line history.
 :30             Replay command number 30 from history.
 :1-10           Replay commands 1 - 10 inclusive.
 :reset          Reset and clear history (also :clear).
 :rm 3-4         Remove commands 3-4 from history.
 :do <path>      Run dofile(path) on a .lua file.
 :source <path>  Re-play Go code from a file.
 :ls             List all global user variables.
 :gls            List all global variables (include __ prefixed).
 = 3 + 4         Calculate the expression after the '=' (one line).
 ==              Multiple entry calculator mode. ':' to exit.
 import "fmt"    Import the binary, pre-compiled package.
 ctrl-d to exit  History is saved in ~/.gitit.hist
`)
		return "", nil
	}

	r.isDo = strings.HasPrefix(low, ":do")
	r.isSource = strings.HasPrefix(low, ":source")
	if r.isDo || r.isSource {
		off := 3
		action := "running dofile"
		nm := "dofiles"
		if r.isSource {
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
			if r.isDo {
				err = LuaDoFiles(r.vm, final)
			} else {
				by, err = sourceGoFiles(final)
				if err != nil {
					fmt.Printf("error during %s: '%v'\n", action, err)
				} else {
					src = string(by)
					return src, nil
				}
			}
			if err != nil {
				fmt.Printf("error during %s: '%v'\n", action, err)
			}
		} else {
			fmt.Printf("nothing to do.\n")
		}
		return "", nil
	}
	return src, nil
}

func (r *Repl) setPrompt() {
	if r.cfg.CalculatorMode {
		r.prompt = r.calcPrompt
		return
	}
	if r.cfg.RawLua {
		r.prompt = r.luaPrompt
		return
	}
	r.prompt = r.goPrompt
}

func (r *Repl) Eval(src string) error {

	var use string
	isContinuation := len(r.prevSrc) > 0
	if !r.cfg.RawLua {
		if isContinuation {
			src = r.prevSrc + "\n" + src
		}
		//fmt.Printf("src = '%s'\n", src)
		//fmt.Printf("prevSrc = '%s'\n", prevSrc)

		eof, syntaxErr, empty, err := front.TopLevelParseGoSource([]byte(src))
		if empty {
			r.prevSrc = ""
			return nil
		}
		//fmt.Printf("eof = %v, syntaxErr = %v\n", eof, syntaxErr)
		if eof && !syntaxErr {
			r.prompt = r.goMorePrompt
			// get another line of input
			r.prevSrc = src
			return nil
		}
		r.prevSrc = ""

		r.setPrompt()
		translation, err := translateAndCatchPanic(r.inc, []byte(src))
		if err != nil {
			fmt.Printf("oops: '%v' on input '%s'\n", err, strings.TrimSpace(src))
			translation = "\n"
			// still write, so we get another prompt

			// hmm, or maybe not
			return err
		} else {
			p("got translation of line from Go into lua: '%s'\n", strings.TrimSpace(string(translation)))
		}
		use = translation

	} else {
		// raw mode, under :r
		use = src
	}

	p("sending use='%v'\n", use)

	// add to history as separate lines
	srcLines := strings.Split(src, "\n")
	//fmt.Printf("appending to history: src='%#v', srcLines='%#v'\n", src, srcLines)
	lensrc := len(srcLines)
	histBeg := len(r.history)
	if lensrc > 1 && strings.TrimSpace(srcLines[lensrc-1]) == "" {
		r.history = append(r.history, srcLines[:lensrc-1]...)
	} else {
		r.history = append(r.history, srcLines[:len(srcLines)]...)
	}
	histEnd := len(r.history)
	if r.histFile != nil {
		for i := histBeg; i < histEnd; i++ {
			fmt.Fprintf(r.histFile, "%s\n", r.history[i])
		}
		r.histFile.Sync()
	}
	r.t0 = time.Now()
	// 	loadstring: returns 0 if there are no errors or 1 in case of errors.
	interr := r.vm.LoadString(use)
	if interr != 0 {
		fmt.Printf("error from Lua vm.LoadString(): supplied lua with: '%s'\nlua stack:\n", use[:len(use)-1])
		DumpLuaStack(r.vm)
		r.vm.Pop(1)
		return nil
	}
	err := r.vm.Call(0, 0)
	if err != nil {
		fmt.Printf("error from Lua vm.Call(0,0): '%v'. supplied lua with: '%s'\nlua stack:\n", err, use[:len(use)-1])
		DumpLuaStack(r.vm)
		r.vm.Pop(1)
		return nil
	}
	r.t1 = time.Now()
	// jea debug:
	//DumpLuaStack(vm)
	fmt.Printf("\n")
	r.reader.Reset(os.Stdin)
	fmt.Printf("elapsed: '%v'\n", r.t1.Sub(r.t0))

	return nil
}

// :ls, :gls, :lst, :glst implementation
func (r *Repl) displayCmd(cmd string) {
	err := LuaDoString(r.vm, `__`+cmd+`()`)
	panicOn(err)
}

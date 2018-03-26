package compiler

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/glycerine/gi/pkg/front"
	"github.com/glycerine/gi/pkg/verb"
	golua "github.com/glycerine/golua/lua"
)

var p = verb.P

// Arrange that main.main runs on main thread.
// We want LuaJIT to always be on the one
// main C thread.
func init() {
	runtime.LockOSThread()
}

// the single queue of work to run in main thread.
var mainQ = make(chan func())
var mainShutdown = make(chan bool)

// doMain runs f on the main thread.
func doMainWait(f func()) {
	done := make(chan bool, 1)
	mainQ <- func() {
		f()
		done <- true
	}
	<-done
}

func doMainAsync(f func()) {
	mainQ <- func() {
		f()
	}
}

// Main runs the main LuaJIT service loop.
// The binary's main.main must call LuajitMain() to run this loop.
// Main does not return. If the binary needs to do other work, it
// must do it in separate goroutines.
func MainCThread() {
	for {
		select {
		case <-mainShutdown:
			return
		case f := <-mainQ:
			f()
		}
	}
}

func (cfg *GIConfig) LuajitMain() {
	if reserveMainThread {
		done := make(chan bool)
		var r *Repl
		go func() {
			r = NewRepl(cfg)

			// in place of defer to cleanup:
			go func() {
				<-done
				r.lvm.Close()
				close(mainShutdown)
			}()
			r.Loop()
			done <- true
		}()

		// this is the C main thread, for LuaJIT
		// to always and only be on... it never
		// returns. Use doMain() to communicate
		// with it via mainQ.
		MainCThread()

	} else {

		r := NewRepl(cfg)
		defer func() {
			r.lvm.Close()
			close(mainShutdown)
		}()
		r.Loop()
	}
}

type Repl struct {
	inc   *IncrState
	vmCfg *GIConfig
	lvm   *LuaVm

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

	// give a config so it knows we are
	// running under `gi` cmd. Otherwise
	// tests will assume nil means they
	// need to start a new goroutine for
	// non main work.
	lvm, err := NewLuaVmWithPrelude(cfg)

	panicOn(err)
	inc := NewIncrState(lvm, cfg)

	r := &Repl{cfg: cfg, lvm: lvm, inc: inc}
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
	case ":stacks":
		showLuaStacks(r.lvm.vm)
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
		err = LuaDoPreludeFiles(r.lvm, files)
		if err != nil {
			fmt.Printf("error during prelude reload: '%v'", err)
		}
		return "", nil
	case ":help", ":?":
		fmt.Printf(`
======================
gijit: a go interpreter, just-in-time
https://github.com/glycerine/gi
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
 :stacks         Show lua stacks for each coroutine.
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
				err = LuaDoUserFiles(r.lvm, final)
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

	useEval := !r.cfg.RawLua
	err := LuaRun(r.lvm, use, useEval)
	if err != nil {
		fmt.Printf("error from LuaRun: supplied lua with: '%s'\nlua stack:\n%v\n", use[:len(use)-1], err)
		return nil
	}
	r.t1 = time.Now()
	fmt.Printf("\n")
	r.reader.Reset(os.Stdin)
	fmt.Printf("elapsed: '%v'\n", r.t1.Sub(r.t0))

	return nil
}

// :ls, :gls, :lst, :glst implementation
func (r *Repl) displayCmd(cmd string) {
	err := LuaRun(r.lvm, `__`+cmd+`()`, true)
	panicOn(err)
}

func showLuaStacks(vm *golua.State) {
	top := vm.GetTop()
	vm.GetGlobal("__all_coro")
	if vm.IsNil(-1) {
		panic("could not locate __all_coro in _G")
	}
	fmt.Printf("\n")
	forEachAllCoroArrayValue(vm, -1, func(i int, name, status string) {
		ignoreTopmost := 1
		if i == 1 {
			if name != "main" {
				panic(fmt.Sprintf("name should have been main for i == 1 thread! not '%s'", name))
			}
			// main thread is where we are iterating from,
			// so it has value, key, __all_coro, __all_coro.
			// Ignore that administrivia, its not interesting.
			ignoreTopmost = 5
		}
		thr := vm.ToThread(-1)
		fmt.Printf("===================================\n")
		fmt.Printf("        __all_coro %v: '%s' (%s)\n", i, name, status)
		fmt.Printf("===================================\n"+
			"%s\n", DumpLuaStackAsString(thr, ignoreTopmost))
	})
	vm.SetTop(top)
}

// Call f with each __all_coro array value in term on the top of
// the stack, the f(i, name) call will have i set to 1, 2, 3, ...
// in turn.
func forEachAllCoroArrayValue(L *golua.State, index int, f func(i int, name, status string)) {

	i := 1
	// Push another reference to the table on top of the stack (so we know
	// where it is, and this function can work for negative, positive and
	// pseudo indices.
	L.PushValue(index)
	// stack now contains: -1 => table
	L.PushNil()
	// stack now contains: -1 => nil; -2 => table
	for L.Next(-2) != 0 {

		// stack now contains: -1 => value; -2 => key; -3 => table

		L.PushValue(-1)
		// stack: value, value, key, table

		// get name from __coro2notes
		L.GetGlobal("__coro2notes")
		// stack: __coro2notes, value, value, key, table
		L.Insert(-2)
		// stack: value, __coro2notes, value, key, table
		L.GetTable(-2)
		// stack: details, __coro2notes, value, key, table
		L.PushString("__name")
		// stack: "__name", details, __coro2notes, value, key, table
		L.GetTable(-2)
		// stack: threadName, details, __coro2notes, value, key, table
		if L.IsNil(-1) {
			panic("thread did not have a name?!?")
		}
		name := L.ToString(-1)
		L.Pop(3)
		// stack: value, key, table
		status := getCoroutineStatus(L, -1)

		// stack: value, key, table
		f(i, name, status)

		// pop value, leaving original key
		L.Pop(1)
		// stack: key, table
		i++
	}
	// stack now contains: -1 => table (when lua_next returns 0 it pops the key
	// but does not push anything.)
	// Pop table
	L.Pop(1)
	// Stack is now the same as it was on entry to this function
	return

}

func getCoroutineStatus(L *golua.State, index int) (status string) {
	top := L.GetTop()
	L.PushValue(index)
	// stack: thread
	L.GetGlobal("coroutine")
	// stack: coroutine table, thread
	L.PushString("status")
	// stack: "status", coroutine table, thread
	L.GetTable(-2)
	// stack: status function, coroutine table, thread
	L.Remove(-2)
	// stack: status function, thread
	L.Insert(-2)
	// stack: thread, status function
	L.Call(1, 1)
	// stack: status string
	status = L.ToString(-1)
	L.SetTop(top)
	return
}

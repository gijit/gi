package compiler

import (
	"fmt"
	"sync"
	"time"

	golua "github.com/glycerine/golua/lua"
	"github.com/glycerine/idem"
	"github.com/glycerine/luar"
)

var _ = time.Now

// act as a thread-safe proxy to a lua state
// vm running on its own goroutine.
type Goro struct {
	cfg      *GoroConfig
	lvm      *LuaVm
	vm       *golua.State
	halt     *idem.Halter
	beat     time.Duration
	doticket chan *ticket

	mut     sync.Mutex
	started bool
}

type GoroConfig struct {
	GiCfg *GIConfig
	off2  bool // don't run in separate goroutine.
}

type ticket struct {
	myGoro *Goro

	//input
	// register these vars, optional
	regns  string
	regmap luar.Map

	//input
	// what to do after any registrations, optional
	run []byte

	//input
	// what to fetch for return, optional;
	// one of register, run, or varname should be
	// set. Else no point in making the Do ticket call.
	varname          map[string]interface{}
	gettyp           GetType
	leaveOnTop       bool
	useEvalCoroutine bool

	//output
	runErr error
	getErr error
	// the values in the varname map are output too.

	// closed when txn complete
	done chan struct{}
}

func (r *Goro) newTicket(run string, useEvalCoroutine bool) *ticket {
	//fmt.Printf("goro.newTicket: top \n")

	if r == nil {
		panic("newTicket cannot be called on nil Goro r")
	}
	t := &ticket{
		myGoro:           r,
		regmap:           make(luar.Map),
		varname:          make(map[string]interface{}),
		done:             make(chan struct{}),
		run:              []byte(run),
		useEvalCoroutine: useEvalCoroutine,
	}
	return t
}

// There's only one incremental compiler state; it
// typechecks all new code, and generates any new Lua.
//
// So any issues around typechecking and variable
// types in the front end are all on the main
// repl goroutine.
//
// The multiple backend lua states, each on their
// own goroutine, only know about (translated) lua code.

type GetType int

const (
	GetInt64 GetType = iota
	GetString
	GetChan
)

func NewGoro(lvm *LuaVm, cfg *GoroConfig) (*Goro, error) {
	//fmt.Printf("\n\n NewGoro: top lvm='%#v'; lvm.vm='%#v'\n", lvm, lvm.vm)
	if cfg == nil {
		cfg = &GoroConfig{}
	}

	//var err error
	if lvm == nil {
		panic("lvm cannot be nil") // try to root out duplicate starts.
		/*lvm, err = NewLuaVmWithPrelude(cfg.GiCfg)
		if err != nil {
			return nil, err
		}
		*/
	}

	r := &Goro{
		cfg:      cfg,
		lvm:      lvm,
		vm:       lvm.vm,
		halt:     idem.NewHalter(),
		doticket: make(chan *ticket),
		beat:     2000 * time.Millisecond,
	}

	if !cfg.off2 {
		r.Start()
	}
	return r, nil
}

func (r *Goro) Start() {
	r.mut.Lock()
	if r.started {
		panic("cannot start goro more than once!")
	}
	r.started = true
	fmt.Printf("\n 0000011111 Goro.Start() succeeded on %p\n\n", r)
	r.mut.Unlock()

	go func() {
		defer func() {
			fmt.Printf("\n 0000022222 goro Start() finished. closing r= %p \n", r)
			r.lvm.vm.Close()
			r.halt.MarkDone()
		}()

		fmt.Printf("r.beat is %v\n", r.beat)
		heartbeat := time.After(r.beat)
		for {
			select {
			case <-heartbeat:
				r.handleHeartbeat()
				heartbeat = time.After(r.beat)
			case <-r.halt.ReqStop.Chan:
				return
			case t := <-r.doticket:
				r.handleTicket(t)
			}
		}
	}()
}

var resumeSchedBytes = []byte("__task.resume_scheduler();")

func (r *Goro) handleHeartbeat() {
	fmt.Printf("goro heartbeat! on r=%p, at %v\n", r, time.Now())
	err := r.privateRun(resumeSchedBytes, true)
	panicOn(err)
}

func (r *Goro) handleTicket(t *ticket) {
	//fmt.Printf("goro.handleTicket: top \n")

	if len(t.regmap) > 0 {
		luar.Register(r.vm, t.regns, t.regmap)
		//fmt.Printf("jea debug, back from luar.Register with regns: '%s', map: '%#v'\n", t.regns, t.regmap)
	}

	if len(t.run) > 0 {
		t.runErr = r.privateRun(t.run, t.useEvalCoroutine)
	}
	if t.runErr == nil && len(t.varname) > 0 {
		for key := range t.varname {
			if key == "" {
				continue
			}
			r.vm.GetGlobal(key)
			if r.vm.IsNil(-1) {
				r.vm.Pop(1)
				t.getErr = fmt.Errorf("not found: '%s'", t.varname)
				break
			} else {
				switch t.gettyp {
				case GetInt64:
					t.varname[key] = r.vm.CdataToInt64(-1)
				case GetString:
					t.varname[key] = r.vm.ToString(-1)
				case GetChan:
					r.vm.Pop(1)
					t.varname[key], t.getErr = getChannelFromGlobal(r.lvm, key, true)
				}
				if !t.leaveOnTop {
					r.vm.Pop(1)
				}
			}
		}
	}
	close(t.done)
}

func (r *Goro) do(t *ticket) {
	if r.cfg.off2 {
		r.handleTicket(t)
	} else {
		r.doticket <- t
		<-t.done
	}
}

func (t *ticket) Do() error {
	t.myGoro.do(t)
	if t.runErr != nil {
		return t.runErr
	}
	if t.getErr != nil {
		return t.getErr
	}
	return nil
}

// privateRun should only be called by Goro, to provide
// appropriate synchronization (we must only have one thread at a time
// should be calling the LuaJIT vm). This is
// where code actually gets run on the vm.
func (goro *Goro) privateRun(run []byte, useEvalCoroutine bool) error {

	lvm := goro.lvm

	s := string(run)
	//pp("LuaRun top. s='%s'. stack='%s'", s, string(debug.Stack()))
	vm := lvm.vm
	startTop := vm.GetTop()
	defer vm.SetTop(startTop)

	if useEvalCoroutine {
		// get the eval function. it will spawn us a new coroutine
		// for each evaluation.

		vm.GetGlobal("__eval")
		if vm.IsNil(-1) {
			panic("could not locate __eval in _G")
		}
		eval := vm.ToPointer(-1)
		_ = eval
		vm.PushString(s)

		/*p1("good: found __eval (0x%x). it is at -2 of the stack, our running code at -1. running '%s'\n", eval, s)
		p1("before vm.Call(1,0), stacks are:")
		if verb.Verbose {
			showLuaStacks(vm)
		}
		*/
		vm.Call(1, 0)
		// if things crash, this is the first place
		// to check for an error: dump the Lua stack.
		// With high probability, it will yield clues to the problem.
		/*
			p1("after vm.Call(1,0), stacks are:")
			if verb.Verbose {
				showLuaStacks(vm)
			}
		*/
		return nil
	} else {

		// not using the __eval coroutine.

		interr := vm.LoadString(s)
		if interr != 0 {
			loadErr := fmt.Errorf("%s", DumpLuaStackAsString(vm, 0))
			return loadErr
		} else {
			err := vm.Call(0, 0)
			if err != nil {
				runErr := fmt.Errorf("%s", DumpLuaStackAsString(vm, 0))
				return runErr
			}
		}
	}
	return nil
}

package compiler

import (
	"fmt"

	golua "github.com/glycerine/golua/lua"
	"github.com/glycerine/idem"
	//"github.com/glycerine/luar"
)

type ticket struct {
	// what to do first, optional
	runme  []byte
	runErr error

	// what to fetch for return, optional
	varname    string
	gettyp     GetType
	leaveOnTop bool
	getErr     error
	getResult  interface{}

	// closed when txn complete
	done chan struct{}
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

// act as a thread-safe proxy to a lua state
// vm running on its own goroutine.
type Goro struct {
	cfg  *GoroConfig
	vm   *golua.State
	halt *idem.Halter
	Do   chan *ticket
}

type GoroConfig struct {
	GiCfg *GIConfig
}

func NewGoro(cfg *GoroConfig) (*Goro, error) {
	if cfg == nil {
		cfg = &GoroConfig{}
	}
	vm, err := NewLuaVmWithPrelude(cfg.GiCfg)
	if err != nil {
		return nil, err
	}
	r := &Goro{
		cfg:  cfg,
		vm:   vm,
		halt: idem.NewHalter(),
		Do:   make(chan *ticket),
	}
	r.Start()
	return r, nil
}

func (r *Goro) Start() {
	go func() {
		defer func() {
			r.halt.MarkDone()
		}()
		for {
			select {
			case <-r.halt.ReqStop.Chan:
				return
			case t := <-r.Do:
				startTop := r.vm.GetTop()
				if len(t.runme) > 0 {
					s := string(t.runme)

					interr := r.vm.LoadString(s)
					if interr != 0 {
						fmt.Printf("error from Lua vm.LoadString(): supplied lua with: '%s'\nlua stack:\n", s)
						DumpLuaStack(r.vm)
						t.runErr = fmt.Errorf(r.vm.ToString(-1))
						r.vm.SetTop(startTop)
					} else {
						err := r.vm.Call(0, 0)
						if err != nil {
							fmt.Printf("error from Lua vm.Call(0,0): '%v'. supplied lua with: '%s'\nlua stack:\n", err, s)
							DumpLuaStack(r.vm)
							r.vm.Pop(1)
							t.runErr = err
						}
					}

				}
				if t.varname != "" {
					r.vm.GetGlobal(t.varname)
					if r.vm.IsNil(-1) {
						r.vm.Pop(1)
						t.getErr = fmt.Errorf("not found: '%s'", t.varname)
					} else {
						switch t.gettyp {
						case GetInt64:
							t.getResult = r.vm.CdataToInt64(-1)
						case GetString:
							t.getResult = r.vm.ToString(-1)
						case GetChan:
							t.getResult, t.getErr = getChannelFromGlobal(r.vm, t.varname, true)
						}
						if !t.leaveOnTop {
							r.vm.Pop(1)
						}
					}
				}
				close(t.done)
			}
		}
	}()
}

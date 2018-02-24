package compiler

import (
	"fmt"

	golua "github.com/glycerine/golua/lua"
	"github.com/glycerine/idem"
	"github.com/glycerine/luar"
)

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
	varname    map[string]interface{}
	gettyp     GetType
	leaveOnTop bool

	//output
	runErr error
	getErr error
	// the values in the varname map are output too.

	// closed when txn complete
	done chan struct{}
}

func (r *Goro) newTicket() *ticket {
	if r == nil {
		panic("newTicket cannot be called on nil Goro r")
	}
	t := &ticket{
		myGoro:  r,
		regmap:  make(luar.Map),
		varname: make(map[string]interface{}),
		done:    make(chan struct{}),
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

// act as a thread-safe proxy to a lua state
// vm running on its own goroutine.
type Goro struct {
	cfg      *GoroConfig
	vm       *golua.State
	halt     *idem.Halter
	doticket chan *ticket
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
		cfg:      cfg,
		vm:       vm,
		halt:     idem.NewHalter(),
		doticket: make(chan *ticket),
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
			case t := <-r.doticket:
				startTop := r.vm.GetTop()

				if len(t.regmap) > 0 {
					luar.Register(r.vm, t.regns, t.regmap)
				}

				if len(t.run) > 0 {
					s := string(t.run)

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
				if len(t.varname) > 0 {
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
								t.varname[key], t.getErr = getChannelFromGlobal(r.vm, key, true)
							}
							if !t.leaveOnTop {
								r.vm.Pop(1)
							}
						}
					}
				}
				close(t.done)
			}
		}
	}()
}

func (r *Goro) do(t *ticket) {
	r.doticket <- t
	<-t.done
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

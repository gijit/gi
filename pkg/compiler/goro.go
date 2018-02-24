package compiler

import (
	golua "github.com/glycerine/golua/lua"
	"github.com/glycerine/idem"
	//"github.com/glycerine/luar"
)

// act as a thread-safe proxy to a lua state
// vm running on its own goroutine.
type Goro struct {
	cfg  *GoroConfig
	vm   *golua.State
	halt *idem.Halter
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

			}
		}
	}()
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

func (r *Goro) LuaRunAndReport(code string)                     {}
func (r *Goro) GetInt64(varname string) int64                   { return 0 }
func (r *Goro) GetString(varname string) string                 { return "" }
func (r *Goro) GetVar(varname string) (interface{}, int, error) { return nil, 0, nil }
func (r *Goro) GetChannel(varname string, leaveOnTop bool) (interface{}, error) {
	return nil, nil
}

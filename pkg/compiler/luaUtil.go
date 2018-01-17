package compiler

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	golua "github.com/glycerine/golua/lua"
	"github.com/glycerine/luar"
)

type VmConfig struct {
	PreludePath string
	Quiet       bool
}

func NewVmConfig() *VmConfig {
	return &VmConfig{}
}

func NewLuaVmWithPrelude(cfg *VmConfig) (*golua.State, error) {
	vm := luar.Init() // does vm.OpenLibs() for us, adds luar. functions.

	if cfg == nil {
		cfg = NewVmConfig()
		cfg.PreludePath = "."
	}

	// load prelude
	files, err := FetchPreludeFilenames(cfg.PreludePath, cfg.Quiet)
	if err != nil {
		return nil, err
	}
	err = LuaDoFiles(vm, files)
	return vm, err
}

func LuaDoFiles(vm *golua.State, files []string) error {
	for _, f := range files {
		pp("LuaDoFiles, f = '%s'", f)
		if f == "lua.help.lua" {
			panic("where lua.help.lua?")
		}
		interr := vm.LoadString(fmt.Sprintf(`dofile("%s")`, f))
		if interr != 0 {
			pp("interr %v on vm.LoadString for dofile on '%s'", interr, f)
			msg := DumpLuaStackAsString(vm)
			vm.Pop(1)
			return fmt.Errorf("error in setupPrelude during LoadString on file '%s': Details: '%s'", f, msg)
		}
		err := vm.Call(0, 0)
		if err != nil {
			msg := DumpLuaStackAsString(vm)
			vm.Pop(1)
			return fmt.Errorf("error in setupPrelude during Call on file '%s': '%v'. Details: '%s'", f, err, msg)
		}
	}
	return nil
}

func DumpLuaStack(L *golua.State) {
	var top int

	top = L.GetTop()
	fmt.Printf("========== begin DumpLuaStack: top = %v\n", top)
	for i := top; i >= 1; i-- {
		t := L.Type(i)
		fmt.Printf("DumpLuaStack: i=%v, t= %v\n", i, t)
		switch t {
		case golua.LUA_TSTRING:
			fmt.Printf("String : \t%v\n", L.ToString(i))
		case golua.LUA_TBOOLEAN:
			fmt.Printf("Bool : \t\t%v\n", L.ToBoolean(i))
		case golua.LUA_TNUMBER:
			fmt.Printf("Number : \t%v\n", L.ToNumber(i))
		default:
			fmt.Printf("Type(number %v) : has type name \t%v\n", t, L.Typename(i))
		}
	}
	fmt.Printf("========= end of DumpLuaStack\n")
}

func DumpLuaStackAsString(L *golua.State) string {
	var top int
	s := ""
	top = L.GetTop()
	pp("DumpLuaStackAsString: top = %v", top)
	for i := top; i >= 1; i-- {
		pp("i=%v out of top = %v", i, top)
		t := L.Type(i)
		switch t {
		case golua.LUA_TSTRING:
			s += fmt.Sprintf("String : \t%v", L.ToString(i))
		case golua.LUA_TBOOLEAN:
			s += fmt.Sprintf("Bool : \t\t%v", L.ToBoolean(i))
		case golua.LUA_TNUMBER:
			s += fmt.Sprintf("Number : \t%v", L.ToNumber(i))
		default:
			s += fmt.Sprintf("Type : \t\t%v", L.Typename(i))
		}
	}
	return s
}

func FetchPreludeFilenames(preludePath string, quiet bool) ([]string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	pp("FetchPrelude called on path '%s', where cwd = '%s'", preludePath, cwd)
	if !DirExists(preludePath) {
		return nil, fmt.Errorf("-prelude dir does not exist: '%s'", preludePath)
	}
	files, err := filepath.Glob(fmt.Sprintf("%s/*.lua", preludePath))
	if err != nil {
		return nil, fmt.Errorf("-prelude dir '%s' open problem: '%v'", preludePath, err)
	}
	if len(files) < 1 {
		return nil, fmt.Errorf("-prelude dir '%s' had no lua files in it.", preludePath)
	}
	if !quiet {
		fmt.Printf("using this prelude directory: '%s'\n", preludePath)
		shortFn := make([]string, len(files))
		for i, fn := range files {
			shortFn[i] = path.Base(fn)
		}
		fmt.Printf("using these files as prelude: %s\n", strings.Join(shortFn, ", "))
	}
	return files, nil
}

func LuaMustInt(vm *golua.State, varname string, expect int) {

	vm.GetGlobal(varname)
	top := vm.GetTop()
	value_int := vm.ToInteger(top)

	pp("value_int=%v", value_int)
	if value_int != expect {
		panic(fmt.Sprintf("expected %v, got %v for '%v'", expect, value_int, varname))
	}
}

func LuaMustString(vm *golua.State, varname string, expect string) {

	vm.GetGlobal(varname)
	top := vm.GetTop()
	value_string := vm.ToString(top)

	pp("value_string=%v", value_string)
	if value_string != expect {
		panic(fmt.Sprintf("expected %v, got value '%s' -> '%v'", expect, varname, value_string))
	}
}

func LuaMustBool(vm *golua.State, varname string, expect bool) {

	vm.GetGlobal(varname)
	top := vm.GetTop()
	value_bool := vm.ToBoolean(top)

	pp("value_bool=%v", value_bool)
	if value_bool != expect {
		panic(fmt.Sprintf("expected %v, got value '%s' -> '%v'", expect, varname, value_bool))
	}
}

func LuaRunAndReport(vm *golua.State, s string) {
	interr := vm.LoadString(s)
	if interr != 0 {
		fmt.Printf("error from Lua vm.LoadString(): supplied lua with: '%s'\nlua stack:\n", s)
		DumpLuaStack(vm)
		vm.Pop(1)
	} else {
		err := vm.Call(0, 0)
		if err != nil {
			fmt.Printf("error from Lua vm.Call(0,0): '%v'. supplied lua with: '%s'\nlua stack:\n", err, s)
			DumpLuaStack(vm)
			vm.Pop(1)
		}
	}
}

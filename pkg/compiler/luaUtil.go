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
	fmt.Printf("\n%s\n", DumpLuaStackAsString(L))
}

func DumpLuaStackAsString(L *golua.State) (s string) {
	var top int

	top = L.GetTop()
	s += fmt.Sprintf("========== begin DumpLuaStack: top = %v\n", top)
	for i := top; i >= 1; i-- {
		t := L.Type(i)
		s += fmt.Sprintf("DumpLuaStack: i=%v, t= %v\n", i, t)
		switch t {
		case golua.LUA_TSTRING:
			s += fmt.Sprintf(" String : \t%v\n", L.ToString(i))
		case golua.LUA_TBOOLEAN:
			s += fmt.Sprintf(" Bool : \t\t%v\n", L.ToBoolean(i))
		case golua.LUA_TNUMBER:
			s += fmt.Sprintf(" Number : \t%v\n", L.ToNumber(i))
		case golua.LUA_TTABLE:
			s += fmt.Sprintf(" Table : \n%s\n", dumpTableString(L, i))

		case 10: // LUA_TCDATA aka cdata
			//pp("Dump cdata case, L.Type(idx) = '%v'", L.Type(i))
			ctype := L.LuaJITctypeID(i)
			//pp("luar.go Dump sees ctype = %v", ctype)
			switch ctype {
			case 5: //  int8
			case 6: //  uint8
			case 7: //  int16
			case 8: //  uint16
			case 9: //  int32
			case 10: //  uint32
			case 11: //  int64
				val := L.CdataToInt64(i)
				s += fmt.Sprintf(" int64: '%v'\n", val)
			case 12: //  uint64
				val := L.CdataToUint64(i)
				s += fmt.Sprintf(" uint64: '%v'\n", val)
			case 13: //  float32
			case 14: //  float64

			case 0: // means it wasn't a ctype
			}

		default:
			s += fmt.Sprintf(" Type(code %v) : no auto-print available.\n", t)
		}
	}
	s += fmt.Sprintf("========= end of DumpLuaStack\n")
	return
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

// prefer below LuaMustInt64
func LuaMustInt(vm *golua.State, varname string, expect int) {

	vm.GetGlobal(varname)
	top := vm.GetTop()
	value_int := vm.ToInteger(top) // lossy for 64-bit int64, use vm.CdataToInt64() instead.

	pp("LuaMustInt, expect='%v'; observe value_int='%v'", expect, value_int)
	if value_int != expect {
		DumpLuaStack(vm)
		panic(fmt.Sprintf("expected %v, got %v for '%v'", expect, value_int, varname))
	}
}

func LuaMustInt64(vm *golua.State, varname string, expect int64) {

	vm.GetGlobal(varname)
	top := vm.GetTop()
	value_int := vm.CdataToInt64(top)

	pp("LuaMustInt64, expect='%v'; observe value_int='%v'", expect, value_int)
	if value_int != expect {
		DumpLuaStack(vm)
		panic(fmt.Sprintf("expected %v, got %v for '%v'", expect, value_int, varname))
	}
}

func LuaMustString(vm *golua.State, varname string, expect string) {

	vm.GetGlobal(varname)
	top := vm.GetTop()
	value_string := vm.ToString(top)

	pp("value_string=%v", value_string)
	if value_string != expect {
		DumpLuaStack(vm)
		panic(fmt.Sprintf("expected %v, got value '%s' -> '%v'", expect, varname, value_string))
	}
}

func LuaMustBool(vm *golua.State, varname string, expect bool) {

	vm.GetGlobal(varname)
	top := vm.GetTop()
	value_bool := vm.ToBoolean(top)

	pp("value_bool=%v", value_bool)
	if value_bool != expect {
		DumpLuaStack(vm)
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

func dumpTableString(L *golua.State, index int) (s string) {

	// Push another reference to the table on top of the stack (so we know
	// where it is, and this function can work for negative, positive and
	// pseudo indices
	L.PushValue(index)
	// stack now contains: -1 => table
	L.PushNil()
	// stack now contains: -1 => nil; -2 => table
	for L.Next(-2) != 0 {

		// stack now contains: -1 => value; -2 => key; -3 => table
		// copy the key so that lua_tostring does not modify the original
		L.PushValue(-2)
		// stack now contains: -1 => key; -2 => value; -3 => key; -4 => table
		key := L.ToString(-1)
		value := L.ToString(-2)
		s += fmt.Sprintf("'%s' => '%s'\n", key, value)
		// pop value + copy of key, leaving original key
		L.Pop(2)
		// stack now contains: -1 => key; -2 => table
	}
	// stack now contains: -1 => table (when lua_next returns 0 it pops the key
	// but does not push anything.)
	// Pop table
	L.Pop(1)
	// Stack is now the same as it was on entry to this function
	return
}

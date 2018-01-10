package compiler

import (
	"fmt"
	luajit "github.com/glycerine/golua/lua"
	"path/filepath"
)

func LuaDoFiles(vm *luajit.State, files []string) error {
	for _, f := range files {
		pp("LuaDoFiles, f = '%s'", f)
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

func DumpLuaStack(L *luajit.State) {
	var top int

	top = L.GetTop()
	pp("DumpLuaStack: top = %v", top)
	for i := 1; i <= top; i++ {
		t := L.Type(i)
		switch t {
		case luajit.LUA_TSTRING:
			fmt.Println("String : \t", L.ToString(i))
		case luajit.LUA_TBOOLEAN:
			fmt.Println("Bool : \t\t", L.ToBoolean(i))
		case luajit.LUA_TNUMBER:
			fmt.Println("Number : \t", L.ToNumber(i))
		default:
			fmt.Println("Type : \t\t", L.Typename(i))
		}
	}
	print("\n")
}

func DumpLuaStackAsString(L *luajit.State) string {
	var top int
	s := ""
	top = L.GetTop()
	pp("DumpLuaStackAsString: top = %v", top)
	for i := 1; i <= top; i++ {
		pp("i=%v out of top = %v", i, top)
		t := L.Type(i)
		switch t {
		case luajit.LUA_TSTRING:
			s += fmt.Sprintf("String : \t%v", L.ToString(i))
		case luajit.LUA_TBOOLEAN:
			s += fmt.Sprintf("Bool : \t\t%v", L.ToBoolean(i))
		case luajit.LUA_TNUMBER:
			s += fmt.Sprintf("Number : \t%v", L.ToNumber(i))
		default:
			s += fmt.Sprintf("Type : \t\t%v", L.Typename(i))
		}
	}
	return s
}

func FetchPrelude(path string) ([]string, error) {
	if !DirExists(path) {
		return nil, fmt.Errorf("-prelude dir does not exist: '%s'", path)
	}
	files, err := filepath.Glob(fmt.Sprintf("%s/*.lua", path))
	if err != nil {
		return nil, fmt.Errorf("-prelude dir '%s' open problem: '%v'", path, err)
	}
	if len(files) < 1 {
		return nil, fmt.Errorf("-prelude dir '%s' had no lua files in it.", path)
	}
	return files, nil
}

func LuaMustInt(vm *luajit.State, varname string, expect int) {

	vm.GetGlobal(varname)
	top := vm.GetTop()
	value_int := vm.ToInteger(top)

	pp("value_int=%v", value_int)
	if value_int != expect {
		panic(fmt.Sprintf("expected %v, got %v for '%v'", expect, value_int, varname))
	}
}

func LuaMustString(vm *luajit.State, varname string, expect string) {

	vm.GetGlobal(varname)
	top := vm.GetTop()
	value_string := vm.ToString(top)

	pp("value_string=%v", value_string)
	if value_string != expect {
		panic(fmt.Sprintf("expected %v, got value '%s' -> '%v'", expect, varname, value_string))
	}
}

func LuaMustBool(vm *luajit.State, varname string, expect bool) {

	vm.GetGlobal(varname)
	top := vm.GetTop()
	value_bool := vm.ToBoolean(top)

	pp("value_bool=%v", value_bool)
	if value_bool != expect {
		panic(fmt.Sprintf("expected %v, got value '%s' -> '%v'", expect, varname, value_bool))
	}
}

func LuaRunAndReport(vm *luajit.State, s string) {
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

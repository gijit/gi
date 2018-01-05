package compiler

import (
	"fmt"
	"github.com/glycerine/luajit"
	"path/filepath"
)

func SetupPrelude(vm *luajit.State, files []string) error {
	for _, f := range files {
		err := vm.Loadstring(fmt.Sprintf(`dofile("%s")`, f))
		if err != nil {
			msg := DumpLuaStackAsString(vm)
			vm.Pop(1)
			return fmt.Errorf("error in setupPrelude during LoadString on file '%s': '%v'. Details: '%s'", f, err, msg)
		}
		err = vm.Pcall(0, 0, 0)
		if err != nil {
			msg := DumpLuaStackAsString(vm)
			vm.Pop(1)
			return fmt.Errorf("error in setupPrelude during Pcall on file '%s': '%v'. Details: '%s'", f, err, msg)
		}
	}
	return nil
}

func DumpLuaStack(L *luajit.State) {
	var top int

	top = L.Gettop()
	for i := 1; i <= top; i++ {
		t := L.Type(i)
		switch t {
		case luajit.Tstring:
			fmt.Println("String : \t", L.Tostring(i))
		case luajit.Tboolean:
			fmt.Println("Bool : \t\t", L.Toboolean(i))
		case luajit.Tnumber:
			fmt.Println("Number : \t", L.Tonumber(i))
		default:
			fmt.Println("Type : \t\t", L.Typename(i))
		}
	}
	print("\n")
}

func DumpLuaStackAsString(L *luajit.State) string {
	var top int
	s := ""
	top = L.Gettop()
	for i := 1; i <= top; i++ {
		t := L.Type(i)
		switch t {
		case luajit.Tstring:
			s += fmt.Sprintf("String : \t%v", L.Tostring(i))
		case luajit.Tboolean:
			s += fmt.Sprintf("Bool : \t\t%v", L.Toboolean(i))
		case luajit.Tnumber:
			s += fmt.Sprintf("Number : \t%v", L.Tonumber(i))
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

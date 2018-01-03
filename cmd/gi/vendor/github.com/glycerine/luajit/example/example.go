package main

import (
	"fmt"
	"github.com/glycerine/luajit"
)

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

func PlayWithStack() {
	L := luajit.Newstate()
	defer L.Close()

	L.Pushboolean(true)
	L.Pushnumber(10)
	L.Pushnil()
	L.Pushstring("hello")

	DumpLuaStack(L)

	L.Pushvalue(-4)
	DumpLuaStack(L)

	L.Replace(3)
	DumpLuaStack(L)

	L.Settop(6)
	DumpLuaStack(L)

	L.Remove(-3)
	DumpLuaStack(L)

	L.Settop(-5)
	DumpLuaStack(L)
}

func PlayWithGlobal() {
	L := luajit.Newstate()
	defer L.Close()

	L.Openlibs()

	L.Loadfile("example.lua")
	L.Pcall(0, 0, 0)

	L.Getglobal("width")
	L.Getglobal("height")

	if L.Isnumber(-2) != true {
		print("width should be number\n")
	}
	if L.Isnumber(-1) != true {
		print("height should be number\n")
	}

	DumpLuaStack(L)
}

func GetField(L *luajit.State, f string) (float64, bool) {
	var result float64

	L.Pushstring(f)
	L.Gettable(-2)
	if L.Isnumber(-1) != true {
		print("invalid field in table")
		return 0, false
	}

	result = L.Tonumber(-1)
	L.Pop(1)

	return result, true

}

func SetField(L *luajit.State, f string, v float64) {
	L.Pushstring(f)
	L.Pushnumber(v)
	L.Settable(-3)
}

func PlayWithTables() {
	L := luajit.Newstate()
	defer L.Close()

	L.Openlibs()

	L.Loadfile("example.lua")
	L.Pcall(0, 0, 0)

	L.Getglobal("background")
	if L.Istable(-1) != true {
		print("'background' is not defined in lua script")
	}

	DumpLuaStack(L)

	red, _ := GetField(L, "r")
	green, _ := GetField(L, "g")
	blue, _ := GetField(L, "b")
	DumpLuaStack(L)

	print("background: ", red, " ", green, " ", blue, "\n")

	L.Newtable()
	SetField(L, "red", 50.0)
	SetField(L, "green", 30.0)
	SetField(L, "blue", 20.0)
	L.Setglobal("foreground")

	L.Getglobal("foreground")
	if L.Istable(-1) != true {
		print("'foreground' is not defined in lua script")
	}

	DumpLuaStack(L)

	red, _ = GetField(L, "red")
	green, _ = GetField(L, "green")
	blue, _ = GetField(L, "blue")
	DumpLuaStack(L)

	print("foreground: ", red, " ", green, " ", blue, "\n")

	L.Getglobal("f")
	L.Pushnumber(1)
	L.Pushnumber(2)
	L.Pcall(2, 1, 0)

	if L.Isnumber(-1) != true {
		print("result is not number")
	}
	r := L.Tonumber(-1)
	L.Pop(1)

	print("result of f(x,y): ", r, "\n")

	L.Getglobal("print_foreground")
	L.Pcall(0, 0, 0)

	L.Getglobal("print_background")
	L.Pcall(0, 0, 0)
}

func summator(L *luajit.State) int {
	a := L.Tonumber(1)
	b := L.Tonumber(2)
	L.Pushnumber(a + b)
	return 1
}

func PlayWithGoFunction() {
	L := luajit.Newstate()
	defer L.Close()

	L.Openlibs()

	L.Loadfile("example.lua")
	L.Pcall(0, 0, 0)

	L.Pushfunction(summator)
	L.Setglobal("summator")

	L.Getglobal("print_summator")
	L.Pcall(0, 0, 0)

}

func main() {
	//	PlayWithStack()

	//	PlayWithGlobal()

	//  PlayWithTables()

	PlayWithGoFunction()
}

package compiler

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

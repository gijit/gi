package compiler

import (
	"fmt"

	//"github.com/glycerine/golua/lua"
	"github.com/glycerine/luar"
)

// call fmt.Sprintf from Lua

func CallSprintf() {
	const test = `
for i = 1, 3 do
		print(msg, i)
end
print(user)
print(user.Name, user.Age)
`

	type person struct {
		Name string
		Age  int
	}

	L := luar.Init()
	defer L.Close()

	user := &person{"Dolly", 46}

	luar.Register(L, "", luar.Map{
		// Go functions may be registered directly.
		"print": fmt.Println,
		// Constants can be registered.
		"msg": "foo",
		// And other values as well.
		"user": user,
	})

	L.DoString(test)
}

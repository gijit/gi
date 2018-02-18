package luar_test

// TODO: Lua's print() function will output to stdout instead of being compared
// to the desired result. Workaround: register one of Go's printing functions
// from the fmt library.

import (
	"fmt"
	"log"
	"sort"
	//"strconv"
	"sync"

	"github.com/glycerine/golua/lua"
	"github.com/glycerine/luar"
)

func Example() {
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
	// Output:
	// foo 1
	// foo 2
	// foo 3
	// &{Dolly 46}
	// Dolly 46
}

// Pointers to structs and structs within pointers are automatically dereferenced.
func Example_pointers() {
	const test = `
local t = newRef()
print(t.Index, t.Number, t.Title)
`

	type Ref struct {
		Index  int
		Number *int
		Title  *string
	}

	newRef := func() *Ref {
		n := new(int)
		*n = 10
		t := new(string)
		*t = "foo"
		return &Ref{Index: 17, Number: n, Title: t}
	}

	L := luar.Init()
	defer L.Close()

	luar.Register(L, "", luar.Map{
		"print":  fmt.Println,
		"newRef": newRef,
	})

	L.DoString(test)
	// Output:
	// 17 10 foo
}

// Slices must be looped over with 'ipairs'.
func Example_slices() {
	const test = `
for i, v in ipairs(names) do
	 print(i, v)
end
`

	L := luar.Init()
	defer L.Close()

	names := []string{"alfred", "alice", "bob", "frodo"}

	luar.Register(L, "", luar.Map{
		"print": fmt.Println,
		"names": names,
	})

	L.DoString(test)
	// Output:
	// 1 alfred
	// 2 alice
	// 3 bob
	// 4 frodo
}

// Init() is not needed when no proxy is involved.
func ExampleGoToLua() {
	L := lua.NewState()
	defer L.Close()
	L.OpenLibs()

	input := "Hello world!"
	luar.GoToLua(L, input)
	L.SetGlobal("input")

	luar.GoToLua(L, fmt.Println)
	L.SetGlobal("print")
	L.DoString("print(input)")
	// Output:
	// Hello world!
}

// jea: ExampleInit output based on maps can
//      change the order of output lines, so
//  it isn't deterministic. We comment it out
//  to get repeatable tests.
/*
// This example shows how Go slices and maps are marshalled to Lua tables and
// vice versa.
//
// An arbitrary Go function is callable from Lua, and list-like tables become
// slices on the Go side. The Go function returns a map, which is wrapped as a
// proxy object. You can copy this proxy to a Lua table explicitly. There is
// also `luar.unproxify` on the Lua side to do this.
//
// We initialize the state with Init() to register Unproxify() as 'luar.unproxify()'.
func ExampleInit() {
	const code = `
-- Lua tables auto-convert to slices in Go-function calls.
local res = foo {10, 20, 30, 40}

-- The result is a map-proxy.
print(res['1'], res['2'])

-- Which we may explicitly convert to a table.
res = luar.unproxify(res)
for k,v in pairs(res) do
	print(k,v)
end
`

	foo := func(args []int) (res map[string]int) {
		res = make(map[string]int)
		for i, val := range args {
			res[strconv.Itoa(i)] = val * val
		}
		return
	}

	L := luar.Init()
	defer L.Close()

	luar.Register(L, "", luar.Map{
		"foo":   foo,
		"print": fmt.Println,
	})

	res := L.DoString(code)
	if res != nil {
		fmt.Println("Error:", res)
	}
	// Output:
	// 400 900
	// 1 400
	// 0 100
	// 3 1600
	// 2 900
}
*/

func ExampleLuaObject_Call() {
	L := luar.Init()
	defer L.Close()

	const code = `
function return_strings()
	return 'one', luar.null, 'three'
end`

	err := L.DoString(code)
	if err != nil {
		log.Fatal(err)
	}

	fun := luar.NewLuaObjectFromName(L, "return_strings")
	defer fun.Close()

	// Using `Call` we would get a generic `[]interface{}`, which is awkward to
	// work with. But the return type can be specified:
	results := []string{}
	err = fun.Call(&results)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(results[0])
	// We get an empty string corresponding to a luar.null in a table,
	// since that's the empty 'zero' value for a string.
	fmt.Println(results[1])
	fmt.Println(results[2])
	// Output:
	// one
	//
	// three
}

func ExampleMap() {
	const code = `
print(#M)
print(M.one)
print(M.two)
print(M.three)
`

	L := luar.Init()
	defer L.Close()

	M := luar.Map{
		"one":   "eins",
		"two":   "zwei",
		"three": "drei",
	}

	luar.Register(L, "", luar.Map{
		"M":     M,
		"print": fmt.Println,
	})

	err := L.DoString(code)
	if err != nil {
		fmt.Println("error", err.Error())
	}
	// Output:
	// 3
	// eins
	// zwei
	// drei
}

func ExampleMakeChan() {
	L1 := luar.Init()
	defer L1.Close()
	L2 := luar.Init()
	defer L2.Close()

	luar.MakeChan(L1)
	L1.SetGlobal("c")
	L1.GetGlobal("c")
	var c interface{}
	_, err := luar.LuaToGo(L1, -1, &c)
	L1.Pop(1)
	if err != nil {
		log.Fatal(err)
	}

	luar.Register(L2, "", luar.Map{"c": c})
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		err := L1.DoString(`c.send(17)`)
		if err != nil {
			fmt.Println(err)
		}
		wg.Done()
	}()

	err = L2.DoString(`return c.recv()`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(L2.ToNumber(-1))
	L2.Pop(1)
	wg.Wait()
	// Output:
	// 17
}

// Another way to do parse configs: use a LuaObject to manipulate the table.
func ExampleNewLuaObject() {
	L := luar.Init()
	defer L.Close()

	const config = `return {
	baggins = true,
	age = 24,
	name = 'dumbo' ,
	marked = {1,2},
	options = {
		leave = true,
		cancel = 'always',
		tags = {strong=true, foolish=true},
	}
}`

	err := L.DoString(config)
	if err != nil {
		log.Fatal(err)
	}

	lo := luar.NewLuaObject(L, -1)
	defer lo.Close()
	// Can get the field itself as a Lua object, and so forth.
	opts, _ := lo.GetObject("options")
	marked, _ := lo.GetObject("marked")

	loPrint := func(lo *luar.LuaObject, subfields ...interface{}) {
		var a interface{}
		lo.Get(&a, subfields...)
		fmt.Printf("%#v\n", a)
	}
	loPrinti := func(lo *luar.LuaObject, idx int64) {
		var a interface{}
		lo.Get(&a, idx)
		fmt.Printf("%.1f\n", a)
	}
	loPrint(lo, "baggins")
	loPrint(lo, "name")
	loPrint(opts, "leave")
	// Note that these Get methods understand nested fields.
	loPrint(lo, "options", "leave")
	loPrint(lo, "options", "tags", "strong")
	// Non-existent nested fields don't panic but return nil.
	loPrint(lo, "options", "tags", "extra", "flakey")
	loPrinti(marked, 1)

	iter, err := lo.Iter()
	if err != nil {
		log.Fatal(err)
	}
	keys := []string{}
	for key := ""; iter.Next(&key, nil); {
		keys = append(keys, key)
	}
	if iter.Error() != nil {
		log.Fatal(iter.Error())
	}
	sort.Strings(keys)

	fmt.Println("Keys:")
	for _, v := range keys {
		fmt.Println(v)
	}
	// Output:
	// true
	// "dumbo"
	// true
	// true
	// true
	// <nil>
	// 1.0
	// Keys:
	// age
	// baggins
	// marked
	// name
	// options

}

func ExampleNewLuaObjectFromValue() {
	L := luar.Init()
	defer L.Close()

	gsub := luar.NewLuaObjectFromName(L, "string", "gsub")
	defer gsub.Close()

	// We do have to explicitly copy the map to a Lua table, because `gsub`
	// will not handle userdata types.
	gmap := luar.NewLuaObjectFromValue(L, luar.Map{
		"NAME": "Dolly",
		"HOME": "where you belong",
	})
	defer gmap.Close()

	var res = new(string)
	err := gsub.Call(&res, "Hello $NAME, go $HOME", "%$(%u+)", gmap)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*res)
	// Output:
	// Hello Dolly, go where you belong
}

func ExampleLuaTableIter_Next() {
	const code = `
return {
	foo = 17,
	bar = 18,
}
`

	L := luar.Init()
	defer L.Close()

	err := L.DoString(code)
	if err != nil {
		log.Fatal(err)
	}

	lo := luar.NewLuaObject(L, -1)
	defer lo.Close()

	iter, err := lo.Iter()
	if err != nil {
		log.Fatal(err)
	}
	keys := []string{}
	values := map[string]float64{}
	for key, value := "", 0.0; iter.Next(&key, &value); {
		keys = append(keys, key)
		values[key] = value
	}
	sort.Strings(keys)

	for _, v := range keys {
		fmt.Println(v, values[v])
	}
	// Output:
	// bar 18
	// foo 17
}

func ExampleRegister_sandbox() {
	const code = `
print("foo")
print(io ~= nil)
print(os == nil)
`

	L := luar.Init()
	defer L.Close()

	res := L.LoadString(code)
	if res != 0 {
		msg := L.ToString(-1)
		fmt.Println("could not compile", msg)
	}

	// Create a empty sandbox.
	L.NewTable()
	// "*" means "use table on top of the stack."
	luar.Register(L, "*", luar.Map{
		"print": fmt.Println,
	})
	env := luar.NewLuaObject(L, -1)
	G := luar.NewLuaObjectFromName(L, "_G")
	defer env.Close()
	defer G.Close()

	// We can copy any Lua object from "G" to env with 'Set', e.g.:
	//   env.Set("print", G.Get("print"))
	// A more convenient and efficient way is to do a bulk copy with 'Setv':
	env.Setv(G, "type", "io", "table")

	// Set up sandbox.
	L.SetfEnv(-2)

	// Run 'code' chunk.
	err := L.Call(0, 0)
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// foo
	// true
	// true
}

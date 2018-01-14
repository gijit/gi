package luar

import (
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/glycerine/golua/lua"
)

type luaTestData struct {
	input string
	want  string
}

// From http://stackoverflow.com/questions/25922437/how-can-i-deep-compare-2-lua-tables-which-may-or-may-not-have-tables-as-keys
const luaDeepEqual = `
function deep_equal(table1, table2)
	local avoid_loops = {}
	local function recurse(t1, t2)
		-- compare value types
		if type(t1) ~= type(t2) then return false end
		-- Base case: compare simple values
		if type(t1) ~= "table" then return t1 == t2 end
		-- Now, on to tables.
		-- First, let's avoid looping forever.
		if avoid_loops[t1] then return avoid_loops[t1] == t2 end
		avoid_loops[t1] = t2
		-- Copy keys from t2
		local t2keys = {}
		local t2tablekeys = {}
		for k, _ in pairs(t2) do
			if type(k) == "table" then table.insert(t2tablekeys, k) end
			t2keys[k] = true
		end
		-- Let's iterate keys from t1
		for k1, v1 in pairs(t1) do
			local v2 = t2[k1]
			if type(k1) == "table" then
				-- if key is a table, we need to find an equivalent one.
				local ok = false
				for i, tk in ipairs(t2tablekeys) do
					if table_eq(k1, tk) and recurse(v1, t2[tk]) then
						table.remove(t2tablekeys, i)
						t2keys[tk] = nil
						ok = true
						break
					end
				end
				if not ok then return false end
			else
				-- t1 has a key which t2 doesn't have, fail.
				if v2 == nil then return false end
				t2keys[k1] = nil
				if not recurse(v1, v2) then return false end
			end
		end
		-- if t2 has a key which t1 doesn't have, fail.
		if next(t2keys) then return false end
		return true
	end
	return recurse(table1, table2)
end
`

// We could use luaDump as a deepEqual but it is not as precise.
const luaDump = `
function dump(t, name)
	local visited = {}
	local function recurse(t, name)
		if type(t) ~= "table" then return tostring(t) end
		-- Let's avoid looping forever.
		if visited[t] then return visited[t] end
		visited[t] = name or ""
		local output = {}
		output[#output+1] = '{'
		-- Sort the keys to have deterministic output.
		local keys = {}
		for k in pairs(t) do
			local name
			if type(k) == 'number' then
				name = '#' .. k
			else
				name = tostring(k)
			end
			table.insert(keys, {tostring(name), k})
		end
		table.sort(keys, function(t1, t2) return t1[1]<t2[1] end)
		-- Dump the table.
		for _, key in ipairs(keys) do
			local k = key[2]
			local v = t[k]
			output[#output+1] = key[1]
			output[#output+1] = '='
			output[#output+1] = recurse(v, key[1])
			output[#output+1] = ', '
		end
		if output[#output] == ', ' then output[#output] = nil end
		output[#output+1] = '}'
		return table.concat(output)
	end
	return recurse(t, name)
end
`

func mustDoString(t *testing.T, L *lua.State, code string) {
	err := L.DoString(code)
	if err != nil {
		t.Fatal(err)
	}
}

func checkStack(t *testing.T, L *lua.State) {
	if L.GetTop() != 0 {
		t.Error("unbalanced stack:", L.GetTop())
	}
}

func runLuaTest(t *testing.T, L *lua.State, tdt []luaTestData) {
	mustDoString(t, L, luaDeepEqual)
	mustDoString(t, L, luaDump)
	for _, test := range tdt {
		mustDoString(t, L, `return deep_equal(`+test.input+`,`+test.want+`)`)
		result := L.ToBoolean(-1)
		L.Pop(1)
		if !result {
			mustDoString(t, L, `return dump(`+test.input+`)`)
			got := L.ToString(-1)
			L.Pop(1)
			mustDoString(t, L, `return dump(`+test.want+`)`)
			want := L.ToString(-1)
			L.Pop(1)
			t.Errorf("got %q, want %q from %q", got, want, test.input)
		}
		checkStack(t, L)
	}
}

type goTestData struct {
	input string
	want  interface{}
	err   string
}

// The 'tdt' contains the input Lua expression that has to be converted to the
// 'want' type. If successful, the result is compared to the 'want' value. If
// not, the error message is compared with 'err'.
//
// Warning:'want' should not be referenced.
func runGoTest(t *testing.T, L *lua.State, tdt []goTestData) {
	for _, test := range tdt {
		mustDoString(t, L, `return `+test.input)
		got := reflect.New(reflect.TypeOf(test.want))
		err := LuaToGo(L, -1, got.Interface())
		L.Pop(1)
		checkStack(t, L)
		if test.err == "" && err != nil {
			t.Error(err)
			continue
		}
		if test.err != "" {
			if err == nil {
				t.Errorf("missing error %q from Lua->Go conversion `%v->%#v`", test.err, test.input, test.want)
			} else if !strings.Contains(err.Error(), test.err) {
				t.Errorf("wrong error %q, want %q from Lua->Go conversion `%v->%#v`", err, test.err, test.input, test.want)
			}
			continue
		}
		got = got.Elem()
		if !reflect.DeepEqual(got.Interface(), test.want) {
			t.Errorf("got %#v, want %#v from Lua->Go conversion of `%v`", got, test.want, test.input)
		}
	}
}

func TestArray(t *testing.T) {
	L := Init()
	defer L.Close()

	a := [2]int{17, 18}
	Register(L, "", Map{"a": a})

	runLuaTest(t, L, []luaTestData{{`a`, `{17, 18}`}})

	// Conversion from sub-type should fail.
	runGoTest(t, L, []goTestData{{`a`, new([2]string), ErrTableConv.Error()}})

	mustDoString(t, L, `a[2] = 180`)
	runGoTest(t, L, []goTestData{
		{`a`, [2]int{17, 180}, ""},
		{`a`, [1]int{17}, ""},
		{`a`, [3]int{17, 180, 0}, ""},
	})
}

func TestChan(t *testing.T) {
	L1 := Init()
	defer L1.Close()
	L2 := Init()
	defer L2.Close()

	c := make(chan int)
	Register(L1, "", Map{"c": c})
	Register(L2, "", Map{"c": c})

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		mustDoString(t, L1, `c.send(17)`)
		wg.Done()
	}()

	mustDoString(t, L2, `return c.recv()`)
	got := L2.ToNumber(-1)
	want := 17.0
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	L2.Pop(1)

	wg.Wait()
	checkStack(t, L1)
	checkStack(t, L2)
}

func TestComplex(t *testing.T) {
	L := Init()
	defer L.Close()

	c := 2 + 3i
	a := newIntA(32)

	Register(L, "", Map{
		"c": c,
		"a": a,
	})

	tdt := []luaTestData{
		{`c`, `luar.complex(2, 3)`},
		{`{c.real, c.imag}`, `{2, 3}`},
		{`c+c`, `luar.complex(4, 6)`},
		{`c-c`, `luar.complex(0, 0)`},
		{`-c`, `luar.complex(-2, -3)`},
		{`2*c`, `luar.complex(4, 6)`},
		// {`c^2`, `luar.complex(4, 6)`},
		{`c / a`, `luar.complex(0.0625, 0.09375)`},
	}

	runLuaTest(t, L, tdt)
}

type list struct {
	V    int
	Next *list
}

func TestCycleGoToLua(t *testing.T) {
	L := Init()
	defer L.Close()

	{
		s := make([]interface{}, 2)
		s[0] = 17
		s[1] = s
		GoToLua(L, s)
		output := L.ToPointer(-1)
		L.RawGeti(-1, 2)
		output_1 := L.ToPointer(-1)
		L.SetTop(0)
		if output != output_1 {
			t.Error("address of repeated element differs")
		}
	}

	{
		s := make([]interface{}, 2)
		s[0] = 17
		s2 := make([]interface{}, 2)
		s2[0] = 18
		s2[1] = s
		s[1] = s2
		GoToLua(L, s)
		output := L.ToPointer(-1)
		L.RawGeti(-1, 2)
		L.RawGeti(-1, 2)
		output_1_1 := L.ToPointer(-1)
		L.SetTop(0)
		if output != output_1_1 {
			t.Error("address of repeated element differs")
		}
	}

	{
		s := map[string]interface{}{}
		s["foo"] = 17
		s["bar"] = s
		GoToLua(L, s)
		output := L.ToPointer(-1)
		L.GetField(-1, "bar")
		output_bar := L.ToPointer(-1)
		L.SetTop(0)
		if output != output_bar {
			t.Error("address of repeated element differs")
		}
	}

	{
		s := map[string]interface{}{}
		s["foo"] = 17
		s2 := map[string]interface{}{}
		s2["bar"] = 18
		s2["baz"] = s
		s["qux"] = s2
		GoToLua(L, s)
		output := L.ToPointer(-1)
		L.GetField(-1, "qux")
		L.GetField(-1, "baz")
		output_qux_baz := L.ToPointer(-1)
		L.SetTop(0)
		if output != output_qux_baz {
			t.Error("address of repeated element differs")
		}
	}

	{
		l1 := &list{V: 17}
		l2 := &list{V: 18}
		l1.Next = l2
		l2.Next = l1
		GoToLua(L, l1)
		output_l1 := L.ToPointer(-1)
		L.GetField(-1, "Next")
		L.GetField(-1, "Next")
		output_l1_l2_l1 := L.ToPointer(-1)
		L.SetTop(0)
		if output_l1 != output_l1_l2_l1 {
			t.Error("address of repeated element differs")
		}
	}

	{
		l1 := &list{V: 17}
		l2 := &list{V: 18}
		l1.Next = l2
		l2.Next = l1
		GoToLua(L, l1)
		// Note that root table is only repeated if we call CopyStructToTable on the
		// pointer.
		output_l1 := L.ToPointer(-1)
		L.GetField(-1, "Next")
		output_l1_l2 := L.ToPointer(-1)
		L.GetField(-1, "Next")
		output_l1_l2_l1 := L.ToPointer(-1)
		L.GetField(-1, "Next")
		output_l1_l2_l1_l2 := L.ToPointer(-1)

		L.SetTop(0)
		if output_l1 != output_l1_l2_l1 || output_l1_l2 != output_l1_l2_l1_l2 {
			t.Error("address of repeated element differs")
		}
	}

	{
		a := [2]interface{}{}
		a[0] = 17
		a[1] = &a

		// Pass reference so that first element can be part of the cycle.
		GoToLua(L, &a)

		L.RawGeti(-1, 1)
		got := L.ToInteger(-1)
		L.Pop(1)
		if got != 17 {
			t.Errorf("got %v, want 17", got)
		}

		p := L.ToPointer(-1)
		L.RawGeti(-1, 2)
		pp := L.ToPointer(-1)
		if p != pp {
			t.Error("address of repeated element differs")
		}
		L.Pop(2)
		checkStack(t, L)
	}
}

// Arrays of interface{} cannot cycle since a conversion over an interface{}
// will yield a slice.
func TestCycleLuaToGo(t *testing.T) {
	L := Init()
	defer L.Close()

	{
		var output []interface{}
		L.DoString(`t = {17}; t[2] = t`)
		L.GetGlobal("t")
		LuaToGo(L, -1, &output)
		L.Pop(1)
		output_1 := output[1].([]interface{})
		if &output_1[0] != &output[0] {
			t.Error("address of repeated element differs")
		}
	}

	{
		var output []interface{}
		L.DoString(`t = {17}; v = {t}; t[2] = v`)
		L.GetGlobal("t")
		err := LuaToGo(L, -1, &output)
		L.Pop(1)
		if err != nil {
			t.Error(err)
		}
		output_1 := output[1].([]interface{})
		output_1_0 := output_1[0].([]interface{})
		if &output_1_0[0] != &output[0] {
			t.Error("address of repeated element differs")
		}
	}

	{
		var output []interface{}
		L.DoString(`t = {17}; v = {t, t}; t[2] = v; t[3] = v; t[4] = t`)
		L.GetGlobal("t")
		err := LuaToGo(L, -1, &output)
		L.Pop(1)
		if err != nil {
			t.Error(err)
		}
		output_2 := output[2].([]interface{})
		output_2_0 := output_2[0].([]interface{})
		if &output_2_0[0] != &output[0] {
			t.Error("address of repeated element differs")
		}
	}

	{
		var output map[string]interface{}
		L.DoString(`t = {foo=17}; t["bar"] = t`)
		L.GetGlobal("t")
		err := LuaToGo(L, -1, &output)
		L.Pop(1)
		if err != nil {
			t.Error(err)
		}
		output_1 := output["bar"].(map[string]interface{})
		output["foo"] = 18
		if output["foo"] != output_1["foo"] {
			t.Error("address of repeated element differs")
		}
	}

	{
		var output map[string]interface{}
		L.DoString(`t = {foo=17}; v = {baz=t}; t["bar"] = v`)
		L.GetGlobal("t")
		err := LuaToGo(L, -1, &output)
		L.Pop(1)
		if err != nil {
			t.Error(err)
		}
		output_bar := output["bar"].(map[string]interface{})
		output_bar_baz := output_bar["baz"].(map[string]interface{})
		output["foo"] = 18
		if output["foo"] != output_bar_baz["foo"] {
			t.Error("address of repeated element differs")
		}
	}

	{
		L.DoString(`t = {V=17}; t.Next = t`)
		L.GetGlobal("t")
		var output *list
		err := LuaToGo(L, -1, &output)
		L.Pop(1)
		if err != nil {
			t.Error(err)
		}
		if output.Next != output {
			t.Error("address of repeated element differs")
		}
	}

	{
		L.DoString(`t1 = {V=17}; t2 = {V=18, Next=t1}; t1.Next=t2`)
		L.GetGlobal("t1")
		var output = list{}
		err := LuaToGo(L, -1, &output)
		if err != nil {
			t.Error(err)
		}
		L.Pop(1)
		if output.Next.Next != &output {
			t.Error("address of repeated element differs")
		}
	}
}

// See if Go values are not garbage collected.
func TestGC(t *testing.T) {
	L := Init()
	defer L.Close()

	Register(L, "", Map{"gc": runtime.GC})
	mustDoString(t, L, `s = luar.slice(2)
s[1] = 10
s[2] = 20
gc()`)
	runLuaTest(t, L, []luaTestData{
		{`s[1]`, `10`},
		{`s[2]`, `20`},
	})
}

func TestGoToLuaFunction(t *testing.T) {
	L := Init()
	defer L.Close()

	multiresult := func(x float32, a string) (float32, string) {
		return x, a
	}

	sum := func(args []float64) float64 {
		res := 0.0
		for _, val := range args {
			res += val
		}
		return res
	}

	sumv := func(args ...float64) float64 {
		return sum(args)
	}

	// [10,20] -> {'0':100, '1':400}
	squares := func(args []int) (res map[string]int) {
		res = make(map[string]int)
		for i, val := range args {
			res[strconv.Itoa(i)] = val * val
		}
		return
	}

	IsNilInterface := func(v interface{}) bool {
		return v == nil
	}

	IsNilPointer := func(v *person) bool {
		return v == nil
	}

	// Trick here: we do not return a pointer to 'person' while GetName() is a
	// method on pointer.
	newDirectPerson := func(name string) person {
		return person{Name: name}
	}

	Register(L, "", Map{
		"multiresult":     multiresult,
		"sum":             sum,
		"sumv":            sumv,
		"squares":         squares,
		"IsNilInterface":  IsNilInterface,
		"IsNilPointer":    IsNilPointer,
		"newDirectPerson": newDirectPerson,
	})

	runLuaTest(t, L, []luaTestData{
		{`{multiresult(42, 'foo')}`, `{42, 'foo'}`},
		{`sum{1, 10, 100}`, `111`},      // Auto-convert table to slice.
		{`sumv(1, 10, 100)`, `111`},     // Variadic call table to slice.
		{`squares{10, 20}['0']`, `100`}, // Proxy return value.
		{`squares{10, 20}['1']`, `400`}, // Proxy return value.
		{`IsNilInterface(nil)`, `true`},
		{`IsNilPointer(nil)`, `true`},
		{`newDirectPerson("Charly").GetName()`, `"Charly"`},
	})
}

func TestLuaObject(t *testing.T) {
	L := Init()
	defer L.Close()

	for _, name := range [][]interface{}{{""}, {"dummy"}, {"table.concat"}} {
		a := NewLuaObjectFromName(L, name...)
		a.Push()
		if L.LTypename(-1) != "nil" {
			t.Errorf(`got %q, want "nil"`, L.LTypename(-1))
		}
		L.Pop(1)
	}
	checkStack(t, L)

	mustDoString(t, L, `t = {10, 20, 30}`)
	a := NewLuaObjectFromName(L, "t")
	defer a.Close()
	err := a.Set(200, 2)
	if err != nil {
		t.Error(err)
	}
	checkStack(t, L)
	res := 0
	err = a.Get(&res, 2)
	if res != 200 {
		t.Errorf(`got %v, want 200`, res)
	}
	checkStack(t, L)

	mustDoString(t, L, `t = {foo=17, bar=18, ["1"]=19, 20, ["qux.quuz"]=21, qux={quuz=22}}`)
	a = NewLuaObjectFromName(L, "t")
	err = a.Set(200, 1)
	if err != nil {
		t.Error(err)
	}
	err = a.Set(190, "1")
	if err != nil {
		t.Error(err)
	}
	checkStack(t, L)

	res = 0
	err = a.Get(&res, 1)
	if res != 200 {
		t.Errorf(`got %v, want 200`, res)
	}
	err = a.Get(&res, "1")
	if res != 190 {
		t.Errorf(`got %v, want 190`, res)
	}
	err = a.Get(&res, "qux.quuz")
	if res != 21 {
		t.Errorf(`got %v, want 21`, res)
	}
	err = a.Get(&res, "qux", "quuz")
	if res != 22 {
		t.Errorf(`got %v, want 22`, res)
	}
	checkStack(t, L)
}

func TestLuaObjectMT(t *testing.T) {
	L := Init()
	defer L.Close()

	string__index := func(L *lua.State) int {
		// Note: we skip error checking.
		v, _ := valueOfProxy(L, 1)
		k := L.ToInteger(2)
		L.PushString(string(v.Index(k).Int()))
		return 1
	}

	string__newindex := func(L *lua.State) int {
		// Note: we skip error checking.
		v, _ := valueOfProxy(L, 1)
		k := L.ToInteger(2)
		val := rune(L.ToInteger(3))
		v.Index(k).Set(reflect.ValueOf(val))
		return 0
	}

	Register(L, "", Map{"a": []rune("foobar")})

	L.GetGlobal("a")
	L.NewTable()
	L.SetMetaMethod("__index", string__index)
	L.SetMetaMethod("__newindex", string__newindex)
	L.SetMetaTable(-2)
	L.Pop(1)

	a := NewLuaObjectFromName(L, "a")
	res := ""
	err := a.Set(rune('F'), 1)
	if err != nil {
		t.Fatal(err)
	}
	checkStack(t, L)

	err = a.Get(&res, 1)
	if err != nil {
		t.Fatal(err)
	}
	if res != "F" {
		t.Errorf(`got %v, want 'F'`, res)
	}
	checkStack(t, L)
}

func TestLuaObjectCall(t *testing.T) {
	L := Init()
	defer L.Close()

	const code = `
function id(...)
	return ...
end
`

	mustDoString(t, L, code)
	arg1 := []string{"a", "b"}
	arg2 := Null
	arg3 := "foobar"

	id := NewLuaObjectFromName(L, "id")

	{
		want := new([]string)
		*want = arg1

		got := new([]string)
		err := id.Call(&got, arg1, arg2, arg3)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
		checkStack(t, L)

		dummy := NewLuaObjectFromName(L, "dummy")
		gotErr := dummy.Call(&got, arg1, arg2, arg3)
		wantErr := ErrLuaObjectCallable
		if gotErr == nil || gotErr != wantErr {
			t.Fatalf("got error %q, want %q", gotErr, wantErr)
		}
		checkStack(t, L)
	}

	{
		want := []interface{}{arg1, nil, arg3}
		got := []interface{}{}
		err := id.Call(&got, arg1, arg2, arg3)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
		checkStack(t, L)

		got = make([]interface{}, 4)
		err = id.Call(&got, arg1, arg2, arg3)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
		checkStack(t, L)
	}

	{
		want := struct {
			Res1 []string
			Res2 int
			Res3 string
		}{
			Res1: arg1,
			Res2: 0,
			Res3: arg3,
		}

		got := struct {
			Res1 []string
			Res2 int
			Res3 string
		}{}
		err := id.Call(&got, arg1, arg2, arg3)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
		checkStack(t, L)
	}
}

func TestLuaObjectCallMT(t *testing.T) {
	L := Init()
	defer L.Close()

	const code = `
a = {17}
setmetatable(a, { __call = function(arg) a[1] = a[1] + arg end })
`

	mustDoString(t, L, code)
	a := NewLuaObjectFromName(L, "a")
	err := a.Call(nil, 2)
	if err != nil {
		t.Fatal(err)
	}

	res := 0
	a.Get(&res, 1)
	if res != 19 {
		t.Fatalf("got %q, want 19", res)
	}
	checkStack(t, L)

	err = a.Call(nil)
	wantErr := "[string \"...\"]:3: attempt to perform arithmetic on local 'arg' (a nil value)"
	if err == nil || err.Error() != wantErr {
		t.Fatalf("got error %q, want %q", err, wantErr)
	}
	checkStack(t, L)
}

func TestLuaObjectIter(t *testing.T) {
	L := Init()
	defer L.Close()

	mustDoString(t, L, `a = {foo=10, bar=20}`)

	a := NewLuaObjectFromName(L, "a")
	iter, err := a.Iter()
	if err != nil {
		t.Fatal(err)
	}
	checkStack(t, L)

	keys := []string{}
	values := map[string]float64{}
	for key, value := "", 0.0; iter.Next(&key, &value); {
		keys = append(keys, key)
		values[key] = value
	}
	sort.Strings(keys)

	wantKeys := []string{"bar", "foo"}
	if !reflect.DeepEqual(keys, wantKeys) {
		t.Errorf("got %q, want %q", keys, wantKeys)
	}

	wantValues := map[string]float64{"foo": 10, "bar": 20}
	if !reflect.DeepEqual(values, wantValues) {
		t.Errorf("got %q, want %q", keys, wantValues)
	}

	checkStack(t, L)
}

func TestLuaObjectIterMT(t *testing.T) {
	L := Init()
	defer L.Close()

	Register(L, "", Map{"a": map[string]float64{"foo": 10, "bar": 20}})

	a := NewLuaObjectFromName(L, "a")

	iter, err := a.Iter()
	if err != nil {
		t.Fatal(err)
	}
	checkStack(t, L)

	// Remove the metatable to see if it does not corrupt the iterator.
	L.GetGlobal("a")
	L.PushNil()
	L.SetMetaTable(-2)
	L.Pop(1)

	keys := []string{}
	values := map[string]float64{}
	for key, value := "", 0.0; iter.Next(&key, &value); {
		keys = append(keys, key)
		values[key] = value
	}
	sort.Strings(keys)
	checkStack(t, L)

	if iter.Error() != nil {
		t.Errorf("%q", iter.Error())
	}

	wantKeys := []string{"bar", "foo"}
	if !reflect.DeepEqual(keys, wantKeys) {
		t.Errorf("got %q, want %q", keys, wantKeys)
	}

	wantValues := map[string]float64{"foo": 10, "bar": 20}
	if !reflect.DeepEqual(values, wantValues) {
		t.Errorf("got %q, want %q", keys, wantValues)
	}

	checkStack(t, L)
}

func TestLuaToGoPointers(t *testing.T) {
	L := Init()
	defer L.Close()

	L.PushInteger(17)
	var err error

	// TODO: Factor with runGoTest?
	printError := func(want string) {
		if want == "" && err != nil {
			t.Error(err)
		}
		if want != "" {
			if err == nil {
				t.Errorf("missing error %q", want)
			} else if !strings.Contains(err.Error(), want) {
				t.Errorf("wrong error %q, want %q", err, want)
			}
		}
	}

	// nil pointer
	err = LuaToGo(L, -1, nil)
	printError("not a pointer")

	// pointer to nil
	var ip *int
	err = LuaToGo(L, -1, ip)
	printError("nil pointer")

	// pointer to zero
	var i int
	ip = &i
	err = LuaToGo(L, -1, ip)
	printError("")

	// pointer to pointer to nil
	var ipp **int
	err = LuaToGo(L, -1, ipp)
	printError("nil pointer")

	// pointer to pointer to zero
	ipp = &ip
	err = LuaToGo(L, -1, ipp)
	printError("")

	L.Pop(1)

	// Test pointer in function arguments.
	foo := func(i *int) int {
		return *i
	}
	Register(L, "", Map{"foo": foo})
	runLuaTest(t, L, []luaTestData{{`foo(17)`, `17`}})
}

type myMap map[string]int

func (m *myMap) Foo() int {
	return len(*m)
}

type myIntMap map[int]int

func (m *myIntMap) Foo() int {
	return len(*m)
}

func TestMap(t *testing.T) {
	L := Init()
	defer L.Close()

	want := map[string]int{
		"foo":  170,
		"qux":  18,
		"quux": 19,
	}
	GoToLua(L, want)
	L.SetGlobal("a")

	runLuaTest(t, L, []luaTestData{
		{`a`, `{foo=170, qux=18, quux=19}`},
	})

	input := `{foo=170, bar="baz", qux=18, "idx1"}`
	mustDoString(t, L, `return `+input)
	got := map[string]int{
		"foo":  17,
		"quux": 19,
	}

	err := LuaToGo(L, -1, &got)
	if err != ErrTableConv {
		t.Errorf("wrong error %q, want %q", err, ErrTableConv)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v from Lua->Go conversion of `%v`", got, want, input)
	}

	got = nil
	want = map[string]int{
		"foo": 170,
		"qux": 18,
	}
	err = LuaToGo(L, -1, &got)
	if err != ErrTableConv {
		t.Errorf("wrong error %q, want %q", err, ErrTableConv)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v from Lua->Go conversion of `%v`", got, want, input)
	}

	var i interface{}
	want2 := map[string]interface{}{
		"foo": 170.0,
		"bar": "baz",
		"qux": 18.0,
	}
	err = LuaToGo(L, -1, &i)
	if err != ErrTableConv {
		t.Errorf("wrong error %q, want %q", err, ErrTableConv)
	}
	if !reflect.DeepEqual(i, want2) {
		t.Errorf("got %#v, want %#v from Lua->Go conversion of `%v`", i, want2, input)
	}

	i = map[string]interface{}{
		"foo":  17.0,
		"quux": 19.0,
	}
	want3 := map[string]interface{}{
		"foo":  170.0,
		"bar":  "baz",
		"qux":  18.0,
		"quux": 19.0,
	}
	err = LuaToGo(L, -1, &i)
	if err != ErrTableConv {
		t.Errorf("wrong error %q, want %q", err, ErrTableConv)
	}
	if !reflect.DeepEqual(i, want3) {
		t.Errorf("got %#v, want %#v from Lua->Go conversion of `%v`", i, want3, input)
	}
}

type hasName interface {
	GetName() string
}

func (p *person) GetName() string {
	return p.Name
}

func newPerson(name string, age int) *person {
	return &person{name, age}
}

func newName(t *person) hasName {
	return t
}

func getName(o hasName) string {
	return o.GetName()
}

func TestProxy(t *testing.T) {
	L := Init()
	defer L.Close()

	Register(L, "", Map{"a": myIntA(17)})

	runGoTest(t, L, []goTestData{
		{`a`, myIntA(17), ""},
		{`a`, myIntB(17), ""},
		{`a`, 17, ""},
		{`a`, uint(17), ""},
	})

	// runGoTest(t, L, []goTestData{{`a`, i, ""}})
	mustDoString(t, L, `return a`)
	var i interface{}
	want := myIntA(17)

	err := LuaToGo(L, -1, &i)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(i, want) {
		t.Errorf("got %#v, want %#v from Lua->Go interface conversion", i, want)
	}

	i = "foo"
	err = LuaToGo(L, -1, &i)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(i, want) {
		t.Errorf("got %#v, want %#v from Lua->Go interface conversion", i, want)
	}

	L.Pop(1)

	a := myIntA(17)
	pa := &a
	ppa := &pa

	b := myIntB(17)
	pb := &b
	ppb := &pb

	s := true
	ps := &s
	pps := &ps
	runGoTest(t, L, []goTestData{
		{`a`, a, ""},
		{`a`, pa, ""},
		{`a`, ppa, ""},
		{`a`, b, ""},
		{`a`, pb, ""},
		{`a`, ppb, ""},
		{`a`, s, "cannot convert proxy (luar.myIntA) to bool"},
		{`a`, ps, "cannot convert proxy (luar.myIntA) to bool"},
		{`a`, pps, "cannot convert proxy (luar.myIntA) to bool"},
	})

	runGoTest(t, L, []goTestData{
		{`luar.null`, myIntA(0), ""},
		{`luar.null`, 0, ""},
		{`luar.null`, "", ""},
	})
}

func TestProxyArray(t *testing.T) {
	L := Init()
	defer L.Close()

	a := [2]int{17, 18}
	Register(L, "", Map{"a": &a})

	tdt := []luaTestData{
		{`#a`, `2`},
		{`type(a)`, `'table<[2]int>'`},
		{`a[1]`, `17`},
		{`a[2]`, `18`},
	}

	runLuaTest(t, L, tdt)

	mustDoString(t, L, `a1 = {ipairs(a)(a, 0)}; a1 = a1[2]`)
	runLuaTest(t, L, []luaTestData{{`a1`, `17`}})

	mustDoString(t, L, `a[1] = 170`)
	runGoTest(t, L, []goTestData{
		{`a`, &[2]int{170, 18}, ""},
	})

	// 'ipairs' on regular tables.
	mustDoString(t, L, `t = {37}; t1 = {ipairs(t)(t, 0)}; t1 = t1[2]`)
	runLuaTest(t, L, []luaTestData{{`t1`, `37`}})
}

type myIntA int

func newIntA(i int) myIntA {
	return myIntA(i)
}

func (a myIntA) FooIntA() string {
	return "FooIntA"
}

type myIntB int

func NewIntB(i int) myIntB {
	return myIntB(i)
}

type myStringA string

func newStringA(s string) myStringA {
	return myStringA(s)
}

func (c myStringA) FooStringA() string {
	return "FooStringA"
}

type myStringB string

func newStringB(s string) myStringB {
	return myStringB(s)
}

func (d myStringB) FooStringB() string {
	return "FooStringB"
}

func TestProxyScalars(t *testing.T) {
	L := Init()
	defer L.Close()

	i := myIntA(3)
	j := myIntB(17)
	s1 := myStringA("foo")
	s2 := myStringB("bar")

	Register(L, "", Map{
		"i":          i,
		"j":          j,
		"s1":         s1,
		"s2":         s2,
		"newIntA":    newIntA,
		"newStringA": newStringA,
	})

	runLuaTest(t, L, []luaTestData{
		// Number proxy
		{`i`, `newIntA(3)`},
		// {`i`, `17`}, // Not equal.
		{`i+i`, `newIntA(6)`},
		{`i-i`, `newIntA(0)`},
		{`-i`, `newIntA(-3)`},
		{`i*i`, `newIntA(9)`},
		{`i/i`, `newIntA(1)`},
		{`i^i`, `newIntA(27)`},
		{`i%i`, `newIntA(0)`},
		{`i==i`, `true`},
		{`i~=i`, `false`},
		{`i<i`, `false`},
		{`i<=i`, `true`},
		{`i>i`, `false`},
		{`i>=i`, `true`},
		// Number proxy & number
		{`i+2`, `newIntA(5)`},
		{`i*2`, `newIntA(6)`},
		{`i/2`, `newIntA(1)`},
		{`i^2`, `newIntA(9)`},
		{`i%2`, `newIntA(1)`},
		{`i==3`, `false`},
		{`i~=3`, `true`},
		// <, >, <= and >= do not work between number/string and userdata.
		// Number & number proxy
		{`2+i`, `newIntA(5)`},
		{`2*i`, `newIntA(6)`},
		{`2/i`, `newIntA(0)`},
		{`2^i`, `newIntA(8)`},
		{`2%i`, `newIntA(2)`},
		{`3==i`, `false`},
		{`3~=i`, `true`},
		// Proxy A & proxy B
		{`i+j`, `20`},
		{`i<j`, `true`},
		// Proxy B & proxy A
		{`j+i`, `20`},
		{`j<i`, `false`},
		{`i.FooIntA()`, `"FooIntA"`},
		// Strings.
		{`s1`, `newStringA("foo")`},
		// {`s1`, `"foo"`}, // Not equal.
		{`#s1`, `3`},
		{`s1 .. "bar"`, `newStringA("foobar")`},
		{`"bar" .. s1`, `newStringA("barfoo")`},
		{`s1 < s2`, `false`},
		{`s1 .. s1`, `newStringA("foofoo")`},
		{`s1 .. s2`, `"foobar"`},
		{`s1 .. 17`, `newStringA("foo17")`},
		{`s1 .. i`, `"foo3"`},
		{`s1.FooStringA()`, `"FooStringA"`},
	})
}

func TestProxyMap(t *testing.T) {
	L := Init()
	defer L.Close()

	a := myMap{"foo": 17, "bar": 170}
	b := myMap{"Foo": 17}
	c := myIntMap{1: 10, 2: 20}
	m := map[interface{}]string{
		-1:  "ko",
		0:   "ko",
		1:   "foo",
		2:   "bar",
		"3": "baz",
	}
	Register(L, "", Map{"a": a, "b": b, "c": c, "m": m})

	runLuaTest(t, L, []luaTestData{
		{`a.Foo()`, `2`},
		{`a.foo`, `17`},
		{`a.empty`, `nil`},
		{`b.Foo`, `17`},
		{`luar.method(b, "Foo")()`, `1`},
		{`luar.method(b, "Nonexistent")`, `nil`},
		{`luar.method(nil, "Nonproxy")`, `nil`},
		{`c[1]`, `10`},
		{`c[2]`, `20`},
		{`c.Foo()`, `2`},
		{`c.bar`, `nil`},
	})

	mustDoString(t, L, `t = {}
for k, v in ipairs(m) do
t[k] = v
end`)
	runLuaTest(t, L, []luaTestData{{`t`, `{'foo', 'bar'}`}})

	mustDoString(t, L, `
n = luar.map()
n.foo = "bar"
n.baz = "qux"
u = luar.unproxify(n)
`)
	runLuaTest(t, L, []luaTestData{{`u`, `{foo="bar", baz="qux"}`}})

	mustDoString(t, L, `
p = {}
for k, v in pairs(n) do
p[k] = v
end
`)
	runLuaTest(t, L, []luaTestData{{`p`, `{foo="bar", baz="qux"}`}})
}

type mySlice []int

func (m *mySlice) Foo() int {
	return len(*m)
}

func TestProxySlice(t *testing.T) {
	L := Init()
	defer L.Close()

	a := mySlice{17, 170}
	Register(L, "", Map{"a": a})
	mustDoString(t, L, `a = a.append(18.5, 19)`)
	mustDoString(t, L, `a = a.append(unpack({3, 2}))`)
	runLuaTest(t, L, []luaTestData{
		{`a.Foo()`, `6`},
		{`a[1]`, `17`},
		{`a.slice(1, 2)[1]`, `17`},
		{`a.slice(#a, #a+1)[1]`, `2`},
		{`a.slice(3, 5)[1]`, `18`},
		{`a.slice(3, 5)[2]`, `19`},
	})
}

func TestProxyString(t *testing.T) {
	L := Init()
	defer L.Close()

	a := myStringA("naïveté")

	Register(L, "", Map{
		"a":          a,
		"newStringA": newStringA,
	})

	const code = `
for k, v in ipairs(a) do
if k == 3 then
a3 = v
break
end
end
`

	mustDoString(t, L, code)
	runLuaTest(t, L, []luaTestData{
		{`a3`, `'ï'`},
		{`a[1]`, `'n'`}, // Go string indexing does not support unicode.
		{`a.slice(2, 3)`, `newStringA('a')`},
		{`a.slice(2, 3)[1]`, `'a'`},
	})
}

// Get and set public fields in struct proxies.
// Test interface conversion and calls.
func TestProxyStruct(t *testing.T) {
	L := Init()
	defer L.Close()

	Register(L, "", Map{
		"NewPerson": newPerson,
		"NewName":   newName,
		"GetName":   getName,
	})

	mustDoString(t, L, `t = NewPerson("Alice", 17)`)
	runLuaTest(t, L, []luaTestData{
		{`t.Name`, `'Alice'`},
		{`t.Age`, `17`},
	})

	mustDoString(t, L, `t.Name = 'Bob'`)
	runLuaTest(t, L, []luaTestData{
		{`t.GetName()`, `'Bob'`},
	})

	mustDoString(t, L, `it = NewName(t)`)
	runLuaTest(t, L, []luaTestData{
		{`it.GetName()`, `'Bob'`},
		{`GetName(it)`, `'Bob'`},
		{`GetName(t)`, `'Bob'`},
		{`type(t)`, `'table<luar.person>'`},
		{`type(it)`, `'table<luar.person>'`},
	})
}

// nil, bool, number, string
func TestScalar(t *testing.T) {
	L := Init()
	defer L.Close()

	type bar int
	runGoTest(t, L, []goTestData{
		{`nil`, "", ""},
		{`nil`, 0, ""},
		{`true`, true, ""},
		{`17`, 17, ""},
		{`17`, int16(17), ""},
		{`17`, float32(17), ""},
		{`17`, bar(17), ""},
		{`-1.0`, -1, ""},
		{`-1.7`, -1.7, ""},
		{`"foo"`, "foo", ""},
		{`true`, int16(17), "cannot convert Lua value 'true' (boolean) to int16"},
		{`17`, "17", "cannot convert Lua value '17' (number) to string"},
	})

	var i interface{}

	L.PushNil()
	err := LuaToGo(L, -1, &i)
	L.Pop(1)
	if err != nil {
		t.Error(err)
	}
	if i != nil {
		t.Errorf("got %T, expected 'nil' from Lua conversion to Go interface", i)
	}

	L.PushBoolean(true)
	err = LuaToGo(L, -1, &i)
	L.Pop(1)
	if err != nil {
		t.Error(err)
	}
	ibool, ok := i.(bool)
	if !ok {
		t.Errorf("got %T, expected type 'bool' from Lua conversion to Go interface", ibool)
	} else if !ibool {
		t.Errorf("got %v, expected 'true' from Lua conversion to Go interface", ibool)
	}

	L.PushNumber(17)
	err = LuaToGo(L, -1, &i)
	L.Pop(1)
	if err != nil {
		t.Error(err)
	}
	ifloat, ok := i.(float64)
	if !ok {
		t.Errorf("got %T, expected type 'float64' from Lua conversion to Go interface", ifloat)
	} else if ifloat != 17 {
		t.Errorf("got %v, expected '17' from Lua conversion to Go interface", i)
	}
}

func TestSlice(t *testing.T) {
	L := Init()
	defer L.Close()

	want := []string{"idx1", "idx2", "idx3"}
	GoToLua(L, want)
	L.SetGlobal("a")

	runLuaTest(t, L, []luaTestData{
		{`a`, `{"idx1", "idx2", "idx3"}`},
	})

	input := `{"idx1", "idx2", "idx3"}`
	mustDoString(t, L, `return `+input)

	got := []string{"tooshort"}
	err := LuaToGo(L, -1, &got)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v from Lua->Go conversion of `%v`", got, want, input)
	}

	got = []string{"more", "than", "three", "elements"}
	err = LuaToGo(L, -1, &got)
	if err != nil {
		t.Error(nil)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v from Lua->Go conversion of `%v`", got, want, input)
	}

	got = nil
	err = LuaToGo(L, -1, &got)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v from Lua->Go conversion of `%v`", got, want, input)
	}

	var i interface{}
	want2 := []interface{}{"idx1", "idx2", "idx3"}
	err = LuaToGo(L, -1, &i)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(i, want2) {
		t.Errorf("got %#v, want %#v from Lua->Go conversion of `%v`", i, want2, input)
	}

	L.Pop(1)
	input = `{"idx1", "idx2", 17, "idx4"}`
	mustDoString(t, L, `return `+input)
	i = []string{"foo", "bar"}
	want3 := []string{"idx1", "idx2", "", "idx4"}
	err = LuaToGo(L, -1, &i)
	if err != ErrTableConv {
		t.Errorf("wrong error %q, want %q", err, ErrTableConv)
	}
	if !reflect.DeepEqual(i, want3) {
		t.Errorf("got %#v, want %#v from Lua->Go conversion of `%v`", i, want3, input)
	}

}

type person struct {
	Name string
	Age  int
}

type personWithTags struct {
	Name string `lua:"name"`
	Age  int    `lua:"year"`
}

func TestStruct(t *testing.T) {
	L := Init()
	defer L.Close()

	want := person{Name: "foo", Age: 17}
	Register(L, "", Map{"a": want})
	runLuaTest(t, L, []luaTestData{{`a`, `{Name='foo', Age=17}`}})

	wantTags := personWithTags{Name: "foo", Age: 17}
	Register(L, "", Map{"a": want, "atags": wantTags})
	runLuaTest(t, L, []luaTestData{{`atags`, `{name='foo', year=17}`}})
	runGoTest(t, L, []goTestData{{`atags`, wantTags, ""}})

	input := `{Name="foo", Ignored="baz"}`
	mustDoString(t, L, `return `+input)
	got := person{Name: "bar", Age: 17}
	err := LuaToGo(L, -1, &got)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v from Lua->Go conversion of `%v`", got, want, input)
	}

	L.Pop(1)
	input = `{Name="foo", Age="17yo", Ignored="baz"}`
	mustDoString(t, L, `return `+input)
	got = person{Name: "bar", Age: 17}
	err = LuaToGo(L, -1, &got)
	if err != ErrTableConv {
		t.Errorf("wrong error %q, want %q", err, ErrTableConv)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v from Lua->Go conversion of `%v`", got, want, input)
	}

	got = person{}
	want = person{Name: "foo", Age: 0}
	err = LuaToGo(L, -1, &got)
	if err != ErrTableConv {
		t.Errorf("wrong error %q, want %q", err, ErrTableConv)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v from Lua->Go conversion of `%v`", got, want, input)
	}
}

// 'nil' in Go slices and maps is represented by luar.null.
func TestUnproxify(t *testing.T) {
	L := Init()
	defer L.Close()

	s := [][]int{nil, {1, 2}, nil, {10, 20}}
	m := map[string][]int{
		"a": {1, 2},
		"b": nil,
		"c": {10, 20},
		"d": nil,
	}
	Register(L, "", Map{"s": s, "m": m})

	mustDoString(t, L, `ts = luar.unproxify(s)`)
	runLuaTest(t, L, []luaTestData{{`ts`, `{luar.null, {1, 2}, luar.null, {10, 20}}`}})

	mustDoString(t, L, `tm = luar.unproxify(m)`)
	runLuaTest(t, L, []luaTestData{{`tm`, `{a={1, 2}, b=luar.null, c={10, 20}, d=luar.null}`}})
}

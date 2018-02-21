package luar

// TODO: v2 seems to be somewhat slower than v1. Profile and optimize.

import (
	"fmt"
	"testing"

	"github.com/aarzilli/golua/lua"
)

func BenchmarkLuaToGoSliceInt(b *testing.B) {
	L := Init()
	defer L.Close()

	var output []interface{}
	L.DoString(`t={}; for i = 1,100 do t[i]=i; end`)
	L.GetGlobal("t")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		LuaToGo(L, -1, &output)
	}
}

func BenchmarkLuaToGoSliceMap(b *testing.B) {
	L := Init()
	defer L.Close()

	var output []interface{}
	L.DoString(`t={}; s={17}; for i = 1,100 do t[i]=s; end`)
	L.GetGlobal("t")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		LuaToGo(L, -1, &output)
	}
}

func BenchmarkLuaToGoSliceMapUnique(b *testing.B) {
	L := Init()
	defer L.Close()

	var output []interface{}
	L.DoString(`t={}`)
	for i := 0; i < 100; i++ {
		L.DoString(fmt.Sprintf(`s%[1]d={17}; t[%[1]d]=s%[1]d`, i))
	}
	L.GetGlobal("t")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		LuaToGo(L, -1, &output)
	}
}

func BenchmarkLuaToGoMapInt(b *testing.B) {
	L := Init()
	defer L.Close()

	var output map[string]interface{}
	L.DoString(`t={}; for i = 1,100 do t[tostring(i)]=i; end`)
	L.GetGlobal("t")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		LuaToGo(L, -1, &output)
	}
}

func BenchmarkLuaToGoMapSlice(b *testing.B) {
	L := Init()
	defer L.Close()

	var output map[string]interface{}
	L.DoString(`t={}; s={17}; for i = 1,100 do t[tostring(i)]=s; end`)
	L.GetGlobal("t")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		LuaToGo(L, -1, &output)
	}
}

func BenchmarkLuaToGoMapSliceUnique(b *testing.B) {
	L := Init()
	defer L.Close()

	var output map[string]interface{}
	L.DoString(`t={}`)
	for i := 0; i < 100; i++ {
		L.DoString(fmt.Sprintf(`s%[1]d={17}; t["%[1]d"]=s%[1]d`, i))
	}
	L.GetGlobal("t")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		LuaToGo(L, -1, &output)
	}
}

func BenchmarkGoToLuaSliceInt(b *testing.B) {
	L := Init()
	defer L.Close()

	input := make([]int, 100)
	for i := 0; i < 100; i++ {
		input[i] = i
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		GoToLua(L, input)
		L.SetTop(0)
	}
}

func BenchmarkGoToLuaSliceSlice(b *testing.B) {
	L := Init()
	defer L.Close()

	sub := []int{17}
	input := make([][]int, 100)
	for i := 0; i < 100; i++ {
		input[i] = sub
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		GoToLua(L, input)
		L.SetTop(0)
	}
}

func BenchmarkGoToLuaSliceSliceUnique(b *testing.B) {
	L := Init()
	defer L.Close()

	input := make([][]int, 100)
	for i := 0; i < 100; i++ {
		input[i] = []int{17}
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		GoToLua(L, input)
		L.SetTop(0)
	}
}

func BenchmarkGoToLuaMapInt(b *testing.B) {
	L := Init()
	defer L.Close()

	input := map[int]int{}
	for i := 0; i < 100; i++ {
		input[i] = i
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		GoToLua(L, input)
		L.SetTop(0)
	}
}

func BenchmarkGoToLuaMapSlice(b *testing.B) {
	L := Init()
	defer L.Close()

	sub := []int{17}
	input := map[int][]int{}
	for i := 0; i < 100; i++ {
		input[i] = sub
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		GoToLua(L, input)
		L.SetTop(0)
	}
}

func BenchmarkGoToLuaMapSliceUnique(b *testing.B) {
	L := Init()
	defer L.Close()

	input := map[int][]int{}
	for i := 0; i < 100; i++ {
		input[i] = []int{17}
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		GoToLua(L, input)
		L.SetTop(0)
	}
}

const codePairs = `
local t = {
	a = "a",
	b = "b",
	c = "c",
	d = "d",
}
for i = 1, 200 do
	t[i] = i
end

local pairs = pairs
local ipairs = ipairs
local luar_pairs = luar.pairs
local luar_ipairs = luar.ipairs
print("pairs", pairs, "luar_pairs", luar_pairs)
print("ipairs", ipairs, "luar_ipairs", luar_ipairs)

function pairs_test()
	local tmp
	for i = 1, 1000 do
		for k, v in pairs(t) do
			tmp = v
		end
	end
end

function luar_pairs_test()
	local tmp
	for i = 1, 1000 do
		for k, v in luar_pairs(t) do
			tmp = v
		end
	end
end

function ipairs_test()
	local tmp
	for i = 1, 1000 do
		for k, v in ipairs(t) do
			tmp = v
		end
	end
end

function luar_ipairs_test()
	local tmp
	for i = 1, 1000 do
		for k, v in luar_ipairs(t) do
			tmp = v
		end
	end
end
`

func BenchmarkPairs(b *testing.B) {
	var L = lua.NewState()
	defer L.Close()
	L.OpenLibs()
	Register(L, "luar", Map{
		"pairs": ProxyPairs,
	})
	L.DoString(codePairs)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = L.DoString("pairs_test()")
	}
}

func BenchmarkLuarPairs(b *testing.B) {
	var L = lua.NewState()
	defer L.Close()
	L.OpenLibs()
	Register(L, "luar", Map{
		"pairs": ProxyPairs,
	})
	L.DoString(codePairs)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = L.DoString("luar_pairs_test()")
	}
}

func BenchmarkIpairs(b *testing.B) {
	var L = lua.NewState()
	defer L.Close()
	L.OpenLibs()
	RegProxyIpairs(L, "luar", "ipairs")
	L.DoString(codePairs)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = L.DoString("ipairs_test()")
	}
}

func BenchmarkLuarIpairs(b *testing.B) {
	var L = lua.NewState()
	defer L.Close()
	L.OpenLibs()
	RegProxyIpairs(L, "luar", "ipairs")
	L.DoString(codePairs)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = L.DoString("luar_ipairs_test()")
	}
}

package luar

// TODO: v2 seems to be somewhat slower than v1. Profile and optimize.

import (
	"fmt"
	"testing"
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

package compiler

import (
	"fmt"

	"github.com/gijit/gi/pkg/token"
	"github.com/gijit/gi/pkg/types"
	"github.com/glycerine/luar"
)

func (ic *IncrState) GiImportFunc(path string) (*Archive, error) {

	// `import "fmt"` means that path == "fmt", for example.
	pp("GiImportFunc called with path = '%s'", path)

	//panic("where import?")

	pkg := types.NewPackage("fmt", "fmt")
	pkg.MarkComplete()
	scope := pkg.Scope()

	fun := getFunForSprintf(pkg)

	// As it should, scope.Insert(fun)
	// gets rid of 'Sprintf not declared by package fmt'
	// from types/call.go:302.
	scope.Insert(fun)

	summer := getFunForSummer(pkg)
	scope.Insert(summer)

	summerAny := getFunForSummerAny(pkg)
	scope.Insert(summerAny)

	incr := getFunForIncr(pkg)
	scope.Insert(incr)

	ic.importContext.Packages[path] = pkg

	// implementation via luar-based reflection

	// fmt
	luar.Register(ic.vm, "fmt", luar.Map{
		// Go functions may be registered directly.
		"Sprintf":   fmt.Sprintf,
		"Summer":    Summer,
		"SummerAny": SummerAny,
		"Incr":      Incr,
	})

	return &Archive{
		ImportPath: path,
		pkg:        pkg,
	}, nil
}

func getFunForSprintf(pkg *types.Package) *types.Func {
	// func Sprintf(format string, a ...interface{}) string
	var recv *types.Var
	var T types.Type = &types.Interface{}
	str := types.Typ[types.String]
	results := types.NewTuple(types.NewVar(token.NoPos, pkg, "", str))
	params := types.NewTuple(types.NewVar(token.NoPos, pkg, "format", str),
		types.NewVar(token.NoPos, pkg, "a", types.NewSlice(T)))
	variadic := true
	sig := types.NewSignature(recv, params, results, variadic)
	fun := types.NewFunc(token.NoPos, pkg, "Sprintf", sig)
	return fun
}

func Summer(a, b int) int {
	return a + b
}

func getFunForSummer(pkg *types.Package) *types.Func {
	// func Summer(a, b int) int
	var recv *types.Var
	nt := types.Typ[types.Int]
	results := types.NewTuple(types.NewVar(token.NoPos, pkg, "", nt))
	params := types.NewTuple(types.NewVar(token.NoPos, pkg, "a", nt),
		types.NewVar(token.NoPos, pkg, "b", nt))
	variadic := false
	sig := types.NewSignature(recv, params, results, variadic)
	fun := types.NewFunc(token.NoPos, pkg, "Summer", sig)
	return fun
}

func SummerAny(a ...int) int {
	tot := 0
	for i := range a {
		tot += a[i]
	}
	return tot
}

func getFunForSummerAny(pkg *types.Package) *types.Func {
	// func Summer(a, b int) int
	var recv *types.Var
	nt := types.Typ[types.Int]
	results := types.NewTuple(types.NewVar(token.NoPos, pkg, "", nt))
	params := types.NewTuple(types.NewVar(token.NoPos, pkg, "a", types.NewSlice(nt)))
	variadic := true
	sig := types.NewSignature(recv, params, results, variadic)
	fun := types.NewFunc(token.NoPos, pkg, "SummerAny", sig)
	return fun
}

func Incr(a int) int {
	fmt.Printf("\nYAY Incr(a) called! with a = '%v'\n", a)
	return a + 1
}

func getFunForIncr(pkg *types.Package) *types.Func {
	// func Incr(a int) int
	var recv *types.Var
	nt := types.Typ[types.Int]
	results := types.NewTuple(types.NewVar(token.NoPos, pkg, "", nt))
	params := types.NewTuple(types.NewVar(token.NoPos, pkg, "a", nt))
	variadic := false
	sig := types.NewSignature(recv, params, results, variadic)
	fun := types.NewFunc(token.NoPos, pkg, "Incr", sig)
	return fun
}

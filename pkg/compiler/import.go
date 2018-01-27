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
	var pkg *types.Package

	switch path {
	case "fmt":
		pkg = types.NewPackage("fmt", "fmt")
		pkg.MarkComplete()
		scope := pkg.Scope()

		// These scope.Insert() calls let us get
		// past the Go type checker.

		// As it should, scope.Insert(fun)
		// gets rid of 'Sprintf not declared by package fmt'
		// from types/call.go:302.
		fun := getFunForSprintf(pkg)
		scope.Insert(fun)

		scope.Insert(getFunForPrintf(pkg))

		summer := getFunForSummer(pkg)
		scope.Insert(summer)

		summerAny := getFunForSummerAny(pkg)
		scope.Insert(summerAny)

		incr := getFunForIncr(pkg)
		scope.Insert(incr)

		// implementation via luar-based reflection

		// fmt
		luar.Register(ic.vm, "fmt", luar.Map{
			// Go functions may be registered directly.
			"Sprintf":   fmt.Sprintf,
			"Printf":    fmt.Printf,
			"Summer":    Summer,
			"SummerAny": SummerAny,
			"Incr":      Incr,
		})

	case "gitesting":
		// test only:
		if !ic.vmCfg.NotTestMode {
			fmt.Print("\n registering gitesting.SumArrayInt64! \n")
			pkg = types.NewPackage("gitesting", "gitesting")
			pkg.MarkComplete()
			scope := pkg.Scope()

			fun := getFunForSumArrayInt64(pkg)
			scope.Insert(fun)

			luar.Register(ic.vm, "gitesting", luar.Map{
				"SumArrayInt64": sumArrayInt64,
				//"__giClone":     __giClone,
			})
		}
	} // end switch on path

	ic.importContext.Packages[path] = pkg
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

func getFunForPrintf(pkg *types.Package) *types.Func {
	// func Sprintf(format string, a ...interface{}) string
	var recv *types.Var
	var T types.Type = &types.Interface{}
	str := types.Typ[types.String]
	nt := types.Typ[types.Int]
	errt := types.Universe.Lookup("error")
	if errt == nil {
		panic("could not locate error interface in types.Universe")
	}
	results := types.NewTuple(types.NewVar(token.NoPos, pkg, "", nt),
		types.NewVar(token.NoPos, pkg, "", errt.Type()))
	params := types.NewTuple(types.NewVar(token.NoPos, pkg, "format", str),
		types.NewVar(token.NoPos, pkg, "a", types.NewSlice(T)))
	variadic := true
	sig := types.NewSignature(recv, params, results, variadic)
	fun := types.NewFunc(token.NoPos, pkg, "Printf", sig)
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
	fmt.Printf("top of SummaryAny, a is len %v\n", len(a))
	tot := 0
	for i := range a {
		tot += a[i]
	}
	fmt.Printf("end of SummaryAny, returning tot=%v\n", tot)
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

func getFunForSumArrayInt64(pkg *types.Package) *types.Func {
	// func sumArrayInt64(a [3]int64) (tot int64)
	var recv *types.Var
	nt64 := types.Typ[types.Int64]
	results := types.NewTuple(types.NewVar(token.NoPos, pkg, "tot", nt64))
	params := types.NewTuple(types.NewVar(token.NoPos, pkg, "a", types.NewArray(nt64, 3)))
	variadic := false
	sig := types.NewSignature(recv, params, results, variadic)
	fun := types.NewFunc(token.NoPos, pkg, "SumArrayInt64", sig)
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

//"__giClone", __giClone,

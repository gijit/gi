package compiler

import (
	"fmt"

	"github.com/gijit/gi/pkg/importer"
	"github.com/gijit/gi/pkg/token"
	"github.com/gijit/gi/pkg/types"
	"github.com/glycerine/luar"

	// shadow_ imports: available inside the REPL
	"github.com/gijit/gi/pkg/compiler/shadow/bytes"
	"github.com/gijit/gi/pkg/compiler/shadow/fmt"
	"github.com/gijit/gi/pkg/compiler/shadow/io"
	"github.com/gijit/gi/pkg/compiler/shadow/math"
	shadow_math_rand "github.com/gijit/gi/pkg/compiler/shadow/math/rand"
	"github.com/gijit/gi/pkg/compiler/shadow/os"
	"github.com/gijit/gi/pkg/compiler/shadow/regexp"
	"github.com/gijit/gi/pkg/compiler/shadow/time"

	// actuals
	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/integrate"
	"gonum.org/v1/gonum/lapack"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/unit"
)

func init() {
	a := 1
	b := interface{}(&a)
	_, ok := b.(blas.Complex128)
	_ = ok
}

var _ = fd.Backward
var _ = floats.Add
var _ = graph.Copy
var _ = integrate.Trapezoidal
var _ = lapack.None
var _ = mat.Norm
var _ = optimize.ArmijoConditionMet
var _ = stat.CDF
var _ = unit.Atto

func (ic *IncrState) EnableImportsFromLua() {

	goImportFromLua := func(path string) {
		ic.GiImportFunc(path)
	}
	luar.Register(ic.vm, "", luar.Map{
		"__go_import": goImportFromLua,
	})
}

func (ic *IncrState) GiImportFunc(path string) (*Archive, error) {

	// `import "fmt"` means that path == "fmt", for example.
	pp("GiImportFunc called with path = '%s'", path)

	var pkg *types.Package

	switch path {
	case "gitesting":
		// test only:
		if !ic.vmCfg.NotTestMode {
			fmt.Print("\n registering gitesting.SumArrayInt64! \n")
			pkg = types.NewPackage("gitesting", "gitesting")
			pkg.MarkComplete()
			scope := pkg.Scope()

			fun := getFunForSumArrayInt64(pkg)
			scope.Insert(fun)

			summer := getFunForSummer(pkg)
			scope.Insert(summer)

			summerAny := getFunForSummerAny(pkg)
			scope.Insert(summerAny)

			incr := getFunForIncr(pkg)
			scope.Insert(incr)

			luar.Register(ic.vm, "gitesting", luar.Map{
				"SumArrayInt64": sumArrayInt64,
				//"__giClone":     __giClone,
				"Summer":    Summer,
				"SummerAny": SummerAny,
				"Incr":      Incr,
			})

			ic.CurPkg.importContext.Packages[path] = pkg
			return &Archive{
				ImportPath: path,
				Pkg:        pkg,
			}, nil
		}

		// gen-gijit-shadow outputs to pkg/compiler/shadow/...
	case "bytes":
		luar.Register(ic.vm, "bytes", shadow_bytes.Pkg)
	case "fmt":
		luar.Register(ic.vm, "fmt", shadow_fmt.Pkg)
	case "io":
		luar.Register(ic.vm, "io", shadow_io.Pkg)
	case "math":
		luar.Register(ic.vm, "math", shadow_math.Pkg)
	case "math/rand":
		luar.Register(ic.vm, "rand", shadow_math_rand.Pkg)
	case "os":
		luar.Register(ic.vm, "os", shadow_os.Pkg)
	case "regexp":
		luar.Register(ic.vm, "regexp", shadow_regexp.Pkg)
	case "time":
		luar.Register(ic.vm, "time", shadow_time.Pkg)

	default:
		// need to run gen-gijit-shadow-import
		return nil, fmt.Errorf("erro: package '%s' unknown, or not shadowed. To shadow it, run gen-gijit-shadow-import on the package, add a case and import above, and recompile gijit.", path)
	}

	// loading from real GOROOT/GOPATH.
	// Omit vendor support for now, for sanity.
	shadowPath := "github.com/gijit/gi/pkg/compiler/shadow/" + path
	return ic.ActuallyImportPackage(path, "", shadowPath)
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

// We use the go/importer to load the compiled form of
// the package. This reads from the
// last built binary .a lib on disk. Warning: this might
// be out of date. Later we might read source using the
// go/loader from tools/x, to be most up to date.
// However, the binary loader is *much* faster.
//
// dir provides where to import from, to honor vendored packages.
func (ic *IncrState) ActuallyImportPackage(path, dir, shadowPath string) (*Archive, error) {
	var pkg *types.Package

	//imp := importer.For("source", nil) // Default()
	// faster than source importing is reading the binary.
	imp := importer.Default()
	imp2, ok := imp.(types.ImporterFrom)
	if !ok {
		panic("importer.ImportFrom not available, vendored packages would be lost")
	}
	var mode types.ImportMode
	var err error
	pkg, err = imp2.ImportFrom(path, dir, mode)

	if err != nil {
		return nil, err
	}

	pkgName := pkg.Name()

	res := &Archive{
		Name:       pkgName,
		ImportPath: path,
		Pkg:        pkg,
	}

	pkg.SetPath(shadowPath)

	// very important, must do this or we won't locate the package!
	ic.CurPkg.importContext.Packages[path] = pkg

	return res, nil
}

// __gijit_printQuoted(a ...interface{})
func getFunForGijitPrintQuoted(pkg *types.Package) *types.Func {
	// func __gijit_printQuoted(a ...interface{})
	var recv *types.Var
	var T types.Type = &types.Interface{}
	results := types.NewTuple()
	params := types.NewTuple(types.NewVar(token.NoPos, pkg, "a", types.NewSlice(T)))
	variadic := true
	sig := types.NewSignature(recv, params, results, variadic)
	fun := types.NewFunc(token.NoPos, pkg, "__gijit_printQuoted", sig)
	return fun
}

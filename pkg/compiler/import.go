package compiler

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gijit/gi/pkg/importer"
	"github.com/gijit/gi/pkg/token"
	"github.com/gijit/gi/pkg/types"
	golua "github.com/glycerine/golua/lua"
	"github.com/glycerine/luar"

	// shadow_ imports: available inside the REPL
	"github.com/gijit/gi/pkg/compiler/shadow/bytes"
	"github.com/gijit/gi/pkg/compiler/shadow/fmt"
	"github.com/gijit/gi/pkg/compiler/shadow/io"

	shadow_sync "github.com/gijit/gi/pkg/compiler/shadow/sync"
	shadow_sync_atomic "github.com/gijit/gi/pkg/compiler/shadow/sync/atomic"

	shadow_io_ioutil "github.com/gijit/gi/pkg/compiler/shadow/io/ioutil"
	"io/ioutil"

	"github.com/gijit/gi/pkg/compiler/shadow/math"
	shadow_math_rand "github.com/gijit/gi/pkg/compiler/shadow/math/rand"
	"github.com/gijit/gi/pkg/compiler/shadow/os"
	shadow_reflect "github.com/gijit/gi/pkg/compiler/shadow/reflect"
	"github.com/gijit/gi/pkg/compiler/shadow/regexp"
	shadow_runtime "github.com/gijit/gi/pkg/compiler/shadow/runtime"
	shadow_runtime_debug "github.com/gijit/gi/pkg/compiler/shadow/runtime/debug"
	shadow_strconv "github.com/gijit/gi/pkg/compiler/shadow/strconv"
	"github.com/gijit/gi/pkg/compiler/shadow/time"

	"runtime/debug"

	// gonum
	shadow_blas "github.com/gijit/gi/pkg/compiler/shadow/gonum.org/v1/gonum/blas"
	shadow_fd "github.com/gijit/gi/pkg/compiler/shadow/gonum.org/v1/gonum/diff/fd"
	shadow_floats "github.com/gijit/gi/pkg/compiler/shadow/gonum.org/v1/gonum/floats"
	shadow_graph "github.com/gijit/gi/pkg/compiler/shadow/gonum.org/v1/gonum/graph"
	shadow_integrate "github.com/gijit/gi/pkg/compiler/shadow/gonum.org/v1/gonum/integrate"
	shadow_lapack "github.com/gijit/gi/pkg/compiler/shadow/gonum.org/v1/gonum/lapack"
	shadow_mat "github.com/gijit/gi/pkg/compiler/shadow/gonum.org/v1/gonum/mat"
	shadow_optimize "github.com/gijit/gi/pkg/compiler/shadow/gonum.org/v1/gonum/optimize"
	shadow_stat "github.com/gijit/gi/pkg/compiler/shadow/gonum.org/v1/gonum/stat"
	shadow_unit "github.com/gijit/gi/pkg/compiler/shadow/gonum.org/v1/gonum/unit"

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

var _ = ioutil.Discard
var _ = shadow_blas.GijitShadow_InterfaceConvertTo2_Float64
var _ = fd.Backward
var _ = floats.Add
var _ = graph.Copy
var _ = integrate.Trapezoidal
var _ = lapack.None
var _ = mat.Norm
var _ = optimize.ArmijoConditionMet
var _ = stat.CDF
var _ = unit.Atto

func registerLuarReqs(vm *golua.State) {
	// channel ops need reflect, so import it always.

	luar.Register(vm, "reflect", shadow_reflect.Pkg)
	luar.Register(vm, "fmt", shadow_fmt.Pkg)
	//fmt.Printf("reflect/fmt registered\n")

	// give goroutines.lua something to clone
	// to generate select cases.
	refSelCaseVal := reflect.SelectCase{}

	luar.Register(vm, "", luar.Map{
		"__refSelCaseVal": refSelCaseVal,
	})

	registerBasicReflectTypes(vm)

}

func (ic *IncrState) EnableImportsFromLua() {

	// minimize luar stuff for now, focus on pure Lua runtime.
	goImportFromLua := func(path string) {
		ic.GiImportFunc(path, "")
	}
	stacksClosure := func() {
		showLuaStacks(ic.goro.vm)
	}
	luar.Register(ic.goro.vm, "", luar.Map{
		"__go_import": goImportFromLua,
		"__stacks":    stacksClosure,
	})
}

func (ic *IncrState) GiImportFunc(path, pkgDir string) (*Archive, error) {

	// `import "fmt"` means that path == "fmt", for example.
	pp("GiImportFunc called with path = '%s'... TODO: pure Lua packages. No go/binary/luar based stuff for now\n", path)

	//return nil, nil

	var pkg *types.Package
	t0 := ic.goro.newTicket("", true)

	switch path {
	case "gitesting":
		// test only:
		fmt.Printf("ic.cfg.IsTestMode = %v\n", ic.cfg.IsTestMode)
		if ic.cfg.IsTestMode {
			//fmt.Print("\n registering gitesting.SumArrayInt64! \n")
			pkg = types.NewPackage("gitesting", "gitesting")
			pkg.MarkComplete()
			scope := pkg.Scope()

			suma := getFunForSumArrayInt64(pkg)
			scope.Insert(suma)

			summer := getFunForSummer(pkg)
			scope.Insert(summer)

			summerAny := getFunForSummerAny(pkg)
			scope.Insert(summerAny)

			incr := getFunForIncr(pkg)
			scope.Insert(incr)

			t0.regns = "gitesting"
			t0.regmap["SumArrayInt64"] = sumArrayInt64
			t0.regmap["Summer"] = Summer
			t0.regmap["SummerAny"] = SummerAny
			t0.regmap["Incr"] = Incr
			panicOn(t0.Do())

			ic.CurPkg.importContext.Packages[path] = pkg
			return &Archive{
				SavedArchive: SavedArchive{
					ImportPath: path,
				},
				Pkg: pkg,
			}, nil
		}

		// gen-gijit-shadow outputs to pkg/compiler/shadow/...
	case "bytes":
		t0.regmap["bytes"] = shadow_bytes.Pkg
		t0.regmap["__ctor__bytes"] = shadow_bytes.Ctor
		t0.run = append(t0.run, shadow_bytes.InitLua()...)
	case "fmt":
		t0.regmap["fmt"] = shadow_fmt.Pkg
		t0.regmap["__ctor__fmt"] = shadow_fmt.Ctor
		t0.run = append(t0.run, shadow_fmt.InitLua()...)
	case "io":
		t0.regmap["io"] = shadow_io.Pkg
		t0.regmap["__ctor__io"] = shadow_io.Ctor
		t0.run = append(t0.run, shadow_io.InitLua()...)
	case "math":
		t0.regmap["math"] = shadow_math.Pkg
		t0.regmap["__ctor__math"] = shadow_math.Ctor
		t0.run = append(t0.run, shadow_math.InitLua()...)
	case "math/rand":
		t0.regmap["rand"] = shadow_math_rand.Pkg
		t0.regmap["__ctor__math_rand"] = shadow_math_rand.Ctor
		t0.run = append(t0.run, shadow_math_rand.InitLua()...)
	case "os":
		t0.regmap["os"] = shadow_os.Pkg
		t0.regmap["__ctor__os"] = shadow_os.Ctor
		t0.run = append(t0.run, shadow_os.InitLua()...)

	case "reflect":
		t0.regmap["reflect"] = shadow_reflect.Pkg
		t0.regmap["__ctor__reflect"] = shadow_reflect.Ctor
		t0.run = append(t0.run, shadow_reflect.InitLua()...)

	case "regexp":
		t0.regmap["regexp"] = shadow_regexp.Pkg
		t0.regmap["__ctor__regexp"] = shadow_regexp.Ctor
		t0.run = append(t0.run, shadow_regexp.InitLua()...)

	case "sync":
		t0.regmap["sync"] = shadow_sync.Pkg
		t0.regmap["__ctor__sync"] = shadow_sync.Ctor
		t0.run = append(t0.run, shadow_sync.InitLua()...)

	case "sync/atomic":
		t0.regmap["atomic"] = shadow_sync_atomic.Pkg
		t0.regmap["__ctor__atomic"] = shadow_sync_atomic.Ctor
		t0.run = append(t0.run, shadow_sync_atomic.InitLua()...)

	case "time":
		t0.regmap["time"] = shadow_time.Pkg
		t0.regmap["__ctor__time"] = shadow_time.Ctor
		t0.run = append(t0.run, shadow_time.InitLua()...)

	case "runtime":
		t0.regmap["runtime"] = shadow_runtime.Pkg
		t0.regmap["__ctor__runtime"] = shadow_runtime.Ctor
		t0.run = append(t0.run, shadow_runtime.InitLua()...)

	case "runtime/debug":
		t0.regmap["debug"] = shadow_runtime_debug.Pkg
		t0.regmap["__ctor__debug"] = shadow_runtime_debug.Ctor
		t0.run = append(t0.run, shadow_runtime_debug.InitLua()...)

	case "strconv":
		t0.regmap["strconv"] = shadow_strconv.Pkg
		t0.regmap["__ctor__strconv"] = shadow_strconv.Ctor
		t0.run = append(t0.run, shadow_strconv.InitLua()...)

	case "io/ioutil":
		t0.regmap["ioutil"] = shadow_io_ioutil.Pkg
		t0.regmap["__ctor__ioutil"] = shadow_io_ioutil.Ctor
		t0.run = append(t0.run, shadow_io_ioutil.InitLua()...)

		// gonum:
	case "gonum.org/v1/gonum/blas":
		t0.regmap["blas"] = shadow_blas.Pkg
	case "gonum.org/v1/gonum/fd":
		t0.regmap["fd"] = shadow_fd.Pkg
	case "gonum.org/v1/gonum/floats":
		t0.regmap["floats"] = shadow_floats.Pkg
	case "gonum.org/v1/gonum/graph":
		t0.regmap["graph"] = shadow_graph.Pkg
	case "gonum.org/v1/gonum/integrate":
		t0.regmap["integrate"] = shadow_integrate.Pkg
	case "gonum.org/v1/gonum/lapack":
		t0.regmap["lapack"] = shadow_lapack.Pkg
	case "gonum.org/v1/gonum/mat":
		t0.regmap["mat"] = shadow_mat.Pkg
	case "gonum.org/v1/gonum/optimize":
		t0.regmap["optimize"] = shadow_optimize.Pkg
	case "gonum.org/v1/gonum/stat":
		t0.regmap["stat"] = shadow_stat.Pkg
	case "gonum.org/v1/gonum/unit":
		t0.regmap["unit"] = shadow_unit.Pkg

	default:
		// try a source import.
		archive, err := ic.ImportSourcePackage(path, pkgDir)
		if err == nil {
			if archive == nil {
				panic("why was archive nil if err was nil?")
			}
			if archive.Pkg == nil {
				panic("why was archive.Pkg nil if err was nil?")
			}
			// success. execute the code to define
			// functions in the lua namespace.
			t0.regns = archive.Pkg.Name()
			pp("calling WriteCommandPackage")
			isMain := false
			code, err := ic.CurPkg.Session.WriteCommandPackage(archive, "", isMain)
			pp("back from WriteCommandPackage for path='%s', err='%v', code is\n'%s'", path, err, string(code))
			if err != nil {
				return nil, err
			}
			archive.NewCodeText = [][]byte{code}
			t0.run = code
			panicOn(t0.Do())

			return archive, err
		}
		// source import failed.
		pp("source import of '%s' failed: '%v'", path, err)

		// need to run gen-gijit-shadow-import
		return nil, fmt.Errorf("error on import: problem with package '%s' (not shadowed? [1]): '%v'. ... [footnote 1] To shadow it, run gen-gijit-shadow-import on the package, add a case and import above, and recompile gijit.", path, err)

	}
	panicOn(t0.Do())

	// loading from real GOROOT/GOPATH.
	// Omit vendor support for now, for sanity.
	shadowPath := "github.com/gijit/gi/pkg/compiler/shadow/" + path
	return ic.ActuallyImportPackage(path, "", shadowPath)
}

func omitAnyShadowPathPrefix(path string) string {
	const prefix = "github.com/gijit/gi/pkg/compiler/shadow/"
	if strings.HasPrefix(path, prefix) {
		return path[len(prefix):]
	}
	return path
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

func getFunFor__tostring(pkg *types.Package) *types.Func {
	// func Tostring(a interface{}) string
	var recv *types.Var
	str := types.Typ[types.String]
	results := types.NewTuple(types.NewVar(token.NoPos, pkg, "", str))
	emptyInterface := types.NewInterface(nil, nil)
	params := types.NewTuple(types.NewVar(token.NoPos, pkg, "a", emptyInterface))
	variadic := false
	sig := types.NewSignature(recv, params, results, variadic)
	fun := types.NewFunc(token.NoPos, pkg, "__tostring", sig)
	return fun
}

// make the lua __st (show table) utility available in Go land.
func getFunFor__st(pkg *types.Package) *types.Func {
	// func __st(a interface{}) string
	var recv *types.Var
	str := types.Typ[types.String]
	results := types.NewTuple(types.NewVar(token.NoPos, pkg, "", str))
	emptyInterface := types.NewInterface(nil, nil)
	params := types.NewTuple(types.NewVar(token.NoPos, pkg, "a", emptyInterface))
	variadic := false
	sig := types.NewSignature(recv, params, results, variadic)
	fun := types.NewFunc(token.NoPos, pkg, "__st", sig)
	return fun
}

// and __ls, __gls, __lst, __glst; all functions with no arguments
// and no results, just side effect of displaying info.
func getFunFor__replUtil(cmd string, pkg *types.Package) *types.Func {
	// func __ls()
	sig := types.NewSignature(nil, nil, nil, false)
	fun := types.NewFunc(token.NoPos, pkg, cmd, sig)
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
	pp("IncrState.ActuallyImportPackage(path='%s', dir='%s', shadowPath='%s'", path, dir, shadowPath)
	pp("stack='%s'", string(debug.Stack()))
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
		SavedArchive: SavedArchive{
			Name:       pkgName,
			ImportPath: path,
		},
		Pkg: pkg,
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

/* time stuff

var $setTimeout = function(f, t) {
  $awakeGoroutines++;
  return setTimeout(function() {
    $awakeGoroutines--;
    f();
  }, t);
};

func Sleep(d Duration) {
	c := make(chan struct{})
	js.Global.Call("$setTimeout", js.InternalObject(func() { close(c) }), int(d/Millisecond))
	<-c
}


func startTimer(t *runtimeTimer) {
	t.active = true
	diff := (t.when - runtimeNano()) / int64(Millisecond)
	if diff > 1<<31-1 { // math.MaxInt32
		return
	}
	if diff < 0 {
		diff = 0
	}
	t.timeout = js.Global.Call("$setTimeout", js.InternalObject(func() {
		t.active = false
		if t.period != 0 {
			t.when += t.period
			startTimer(t)
		}
		go t.f(t.arg, 0)
	}), diff+1)
}

func stopTimer(t *runtimeTimer) bool {
	js.Global.Call("clearTimeout", t.timeout)
	wasActive := t.active
	t.active = false
	return wasActive
}

*/

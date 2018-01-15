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

	pack := types.NewPackage("fmt", "fmt")
	pack.MarkComplete()
	scope := pack.Scope()

	fun := getFunForSprintf(pack)

	// As it should, scope.Insert(fun)
	// gets rid of 'Sprintf not declared by package fmt'
	// from types/call.go:302.
	scope.Insert(fun)

	ic.importContext.Packages[path] = pack

	// implementation via luar-based reflection

	// fmt
	luar.Register(ic.vm, "fmt", luar.Map{
		// Go functions may be registered directly.
		"Sprintf": fmt.Sprintf,
	})

	return &Archive{
		ImportPath: path,
		pkg:        pack,
	}, nil
}

func getFunForSprintf(pack *types.Package) *types.Func {
	// func Sprintf(format string, a ...interface{}) string
	var recv *types.Var
	var T types.Type = &types.Interface{}
	str := types.Typ[types.String]
	results := types.NewTuple(types.NewVar(token.NoPos, pack, "", str))
	params := types.NewTuple(types.NewVar(token.NoPos, pack, "format", str),
		types.NewVar(token.NoPos, pack, "a", types.NewSlice(T)))
	variadic := true
	sig := types.NewSignature(recv, params, results, variadic)
	fun := types.NewFunc(token.NoPos, pack, "Sprintf", sig)
	return fun
}

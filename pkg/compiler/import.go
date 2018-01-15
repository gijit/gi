package compiler

import (
	"github.com/gijit/gi/pkg/token"
	"github.com/gijit/gi/pkg/types"
)

func (ic *IncrState) GiImportFunc(path string) (*Archive, error) {

	pp("GiImportFunc called with path = '%s'", path)

	//panic("where import?")
	// gotta return a such that
	// importContext.Packages[a.ImportPath]
	// gives a *types.Package for path.

	// jea: mvp hack?

	pack := types.NewPackage("fmt", "fmt")
	pack.MarkComplete()

	//var parent *types.Scope = ic.archive.pkg.Scope()
	pos := token.NoPos
	//end := token.NoPos
	scope := pack.Scope()

	// func Sprintf(format string, a ...interface{}) string
	var recv *types.Var
	var T types.Type = &types.Interface{}
	str := types.Typ[types.String]
	results := types.NewTuple(types.NewVar(pos, pack, "", str))
	params := types.NewTuple(types.NewVar(pos, pack, "format", str),
		types.NewVar(pos, pack, "a", types.NewSlice(T)))
	variadic := true
	sig := types.NewSignature(recv, params, results, variadic)
	fun := types.NewFunc(pos, pack, "Sprintf", sig)

	scope.Insert(fun)

	// try to get "fmt" to resolve

	// ic.archive is nil on first use
	/*
		if ic.archive == nil {

		} else {
			imps := ic.archive.pkg.Imports()
			uniq := make(map[string]*types.Package)
			for _, im := range imps {
				uniq[im.Path()] = im
			}
			uniq[pack.Path()] = pack // insert or update
			imps = imps[:0]
			for _, im := range uniq {
				imps = append(imps, im)
			}
			ic.archive.pkg.SetImports(imps)
		}
	*/
	//	ic.archive.pkg.Scope().Insert(scope)

	ic.importContext.Packages[path] = pack

	return &Archive{
		ImportPath: path,
		pkg:        pack,
		//Declarations: []*Decl{},
	}, nil
}

package compiler

import (
	"github.com/gijit/gi/pkg/token"
	"github.com/gijit/gi/pkg/types"
)

func (ic *IncrState) MakeGiImportFunc() func(string) (*Archive, error) {

	return func(path string) (*Archive, error) {
		// gotta return a such that
		// importContext.Packages[a.ImportPath]
		// gives a *types.Package for path.

		// jea: mvp hack?

		pack := &types.Package{}
		pack.SetName("fmt")
		pack.SetPath("fmt")
		comment := ""
		var parent *types.Scope
		pos := token.NoPos
		end := token.NoPos
		scope := types.NewScope(parent, pos, end, comment)

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
		ic.importContext.Packages[path] = pack

		return &Archive{ImportPath: path}, nil
	}
}

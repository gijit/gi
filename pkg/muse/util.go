package muse

import (
	"fmt"
	"github.com/gijit/gi/pkg/ast"
	//"github.com/gijit/gi/pkg/parser"
	"github.com/gijit/gi/pkg/token"
	"github.com/gijit/gi/pkg/types"
)

var sizes64 = &types.StdSizes{WordSize: 8, MaxAlign: 8}

func typeCheck(typ ast.Expr, fileSet *token.FileSet, file *ast.File) types.Type {

	file.Name = &ast.Ident{
		Name: "",
	}
	files := []*ast.File{file}

	config := &types.Config{
		DisableUnusedImportCheck: true,
		Sizes: sizes64,
		Error: func(err error) {
			panic(fmt.Sprintf("where error? err = '%v'", err))
		},
	}

	typesInfo := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object), // imports, but those without renames?
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}

	importPath := ""
	pkg, check, err := config.Check(nil, nil, importPath, fileSet, files, typesInfo, nil)
	panicOn(err)

	pp("check: '%#v'", check)
	pp("pkg: '%#v'", pkg)

	return nil
}

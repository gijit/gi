package muse

import (
	"fmt"
	"github.com/gijit/gi/pkg/ast"
	"github.com/gijit/gi/pkg/parser"
	"github.com/gijit/gi/pkg/token"
	"github.com/gijit/gi/pkg/types"
)

func check(typ *ast.Expr, fileSet *ast.FileSet, files []*ast.File) types.Type {

	config := &types.Config{
		DisableUnusedImportCheck: true,
		Sizes: sizes32,
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
	pkg, check, err := config.Check(nil, check, importPath, fileSet, files, typesInfo, nil)
	panicOn(err)

	pp("got past config.Check")
}

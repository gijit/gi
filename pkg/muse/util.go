package muse

import (
	"fmt"
	"github.com/gijit/gi/pkg/ast"
	"github.com/gijit/gi/pkg/compiler/analysis"
	"github.com/gijit/gi/pkg/token"
	"github.com/gijit/gi/pkg/types"
	"github.com/neelance/astrewrite"

	"golang.org/x/tools/go/types/typeutil"
)

var sizes64 = &types.StdSizes{WordSize: 8, MaxAlign: 8}

var reservedKeywords = make(map[string]bool)

func init() {
	for _, keyword := range []string{"abstract", "arguments", "boolean", "break", "byte", "case", "catch", "char", "class", "const", "continue", "debugger", "default", "delete", "do", "double", "else", "enum", "eval", "export", "extends", "false", "final", "finally", "float", "for", "function", "goto", "if", "implements", "import", "in", "instanceof", "int", "interface", "let", "long", "native", "new", "null", "package", "private", "protected", "public", "return", "short", "static", "super", "switch", "synchronized", "this", "throw", "throws", "transient", "true", "try", "typeof", "undefined", "var", "void", "volatile", "while", "with", "yield"} {
		reservedKeywords[keyword] = true
	}
}

type ErrorList []error

type selection interface {
	Kind() types.SelectionKind
	Recv() types.Type
	Index() []int
	Obj() types.Object
	Type() types.Type
}

type pkgContext struct {
	*analysis.Info
	additionalSelections map[*ast.SelectorExpr]selection

	typeNames    []*types.TypeName
	pkgVars      map[string]string
	objectNames  map[types.Object]string
	varPtrNames  map[*types.Var]string
	anonTypes    []*types.TypeName
	anonTypeMap  typeutil.Map
	escapingVars map[*types.Var]bool
	indentation  int
	dependencies map[types.Object]bool
	minify       bool
	fileSet      *token.FileSet
	files        []*ast.File
	errList      ErrorList
}

type flowData struct {
	postStmt  func()
	beginCase int
	endCase   int
}

type funcContext struct {
	*analysis.FuncInfo
	p             *pkgContext
	parent        *funcContext
	sig           *types.Signature
	allVars       map[string]int
	localVars     []string
	resultNames   []ast.Expr
	flowDatas     map[*types.Label]*flowData
	caseCounter   int
	labelCases    map[*types.Label]int
	output        []byte
	delayedOutput []byte
	posAvailable  bool
	pos           token.Pos

	genSymCounter int64

	intType types.Type
}

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

	simplifiedFiles := make([]*ast.File, len(files))
	for i, file := range files {
		simplifiedFiles[i] = astrewrite.Simplify(file, typesInfo, false)
	}

	isBlocking := func(f *types.Func) bool {
		return false
	}

	//pp("about to call AnalyzePkg")
	pkgInfo := analysis.AnalyzePkg(simplifiedFiles, fileSet, typesInfo, pkg, isBlocking)
	c := &funcContext{
		FuncInfo: pkgInfo.InitFuncInfo,
		p: &pkgContext{
			Info:                 pkgInfo,
			additionalSelections: make(map[*ast.SelectorExpr]selection),

			pkgVars:      make(map[string]string),
			objectNames:  make(map[types.Object]string),
			varPtrNames:  make(map[*types.Var]string),
			escapingVars: make(map[*types.Var]bool),
			indentation:  1,
			dependencies: make(map[types.Object]bool),
			minify:       false,
			fileSet:      fileSet,
			files:        files,
		},
		allVars:     make(map[string]int),
		flowDatas:   map[*types.Label]*flowData{nil: {}},
		caseCounter: 1,
		labelCases:  make(map[*types.Label]int),
	}
	for name := range reservedKeywords {
		c.allVars[name] = 1
	}
	pp("got past AnalyzePkg")

	for _, file := range simplifiedFiles {
		for _, decl := range file.Nodes {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				pp("next decl from file.Nodes is a funcDecl: '%#v'", d)

				sig := c.p.Defs[d.Name].(*types.Func).Type().(*types.Signature)
				var recvType types.Type
				if sig.Recv() != nil {
					recvType = sig.Recv().Type()
					if ptr, isPtr := recvType.(*types.Pointer); isPtr {
						recvType = ptr.Elem()
					}
				}
				if sig.Recv() == nil {
					//c.objectName(c.p.Defs[d.Name].(*types.Func)) // register toplevel name
				}
				_ = recvType
			case *ast.GenDecl:
				// jea: could also be *ast.TypeSpec here when declaring a struct!
				switch ds := d.Specs[0].(type) {
				case *ast.TypeSpec:
					pp("next decl from file.Nodes is a *ast.TypeSpec: '%#v'", ds)
				case *ast.ValueSpec:
					pp("next decl from file.Nodes is a *ast.VaueSpec: '%#v'", ds.Names[0])
				default:
					pp("next decl from file.Nodes is an unrecognized *ast.GenDecl: '%#v'", ds)
				}
				switch d.Tok {
				case token.TYPE:
					pp("we're in the token.TYPE!")
					for _, spec := range d.Specs {
						obj := c.p.Defs[spec.(*ast.TypeSpec).Name]
						o := obj.(*types.TypeName)
						pp("o='%#v'", o)
						return obj.Type()
					}
				case token.VAR:
					for _, spec := range d.Specs {
						for _, name := range spec.(*ast.ValueSpec).Names {
							_ = name
						}
					}
				case token.CONST:

				}
			default:
				pp("next decl from file.Nodes is an unknown/default type: '%#v'", decl)
				//_, isExprStmt := decl.(*ast.ExprStmt)

			}
		}
	}

	return nil
}

package compiler

import (
	"bytes"
	"fmt"
	"github.com/gijit/gi/pkg/ast"
	"github.com/gijit/gi/pkg/gostd/build"
	"github.com/gijit/gi/pkg/parser"
	"github.com/gijit/gi/pkg/token"
	"github.com/gijit/gi/pkg/types"
	"github.com/glycerine/zygomys/zygo"
	//"github.com/gijit/gi/pkg/verb"
	"unicode"
	//luajit "github.com/glycerine/golua/lua"
	//"github.com/kisielk/gotool"
	//"github.com/shurcooL/go-goon"
	//gbuild "github.com/gijit/gi/pkg/gostd/build"
	//gibuild "github.com/gijit/gi/pkg/gibuild"
)

// IncrState holds the incremental translation
// (Golang-source-to-Lua-souce compilation) state.
type IncrState struct {

	// allow multiple packages to
	// be worked on at once.
	pkgMap map[string]*IncrPkg

	CurPkg *IncrPkg

	// the vm lets us add import bindings
	// like `import "fmt"` on demand.
	// Update: But this is now per-goroutine,
	// needing syncrhonization, so front end doesn't get to touch.
	// vm *luajit.State
	goro *Goro

	cfg *GIConfig

	minify   bool
	PrintAST bool

	// default to no import caching
	//AllowImportCaching bool

	Session *Session

	zlisp *zygo.Zlisp
}

func NewIncrState(lvm *LuaVm, cfg *GIConfig) *IncrState {

	if lvm == nil {
		panic("NewIncrState(): lvm cannot be nil")
	}

	if cfg == nil {
		cfg = NewGIConfig()
		lvm.cfg = cfg
	}
	ic := &IncrState{
		goro:   lvm.goro,
		pkgMap: make(map[string]*IncrPkg),
		cfg:    cfg,
	}
	ic.Session = NewSession(&Options{}, ic)

	pack := &build.Package{
		Name:       "main",
		ImportPath: "main",
		Dir:        ".",
	}
	fileSet := token.NewFileSet() // positions are relative to fileSet
	importContext := &ImportContext{
		Packages: make(map[string]*types.Package),
		Import:   ic.CompileTimeGiImportFunc,
	}

	key := "main"
	pk := newIncrPkg(key, pack, fileSet, importContext, nil, ic)

	ic.pkgMap[key] = pk
	ic.CurPkg = pk

	ic.EnableImportsFromLua() // from Lua, use __go_import("fmt");

	return ic
}

// an incrementally built package,
// stored in IncrState.pkgMap
//
type IncrPkg struct {
	key string

	pack          *build.Package
	fileSet       *token.FileSet
	importContext *ImportContext // has Packages map[string]*types.Package
	Arch          *Archive

	// use ic.Session.Archives instead.
	//localImportPathCache map[string]*Archive

	ic *IncrState
}

func newIncrPkg(key string,
	pack *build.Package,
	fileSet *token.FileSet,
	importContext *ImportContext,
	archive *Archive,
	ic *IncrState,
) *IncrPkg {

	return &IncrPkg{
		key:           key,
		pack:          pack,
		fileSet:       fileSet,
		importContext: importContext,
		Arch:          archive,

		// use ic.Session.Archives instead.
		//localImportPathCache: make(map[string]*Archive),

		ic: ic,
	}
}

type UniqPkgPath string

func (tr *IncrState) Close() {
	tr.goro.halt.RequestStop()
	<-tr.goro.halt.Done.Chan
	tr.goro.vm.Close()
}

//  panic on errors, a test helper
func (tr *IncrState) trMust(src []byte) []byte {
	by, err := tr.Tr(src)
	panicOn(err)
	return by
}

func (tr *IncrState) trMustPre(src []byte, pre bool) []byte {
	by, err := tr.TrWithPrepend(src, pre)
	panicOn(err)
	return by
}

// Tr: translate from go to Lua, statement by statement or
// expression by expression
func (tr *IncrState) Tr(src []byte) (by []byte, err error) {
	return tr.TrWithPrepend(src, true)
}

func (tr *IncrState) TrWithPrepend(src []byte, prependOk bool) (by []byte, err error) {

	defer func() {
		r := recover()
		if r != nil {
			pp("Tr sees panic of: '%v'", r)
			err = fmt.Errorf("%v", r)
		}
	}()

	// detect the leading '=' and turn it into
	// __gijit_ans :=
	var didPrepend, shouldPrepend bool
	if prependOk {
		src, didPrepend = tr.prependAns(src)
		pp("after prependAns, src = '%s'", src)
	}

	// classic
	file, err := parser.ParseFile(tr.CurPkg.fileSet, "", src, 0)
	if err != nil {
		pp("we got an error on the ParseFile: '%v'", err)
	}
	panicOn(err)

	// expression 1st thing? re-parse with ans prepended.
	if prependOk && !didPrepend {
		if len(file.IsExpr) > 0 && file.IsExpr[0] {
			shouldPrepend = true
		}

		// unfortunately for simple function calls, because
		// they could return 2 or more values, we can't just
		// blindly wrap them assuming that the return only 1.
		if len(file.IsStmt) == 1 && file.IsStmt[0] {
			shouldPrepend = true
		}

		if shouldPrepend {
			pp("should prepend is true, trying parse with ans prepend")
			prev := tr.cfg.CalculatorMode
			tr.cfg.CalculatorMode = true // force pre-pending of "ans = "
			src, didPrepend = tr.prependAns(src)
			tr.cfg.CalculatorMode = prev

			file2, err := parser.ParseFile(tr.CurPkg.fileSet, "", src, 0)
			if err == nil {
				file = file2
			} // else we leave file as in, since it parsed without the prepend..
			pp("after shouldPrepend, ParseFile gave err = '%v'; src='%s'", err, src)
		}
	}
	pp("we got past the ParseFile !")

	if tr.PrintAST {
		ast.Print(tr.CurPkg.fileSet, file)
	}

	files := []*ast.File{file}
	pp("file='%#v'", file)
	pp("file.Name='%#v'", file.Name)
	file.Name = &ast.Ident{
		Name: "", // jea: was "/repl", but that seemed to cause scope issues.
	}

	hasBadId, whichBad := checkAllowedIdents(file)
	if hasBadId {
		msg := fmt.Sprintf("bad identifier: cannot "+
			"use '%s' as an identifier in gijit, as this may confuse the online type checker.",
			whichBad)
		panic(msg)
		return nil, fmt.Errorf(msg)
	}
	depth := 0
	tr.CurPkg.Arch, err = IncrementallyCompile(tr.CurPkg.Arch, tr.CurPkg.pack.ImportPath, files, tr.CurPkg.fileSet, tr.CurPkg.importContext, tr.minify, depth)
	panicOn(err)
	//pp("archive = '%#v'", tr.CurPkg.Arch)
	//pp("len(tr.CurPkg.Arch.Declarations)= '%v'", len(tr.CurPkg.Arch.Declarations))
	//pp("len(tr.CurPkg.Arch.NewCode)= '%v'", len(tr.CurPkg.Arch.NewCodeText))

	pp("got past config.Check")

	var res bytes.Buffer
	for i, d := range tr.CurPkg.Arch.NewCodeText {
		pp("writing tr.CurPkg.Arch.NewCode[i=%v].Code = '%v'", i, string(tr.CurPkg.Arch.NewCodeText[i]))
		res.Write(d)
	}
	tr.CurPkg.Arch.NewCodeText = nil

	return res.Bytes(), nil
}

var gijitAnsPrefix = []byte("__gijit_ans := []interface{}{")
var gijitAnsSuffix = []byte("}\n __gijit_printQuoted(__gijit_ans...);")

// at the beginning of src, transform a first '='[^=] into
// "__gijit_ans := "
func (tr *IncrState) prependAns(src []byte) (ret []byte, didPrepend bool) {
	//fmt.Printf("prependAns starting with '%s'\n", string(src))
	//defer func() {
	//	fmt.Printf("prependAns returning ret='%s'\n", string(ret))
	//}()
	nsrc := len(src)
	leftTrimmed := bytes.TrimLeftFunc(src, unicode.IsSpace)
	trimmed := bytes.TrimFunc(src, unicode.IsSpace)
	n := len(leftTrimmed)
	leftdiff := nsrc - n
	if tr.cfg.CalculatorMode {
		middle := removeTrailingSemicolon(trimmed[leftdiff:])
		return append(gijitAnsPrefix, append(middle, gijitAnsSuffix...)...), true
	}
	if n > 1 && leftTrimmed[0] == '=' && leftTrimmed[1] != '=' {
		middle := removeTrailingSemicolon(trimmed[leftdiff+1:])
		return append(gijitAnsPrefix, append(middle, gijitAnsSuffix...)...), true
	}
	return src, false
}

// full package

// FullPackage: translate a full package from go to Lua.
func (tr *IncrState) FullPackage(src []byte, importPath string, depth int) ([]byte, error) {
	pp("FullPackage top.")

	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "", src, 0)
	if err != nil {
		pp("we got an error on the ParseFile: '%v'", err)
	}
	panicOn(err)
	pp("we got past the ParseFile !")

	files := []*ast.File{file}
	file.Name = &ast.Ident{
		Name: "main",
	}

	//tr.CurPkg.Arch,
	arch, err := FullPackageCompile(importPath, files, fileSet, tr.CurPkg.importContext, tr.minify, depth)
	panicOn(err)
	arch.ImportPath = "main"
	//pp("archive = '%#v'", tr.CurPkg.Arch)
	//pp("len(tr.CurPkg.Arch.Declarations)= '%v'", len(tr.CurPkg.Arch.Declarations))
	//pp("len(tr.CurPkg.Arch.NewCode)= '%v'", len(tr.CurPkg.Arch.NewCodeText))

	pp("got past FullPackageCompile")

	var res bytes.Buffer
	w := &SourceMapFilter{
		Writer: &res,
	}
	isMain := true
	err = WriteProgramCode([]*Archive{arch}, w, isMain)

	return res.Bytes(), err
}

var semi = []byte{';'}

func removeTrailingSemicolon(src []byte) []byte {
	tr := bytes.TrimSpace(src)
	for bytes.HasSuffix(tr, semi) {
		tr = tr[:len(tr)-1]
	}
	return tr
}

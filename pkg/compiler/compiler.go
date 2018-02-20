package compiler

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/gijit/gi/pkg/token"
	"github.com/gijit/gi/pkg/types"
	"io"
	"strings"

	prelude "github.com/gijit/gi/pkg/compiler/prelude_lua"
	gcimporter "golang.org/x/tools/go/gcimporter15"
)

//jea comment sizes32 in favor of sizes64
//var sizes32 = &types.StdSizes{WordSize: 4, MaxAlign: 8}
var sizes64 = &types.StdSizes{WordSize: 8, MaxAlign: 8}
var reservedKeywords = make(map[string]bool)
var predeclared = make(map[string]bool)

func init() {
	// javascript reserved words
	for _, w := range []string{"abstract", "arguments", "boolean", "break", "byte", "case", "catch", "char", "class", "const", "continue", "debugger", "default", "delete", "do", "double", "else", "enum", "eval", "export", "extends", "false", "final", "finally", "float", "for", "function", "goto", "if", "implements", "import", "in", "instanceof", "int", "interface", "let", "long", "native", "new", "null", "package", "private", "protected", "public", "return", "short", "static", "super", "switch", "synchronized", "this", "throw", "throws", "transient", "true", "try", "typeof", "undefined", "var", "void", "volatile", "while", "with", "yield"} {
		reservedKeywords[w] = true
	}

	// predeclared numeric types
	for _, w := range []string{"int", "uint", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64", "complex64", "complex128", "byte", "rune", "uintptr"} {
		reservedKeywords[w] = true
		predeclared[w] = true
	}

	// lua reserved words
	for _, w := range []string{"and", "break", "do", "else", "elseif", "", "end", "false", "for", "function", "if", "in", "local", "nil", "not", "or", "repeat", "return", "then", "true", "until", "while"} {
		reservedKeywords[w] = true
	}
}

type ErrorList []error

func (err ErrorList) Error() string {
	return err[0].Error()
}

type Archive struct {
	ImportPath   string
	Name         string
	Imports      []string
	ExportData   []byte
	Declarations []*Decl
	IncJSCode    []byte
	FileSet      []byte
	Minified     bool
	NewCodeText  [][]byte

	// save state so we can type incrementally
	TypesInfo *types.Info
	Config    *types.Config
	Pkg       *types.Package
	Check     *types.Checker

	FuncSrcCache map[string]string
}

type Decl struct {
	FullName        string
	Vars            []string
	DeclCode        []byte
	MethodListCode  []byte
	TypeInitCode    []byte
	InitCode        []byte
	DceObjectFilter string
	DceMethodFilter string
	DceDeps         []string
	Blocking        bool
}

type Stmt struct {
	Code []byte
}

type Dependency struct {
	Pkg    string
	Type   string
	Method string
}

func ImportDependencies(archive *Archive, importPkg func(string) (*Archive, error)) ([]*Archive, error) {
	var deps []*Archive
	paths := make(map[string]bool)
	var collectDependencies func(path string) error
	collectDependencies = func(path string) error {
		if paths[path] {
			return nil
		}
		dep, err := importPkg(path)
		if err != nil {
			return err
		}
		for _, imp := range dep.Imports {
			if err := collectDependencies(imp); err != nil {
				return err
			}
		}
		deps = append(deps, dep)
		paths[dep.ImportPath] = true
		return nil
	}

	if err := collectDependencies("runtime"); err != nil {
		return nil, err
	}
	for _, imp := range archive.Imports {
		if err := collectDependencies(imp); err != nil {
			return nil, err
		}
	}

	deps = append(deps, archive)
	return deps, nil
}

type dceInfo struct {
	decl         *Decl
	objectFilter string
	methodFilter string
}

func WriteProgramCode(pkgs []*Archive, w *SourceMapFilter) error {
	mainPkg := pkgs[len(pkgs)-1]
	minify := mainPkg.Minified

	byFilter := make(map[string][]*dceInfo)
	var pendingDecls []*Decl
	for _, pkg := range pkgs {
		for _, d := range pkg.Declarations {
			if d.DceObjectFilter == "" && d.DceMethodFilter == "" {
				pendingDecls = append(pendingDecls, d)
				continue
			}
			info := &dceInfo{decl: d}
			if d.DceObjectFilter != "" {
				info.objectFilter = pkg.ImportPath + "." + d.DceObjectFilter
				byFilter[info.objectFilter] = append(byFilter[info.objectFilter], info)
			}
			if d.DceMethodFilter != "" {
				info.methodFilter = pkg.ImportPath + "." + d.DceMethodFilter
				byFilter[info.methodFilter] = append(byFilter[info.methodFilter], info)
			}
		}
	}

	dceSelection := make(map[*Decl]struct{})
	for len(pendingDecls) != 0 {
		d := pendingDecls[len(pendingDecls)-1]
		pendingDecls = pendingDecls[:len(pendingDecls)-1]

		dceSelection[d] = struct{}{}

		for _, dep := range d.DceDeps {
			if infos, ok := byFilter[dep]; ok {
				delete(byFilter, dep)
				for _, info := range infos {
					if info.objectFilter == dep {
						info.objectFilter = ""
					}
					if info.methodFilter == dep {
						info.methodFilter = ""
					}
					if info.objectFilter == "" && info.methodFilter == "" {
						pendingDecls = append(pendingDecls, info.decl)
					}
				}
			}
		}
	}

	if _, err := w.Write([]byte("\n(function()\n\n")); err != nil {
		return err
	}
	usePrecompiledPrelude := false // quickly out of date.
	if usePrecompiledPrelude {
		if _, err := w.Write(removeWhitespace([]byte(prelude.Prelude), minify)); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte("\n")); err != nil {
		return err
	}

	// write packages
	for _, pkg := range pkgs {
		if err := WritePkgCode(pkg, dceSelection, minify, w); err != nil {
			return err
		}
	}

	// leave off the   __go(__mainPkg._init, {});
	if _, err := w.Write([]byte(`
  __synthesizeMethods();

  local __mainPkg = __packages["` + string(mainPkg.ImportPath) + `"];
  local __rtTemp = __packages["runtime"]; 
  if __rtTemp ~= nil then _init(); end;

  __mainPkg._init()
end)();
`)); err != nil {
		return err
	}

	return nil
}

func WritePkgCode(pkg *Archive, dceSelection map[*Decl]struct{}, minify bool, w *SourceMapFilter) error {
	if w.MappingCallback != nil && pkg.FileSet != nil {
		w.fileSet = token.NewFileSet()
		if err := w.fileSet.Read(json.NewDecoder(bytes.NewReader(pkg.FileSet)).Decode); err != nil {
			panic(err)
		}
	}
	if _, err := w.Write(pkg.IncJSCode); err != nil {
		return err
	}
	if _, err := w.Write(removeWhitespace([]byte(fmt.Sprintf("__packages[\"%s\"] = (function()\n", pkg.ImportPath)), minify)); err != nil {
		return err
	}
	vars := []string{"_pkg = {}", "_init"}
	var filteredDecls []*Decl
	for _, d := range pkg.Declarations {
		if _, ok := dceSelection[d]; ok {
			vars = append(vars, d.Vars...)
			filteredDecls = append(filteredDecls, d)
		}
	}
	if _, err := w.Write(removeWhitespace([]byte(fmt.Sprintf("\tlocal %s;\n", strings.Join(vars, ", "))), minify)); err != nil {
		return err
	}
	for _, d := range filteredDecls {
		if _, err := w.Write(d.DeclCode); err != nil {
			return err
		}
	}
	for _, d := range filteredDecls {
		if _, err := w.Write(d.MethodListCode); err != nil {
			return err
		}
	}
	for _, d := range filteredDecls {
		if _, err := w.Write(d.TypeInitCode); err != nil {
			return err
		}
	}

	_, err := w.Write(removeWhitespace([]byte(`
   _init = function(self)
   _pkg._init = function() end;
    -- jea compiler.go:245
    local __f
    local __c = false;
    local __s = 0;
    local __r;
    if self ~= nil and self.__blk ~= nil then
         __f = self;
         __c = true;
         __s = __f.__s;
         __r = __f.__r;
    end;
     ::s::
    while (true) do
         --switch (__s)
    if __s == 0 then
`), minify))

	//_, err := w.Write(removeWhitespace([]byte("\t_init = function(self)\n\t\t _pkg._init = function() end;\n\t\t-- jea compiler.go:245\n\t\t local __f\n\t\t local __c = false;\n\t\t local __s = 0;\n\t\t local __r;\n\t\t if self ~= nil and self.__blk ~= nil then\n\t\t  __f = self;\n\t\t __c = true;\n\t\t __s = __f.__s;\n\t\t __r = __f.__r;\n\t\t end;\n\t\t ::s::\n\t\t while (true) do\n  --switch (__s)\n if __s == 0 then\n"), minify))
	if err != nil {
		return err
	}
	for _, d := range filteredDecls {
		if _, err := w.Write(d.InitCode); err != nil {
			return err
		}
	}
	if _, err := w.Write(removeWhitespace([]byte("\t\t--jea compiler.go:238\n end;\n\t\t return;\n\t end;\n\t\t if __f == nil then\n\t __f = { __blk= _init };\n\t end;\n\t  __f.__s = __s;\n\t __f.__r = __r;\n\t return __f;\n\t end;\n\t__pkg._init = _init;\n\t return __pkg;\n end)();\n"), minify)); err != nil {
		return err
	}
	if _, err := w.Write([]byte("\n")); err != nil { // keep this \n even when minified
		return err
	}
	return nil
}

func ReadArchive(filename, path string, r io.Reader, packages map[string]*types.Package) (*Archive, error) {
	var a Archive
	if err := gob.NewDecoder(r).Decode(&a); err != nil {
		return nil, err
	}

	var err error
	_, packages[path], err = gcimporter.BImportData(token.NewFileSet(), packages, a.ExportData, path)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func WriteArchive(a *Archive, w io.Writer) error {
	return gob.NewEncoder(w).Encode(a)
}

type SourceMapFilter struct {
	Writer          io.Writer
	MappingCallback func(generatedLine, generatedColumn int, originalPos token.Position)
	line            int
	column          int
	fileSet         *token.FileSet
}

func (f *SourceMapFilter) Write(p []byte) (n int, err error) {
	var n2 int
	for {
		i := bytes.IndexByte(p, '\b')
		w := p
		if i != -1 {
			w = p[:i]
		}

		n2, err = f.Writer.Write(w)
		n += n2
		for {
			i := bytes.IndexByte(w, '\n')
			if i == -1 {
				f.column += len(w)
				break
			}
			f.line++
			f.column = 0
			w = w[i+1:]
		}

		if err != nil || i == -1 {
			return
		}
		if f.MappingCallback != nil {
			f.MappingCallback(f.line+1, f.column, f.fileSet.Position(token.Pos(binary.BigEndian.Uint32(p[i+1:i+5]))))
		}
		p = p[i+5:]
		n += 5
	}
}

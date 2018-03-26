package compiler

import (
	"bytes"
	"encoding/binary"
	"strings"

	"encoding/gob"
	"github.com/glycerine/gi/pkg/ast"

	"encoding/json"
	"fmt"
	"github.com/glycerine/gi/pkg/token"
	"github.com/glycerine/gi/pkg/types"
	"io"

	prelude "github.com/glycerine/gi/pkg/compiler/prelude_lua"
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

	// lua/gijit system already in-use method names
	for _, w := range []string{"_G", "_VERSION", "assert", "bit", "byte", "cmath", "collectgarbage", "complex", "complex128", "complex64", "coroutine", "debug", "dofile", "error", "float32", "float64", "gcinfo", "getfenv", "getmetatable", "golua_default_msghandler", "imag", "int", "int16", "int32", "int64", "int8", "io", "ipairs", "jit", "load", "loadfile", "loadstring", "luar", "math", "module", "newproxy", "next", "os", "package", "pairs", "panic", "pcall", "print", "rawequal", "rawget", "rawlen", "rawset", "real", "recover", "require", "select", "setfenv", "setmetatable", "string", "table", "tonumber", "tostring", "type", "uint", "uint16", "uint32", "uint64", "uint8", "unpack", "xpcall"} {
		reservedKeywords[w] = true
	}

	// register for gob
	var iden ast.Ident
	gob.Register(iden)
	var fld ast.Field
	gob.Register(fld)

}

type ErrorList []error

func (err ErrorList) Error() string {
	return err[0].Error()
}

type SavedArchive struct {
	ImportPath   string
	Name         string
	Imports      []string
	ExportData   []byte
	Declarations []*Decl
	IncJSCode    []byte
	FileSet      []byte
	Minified     bool
}

type Archive struct {
	SavedArchive

	// above from GopherJS, below added for gijit.

	NewCodeText [][]byte

	// save state so we can type incrementally
	Pkg       *types.Package
	TypesInfo *types.Info
	Config    *types.Config
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

func ImportDependencies(archive *Archive, importPkg func(pth string, depth int) (*Archive, error), depth int) ([]*Archive, error) {
	var deps []*Archive
	paths := make(map[string]bool)
	var collectDependencies func(path string, depth int) error
	collectDependencies = func(path string, depth int) error {
		if paths[path] {
			return nil
		}
		dep, err := importPkg(path, depth)
		if err != nil {
			return err
		}
		for _, imp := range dep.Imports {
			if err := collectDependencies(imp, depth+1); err != nil {
				return err
			}
		}
		deps = append(deps, dep)
		paths[dep.ImportPath] = true
		return nil
	}
	// jea: temp disable to see if 1000 passes: yes, this makes 1000 green. commenting it in goes red.
	// but 1002 still reports red.
	// /Users/jaten/go/src/github.com/glycerine/gi/pkg/luaapi/luaapi.go:91:6: Error should have been declared
	//
	// with collectDependencies("runtime") in place, all 1000, 1001, and 1002 tests fail
	//  with the "Error should have been declared" error above.
	//
	// with it commented out, tests 1000 and 1001 go green, but 1002 stays red.
	//
	/*
		if err := collectDependencies("runtime"); err != nil {
			return nil, err
		}
	*/
	for _, imp := range archive.Imports {
		if err := collectDependencies(imp, depth+1); err != nil {
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

func WriteProgramCode(pkgs []*Archive, w *SourceMapFilter, isMain bool) error {
	// jea: notice this assumption that main is the last package...
	// This was right for GopherJS, but is now probably wrong in gijit.
	mainPkg := pkgs[len(pkgs)-1]
	minify := mainPkg.Minified

	dceSelection := make(map[*Decl]struct{})

	// jea: debug try back in. TODO fix/revert if need be.
	if false { // len(pkgs) == 1 && !isMain {
		// skip this filtering stuff

	} else {

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

					pp("appending to byFilter[info.objectFiler='%#v'] the following info='%#v'", info.objectFilter, info)
				}
				pp("d.DceMethodFilter is '%s'", d.DceMethodFilter)
				pp("d.DceObjectFilter is '%s'", d.DceObjectFilter)

				if d.DceMethodFilter != "" {

					info.methodFilter = pkg.ImportPath + "." + d.DceMethodFilter
					byFilter[info.methodFilter] = append(byFilter[info.methodFilter], info)
				}
			}
		}

		pp("len(pendingDecls) is '%v'", len(pendingDecls)) // spkg_tst, 0 here, so nothing being allowed through.
		for len(pendingDecls) != 0 {
			d := pendingDecls[len(pendingDecls)-1]
			pendingDecls = pendingDecls[:len(pendingDecls)-1]

			pp("adding pendingDecls d to nil dceSelection map: d.='%#v'", d)
			dceSelection[d] = struct{}{}

			for _, dep := range d.DceDeps {
				pp("considering d.DceDeps, dep='%#v'", dep)
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
	} //end filtering stuff

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

		if !isMain {

			// not Main: need to write the package to the global Lua env.

			// avoid name conflict between types and package name by assigning
			// directly to _G (types may have prior 'local' declarations?)
			_, err := w.Write([]byte(fmt.Sprintf("\n _G.%s = __packages[\"%s\"];\n",
				pkg.Name, string(pkg.ImportPath))))

			if err != nil {
				return err
			}
		}
	}

	_, err := w.Write([]byte(`
  __synthesizeMethods();
`))
	if err != nil {
		return err
	}

	if isMain {
		_, err := w.Write([]byte(`

  local __mainPkg = __packages["` + string(mainPkg.ImportPath) + `"];
  local __rtTemp = __packages["runtime"]; 
  if __rtTemp ~= nil and __rtTemp.__init ~= nil then 
     __rtTemp.__init(); 
  end;

  __go(__mainPkg.__init, {});
`))

		if err != nil {
			return err
		}
	}
	_, err = w.Write([]byte(`
end)();
`))
	if err != nil {
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
	vars := []string{"__pkg = {}", "__init"}
	var filteredDecls []*Decl
	for _, d := range pkg.Declarations {
		// jea TODO figure out why this dceSelection[d] was over filtering
		// our spkg_tst.Fish method and it never showed up in the output.

		// jea: gotta not mix our types into our variables...
		//pp("d.Vars is '%#v'; dceSelection[d]='%v'", d.Vars, dceSelection[d])
		if true { // _, ok := dceSelection[d]; ok {

			// jea: hack, exclude those with '.', since they won't compile anyway...
			for _, v := range d.Vars {
				if !strings.Contains(v, ".") {
					vars = append(vars, v)
				}
			}
			filteredDecls = append(filteredDecls, d)
		} else {
			pp("rejected, filtered out Declaration d='%#v'  \n   because dceSelection[d]='%#v'", d, dceSelection[d])
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
   __init = function(self)
   __pkg.__init = function() end;
    -- jea compiler.go:260
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

	if err != nil {
		return err
	}
	for _, d := range filteredDecls {
		if _, err := w.Write(d.InitCode); err != nil {
			return err
		}
	}
	if _, err := w.Write(removeWhitespace([]byte("\t\t--jea compiler.go:238\n end;\n\t\t return;\n\t end;\n\t\t if __f == nil then\n\t __f = { __blk= __init };\n\t end;\n\t  __f.__s = __s;\n\t __f.__r = __r;\n\t return __f;\n\t end;\n\t__pkg.__init = __init;\n\t return __pkg;\n end)();\n"), minify)); err != nil {
		return err
	}
	if _, err := w.Write([]byte("\n")); err != nil { // keep this \n even when minified
		return err
	}
	return nil
}

func ReadArchive(filename, path string, r io.Reader, packages map[string]*types.Package) (*Archive, error) {
	var a SavedArchive
	if err := gob.NewDecoder(r).Decode(&a); err != nil {
		return nil, err
	}

	var err error
	_, packages[path], err = gcimporter.BImportData(token.NewFileSet(), packages, a.ExportData, path)
	if err != nil {
		return nil, err
	}

	return &Archive{SavedArchive: a}, nil
}

func WriteArchive(a *Archive, w io.Writer) error {
	return gob.NewEncoder(w).Encode(a.SavedArchive)
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

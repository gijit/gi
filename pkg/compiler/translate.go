package compiler

import (
	"bytes"
	"fmt"
	"github.com/gijit/gi/pkg/ast"
	"github.com/gijit/gi/pkg/gostd/build"
	"github.com/gijit/gi/pkg/parser"
	"github.com/gijit/gi/pkg/token"
	"github.com/gijit/gi/pkg/types"
	//"github.com/gijit/gi/pkg/verb"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/fsnotify/fsnotify"
	luajit "github.com/glycerine/golua/lua"
	//"github.com/shurcooL/go-goon"
	//gbuild "github.com/gijit/gi/pkg/gostd/build"
)

// the incremental translation state
func NewIncrState(vm *luajit.State, vmCfg *VmConfig) *IncrState {

	if vmCfg == nil {
		vmCfg = NewVmConfig()
	}
	//	options := &gbuild.Options{CreateMapFile: true}
	//	s := gbuild.NewSession(options)

	ic := &IncrState{
		vm:    vm,
		vmCfg: vmCfg,
		pack: &build.Package{
			Name:       "main",
			ImportPath: "main",
			Dir:        ".",
		},
		fileSet: token.NewFileSet(), // positions are relative to fileSet
		importContext: &ImportContext{
			Packages: make(map[string]*types.Package),

			// from GopherJS:
			/*
				Import: func(path string) (*Archive, error) {
					if path == pkg.ImportPath || path == pkg.ImportPath+"_test" {
						return s.Archives[path], nil
					}
					return s.BuildImportPath(path)
				},
			*/
		},
	}

	ic.importContext.Import = ic.GiImportFunc

	return ic
}

type IncrState struct {
	pack          *build.Package
	fileSet       *token.FileSet
	archive       *Archive
	importContext *ImportContext

	// the vm lets us add import bindings
	// like `import "fmt"` on demand.
	vm *luajit.State

	vmCfg *VmConfig

	minify   bool
	PrintAST bool
}

// Tr: translate from go to javascript, statement by statement or
// expression by expression
func (tr *IncrState) Tr(src []byte) []byte {

	// detect the leading '=' and turn it into
	// __gijit_ans :=
	src = prependAns(src)

	pp("after prependAns, src = '%s'", src)

	// classic
	file, err := parser.ParseFile(tr.fileSet, "", src, 0)
	if err != nil {
		pp("we got an error on the ParseFile: '%v'", err)
	}
	panicOn(err)
	pp("we got past the ParseFile !")

	// Print the AST.
	if tr.PrintAST {
		ast.Print(tr.fileSet, file)
	}

	importPath := "" // was "/repl" but scope issues?
	files := []*ast.File{file}
	pp("file='%#v'", file)
	pp("file.Name='%#v'", file.Name)
	file.Name = &ast.Ident{
		Name: "", // jea: was "/repl", but that seemed to cause scope issues.
	}

	hasBadId, whichBad := checkAllowedIdents(file)
	if hasBadId {
		msg := fmt.Sprintf("bad identifier: cannot "+
			"use '%s' as an identifier in gijit.\n",
			whichBad)
		fmt.Printf(msg)
		panic(msg)
		return nil
	}

	tr.archive, err = IncrementallyCompile(tr.archive, importPath, files, tr.fileSet, tr.importContext, tr.minify)
	panicOn(err)
	pp("archive = '%#v'", tr.archive)
	//pp("len(tr.archive.Declarations)= '%v'", len(tr.archive.Declarations))
	//pp("len(tr.archive.NewCode)= '%v'", len(tr.archive.NewCodeText))

	pp("got past config.Check")

	var res bytes.Buffer
	for i, d := range tr.archive.NewCodeText {
		pp("writing tr.archive.NewCode[i=%v].Code = '%v'", i, string(tr.archive.NewCodeText[i]))
		res.Write(d)
	}
	tr.archive.NewCodeText = nil

	return res.Bytes()
}

type ImportCError struct {
	pkgPath string
}

func (e *ImportCError) Error() string {
	return e.pkgPath + `: importing "C" is not supported by GopherJS`
}

func NewReplContext(installSuffix string, buildTags []string) *build.Context {
	return &build.Context{
		GOROOT:        build.Default.GOROOT,
		GOPATH:        build.Default.GOPATH,
		GOOS:          build.Default.GOOS,
		GOARCH:        "js",
		InstallSuffix: installSuffix,
		Compiler:      "gc",
		BuildTags:     append(buildTags, "netgo"),
		ReleaseTags:   build.Default.ReleaseTags,
		CgoEnabled:    true, // detect `import "C"` to throw proper error
	}
}

// Import returns details about the Go package named by the import path. If the
// path is a local import path naming a package that can be imported using
// a standard import path, the returned package will set p.ImportPath to
// that path.
//
// In the directory containing the package, .go and .inc.js files are
// considered part of the package except for:
//
//    - .go files in package documentation
//    - files starting with _ or . (likely editor temporary files)
//    - files with build constraints not satisfied by the context
//
// If an error occurs, Import returns a non-nil error and a nil
// *PackageData.
func Import(path string, mode build.ImportMode, installSuffix string, buildTags []string) (*PackageData, error) {
	wd, err := os.Getwd()
	if err != nil {
		// Getwd may fail if we're in GOARCH=js mode. That's okay, handle
		// it by falling back to empty working directory. It just means
		// Import will not be able to resolve relative import paths.
		wd = ""
	}
	return importWithSrcDir(path, wd, mode, installSuffix, buildTags)
}

func importWithSrcDir(path string, srcDir string, mode build.ImportMode, installSuffix string, buildTags []string) (*PackageData, error) {
	bctx := NewReplContext(installSuffix, buildTags)
	switch path {
	case "syscall":
		// syscall needs to use a typical GOARCH like amd64 to pick up definitions for _Socklen, BpfInsn, IFNAMSIZ, Timeval, BpfStat, SYS_FCNTL, Flock_t, etc.
		bctx.GOARCH = runtime.GOARCH
		bctx.InstallSuffix = "js"
		if installSuffix != "" {
			bctx.InstallSuffix += "_" + installSuffix
		}
	case "math/big":
		// Use pure Go version of math/big; we don't want non-Go assembly versions.
		bctx.BuildTags = append(bctx.BuildTags, "math_big_pure_go")
	case "crypto/x509", "os/user":
		// These stdlib packages have cgo and non-cgo versions (via build tags); we want the latter.
		bctx.CgoEnabled = false
	}
	pkg, err := bctx.Import(path, srcDir, mode)
	if err != nil {
		return nil, err
	}

	// TODO: Resolve issue #415 and remove this temporary workaround.
	if strings.HasSuffix(pkg.ImportPath, "/vendor/github.com/glycerine/gofront/incr/js") {
		return nil, fmt.Errorf("vendoring github.com/glycerine/gofront/incr/js package is not supported, see https://github.com/gopherjs/gopherjs/issues/415")
	}

	switch path {
	case "os":
		pkg.GoFiles = excludeExecutable(pkg.GoFiles) // Need to exclude executable implementation files, because some of them contain package scope variables that perform (indirectly) syscalls on init.
	case "runtime":
		pkg.GoFiles = []string{"error.go"}
	case "runtime/internal/sys":
		pkg.GoFiles = []string{fmt.Sprintf("zgoos_%s.go", bctx.GOOS), "zversion.go"}
	case "runtime/pprof":
		pkg.GoFiles = nil
	case "internal/poll":
		pkg.GoFiles = exclude(pkg.GoFiles, "fd_poll_runtime.go")
	case "crypto/rand":
		pkg.GoFiles = []string{"rand.go", "util.go"}
	}

	if len(pkg.CgoFiles) > 0 {
		return nil, &ImportCError{path}
	}

	if pkg.IsCommand() {
		pkg.PkgObj = filepath.Join(pkg.BinDir, filepath.Base(pkg.ImportPath)+".js")
	}

	if _, err := os.Stat(pkg.PkgObj); os.IsNotExist(err) && strings.HasPrefix(pkg.PkgObj, build.Default.GOROOT) {
		// fall back to GOPATH
		firstGopathWorkspace := filepath.SplitList(build.Default.GOPATH)[0] // TODO: Need to check inside all GOPATH workspaces.
		gopathPkgObj := filepath.Join(firstGopathWorkspace, pkg.PkgObj[len(build.Default.GOROOT):])
		if _, err := os.Stat(gopathPkgObj); err == nil {
			pkg.PkgObj = gopathPkgObj
		}
	}

	jsFiles, err := jsFilesFromDir(pkg.Dir)
	if err != nil {
		return nil, err
	}

	return &PackageData{Package: pkg, JSFiles: jsFiles}, nil
}

// excludeExecutable excludes all executable implementation .go files.
// They have "executable_" prefix.
func excludeExecutable(goFiles []string) []string {
	var s []string
	for _, f := range goFiles {
		if strings.HasPrefix(f, "executable_") {
			continue
		}
		s = append(s, f)
	}
	return s
}

// exclude returns files, excluding specified files.
func exclude(files []string, exclude ...string) []string {
	var s []string
Outer:
	for _, f := range files {
		for _, e := range exclude {
			if f == e {
				continue Outer
			}
		}
		s = append(s, f)
	}
	return s
}

// ImportDir is like Import but processes the Go package found in the named
// directory.
func ImportDir(dir string, mode build.ImportMode, installSuffix string, buildTags []string) (*PackageData, error) {
	pkg, err := NewReplContext(installSuffix, buildTags).ImportDir(dir, mode)
	if err != nil {
		return nil, err
	}

	jsFiles, err := jsFilesFromDir(pkg.Dir)
	if err != nil {
		return nil, err
	}

	return &PackageData{Package: pkg, JSFiles: jsFiles}, nil
}

type Options struct {
	GOROOT         string
	GOPATH         string
	Verbose        bool
	Quiet          bool
	Watch          bool
	CreateMapFile  bool
	MapToLocalDisk bool
	Minify         bool
	Color          bool
	BuildTags      []string
}

func (o *Options) PrintError(format string, a ...interface{}) {
	if o.Color {
		format = "\x1B[31m" + format + "\x1B[39m"
	}
	fmt.Fprintf(os.Stderr, format, a...)
}

func (o *Options) PrintSuccess(format string, a ...interface{}) {
	if o.Color {
		format = "\x1B[32m" + format + "\x1B[39m"
	}
	fmt.Fprintf(os.Stderr, format, a...)
}

type PackageData struct {
	*build.Package
	JSFiles    []string
	IsTest     bool // IsTest is true if the package is being built for running tests.
	SrcModTime time.Time
	UpToDate   bool
}

type Session struct {
	options  *Options
	Archives map[string]*Archive
	Types    map[string]*types.Package
	Watcher  *fsnotify.Watcher
}

func NewSession(options *Options) *Session {
	if options.GOROOT == "" {
		options.GOROOT = build.Default.GOROOT
	}
	if options.GOPATH == "" {
		options.GOPATH = build.Default.GOPATH
	}
	options.Verbose = options.Verbose || options.Watch

	s := &Session{
		options:  options,
		Archives: make(map[string]*Archive),
	}
	s.Types = make(map[string]*types.Package)
	if options.Watch {
		if out, err := exec.Command("ulimit", "-n").Output(); err == nil {
			if n, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil && n < 1024 {
				fmt.Printf("Warning: The maximum number of open file descriptors is very low (%d). Change it with 'ulimit -n 8192'.\n", n)
			}
		}

		var err error
		s.Watcher, err = fsnotify.NewWatcher()
		if err != nil {
			panic(err)
		}
	}
	return s
}

func (s *Session) InstallSuffix() string {
	if s.options.Minify {
		return "min"
	}
	return ""
}

func jsFilesFromDir(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var jsFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".inc.js") && file.Name()[0] != '_' && file.Name()[0] != '.' {
			jsFiles = append(jsFiles, file.Name())
		}
	}
	return jsFiles, nil
}

// hasGopathPrefix returns true and the length of the matched GOPATH workspace,
// iff file has a prefix that matches one of the GOPATH workspaces.
func hasGopathPrefix(file, gopath string) (hasGopathPrefix bool, prefixLen int) {
	gopathWorkspaces := filepath.SplitList(gopath)
	for _, gopathWorkspace := range gopathWorkspaces {
		gopathWorkspace = filepath.Clean(gopathWorkspace)
		if strings.HasPrefix(file, gopathWorkspace) {
			return true, len(gopathWorkspace)
		}
	}
	return false, 0
}

func (s *Session) WaitForChange() {
	s.options.PrintSuccess("watching for changes...\n")
	for {
		select {
		case ev := <-s.Watcher.Events:
			if ev.Op&(fsnotify.Create|fsnotify.Write|fsnotify.Remove|fsnotify.Rename) == 0 || filepath.Base(ev.Name)[0] == '.' {
				continue
			}
			if !strings.HasSuffix(ev.Name, ".go") && !strings.HasSuffix(ev.Name, ".inc.js") {
				continue
			}
			s.options.PrintSuccess("change detected: %s\n", ev.Name)
		case err := <-s.Watcher.Errors:
			s.options.PrintError("watcher error: %s\n", err.Error())
		}
		break
	}

	go func() {
		for range s.Watcher.Events {
			// consume, else Close() may deadlock
		}
	}()
	s.Watcher.Close()
}

var gijitAnsPrefix = []byte("__gijit_ans := []interface{}{")
var gijitAnsSuffix = []byte("}\n __gijit_printQuoted(__gijit_ans...);")

// at the beginning of src, transform a first '='[^=] into
// "__gijit_ans := "
func prependAns(src []byte) []byte {
	nsrc := len(src)
	leftTrimmed := bytes.TrimLeftFunc(src, unicode.IsSpace)
	trimmed := bytes.TrimFunc(src, unicode.IsSpace)
	n := len(leftTrimmed)
	leftdiff := nsrc - n
	if n > 1 && leftTrimmed[0] == '=' && leftTrimmed[1] != '=' {
		return append(gijitAnsPrefix, append(trimmed[leftdiff+1:], gijitAnsSuffix...)...)
	}
	return src
}

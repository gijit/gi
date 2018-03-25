package compiler

import (
	"bytes"
	"fmt"
	"github.com/gijit/gi/pkg/ast"
	"github.com/gijit/gi/pkg/gostd/build"
	"github.com/gijit/gi/pkg/parser"
	"github.com/gijit/gi/pkg/scanner"
	"github.com/gijit/gi/pkg/token"
	//"github.com/gijit/gi/pkg/types"
	"io"
	"io/ioutil"
	"os"
	//"os/exec"
	"path"
	"path/filepath"
	//"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gijit/gi/pkg/compiler/natives"
	"github.com/neelance/sourcemap"
)

// from ~/go/src/github.com/gopherjs/gopherjs/tool.go:
// currentDirectory and the init() that follows.
var currentDirectory string

func init() {
	var err error
	currentDirectory, err = os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	currentDirectory, err = filepath.EvalSymlinks(currentDirectory)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	gopaths := filepath.SplitList(build.Default.GOPATH)
	if len(gopaths) == 0 {
		fmt.Fprintf(os.Stderr, "$GOPATH not set. For more details see: go help gopath\n")
		os.Exit(1)
	}
}

// original location of this file was
// ~/go/src/github.com/gopherjs/gopherjs/build/build.go
//
// jea: now it is half in translate.go, half here.

type ImportCError struct {
	pkgPath string
}

func (e *ImportCError) Error() string {
	return e.pkgPath + `: importing "C" is not supported by gijit`
}

// parseAndAugment parses and returns all .go files of given pkg.
// Standard Go library packages are augmented with files in compiler/natives folder.
// If isTest is true and pkg.ImportPath has no _test suffix, package is built for running internal tests.
// If isTest is true and pkg.ImportPath has _test suffix, package is built for running external tests.
//
// The native packages are augmented by the contents of natives.FS in the following way.
// The file names do not matter except the usual `_test` suffix. The files for
// native overrides get added to the package (even if they have the same name
// as an existing file from the standard library). For all identifiers that exist
// in the original AND the overrides, the original identifier in the AST gets
// replaced by `_`. New identifiers that don't exist in original package get added.
func parseAndAugment(pkg *build.Package, isTest bool, fileSet *token.FileSet) ([]*ast.File, error) {
	var files []*ast.File
	replacedDeclNames := make(map[string]bool)
	funcName := func(d *ast.FuncDecl) string {
		if d.Recv == nil || len(d.Recv.List) == 0 {
			return d.Name.Name
		}
		recv := d.Recv.List[0].Type
		if star, ok := recv.(*ast.StarExpr); ok {
			recv = star.X
		}
		return recv.(*ast.Ident).Name + "." + d.Name.Name
	}
	isXTest := strings.HasSuffix(pkg.ImportPath, "_test")
	importPath := pkg.ImportPath
	if isXTest {
		importPath = importPath[:len(importPath)-5]
	}

	nativesContext := &build.Context{
		GOROOT:   "/",
		GOOS:     build.Default.GOOS,
		GOARCH:   "gijit",
		Compiler: "gc",
		JoinPath: path.Join,
		SplitPathList: func(list string) []string {
			if list == "" {
				return nil
			}
			return strings.Split(list, "/")
		},
		IsAbsPath: path.IsAbs,
		IsDir: func(name string) bool {
			dir, err := natives.FS.Open(name)
			if err != nil {
				return false
			}
			defer dir.Close()
			info, err := dir.Stat()
			if err != nil {
				return false
			}
			return info.IsDir()
		},
		HasSubdir: func(root, name string) (rel string, ok bool) {
			panic("not implemented")
		},
		ReadDir: func(name string) (fi []os.FileInfo, err error) {
			dir, err := natives.FS.Open(name)
			if err != nil {
				return nil, err
			}
			defer dir.Close()
			return dir.Readdir(0)
		},
		OpenFile: func(name string) (r io.ReadCloser, err error) {
			return natives.FS.Open(name)
		},
	}
	if nativesPkg, err := nativesContext.Import(importPath, "", 0); err == nil {
		names := nativesPkg.GoFiles
		if isTest {
			names = append(names, nativesPkg.TestGoFiles...)
		}
		if isXTest {
			names = nativesPkg.XTestGoFiles
		}
		for _, name := range names {
			fullPath := path.Join(nativesPkg.Dir, name)
			r, err := nativesContext.OpenFile(fullPath)
			if err != nil {
				panic(err)
			}
			file, err := parser.ParseFile(fileSet, fullPath, r, parser.ParseComments)
			if err != nil {
				panic(err)
			}
			r.Close()
			for _, decl := range file.Nodes {
				switch d := decl.(type) {
				case *ast.FuncDecl:
					replacedDeclNames[funcName(d)] = true
				case *ast.GenDecl:
					switch d.Tok {
					case token.TYPE:
						for _, spec := range d.Specs {
							replacedDeclNames[spec.(*ast.TypeSpec).Name.Name] = true
						}
					case token.VAR, token.CONST:
						for _, spec := range d.Specs {
							for _, name := range spec.(*ast.ValueSpec).Names {
								replacedDeclNames[name.Name] = true
							}
						}
					}
				}
			}
			files = append(files, file)
		}
	}
	delete(replacedDeclNames, "init")

	var errList ErrorList
	for _, name := range pkg.GoFiles {
		if !filepath.IsAbs(name) {
			name = filepath.Join(pkg.Dir, name)
		}
		r, err := os.Open(name)
		if err != nil {
			return nil, err
		}
		file, err := parser.ParseFile(fileSet, name, r, parser.ParseComments)
		r.Close()
		if err != nil {
			if list, isList := err.(scanner.ErrorList); isList {
				if len(list) > 10 {
					list = append(list[:10], &scanner.Error{Pos: list[9].Pos, Msg: "too many errors"})
				}
				for _, entry := range list {
					errList = append(errList, entry)
				}
				continue
			}
			errList = append(errList, err)
			continue
		}

		switch pkg.ImportPath {
		case "crypto/rand", "encoding/gob", "encoding/json", "expvar", "go/token", "log", "math/big", "math/rand", "regexp", "testing", "time":
			for _, spec := range file.Imports {
				path, _ := strconv.Unquote(spec.Path.Value)
				if path == "sync" {
					if spec.Name == nil {
						spec.Name = ast.NewIdent("sync")
					}
					spec.Path.Value = `"github.com/gijit/gi/pkg/nosync"`
				}
			}
		}

		for _, decl := range file.Nodes {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if replacedDeclNames[funcName(d)] {
					d.Name = ast.NewIdent("_")
				}
			case *ast.GenDecl:
				switch d.Tok {
				case token.TYPE:
					for _, spec := range d.Specs {
						s := spec.(*ast.TypeSpec)
						if replacedDeclNames[s.Name.Name] {
							s.Name = ast.NewIdent("_")
						}
					}
				case token.VAR, token.CONST:
					for _, spec := range d.Specs {
						s := spec.(*ast.ValueSpec)
						for i, name := range s.Names {
							if replacedDeclNames[name.Name] {
								s.Names[i] = ast.NewIdent("_")
							}
						}
					}
				}
			}
		}
		files = append(files, file)
	}
	if errList != nil {
		return nil, errList
	}
	return files, nil
}

func (s *Session) BuildDir(packagePath string, importPath string, pkgObj string, isMain bool) (out []byte, err error) {

	buildPkg, err := NewReplContext(s.InstallSuffix(), s.options.BuildTags).ImportDir(packagePath, 0)
	if err != nil {
		return nil, err
	}
	pkg := &PackageData{Package: buildPkg}
	jsFiles, err := jsFilesFromDir(pkg.Dir)
	if err != nil {
		return nil, err
	}
	pkg.JSFiles = jsFiles
	archive, err := s.BuildPackage(pkg)
	if err != nil {
		return nil, err
	}
	if pkgObj == "" {
		pkgObj = filepath.Base(packagePath) + ".gijit"
	}
	if pkg.IsCommand() && !pkg.UpToDate {
		out, err = s.WriteCommandPackage(archive, pkgObj, isMain)
		if err != nil {
			return nil, err
		}
	}
	return
}

// filenames is the set of .go source files to compile,
// pkgObj is the output file path,
// packagePath is the current directory.
func (s *Session) BuildFiles(filenames []string, pkgObj string, packagePath string) (out []byte, err error) {
	pkg := &PackageData{
		Package: &build.Package{
			Name:       "main",
			ImportPath: "main",
			Dir:        packagePath,
		},
	}

	for _, file := range filenames {
		if strings.HasSuffix(file, ".inc.gijit") {
			pkg.JSFiles = append(pkg.JSFiles, file)
			continue
		}
		pkg.GoFiles = append(pkg.GoFiles, file)
	}

	archive, err := s.BuildPackage(pkg)
	if err != nil {
		return nil, err
	}
	if s.Types["main"].Name() != "main" {
		return nil, fmt.Errorf("cannot build/run non-main package")
	}
	isMain := true
	out, err = s.WriteCommandPackage(archive, pkgObj, isMain)
	return
}

func (s *Session) BuildImportPath(path string) (*Archive, error) {
	_, archive, err := s.BuildImportPathWithSrcDir(path, "")
	return archive, err
}

func (s *Session) BuildImportPathWithSrcDir(path string, srcDir string) (*PackageData, *Archive, error) {
	pp("Session.BuildImportPathWithSrcDir() top. path='%s', srcDir='%s'", path, srcDir)

	pkg, err := importWithSrcDir(path, srcDir, 0, s.InstallSuffix(), s.options.BuildTags)

	if err != nil {
		return nil, nil, err
	}

	archive, err := s.BuildPackage(pkg)
	if err != nil {
		return nil, nil, err
	}

	return pkg, archive, nil
}

func (s *Session) BuildPackage(pkg *PackageData) (*Archive, error) {

	if s.AllowImportCaching {
		pp("Session.BuildPackage() top. pkg.ImportPath='%s'", pkg.ImportPath)
		if archive, ok := s.Archives[pkg.ImportPath]; ok {
			pp("build.go:234 using cached version of archive for path '%s'\n", pkg.ImportPath)
			return archive, nil
		}
	}

	if pkg.PkgObj != "" {
		var fileInfo os.FileInfo
		gijitBinary, err := os.Executable()
		if err == nil {
			fileInfo, err = os.Stat(gijitBinary)
			if err == nil {
				pkg.SrcModTime = fileInfo.ModTime()
			}
		}
		if err != nil {
			os.Stderr.WriteString("Could not get gijit_build binary's modification timestamp. Please report issue.\n")
			pkg.SrcModTime = time.Now()
		}

		for _, importedPkgPath := range pkg.Imports {
			// Ignore all imports that aren't mentioned in import specs of pkg.
			// For example, this ignores imports such as runtime/internal/sys and runtime/internal/atomic.
			ignored := true
			for _, pos := range pkg.ImportPos[importedPkgPath] {
				importFile := filepath.Base(pos.Filename)
				for _, file := range pkg.GoFiles {
					if importFile == file {
						ignored = false
						break
					}
				}
				if !ignored {
					break
				}
			}

			if importedPkgPath == "unsafe" || ignored {
				continue
			}
			importedPkg, _, err := s.BuildImportPathWithSrcDir(importedPkgPath, pkg.Dir)
			if err != nil {
				return nil, err
			}
			impModTime := importedPkg.SrcModTime
			if impModTime.After(pkg.SrcModTime) {
				pkg.SrcModTime = impModTime
			}
		}

		for _, name := range append(pkg.GoFiles, pkg.JSFiles...) {
			fileInfo, err := os.Stat(filepath.Join(pkg.Dir, name))
			if err != nil {
				return nil, err
			}
			if fileInfo.ModTime().After(pkg.SrcModTime) {
				pkg.SrcModTime = fileInfo.ModTime()
			}
		}

		cachingAllowed := false
		if cachingAllowed {
			pkgObjFileInfo, err := os.Stat(pkg.PkgObj)
			if err == nil && !pkg.SrcModTime.After(pkgObjFileInfo.ModTime()) {
				pp("\n\n package object is up to date, load from disk if library\n\n")
				// package object is up to date, load from disk if library
				pkg.UpToDate = true
				if pkg.IsCommand() {
					return nil, nil
				}

				objFile, err := os.Open(pkg.PkgObj)
				if err != nil {
					return nil, err
				}
				defer objFile.Close()

				archive, err := ReadArchive(pkg.PkgObj, pkg.ImportPath, objFile, s.Types)
				if err != nil {
					return nil, err
				}

				pp("\n\n reading Archive from disk, filename='%s', import path='%s'. archive.Pkg='%#v'\n and archive='%#v'", pkg.PkgObj, pkg.ImportPath, archive.Pkg, archive)
				s.Archives[pkg.ImportPath] = archive
				return archive, err
			} // end ReadArchive from disk

		} // end if cachingAllowe
	}

	fileSet := token.NewFileSet()
	files, err := parseAndAugment(pkg.Package, pkg.IsTest, fileSet)
	if err != nil {
		return nil, err
	}

	localImportPathCache := make(map[string]*Archive)
	importContext := &ImportContext{
		Packages: s.Types,
		Import: func(path string) (*Archive, error) {
			if archive, ok := localImportPathCache[path]; ok {
				pp("\n\n using cached archive from localImportPathCache... archive.Pkg='%#v'\n", archive.Pkg)
				return archive, nil
			}
			_, archive, err := s.BuildImportPathWithSrcDir(path, pkg.Dir)
			if err != nil {
				return nil, err
			}
			localImportPathCache[path] = archive
			pp("\n\n saving archive into localImportPathCache. archive.Pkg='%#v'\n", archive.Pkg)
			return archive, nil
		},
	}
	archive, err := FullPackageCompile(pkg.ImportPath, files, fileSet, importContext, s.options.Minify)
	if err != nil {
		return nil, err
	}
	pp("\n\n archive back from FullPackageCompile, archive.Pkg='%#v'\n", archive.Pkg)

	for _, jsFile := range pkg.JSFiles {
		code, err := ioutil.ReadFile(filepath.Join(pkg.Dir, jsFile))
		if err != nil {
			return nil, err
		}
		archive.IncJSCode = append(archive.IncJSCode, []byte("\t(function() \n")...)
		archive.IncJSCode = append(archive.IncJSCode, code...)
		archive.IncJSCode = append(archive.IncJSCode, []byte("\n\t end)(_global);\n")...)
	}

	if s.options.Verbose {
		fmt.Println(pkg.ImportPath)
	}

	s.Archives[pkg.ImportPath] = archive

	if pkg.PkgObj == "" || pkg.IsCommand() {
		pp("\n\n returning early, pkg.PkgObj==\"\" or pkg.IsCommand()=%v, archive.Pkg='%#v'\n", pkg.IsCommand(), archive.Pkg)
		return archive, nil
	}

	if err := s.writeLibraryPackage(archive, pkg.PkgObj); err != nil {
		if strings.HasPrefix(pkg.PkgObj, s.options.GOROOT) {
			// fall back to first GOPATH workspace
			firstGopathWorkspace := filepath.SplitList(s.options.GOPATH)[0]
			if err := s.writeLibraryPackage(archive, filepath.Join(firstGopathWorkspace, pkg.PkgObj[len(s.options.GOROOT):])); err != nil {
				return nil, err
			}
			return archive, nil
		}
		return nil, err
	}

	return archive, nil
}

func (s *Session) writeLibraryPackage(archive *Archive, pkgObj string) error {
	if err := os.MkdirAll(filepath.Dir(pkgObj), 0777); err != nil {
		return err
	}

	objFile, err := os.Create(pkgObj)
	if err != nil {
		return err
	}
	defer objFile.Close()

	return WriteArchive(archive, objFile)
}

func (s *Session) WriteCommandPackage(archive *Archive, pkgObj string, isMain bool) (out []byte, err error) {
	var codeFile io.Writer
	if s.options.WriteToFile {
		if err := os.MkdirAll(filepath.Dir(pkgObj), 0777); err != nil {
			return nil, err
		}

		fd, err := os.Create(pkgObj)
		if err != nil {
			return nil, err
		}
		defer fd.Close()
		codeFile = fd
	} else {
		by := bytes.NewBuffer(nil)
		codeFile = by
		defer func() {
			out = by.Bytes()
		}()
	}

	sourceMapFilter := &SourceMapFilter{Writer: codeFile}
	if s.options.WriteToFile && s.options.CreateMapFile {

		m := &sourcemap.Map{File: filepath.Base(pkgObj)}
		mapFile, err := os.Create(pkgObj + ".map")
		if err != nil {
			return nil, err
		}

		defer func() {
			m.WriteTo(mapFile)
			mapFile.Close()
			fmt.Fprintf(codeFile, "--# sourceMappingURL=%s.map\n", filepath.Base(pkgObj))
		}()

		sourceMapFilter.MappingCallback = NewMappingCallback(m, s.options.GOROOT, s.options.GOPATH, s.options.MapToLocalDisk)
	}

	deps, err := ImportDependencies(archive, func(path string) (*Archive, error) {
		if archive, ok := s.Archives[path]; ok {
			return archive, nil
		}
		_, archive, err := s.BuildImportPathWithSrcDir(path, "")
		return archive, err
	})
	if err != nil {
		return nil, err
	}
	err = WriteProgramCode(deps, sourceMapFilter, isMain)
	return
}

func NewMappingCallback(m *sourcemap.Map, goroot, gopath string, localMap bool) func(generatedLine, generatedColumn int, originalPos token.Position) {
	return func(generatedLine, generatedColumn int, originalPos token.Position) {
		if !originalPos.IsValid() {
			m.AddMapping(&sourcemap.Mapping{GeneratedLine: generatedLine, GeneratedColumn: generatedColumn})
			return
		}

		file := originalPos.Filename

		switch hasGopathPrefix, prefixLen := hasGopathPrefix(file, gopath); {
		case localMap:
			// no-op:  keep file as-is
		case hasGopathPrefix:
			file = filepath.ToSlash(file[prefixLen+4:])
		case strings.HasPrefix(file, goroot):
			file = filepath.ToSlash(file[len(goroot)+4:])
		default:
			file = filepath.Base(file)
		}

		m.AddMapping(&sourcemap.Mapping{GeneratedLine: generatedLine, GeneratedColumn: generatedColumn, OriginalFile: file, OriginalLine: originalPos.Line, OriginalColumn: originalPos.Column})
	}
}

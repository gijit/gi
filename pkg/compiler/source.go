package compiler

import (
	"fmt"
	//"github.com/gijit/gi/pkg/token"
	//"github.com/gijit/gi/pkg/types"
	//golua "github.com/glycerine/golua/lua"
	//"github.com/glycerine/luar"
)

func (ic *IncrState) ImportSourcePackage(path, pkgDir string) (res *Archive, err error) {

	// from ~/go/src/github.com/gopherjs/gopherjs/build/build.go
	if archive, ok := ic.CurPkg.localImportPathCache[path]; ok {
		return archive, nil
	}
	_, archive, err := ic.CurPkg.Session.BuildImportPathWithSrcDir(path, pkgDir)
	if err != nil {
		return nil, err
	}
	ic.CurPkg.localImportPathCache[path] = archive
	return archive, nil

	/*
		var pkg *types.Package

		pkgName := pkg.Name()

		res = &Archive{
			Name:       pkgName,
			ImportPath: path,
			Pkg:        pkg,
		}

		pkg.SetPath(path)

		// very important, must do this or we won't locate the package!
		ic.CurPkg.importContext.Packages[path] = pkg
	*/
	return nil, fmt.Errorf("TODO: source import")
}

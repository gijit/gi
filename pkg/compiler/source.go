package compiler

import (
//"fmt"
//"github.com/gijit/gi/pkg/token"
//"github.com/gijit/gi/pkg/types"
//golua "github.com/glycerine/golua/lua"
//"github.com/glycerine/luar"
)

func (ic *IncrState) ImportSourcePackage(path, pkgDir string) (res *Archive, err error) {
	pp("IncrState.ImportSourcePackage() top. path='%s', pkgDir='%s'", path, pkgDir)

	// from ~/go/src/github.com/gopherjs/gopherjs/build/build.go
	if archive, ok := ic.CurPkg.localImportPathCache[path]; ok {
		return archive, nil
	}
	_, archive, err := ic.CurPkg.Session.BuildImportPathWithSrcDir(path, pkgDir)
	if err != nil {
		return nil, err
	}
	ic.CurPkg.localImportPathCache[path] = archive

	// very important, must do this or we won't locate the package!
	ic.CurPkg.importContext.Packages[path] = archive.Pkg

	return archive, nil
}

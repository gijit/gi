package compiler

import (
	"fmt"
	//"github.com/gijit/gi/pkg/token"
	//"github.com/gijit/gi/pkg/types"
	//golua "github.com/glycerine/golua/lua"
	//"github.com/glycerine/luar"
)

func (ic *IncrState) ImportSourcePackage(path string) (res *Archive, err error) {
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

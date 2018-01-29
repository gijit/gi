/*
genimport: create a map from string (pkg.FuncName) -> function pointer
*/
package main

import (
	"fmt"
	"os"

	"github.com/gijit/gi/pkg/importer"
	"github.com/gijit/gi/pkg/types"
	"github.com/gijit/gi/pkg/verb"
)

var pp = verb.PP

func main() {
	readImport("fmt", "", "main")
}

func readImport(importPath, dir string, residentPkg string) error {
	var pkg *types.Package

	imp := importer.Default()
	imp2, ok := imp.(types.ImporterFrom)
	if !ok {
		panic("importer.ImportFrom not available, vendored packages would be lost")
	}
	var mode types.ImportMode
	var err error
	pkg, err = imp2.ImportFrom(importPath, dir, mode)

	if err != nil {
		return err
	}

	pkgName := pkg.Name()

	o, err := os.Create(pkgName + ".genimp.go")
	if err != nil {
		return err
	}

	fmt.Fprintf(o, `package %s

import "%s"

var Pkg = make(map[string]interface{})
func init() {
`, residentPkg, importPath)

	scope := pkg.Scope()
	nms := scope.Names()

	for _, nm := range nms {
		obj := scope.Lookup(nm)
		if !obj.Exported() {
			continue
		}
		switch obj.Type().(type) {
		case *types.Signature:
			pp("in package '%s', registering nm='%s' -> '%#v'", pkgName, nm, obj)
			fmt.Fprintf(o, `    Pkg["%s"] = %s.%s
`, nm, pkgName, nm)
		}
	}
	fmt.Fprintf(o, "\n}")

	return nil
}

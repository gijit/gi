package compiler

import (
	"fmt"
	"os"

	"github.com/gijit/gi/pkg/importer"
	"github.com/gijit/gi/pkg/types"
)

//
//GenShadowImport: create a map from string (pkg.FuncName) -> function pointer
//  that can be used inside the "shadow" REPL environment that Luar can call.
//
func GenShadowImport(importPath, dirForVendor, residentPkg, outDir string) error {
	var pkg *types.Package

	imp := importer.Default()
	imp2, ok := imp.(types.ImporterFrom)
	if !ok {
		panic("importer.ImportFrom not available, vendored packages would be lost")
	}
	var mode types.ImportMode
	var err error
	pkg, err = imp2.ImportFrom(importPath, dirForVendor, mode)

	if err != nil {
		return err
	}

	pkgName := pkg.Name()

	o, err := os.Create(outDir + string(os.PathSeparator) + pkgName + ".genimp.go")
	if err != nil {
		return err
	}

	fmt.Fprintf(o, `package shadow_%s

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
		case *types.Tuple:
			continue
		default:
			//case *types.Signature:
			pp("in package '%s', registering nm='%s' -> '%#v'", pkgName, nm, obj)
			fmt.Fprintf(o, `    Pkg["%s"] = %s.%s
`, nm, pkgName, nm)
		}
	}
	fmt.Fprintf(o, "\n}")

	return nil
}

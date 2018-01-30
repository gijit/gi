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

	atEnd := []string{}

	for _, nm := range nms {

		obj := scope.Lookup(nm)
		if !obj.Exported() {
			continue
		}

		oty := obj.Type()
		under := oty.Underlying()

		switch oty.(type) {
		default:
			panic(fmt.Sprintf("genshadow: "+
				"unhandled type! what oty type? '%T'", oty))
		case *types.Basic:
			// these all compile
			direct(o, nm, pkgName)

		case *types.Signature:
			// these all compile
			direct(o, nm, pkgName)

		case *types.Named:

			switch under.(type) {
			case *types.Interface:
				switch obj.(type) {
				case *types.TypeName:
					ifaceTemplate(o, obj, nm, pkgName, oty, under, &atEnd)
				case *types.Var:
					direct(o, nm, pkgName)
				}
			case *types.Struct:
				// none of these in "io"
				structTemplate(o, obj, nm, pkgName, oty, under)
			}
		}
	}
	fmt.Fprintf(o, "\n}")

	for _, s := range atEnd {
		fmt.Fprintf(o, "%s\n", s)
	}
	return nil
}

/* make a function like:
func __gi_ConvertTo_Reader(x interface{}) (y io.Reader, b bool) {
	y, b = x.(io.Reader)
	return
}
*/
func GenInterfaceConverter(pkg, name string) (funcName, decl string) {

	funcName = fmt.Sprintf("__gi_ConvertTo_%s", name)

	decl = fmt.Sprintf(`
func %s(x interface{}) (y %s.%s, b bool) {
	y, b = x.(%s.%s)
	return
}
`, funcName, pkg, name, pkg, name)
	return
}

func direct(o *os.File, nm, pkgName string) {
	fmt.Fprintf(o, "    Pkg[\"%s\"] = %s.%s\n", nm, pkgName, nm)
}

func structTemplate(o *os.File, obj types.Object, nm, pkgName string, oty, under types.Type) {
	// example from "io":
	/*
		type PipeReader struct {
			p *pipe
		}
	*/

	fmt.Printf("we have under as a struct type, for nm='%s', obj='%#v'\n", nm, obj)
	pp("ignoring nm='%v' as under is *types.Struct", nm)
}

func ifaceTemplate(o *os.File, obj types.Object, nm, pkgName string, oty, under types.Type, atEnd *[]string) {

	//pp("ifaceTemplate:: we see Named '%s'\n. oty:'%#v',\n under:'%#v',\n, obj='%#v', \n", nm, oty, under, obj)

	funcName, decl := GenInterfaceConverter(pkgName, nm)
	*atEnd = append((*atEnd), decl)
	fmt.Fprintf(o, "    Pkg[\"%s\"] = %s\n", nm, funcName)
}

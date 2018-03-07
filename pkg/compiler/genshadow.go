package compiler

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/gijit/gi/pkg/importer"
	"github.com/gijit/gi/pkg/types"
)

//
//GenShadowImport: create a map from string (pkg.FuncName) -> function pointer
//  that can be used inside the "shadow" REPL environment that Luar can call.
//
func GenShadowImport(importPath, dirForVendor, residentPkg, outDir string) error {
	var pkg *types.Package

	structs := []string{}
	base := filepath.Base(residentPkg)
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
var Ctor = make(map[string]interface{})

func init() {
`, base, importPath)

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

		pp("oty = '%#v', under='%#v'", oty, under)

		switch oty.(type) {
		default:
			panic(fmt.Sprintf("genshadow: "+
				"unhandled type! what oty type? '%T' for nm='%v'", oty, nm))
		case *types.Slice:
			//fmt.Printf("Slice: nm = '%s', obj='%#v', oty='%#v', under='%#v'\n", nm, obj, oty, under)
			// ex: os.Args
			switch obj.(type) {
			case *types.Var:
				direct(o, nm, pkgName)
			default:
				panic(fmt.Sprintf("genshadow: "+
					"unhandled type! what oty type? '%T' for nm='%v'", oty, nm))
			}
		case *types.Pointer:
			//fmt.Printf("Pointer: nm = '%s', obj='%#v', oty='%#v', under='%#v'\n", nm, obj, oty, under)
			// ex: os.Stderr, os.Stdin, os.Stdout
			switch obj.(type) {
			case *types.Var:
				direct(o, nm, pkgName)
			default:
				panic(fmt.Sprintf("genshadow: "+
					"unhandled type! what oty type? '%T' for nm='%v'", oty, nm))
			}
		case *types.Basic:
			// these all compile
			direct(o, nm, pkgName)

		case *types.Signature:
			// these all compile
			direct(o, nm, pkgName)

		case *types.Named:
			pp("oty is types.Named...")
			switch under.(type) {
			case *types.Interface:
				switch obj.(type) {
				case *types.TypeName:
					ifaceTemplate(o, obj, nm, pkgName, oty, under, &atEnd)
				case *types.Var:
					direct(o, nm, pkgName)
				default:
					panic(fmt.Sprintf("genshadow: "+
						"unhandled type! what oty type? '%T' for nm='%v'", oty, nm))
				}
			case *types.Struct:
				// none of these in "io"
				structTemplate(o, obj, nm, pkgName, oty, under, &atEnd)
				structs = append(structs, nm)
			}
		}
	}
	fmt.Fprintf(o, "\n}")

	for _, s := range atEnd {
		fmt.Fprintf(o, "%s\n", s)
	}

	pp("structs = '%#v'", structs)

	// write the InitLua() function that
	// sets up the native struct (copy) constructors.
	fmt.Fprintf(o, "%s", genInitLuaStart(pkgName))
	for _, r := range structs {
		fmt.Fprintf(o, "%s", perStructInitLua(pkgName, r))
	}
	fmt.Fprintf(o, "%s", genInitLuaFinish(pkgName))

	return nil
}

/* make a function like:
func __gi_ConvertTo2_Reader(x interface{}) (y io.Reader, b bool) {
	y, b = x.(io.Reader)
	return
}

and one like:
func __gi_ConvertTo1_Reader(x interface{}) io.Reader {
	return x.(io.Reader)
}

*/
func GenInterfaceConverter(pkg, name string) (funcName1, funcName2, decl string) {

	funcName2 = fmt.Sprintf("GijitShadow_InterfaceConvertTo2_%s", name)
	funcName1 = fmt.Sprintf("GijitShadow_InterfaceConvertTo1_%s", name)

	decl = fmt.Sprintf(`
func %s(x interface{}) (y %s.%s, b bool) {
	y, b = x.(%s.%s)
	return
}
`, funcName2, pkg, name, pkg, name)

	decl += fmt.Sprintf(`
func %s(x interface{}) %s.%s {
	return x.(%s.%s)
}
`, funcName1, pkg, name, pkg, name)

	return
}

func direct(o *os.File, nm, pkgName string) {
	fmt.Fprintf(o, "    Pkg[\"%s\"] = %s.%s\n", nm, pkgName, nm)
}

func ctor(o *os.File, nm, pkgName string) {
	fmt.Fprintf(o, "    Ctor[\"%[1]s\"] = GijitShadow_NewStruct_%[1]s\n", nm)
}

/* make a function like:
func GijitShadow_CtorStruct_Time(src *time.Time) *time.Time {
	if src == nil {
		return &time.Time{}
	}
	a := *src
	return &a
}
*/
func structTemplate(o *os.File, obj types.Object, nm, pkgName string, oty, under types.Type, atEnd *[]string) {
	// example from "io":
	/*
		type PipeReader struct {
			p *pipe
		}
	*/
	switch obj.(type) {
	case *types.TypeName:
		*atEnd = append(*atEnd, fmt.Sprintf(`
func GijitShadow_NewStruct_%[2]s(src *%[1]s.%[2]s) *%[1]s.%[2]s {
    if src == nil {
	   return &%[1]s.%[2]s{}
    }
    a := *src
    return &a
}
`, pkgName, nm))
		ctor(o, nm, pkgName)
	default:
		direct(o, nm, pkgName)
	}

}

func ifaceTemplate(o *os.File, obj types.Object, nm, pkgName string, oty, under types.Type, atEnd *[]string) {

	//pp("ifaceTemplate:: we see Named '%s'\n. oty:'%#v',\n under:'%#v',\n, obj='%#v', \n", nm, oty, under, obj)

	_, funcName2, decl := GenInterfaceConverter(pkgName, nm)
	*atEnd = append((*atEnd), decl)
	fmt.Fprintf(o, "    Pkg[\"%s\"] = %s\n", nm, funcName2)
	//fmt.Fprintf(o, "    Pkg[\"%s\"] = %s\n", nm, funcName1)
}

func genInitLuaStart(shortPkg string) string {

	return fmt.Sprintf("\n\n func InitLua() string {\n  "+
		"return `\n__type__.%s ={};\n", shortPkg)
}

func genInitLuaFinish(shortPkg string) string {
	return "\n`}"
}

func perStructInitLua(shortPkg, structName string) string {

	return fmt.Sprintf(`
-----------------
-- struct %[2]s
-----------------

__type__.%[1]s.%[2]s = {
 __name = "native_Go_struct_type_wrapper",
 __native_type = "%[2]s",
 __call = function(t, src)
   return __ctor__%[1]s.%[2]s(src)
 end,
};
setmetatable(__type__.%[1]s.%[2]s, __type__.%[1]s.%[2]s);

`, shortPkg, structName)
}

package typesutil

import (
	//"fmt"
	"github.com/glycerine/gi/pkg/types"
	"strings"
)

func IsJsPackage(pkg *types.Package) bool {
	return pkg != nil && (pkg.Path() == "github.com/glycerine/gi/pkg/luaapi" || strings.HasSuffix(pkg.Path(), "/vendor/github.com/glycerine/gi/pkg/luaapi"))
}

func IsJsObject(t types.Type) bool {
	ptr, isPtr := t.(*types.Pointer)
	if !isPtr {
		return false
	}
	named, isNamed := ptr.Elem().(*types.Named)
	return isNamed && IsJsPackage(named.Obj().Pkg()) && named.Obj().Name() == "Object"
}

func IsLuarPackage(pkg *types.Package) bool {
	//fmt.Printf("\n pkg.Path()='%s'\n", pkg.Path())
	return pkg != nil && strings.HasPrefix(pkg.Path(), "github.com/glycerine/gi/pkg/compiler/shadow")
}

func IsLuarObject(t types.Type) bool {
	ptr, isPtr := t.(*types.Pointer)
	if !isPtr {
		//fmt.Printf("\n IsLuarObject, not b/c !isPtr\n")
		return false
	}
	named, isNamed := ptr.Elem().(*types.Named)
	if isNamed {
		//fmt.Printf(" named.Obj().Name()= '%#v', named.Obj().Pkg()='%#v'\n", named.Obj().Name(), named.Obj().Pkg())
	}
	return isNamed && IsLuarPackage(named.Obj().Pkg())
}

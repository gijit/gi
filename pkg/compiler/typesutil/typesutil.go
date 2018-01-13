package typesutil

import (
	"github.com/gijit/gi/pkg/types"
	"strings"
)

func IsJsPackage(pkg *types.Package) bool {
	return pkg != nil && (pkg.Path() == "github.com/gijit/gi/pkg/gopherjs/js" || strings.HasSuffix(pkg.Path(), "/vendor/github.com/glycerine/gofront/gopherjs/js"))
}

func IsJsObject(t types.Type) bool {
	ptr, isPtr := t.(*types.Pointer)
	if !isPtr {
		return false
	}
	named, isNamed := ptr.Elem().(*types.Named)
	return isNamed && IsJsPackage(named.Obj().Pkg()) && named.Obj().Name() == "Object"
}

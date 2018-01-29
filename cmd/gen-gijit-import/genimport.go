/*
genimport: create a map from string (pkg.FuncName) -> function pointer
*/
package main

import (
	"github.com/gijit/gi/pkg/compiler"
)

func main() {
	compiler.GenShadowImport("fmt", "", "main", ".")
}

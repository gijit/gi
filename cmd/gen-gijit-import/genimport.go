/*
genimport: create a map from string (pkg.FuncName) -> function pointer
*/
package main

import (
	"os"

	"github.com/gijit/gi/pkg/compiler"
)

func main() {
	compiler.GenShadowImport(os.Args[1], "", "main", ".")
}

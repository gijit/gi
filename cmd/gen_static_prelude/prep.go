// prep the prelude for static inclusion with the
// `gi` binary.
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/shurcooL/vfsgen"
)

func main() {

	gopath := os.Getenv("GOPATH")
	compiler := gopath + "/src/github.com/gijit/gi/pkg/compiler"
	prelude := compiler + "/prelude"
	gentarget := compiler + "/prelude_static.go"
	var fs http.FileSystem = http.Dir(prelude)

	err := vfsgen.Generate(fs, vfsgen.Options{
		Filename:    gentarget,
		PackageName: "compiler",
		//BuildTags: "!dev",
		VariableName: "preludeFiles",
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("gen_static_prelude '%s' ->\n   '%s'\n", prelude, gentarget)
}

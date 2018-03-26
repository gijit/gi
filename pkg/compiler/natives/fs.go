// plus build gopherjsdev

package natives

import (
	"github.com/gijit/gi/pkg/gostd/build"
	"log"
	"net/http"
	"os"
	"strings"

	//"github.com/gijit/gi/pkg/verb"
	"github.com/shurcooL/httpfs/filter"
)

func importPathToDir(importPath string) string {
	p, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		log.Fatalln(err)
	}
	return p.Dir
}

// FS is a virtual filesystem that contains native packages.
var FS = filter.Keep(
	http.Dir(importPathToDir("github.com/gijit/gi/pkg/compiler/natives")),
	func(path string, fi os.FileInfo) bool {
		if strings.HasSuffix(path, "~") {
			//verb.VV("jea debug: native.FS rejecting path ending in ~tilde~: '%s'", path)
			return false
		}
		res := path == "/" || path == "/src" || strings.HasPrefix(path, "/src/")
		/*		if res {
					verb.VV("jea debug: native.FS is keeping path: '%s' -> %v", path, res)
				} else {
					verb.VV("jea debug: native.FS is rejecting path: '%s' because res was %v", path, res)
				}
		*/
		return res
	},
)

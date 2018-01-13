// +build go1.4,!go1.6

package gotool

import (
	"github.com/gijit/gi/pkg/build"
	"path/filepath"
	"runtime"
)

var gorootSrc = filepath.Join(runtime.GOROOT(), "src")

func shouldIgnoreImport(p *build.Package) bool {
	return true
}

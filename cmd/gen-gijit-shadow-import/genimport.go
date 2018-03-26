/*
genimport: create a map from string (pkg.FuncName) -> function pointer
*/
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/glycerine/gi/pkg/compiler"
)

var defaultPreludePath = "src/github.com/glycerine/gi/pkg/compiler/shadow"

var defaultPreludePathParts []string

func init() {
	defaultPreludePathParts = strings.Split(defaultPreludePath, "/")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "supply package name to shadow as only argument.\n")
		os.Exit(1)
	}
	pkg := os.Args[1]

	odir := "."
	dir := os.Getenv("GOINTERP_PRELUDE_DIR")
	if dir != "" {
		odir = dir + "/shadow"
	} else {
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			// try $HOME/go
			home := os.Getenv("HOME")
			proposed := filepath.Join(home, "go")
			if !DirExists(home) || !DirExists(proposed) {
				fmt.Fprintf(os.Stderr, "gen-gijit-shadow-import error: "+
					"could not locate output pkg/compiler/shadow dir.\n")
				os.Exit(1)
			}
			gopath = proposed
		}

		odir = filepath.Join(append([]string{gopath}, defaultPreludePathParts...)...)
	}

	odir += string(os.PathSeparator) + pkg
	os.MkdirAll(odir, 0777)
	fmt.Printf("writing to odir '%s'\n", odir)
	cwd, err := os.Getwd()
	panicOn(err)

	err = compiler.GenShadowImport(pkg, cwd, pkg, odir)
	panicOn(err)
}

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}

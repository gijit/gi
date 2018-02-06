package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gijit/gi/pkg/verb"
)

var ProgramName string = path.Base(os.Args[0])

type GIConfig struct {
	Quiet          bool
	Verbose        bool
	VerboseVerbose bool
	RawLua         bool
	PreludePath    string
	IsTestMode     bool
	NoLiner        bool // for under test/emacs
}

// call DefineFlags before myflags.Parse()
func (c *GIConfig) DefineFlags(fs *flag.FlagSet) {
	fs.BoolVar(&c.Quiet, "q", false, "don't show banner on startup")
	fs.BoolVar(&c.Verbose, "v", false, "show debug prints")
	fs.BoolVar(&c.VerboseVerbose, "vv", false, "show even more verbose debug prints")
	fs.BoolVar(&c.RawLua, "r", false, "raw mode: skip all translation, type raw Lua to LuaJIT with our prelude installed")
	fs.StringVar(&c.PreludePath, "prelude", "", "path to the prelude directory. All .lua files are sourced before startup from this directory. Default is to to read from 'GOINTERP_PRELUDE_DIR' env var. -prelude overrides this.")
	fs.BoolVar(&c.IsTestMode, "t", true, "load test mode functions and types")
	fs.BoolVar(&c.NoLiner, "no-liner", false, "turn off liner, e.g. under emacs")
}

var defaultPreludePath = "src/github.com/gijit/gi/pkg/compiler"

var defaultPreludePathParts []string

func init() {
	defaultPreludePathParts = strings.Split(defaultPreludePath, "/")
}

// call c.ValidateConfig() after myflags.Parse()
func (c *GIConfig) ValidateConfig() error {

	if c.PreludePath == "" {
		dir := os.Getenv("GOINTERP_PRELUDE_DIR")
		if dir != "" {
			c.PreludePath = dir
		} else {
			// try hard... try $GOPATH/src/github.com/gijit/gi/pkg/compiler
			// by default.
			gopath := os.Getenv("GOPATH")
			if gopath == "" {
				// try $HOME/go
				home := os.Getenv("HOME")
				proposed := filepath.Join(home, "go")
				if !DirExists(home) || !DirExists(proposed) {
					return preludeError()
				}
				gopath = proposed
			}

			c.PreludePath = filepath.Join(append([]string{gopath}, defaultPreludePathParts...)...)
		}
	}
	verb.Verbose = c.Verbose || c.VerboseVerbose
	verb.VerboseVerbose = c.VerboseVerbose

	return nil
}

func preludeError() error {
	return fmt.Errorf("setenv GOINTERP_PRELUDE_DIR to point to your prelude dir. This is typically $GOPATH/src/github.com/gijit/gi/pkg/compiler but GOINTERP_PRELUDE_DIR was not set and -prelude was not specified.")
}

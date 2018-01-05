package main

import (
	"flag"
	"github.com/go-interpreter/gi/pkg/compiler"
	"os"
	"path"
)

var ProgramName string = path.Base(os.Args[0])

type GIConfig struct {
	Quiet          bool
	Verbose        bool
	VerboseVerbose bool
	RawLua         bool
	PreludePath    string

	preludeFiles []string
}

// call DefineFlags before myflags.Parse()
func (c *GIConfig) DefineFlags(fs *flag.FlagSet) {
	fs.BoolVar(&c.Quiet, "q", false, "don't show banner on startup")
	fs.BoolVar(&c.Verbose, "v", false, "show debug prints")
	fs.BoolVar(&c.VerboseVerbose, "vv", false, "show even more verbose debug prints")
	fs.BoolVar(&c.RawLua, "raw", false, "skip all translation, type raw Lua to LuaJIT with our prelude installed")
	fs.StringVar(&c.PreludePath, "prelude", "./prelude", "path to the prelude directory. All .lua files are sourced before startup from this directory.")
}

// call c.ValidateConfig() after myflags.Parse()
func (c *GIConfig) ValidateConfig() error {

	files, err := compiler.FetchPrelude(c.PreludePath)
	if err != nil {
		return err
	}
	c.preludeFiles = files
	return nil
}

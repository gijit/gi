package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

var ProgramName string = path.Base(os.Args[0])

type GIConfig struct {
	Quiet          bool
	Verbose        bool
	VerboseVerbose bool
	RawLua         bool
	PreludePath    string
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
	if !DirExists(c.PreludePath) {
		return fmt.Errorf("-prelude dir does not exist: '%s'", c.PreludePath)
	}
	files, err := filepath.Glob("*.lua")
	if err != nil {
		return fmt.Errorf("-prelude dir '%s' open problem: '%v'", c.PreludePath, err)
	}
	if len(files) < 1 {
		return fmt.Errorf("-prelude dir '%s' had no lua files in it.", c.PreludePath)
	}
	return nil
}

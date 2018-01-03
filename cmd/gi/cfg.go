package main

import (
	"flag"
	"os"
	"path"
)

var ProgramName string = path.Base(os.Args[0])

type GIConfig struct {
	Quiet bool
}

// call DefineFlags before myflags.Parse()
func (c *GIConfig) DefineFlags(fs *flag.FlagSet) {
	fs.BoolVar(&c.Quiet, "q", false, "don't show banner on startup")
}

// call c.ValidateConfig() after myflags.Parse()
func (c *GIConfig) ValidateConfig() error {

	return nil
}

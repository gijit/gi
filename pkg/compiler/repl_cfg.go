package compiler

import (
	"flag"

	"github.com/gijit/gi/pkg/verb"
)

type GIConfig struct {
	Quiet          bool
	Verbose        bool
	VerboseVerbose bool
	RawLua         bool
	CalculatorMode bool
	PreludePath    string
	IsTestMode     bool
	NoLiner        bool // for under test/emacs
	NoPrelude      bool
	NoLuar         bool

	Dev bool // dev mode, don't use statically cached prelude
}

var defaultTestMode bool // set to true by init() for tests, in repl_test.go.

func NewGIConfig() *GIConfig {
	return &GIConfig{
		IsTestMode: defaultTestMode, // under tests, is set to true
	}
}

// call DefineFlags before myflags.Parse()
func (c *GIConfig) DefineFlags(fs *flag.FlagSet) {
	fs.BoolVar(&c.Quiet, "q", false, "don't show banner on startup")
	fs.BoolVar(&c.Verbose, "v", false, "show debug prints")
	fs.BoolVar(&c.VerboseVerbose, "vv", false, "show even more verbose debug prints")
	fs.BoolVar(&c.RawLua, "r", false, "raw mode: skip all translation, type raw Lua to LuaJIT with our prelude installed")
	fs.StringVar(&c.PreludePath, "prelude", "", "path to the prelude directory. All *.lua files are sourced before startup from this directory. Default is to use the statically embedded version.")
	fs.BoolVar(&c.IsTestMode, "t", false, "load test mode functions and types")
	fs.BoolVar(&c.NoLiner, "no-liner", false, "turn off liner, e.g. under emacs")
	fs.BoolVar(&c.NoPrelude, "np", false, "no prelude; skip loading the prelude .lua files and Luar. implies -r raw mode too.")
	fs.BoolVar(&c.Dev, "d", false, "dev mode uses the pkg/compiler/prelude/*.lua files, skipping the statically cached pkg/compiler/prelude_static.go version.")
}

// call c.ValidateConfig() after myflags.Parse()
func (c *GIConfig) ValidateConfig() error {

	if c.NoPrelude {
		c.NoLuar = true
		c.RawLua = true
	}

	if c.PreludePath == "" {
		// just use the statically embedded prelude from build time.
	}
	verb.Verbose = c.Verbose || c.VerboseVerbose
	verb.VerboseVerbose = c.VerboseVerbose

	return nil
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gijit/gi/pkg/verb"
)

var p = verb.P
var pp = verb.PP

func main() {

	myflags := flag.NewFlagSet("gi", flag.ExitOnError)
	cfg := &GIConfig{}
	cfg.DefineFlags(myflags)

	err := myflags.Parse(os.Args[1:])
	err = cfg.ValidateConfig()
	if err != nil {
		log.Fatalf("%s command line flag error: '%s'", ProgramName, err)
	}

	verb.Verbose = cfg.Verbose || cfg.VerboseVerbose
	verb.VerboseVerbose = cfg.VerboseVerbose

	if !cfg.Quiet {
		fmt.Printf(
			`====================
gi: a go interpreter
====================
https://github.com/gijit/gi
Copyright (c) 2018, Jason E. Aten, Ph.D.
License: 3-clause BSD. See the LICENSE file at
https://github.com/gijit/gi/blob/master/LICENSE
====================
  [ gi is an interactive Golang environment,
    also known as a REPL or Read-Eval-Print-Loop ]
  [ type ctrl-d to exit ]
  [ type :help for help ]
  [ gi -h for flag help ]
  [ gi -q to start quietly ]
====================
%s
==================
`, Version())
	}

	cfg.LuajitMain()
}

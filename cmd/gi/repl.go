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

	if !cfg.Quiet {
		fmt.Printf(
			`====================
gijit: a go interpreter, just-in-time.
====================
https://github.com/gijit/gi
Copyright (c) 2018, Jason E. Aten. All rights reserved.
License: 3-clause BSD. See the LICENSE file at
https://github.com/gijit/gi/blob/master/LICENSE
====================
  [ gigit/gi is an interactive Golang environment,
    also known as a REPL or Read-Eval-Print-Loop.]
  [ at the gi> prompt, type ctrl-d to exit.]
  [ at the gi> prompt, type :help for special commands.]
  [ $ gi -h for flag help, when first launching gijit.]
  [ $ gi -q to start quietly, without this banner.]
====================
%s
==================
`, Version())
	}

	cfg.LuajitMain()
}

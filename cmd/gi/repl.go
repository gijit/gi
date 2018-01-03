package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

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
			`==================
gi: go interactive (https://github.com/go-interpreter/gi)
==================
  [gi is an interactive Golang environment,
   also known as a REPL or Read-Eval-Print-Loop.]
  [type ctrl-d to exit]
  [type :help for help]
  [gi -h for flag help]
  [gi -q to start quietly]
==================
%s
==================
`, Version())
	}

	LuajitMain()
	//NodeChildMain()
	//OttoReplMain()
}

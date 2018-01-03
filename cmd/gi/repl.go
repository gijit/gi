package main

import "fmt"

func main() {
	fmt.Printf(`
gi: go interactive
  [an interactive Golang environment (aka a REPL)]
  [type ctrl-d to exit]
==================
%s
==================
`, Version())

	LuajitMain()
	//NodeChildMain()
	//OttoReplMain()
}

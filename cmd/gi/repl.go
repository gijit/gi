package main

import "fmt"

func main() {
	fmt.Printf(`
gi: go interactive
==================
%s
==================
`, Version())

	LuajitMain()
	//NodeChildMain()
	//OttoReplMain()
}

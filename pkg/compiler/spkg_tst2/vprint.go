package spkg_tst2

import (
	"fmt"
)

var Verbose bool = false

var A, B = true, false

func P(format string, a ...interface{}) {
	fmt.Printf("package spkg_tst2 var Verbose is %v\n", Verbose)
	if Verbose {
		fmt.Printf(format, a...)
	}
}

func ToString(format string, a ...interface{}) string {
	s := fmt.Sprintf("Verbose=%v", Verbose)
	return s + fmt.Sprintf(format, a...)
}

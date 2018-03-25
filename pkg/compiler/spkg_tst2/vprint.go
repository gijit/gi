package spkg_tst2

import (
	"fmt"
)

var Verbose bool = false

func P(format string, a ...interface{}) {
	if Verbose {
		fmt.Printf(format, a...)
	}
}

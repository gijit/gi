package compiler

import (
	"fmt"
	"io"
	"os"
	"time"
)

// for tons of debug output (see also WorkerVerbose)
var Verbose bool = true

func p(format string, a ...interface{}) {
	if Verbose {
		TSPrintf(format, a...)
	}
}

func pp(format string, a ...interface{}) {
	TSPrintf(format, a...)
}

func pb(w io.Writer, format string, a ...interface{}) {
	if Verbose {
		fmt.Fprintf(w, "\n"+format+"\n", a...)
	}
}

// time-stamped printf
func TSPrintf(format string, a ...interface{}) {
	Printf("\n%s ", ts())
	Printf(format+"\n", a...)
}

// get timestamp for logging purposes
func ts() string {
	return time.Now().Format("2006-01-02 15:04:05.999 -0700 MST")
}

// so we can multi write easily, use our own printf
var OurStdout io.Writer = os.Stdout

// Printf formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
func Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(OurStdout, format, a...)
}

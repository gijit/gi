package verb

import (
	"fmt"
	"io"
	"os"
	"time"
)

// for tons of debug output
var Verbose bool = true
var VerboseVerbose bool = true

func P(format string, a ...interface{}) {
	if Verbose {
		TSPrintf(format, a...)
	}
}

func PP(format string, a ...interface{}) {
	if VerboseVerbose {
		TSPrintf(format, a...)
	}
}

func PB(w io.Writer, format string, a ...interface{}) {
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

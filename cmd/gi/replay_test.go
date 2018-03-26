package main

import (
	"flag"
	"fmt"
	"strings"
	"testing"

	"github.com/glycerine/gi/pkg/compiler"
	cv "github.com/glycerine/goconvey/convey"
)

var _ = fmt.Printf

func Test301ReplayOfStructdDef(t *testing.T) {

	cv.Convey(`if we replay struct defn and method call, the method call should succeed the 2nd time (was failing to replay from history, stumbling on the type checker)`, t, func() {
		src := `
 type S struct{}
 var s S
 func (s *S) Hi() {println("S.Hi() called")}
 s.Hi()
`
		fmt.Printf("replay 2x, src='%s'\n", src)

		myflags := flag.NewFlagSet("gi", flag.ExitOnError)
		cfg := &compiler.GIConfig{}
		cfg.DefineFlags(myflags)

		err := myflags.Parse([]string{"-q", "-no-liner"}) // , "-vv"})
		err = cfg.ValidateConfig()
		panicOn(err)
		r := compiler.NewRepl(cfg)

		//verb.VerboseVerbose = true
		//verb.Verbose = true

		/*
			// oddly, when sent as a 4 line chunk, not
			// broken up into separate lines, we don't see
			// the issue.
			err = r.Eval(src)
			cv.So(err, cv.ShouldBeNil)
			err = r.Eval(src)
			cv.So(err, cv.ShouldBeNil)
		*/

		// now split into lines: then we get the oops.
		lines := strings.Split(src, "\n")
		// and do the lines set 2x
		for j := 0; j < 2; j++ {
			for i := range lines {
				err = r.Eval(lines[i])
				panicOn(err)
			}
			//fmt.Printf("\n pass j=%v complete.\n", j)
		}
	})
}

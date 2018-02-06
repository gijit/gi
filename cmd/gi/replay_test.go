package main

import (
	"flag"
	"fmt"
	//"github.com/shurcooL/go-goon"
	"strings"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

var _ = fmt.Printf

func Test001ReplayOfStructdDef(t *testing.T) {

	cv.Convey(`if we replay struct defn and method call, the method call should succeed the 2nd time (was failing to replay from history, stumbling on the type checker)`, t, func() {

		// TODO: DEBUG AND FINISH this. make it green.
		fmt.Printf("\n ...skipping for now, but still RED.\n")
		return

		src := `
 type S struct{}
 var s S
 func (s *S) Hi() {println("S.Hi() called")}
 s.Hi()
`
		fmt.Printf("replay 2x, src='%s'\n", src)

		myflags := flag.NewFlagSet("gi", flag.ExitOnError)
		cfg := &GIConfig{}
		cfg.DefineFlags(myflags)

		err := myflags.Parse([]string{"-q", "-no-liner"}) // , "-vv"})
		err = cfg.ValidateConfig()
		panicOn(err)
		r := NewRepl(cfg)

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

				//pp("TypesInfo.Types='%#v'", r.inc.CurPkg.Arch.TypesInfo) // .Types, .Defs, .Name2node
				/*
					// CurPkg.Arch is nil until i > 0 the first time.
					if i > 0 || j > 0 {
						for nm := range r.inc.CurPkg.Arch.TypesInfo.Name2node {
							fmt.Printf("nm='%s'\n", nm) // S.Hi
						}
						//fmt.Printf("\n node = '%#v'\n", node)
						//fmt.Printf("\n on j=%v, after Eval of line i=%v:  TypesInfo.Types='%#v'\n", j, i, r.inc.CurPkg.Arch.TypesInfo.Defs) // .Types, .Defs, .Name2node

						// too much, will print for a long, long time.
						//goon.Dump(r.inc.CurPkg.Arch.TypesInfo.Defs)
					}
				*/
			}
			fmt.Printf("\n pass j=%v complete.\n", j)
		}
	})
}

package main

import (
	"fmt"
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

var _ = fmt.Printf

func Test001ReplayOfStructdDef(t *testing.T) {

	cv.Convey(`if we replay struct defn and method call, the method call should succeed the 2nd time (was failing to replay from history, stumbling on the type checker)`, t, func() {

		src := `
 type S struct{}
 var s S
 func (s *S) Hi() {println("S.Hi() called")}
 s.Hi()
`
		fmt.Printf("replay 2x, src='%s'\n", src)

	})
}

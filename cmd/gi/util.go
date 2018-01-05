package main

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/go-interpreter/gi/pkg/compiler"
)

func translateAndCatchPanic(inc *compiler.IncrState, src []byte) (translation string, err error) {
	defer func() {
		recov := recover()
		if recov != nil {
			err = fmt.Errorf("caught panic: '%v'\n%s\n", recov, string(debug.Stack()))
		}
	}()
	translation = string(inc.Tr([]byte(src)))
	t2 := strings.TrimSpace(translation)
	nt2 := len(t2)
	if nt2 > 0 {
		if t2[nt2-1] == '\n' {
			t2 = t2[:nt2-1]
		}
	}
	p("go:'%s'  -->  '%s'\n", src, t2)
	return
}

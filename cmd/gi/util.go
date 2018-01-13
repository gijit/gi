package main

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/gijit/gi/pkg/compiler"
	"github.com/gijit/gi/pkg/verb"
)

func translateAndCatchPanic(inc *compiler.IncrState, src []byte) (translation string, err error) {
	defer func() {
		recov := recover()
		if recov != nil {
			msg := fmt.Sprintf("problem detected during Go static type checking: '%v'", recov)
			if verb.Verbose {
				msg += fmt.Sprintf("\n%s\n", string(debug.Stack()))
			}
			err = fmt.Errorf(msg)
		}
	}()
	pp("about to translate Go source '%s'", string(src))
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

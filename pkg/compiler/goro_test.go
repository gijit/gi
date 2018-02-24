package compiler

import (
	//"fmt"
	"testing"

	//"github.com/gijit/gi/pkg/token"
	//"github.com/gijit/gi/pkg/types"
	cv "github.com/glycerine/goconvey/convey"
	//"github.com/glycerine/luar"
)

func Test700StartGoroutine(t *testing.T) {

	cv.Convey(`start a new goroutine that gets its own *golua.State`, t, func() {

		r, err := NewGoro(nil)
		panicOn(err)

		//r2, err := NewGoro(nil)
		//panicOn(err)

		t0 := r.newTicket()

		ch := make(chan int, 1)
		ch <- 57

		t0.regmap["ch"] = ch

		code := `a := <- ch;`

		inc := NewIncrState(r.vm, nil)
		translation, err := inc.Tr([]byte(code))
		panicOn(err)

		pp("translation='%s'", string(translation))

		t0.run = translation

		t0.varname["a"] = true

		// execute the `a := <- ch;`
		panicOn(t0.Do())

		ai := t0.varname["a"].(int)

		cv.So(ai, cv.ShouldEqual, 57)

	})
}

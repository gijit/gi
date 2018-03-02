package front

import (
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test001EOF_versus_syntaxError_versus_complete_statement(t *testing.T) {

	cv.Convey("`more` parse should accept '3 * 4' as a compete statement, no syntax error, no further input needed", t, func() {

		var eof, syntaxErr, empty bool

		pp("complete statement: 3 * 4")

		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte("3 * 4"))
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)
		cv.So(empty, cv.ShouldBeFalse)
	})
}

func Test002MoreVsSyntaxErr(t *testing.T) {

	cv.Convey("`more` parse should ask for more input on '3 *'", t, func() {

		var eof, syntaxErr, empty bool

		pp("0th test for EOF: '3 * '")

		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte("3 * "))
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)
		cv.So(empty, cv.ShouldBeFalse)
	})
}

func Test003MoreVsSyntaxErr(t *testing.T) {

	cv.Convey("`more` parse should give more input on '3 * *' because this could be the start of multiplication times a pointer dereference", t, func() {

		var eof, syntaxErr, empty bool

		pp("syntax error: '3 * *'")

		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte("3 * *"))
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)
		cv.So(empty, cv.ShouldBeFalse)
	})
}

func Test004MoreVsSyntaxErr(t *testing.T) {

	cv.Convey("`more` parse:  '3' is complete by itself, no more needed, no syntax error", t, func() {

		var eof, syntaxErr, empty bool

		pp("complete statement: 3")

		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte("3"))
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)
		cv.So(empty, cv.ShouldBeFalse)
	})
}

func Test005MoreVsSyntaxErr(t *testing.T) {

	cv.Convey("`more` parse:  '' empty statement should indicate empty string, not more needed, no syntax error", t, func() {

		var eof, syntaxErr, empty bool

		pp("empty statement: ``")

		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte(""))
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)
		cv.So(empty, cv.ShouldBeTrue)
	})
}

func Test006MoreVsSyntaxErr(t *testing.T) {

	cv.Convey("`more` parse:  '3 4' at top level is a syntax error", t, func() {

		var eof, syntaxErr, empty bool

		pp("should be a syntax error: '3 4'")

		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte("3 4"))
		cv.So(syntaxErr, cv.ShouldBeTrue)
		cv.So(eof, cv.ShouldBeFalse)
		cv.So(empty, cv.ShouldBeFalse)
	})
}

func Test007MoreVsSyntaxErr(t *testing.T) {

	cv.Convey("`more` parse:  '2 / 3 * newline ++' is a syntax error ", t, func() {

		var eof, syntaxErr, empty bool

		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte("2 / 3 * \n ++"))
		cv.So(syntaxErr, cv.ShouldBeTrue)
		cv.So(eof, cv.ShouldBeFalse)
		cv.So(empty, cv.ShouldBeFalse)
	})
}

func Test008MoreVsSyntaxErr(t *testing.T) {

	cv.Convey("`more` parse:  '2 / ++' is a syntax error ", t, func() {

		var eof, syntaxErr, empty bool

		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte("2 / ++"))
		cv.So(syntaxErr, cv.ShouldBeTrue)
		cv.So(eof, cv.ShouldBeFalse)
		cv.So(empty, cv.ShouldBeFalse)
	})
}

func Test009MoreVsSyntaxErr(t *testing.T) {

	cv.Convey("`more` parse:  '2 / 3 *' needs more input, no syntax error", t, func() {

		var eof, syntaxErr, empty bool

		pp("2 / 3 *")

		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte(`2 / 3 * `))
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)
		cv.So(empty, cv.ShouldBeFalse)
	})
}

func Test010MoreVsSyntaxErr(t *testing.T) {

	cv.Convey("`more` parse:  '2 / 3 * newline 1' is complete by itself, no more needed, no syntax error", t, func() {

		var eof, syntaxErr, empty bool

		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte("2 / 3 * \n 1"))
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)
		cv.So(empty, cv.ShouldBeFalse)
	})
}

func Test011MoreVsSyntaxErr(t *testing.T) {

	cv.Convey("`more` parse:  '2 / 3 * newline' needs more input, no syntax error", t, func() {

		var eof, syntaxErr, empty bool

		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte("2 / 3 * \n "))
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)
		cv.So(empty, cv.ShouldBeFalse)
	})
}

func Test012MoreVsSyntaxErr(t *testing.T) {

	cv.Convey("`more` parse:  '2 / 3 * newline *' is a need more input error -- since the last * could be the start of a dereference", t, func() {

		var eof, syntaxErr, empty bool

		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte("2 / 3 * \n *"))
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)
		cv.So(empty, cv.ShouldBeFalse)
	})
}

func Test013MoreVsSyntaxErr(t *testing.T) {

	cv.Convey("multiple assignment a, b := 3, 4 is complete, no more input needed, no syntax err", t, func() {

		var eof, syntaxErr, empty bool

		pp("multiple assign: `a,b := 3,4`")

		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte("a, b := 3, 4"))
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)
		cv.So(empty, cv.ShouldBeFalse)
	})
}

func Test014MoreVsSyntaxErr(t *testing.T) {

	cv.Convey("`more` parse:  'a,b := 3, \n 4' is complete by itself, no more needed, no syntax error", t, func() {

		var eof, syntaxErr, empty bool

		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte("a,b :=3, \n 4"))
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)
		cv.So(empty, cv.ShouldBeFalse)
	})
}

func Test015MoreVsSyntaxErr(t *testing.T) {

	cv.Convey("`more` parse:  'a,b := 3,' needs more input, no syntax error", t, func() {

		var eof, syntaxErr, empty bool

		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte("a,b :=3,"))
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)
		cv.So(empty, cv.ShouldBeFalse)
	})
}

func Test016MoreVsSyntaxErr(t *testing.T) {

	cv.Convey("`more` parse:  'a,b := newline ,' is a syntax error", t, func() {

		var eof, syntaxErr, empty bool

		pp("expect EOF: `a,b := newline ,`")

		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte("a,b := \n ,"))
		cv.So(syntaxErr, cv.ShouldBeTrue)
		cv.So(eof, cv.ShouldBeFalse)
		cv.So(empty, cv.ShouldBeFalse)

	})
}

func Test017MoreVsSyntaxErr(t *testing.T) {

	cv.Convey("`more` parse:  'type A struct {'  needs more input", t, func() {

		var eof, syntaxErr, empty bool

		pp("expect need more: `type A struct {`")

		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte("type A struct {"))
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)
		cv.So(empty, cv.ShouldBeFalse)

	})
}

func Test018MoreVsSyntaxErr(t *testing.T) {

	cv.Convey("`more` parse:  'func f { nl x float nl y float64, nl } (float64,'  needs more input", t, func() {

		var eof, syntaxErr, empty bool

		src := `func f (
    x float64,
    y float64,
    ) (float64,`

		pp("expect need more: `%s`", src)
		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte(src))
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)
		cv.So(empty, cv.ShouldBeFalse)

	})
}

func Test019MoreVsSyntaxErr(t *testing.T) {

	cv.Convey("more input examples:  'switch {'  needs more "+
		"input. Same with 'for {', 'if true {', and 'select {'", t, func() {

		var eof, syntaxErr, empty bool

		src := `switch {`
		pp("expect need more: `%s`", src)
		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte(src))
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)
		cv.So(empty, cv.ShouldBeFalse)

		src = `select {`
		pp("expect need more: `%s`", src)
		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte(src))
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)
		cv.So(empty, cv.ShouldBeFalse)

		src = `if true {`
		pp("expect need more: `%s`", src)
		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte(src))
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)
		cv.So(empty, cv.ShouldBeFalse)

		src = `for {`
		pp("expect need more: `%s`", src)
		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte(src))
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)
		cv.So(empty, cv.ShouldBeFalse)

	})
}

func Test020MoreVsSyntaxErr(t *testing.T) {

	cv.Convey("multiline raw strings need more input", t, func() {

		var eof, syntaxErr, empty bool

		src := " a := `"
		pp("expect need more: `%s`", src)
		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte(src))
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)
		cv.So(empty, cv.ShouldBeFalse)
	})
}

func Test021AssignmenThoughPointer(t *testing.T) {

	cv.Convey("After `a:=1; b:= &a; `, `*b = 33` must be allowed", t, func() {

		var eof, syntaxErr, empty bool

		src := `a:=1; b:= &a;`
		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte(src))
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)
		cv.So(empty, cv.ShouldBeFalse)

		src2 := "*b = 33"
		pp("parsing '*b = 33'")
		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte(src2))
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)
		cv.So(empty, cv.ShouldBeFalse)
	})
}

func Test022GoroutineLaunchAnonFunc(t *testing.T) {

	cv.Convey("`go func() {` should return eof and ask for more input", t, func() {

		var eof, syntaxErr, empty bool

		src := `go func() {`
		eof, syntaxErr, empty, _ = TopLevelParseGoSource([]byte(src))
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)
		cv.So(empty, cv.ShouldBeFalse)
	})
}

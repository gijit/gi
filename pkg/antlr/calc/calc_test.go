package calc

import (
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test001EOF_versus_syntaxError_versus_complete_statement(t *testing.T) {

	cv.Convey("Antlr grammar processing should distinguish between these 3: a complete statement, early EOF, and syntax error", t, func() {

		var eof, syntaxErr bool

		pp("complete statement: 3")

		eof, syntaxErr = TopLevelParseSquibbleSource("3")
		cv.So(syntaxErr, cv.ShouldBeFalse) // failing here
		cv.So(eof, cv.ShouldBeFalse)

		pp("complete statement: 3 * 4")

		eof, syntaxErr = TopLevelParseSquibbleSource("3 * 4")
		cv.So(syntaxErr, cv.ShouldBeFalse) // failing here
		cv.So(eof, cv.ShouldBeFalse)

		pp("empty statement: ``")

		eof, syntaxErr = TopLevelParseSquibbleSource("")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)

		pp("complete statement: 3")

		eof, syntaxErr = TopLevelParseSquibbleSource("3")
		cv.So(syntaxErr, cv.ShouldBeFalse) // original Gofront parser doesn't think this is enough.
		cv.So(eof, cv.ShouldBeFalse)

		// this one may be too hard, let the go parser catch it.
		pp("syntax error 0: 3 * *")

		eof, syntaxErr = TopLevelParseSquibbleSource("3 * *")
		cv.So(syntaxErr, cv.ShouldBeTrue)
		cv.So(eof, cv.ShouldBeFalse)

		pp("syntax error 0: 3 4")

		eof, syntaxErr = TopLevelParseSquibbleSource("3 4")
		cv.So(syntaxErr, cv.ShouldBeTrue)
		cv.So(eof, cv.ShouldBeFalse)

		pp("syntax error 0: 3 newline 4")

		eof, syntaxErr = TopLevelParseSquibbleSource("3 \n 4")
		cv.So(syntaxErr, cv.ShouldBeTrue)
		cv.So(eof, cv.ShouldBeFalse)

		pp("test of '3 *' should give eof")

		eof, syntaxErr = TopLevelParseSquibbleSource("3 * ")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)

		pp("newline test")

		eof, syntaxErr = TopLevelParseSquibbleSource("3 * \n 4")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)

		pp("first test")

		eof, syntaxErr = TopLevelParseSquibbleSource(`3 * 4`)
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)

		pp("2nd test")

		eof, syntaxErr = TopLevelParseSquibbleSource(`2 / 3 * `)
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)

		eof, syntaxErr = TopLevelParseSquibbleSource("2 / 3 * \n 1")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)

		eof, syntaxErr = TopLevelParseSquibbleSource("2 / 3 * \n ")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)

		eof, syntaxErr = TopLevelParseSquibbleSource("2 / 3 * \n *")
		cv.So(syntaxErr, cv.ShouldBeTrue)
		cv.So(eof, cv.ShouldBeFalse)
	})
}

func Test002EOF_versus_syntaxError_versus_complete_statement(t *testing.T) {

	cv.Convey("complete statement, early EOF, and syntax error distinction: should work for multiple assignment a,b := 3 newline 4; and for func a(\n\n", t, func() {

		var eof, syntaxErr bool

		pp("multiple assign: `a,b := 3,4`")

		eof, syntaxErr = TopLevelParseSquibbleSource("a, b := 3, 4")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)

		eof, syntaxErr = TopLevelParseSquibbleSource("a,b :=3, \n 4")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)

		eof, syntaxErr = TopLevelParseSquibbleSource("a,b :=3,")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)

		eof, syntaxErr = TopLevelParseSquibbleSource("a,b := \n ,")
		cv.So(syntaxErr, cv.ShouldBeTrue)
		cv.So(eof, cv.ShouldBeFalse)

	})
}

func Test003_focused_on_one_EOF_versus_syntaxError_versus_complete_statement(t *testing.T) {

	cv.Convey("Antlr grammar processing should distinguish between these 3: a complete statement, early EOF, and syntax error", t, func() {

		var eof, syntaxErr bool

		pp("test of '3 *' should give eof")

		eof, syntaxErr = TopLevelParseSquibbleSource("3 * ")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)

	})
}

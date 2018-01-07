package calc

import (
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test001EOF_versus_syntaxError_versus_complete_statement(t *testing.T) {

	cv.Convey("Antlr grammar processing should distinguish between these 3: a complete statement, early EOF, and syntax error", t, func() {

		var eof, syntaxErr bool

		pp("empty statement: ``")

		eof, syntaxErr = TopLevelParseQuibbleSource("")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)

		pp("complete statement: 3")

		eof, syntaxErr = TopLevelParseQuibbleSource("3")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)

		pp("complete statement: 3 * 4")

		eof, syntaxErr = TopLevelParseQuibbleSource("3 * 4")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)

		pp("syntax error 0: 3 * *")

		eof, syntaxErr = TopLevelParseQuibbleSource("3 * *")
		cv.So(syntaxErr, cv.ShouldBeTrue)
		cv.So(eof, cv.ShouldBeFalse)

		pp("syntax error 0: 3 4")

		eof, syntaxErr = TopLevelParseQuibbleSource("3 4")
		cv.So(syntaxErr, cv.ShouldBeTrue)
		cv.So(eof, cv.ShouldBeFalse)

		pp("syntax error 0: 3 newline 4")

		eof, syntaxErr = TopLevelParseQuibbleSource("3 \n 4")
		cv.So(syntaxErr, cv.ShouldBeTrue)
		cv.So(eof, cv.ShouldBeFalse)

		pp("0th test: '3 * ' should give eof, not syntax error.")

		eof, syntaxErr = TopLevelParseQuibbleSource("3 * ")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)

		pp("newline test")

		eof, syntaxErr = TopLevelParseQuibbleSource("3 * \n 4")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)

		pp("first test")

		eof, syntaxErr = TopLevelParseQuibbleSource(`3 * 4`)
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)

		pp("2nd test")

		eof, syntaxErr = TopLevelParseQuibbleSource(`2 / 3 * `)
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)

		eof, syntaxErr = TopLevelParseQuibbleSource("2 / 3 * \n 1")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)

		eof, syntaxErr = TopLevelParseQuibbleSource("2 / 3 * \n ")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)

		eof, syntaxErr = TopLevelParseQuibbleSource("2 / 3 * \n *")
		cv.So(syntaxErr, cv.ShouldBeTrue)
		cv.So(eof, cv.ShouldBeFalse)
	})
}

func Test002EOF_versus_syntaxError_versus_complete_statement(t *testing.T) {

	cv.Convey("complete statement, early EOF, and syntax error distinction: should work for multiple assignment a,b := 3 newline 4; and for func a(\n\n", t, func() {

		var eof, syntaxErr bool

		pp("multiple assign: `a,b := 3,4`")

		eof, syntaxErr = TopLevelParseQuibbleSource("a, b := 3, 4")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)

		eof, syntaxErr = TopLevelParseQuibbleSource("a,b :=3, \n 4")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)

		pp("expect EOF: `a,b := 3,`")

		eof, syntaxErr = TopLevelParseQuibbleSource("a,b :=3,") // why failing now?
		cv.So(syntaxErr, cv.ShouldBeFalse)                      // getting true but should be getting false.
		cv.So(eof, cv.ShouldBeTrue)

		eof, syntaxErr = TopLevelParseQuibbleSource("a,b := \n ,")
		cv.So(syntaxErr, cv.ShouldBeTrue)
		cv.So(eof, cv.ShouldBeFalse)

	})
}

func Test004_EOF_versus_syntaxError_versus_complete_statement(t *testing.T) {

	cv.Convey("eof vs syntax error: '3 * ' should register as EOF so we request more data from the user", t, func() {

		pp("0th test: '3 * ' should give eof, not syntax error.")

		var eof, syntaxErr bool

		eof, syntaxErr = TopLevelParseQuibbleSource("3 * ")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)
	})
}

package more

import (
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test001EOF_versus_syntaxError_versus_complete_statement(t *testing.T) {

	cv.Convey("`more` parse should distinguish between these 3: a complete statement, early EOF, and syntax error", t, func() {

		var eof, syntaxErr bool

		pp("syntax error 0: 3 * *")

		eof, syntaxErr = TopLevelParseWorksSource("3 * *")
		cv.So(syntaxErr, cv.ShouldBeTrue)
		cv.So(eof, cv.ShouldBeFalse)

		pp("0th test for EOF: '3 * '")

		eof, syntaxErr = TopLevelParseWorksSource("3 * ")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)

		pp("complete statement: 3")

		eof, syntaxErr = TopLevelParseWorksSource("3")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)

		pp("empty statement: ``")

		eof, syntaxErr = TopLevelParseWorksSource("")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)

		pp("complete statement: 3 * 4")

		eof, syntaxErr = TopLevelParseWorksSource("3 * 4")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)

		pp("syntax error 0: 3 4")

		eof, syntaxErr = TopLevelParseWorksSource("3 4")
		cv.So(syntaxErr, cv.ShouldBeTrue)
		cv.So(eof, cv.ShouldBeFalse)

		pp("newline test")

		eof, syntaxErr = TopLevelParseWorksSource("3 * \n 4")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)

		pp("first test")

		eof, syntaxErr = TopLevelParseWorksSource(`3 * 4`)
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)

		pp("2nd test")

		eof, syntaxErr = TopLevelParseWorksSource(`2 / 3 * `)
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)

		eof, syntaxErr = TopLevelParseWorksSource("2 / 3 * \n 1")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)

		eof, syntaxErr = TopLevelParseWorksSource("2 / 3 * \n ")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)

		eof, syntaxErr = TopLevelParseWorksSource("2 / 3 * \n *")
		cv.So(syntaxErr, cv.ShouldBeTrue)
		cv.So(eof, cv.ShouldBeFalse)
	})
}

func Test002EOF_versus_syntaxError_versus_complete_statement(t *testing.T) {

	cv.Convey("complete statement, early EOF, and syntax error distinction: should work for multiple assignment a,b := 3 newline 4; and for func a(\n\n", t, func() {

		var eof, syntaxErr bool

		pp("multiple assign: `a,b := 3,4`")

		eof, syntaxErr = TopLevelParseWorksSource("a, b := 3, 4")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)

		eof, syntaxErr = TopLevelParseWorksSource("a,b :=3, \n 4")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)

		eof, syntaxErr = TopLevelParseWorksSource("a,b :=3,")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)

		pp("expect EOF: `a,b := 3,`")

		eof, syntaxErr = TopLevelParseWorksSource("a,b := \n ,")
		cv.So(syntaxErr, cv.ShouldBeTrue)
		cv.So(eof, cv.ShouldBeFalse)

	})
}

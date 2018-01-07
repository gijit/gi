package multigo

import (
	"testing"

	cv "github.com/glycerine/goconvey/convey"
)

func Test001FunctionAsAWhole(t *testing.T) {

	cv.Convey("Given that we are parsing successive lines of input from the gi repl, when we encounter the simple function definition `func a() { (newline) }` we should gather the 2nd line before returning the text of the whole function, thus enabling multiline input\n\n", t, func() {

		var eof, syntaxErr bool

		eof, syntaxErr = TopLevelParseGoSource(`a, b := 3, 4`)
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)

		eof, syntaxErr = TopLevelParseGoSource(`a := 3`)
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)

		eof, syntaxErr = TopLevelParseGoSource(`a, b := 3,
     4
`)
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeFalse)

		eof, syntaxErr = TopLevelParseGoSource(`a, b := 3,`)
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)

		eof, syntaxErr = TopLevelParseGoSource("func f( a int,")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)

		eof, syntaxErr = TopLevelParseGoSource("func f(")
		cv.So(syntaxErr, cv.ShouldBeFalse)
		cv.So(eof, cv.ShouldBeTrue)

		eof, syntaxErr = TopLevelParseGoSource(`
func hello(
	a string,
	b string,
) {
}
`)
		cv.So(eof, cv.ShouldBeFalse)
		cv.So(syntaxErr, cv.ShouldBeFalse)

		eof, syntaxErr = TopLevelParseGoSource("func f() {")
		cv.So(eof, cv.ShouldBeTrue)
		cv.So(syntaxErr, cv.ShouldBeFalse)

		eof, syntaxErr = TopLevelParseGoSource("func f() {\n}")
		cv.So(eof, cv.ShouldBeFalse)
		cv.So(syntaxErr, cv.ShouldBeFalse)

		eof, syntaxErr = TopLevelParseGoSource("func a f() {\n}")
		cv.So(eof, cv.ShouldBeFalse)
		cv.So(syntaxErr, cv.ShouldBeTrue)
	})
}

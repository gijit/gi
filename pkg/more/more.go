// accept multiline input, only feeding
// it to the interpreter when we have
// a full statement/declaration/expression.
package more

import (
	"github.com/go-interpreter/gi/pkg/more/parser"
	"github.com/go-interpreter/gi/pkg/more/token"
	"github.com/go-interpreter/gi/pkg/verb"
)

var pp = verb.PP

func TopLevelParseWorksSource(sourceCode string) (eofSeen, errorSeen bool) {

	fileSet := token.NewFileSet() // positions are relative to fileSet
	_, err := parser.ParseFile(fileSet, "", sourceCode, parser.Trace)
	if err == parser.ErrSyntax {
		errorSeen = true
		err = nil
	}
	if err == parser.ErrMoreInput {
		eofSeen = true
		err = nil
	}
	panicOn(err)
	pp("we got past the ParseFile !")

	// set eofSeen, errorSeen !!
	return
}

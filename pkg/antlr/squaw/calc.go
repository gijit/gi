// accept multiline input, only feeding
// it to the interpreter when we have
// a full statement/declaration/expression.
package calc

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	pg "github.com/go-interpreter/gi/pkg/antlr/squaw/parse_calc"
)

func TopLevelParseQuibbleSource(sourceCode string) (eofSeen, errorSeen bool) {

	input := antlr.NewInputStream(sourceCode)
	lexer := pg.NewQuibbleLexer(input)

	tokenStream := antlr.NewCommonTokenStream(lexer, 0)
	p := pg.NewQuibbleParser(tokenStream)

	tl := antlr.NewTraceListener(p.BaseParser)
	p.SetTrace(tl)

	es := p.NewQuibbleErrorStrategy()
	p.SetErrorHandler(es)

	lsn := p.InitOurErrorHandler(es)
	_ = lsn

	defer func() {
		eofSeen = es.EofSeen
		errorSeen = es.ErrorSeen
		recov := recover()
		switch recov {
		case nil:
		case pg.ErrReplSyntax:
		case pg.ErrReplEOF:
		default:
			panic(recov)
		}
	}()

	// all the parsing actually happens here, with none
	// of the listener callbacks.
	p.ReplStuff()

	return es.EofSeen, es.ErrorSeen
}

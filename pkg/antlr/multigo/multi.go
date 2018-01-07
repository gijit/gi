// accept multiline input, only feeding
// it to the interpreter when we have
// a full statement/declaration/expression.
package multigo

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	pg "github.com/go-interpreter/gi/pkg/antlr/parse_gofront"
)

func TopLevelParseGoSource(sourceCode string) (eofSeen, errorSeen bool) {

	input := antlr.NewInputStream(sourceCode)
	lexer := pg.NewGofrontLexer(input)

	tokenStream := antlr.NewCommonTokenStream(lexer, 0)
	p := pg.NewGofrontParser(tokenStream)

	es := p.NewGofrontErrorStrategy()
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

	p.ReplStuff()

	return es.EofSeen, es.ErrorSeen
}

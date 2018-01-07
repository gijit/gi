package multigo

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	pg "github.com/go-interpreter/gi/pkg/antlr/parse_gofront"
)

var _ antlr.ParserRuleContext

type ReplListener struct {
	*pg.BaseGofrontListener
}

func NewReplListener(packageName string) *ReplListener {
	return &ReplListener{}
}

// EveryRule is useful for debugging, but too verbose.
/*
func (s *ReplListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	fmt.Printf("EnterEveryRule: %s\n", ctx.GetText())
}

func (s *ReplListener) ExitEveryRule(ctx antlr.ParserRuleContext) {
	fmt.Printf("ExitEveryRule: %s\n", ctx.GetText())
}


func (s *ReplListener) VisitErrorNode(node antlr.ErrorNode) {
	fmt.Printf("VisitErrorNode called\n")
}

func (s *ReplListener) EnterExpression(c *pg.ExpressionContext) {
	fmt.Printf("EnterExpression: %s\n", c.GetText())
}

func (s *ReplListener) ExitExpression(c *pg.ExpressionContext) {
	fmt.Printf("ExitExpression: %s\n", c.GetText())
}

func (s *ReplListener) ExitReplStuff(ctx *pg.ReplStuffContext) {
	s.replDepth--
	fmt.Printf("ExitReplStuff at depth %v: %s\n", s.replDepth, ctx.GetText())
	if s.replDepth == 0 {
		fmt.Printf("exprDepth = %v, we have a full expression", s.replDepth)
	}
}

func (s *ReplListener) EnterReplStuff(ctx *pg.ReplStuffContext) {
	s.replDepth++
	fmt.Printf("EnterReplStuff at depth %v: %s\n", s.replDepth, ctx.GetText())
}
*/

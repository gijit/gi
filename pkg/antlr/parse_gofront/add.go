package parser // Gofront

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"os"
	"strconv"
)

var ErrReplSyntax = fmt.Errorf("syntax error, stopping early.")
var ErrReplEOF = fmt.Errorf("EOF encountered before parse finished.")

func (s *GofrontParser) InitOurErrorHandler(strat *GofrontErrorStrategy) *GofrontErrorListener {
	e := NewGofrontErrorListener(strat)
	//s.RemoveErrorListeners()
	s.AddErrorListener(e)
	return e
}

type GofrontErrorStrategy struct {
	*antlr.DefaultErrorStrategy
	EofSeen   bool
	ErrorSeen bool
}

func (s *GofrontParser) NewGofrontErrorStrategy() *GofrontErrorStrategy {
	return &GofrontErrorStrategy{
		DefaultErrorStrategy: antlr.NewDefaultErrorStrategy(),
	}
}

func (d *GofrontErrorStrategy) RecoverInline(recognizer antlr.Parser) antlr.Token {
	tok := recognizer.GetCurrentToken()
	if d.EofSeen || d.ErrorSeen {
		return tok
	}
	fmt.Printf("GofrontErrorStrategy.RecoverInline called!\n")

	o := recognizer.GetCurrentToken()
	if o.GetTokenType() == antlr.TokenEOF {
		d.EofSeen = true
		panic(ErrReplEOF)
	}

	d.ErrorSeen = true
	panic(ErrReplSyntax)
	return tok
}

func (d *GofrontErrorStrategy) ReportUnwantedToken(recognizer antlr.Parser) {
	if d.EofSeen || d.ErrorSeen {
		return
	}
	fmt.Printf("GofrontErrorStrategy.ReportUnwantedToken() called!\n")
	d.ErrorSeen = true
	panic(ErrReplSyntax)

	return
}

func (d *GofrontErrorStrategy) Sync(recognizer antlr.Parser) {
	//fmt.Printf("GofrontErrorStrategy.Sync() called\n")
}

func (d *GofrontErrorStrategy) ReportError(recognizer antlr.Parser, e antlr.RecognitionException) {
	fmt.Printf("GofrontErrorStrategy.ReportError() called.\n")

	o := recognizer.GetCurrentToken()
	if o.GetTokenType() == antlr.TokenEOF {
		d.EofSeen = true
		panic(ErrReplEOF)
	}

	if d.EofSeen || d.ErrorSeen {
		return
	}
	d.ErrorSeen = true
	panic(ErrReplSyntax)

}

func (d *GofrontErrorStrategy) ReportMatch(recognizer antlr.Parser) {
	fmt.Printf("GofrontErrorStrategy.ReportMatch() called.\n")
}

type GofrontErrorListener struct {
	*antlr.DefaultErrorListener
	strat *GofrontErrorStrategy
}

func NewGofrontErrorListener(strat *GofrontErrorStrategy) *GofrontErrorListener {
	return &GofrontErrorListener{
		strat: strat,
	}
}

func (c *GofrontErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	// msg can be: oops-antlr-missing-token: missing '}' at '<EOF>'
	fmt.Fprintln(os.Stderr, "YOWZA! SyntaxError here: line "+strconv.Itoa(line)+":"+strconv.Itoa(column)+" "+msg)
	c.strat.ErrorSeen = true
	panic(ErrReplSyntax)
}

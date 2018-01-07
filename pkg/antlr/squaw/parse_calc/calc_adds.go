package parser // Squibble.g4

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"os"
	"strconv"
)

var ErrReplSyntax = fmt.Errorf("syntax error, stopping early.")
var ErrReplEOF = fmt.Errorf("EOF encountered before parse finished.")

func (s *QuibbleParser) InitOurErrorHandler(strat *QuibbleErrorStrategy) *QuibbleErrorListener {
	e := NewQuibbleErrorListener(strat)
	//s.RemoveErrorListeners()
	s.AddErrorListener(e)
	return e
}

type QuibbleErrorStrategy struct {
	*antlr.DefaultErrorStrategy
	EofSeen   bool
	ErrorSeen bool
}

func (s *QuibbleParser) NewQuibbleErrorStrategy() *QuibbleErrorStrategy {
	return &QuibbleErrorStrategy{
		DefaultErrorStrategy: antlr.NewDefaultErrorStrategy(),
	}
}

func (d *QuibbleErrorStrategy) RecoverInline(recognizer antlr.Parser) antlr.Token {
	tok := recognizer.GetCurrentToken()
	if d.EofSeen || d.ErrorSeen {
		return tok
	}
	fmt.Printf("QuibbleErrorStrategy.RecoverInline called!\n")

	o := recognizer.GetCurrentToken()
	if o.GetTokenType() == antlr.TokenEOF {
		d.EofSeen = true
		panic(ErrReplEOF)
	}

	d.ErrorSeen = true
	panic(ErrReplSyntax)
	return tok
}

func (d *QuibbleErrorStrategy) ReportUnwantedToken(recognizer antlr.Parser) {
	if d.EofSeen || d.ErrorSeen {
		return
	}
	fmt.Printf("QuibbleErrorStrategy.ReportUnwantedToken() called!\n")
	d.ErrorSeen = true
	panic(ErrReplSyntax)

	return
}

func (d *QuibbleErrorStrategy) Sync(recognizer antlr.Parser) {
	//fmt.Printf("QuibbleErrorStrategy.Sync() called\n")
}

func (d *QuibbleErrorStrategy) ReportError(recognizer antlr.Parser, e antlr.RecognitionException) {
	fmt.Printf("QuibbleErrorStrategy.ReportError() called.\n")

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

//func (d *QuibbleErrorStrategy) ReportMatch(recognizer antlr.Parser) {
//	//fmt.Printf("QuibbleErrorStrategy.ReportMatch() called.\n")
//}

type QuibbleErrorListener struct {
	*antlr.DefaultErrorListener
	strat *QuibbleErrorStrategy
}

func NewQuibbleErrorListener(strat *QuibbleErrorStrategy) *QuibbleErrorListener {
	return &QuibbleErrorListener{
		strat: strat,
	}
}

func (c *QuibbleErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	// msg can be: oops-antlr-missing-token: missing '}' at '<EOF>'
	fmt.Fprintln(os.Stderr, "YOWZA! SyntaxError here: line "+strconv.Itoa(line)+":"+strconv.Itoa(column)+" "+msg)
	c.strat.ErrorSeen = true
	panic(ErrReplSyntax)
}

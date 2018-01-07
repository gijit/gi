package parser // Squibble.g4

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"os"
	"strconv"
)

var ErrReplSyntax = fmt.Errorf("syntax error, stopping early.")
var ErrReplEOF = fmt.Errorf("EOF encountered before parse finished.")

func (s *WorksParser) InitOurErrorHandler(strat *WorksErrorStrategy) *WorksErrorListener {
	e := NewWorksErrorListener(strat)
	//s.RemoveErrorListeners()
	s.AddErrorListener(e)
	return e
}

type WorksErrorStrategy struct {
	*antlr.DefaultErrorStrategy
	EofSeen   bool
	ErrorSeen bool
}

func (s *WorksParser) NewWorksErrorStrategy() *WorksErrorStrategy {
	return &WorksErrorStrategy{
		DefaultErrorStrategy: antlr.NewDefaultErrorStrategy(),
	}
}

func (d *WorksErrorStrategy) RecoverInline(recognizer antlr.Parser) antlr.Token {
	tok := recognizer.GetCurrentToken()
	if d.EofSeen || d.ErrorSeen {
		return tok
	}
	fmt.Printf("WorksErrorStrategy.RecoverInline called!\n")

	o := recognizer.GetCurrentToken()
	if o.GetTokenType() == antlr.TokenEOF {
		d.EofSeen = true
		panic(ErrReplEOF)
	}

	d.ErrorSeen = true
	panic(ErrReplSyntax)
	return tok
}

func (d *WorksErrorStrategy) ReportUnwantedToken(recognizer antlr.Parser) {
	if d.EofSeen || d.ErrorSeen {
		return
	}
	fmt.Printf("WorksErrorStrategy.ReportUnwantedToken() called!\n")
	d.ErrorSeen = true
	panic(ErrReplSyntax)

	return
}

func (d *WorksErrorStrategy) Sync(recognizer antlr.Parser) {
	//fmt.Printf("WorksErrorStrategy.Sync() called\n")
}

func (d *WorksErrorStrategy) ReportError(recognizer antlr.Parser, e antlr.RecognitionException) {
	fmt.Printf("WorksErrorStrategy.ReportError() called.\n")

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

//func (d *WorksErrorStrategy) ReportMatch(recognizer antlr.Parser) {
//	//fmt.Printf("WorksErrorStrategy.ReportMatch() called.\n")
//}

type WorksErrorListener struct {
	*antlr.DefaultErrorListener
	strat *WorksErrorStrategy
}

func NewWorksErrorListener(strat *WorksErrorStrategy) *WorksErrorListener {
	return &WorksErrorListener{
		strat: strat,
	}
}

func (c *WorksErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	// msg can be: oops-antlr-missing-token: missing '}' at '<EOF>'
	fmt.Fprintln(os.Stderr, "YOWZA! SyntaxError here: line "+strconv.Itoa(line)+":"+strconv.Itoa(column)+" "+msg)
	c.strat.ErrorSeen = true
	panic(ErrReplSyntax)
}

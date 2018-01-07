package parser // Squibble.g4

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"os"
	"strconv"
)

var ErrReplSyntax = fmt.Errorf("syntax error, stopping early.")
var ErrReplEOF = fmt.Errorf("EOF encountered before parse finished.")

func (s *SquibbleParser) InitOurErrorHandler(strat *SquibbleErrorStrategy) *SquibbleErrorListener {
	e := NewSquibbleErrorListener(strat)
	//s.RemoveErrorListeners()
	s.AddErrorListener(e)
	return e
}

type SquibbleErrorStrategy struct {
	*antlr.DefaultErrorStrategy
	EofSeen   bool
	ErrorSeen bool
}

func (s *SquibbleParser) NewSquibbleErrorStrategy() *SquibbleErrorStrategy {
	return &SquibbleErrorStrategy{
		DefaultErrorStrategy: antlr.NewDefaultErrorStrategy(),
	}
}

/*
https://stackoverflow.com/questions/14152978/how-can-i-use-antlr-to-make-an-interpreter-that-allows-multiple-lines

This is only a set of guidelines, but should provide a general idea.

You'll want to make sure the parser rule you're using to parse statements ends with an explicit EOF.

start : stmt EOF; // good

start : stmt;     // bad! might not parse the whole buffer!

1. After any line of input produces a complete statement which gets executed, clear the multiline input "buffer".

2. After a line of data is input from the user, add the data to the buffer, then attempt to parse the entire buffer.

2. a) If the input buffer is successfully parsed, execute the input (i.e. go to item #1 above).

2. b) If a syntax error occurs but only after the parser has consumed the EOF symbol (either during prediction or matching) from the token stream, the current input is incomplete. Do not execute anything and wait for the user to enter more text.

2. c) If a syntax error occurs before the parser has consumed the EOF symbol, the input contains a syntax error. Report and/or handle the error as you see fit (including probably clear the input buffer without executing anything).

*/
func (d *SquibbleErrorStrategy) RecoverInline(recognizer antlr.Parser) antlr.Token {
	tok := recognizer.GetCurrentToken()
	if d.EofSeen || d.ErrorSeen {
		return tok
	}

	expected := recognizer.GetExpectedTokens()
	ctx := recognizer.GetParserRuleContext()

	la1 := recognizer.GetTokenStream().LA(1)
	fmt.Printf("SquibbleErrorStrategy.RecoverInline called!.  current token = '%v'/'%v', and p.GetTokenStream().LA(1) = '%v'   ....   expected='%#v',   ctx='%#v', ctx.GetRuleIndex = '%v'\n", tok.GetTokenType(), tok.GetText(), la1, expected, ctx, ruleNames[ctx.GetRuleIndex()])
	//	panic("where")

	// "3 * " is getting here.
	// jea: try treating all RecoverInline as EOF? // no then "3 * *" is called eof.
	o := recognizer.GetCurrentToken()
	if o.GetTokenType() == antlr.TokenEOF {
		d.EofSeen = true
		fmt.Printf("see TokenEOF\n")
		panic(ErrReplEOF)
	}
	fmt.Printf("NOT TokenEOF\n")

	d.ErrorSeen = true
	panic(ErrReplSyntax)
	return tok
}

func (d *SquibbleErrorStrategy) ReportUnwantedToken(recognizer antlr.Parser) {
	if d.EofSeen || d.ErrorSeen {
		return
	}
	fmt.Printf("SquibbleErrorStrategy.ReportUnwantedToken() called!\n")
	d.ErrorSeen = true
	panic(ErrReplSyntax)

	return
}

func (d *SquibbleErrorStrategy) Sync(recognizer antlr.Parser) {
	//fmt.Printf("SquibbleErrorStrategy.Sync() called\n")
}

func (d *SquibbleErrorStrategy) ReportError(recognizer antlr.Parser, e antlr.RecognitionException) {
	fmt.Printf("SquibbleErrorStrategy.ReportError() called with e='%s', token='%s', msg='%s'.\n", e, e.GetOffendingToken(), e.GetMessage())

	tok := recognizer.GetCurrentToken()
	expected := recognizer.GetExpectedTokens()
	ctx := recognizer.GetParserRuleContext()

	la1 := recognizer.GetTokenStream().LA(1)
	ruleIndex := ctx.GetRuleIndex()
	fmt.Printf("SquibbleErrorStrategy.ReportError detail:  current token = '%v'/'%v', and p.GetTokenStream().LA(1) = '%v'   ....   expected='%#v',   ctx='%#v', ctx.GetRuleIndex = '%v', ctx.GetRuleIndex translation by rulesNames='%s'\n", tok.GetTokenType(), tok.GetText(), la1, expected, ctx, ruleIndex, ruleNames[ctx.GetRuleIndex()])

	//	if ruleNames[ruleIndex] == "replEntry" {
	//
	//		d.EofSeen = true
	//		fmt.Printf("rule is replEntry, treating this as EOF...\n")
	//		panic(ErrReplEOF)
	//	}

	o := recognizer.GetCurrentToken()
	if o.GetTokenType() == antlr.TokenEOF {
		d.EofSeen = true
		fmt.Printf("see TokenEOF\n")
		panic(ErrReplEOF)
	}
	if d.EofSeen || d.ErrorSeen {
		return
	}
	fmt.Printf("NOT TokenEOF, so setting ErrorSeen\n")
	d.ErrorSeen = true
	panic(ErrReplSyntax)

}

//func (d *SquibbleErrorStrategy) ReportMatch(recognizer antlr.Parser) {
//	//fmt.Printf("SquibbleErrorStrategy.ReportMatch() called.\n")
//}

type SquibbleErrorListener struct {
	*antlr.DefaultErrorListener
	strat *SquibbleErrorStrategy
}

func NewSquibbleErrorListener(strat *SquibbleErrorStrategy) *SquibbleErrorListener {
	return &SquibbleErrorListener{
		strat: strat,
	}
}

func (c *SquibbleErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	// msg can be: oops-antlr-missing-token: missing '}' at '<EOF>'
	fmt.Fprintln(os.Stderr, "YOWZA! SyntaxError here: line "+strconv.Itoa(line)+":"+strconv.Itoa(column)+" "+msg)
	c.strat.ErrorSeen = true
	panic(ErrReplSyntax)
}

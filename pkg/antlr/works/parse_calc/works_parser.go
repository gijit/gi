// Code generated from Works.g4 by ANTLR 4.7.1. DO NOT EDIT.

package parser // Works

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

import "strings"

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 26, 123,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 3, 2, 3, 2, 3, 2, 3, 2, 5, 2, 35, 10,
	2, 3, 3, 3, 3, 3, 3, 6, 3, 40, 10, 3, 13, 3, 14, 3, 41, 3, 4, 3, 4, 3,
	4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 5, 4, 56, 10,
	4, 3, 5, 3, 5, 3, 5, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 7, 6, 67, 10,
	6, 12, 6, 14, 6, 70, 11, 6, 3, 6, 5, 6, 73, 10, 6, 3, 7, 5, 7, 76, 10,
	7, 3, 7, 3, 7, 3, 8, 3, 8, 3, 9, 3, 9, 3, 9, 3, 9, 3, 10, 3, 10, 3, 10,
	7, 10, 89, 10, 10, 12, 10, 14, 10, 92, 11, 10, 3, 11, 3, 11, 3, 11, 7,
	11, 97, 10, 11, 12, 11, 14, 11, 100, 11, 11, 3, 12, 3, 12, 3, 12, 3, 12,
	3, 12, 3, 12, 7, 12, 108, 10, 12, 12, 12, 14, 12, 111, 11, 12, 3, 13, 3,
	13, 3, 14, 3, 14, 5, 14, 117, 10, 14, 3, 15, 3, 15, 3, 15, 3, 15, 3, 15,
	2, 3, 22, 16, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 2, 5,
	4, 2, 8, 8, 13, 13, 3, 2, 11, 12, 4, 2, 16, 19, 22, 22, 2, 122, 2, 34,
	3, 2, 2, 2, 4, 39, 3, 2, 2, 2, 6, 55, 3, 2, 2, 2, 8, 57, 3, 2, 2, 2, 10,
	60, 3, 2, 2, 2, 12, 75, 3, 2, 2, 2, 14, 79, 3, 2, 2, 2, 16, 81, 3, 2, 2,
	2, 18, 85, 3, 2, 2, 2, 20, 93, 3, 2, 2, 2, 22, 101, 3, 2, 2, 2, 24, 112,
	3, 2, 2, 2, 26, 116, 3, 2, 2, 2, 28, 118, 3, 2, 2, 2, 30, 35, 7, 3, 2,
	2, 31, 35, 7, 2, 2, 3, 32, 35, 6, 2, 2, 2, 33, 35, 6, 2, 3, 2, 34, 30,
	3, 2, 2, 2, 34, 31, 3, 2, 2, 2, 34, 32, 3, 2, 2, 2, 34, 33, 3, 2, 2, 2,
	35, 3, 3, 2, 2, 2, 36, 37, 5, 6, 4, 2, 37, 38, 5, 2, 2, 2, 38, 40, 3, 2,
	2, 2, 39, 36, 3, 2, 2, 2, 40, 41, 3, 2, 2, 2, 41, 39, 3, 2, 2, 2, 41, 42,
	3, 2, 2, 2, 42, 5, 3, 2, 2, 2, 43, 44, 5, 22, 12, 2, 44, 45, 7, 2, 2, 3,
	45, 56, 3, 2, 2, 2, 46, 47, 5, 16, 9, 2, 47, 48, 7, 2, 2, 3, 48, 56, 3,
	2, 2, 2, 49, 50, 5, 8, 5, 2, 50, 51, 7, 2, 2, 3, 51, 56, 3, 2, 2, 2, 52,
	53, 5, 10, 6, 2, 53, 54, 7, 2, 2, 3, 54, 56, 3, 2, 2, 2, 55, 43, 3, 2,
	2, 2, 55, 46, 3, 2, 2, 2, 55, 49, 3, 2, 2, 2, 55, 52, 3, 2, 2, 2, 56, 7,
	3, 2, 2, 2, 57, 58, 7, 4, 2, 2, 58, 59, 7, 13, 2, 2, 59, 9, 3, 2, 2, 2,
	60, 72, 7, 5, 2, 2, 61, 73, 5, 12, 7, 2, 62, 68, 7, 6, 2, 2, 63, 64, 5,
	12, 7, 2, 64, 65, 5, 2, 2, 2, 65, 67, 3, 2, 2, 2, 66, 63, 3, 2, 2, 2, 67,
	70, 3, 2, 2, 2, 68, 66, 3, 2, 2, 2, 68, 69, 3, 2, 2, 2, 69, 71, 3, 2, 2,
	2, 70, 68, 3, 2, 2, 2, 71, 73, 7, 7, 2, 2, 72, 61, 3, 2, 2, 2, 72, 62,
	3, 2, 2, 2, 73, 11, 3, 2, 2, 2, 74, 76, 9, 2, 2, 2, 75, 74, 3, 2, 2, 2,
	75, 76, 3, 2, 2, 2, 76, 77, 3, 2, 2, 2, 77, 78, 5, 14, 8, 2, 78, 13, 3,
	2, 2, 2, 79, 80, 7, 22, 2, 2, 80, 15, 3, 2, 2, 2, 81, 82, 5, 18, 10, 2,
	82, 83, 7, 9, 2, 2, 83, 84, 5, 20, 11, 2, 84, 17, 3, 2, 2, 2, 85, 90, 7,
	13, 2, 2, 86, 87, 7, 10, 2, 2, 87, 89, 7, 13, 2, 2, 88, 86, 3, 2, 2, 2,
	89, 92, 3, 2, 2, 2, 90, 88, 3, 2, 2, 2, 90, 91, 3, 2, 2, 2, 91, 19, 3,
	2, 2, 2, 92, 90, 3, 2, 2, 2, 93, 98, 5, 22, 12, 2, 94, 95, 7, 10, 2, 2,
	95, 97, 5, 22, 12, 2, 96, 94, 3, 2, 2, 2, 97, 100, 3, 2, 2, 2, 98, 96,
	3, 2, 2, 2, 98, 99, 3, 2, 2, 2, 99, 21, 3, 2, 2, 2, 100, 98, 3, 2, 2, 2,
	101, 102, 8, 12, 1, 2, 102, 103, 5, 24, 13, 2, 103, 109, 3, 2, 2, 2, 104,
	105, 12, 4, 2, 2, 105, 106, 9, 3, 2, 2, 106, 108, 5, 22, 12, 5, 107, 104,
	3, 2, 2, 2, 108, 111, 3, 2, 2, 2, 109, 107, 3, 2, 2, 2, 109, 110, 3, 2,
	2, 2, 110, 23, 3, 2, 2, 2, 111, 109, 3, 2, 2, 2, 112, 113, 9, 4, 2, 2,
	113, 25, 3, 2, 2, 2, 114, 117, 7, 13, 2, 2, 115, 117, 5, 28, 15, 2, 116,
	114, 3, 2, 2, 2, 116, 115, 3, 2, 2, 2, 117, 27, 3, 2, 2, 2, 118, 119, 7,
	13, 2, 2, 119, 120, 7, 8, 2, 2, 120, 121, 7, 13, 2, 2, 121, 29, 3, 2, 2,
	2, 12, 34, 41, 55, 68, 72, 75, 90, 98, 109, 116,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "';'", "'package'", "'import'", "'('", "')'", "'.'", "':='", "','",
	"'*'", "'/'",
}
var symbolicNames = []string{
	"", "", "", "", "", "", "", "", "", "", "", "IDENTIFIER", "KEYWORD", "BINARY_OP",
	"INT_LIT", "FLOAT_LIT", "IMAGINARY_LIT", "RUNE_LIT", "LITTLE_U_VALUE",
	"BIG_U_VALUE", "STRING_LIT", "WS", "COMMENT", "TERMINATOR", "LINE_COMMENT",
}

var ruleNames = []string{
	"eos", "replStuff", "replEntry", "packageClause", "importDecl", "importSpec",
	"importPath", "assignment", "identifierList", "expressionList", "expression",
	"basicLit", "operandName", "qualifiedIdent",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type WorksParser struct {
	*antlr.BaseParser
}

func NewWorksParser(input antlr.TokenStream) *WorksParser {
	this := new(WorksParser)

	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "Works.g4"

	return this
}

// here returns `true` iff on the current index of the parser's
// token stream a token of the given `type` exists on the
// `HIDDEN` channel.
//
// Args:
//  type (int): the type of the token on the `HIDDEN` channel
//      to check.
//
//  Returns:
//      `True` iff on the current index of the parser's
//      token stream a token of the given `type` exists on the
//      `HIDDEN` channel.
func (p *WorksParser) here(tokenType int) bool {
	// Get the token ahead of the current index.
	possibleIndexEosToken := p.GetCurrentToken().GetTokenIndex() - 1
	ahead := p.GetTokenStream().Get(possibleIndexEosToken)

	// Check if the token resides on the HIDDEN channel and if it is of the
	// provided tokenType.
	return (ahead.GetChannel() == antlr.LexerHidden) && (ahead.GetTokenType() == tokenType)
}

/**
 * Returns {@code true} iff on the current index of the parser's
 * token stream a token exists on the {@code HIDDEN} channel which
 * either is a line terminator, or is a multi line comment that
 * contains a line terminator.
 *
 * @return {@code true} iff on the current index of the parser's
 * token stream a token exists on the {@code HIDDEN} channel which
 * either is a line terminator, or is a multi line comment that
 * contains a line terminator.
 */
func (p *WorksParser) lineTerminatorAhead() bool {
	possibleIndexEosToken := p.GetCurrentToken().GetTokenIndex() - 1
	ahead := p.GetTokenStream().Get(possibleIndexEosToken)

	if ahead.GetChannel() != antlr.LexerHidden {
		// We're only interested in tokens on the HIDDEN channel.
		return false
	}

	if ahead.GetTokenType() == WorksParserTERMINATOR {
		// There is definitely a line terminator ahead.
		return true
	}

	if ahead.GetTokenType() == WorksParserWS {
		// Get the token ahead of the current whitespaces.
		possibleIndexEosToken = p.GetCurrentToken().GetTokenIndex() - 2
		ahead = p.GetTokenStream().Get(possibleIndexEosToken)
	}

	// Get the token's text and type.
	text := ahead.GetText()
	tokenType := ahead.GetTokenType()

	// Check if the token is, or contains a line terminator.
	if (tokenType == WorksParserCOMMENT ||
		tokenType == WorksParserLINE_COMMENT) &&
		strings.ContainsAny(text, "\r\n") {

		return true
	}

	return tokenType == WorksParserTERMINATOR
}

// WorksParser tokens.
const (
	WorksParserEOF            = antlr.TokenEOF
	WorksParserT__0           = 1
	WorksParserT__1           = 2
	WorksParserT__2           = 3
	WorksParserT__3           = 4
	WorksParserT__4           = 5
	WorksParserT__5           = 6
	WorksParserT__6           = 7
	WorksParserT__7           = 8
	WorksParserT__8           = 9
	WorksParserT__9           = 10
	WorksParserIDENTIFIER     = 11
	WorksParserKEYWORD        = 12
	WorksParserBINARY_OP      = 13
	WorksParserINT_LIT        = 14
	WorksParserFLOAT_LIT      = 15
	WorksParserIMAGINARY_LIT  = 16
	WorksParserRUNE_LIT       = 17
	WorksParserLITTLE_U_VALUE = 18
	WorksParserBIG_U_VALUE    = 19
	WorksParserSTRING_LIT     = 20
	WorksParserWS             = 21
	WorksParserCOMMENT        = 22
	WorksParserTERMINATOR     = 23
	WorksParserLINE_COMMENT   = 24
)

// WorksParser rules.
const (
	WorksParserRULE_eos            = 0
	WorksParserRULE_replStuff      = 1
	WorksParserRULE_replEntry      = 2
	WorksParserRULE_packageClause  = 3
	WorksParserRULE_importDecl     = 4
	WorksParserRULE_importSpec     = 5
	WorksParserRULE_importPath     = 6
	WorksParserRULE_assignment     = 7
	WorksParserRULE_identifierList = 8
	WorksParserRULE_expressionList = 9
	WorksParserRULE_expression     = 10
	WorksParserRULE_basicLit       = 11
	WorksParserRULE_operandName    = 12
	WorksParserRULE_qualifiedIdent = 13
)

// IEosContext is an interface to support dynamic dispatch.
type IEosContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEosContext differentiates from other interfaces.
	IsEosContext()
}

type EosContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEosContext() *EosContext {
	var p = new(EosContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = WorksParserRULE_eos
	return p
}

func (*EosContext) IsEosContext() {}

func NewEosContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EosContext {
	var p = new(EosContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = WorksParserRULE_eos

	return p
}

func (s *EosContext) GetParser() antlr.Parser { return s.parser }

func (s *EosContext) EOF() antlr.TerminalNode {
	return s.GetToken(WorksParserEOF, 0)
}

func (s *EosContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EosContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EosContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.EnterEos(s)
	}
}

func (s *EosContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.ExitEos(s)
	}
}

func (p *WorksParser) Eos() (localctx IEosContext) {
	localctx = NewEosContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, WorksParserRULE_eos)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(32)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(28)
			p.Match(WorksParserT__0)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(29)
			p.Match(WorksParserEOF)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		p.SetState(30)

		if !(p.lineTerminatorAhead()) {
			panic(antlr.NewFailedPredicateException(p, "p.lineTerminatorAhead()", ""))
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		p.SetState(31)

		if !(p.GetTokenStream().LT(1).GetText() == "}") {
			panic(antlr.NewFailedPredicateException(p, "p.GetTokenStream().LT(1).GetText() == \"}\" ", ""))
		}

	}

	return localctx
}

// IReplStuffContext is an interface to support dynamic dispatch.
type IReplStuffContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsReplStuffContext differentiates from other interfaces.
	IsReplStuffContext()
}

type ReplStuffContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyReplStuffContext() *ReplStuffContext {
	var p = new(ReplStuffContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = WorksParserRULE_replStuff
	return p
}

func (*ReplStuffContext) IsReplStuffContext() {}

func NewReplStuffContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ReplStuffContext {
	var p = new(ReplStuffContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = WorksParserRULE_replStuff

	return p
}

func (s *ReplStuffContext) GetParser() antlr.Parser { return s.parser }

func (s *ReplStuffContext) AllReplEntry() []IReplEntryContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IReplEntryContext)(nil)).Elem())
	var tst = make([]IReplEntryContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IReplEntryContext)
		}
	}

	return tst
}

func (s *ReplStuffContext) ReplEntry(i int) IReplEntryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IReplEntryContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IReplEntryContext)
}

func (s *ReplStuffContext) AllEos() []IEosContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IEosContext)(nil)).Elem())
	var tst = make([]IEosContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IEosContext)
		}
	}

	return tst
}

func (s *ReplStuffContext) Eos(i int) IEosContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEosContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IEosContext)
}

func (s *ReplStuffContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ReplStuffContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ReplStuffContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.EnterReplStuff(s)
	}
}

func (s *ReplStuffContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.ExitReplStuff(s)
	}
}

func (p *WorksParser) ReplStuff() (localctx IReplStuffContext) {
	localctx = NewReplStuffContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, WorksParserRULE_replStuff)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(37)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<WorksParserT__1)|(1<<WorksParserT__2)|(1<<WorksParserIDENTIFIER)|(1<<WorksParserINT_LIT)|(1<<WorksParserFLOAT_LIT)|(1<<WorksParserIMAGINARY_LIT)|(1<<WorksParserRUNE_LIT)|(1<<WorksParserSTRING_LIT))) != 0) {
		{
			p.SetState(34)
			p.ReplEntry()
		}
		{
			p.SetState(35)
			p.Eos()
		}

		p.SetState(39)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IReplEntryContext is an interface to support dynamic dispatch.
type IReplEntryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsReplEntryContext differentiates from other interfaces.
	IsReplEntryContext()
}

type ReplEntryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyReplEntryContext() *ReplEntryContext {
	var p = new(ReplEntryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = WorksParserRULE_replEntry
	return p
}

func (*ReplEntryContext) IsReplEntryContext() {}

func NewReplEntryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ReplEntryContext {
	var p = new(ReplEntryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = WorksParserRULE_replEntry

	return p
}

func (s *ReplEntryContext) GetParser() antlr.Parser { return s.parser }

func (s *ReplEntryContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ReplEntryContext) EOF() antlr.TerminalNode {
	return s.GetToken(WorksParserEOF, 0)
}

func (s *ReplEntryContext) Assignment() IAssignmentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAssignmentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAssignmentContext)
}

func (s *ReplEntryContext) PackageClause() IPackageClauseContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPackageClauseContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPackageClauseContext)
}

func (s *ReplEntryContext) ImportDecl() IImportDeclContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IImportDeclContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IImportDeclContext)
}

func (s *ReplEntryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ReplEntryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ReplEntryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.EnterReplEntry(s)
	}
}

func (s *ReplEntryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.ExitReplEntry(s)
	}
}

func (p *WorksParser) ReplEntry() (localctx IReplEntryContext) {
	localctx = NewReplEntryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, WorksParserRULE_replEntry)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(53)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case WorksParserINT_LIT, WorksParserFLOAT_LIT, WorksParserIMAGINARY_LIT, WorksParserRUNE_LIT, WorksParserSTRING_LIT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(41)
			p.expression(0)
		}
		{
			p.SetState(42)
			p.Match(WorksParserEOF)
		}

	case WorksParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(44)
			p.Assignment()
		}
		{
			p.SetState(45)
			p.Match(WorksParserEOF)
		}

	case WorksParserT__1:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(47)
			p.PackageClause()
		}
		{
			p.SetState(48)
			p.Match(WorksParserEOF)
		}

	case WorksParserT__2:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(50)
			p.ImportDecl()
		}
		{
			p.SetState(51)
			p.Match(WorksParserEOF)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IPackageClauseContext is an interface to support dynamic dispatch.
type IPackageClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPackageClauseContext differentiates from other interfaces.
	IsPackageClauseContext()
}

type PackageClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPackageClauseContext() *PackageClauseContext {
	var p = new(PackageClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = WorksParserRULE_packageClause
	return p
}

func (*PackageClauseContext) IsPackageClauseContext() {}

func NewPackageClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PackageClauseContext {
	var p = new(PackageClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = WorksParserRULE_packageClause

	return p
}

func (s *PackageClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *PackageClauseContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(WorksParserIDENTIFIER, 0)
}

func (s *PackageClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PackageClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PackageClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.EnterPackageClause(s)
	}
}

func (s *PackageClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.ExitPackageClause(s)
	}
}

func (p *WorksParser) PackageClause() (localctx IPackageClauseContext) {
	localctx = NewPackageClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, WorksParserRULE_packageClause)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(55)
		p.Match(WorksParserT__1)
	}
	{
		p.SetState(56)
		p.Match(WorksParserIDENTIFIER)
	}

	return localctx
}

// IImportDeclContext is an interface to support dynamic dispatch.
type IImportDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsImportDeclContext differentiates from other interfaces.
	IsImportDeclContext()
}

type ImportDeclContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyImportDeclContext() *ImportDeclContext {
	var p = new(ImportDeclContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = WorksParserRULE_importDecl
	return p
}

func (*ImportDeclContext) IsImportDeclContext() {}

func NewImportDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImportDeclContext {
	var p = new(ImportDeclContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = WorksParserRULE_importDecl

	return p
}

func (s *ImportDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *ImportDeclContext) AllImportSpec() []IImportSpecContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IImportSpecContext)(nil)).Elem())
	var tst = make([]IImportSpecContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IImportSpecContext)
		}
	}

	return tst
}

func (s *ImportDeclContext) ImportSpec(i int) IImportSpecContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IImportSpecContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IImportSpecContext)
}

func (s *ImportDeclContext) AllEos() []IEosContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IEosContext)(nil)).Elem())
	var tst = make([]IEosContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IEosContext)
		}
	}

	return tst
}

func (s *ImportDeclContext) Eos(i int) IEosContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEosContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IEosContext)
}

func (s *ImportDeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ImportDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ImportDeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.EnterImportDecl(s)
	}
}

func (s *ImportDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.ExitImportDecl(s)
	}
}

func (p *WorksParser) ImportDecl() (localctx IImportDeclContext) {
	localctx = NewImportDeclContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, WorksParserRULE_importDecl)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(58)
		p.Match(WorksParserT__2)
	}
	p.SetState(70)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case WorksParserT__5, WorksParserIDENTIFIER, WorksParserSTRING_LIT:
		{
			p.SetState(59)
			p.ImportSpec()
		}

	case WorksParserT__3:
		{
			p.SetState(60)
			p.Match(WorksParserT__3)
		}
		p.SetState(66)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<WorksParserT__5)|(1<<WorksParserIDENTIFIER)|(1<<WorksParserSTRING_LIT))) != 0 {
			{
				p.SetState(61)
				p.ImportSpec()
			}
			{
				p.SetState(62)
				p.Eos()
			}

			p.SetState(68)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(69)
			p.Match(WorksParserT__4)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IImportSpecContext is an interface to support dynamic dispatch.
type IImportSpecContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsImportSpecContext differentiates from other interfaces.
	IsImportSpecContext()
}

type ImportSpecContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyImportSpecContext() *ImportSpecContext {
	var p = new(ImportSpecContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = WorksParserRULE_importSpec
	return p
}

func (*ImportSpecContext) IsImportSpecContext() {}

func NewImportSpecContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImportSpecContext {
	var p = new(ImportSpecContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = WorksParserRULE_importSpec

	return p
}

func (s *ImportSpecContext) GetParser() antlr.Parser { return s.parser }

func (s *ImportSpecContext) ImportPath() IImportPathContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IImportPathContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IImportPathContext)
}

func (s *ImportSpecContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(WorksParserIDENTIFIER, 0)
}

func (s *ImportSpecContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ImportSpecContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ImportSpecContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.EnterImportSpec(s)
	}
}

func (s *ImportSpecContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.ExitImportSpec(s)
	}
}

func (p *WorksParser) ImportSpec() (localctx IImportSpecContext) {
	localctx = NewImportSpecContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, WorksParserRULE_importSpec)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(73)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == WorksParserT__5 || _la == WorksParserIDENTIFIER {
		{
			p.SetState(72)
			_la = p.GetTokenStream().LA(1)

			if !(_la == WorksParserT__5 || _la == WorksParserIDENTIFIER) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	}
	{
		p.SetState(75)
		p.ImportPath()
	}

	return localctx
}

// IImportPathContext is an interface to support dynamic dispatch.
type IImportPathContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsImportPathContext differentiates from other interfaces.
	IsImportPathContext()
}

type ImportPathContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyImportPathContext() *ImportPathContext {
	var p = new(ImportPathContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = WorksParserRULE_importPath
	return p
}

func (*ImportPathContext) IsImportPathContext() {}

func NewImportPathContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImportPathContext {
	var p = new(ImportPathContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = WorksParserRULE_importPath

	return p
}

func (s *ImportPathContext) GetParser() antlr.Parser { return s.parser }

func (s *ImportPathContext) STRING_LIT() antlr.TerminalNode {
	return s.GetToken(WorksParserSTRING_LIT, 0)
}

func (s *ImportPathContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ImportPathContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ImportPathContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.EnterImportPath(s)
	}
}

func (s *ImportPathContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.ExitImportPath(s)
	}
}

func (p *WorksParser) ImportPath() (localctx IImportPathContext) {
	localctx = NewImportPathContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, WorksParserRULE_importPath)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(77)
		p.Match(WorksParserSTRING_LIT)
	}

	return localctx
}

// IAssignmentContext is an interface to support dynamic dispatch.
type IAssignmentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAssignmentContext differentiates from other interfaces.
	IsAssignmentContext()
}

type AssignmentContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAssignmentContext() *AssignmentContext {
	var p = new(AssignmentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = WorksParserRULE_assignment
	return p
}

func (*AssignmentContext) IsAssignmentContext() {}

func NewAssignmentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AssignmentContext {
	var p = new(AssignmentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = WorksParserRULE_assignment

	return p
}

func (s *AssignmentContext) GetParser() antlr.Parser { return s.parser }

func (s *AssignmentContext) IdentifierList() IIdentifierListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIdentifierListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIdentifierListContext)
}

func (s *AssignmentContext) ExpressionList() IExpressionListContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionListContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionListContext)
}

func (s *AssignmentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AssignmentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AssignmentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.EnterAssignment(s)
	}
}

func (s *AssignmentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.ExitAssignment(s)
	}
}

func (p *WorksParser) Assignment() (localctx IAssignmentContext) {
	localctx = NewAssignmentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, WorksParserRULE_assignment)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(79)
		p.IdentifierList()
	}
	{
		p.SetState(80)
		p.Match(WorksParserT__6)
	}
	{
		p.SetState(81)
		p.ExpressionList()
	}

	return localctx
}

// IIdentifierListContext is an interface to support dynamic dispatch.
type IIdentifierListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIdentifierListContext differentiates from other interfaces.
	IsIdentifierListContext()
}

type IdentifierListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIdentifierListContext() *IdentifierListContext {
	var p = new(IdentifierListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = WorksParserRULE_identifierList
	return p
}

func (*IdentifierListContext) IsIdentifierListContext() {}

func NewIdentifierListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IdentifierListContext {
	var p = new(IdentifierListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = WorksParserRULE_identifierList

	return p
}

func (s *IdentifierListContext) GetParser() antlr.Parser { return s.parser }

func (s *IdentifierListContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(WorksParserIDENTIFIER)
}

func (s *IdentifierListContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(WorksParserIDENTIFIER, i)
}

func (s *IdentifierListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentifierListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IdentifierListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.EnterIdentifierList(s)
	}
}

func (s *IdentifierListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.ExitIdentifierList(s)
	}
}

func (p *WorksParser) IdentifierList() (localctx IIdentifierListContext) {
	localctx = NewIdentifierListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, WorksParserRULE_identifierList)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(83)
		p.Match(WorksParserIDENTIFIER)
	}
	p.SetState(88)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == WorksParserT__7 {
		{
			p.SetState(84)
			p.Match(WorksParserT__7)
		}
		{
			p.SetState(85)
			p.Match(WorksParserIDENTIFIER)
		}

		p.SetState(90)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IExpressionListContext is an interface to support dynamic dispatch.
type IExpressionListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExpressionListContext differentiates from other interfaces.
	IsExpressionListContext()
}

type ExpressionListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionListContext() *ExpressionListContext {
	var p = new(ExpressionListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = WorksParserRULE_expressionList
	return p
}

func (*ExpressionListContext) IsExpressionListContext() {}

func NewExpressionListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionListContext {
	var p = new(ExpressionListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = WorksParserRULE_expressionList

	return p
}

func (s *ExpressionListContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionListContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *ExpressionListContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ExpressionListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.EnterExpressionList(s)
	}
}

func (s *ExpressionListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.ExitExpressionList(s)
	}
}

func (p *WorksParser) ExpressionList() (localctx IExpressionListContext) {
	localctx = NewExpressionListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, WorksParserRULE_expressionList)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(91)
		p.expression(0)
	}
	p.SetState(96)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == WorksParserT__7 {
		{
			p.SetState(92)
			p.Match(WorksParserT__7)
		}
		{
			p.SetState(93)
			p.expression(0)
		}

		p.SetState(98)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}

type ExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = WorksParserRULE_expression
	return p
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = WorksParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) CopyFrom(ctx *ExpressionContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type BasicLiteralContext struct {
	*ExpressionContext
}

func NewBasicLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BasicLiteralContext {
	var p = new(BasicLiteralContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *BasicLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BasicLiteralContext) BasicLit() IBasicLitContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBasicLitContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBasicLitContext)
}

func (s *BasicLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.EnterBasicLiteral(s)
	}
}

func (s *BasicLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.ExitBasicLiteral(s)
	}
}

type MulDivContext struct {
	*ExpressionContext
	op antlr.Token
}

func NewMulDivContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MulDivContext {
	var p = new(MulDivContext)

	p.ExpressionContext = NewEmptyExpressionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExpressionContext))

	return p
}

func (s *MulDivContext) GetOp() antlr.Token { return s.op }

func (s *MulDivContext) SetOp(v antlr.Token) { s.op = v }

func (s *MulDivContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MulDivContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *MulDivContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *MulDivContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.EnterMulDiv(s)
	}
}

func (s *MulDivContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.ExitMulDiv(s)
	}
}

func (p *WorksParser) Expression() (localctx IExpressionContext) {
	return p.expression(0)
}

func (p *WorksParser) expression(_p int) (localctx IExpressionContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExpressionContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 20
	p.EnterRecursionRule(localctx, 20, WorksParserRULE_expression, _p)
	var _la int

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	localctx = NewBasicLiteralContext(p, localctx)
	p.SetParserRuleContext(localctx)
	_prevctx = localctx

	{
		p.SetState(100)
		p.BasicLit()
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(107)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewMulDivContext(p, NewExpressionContext(p, _parentctx, _parentState))
			p.PushNewRecursionContext(localctx, _startState, WorksParserRULE_expression)
			p.SetState(102)

			if !(p.Precpred(p.GetParserRuleContext(), 2)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
			}
			{
				p.SetState(103)

				var _lt = p.GetTokenStream().LT(1)

				localctx.(*MulDivContext).op = _lt

				_la = p.GetTokenStream().LA(1)

				if !(_la == WorksParserT__8 || _la == WorksParserT__9) {
					var _ri = p.GetErrorHandler().RecoverInline(p)

					localctx.(*MulDivContext).op = _ri
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(104)
				p.expression(3)
			}

		}
		p.SetState(109)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext())
	}

	return localctx
}

// IBasicLitContext is an interface to support dynamic dispatch.
type IBasicLitContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBasicLitContext differentiates from other interfaces.
	IsBasicLitContext()
}

type BasicLitContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBasicLitContext() *BasicLitContext {
	var p = new(BasicLitContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = WorksParserRULE_basicLit
	return p
}

func (*BasicLitContext) IsBasicLitContext() {}

func NewBasicLitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BasicLitContext {
	var p = new(BasicLitContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = WorksParserRULE_basicLit

	return p
}

func (s *BasicLitContext) GetParser() antlr.Parser { return s.parser }

func (s *BasicLitContext) INT_LIT() antlr.TerminalNode {
	return s.GetToken(WorksParserINT_LIT, 0)
}

func (s *BasicLitContext) FLOAT_LIT() antlr.TerminalNode {
	return s.GetToken(WorksParserFLOAT_LIT, 0)
}

func (s *BasicLitContext) IMAGINARY_LIT() antlr.TerminalNode {
	return s.GetToken(WorksParserIMAGINARY_LIT, 0)
}

func (s *BasicLitContext) RUNE_LIT() antlr.TerminalNode {
	return s.GetToken(WorksParserRUNE_LIT, 0)
}

func (s *BasicLitContext) STRING_LIT() antlr.TerminalNode {
	return s.GetToken(WorksParserSTRING_LIT, 0)
}

func (s *BasicLitContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BasicLitContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BasicLitContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.EnterBasicLit(s)
	}
}

func (s *BasicLitContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.ExitBasicLit(s)
	}
}

func (p *WorksParser) BasicLit() (localctx IBasicLitContext) {
	localctx = NewBasicLitContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, WorksParserRULE_basicLit)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(110)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<WorksParserINT_LIT)|(1<<WorksParserFLOAT_LIT)|(1<<WorksParserIMAGINARY_LIT)|(1<<WorksParserRUNE_LIT)|(1<<WorksParserSTRING_LIT))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IOperandNameContext is an interface to support dynamic dispatch.
type IOperandNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsOperandNameContext differentiates from other interfaces.
	IsOperandNameContext()
}

type OperandNameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOperandNameContext() *OperandNameContext {
	var p = new(OperandNameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = WorksParserRULE_operandName
	return p
}

func (*OperandNameContext) IsOperandNameContext() {}

func NewOperandNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OperandNameContext {
	var p = new(OperandNameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = WorksParserRULE_operandName

	return p
}

func (s *OperandNameContext) GetParser() antlr.Parser { return s.parser }

func (s *OperandNameContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(WorksParserIDENTIFIER, 0)
}

func (s *OperandNameContext) QualifiedIdent() IQualifiedIdentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IQualifiedIdentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IQualifiedIdentContext)
}

func (s *OperandNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OperandNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OperandNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.EnterOperandName(s)
	}
}

func (s *OperandNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.ExitOperandName(s)
	}
}

func (p *WorksParser) OperandName() (localctx IOperandNameContext) {
	localctx = NewOperandNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, WorksParserRULE_operandName)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(114)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(112)
			p.Match(WorksParserIDENTIFIER)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(113)
			p.QualifiedIdent()
		}

	}

	return localctx
}

// IQualifiedIdentContext is an interface to support dynamic dispatch.
type IQualifiedIdentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsQualifiedIdentContext differentiates from other interfaces.
	IsQualifiedIdentContext()
}

type QualifiedIdentContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQualifiedIdentContext() *QualifiedIdentContext {
	var p = new(QualifiedIdentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = WorksParserRULE_qualifiedIdent
	return p
}

func (*QualifiedIdentContext) IsQualifiedIdentContext() {}

func NewQualifiedIdentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QualifiedIdentContext {
	var p = new(QualifiedIdentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = WorksParserRULE_qualifiedIdent

	return p
}

func (s *QualifiedIdentContext) GetParser() antlr.Parser { return s.parser }

func (s *QualifiedIdentContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(WorksParserIDENTIFIER)
}

func (s *QualifiedIdentContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(WorksParserIDENTIFIER, i)
}

func (s *QualifiedIdentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QualifiedIdentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QualifiedIdentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.EnterQualifiedIdent(s)
	}
}

func (s *QualifiedIdentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WorksListener); ok {
		listenerT.ExitQualifiedIdent(s)
	}
}

func (p *WorksParser) QualifiedIdent() (localctx IQualifiedIdentContext) {
	localctx = NewQualifiedIdentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, WorksParserRULE_qualifiedIdent)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(116)
		p.Match(WorksParserIDENTIFIER)
	}
	{
		p.SetState(117)
		p.Match(WorksParserT__5)
	}
	{
		p.SetState(118)
		p.Match(WorksParserIDENTIFIER)
	}

	return localctx
}

func (p *WorksParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 0:
		var t *EosContext = nil
		if localctx != nil {
			t = localctx.(*EosContext)
		}
		return p.Eos_Sempred(t, predIndex)

	case 10:
		var t *ExpressionContext = nil
		if localctx != nil {
			t = localctx.(*ExpressionContext)
		}
		return p.Expression_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *WorksParser) Eos_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.lineTerminatorAhead()

	case 1:
		return p.GetTokenStream().LT(1).GetText() == "}"

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *WorksParser) Expression_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 2:
		return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

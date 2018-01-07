// Code generated from Works.g4 by ANTLR 4.7.1. DO NOT EDIT.

package parser // Works

import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseWorksListener is a complete listener for a parse tree produced by WorksParser.
type BaseWorksListener struct{}

var _ WorksListener = &BaseWorksListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseWorksListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseWorksListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseWorksListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseWorksListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterEos is called when production eos is entered.
func (s *BaseWorksListener) EnterEos(ctx *EosContext) {}

// ExitEos is called when production eos is exited.
func (s *BaseWorksListener) ExitEos(ctx *EosContext) {}

// EnterReplStuff is called when production replStuff is entered.
func (s *BaseWorksListener) EnterReplStuff(ctx *ReplStuffContext) {}

// ExitReplStuff is called when production replStuff is exited.
func (s *BaseWorksListener) ExitReplStuff(ctx *ReplStuffContext) {}

// EnterReplEntry is called when production replEntry is entered.
func (s *BaseWorksListener) EnterReplEntry(ctx *ReplEntryContext) {}

// ExitReplEntry is called when production replEntry is exited.
func (s *BaseWorksListener) ExitReplEntry(ctx *ReplEntryContext) {}

// EnterPackageClause is called when production packageClause is entered.
func (s *BaseWorksListener) EnterPackageClause(ctx *PackageClauseContext) {}

// ExitPackageClause is called when production packageClause is exited.
func (s *BaseWorksListener) ExitPackageClause(ctx *PackageClauseContext) {}

// EnterImportDecl is called when production importDecl is entered.
func (s *BaseWorksListener) EnterImportDecl(ctx *ImportDeclContext) {}

// ExitImportDecl is called when production importDecl is exited.
func (s *BaseWorksListener) ExitImportDecl(ctx *ImportDeclContext) {}

// EnterImportSpec is called when production importSpec is entered.
func (s *BaseWorksListener) EnterImportSpec(ctx *ImportSpecContext) {}

// ExitImportSpec is called when production importSpec is exited.
func (s *BaseWorksListener) ExitImportSpec(ctx *ImportSpecContext) {}

// EnterImportPath is called when production importPath is entered.
func (s *BaseWorksListener) EnterImportPath(ctx *ImportPathContext) {}

// ExitImportPath is called when production importPath is exited.
func (s *BaseWorksListener) ExitImportPath(ctx *ImportPathContext) {}

// EnterAssignment is called when production assignment is entered.
func (s *BaseWorksListener) EnterAssignment(ctx *AssignmentContext) {}

// ExitAssignment is called when production assignment is exited.
func (s *BaseWorksListener) ExitAssignment(ctx *AssignmentContext) {}

// EnterIdentifierList is called when production identifierList is entered.
func (s *BaseWorksListener) EnterIdentifierList(ctx *IdentifierListContext) {}

// ExitIdentifierList is called when production identifierList is exited.
func (s *BaseWorksListener) ExitIdentifierList(ctx *IdentifierListContext) {}

// EnterExpressionList is called when production expressionList is entered.
func (s *BaseWorksListener) EnterExpressionList(ctx *ExpressionListContext) {}

// ExitExpressionList is called when production expressionList is exited.
func (s *BaseWorksListener) ExitExpressionList(ctx *ExpressionListContext) {}

// EnterBasicLiteral is called when production BasicLiteral is entered.
func (s *BaseWorksListener) EnterBasicLiteral(ctx *BasicLiteralContext) {}

// ExitBasicLiteral is called when production BasicLiteral is exited.
func (s *BaseWorksListener) ExitBasicLiteral(ctx *BasicLiteralContext) {}

// EnterMulDiv is called when production MulDiv is entered.
func (s *BaseWorksListener) EnterMulDiv(ctx *MulDivContext) {}

// ExitMulDiv is called when production MulDiv is exited.
func (s *BaseWorksListener) ExitMulDiv(ctx *MulDivContext) {}

// EnterBasicLit is called when production basicLit is entered.
func (s *BaseWorksListener) EnterBasicLit(ctx *BasicLitContext) {}

// ExitBasicLit is called when production basicLit is exited.
func (s *BaseWorksListener) ExitBasicLit(ctx *BasicLitContext) {}

// EnterOperandName is called when production operandName is entered.
func (s *BaseWorksListener) EnterOperandName(ctx *OperandNameContext) {}

// ExitOperandName is called when production operandName is exited.
func (s *BaseWorksListener) ExitOperandName(ctx *OperandNameContext) {}

// EnterQualifiedIdent is called when production qualifiedIdent is entered.
func (s *BaseWorksListener) EnterQualifiedIdent(ctx *QualifiedIdentContext) {}

// ExitQualifiedIdent is called when production qualifiedIdent is exited.
func (s *BaseWorksListener) ExitQualifiedIdent(ctx *QualifiedIdentContext) {}

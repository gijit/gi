// Code generated from Works.g4 by ANTLR 4.7.1. DO NOT EDIT.

package parser // Works

import "github.com/antlr/antlr4/runtime/Go/antlr"

// WorksListener is a complete listener for a parse tree produced by WorksParser.
type WorksListener interface {
	antlr.ParseTreeListener

	// EnterEos is called when entering the eos production.
	EnterEos(c *EosContext)

	// EnterReplStuff is called when entering the replStuff production.
	EnterReplStuff(c *ReplStuffContext)

	// EnterReplEntry is called when entering the replEntry production.
	EnterReplEntry(c *ReplEntryContext)

	// EnterPackageClause is called when entering the packageClause production.
	EnterPackageClause(c *PackageClauseContext)

	// EnterImportDecl is called when entering the importDecl production.
	EnterImportDecl(c *ImportDeclContext)

	// EnterImportSpec is called when entering the importSpec production.
	EnterImportSpec(c *ImportSpecContext)

	// EnterImportPath is called when entering the importPath production.
	EnterImportPath(c *ImportPathContext)

	// EnterAssignment is called when entering the assignment production.
	EnterAssignment(c *AssignmentContext)

	// EnterIdentifierList is called when entering the identifierList production.
	EnterIdentifierList(c *IdentifierListContext)

	// EnterExpressionList is called when entering the expressionList production.
	EnterExpressionList(c *ExpressionListContext)

	// EnterBasicLiteral is called when entering the BasicLiteral production.
	EnterBasicLiteral(c *BasicLiteralContext)

	// EnterMulDiv is called when entering the MulDiv production.
	EnterMulDiv(c *MulDivContext)

	// EnterBasicLit is called when entering the basicLit production.
	EnterBasicLit(c *BasicLitContext)

	// EnterOperandName is called when entering the operandName production.
	EnterOperandName(c *OperandNameContext)

	// EnterQualifiedIdent is called when entering the qualifiedIdent production.
	EnterQualifiedIdent(c *QualifiedIdentContext)

	// ExitEos is called when exiting the eos production.
	ExitEos(c *EosContext)

	// ExitReplStuff is called when exiting the replStuff production.
	ExitReplStuff(c *ReplStuffContext)

	// ExitReplEntry is called when exiting the replEntry production.
	ExitReplEntry(c *ReplEntryContext)

	// ExitPackageClause is called when exiting the packageClause production.
	ExitPackageClause(c *PackageClauseContext)

	// ExitImportDecl is called when exiting the importDecl production.
	ExitImportDecl(c *ImportDeclContext)

	// ExitImportSpec is called when exiting the importSpec production.
	ExitImportSpec(c *ImportSpecContext)

	// ExitImportPath is called when exiting the importPath production.
	ExitImportPath(c *ImportPathContext)

	// ExitAssignment is called when exiting the assignment production.
	ExitAssignment(c *AssignmentContext)

	// ExitIdentifierList is called when exiting the identifierList production.
	ExitIdentifierList(c *IdentifierListContext)

	// ExitExpressionList is called when exiting the expressionList production.
	ExitExpressionList(c *ExpressionListContext)

	// ExitBasicLiteral is called when exiting the BasicLiteral production.
	ExitBasicLiteral(c *BasicLiteralContext)

	// ExitMulDiv is called when exiting the MulDiv production.
	ExitMulDiv(c *MulDivContext)

	// ExitBasicLit is called when exiting the basicLit production.
	ExitBasicLit(c *BasicLitContext)

	// ExitOperandName is called when exiting the operandName production.
	ExitOperandName(c *OperandNameContext)

	// ExitQualifiedIdent is called when exiting the qualifiedIdent production.
	ExitQualifiedIdent(c *QualifiedIdentContext)
}

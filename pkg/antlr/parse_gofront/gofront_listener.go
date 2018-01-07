// Code generated from Gofront.g4 by ANTLR 4.7.1. DO NOT EDIT.

package parser // Gofront

import "github.com/antlr/antlr4/runtime/Go/antlr"

// GofrontListener is a complete listener for a parse tree produced by GofrontParser.
type GofrontListener interface {
	antlr.ParseTreeListener

	// EnterReplStuff is called when entering the replStuff production.
	EnterReplStuff(c *ReplStuffContext)

	// EnterReplEntry is called when entering the replEntry production.
	EnterReplEntry(c *ReplEntryContext)

	// EnterSourceFile is called when entering the sourceFile production.
	EnterSourceFile(c *SourceFileContext)

	// EnterPackageClause is called when entering the packageClause production.
	EnterPackageClause(c *PackageClauseContext)

	// EnterImportDecl is called when entering the importDecl production.
	EnterImportDecl(c *ImportDeclContext)

	// EnterImportSpec is called when entering the importSpec production.
	EnterImportSpec(c *ImportSpecContext)

	// EnterImportPath is called when entering the importPath production.
	EnterImportPath(c *ImportPathContext)

	// EnterTopLevelDecl is called when entering the topLevelDecl production.
	EnterTopLevelDecl(c *TopLevelDeclContext)

	// EnterDeclaration is called when entering the declaration production.
	EnterDeclaration(c *DeclarationContext)

	// EnterConstDecl is called when entering the constDecl production.
	EnterConstDecl(c *ConstDeclContext)

	// EnterConstSpec is called when entering the constSpec production.
	EnterConstSpec(c *ConstSpecContext)

	// EnterIdentifierList is called when entering the identifierList production.
	EnterIdentifierList(c *IdentifierListContext)

	// EnterExpressionList is called when entering the expressionList production.
	EnterExpressionList(c *ExpressionListContext)

	// EnterTypeDecl is called when entering the typeDecl production.
	EnterTypeDecl(c *TypeDeclContext)

	// EnterTypeSpec is called when entering the typeSpec production.
	EnterTypeSpec(c *TypeSpecContext)

	// EnterFunctionDecl is called when entering the functionDecl production.
	EnterFunctionDecl(c *FunctionDeclContext)

	// EnterFunction is called when entering the function production.
	EnterFunction(c *FunctionContext)

	// EnterMethodDecl is called when entering the methodDecl production.
	EnterMethodDecl(c *MethodDeclContext)

	// EnterReceiver is called when entering the receiver production.
	EnterReceiver(c *ReceiverContext)

	// EnterVarDecl is called when entering the varDecl production.
	EnterVarDecl(c *VarDeclContext)

	// EnterVarSpec is called when entering the varSpec production.
	EnterVarSpec(c *VarSpecContext)

	// EnterBlock is called when entering the block production.
	EnterBlock(c *BlockContext)

	// EnterStatementList is called when entering the statementList production.
	EnterStatementList(c *StatementListContext)

	// EnterStatement is called when entering the statement production.
	EnterStatement(c *StatementContext)

	// EnterSimpleStmt is called when entering the simpleStmt production.
	EnterSimpleStmt(c *SimpleStmtContext)

	// EnterExpressionStmt is called when entering the expressionStmt production.
	EnterExpressionStmt(c *ExpressionStmtContext)

	// EnterSendStmt is called when entering the sendStmt production.
	EnterSendStmt(c *SendStmtContext)

	// EnterIncDecStmt is called when entering the incDecStmt production.
	EnterIncDecStmt(c *IncDecStmtContext)

	// EnterAssignment is called when entering the assignment production.
	EnterAssignment(c *AssignmentContext)

	// EnterAssign_op is called when entering the assign_op production.
	EnterAssign_op(c *Assign_opContext)

	// EnterShortVarDecl is called when entering the shortVarDecl production.
	EnterShortVarDecl(c *ShortVarDeclContext)

	// EnterEmptyStmt is called when entering the emptyStmt production.
	EnterEmptyStmt(c *EmptyStmtContext)

	// EnterLabeledStmt is called when entering the labeledStmt production.
	EnterLabeledStmt(c *LabeledStmtContext)

	// EnterReturnStmt is called when entering the returnStmt production.
	EnterReturnStmt(c *ReturnStmtContext)

	// EnterBreakStmt is called when entering the breakStmt production.
	EnterBreakStmt(c *BreakStmtContext)

	// EnterContinueStmt is called when entering the continueStmt production.
	EnterContinueStmt(c *ContinueStmtContext)

	// EnterGotoStmt is called when entering the gotoStmt production.
	EnterGotoStmt(c *GotoStmtContext)

	// EnterFallthroughStmt is called when entering the fallthroughStmt production.
	EnterFallthroughStmt(c *FallthroughStmtContext)

	// EnterDeferStmt is called when entering the deferStmt production.
	EnterDeferStmt(c *DeferStmtContext)

	// EnterIfStmt is called when entering the ifStmt production.
	EnterIfStmt(c *IfStmtContext)

	// EnterSwitchStmt is called when entering the switchStmt production.
	EnterSwitchStmt(c *SwitchStmtContext)

	// EnterExprSwitchStmt is called when entering the exprSwitchStmt production.
	EnterExprSwitchStmt(c *ExprSwitchStmtContext)

	// EnterExprCaseClause is called when entering the exprCaseClause production.
	EnterExprCaseClause(c *ExprCaseClauseContext)

	// EnterExprSwitchCase is called when entering the exprSwitchCase production.
	EnterExprSwitchCase(c *ExprSwitchCaseContext)

	// EnterTypeSwitchStmt is called when entering the typeSwitchStmt production.
	EnterTypeSwitchStmt(c *TypeSwitchStmtContext)

	// EnterTypeSwitchGuard is called when entering the typeSwitchGuard production.
	EnterTypeSwitchGuard(c *TypeSwitchGuardContext)

	// EnterTypeCaseClause is called when entering the typeCaseClause production.
	EnterTypeCaseClause(c *TypeCaseClauseContext)

	// EnterTypeSwitchCase is called when entering the typeSwitchCase production.
	EnterTypeSwitchCase(c *TypeSwitchCaseContext)

	// EnterTypeList is called when entering the typeList production.
	EnterTypeList(c *TypeListContext)

	// EnterSelectStmt is called when entering the selectStmt production.
	EnterSelectStmt(c *SelectStmtContext)

	// EnterCommClause is called when entering the commClause production.
	EnterCommClause(c *CommClauseContext)

	// EnterCommCase is called when entering the commCase production.
	EnterCommCase(c *CommCaseContext)

	// EnterRecvStmt is called when entering the recvStmt production.
	EnterRecvStmt(c *RecvStmtContext)

	// EnterForStmt is called when entering the forStmt production.
	EnterForStmt(c *ForStmtContext)

	// EnterForClause is called when entering the forClause production.
	EnterForClause(c *ForClauseContext)

	// EnterRangeClause is called when entering the rangeClause production.
	EnterRangeClause(c *RangeClauseContext)

	// EnterGoStmt is called when entering the goStmt production.
	EnterGoStmt(c *GoStmtContext)

	// EnterLtype is called when entering the ltype production.
	EnterLtype(c *LtypeContext)

	// EnterTypeName is called when entering the typeName production.
	EnterTypeName(c *TypeNameContext)

	// EnterTypeLit is called when entering the typeLit production.
	EnterTypeLit(c *TypeLitContext)

	// EnterArrayType is called when entering the arrayType production.
	EnterArrayType(c *ArrayTypeContext)

	// EnterArrayLength is called when entering the arrayLength production.
	EnterArrayLength(c *ArrayLengthContext)

	// EnterElementType is called when entering the elementType production.
	EnterElementType(c *ElementTypeContext)

	// EnterPointerType is called when entering the pointerType production.
	EnterPointerType(c *PointerTypeContext)

	// EnterInterfaceType is called when entering the interfaceType production.
	EnterInterfaceType(c *InterfaceTypeContext)

	// EnterSliceType is called when entering the sliceType production.
	EnterSliceType(c *SliceTypeContext)

	// EnterMapType is called when entering the mapType production.
	EnterMapType(c *MapTypeContext)

	// EnterChannelType is called when entering the channelType production.
	EnterChannelType(c *ChannelTypeContext)

	// EnterMethodSpec is called when entering the methodSpec production.
	EnterMethodSpec(c *MethodSpecContext)

	// EnterFunctionType is called when entering the functionType production.
	EnterFunctionType(c *FunctionTypeContext)

	// EnterSignature is called when entering the signature production.
	EnterSignature(c *SignatureContext)

	// EnterResult is called when entering the result production.
	EnterResult(c *ResultContext)

	// EnterParameters is called when entering the parameters production.
	EnterParameters(c *ParametersContext)

	// EnterParameterList is called when entering the parameterList production.
	EnterParameterList(c *ParameterListContext)

	// EnterParameterDecl is called when entering the parameterDecl production.
	EnterParameterDecl(c *ParameterDeclContext)

	// EnterOperand is called when entering the operand production.
	EnterOperand(c *OperandContext)

	// EnterLiteral is called when entering the literal production.
	EnterLiteral(c *LiteralContext)

	// EnterBasicLit is called when entering the basicLit production.
	EnterBasicLit(c *BasicLitContext)

	// EnterOperandName is called when entering the operandName production.
	EnterOperandName(c *OperandNameContext)

	// EnterQualifiedIdent is called when entering the qualifiedIdent production.
	EnterQualifiedIdent(c *QualifiedIdentContext)

	// EnterCompositeLit is called when entering the compositeLit production.
	EnterCompositeLit(c *CompositeLitContext)

	// EnterLiteralType is called when entering the literalType production.
	EnterLiteralType(c *LiteralTypeContext)

	// EnterLiteralValue is called when entering the literalValue production.
	EnterLiteralValue(c *LiteralValueContext)

	// EnterElementList is called when entering the elementList production.
	EnterElementList(c *ElementListContext)

	// EnterKeyedElement is called when entering the keyedElement production.
	EnterKeyedElement(c *KeyedElementContext)

	// EnterKey is called when entering the key production.
	EnterKey(c *KeyContext)

	// EnterElement is called when entering the element production.
	EnterElement(c *ElementContext)

	// EnterStructType is called when entering the structType production.
	EnterStructType(c *StructTypeContext)

	// EnterFieldDecl is called when entering the fieldDecl production.
	EnterFieldDecl(c *FieldDeclContext)

	// EnterAnonymousField is called when entering the anonymousField production.
	EnterAnonymousField(c *AnonymousFieldContext)

	// EnterFunctionLit is called when entering the functionLit production.
	EnterFunctionLit(c *FunctionLitContext)

	// EnterPrimaryExpr is called when entering the primaryExpr production.
	EnterPrimaryExpr(c *PrimaryExprContext)

	// EnterSelector is called when entering the selector production.
	EnterSelector(c *SelectorContext)

	// EnterIndex is called when entering the index production.
	EnterIndex(c *IndexContext)

	// EnterSlice is called when entering the slice production.
	EnterSlice(c *SliceContext)

	// EnterTypeAssertion is called when entering the typeAssertion production.
	EnterTypeAssertion(c *TypeAssertionContext)

	// EnterArguments is called when entering the arguments production.
	EnterArguments(c *ArgumentsContext)

	// EnterMethodExpr is called when entering the methodExpr production.
	EnterMethodExpr(c *MethodExprContext)

	// EnterReceiverType is called when entering the receiverType production.
	EnterReceiverType(c *ReceiverTypeContext)

	// EnterExpression is called when entering the expression production.
	EnterExpression(c *ExpressionContext)

	// EnterUnaryExpr is called when entering the unaryExpr production.
	EnterUnaryExpr(c *UnaryExprContext)

	// EnterConversion is called when entering the conversion production.
	EnterConversion(c *ConversionContext)

	// EnterEos is called when entering the eos production.
	EnterEos(c *EosContext)

	// ExitReplStuff is called when exiting the replStuff production.
	ExitReplStuff(c *ReplStuffContext)

	// ExitReplEntry is called when exiting the replEntry production.
	ExitReplEntry(c *ReplEntryContext)

	// ExitSourceFile is called when exiting the sourceFile production.
	ExitSourceFile(c *SourceFileContext)

	// ExitPackageClause is called when exiting the packageClause production.
	ExitPackageClause(c *PackageClauseContext)

	// ExitImportDecl is called when exiting the importDecl production.
	ExitImportDecl(c *ImportDeclContext)

	// ExitImportSpec is called when exiting the importSpec production.
	ExitImportSpec(c *ImportSpecContext)

	// ExitImportPath is called when exiting the importPath production.
	ExitImportPath(c *ImportPathContext)

	// ExitTopLevelDecl is called when exiting the topLevelDecl production.
	ExitTopLevelDecl(c *TopLevelDeclContext)

	// ExitDeclaration is called when exiting the declaration production.
	ExitDeclaration(c *DeclarationContext)

	// ExitConstDecl is called when exiting the constDecl production.
	ExitConstDecl(c *ConstDeclContext)

	// ExitConstSpec is called when exiting the constSpec production.
	ExitConstSpec(c *ConstSpecContext)

	// ExitIdentifierList is called when exiting the identifierList production.
	ExitIdentifierList(c *IdentifierListContext)

	// ExitExpressionList is called when exiting the expressionList production.
	ExitExpressionList(c *ExpressionListContext)

	// ExitTypeDecl is called when exiting the typeDecl production.
	ExitTypeDecl(c *TypeDeclContext)

	// ExitTypeSpec is called when exiting the typeSpec production.
	ExitTypeSpec(c *TypeSpecContext)

	// ExitFunctionDecl is called when exiting the functionDecl production.
	ExitFunctionDecl(c *FunctionDeclContext)

	// ExitFunction is called when exiting the function production.
	ExitFunction(c *FunctionContext)

	// ExitMethodDecl is called when exiting the methodDecl production.
	ExitMethodDecl(c *MethodDeclContext)

	// ExitReceiver is called when exiting the receiver production.
	ExitReceiver(c *ReceiverContext)

	// ExitVarDecl is called when exiting the varDecl production.
	ExitVarDecl(c *VarDeclContext)

	// ExitVarSpec is called when exiting the varSpec production.
	ExitVarSpec(c *VarSpecContext)

	// ExitBlock is called when exiting the block production.
	ExitBlock(c *BlockContext)

	// ExitStatementList is called when exiting the statementList production.
	ExitStatementList(c *StatementListContext)

	// ExitStatement is called when exiting the statement production.
	ExitStatement(c *StatementContext)

	// ExitSimpleStmt is called when exiting the simpleStmt production.
	ExitSimpleStmt(c *SimpleStmtContext)

	// ExitExpressionStmt is called when exiting the expressionStmt production.
	ExitExpressionStmt(c *ExpressionStmtContext)

	// ExitSendStmt is called when exiting the sendStmt production.
	ExitSendStmt(c *SendStmtContext)

	// ExitIncDecStmt is called when exiting the incDecStmt production.
	ExitIncDecStmt(c *IncDecStmtContext)

	// ExitAssignment is called when exiting the assignment production.
	ExitAssignment(c *AssignmentContext)

	// ExitAssign_op is called when exiting the assign_op production.
	ExitAssign_op(c *Assign_opContext)

	// ExitShortVarDecl is called when exiting the shortVarDecl production.
	ExitShortVarDecl(c *ShortVarDeclContext)

	// ExitEmptyStmt is called when exiting the emptyStmt production.
	ExitEmptyStmt(c *EmptyStmtContext)

	// ExitLabeledStmt is called when exiting the labeledStmt production.
	ExitLabeledStmt(c *LabeledStmtContext)

	// ExitReturnStmt is called when exiting the returnStmt production.
	ExitReturnStmt(c *ReturnStmtContext)

	// ExitBreakStmt is called when exiting the breakStmt production.
	ExitBreakStmt(c *BreakStmtContext)

	// ExitContinueStmt is called when exiting the continueStmt production.
	ExitContinueStmt(c *ContinueStmtContext)

	// ExitGotoStmt is called when exiting the gotoStmt production.
	ExitGotoStmt(c *GotoStmtContext)

	// ExitFallthroughStmt is called when exiting the fallthroughStmt production.
	ExitFallthroughStmt(c *FallthroughStmtContext)

	// ExitDeferStmt is called when exiting the deferStmt production.
	ExitDeferStmt(c *DeferStmtContext)

	// ExitIfStmt is called when exiting the ifStmt production.
	ExitIfStmt(c *IfStmtContext)

	// ExitSwitchStmt is called when exiting the switchStmt production.
	ExitSwitchStmt(c *SwitchStmtContext)

	// ExitExprSwitchStmt is called when exiting the exprSwitchStmt production.
	ExitExprSwitchStmt(c *ExprSwitchStmtContext)

	// ExitExprCaseClause is called when exiting the exprCaseClause production.
	ExitExprCaseClause(c *ExprCaseClauseContext)

	// ExitExprSwitchCase is called when exiting the exprSwitchCase production.
	ExitExprSwitchCase(c *ExprSwitchCaseContext)

	// ExitTypeSwitchStmt is called when exiting the typeSwitchStmt production.
	ExitTypeSwitchStmt(c *TypeSwitchStmtContext)

	// ExitTypeSwitchGuard is called when exiting the typeSwitchGuard production.
	ExitTypeSwitchGuard(c *TypeSwitchGuardContext)

	// ExitTypeCaseClause is called when exiting the typeCaseClause production.
	ExitTypeCaseClause(c *TypeCaseClauseContext)

	// ExitTypeSwitchCase is called when exiting the typeSwitchCase production.
	ExitTypeSwitchCase(c *TypeSwitchCaseContext)

	// ExitTypeList is called when exiting the typeList production.
	ExitTypeList(c *TypeListContext)

	// ExitSelectStmt is called when exiting the selectStmt production.
	ExitSelectStmt(c *SelectStmtContext)

	// ExitCommClause is called when exiting the commClause production.
	ExitCommClause(c *CommClauseContext)

	// ExitCommCase is called when exiting the commCase production.
	ExitCommCase(c *CommCaseContext)

	// ExitRecvStmt is called when exiting the recvStmt production.
	ExitRecvStmt(c *RecvStmtContext)

	// ExitForStmt is called when exiting the forStmt production.
	ExitForStmt(c *ForStmtContext)

	// ExitForClause is called when exiting the forClause production.
	ExitForClause(c *ForClauseContext)

	// ExitRangeClause is called when exiting the rangeClause production.
	ExitRangeClause(c *RangeClauseContext)

	// ExitGoStmt is called when exiting the goStmt production.
	ExitGoStmt(c *GoStmtContext)

	// ExitLtype is called when exiting the ltype production.
	ExitLtype(c *LtypeContext)

	// ExitTypeName is called when exiting the typeName production.
	ExitTypeName(c *TypeNameContext)

	// ExitTypeLit is called when exiting the typeLit production.
	ExitTypeLit(c *TypeLitContext)

	// ExitArrayType is called when exiting the arrayType production.
	ExitArrayType(c *ArrayTypeContext)

	// ExitArrayLength is called when exiting the arrayLength production.
	ExitArrayLength(c *ArrayLengthContext)

	// ExitElementType is called when exiting the elementType production.
	ExitElementType(c *ElementTypeContext)

	// ExitPointerType is called when exiting the pointerType production.
	ExitPointerType(c *PointerTypeContext)

	// ExitInterfaceType is called when exiting the interfaceType production.
	ExitInterfaceType(c *InterfaceTypeContext)

	// ExitSliceType is called when exiting the sliceType production.
	ExitSliceType(c *SliceTypeContext)

	// ExitMapType is called when exiting the mapType production.
	ExitMapType(c *MapTypeContext)

	// ExitChannelType is called when exiting the channelType production.
	ExitChannelType(c *ChannelTypeContext)

	// ExitMethodSpec is called when exiting the methodSpec production.
	ExitMethodSpec(c *MethodSpecContext)

	// ExitFunctionType is called when exiting the functionType production.
	ExitFunctionType(c *FunctionTypeContext)

	// ExitSignature is called when exiting the signature production.
	ExitSignature(c *SignatureContext)

	// ExitResult is called when exiting the result production.
	ExitResult(c *ResultContext)

	// ExitParameters is called when exiting the parameters production.
	ExitParameters(c *ParametersContext)

	// ExitParameterList is called when exiting the parameterList production.
	ExitParameterList(c *ParameterListContext)

	// ExitParameterDecl is called when exiting the parameterDecl production.
	ExitParameterDecl(c *ParameterDeclContext)

	// ExitOperand is called when exiting the operand production.
	ExitOperand(c *OperandContext)

	// ExitLiteral is called when exiting the literal production.
	ExitLiteral(c *LiteralContext)

	// ExitBasicLit is called when exiting the basicLit production.
	ExitBasicLit(c *BasicLitContext)

	// ExitOperandName is called when exiting the operandName production.
	ExitOperandName(c *OperandNameContext)

	// ExitQualifiedIdent is called when exiting the qualifiedIdent production.
	ExitQualifiedIdent(c *QualifiedIdentContext)

	// ExitCompositeLit is called when exiting the compositeLit production.
	ExitCompositeLit(c *CompositeLitContext)

	// ExitLiteralType is called when exiting the literalType production.
	ExitLiteralType(c *LiteralTypeContext)

	// ExitLiteralValue is called when exiting the literalValue production.
	ExitLiteralValue(c *LiteralValueContext)

	// ExitElementList is called when exiting the elementList production.
	ExitElementList(c *ElementListContext)

	// ExitKeyedElement is called when exiting the keyedElement production.
	ExitKeyedElement(c *KeyedElementContext)

	// ExitKey is called when exiting the key production.
	ExitKey(c *KeyContext)

	// ExitElement is called when exiting the element production.
	ExitElement(c *ElementContext)

	// ExitStructType is called when exiting the structType production.
	ExitStructType(c *StructTypeContext)

	// ExitFieldDecl is called when exiting the fieldDecl production.
	ExitFieldDecl(c *FieldDeclContext)

	// ExitAnonymousField is called when exiting the anonymousField production.
	ExitAnonymousField(c *AnonymousFieldContext)

	// ExitFunctionLit is called when exiting the functionLit production.
	ExitFunctionLit(c *FunctionLitContext)

	// ExitPrimaryExpr is called when exiting the primaryExpr production.
	ExitPrimaryExpr(c *PrimaryExprContext)

	// ExitSelector is called when exiting the selector production.
	ExitSelector(c *SelectorContext)

	// ExitIndex is called when exiting the index production.
	ExitIndex(c *IndexContext)

	// ExitSlice is called when exiting the slice production.
	ExitSlice(c *SliceContext)

	// ExitTypeAssertion is called when exiting the typeAssertion production.
	ExitTypeAssertion(c *TypeAssertionContext)

	// ExitArguments is called when exiting the arguments production.
	ExitArguments(c *ArgumentsContext)

	// ExitMethodExpr is called when exiting the methodExpr production.
	ExitMethodExpr(c *MethodExprContext)

	// ExitReceiverType is called when exiting the receiverType production.
	ExitReceiverType(c *ReceiverTypeContext)

	// ExitExpression is called when exiting the expression production.
	ExitExpression(c *ExpressionContext)

	// ExitUnaryExpr is called when exiting the unaryExpr production.
	ExitUnaryExpr(c *UnaryExprContext)

	// ExitConversion is called when exiting the conversion production.
	ExitConversion(c *ConversionContext)

	// ExitEos is called when exiting the eos production.
	ExitEos(c *EosContext)
}

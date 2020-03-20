package ast

import "github.com/dapperlabs/cadence/runtime/common"

type FunctionDeclaration struct {
	Access               Access
	Identifier           Identifier
	ParameterList        *ParameterList
	ReturnTypeAnnotation *TypeAnnotation
	FunctionBlock        *FunctionBlock
	StartPos             Position
}

func (f *FunctionDeclaration) StartPosition() Position {
	return f.StartPos
}

func (f *FunctionDeclaration) EndPosition() Position {
	if f.FunctionBlock != nil {
		return f.FunctionBlock.EndPosition()
	}
	if f.ReturnTypeAnnotation != nil {
		return f.ReturnTypeAnnotation.EndPosition()
	}
	return f.ParameterList.EndPosition()
}

func (f *FunctionDeclaration) Accept(visitor Visitor) Repr {
	return visitor.VisitFunctionDeclaration(f)
}

func (*FunctionDeclaration) isDeclaration() {}
func (*FunctionDeclaration) isStatement()   {}

func (f *FunctionDeclaration) DeclarationIdentifier() *Identifier {
	return &f.Identifier
}

func (f *FunctionDeclaration) DeclarationKind() common.DeclarationKind {
	return common.DeclarationKindFunction
}

func (f *FunctionDeclaration) DeclarationAccess() Access {
	return f.Access
}

func (f *FunctionDeclaration) ToExpression() *FunctionExpression {
	return &FunctionExpression{
		ParameterList:        f.ParameterList,
		ReturnTypeAnnotation: f.ReturnTypeAnnotation,
		FunctionBlock:        f.FunctionBlock,
		StartPos:             f.StartPos,
	}
}

// SpecialFunctionDeclaration

type SpecialFunctionDeclaration struct {
	DeclarationKind common.DeclarationKind
	*FunctionDeclaration
}

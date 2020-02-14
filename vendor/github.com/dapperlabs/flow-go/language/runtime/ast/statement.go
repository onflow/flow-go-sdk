package ast

type Statement interface {
	Element
	isStatement()
}

// ReturnStatement

type ReturnStatement struct {
	Expression Expression
	Range
}

func (*ReturnStatement) isStatement() {}

func (s *ReturnStatement) Accept(visitor Visitor) Repr {
	return visitor.VisitReturnStatement(s)
}

// BreakStatement

type BreakStatement struct {
	Range
}

func (*BreakStatement) isStatement() {}

func (s *BreakStatement) Accept(visitor Visitor) Repr {
	return visitor.VisitBreakStatement(s)
}

// ContinueStatement

type ContinueStatement struct {
	Range
}

func (*ContinueStatement) isStatement() {}

func (s *ContinueStatement) Accept(visitor Visitor) Repr {
	return visitor.VisitContinueStatement(s)
}

// IfStatementTest

type IfStatementTest interface {
	isIfStatementTest()
}

// IfStatement

type IfStatement struct {
	Test     IfStatementTest
	Then     *Block
	Else     *Block
	StartPos Position
}

func (s *IfStatement) StartPosition() Position {
	return s.StartPos
}

func (s *IfStatement) EndPosition() Position {
	if s.Else != nil {
		return s.Else.EndPosition()
	}
	return s.Then.EndPosition()
}

func (*IfStatement) isStatement() {}

func (s *IfStatement) Accept(visitor Visitor) Repr {
	return visitor.VisitIfStatement(s)
}

// WhileStatement

type WhileStatement struct {
	Test  Expression
	Block *Block
	Range
}

func (*WhileStatement) isStatement() {}

func (s *WhileStatement) Accept(visitor Visitor) Repr {
	return visitor.VisitWhileStatement(s)
}

// EmitStatement

type EmitStatement struct {
	InvocationExpression *InvocationExpression
	StartPos             Position
}

func (s *EmitStatement) StartPosition() Position {
	return s.StartPos
}

func (s *EmitStatement) EndPosition() Position {
	return s.InvocationExpression.EndPosition()
}

func (*EmitStatement) isStatement() {}

func (s *EmitStatement) Accept(visitor Visitor) Repr {
	return visitor.VisitEmitStatement(s)
}

// AssignmentStatement

type AssignmentStatement struct {
	Target   Expression
	Transfer *Transfer
	Value    Expression
}

func (s *AssignmentStatement) StartPosition() Position {
	return s.Target.StartPosition()
}

func (s *AssignmentStatement) EndPosition() Position {
	return s.Value.EndPosition()
}

func (*AssignmentStatement) isStatement() {}

func (s *AssignmentStatement) Accept(visitor Visitor) Repr {
	return visitor.VisitAssignmentStatement(s)
}

// SwapStatement

type SwapStatement struct {
	Left  Expression
	Right Expression
}

func (s *SwapStatement) StartPosition() Position {
	return s.Left.StartPosition()
}

func (s *SwapStatement) EndPosition() Position {
	return s.Right.EndPosition()
}

func (*SwapStatement) isStatement() {}

func (s *SwapStatement) Accept(visitor Visitor) Repr {
	return visitor.VisitSwapStatement(s)
}

// ExpressionStatement

type ExpressionStatement struct {
	Expression Expression
}

func (s *ExpressionStatement) StartPosition() Position {
	return s.Expression.StartPosition()
}

func (s *ExpressionStatement) EndPosition() Position {
	return s.Expression.EndPosition()
}

func (*ExpressionStatement) isStatement() {}

func (s *ExpressionStatement) Accept(visitor Visitor) Repr {
	return visitor.VisitExpressionStatement(s)
}

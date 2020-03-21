package sema

import (
	"github.com/dapperlabs/cadence/runtime/ast"
	"github.com/dapperlabs/cadence/runtime/errors"
)

//go:generate stringer -type=BinaryOperationKind

type BinaryOperationKind int

const (
	BinaryOperationKindUnknown BinaryOperationKind = iota
	BinaryOperationKindArithmetic
	BinaryOperationKindNonEqualityComparison
	BinaryOperationKindBooleanLogic
	BinaryOperationKindEquality
	BinaryOperationKindNilCoalescing
	BinaryOperationKindConcatenation
)

func binaryOperationKind(operation ast.Operation) BinaryOperationKind {
	switch operation {
	case ast.OperationPlus,
		ast.OperationMinus,
		ast.OperationMod,
		ast.OperationMul,
		ast.OperationDiv:

		return BinaryOperationKindArithmetic

	case ast.OperationLess,
		ast.OperationLessEqual,
		ast.OperationGreater,
		ast.OperationGreaterEqual:

		return BinaryOperationKindNonEqualityComparison

	case ast.OperationOr,
		ast.OperationAnd:

		return BinaryOperationKindBooleanLogic

	case ast.OperationEqual,
		ast.OperationUnequal:

		return BinaryOperationKindEquality

	case ast.OperationNilCoalesce:
		return BinaryOperationKindNilCoalescing

	case ast.OperationConcat:
		return BinaryOperationKindConcatenation
	}

	panic(errors.NewUnreachableError())
}

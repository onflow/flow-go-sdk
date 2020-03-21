package common

import (
	"github.com/dapperlabs/cadence/runtime/errors"
)

//go:generate stringer -type=OperationKind

type OperationKind int

const (
	OperationKindUnknown OperationKind = iota
	OperationKindUnary
	OperationKindBinary
	OperationKindTernary
)

func (k OperationKind) Name() string {
	switch k {
	case OperationKindUnary:
		return "unary"
	case OperationKindBinary:
		return "binary"
	case OperationKindTernary:
		return "ternary"
	}

	panic(errors.NewUnreachableError())
}

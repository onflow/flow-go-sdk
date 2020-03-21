package sema

import (
	"github.com/dapperlabs/cadence/runtime/errors"
)

//go:generate stringer -type=ResourceInvalidationKind

type ResourceInvalidationKind int

const (
	ResourceInvalidationKindUnknown ResourceInvalidationKind = iota
	ResourceInvalidationKindMove
	ResourceInvalidationKindDestroy
)

func (k ResourceInvalidationKind) Name() string {
	switch k {
	case ResourceInvalidationKindMove:
		return "move"
	case ResourceInvalidationKindDestroy:
		return "destroy"
	}

	panic(errors.NewUnreachableError())
}

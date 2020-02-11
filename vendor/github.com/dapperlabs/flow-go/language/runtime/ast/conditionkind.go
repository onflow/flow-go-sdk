package ast

import (
	"github.com/dapperlabs/flow-go/language/runtime/errors"
)

//go:generate stringer -type=ConditionKind

type ConditionKind int

const (
	ConditionKindUnknown ConditionKind = iota
	ConditionKindPre
	ConditionKindPost
)

func (k ConditionKind) Name() string {
	switch k {
	case ConditionKindPre:
		return "pre-condition"
	case ConditionKindPost:
		return "post-condition"
	}

	panic(errors.NewUnreachableError())
}

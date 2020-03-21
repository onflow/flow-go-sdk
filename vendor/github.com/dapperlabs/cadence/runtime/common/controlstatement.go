package common

import (
	"github.com/dapperlabs/cadence/runtime/errors"
)

//go:generate stringer -type=ControlStatement

type ControlStatement int

const (
	ControlStatementUnknown ControlStatement = iota
	ControlStatementBreak
	ControlStatementContinue
)

func (s ControlStatement) Symbol() string {
	switch s {
	case ControlStatementBreak:
		return "break"
	case ControlStatementContinue:
		return "continue"
	}

	panic(errors.NewUnreachableError())
}

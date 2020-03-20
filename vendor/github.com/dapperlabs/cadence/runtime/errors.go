package runtime

import (
	"fmt"
	"strings"

	"github.com/dapperlabs/cadence/runtime/errors"
	"github.com/dapperlabs/cadence/runtime/sema"
)

// Error is the containing type for all errors produced by the runtime.
type Error struct {
	Err error
}

func newError(err error) Error {
	return Error{Err: err}
}

func (e Error) Unwrap() error {
	return e.Err
}

func (e Error) Error() string {
	var sb strings.Builder
	sb.WriteString("Execution failed:\n")
	sb.WriteString(errors.UnrollChildErrors(e.Err))
	sb.WriteString("\n")
	return sb.String()
}

// InvalidTransactionCountError

type InvalidTransactionCountError struct {
	Count int
}

func (e InvalidTransactionCountError) Error() string {
	if e.Count == 0 {
		return "no transaction declared: expected 1, got 0"
	}

	return fmt.Sprintf(
		"multiple transactions declared: expected 1, got %d",
		e.Count,
	)
}

// InvalidTransactionParameterCountError

type InvalidTransactionParameterCountError struct {
	Expected int
	Actual   int
}

func (e InvalidTransactionParameterCountError) Error() string {
	return fmt.Sprintf(
		"parameter count mismatch for transaction: expected %d, got %d",
		e.Expected,
		e.Actual,
	)
}

// InvalidTransactionParameterTypeError

type InvalidTransactionParameterTypeError struct {
	Actual sema.Type
}

func (e InvalidTransactionParameterTypeError) Error() string {
	return fmt.Sprintf(
		"parameter type mismatch for transaction: expected `%s`, got `%s`",
		&sema.AuthAccountType{},
		e.Actual,
	)
}

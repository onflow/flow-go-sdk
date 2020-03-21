package parser

import (
	"fmt"
	"strings"

	"github.com/dapperlabs/cadence/runtime/ast"
	"github.com/dapperlabs/cadence/runtime/errors"
)

// Error

type Error struct {
	Errors []error
}

func (e Error) Error() string {
	var sb strings.Builder
	sb.WriteString("Parsing failed:\n")
	for _, err := range e.Errors {
		sb.WriteString(err.Error())
		if err, ok := err.(errors.SecondaryError); ok {
			sb.WriteString(". ")
			sb.WriteString(err.SecondaryError())
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (e Error) ChildErrors() []error {
	return e.Errors
}

// ParserError

type ParseError interface {
	error
	ast.HasPosition
	isParseError()
}

// SyntaxError

type SyntaxError struct {
	Pos     ast.Position
	Message string
}

func (*SyntaxError) isParseError() {}

func (e *SyntaxError) StartPosition() ast.Position {
	return e.Pos
}

func (e *SyntaxError) EndPosition() ast.Position {
	return e.Pos
}

func (e *SyntaxError) Error() string {
	return e.Message
}

// JuxtaposedUnaryOperatorsError

type JuxtaposedUnaryOperatorsError struct {
	Pos ast.Position
}

func (*JuxtaposedUnaryOperatorsError) isParseError() {}

func (e *JuxtaposedUnaryOperatorsError) StartPosition() ast.Position {
	return e.Pos
}

func (e *JuxtaposedUnaryOperatorsError) EndPosition() ast.Position {
	return e.Pos
}

func (e *JuxtaposedUnaryOperatorsError) Error() string {
	return "unary operators must not be juxtaposed; parenthesize inner expression"
}

// InvalidIntegerLiteralError

type InvalidIntegerLiteralError struct {
	Literal                   string
	IntegerLiteralKind        IntegerLiteralKind
	InvalidIntegerLiteralKind InvalidNumberLiteralKind
	ast.Range
}

func (*InvalidIntegerLiteralError) isParseError() {}

func (e *InvalidIntegerLiteralError) Error() string {
	if e.IntegerLiteralKind == IntegerLiteralKindUnknown {
		return fmt.Sprintf(
			"invalid integer literal `%s`: %s",
			e.Literal,
			e.InvalidIntegerLiteralKind.Description(),
		)
	}

	return fmt.Sprintf(
		"invalid %s integer literal `%s`: %s",
		e.IntegerLiteralKind.Name(),
		e.Literal,
		e.InvalidIntegerLiteralKind.Description(),
	)
}

func (e *InvalidIntegerLiteralError) SecondaryError() string {
	switch e.InvalidIntegerLiteralKind {
	case InvalidNumberLiteralKindUnknown:
		return ""
	case InvalidNumberLiteralKindLeadingUnderscore:
		return "remove the leading underscore"
	case InvalidNumberLiteralKindTrailingUnderscore:
		return "remove the trailing underscore"
	case InvalidNumberLiteralKindUnknownPrefix:
		return "did you mean `0x` (hexadecimal), `0b` (binary), or `0o` (octal)?"
	}

	panic(errors.NewUnreachableError())
}

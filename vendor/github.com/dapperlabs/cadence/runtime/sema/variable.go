package sema

import (
	"github.com/dapperlabs/cadence/runtime/ast"
	"github.com/dapperlabs/cadence/runtime/common"
)

type Variable struct {
	Identifier      string
	DeclarationKind common.DeclarationKind
	// Type is the type of the variable
	Type Type
	// Access is the access modifier
	Access ast.Access
	// IsConstant indicates if the variable is read-only
	IsConstant bool
	// Depth is the depth of scopes in which the variable was declared
	Depth int
	// ArgumentLabels are the argument labels that must be used in an invocation of the variable
	ArgumentLabels []string
	// Pos is the position where the variable was declared
	Pos *ast.Position
	// IsBaseValue indicates if the variable is a base value,
	// i.e. it is defined by the checker and not the program
	IsBaseValue bool
}

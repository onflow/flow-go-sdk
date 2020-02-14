package stdlib

import (
	"github.com/dapperlabs/flow-go/language/runtime/ast"
	"github.com/dapperlabs/flow-go/language/runtime/common"
	"github.com/dapperlabs/flow-go/language/runtime/sema"
)

type StandardLibraryValue struct {
	Name       string
	Type       sema.Type
	Kind       common.DeclarationKind
	IsConstant bool
}

func (v StandardLibraryValue) ValueDeclarationType() sema.Type {
	return v.Type
}

func (v StandardLibraryValue) ValueDeclarationKind() common.DeclarationKind {
	if v.IsConstant {
		return common.DeclarationKindConstant
	}
	return common.DeclarationKindVariable
}

func (StandardLibraryValue) ValueDeclarationPosition() ast.Position {
	return ast.Position{}
}

func (v StandardLibraryValue) ValueDeclarationIsConstant() bool {
	return v.IsConstant
}

func (StandardLibraryValue) ValueDeclarationArgumentLabels() []string {
	return nil
}

// StandardLibraryValues

type StandardLibraryValues []StandardLibraryValue

func (functions StandardLibraryValues) ToValueDeclarations() map[string]sema.ValueDeclaration {
	valueDeclarations := make(map[string]sema.ValueDeclaration, len(functions))
	for _, function := range functions {
		valueDeclarations[function.Name] = function
	}
	return valueDeclarations
}

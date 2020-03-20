package sema

import (
	"github.com/dapperlabs/cadence/runtime/ast"
	"github.com/dapperlabs/cadence/runtime/common"
)

func (checker *Checker) VisitPathExpression(expression *ast.PathExpression) ast.Repr {

	// Check that the domain is valid

	domain := expression.Domain

	if _, ok := common.AllPathDomainsByIdentifier[domain.Identifier]; !ok {
		checker.report(
			&InvalidPathDomainError{
				ActualDomain: domain.Identifier,
				Range:        ast.NewRangeFromPositioned(domain),
			},
		)
	}

	return &PathType{}
}

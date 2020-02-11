package sema

import "github.com/dapperlabs/flow-go/language/runtime/ast"

func (checker *Checker) VisitArrayExpression(expression *ast.ArrayExpression) ast.Repr {

	// visit all elements, ensure they are all the same type

	var elementType Type

	argumentTypes := make([]Type, len(expression.Values))

	for i, value := range expression.Values {
		valueType := value.Accept(checker).(Type)

		argumentTypes[i] = valueType

		checker.checkVariableMove(value)
		checker.checkResourceMoveOperation(value, valueType)

		// infer element type from first element
		// TODO: find common super type?
		if elementType == nil {
			elementType = valueType
		} else if !valueType.IsInvalidType() &&
			!IsSubType(valueType, elementType) {

			checker.report(
				&TypeMismatchError{
					ExpectedType: elementType,
					ActualType:   valueType,
					Range:        ast.NewRangeFromPositioned(value),
				},
			)
		}
	}

	checker.Elaboration.ArrayExpressionArgumentTypes[expression] = argumentTypes

	if elementType == nil {
		elementType = &NeverType{}
	}

	checker.Elaboration.ArrayExpressionElementType[expression] = elementType

	return &VariableSizedType{
		Type: elementType,
	}
}

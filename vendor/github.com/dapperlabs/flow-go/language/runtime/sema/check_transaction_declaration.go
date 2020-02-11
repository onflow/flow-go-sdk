package sema

import (
	"github.com/dapperlabs/flow-go/language/runtime/ast"
	"github.com/dapperlabs/flow-go/language/runtime/common"
	"github.com/dapperlabs/flow-go/language/runtime/errors"
)

func (checker *Checker) VisitTransactionDeclaration(declaration *ast.TransactionDeclaration) ast.Repr {
	transactionType := checker.Elaboration.TransactionDeclarationTypes[declaration]
	if transactionType == nil {
		panic(errors.NewUnreachableError())
	}

	checker.containerTypes[transactionType] = true
	defer func() {
		checker.containerTypes[transactionType] = false
	}()

	fieldMembers := map[*Member]*ast.FieldDeclaration{}

	for _, field := range declaration.Fields {
		fieldName := field.Identifier.Identifier
		member := transactionType.Members[fieldName]
		fieldMembers[member] = field
	}

	checker.checkTransactionFields(declaration)
	checker.checkTransactionBlocks(declaration)

	// enter a new scope for this transaction
	checker.enterValueScope()
	defer checker.leaveValueScope(true)

	checker.declareSelfValue(transactionType)

	checker.visitTransactionPrepareFunction(declaration.Prepare, transactionType, fieldMembers)

	if declaration.PreConditions != nil {
		checker.visitConditions(*declaration.PreConditions)
	}

	checker.visitWithPostConditions(
		declaration.PostConditions,
		&VoidType{},
		func() {
			checker.withSelfResourceInvalidationAllowed(func() {
				checker.visitTransactionExecuteFunction(declaration.Execute, transactionType)
			})
		},
	)

	checker.checkResourceFieldsInvalidated(transactionType.String(), transactionType.Members)

	return nil
}

// checkTransactionFields validates the field declarations for a transaction.
//
func (checker *Checker) checkTransactionFields(declaration *ast.TransactionDeclaration) {
	for _, field := range declaration.Fields {
		if field.Access != ast.AccessNotSpecified {
			checker.report(
				&InvalidTransactionFieldAccessModifierError{
					Name:   field.Identifier.Identifier,
					Access: field.Access,
					Pos:    field.StartPosition(),
				},
			)
		}
	}
}

// checkTransactionBlocks checks that a transaction contains the required prepare and execute blocks.
//
// An execute block is always required, but a prepare block is only required if fields are present.
//
func (checker *Checker) checkTransactionBlocks(declaration *ast.TransactionDeclaration) {
	if declaration.Prepare != nil {
		// parser allows any identifier so it must be checked here
		prepareIdentifier := declaration.Prepare.Identifier
		if prepareIdentifier.Identifier != common.DeclarationKindPrepare.Keywords() {
			checker.report(&InvalidTransactionBlockError{
				Name: prepareIdentifier.Identifier,
				Pos:  prepareIdentifier.Pos,
			})
		}
	} else if len(declaration.Fields) != 0 {
		// report an error if fields are defined but no prepare statement exists
		// note: field initialization is checked later

		// report error for first field
		firstField := declaration.Fields[0]

		checker.report(
			&TransactionMissingPrepareError{
				FirstFieldName: firstField.Identifier.Identifier,
				FirstFieldPos:  firstField.Identifier.Pos,
			},
		)
	}

	if declaration.Execute != nil {
		// parser allows any identifier so it must be checked here
		executeIdentifier := declaration.Execute.Identifier
		if executeIdentifier.Identifier != common.DeclarationKindExecute.Keywords() {
			checker.report(&InvalidTransactionBlockError{
				Name: executeIdentifier.Identifier,
				Pos:  executeIdentifier.Pos,
			})
		}
	}
}

// visitTransactionPrepareFunction visits and checks the prepare function of a transaction.
//
func (checker *Checker) visitTransactionPrepareFunction(
	prepareFunction *ast.SpecialFunctionDeclaration,
	transactionType *TransactionType,
	fieldMembers map[*Member]*ast.FieldDeclaration,
) {
	if prepareFunction == nil {
		return
	}

	initializationInfo := NewInitializationInfo(transactionType, fieldMembers)

	prepareFunctionType := transactionType.PrepareFunctionType().InvocationFunctionType()

	checker.checkFunction(
		prepareFunction.ParameterList,
		ast.Position{},
		prepareFunctionType,
		prepareFunction.FunctionBlock,
		true,
		initializationInfo,
		true,
	)

	checker.checkTransactionPrepareFunctionParameters(
		prepareFunction.ParameterList,
		prepareFunctionType.Parameters,
	)
}

// checkTransactionPrepareFunctionParameters checks that the parameters are each of type Account.
//
func (checker *Checker) checkTransactionPrepareFunctionParameters(
	parameterList *ast.ParameterList,
	parameters []*Parameter,
) {
	for i, parameter := range parameterList.Parameters {
		parameterType := parameters[i].TypeAnnotation.Type

		if !IsSubType(parameterType, &AccountType{}) {
			checker.report(
				&InvalidTransactionPrepareParameterTypeError{
					Type:  parameterType,
					Range: ast.NewRangeFromPositioned(parameter.TypeAnnotation),
				},
			)
		}
	}

}

// visitTransactionExecuteFunction visits and checks the execute function of a transaction.
func (checker *Checker) visitTransactionExecuteFunction(
	executeFunction *ast.SpecialFunctionDeclaration,
	transactionType *TransactionType,
) {
	if executeFunction == nil {
		return
	}

	executeFunctionType := transactionType.ExecuteFunctionType().InvocationFunctionType()

	checker.checkFunction(
		&ast.ParameterList{},
		ast.Position{},
		executeFunctionType,
		executeFunction.FunctionBlock,
		true,
		nil,
		true,
	)
}

func (checker *Checker) declareTransactionDeclaration(declaration *ast.TransactionDeclaration) {
	transactionType := &TransactionType{}

	members, origins := checker.nonEventMembersAndOrigins(
		transactionType,
		declaration.Fields,
		nil,
		true,
	)

	checker.memberOrigins[transactionType] = origins

	var prepareParameters []*Parameter
	if declaration.Prepare != nil {
		prepareParameters = checker.parameters(declaration.Prepare.ParameterList)
	}

	transactionType.Members = members
	transactionType.prepareParameters = prepareParameters

	checker.Elaboration.TransactionDeclarationTypes[declaration] = transactionType
	checker.TransactionTypes = append(checker.TransactionTypes, transactionType)
}

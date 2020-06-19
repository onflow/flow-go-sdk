package client

import (
	"fmt"
)

const errorMessagePrefix = "client: "

func errorMessage(format string, a ...interface{}) string {
	return errorMessagePrefix + fmt.Sprintf(format, a...)
}

const (
	entityBlock             = "flow.Block"
	entityBlockHeader       = "flow.BlockHeader"
	entityCollection        = "flow.Collection"
	entityTransaction       = "flow.Transaction"
	entityTransactionResult = "flow.TransactionResult"
	entityAccount           = "flow.Account"
	entityEvent             = "flow.Event"
	entityCadenceValue      = "cadence.Value"
)

type EntityToMessageError struct {
	Entity string
	Err    error
}

func (e EntityToMessageError) Error() string {
	return errorMessage("failed to construct protobuf message from %s entity: %s", e.Entity, e.Err.Error())
}

func newEntityToMessageError(entity string, err error) EntityToMessageError {
	return EntityToMessageError{
		Entity: entity,
		Err:    err,
	}
}

type MessageToEntityError struct {
	Entity string
	Err    error
}

func (e MessageToEntityError) Error() string {
	return errorMessage("failed to construct %s entity from protobuf value: %s", e.Entity, e.Err.Error())
}

func newMessageToEntityError(entity string, err error) MessageToEntityError {
	return MessageToEntityError{
		Entity: entity,
		Err:    err,
	}
}

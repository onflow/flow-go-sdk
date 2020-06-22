package client

import (
	"fmt"

	"google.golang.org/grpc/status"
)

const errorMessagePrefix = "client: "

func errorMessage(format string, a ...interface{}) string {
	return errorMessagePrefix + fmt.Sprintf(format, a...)
}

type RPCError struct {
	GRPCError error
}

func newRPCError(gRPCErr error) RPCError {
	return RPCError{GRPCError: gRPCErr}
}

func (e RPCError) Error() string {
	return errorMessage(e.GRPCError.Error())
}

func (e RPCError) Unwrap() error {
	return e.GRPCError
}

func (e RPCError) GRPCStatus() *status.Status {
	s, _ := status.FromError(e.GRPCError)
	return s
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

func newEntityToMessageError(entity string, err error) EntityToMessageError {
	return EntityToMessageError{
		Entity: entity,
		Err:    err,
	}
}

func (e EntityToMessageError) Error() string {
	return errorMessage("failed to construct protobuf message from %s entity: %s", e.Entity, e.Err.Error())
}

func (e EntityToMessageError) Unwrap() error {
	return e.Err
}

type MessageToEntityError struct {
	Entity string
	Err    error
}

func newMessageToEntityError(entity string, err error) MessageToEntityError {
	return MessageToEntityError{
		Entity: entity,
		Err:    err,
	}
}

func (e MessageToEntityError) Error() string {
	return errorMessage("failed to construct %s entity from protobuf value: %s", e.Entity, e.Err.Error())
}

func (e MessageToEntityError) Unwrap() error {
	return e.Err
}

/*
 * Flow Go SDK
 *
 * Copyright Flow Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package grpc

import (
	"fmt"

	"google.golang.org/grpc/status"
)

const errorMessagePrefix = "client: "

// An RPCError is an error returned by an RPC call to an Access API.
//
// An RPC error can be unwrapped to produce the original gRPC error.
type RPCError struct {
	GRPCErr error
}

func newRPCError(gRPCErr error) RPCError {
	return RPCError{GRPCErr: gRPCErr}
}

func (e RPCError) Error() string {
	return errorMessagePrefix + e.GRPCErr.Error()
}

func (e RPCError) Unwrap() error {
	return e.GRPCErr
}

// GRPCStatus returns the gRPC status for this error.
//
// This function satisfies the interface defined in the status.FromError function.
func (e RPCError) GRPCStatus() *status.Status {
	s, _ := status.FromError(e.GRPCErr)
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

// An EntityToMessageError indicates that an entity could not be converted to a protobuf message.
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
	return errorMessagePrefix + fmt.Sprintf(
		"failed to construct protobuf message from %s entity: %s",
		e.Entity,
		e.Err.Error(),
	)
}

func (e EntityToMessageError) Unwrap() error {
	return e.Err
}

// A MessageToEntityError indicates that a protobuf message could not be converted to an SDK entity.
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
	return errorMessagePrefix + fmt.Sprintf(
		"failed to construct %s entity from protobuf value: %s",
		e.Entity,
		e.Err.Error(),
	)
}

func (e MessageToEntityError) Unwrap() error {
	return e.Err
}

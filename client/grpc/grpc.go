/*
 * Flow Go SDK
 *
 * Copyright 2019 Dapper Labs, Inc.
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

//go:generate go run github.com/vektra/mockery/cmd/mockery --name RPCClient --filename mock_client_test.go --inpkg

import (
	"context"

	"github.com/onflow/cadence"
	"github.com/onflow/flow/protobuf/go/flow/access"
	"google.golang.org/grpc"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client/convert"
)

// RPCClient is an RPC client for the Flow Access API.
type RPCClient interface {
	access.AccessAPIClient
}

// GRPCClient is a gRPC client for the Flow Access API exposing all grpc specific methods.
type GRPCClient struct {
	rpcClient RPCClient
	close     func() error
}

// NewGRPCClient creates a new gRPC handler for network communication.
func NewGRPCClient(url string, opts ...grpc.DialOption) (*GRPCClient, error) {
	conn, err := grpc.Dial(url, opts...)
	if err != nil {
		return nil, err
	}

	grpcClient := access.NewAccessAPIClient(conn)

	return &GRPCClient{
		rpcClient: grpcClient,
		close:     func() error { return conn.Close() },
	}, nil
}

// NewFromRPCClient initializes a Flow client using a pre-configured gRPC provider.
func NewFromRPCClient(rpcClient RPCClient) *GRPCClient {
	return &GRPCClient{
		rpcClient: rpcClient,
		close:     func() error { return nil },
	}
}

// Close closes the client connection.
func (c *GRPCClient) Close() error {
	return c.close()
}

func (c *GRPCClient) Ping(ctx context.Context, opts ...grpc.CallOption) error {
	_, err := c.rpcClient.Ping(ctx, &access.PingRequest{}, opts...)
	return err
}

func (c *GRPCClient) GetLatestBlockHeader(
	ctx context.Context,
	isSealed bool,
	opts ...grpc.CallOption,
) (*flow.BlockHeader, error) {

	req := &access.GetLatestBlockHeaderRequest{
		IsSealed: isSealed,
	}

	res, err := c.rpcClient.GetLatestBlockHeader(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	return getBlockHeaderResult(res)
}

func (c *GRPCClient) GetBlockHeaderByID(
	ctx context.Context,
	blockID flow.Identifier,
	opts ...grpc.CallOption,
) (*flow.BlockHeader, error) {
	req := &access.GetBlockHeaderByIDRequest{
		Id: blockID.Bytes(),
	}

	res, err := c.rpcClient.GetBlockHeaderByID(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	return getBlockHeaderResult(res)
}

func (c *GRPCClient) GetBlockHeaderByHeight(
	ctx context.Context,
	height uint64,
	opts ...grpc.CallOption,
) (*flow.BlockHeader, error) {
	req := &access.GetBlockHeaderByHeightRequest{
		Height: height,
	}

	res, err := c.rpcClient.GetBlockHeaderByHeight(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	return getBlockHeaderResult(res)
}

func getBlockHeaderResult(res *access.BlockHeaderResponse) (*flow.BlockHeader, error) {
	header, err := convert.MessageToBlockHeader(res.GetBlock())
	if err != nil {
		return nil, newMessageToEntityError(entityBlockHeader, err)
	}

	return &header, nil
}

func (c *GRPCClient) GetLatestBlock(
	ctx context.Context,
	isSealed bool,
	opts ...grpc.CallOption,
) (*flow.Block, error) {
	req := &access.GetLatestBlockRequest{
		IsSealed: isSealed,
	}

	res, err := c.rpcClient.GetLatestBlock(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	return getBlockResult(res)
}

func (c *GRPCClient) GetBlockByID(
	ctx context.Context,
	blockID flow.Identifier,
	opts ...grpc.CallOption,
) (*flow.Block, error) {
	req := &access.GetBlockByIDRequest{
		Id: blockID.Bytes(),
	}

	res, err := c.rpcClient.GetBlockByID(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	return getBlockResult(res)
}

func (c *GRPCClient) GetBlockByHeight(
	ctx context.Context,
	height uint64,
	opts ...grpc.CallOption,
) (*flow.Block, error) {
	req := &access.GetBlockByHeightRequest{
		Height: height,
	}

	res, err := c.rpcClient.GetBlockByHeight(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	return getBlockResult(res)
}

func getBlockResult(res *access.BlockResponse) (*flow.Block, error) {
	block, err := convert.MessageToBlock(res.GetBlock())
	if err != nil {
		return nil, newMessageToEntityError(entityBlock, err)
	}

	return &block, nil
}

func (c *GRPCClient) GetCollection(
	ctx context.Context,
	colID flow.Identifier,
	opts ...grpc.CallOption,
) (*flow.Collection, error) {
	req := &access.GetCollectionByIDRequest{
		Id: colID.Bytes(),
	}

	res, err := c.rpcClient.GetCollectionByID(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	result, err := convert.MessageToCollection(res.GetCollection())
	if err != nil {
		return nil, newMessageToEntityError(entityCollection, err)
	}

	return &result, nil
}

func (c *GRPCClient) SendTransaction(
	ctx context.Context,
	tx flow.Transaction,
	opts ...grpc.CallOption,
) error {
	txMsg, err := convert.TransactionToMessage(tx)
	if err != nil {
		return newEntityToMessageError(entityTransaction, err)
	}

	req := &access.SendTransactionRequest{
		Transaction: txMsg,
	}

	_, err = c.rpcClient.SendTransaction(ctx, req, opts...)
	if err != nil {
		return newRPCError(err)
	}

	return nil
}

func (c *GRPCClient) GetTransaction(
	ctx context.Context,
	txID flow.Identifier,
	opts ...grpc.CallOption,
) (*flow.Transaction, error) {
	req := &access.GetTransactionRequest{
		Id: txID.Bytes(),
	}

	res, err := c.rpcClient.GetTransaction(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	result, err := convert.MessageToTransaction(res.GetTransaction())
	if err != nil {
		return nil, newMessageToEntityError(entityTransaction, err)
	}

	return &result, nil
}

func (c *GRPCClient) GetTransactionResult(
	ctx context.Context,
	txID flow.Identifier,
	opts ...grpc.CallOption,
) (*flow.TransactionResult, error) {
	req := &access.GetTransactionRequest{
		Id: txID.Bytes(),
	}

	res, err := c.rpcClient.GetTransactionResult(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	result, err := convert.MessageToTransactionResult(res)
	if err != nil {
		return nil, newMessageToEntityError(entityTransactionResult, err)
	}

	return &result, nil
}

func (c *GRPCClient) GetAccount(ctx context.Context, address flow.Address, opts ...grpc.CallOption) (*flow.Account, error) {
	return c.GetAccountAtLatestBlock(ctx, address, opts...)
}

func (c *GRPCClient) GetAccountAtLatestBlock(
	ctx context.Context,
	address flow.Address,
	opts ...grpc.CallOption,
) (*flow.Account, error) {
	req := &access.GetAccountAtLatestBlockRequest{
		Address: address.Bytes(),
	}

	res, err := c.rpcClient.GetAccountAtLatestBlock(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	account, err := convert.MessageToAccount(res.GetAccount())
	if err != nil {
		return nil, newMessageToEntityError(entityAccount, err)
	}

	return &account, nil
}

func (c *GRPCClient) GetAccountAtBlockHeight(
	ctx context.Context,
	address flow.Address,
	blockHeight uint64,
	opts ...grpc.CallOption,
) (*flow.Account, error) {
	req := &access.GetAccountAtBlockHeightRequest{
		Address:     address.Bytes(),
		BlockHeight: blockHeight,
	}

	res, err := c.rpcClient.GetAccountAtBlockHeight(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	account, err := convert.MessageToAccount(res.GetAccount())
	if err != nil {
		return nil, newMessageToEntityError(entityAccount, err)
	}

	return &account, nil
}

func (c *GRPCClient) ExecuteScriptAtLatestBlock(
	ctx context.Context,
	script []byte,
	arguments []cadence.Value,
	opts ...grpc.CallOption,
) (cadence.Value, error) {

	args, err := convert.CadenceValuesToMessages(arguments)
	if err != nil {
		return nil, newEntityToMessageError(entityCadenceValue, err)
	}

	req := &access.ExecuteScriptAtLatestBlockRequest{
		Script:    script,
		Arguments: args,
	}

	res, err := c.rpcClient.ExecuteScriptAtLatestBlock(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	return executeScriptResult(res)
}

func (c *GRPCClient) ExecuteScriptAtBlockID(
	ctx context.Context,
	blockID flow.Identifier,
	script []byte,
	arguments []cadence.Value,
	opts ...grpc.CallOption,
) (cadence.Value, error) {

	args, err := convert.CadenceValuesToMessages(arguments)
	if err != nil {
		return nil, newEntityToMessageError(entityCadenceValue, err)
	}

	req := &access.ExecuteScriptAtBlockIDRequest{
		BlockId:   blockID.Bytes(),
		Script:    script,
		Arguments: args,
	}

	res, err := c.rpcClient.ExecuteScriptAtBlockID(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	return executeScriptResult(res)
}

func (c *GRPCClient) ExecuteScriptAtBlockHeight(
	ctx context.Context,
	height uint64,
	script []byte,
	arguments []cadence.Value,
	opts ...grpc.CallOption,
) (cadence.Value, error) {

	args, err := convert.CadenceValuesToMessages(arguments)
	if err != nil {
		return nil, newEntityToMessageError(entityCadenceValue, err)
	}

	req := &access.ExecuteScriptAtBlockHeightRequest{
		BlockHeight: height,
		Script:      script,
		Arguments:   args,
	}

	res, err := c.rpcClient.ExecuteScriptAtBlockHeight(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	return executeScriptResult(res)
}

func executeScriptResult(res *access.ExecuteScriptResponse) (cadence.Value, error) {
	value, err := convert.MessageToCadenceValue(res.GetValue())
	if err != nil {
		return nil, newMessageToEntityError(entityCadenceValue, err)
	}

	return value, nil
}

// EventRangeQuery defines a query for Flow events.
type EventRangeQuery struct {
	// The event type to search for.
	Type string
	// The block height to begin looking for events (inclusive).
	StartHeight uint64
	// The block height to end looking for events (inclusive).
	EndHeight uint64
}

func (c *GRPCClient) GetEventsForHeightRange(
	ctx context.Context,
	query EventRangeQuery,
	opts ...grpc.CallOption,
) ([]flow.BlockEvents, error) {
	req := &access.GetEventsForHeightRangeRequest{
		Type:        query.Type,
		StartHeight: query.StartHeight,
		EndHeight:   query.EndHeight,
	}

	res, err := c.rpcClient.GetEventsForHeightRange(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	return getEventsResult(res)
}

func (c *GRPCClient) GetEventsForBlockIDs(
	ctx context.Context,
	eventType string,
	blockIDs []flow.Identifier,
	opts ...grpc.CallOption,
) ([]flow.BlockEvents, error) {
	req := &access.GetEventsForBlockIDsRequest{
		Type:     eventType,
		BlockIds: convert.IdentifiersToMessages(blockIDs),
	}

	res, err := c.rpcClient.GetEventsForBlockIDs(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	return getEventsResult(res)
}

func getEventsResult(res *access.EventsResponse) ([]flow.BlockEvents, error) {
	resultMessages := res.GetResults()

	results := make([]flow.BlockEvents, len(resultMessages))
	for i, result := range resultMessages {
		eventMessages := result.GetEvents()

		events := make([]flow.Event, len(eventMessages))

		for i, m := range eventMessages {
			evt, err := convert.MessageToEvent(m)
			if err != nil {
				return nil, newMessageToEntityError(entityEvent, err)
			}

			events[i] = evt
		}

		blockTimestamp := result.BlockTimestamp.AsTime()
		results[i] = flow.BlockEvents{
			BlockID:        flow.HashToID(result.GetBlockId()),
			Height:         result.GetBlockHeight(),
			BlockTimestamp: blockTimestamp,
			Events:         events,
		}
	}

	return results, nil
}

func (c *GRPCClient) GetLatestProtocolStateSnapshot(ctx context.Context, opts ...grpc.CallOption) ([]byte, error) {
	res, err := c.rpcClient.GetLatestProtocolStateSnapshot(ctx, &access.GetLatestProtocolStateSnapshotRequest{}, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	return res.GetSerializedSnapshot(), nil
}

func (c *GRPCClient) GetExecutionResultForBlockID(ctx context.Context, blockID flow.Identifier, opts ...grpc.CallOption) (*flow.ExecutionResult, error) {
	er, err := c.rpcClient.GetExecutionResultForBlockID(ctx, &access.GetExecutionResultForBlockIDRequest{
		BlockId: convert.IdentifierToMessage(blockID),
	}, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	chunks := make([]*flow.Chunk, len(er.ExecutionResult.Chunks))
	serviceEvents := make([]*flow.ServiceEvent, len(er.ExecutionResult.ServiceEvents))

	for i, chunk := range er.ExecutionResult.Chunks {
		chunks[i] = &flow.Chunk{
			CollectionIndex:      uint(chunk.CollectionIndex),
			StartState:           flow.BytesToStateCommitment(chunk.StartState),
			EventCollection:      flow.BytesToHash(chunk.EventCollection),
			BlockID:              flow.BytesToID(chunk.BlockId),
			TotalComputationUsed: chunk.TotalComputationUsed,
			NumberOfTransactions: uint16(chunk.NumberOfTransactions),
			Index:                chunk.Index,
			EndState:             flow.BytesToStateCommitment(chunk.EndState),
		}
	}

	for i, serviceEvent := range er.ExecutionResult.ServiceEvents {
		serviceEvents[i] = &flow.ServiceEvent{
			Type:    serviceEvent.Type,
			Payload: serviceEvent.Payload,
		}
	}

	return &flow.ExecutionResult{
		PreviousResultID: flow.BytesToID(er.ExecutionResult.PreviousResultId),
		BlockID:          flow.BytesToID(er.ExecutionResult.BlockId),
		Chunks:           chunks,
		ServiceEvents:    serviceEvents,
	}, nil
}

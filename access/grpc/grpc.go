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

//go:generate go run github.com/vektra/mockery/cmd/mockery --name RPCClient --structname MockRPCClient --output mocks
//go:generate go run github.com/vektra/mockery/cmd/mockery --name ExecutionDataRPCClient --structname MockExecutionDataRPCClient --output mocks

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"

	"github.com/onflow/cadence"
	"github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow/protobuf/go/flow/access"
	"github.com/onflow/flow/protobuf/go/flow/entities"
	"github.com/onflow/flow/protobuf/go/flow/executiondata"

	"github.com/onflow/flow-go-sdk"
)

// RPCClient is an RPC client for the Flow Access API.
type RPCClient interface {
	access.AccessAPIClient
}

// ExecutionDataRPCClient is an RPC client for the Flow ExecutionData API.
type ExecutionDataRPCClient interface {
	executiondata.ExecutionDataAPIClient
}

type SubscribeOption func(*SubscribeConfig)

type SubscribeConfig struct {
	heartbeatInterval uint64
	grpcOpts          []grpc.CallOption
}

func DefaultSubscribeConfig() *SubscribeConfig {
	return &SubscribeConfig{
		heartbeatInterval: 100,
	}
}

func WithHeartbeatInterval(interval uint64) SubscribeOption {
	return func(config *SubscribeConfig) {
		config.heartbeatInterval = interval
	}
}

func WithGRPCOptions(grpcOpts ...grpc.CallOption) SubscribeOption {
	return func(config *SubscribeConfig) {
		config.grpcOpts = grpcOpts
	}
}

// BaseClient is a gRPC client for the Flow Access API exposing all grpc specific methods.
//
// Use this client if you need advance access to the HTTP API. If you
// don't require special methods use the Client instead.
type BaseClient struct {
	rpcClient           RPCClient
	executionDataClient ExecutionDataRPCClient
	close               func() error
	jsonOptions         []json.Option
}

// NewBaseClient creates a new gRPC handler for network communication.
func NewBaseClient(url string, opts ...grpc.DialOption) (*BaseClient, error) {
	conn, err := grpc.Dial(url, opts...)
	if err != nil {
		return nil, err
	}

	grpcClient := access.NewAccessAPIClient(conn)

	execDataClient := executiondata.NewExecutionDataAPIClient(conn)

	return &BaseClient{
		rpcClient:           grpcClient,
		executionDataClient: execDataClient,
		close:               func() error { return conn.Close() },
		jsonOptions:         []json.Option{json.WithAllowUnstructuredStaticTypes(true)},
	}, nil
}

// NewFromRPCClient initializes a Flow client using a pre-configured gRPC provider.
func NewFromRPCClient(rpcClient RPCClient) *BaseClient {
	return &BaseClient{
		rpcClient: rpcClient,
		close:     func() error { return nil },
	}
}

// NewFromExecutionDataRPCClient initializes a Flow client using a pre-configured gRPC provider.
func NewFromExecutionDataRPCClient(rpcClient ExecutionDataRPCClient) *BaseClient {
	return &BaseClient{
		executionDataClient: rpcClient,
		close:               func() error { return nil },
	}
}

func (c *BaseClient) SetJSONOptions(options []json.Option) {
	c.jsonOptions = options
}

// Close closes the client connection.
func (c *BaseClient) Close() error {
	return c.close()
}

func (c *BaseClient) Ping(ctx context.Context, opts ...grpc.CallOption) error {
	_, err := c.rpcClient.Ping(ctx, &access.PingRequest{}, opts...)
	return err
}

func (c *BaseClient) GetNetworkParameters(ctx context.Context, opts ...grpc.CallOption) (*flow.NetworkParameters, error) {
	res, err := c.rpcClient.GetNetworkParameters(ctx, &access.GetNetworkParametersRequest{}, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}
	return &flow.NetworkParameters{
		ChainID: flow.ChainID(res.ChainId),
	}, nil
}

func (c *BaseClient) GetLatestBlockHeader(
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

func (c *BaseClient) GetBlockHeaderByID(
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

func (c *BaseClient) GetBlockHeaderByHeight(
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
	header, err := messageToBlockHeader(res.GetBlock())
	if err != nil {
		return nil, newMessageToEntityError(entityBlockHeader, err)
	}
	header.Status = flow.BlockStatus(res.GetBlockStatus())
	return &header, nil
}

func (c *BaseClient) GetLatestBlock(
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

func (c *BaseClient) GetBlockByID(
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

func (c *BaseClient) GetBlockByHeight(
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
	block, err := messageToBlock(res.GetBlock())
	if err != nil {
		return nil, newMessageToEntityError(entityBlock, err)
	}
	block.BlockHeader.Status = flow.BlockStatus(res.GetBlockStatus())
	return &block, nil
}

func (c *BaseClient) GetCollection(
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

	result, err := messageToCollection(res.GetCollection())
	if err != nil {
		return nil, newMessageToEntityError(entityCollection, err)
	}

	return &result, nil
}

func (c *BaseClient) SendTransaction(
	ctx context.Context,
	tx flow.Transaction,
	opts ...grpc.CallOption,
) error {
	txMsg, err := transactionToMessage(tx)
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

func (c *BaseClient) GetTransaction(
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

	result, err := messageToTransaction(res.GetTransaction())
	if err != nil {
		return nil, newMessageToEntityError(entityTransaction, err)
	}

	return &result, nil
}

func (c *BaseClient) GetTransactionsByBlockID(
	ctx context.Context,
	blockID flow.Identifier,
	opts ...grpc.CallOption,
) ([]*flow.Transaction, error) {
	req := &access.GetTransactionsByBlockIDRequest{
		BlockId: blockID.Bytes(),
	}

	res, err := c.rpcClient.GetTransactionsByBlockID(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	unparsedResults := res.GetTransactions()
	results := make([]*flow.Transaction, 0, len(unparsedResults))
	for _, result := range unparsedResults {
		parsed, err := messageToTransaction(result)
		if err != nil {
			return nil, newMessageToEntityError(entityTransaction, err)
		}
		results = append(results, &parsed)
	}

	return results, nil
}

func (c *BaseClient) GetTransactionResult(
	ctx context.Context,
	txID flow.Identifier,
	opts ...grpc.CallOption,
) (*flow.TransactionResult, error) {
	req := &access.GetTransactionRequest{
		Id:                   txID.Bytes(),
		EventEncodingVersion: entities.EventEncodingVersion_CCF_V0,
	}

	res, err := c.rpcClient.GetTransactionResult(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	result, err := messageToTransactionResult(res, c.jsonOptions)
	if err != nil {
		return nil, newMessageToEntityError(entityTransactionResult, err)
	}

	return &result, nil
}

func (c *BaseClient) GetTransactionResultByIndex(
	ctx context.Context,
	blockID flow.Identifier,
	index uint32,
	opts ...grpc.CallOption,
) (*flow.TransactionResult, error) {

	req := &access.GetTransactionByIndexRequest{
		BlockId:              blockID.Bytes(),
		Index:                index,
		EventEncodingVersion: entities.EventEncodingVersion_CCF_V0,
	}

	res, err := c.rpcClient.GetTransactionResultByIndex(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	parsed, err := messageToTransactionResult(res, c.jsonOptions)
	if err != nil {
		return nil, newMessageToEntityError(entityTransactionResult, err)
	}
	return &parsed, nil
}

func (c *BaseClient) GetTransactionResultsByBlockID(
	ctx context.Context,
	blockID flow.Identifier,
	opts ...grpc.CallOption,
) ([]*flow.TransactionResult, error) {

	req := &access.GetTransactionsByBlockIDRequest{
		BlockId:              blockID.Bytes(),
		EventEncodingVersion: entities.EventEncodingVersion_CCF_V0,
	}

	res, err := c.rpcClient.GetTransactionResultsByBlockID(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	unparsedResults := res.GetTransactionResults()
	results := make([]*flow.TransactionResult, 0, len(unparsedResults))
	for _, result := range unparsedResults {
		parsed, err := messageToTransactionResult(result, c.jsonOptions)
		if err != nil {
			return nil, newMessageToEntityError(entityTransactionResult, err)
		}
		results = append(results, &parsed)
	}

	return results, nil
}

func (c *BaseClient) GetAccount(ctx context.Context, address flow.Address, opts ...grpc.CallOption) (*flow.Account, error) {
	return c.GetAccountAtLatestBlock(ctx, address, opts...)
}

func (c *BaseClient) GetAccountAtLatestBlock(
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

	account, err := messageToAccount(res.GetAccount())
	if err != nil {
		return nil, newMessageToEntityError(entityAccount, err)
	}

	return &account, nil
}

func (c *BaseClient) GetAccountAtBlockHeight(
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

	account, err := messageToAccount(res.GetAccount())
	if err != nil {
		return nil, newMessageToEntityError(entityAccount, err)
	}

	return &account, nil
}

func (c *BaseClient) ExecuteScriptAtLatestBlock(
	ctx context.Context,
	script []byte,
	arguments []cadence.Value,
	opts ...grpc.CallOption,
) (cadence.Value, error) {

	args, err := cadenceValuesToMessages(arguments)
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

	return executeScriptResult(res, c.jsonOptions)
}

func (c *BaseClient) ExecuteScriptAtBlockID(
	ctx context.Context,
	blockID flow.Identifier,
	script []byte,
	arguments []cadence.Value,
	opts ...grpc.CallOption,
) (cadence.Value, error) {

	args, err := cadenceValuesToMessages(arguments)
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

	return executeScriptResult(res, c.jsonOptions)
}

func (c *BaseClient) ExecuteScriptAtBlockHeight(
	ctx context.Context,
	height uint64,
	script []byte,
	arguments []cadence.Value,
	opts ...grpc.CallOption,
) (cadence.Value, error) {

	args, err := cadenceValuesToMessages(arguments)
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

	return executeScriptResult(res, c.jsonOptions)
}

func executeScriptResult(res *access.ExecuteScriptResponse, options []json.Option) (cadence.Value, error) {
	value, err := messageToCadenceValue(res.GetValue(), options)
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

func (c *BaseClient) GetEventsForHeightRange(
	ctx context.Context,
	query EventRangeQuery,
	opts ...grpc.CallOption,
) ([]flow.BlockEvents, error) {
	req := &access.GetEventsForHeightRangeRequest{
		Type:                 query.Type,
		StartHeight:          query.StartHeight,
		EndHeight:            query.EndHeight,
		EventEncodingVersion: entities.EventEncodingVersion_CCF_V0,
	}

	res, err := c.rpcClient.GetEventsForHeightRange(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	return getEventsResult(res, c.jsonOptions)
}

func (c *BaseClient) GetEventsForBlockIDs(
	ctx context.Context,
	eventType string,
	blockIDs []flow.Identifier,
	opts ...grpc.CallOption,
) ([]flow.BlockEvents, error) {
	req := &access.GetEventsForBlockIDsRequest{
		Type:                 eventType,
		BlockIds:             identifiersToMessages(blockIDs),
		EventEncodingVersion: entities.EventEncodingVersion_CCF_V0,
	}

	res, err := c.rpcClient.GetEventsForBlockIDs(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	return getEventsResult(res, c.jsonOptions)
}

func getEventsResult(res *access.EventsResponse, options []json.Option) ([]flow.BlockEvents, error) {
	resultMessages := res.GetResults()

	results := make([]flow.BlockEvents, len(resultMessages))
	for i, result := range resultMessages {
		eventMessages := result.GetEvents()

		events := make([]flow.Event, len(eventMessages))

		for i, m := range eventMessages {
			evt, err := messageToEvent(m, options)
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

func (c *BaseClient) GetLatestProtocolStateSnapshot(ctx context.Context, opts ...grpc.CallOption) ([]byte, error) {
	res, err := c.rpcClient.GetLatestProtocolStateSnapshot(ctx, &access.GetLatestProtocolStateSnapshotRequest{}, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	return res.GetSerializedSnapshot(), nil
}

func (c *BaseClient) GetExecutionResultForBlockID(ctx context.Context, blockID flow.Identifier, opts ...grpc.CallOption) (*flow.ExecutionResult, error) {
	er, err := c.rpcClient.GetExecutionResultForBlockID(ctx, &access.GetExecutionResultForBlockIDRequest{
		BlockId: identifierToMessage(blockID),
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

func (c *BaseClient) GetExecutionDataByBlockID(
	ctx context.Context,
	blockID flow.Identifier,
	opts ...grpc.CallOption,
) (*flow.ExecutionData, error) {

	ed, err := c.executionDataClient.GetExecutionDataByBlockID(ctx, &executiondata.GetExecutionDataByBlockIDRequest{
		BlockId:              identifierToMessage(blockID),
		EventEncodingVersion: entities.EventEncodingVersion_CCF_V0,
	}, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	return messageToBlockExecutionData(ed.GetBlockExecutionData())

}

func (c *BaseClient) SubscribeExecutionDataByBlockID(
	ctx context.Context,
	startBlockID flow.Identifier,
	opts ...grpc.CallOption,
) (<-chan flow.ExecutionDataStreamResponse, <-chan error, error) {
	req := executiondata.SubscribeExecutionDataRequest{
		StartBlockId:         startBlockID[:],
		EventEncodingVersion: entities.EventEncodingVersion_CCF_V0,
	}
	return c.subscribeExecutionData(ctx, &req, opts...)
}

func (c *BaseClient) SubscribeExecutionDataByBlockHeight(
	ctx context.Context,
	startHeight uint64,
	opts ...grpc.CallOption,
) (<-chan flow.ExecutionDataStreamResponse, <-chan error, error) {
	req := executiondata.SubscribeExecutionDataRequest{
		StartBlockHeight:     startHeight,
		EventEncodingVersion: entities.EventEncodingVersion_CCF_V0,
	}
	return c.subscribeExecutionData(ctx, &req, opts...)
}

func (c *BaseClient) subscribeExecutionData(
	ctx context.Context,
	req *executiondata.SubscribeExecutionDataRequest,
	opts ...grpc.CallOption,
) (<-chan flow.ExecutionDataStreamResponse, <-chan error, error) {
	stream, err := c.executionDataClient.SubscribeExecutionData(ctx, req, opts...)
	if err != nil {
		return nil, nil, err
	}

	sub := make(chan flow.ExecutionDataStreamResponse)
	errChan := make(chan error)

	sendErr := func(err error) {
		select {
		case <-ctx.Done():
		case errChan <- err:
		}
	}

	go func() {
		defer close(sub)
		defer close(errChan)

		for {
			resp, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					return
				}

				sendErr(fmt.Errorf("error receiving execution data: %w", err))
				return
			}

			execData, err := messageToBlockExecutionData(resp.GetBlockExecutionData())
			if err != nil {
				sendErr(fmt.Errorf("error converting execution data for block %d: %w", resp.GetBlockHeight(), err))
				return
			}

			response := flow.ExecutionDataStreamResponse{
				Height:         resp.BlockHeight,
				ExecutionData:  execData,
				BlockTimestamp: resp.BlockTimestamp.AsTime(),
			}

			select {
			case <-ctx.Done():
				return
			case sub <- response:
			}
		}
	}()

	return sub, errChan, nil
}

func (c *BaseClient) SubscribeEventsByBlockID(
	ctx context.Context,
	startBlockID flow.Identifier,
	filter flow.EventFilter,
	opts ...SubscribeOption,
) (<-chan flow.BlockEvents, <-chan error, error) {
	req := executiondata.SubscribeEventsRequest{
		StartBlockId:         startBlockID[:],
		EventEncodingVersion: entities.EventEncodingVersion_CCF_V0,
	}
	return c.subscribeEvents(ctx, &req, filter, opts...)
}

func (c *BaseClient) SubscribeEventsByBlockHeight(
	ctx context.Context,
	startHeight uint64,
	filter flow.EventFilter,
	opts ...SubscribeOption,
) (<-chan flow.BlockEvents, <-chan error, error) {
	req := executiondata.SubscribeEventsRequest{
		StartBlockHeight:     startHeight,
		EventEncodingVersion: entities.EventEncodingVersion_CCF_V0,
	}
	return c.subscribeEvents(ctx, &req, filter, opts...)
}

func (c *BaseClient) subscribeEvents(
	ctx context.Context,
	req *executiondata.SubscribeEventsRequest,
	filter flow.EventFilter,
	opts ...SubscribeOption,
) (<-chan flow.BlockEvents, <-chan error, error) {
	conf := DefaultSubscribeConfig()
	for _, apply := range opts {
		apply(conf)
	}

	req.Filter = &executiondata.EventFilter{
		EventType: filter.EventTypes,
		Address:   filter.Addresses,
		Contract:  filter.Contracts,
	}
	req.HeartbeatInterval = conf.heartbeatInterval

	stream, err := c.executionDataClient.SubscribeEvents(ctx, req, conf.grpcOpts...)
	if err != nil {
		return nil, nil, err
	}

	sub := make(chan flow.BlockEvents)
	errChan := make(chan error)

	sendErr := func(err error) {
		select {
		case <-ctx.Done():
		case errChan <- err:
		}
	}

	go func() {
		defer close(sub)
		defer close(errChan)

		for {
			resp, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					return
				}

				sendErr(fmt.Errorf("error receiving event: %w", err))
				return
			}

			events, err := messagesToEvents(resp.GetEvents(), c.jsonOptions)
			if err != nil {
				sendErr(fmt.Errorf("error converting event for block %d: %w", resp.GetBlockHeight(), err))
				return
			}

			response := flow.BlockEvents{
				Height:         resp.GetBlockHeight(),
				BlockID:        messageToIdentifier(resp.GetBlockId()),
				Events:         events,
				BlockTimestamp: resp.GetBlockTimestamp().AsTime(),
			}

			select {
			case <-ctx.Done():
				return
			case sub <- response:
			}
		}
	}()

	return sub, errChan, nil
}

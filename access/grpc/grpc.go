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

//go:generate go run github.com/vektra/mockery/cmd/mockery --name RPCClient --structname MockRPCClient --output mocks
//go:generate go run github.com/vektra/mockery/cmd/mockery --name ExecutionDataRPCClient --structname MockExecutionDataRPCClient --output mocks

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/onflow/flow/protobuf/go/flow/entities"
	"google.golang.org/grpc"

	"github.com/onflow/cadence"
	"github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow/protobuf/go/flow/access"
	"github.com/onflow/flow/protobuf/go/flow/executiondata"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access/grpc/convert"
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
	eventEncoding       flow.EventEncodingVersion
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
		eventEncoding:       flow.EventEncodingVersionCCF,
	}, nil
}

// NewFromRPCClient initializes a Flow client using a pre-configured gRPC provider.
func NewFromRPCClient(rpcClient RPCClient) *BaseClient {
	return &BaseClient{
		rpcClient:     rpcClient,
		close:         func() error { return nil },
		eventEncoding: flow.EventEncodingVersionCCF,
	}
}

// NewFromExecutionDataRPCClient initializes a Flow client using a pre-configured gRPC provider.
func NewFromExecutionDataRPCClient(rpcClient ExecutionDataRPCClient) *BaseClient {
	return &BaseClient{
		executionDataClient: rpcClient,
		close:               func() error { return nil },
		eventEncoding:       flow.EventEncodingVersionCCF,
	}
}

func (c *BaseClient) SetJSONOptions(options []json.Option) {
	c.jsonOptions = options
}

func (c *BaseClient) SetEventEncoding(version flow.EventEncodingVersion) {
	c.eventEncoding = version
}

func (c *BaseClient) RPCClient() RPCClient {
	return c.rpcClient
}

func (c *BaseClient) ExecutionDataRPCClient() ExecutionDataRPCClient {
	return c.executionDataClient
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

func (c *BaseClient) GetNodeVersionInfo(ctx context.Context, opts ...grpc.CallOption) (*flow.NodeVersionInfo, error) {
	res, err := c.rpcClient.GetNodeVersionInfo(ctx, &access.GetNodeVersionInfoRequest{}, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	info := res.GetInfo()
	return &flow.NodeVersionInfo{
		Semver:               info.Semver,
		Commit:               info.Commit,
		SporkId:              flow.BytesToID(info.SporkId),
		ProtocolVersion:      info.ProtocolVersion,
		SporkRootBlockHeight: info.SporkRootBlockHeight,
		NodeRootBlockHeight:  info.NodeRootBlockHeight,
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
	header, err := convert.MessageToBlockHeader(res.GetBlock())
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
	block, err := convert.MessageToBlock(res.GetBlock())
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

	result, err := convert.MessageToCollection(res.GetCollection())
	if err != nil {
		return nil, newMessageToEntityError(entityCollection, err)
	}

	return &result, nil
}

func (c *BaseClient) GetLightCollectionByID(
	ctx context.Context,
	id flow.Identifier,
	opts ...grpc.CallOption,
) (*flow.Collection, error) {
	req := &access.GetCollectionByIDRequest{
		Id: id.Bytes(),
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

func (c *BaseClient) GetFullCollectionByID(
	ctx context.Context,
	id flow.Identifier,
	opts ...grpc.CallOption,
) (*flow.FullCollection, error) {
	req := &access.GetFullCollectionByIDRequest{
		Id: id.Bytes(),
	}

	res, err := c.rpcClient.GetFullCollectionByID(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	result, err := convert.MessageToFullCollection(res.GetTransactions())
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

	result, err := convert.MessageToTransaction(res.GetTransaction())
	if err != nil {
		return nil, newMessageToEntityError(entityTransaction, err)
	}

	return &result, nil
}

func (c *BaseClient) GetSystemTransaction(
	ctx context.Context,
	blockID flow.Identifier,
	opts ...grpc.CallOption,
) (*flow.Transaction, error) {
	req := &access.GetSystemTransactionRequest{
		BlockId: blockID.Bytes(),
	}

	res, err := c.rpcClient.GetSystemTransaction(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	result, err := convert.MessageToTransaction(res.GetTransaction())
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
		parsed, err := convert.MessageToTransaction(result)
		if err != nil {
			return nil, newMessageToEntityError(entityTransaction, err)
		}
		results = append(results, &parsed)
	}

	return results, nil
}

func (c *BaseClient) GetSystemTransactionResult(
	ctx context.Context,
	blockID flow.Identifier,
	opts ...grpc.CallOption,
) (*flow.TransactionResult, error) {
	req := &access.GetSystemTransactionResultRequest{
		BlockId:              blockID.Bytes(),
		EventEncodingVersion: c.eventEncoding,
	}

	res, err := c.rpcClient.GetSystemTransactionResult(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	result, err := convert.MessageToTransactionResult(res, c.jsonOptions)
	if err != nil {
		return nil, newMessageToEntityError(entityTransactionResult, err)
	}

	return &result, nil
}

func (c *BaseClient) GetTransactionResult(
	ctx context.Context,
	txID flow.Identifier,
	opts ...grpc.CallOption,
) (*flow.TransactionResult, error) {
	req := &access.GetTransactionRequest{
		Id:                   txID.Bytes(),
		EventEncodingVersion: c.eventEncoding,
	}

	res, err := c.rpcClient.GetTransactionResult(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	result, err := convert.MessageToTransactionResult(res, c.jsonOptions)
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
		EventEncodingVersion: c.eventEncoding,
	}

	res, err := c.rpcClient.GetTransactionResultByIndex(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	parsed, err := convert.MessageToTransactionResult(res, c.jsonOptions)
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
		EventEncodingVersion: c.eventEncoding,
	}

	res, err := c.rpcClient.GetTransactionResultsByBlockID(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	unparsedResults := res.GetTransactionResults()
	results := make([]*flow.TransactionResult, 0, len(unparsedResults))
	for _, result := range unparsedResults {
		parsed, err := convert.MessageToTransactionResult(result, c.jsonOptions)
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

	account, err := convert.MessageToAccount(res.GetAccount())
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

	account, err := convert.MessageToAccount(res.GetAccount())
	if err != nil {
		return nil, newMessageToEntityError(entityAccount, err)
	}

	return &account, nil
}

func (c *BaseClient) GetAccountBalanceAtLatestBlock(
	ctx context.Context,
	address flow.Address,
	opts ...grpc.CallOption,
) (uint64, error) {
	request := &access.GetAccountBalanceAtLatestBlockRequest{
		Address: address.Bytes(),
	}

	response, err := c.rpcClient.GetAccountBalanceAtLatestBlock(ctx, request, opts...)
	if err != nil {
		return 0, newRPCError(err)
	}

	return response.GetBalance(), nil
}

func (c *BaseClient) GetAccountBalanceAtBlockHeight(
	ctx context.Context,
	address flow.Address,
	blockHeight uint64,
	opts ...grpc.CallOption,
) (uint64, error) {
	request := &access.GetAccountBalanceAtBlockHeightRequest{
		Address:     address.Bytes(),
		BlockHeight: blockHeight,
	}

	response, err := c.rpcClient.GetAccountBalanceAtBlockHeight(ctx, request, opts...)
	if err != nil {
		return 0, newRPCError(err)
	}

	return response.GetBalance(), nil
}

func (c *BaseClient) GetAccountKeyAtLatestBlock(
	ctx context.Context,
	address flow.Address,
	keyIndex uint32,
) (*flow.AccountKey, error) {
	request := &access.GetAccountKeyAtLatestBlockRequest{
		Address: address.Bytes(),
		Index:   keyIndex,
	}

	response, err := c.rpcClient.GetAccountKeyAtLatestBlock(ctx, request)
	if err != nil {
		return nil, newRPCError(err)
	}

	accountKey, err := convert.MessageToAccountKey(response.GetAccountKey())
	if err != nil {
		return nil, newMessageToEntityError(entityAccount, err)
	}

	return accountKey, nil
}

func (c *BaseClient) GetAccountKeyAtBlockHeight(
	ctx context.Context,
	address flow.Address,
	keyIndex uint32,
	height uint64,
) (*flow.AccountKey, error) {
	request := &access.GetAccountKeyAtBlockHeightRequest{
		Address:     address.Bytes(),
		Index:       keyIndex,
		BlockHeight: height,
	}

	response, err := c.rpcClient.GetAccountKeyAtBlockHeight(ctx, request)
	if err != nil {
		return nil, newRPCError(err)
	}

	accountKey, err := convert.MessageToAccountKey(response.GetAccountKey())
	if err != nil {
		return nil, newMessageToEntityError(entityAccount, err)
	}

	return accountKey, nil
}

func (c *BaseClient) GetAccountKeysAtLatestBlock(
	ctx context.Context,
	address flow.Address,
) ([]flow.AccountKey, error) {
	request := &access.GetAccountKeysAtLatestBlockRequest{
		Address: address.Bytes(),
	}

	response, err := c.rpcClient.GetAccountKeysAtLatestBlock(ctx, request)
	if err != nil {
		return nil, newRPCError(err)
	}

	accountKeys, err := convert.MessageToAccountKeys(response.GetAccountKeys())
	if err != nil {
		return nil, newMessageToEntityError(entityAccount, err)
	}

	return accountKeys, nil
}

func (c *BaseClient) GetAccountKeysAtBlockHeight(
	ctx context.Context,
	address flow.Address,
	height uint64,
) ([]flow.AccountKey, error) {
	request := &access.GetAccountKeysAtBlockHeightRequest{
		Address:     address.Bytes(),
		BlockHeight: height,
	}

	response, err := c.rpcClient.GetAccountKeysAtBlockHeight(ctx, request)
	if err != nil {
		return nil, newRPCError(err)
	}

	accountKeys, err := convert.MessageToAccountKeys(response.GetAccountKeys())
	if err != nil {
		return nil, newMessageToEntityError(entityAccount, err)
	}

	return accountKeys, nil
}

func (c *BaseClient) ExecuteScriptAtLatestBlock(
	ctx context.Context,
	script []byte,
	arguments []cadence.Value,
	opts ...grpc.CallOption,
) (cadence.Value, error) {

	args, err := convert.CadenceValuesToMessages(arguments, flow.EventEncodingVersionJSONCDC)
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

	args, err := convert.CadenceValuesToMessages(arguments, flow.EventEncodingVersionJSONCDC)
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

	args, err := convert.CadenceValuesToMessages(arguments, flow.EventEncodingVersionJSONCDC)
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
	value, err := convert.MessageToCadenceValue(res.GetValue(), options)
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
		EventEncodingVersion: c.eventEncoding,
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
		BlockIds:             convert.IdentifiersToMessages(blockIDs),
		EventEncodingVersion: c.eventEncoding,
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
			evt, err := convert.MessageToEvent(m, options)
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

func (c *BaseClient) GetProtocolStateSnapshotByBlockID(ctx context.Context, blockID flow.Identifier, opts ...grpc.CallOption) ([]byte, error) {
	req := &access.GetProtocolStateSnapshotByBlockIDRequest{
		BlockId: blockID.Bytes(),
	}

	res, err := c.rpcClient.GetProtocolStateSnapshotByBlockID(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	return res.GetSerializedSnapshot(), nil
}

func (c *BaseClient) GetProtocolStateSnapshotByHeight(ctx context.Context, blockHeight uint64, opts ...grpc.CallOption) ([]byte, error) {
	req := &access.GetProtocolStateSnapshotByHeightRequest{
		BlockHeight: blockHeight,
	}

	res, err := c.rpcClient.GetProtocolStateSnapshotByHeight(ctx, req, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	return res.GetSerializedSnapshot(), nil
}

func (c *BaseClient) GetExecutionResultForBlockID(ctx context.Context, blockID flow.Identifier, opts ...grpc.CallOption) (*flow.ExecutionResult, error) {
	er, err := c.rpcClient.GetExecutionResultForBlockID(ctx, &access.GetExecutionResultForBlockIDRequest{
		BlockId: convert.IdentifierToMessage(blockID),
	}, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	return convert.MessageToExecutionResult(er.ExecutionResult)
}

func (c *BaseClient) GetExecutionResultByID(ctx context.Context, id flow.Identifier, opts ...grpc.CallOption) (*flow.ExecutionResult, error) {
	er, err := c.rpcClient.GetExecutionResultByID(ctx, &access.GetExecutionResultByIDRequest{
		Id: convert.IdentifierToMessage(id),
	}, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	return convert.MessageToExecutionResult(er.ExecutionResult)
}

func (c *BaseClient) GetExecutionDataByBlockID(
	ctx context.Context,
	blockID flow.Identifier,
	opts ...grpc.CallOption,
) (*flow.ExecutionData, error) {

	ed, err := c.executionDataClient.GetExecutionDataByBlockID(ctx, &executiondata.GetExecutionDataByBlockIDRequest{
		BlockId:              convert.IdentifierToMessage(blockID),
		EventEncodingVersion: c.eventEncoding,
	}, opts...)
	if err != nil {
		return nil, newRPCError(err)
	}

	return convert.MessageToBlockExecutionData(ed.GetBlockExecutionData())

}

func (c *BaseClient) SubscribeExecutionDataByBlockID(
	ctx context.Context,
	startBlockID flow.Identifier,
	opts ...grpc.CallOption,
) (<-chan flow.ExecutionDataStreamResponse, <-chan error, error) {
	req := executiondata.SubscribeExecutionDataRequest{
		StartBlockId:         startBlockID[:],
		EventEncodingVersion: c.eventEncoding,
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
		EventEncodingVersion: c.eventEncoding,
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

			execData, err := convert.MessageToBlockExecutionData(resp.GetBlockExecutionData())
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
		EventEncodingVersion: c.eventEncoding,
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
		EventEncodingVersion: c.eventEncoding,
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

			events, err := convert.MessagesToEvents(resp.GetEvents(), c.jsonOptions)
			if err != nil {
				sendErr(fmt.Errorf("error converting event for block %d: %w", resp.GetBlockHeight(), err))
				return
			}

			response := flow.BlockEvents{
				Height:         resp.GetBlockHeight(),
				BlockID:        convert.MessageToIdentifier(resp.GetBlockId()),
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

func (c *BaseClient) SubscribeBlocksFromStartBlockID(
	ctx context.Context,
	startBlockID flow.Identifier,
	blockStatus flow.BlockStatus,
	opts ...grpc.CallOption,
) (<-chan flow.Block, <-chan error, error) {
	status := convert.BlockStatusToEntity(blockStatus)
	if status == entities.BlockStatus_BLOCK_UNKNOWN {
		return nil, nil, newRPCError(errors.New("unknown block status"))
	}

	request := &access.SubscribeBlocksFromStartBlockIDRequest{
		StartBlockId: startBlockID.Bytes(),
		BlockStatus:  status,
	}

	subscribeClient, err := c.rpcClient.SubscribeBlocksFromStartBlockID(ctx, request, opts...)
	if err != nil {
		return nil, nil, newRPCError(err)
	}

	blocksChan := make(chan flow.Block)
	errChan := make(chan error)

	go func() {
		defer close(blocksChan)
		defer close(errChan)
		receiveBlocksFromClient(ctx, subscribeClient, blocksChan, errChan)
	}()

	return blocksChan, errChan, nil
}

func (c *BaseClient) SubscribeBlocksFromStartHeight(
	ctx context.Context,
	startHeight uint64,
	blockStatus flow.BlockStatus,
	opts ...grpc.CallOption,
) (<-chan flow.Block, <-chan error, error) {
	status := convert.BlockStatusToEntity(blockStatus)
	if status == entities.BlockStatus_BLOCK_UNKNOWN {
		return nil, nil, newRPCError(errors.New("unknown block status"))
	}

	request := &access.SubscribeBlocksFromStartHeightRequest{
		StartBlockHeight: startHeight,
		BlockStatus:      status,
	}

	subscribeClient, err := c.rpcClient.SubscribeBlocksFromStartHeight(ctx, request, opts...)
	if err != nil {
		return nil, nil, newRPCError(err)
	}

	blocksChan := make(chan flow.Block)
	errChan := make(chan error)

	go func() {
		defer close(blocksChan)
		defer close(errChan)
		receiveBlocksFromClient(ctx, subscribeClient, blocksChan, errChan)
	}()

	return blocksChan, errChan, nil
}

func (c *BaseClient) SubscribeBlocksFromLatest(
	ctx context.Context,
	blockStatus flow.BlockStatus,
	opts ...grpc.CallOption,
) (<-chan flow.Block, <-chan error, error) {
	status := convert.BlockStatusToEntity(blockStatus)
	if status == entities.BlockStatus_BLOCK_UNKNOWN {
		return nil, nil, newRPCError(errors.New("unknown block status"))
	}

	request := &access.SubscribeBlocksFromLatestRequest{
		BlockStatus: status,
	}

	subscribeClient, err := c.rpcClient.SubscribeBlocksFromLatest(ctx, request, opts...)
	if err != nil {
		return nil, nil, newRPCError(err)
	}

	blocksChan := make(chan flow.Block)
	errChan := make(chan error)

	go func() {
		defer close(blocksChan)
		defer close(errChan)
		receiveBlocksFromClient(ctx, subscribeClient, blocksChan, errChan)
	}()

	return blocksChan, errChan, nil
}

func receiveBlocksFromClient[Client interface {
	Recv() (*access.SubscribeBlocksResponse, error)
}](
	ctx context.Context,
	client Client,
	blocksChan chan<- flow.Block,
	errChan chan<- error,
) {
	sendErr := func(err error) {
		select {
		case <-ctx.Done():
		case errChan <- err:
		}
	}

	for {
		// Receive the next block response
		blockResponse, err := client.Recv()
		if err != nil {
			if err == io.EOF {
				// End of stream, return gracefully
				return
			}

			sendErr(fmt.Errorf("error receiving block: %w", err))
			return
		}

		block, err := convert.MessageToBlock(blockResponse.GetBlock())
		if err != nil {
			sendErr(fmt.Errorf("error converting message to block: %w", err))
			return
		}

		select {
		case <-ctx.Done():
			return
		case blocksChan <- block:
		}
	}
}

func (c *BaseClient) SubscribeAccountStatusesFromStartHeight(
	ctx context.Context,
	startHeight uint64,
	filter flow.AccountStatusFilter,
	opts ...grpc.CallOption,
) (<-chan flow.AccountStatus, <-chan error, error) {
	request := &executiondata.SubscribeAccountStatusesFromStartHeightRequest{
		StartBlockHeight:     startHeight,
		EventEncodingVersion: c.eventEncoding,
	}
	request.Filter = &executiondata.StatusFilter{
		EventType: filter.EventTypes,
		Address:   filter.Addresses,
	}

	subscribeClient, err := c.executionDataClient.SubscribeAccountStatusesFromStartHeight(ctx, request, opts...)
	if err != nil {
		return nil, nil, newRPCError(err)
	}

	accountStatutesChan := make(chan flow.AccountStatus)
	errChan := make(chan error)

	go func() {
		defer close(accountStatutesChan)
		defer close(errChan)
		receiveAccountStatusesFromStream(ctx, subscribeClient, accountStatutesChan, errChan)
	}()

	return accountStatutesChan, errChan, nil
}

func (c *BaseClient) SubscribeAccountStatusesFromStartBlockID(
	ctx context.Context,
	startBlockID flow.Identifier,
	filter flow.AccountStatusFilter,
	opts ...grpc.CallOption,
) (<-chan flow.AccountStatus, <-chan error, error) {
	request := &executiondata.SubscribeAccountStatusesFromStartBlockIDRequest{
		StartBlockId:         startBlockID.Bytes(),
		EventEncodingVersion: c.eventEncoding,
	}
	request.Filter = &executiondata.StatusFilter{
		EventType: filter.EventTypes,
		Address:   filter.Addresses,
	}

	subscribeClient, err := c.executionDataClient.SubscribeAccountStatusesFromStartBlockID(ctx, request, opts...)
	if err != nil {
		return nil, nil, newRPCError(err)
	}

	accountStatutesChan := make(chan flow.AccountStatus)
	errChan := make(chan error)

	go func() {
		defer close(accountStatutesChan)
		defer close(errChan)
		receiveAccountStatusesFromStream(ctx, subscribeClient, accountStatutesChan, errChan)
	}()

	return accountStatutesChan, errChan, nil
}

func (c *BaseClient) SubscribeAccountStatusesFromLatestBlock(
	ctx context.Context,
	filter flow.AccountStatusFilter,
	opts ...grpc.CallOption,
) (<-chan flow.AccountStatus, <-chan error, error) {
	request := &executiondata.SubscribeAccountStatusesFromLatestBlockRequest{
		EventEncodingVersion: c.eventEncoding,
	}
	request.Filter = &executiondata.StatusFilter{
		EventType: filter.EventTypes,
		Address:   filter.Addresses,
	}

	subscribeClient, err := c.executionDataClient.SubscribeAccountStatusesFromLatestBlock(ctx, request, opts...)
	if err != nil {
		return nil, nil, newRPCError(err)
	}

	accountStatutesChan := make(chan flow.AccountStatus)
	errChan := make(chan error)

	go func() {
		defer close(accountStatutesChan)
		defer close(errChan)
		receiveAccountStatusesFromStream(ctx, subscribeClient, accountStatutesChan, errChan)
	}()

	return accountStatutesChan, errChan, nil
}

func receiveAccountStatusesFromStream[Stream interface {
	Recv() (*executiondata.SubscribeAccountStatusesResponse, error)
}](
	ctx context.Context,
	stream Stream,
	accountStatutesChan chan<- flow.AccountStatus,
	errChan chan<- error,
) {
	sendErr := func(err error) {
		select {
		case <-ctx.Done():
		case errChan <- err:
		}
	}

	var nextExpectedMsgIndex uint64
	for {
		accountStatusResponse, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// End of stream, return gracefully
				return
			}

			sendErr(fmt.Errorf("error receiving account status: %w", err))
			return
		}

		accountStatus, err := convert.MessageToAccountStatus(accountStatusResponse)
		if err != nil {
			sendErr(fmt.Errorf("error converting message to account status: %w", err))
			return
		}

		if accountStatus.MessageIndex != nextExpectedMsgIndex {
			sendErr(fmt.Errorf("messages are not ordered"))
			return
		}
		nextExpectedMsgIndex = accountStatus.MessageIndex + 1

		select {
		case <-ctx.Done():
			return
		case accountStatutesChan <- accountStatus:
		}
	}
}

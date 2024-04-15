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

// Package grpc client provides a Go client for the Flow Access gRPC API.
//
// The Access API provides a set of methods that can be used to submit transactions
// and read state from Flow. This client is compatible with the Access API implemented by the
// Access Node role, as well as the mock Access API exposed by the Flow Emulator.
//
// The full Access API specification is here: https://docs.onflow.org/access-api/
package grpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/onflow/cadence"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/onflow/flow-go-sdk"
)

const EmulatorHost = "127.0.0.1:3569"
const TestnetHost = "access.devnet.nodes.onflow.org:9000"
const CanarynetHost = "access.canary.nodes.onflow.org:9000"
const MainnetHost = "access.mainnet.nodes.onflow.org:9000"
const PreviewnetHost = "access.previewnet.nodes.onflow.org:9000"

// NewClient creates an gRPC client exposing all the common access APIs.
// Client will use provided host for connection.
func NewClient(host string, opts ...grpc.DialOption) (*Client, error) {
	var client *BaseClient
	var err error
	if len(opts) > 0 {
		client, err = NewBaseClient(host, opts...)
	} else {
		client, err = NewBaseClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	if err != nil {
		return nil, err
	}

	return &Client{client}, nil
}

// Client implements all common gRPC methods providing a network agnostic API.
type Client struct {
	grpc *BaseClient
}

func (c *Client) Ping(ctx context.Context) error {
	return c.grpc.Ping(ctx)
}

func (c *Client) WaitServer(ctx context.Context) error {
	return c.grpc.Ping(ctx, grpc.WaitForReady(true))
}

func (c *Client) GetNetworkParameters(ctx context.Context) (*flow.NetworkParameters, error) {
	return c.grpc.GetNetworkParameters(ctx)
}

func (c *Client) GetLatestBlockHeader(ctx context.Context, isSealed bool) (*flow.BlockHeader, error) {
	return c.grpc.GetLatestBlockHeader(ctx, isSealed)
}

func (c *Client) GetBlockHeaderByID(ctx context.Context, blockID flow.Identifier) (*flow.BlockHeader, error) {
	return c.grpc.GetBlockHeaderByID(ctx, blockID)
}

func (c *Client) GetBlockHeaderByHeight(ctx context.Context, height uint64) (*flow.BlockHeader, error) {
	return c.grpc.GetBlockHeaderByHeight(ctx, height)
}

func (c *Client) GetLatestBlock(ctx context.Context, isSealed bool) (*flow.Block, error) {
	return c.grpc.GetLatestBlock(ctx, isSealed)
}

func (c *Client) GetBlockByID(ctx context.Context, blockID flow.Identifier) (*flow.Block, error) {
	return c.grpc.GetBlockByID(ctx, blockID)
}

func (c *Client) GetBlockByHeight(ctx context.Context, height uint64) (*flow.Block, error) {
	return c.grpc.GetBlockByHeight(ctx, height)
}

func (c *Client) GetCollection(ctx context.Context, colID flow.Identifier) (*flow.Collection, error) {
	return c.grpc.GetCollection(ctx, colID)
}

func (c *Client) SendTransaction(ctx context.Context, tx flow.Transaction) error {
	return c.grpc.SendTransaction(ctx, tx)
}

func (c *Client) GetTransaction(ctx context.Context, txID flow.Identifier) (*flow.Transaction, error) {
	return c.grpc.GetTransaction(ctx, txID)
}

func (c *Client) GetTransactionsByBlockID(ctx context.Context, blockID flow.Identifier) ([]*flow.Transaction, error) {
	return c.grpc.GetTransactionsByBlockID(ctx, blockID)
}

func (c *Client) GetTransactionResult(ctx context.Context, txID flow.Identifier) (*flow.TransactionResult, error) {
	return c.grpc.GetTransactionResult(ctx, txID)
}

func (c *Client) GetTransactionResultByIndex(ctx context.Context, blockID flow.Identifier, index uint32) (*flow.TransactionResult, error) {
	return c.grpc.GetTransactionResultByIndex(ctx, blockID, index)
}
func (c *Client) GetTransactionResultsByBlockID(ctx context.Context, blockID flow.Identifier) ([]*flow.TransactionResult, error) {
	return c.grpc.GetTransactionResultsByBlockID(ctx, blockID)
}

func (c *Client) GetAccount(ctx context.Context, address flow.Address) (*flow.Account, error) {
	return c.grpc.GetAccount(ctx, address)
}

func (c *Client) GetAccountAtLatestBlock(ctx context.Context, address flow.Address) (*flow.Account, error) {
	return c.grpc.GetAccountAtLatestBlock(ctx, address)
}

func (c *Client) GetAccountAtBlockHeight(ctx context.Context, address flow.Address, blockHeight uint64) (*flow.Account, error) {
	return c.grpc.GetAccountAtBlockHeight(ctx, address, blockHeight)
}

func (c *Client) ExecuteScriptAtLatestBlock(ctx context.Context, script []byte, arguments []cadence.Value) (cadence.Value, error) {
	return c.grpc.ExecuteScriptAtLatestBlock(ctx, script, arguments)
}

func (c *Client) ExecuteScriptAtBlockID(ctx context.Context, blockID flow.Identifier, script []byte, arguments []cadence.Value) (cadence.Value, error) {
	return c.grpc.ExecuteScriptAtBlockID(ctx, blockID, script, arguments)
}

func (c *Client) ExecuteScriptAtBlockHeight(ctx context.Context, height uint64, script []byte, arguments []cadence.Value) (cadence.Value, error) {
	return c.grpc.ExecuteScriptAtBlockHeight(ctx, height, script, arguments)
}

func (c *Client) GetEventsForHeightRange(ctx context.Context, eventType string, startHeight uint64, endHeight uint64) ([]flow.BlockEvents, error) {
	return c.grpc.GetEventsForHeightRange(ctx, EventRangeQuery{
		Type:        eventType,
		StartHeight: startHeight,
		EndHeight:   endHeight,
	})
}

func (c *Client) GetEventsForBlockIDs(ctx context.Context, eventType string, blockIDs []flow.Identifier) ([]flow.BlockEvents, error) {
	return c.grpc.GetEventsForBlockIDs(ctx, eventType, blockIDs)
}

func (c *Client) GetLatestProtocolStateSnapshot(ctx context.Context) ([]byte, error) {
	return c.grpc.GetLatestProtocolStateSnapshot(ctx)
}

func (c *Client) GetExecutionResultForBlockID(ctx context.Context, blockID flow.Identifier) (*flow.ExecutionResult, error) {
	return c.grpc.GetExecutionResultForBlockID(ctx, blockID)
}

func (c *Client) GetExecutionDataByBlockID(ctx context.Context, blockID flow.Identifier) (*flow.ExecutionData, error) {
	return c.grpc.GetExecutionDataByBlockID(ctx, blockID)
}

func (c *Client) SubscribeExecutionDataByBlockID(ctx context.Context, startBlockID flow.Identifier) (<-chan flow.ExecutionDataStreamResponse, <-chan error, error) {
	return c.grpc.SubscribeExecutionDataByBlockID(ctx, startBlockID)
}

func (c *Client) SubscribeExecutionDataByBlockHeight(ctx context.Context, startHeight uint64) (<-chan flow.ExecutionDataStreamResponse, <-chan error, error) {
	return c.grpc.SubscribeExecutionDataByBlockHeight(ctx, startHeight)
}

func (c *Client) SubscribeEventsByBlockID(ctx context.Context, startBlockID flow.Identifier, filter flow.EventFilter) (<-chan flow.BlockEvents, <-chan error, error) {
	return c.grpc.SubscribeEventsByBlockID(ctx, startBlockID, filter)
}

func (c *Client) SubscribeEventsByBlockHeight(ctx context.Context, startHeight uint64, filter flow.EventFilter) (<-chan flow.BlockEvents, <-chan error, error) {
	return c.grpc.SubscribeEventsByBlockHeight(ctx, startHeight, filter)
}

func (c *Client) Close() error {
	return c.grpc.Close()
}

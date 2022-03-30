/*
 * Flow Go SDK
 *
 * Copyright 2019-2020 Dapper Labs, Inc.
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

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
)

const EMULATOR_API = "127.0.0.1:3569"
const TESTNET_API = "access.devnet.nodes.onflow.org:9000"
const CANARYNET_API = "access.canary.nodes.onflow.org:9000"
const MAINNET_API = "access.mainnet.nodes.onflow.org:9000"

// NewClient create a client by passing the gRPC handler.
func NewClient(handler *GRPCClient) *BaseClient {
	return &BaseClient{
		grpc: handler,
	}
}

// BaseClient complies with the client interface and hides any gRPC specific options.
type BaseClient struct {
	grpc *GRPCClient
}

// NewDefaultEmulatorClient creates a client for accessing default local emulator network using gRPC.
func NewDefaultEmulatorClient() (*BaseClient, error) {
	handler, err := New(EMULATOR_API)
	if err != nil {
		return nil, err
	}

	return NewClient(handler), nil
}

// NewDefaultTestnetClient creates a client for accessing default testnet AN using gRPC.
func NewDefaultTestnetClient() (*BaseClient, error) {
	handler, err := New(TESTNET_API)
	if err != nil {
		return nil, err
	}

	return NewClient(handler), nil
}

// NewDefaultCanaryClient creates a client for accessing default canary AN using gRPC.
func NewDefaultCanaryClient() (*BaseClient, error) {
	handler, err := New(CANARYNET_API)
	if err != nil {
		return nil, err
	}

	return NewClient(handler), nil
}

// NewDefaultMainnetClient creates a client for accessing default mainnet AN using gRPC.
func NewDefaultMainnetClient() (*BaseClient, error) {
	handler, err := New(MAINNET_API)
	if err != nil {
		return nil, err
	}

	return NewClient(handler), nil
}

func (c *BaseClient) Ping(ctx context.Context) error {
	return c.grpc.Ping(ctx)
}

func (c *BaseClient) GetLatestBlockHeader(ctx context.Context, isSealed bool) (*flow.BlockHeader, error) {
	return c.grpc.GetLatestBlockHeader(ctx, isSealed)
}

func (c *BaseClient) GetBlockHeaderByID(ctx context.Context, blockID flow.Identifier) (*flow.BlockHeader, error) {
	return c.grpc.GetBlockHeaderByID(ctx, blockID)
}

func (c *BaseClient) GetBlockHeaderByHeight(ctx context.Context, height uint64) (*flow.BlockHeader, error) {
	return c.grpc.GetBlockHeaderByHeight(ctx, height)
}

func (c *BaseClient) GetLatestBlock(ctx context.Context, isSealed bool) (*flow.Block, error) {
	return c.grpc.GetLatestBlock(ctx, isSealed)
}

func (c *BaseClient) GetBlockByID(ctx context.Context, blockID flow.Identifier) (*flow.Block, error) {
	return c.grpc.GetBlockByID(ctx, blockID)
}

func (c *BaseClient) GetBlockByHeight(ctx context.Context, height uint64) (*flow.Block, error) {
	return c.grpc.GetBlockByHeight(ctx, height)
}

func (c *BaseClient) GetCollection(ctx context.Context, colID flow.Identifier) (*flow.Collection, error) {
	return c.grpc.GetCollection(ctx, colID)
}

func (c *BaseClient) SendTransaction(ctx context.Context, tx flow.Transaction) error {
	return c.grpc.SendTransaction(ctx, tx)
}

func (c *BaseClient) GetTransaction(ctx context.Context, txID flow.Identifier) (*flow.Transaction, error) {
	return c.grpc.GetTransaction(ctx, txID)
}

func (c *BaseClient) GetTransactionResult(ctx context.Context, txID flow.Identifier) (*flow.TransactionResult, error) {
	return c.grpc.GetTransactionResult(ctx, txID)
}

func (c *BaseClient) GetAccount(ctx context.Context, address flow.Address) (*flow.Account, error) {
	return c.grpc.GetAccount(ctx, address)
}

func (c *BaseClient) GetAccountAtLatestBlock(ctx context.Context, address flow.Address) (*flow.Account, error) {
	return c.grpc.GetAccountAtLatestBlock(ctx, address)
}

func (c *BaseClient) GetAccountAtBlockHeight(ctx context.Context, address flow.Address, blockHeight uint64) (*flow.Account, error) {
	return c.grpc.GetAccountAtBlockHeight(ctx, address, blockHeight)
}

func (c *BaseClient) ExecuteScriptAtLatestBlock(ctx context.Context, script []byte, arguments []cadence.Value) (cadence.Value, error) {
	return c.grpc.ExecuteScriptAtLatestBlock(ctx, script, arguments)
}

func (c *BaseClient) ExecuteScriptAtBlockID(ctx context.Context, blockID flow.Identifier, script []byte, arguments []cadence.Value) (cadence.Value, error) {
	return c.grpc.ExecuteScriptAtBlockID(ctx, blockID, script, arguments)
}

func (c *BaseClient) ExecuteScriptAtBlockHeight(ctx context.Context, height uint64, script []byte, arguments []cadence.Value) (cadence.Value, error) {
	return c.grpc.ExecuteScriptAtBlockHeight(ctx, height, script, arguments)
}

func (c *BaseClient) GetEventsForHeightRange(ctx context.Context, eventType string, startHeight uint64, endHeight uint64) ([]flow.BlockEvents, error) {
	return c.grpc.GetEventsForHeightRange(ctx, EventRangeQuery{
		Type:        eventType,
		StartHeight: startHeight,
		EndHeight:   endHeight,
	})
}

func (c *BaseClient) GetEventsForBlockIDs(ctx context.Context, eventType string, blockIDs []flow.Identifier) ([]flow.BlockEvents, error) {
	return c.grpc.GetEventsForBlockIDs(ctx, eventType, blockIDs)
}

func (c *BaseClient) GetLatestProtocolStateSnapshot(ctx context.Context) ([]byte, error) {
	return c.grpc.GetLatestProtocolStateSnapshot(ctx)
}

func (c *BaseClient) GetExecutionResultForBlockID(ctx context.Context, blockID flow.Identifier) (*flow.ExecutionResult, error) {
	return c.grpc.GetExecutionResultForBlockID(ctx, blockID)
}

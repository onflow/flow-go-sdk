/*
 * Flow Go SDK
 *
 * Copyright 2019-2022 Dapper Labs, Inc.
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

package http

import (
	"context"

	"github.com/onflow/cadence"

	"github.com/onflow/flow-go-sdk"
)

const (
	EMULATOR_URL  = "http://127.0.0.1:8888/v1"
	TESTNET_URL   = "https://rest-testnet.onflow.org/v1/"
	MAINNET_URL   = "https://rest-mainnet.onflow.org/v1/"
	CANARYNET_URL = "" // todo define
)

// NewClient creates an instance of the client with the provided http handler.
func NewClient(url string) (*BaseClient, error) {
	client, err := NewHTTPClient(url)
	if err != nil {
		return nil, err
	}
	return &BaseClient{client}, nil
}

// NewDefaultEmulatorClient creates a new client for connecting to the emulator AN API.
func NewDefaultEmulatorClient(debug bool) (*BaseClient, error) {
	return NewClient(EMULATOR_URL)
}

// NewDefaultTestnetClient creates a new client for connecting to the testnet AN API.
func NewDefaultTestnetClient() (*BaseClient, error) {
	return NewClient(TESTNET_URL)
}

// NewDefaultCanaryClient creates a new client for connecting to the canary AN API.
func NewDefaultCanaryClient() (*BaseClient, error) {
	return NewClient(CANARYNET_URL)
}

// NewDefaultMainnetClient creates a new client for connecting to the mainnet AN API.
func NewDefaultMainnetClient() (*BaseClient, error) {
	return NewClient(MAINNET_URL)
}

// BaseClient implementing all the network interactions according to the client interface.
type BaseClient struct {
	httpClient *HTTPClient
}

func (c *BaseClient) Ping(ctx context.Context) error {
	panic("implement me")
}

func (c *BaseClient) GetBlockByID(ctx context.Context, blockID flow.Identifier) (*flow.Block, error) {
	return c.httpClient.GetBlockByID(ctx, blockID)
}

func (c *BaseClient) GetLatestBlockHeader(ctx context.Context, isSealed bool) (*flow.BlockHeader, error) {
	block, err := c.GetLatestBlock(ctx, isSealed)
	if err != nil {
		return nil, err
	}

	return &block.BlockHeader, nil
}

func (c *BaseClient) GetBlockHeaderByID(ctx context.Context, blockID flow.Identifier) (*flow.BlockHeader, error) {
	block, err := c.GetBlockByID(ctx, blockID) // todo optimization: passing the 'select' option to only get the header
	if err != nil {
		return nil, err
	}

	return &block.BlockHeader, nil
}

func (c *BaseClient) GetBlockHeaderByHeight(ctx context.Context, height uint64) (*flow.BlockHeader, error) {
	block, err := c.GetBlockByHeight(ctx, height) // todo optimization: passing the 'select' option to only get the header
	if err != nil {
		return nil, err
	}

	return &block.BlockHeader, nil
}

func (c *BaseClient) GetLatestBlock(ctx context.Context, isSealed bool) (*flow.Block, error) {
	height := FINAL
	if isSealed {
		height = SEALED
	}

	blocks, err := c.httpClient.GetBlocksByHeights(
		ctx,
		HeightQuery{Heights: []uint64{height}},
	)
	if err != nil {
		return nil, err
	}

	return blocks[0], nil
}

func (c *BaseClient) GetBlockByHeight(ctx context.Context, height uint64) (*flow.Block, error) {
	blocks, err := c.httpClient.GetBlocksByHeights(ctx, HeightQuery{Heights: []uint64{height}})
	if err != nil {
		return nil, err
	}

	return blocks[0], nil
}

func (c *BaseClient) GetCollection(ctx context.Context, ID flow.Identifier) (*flow.Collection, error) {
	return c.httpClient.GetCollection(ctx, ID)
}

func (c *BaseClient) SendTransaction(ctx context.Context, tx flow.Transaction) error {
	return c.httpClient.SendTransaction(ctx, tx)
}

func (c *BaseClient) GetTransaction(ctx context.Context, ID flow.Identifier) (*flow.Transaction, error) {
	return c.httpClient.GetTransaction(ctx, ID)
}

func (c *BaseClient) GetTransactionResult(ctx context.Context, ID flow.Identifier) (*flow.TransactionResult, error) {
	return c.httpClient.GetTransactionResult(ctx, ID)
}

// GetAccount is an alias for GetAccountAtLatestBlock.
func (c *BaseClient) GetAccount(ctx context.Context, address flow.Address) (*flow.Account, error) {
	return c.GetAccountAtLatestBlock(ctx, address)
}

func (c *BaseClient) GetAccountAtLatestBlock(ctx context.Context, address flow.Address) (*flow.Account, error) {
	return c.httpClient.GetAccountAtBlockHeight(
		ctx,
		address, HeightQuery{Heights: []uint64{SEALED}},
	)
}

func (c *BaseClient) GetAccountAtBlockHeight(
	ctx context.Context,
	address flow.Address,
	blockHeight uint64,
) (*flow.Account, error) {
	return c.httpClient.GetAccountAtBlockHeight(
		ctx,
		address,
		HeightQuery{Heights: []uint64{blockHeight}},
	)
}

func (c *BaseClient) ExecuteScriptAtLatestBlock(
	ctx context.Context,
	script []byte,
	arguments []cadence.Value,
) (cadence.Value, error) {
	return c.httpClient.ExecuteScriptAtBlockHeight(
		ctx,
		HeightQuery{Heights: []uint64{SEALED}},
		script,
		arguments,
	)
}

func (c *BaseClient) ExecuteScriptAtBlockID(
	ctx context.Context,
	blockID flow.Identifier,
	script []byte,
	arguments []cadence.Value,
) (cadence.Value, error) {
	return c.httpClient.ExecuteScriptAtBlockID(ctx, blockID, script, arguments)
}

func (c *BaseClient) ExecuteScriptAtBlockHeight(
	ctx context.Context,
	height uint64,
	script []byte,
	arguments []cadence.Value,
) (cadence.Value, error) {
	return c.httpClient.ExecuteScriptAtBlockHeight(
		ctx,
		HeightQuery{Heights: []uint64{height}},
		script,
		arguments,
	)
}

func (c *BaseClient) GetEventsForHeightRange(
	ctx context.Context,
	eventType string,
	startHeight uint64,
	endHeight uint64,
) ([]flow.BlockEvents, error) {
	return c.httpClient.GetEventsForHeightRange(
		ctx,
		eventType,
		HeightQuery{
			Start: startHeight,
			End:   endHeight,
		},
	)
}

func (c *BaseClient) GetEventsForBlockIDs(
	ctx context.Context,
	eventType string,
	blockIDs []flow.Identifier,
) ([]flow.BlockEvents, error) {
	return c.httpClient.GetEventsForBlockIDs(ctx, eventType, blockIDs)
}

func (c *BaseClient) GetLatestProtocolStateSnapshot(ctx context.Context) ([]byte, error) {
	panic("implement me")
}

func (c *BaseClient) GetExecutionResultForBlockID(ctx context.Context, blockID flow.Identifier) (*flow.ExecutionResult, error) {
	panic("implement me")
}

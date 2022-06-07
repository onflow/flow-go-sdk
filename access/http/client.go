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

package http

import (
	"context"

	"github.com/onflow/cadence"

	"github.com/onflow/flow-go-sdk"
)

const (
	EmulatorHost  = "http://127.0.0.1:8888/v1"
	TestnetHost   = "https://rest-testnet.onflow.org/v1/"
	MainnetHost   = "https://rest-mainnet.onflow.org/v1/"
	CanarynetHost = "https://rest-canary.onflow.org/v1/"
)

// NewClient creates an HTTP client exposing all the common access APIs.
// Client will use provided host for connection.
func NewClient(host string) (*Client, error) {
	client, err := NewBaseClient(host)
	if err != nil {
		return nil, err
	}
	return &Client{client}, nil
}

// Client implements all common HTTP methods providing a network agnostic API.
type Client struct {
	httpClient *BaseClient
}

func (c *Client) Ping(ctx context.Context) error {
	return c.httpClient.Ping(ctx)
}

func (c *Client) GetBlockByID(ctx context.Context, blockID flow.Identifier) (*flow.Block, error) {
	return c.httpClient.GetBlockByID(ctx, blockID)
}

func (c *Client) GetLatestBlockHeader(ctx context.Context, isSealed bool) (*flow.BlockHeader, error) {
	block, err := c.GetLatestBlock(ctx, isSealed)
	if err != nil {
		return nil, err
	}

	return &block.BlockHeader, nil
}

func (c *Client) GetBlockHeaderByID(ctx context.Context, blockID flow.Identifier) (*flow.BlockHeader, error) {
	block, err := c.GetBlockByID(ctx, blockID) // todo optimization: passing the 'select' option to only get the header
	if err != nil {
		return nil, err
	}

	return &block.BlockHeader, nil
}

func (c *Client) GetBlockHeaderByHeight(ctx context.Context, height uint64) (*flow.BlockHeader, error) {
	block, err := c.GetBlockByHeight(ctx, height) // todo optimization: passing the 'select' option to only get the header
	if err != nil {
		return nil, err
	}

	return &block.BlockHeader, nil
}

func (c *Client) GetLatestBlock(ctx context.Context, isSealed bool) (*flow.Block, error) {
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

func (c *Client) GetBlockByHeight(ctx context.Context, height uint64) (*flow.Block, error) {
	blocks, err := c.httpClient.GetBlocksByHeights(ctx, HeightQuery{Heights: []uint64{height}})
	if err != nil {
		return nil, err
	}

	return blocks[0], nil
}

func (c *Client) GetCollection(ctx context.Context, ID flow.Identifier) (*flow.Collection, error) {
	return c.httpClient.GetCollection(ctx, ID)
}

func (c *Client) SendTransaction(ctx context.Context, tx flow.Transaction) error {
	return c.httpClient.SendTransaction(ctx, tx)
}

func (c *Client) GetTransaction(ctx context.Context, ID flow.Identifier) (*flow.Transaction, error) {
	return c.httpClient.GetTransaction(ctx, ID)
}

func (c *Client) GetTransactionResult(ctx context.Context, ID flow.Identifier) (*flow.TransactionResult, error) {
	return c.httpClient.GetTransactionResult(ctx, ID)
}

// GetAccount is an alias for GetAccountAtLatestBlock.
func (c *Client) GetAccount(ctx context.Context, address flow.Address) (*flow.Account, error) {
	return c.GetAccountAtLatestBlock(ctx, address)
}

func (c *Client) GetAccountAtLatestBlock(ctx context.Context, address flow.Address) (*flow.Account, error) {
	return c.httpClient.GetAccountAtBlockHeight(
		ctx,
		address, HeightQuery{Heights: []uint64{SEALED}},
	)
}

func (c *Client) GetAccountAtBlockHeight(
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

func (c *Client) ExecuteScriptAtLatestBlock(
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

func (c *Client) ExecuteScriptAtBlockID(
	ctx context.Context,
	blockID flow.Identifier,
	script []byte,
	arguments []cadence.Value,
) (cadence.Value, error) {
	return c.httpClient.ExecuteScriptAtBlockID(ctx, blockID, script, arguments)
}

func (c *Client) ExecuteScriptAtBlockHeight(
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

func (c *Client) GetEventsForHeightRange(
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

func (c *Client) GetEventsForBlockIDs(
	ctx context.Context,
	eventType string,
	blockIDs []flow.Identifier,
) ([]flow.BlockEvents, error) {
	return c.httpClient.GetEventsForBlockIDs(ctx, eventType, blockIDs)
}

func (c *Client) GetLatestProtocolStateSnapshot(ctx context.Context) ([]byte, error) {
	return c.httpClient.GetLatestProtocolStateSnapshot(ctx)
}

func (c *Client) GetExecutionResultForBlockID(ctx context.Context, blockID flow.Identifier) (*flow.ExecutionResult, error) {
	return c.httpClient.GetExecutionResultForBlockID(ctx, blockID)
}

func (c *Client) Close() error {
	// Close method is not required by the HTTP as the connection is setup and tear down with every request.
	return nil
}

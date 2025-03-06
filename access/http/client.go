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

package http

import (
	"context"
	"fmt"

	"github.com/onflow/flow-go-sdk/access"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"

	"github.com/onflow/flow-go-sdk"
)

const (
	EmulatorHost  = "http://127.0.0.1:8888/v1"
	TestnetHost   = "https://rest-testnet.onflow.org/v1"
	MainnetHost   = "https://rest-mainnet.onflow.org/v1"
	CanarynetHost = "https://rest-canary.onflow.org/v1"
)

// ClientOption is a configuration option for the client.
type ClientOption func(*options)

type options struct {
	jsonOptions []jsoncdc.Option
}

func DefaultClientOptions() *options {
	return &options{
		jsonOptions: []jsoncdc.Option{
			jsoncdc.WithAllowUnstructuredStaticTypes(true),
		},
	}
}

// WithJSONOptions wraps a json.Option into a ClientOption.
func WithJSONOptions(jsonOpts ...jsoncdc.Option) ClientOption {
	return func(opts *options) {
		opts.jsonOptions = append(opts.jsonOptions, jsonOpts...)
	}
}

// NewClient creates an HTTP client exposing all the common access APIs.
// Client will use provided host for connection.
func NewClient(host string, opts ...ClientOption) (*Client, error) {
	cfg := DefaultClientOptions()
	for _, apply := range opts {
		apply(cfg)
	}

	client, err := NewBaseClient(host)
	if err != nil {
		return nil, err
	}

	client.SetJSONOptions(cfg.jsonOptions)

	return &Client{
		UnimplementedClient: &access.UnimplementedClient{},
		httpClient:          client,
	}, nil
}

var _ access.Client = &Client{}

// Client implements all common HTTP methods providing a network agnostic API.
type Client struct {
	*access.UnimplementedClient
	httpClient *BaseClient
}

func (c *Client) Ping(ctx context.Context) error {
	return c.httpClient.Ping(ctx)
}

func (c *Client) GetNetworkParameters(ctx context.Context) (*flow.NetworkParameters, error) {
	return c.httpClient.GetNetworkParameters(ctx)
}

func (c *Client) GetNodeVersionInfo(ctx context.Context) (*flow.NodeVersionInfo, error) {
	return c.httpClient.GetNodeVersionInfo(ctx)
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

func (c *Client) GetCollectionByID(ctx context.Context, ID flow.Identifier) (*flow.Collection, error) {
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

func (c *Client) GetAccountBalanceAtLatestBlock(ctx context.Context, address flow.Address) (uint64, error) {
	account, err := c.GetAccountAtLatestBlock(ctx, address)
	if err != nil {
		return 0, err
	}
	return account.Balance, nil
}

func (c *Client) GetAccountBalanceAtBlockHeight(ctx context.Context, address flow.Address, blockHeight uint64) (uint64, error) {
	account, err := c.GetAccountAtBlockHeight(ctx, address, blockHeight)
	if err != nil {
		return 0, err
	}
	return account.Balance, nil
}

func (c *Client) GetAccountKeyAtLatestBlock(ctx context.Context, address flow.Address, keyIndex uint32) (*flow.AccountKey, error) {
	account, err := c.GetAccountAtLatestBlock(ctx, address)
	if err != nil {
		return nil, err
	}
	if keyIndex >= uint32(len(account.Keys)) {
		return nil, fmt.Errorf("key index out of bounds")
	}

	return account.Keys[keyIndex], nil
}

func (c *Client) GetAccountKeyAtBlockHeight(ctx context.Context, address flow.Address, keyIndex uint32, height uint64) (*flow.AccountKey, error) {
	account, err := c.GetAccountAtBlockHeight(ctx, address, height)
	if err != nil {
		return nil, err
	}
	if keyIndex >= uint32(len(account.Keys)) {
		return nil, fmt.Errorf("key index out of bounds")
	}

	return account.Keys[keyIndex], nil
}

func (c *Client) GetAccountKeysAtLatestBlock(ctx context.Context, address flow.Address) ([]*flow.AccountKey, error) {
	account, err := c.GetAccountAtLatestBlock(ctx, address)
	if err != nil {
		return nil, err
	}
	return account.Keys, nil
}

func (c *Client) GetAccountKeysAtBlockHeight(ctx context.Context, address flow.Address, height uint64) ([]*flow.AccountKey, error) {
	account, err := c.GetAccountAtBlockHeight(ctx, address, height)
	if err != nil {
		return nil, err
	}
	return account.Keys, nil
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

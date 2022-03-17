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
	"fmt"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk/client/convert"

	"github.com/onflow/flow-go-sdk"
)

const SEALED_HEIGHT = "sealed"
const EMULATOR_API = "http://127.0.0.1:8888/v1"
const TESTNET_API = "https://rest-testnet.onflow.org/v1/"
const MAINNET_API = "https://rest-mainnet.onflow.org/v1/"
const CANARYNET_API = ""

func NewClient(handler *handler) *Client {
	return &Client{handler}
}

func NewDefaultEmulatorClient(debug bool) (*Client, error) {
	handler, err := newHandler(EMULATOR_API, debug)
	if err != nil {
		return nil, err
	}

	return NewClient(handler), nil
}

func NewDefaultTestnetClient() (*Client, error) {
	handler, err := newHandler(TESTNET_API, false)
	if err != nil {
		return nil, err
	}

	return NewClient(handler), nil
}

func NewDefaultCanaryClient() (*Client, error) {
	handler, err := newHandler(CANARYNET_API, false)
	if err != nil {
		return nil, err
	}

	return NewClient(handler), nil
}

func NewDefaultMainnetClient() (*Client, error) {
	handler, err := newHandler(MAINNET_API, false)
	if err != nil {
		return nil, err
	}

	return NewClient(handler), nil
}

type Client struct {
	handler *handler
}

func (c *Client) Ping(ctx context.Context) error {
	panic("implement me")
}

func (c *Client) GetBlockByID(ctx context.Context, blockID flow.Identifier) (*flow.Block, error) {
	block, err := c.handler.getBlockByID(ctx, blockID.String())
	if err != nil {
		return nil, err
	}

	return convert.HTTPToBlock(block), nil
}

func (c *Client) GetLatestBlockHeader(ctx context.Context, isSealed bool) (*flow.BlockHeader, error) {
	block, err := c.GetLatestBlock(ctx, isSealed)
	if err != nil {
		return nil, err
	}

	return &block.BlockHeader, nil
}

func (c *Client) GetBlockHeaderByID(ctx context.Context, blockID flow.Identifier) (*flow.BlockHeader, error) {
	block, err := c.GetBlockByID(ctx, blockID)
	if err != nil {
		return nil, err
	}

	return &block.BlockHeader, nil
}

func (c *Client) GetBlockHeaderByHeight(ctx context.Context, height uint64) (*flow.BlockHeader, error) {
	block, err := c.GetBlockByHeight(ctx, height)
	if err != nil {
		return nil, err
	}

	return &block.BlockHeader, nil
}

func (c *Client) GetLatestBlock(ctx context.Context, isSealed bool) (*flow.Block, error) {
	blocks, err := c.handler.getBlockByHeight(ctx, convert.SealedToHTTP(isSealed))
	if err != nil {
		return nil, err
	}

	return convert.HTTPToBlock(blocks[0]), nil
}

func (c *Client) GetBlockByHeight(ctx context.Context, height uint64) (*flow.Block, error) {
	blocks, err := c.handler.getBlockByHeight(ctx, fmt.Sprintf("%d", height))
	if err != nil {
		return nil, err
	}

	return convert.HTTPToBlock(blocks[0]), nil
}

func (c *Client) GetCollection(ctx context.Context, ID flow.Identifier) (*flow.Collection, error) {
	collection, err := c.handler.getCollection(ctx, ID.String())
	if err != nil {
		return nil, err
	}

	return convert.HTTPToCollection(collection), nil
}

func (c *Client) SendTransaction(ctx context.Context, tx flow.Transaction) error {
	convertedTx, err := convert.TransactionToHTTP(tx)
	if err != nil {
		return err
	}

	return c.handler.sendTransaction(ctx, convertedTx)
}

func (c *Client) GetTransaction(ctx context.Context, ID flow.Identifier) (*flow.Transaction, error) {
	tx, err := c.handler.getTransaction(ctx, ID.String(), false)
	if err != nil {
		return nil, err
	}

	return convert.HTTPToTransaction(tx)
}

func (c *Client) GetTransactionResult(ctx context.Context, ID flow.Identifier) (*flow.TransactionResult, error) {
	tx, err := c.handler.getTransaction(ctx, ID.String(), true)
	if err != nil {
		return nil, err
	}

	return convert.HTTPToTransactionResult(tx.Result), nil
}

func (c *Client) GetAccount(ctx context.Context, address flow.Address) (*flow.Account, error) {
	account, err := c.handler.getAccount(ctx, address.String(), SEALED_HEIGHT)
	if err != nil {
		return nil, err
	}

	return convert.HTTPToAccount(account)
}

func (c *Client) GetAccountAtLatestBlock(ctx context.Context, address flow.Address) (*flow.Account, error) {
	return c.GetAccount(ctx, address)
}

func (c *Client) GetAccountAtBlockHeight(ctx context.Context, address flow.Address, blockHeight uint64) (*flow.Account, error) {
	account, err := c.handler.getAccount(ctx, address.String(), fmt.Sprintf("%d", blockHeight))
	if err != nil {
		return nil, err
	}

	return convert.HTTPToAccount(account)
}

func (c *Client) ExecuteScriptAtLatestBlock(ctx context.Context, script []byte, arguments []cadence.Value) (cadence.Value, error) {
	result, err := c.handler.executeScriptAtBlockHeight(
		ctx,
		SEALED_HEIGHT,
		convert.ScriptToHTTP(script),
		convert.CadenceArgsToHTTP(arguments),
	)
	if err != nil {
		return nil, err
	}

	return cadence.NewString(result)
}

func (c *Client) ExecuteScriptAtBlockID(ctx context.Context, blockID flow.Identifier, script []byte, arguments []cadence.Value) (cadence.Value, error) {
	result, err := c.handler.executeScriptAtBlockID(
		ctx,
		blockID.String(),
		convert.ScriptToHTTP(script),
		convert.CadenceArgsToHTTP(arguments),
	)
	if err != nil {
		return nil, err
	}

	return cadence.NewString(result)
}

func (c *Client) ExecuteScriptAtBlockHeight(ctx context.Context, height uint64, script []byte, arguments []cadence.Value) (cadence.Value, error) {
	result, err := c.handler.executeScriptAtBlockHeight(
		ctx,
		fmt.Sprintf("%d", height),
		convert.ScriptToHTTP(script),
		convert.CadenceArgsToHTTP(arguments),
	)
	if err != nil {
		return nil, err
	}

	return cadence.NewString(result)
}

func (c *Client) GetEventsForHeightRange(ctx context.Context, eventType string, startHeight uint64, endHeight uint64) ([]flow.BlockEvents, error) {
	panic("implement me")
}

func (c *Client) GetEventsForBlockIDs(ctx context.Context, eventType string, blockIDs []flow.Identifier) ([]flow.BlockEvents, error) {
	panic("implement me")
}

func (c *Client) GetLatestProtocolStateSnapshot(ctx context.Context) ([]byte, error) {
	panic("implement me")
}

func (c *Client) GetExecutionResultForBlockID(ctx context.Context, blockID flow.Identifier) (*flow.ExecutionResult, error) {
	panic("implement me")
}

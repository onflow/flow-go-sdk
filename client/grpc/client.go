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

func NewClient(handler *Handler) *Client {
	return &Client{
		handler: handler,
	}
}

type Client struct {
	handler *Handler
}

func NewDefaultEmulatorClient() (*Client, error) {
	handler, err := New(EMULATOR_API)
	if err != nil {
		return nil, err
	}

	return NewClient(handler), nil
}

func NewDefaultTestnetClient() (*Client, error) {
	handler, err := New(TESTNET_API)
	if err != nil {
		return nil, err
	}

	return NewClient(handler), nil
}

func NewDefaultCanaryClient() (*Client, error) {
	handler, err := New(CANARYNET_API)
	if err != nil {
		return nil, err
	}

	return NewClient(handler), nil
}

func NewDefaultMainnetClient() (*Client, error) {
	handler, err := New(MAINNET_API)
	if err != nil {
		return nil, err
	}

	return NewClient(handler), nil
}

func (c *Client) Ping(ctx context.Context) error {
	return c.handler.Ping(ctx)
}

func (c *Client) GetLatestBlockHeader(ctx context.Context, isSealed bool) (*flow.BlockHeader, error) {
	return c.handler.GetLatestBlockHeader(ctx, isSealed)
}

func (c *Client) GetBlockHeaderByID(ctx context.Context, blockID flow.Identifier) (*flow.BlockHeader, error) {
	return c.handler.GetBlockHeaderByID(ctx, blockID)
}

func (c *Client) GetBlockHeaderByHeight(ctx context.Context, height uint64) (*flow.BlockHeader, error) {
	return c.handler.GetBlockHeaderByHeight(ctx, height)
}

func (c *Client) GetLatestBlock(ctx context.Context, isSealed bool) (*flow.Block, error) {
	return c.handler.GetLatestBlock(ctx, isSealed)
}

func (c *Client) GetBlockByID(ctx context.Context, blockID flow.Identifier) (*flow.Block, error) {
	return c.handler.GetBlockByID(ctx, blockID)
}

func (c *Client) GetBlockByHeight(ctx context.Context, height uint64) (*flow.Block, error) {
	return c.handler.GetBlockByHeight(ctx, height)
}

func (c *Client) GetCollection(ctx context.Context, colID flow.Identifier) (*flow.Collection, error) {
	return c.handler.GetCollection(ctx, colID)
}

func (c *Client) SendTransaction(ctx context.Context, tx flow.Transaction) error {
	return c.handler.SendTransaction(ctx, tx)
}

func (c *Client) GetTransaction(ctx context.Context, txID flow.Identifier) (*flow.Transaction, error) {
	return c.handler.GetTransaction(ctx, txID)
}

func (c *Client) GetTransactionResult(ctx context.Context, txID flow.Identifier) (*flow.TransactionResult, error) {
	return c.handler.GetTransactionResult(ctx, txID)
}

func (c *Client) GetAccount(ctx context.Context, address flow.Address) (*flow.Account, error) {
	return c.handler.GetAccount(ctx, address)
}

func (c *Client) GetAccountAtLatestBlock(ctx context.Context, address flow.Address) (*flow.Account, error) {
	return c.handler.GetAccountAtLatestBlock(ctx, address)
}

func (c *Client) GetAccountAtBlockHeight(ctx context.Context, address flow.Address, blockHeight uint64) (*flow.Account, error) {
	return c.handler.GetAccountAtBlockHeight(ctx, address, blockHeight)
}

func (c *Client) ExecuteScriptAtLatestBlock(ctx context.Context, script []byte, arguments []cadence.Value) (cadence.Value, error) {
	return c.handler.ExecuteScriptAtLatestBlock(ctx, script, arguments)
}

func (c *Client) ExecuteScriptAtBlockID(ctx context.Context, blockID flow.Identifier, script []byte, arguments []cadence.Value) (cadence.Value, error) {
	return c.handler.ExecuteScriptAtBlockID(ctx, blockID, script, arguments)
}

func (c *Client) ExecuteScriptAtBlockHeight(ctx context.Context, height uint64, script []byte, arguments []cadence.Value) (cadence.Value, error) {
	return c.handler.ExecuteScriptAtBlockHeight(ctx, height, script, arguments)
}

func (c *Client) GetEventsForHeightRange(ctx context.Context, eventType string, startHeight uint64, endHeight uint64) ([]flow.BlockEvents, error) {
	return c.handler.GetEventsForHeightRange(ctx, EventRangeQuery{
		Type:        eventType,
		StartHeight: startHeight,
		EndHeight:   endHeight,
	})
}

func (c *Client) GetEventsForBlockIDs(ctx context.Context, eventType string, blockIDs []flow.Identifier) ([]flow.BlockEvents, error) {
	return c.handler.GetEventsForBlockIDs(ctx, eventType, blockIDs)
}

func (c *Client) GetLatestProtocolStateSnapshot(ctx context.Context) ([]byte, error) {
	return c.handler.GetLatestProtocolStateSnapshot(ctx)
}

func (c *Client) GetExecutionResultForBlockID(ctx context.Context, blockID flow.Identifier) (*flow.ExecutionResult, error) {
	return c.handler.GetExecutionResultForBlockID(ctx, blockID)
}

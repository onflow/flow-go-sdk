package grpc

import (
	"context"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
)

type Client struct {
	Handler Handler
}

func (c Client) Ping(ctx context.Context) error {
	return c.Handler.Ping(ctx)
}

func (c Client) GetLatestBlockHeader(ctx context.Context, isSealed bool) (*flow.BlockHeader, error) {
	return c.Handler.GetLatestBlockHeader(ctx, isSealed)
}

func (c Client) GetBlockHeaderByID(ctx context.Context, blockID flow.Identifier) (*flow.BlockHeader, error) {
	return c.Handler.GetBlockHeaderByID(ctx, blockID)
}

func (c Client) GetBlockHeaderByHeight(ctx context.Context, height uint64) (*flow.BlockHeader, error) {
	return c.Handler.GetBlockHeaderByHeight(ctx, height)
}

func (c Client) GetLatestBlock(ctx context.Context, isSealed bool) (*flow.Block, error) {
	return c.Handler.GetLatestBlock(ctx, isSealed)
}

func (c Client) GetBlockByID(ctx context.Context, blockID flow.Identifier) (*flow.Block, error) {
	return c.Handler.GetBlockByID(ctx, blockID)
}

func (c Client) GetBlockByHeight(ctx context.Context, height uint64) (*flow.Block, error) {
	return c.Handler.GetBlockByHeight(ctx, height)
}

func (c Client) GetCollection(ctx context.Context, colID flow.Identifier) (*flow.Collection, error) {
	return c.Handler.GetCollection(ctx, colID)
}

func (c Client) SendTransaction(ctx context.Context, tx flow.Transaction) error {
	return c.Handler.SendTransaction(ctx, tx)
}

func (c Client) GetTransaction(ctx context.Context, txID flow.Identifier) (*flow.Transaction, error) {
	return c.Handler.GetTransaction(ctx, txID)
}

func (c Client) GetTransactionResult(ctx context.Context, txID flow.Identifier) (*flow.TransactionResult, error) {
	return c.Handler.GetTransactionResult(ctx, txID)
}

func (c Client) GetAccount(ctx context.Context, address flow.Address) (*flow.Account, error) {
	return c.Handler.GetAccount(ctx, address)
}

func (c Client) GetAccountAtLatestBlock(ctx context.Context, address flow.Address) (*flow.Account, error) {
	return c.Handler.GetAccountAtLatestBlock(ctx, address)
}

func (c Client) GetAccountAtBlockHeight(ctx context.Context, address flow.Address, blockHeight uint64) (*flow.Account, error) {
	return c.Handler.GetAccountAtBlockHeight(ctx, address, blockHeight)
}

func (c Client) ExecuteScriptAtLatestBlock(ctx context.Context, script []byte, arguments []cadence.Value) (cadence.Value, error) {
	return c.Handler.ExecuteScriptAtLatestBlock(ctx, script, arguments)
}

func (c Client) ExecuteScriptAtBlockID(ctx context.Context, blockID flow.Identifier, script []byte, arguments []cadence.Value) (cadence.Value, error) {
	return c.Handler.ExecuteScriptAtBlockID(ctx, blockID, script, arguments)
}

func (c Client) ExecuteScriptAtBlockHeight(ctx context.Context, height uint64, script []byte, arguments []cadence.Value) (cadence.Value, error) {
	return c.Handler.ExecuteScriptAtBlockHeight(ctx, height, script, arguments)
}

func (c Client) GetEventsForHeightRange(ctx context.Context, eventType string, startHeight uint64, endHeight uint64) ([]flow.BlockEvents, error) {
	return c.Handler.GetEventsForHeightRange(ctx, EventRangeQuery{
		Type:        eventType,
		StartHeight: startHeight,
		EndHeight:   endHeight,
	})
}

func (c Client) GetEventsForBlockIDs(ctx context.Context, eventType string, blockIDs []flow.Identifier) ([]flow.BlockEvents, error) {
	return c.Handler.GetEventsForBlockIDs(ctx, eventType, blockIDs)
}

func (c Client) GetLatestProtocolStateSnapshot(ctx context.Context) ([]byte, error) {
	return c.Handler.GetLatestProtocolStateSnapshot(ctx)
}

func (c Client) GetExecutionResultForBlockID(ctx context.Context, blockID flow.Identifier) (*flow.ExecutionResult, error) {
	return c.Handler.GetExecutionResultForBlockID(ctx, blockID)
}

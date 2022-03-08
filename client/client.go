package client

import (
	"context"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
)

type Client interface {
	Ping(ctx context.Context) error
	GetLatestBlockHeader(ctx context.Context, isSealed bool) (*flow.BlockHeader, error)
	GetBlockHeaderByID(ctx context.Context, blockID flow.Identifier) (*flow.BlockHeader, error)
	GetBlockHeaderByHeight(ctx context.Context, height uint64) (*flow.BlockHeader, error)
	GetLatestBlock(ctx context.Context, isSealed bool) (*flow.Block, error)
	GetBlockByID(ctx context.Context, blockID flow.Identifier) (*flow.Block, error)
	GetBlockByHeight(ctx context.Context, height uint64) (*flow.Block, error)
	GetCollection(ctx context.Context, colID flow.Identifier) (*flow.Collection, error)
	SendTransaction(ctx context.Context, tx flow.Transaction) error
	GetTransaction(ctx context.Context, txID flow.Identifier) (*flow.Transaction, error)
	GetTransactionResult(ctx context.Context, txID flow.Identifier) (*flow.TransactionResult, error)
	GetAccount(ctx context.Context, address flow.Address) (*flow.Account, error)
	GetAccountAtLatestBlock(ctx context.Context, address flow.Address) (*flow.Account, error)
	GetAccountAtBlockHeight(ctx context.Context, address flow.Address, blockHeight uint64) (*flow.Account, error)
	ExecuteScriptAtLatestBlock(ctx context.Context, script []byte, arguments []cadence.Value) (cadence.Value, error)
	ExecuteScriptAtBlockID(ctx context.Context, blockID flow.Identifier, script []byte, arguments []cadence.Value) (cadence.Value, error)
	ExecuteScriptAtBlockHeight(ctx context.Context, height uint64, script []byte, arguments []cadence.Value) (cadence.Value, error)
	GetEventsForHeightRange(ctx context.Context, eventType string, startHeight uint64, endHeight uint64) ([]flow.BlockEvents, error)
	GetEventsForBlockIDs(ctx context.Context, eventType string, blockIDs []flow.Identifier) ([]flow.BlockEvents, error)
	GetLatestProtocolStateSnapshot(ctx context.Context) ([]byte, error)
	GetExecutionResultForBlockID(ctx context.Context, blockID flow.Identifier) (*flow.ExecutionResult, error)
}

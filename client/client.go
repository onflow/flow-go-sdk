package client

import (
	"context"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
)

type Client interface {
	// Ping is used to check if the access node is alive and healthy.
	Ping(ctx context.Context) error

	// GetLatestBlockHeader gets the latest sealed or unsealed block header.
	GetLatestBlockHeader(ctx context.Context, isSealed bool) (*flow.BlockHeader, error)

	// GetBlockHeaderByID gets a block header by ID.
	GetBlockHeaderByID(ctx context.Context, blockID flow.Identifier) (*flow.BlockHeader, error)

	// GetBlockHeaderByHeight gets a block header by height.
	GetBlockHeaderByHeight(ctx context.Context, height uint64) (*flow.BlockHeader, error)

	// GetLatestBlock gets the full payload of the latest sealed or unsealed block.
	GetLatestBlock(ctx context.Context, isSealed bool) (*flow.Block, error)

	// GetBlockByID gets a full block by ID.
	GetBlockByID(ctx context.Context, blockID flow.Identifier) (*flow.Block, error)

	// GetBlockByHeight gets a full block by height.
	GetBlockByHeight(ctx context.Context, height uint64) (*flow.Block, error)

	// GetCollection gets a collection by ID.
	GetCollection(ctx context.Context, colID flow.Identifier) (*flow.Collection, error)

	// SendTransaction submits a transaction to the network.
	SendTransaction(ctx context.Context, tx flow.Transaction) error

	// GetTransaction gets a transaction by ID.
	GetTransaction(ctx context.Context, txID flow.Identifier) (*flow.Transaction, error)

	// GetTransactionResult gets the result of a transaction.
	GetTransactionResult(ctx context.Context, txID flow.Identifier) (*flow.TransactionResult, error)

	// GetAccount is an alias for GetAccountAtLatestBlock.
	GetAccount(ctx context.Context, address flow.Address) (*flow.Account, error)

	// GetAccountAtLatestBlock gets an account by address at the latest sealed block.
	GetAccountAtLatestBlock(ctx context.Context, address flow.Address) (*flow.Account, error)

	// GetAccountAtBlockHeight gets an account by address at the given block height
	GetAccountAtBlockHeight(ctx context.Context, address flow.Address, blockHeight uint64) (*flow.Account, error)

	// ExecuteScriptAtLatestBlock executes a read-only Cadence script against the latest sealed execution state.
	ExecuteScriptAtLatestBlock(ctx context.Context, script []byte, arguments []cadence.Value) (cadence.Value, error)

	// ExecuteScriptAtBlockID executes a ready-only Cadence script against the execution state at the block with the given ID.
	ExecuteScriptAtBlockID(ctx context.Context, blockID flow.Identifier, script []byte, arguments []cadence.Value) (cadence.Value, error)

	// ExecuteScriptAtBlockHeight executes a ready-only Cadence script against the execution state at the given block height.
	ExecuteScriptAtBlockHeight(ctx context.Context, height uint64, script []byte, arguments []cadence.Value) (cadence.Value, error)

	// GetEventsForHeightRange retrieves events for all sealed blocks between the start and end block heights (inclusive) with the given type.
	GetEventsForHeightRange(ctx context.Context, eventType string, startHeight uint64, endHeight uint64) ([]flow.BlockEvents, error)

	// GetEventsForBlockIDs retrieves events with the given type from the specified block IDs.
	GetEventsForBlockIDs(ctx context.Context, eventType string, blockIDs []flow.Identifier) ([]flow.BlockEvents, error)

	// GetLatestProtocolStateSnapshot retrieves the latest snapshot of the protocol
	// state in serialized form. This is used to generate a root snapshot file
	// used by Flow nodes to bootstrap their local protocol state database.
	GetLatestProtocolStateSnapshot(ctx context.Context) ([]byte, error)

	// GetExecutionResultForBlockID gets the execution results at the block ID.
	GetExecutionResultForBlockID(ctx context.Context, blockID flow.Identifier) (*flow.ExecutionResult, error)
}

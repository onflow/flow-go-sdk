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

// Package access contains an interface defining functions for the API clients.
//
// The Access API provides a set of methods that can be used to submit transactions
// and read state from Flow. Clients are compatible with the Access API implemented by the
// Access Node role, as well as the mock Access API exposed by the Flow Emulator.
package access

import (
	"context"

	"github.com/onflow/cadence"

	"github.com/onflow/flow-go-sdk"
)

//go:generate go run github.com/vektra/mockery/cmd/mockery --name Client --structname Client --output mocks

type Client interface {
	// Ping is used to check if the access node is alive and healthy.
	Ping(ctx context.Context) error

	// GetNetworkParameters returns the network parameters.
	GetNetworkParameters(ctx context.Context) (*flow.NetworkParameters, error)

	// GetNodeVersionInfo returns the node information about the network.
	GetNodeVersionInfo(ctx context.Context) (*flow.NodeVersionInfo, error)

	// GetLatestBlockHeader returns the latest sealed or unsealed block header.
	GetLatestBlockHeader(ctx context.Context, isSealed bool) (*flow.BlockHeader, error)

	// GetBlockHeaderByID returns a block header by ID.
	GetBlockHeaderByID(ctx context.Context, blockID flow.Identifier) (*flow.BlockHeader, error)

	// GetBlockHeaderByHeight returns a block header by height.
	GetBlockHeaderByHeight(ctx context.Context, height uint64) (*flow.BlockHeader, error)

	// GetLatestBlock returns the full payload of the latest sealed or unsealed block.
	GetLatestBlock(ctx context.Context, isSealed bool) (*flow.Block, error)

	// GetBlockByID returns a full block by ID.
	GetBlockByID(ctx context.Context, blockID flow.Identifier) (*flow.Block, error)

	// GetBlockByHeight returns a full block by height.
	GetBlockByHeight(ctx context.Context, height uint64) (*flow.Block, error)

	// GetCollection returns a collection by ID.
	GetCollection(ctx context.Context, colID flow.Identifier) (*flow.Collection, error)

	// GetCollectionByID returns a collection by ID.
	GetCollectionByID(ctx context.Context, id flow.Identifier) (*flow.Collection, error)

	// GetFullCollectionByID returns a full collection including transaction bodies by ID.
	GetFullCollectionByID(ctx context.Context, id flow.Identifier) (*flow.FullCollection, error)

	// SendTransaction submits a transaction to the network.
	SendTransaction(ctx context.Context, tx flow.Transaction) error

	// GetTransaction returns a transaction by ID.
	GetTransaction(ctx context.Context, txID flow.Identifier) (*flow.Transaction, error)

	// GetTransactionsByBlockID returns all the transactions for a specified block.
	GetTransactionsByBlockID(ctx context.Context, blockID flow.Identifier) ([]*flow.Transaction, error)

	// GetTransactionResult returns the result of a transaction.
	GetTransactionResult(ctx context.Context, txID flow.Identifier) (*flow.TransactionResult, error)

	// GetTransactionResultByIndex returns a transaction result by transaction index for the given block ID.
	GetTransactionResultByIndex(ctx context.Context, blockID flow.Identifier, index uint32) (*flow.TransactionResult, error)

	// GetTransactionResultsByBlockID returns all the transaction results for a specified block.
	GetTransactionResultsByBlockID(ctx context.Context, blockID flow.Identifier) ([]*flow.TransactionResult, error)

	// GetSystemTransaction returns the system transaction for the given block ID.
	GetSystemTransaction(ctx context.Context, blockID flow.Identifier) (*flow.Transaction, error)

	// GetSystemTransactionResult returns the transaction result of the system transaction for the given block ID.
	GetSystemTransactionResult(ctx context.Context, blockID flow.Identifier) (*flow.TransactionResult, error)

	// GetAccount is an alias for GetAccountAtLatestBlock.
	GetAccount(ctx context.Context, address flow.Address) (*flow.Account, error)

	// GetAccountAtLatestBlock returns an account by address at the latest sealed block.
	GetAccountAtLatestBlock(ctx context.Context, address flow.Address) (*flow.Account, error)

	// GetAccountAtBlockHeight returns an account by address at the given block height
	GetAccountAtBlockHeight(ctx context.Context, address flow.Address, blockHeight uint64) (*flow.Account, error)

	// GetAccountBalanceAtLatestBlock returns the balance of an account at the latest sealed block.
	GetAccountBalanceAtLatestBlock(ctx context.Context, address flow.Address) (uint64, error)

	// GetAccountBalanceAtBlockHeight returns the balance of an account at the given block height.
	GetAccountBalanceAtBlockHeight(ctx context.Context, address flow.Address, blockHeight uint64) (uint64, error)

	// GetAccountKeyAtLatestBlock returns the account key with the provided index at the latest sealed block.
	GetAccountKeyAtLatestBlock(ctx context.Context, address flow.Address, keyIndex uint32) (*flow.AccountKey, error)

	// GetAccountKeyAtBlockHeight returns the account key with the provided index at the given block height.
	GetAccountKeyAtBlockHeight(ctx context.Context, address flow.Address, keyIndex uint32, height uint64) (*flow.AccountKey, error)

	// GetAccountKeysAtLatestBlock returns all account keys at the latest sealed block.
	GetAccountKeysAtLatestBlock(ctx context.Context, address flow.Address) ([]*flow.AccountKey, error)

	// GetAccountKeysAtBlockHeight returns all account keys at the given block height.
	GetAccountKeysAtBlockHeight(ctx context.Context, address flow.Address, height uint64) ([]*flow.AccountKey, error)

	// ExecuteScriptAtLatestBlock executes a read-only Cadence script against the latest sealed execution state.
	ExecuteScriptAtLatestBlock(ctx context.Context, script []byte, arguments []cadence.Value) (cadence.Value, error)

	// ExecuteScriptAtBlockID executes a ready-only Cadence script against the execution state at the block with the given ID.
	ExecuteScriptAtBlockID(ctx context.Context, blockID flow.Identifier, script []byte, arguments []cadence.Value) (cadence.Value, error)

	// ExecuteScriptAtBlockHeight executes a ready-only Cadence script against the execution state at the given block height.
	ExecuteScriptAtBlockHeight(ctx context.Context, height uint64, script []byte, arguments []cadence.Value) (cadence.Value, error)

	// GetEventsForHeightRange returns events for all sealed blocks between the start and end block heights (inclusive) with the given type.
	GetEventsForHeightRange(ctx context.Context, eventType string, startHeight uint64, endHeight uint64) ([]flow.BlockEvents, error)

	// GetEventsForBlockIDs returns events with the given type from the specified block IDs.
	GetEventsForBlockIDs(ctx context.Context, eventType string, blockIDs []flow.Identifier) ([]flow.BlockEvents, error)

	// GetLatestProtocolStateSnapshot returns the protocol state snapshot in serialized form at latest sealed block.
	// This is used to generate a root snapshot file used by Flow nodes to bootstrap their local protocol state database.
	GetLatestProtocolStateSnapshot(ctx context.Context) ([]byte, error)

	// GetProtocolStateSnapshotByBlockID returns the protocol state snapshot in serialized form at the given block ID.
	// This is used to generate a root snapshot file used by Flow nodes to bootstrap their local protocol state database.
	GetProtocolStateSnapshotByBlockID(ctx context.Context, blockID flow.Identifier) ([]byte, error)

	// GetProtocolStateSnapshotByHeight returns the protocol state snapshot in serialized form at the given block height.
	// This is used to generate a root snapshot file used by Flow nodes to bootstrap their local protocol state database.
	GetProtocolStateSnapshotByHeight(ctx context.Context, blockHeight uint64) ([]byte, error)

	// GetExecutionResultByID returns the execution result by ID.
	GetExecutionResultByID(ctx context.Context, id flow.Identifier) (*flow.ExecutionResult, error)

	// GetExecutionResultForBlockID returns the execution results at the block ID.
	GetExecutionResultForBlockID(ctx context.Context, blockID flow.Identifier) (*flow.ExecutionResult, error)

	// GetExecutionDataByBlockID returns execution data for a specific block ID.
	GetExecutionDataByBlockID(ctx context.Context, blockID flow.Identifier) (*flow.ExecutionData, error)

	// SubscribeExecutionDataByBlockID subscribes to execution data updates starting at the given block ID.
	SubscribeExecutionDataByBlockID(
		ctx context.Context,
		startBlockID flow.Identifier,
	) (<-chan *flow.ExecutionDataStreamResponse, <-chan error, error)

	// SubscribeExecutionDataByBlockHeight subscribes to execution data updates starting at the given block height.
	SubscribeExecutionDataByBlockHeight(
		ctx context.Context,
		startHeight uint64,
	) (<-chan *flow.ExecutionDataStreamResponse, <-chan error, error)

	// SubscribeEventsByBlockID subscribes to events starting at the given block ID.
	SubscribeEventsByBlockID(
		ctx context.Context,
		startBlockID flow.Identifier,
		filter flow.EventFilter,
		opts ...SubscribeOption,
	) (<-chan flow.BlockEvents, <-chan error, error)

	// SubscribeEventsByBlockHeight subscribes to events starting at the given block height.
	SubscribeEventsByBlockHeight(
		ctx context.Context,
		startHeight uint64,
		filter flow.EventFilter,
		opts ...SubscribeOption,
	) (<-chan flow.BlockEvents, <-chan error, error)

	// SubscribeBlockDigestsFromStartBlockID subscribes to block digests with the given status
	// starting at the given block ID.
	// The status may be either flow.BlockStatusFinalized or flow.BlockStatusSealed
	SubscribeBlockDigestsFromStartBlockID(
		ctx context.Context,
		startBlockID flow.Identifier,
		blockStatus flow.BlockStatus,
	) (<-chan *flow.BlockDigest, <-chan error, error)

	// SubscribeBlockDigestsFromStartHeight subscribes to block digests with the given status
	// starting at the given block height.
	// The status may be either flow.BlockStatusFinalized or flow.BlockStatusSealed
	SubscribeBlockDigestsFromStartHeight(
		ctx context.Context,
		startHeight uint64,
		blockStatus flow.BlockStatus,
	) (<-chan *flow.BlockDigest, <-chan error, error)

	// SubscribeBlockDigestsFromLatest subscribes to block digests with the given status
	// starting at the latest block.
	// The status may be either flow.BlockStatusFinalized or flow.BlockStatusSealed
	SubscribeBlockDigestsFromLatest(
		ctx context.Context,
		blockStatus flow.BlockStatus,
	) (<-chan *flow.BlockDigest, <-chan error, error)

	// SubscribeBlocksFromStartBlockID subscribes to blocks with the given status starting at the
	// given block ID.
	// The status may be either flow.BlockStatusFinalized or flow.BlockStatusSealed
	SubscribeBlocksFromStartBlockID(
		ctx context.Context,
		startBlockID flow.Identifier,
		blockStatus flow.BlockStatus,
	) (<-chan *flow.Block, <-chan error, error)

	// SubscribeBlocksFromStartHeight subscribes to blocks with the given status starting at the
	// given block height.
	// The status may be either flow.BlockStatusFinalized or flow.BlockStatusSealed
	SubscribeBlocksFromStartHeight(
		ctx context.Context,
		startHeight uint64,
		blockStatus flow.BlockStatus,
	) (<-chan *flow.Block, <-chan error, error)

	// SubscribeBlocksFromLatest subscribes to blocks with the given status starting at the latest block.
	// The status may be either flow.BlockStatusFinalized or flow.BlockStatusSealed
	SubscribeBlocksFromLatest(
		ctx context.Context,
		blockStatus flow.BlockStatus,
	) (<-chan *flow.Block, <-chan error, error)

	// SubscribeBlockHeadersFromStartBlockID subscribes to block headers with the given status starting
	// at the given block ID.
	// The status may be either flow.BlockStatusFinalized or flow.BlockStatusSealed
	SubscribeBlockHeadersFromStartBlockID(
		ctx context.Context,
		startBlockID flow.Identifier,
		blockStatus flow.BlockStatus,
	) (<-chan *flow.BlockHeader, <-chan error, error)

	// SubscribeBlockHeadersFromStartHeight subscribes to block headers with the given status starting
	// at the given block height.
	// The status may be either flow.BlockStatusFinalized or flow.BlockStatusSealed
	SubscribeBlockHeadersFromStartHeight(
		ctx context.Context,
		startHeight uint64,
		blockStatus flow.BlockStatus,
	) (<-chan *flow.BlockHeader, <-chan error, error)

	// SubscribeBlockHeadersFromLatest subscribes to block headers with the given status starting
	// at the latest block.
	// The status may be either flow.BlockStatusFinalized or flow.BlockStatusSealed
	SubscribeBlockHeadersFromLatest(
		ctx context.Context,
		blockStatus flow.BlockStatus,
	) (<-chan *flow.BlockHeader, <-chan error, error)

	// SubscribeAccountStatusesFromStartHeight subscribes to account status events starting at the
	// given block height.
	SubscribeAccountStatusesFromStartHeight(
		ctx context.Context,
		startBlockHeight uint64,
		filter flow.AccountStatusFilter,
	) (<-chan *flow.AccountStatus, <-chan error, error)

	// SubscribeAccountStatusesFromStartBlockID subscribes to account status events starting at the
	// given block ID.
	SubscribeAccountStatusesFromStartBlockID(
		ctx context.Context,
		startBlockID flow.Identifier,
		filter flow.AccountStatusFilter,
	) (<-chan *flow.AccountStatus, <-chan error, error)

	// SubscribeAccountStatusesFromLatestBlock subscribes to account status events starting at the
	// latest block.
	SubscribeAccountStatusesFromLatestBlock(
		ctx context.Context,
		filter flow.AccountStatusFilter,
	) (<-chan *flow.AccountStatus, <-chan error, error)

	// SendAndSubscribeTransactionStatuses submits a transaction to the network and subscribes to the
	// transaction status updates.
	SendAndSubscribeTransactionStatuses(
		ctx context.Context,
		tx flow.Transaction,
	) (<-chan *flow.TransactionResult, <-chan error, error)

	// Close stops the client connection to the access node.
	Close() error
}

type SubscribeOption func(*SubscribeConfig)

type SubscribeConfig struct {
	HeartbeatInterval uint64
}

func WithHeartbeatInterval(interval uint64) SubscribeOption {
	return func(config *SubscribeConfig) {
		config.HeartbeatInterval = interval
	}
}

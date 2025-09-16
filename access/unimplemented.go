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
	"errors"

	"github.com/onflow/cadence"

	"github.com/onflow/flow-go-sdk"
)

var ErrUnimplemented = errors.New("not implemented")

var _ Client = &UnimplementedClient{}

type UnimplementedClient struct {
}

// Ping is used to check if the access node is alive and healthy.
func (c *UnimplementedClient) Ping(ctx context.Context) error {
	return ErrUnimplemented
}

// GetNetworkParameters returns the network parameters.
func (c *UnimplementedClient) GetNetworkParameters(ctx context.Context) (*flow.NetworkParameters, error) {
	return nil, ErrUnimplemented
}

// GetNodeVersionInfo returns the node information about the network.
func (c *UnimplementedClient) GetNodeVersionInfo(ctx context.Context) (*flow.NodeVersionInfo, error) {
	return nil, ErrUnimplemented
}

// GetLatestBlockHeader returns the latest sealed or unsealed block header.
func (c *UnimplementedClient) GetLatestBlockHeader(ctx context.Context, isSealed bool) (*flow.BlockHeader, error) {
	return nil, ErrUnimplemented
}

// GetBlockHeaderByID returns a block header by ID.
func (c *UnimplementedClient) GetBlockHeaderByID(ctx context.Context, blockID flow.Identifier) (*flow.BlockHeader, error) {
	return nil, ErrUnimplemented
}

// GetBlockHeaderByHeight returns a block header by height.
func (c *UnimplementedClient) GetBlockHeaderByHeight(ctx context.Context, height uint64) (*flow.BlockHeader, error) {
	return nil, ErrUnimplemented
}

// GetLatestBlock returns the full payload of the latest sealed or unsealed block.
func (c *UnimplementedClient) GetLatestBlock(ctx context.Context, isSealed bool) (*flow.Block, error) {
	return nil, ErrUnimplemented
}

// GetBlockByID returns a full block by ID.
func (c *UnimplementedClient) GetBlockByID(ctx context.Context, blockID flow.Identifier) (*flow.Block, error) {
	return nil, ErrUnimplemented
}

// GetBlockByHeight returns a full block by height.
func (c *UnimplementedClient) GetBlockByHeight(ctx context.Context, height uint64) (*flow.Block, error) {
	return nil, ErrUnimplemented
}

// GetCollection returns a collection by ID.
func (c *UnimplementedClient) GetCollection(ctx context.Context, colID flow.Identifier) (*flow.Collection, error) {
	return nil, ErrUnimplemented
}

// GetCollectionByID returns a collection by ID.
func (c *UnimplementedClient) GetCollectionByID(ctx context.Context, id flow.Identifier) (*flow.Collection, error) {
	return nil, ErrUnimplemented
}

// GetFullCollectionByID returns a full collection including transaction bodies by ID.
func (c *UnimplementedClient) GetFullCollectionByID(ctx context.Context, id flow.Identifier) (*flow.FullCollection, error) {
	return nil, ErrUnimplemented
}

// SendTransaction submits a transaction to the network.
func (c *UnimplementedClient) SendTransaction(ctx context.Context, tx flow.Transaction) error {
	return ErrUnimplemented
}

// GetTransaction returns a transaction by ID.
func (c *UnimplementedClient) GetTransaction(ctx context.Context, txID flow.Identifier) (*flow.Transaction, error) {
	return nil, ErrUnimplemented
}

// GetTransactionsByBlockID returns all the transactions for a specified block.
func (c *UnimplementedClient) GetTransactionsByBlockID(ctx context.Context, blockID flow.Identifier) ([]*flow.Transaction, error) {
	return nil, ErrUnimplemented
}

// GetTransactionResult returns the result of a transaction.
func (c *UnimplementedClient) GetTransactionResult(ctx context.Context, txID flow.Identifier) (*flow.TransactionResult, error) {
	return nil, ErrUnimplemented
}

// GetTransactionResultByIndex returns a transaction result by transaction index for the given block ID.
func (c *UnimplementedClient) GetTransactionResultByIndex(ctx context.Context, blockID flow.Identifier, index uint32) (*flow.TransactionResult, error) {
	return nil, ErrUnimplemented
}

// GetTransactionResultsByBlockID returns all the transaction results for a specified block.
func (c *UnimplementedClient) GetTransactionResultsByBlockID(ctx context.Context, blockID flow.Identifier) ([]*flow.TransactionResult, error) {
	return nil, ErrUnimplemented
}

// GetSystemTransaction returns the system transaction for the given block ID.
func (c *UnimplementedClient) GetSystemTransaction(ctx context.Context, blockID flow.Identifier) (*flow.Transaction, error) {
	return nil, ErrUnimplemented
}

// GetSystemTransactionResult returns the transaction result of the system transaction for the given block ID.
func (c *UnimplementedClient) GetSystemTransactionResult(ctx context.Context, blockID flow.Identifier) (*flow.TransactionResult, error) {
	return nil, ErrUnimplemented
}

// GetAccount is an alias for GetAccountAtLatestBlock.
func (c *UnimplementedClient) GetAccount(ctx context.Context, address flow.Address) (*flow.Account, error) {
	return nil, ErrUnimplemented
}

// GetAccountAtLatestBlock returns an account by address at the latest sealed block.
func (c *UnimplementedClient) GetAccountAtLatestBlock(ctx context.Context, address flow.Address) (*flow.Account, error) {
	return nil, ErrUnimplemented
}

// GetAccountAtBlockHeight returns an account by address at the given block height
func (c *UnimplementedClient) GetAccountAtBlockHeight(ctx context.Context, address flow.Address, blockHeight uint64) (*flow.Account, error) {
	return nil, ErrUnimplemented
}

// GetAccountBalanceAtLatestBlock returns the balance of an account at the latest sealed block.
func (c *UnimplementedClient) GetAccountBalanceAtLatestBlock(ctx context.Context, address flow.Address) (uint64, error) {
	return 0, ErrUnimplemented
}

// GetAccountBalanceAtBlockHeight returns the balance of an account at the given block height.
func (c *UnimplementedClient) GetAccountBalanceAtBlockHeight(ctx context.Context, address flow.Address, blockHeight uint64) (uint64, error) {
	return 0, ErrUnimplemented
}

// GetAccountKeyAtLatestBlock returns the account key with the provided index at the latest sealed block.
func (c *UnimplementedClient) GetAccountKeyAtLatestBlock(ctx context.Context, address flow.Address, keyIndex uint32) (*flow.AccountKey, error) {
	return nil, ErrUnimplemented
}

// GetAccountKeyAtBlockHeight returns the account key with the provided index at the given block height.
func (c *UnimplementedClient) GetAccountKeyAtBlockHeight(ctx context.Context, address flow.Address, keyIndex uint32, height uint64) (*flow.AccountKey, error) {
	return nil, ErrUnimplemented
}

// GetAccountKeysAtLatestBlock returns all account keys at the latest sealed block.
func (c *UnimplementedClient) GetAccountKeysAtLatestBlock(ctx context.Context, address flow.Address) ([]*flow.AccountKey, error) {
	return nil, ErrUnimplemented
}

// GetAccountKeysAtBlockHeight returns all account keys at the given block height.
func (c *UnimplementedClient) GetAccountKeysAtBlockHeight(ctx context.Context, address flow.Address, height uint64) ([]*flow.AccountKey, error) {
	return nil, ErrUnimplemented
}

// ExecuteScriptAtLatestBlock executes a read-only Cadence script against the latest sealed execution state.
func (c *UnimplementedClient) ExecuteScriptAtLatestBlock(ctx context.Context, script []byte, arguments []cadence.Value) (cadence.Value, error) {
	return nil, ErrUnimplemented
}

// ExecuteScriptAtBlockID executes a ready-only Cadence script against the execution state at the block with the given ID.
func (c *UnimplementedClient) ExecuteScriptAtBlockID(ctx context.Context, blockID flow.Identifier, script []byte, arguments []cadence.Value) (cadence.Value, error) {
	return nil, ErrUnimplemented
}

// ExecuteScriptAtBlockHeight executes a ready-only Cadence script against the execution state at the given block height.
func (c *UnimplementedClient) ExecuteScriptAtBlockHeight(ctx context.Context, height uint64, script []byte, arguments []cadence.Value) (cadence.Value, error) {
	return nil, ErrUnimplemented
}

// GetEventsForHeightRange returns events for all sealed blocks between the start and end block heights (inclusive) with the given type.
func (c *UnimplementedClient) GetEventsForHeightRange(ctx context.Context, eventType string, startHeight uint64, endHeight uint64) ([]flow.BlockEvents, error) {
	return nil, ErrUnimplemented
}

// GetEventsForBlockIDs returns events with the given type from the specified block IDs.
func (c *UnimplementedClient) GetEventsForBlockIDs(ctx context.Context, eventType string, blockIDs []flow.Identifier) ([]flow.BlockEvents, error) {
	return nil, ErrUnimplemented
}

// GetLatestProtocolStateSnapshot returns the protocol state snapshot in serialized form at latest sealed block.
// This is used to generate a root snapshot file used by Flow nodes to bootstrap their local protocol state database.
func (c *UnimplementedClient) GetLatestProtocolStateSnapshot(ctx context.Context) ([]byte, error) {
	return nil, ErrUnimplemented
}

// GetProtocolStateSnapshotByBlockID returns the protocol state snapshot in serialized form at the given block ID.
// This is used to generate a root snapshot file used by Flow nodes to bootstrap their local protocol state database.
func (c *UnimplementedClient) GetProtocolStateSnapshotByBlockID(ctx context.Context, blockID flow.Identifier) ([]byte, error) {
	return nil, ErrUnimplemented
}

// GetProtocolStateSnapshotByHeight returns the protocol state snapshot in serialized form at the given block height.
// This is used to generate a root snapshot file used by Flow nodes to bootstrap their local protocol state database.
func (c *UnimplementedClient) GetProtocolStateSnapshotByHeight(ctx context.Context, blockHeight uint64) ([]byte, error) {
	return nil, ErrUnimplemented
}

// GetExecutionResultByID returns the execution result by ID.
func (c *UnimplementedClient) GetExecutionResultByID(ctx context.Context, id flow.Identifier) (*flow.ExecutionResult, error) {
	return nil, ErrUnimplemented
}

// GetExecutionResultForBlockID returns the execution results at the block ID.
func (c *UnimplementedClient) GetExecutionResultForBlockID(ctx context.Context, blockID flow.Identifier) (*flow.ExecutionResult, error) {
	return nil, ErrUnimplemented
}

// GetExecutionDataByBlockID returns execution data for a specific block ID.
func (c *UnimplementedClient) GetExecutionDataByBlockID(ctx context.Context, blockID flow.Identifier) (*flow.ExecutionData, error) {
	return nil, ErrUnimplemented
}

// SubscribeExecutionDataByBlockID subscribes to execution data updates starting at the given block ID.
func (c *UnimplementedClient) SubscribeExecutionDataByBlockID(
	ctx context.Context,
	startBlockID flow.Identifier,
) (<-chan *flow.ExecutionDataStreamResponse, <-chan error, error) {
	return nil, nil, ErrUnimplemented
}

// SubscribeExecutionDataByBlockHeight subscribes to execution data updates starting at the given block height.
func (c *UnimplementedClient) SubscribeExecutionDataByBlockHeight(
	ctx context.Context,
	startHeight uint64,
) (<-chan *flow.ExecutionDataStreamResponse, <-chan error, error) {
	return nil, nil, ErrUnimplemented
}

// SubscribeEventsByBlockID subscribes to events starting at the given block ID.
func (c *UnimplementedClient) SubscribeEventsByBlockID(
	ctx context.Context,
	startBlockID flow.Identifier,
	filter flow.EventFilter,
	opts ...SubscribeOption,
) (<-chan flow.BlockEvents, <-chan error, error) {
	return nil, nil, ErrUnimplemented
}

// SubscribeEventsByBlockHeight subscribes to events starting at the given block height.
func (c *UnimplementedClient) SubscribeEventsByBlockHeight(
	ctx context.Context,
	startHeight uint64,
	filter flow.EventFilter,
	opts ...SubscribeOption,
) (<-chan flow.BlockEvents, <-chan error, error) {
	return nil, nil, ErrUnimplemented
}

// SubscribeBlockDigestsFromStartBlockID subscribes to block digests with the given status
// starting at the given block ID.
// The status may be either flow.BlockStatusFinalized or flow.BlockStatusSealed
func (c *UnimplementedClient) SubscribeBlockDigestsFromStartBlockID(
	ctx context.Context,
	startBlockID flow.Identifier,
	blockStatus flow.BlockStatus,
) (<-chan *flow.BlockDigest, <-chan error, error) {
	return nil, nil, ErrUnimplemented
}

// SubscribeBlockDigestsFromStartHeight subscribes to block digests with the given status
// starting at the given block height.
// The status may be either flow.BlockStatusFinalized or flow.BlockStatusSealed
func (c *UnimplementedClient) SubscribeBlockDigestsFromStartHeight(
	ctx context.Context,
	startHeight uint64,
	blockStatus flow.BlockStatus,
) (<-chan *flow.BlockDigest, <-chan error, error) {
	return nil, nil, ErrUnimplemented
}

// SubscribeBlockDigestsFromLatest subscribes to block digests with the given status
// starting at the latest block.
// The status may be either flow.BlockStatusFinalized or flow.BlockStatusSealed
func (c *UnimplementedClient) SubscribeBlockDigestsFromLatest(
	ctx context.Context,
	blockStatus flow.BlockStatus,
) (<-chan *flow.BlockDigest, <-chan error, error) {
	return nil, nil, ErrUnimplemented
}

// SubscribeBlocksFromStartBlockID subscribes to blocks with the given status starting at the
// given block ID.
// The status may be either flow.BlockStatusFinalized or flow.BlockStatusSealed
func (c *UnimplementedClient) SubscribeBlocksFromStartBlockID(
	ctx context.Context,
	startBlockID flow.Identifier,
	blockStatus flow.BlockStatus,
) (<-chan *flow.Block, <-chan error, error) {
	return nil, nil, ErrUnimplemented
}

// SubscribeBlocksFromStartHeight subscribes to blocks with the given status starting at the
// given block height.
// The status may be either flow.BlockStatusFinalized or flow.BlockStatusSealed
func (c *UnimplementedClient) SubscribeBlocksFromStartHeight(
	ctx context.Context,
	startHeight uint64,
	blockStatus flow.BlockStatus,
) (<-chan *flow.Block, <-chan error, error) {
	return nil, nil, ErrUnimplemented
}

// SubscribeBlocksFromLatest subscribes to blocks with the given status starting at the latest block.
// The status may be either flow.BlockStatusFinalized or flow.BlockStatusSealed
func (c *UnimplementedClient) SubscribeBlocksFromLatest(
	ctx context.Context,
	blockStatus flow.BlockStatus,
) (<-chan *flow.Block, <-chan error, error) {
	return nil, nil, ErrUnimplemented
}

// SubscribeBlockHeadersFromStartBlockID subscribes to block headers with the given status starting
// at the given block ID.
// The status may be either flow.BlockStatusFinalized or flow.BlockStatusSealed
func (c *UnimplementedClient) SubscribeBlockHeadersFromStartBlockID(
	ctx context.Context,
	startBlockID flow.Identifier,
	blockStatus flow.BlockStatus,
) (<-chan *flow.BlockHeader, <-chan error, error) {
	return nil, nil, ErrUnimplemented
}

// SubscribeBlockHeadersFromStartHeight subscribes to block headers with the given status starting
// at the given block height.
// The status may be either flow.BlockStatusFinalized or flow.BlockStatusSealed
func (c *UnimplementedClient) SubscribeBlockHeadersFromStartHeight(
	ctx context.Context,
	startHeight uint64,
	blockStatus flow.BlockStatus,
) (<-chan *flow.BlockHeader, <-chan error, error) {
	return nil, nil, ErrUnimplemented
}

// SubscribeBlockHeadersFromLatest subscribes to block headers with the given status starting
// at the latest block.
// The status may be either flow.BlockStatusFinalized or flow.BlockStatusSealed
func (c *UnimplementedClient) SubscribeBlockHeadersFromLatest(
	ctx context.Context,
	blockStatus flow.BlockStatus,
) (<-chan *flow.BlockHeader, <-chan error, error) {
	return nil, nil, ErrUnimplemented
}

// SubscribeAccountStatusesFromStartHeight subscribes to account status events starting at the
// given block height.
func (c *UnimplementedClient) SubscribeAccountStatusesFromStartHeight(
	ctx context.Context,
	startBlockHeight uint64,
	filter flow.AccountStatusFilter,
) (<-chan *flow.AccountStatus, <-chan error, error) {
	return nil, nil, ErrUnimplemented
}

// SubscribeAccountStatusesFromStartBlockID subscribes to account status events starting at the
// given block ID.
func (c *UnimplementedClient) SubscribeAccountStatusesFromStartBlockID(
	ctx context.Context,
	startBlockID flow.Identifier,
	filter flow.AccountStatusFilter,
) (<-chan *flow.AccountStatus, <-chan error, error) {
	return nil, nil, ErrUnimplemented
}

// SubscribeAccountStatusesFromLatestBlock subscribes to account status events starting at the
// latest block.
func (c *UnimplementedClient) SubscribeAccountStatusesFromLatestBlock(
	ctx context.Context,
	filter flow.AccountStatusFilter,
) (<-chan *flow.AccountStatus, <-chan error, error) {
	return nil, nil, ErrUnimplemented
}

// SendAndSubscribeTransactionStatuses submits a transaction to the network and subscribes to the
// transaction status updates.
func (c *UnimplementedClient) SendAndSubscribeTransactionStatuses(
	ctx context.Context,
	tx flow.Transaction,
) (<-chan *flow.TransactionResult, <-chan error, error) {
	return nil, nil, ErrUnimplemented
}

// Close stops the client connection to the access node.
func (c *UnimplementedClient) Close() error {
	return ErrUnimplemented
}

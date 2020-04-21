/*
 * Flow Go SDK
 *
 * Copyright 2019-2020 Dapper Labs, Inc.
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

package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/onflow/cadence"
	encoding "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow/protobuf/go/flow/access"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client/convert"
)

// An RPCClient is an RPC client for the Flow Access API.
type RPCClient interface {
	access.AccessAPIClient
}

// A Client is a gRPC Client for the Flow Access API.
type Client struct {
	rpcClient RPCClient
	close     func() error
}

// New initializes a Flow client with the default gRPC provider.
//
// An error will be returned if the host is unreachable.
func New(addr string, opts ...grpc.DialOption) (*Client, error) {
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}

	grpcClient := access.NewAccessAPIClient(conn)

	return &Client{
		rpcClient: grpcClient,
		close:     func() error { return conn.Close() },
	}, nil
}

// NewFromRPCClient initializes a Flow client using a pre-configured gRPC provider.
func NewFromRPCClient(rpcClient RPCClient) *Client {
	return &Client{
		rpcClient: rpcClient,
		close:     func() error { return nil },
	}
}

// Close closes the client connection.
func (c *Client) Close() error {
	return c.close()
}

// Ping is used to check if the access node is alive and healthy.
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.rpcClient.Ping(ctx, &access.PingRequest{})
	return err
}

// GetLatestBlockHeader gets the latest sealed or unsealed block header.
func (c *Client) GetLatestBlockHeader(
	ctx context.Context,
	isSealed bool,
) (*flow.BlockHeader, error) {
	req := &access.GetLatestBlockHeaderRequest{
		IsSealed: isSealed,
	}

	res, err := c.rpcClient.GetLatestBlockHeader(ctx, req)
	if err != nil {
		// TODO: improve errors
		return nil, fmt.Errorf("client: %w", err)
	}

	return getBlockHeaderResult(res)
}

// GetBlockHeaderByID gets a block header by ID.
func (c *Client) GetBlockHeaderByID(ctx context.Context, blockID flow.Identifier) (*flow.BlockHeader, error) {
	req := &access.GetBlockHeaderByIDRequest{
		Id: blockID.Bytes(),
	}

	res, err := c.rpcClient.GetBlockHeaderByID(ctx, req)
	if err != nil {
		// TODO: improve errors
		return nil, fmt.Errorf("client: %w", err)
	}

	return getBlockHeaderResult(res)
}

// GetBlockHeaderByHeight gets a block header by height.
func (c *Client) GetBlockHeaderByHeight(ctx context.Context, height uint64) (*flow.BlockHeader, error) {
	req := &access.GetBlockHeaderByHeightRequest{
		Height: height,
	}

	res, err := c.rpcClient.GetBlockHeaderByHeight(ctx, req)
	if err != nil {
		// TODO: improve errors
		return nil, fmt.Errorf("client: %w", err)
	}

	return getBlockHeaderResult(res)
}

func getBlockHeaderResult(res *access.BlockHeaderResponse) (*flow.BlockHeader, error) {
	result, err := convert.MessageToBlockHeader(res.GetBlock())
	if err != nil {
		// TODO: improve errors
		return nil, fmt.Errorf("client: %w", err)
	}

	return &result, nil
}

// GetLatestBlock gets the full payload of the latest sealed or unsealed block.
func (c *Client) GetLatestBlock(ctx context.Context, isSealed bool) (*flow.Block, error) {
	req := &access.GetLatestBlockRequest{
		IsSealed: isSealed,
	}

	res, err := c.rpcClient.GetLatestBlock(ctx, req)
	if err != nil {
		// TODO: improve errors
		return nil, fmt.Errorf("client: %w", err)
	}

	return getBlockResult(res)
}

// GetBlockByID gets a full block by ID.
func (c *Client) GetBlockByID(ctx context.Context, blockID flow.Identifier) (*flow.Block, error) {
	req := &access.GetBlockByIDRequest{
		Id: blockID.Bytes(),
	}

	res, err := c.rpcClient.GetBlockByID(ctx, req)
	if err != nil {
		// TODO: improve errors
		return nil, fmt.Errorf("client: %w", err)
	}

	return getBlockResult(res)
}

// GetBlockByHeight gets a full block by height.
func (c *Client) GetBlockByHeight(ctx context.Context, height uint64) (*flow.Block, error) {
	req := &access.GetBlockByHeightRequest{
		Height: height,
	}

	res, err := c.rpcClient.GetBlockByHeight(ctx, req)
	if err != nil {
		// TODO: improve errors
		return nil, fmt.Errorf("client: %w", err)
	}

	return getBlockResult(res)
}

func getBlockResult(res *access.BlockResponse) (*flow.Block, error) {
	result, err := convert.MessageToBlock(res.GetBlock())
	if err != nil {
		// TODO: improve errors
		return nil, fmt.Errorf("client: %w", err)
	}

	return &result, nil
}

// GetCollection gets a collection by ID.
func (c *Client) GetCollection(ctx context.Context, colID flow.Identifier) (*flow.Collection, error) {
	req := &access.GetCollectionByIDRequest{
		Id: colID.Bytes(),
	}

	res, err := c.rpcClient.GetCollectionByID(ctx, req)
	if err != nil {
		// TODO: improve errors
		return nil, fmt.Errorf("client: %w", err)
	}

	result, err := convert.MessageToCollection(res.GetCollection())
	if err != nil {
		// TODO: improve errors
		return nil, fmt.Errorf("client: %w", err)
	}

	return &result, nil
}

// SendTransaction submits a transaction to the network.
func (c *Client) SendTransaction(ctx context.Context, transaction flow.Transaction) error {
	req := &access.SendTransactionRequest{
		Transaction: convert.TransactionToMessage(transaction),
	}

	_, err := c.rpcClient.SendTransaction(ctx, req)
	if err != nil {
		// TODO: improve errors
		return fmt.Errorf("client: %w", err)
	}

	return nil
}

// GetTransaction gets a transaction by ID.
func (c *Client) GetTransaction(ctx context.Context, txID flow.Identifier) (*flow.Transaction, error) {
	req := &access.GetTransactionRequest{
		Id: txID.Bytes(),
	}

	res, err := c.rpcClient.GetTransaction(ctx, req)
	if err != nil {
		// TODO: improve errors
		return nil, fmt.Errorf("client: %w", err)
	}

	result, err := convert.MessageToTransaction(res.GetTransaction())
	if err != nil {
		// TODO: improve errors
		return nil, fmt.Errorf("client: %w", err)
	}

	return &result, nil
}

// GetTransactionResult gets the result of a transaction.
func (c *Client) GetTransactionResult(ctx context.Context, txID flow.Identifier) (*flow.TransactionResult, error) {
	req := &access.GetTransactionRequest{
		Id: txID.Bytes(),
	}

	res, err := c.rpcClient.GetTransactionResult(ctx, req)
	if err != nil {
		// TODO: improve errors
		return nil, fmt.Errorf("client: %w", err)
	}

	result, err := convert.MessageToTransactionResult(res)
	if err != nil {
		// TODO: improve errors
		return nil, fmt.Errorf("client: %w", err)
	}

	return &result, nil
}

// GetAccount gets an account by address.
func (c *Client) GetAccount(ctx context.Context, address flow.Address) (*flow.Account, error) {
	res, err := c.rpcClient.GetAccount(
		ctx,
		&access.GetAccountRequest{Address: address.Bytes()},
	)
	if err != nil {
		return nil, err
	}

	account, err := convert.MessageToAccount(res.GetAccount())
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// ExecuteScriptAtLatestBlock executes a read-only Cadence script against the latest sealed execution state.
func (c *Client) ExecuteScriptAtLatestBlock(ctx context.Context, script []byte) (cadence.Value, error) {
	res, err := c.rpcClient.ExecuteScriptAtLatestBlock(ctx, &access.ExecuteScriptAtLatestBlockRequest{Script: script})
	if err != nil {
		return nil, err
	}

	value, err := encoding.Decode(res.GetValue())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return value, nil
}

// ExecuteScriptAtBlockID executes a ready-only Cadence script against the execution state at the block with the given ID.
func (c *Client) ExecuteScriptAtBlockID(ctx context.Context) error {
	panic("not implemented")
}

// ExecuteScriptAtBlockHeight executes a ready-only Cadence script against the execution state at the given block height.
func (c *Client) ExecuteScriptAtBlockHeight(ctx context.Context) error {
	panic("not implemented")
}

// EventRangeQuery defines a query for Flow events.
type EventRangeQuery struct {
	// The event type to search for. If empty, no filtering by type is done.
	Type string
	// The block height to begin looking for events
	StartHeight uint64
	// The block height to end looking for events (inclusive)
	EndHeight uint64
}

// BlockEvents are the events that occurred in a specific block.
type BlockEvents struct {
	BlockID flow.Identifier
	Height  uint64
	Events  []flow.Event
}

// GetEventsForHeightRange retrieves events for all sealed blocks between the start and end block
// heights (inclusive) with the given type.
func (c *Client) GetEventsForHeightRange(ctx context.Context, query EventRangeQuery) ([]BlockEvents, error) {
	req := &access.GetEventsForHeightRangeRequest{
		Type:        query.Type,
		StartHeight: query.StartHeight,
		EndHeight:   query.EndHeight,
	}

	res, err := c.rpcClient.GetEventsForHeightRange(ctx, req)
	if err != nil {
		// TODO: improve errors
		return nil, fmt.Errorf("client: %w", err)
	}

	return getEventsResult(res)
}

// GetEventsForBlockIDs retrieves events with the given type from the specified block IDs.
func (c *Client) GetEventsForBlockIDs(
	ctx context.Context,
	eventType string,
	blockIDs []flow.Identifier,
) ([]BlockEvents, error) {
	req := &access.GetEventsForBlockIDsRequest{
		Type:     eventType,
		BlockIds: convert.IDsToMessages(blockIDs),
	}

	res, err := c.rpcClient.GetEventsForBlockIDs(ctx, req)
	if err != nil {
		// TODO: improve errors
		return nil, fmt.Errorf("client: %w", err)
	}

	return getEventsResult(res)
}

func getEventsResult(res *access.EventsResponse) ([]BlockEvents, error) {
	resultMessages := res.GetResults()

	results := make([]BlockEvents, len(resultMessages))
	for i, result := range resultMessages {
		eventMessages := result.GetEvents()

		events := make([]flow.Event, len(eventMessages))

		for i, m := range eventMessages {
			evt, err := convert.MessageToEvent(m)
			if err != nil {
				// TODO: improve errors
				return nil, fmt.Errorf("client: %w", err)
			}

			events[i] = evt
		}

		results[i] = BlockEvents{
			BlockID: flow.HashToID(result.GetBlockId()),
			Height:  result.GetBlockHeight(),
			Events:  events,
		}
	}

	return results, nil
}

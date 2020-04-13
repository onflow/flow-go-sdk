package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dapperlabs/cadence"
	encoding "github.com/dapperlabs/cadence/encoding/json"
	"github.com/dapperlabs/flow/protobuf/go/flow/access"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/client/convert"
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
func (c *Client) GetLatestBlockHeader(ctx context.Context) error {
	panic("not implemented")
}

// GetBlockHeaderByID gets a block header by ID.
func (c *Client) GetBlockHeaderByID(ctx context.Context) error {
	panic("not implemented")
}

// GetBlockHeaderByHeight gets a block header by height.
func (c *Client) GetBlockHeaderByHeight(ctx context.Context) error {
	panic("not implemented")
}

// GetLatestBlock gets the full payload of the latest sealed or unsealed block.
func (c *Client) GetLatestBlock(ctx context.Context) error {
	panic("not implemented")
}

// GetBlockByID gets a full block by ID.
func (c *Client) GetBlockByID(ctx context.Context) error {
	panic("not implemented")
}

// GetBlockByHeight gets a full block by height.
func (c *Client) GetBlockByHeight(ctx context.Context) error {
	panic("not implemented")
}

// GetCollectionByID gets a collection by ID.
func (c *Client) GetCollectionByID(ctx context.Context) error {
	panic("not implemented")
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
func (c *Client) GetTransaction(ctx context.Context, id flow.Identifier) (*flow.Transaction, error) {
	panic("not implemented")
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

// ExecuteScriptAtLatestBlock executes a read-only Cadance script against the latest sealed execution state.
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

type EventRangeResult struct {
	BlockID flow.Identifier
	Height  uint64
	Events  []flow.Event
}

// GetEventsForHeightRange retrieves events for all sealed blocks between the start and end block
// heights (inclusive) with the given type.
func (c *Client) GetEventsForHeightRange(ctx context.Context, query EventRangeQuery) ([]EventRangeResult, error) {
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

	resultMessages := res.GetResults()

	results := make([]EventRangeResult, len(resultMessages))
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

		results[i] = EventRangeResult{
			BlockID: flow.HashToID(result.GetBlockId()),
			Height:  result.GetBlockHeight(),
			Events:  events,
		}
	}

	return results, nil
}

// GetEventsForBlockIDs retrieves events for all the specified block IDs that have the given type
func (c *Client) GetEventsForBlockIDs(ctx context.Context, blockIDs [][]byte) ([]flow.Event, error) {
	panic("not implemented")
}

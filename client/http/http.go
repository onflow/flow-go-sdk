package http

import (
	"context"
	"fmt"
	"strings"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client/convert"
	"github.com/onflow/flow-go/engine/access/rest/models"
)

// handler interface defines methods needed to be offered by a specific http network implementation.
type handler interface {
	getBlockByID(ctx context.Context, ID string, opts ...queryOpts) (*models.Block, error)
	getBlocksByHeights(ctx context.Context, heights string, startHeight string, endHeight string, opts ...queryOpts) ([]*models.Block, error)
	getAccount(ctx context.Context, address string, height string, opts ...queryOpts) (*models.Account, error)
	getCollection(ctx context.Context, ID string, opts ...queryOpts) (*models.Collection, error)
	executeScriptAtBlockHeight(ctx context.Context, height string, script string, arguments []string, opts ...queryOpts) (string, error)
	executeScriptAtBlockID(ctx context.Context, ID string, script string, arguments []string, opts ...queryOpts) (string, error)
	getTransaction(ctx context.Context, ID string, includeResult bool, opts ...queryOpts) (*models.Transaction, error)
	sendTransaction(ctx context.Context, transaction []byte, opts ...queryOpts) error
	getEvents(ctx context.Context, eventType string, start string, end string, blockIDs []string, opts ...queryOpts) ([]models.BlockEvents, error)
}

// ExpandOpts allows you to define a list of fields that you want to retrieve as extra data in the response.
//
// Be sure to follow the documentation for allowed values found here https://docs.onflow.org/http-api/
type ExpandOpts struct {
	Expands []string
}

func (e *ExpandOpts) toQuery() (string, string) {
	return "expands", strings.Join(e.Expands, ",")
}

// SelectOpts allows you to define a list of fields that you only want to fetch in the response filtering out any other data.
//
// Be sure to follow the documentation for allowed values found here https://docs.onflow.org/http-api/
type SelectOpts struct {
	Selects []string
}

func (e *SelectOpts) toQuery() (string, string) {
	return "select", strings.Join(e.Selects, ",")
}

func NewHTTPClient(handler handler) *HTTPClient {
	return &HTTPClient{handler: handler}
}

// HTTPClient exposes methods specific to the http clients exposing all capabilities of the network implementation.
type HTTPClient struct {
	handler handler
}

func (c *HTTPClient) Ping(ctx context.Context) error {
	panic("implement me")
}

func (c *HTTPClient) GetBlockByID(ctx context.Context, blockID flow.Identifier, opts ...queryOpts) (*flow.Block, error) {
	block, err := c.handler.getBlockByID(ctx, blockID.String())
	if err != nil {
		return nil, err
	}

	return convert.HTTPToBlock(block)
}

// SpecialHeight defines two special height values.
type SpecialHeight string

const (
	// FINAL points to latest finalised block height.
	FINAL SpecialHeight = "final"
	// SEALED points to latest sealed block height.
	SEALED = "sealed"
)

// BlockHeightQuery defines all the possible heights you can pass when fetching blocks.
//
// Make sure you only pass either heights or special heights or start and end height else an
// error will be returned. You can refer to the docs for querying blocks found here https://docs.onflow.org/http-api/#tag/Blocks/paths/~1blocks/get
type BlockHeightQuery struct {
	Heights []uint64
	Special SpecialHeight
	Start   uint64
	End     uint64
}

// GetBlocksByHeights
func (c *HTTPClient) GetBlocksByHeights(ctx context.Context, blockQuery BlockHeightQuery, opts ...queryOpts) ([]*flow.Block, error) {
	var heights, start, end string
	if len(blockQuery.Heights) > 0 {
		heights = convert.HeightsToHTTP(blockQuery.Heights)
	} else if blockQuery.Special != "" {
		heights = string(blockQuery.Special)
	} else if blockQuery.Start != 0 && blockQuery.End != 0 {
		start = fmt.Sprintf("%d", blockQuery.Start)
		end = fmt.Sprintf("%d", blockQuery.End)
	} else {
		return nil, fmt.Errorf("must either provide heights or start and end height range")
	}

	httpBlocks, err := c.handler.getBlocksByHeights(ctx, heights, start, end, opts...)
	if err != nil {
		return nil, err
	}

	return convert.HTTPToBlocks(httpBlocks)
}

func (c *HTTPClient) GetCollection(ctx context.Context, ID flow.Identifier) (*flow.Collection, error) {
	collection, err := c.handler.getCollection(ctx, ID.String())
	if err != nil {
		return nil, err
	}

	return convert.HTTPToCollection(collection), nil
}

func (c *HTTPClient) SendTransaction(ctx context.Context, tx flow.Transaction) error {
	convertedTx, err := convert.TransactionToHTTP(tx)
	if err != nil {
		return err
	}

	return c.handler.sendTransaction(ctx, convertedTx)
}

func (c *HTTPClient) GetTransaction(ctx context.Context, ID flow.Identifier) (*flow.Transaction, error) {
	tx, err := c.handler.getTransaction(ctx, ID.String(), false)
	if err != nil {
		return nil, err
	}

	return convert.HTTPToTransaction(tx)
}

func (c *HTTPClient) GetTransactionResult(ctx context.Context, ID flow.Identifier) (*flow.TransactionResult, error) {
	tx, err := c.handler.getTransaction(ctx, ID.String(), true)
	if err != nil {
		return nil, err
	}

	return convert.HTTPToTransactionResult(tx.Result)
}

func (c *HTTPClient) GetAccount(ctx context.Context, address flow.Address) (*flow.Account, error) {
	account, err := c.handler.getAccount(ctx, address.String(), SEALED_HEIGHT)
	if err != nil {
		return nil, err
	}

	return convert.HTTPToAccount(account)
}

func (c *HTTPClient) GetAccountAtLatestBlock(ctx context.Context, address flow.Address) (*flow.Account, error) {
	return c.GetAccount(ctx, address)
}

func (c *HTTPClient) GetAccountAtBlockHeight(
	ctx context.Context,
	address flow.Address,
	blockHeight uint64,
) (*flow.Account, error) {
	account, err := c.handler.getAccount(ctx, address.String(), fmt.Sprintf("%d", blockHeight))
	if err != nil {
		return nil, err
	}

	return convert.HTTPToAccount(account)
}

func (c *HTTPClient) ExecuteScriptAtLatestBlock(
	ctx context.Context,
	script []byte,
	arguments []cadence.Value,
) (cadence.Value, error) {
	args, err := convert.CadenceArgsToHTTP(arguments)
	if err != nil {
		return nil, err
	}

	result, err := c.handler.executeScriptAtBlockHeight(ctx, SEALED_HEIGHT, convert.ScriptToHTTP(script), args)
	if err != nil {
		return nil, err
	}

	return convert.HTTPToCadenceValue(result)
}

func (c *HTTPClient) ExecuteScriptAtBlockID(
	ctx context.Context,
	blockID flow.Identifier,
	script []byte,
	arguments []cadence.Value,
) (cadence.Value, error) {
	args, err := convert.CadenceArgsToHTTP(arguments)
	if err != nil {
		return nil, err
	}

	result, err := c.handler.executeScriptAtBlockID(ctx, blockID.String(), convert.ScriptToHTTP(script), args)
	if err != nil {
		return nil, err
	}

	return convert.HTTPToCadenceValue(result)
}

func (c *HTTPClient) ExecuteScriptAtBlockHeight(
	ctx context.Context,
	height uint64,
	script []byte,
	arguments []cadence.Value,
) (cadence.Value, error) {
	args, err := convert.CadenceArgsToHTTP(arguments)
	if err != nil {
		return nil, err
	}

	result, err := c.handler.executeScriptAtBlockHeight(ctx, fmt.Sprintf("%d", height), convert.ScriptToHTTP(script), args)
	if err != nil {
		return nil, err
	}

	return convert.HTTPToCadenceValue(result)
}

func (c *HTTPClient) GetEventsForHeightRange(
	ctx context.Context,
	eventType string,
	startHeight uint64,
	endHeight uint64,
) ([]flow.BlockEvents, error) {
	events, err := c.handler.getEvents(
		ctx,
		eventType,
		fmt.Sprintf("%d", startHeight),
		fmt.Sprintf("%d", endHeight),
		nil,
	)
	if err != nil {
		return nil, err
	}

	return convert.HTTPToBlockEvents(events)
}

func (c *HTTPClient) GetEventsForBlockIDs(
	ctx context.Context,
	eventType string,
	blockIDs []flow.Identifier,
) ([]flow.BlockEvents, error) {
	ids := make([]string, len(blockIDs))
	for i, id := range blockIDs {
		ids[i] = id.String()
	}

	events, err := c.handler.getEvents(ctx, eventType, "", "", ids)
	if err != nil {
		return nil, err
	}

	return convert.HTTPToBlockEvents(events)
}

func (c *HTTPClient) GetLatestProtocolStateSnapshot(ctx context.Context) ([]byte, error) {
	panic("implement me")
}

func (c *HTTPClient) GetExecutionResultForBlockID(ctx context.Context, blockID flow.Identifier) (*flow.ExecutionResult, error) {
	panic("implement me")
}

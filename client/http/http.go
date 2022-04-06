package http

import (
	"context"
	"fmt"
	"math"
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

// special height values definition.
const (
	// FINAL points to latest finalised block height.
	FINAL uint64 = math.MaxUint64 - 1
	// SEALED points to latest sealed block height.
	SEALED uint64 = math.MaxUint64 - 2
)

var specialHeightMap = map[uint64]string{
	FINAL:  "final",
	SEALED: "sealed",
}

// HeightQuery defines all the possible heights you can pass when fetching blocks.
//
// Make sure you only pass either heights or special heights or start and end height else an
// error will be returned. You can refer to the docs for querying blocks found here https://docs.onflow.org/http-api/#tag/Blocks/paths/~1blocks/get
type HeightQuery struct {
	Heights []uint64
	Start   uint64
	End     uint64
}

// heightToString is a helper method to get first height as string.
func (b *HeightQuery) heightsString() string {
	converted := ""
	for _, h := range b.Heights {
		str := fmt.Sprintf("%d", h)
		if h == FINAL || h == SEALED {
			str = specialHeightMap[h]
		}

		if converted == "" {
			converted = str
		} else {
			converted = fmt.Sprintf("%s,%s", converted, str)
		}
	}
	return converted
}

func (b *HeightQuery) startString() string {
	if b.Start == 0 {
		return ""
	}
	return fmt.Sprintf("%d", b.Start)
}

func (b *HeightQuery) endString() string {
	if b.End == 0 {
		return ""
	}
	return fmt.Sprintf("%d", b.End)
}

func (b *HeightQuery) rangeDefined() bool {
	return b.Start != 0 && b.End != 0
}

func (b *HeightQuery) heightsDefined() bool {
	return len(b.Heights) > 0
}

func (b *HeightQuery) singleHeightDefined() bool {
	return len(b.Heights) == 1
}

func NewHTTPClient(url string) (*HTTPClient, error) {
	handler, err := newHandler(url, false)
	if err != nil {
		return nil, err
	}

	return &HTTPClient{handler}, nil
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

// GetBlocksByHeights requests the blocks by the specified block query.
func (c *HTTPClient) GetBlocksByHeights(
	ctx context.Context,
	heightQuery HeightQuery,
	opts ...queryOpts,
) ([]*flow.Block, error) {

	if !heightQuery.heightsDefined() && !heightQuery.rangeDefined() {
		return nil, fmt.Errorf("must either provide heights or start and end height range")
	}

	httpBlocks, err := c.handler.getBlocksByHeights(
		ctx,
		heightQuery.heightsString(),
		heightQuery.startString(),
		heightQuery.endString(),
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return convert.HTTPToBlocks(httpBlocks)
}

func (c *HTTPClient) GetCollection(
	ctx context.Context,
	ID flow.Identifier,
	opts ...queryOpts,
) (*flow.Collection, error) {
	collection, err := c.handler.getCollection(ctx, ID.String(), opts...)
	if err != nil {
		return nil, err
	}

	return convert.HTTPToCollection(collection), nil
}

func (c *HTTPClient) SendTransaction(
	ctx context.Context,
	tx flow.Transaction,
	opts ...queryOpts,
) error {
	convertedTx, err := convert.TransactionToHTTP(tx)
	if err != nil {
		return err
	}

	return c.handler.sendTransaction(ctx, convertedTx, opts...)
}

func (c *HTTPClient) GetTransaction(
	ctx context.Context,
	ID flow.Identifier,
	opts ...queryOpts,
) (*flow.Transaction, error) {
	tx, err := c.handler.getTransaction(ctx, ID.String(), false, opts...)
	if err != nil {
		return nil, err
	}

	return convert.HTTPToTransaction(tx)
}

func (c *HTTPClient) GetTransactionResult(
	ctx context.Context,
	ID flow.Identifier,
	opts ...queryOpts,
) (*flow.TransactionResult, error) {
	tx, err := c.handler.getTransaction(ctx, ID.String(), true, opts...)
	if err != nil {
		return nil, err
	}

	return convert.HTTPToTransactionResult(tx.Result)
}

func (c *HTTPClient) GetAccountAtBlockHeight(
	ctx context.Context,
	address flow.Address,
	blockQuery HeightQuery,
	opts ...queryOpts,
) (*flow.Account, error) {
	if !blockQuery.singleHeightDefined() {
		return nil, fmt.Errorf("can only provide one block height at a time")
	}

	account, err := c.handler.getAccount(ctx, address.String(), blockQuery.heightsString(), opts...)
	if err != nil {
		return nil, err
	}

	return convert.HTTPToAccount(account)
}

func (c *HTTPClient) ExecuteScriptAtBlockID(
	ctx context.Context,
	blockID flow.Identifier,
	script []byte,
	arguments []cadence.Value,
	opts ...queryOpts,
) (cadence.Value, error) {
	args, err := convert.CadenceArgsToHTTP(arguments)
	if err != nil {
		return nil, err
	}

	result, err := c.handler.executeScriptAtBlockID(ctx, blockID.String(), convert.ScriptToHTTP(script), args, opts...)
	if err != nil {
		return nil, err
	}

	return convert.HTTPToCadenceValue(result)
}

func (c *HTTPClient) ExecuteScriptAtBlockHeight(
	ctx context.Context,
	blockQuery HeightQuery,
	script []byte,
	arguments []cadence.Value,
	opts ...queryOpts,
) (cadence.Value, error) {
	args, err := convert.CadenceArgsToHTTP(arguments)
	if err != nil {
		return nil, err
	}

	if !blockQuery.singleHeightDefined() {
		return nil, fmt.Errorf("must only provide one height at a time")
	}

	result, err := c.handler.executeScriptAtBlockHeight(
		ctx,
		blockQuery.heightsString(),
		convert.ScriptToHTTP(script),
		args,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return convert.HTTPToCadenceValue(result)
}

func (c *HTTPClient) GetEventsForHeightRange(
	ctx context.Context,
	eventType string,
	heightQuery HeightQuery,
) ([]flow.BlockEvents, error) {
	if !heightQuery.rangeDefined() {
		return nil, fmt.Errorf("must provide start and end height range")
	}

	events, err := c.handler.getEvents(
		ctx,
		eventType,
		heightQuery.startString(),
		heightQuery.endString(),
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

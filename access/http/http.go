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

package http

//go:generate mockery --name handler --structname mockHandler --filename=mock_handler.go --inpackage

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/onflow/cadence/encoding/json"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access/http/convert"
	"github.com/onflow/flow-go-sdk/access/http/models"

	"github.com/pkg/errors"

	"github.com/onflow/cadence"
)

// handler interface defines methods needed to be offered by a specific http network implementation.
type handler interface {
	getNetworkParameters(ctx context.Context, opts ...queryOpts) (*models.NetworkParameters, error)
	getNodeVersionInfo(ctx context.Context, opts ...queryOpts) (*models.NodeVersionInfo, error)
	getBlockByID(ctx context.Context, ID string, opts ...queryOpts) (*models.Block, error)
	getBlocksByHeights(ctx context.Context, heights string, startHeight string, endHeight string, opts ...queryOpts) ([]*models.Block, error)
	getAccount(ctx context.Context, address string, height string, opts ...queryOpts) (*models.Account, error)
	getCollection(ctx context.Context, ID string, opts ...queryOpts) (*models.Collection, error)
	executeScriptAtBlockHeight(ctx context.Context, height string, script string, arguments []string, opts ...queryOpts) (string, error)
	executeScriptAtBlockID(ctx context.Context, ID string, script string, arguments []string, opts ...queryOpts) (string, error)
	getTransaction(ctx context.Context, ID string, includeResult bool, opts ...queryOpts) (*models.Transaction, error)
	sendTransaction(ctx context.Context, transaction []byte, opts ...queryOpts) error
	getEvents(ctx context.Context, eventType string, start string, end string, blockIDs []string, opts ...queryOpts) ([]models.BlockEvents, error)
	getExecutionResultByID(ctx context.Context, id string, opts ...queryOpts) (*models.ExecutionResult, error)
	getExecutionResults(ctx context.Context, blockIDs []string, opts ...queryOpts) ([]models.ExecutionResult, error)
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
	if b.Start == 0 && b.End == 0 { // start height can be 0 if end height is not
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
	return b.End != 0 || b.Start != 0
}

func (b *HeightQuery) validateRange() error {
	if b.rangeDefined() && b.Start > b.End {
		return fmt.Errorf("start height (%d) must be smaller than end height (%d)", b.Start, b.End)
	}
	return nil
}

func (b *HeightQuery) heightsDefined() bool {
	return len(b.Heights) > 0
}

func (b *HeightQuery) singleHeightDefined() bool {
	return len(b.Heights) == 1
}

// NewBaseClient creates a new BaseClient. BaseClient provides an API specific to the HTTP.
//
// Use this client if you need advance access to the HTTP API. If you
// don't require special methods use the Client instead.
func NewBaseClient(host string) (*BaseClient, error) {
	handler, err := newHandler(host, false)
	if err != nil {
		return nil, err
	}

	return &BaseClient{
		handler: handler,
		jsonOptions: []json.Option{
			json.WithAllowUnstructuredStaticTypes(true),
		},
	}, nil
}

// BaseClient provides an API specific to the HTTP.
//
// Use this client if you need advance access to the HTTP API. If you
// don't require special methods use the Client instead.
type BaseClient struct {
	handler     handler
	jsonOptions []json.Option
}

func (c *BaseClient) SetJSONOptions(options []json.Option) {
	c.jsonOptions = options
}

func (c *BaseClient) Ping(ctx context.Context) error {
	_, err := c.handler.getBlocksByHeights(ctx, specialHeightMap[SEALED], "", "")
	if err != nil {
		return errors.Wrap(err, "ping error")
	}

	return nil
}

func (c *BaseClient) GetNetworkParameters(ctx context.Context) (*flow.NetworkParameters, error) {
	params, err := c.handler.getNetworkParameters(ctx)
	if err != nil {
		return nil, err
	}

	return convert.ToNetworkParameters(params), nil
}

func (c *BaseClient) GetNodeVersionInfo(ctx context.Context) (*flow.NodeVersionInfo, error) {
	info, err := c.handler.getNodeVersionInfo(ctx)
	if err != nil {
		return nil, err
	}

	return convert.ToNodeVersionInfo(info)
}

func (c *BaseClient) GetBlockByID(ctx context.Context, blockID flow.Identifier, opts ...queryOpts) (*flow.Block, error) {
	block, err := c.handler.getBlockByID(ctx, blockID.String())
	if err != nil {
		return nil, err
	}

	return convert.ToBlock(block)
}

// GetBlocksByHeights requests the blocks by the specified block query.
func (c *BaseClient) GetBlocksByHeights(
	ctx context.Context,
	heightQuery HeightQuery,
	opts ...queryOpts,
) ([]*flow.Block, error) {

	if !heightQuery.heightsDefined() && !heightQuery.rangeDefined() {
		return nil, fmt.Errorf("must either provide heights or start and end height range")
	}

	err := heightQuery.validateRange()
	if err != nil {
		return nil, err
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

	return convert.ToBlocks(httpBlocks)
}

func (c *BaseClient) GetCollection(
	ctx context.Context,
	ID flow.Identifier,
	opts ...queryOpts,
) (*flow.Collection, error) {
	collection, err := c.handler.getCollection(ctx, ID.String(), opts...)
	if err != nil {
		return nil, err
	}

	return convert.ToCollection(collection), nil
}

func (c *BaseClient) SendTransaction(
	ctx context.Context,
	tx flow.Transaction,
	opts ...queryOpts,
) error {
	convertedTx, err := convert.TncodeTransaction(tx)
	if err != nil {
		return err
	}

	return c.handler.sendTransaction(ctx, convertedTx, opts...)
}

func (c *BaseClient) GetTransaction(
	ctx context.Context,
	ID flow.Identifier,
	opts ...queryOpts,
) (*flow.Transaction, error) {
	tx, err := c.handler.getTransaction(ctx, ID.String(), false, opts...)
	if err != nil {
		return nil, err
	}

	return convert.ToTransaction(tx)
}

func (c *BaseClient) GetTransactionResult(
	ctx context.Context,
	ID flow.Identifier,
	opts ...queryOpts,
) (*flow.TransactionResult, error) {
	tx, err := c.handler.getTransaction(ctx, ID.String(), true, opts...)
	if err != nil {
		return nil, err
	}

	return convert.ToTransactionResult(tx.Result, c.jsonOptions)
}

func (c *BaseClient) GetAccountAtBlockHeight(
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

	return convert.ToAccount(account)
}

func (c *BaseClient) ExecuteScriptAtBlockID(
	ctx context.Context,
	blockID flow.Identifier,
	script []byte,
	arguments []cadence.Value,
	opts ...queryOpts,
) (cadence.Value, error) {
	args, err := convert.EncodeCadenceArgs(arguments)
	if err != nil {
		return nil, err
	}

	result, err := c.handler.executeScriptAtBlockID(
		ctx,
		blockID.String(),
		convert.EncodeScript(script),
		args,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return convert.DecodeCadenceValue(result, c.jsonOptions)
}

func (c *BaseClient) ExecuteScriptAtBlockHeight(
	ctx context.Context,
	blockQuery HeightQuery,
	script []byte,
	arguments []cadence.Value,
	opts ...queryOpts,
) (cadence.Value, error) {
	args, err := convert.EncodeCadenceArgs(arguments)
	if err != nil {
		return nil, err
	}

	if !blockQuery.singleHeightDefined() {
		return nil, fmt.Errorf("must only provide one height at a time")
	}

	result, err := c.handler.executeScriptAtBlockHeight(
		ctx,
		blockQuery.heightsString(),
		convert.EncodeScript(script),
		args,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return convert.DecodeCadenceValue(result, c.jsonOptions)
}

func (c *BaseClient) GetEventsForHeightRange(
	ctx context.Context,
	eventType string,
	heightQuery HeightQuery,
) ([]flow.BlockEvents, error) {
	if !heightQuery.rangeDefined() {
		return nil, fmt.Errorf("must provide start and end height range")
	}

	err := heightQuery.validateRange()
	if err != nil {
		return nil, err
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

	return convert.ToBlockEvents(events, c.jsonOptions)
}

func (c *BaseClient) GetEventsForBlockIDs(
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

	return convert.ToBlockEvents(events, c.jsonOptions)
}

func (c *BaseClient) GetLatestProtocolStateSnapshot(ctx context.Context) ([]byte, error) {
	return nil, fmt.Errorf("get latest protocol snapshot is currently not supported for HTTP API, if you require this functionality please open an issue on the flow-go-sdk github")
}

func (c *BaseClient) GetExecutionResultForBlockID(ctx context.Context, blockID flow.Identifier) (*flow.ExecutionResult, error) {
	results, err := c.handler.getExecutionResults(ctx, []string{blockID.String()})
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("results not found") // sanity check
	}

	return convert.ToExecutionResults(results[0]), nil
}

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

package client_test

import (
	"context"
	"errors"
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow/protobuf/go/flow/access"
	"github.com/onflow/flow/protobuf/go/flow/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/client/convert"
	"github.com/onflow/flow-go-sdk/client/mocks"
	"github.com/onflow/flow-go-sdk/test"
)

var mockRPCError = errors.New("rpc error")

func TestClient_Ping(t *testing.T) {
	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		response := &access.PingResponse{}

		rpc.On("Ping", ctx, mock.Anything).Return(response, nil)

		err := c.Ping(ctx)
		assert.NoError(t, err)
	}))

	t.Run("Error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		rpc.On("Ping", ctx, mock.Anything).
			Return(nil, mockRPCError)

		err := c.Ping(ctx)
		assert.Error(t, err)
	}))
}

func TestClient_GetLatestBlockHeader(t *testing.T) {
	blocks := test.BlockGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		expectedHeader := blocks.New().BlockHeader

		response := &access.BlockHeaderResponse{
			Block: convert.BlockHeaderToMessage(expectedHeader),
		}

		rpc.On("GetLatestBlockHeader", ctx, mock.Anything).Return(response, nil)

		header, err := c.GetLatestBlockHeader(ctx, true)
		require.NoError(t, err)

		assert.Equal(t, expectedHeader, *header)
	}))

	t.Run("Error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		rpc.On("GetLatestBlockHeader", ctx, mock.Anything).
			Return(nil, mockRPCError)

		result, err := c.GetLatestBlockHeader(ctx, true)
		assert.Error(t, err)
		assert.Nil(t, result)
	}))
}

func TestClient_GetBlockHeaderByID(t *testing.T) {
	blocks := test.BlockGenerator()
	ids := test.IdentifierGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		blockID := ids.New()
		expectedHeader := blocks.New().BlockHeader
		response := &access.BlockHeaderResponse{
			Block: convert.BlockHeaderToMessage(expectedHeader),
		}

		rpc.On("GetBlockHeaderByID", ctx, mock.Anything).Return(response, nil)

		header, err := c.GetBlockHeaderByID(ctx, blockID)
		require.NoError(t, err)

		assert.Equal(t, expectedHeader, *header)
	}))

	t.Run("Error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		blockID := ids.New()

		rpc.On("GetBlockHeaderByID", ctx, mock.Anything).
			Return(nil, mockRPCError)

		result, err := c.GetBlockHeaderByID(ctx, blockID)
		assert.Error(t, err)
		assert.Nil(t, result)
	}))
}

func TestClient_GetBlockHeaderByHeight(t *testing.T) {
	blocks := test.BlockGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		expectedHeader := blocks.New().BlockHeader
		response := &access.BlockHeaderResponse{
			Block: convert.BlockHeaderToMessage(expectedHeader),
		}

		rpc.On("GetBlockHeaderByHeight", ctx, mock.Anything).Return(response, nil)

		header, err := c.GetBlockHeaderByHeight(ctx, 42)
		require.NoError(t, err)

		assert.Equal(t, expectedHeader, *header)
	}))

	t.Run("Error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		rpc.On("GetBlockHeaderByHeight", ctx, mock.Anything).
			Return(nil, mockRPCError)

		result, err := c.GetBlockHeaderByHeight(ctx, 42)
		assert.Error(t, err)
		assert.Nil(t, result)
	}))
}

func TestClient_GetLatestBlock(t *testing.T) {
	blocks := test.BlockGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		expectedBlock := blocks.New()
		response := &access.BlockResponse{
			Block: convert.BlockToMessage(*expectedBlock),
		}

		rpc.On("GetLatestBlock", ctx, mock.Anything).Return(response, nil)

		block, err := c.GetLatestBlock(ctx, true)
		require.NoError(t, err)

		assert.Equal(t, expectedBlock, block)
	}))

	t.Run("Error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		rpc.On("GetLatestBlock", ctx, mock.Anything).
			Return(nil, mockRPCError)

		result, err := c.GetLatestBlock(ctx, true)
		assert.Error(t, err)
		assert.Nil(t, result)
	}))
}

func TestClient_GetBlockByID(t *testing.T) {
	blocks := test.BlockGenerator()
	ids := test.IdentifierGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		blockID := ids.New()
		expectedBlock := blocks.New()
		response := &access.BlockResponse{
			Block: convert.BlockToMessage(*expectedBlock),
		}

		rpc.On("GetBlockByID", ctx, mock.Anything).Return(response, nil)

		block, err := c.GetBlockByID(ctx, blockID)
		require.NoError(t, err)

		assert.Equal(t, expectedBlock, block)
	}))

	t.Run("Error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		blockID := ids.New()

		rpc.On("GetBlockByID", ctx, mock.Anything).
			Return(nil, mockRPCError)

		result, err := c.GetBlockByID(ctx, blockID)
		assert.Error(t, err)
		assert.Nil(t, result)
	}))
}

func TestClient_GetBlockByHeight(t *testing.T) {
	blocks := test.BlockGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		expectedBlock := blocks.New()
		response := &access.BlockResponse{
			Block: convert.BlockToMessage(*expectedBlock),
		}

		rpc.On("GetBlockByHeight", ctx, mock.Anything).Return(response, nil)

		block, err := c.GetBlockByHeight(ctx, 42)
		require.NoError(t, err)

		assert.Equal(t, expectedBlock, block)
	}))

	t.Run("Error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		rpc.On("GetBlockByHeight", ctx, mock.Anything).
			Return(nil, mockRPCError)

		result, err := c.GetBlockByHeight(ctx, 42)
		assert.Error(t, err)
		assert.Nil(t, result)
	}))
}

func TestClient_GetCollection(t *testing.T) {
	cols := test.CollectionGenerator()
	ids := test.IdentifierGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		colID := ids.New()
		expectedCol := cols.New()
		response := &access.CollectionResponse{
			Collection: convert.CollectionToMessage(*expectedCol),
		}

		rpc.On("GetCollectionByID", ctx, mock.Anything).Return(response, nil)

		col, err := c.GetCollection(ctx, colID)
		require.NoError(t, err)

		assert.Equal(t, expectedCol, col)
	}))

	t.Run("Error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		colID := ids.New()

		rpc.On("GetCollectionByID", ctx, mock.Anything).
			Return(nil, mockRPCError)

		result, err := c.GetCollection(ctx, colID)
		assert.Error(t, err)
		assert.Nil(t, result)
	}))
}

func TestClient_SendTransaction(t *testing.T) {
	transactions := test.TransactionGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		tx := transactions.New()

		response := &access.SendTransactionResponse{
			Id: tx.ID().Bytes(),
		}

		rpc.On("SendTransaction", ctx, mock.Anything).Return(response, nil)

		err := c.SendTransaction(ctx, *tx)
		require.NoError(t, err)
	}))

	t.Run("Error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		tx := transactions.New()

		rpc.On("SendTransaction", ctx, mock.Anything).
			Return(nil, mockRPCError)

		err := c.SendTransaction(ctx, *tx)
		assert.Error(t, err)
	}))
}

func TestClient_GetTransaction(t *testing.T) {
	txs := test.TransactionGenerator()
	ids := test.IdentifierGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		txID := ids.New()
		expectedTx := txs.New()
		response := &access.TransactionResponse{
			Transaction: convert.TransactionToMessage(*expectedTx),
		}

		rpc.On("GetTransaction", ctx, mock.Anything).Return(response, nil)

		tx, err := c.GetTransaction(ctx, txID)
		require.NoError(t, err)

		assert.Equal(t, expectedTx, tx)
	}))

	t.Run("Error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		txID := ids.New()

		rpc.On("GetTransaction", ctx, mock.Anything).
			Return(nil, mockRPCError)

		result, err := c.GetTransaction(ctx, txID)
		assert.Error(t, err)
		assert.Nil(t, result)
	}))
}

func TestClient_GetTransactionResult(t *testing.T) {
	results := test.TransactionResultGenerator()
	ids := test.IdentifierGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		txID := ids.New()
		expectedResult := results.New()
		response, _ := convert.TransactionResultToMessage(expectedResult)

		rpc.On("GetTransactionResult", ctx, mock.Anything).Return(response, nil)

		result, err := c.GetTransactionResult(ctx, txID)
		require.NoError(t, err)

		assert.Equal(t, expectedResult, *result)

	}))

	t.Run("Error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		txID := ids.New()

		rpc.On("GetTransactionResult", ctx, mock.Anything).
			Return(nil, mockRPCError)

		result, err := c.GetTransactionResult(ctx, txID)
		assert.Error(t, err)
		assert.Nil(t, result)
	}))
}

func TestClient_GetAccount(t *testing.T) {
	accounts := test.AccountGenerator()
	addresses := test.AddressGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		expectedAccount := accounts.New()
		response := &access.GetAccountResponse{
			Account: convert.AccountToMessage(*expectedAccount),
		}

		rpc.On("GetAccount", ctx, mock.Anything).Return(response, nil)

		account, err := c.GetAccount(ctx, expectedAccount.Address)
		require.NoError(t, err)

		assert.Equal(t, expectedAccount, account)
	}))

	t.Run("Error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		address := addresses.New()

		rpc.On("GetAccount", ctx, mock.Anything).
			Return(nil, mockRPCError)

		account, err := c.GetAccount(ctx, address)
		assert.Error(t, err)
		assert.Nil(t, account)
	}))
}

func TestClient_ExecuteScriptAtLatestBlock(t *testing.T) {
	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		expectedValue := cadence.NewInt(42)
		encodedValue, err := jsoncdc.Encode(expectedValue)
		require.NoError(t, err)

		response := &access.ExecuteScriptResponse{
			Value: encodedValue,
		}

		rpc.On("ExecuteScriptAtLatestBlock", ctx, mock.Anything).Return(response, nil)

		value, err := c.ExecuteScriptAtLatestBlock(ctx, []byte("foo"))
		require.NoError(t, err)

		assert.Equal(t, expectedValue, value)
	}))

	t.Run(
		"Invalid JSON-CDC",
		clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
			response := &access.ExecuteScriptResponse{
				Value: []byte("invalid JSON-CDC bytes"),
			}

			rpc.On("ExecuteScriptAtLatestBlock", ctx, mock.Anything).Return(response, nil)

			value, err := c.ExecuteScriptAtLatestBlock(ctx, []byte("foo"))
			assert.Error(t, err)
			assert.Nil(t, value)
		}),
	)

	t.Run("Error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		rpc.On("ExecuteScriptAtLatestBlock", ctx, mock.Anything).
			Return(nil, mockRPCError)

		value, err := c.ExecuteScriptAtLatestBlock(ctx, []byte("foo"))
		assert.Error(t, err)
		assert.Nil(t, value)
	}))
}

func TestClient_ExecuteScriptAtBlockID(t *testing.T) {
	ids := test.IdentifierGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		expectedValue := cadence.NewInt(42)
		encodedValue, err := jsoncdc.Encode(expectedValue)
		require.NoError(t, err)

		response := &access.ExecuteScriptResponse{
			Value: encodedValue,
		}

		rpc.On("ExecuteScriptAtBlockID", ctx, mock.Anything).Return(response, nil)

		value, err := c.ExecuteScriptAtBlockID(ctx, ids.New(), []byte("foo"))
		require.NoError(t, err)

		assert.Equal(t, expectedValue, value)
	}))

	t.Run(
		"Invalid JSON-CDC",
		clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
			response := &access.ExecuteScriptResponse{
				Value: []byte("invalid JSON-CDC bytes"),
			}

			rpc.On("ExecuteScriptAtBlockID", ctx, mock.Anything).Return(response, nil)

			value, err := c.ExecuteScriptAtBlockID(ctx, ids.New(), []byte("foo"))
			assert.Error(t, err)
			assert.Nil(t, value)
		}),
	)

	t.Run("Error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		rpc.On("ExecuteScriptAtBlockID", ctx, mock.Anything).
			Return(nil, mockRPCError)

		value, err := c.ExecuteScriptAtBlockID(ctx, ids.New(), []byte("foo"))
		assert.Error(t, err)
		assert.Nil(t, value)
	}))
}

func TestClient_ExecuteScriptAtBlockHeight(t *testing.T) {
	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		expectedValue := cadence.NewInt(42)
		encodedValue, err := jsoncdc.Encode(expectedValue)
		require.NoError(t, err)

		response := &access.ExecuteScriptResponse{
			Value: encodedValue,
		}

		rpc.On("ExecuteScriptAtBlockHeight", ctx, mock.Anything).Return(response, nil)

		value, err := c.ExecuteScriptAtBlockHeight(ctx, 42, []byte("foo"))
		require.NoError(t, err)

		assert.Equal(t, expectedValue, value)
	}))

	t.Run(
		"Invalid JSON-CDC",
		clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
			response := &access.ExecuteScriptResponse{
				Value: []byte("invalid JSON-CDC bytes"),
			}

			rpc.On("ExecuteScriptAtBlockHeight", ctx, mock.Anything).Return(response, nil)

			value, err := c.ExecuteScriptAtBlockHeight(ctx, 42, []byte("foo"))
			assert.Error(t, err)
			assert.Nil(t, value)
		}),
	)

	t.Run("Error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		rpc.On("ExecuteScriptAtBlockHeight", ctx, mock.Anything).
			Return(nil, mockRPCError)

		value, err := c.ExecuteScriptAtBlockHeight(ctx, 42, []byte("foo"))
		assert.Error(t, err)
		assert.Nil(t, value)
	}))
}

func TestClient_GetEventsForHeightRange(t *testing.T) {
	ids := test.IdentifierGenerator()
	events := test.EventGenerator()

	t.Run(
		"Empty result",
		clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
			response := &access.EventsResponse{
				Results: []*access.EventsResponse_Result{},
			}

			rpc.On("GetEventsForHeightRange", ctx, mock.Anything).Return(response, nil)

			blocks, err := c.GetEventsForHeightRange(ctx, client.EventRangeQuery{
				Type:        "foo",
				StartHeight: 1,
				EndHeight:   10,
			})
			require.NoError(t, err)

			assert.Empty(t, blocks)
		}),
	)

	t.Run(
		"Non-empty result",
		clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
			eventA, eventB, eventC, eventD := events.New(), events.New(), events.New(), events.New()

			eventAMsg, _ := convert.EventToMessage(eventA)
			eventBMsg, _ := convert.EventToMessage(eventB)
			eventCMsg, _ := convert.EventToMessage(eventC)
			eventDMsg, _ := convert.EventToMessage(eventD)

			response := &access.EventsResponse{
				Results: []*access.EventsResponse_Result{
					{
						BlockId:     ids.New().Bytes(),
						BlockHeight: 1,
						Events: []*entities.Event{
							eventAMsg,
							eventBMsg,
						},
					},
					{
						BlockId:     ids.New().Bytes(),
						BlockHeight: 2,
						Events: []*entities.Event{
							eventCMsg,
							eventDMsg,
						},
					},
				},
			}

			rpc.On("GetEventsForHeightRange", ctx, mock.Anything).Return(response, nil)

			blocks, err := c.GetEventsForHeightRange(ctx, client.EventRangeQuery{
				Type:        "foo",
				StartHeight: 1,
				EndHeight:   10,
			})
			require.NoError(t, err)

			assert.Len(t, blocks, len(response.Results))

			assert.Equal(t, response.Results[0].BlockId, blocks[0].BlockID.Bytes())
			assert.Equal(t, response.Results[0].BlockHeight, blocks[0].Height)

			assert.Equal(t, response.Results[1].BlockId, blocks[1].BlockID.Bytes())
			assert.Equal(t, response.Results[1].BlockHeight, blocks[1].Height)

			assert.Equal(t, eventA, blocks[0].Events[0])
			assert.Equal(t, eventB, blocks[0].Events[1])
			assert.Equal(t, eventC, blocks[1].Events[0])
			assert.Equal(t, eventD, blocks[1].Events[1])
		}),
	)

	t.Run("Error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		rpc.On("GetEventsForHeightRange", ctx, mock.Anything).
			Return(nil, mockRPCError)

		blocks, err := c.GetEventsForHeightRange(ctx, client.EventRangeQuery{
			Type:        "foo",
			StartHeight: 1,
			EndHeight:   10,
		})
		assert.Error(t, err)
		assert.Empty(t, blocks)
	}))
}

func TestClient_GetEventsForBlockIDs(t *testing.T) {
	ids := test.IdentifierGenerator()
	events := test.EventGenerator()

	t.Run(
		"Empty result",
		clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
			blockIDs := []flow.Identifier{ids.New(), ids.New()}

			response := &access.EventsResponse{
				Results: []*access.EventsResponse_Result{},
			}

			rpc.On("GetEventsForBlockIDs", ctx, mock.Anything).Return(response, nil)

			blocks, err := c.GetEventsForBlockIDs(ctx, "foo", blockIDs)
			require.NoError(t, err)

			assert.Empty(t, blocks)
		}),
	)

	t.Run(
		"Non-empty result",
		clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
			blockIDA, blockIDB := ids.New(), ids.New()
			eventA, eventB, eventC, eventD := events.New(), events.New(), events.New(), events.New()

			eventAMsg, _ := convert.EventToMessage(eventA)
			eventBMsg, _ := convert.EventToMessage(eventB)
			eventCMsg, _ := convert.EventToMessage(eventC)
			eventDMsg, _ := convert.EventToMessage(eventD)

			response := &access.EventsResponse{
				Results: []*access.EventsResponse_Result{
					{
						BlockId:     blockIDA.Bytes(),
						BlockHeight: 1,
						Events: []*entities.Event{
							eventAMsg,
							eventBMsg,
						},
					},
					{
						BlockId:     blockIDB.Bytes(),
						BlockHeight: 2,
						Events: []*entities.Event{
							eventCMsg,
							eventDMsg,
						},
					},
				},
			}

			rpc.On("GetEventsForBlockIDs", ctx, mock.Anything).Return(response, nil)

			blocks, err := c.GetEventsForBlockIDs(ctx, "foo", []flow.Identifier{blockIDA, blockIDB})
			require.NoError(t, err)

			assert.Len(t, blocks, len(response.Results))

			assert.Equal(t, response.Results[0].BlockId, blocks[0].BlockID.Bytes())
			assert.Equal(t, response.Results[0].BlockHeight, blocks[0].Height)

			assert.Equal(t, response.Results[1].BlockId, blocks[1].BlockID.Bytes())
			assert.Equal(t, response.Results[1].BlockHeight, blocks[1].Height)

			assert.Equal(t, eventA, blocks[0].Events[0])
			assert.Equal(t, eventB, blocks[0].Events[1])
			assert.Equal(t, eventC, blocks[1].Events[0])
			assert.Equal(t, eventD, blocks[1].Events[1])
		}),
	)

	t.Run("Error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, c *client.Client) {
		blockIDA, blockIDB := ids.New(), ids.New()

		rpc.On("GetEventsForBlockIDs", ctx, mock.Anything).
			Return(nil, mockRPCError)

		blocks, err := c.GetEventsForBlockIDs(ctx, "foo", []flow.Identifier{blockIDA, blockIDB})
		assert.Error(t, err)
		assert.Empty(t, blocks)
	}))
}

func clientTest(
	f func(t *testing.T, ctx context.Context, rpc *mocks.RPCClient, client *client.Client),
) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		rpc := &mocks.RPCClient{}
		c := client.NewFromRPCClient(rpc)
		f(t, ctx, rpc, c)
		rpc.AssertExpectations(t)
	}
}

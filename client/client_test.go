package client_test

import (
	"context"
	"errors"
	"testing"

	"github.com/onflow/flow/protobuf/go/flow/access"
	"github.com/onflow/flow/protobuf/go/flow/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/client/convert"
	"github.com/onflow/flow-go-sdk/client/mocks"
	"github.com/onflow/flow-go-sdk/test"
)

func TestClient_GetLatestBlockHeader(t *testing.T) {
	blocks := test.BlockGenerator()

	t.Run("Success", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		expectedHeader := blocks.New().BlockHeader
		response := &access.BlockHeaderResponse{
			Block: convert.BlockHeaderToMessage(expectedHeader),
		}

		rpc.On("GetLatestBlockHeader", ctx, mock.Anything).Return(response, nil)

		c := client.NewFromRPCClient(rpc)

		header, err := c.GetLatestBlockHeader(ctx, true)
		assert.NoError(t, err)

		assert.Equal(t, expectedHeader, *header)

		rpc.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		rpc.On("GetLatestBlockHeader", ctx, mock.Anything).
			Return(nil, errors.New("rpc error"))

		c := client.NewFromRPCClient(rpc)

		result, err := c.GetLatestBlockHeader(ctx, true)
		assert.Error(t, err)
		assert.Nil(t, result)

		rpc.AssertExpectations(t)
	})
}

func TestClient_GetBlockHeaderByID(t *testing.T) {
	blocks := test.BlockGenerator()
	ids := test.IdentifierGenerator()

	t.Run("Success", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		blockID := ids.New()
		expectedHeader := blocks.New().BlockHeader
		response := &access.BlockHeaderResponse{
			Block: convert.BlockHeaderToMessage(expectedHeader),
		}

		rpc.On("GetBlockHeaderByID", ctx, mock.Anything).Return(response, nil)

		c := client.NewFromRPCClient(rpc)

		header, err := c.GetBlockHeaderByID(ctx, blockID)
		assert.NoError(t, err)

		assert.Equal(t, expectedHeader, *header)

		rpc.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		blockID := ids.New()

		rpc.On("GetBlockHeaderByID", ctx, mock.Anything).
			Return(nil, errors.New("rpc error"))

		c := client.NewFromRPCClient(rpc)

		result, err := c.GetBlockHeaderByID(ctx, blockID)
		assert.Error(t, err)
		assert.Nil(t, result)

		rpc.AssertExpectations(t)
	})
}

func TestClient_GetBlockHeaderByHeight(t *testing.T) {
	blocks := test.BlockGenerator()

	t.Run("Success", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		expectedHeader := blocks.New().BlockHeader
		response := &access.BlockHeaderResponse{
			Block: convert.BlockHeaderToMessage(expectedHeader),
		}

		rpc.On("GetBlockHeaderByHeight", ctx, mock.Anything).Return(response, nil)

		c := client.NewFromRPCClient(rpc)

		header, err := c.GetBlockHeaderByHeight(ctx, 42)
		assert.NoError(t, err)

		assert.Equal(t, expectedHeader, *header)

		rpc.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		rpc.On("GetBlockHeaderByHeight", ctx, mock.Anything).
			Return(nil, errors.New("rpc error"))

		c := client.NewFromRPCClient(rpc)

		result, err := c.GetBlockHeaderByHeight(ctx, 42)
		assert.Error(t, err)
		assert.Nil(t, result)

		rpc.AssertExpectations(t)
	})
}

func TestClient_GetLatestBlock(t *testing.T) {
	blocks := test.BlockGenerator()

	t.Run("Success", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		expectedBlock := blocks.New()
		response := &access.BlockResponse{
			Block: convert.BlockToMessage(*expectedBlock),
		}

		rpc.On("GetLatestBlock", ctx, mock.Anything).Return(response, nil)

		c := client.NewFromRPCClient(rpc)

		block, err := c.GetLatestBlock(ctx, true)
		assert.NoError(t, err)

		assert.Equal(t, expectedBlock, block)

		rpc.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		rpc.On("GetLatestBlock", ctx, mock.Anything).
			Return(nil, errors.New("rpc error"))

		c := client.NewFromRPCClient(rpc)

		result, err := c.GetLatestBlock(ctx, true)
		assert.Error(t, err)
		assert.Nil(t, result)

		rpc.AssertExpectations(t)
	})
}

func TestClient_GetBlockByID(t *testing.T) {
	blocks := test.BlockGenerator()
	ids := test.IdentifierGenerator()

	t.Run("Success", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		blockID := ids.New()
		expectedBlock := blocks.New()
		response := &access.BlockResponse{
			Block: convert.BlockToMessage(*expectedBlock),
		}

		rpc.On("GetBlockByID", ctx, mock.Anything).Return(response, nil)

		c := client.NewFromRPCClient(rpc)

		block, err := c.GetBlockByID(ctx, blockID)
		assert.NoError(t, err)

		assert.Equal(t, expectedBlock, block)

		rpc.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		blockID := ids.New()

		rpc.On("GetBlockByID", ctx, mock.Anything).
			Return(nil, errors.New("rpc error"))

		c := client.NewFromRPCClient(rpc)

		result, err := c.GetBlockByID(ctx, blockID)
		assert.Error(t, err)
		assert.Nil(t, result)

		rpc.AssertExpectations(t)
	})
}

func TestClient_GetBlockByHeight(t *testing.T) {
	blocks := test.BlockGenerator()

	t.Run("Success", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		expectedBlock := blocks.New()
		response := &access.BlockResponse{
			Block: convert.BlockToMessage(*expectedBlock),
		}

		rpc.On("GetBlockByHeight", ctx, mock.Anything).Return(response, nil)

		c := client.NewFromRPCClient(rpc)

		block, err := c.GetBlockByHeight(ctx, 42)
		assert.NoError(t, err)

		assert.Equal(t, expectedBlock, block)

		rpc.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		rpc.On("GetBlockByHeight", ctx, mock.Anything).
			Return(nil, errors.New("rpc error"))

		c := client.NewFromRPCClient(rpc)

		result, err := c.GetBlockByHeight(ctx, 42)
		assert.Error(t, err)
		assert.Nil(t, result)

		rpc.AssertExpectations(t)
	})
}

func TestClient_GetCollection(t *testing.T) {
	cols := test.CollectionGenerator()
	ids := test.IdentifierGenerator()

	t.Run("Success", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		colID := ids.New()
		expectedCol := cols.New()
		response := &access.CollectionResponse{
			Collection: convert.CollectionToMessage(*expectedCol),
		}

		rpc.On("GetCollectionByID", ctx, mock.Anything).Return(response, nil)

		c := client.NewFromRPCClient(rpc)

		col, err := c.GetCollection(ctx, colID)
		assert.NoError(t, err)

		assert.Equal(t, expectedCol, col)

		rpc.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		colID := ids.New()

		rpc.On("GetCollectionByID", ctx, mock.Anything).
			Return(nil, errors.New("rpc error"))

		c := client.NewFromRPCClient(rpc)

		result, err := c.GetCollection(ctx, colID)
		assert.Error(t, err)
		assert.Nil(t, result)

		rpc.AssertExpectations(t)
	})
}

func TestClient_SendTransaction(t *testing.T) {
	transactions := test.TransactionGenerator()

	t.Run("Success", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		tx := transactions.New()

		response := &access.SendTransactionResponse{
			Id: tx.ID().Bytes(),
		}

		rpc.On("SendTransaction", ctx, mock.Anything).Return(response, nil)

		c := client.NewFromRPCClient(rpc)

		err := c.SendTransaction(ctx, *tx)

		assert.NoError(t, err)

		rpc.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		tx := transactions.New()

		rpc.On("SendTransaction", ctx, mock.Anything).
			Return(nil, errors.New("rpc error"))

		c := client.NewFromRPCClient(rpc)

		err := c.SendTransaction(ctx, *tx)
		assert.Error(t, err)

		rpc.AssertExpectations(t)
	})
}

func TestClient_GetTransaction(t *testing.T) {
	txs := test.TransactionGenerator()
	ids := test.IdentifierGenerator()

	t.Run("Success", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		txID := ids.New()
		expectedTx := txs.New()
		response := &access.TransactionResponse{
			Transaction: convert.TransactionToMessage(*expectedTx),
		}

		rpc.On("GetTransaction", ctx, mock.Anything).Return(response, nil)

		c := client.NewFromRPCClient(rpc)

		tx, err := c.GetTransaction(ctx, txID)
		assert.NoError(t, err)

		assert.Equal(t, expectedTx, tx)

		rpc.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		txID := ids.New()

		rpc.On("GetTransaction", ctx, mock.Anything).
			Return(nil, errors.New("rpc error"))

		c := client.NewFromRPCClient(rpc)

		result, err := c.GetTransaction(ctx, txID)
		assert.Error(t, err)
		assert.Nil(t, result)

		rpc.AssertExpectations(t)
	})
}

func TestClient_GetTransactionResult(t *testing.T) {
	results := test.TransactionResultGenerator()
	ids := test.IdentifierGenerator()

	t.Run("Success", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		txID := ids.New()
		expectedResult := results.New()
		response, _ := convert.TransactionResultToMessage(expectedResult)

		rpc.On("GetTransactionResult", ctx, mock.Anything).Return(response, nil)

		c := client.NewFromRPCClient(rpc)

		result, err := c.GetTransactionResult(ctx, txID)
		assert.NoError(t, err)

		assert.Equal(t, expectedResult, *result)

		rpc.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		txID := ids.New()

		rpc.On("GetTransactionResult", ctx, mock.Anything).
			Return(nil, errors.New("rpc error"))

		c := client.NewFromRPCClient(rpc)

		result, err := c.GetTransactionResult(ctx, txID)
		assert.Error(t, err)
		assert.Nil(t, result)

		rpc.AssertExpectations(t)
	})
}

func TestClient_GetEventsForHeightRange(t *testing.T) {
	ids := test.IdentifierGenerator()
	events := test.EventGenerator()

	t.Run("Empty result", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		response := &access.EventsResponse{
			Results: []*access.EventsResponse_Result{},
		}

		rpc.On("GetEventsForHeightRange", ctx, mock.Anything).Return(response, nil)

		c := client.NewFromRPCClient(rpc)

		blocks, err := c.GetEventsForHeightRange(ctx, client.EventRangeQuery{
			Type:        "foo",
			StartHeight: 1,
			EndHeight:   10,
		})
		assert.NoError(t, err)

		assert.Empty(t, blocks)

		rpc.AssertExpectations(t)
	})

	t.Run("Non-empty result", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

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

		c := client.NewFromRPCClient(rpc)

		blocks, err := c.GetEventsForHeightRange(ctx, client.EventRangeQuery{
			Type:        "foo",
			StartHeight: 1,
			EndHeight:   10,
		})
		assert.NoError(t, err)

		assert.Len(t, blocks, len(response.Results))

		assert.Equal(t, response.Results[0].BlockId, blocks[0].BlockID.Bytes())
		assert.Equal(t, response.Results[0].BlockHeight, blocks[0].Height)

		assert.Equal(t, response.Results[1].BlockId, blocks[1].BlockID.Bytes())
		assert.Equal(t, response.Results[1].BlockHeight, blocks[1].Height)

		assert.Equal(t, eventA, blocks[0].Events[0])
		assert.Equal(t, eventB, blocks[0].Events[1])
		assert.Equal(t, eventC, blocks[1].Events[0])
		assert.Equal(t, eventD, blocks[1].Events[1])

		rpc.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		rpc.On("GetEventsForHeightRange", ctx, mock.Anything).
			Return(nil, errors.New("rpc error"))

		c := client.NewFromRPCClient(rpc)

		blocks, err := c.GetEventsForHeightRange(ctx, client.EventRangeQuery{
			Type:        "foo",
			StartHeight: 1,
			EndHeight:   10,
		})

		assert.Error(t, err)
		assert.Empty(t, blocks)

		rpc.AssertExpectations(t)
	})
}

func TestClient_GetEventsForBlockIDs(t *testing.T) {
	ids := test.IdentifierGenerator()
	events := test.EventGenerator()

	t.Run("Empty result", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		blockIDs := []flow.Identifier{ids.New(), ids.New()}

		response := &access.EventsResponse{
			Results: []*access.EventsResponse_Result{},
		}

		rpc.On("GetEventsForBlockIDs", ctx, mock.Anything).Return(response, nil)

		c := client.NewFromRPCClient(rpc)

		blocks, err := c.GetEventsForBlockIDs(ctx, "foo", blockIDs)
		assert.NoError(t, err)

		assert.Empty(t, blocks)

		rpc.AssertExpectations(t)
	})

	t.Run("Non-empty result", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

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

		c := client.NewFromRPCClient(rpc)

		blocks, err := c.GetEventsForBlockIDs(ctx, "foo", []flow.Identifier{blockIDA, blockIDB})
		assert.NoError(t, err)

		assert.Len(t, blocks, len(response.Results))

		assert.Equal(t, response.Results[0].BlockId, blocks[0].BlockID.Bytes())
		assert.Equal(t, response.Results[0].BlockHeight, blocks[0].Height)

		assert.Equal(t, response.Results[1].BlockId, blocks[1].BlockID.Bytes())
		assert.Equal(t, response.Results[1].BlockHeight, blocks[1].Height)

		assert.Equal(t, eventA, blocks[0].Events[0])
		assert.Equal(t, eventB, blocks[0].Events[1])
		assert.Equal(t, eventC, blocks[1].Events[0])
		assert.Equal(t, eventD, blocks[1].Events[1])

		rpc.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		blockIDA, blockIDB := ids.New(), ids.New()

		rpc.On("GetEventsForBlockIDs", ctx, mock.Anything).
			Return(nil, errors.New("rpc error"))

		c := client.NewFromRPCClient(rpc)

		blocks, err := c.GetEventsForBlockIDs(ctx, "foo", []flow.Identifier{blockIDA, blockIDB})

		assert.Error(t, err)
		assert.Empty(t, blocks)

		rpc.AssertExpectations(t)
	})
}

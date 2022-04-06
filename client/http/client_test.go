package http

import (
	"context"
	"testing"

	"github.com/onflow/flow-go/engine/access/rest/models"

	"github.com/onflow/flow-go-sdk/client/convert"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func clientTest(
	f func(ctx context.Context, t *testing.T, handler *mockHandler, client *BaseClient),
) func(t *testing.T) {
	return func(t *testing.T) {
		h := &mockHandler{}
		client := &BaseClient{
			&HTTPClient{h},
		}
		f(context.Background(), t, h, client)
		h.AssertExpectations(t)
	}
}

func TestBaseClient_GetBlockByID(t *testing.T) {
	const handlerName = "getBlockByID"
	t.Run("Success", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *BaseClient) {
		httpBlock := test.BlockHTTP()
		expectedBlock, err := convert.HTTPToBlock(&httpBlock)
		assert.NoError(t, err)

		handler.
			On(handlerName, mock.Anything, httpBlock.Header.Id).
			Return(&httpBlock, nil)

		block, err := client.GetBlockByID(ctx, flow.HexToID(httpBlock.Header.Id))
		assert.NoError(t, err)
		assert.Equal(t, block, expectedBlock)
	}))

	t.Run("Not found", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *BaseClient) {
		handler.
			On(handlerName, mock.Anything, mock.Anything).
			Return(nil, HttpError{
				Url:     "/",
				Code:    404,
				Message: "block not found",
			})
		_, err := client.GetBlockByID(ctx, flow.HexToID("0x1"))
		assert.EqualError(t, err, "block not found")
	}))
}

func TestBaseClient_GetBlockByHeight(t *testing.T) {
	const handlerName = "getBlocksByHeights"

	t.Run("Success", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *BaseClient) {
		httpBlock := test.BlockHTTP()
		expectedBlock, err := convert.HTTPToBlock(&httpBlock)
		assert.NoError(t, err)

		handler.
			On(handlerName, mock.Anything, httpBlock.Header.Height, "", "").
			Return([]*models.Block{&httpBlock}, nil)

		block, err := client.GetBlockByHeight(ctx, expectedBlock.Height)
		assert.NoError(t, err)
		assert.Equal(t, block, expectedBlock)
	}))

	t.Run("Not found", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *BaseClient) {
		handler.
			On(handlerName, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, HttpError{
				Url:     "/",
				Code:    404,
				Message: "block not found",
			})

		_, err := client.GetBlockByHeight(ctx, 10)
		assert.EqualError(t, err, "block not found")
	}))
}

func TestBaseClient_GetCollection(t *testing.T) {
	const handlerName = "getCollection"

	t.Run("Success", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *BaseClient) {
		httpCollection := test.CollectionHTTP()
		expectedCollection := convert.HTTPToCollection(&httpCollection)

		handler.
			On(handlerName, mock.Anything, expectedCollection.ID().String()).
			Return(&httpCollection, nil)

		collection, err := client.GetCollection(ctx, expectedCollection.ID())

		assert.NoError(t, err)
		assert.Equal(t, collection, expectedCollection)
	}))

	t.Run("Not Found", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *BaseClient) {
		handler.
			On(handlerName, mock.Anything, mock.Anything).
			Return(nil, HttpError{
				Url:     "/",
				Code:    404,
				Message: "collection not found",
			})

		_, err := client.GetCollection(ctx, flow.HexToID("0x1"))
		assert.EqualError(t, err, "collection not found")
	}))
}

func TestBaseClient_SendTransaction(t *testing.T) {
	const handlerName = "sendTransaction"

	t.Run("Success", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *BaseClient) {
		httpTx := test.TransactionHTTP()
		expectedTx, err := convert.HTTPToTransaction(&httpTx)
		assert.NoError(t, err)

		sentTx, err := convert.TransactionToHTTP(*expectedTx)
		assert.NoError(t, err)

		handler.
			On(handlerName, mock.Anything, sentTx).
			Return(nil)

		err = client.SendTransaction(ctx, *expectedTx)
		assert.NoError(t, err)
	}))

	t.Run("Not Found", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *BaseClient) {
		handler.On(handlerName, mock.Anything, mock.Anything).Return(HttpError{
			Url:     "/",
			Code:    400,
			Message: "invalid payload",
		})

		tx := test.TransactionGenerator().New()
		err := client.SendTransaction(ctx, *tx)
		assert.EqualError(t, err, "invalid payload")
	}))
}

func TestBaseClient_GetTransaction(t *testing.T) {
	const handlerName = "getTransaction"

	t.Run("Success", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *BaseClient) {
		httpTx := test.TransactionHTTP()
		expectedTx, err := convert.HTTPToTransaction(&httpTx)
		assert.NoError(t, err)

		handler.
			On(handlerName, mock.Anything, expectedTx.ID().String(), false).
			Return(&httpTx, nil)

		tx, err := client.GetTransaction(ctx, expectedTx.ID())
		assert.NoError(t, err)
		assert.Equal(t, tx, expectedTx)
	}))

	t.Run("Not Found", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *BaseClient) {
		handler.On(handlerName, mock.Anything, mock.Anything, mock.Anything).Return(nil, HttpError{
			Url:     "/",
			Code:    404,
			Message: "tx not found",
		})

		_, err := client.GetTransaction(ctx, flow.HexToID("0x1"))
		assert.EqualError(t, err, "tx not found")
	}))
}

func TestBaseClient_GetTransactionResult(t *testing.T) {
	const handlerName = "getTransaction"

	t.Run("Success", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *BaseClient) {
		httpTx := test.TransactionHTTP()
		httpTxRes := test.TransactionResultHTTP()
		httpTx.Result = &httpTxRes
		expectedTx, err := convert.HTTPToTransaction(&httpTx)
		assert.NoError(t, err)

		expectedTxRes, err := convert.HTTPToTransactionResult(&httpTxRes)
		assert.NoError(t, err)

		handler.
			On(handlerName, mock.Anything, expectedTx.ID().String(), true).
			Return(&httpTx, nil)

		txRes, err := client.GetTransactionResult(ctx, expectedTx.ID())
		assert.NoError(t, err)
		assert.Equal(t, txRes, expectedTxRes)
	}))

	t.Run("Not Found", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *BaseClient) {
		handler.On(handlerName, mock.Anything, mock.Anything, true).Return(nil, HttpError{
			Url:     "/",
			Code:    404,
			Message: "tx result not found",
		})

		_, err := client.GetTransactionResult(ctx, flow.HexToID("0x1"))
		assert.EqualError(t, err, "tx result not found")
	}))
}

func TestBaseClient_GetAccount(t *testing.T) {
	const handlerName = "getAccount"

	t.Run("Success", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *BaseClient) {
		httpAccount := test.AccountHTTP()
		expectedAccount, err := convert.HTTPToAccount(&httpAccount)
		assert.NoError(t, err)

		handler.
			On(handlerName, mock.Anything, httpAccount.Address, "sealed").
			Return(&httpAccount, nil)

		account, err := client.GetAccount(ctx, expectedAccount.Address)
		assert.NoError(t, err)
		assert.Equal(t, account, expectedAccount)
	}))

	t.Run("Not Found", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *BaseClient) {
		handler.On(handlerName, mock.Anything, mock.Anything, mock.Anything).Return(nil, HttpError{
			Url:     "/",
			Code:    404,
			Message: "account not found",
		})
	}))
}

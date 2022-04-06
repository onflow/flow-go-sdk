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
}

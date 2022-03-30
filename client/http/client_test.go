package http

import (
	"context"
	"testing"

	"github.com/onflow/flow-go-sdk/client/convert"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func clientTest(
	f func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client),
) func(t *testing.T) {
	return func(t *testing.T) {
		h := &mockHandler{}
		client := NewClient(h)
		f(context.Background(), t, h, client)
		h.AssertExpectations(t)
	}
}

func Test_GetBlockByID(t *testing.T) {
	const handlerName = "getBlockByID"
	t.Run("Success", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		httpBlock := test.BlockHTTP()
		expectedBlock, err := convert.HTTPToBlock(&httpBlock)
		assert.NoError(t, err)

		handler.On(handlerName, mock.Anything, httpBlock.Header.Id).Return(&httpBlock, nil)

		block, err := client.GetBlockByID(ctx, flow.HexToID(httpBlock.Header.Id))
		assert.NoError(t, err)
		assert.Equal(t, block, expectedBlock)
	}))

	t.Run("Not found", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		handler.On(handlerName, mock.Anything, mock.Anything).Return(nil)
	}))

	t.Run("Bad request", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {

	}))
}

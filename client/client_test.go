package client_test

import (
	"context"
	"errors"
	"testing"

	"github.com/dapperlabs/flow/protobuf/go/flow/access"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/dapperlabs/flow-go-sdk/client"
	"github.com/dapperlabs/flow-go-sdk/client/convert"
	"github.com/dapperlabs/flow-go-sdk/client/mocks"
	"github.com/dapperlabs/flow-go-sdk/test"
)

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

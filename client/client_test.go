package client_test

import (
	"context"
	"errors"
	"testing"

	"github.com/dapperlabs/flow/protobuf/go/flow/access"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/dapperlabs/flow-go-sdk/client"
	"github.com/dapperlabs/flow-go-sdk/client/mocks"
	"github.com/dapperlabs/flow-go-sdk/test"
)

func TestClient_SendTransaction(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		tx := test.TransactionGenerator().New()

		rpc.On("SendTransaction", ctx, mock.Anything).
			Return(
				&access.SendTransactionResponse{
					Id: tx.ID().Bytes(),
				},
				nil,
			)

		c := client.NewFromRPCClient(rpc)

		err := c.SendTransaction(ctx, *tx)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		rpc := &mocks.RPCClient{}

		ctx := context.Background()

		tx := test.TransactionGenerator().New()

		rpc.On("SendTransaction", ctx, mock.Anything).
			Return(nil, errors.New("rpc error"))

		c := client.NewFromRPCClient(rpc)

		err := c.SendTransaction(ctx, *tx)

		assert.Error(t, err)
	})
}

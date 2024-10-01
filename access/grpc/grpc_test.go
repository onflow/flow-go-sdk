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

package grpc

import (
	"context"
	"io"
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow/protobuf/go/flow/access"
	"github.com/onflow/flow/protobuf/go/flow/entities"
	"github.com/onflow/flow/protobuf/go/flow/executiondata"

	"github.com/onflow/flow-go-sdk/access/grpc/convert"
	"github.com/onflow/flow-go-sdk/access/grpc/mocks"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/test"
)

var (
	errInternal = status.Error(codes.Internal, "internal server error")
	errNotFound = status.Error(codes.NotFound, "not found")
)

func clientTest(
	f func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, client *BaseClient),
) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		rpc := new(mocks.MockRPCClient)
		c := NewFromRPCClient(rpc)
		f(t, ctx, rpc, c)
		rpc.AssertExpectations(t)
	}
}

func executionDataClientTest(
	f func(t *testing.T, ctx context.Context, rpc *mocks.MockExecutionDataRPCClient, client *BaseClient),
) func(t *testing.T) {
	return func(t *testing.T) {
		ctx := context.Background()
		rpc := new(mocks.MockExecutionDataRPCClient)
		c := NewFromExecutionDataRPCClient(rpc)
		f(t, ctx, rpc, c)
		rpc.AssertExpectations(t)
	}
}

func Test_ClientOptions(t *testing.T) {
	t.Run("WithJSONOptions", func(t *testing.T) {

		expectedJsonOption := []jsoncdc.Option{jsoncdc.WithBackwardsCompatibility()}

		// Confirm that the options are set
		options := WithJSONOptions(expectedJsonOption...)
		cfg := DefaultClientOptions()
		options(cfg)

		// hard to run a contains check on the options due to it comparing functions, so just check the length
		assert.Equal(t, len(cfg.jsonOptions), len(expectedJsonOption)+len(DefaultClientOptions().jsonOptions))
	})

	t.Run("WithGRPCDialOptions", func(t *testing.T) {

		expectedDialOption := []grpc.DialOption{grpc.WithAuthority("test"), grpc.WithBlock()}

		options := WithGRPCDialOptions(expectedDialOption...)
		cfg := DefaultClientOptions()
		options(cfg)

		// hard to run a contains check on the options due to it comparing functions, so just check the length
		assert.Equal(t, len(cfg.dialOptions), len(expectedDialOption)+len(DefaultClientOptions().dialOptions))
	})
}

func TestClient_Ping(t *testing.T) {
	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		response := &access.PingResponse{}

		rpc.On("Ping", ctx, mock.Anything).Return(response, nil)

		err := c.Ping(ctx)
		assert.NoError(t, err)
	}))

	t.Run("Internal error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		rpc.On("Ping", ctx, mock.Anything).
			Return(nil, errInternal)

		err := c.Ping(ctx)
		assert.Error(t, err)
		assert.Equal(t, codes.Internal, status.Code(err))
	}))
}

func TestClient_GetNetworkParameters(t *testing.T) {
	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		response := &access.GetNetworkParametersResponse{
			ChainId: "flow-testnet",
		}
		expectedParams := &flow.NetworkParameters{
			ChainID: flow.ChainID("flow-testnet"),
		}

		rpc.On("GetNetworkParameters", ctx, mock.Anything).Return(response, nil)

		params, err := c.GetNetworkParameters(ctx)
		require.NoError(t, err)

		assert.Equal(t, params, expectedParams)
	}))

	t.Run("Internal error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		rpc.On("GetNetworkParameters", ctx, mock.Anything).
			Return(nil, errInternal)

		params, err := c.GetNetworkParameters(ctx)
		assert.Error(t, err)
		assert.Equal(t, codes.Internal, status.Code(err))
		assert.Nil(t, params)
	}))
}

func TestClient_GetNodeInfo(t *testing.T) {
	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		id := flow.HexToID("0x01")
		ver := uint64(1)
		spork := uint64(2)
		root := uint64(3)

		response := &access.GetNodeVersionInfoResponse{
			Info: &entities.NodeVersionInfo{
				Semver:               "1.0",
				Commit:               "123",
				SporkId:              id.Bytes(),
				ProtocolVersion:      ver,
				SporkRootBlockHeight: spork,
				NodeRootBlockHeight:  root,
			},
		}

		expected := &flow.NodeVersionInfo{
			Semver:               "1.0",
			Commit:               "123",
			SporkId:              id,
			ProtocolVersion:      ver,
			SporkRootBlockHeight: spork,
			NodeRootBlockHeight:  root,
		}

		rpc.On("GetNodeVersionInfo", ctx, mock.Anything).Return(response, nil)

		info, err := c.GetNodeVersionInfo(ctx)
		require.NoError(t, err)
		require.Equal(t, expected, info)
	}))
}

func TestClient_GetLatestBlockHeader(t *testing.T) {
	blocks := test.BlockGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		expectedHeader := blocks.New().BlockHeader

		b, err := convert.BlockHeaderToMessage(expectedHeader)
		require.NoError(t, err)

		response := &access.BlockHeaderResponse{
			Block: b,
		}

		rpc.On("GetLatestBlockHeader", ctx, mock.Anything).Return(response, nil)

		header, err := c.GetLatestBlockHeader(ctx, true)
		require.NoError(t, err)

		assert.Equal(t, expectedHeader, *header)
	}))

	t.Run("Internal error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		rpc.On("GetLatestBlockHeader", ctx, mock.Anything).
			Return(nil, errInternal)

		header, err := c.GetLatestBlockHeader(ctx, true)
		assert.Error(t, err)
		assert.Equal(t, codes.Internal, status.Code(err))
		assert.Nil(t, header)
	}))
}

func TestClient_GetBlockHeaderByID(t *testing.T) {
	blocks := test.BlockGenerator()
	ids := test.IdentifierGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		blockID := ids.New()
		expectedHeader := blocks.New().BlockHeader

		b, err := convert.BlockHeaderToMessage(expectedHeader)
		require.NoError(t, err)

		response := &access.BlockHeaderResponse{
			Block: b,
		}

		rpc.On("GetBlockHeaderByID", ctx, mock.Anything).Return(response, nil)

		header, err := c.GetBlockHeaderByID(ctx, blockID)
		require.NoError(t, err)

		assert.Equal(t, expectedHeader, *header)
	}))

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		blockID := ids.New()

		rpc.On("GetBlockHeaderByID", ctx, mock.Anything).
			Return(nil, errNotFound)

		header, err := c.GetBlockHeaderByID(ctx, blockID)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Nil(t, header)
	}))
}

func TestClient_GetBlockHeaderByHeight(t *testing.T) {
	blocks := test.BlockGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		expectedHeader := blocks.New().BlockHeader

		b, err := convert.BlockHeaderToMessage(expectedHeader)
		require.NoError(t, err)

		response := &access.BlockHeaderResponse{
			Block: b,
		}

		rpc.On("GetBlockHeaderByHeight", ctx, mock.Anything).Return(response, nil)

		header, err := c.GetBlockHeaderByHeight(ctx, 42)
		require.NoError(t, err)

		assert.Equal(t, expectedHeader, *header)
	}))

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		rpc.On("GetBlockHeaderByHeight", ctx, mock.Anything).
			Return(nil, errNotFound)

		header, err := c.GetBlockHeaderByHeight(ctx, 42)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Nil(t, header)
	}))
}

func TestClient_GetLatestBlock(t *testing.T) {
	blocks := test.BlockGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		expectedBlock := blocks.New()

		b, err := convert.BlockToMessage(*expectedBlock)
		require.NoError(t, err)

		response := &access.BlockResponse{
			Block: b,
		}

		rpc.On("GetLatestBlock", ctx, mock.Anything).Return(response, nil)

		block, err := c.GetLatestBlock(ctx, true)
		require.NoError(t, err)

		assert.Equal(t, expectedBlock, block)
	}))

	t.Run("Internal error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		rpc.On("GetLatestBlock", ctx, mock.Anything).
			Return(nil, errInternal)

		block, err := c.GetLatestBlock(ctx, true)
		assert.Error(t, err)
		assert.Equal(t, codes.Internal, status.Code(err))
		assert.Nil(t, block)
	}))
}

func TestClient_GetBlockByID(t *testing.T) {
	blocks := test.BlockGenerator()
	ids := test.IdentifierGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		blockID := ids.New()
		expectedBlock := blocks.New()

		b, err := convert.BlockToMessage(*expectedBlock)
		require.NoError(t, err)

		response := &access.BlockResponse{
			Block: b,
		}

		rpc.On("GetBlockByID", ctx, mock.Anything).Return(response, nil)

		block, err := c.GetBlockByID(ctx, blockID)
		require.NoError(t, err)

		assert.Equal(t, expectedBlock, block)
	}))

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		blockID := ids.New()

		rpc.On("GetBlockByID", ctx, mock.Anything).
			Return(nil, errNotFound)

		block, err := c.GetBlockByID(ctx, blockID)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Nil(t, block)
	}))
}

func TestClient_GetBlockByHeight(t *testing.T) {
	blocks := test.BlockGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		expectedBlock := blocks.New()

		b, err := convert.BlockToMessage(*expectedBlock)
		require.NoError(t, err)

		response := &access.BlockResponse{
			Block: b,
		}

		rpc.On("GetBlockByHeight", ctx, mock.Anything).Return(response, nil)

		block, err := c.GetBlockByHeight(ctx, 42)
		require.NoError(t, err)

		assert.Equal(t, expectedBlock, block)
	}))

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		rpc.On("GetBlockByHeight", ctx, mock.Anything).
			Return(nil, errNotFound)

		block, err := c.GetBlockByHeight(ctx, 42)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Nil(t, block)
	}))
}

func TestClient_GetCollection(t *testing.T) {
	cols := test.LightCollectionGenerator()
	ids := test.IdentifierGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
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

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		colID := ids.New()

		rpc.On("GetCollectionByID", ctx, mock.Anything).
			Return(nil, errNotFound)

		col, err := c.GetCollection(ctx, colID)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Nil(t, col)
	}))
}

func TestClient_GetFullCollectionById(t *testing.T) {
	collections := test.FullCollectionGenerator()
	ids := test.IdentifierGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		expectedCollection := collections.New()
		txs, err := convert.FullCollectionToTransactionsMessage(*expectedCollection)
		require.NoError(t, err)

		response := &access.FullCollectionResponse{
			Transactions: txs,
		}

		rpc.
			On("GetFullCollectionByID", ctx, mock.Anything).
			Return(response, nil)

		id := ids.New()
		actualCollection, err := c.GetFullCollectionByID(ctx, id)
		require.NoError(t, err)

		require.Equal(t, expectedCollection, actualCollection)

	}))

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		rpc.
			On("GetFullCollectionByID", ctx, mock.Anything).
			Return(nil, errNotFound)

		id := ids.New()
		col, err := c.GetFullCollectionByID(ctx, id)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Nil(t, col)
	}))
}

func TestClient_SendTransaction(t *testing.T) {
	transactions := test.TransactionGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		tx := transactions.New()

		response := &access.SendTransactionResponse{
			Id: tx.ID().Bytes(),
		}

		rpc.On("SendTransaction", ctx, mock.Anything).Return(response, nil)

		err := c.SendTransaction(ctx, *tx)
		require.NoError(t, err)
	}))

	t.Run("Internal error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		tx := transactions.New()

		rpc.On("SendTransaction", ctx, mock.Anything).
			Return(nil, errInternal)

		err := c.SendTransaction(ctx, *tx)
		assert.Error(t, err)
		assert.Equal(t, codes.Internal, status.Code(err))
	}))
}

func TestClient_GetSystemTransaction(t *testing.T) {
	txs := test.TransactionGenerator()
	ids := test.IdentifierGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		blockID := ids.New()
		expectedTx := txs.New()

		txMsg, err := convert.TransactionToMessage(*expectedTx)
		require.NoError(t, err)

		response := &access.TransactionResponse{
			Transaction: txMsg,
		}

		rpc.On("GetSystemTransaction", ctx, mock.Anything).Return(response, nil)

		tx, err := c.GetSystemTransaction(ctx, blockID)
		require.NoError(t, err)

		assert.Equal(t, expectedTx, tx)
	}))

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		blockID := ids.New()

		rpc.On("GetSystemTransaction", ctx, mock.Anything).
			Return(nil, errNotFound)

		tx, err := c.GetSystemTransaction(ctx, blockID)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Nil(t, tx)
	}))
}

func TestClient_GetTransaction(t *testing.T) {
	txs := test.TransactionGenerator()
	ids := test.IdentifierGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		txID := ids.New()
		expectedTx := txs.New()

		txMsg, err := convert.TransactionToMessage(*expectedTx)
		require.NoError(t, err)

		response := &access.TransactionResponse{
			Transaction: txMsg,
		}

		rpc.On("GetTransaction", ctx, mock.Anything).Return(response, nil)

		tx, err := c.GetTransaction(ctx, txID)
		require.NoError(t, err)

		assert.Equal(t, expectedTx, tx)
	}))

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		txID := ids.New()

		rpc.On("GetTransaction", ctx, mock.Anything).
			Return(nil, errNotFound)

		tx, err := c.GetTransaction(ctx, txID)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Nil(t, tx)
	}))
}

func TestClient_GetTransactionsByBlockID(t *testing.T) {
	txs := test.TransactionGenerator()
	ids := test.IdentifierGenerator()
	blockID := ids.New()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		expectedTx := txs.New()

		txMsg, err := convert.TransactionToMessage(*expectedTx)
		require.NoError(t, err)

		responses := &access.TransactionsResponse{
			Transactions: []*entities.Transaction{txMsg},
		}

		rpc.On("GetTransactionsByBlockID", ctx, mock.Anything).Return(responses, nil)

		txs, err := c.GetTransactionsByBlockID(ctx, blockID)
		require.NoError(t, err)

		assert.Equal(t, len(txs), 1)
		assert.Equal(t, expectedTx, txs[0])
	}))

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		rpc.On("GetTransactionsByBlockID", ctx, mock.Anything).
			Return(nil, errNotFound)

		tx, err := c.GetTransactionsByBlockID(ctx, blockID)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Nil(t, tx)
	}))
}

func TestClient_GetSystemTransactionResult(t *testing.T) {
	ids := test.IdentifierGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		results := test.TransactionResultGenerator(flow.EventEncodingVersionCCF)
		blockID := ids.New()
		expectedResult := results.New()
		response, _ := convert.TransactionResultToMessage(expectedResult, flow.EventEncodingVersionCCF)

		rpc.On("GetSystemTransactionResult", ctx, mock.Anything).Return(response, nil)

		result, err := c.GetSystemTransactionResult(ctx, blockID)
		require.NoError(t, err)

		// Force evaluation of type ID, which is cached in type.
		// Necessary for equality check below
		for _, event := range result.Events {
			_ = event.Value.Type().ID()
		}

		assert.Equal(t, expectedResult, *result)
	}))

	t.Run("Success with jsoncdc", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		results := test.TransactionResultGenerator(flow.EventEncodingVersionJSONCDC)
		blockID := ids.New()
		expectedResult := results.New()
		response, _ := convert.TransactionResultToMessage(expectedResult, flow.EventEncodingVersionJSONCDC)

		rpc.On("GetSystemTransactionResult", ctx, mock.Anything).Return(response, nil)

		result, err := c.GetSystemTransactionResult(ctx, blockID)
		require.NoError(t, err)

		// Force evaluation of type ID, which is cached in type.
		// Necessary for equality check below
		for _, event := range result.Events {
			_ = event.Value.Type().ID()
		}

		assert.Equal(t, expectedResult, *result)
	}))

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		blockID := ids.New()

		rpc.On("GetSystemTransactionResult", ctx, mock.Anything).
			Return(nil, errNotFound)

		result, err := c.GetSystemTransactionResult(ctx, blockID)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Nil(t, result)
	}))
}

func TestClient_GetTransactionResult(t *testing.T) {
	ids := test.IdentifierGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		results := test.TransactionResultGenerator(flow.EventEncodingVersionCCF)
		txID := ids.New()
		expectedResult := results.New()
		response, _ := convert.TransactionResultToMessage(expectedResult, flow.EventEncodingVersionCCF)

		rpc.On("GetTransactionResult", ctx, mock.Anything).Return(response, nil)

		result, err := c.GetTransactionResult(ctx, txID)
		require.NoError(t, err)

		// Force evaluation of type ID, which is cached in type.
		// Necessary for equality check below
		for _, event := range result.Events {
			_ = event.Value.Type().ID()
		}

		assert.Equal(t, expectedResult, *result)

	}))

	t.Run("Success with jsoncdc", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		results := test.TransactionResultGenerator(flow.EventEncodingVersionJSONCDC)
		txID := ids.New()
		expectedResult := results.New()
		response, _ := convert.TransactionResultToMessage(expectedResult, flow.EventEncodingVersionJSONCDC)

		rpc.On("GetTransactionResult", ctx, mock.Anything).Return(response, nil)

		result, err := c.GetTransactionResult(ctx, txID)
		require.NoError(t, err)

		// Force evaluation of type ID, which is cached in type.
		// Necessary for equality check below
		for _, event := range result.Events {
			_ = event.Value.Type().ID()
		}

		assert.Equal(t, expectedResult, *result)

	}))

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		txID := ids.New()

		rpc.On("GetTransactionResult", ctx, mock.Anything).
			Return(nil, errNotFound)

		result, err := c.GetTransactionResult(ctx, txID)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Nil(t, result)
	}))
}

func TestClient_GetTransactionResultsByBlockID(t *testing.T) {
	ids := test.IdentifierGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		resultGenerator := test.TransactionResultGenerator(flow.EventEncodingVersionCCF)
		blockID := ids.New()
		expectedResult := resultGenerator.New()
		response, err := convert.TransactionResultToMessage(expectedResult, flow.EventEncodingVersionCCF)
		require.NoError(t, err)

		responses := &access.TransactionResultsResponse{
			TransactionResults: []*access.TransactionResultResponse{response},
		}

		rpc.On("GetTransactionResultsByBlockID", ctx, mock.Anything).Return(responses, nil)

		results, err := c.GetTransactionResultsByBlockID(ctx, blockID)
		require.NoError(t, err)

		// Force evaluation of type ID, which is cached in type.
		// Necessary for equality check below
		for _, result := range results {
			for _, event := range result.Events {
				_ = event.Value.Type().ID()
			}
		}

		assert.Equal(t, len(results), 1)
		assert.Equal(t, expectedResult, *results[0])
	}))

	t.Run("Success with jsoncdc", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		resultGenerator := test.TransactionResultGenerator(flow.EventEncodingVersionJSONCDC)
		blockID := ids.New()
		expectedResult := resultGenerator.New()
		response, err := convert.TransactionResultToMessage(expectedResult, flow.EventEncodingVersionJSONCDC)
		require.NoError(t, err)

		responses := &access.TransactionResultsResponse{
			TransactionResults: []*access.TransactionResultResponse{response},
		}

		rpc.On("GetTransactionResultsByBlockID", ctx, mock.Anything).Return(responses, nil)

		results, err := c.GetTransactionResultsByBlockID(ctx, blockID)
		require.NoError(t, err)

		// Force evaluation of type ID, which is cached in type.
		// Necessary for equality check below
		for _, result := range results {
			for _, event := range result.Events {
				_ = event.Value.Type().ID()
			}
		}

		assert.Equal(t, len(results), 1)
		assert.Equal(t, expectedResult, *results[0])
	}))

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		blockID := ids.New()

		rpc.On("GetTransactionResultsByBlockID", ctx, mock.Anything).
			Return(nil, errNotFound)

		result, err := c.GetTransactionResultsByBlockID(ctx, blockID)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Nil(t, result)
	}))
}

func TestClient_GetAccountAtLatestBlock(t *testing.T) {
	accounts := test.AccountGenerator()
	addresses := test.AddressGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		expectedAccount := accounts.New()
		response := &access.AccountResponse{
			Account: convert.AccountToMessage(*expectedAccount),
		}

		rpc.On("GetAccountAtLatestBlock", ctx, mock.Anything).Return(response, nil)

		account, err := c.GetAccountAtLatestBlock(ctx, expectedAccount.Address)
		require.NoError(t, err)

		assert.Equal(t, expectedAccount, account)
	}))

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		address := addresses.New()

		rpc.On("GetAccountAtLatestBlock", ctx, mock.Anything).
			Return(nil, errNotFound)

		account, err := c.GetAccountAtLatestBlock(ctx, address)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Nil(t, account)
	}))
}

func TestClient_GetAccountAtBlockHeight(t *testing.T) {
	accounts := test.AccountGenerator()
	addresses := test.AddressGenerator()
	height := uint64(42)

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		expectedAccount := accounts.New()
		response := &access.AccountResponse{
			Account: convert.AccountToMessage(*expectedAccount),
		}

		rpc.On("GetAccountAtBlockHeight", ctx, mock.Anything).Return(response, nil)

		account, err := c.GetAccountAtBlockHeight(ctx, expectedAccount.Address, height)
		require.NoError(t, err)

		assert.Equal(t, expectedAccount, account)
	}))

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		address := addresses.New()

		rpc.On("GetAccountAtBlockHeight", ctx, mock.Anything).
			Return(nil, errNotFound)

		account, err := c.GetAccountAtBlockHeight(ctx, address, height)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Nil(t, account)
	}))
}

func TestClient_GetAccountBalanceAtLatestBlock(t *testing.T) {
	accounts := test.AccountGenerator()
	addresses := test.AddressGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		account := accounts.New()

		response := &access.AccountBalanceResponse{
			Balance: account.Balance,
		}

		rpc.On("GetAccountBalanceAtLatestBlock", ctx, mock.Anything).Return(response, nil)

		balance, err := c.GetAccountBalanceAtLatestBlock(ctx, account.Address)
		require.NoError(t, err)

		assert.Equal(t, account.Balance, balance)

	}))

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		address := addresses.New()

		rpc.On("GetAccountBalanceAtLatestBlock", ctx, mock.Anything).
			Return(nil, errNotFound)

		balance, err := c.GetAccountBalanceAtLatestBlock(ctx, address)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Equal(t, balance, uint64(0))
	}))
}

func TestClient_GetAccountBalanceAtBlockHeight(t *testing.T) {
	accounts := test.AccountGenerator()
	addresses := test.AddressGenerator()
	blockHeight := uint64(42)

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		account := accounts.New()

		response := &access.AccountBalanceResponse{
			Balance: account.Balance,
		}

		rpc.On("GetAccountBalanceAtBlockHeight", ctx, mock.Anything).Return(response, nil)

		balance, err := c.GetAccountBalanceAtBlockHeight(ctx, account.Address, blockHeight)
		require.NoError(t, err)

		assert.Equal(t, account.Balance, balance)

	}))

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		address := addresses.New()

		rpc.On("GetAccountBalanceAtBlockHeight", ctx, mock.Anything).
			Return(nil, errNotFound)

		balance, err := c.GetAccountBalanceAtBlockHeight(ctx, address, blockHeight)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Equal(t, balance, uint64(0))
	}))
}

func TestClient_ExecuteScriptAtLatestBlock(t *testing.T) {
	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		expectedValue := cadence.NewInt(42)
		encodedValue, err := jsoncdc.Encode(expectedValue)
		require.NoError(t, err)

		response := &access.ExecuteScriptResponse{
			Value: encodedValue,
		}

		rpc.On("ExecuteScriptAtLatestBlock", ctx, mock.Anything).Return(response, nil)

		var value cadence.Value
		value, err = c.ExecuteScriptAtLatestBlock(ctx, []byte("foo"), nil)
		require.NoError(t, err)

		assert.Equal(t, expectedValue, value)
	}))

	t.Run("Arguments", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		expectedValue := cadence.NewInt(42)
		encodedValue, err := jsoncdc.Encode(expectedValue)
		require.NoError(t, err)

		arg := cadence.String("test")
		expectedArgs, err := jsoncdc.Encode(arg)
		require.NoError(t, err)

		rpcReq := &access.ExecuteScriptAtLatestBlockRequest{
			Script:    []byte("foo"),
			Arguments: [][]byte{expectedArgs},
		}

		response := &access.ExecuteScriptResponse{
			Value: encodedValue,
		}

		rpc.On("ExecuteScriptAtLatestBlock", ctx, rpcReq).Return(response, nil)

		value, err := c.ExecuteScriptAtLatestBlock(ctx, []byte("foo"), []cadence.Value{arg})
		require.NoError(t, err)

		assert.Equal(t, expectedValue, value)
	}))

	t.Run(
		"Invalid JSON-CDC",
		clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
			response := &access.ExecuteScriptResponse{
				Value: []byte("invalid JSON-CDC bytes"),
			}

			rpc.On("ExecuteScriptAtLatestBlock", ctx, mock.Anything).Return(response, nil)

			value, err := c.ExecuteScriptAtLatestBlock(ctx, []byte("foo"), nil)
			assert.Error(t, err)
			assert.Nil(t, value)
		}),
	)

	t.Run("Internal error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		rpc.On("ExecuteScriptAtLatestBlock", ctx, mock.Anything).
			Return(nil, errInternal)

		value, err := c.ExecuteScriptAtLatestBlock(ctx, []byte("foo"), nil)
		assert.Error(t, err)
		assert.Equal(t, codes.Internal, status.Code(err))
		assert.Nil(t, value)
	}))
}

func TestClient_GetAccountKeyAtLatestBlock(t *testing.T) {
	accounts := test.AccountGenerator()
	addresses := test.AddressGenerator()
	index := uint32(0)

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		account := accounts.New()
		response := &access.AccountKeyResponse{
			AccountKey: convert.AccountKeyToMessage(account.Keys[index]),
		}

		rpc.
			On("GetAccountKeyAtLatestBlock", ctx, mock.Anything).
			Return(response, nil)

		key, err := c.GetAccountKeyAtLatestBlock(ctx, account.Address, index)
		require.NoError(t, err)

		assert.Equal(t, account.Keys[index], key)
	}))

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		address := addresses.New()

		rpc.
			On("GetAccountKeyAtLatestBlock", ctx, mock.Anything).
			Return(nil, errNotFound)

		key, err := c.GetAccountKeyAtLatestBlock(ctx, address, index)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Nil(t, key)
	}))
}

func TestClient_GetAccountKeyAtBlockHeight(t *testing.T) {
	accounts := test.AccountGenerator()
	addresses := test.AddressGenerator()
	height := uint64(42)
	index := uint32(0)

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		account := accounts.New()
		response := &access.AccountKeyResponse{
			AccountKey: convert.AccountKeyToMessage(account.Keys[index]),
		}

		rpc.
			On("GetAccountKeyAtBlockHeight", ctx, mock.Anything).
			Return(response, nil)

		key, err := c.GetAccountKeyAtBlockHeight(ctx, account.Address, index, height)
		require.NoError(t, err)

		assert.Equal(t, account.Keys[index], key)
	}))

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		address := addresses.New()

		rpc.
			On("GetAccountKeyAtBlockHeight", ctx, mock.Anything).
			Return(nil, errNotFound)

		key, err := c.GetAccountKeyAtBlockHeight(ctx, address, index, height)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Nil(t, key)
	}))
}

func TestClient_GetAccountKeysAtLatestBlock(t *testing.T) {
	accounts := test.AccountGenerator()
	addresses := test.AddressGenerator()
	index := uint32(0)

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		account := accounts.New()

		response := &access.AccountKeysResponse{
			AccountKeys: []*entities.AccountKey{convert.AccountKeyToMessage(account.Keys[index])},
		}

		rpc.
			On("GetAccountKeysAtLatestBlock", ctx, mock.Anything).
			Return(response, nil)

		keys, err := c.GetAccountKeysAtLatestBlock(ctx, account.Address)
		require.NoError(t, err)
		assert.Equal(t, *account.Keys[index], keys[index])
	}))

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		address := addresses.New()

		rpc.
			On("GetAccountKeysAtLatestBlock", ctx, mock.Anything).
			Return(nil, errNotFound)

		keys, err := c.GetAccountKeysAtLatestBlock(ctx, address)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Empty(t, keys)
	}))
}

func TestClient_GetAccountKeysAtBlockHeight(t *testing.T) {
	accounts := test.AccountGenerator()
	addresses := test.AddressGenerator()
	height := uint64(42)
	index := uint32(0)

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		account := accounts.New()

		response := &access.AccountKeysResponse{
			AccountKeys: []*entities.AccountKey{convert.AccountKeyToMessage(account.Keys[index])},
		}

		rpc.
			On("GetAccountKeysAtBlockHeight", ctx, mock.Anything).
			Return(response, nil)

		keys, err := c.GetAccountKeysAtBlockHeight(ctx, account.Address, height)
		require.NoError(t, err)
		assert.Equal(t, *account.Keys[index], keys[index])
	}))

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		address := addresses.New()

		rpc.
			On("GetAccountKeysAtBlockHeight", ctx, mock.Anything).
			Return(nil, errNotFound)

		keys, err := c.GetAccountKeysAtBlockHeight(ctx, address, height)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Empty(t, keys)
	}))
}

func TestClient_ExecuteScriptAtBlockID(t *testing.T) {
	ids := test.IdentifierGenerator()

	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		expectedValue := cadence.NewInt(42)
		encodedValue, err := jsoncdc.Encode(expectedValue)
		require.NoError(t, err)

		response := &access.ExecuteScriptResponse{
			Value: encodedValue,
		}

		rpc.On("ExecuteScriptAtBlockID", ctx, mock.Anything).Return(response, nil)

		value, err := c.ExecuteScriptAtBlockID(ctx, ids.New(), []byte("foo"), nil)
		require.NoError(t, err)

		assert.Equal(t, expectedValue, value)
	}))

	t.Run(
		"Invalid JSON-CDC",
		clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
			response := &access.ExecuteScriptResponse{
				Value: []byte("invalid JSON-CDC bytes"),
			}

			rpc.On("ExecuteScriptAtBlockID", ctx, mock.Anything).Return(response, nil)

			value, err := c.ExecuteScriptAtBlockID(ctx, ids.New(), []byte("foo"), nil)
			assert.Error(t, err)
			assert.Nil(t, value)
		}),
	)

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		rpc.On("ExecuteScriptAtBlockID", ctx, mock.Anything).
			Return(nil, errNotFound)

		value, err := c.ExecuteScriptAtBlockID(ctx, ids.New(), []byte("foo"), nil)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Nil(t, value)
	}))
}

func TestClient_ExecuteScriptAtBlockHeight(t *testing.T) {
	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		expectedValue := cadence.NewInt(42)
		encodedValue, err := jsoncdc.Encode(expectedValue)
		require.NoError(t, err)

		response := &access.ExecuteScriptResponse{
			Value: encodedValue,
		}

		rpc.On("ExecuteScriptAtBlockHeight", ctx, mock.Anything).Return(response, nil)

		value, err := c.ExecuteScriptAtBlockHeight(ctx, 42, []byte("foo"), nil)
		require.NoError(t, err)

		assert.Equal(t, expectedValue, value)
	}))

	t.Run(
		"Invalid JSON-CDC",
		clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
			response := &access.ExecuteScriptResponse{
				Value: []byte("invalid JSON-CDC bytes"),
			}

			rpc.On("ExecuteScriptAtBlockHeight", ctx, mock.Anything).Return(response, nil)

			value, err := c.ExecuteScriptAtBlockHeight(ctx, 42, []byte("foo"), nil)
			assert.Error(t, err)
			assert.Nil(t, value)
		}),
	)

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		rpc.On("ExecuteScriptAtBlockHeight", ctx, mock.Anything).
			Return(nil, errNotFound)

		value, err := c.ExecuteScriptAtBlockHeight(ctx, 42, []byte("foo"), nil)
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Nil(t, value)
	}))
}

func TestClient_GetEventsForHeightRange(t *testing.T) {
	ids := test.IdentifierGenerator()
	events := test.EventGenerator(flow.EventEncodingVersionCCF)

	t.Run(
		"Empty result",
		clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
			response := &access.EventsResponse{
				Results: []*access.EventsResponse_Result{},
			}

			rpc.On("GetEventsForHeightRange", ctx, mock.Anything).Return(response, nil)

			blocks, err := c.GetEventsForHeightRange(ctx, EventRangeQuery{
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
		clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
			eventA, eventB, eventC, eventD := events.New(), events.New(), events.New(), events.New()

			eventAMsg, _ := convert.EventToMessage(eventA, flow.EventEncodingVersionCCF)
			eventBMsg, _ := convert.EventToMessage(eventB, flow.EventEncodingVersionCCF)
			eventCMsg, _ := convert.EventToMessage(eventC, flow.EventEncodingVersionCCF)
			eventDMsg, _ := convert.EventToMessage(eventD, flow.EventEncodingVersionCCF)

			response := &access.EventsResponse{
				Results: []*access.EventsResponse_Result{
					{
						BlockId:        ids.New().Bytes(),
						BlockHeight:    1,
						BlockTimestamp: timestamppb.Now(),
						Events: []*entities.Event{
							eventAMsg,
							eventBMsg,
						},
					},
					{
						BlockId:        ids.New().Bytes(),
						BlockHeight:    2,
						BlockTimestamp: timestamppb.Now(),
						Events: []*entities.Event{
							eventCMsg,
							eventDMsg,
						},
					},
				},
			}

			rpc.On("GetEventsForHeightRange", ctx, mock.Anything).Return(response, nil)

			blocks, err := c.GetEventsForHeightRange(ctx, EventRangeQuery{
				Type:        "foo",
				StartHeight: 1,
				EndHeight:   10,
			})
			require.NoError(t, err)

			// Force evaluation of type ID, which is cached in type.
			// Necessary for equality check below
			for _, block := range blocks {
				for _, event := range block.Events {
					_ = event.Value.Type().ID()
				}
			}

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

	t.Run("Internal error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		rpc.On("GetEventsForHeightRange", ctx, mock.Anything).
			Return(nil, errInternal)

		blocks, err := c.GetEventsForHeightRange(ctx, EventRangeQuery{
			Type:        "foo",
			StartHeight: 1,
			EndHeight:   10,
		})

		assert.Error(t, err)
		assert.Equal(t, codes.Internal, status.Code(err))
		assert.Empty(t, blocks)
	}))
}

func TestClient_GetEventsForBlockIDs(t *testing.T) {
	ids := test.IdentifierGenerator()

	t.Run(
		"Empty result",
		clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
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
		"Non-empty result with ccf encoding",
		clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
			events := test.EventGenerator(flow.EventEncodingVersionCCF)
			blockIDA, blockIDB := ids.New(), ids.New()
			eventA, eventB, eventC, eventD := events.New(), events.New(), events.New(), events.New()

			eventAMsg, _ := convert.EventToMessage(eventA, flow.EventEncodingVersionCCF)
			eventBMsg, _ := convert.EventToMessage(eventB, flow.EventEncodingVersionCCF)
			eventCMsg, _ := convert.EventToMessage(eventC, flow.EventEncodingVersionCCF)
			eventDMsg, _ := convert.EventToMessage(eventD, flow.EventEncodingVersionCCF)

			response := &access.EventsResponse{
				Results: []*access.EventsResponse_Result{
					{
						BlockId:        blockIDA.Bytes(),
						BlockHeight:    1,
						BlockTimestamp: timestamppb.Now(),
						Events: []*entities.Event{
							eventAMsg,
							eventBMsg,
						},
					},
					{
						BlockId:        blockIDB.Bytes(),
						BlockHeight:    2,
						BlockTimestamp: timestamppb.Now(),
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

			// Force evaluation of type ID, which is cached in type.
			// Necessary for equality checks below
			for _, block := range blocks {
				for _, event := range block.Events {
					_ = event.Value.Type().ID()
				}
			}

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

	t.Run(
		"Non-empty result with json encoding",
		clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
			events := test.EventGenerator(flow.EventEncodingVersionJSONCDC)
			blockIDA, blockIDB := ids.New(), ids.New()
			eventA, eventB, eventC, eventD := events.New(), events.New(), events.New(), events.New()

			eventAMsg, _ := convert.EventToMessage(eventA, flow.EventEncodingVersionJSONCDC)
			eventBMsg, _ := convert.EventToMessage(eventB, flow.EventEncodingVersionJSONCDC)
			eventCMsg, _ := convert.EventToMessage(eventC, flow.EventEncodingVersionJSONCDC)
			eventDMsg, _ := convert.EventToMessage(eventD, flow.EventEncodingVersionJSONCDC)

			response := &access.EventsResponse{
				Results: []*access.EventsResponse_Result{
					{
						BlockId:        blockIDA.Bytes(),
						BlockHeight:    1,
						BlockTimestamp: timestamppb.Now(),
						Events: []*entities.Event{
							eventAMsg,
							eventBMsg,
						},
					},
					{
						BlockId:        blockIDB.Bytes(),
						BlockHeight:    2,
						BlockTimestamp: timestamppb.Now(),
						Events: []*entities.Event{
							eventCMsg,
							eventDMsg,
						},
					},
				},
			}

			rpc.On("GetEventsForBlockIDs", ctx, mock.Anything).Return(response, nil)

			c.SetEventEncoding(flow.EventEncodingVersionJSONCDC)
			blocks, err := c.GetEventsForBlockIDs(ctx, "foo", []flow.Identifier{blockIDA, blockIDB})
			require.NoError(t, err)

			// Force evaluation of type ID, which is cached in type.
			// Necessary for equality checks below
			for _, block := range blocks {
				for _, event := range block.Events {
					_ = event.Value.Type().ID()
				}
			}

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

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		blockIDA, blockIDB := ids.New(), ids.New()

		rpc.On("GetEventsForBlockIDs", ctx, mock.Anything).
			Return(nil, errNotFound)

		blocks, err := c.GetEventsForBlockIDs(ctx, "foo", []flow.Identifier{blockIDA, blockIDB})
		assert.Error(t, err)
		assert.Equal(t, codes.NotFound, status.Code(err))
		assert.Empty(t, blocks)
	}))
}

func TestClient_GetLatestProtocolStateSnapshot(t *testing.T) {
	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		expected := &access.ProtocolStateSnapshotResponse{
			SerializedSnapshot: make([]byte, 128),
		}
		_, err := rand.Read(expected.SerializedSnapshot)
		assert.NoError(t, err)

		rpc.On("GetLatestProtocolStateSnapshot", ctx, mock.Anything).Return(expected, nil)

		res, err := c.GetLatestProtocolStateSnapshot(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expected.SerializedSnapshot, res)
	}))

	t.Run("Internal error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		rpc.On("GetLatestProtocolStateSnapshot", ctx, mock.Anything).
			Return(nil, errInternal)

		_, err := c.GetLatestProtocolStateSnapshot(ctx)
		assert.Error(t, err)
		assert.Equal(t, codes.Internal, status.Code(err))
	}))
}

func TestClient_GetExecutionResultForBlockID(t *testing.T) {
	ids := test.IdentifierGenerator()
	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		blockID := ids.New()
		executionResult := &entities.ExecutionResult{
			PreviousResultId: ids.New().Bytes(),
			BlockId:          blockID.Bytes(),
			Chunks: []*entities.Chunk{
				{
					CollectionIndex:      0,
					StartState:           ids.New().Bytes(),
					EventCollection:      ids.New().Bytes(),
					BlockId:              blockID.Bytes(),
					TotalComputationUsed: 22,
					NumberOfTransactions: 33,
					Index:                0,
					EndState:             ids.New().Bytes(),
				},
				{
					CollectionIndex:      1,
					StartState:           ids.New().Bytes(),
					EventCollection:      ids.New().Bytes(),
					BlockId:              blockID.Bytes(),
					TotalComputationUsed: 222,
					NumberOfTransactions: 333,
					Index:                1,
					EndState:             ids.New().Bytes(),
				},
			},
			ServiceEvents: []*entities.ServiceEvent{
				{
					Type:    "serviceEvent",
					Payload: []byte("{\"whatever\":21}"),
				},
			},
		}
		result := &access.ExecutionResultForBlockIDResponse{
			ExecutionResult: executionResult,
		}
		rpc.On("GetExecutionResultForBlockID", ctx, &access.GetExecutionResultForBlockIDRequest{
			BlockId: blockID.Bytes(),
		}).Return(result, nil)

		res, err := c.GetExecutionResultForBlockID(ctx, blockID)
		assert.NoError(t, err)

		require.NotNil(t, res)

		require.Len(t, res.Chunks, len(executionResult.Chunks))
		require.Len(t, res.ServiceEvents, len(executionResult.ServiceEvents))

		assert.Equal(t, res.BlockID.Bytes(), executionResult.BlockId)
		assert.Equal(t, res.PreviousResultID.Bytes(), executionResult.PreviousResultId)

		for i, chunk := range res.Chunks {
			assert.Equal(t, chunk.BlockID[:], executionResult.Chunks[i].BlockId)
			assert.Equal(t, chunk.Index, executionResult.Chunks[i].Index)
			assert.Equal(t, uint32(chunk.CollectionIndex), executionResult.Chunks[i].CollectionIndex)
			assert.Equal(t, chunk.StartState[:], executionResult.Chunks[i].StartState)
			assert.Equal(t, []byte(chunk.EventCollection), executionResult.Chunks[i].EventCollection)
			assert.Equal(t, chunk.TotalComputationUsed, executionResult.Chunks[i].TotalComputationUsed)
			assert.Equal(t, uint32(chunk.NumberOfTransactions), executionResult.Chunks[i].NumberOfTransactions)
			assert.Equal(t, chunk.EndState[:], executionResult.Chunks[i].EndState)
		}

		for i, serviceEvent := range res.ServiceEvents {
			assert.Equal(t, serviceEvent.Type, executionResult.ServiceEvents[i].Type)
			assert.Equal(t, serviceEvent.Payload, executionResult.ServiceEvents[i].Payload)
		}

	}))

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		blockId := ids.New()

		rpc.On("GetExecutionResultForBlockID", ctx, &access.GetExecutionResultForBlockIDRequest{
			BlockId: blockId.Bytes(),
		}).Return(nil, errNotFound)

		res, err := c.GetExecutionResultForBlockID(ctx, blockId)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, codes.NotFound, status.Code(err))
	}))
}

func TestClient_GetExecutionResultByID(t *testing.T) {
	ids := test.IdentifierGenerator()
	t.Run("Success", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		blockID := ids.New()
		executionResult := &entities.ExecutionResult{
			PreviousResultId: ids.New().Bytes(),
			BlockId:          blockID.Bytes(),
			Chunks: []*entities.Chunk{
				{
					CollectionIndex:      0,
					StartState:           ids.New().Bytes(),
					EventCollection:      ids.New().Bytes(),
					BlockId:              blockID.Bytes(),
					TotalComputationUsed: 22,
					NumberOfTransactions: 33,
					Index:                0,
					EndState:             ids.New().Bytes(),
				},
				{
					CollectionIndex:      1,
					StartState:           ids.New().Bytes(),
					EventCollection:      ids.New().Bytes(),
					BlockId:              blockID.Bytes(),
					TotalComputationUsed: 222,
					NumberOfTransactions: 333,
					Index:                1,
					EndState:             ids.New().Bytes(),
				},
			},
			ServiceEvents: []*entities.ServiceEvent{
				{
					Type:    "serviceEvent",
					Payload: []byte("{\"whatever\":21}"),
				},
			},
		}
		result := &access.ExecutionResultByIDResponse{
			ExecutionResult: executionResult,
		}
		rpc.On("GetExecutionResultByID", ctx, &access.GetExecutionResultByIDRequest{
			Id: blockID.Bytes(),
		}).Return(result, nil)

		res, err := c.GetExecutionResultByID(ctx, blockID)
		assert.NoError(t, err)

		require.NotNil(t, res)

		require.Len(t, res.Chunks, len(executionResult.Chunks))
		require.Len(t, res.ServiceEvents, len(executionResult.ServiceEvents))

		assert.Equal(t, res.BlockID.Bytes(), executionResult.BlockId)
		assert.Equal(t, res.PreviousResultID.Bytes(), executionResult.PreviousResultId)

		for i, chunk := range res.Chunks {
			assert.Equal(t, chunk.BlockID[:], executionResult.Chunks[i].BlockId)
			assert.Equal(t, chunk.Index, executionResult.Chunks[i].Index)
			assert.Equal(t, uint32(chunk.CollectionIndex), executionResult.Chunks[i].CollectionIndex)
			assert.Equal(t, chunk.StartState[:], executionResult.Chunks[i].StartState)
			assert.Equal(t, []byte(chunk.EventCollection), executionResult.Chunks[i].EventCollection)
			assert.Equal(t, chunk.TotalComputationUsed, executionResult.Chunks[i].TotalComputationUsed)
			assert.Equal(t, uint32(chunk.NumberOfTransactions), executionResult.Chunks[i].NumberOfTransactions)
			assert.Equal(t, chunk.EndState[:], executionResult.Chunks[i].EndState)
		}

		for i, serviceEvent := range res.ServiceEvents {
			assert.Equal(t, serviceEvent.Type, executionResult.ServiceEvents[i].Type)
			assert.Equal(t, serviceEvent.Payload, executionResult.ServiceEvents[i].Payload)
		}

	}))

	t.Run("Not found error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		id := ids.New()

		rpc.On("GetExecutionResultByID", ctx, &access.GetExecutionResultByIDRequest{
			Id: id.Bytes(),
		}).Return(nil, errNotFound)

		res, err := c.GetExecutionResultByID(ctx, id)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, codes.NotFound, status.Code(err))
	}))
}

func TestClient_SubscribeExecutionData(t *testing.T) {
	ids := test.IdentifierGenerator()

	t.Run("Happy Path - by height", executionDataClientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockExecutionDataRPCClient, c *BaseClient) {
		responseCount := uint64(1000)
		startHeight := uint64(10)

		req := executiondata.SubscribeExecutionDataRequest{
			StartBlockHeight:     startHeight,
			EventEncodingVersion: entities.EventEncodingVersion_CCF_V0,
		}

		ctx, cancel := context.WithCancel(ctx)
		stream := &mockClientStream[executiondata.SubscribeExecutionDataResponse]{ctx: ctx}
		for i := startHeight; i < startHeight+responseCount; i++ {
			stream.responses = append(stream.responses, generateExecutionDataResponse(t, ids.New(), i))
		}

		rpc.On("SubscribeExecutionData", ctx, mock.Anything).
			Return(stream, nil).
			Run(assertSubscribeExecutionDataArgs(t, &req))

		eventCh, errCh, err := c.SubscribeExecutionDataByBlockHeight(ctx, startHeight)
		require.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go assertNoErrors(t, errCh, wg.Done)

		i := 0
		for response := range eventCh {
			assert.Equal(t, stream.responses[i].BlockHeight, response.Height)
			assert.Equal(t, stream.responses[i].BlockExecutionData.BlockId[:], response.ExecutionData.BlockID[:])
			assert.Equal(t, stream.responses[i].BlockTimestamp.AsTime(), response.BlockTimestamp)
			i++
			if i == len(stream.responses) {
				cancel()
				break
			}
		}

		wg.Wait()
	}))

	t.Run("Happy Path - by block ID", executionDataClientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockExecutionDataRPCClient, c *BaseClient) {
		responseCount := uint64(1000)
		startBlockID := ids.New()
		startHeight := uint64(10)

		req := executiondata.SubscribeExecutionDataRequest{
			StartBlockId:         startBlockID[:],
			EventEncodingVersion: entities.EventEncodingVersion_CCF_V0,
		}

		ctx, cancel := context.WithCancel(ctx)
		stream := &mockClientStream[executiondata.SubscribeExecutionDataResponse]{ctx: ctx}
		for i := startHeight; i < startHeight+responseCount; i++ {
			stream.responses = append(stream.responses, generateExecutionDataResponse(t, ids.New(), i))
		}

		rpc.On("SubscribeExecutionData", ctx, mock.Anything).
			Return(stream, nil).
			Run(assertSubscribeExecutionDataArgs(t, &req))

		eventCh, errCh, err := c.SubscribeExecutionDataByBlockID(ctx, startBlockID)
		require.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go assertNoErrors(t, errCh, wg.Done)

		i := 0
		for response := range eventCh {
			assert.Equal(t, stream.responses[i].BlockHeight, response.Height)
			assert.Equal(t, stream.responses[i].BlockExecutionData.BlockId[:], response.ExecutionData.BlockID[:])
			assert.Equal(t, stream.responses[i].BlockTimestamp.AsTime(), response.BlockTimestamp)
			i++
			if i == len(stream.responses) {
				cancel()
				break
			}
		}

		wg.Wait()
	}))

	// Test that SubscribeExecutionData returns an error and closes the subscription if the stream returns an error
	t.Run("Returns error from stream", executionDataClientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockExecutionDataRPCClient, c *BaseClient) {
		startHeight := uint64(10)

		req := executiondata.SubscribeExecutionDataRequest{
			StartBlockHeight:     startHeight,
			EventEncodingVersion: entities.EventEncodingVersion_CCF_V0,
		}

		stream := &mockClientStream[executiondata.SubscribeExecutionDataResponse]{
			err: status.Error(codes.Internal, "internal error"),
		}

		rpc.On("SubscribeExecutionData", ctx, mock.Anything).
			Return(stream, nil).
			Run(assertSubscribeExecutionDataArgs(t, &req))

		eventCh, errCh, err := c.SubscribeExecutionDataByBlockHeight(ctx, startHeight)
		require.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go assertNoData(t, eventCh, wg.Done, "execution data")

		i := 0
		for err := range errCh {
			assert.ErrorIs(t, err, stream.err)
			i++
			if i > 1 {
				t.Fatal("should only receive one error")
			}
		}

		wg.Wait()
	}))

	// Test that SubscribeExecutionData returns an error and closes the subscription if an error is encountered while converting the response
	t.Run("Returns error from convert", executionDataClientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockExecutionDataRPCClient, c *BaseClient) {
		startHeight := uint64(10)

		req := executiondata.SubscribeExecutionDataRequest{
			StartBlockHeight:     startHeight,
			EventEncodingVersion: entities.EventEncodingVersion_CCF_V0,
		}

		stream := &mockClientStream[executiondata.SubscribeExecutionDataResponse]{ctx: ctx}
		stream.responses = append(stream.responses, &executiondata.SubscribeExecutionDataResponse{
			BlockHeight:        startHeight,
			BlockExecutionData: nil, // nil BlockExecutionData should cause an error
		})

		rpc.On("SubscribeExecutionData", ctx, mock.Anything).
			Return(stream, nil).
			Run(assertSubscribeExecutionDataArgs(t, &req))

		eventCh, errCh, err := c.SubscribeExecutionDataByBlockHeight(ctx, startHeight)
		require.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go assertNoData(t, eventCh, wg.Done, "execution data")

		i := 0
		for err := range errCh {
			assert.Error(t, err, stream.err)
			i++
			if i > 1 {
				t.Fatal("should only receive one error")
			}
		}

		wg.Wait()
	}))
}

func TestClient_SubscribeEvents(t *testing.T) {
	ids := test.IdentifierGenerator()
	events := test.EventGenerator(flow.EventEncodingVersionCCF)
	addresses := test.AddressGenerator()

	getEvents := func(count int) []flow.Event {
		res := make([]flow.Event, count)
		for i := 0; i < count; i++ {
			res[i] = events.New()
		}
		return res
	}

	t.Run("Happy Path - by height", executionDataClientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockExecutionDataRPCClient, c *BaseClient) {
		responseCount := uint64(1000)
		startHeight := uint64(10)
		filter := flow.EventFilter{
			EventTypes: []string{events.New().Type, events.New().Type},
			Addresses:  []string{addresses.New().String(), addresses.New().String()},
			Contracts:  []string{"A.0.B", "A.1.C"},
		}

		req := executiondata.SubscribeEventsRequest{
			Filter: &executiondata.EventFilter{
				EventType: filter.EventTypes,
				Address:   filter.Addresses,
				Contract:  filter.Contracts,
			},
			EventEncodingVersion: entities.EventEncodingVersion_CCF_V0,
			HeartbeatInterval:    1234,
			StartBlockHeight:     startHeight,
		}

		ctx, cancel := context.WithCancel(ctx)
		stream := &mockClientStream[executiondata.SubscribeEventsResponse]{ctx: ctx}
		for i := startHeight; i < startHeight+responseCount; i++ {
			stream.responses = append(stream.responses, generateEventResponse(t, ids.New(), i, getEvents(2), flow.EventEncodingVersionCCF))
		}

		rpc.On("SubscribeEvents", ctx, mock.Anything).
			Return(stream, nil).
			Run(assertSubscribeEventsArgs(t, &req))

		eventCh, errCh, err := c.SubscribeEventsByBlockHeight(ctx, startHeight, filter, WithHeartbeatInterval(req.HeartbeatInterval))
		require.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go assertNoErrors(t, errCh, wg.Done)

		i := 0
		for response := range eventCh {
			assert.Equal(t, stream.responses[i].BlockHeight, response.Height)
			assert.Equal(t, stream.responses[i].BlockId, response.BlockID[:])
			assert.Equal(t, len(stream.responses[i].Events), len(response.Events))
			i++
			if i == len(stream.responses) {
				cancel()
				break
			}
		}

		wg.Wait()
	}))

	t.Run("Happy Path - by block ID", executionDataClientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockExecutionDataRPCClient, c *BaseClient) {
		responseCount := uint64(1000)
		startBlockID := ids.New()
		startHeight := uint64(10)
		filter := flow.EventFilter{
			EventTypes: []string{events.New().Type, events.New().Type},
			Addresses:  []string{addresses.New().String(), addresses.New().String()},
			Contracts:  []string{"A.0.B", "A.1.C"},
		}

		req := executiondata.SubscribeEventsRequest{
			Filter: &executiondata.EventFilter{
				EventType: filter.EventTypes,
				Address:   filter.Addresses,
				Contract:  filter.Contracts,
			},
			EventEncodingVersion: entities.EventEncodingVersion_CCF_V0,
			HeartbeatInterval:    1234,
			StartBlockId:         startBlockID[:],
		}

		ctx, cancel := context.WithCancel(ctx)
		stream := &mockClientStream[executiondata.SubscribeEventsResponse]{ctx: ctx}
		for i := startHeight; i < startHeight+responseCount; i++ {
			stream.responses = append(stream.responses, generateEventResponse(t, ids.New(), i, getEvents(2), flow.EventEncodingVersionCCF))
		}

		rpc.On("SubscribeEvents", ctx, mock.Anything).
			Return(stream, nil).
			Run(assertSubscribeEventsArgs(t, &req))

		eventCh, errCh, err := c.SubscribeEventsByBlockID(ctx, startBlockID, filter, WithHeartbeatInterval(req.HeartbeatInterval))
		require.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go assertNoErrors(t, errCh, wg.Done)

		i := 0
		for response := range eventCh {
			assert.Equal(t, stream.responses[i].BlockHeight, response.Height)
			assert.Equal(t, stream.responses[i].BlockId, response.BlockID[:])
			assert.Equal(t, len(stream.responses[i].Events), len(response.Events))
			i++
			if i == len(stream.responses) {
				cancel()
				break
			}
		}

		wg.Wait()
	}))

	// Test that SubscribeEvents returns an error and closes the subscription if the stream returns an error
	t.Run("Returns error from stream", executionDataClientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockExecutionDataRPCClient, c *BaseClient) {
		startHeight := uint64(10)
		filter := flow.EventFilter{}

		req := executiondata.SubscribeEventsRequest{
			Filter: &executiondata.EventFilter{
				EventType: filter.EventTypes,
				Address:   filter.Addresses,
				Contract:  filter.Contracts,
			},
			EventEncodingVersion: entities.EventEncodingVersion_CCF_V0,
			HeartbeatInterval:    100,
			StartBlockHeight:     startHeight,
		}

		stream := &mockClientStream[executiondata.SubscribeEventsResponse]{
			err: status.Error(codes.Internal, "internal error"),
		}

		rpc.On("SubscribeEvents", ctx, mock.Anything).
			Return(stream, nil).
			Run(assertSubscribeEventsArgs(t, &req))

		eventCh, errCh, err := c.SubscribeEventsByBlockHeight(ctx, startHeight, filter)
		require.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go assertNoData(t, eventCh, wg.Done, "events")

		i := 0
		for err := range errCh {
			assert.ErrorIs(t, err, stream.err)
			i++
			if i > 1 {
				t.Fatal("should only receive one error")
			}
		}

		wg.Wait()
	}))

	// Test that SubscribeEvents returns an error and closes the subscription if an error is encountered while converting the response
	t.Run("Returns error from convert", executionDataClientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockExecutionDataRPCClient, c *BaseClient) {
		startHeight := uint64(10)
		filter := flow.EventFilter{}

		req := executiondata.SubscribeEventsRequest{
			Filter: &executiondata.EventFilter{
				EventType: filter.EventTypes,
				Address:   filter.Addresses,
				Contract:  filter.Contracts,
			},
			EventEncodingVersion: entities.EventEncodingVersion_CCF_V0,
			HeartbeatInterval:    100,
			StartBlockHeight:     startHeight,
		}

		stream := &mockClientStream[executiondata.SubscribeEventsResponse]{ctx: ctx}
		stream.responses = append(stream.responses, generateEventResponse(t, ids.New(), startHeight, getEvents(2), flow.EventEncodingVersionCCF))

		// corrupt the event payload
		stream.responses[0].Events[0].Payload[0] = 'x'

		rpc.On("SubscribeEvents", ctx, mock.Anything).
			Return(stream, nil).
			Run(assertSubscribeEventsArgs(t, &req))

		eventCh, errCh, err := c.SubscribeEventsByBlockHeight(ctx, startHeight, filter, WithHeartbeatInterval(req.HeartbeatInterval))
		require.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go assertNoData(t, eventCh, wg.Done, "events")

		i := 0
		for err := range errCh {
			assert.Error(t, err, stream.err)
			i++
			if i > 1 {
				t.Fatal("should only receive one error")
			}
		}

		wg.Wait()
	}))
}

func TestClient_SubscribeAccountStatuses(t *testing.T) {
	filter := flow.AccountStatusFilter{}
	accounts := test.AccountGenerator()
	blocks := test.BlockGenerator()

	generateAccountStatusesResponses := func(count uint64) []*executiondata.SubscribeAccountStatusesResponse {
		var resBlockHeaders []*executiondata.SubscribeAccountStatusesResponse

		for i := uint64(0); i < count; i++ {
			account := convert.AccountToMessage(*accounts.New())
			results := []*executiondata.SubscribeAccountStatusesResponse_Result{
				{
					Address: account.Address,
				},
			}

			block := blocks.New()
			resBlockHeaders = append(resBlockHeaders, &executiondata.SubscribeAccountStatusesResponse{
				BlockId:      block.ID.Bytes(),
				BlockHeight:  block.Height,
				MessageIndex: i,
				Results:      results,
			})
		}

		return resBlockHeaders
	}

	t.Run("Happy Path - from start height", executionDataClientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockExecutionDataRPCClient, c *BaseClient) {
		startHeight := uint64(1)
		responseCount := uint64(100)

		ctx, cancel := context.WithCancel(ctx)
		stream := &mockClientStream[executiondata.SubscribeAccountStatusesResponse]{
			ctx:       ctx,
			responses: generateAccountStatusesResponses(responseCount),
		}

		rpc.
			On("SubscribeAccountStatusesFromStartHeight", ctx, mock.Anything).
			Return(stream, nil)

		accountStatusesCh, errCh, err := c.SubscribeAccountStatusesFromStartHeight(ctx, startHeight, filter)
		require.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go assertNoErrors(t, errCh, wg.Done)

		for i := uint64(0); i < responseCount; i++ {
			actualAccountStatus := <-accountStatusesCh

			blockId := convert.MessageToIdentifier(stream.responses[i].GetBlockId())
			require.Equal(t, blockId, actualAccountStatus.BlockID)

			blockHeight := stream.responses[i].GetBlockHeight()
			require.Equal(t, blockHeight, actualAccountStatus.BlockHeight)

			messageIndex := stream.responses[i].GetMessageIndex()
			require.Equal(t, messageIndex, actualAccountStatus.MessageIndex)

			results, err := convert.MessageToAccountStatusResults(stream.responses[i].GetResults())
			require.NoError(t, err)
			require.Equal(t, results, actualAccountStatus.Results)
		}
		cancel()

		wg.Wait()
	}))

	t.Run("Happy Path - from start block id", executionDataClientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockExecutionDataRPCClient, c *BaseClient) {
		responseCount := uint64(100)

		ctx, cancel := context.WithCancel(ctx)
		stream := &mockClientStream[executiondata.SubscribeAccountStatusesResponse]{
			ctx:       ctx,
			responses: generateAccountStatusesResponses(responseCount),
		}

		rpc.
			On("SubscribeAccountStatusesFromStartBlockID", ctx, mock.Anything).
			Return(stream, nil)

		startBlockId := convert.MessageToIdentifier(stream.responses[0].GetBlockId())
		accountStatusesCh, errCh, err := c.SubscribeAccountStatusesFromStartBlockID(ctx, startBlockId, filter)
		require.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go assertNoErrors(t, errCh, wg.Done)

		for i := uint64(0); i < responseCount; i++ {
			actualAccountStatus := <-accountStatusesCh

			blockId := convert.MessageToIdentifier(stream.responses[i].GetBlockId())
			require.Equal(t, blockId, actualAccountStatus.BlockID)

			blockHeight := stream.responses[i].GetBlockHeight()
			require.Equal(t, blockHeight, actualAccountStatus.BlockHeight)

			messageIndex := stream.responses[i].GetMessageIndex()
			require.Equal(t, messageIndex, actualAccountStatus.MessageIndex)

			results, err := convert.MessageToAccountStatusResults(stream.responses[i].GetResults())
			require.NoError(t, err)
			require.Equal(t, results, actualAccountStatus.Results)
		}
		cancel()

		wg.Wait()
	}))

	t.Run("Happy Path - from start latest block", executionDataClientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockExecutionDataRPCClient, c *BaseClient) {
		responseCount := uint64(100)

		ctx, cancel := context.WithCancel(ctx)
		stream := &mockClientStream[executiondata.SubscribeAccountStatusesResponse]{
			ctx:       ctx,
			responses: generateAccountStatusesResponses(responseCount),
		}

		rpc.
			On("SubscribeAccountStatusesFromLatestBlock", ctx, mock.Anything).
			Return(stream, nil)

		accountStatusesCh, errCh, err := c.SubscribeAccountStatusesFromLatestBlock(ctx, filter)
		require.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go assertNoErrors(t, errCh, wg.Done)

		for i := uint64(0); i < responseCount; i++ {
			actualAccountStatus := <-accountStatusesCh

			blockId := convert.MessageToIdentifier(stream.responses[i].GetBlockId())
			require.Equal(t, blockId, actualAccountStatus.BlockID)

			blockHeight := stream.responses[i].GetBlockHeight()
			require.Equal(t, blockHeight, actualAccountStatus.BlockHeight)

			messageIndex := stream.responses[i].GetMessageIndex()
			require.Equal(t, messageIndex, actualAccountStatus.MessageIndex)

			results, err := convert.MessageToAccountStatusResults(stream.responses[i].GetResults())
			require.NoError(t, err)
			require.Equal(t, results, actualAccountStatus.Results)
		}
		cancel()

		wg.Wait()
	}))

	t.Run("Stream returns error", executionDataClientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockExecutionDataRPCClient, c *BaseClient) {
		ctx, cancel := context.WithCancel(ctx)
		stream := &mockClientStream[executiondata.SubscribeAccountStatusesResponse]{
			ctx: ctx,
			err: status.Error(codes.Internal, "internal error"),
		}

		rpc.
			On("SubscribeAccountStatusesFromLatestBlock", ctx, mock.Anything).
			Return(stream, nil)

		accountStatuses, errCh, err := c.SubscribeAccountStatusesFromLatestBlock(ctx, filter)
		require.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go assertNoData(t, accountStatuses, wg.Done, "account statuses")

		errorCount := 0
		for e := range errCh {
			require.Error(t, e)
			require.ErrorIs(t, e, stream.err)
			errorCount += 1
		}
		cancel()

		require.Equalf(t, 1, errorCount, "only 1 error is expected")

		wg.Wait()
	}))
}

func TestClient_SubscribeBlockHeaders(t *testing.T) {
	blockHeaders := test.BlockHeaderGenerator()

	generateBlockHeaderResponses := func(count uint64) []*access.SubscribeBlockHeadersResponse {
		var resBlockHeaders []*access.SubscribeBlockHeadersResponse

		for i := uint64(0); i < count; i++ {
			header, err := convert.BlockHeaderToMessage(blockHeaders.New())
			require.NoError(t, err)

			resBlockHeaders = append(resBlockHeaders, &access.SubscribeBlockHeadersResponse{
				Header: header,
			})
		}

		return resBlockHeaders
	}

	t.Run("Happy Path - from start height", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		startHeight := uint64(1)
		responseCount := uint64(100)

		ctx, cancel := context.WithCancel(ctx)
		stream := &mockClientStream[access.SubscribeBlockHeadersResponse]{
			ctx:       ctx,
			responses: generateBlockHeaderResponses(responseCount),
		}

		rpc.
			On("SubscribeBlockHeadersFromStartHeight", ctx, mock.Anything).
			Return(stream, nil)

		blockHeadersCh, errCh, err := c.SubscribeBlockHeadersFromStartHeight(ctx, startHeight, flow.BlockStatusUnknown)
		require.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go assertNoErrors(t, errCh, wg.Done)

		for i := uint64(0); i < responseCount; i++ {
			actualHeader := <-blockHeadersCh
			expectedHeader, err := convert.MessageToBlockHeader(stream.responses[i].GetHeader())
			require.NoError(t, err)
			require.Equal(t, expectedHeader, actualHeader)
		}
		cancel()

		wg.Wait()
	}))

	t.Run("Happy Path - from start block id", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		responseCount := uint64(100)

		ctx, cancel := context.WithCancel(ctx)
		stream := &mockClientStream[access.SubscribeBlockHeadersResponse]{
			ctx:       ctx,
			responses: generateBlockHeaderResponses(responseCount),
		}

		rpc.
			On("SubscribeBlockHeadersFromStartBlockID", ctx, mock.Anything).
			Return(stream, nil)

		startBlockID := convert.MessageToIdentifier(stream.responses[0].GetHeader().Id)
		blockHeadersCh, errCh, err := c.SubscribeBlockHeadersFromStartBlockID(ctx, startBlockID, flow.BlockStatusUnknown)
		require.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go assertNoErrors(t, errCh, wg.Done)

		for i := uint64(0); i < responseCount; i++ {
			actualHeader := <-blockHeadersCh
			expectedHeader, err := convert.MessageToBlockHeader(stream.responses[i].GetHeader())
			require.NoError(t, err)
			require.Equal(t, expectedHeader, actualHeader)
		}
		cancel()

		wg.Wait()
	}))

	t.Run("Happy Path - from latest", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		responseCount := uint64(100)

		ctx, cancel := context.WithCancel(ctx)
		stream := &mockClientStream[access.SubscribeBlockHeadersResponse]{
			ctx:       ctx,
			responses: generateBlockHeaderResponses(responseCount),
		}

		rpc.
			On("SubscribeBlockHeadersFromLatest", ctx, mock.Anything).
			Return(stream, nil)

		blockHeadersCh, errCh, err := c.SubscribeBlockHeadersFromLatest(ctx, flow.BlockStatusUnknown)
		require.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go assertNoErrors(t, errCh, wg.Done)

		for i := uint64(0); i < responseCount; i++ {
			actualHeader := <-blockHeadersCh
			expectedHeader, err := convert.MessageToBlockHeader(stream.responses[i].GetHeader())
			require.NoError(t, err)
			require.Equal(t, expectedHeader, actualHeader)
		}
		cancel()

		wg.Wait()
	}))

	t.Run("Stream returns error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		ctx, cancel := context.WithCancel(ctx)
		stream := &mockClientStream[access.SubscribeBlockHeadersResponse]{
			ctx: ctx,
			err: status.Error(codes.Internal, "internal error"),
		}

		rpc.
			On("SubscribeBlockHeadersFromLatest", ctx, mock.Anything).
			Return(stream, nil)

		blockHeadersCh, errCh, err := c.SubscribeBlockHeadersFromLatest(ctx, flow.BlockStatusUnknown)
		require.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go assertNoData(t, blockHeadersCh, wg.Done, "block headers")

		errorCount := 0
		for e := range errCh {
			require.Error(t, e)
			require.ErrorIs(t, e, stream.err)
			errorCount += 1
		}
		cancel()

		require.Equalf(t, 1, errorCount, "only 1 error is expected")

		wg.Wait()
	}))
}

func TestClient_SubscribeBlocks(t *testing.T) {
	blocks := test.BlockGenerator()

	generateBlockResponses := func(count uint64) []*access.SubscribeBlocksResponse {
		var resBlocks []*access.SubscribeBlocksResponse

		for i := uint64(0); i < count; i++ {
			b, err := convert.BlockToMessage(*blocks.New())
			require.NoError(t, err)

			resBlocks = append(resBlocks, &access.SubscribeBlocksResponse{
				Block: b,
			})
		}

		return resBlocks
	}

	t.Run("Happy Path - from start height", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		startHeight := uint64(1)
		responseCount := uint64(100)

		ctx, cancel := context.WithCancel(ctx)
		stream := &mockClientStream[access.SubscribeBlocksResponse]{
			ctx:       ctx,
			responses: generateBlockResponses(responseCount),
		}

		rpc.
			On("SubscribeBlocksFromStartHeight", ctx, mock.Anything).
			Return(stream, nil)

		blockCh, errCh, err := c.SubscribeBlocksFromStartHeight(ctx, startHeight, flow.BlockStatusUnknown)
		require.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go assertNoErrors(t, errCh, wg.Done)

		for i := uint64(0); i < responseCount; i++ {
			actualBlock := <-blockCh
			expectedBlock, err := convert.MessageToBlock(stream.responses[i].GetBlock())
			require.NoError(t, err)
			require.Equal(t, expectedBlock, actualBlock)
		}
		cancel()

		wg.Wait()
	}))

	t.Run("Happy Path - from start block id", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		responseCount := uint64(100)

		ctx, cancel := context.WithCancel(ctx)
		stream := &mockClientStream[access.SubscribeBlocksResponse]{
			ctx:       ctx,
			responses: generateBlockResponses(responseCount),
		}

		rpc.
			On("SubscribeBlocksFromStartBlockID", ctx, mock.Anything).
			Return(stream, nil)

		startBlockID := convert.MessageToIdentifier(stream.responses[0].Block.Id)
		blockCh, errCh, err := c.SubscribeBlocksFromStartBlockID(ctx, startBlockID, flow.BlockStatusUnknown)
		require.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go assertNoErrors(t, errCh, wg.Done)

		for i := uint64(0); i < responseCount; i++ {
			actualBlock := <-blockCh
			expectedBlock, err := convert.MessageToBlock(stream.responses[i].GetBlock())
			require.NoError(t, err)
			require.Equal(t, expectedBlock, actualBlock)
		}
		cancel()

		wg.Wait()
	}))

	t.Run("Happy Path - from latest", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		responseCount := uint64(100)

		ctx, cancel := context.WithCancel(ctx)
		stream := &mockClientStream[access.SubscribeBlocksResponse]{
			ctx:       ctx,
			responses: generateBlockResponses(responseCount),
		}

		rpc.
			On("SubscribeBlocksFromLatest", ctx, mock.Anything).
			Return(stream, nil)

		blockCh, errCh, err := c.SubscribeBlocksFromLatest(ctx, flow.BlockStatusUnknown)
		require.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go assertNoErrors(t, errCh, wg.Done)

		for i := uint64(0); i < responseCount; i++ {
			actualBlock := <-blockCh
			expectedBlock, err := convert.MessageToBlock(stream.responses[i].GetBlock())
			require.NoError(t, err)
			require.Equal(t, expectedBlock, actualBlock)
		}
		cancel()

		wg.Wait()
	}))

	t.Run("Stream returns error", clientTest(func(t *testing.T, ctx context.Context, rpc *mocks.MockRPCClient, c *BaseClient) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		stream := &mockClientStream[access.SubscribeBlocksResponse]{
			ctx: ctx,
			err: status.Error(codes.Internal, "internal error"),
		}

		rpc.
			On("SubscribeBlocksFromLatest", ctx, mock.Anything).
			Return(stream, nil)

		blockCh, errCh, err := c.SubscribeBlocksFromLatest(ctx, flow.BlockStatusUnknown)
		require.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(1)
		go assertNoData(t, blockCh, wg.Done, "blocks")

		errorCount := 0
		for e := range errCh {
			require.Error(t, e)
			require.ErrorIs(t, e, stream.err)
			errorCount += 1
		}
		require.Equalf(t, 1, errorCount, "only 1 error is expected")

		wg.Wait()
	}))
}

func generateEventResponse(t *testing.T, blockID flow.Identifier, height uint64, events []flow.Event, encoding flow.EventEncodingVersion) *executiondata.SubscribeEventsResponse {
	responseEvents := make([]*entities.Event, 0, len(events))
	for _, e := range events {
		eventMsg, err := convert.EventToMessage(e, encoding)
		require.NoError(t, err)
		responseEvents = append(responseEvents, eventMsg)
	}

	return &executiondata.SubscribeEventsResponse{
		BlockHeight: height,
		BlockId:     blockID[:],
		Events:      responseEvents,
	}
}

func generateExecutionDataResponse(t *testing.T, blockID flow.Identifier, height uint64) *executiondata.SubscribeExecutionDataResponse {
	return &executiondata.SubscribeExecutionDataResponse{
		BlockHeight: height,
		BlockExecutionData: &entities.BlockExecutionData{
			BlockId:            blockID[:],
			ChunkExecutionData: []*entities.ChunkExecutionData{},
		},
		BlockTimestamp: timestamppb.Now(),
	}
}

func assertSubscribeEventsArgs(t *testing.T, expected *executiondata.SubscribeEventsRequest) func(args mock.Arguments) {
	return func(args mock.Arguments) {
		actual, ok := args.Get(1).(*executiondata.SubscribeEventsRequest)
		require.True(t, ok)

		assert.Equal(t, expected.Filter, actual.Filter)
		assert.Equal(t, expected.EventEncodingVersion, actual.EventEncodingVersion)
		assert.Equal(t, expected.HeartbeatInterval, actual.HeartbeatInterval)
		assert.Equal(t, expected.StartBlockHeight, actual.StartBlockHeight)
		assert.Equal(t, expected.StartBlockId, actual.StartBlockId)
	}
}

func assertSubscribeExecutionDataArgs(t *testing.T, expected *executiondata.SubscribeExecutionDataRequest) func(args mock.Arguments) {
	return func(args mock.Arguments) {
		actual, ok := args.Get(1).(*executiondata.SubscribeExecutionDataRequest)
		require.True(t, ok)

		assert.Equal(t, expected.EventEncodingVersion, actual.EventEncodingVersion)
		assert.Equal(t, expected.StartBlockHeight, actual.StartBlockHeight)
		assert.Equal(t, expected.StartBlockId, actual.StartBlockId)
	}
}

func assertNoErrors(t *testing.T, errCh <-chan error, done func()) {
	defer done()
	for err := range errCh {
		require.NoError(t, err)
	}
}

func assertNoData[T any](t *testing.T, dataChan <-chan T, done func(), topicNameForError string) {
	defer done()
	for range dataChan {
		t.Fatalf("should not receive %s", topicNameForError)
	}
}

type mockClientStream[Response any] struct {
	grpc.ClientStream

	ctx       context.Context
	err       error
	offset    int
	responses []*Response
}

func (s *mockClientStream[Response]) Recv() (*Response, error) {
	if s.err != nil {
		return nil, s.err
	}

	if s.offset >= len(s.responses) {
		<-s.ctx.Done()
		return nil, io.EOF
	}
	defer func() { s.offset++ }()

	return s.responses[s.offset], nil
}

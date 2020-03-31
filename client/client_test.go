package client_test

//
// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"testing"
//
// 	"github.com/dapperlabs/cadence"
// 	"github.com/dapperlabs/flow-go/crypto"
// 	"github.com/dapperlabs/flow-go/protobuf/sdk/entities"
// 	"github.com/dapperlabs/flow-go/protobuf/services/observation"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
//
// 	"github.com/dapperlabs/cadence/encoding"
//
// 	"github.com/dapperlabs/flow-go-sdk"
// 	"github.com/dapperlabs/flow-go-sdk/client"
// 	"github.com/dapperlabs/flow-go-sdk/client/mocks"
// 	"github.com/dapperlabs/flow-go-sdk/client/protobuf/convert"
// 	"github.com/dapperlabs/flow-go-sdk/utils/unittest"
// )
//
// func TestPing(t *testing.T) {
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()
//
// 	mockRPC := mocks.NewMockRPCClient(mockCtrl)
//
// 	c := client.NewFromRPCClient(mockRPC)
// 	ctx := context.Background()
//
// 	t.Run("Success", func(t *testing.T) {
// 		mockRPC.EXPECT().
// 			Ping(ctx, gomock.Any()).
// 			Return(&observation.PingResponse{}, nil).
// 			Times(1)
//
// 		err := c.Ping(ctx)
// 		assert.NoError(t, err)
// 	})
//
// 	t.Run("ServerError", func(t *testing.T) {
// 		mockRPC.EXPECT().
// 			Ping(ctx, gomock.Any()).
// 			Return(nil, fmt.Errorf("fake error")).
// 			Times(1)
//
// 		err := c.Ping(ctx)
// 		assert.Error(t, err)
// 	})
// }
//
// func TestSendTransaction(t *testing.T) {
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()
//
// 	mockRPC := mocks.NewMockRPCClient(mockCtrl)
//
// 	c := client.NewFromRPCClient(mockRPC)
// 	ctx := context.Background()
//
// 	tx := unittest.TransactionFixture()
//
// 	t.Run("Success", func(t *testing.T) {
// 		// client should return non-error if RPC call succeeds
// 		mockRPC.EXPECT().
// 			SendTransaction(ctx, gomock.Any()).
// 			Return(&observation.SendTransactionResponse{Hash: tx.Hash()}, nil).
// 			Times(1)
//
// 		err := c.SendTransaction(ctx, tx)
// 		assert.NoError(t, err)
// 	})
//
// 	t.Run("Server error", func(t *testing.T) {
// 		// client should return error if RPC call fails
// 		mockRPC.EXPECT().
// 			SendTransaction(ctx, gomock.Any()).
// 			Return(nil, errors.New("dummy error")).
// 			Times(1)
//
// 		// error should be passed to user
// 		err := c.SendTransaction(ctx, tx)
// 		assert.Error(t, err)
// 	})
// }
//
// func TestGetLatestBlock(t *testing.T) {
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()
//
// 	mockRPC := mocks.NewMockRPCClient(mockCtrl)
//
// 	c := client.NewFromRPCClient(mockRPC)
// 	ctx := context.Background()
//
// 	res := &observation.GetLatestBlockResponse{
// 		Block: convert.BlockHeaderToMessage(unittest.BlockHeaderFixture()),
// 	}
//
// 	t.Run("Success", func(t *testing.T) {
// 		// client should return non-error if RPC call succeeds
// 		mockRPC.EXPECT().
// 			GetLatestBlock(ctx, gomock.Any()).
// 			Return(res, nil).
// 			Times(1)
//
// 		blockHeaderA, err := c.GetLatestBlock(ctx, true)
// 		assert.NoError(t, err)
//
// 		blockHeaderB := convert.MessageToBlockHeader(res.GetBlock())
// 		assert.Equal(t, *blockHeaderA, blockHeaderB)
// 	})
//
// 	t.Run("Server error", func(t *testing.T) {
// 		// client should return error if RPC call fails
// 		mockRPC.EXPECT().
// 			GetLatestBlock(ctx, gomock.Any()).
// 			Return(nil, errors.New("dummy error")).
// 			Times(1)
//
// 		_, err := c.GetLatestBlock(ctx, true)
// 		assert.Error(t, err)
// 	})
// }
//
// func TestExecuteScript(t *testing.T) {
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()
//
// 	mockRPC := mocks.NewMockRPCClient(mockCtrl)
//
// 	c := client.NewFromRPCClient(mockRPC)
// 	ctx := context.Background()
//
// 	value := cadence.NewInt(42)
// 	valueBytes, err := encoding.Encode(value)
// 	require.NoError(t, err)
//
// 	t.Run("Success", func(t *testing.T) {
// 		// client should return non-error if RPC call succeeds
// 		mockRPC.EXPECT().
// 			ExecuteScript(ctx, gomock.Any()).
// 			Return(&observation.ExecuteScriptResponse{Value: valueBytes}, nil).
// 			Times(1)
//
// 		b, err := c.ExecuteScript(ctx, []byte("pub fun main(): Int { return 1 }"))
// 		assert.NoError(t, err)
//
// 		value, err := encoding.Decode(cadence.IntType{}, b)
// 		assert.NoError(t, err)
//
// 		assert.Equal(t, cadence.NewInt(42), value)
// 	})
//
// 	t.Run("Server error", func(t *testing.T) {
// 		// client should return error if RPC call fails
// 		mockRPC.EXPECT().
// 			ExecuteScript(ctx, gomock.Any()).
// 			Return(nil, errors.New("dummy error")).
// 			Times(1)
//
// 		// error should be passed to user
// 		_, err := c.ExecuteScript(ctx, []byte("pub fun main(): Int { return 1 }"))
// 		assert.Error(t, err)
// 	})
// }
//
// func TestGetTransaction(t *testing.T) {
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()
//
// 	mockRPC := mocks.NewMockRPCClient(mockCtrl)
//
// 	c := client.NewFromRPCClient(mockRPC)
// 	ctx := context.Background()
//
// 	tx := unittest.TransactionFixture()
//
// 	events := []flow.Event{unittest.EventFixture()}
//
// 	eventMessages := make([]*entities.Event, len(events))
// 	for i, event := range events {
// 		eventMessages[i] = convert.EventToMessage(event)
// 	}
//
// 	t.Run("Success", func(t *testing.T) {
// 		mockRPC.EXPECT().
// 			GetTransaction(ctx, gomock.Any()).
// 			Return(&observation.GetTransactionResponse{
// 				Transaction: convert.TransactionToMessage(tx),
// 				Events:      eventMessages,
// 			}, nil).
// 			Times(1)
//
// 		res, err := c.GetTransaction(ctx, crypto.Hash{})
// 		assert.NoError(t, err)
// 		assert.Len(t, res.Events, 1)
// 		assert.Equal(t, events[0].Type, res.Events[0].Type)
// 	})
//
// 	t.Run("Server error", func(t *testing.T) {
// 		mockRPC.EXPECT().
// 			GetTransaction(ctx, gomock.Any()).
// 			Return(nil, fmt.Errorf("dummy error")).
// 			Times(1)
//
// 		// The client should pass along the error
// 		_, err := c.GetTransaction(ctx, crypto.Hash{})
// 		assert.Error(t, err)
// 	})
// }
//
// func TestGetEvents(t *testing.T) {
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()
//
// 	mockRPC := mocks.NewMockRPCClient(mockCtrl)
//
// 	c := client.NewFromRPCClient(mockRPC)
// 	ctx := context.Background()
//
// 	// declare event type used for decoding event payloads
// 	mockEventType := cadence.EventType{
// 		CompositeType: cadence.CompositeType{
// 			Identifier: "Transfer",
// 			Fields: []cadence.Field{
// 				{
// 					Identifier: "to",
// 					Type:       cadence.AddressType{},
// 				},
// 				{
// 					Identifier: "from",
// 					Type:       cadence.AddressType{},
// 				},
// 				{
// 					Identifier: "amount",
// 					Type:       cadence.IntType{},
// 				},
// 			},
// 		},
// 	}
//
// 	to := cadence.Address(flow.ZeroAddress)
// 	from := cadence.Address(flow.ZeroAddress)
// 	amount := cadence.NewInt(42)
//
// 	mockEventValue := cadence.
// 		NewComposite([]cadence.Value{to, from, amount}).
// 		WithType(mockEventType)
//
// 	// encode event payload from mock value
// 	eventPayload, _ := encoding.Encode(mockEventValue)
//
// 	// Set up a mock event response
// 	mockEvent := flow.Event{
// 		Type:    "Transfer",
// 		Payload: eventPayload,
// 	}
//
// 	t.Run("Success", func(t *testing.T) {
// 		// Set up the mock to return a mocked event response
// 		mockRes := &observation.GetEventsResponse{Events: []*entities.Event{
// 			convert.EventToMessage(mockEvent),
// 		}}
//
// 		mockRPC.EXPECT().
// 			GetEvents(ctx, gomock.Any()).
// 			Return(mockRes, nil).
// 			Times(1)
//
// 		// The client should pass the response to the client
// 		events, err := c.GetEvents(ctx, client.EventQuery{})
// 		assert.Nil(t, err)
// 		require.Len(t, events, 1)
//
// 		actualEvent := events[0]
//
// 		value, err := encoding.Decode(mockEventType, actualEvent.Payload)
// 		eventValue := value.(cadence.Composite)
//
// 		assert.Equal(t, actualEvent.Type, mockEvent.Type)
// 		assert.Equal(t, to, eventValue.Fields[0])
// 		assert.Equal(t, from, eventValue.Fields[1])
// 		assert.Equal(t, amount, eventValue.Fields[2])
// 	})
//
// 	t.Run("Server error", func(t *testing.T) {
// 		// Set up the mock to return an error
// 		mockRPC.EXPECT().
// 			GetEvents(ctx, gomock.Any()).
// 			Return(nil, fmt.Errorf("dummy error")).
// 			Times(1)
//
// 		// The client should pass along the error
// 		_, err := c.GetEvents(ctx, client.EventQuery{})
// 		assert.Error(t, err)
// 	})
// }

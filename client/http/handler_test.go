package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/onflow/flow-go-sdk/test"

	"github.com/onflow/flow-go/engine/access/rest/models"

	"github.com/stretchr/testify/assert"
)

// handlerTest is a helper that builds handler with a http test server
// and exposes a referenced test request instance which can be used to set values to test.
func handlerTest(f func(ctx context.Context, t *testing.T, handler httpHandler, req *testRequest)) func(t *testing.T) {
	testReq := &testRequest{}
	return func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			assert.Equal(t, request.URL.String(), testReq.url.String())

			var err error
			if testReq.err != nil {
				writer.WriteHeader(http.StatusBadRequest)
				_, err = writer.Write(testReq.err)
			} else {
				writer.WriteHeader(http.StatusOK)
				_, err = writer.Write(testReq.res)
			}
			assert.NoError(t, err)
		}))
		defer server.Close()

		h := httpHandler{
			client: server.Client(),
			base:   server.URL,
			debug:  false,
		}

		f(context.Background(), t, h, testReq)
	}
}

type testRequest struct {
	url url.URL
	res []byte
	err []byte
}

// setData set url and response data.
func (t *testRequest) SetData(url url.URL, res interface{}) {
	t.url = url
	bytes, _ := json.Marshal(res)
	t.res = bytes
}

func (t *testRequest) SetErr(url url.URL, err interface{}) {
	t.url = url
	bytes, _ := json.Marshal(err)
	t.err = bytes
}

// addQuery adds query parameters from a map to URL.
func addQuery(u *url.URL, q map[string]string) url.URL {
	query := u.Query()
	for key, value := range q {
		query.Add(key, value)
	}
	u.RawQuery = query.Encode()
	return *u
}

// newBlocksURL is a helper factory for building blocks URLs.
func newBlocksURL(query map[string]string) url.URL {
	u, _ := url.Parse("/blocks")
	if query == nil {
		query = map[string]string{}
	}
	query["expand"] = "payload"

	return addQuery(u, query)
}

func TestHandler_GetBlockByID(t *testing.T) {
	t.Run("Success", handlerTest(func(ctx context.Context, t *testing.T, handler httpHandler, req *testRequest) {
		b := test.BlockHTTP()
		httpBlock := []*models.Block{&b}

		const id = "0x1"
		blockURL := newBlocksURL(nil)

		blockURL.Path = fmt.Sprintf("%s/%s", blockURL.Path, id)
		req.SetData(blockURL, httpBlock)

		block, err := handler.getBlockByID(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, block, httpBlock[0])
	}))

	t.Run("Failed", handlerTest(func(ctx context.Context, t *testing.T, handler httpHandler, req *testRequest) {
		const id = "0x1"
		blockURL := newBlocksURL(nil)

		blockURL.Path = fmt.Sprintf("%s/%s", blockURL.Path, id)
		req.SetData(blockURL, []models.Block{})

		_, err := handler.getBlockByID(ctx, id)
		assert.EqualError(t, err, "get block failed")
	}))
}

func TestHandler_GetBlockByHeights(t *testing.T) {
	const startHeightKey = "start_height"
	const endHeightKey = "end_height"
	const heightKey = "height"

	t.Run("Range Height Success", handlerTest(func(ctx context.Context, t *testing.T, handler httpHandler, req *testRequest) {
		b := test.BlockHTTP()
		httpBlock := []*models.Block{&b}

		const (
			startHeight = "1"
			endHeight   = "2"
		)

		req.SetData(
			newBlocksURL(map[string]string{
				startHeightKey: startHeight,
				endHeightKey:   endHeight,
			}), httpBlock,
		)

		block, err := handler.getBlocksByHeights(ctx, "", startHeight, endHeight)
		assert.NoError(t, err)
		assert.Equal(t, block, httpBlock)
	}))

	t.Run("List Height Success", handlerTest(func(ctx context.Context, t *testing.T, handler httpHandler, req *testRequest) {
		b1 := test.BlockHTTP()
		b2 := test.BlockHTTP()
		httpBlocks := []*models.Block{&b1, &b2}

		const heights = "1,2"

		req.SetData(
			newBlocksURL(map[string]string{
				heightKey: heights,
			}),
			httpBlocks,
		)

		block, err := handler.getBlocksByHeights(ctx, heights, "", "")
		assert.NoError(t, err)
		assert.Equal(t, block, httpBlocks)
	}))

	t.Run("Range Height Failure", handlerTest(func(ctx context.Context, t *testing.T, handler httpHandler, req *testRequest) {
		testVectors := []struct {
			q   map[string]string
			err string
		}{{
			q:   map[string]string{startHeightKey: "1"},
			err: "must provide either heights or start and end height",
		}, {
			q:   map[string]string{endHeightKey: "1"},
			err: "must provide either heights or start and end height",
		}, {
			q:   map[string]string{},
			err: "must provide either heights or start and end height",
		}}

		for _, testVector := range testVectors {
			req.SetData(newBlocksURL(testVector.q), nil)

			_, err := handler.getBlocksByHeights(
				ctx,
				testVector.q[heightKey],
				testVector.q[startHeightKey],
				testVector.q[endHeightKey],
			)
			assert.EqualError(t, err, testVector.err)
		}
	}))

	t.Run("Bad Request", handlerTest(func(ctx context.Context, t *testing.T, handler httpHandler, req *testRequest) {
		errHTTP := models.ModelError{
			Code:    400,
			Message: "invalid height values",
		}

		const heights = "foo,bar" // invalid

		req.SetErr(
			newBlocksURL(map[string]string{
				heightKey: heights,
			}),
			errHTTP,
		)

		_, err := handler.getBlocksByHeights(ctx, heights, "", "")
		assert.EqualError(t, err, fmt.Sprintf("get block by height %s failed: %s", heights, errHTTP.Message))
	}))
}

// newAccountsURL is a helper factory for building accounts URLs.
func newAccountsURL(address string, query map[string]string) url.URL {
	u, _ := url.Parse(fmt.Sprintf("/accounts/%s", address))
	if query == nil {
		query = map[string]string{}
	}
	query["expand"] = "keys,contracts"

	return addQuery(u, query)
}

func TestHandler_GetAccount(t *testing.T) {

	t.Run("Success", handlerTest(func(ctx context.Context, t *testing.T, handler httpHandler, req *testRequest) {
		httpAccount := test.AccountHTTP()

		const height = "sealed"
		req.SetData(
			newAccountsURL(httpAccount.Address, map[string]string{
				"height": height,
			}),
			httpAccount,
		)

		acc, err := handler.getAccount(ctx, httpAccount.Address, height)
		assert.NoError(t, err)
		assert.Equal(t, *acc, httpAccount)
	}))

	t.Run("Failure", handlerTest(func(ctx context.Context, t *testing.T, handler httpHandler, req *testRequest) {
		errHTTP := models.ModelError{
			Code:    400,
			Message: "invalid height value",
		}

		const (
			heights = "foo" // invalid
			address = "0x1"
		)

		req.SetErr(
			newAccountsURL(address, map[string]string{
				"height": heights,
			}),
			errHTTP,
		)

		_, err := handler.getAccount(ctx, address, heights)
		assert.EqualError(t, err, fmt.Sprintf("get account %s failed: %s", address, errHTTP.Message))
	}))
}

func TestHandler_GetCollection(t *testing.T) {
	const collectionURL = "/collections"

	t.Run("Success", handlerTest(func(ctx context.Context, t *testing.T, handler httpHandler, req *testRequest) {
		httpCollection := test.CollectionHTTP()
		id := "0x1"

		collURL, _ := url.Parse(fmt.Sprintf("%s/%s", collectionURL, id))
		req.SetData(*collURL, httpCollection)

		collection, err := handler.getCollection(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, *collection, httpCollection)
	}))
}

package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	"github.com/onflow/flow-go/engine/access/rest/models"
)

type handler struct {
	client   *http.Client
	accounts *endpoint
	blocks   *endpoint
}

func newHandler(baseUrl string) *handler {
	newEndpoint := newBaseEndpoint(baseUrl)

	return &handler{
		client:   http.DefaultClient,
		accounts: newEndpoint("/accounts/%s"),
		blocks:   newEndpoint("/blocks"),
	}
}

func (h *handler) get(ctx context.Context, url *url.URL, model interface{}) error {
	res, err := h.client.Get(url.String())
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("HTTP GET %s failed", url))
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(model)
	if err != nil {
		return errors.Wrap(err, "JSON decoding failed")
	}

	return nil
}

func (h *handler) getBlockByID(ctx context.Context, ID string) (*models.Block, error) {
	u, err := h.blocks.buildURL(ID)
	if err != nil {
		return nil, err
	}

	var block models.Block
	err = h.get(ctx, u, &block)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("get block ID %s failed", ID))
	}

	return &block, nil
}

func (h *handler) getBlockByHeight(ctx context.Context, height string) ([]*models.Block, error) {
	u, err := h.blocks.buildURL()
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("height", height)
	u.RawQuery = q.Encode()

	var blocks []*models.Block
	err = h.get(ctx, u, blocks)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("get block by height %s failed", height))
	}

	return blocks, nil
}

func (h *handler) getAccount(ctx context.Context, address string, height string) (*models.Account, error) {
	u, err := h.accounts.buildURL(address)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("height", height)
	u.RawQuery = q.Encode()

	var account models.Account
	err = h.get(ctx, u, &account)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Get account %s failed", address))
	}

	return &account, nil
}

func (h *handler) getCollection(ctx context.Context, ID string) (*models.Collection, error) {

}

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	"github.com/onflow/flow-go/engine/access/rest/models"
)

type handler struct {
	client      *http.Client
	accounts    *endpoint
	blocks      *endpoint
	collections *endpoint
	scripts     *endpoint
}

func newHandler(baseUrl string) *handler {
	newEndpoint := newBaseEndpoint(baseUrl)

	return &handler{
		client:      http.DefaultClient,
		accounts:    newEndpoint("/accounts/%s"),
		blocks:      newEndpoint("/blocks"),
		collections: newEndpoint("/collections/%s"),
		scripts:     newEndpoint("/scripts"),
	}
}

func (h *handler) get(ctx context.Context, url *url.URL, model interface{}) error {
	res, err := h.client.Get(url.String())
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("HTTP GET %s failed", url.String()))
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(model)
	if err != nil {
		return errors.Wrap(err, "JSON decoding failed")
	}

	return nil
}

func (h *handler) post(ctx context.Context, url *url.URL, body []byte, model interface{}) error {
	res, err := h.client.Post(
		url.String(),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("HTTP POST %s failed", url.String()))
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
	u, _ := h.blocks.buildURL()

	q := u.Query()
	q.Add("height", height)
	u.RawQuery = q.Encode()

	var blocks []*models.Block
	err := h.get(ctx, u, blocks)
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
	u, err := h.collections.buildURL(ID)
	if err != nil {
		return nil, err
	}

	var collection models.Collection
	err = h.get(ctx, u, &collection)
	if err != nil {
		return nil, err
	}

	return &collection, nil
}

func (h *handler) executeScriptAtBlockHeight(
	ctx context.Context,
	height string,
	script string,
	arguments []string,
) (string, error) {
	u, _ := h.scripts.buildURL()

	q := u.Query()
	q.Add("height", height)
	u.RawQuery = q.Encode()

	body, err := json.Marshal(
		models.ScriptsBody{
			Script:    script,
			Arguments: arguments,
		},
	)
	if err != nil {
		return "", err
	}

	var result string
	err = h.post(ctx, u, body, &result)
	if err != nil {
		return "", err
	}

	return result, nil
}

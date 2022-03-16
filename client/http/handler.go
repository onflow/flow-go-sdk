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

type Handler struct {
	client       *http.Client
	account      *endpoint
	blocks       *endpoint
	collection   *endpoint
	scripts      *endpoint
	transactions *endpoint
	transaction  *endpoint
}

func NewHandler(baseUrl string) *Handler {
	newEndpoint := newBaseEndpoint(baseUrl)

	return &Handler{
		client:       http.DefaultClient,
		account:      newEndpoint("/accounts/%s"),
		blocks:       newEndpoint("/blocks"),
		collection:   newEndpoint("/collections/%s"),
		scripts:      newEndpoint("/scripts"),
		transactions: newEndpoint("/transactions"),
		transaction:  newEndpoint("/transactions/%s"),
	}
}

func (h *Handler) get(_ context.Context, url *url.URL, model interface{}) error {
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

func (h *Handler) post(_ context.Context, url *url.URL, body []byte, model interface{}) error {
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

func (h *Handler) getBlockByID(ctx context.Context, ID string) (*models.Block, error) {
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

func (h *Handler) getBlockByHeight(ctx context.Context, height string) ([]*models.Block, error) {
	u, _ := h.blocks.buildURL()

	q := u.Query()
	q.Add("height", height)
	u.RawQuery = q.Encode()

	var blocks []*models.Block
	err := h.get(ctx, u, blocks)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("get block by height %s failed", height))
	}

	if len(blocks) == 0 {
		return nil, fmt.Errorf("blocks not found")
	}

	return blocks, nil
}

func (h *Handler) getAccount(ctx context.Context, address string, height string) (*models.Account, error) {
	u, err := h.account.buildURL(address)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("height", height)
	u.RawQuery = q.Encode()

	var account models.Account
	err = h.get(ctx, u, &account)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("get account %s failed", address))
	}

	return &account, nil
}

func (h *Handler) getCollection(ctx context.Context, ID string) (*models.Collection, error) {
	u, err := h.collection.buildURL(ID)
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

func (h *Handler) executeScriptAt(
	ctx context.Context,
	query map[string]string,
	script string,
	arguments []string,
) (string, error) {
	u, _ := h.scripts.buildURL()

	q := u.Query()
	for k, v := range query {
		q.Add(k, v)
	}
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

func (h *Handler) executeScriptAtBlockHeight(
	ctx context.Context,
	height string,
	script string,
	arguments []string,
) (string, error) {
	return h.executeScriptAt(
		ctx,
		map[string]string{"block_height": height},
		script,
		arguments,
	)
}

func (h *Handler) executeScriptAtBlockID(
	ctx context.Context,
	ID string,
	script string,
	arguments []string,
) (string, error) {
	return h.executeScriptAt(
		ctx,
		map[string]string{"block_id": ID},
		script,
		arguments,
	)
}

func (h *Handler) getTransaction(ctx context.Context, ID string, includeResult bool) (*models.Transaction, error) {
	var transaction models.Transaction
	u, err := h.transaction.buildURL(ID)
	if err != nil {
		return nil, err
	}

	if includeResult {
		q := u.Query()
		q.Add("expand", "result")
		u.RawQuery = q.Encode()
	}

	err = h.get(ctx, u, &transaction)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("get transaction %s failed", ID))
	}

	return &transaction, nil
}

func (h *Handler) sendTransaction(ctx context.Context, transaction []byte) error {
	u, _ := h.transactions.buildURL()
	var tx models.Transaction

	return h.post(ctx, u, transaction, &tx)
}

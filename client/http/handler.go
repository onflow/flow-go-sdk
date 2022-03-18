/*
 * Flow Go SDK
 *
 * Copyright 2019-2022 Dapper Labs, Inc.
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

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/onflow/flow-go/engine/access/rest/models"
)

type handler struct {
	client *http.Client
	base   string
	debug  bool
}

func newHandler(baseUrl string, debug bool) (*handler, error) {
	// todo validate url and return err
	return &handler{
		client: http.DefaultClient,
		base:   baseUrl,
		debug:  debug,
	}, nil
}

func (h *handler) mustBuildURL(path string) *url.URL {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s%s", h.base, path)) // we ignore error because the values are always valid
	return u
}

func (h *handler) get(_ context.Context, url *url.URL, model interface{}) error {
	if h.debug {
		fmt.Printf("\n-> GET %s t=%d", url.String(), time.Now().Unix())
	}

	// todo use a .Do() method and use the context
	res, err := h.client.Get(url.String())
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("HTTP GET %s failed", url.String()))
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("HTTP GET %s failed, status code: %d, response :%s", url.String(), res.StatusCode, body)
	}

	if h.debug {
		fmt.Printf("\n<- GET %s t=%d - %s", url.String(), time.Now().Unix(), body)
	}

	err = json.Unmarshal(body, &model)
	if err != nil {
		return errors.Wrap(err, "JSON decoding failed")
	}

	return nil
}

func (h *handler) post(_ context.Context, url *url.URL, body []byte, model interface{}) error {
	if h.debug {
		fmt.Printf("\n-> POST %s t=%d - %s", url.String(), time.Now().Unix(), string(body))
	}

	res, err := h.client.Post(
		url.String(),
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("HTTP POST %s failed", url.String()))
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("HTTP POST %s failed, status code: %d, response :%s", url.String(), res.StatusCode, responseBody)
	}

	if h.debug {
		fmt.Printf("\n<- POST %s t=%d - %s", url.String(), time.Now().Unix(), string(body))
	}

	err = json.Unmarshal(responseBody, &model)
	if err != nil {
		return errors.Wrap(err, "JSON decoding failed")
	}

	return nil
}

func (h *handler) getBlockByID(ctx context.Context, ID string) (*models.Block, error) {
	u := h.mustBuildURL(fmt.Sprintf("/blocks/%s", ID))

	q := u.Query()
	q.Add("expand", "payload")
	u.RawQuery = q.Encode()

	var blocks []*models.Block
	err := h.get(ctx, u, &blocks)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("get block ID %s failed", ID))
	}

	if len(blocks) == 0 {
		return nil, fmt.Errorf("block ID %s not found", ID)
	}

	return blocks[0], nil
}

func (h *handler) getBlockByHeight(ctx context.Context, height string) ([]*models.Block, error) {
	u := h.mustBuildURL("/blocks")

	q := u.Query()
	q.Add("height", height)
	q.Add("expand", "payload")
	u.RawQuery = q.Encode()

	var blocks []*models.Block
	err := h.get(ctx, u, &blocks)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("get block by height %s failed", height))
	}

	if len(blocks) == 0 {
		return nil, fmt.Errorf("blocks not found")
	}

	return blocks, nil
}

func (h *handler) getAccount(ctx context.Context, address string, height string) (*models.Account, error) {
	u := h.mustBuildURL(fmt.Sprintf("/accounts/%s", address))

	q := u.Query()
	q.Add("height", height)
	q.Add("expand", "keys,contracts")
	u.RawQuery = q.Encode()

	var account models.Account
	err := h.get(ctx, u, &account)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("get account %s failed", address))
	}

	return &account, nil
}

func (h *handler) getCollection(ctx context.Context, ID string) (*models.Collection, error) {
	var collection models.Collection
	err := h.get(
		ctx, h.mustBuildURL(fmt.Sprintf("/collections/%s", ID)),
		&collection,
	)
	if err != nil {
		return nil, err
	}

	return &collection, nil
}

func (h *handler) executeScriptAt(
	ctx context.Context,
	query map[string]string,
	script string,
	arguments []string,
) (string, error) {
	u := h.mustBuildURL("/scripts")

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

func (h *handler) executeScriptAtBlockHeight(
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

func (h *handler) executeScriptAtBlockID(
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

func (h *handler) getTransaction(ctx context.Context, ID string, includeResult bool) (*models.Transaction, error) {
	var transaction models.Transaction
	u := h.mustBuildURL(fmt.Sprintf("/transactions/%s", ID))

	if includeResult {
		q := u.Query()
		q.Add("expand", "result")
		u.RawQuery = q.Encode()
	}

	err := h.get(ctx, u, &transaction)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("get transaction %s failed", ID))
	}

	return &transaction, nil
}

func (h *handler) sendTransaction(ctx context.Context, transaction []byte) error {
	var tx models.Transaction
	return h.post(ctx, h.mustBuildURL("/transactions"), transaction, &tx)
}

func (h *handler) getEvents(
	ctx context.Context,
	eventType string,
	start string,
	end string,
	blockIDs []string,
) ([]models.BlockEvents, error) {
	u := h.mustBuildURL("/events")

	q := u.Query()
	if start != "" && end != "" {
		q.Add("start_height", start)
		q.Add("end_height", end)
	} else if len(blockIDs) != 0 {
		q.Add("block_ids", strings.Join(blockIDs, ","))
	} else {
		return nil, fmt.Errorf("must either provide start and end height or block IDs")
	}

	q.Add("type", eventType)
	u.RawQuery = q.Encode()

	var events []models.BlockEvents
	err := h.get(ctx, u, &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

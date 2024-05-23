/*
 * Flow Go SDK
 *
 * Copyright 2019 Dapper Labs, Inc.
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

	"github.com/onflow/flow-go-sdk/access/http/models"

	"github.com/pkg/errors"
)

type queryOpts interface {
	toQuery() (string, string)
}

type HTTPError struct {
	Url     string
	Code    int
	Message string
}

func (h HTTPError) Error() string {
	return h.Message
}

type httpHandler struct {
	client *http.Client
	base   string
	debug  bool
}

func newHandler(host string, debug bool) (*httpHandler, error) {
	_, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	return &httpHandler{
		client: http.DefaultClient,
		base:   host,
		debug:  debug,
	}, nil
}

func (h *httpHandler) mustBuildURL(path string, opts ...queryOpts) *url.URL {
	u, _ := url.ParseRequestURI(fmt.Sprintf("%s%s", h.base, path))

	for _, opt := range opts {
		q := u.Query()
		q.Add(opt.toQuery())
		u.RawQuery = q.Encode()
	}

	return u
}

func (h *httpHandler) get(_ context.Context, url *url.URL, model interface{}) error {
	if h.debug {
		fmt.Printf("\n-> GET %s t=%d", url.String(), time.Now().Unix())
	}

	// todo use a .Do() method and use the context
	res, err := h.client.Get(url.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode >= http.StatusBadRequest {
		if h.debug {
			fmt.Printf("\n<- FAILED GET %s t=%d status=%d - %s", url.String(), res.StatusCode, time.Now().Unix(), body)
		}

		var httpErr HTTPError
		err = json.Unmarshal(body, &httpErr)
		if err != nil {
			return err
		}

		httpErr.Url = url.String()
		return httpErr
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

func (h *httpHandler) post(_ context.Context, url *url.URL, body []byte, model interface{}) error {
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
		if h.debug {
			fmt.Printf("\n<- POST FAILED %s, status=%d, response: %s", url.String(), res.StatusCode, responseBody)
		}

		var httpErr HTTPError
		err = json.Unmarshal(responseBody, &httpErr)
		if err != nil {
			return err
		}

		httpErr.Url = url.String()
		return httpErr
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

func (h *httpHandler) getNetworkParameters(ctx context.Context, opts ...queryOpts) (*models.NetworkParameters, error) {
	var networkParameters models.NetworkParameters
	err := h.get(ctx, h.mustBuildURL("/network/parameters", opts...), &networkParameters)
	if err != nil {
		return nil, errors.Wrap(err, "get network parameters failed")
	}

	return &networkParameters, nil
}

func (h *httpHandler) getNodeVersionInfo(ctx context.Context, opts ...queryOpts) (*models.NodeVersionInfo, error) {
	var nodeVersionInfo models.NodeVersionInfo
	err := h.get(ctx, h.mustBuildURL("/node_version_info", opts...), &nodeVersionInfo)
	if err != nil {
		return nil, errors.Wrap(err, "get node version info failed")
	}

	return &nodeVersionInfo, nil
}

func (h *httpHandler) getBlockByID(ctx context.Context, ID string, opts ...queryOpts) (*models.Block, error) {
	u := h.mustBuildURL(fmt.Sprintf("/blocks/%s", ID), opts...)

	q := u.Query()
	q.Add("expand", "payload")
	u.RawQuery = q.Encode()

	var blocks []*models.Block
	err := h.get(ctx, u, &blocks)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("get block ID %s failed", ID))
	}

	if len(blocks) == 0 { // sanity check
		return nil, fmt.Errorf("get block failed")
	}

	return blocks[0], nil
}

func (h *httpHandler) getBlocksByHeights(
	ctx context.Context,
	heights string,
	startHeight string,
	endHeight string,
	opts ...queryOpts,
) ([]*models.Block, error) {
	u := h.mustBuildURL("/blocks", opts...)

	q := u.Query()
	if heights != "" {
		q.Add("height", heights)
	} else if startHeight != "" && endHeight != "" {
		q.Add("start_height", startHeight)
		q.Add("end_height", endHeight)
	} else {
		return nil, fmt.Errorf("must provide either heights or start and end height")
	}

	q.Add("expand", "payload")
	u.RawQuery = q.Encode()

	var blocks []*models.Block
	err := h.get(ctx, u, &blocks)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("get block by height %s failed", heights))
	}

	return blocks, nil
}

func (h *httpHandler) getAccount(
	ctx context.Context,
	address string,
	height string,
	opts ...queryOpts,
) (*models.Account, error) {
	u := h.mustBuildURL(fmt.Sprintf("/accounts/%s", address), opts...)

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

func (h *httpHandler) getCollection(ctx context.Context, ID string, opts ...queryOpts) (*models.Collection, error) {
	var collection models.Collection
	err := h.get(
		ctx, h.mustBuildURL(fmt.Sprintf("/collections/%s", ID), opts...),
		&collection,
	)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("get collection ID %s failed", ID))
	}

	return &collection, nil
}

func (h *httpHandler) executeScript(
	ctx context.Context,
	query map[string]string,
	script string,
	arguments []string,
	opts ...queryOpts,
) (string, error) {
	u := h.mustBuildURL("/scripts", opts...)

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
		return "", errors.Wrap(err, "executing script failed")
	}

	var result string
	err = h.post(ctx, u, body, &result)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("executing script %s failed", script))
	}

	return result, nil
}

func (h *httpHandler) executeScriptAtBlockHeight(
	ctx context.Context,
	height string,
	script string,
	arguments []string,
	opts ...queryOpts,
) (string, error) {
	return h.executeScript(
		ctx,
		map[string]string{"block_height": height},
		script,
		arguments,
	)
}

func (h *httpHandler) executeScriptAtBlockID(
	ctx context.Context,
	ID string,
	script string,
	arguments []string,
	opts ...queryOpts,
) (string, error) {
	return h.executeScript(
		ctx,
		map[string]string{"block_id": ID},
		script,
		arguments,
	)
}

func (h *httpHandler) getTransaction(
	ctx context.Context,
	ID string,
	includeResult bool,
	opts ...queryOpts,
) (*models.Transaction, error) {
	var transaction models.Transaction
	u := h.mustBuildURL(fmt.Sprintf("/transactions/%s", ID), opts...)

	if includeResult {
		q := u.Query()
		q.Add("expand", "result")
		u.RawQuery = q.Encode()
	}

	err := h.get(ctx, u, &transaction)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("get transaction ID %s failed", ID))
	}

	return &transaction, nil
}

func (h *httpHandler) sendTransaction(ctx context.Context, transaction []byte, opts ...queryOpts) error {
	var tx models.Transaction
	return h.post(ctx, h.mustBuildURL("/transactions", opts...), transaction, &tx)
}

func (h *httpHandler) getEvents(
	ctx context.Context,
	eventType string,
	start string,
	end string,
	blockIDs []string,
	opts ...queryOpts,
) ([]models.BlockEvents, error) {
	u := h.mustBuildURL("/events", opts...)

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
		return nil, errors.Wrap(err, fmt.Sprintf("get events by type %s failed", eventType))
	}

	return events, nil
}

func (h *httpHandler) getExecutionResults(
	ctx context.Context,
	blockIDs []string,
	opts ...queryOpts,
) ([]models.ExecutionResult, error) {
	u := h.mustBuildURL("/execution_results", opts...)

	q := u.Query()
	q.Add("block_ids", strings.Join(blockIDs, ","))
	u.RawQuery = q.Encode()

	var results []models.ExecutionResult
	err := h.get(ctx, u, &results)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("get execution results by IDs %v failed", blockIDs))
	}

	return results, nil
}

func (h *httpHandler) getExecutionResultByID(ctx context.Context, id string, opts ...queryOpts) (*models.ExecutionResult, error) {
	u := h.mustBuildURL(fmt.Sprintf("/execution_results/%s", id), opts...)

	var result models.ExecutionResult
	err := h.get(ctx, u, &result)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("get execution result by ID %s failed", id))
	}

	return &result, nil
}

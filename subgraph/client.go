package subgraph

import (
	"errors"
	"fmt"
	"log"

	gresty "github.com/go-resty/resty/v2"

	"github.com/onflow/flow-go-sdk/subgraph/transaction_detail"
	"github.com/onflow/flow-go-sdk/subgraph/transactions"
)

const (
	Flow_subgraph_mainnet = "https://api.findlabs.io/flowdiver/v1/graphql"
	Flow_subgraph_testnet = "https://api.findlabs.io/flowdiver_testnet/v1/graphql"
)

var errSubGraphHTTPError = errors.New("SubGraph http error")

type subgraphClient struct {
	sgClient *gresty.Client
}

func NewFlowClient(baseUrl string) (*subgraphClient, error) {
	grestyClient := gresty.New()
	grestyClient.SetBaseURL(baseUrl)
	grestyClient.OnAfterResponse(func(c *gresty.Client, r *gresty.Response) error {
		statusCode := r.StatusCode()
		if statusCode >= 400 {
			method := r.Request.Method
			url := r.Request.URL
			return fmt.Errorf("%d cannot %s %s: %w", statusCode, method, url, errSubGraphHTTPError)
		}
		return nil
	})
	fclient := &subgraphClient{
		sgClient: grestyClient,
	}
	return fclient, nil
}

func (f *subgraphClient) GetTransactionsByAddress(graphRequest *GetSubGraphRequest) (*transactions.FlowTransactionsResp, error) {
	var resultTxList transactions.FlowTransactionsResp
	response, err := f.sgClient.R().SetResult(&resultTxList).SetBody(graphRequest).Post("")
	if err != nil {
		log.Printf("GetTxListByAddress  Error: %+v\n", err)
		return nil, err
	}
	if response.StatusCode() != 200 {
		log.Printf("GetTxListByAddress  Error: %+v\n", err)
		return nil, err
	}
	return &resultTxList, nil
}

func (f *subgraphClient) GetTransactionById(graphRequest *GetSubGraphRequest) (*transaction_detail.FlowTransactionDetailResp, error) {
	var resultTxDetail transaction_detail.FlowTransactionDetailResp
	response, err := f.sgClient.R().SetResult(&resultTxDetail).SetBody(graphRequest).Post("")
	if err != nil {
		log.Printf("GetTransactionById  Error: %+v\n", err)
		return nil, err
	}
	if response.StatusCode() != 200 {
		log.Printf("GetTransactionById  Error: %+v\n", err)
		return nil, err
	}
	return &resultTxDetail, nil
}

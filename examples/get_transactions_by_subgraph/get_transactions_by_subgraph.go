package main

import (
	"encoding/json"
	"fmt"

	"github.com/onflow/flow-go-sdk/subgraph"
)

func main() {
	client, err := subgraph.NewFlowClient(subgraph.Flow_subgraph_mainnet)
	if err != nil {
		panic(err)
	}
	transactionDetailReq := client.GetTransactionDetailReq("c260981bb2d5fc80986436cac01b1f609144f369f4c218005c1795c5a25dbbc8")
	transactionById, err := client.GetTransactionById(transactionDetailReq)
	if err != nil {
		panic(err)
	}

	printJsonStr(transactionById)

	fmt.Println("--------------------------------------------------")

	transactionsByAddressReq := client.GetTransactionsByAddressReq(&subgraph.TransactionListVariables{
		Address:           "0x8f4f599546e2d7eb",
		Limit:             25,
		Offset:            0,
		TimeFilter:        make(map[string]interface{}),
		TypeFilter:        make(map[string]interface{}),
		StatusFilter:      make(map[string]interface{}),
		ProposerFilter:    make(map[string]interface{}),
		PayerFilter:       make(map[string]interface{}),
		GasRangeFilter:    make(map[string]interface{}),
		EventCountFilter:  make(map[string]interface{}),
		BlockHeightFilter: make(map[string]interface{}),
		AuthorizersFilter: make(map[string]interface{}),
		ActorFilter:       make([]interface{}, 0)})
	transactions, err := client.GetTransactionsByAddress(transactionsByAddressReq)
	if err != nil {
		panic(err)
	}
	printJsonStr(transactions)
}

func printJsonStr(param interface{}) {
	marshal, _ := json.Marshal(param)
	fmt.Println(string(marshal))
}

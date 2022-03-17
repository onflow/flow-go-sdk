package main

import (
	"context"
	"fmt"

	"github.com/onflow/flow-go-sdk/client/http"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/examples"
)

func main() {
	demo()
}

func demo() {
	ctx := context.Background()
	flowClient, err := http.NewDefaultEmulatorClient(false)
	examples.Handle(err)

	txID := examples.RandomTransaction(flowClient).ID()

	tx, err := flowClient.GetTransaction(ctx, txID)
	printTransaction(tx, err)

	txr, err := flowClient.GetTransactionResult(ctx, txID)
	printTransactionResult(txr, err)
}

func printTransaction(tx *flow.Transaction, err error) {
	examples.Handle(err)

	fmt.Printf("\nID: %s", tx.ID().String())
	fmt.Printf("\nPayer: %s", tx.Payer.String())
	fmt.Printf("\nProposer: %s", tx.ProposalKey.Address.String())
	fmt.Printf("\nAuthorizers: %s", tx.Authorizers)
}

func printTransactionResult(txr *flow.TransactionResult, err error) {
	examples.Handle(err)

	fmt.Printf("\nStatus: %s", txr.Status.String())
	fmt.Printf("\nError: %v", txr.Error)
}

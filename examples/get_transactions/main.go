package main

import (
	"context"
	"fmt"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/examples"
	"google.golang.org/grpc"
)

func main() {
	txID := prepareDemo()
	demo(txID)
}

func demo(txID flow.Identifier) {
	ctx := context.Background()
	flowClient := examples.NewFlowClient()

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

func prepareDemo() flow.Identifier {
	flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
	examples.Handle(err)
	defer func() {
		err := flowClient.Close()
		if err != nil {
			panic(err)
		}
	}()

	return examples.RandomTransaction(flowClient).ID()
}

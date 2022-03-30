package main

import (
	"context"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/examples"
)

func main() {
	tx := prepareDemo()
	demo(tx)
}

func demo(tx *flow.Transaction) {
	ctx := context.Background()
	flowClient := examples.NewFlowClient()

	err := flowClient.SendTransaction(ctx, *tx)
	if err != nil {
		panic(err)
	}
}

func prepareDemo() *flow.Transaction {
	flowClient := examples.NewFlowClient()
	defer func() {
		err := flowClient.Close()
		if err != nil {
			panic(err)
		}
	}()

	serviceAcctAddr, serviceAcctKey, serviceSigner := examples.ServiceAccount(flowClient)

	tx := flow.NewTransaction().
		SetPayer(serviceAcctAddr).
		SetProposalKey(serviceAcctAddr, serviceAcctKey.Index, serviceAcctKey.SequenceNumber).
		SetScript([]byte("transaction {}")).
		SetReferenceBlockID(examples.GetReferenceBlockId(flowClient))

	err := tx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, serviceSigner)
	examples.Handle(err)

	return tx
}

package main

import (
	"context"

	"github.com/onflow/flow-go-sdk/client"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/examples"
)

func main() {
	flowClient := examples.NewFlowHTTPClient()
	tx := prepareDemo(flowClient)
	demo(flowClient, tx)
}

func demo(flowClient client.Client, tx *flow.Transaction) {
	ctx := context.Background()
	err := flowClient.SendTransaction(ctx, *tx)
	examples.Handle(err)
}

func prepareDemo(flowClient client.Client) *flow.Transaction {
	serviceAcctAddr, serviceAcctKey, serviceSigner := examples.ServiceAccount(flowClient)

	tx := flow.NewTransaction().
		SetPayer(serviceAcctAddr).
		SetProposalKey(serviceAcctAddr, serviceAcctKey.Index, serviceAcctKey.SequenceNumber).
		SetScript([]byte("transaction { prepare(acct: AuthAccount) {} }")).
		AddAuthorizer(serviceAcctAddr).
		SetReferenceBlockID(examples.GetReferenceBlockId(flowClient))

	err := tx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, serviceSigner)
	examples.Handle(err)

	return tx
}

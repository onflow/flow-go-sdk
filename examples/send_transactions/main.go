package main

import (
	"context"

	"github.com/onflow/flow-go-sdk/client/http"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/examples"
)

func main() {
	demo()
}

func demo() {
	flowClient, err := http.NewDefaultEmulatorClient()
	examples.Handle(err)

	serviceAcctAddr, serviceAcctKey, serviceSigner := examples.ServiceAccount(flowClient)

	tx := flow.NewTransaction().
		SetPayer(serviceAcctAddr).
		SetProposalKey(serviceAcctAddr, serviceAcctKey.Index, serviceAcctKey.SequenceNumber).
		SetScript([]byte(`transaction {

  prepare(acc: AuthAccount) {}

  execute {
    log("test")
  }
}`)).
		AddAuthorizer(serviceAcctAddr).
		SetReferenceBlockID(examples.GetReferenceBlockId(flowClient))

	err = tx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, serviceSigner)
	examples.Handle(err)

	err = flowClient.SendTransaction(context.Background(), *tx)
	examples.Handle(err)
}

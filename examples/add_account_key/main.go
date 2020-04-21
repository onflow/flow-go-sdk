package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/examples"
	"github.com/onflow/flow-go-sdk/templates"
)

func main() {
	AddAccountKeyDemo()
}

func AddAccountKeyDemo() {
	ctx := context.Background()
	acctAddr, acctKey, acctSigner := examples.CreateAccount() // Creates a new account and returns the address+key

	// Create the new key to add to your account
	myPrivateKey := examples.RandomPrivateKey()
	myAcctKey := flow.NewAccountKey().
		FromPrivateKey(myPrivateKey).
		SetHashAlgo(crypto.SHA3_256).
		SetWeight(flow.AccountKeyWeightThreshold)

	flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
	examples.Handle(err)

	// Create a Cadence script that will add another key to our account.
	addKeyScript, err := templates.AddAccountKey(myAcctKey)
	examples.Handle(err)

	// Create a transaction to execute the script.
	// The transaction is signed by our account key so it has permission to add keys.
	addKeyTx := flow.NewTransaction().
		SetScript(addKeyScript).
		SetProposalKey(acctAddr, acctKey.ID, acctKey.SequenceNumber).
		SetPayer(acctAddr).
		// This defines which accounts are accessed by this transaction
		AddAuthorizer(acctAddr)

	// Sign the transaction with the new account.
	err = addKeyTx.SignEnvelope(acctAddr, acctKey.ID, acctSigner)
	examples.Handle(err)

	// Send the transaction to the network.
	err = flowClient.SendTransaction(ctx, *addKeyTx)
	examples.Handle(err)

	examples.WaitForSeal(ctx, flowClient, addKeyTx.ID())

	fmt.Println("Public key added to account!")
}

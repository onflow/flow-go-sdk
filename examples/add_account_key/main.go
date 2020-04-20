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
	accountAddr, accountKey, accountPrivateKey := examples.CreateAccount() // Creates a new account and returns the address+key

	// Create the new key to add to your account
	myPrivateKey := examples.RandomPrivateKey()
	myAcctKey := flow.AccountKey{
		PublicKey: myPrivateKey.PublicKey(),
		SigAlgo:   myPrivateKey.Algorithm(),
		HashAlgo:  crypto.SHA3_256,
		Weight:    flow.AccountKeyWeightThreshold,
	}

	flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
	examples.Handle(err)

	// Create a Cadence script that will add another key to our account.
	addKeyScript, err := templates.AddAccountKey(myAcctKey)
	examples.Handle(err)

	// Create a transaction to execute the script.
	// The transaction is signed by our account key so it has permission to add keys.
	addKeyTx := flow.NewTransaction().
		SetScript(addKeyScript).
		SetProposalKey(accountAddr, accountKey.ID, accountKey.SequenceNumber).
		SetPayer(accountAddr).
		// This defines which accounts are accessed by this transaction
		AddAuthorizer(accountAddr)

	// Sign the transaction with the new account.
	err = addKeyTx.SignEnvelope(
		accountAddr,
		accountKey.ID,
		crypto.NewNaiveSigner(accountPrivateKey, accountKey.HashAlgo),
	)
	examples.Handle(err)

	// Send the transaction to the network.
	err = flowClient.SendTransaction(ctx, *addKeyTx)
	examples.Handle(err)

	examples.WaitForSeal(ctx, flowClient, addKeyTx.ID())

	fmt.Println("Public key added to account!")

}

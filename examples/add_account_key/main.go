package main

import (
	"context"
	"fmt"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/client"
	"github.com/dapperlabs/flow-go-sdk/examples"
	"github.com/dapperlabs/flow-go-sdk/keys"
	"github.com/dapperlabs/flow-go-sdk/templates"
)

func main() {
	AddAccountKeyDemo()
}

func AddAccountKeyDemo() {
	ctx := context.Background()
	accountKey, accountAddr := examples.CreateAccount() // Creates a new account and returns the address+key

	newAccountKey := examples.RandomPrivateKey() // Create the new key to add to your account
	flowClient, err := client.New("127.0.0.1:3569")
	examples.Handle(err)

	// Create a Cadence script that will add another key to our account.
	addKeyScript, err := templates.AddAccountKey(newAccountKey.PublicKey(keys.PublicKeyWeightThreshold))
	examples.Handle(err)

	// Create a transaction to execute the script.
	// The transaction is signed by our account key so it has permission to add keys.
	addKeyTx := flow.Transaction{
		Script:       addKeyScript,
		Nonce:        examples.GetNonce(),
		ComputeLimit: 10,
		// This defines which accounts are accessed by this transaction
		ScriptAccounts: []flow.Address{accountAddr},
		PayerAccount:   accountAddr,
	}

	// Sign the transaction and add the signature to the transaction.
	sig, err := keys.SignTransaction(addKeyTx, accountKey)
	examples.Handle(err)
	addKeyTx.AddSignature(accountAddr, sig)

	// Send the transaction to the network.
	err = flowClient.SendTransaction(ctx, addKeyTx)
	examples.Handle(err)

	examples.WaitForSeal(ctx, flowClient, addKeyTx.Hash())

	fmt.Println("Public key added to account!")

}

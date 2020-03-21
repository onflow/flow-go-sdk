package main

import (
	"context"
	"fmt"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/client"
	"github.com/dapperlabs/flow-go-sdk/keys"
	"github.com/dapperlabs/flow-go-sdk/templates"

	utils "github.com/dapperlabs/flow-go-sdk/utils/examples"
)

func main() {
	AddAccountKeyDemo()
}

func AddAccountKeyDemo() {
	ctx := context.Background()
	accountKey, accountAddr := utils.CreateAccount() // Creates a new account and returns the address+key

	newAccountKey := utils.RandomPrivateKey() // Create the new key to add to your account
	flowClient, err := client.New("127.0.0.1:3569")
	utils.Handle(err)

	// Create a Cadence script that will add another key to our account.
	addKeyScript, err := templates.AddAccountKey(newAccountKey.PublicKey(keys.PublicKeyWeightThreshold))
	utils.Handle(err)

	// Create a transaction to execute the script.
	// The transaction is signed by our account key so it has permission to add keys.
	addKeyTx := flow.Transaction{
		Script:       addKeyScript,
		Nonce:        utils.GetNonce(),
		ComputeLimit: 10,
		// This defines which accounts are accessed by this transaction
		ScriptAccounts: []flow.Address{accountAddr},
		PayerAccount:   accountAddr,
	}

	// Sign the transaction and add the signature to the transaction.
	sig, err := keys.SignTransaction(addKeyTx, accountKey)
	utils.Handle(err)
	addKeyTx.AddSignature(accountAddr, sig)

	// Send the transaction to the network.
	err = flowClient.SendTransaction(ctx, addKeyTx)
	utils.Handle(err)

	utils.WaitForSeal(ctx, flowClient, addKeyTx.Hash())

	fmt.Println("Public key added to account!")

}

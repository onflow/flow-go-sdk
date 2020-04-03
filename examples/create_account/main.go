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
	CreateAccountDemo()
}

func CreateAccountDemo() {
	ctx := context.Background()
	flowClient, err := client.New("127.0.0.1:3569")
	examples.Handle(err)

	rootAcctAddr, rootAcctKey := examples.RootAccount()

	myPrivateKey := examples.RandomPrivateKey()
	myPublicKey := myPrivateKey.PublicKey(keys.PublicKeyWeightThreshold)

	// Create a Cadence script which will create an account with one key with weight 1 and
	createAccountScript, err := templates.CreateAccount([]flow.AccountPublicKey{myPublicKey}, nil)
	examples.Handle(err)

	// Create a transaction that will execute the script. The transaction is signed
	// by the root account.
	createAccountTx := flow.Transaction{
		Script:         createAccountScript,
		Nonce:          examples.GetNonce(),
		ComputeLimit:   10,
		ScriptAccounts: nil,
		PayerAccount:   rootAcctAddr,
	}

	// Sign the transaction with the root account, which already exists
	// All new accounts must be created by an existing account
	sig, err := keys.SignTransaction(createAccountTx, rootAcctKey)
	examples.Handle(err)

	// Attach the signature to the transaction
	createAccountTx.AddSignature(rootAcctAddr, sig)

	// Send the transaction to the network
	err = flowClient.SendTransaction(ctx, createAccountTx)
	examples.Handle(err)

	accountCreationTxRes := examples.WaitForSeal(ctx, flowClient, createAccountTx.Hash())

	var myAddress flow.Address

	for _, event := range accountCreationTxRes.Events {
		if event.Type == flow.EventAccountCreated {
			accountCreatedEvent := flow.AccountCreatedEvent(event)
			myAddress = accountCreatedEvent.Address()
		}
	}

	fmt.Println("Account created with address:", myAddress.Hex())
}

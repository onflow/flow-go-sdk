package main

import (
	"context"
	"fmt"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/client"
	"github.com/dapperlabs/flow-go-sdk/keys"
	"github.com/dapperlabs/flow-go-sdk/templates"

	"github.com/dapperlabs/flow-developer-demo/examples/utils"
)

func main() {
	// CreateAccountDemo()
	SendTransaction()
}

func SendTransaction() {
	ctx := context.Background()
	flowClient, err := client.New("127.0.0.1:3569")
	utils.Handle(err)

	rootAcctAddr, rootAcctKey := utils.RootAccount()

	createAccountTx := flow.Transaction{
		Script: []byte(`
			transaction {
				execute {
					log("Initialized marketplace")
					log("Minted 100 tokens, deposited to 0xf5a1c511")
					log("5 tokens transferred from 0xf5a1c511 to 0xf87db879")
					log("5 tokens transferred from 0xf5a1c511 to 0x30770201")
					log("5 tokens transferred from 0xf5a1c511 to a235513c")
				}
			}
		`),
		Nonce:          utils.GetNonce(),
		ComputeLimit:   10,
		ScriptAccounts: nil,
		PayerAccount:   rootAcctAddr,
	}

	sig, err := keys.SignTransaction(createAccountTx, rootAcctKey)
	utils.Handle(err)

	createAccountTx.AddSignature(rootAcctAddr, sig)

	err = flowClient.SendTransaction(ctx, createAccountTx)
	utils.Handle(err)
}

func CreateAccountDemo() {
	ctx := context.Background()
	flowClient, err := client.New("127.0.0.1:3569")
	utils.Handle(err)

	rootAcctAddr, rootAcctKey := utils.RootAccount()

	myPrivateKey := utils.RandomPrivateKey()
	myPublicKey := myPrivateKey.PublicKey(keys.PublicKeyWeightThreshold)

	// Create a Cadence script which will create an account with one key with weight 1 and
	createAccountScript, err := templates.CreateAccount([]flow.AccountPublicKey{myPublicKey}, nil)
	utils.Handle(err)

	// Create a transaction that will execute the script. The transaction is signed
	// by the root account.
	createAccountTx := flow.Transaction{
		Script:         createAccountScript,
		Nonce:          utils.GetNonce(),
		ComputeLimit:   10,
		ScriptAccounts: nil,
		PayerAccount:   rootAcctAddr,
	}

	// Sign the transaction with the root account, which already exists
	// All new accounts must be created by an existing account
	sig, err := keys.SignTransaction(createAccountTx, rootAcctKey)
	utils.Handle(err)

	// Attach the signature to the transaction
	createAccountTx.AddSignature(rootAcctAddr, sig)

	// Send the transaction to the network
	err = flowClient.SendTransaction(ctx, createAccountTx)
	utils.Handle(err)

	accountCreationTxRes := utils.WaitForSeal(ctx, flowClient, createAccountTx.Hash())

	var myAddress flow.Address

	for _, event := range accountCreationTxRes.Events {
		if event.Type == flow.EventAccountCreated {
			accountCreatedEvent, err := flow.DecodeAccountCreatedEvent(event.Payload)
			utils.Handle(err)

			myAddress = accountCreatedEvent.Address()
		}
	}

	fmt.Println("Account created with address:", myAddress.Hex())
}

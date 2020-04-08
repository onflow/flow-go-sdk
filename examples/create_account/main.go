package main

import (
	"context"
	"fmt"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/client"
	"github.com/dapperlabs/flow-go-sdk/crypto"
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

	rootAcctAddr, rootAcctKey, rootPrivateKey := examples.RootAccount()
	rootKeySigner := crypto.NewNaiveSigner(rootPrivateKey, rootAcctKey.HashAlgo)

	myPrivateKey := examples.RandomPrivateKey()
	myPublicKey := flow.AccountKey{
		PublicKey: myPrivateKey.PublicKey(),
		SignAlgo:  keys.ECDSA_P256_SHA3_256.SigningAlgorithm(),
		HashAlgo:  keys.ECDSA_P256_SHA3_256.HashingAlgorithm(),
		Weight:    keys.PublicKeyWeightThreshold,
	}

	// Create a Cadence script which will create an account with one key with weight 1 and
	createAccountScript, err := templates.CreateAccount([]flow.AccountKey{myPublicKey}, nil)
	examples.Handle(err)

	// Create a transaction that will execute the script. The transaction is signed
	// by the root account.
	createAccountTx := flow.NewTransaction().
		SetScript(createAccountScript).
		SetProposalKey(rootAcctAddr, rootAcctKey.ID, rootAcctKey.SequenceNumber).
		SetPayer(rootAcctAddr, rootAcctKey.ID)

	// Sign the transaction with the root account, which already exists
	// All new accounts must be created by an existing account
	err = createAccountTx.SignContainer(
		rootAcctAddr,
		rootAcctKey.ID,
		rootKeySigner,
	)
	examples.Handle(err)

	// Send the transaction to the network
	err = flowClient.SendTransaction(ctx, *createAccountTx)
	examples.Handle(err)

	accountCreationTxRes := examples.WaitForSeal(ctx, flowClient, createAccountTx.ID())

	var myAddress flow.Address

	for _, event := range accountCreationTxRes.Events {
		if event.Type == flow.EventAccountCreated {
			accountCreatedEvent := flow.AccountCreatedEvent(event)
			myAddress = accountCreatedEvent.Address()
		}
	}

	fmt.Println("Account created with address:", myAddress.Hex())
}

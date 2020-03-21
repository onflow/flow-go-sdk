package main

import (
	"context"
	"fmt"

	"github.com/dapperlabs/cadence"
	"github.com/dapperlabs/cadence/encoding"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/client"
	"github.com/dapperlabs/flow-go-sdk/examples/contracts"
	"github.com/dapperlabs/flow-go-sdk/keys"
	"github.com/dapperlabs/flow-go-sdk/templates"
	utils "github.com/dapperlabs/flow-go-sdk/utils/examples"
)

const GreatTokenContractFile = "../contracts/contracts/great-token.cdc"

func main() {
	// Connect to an emulator running locally
	ctx := context.Background()
	flowClient, err := client.New("127.0.0.1:3569")
	utils.Handle(err)

	myPrivateKey := utils.RandomPrivateKey()
	myPublicKey := myPrivateKey.PublicKey(keys.PublicKeyWeightThreshold)

	// Generate an account creation script
	createAccountScript, err := templates.CreateAccount([]flow.AccountPublicKey{myPublicKey}, nil)
	utils.Handle(err)

	rootAcctAddr, rootAcctKey := utils.RootAccount()

	createAccountTx := flow.Transaction{
		Script:       createAccountScript,
		Nonce:        utils.GetNonce(),
		ComputeLimit: 10,
		PayerAccount: rootAcctAddr,
	}

	sig, err := keys.SignTransaction(createAccountTx, rootAcctKey)
	utils.Handle(err)

	createAccountTx.AddSignature(rootAcctAddr, sig)

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

	fmt.Println("My Address:", myAddress.Hex())

	// Deploy the Great NFT contract
	nftCode := utils.ReadFile(GreatTokenContractFile)
	deployScript, err := templates.CreateAccount(nil, nftCode)

	deployContractTx := flow.Transaction{
		Script:       deployScript,
		Nonce:        utils.GetNonce(),
		ComputeLimit: 10,
		PayerAccount: myAddress,
	}
	sig, err = keys.SignTransaction(deployContractTx, myPrivateKey)
	utils.Handle(err)

	deployContractTx.AddSignature(myAddress, sig)

	err = flowClient.SendTransaction(ctx, deployContractTx)
	utils.Handle(err)

	deployContractTxResp := utils.WaitForSeal(ctx, flowClient, deployContractTx.Hash())

	var nftAddress flow.Address

	for _, event := range deployContractTxResp.Events {
		if event.Type == flow.EventAccountCreated {
			accountCreatedEvent, err := flow.DecodeAccountCreatedEvent(event.Payload)
			utils.Handle(err)

			nftAddress = accountCreatedEvent.Address()
		}
	}

	fmt.Println("My Address:", nftAddress.Hex())

	// Next, instantiate the minter
	createMinterTx := flow.Transaction{
		Script:         contracts.GenerateCreateMinterScript(nftAddress, 1, 2),
		Nonce:          utils.GetNonce(),
		ComputeLimit:   10,
		PayerAccount:   myAddress,
		ScriptAccounts: []flow.Address{myAddress},
	}

	sig, err = keys.SignTransaction(createMinterTx, myPrivateKey)
	utils.Handle(err)

	createMinterTx.AddSignature(myAddress, sig)

	err = flowClient.SendTransaction(ctx, createMinterTx)
	utils.Handle(err)

	// Mint the NFT
	mintTx := flow.Transaction{
		Script:         contracts.GenerateMintScript(nftAddress),
		Nonce:          utils.GetNonce(),
		ComputeLimit:   10,
		PayerAccount:   myAddress,
		ScriptAccounts: []flow.Address{myAddress},
	}

	sig, err = keys.SignTransaction(mintTx, myPrivateKey)
	utils.Handle(err)

	mintTx.AddSignature(myAddress, sig)

	err = flowClient.SendTransaction(ctx, mintTx)
	utils.Handle(err)

	utils.WaitForSeal(ctx, flowClient, mintTx.Hash())

	fmt.Println("NFT minted!")

	result, err := flowClient.ExecuteScript(ctx, contracts.GenerateGetNFTIDScript(nftAddress, myAddress))
	utils.Handle(err)

	myTokenID, err := encoding.Decode(cadence.IntType{}, result)
	utils.Handle(err)

	id := myTokenID.(cadence.Int)

	fmt.Printf("You now own the Great NFT with ID: %d\n", id.Int())
}

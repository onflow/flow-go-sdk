package main

import (
	"context"
	"fmt"

	"github.com/dapperlabs/cadence"
	encoding "github.com/dapperlabs/cadence/encoding/xdr"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/client"
	"github.com/dapperlabs/flow-go-sdk/examples"
	"github.com/dapperlabs/flow-go-sdk/keys"
	"github.com/dapperlabs/flow-go-sdk/templates"
)

const GreatTokenContractFile = "../great-token.cdc"

func main() {
	DeployContractDemo()
}

func DeployContractDemo() {
	// Connect to an emulator running locally
	ctx := context.Background()
	flowClient, err := client.New("127.0.0.1:3569")
	examples.Handle(err)

	myPrivateKey := examples.RandomPrivateKey()
	myPublicKey := myPrivateKey.PublicKey(keys.PublicKeyWeightThreshold)

	// Generate an account creation script
	createAccountScript, err := templates.CreateAccount([]flow.AccountPublicKey{myPublicKey}, nil)
	examples.Handle(err)

	rootAcctAddr, rootAcctKey := examples.RootAccount()

	createAccountTx := flow.Transaction{
		Script:       createAccountScript,
		Nonce:        examples.GetNonce(),
		ComputeLimit: 10,
		PayerAccount: rootAcctAddr,
	}

	sig, err := keys.SignTransaction(createAccountTx, rootAcctKey)
	examples.Handle(err)

	createAccountTx.AddSignature(rootAcctAddr, sig)

	err = flowClient.SendTransaction(ctx, createAccountTx)
	examples.Handle(err)

	accountCreationTxRes := examples.WaitForSeal(ctx, flowClient, createAccountTx.Hash())

	var myAddress flow.Address

	for _, event := range accountCreationTxRes.Events {
		if event.Type == flow.EventAccountCreated {
			accountCreatedEvent, err := flow.DecodeAccountCreatedEvent(event.Payload)
			examples.Handle(err)

			myAddress = accountCreatedEvent.Address()
		}
	}

	fmt.Println("My Address:", myAddress.Hex())

	// Deploy the Great NFT contract
	nftCode := examples.ReadFile(GreatTokenContractFile)
	deployScript, err := templates.CreateAccount(nil, nftCode)

	deployContractTx := flow.Transaction{
		Script:       deployScript,
		Nonce:        examples.GetNonce(),
		ComputeLimit: 10,
		PayerAccount: myAddress,
	}
	sig, err = keys.SignTransaction(deployContractTx, myPrivateKey)
	examples.Handle(err)

	deployContractTx.AddSignature(myAddress, sig)

	err = flowClient.SendTransaction(ctx, deployContractTx)
	examples.Handle(err)

	deployContractTxResp := examples.WaitForSeal(ctx, flowClient, deployContractTx.Hash())

	var nftAddress flow.Address

	for _, event := range deployContractTxResp.Events {
		if event.Type == flow.EventAccountCreated {
			accountCreatedEvent, err := flow.DecodeAccountCreatedEvent(event.Payload)
			examples.Handle(err)

			nftAddress = accountCreatedEvent.Address()
		}
	}

	fmt.Println("My Address:", nftAddress.Hex())

	// Next, instantiate the minter
	createMinterTx := flow.Transaction{
		Script:         GenerateCreateMinterScript(nftAddress, 1, 2),
		Nonce:          examples.GetNonce(),
		ComputeLimit:   10,
		PayerAccount:   myAddress,
		ScriptAccounts: []flow.Address{myAddress},
	}

	sig, err = keys.SignTransaction(createMinterTx, myPrivateKey)
	examples.Handle(err)

	createMinterTx.AddSignature(myAddress, sig)

	err = flowClient.SendTransaction(ctx, createMinterTx)
	examples.Handle(err)

	// Mint the NFT
	mintTx := flow.Transaction{
		Script:         GenerateMintScript(nftAddress),
		Nonce:          examples.GetNonce(),
		ComputeLimit:   10,
		PayerAccount:   myAddress,
		ScriptAccounts: []flow.Address{myAddress},
	}

	sig, err = keys.SignTransaction(mintTx, myPrivateKey)
	examples.Handle(err)

	mintTx.AddSignature(myAddress, sig)

	err = flowClient.SendTransaction(ctx, mintTx)
	examples.Handle(err)

	examples.WaitForSeal(ctx, flowClient, mintTx.Hash())

	fmt.Println("NFT minted!")

	result, err := flowClient.ExecuteScript(ctx, GenerateGetNFTIDScript(nftAddress, myAddress))
	examples.Handle(err)

	myTokenID, err := encoding.Decode(cadence.IntType{}, result)
	examples.Handle(err)

	id := myTokenID.(cadence.Int)

	fmt.Printf("You now own the Great NFT with ID: %d\n", id.Int())
}

// GenerateCreateMinterScript Creates a script that instantiates
// a new GreatNFTMinter instance and stores it in memory.
// Initial ID and special mod are arguments to the GreatNFTMinter constructor.
// The GreatNFTMinter must have been deployed already.
func GenerateCreateMinterScript(nftAddr flow.Address, initialID, specialMod int) []byte {
	template := `
		import GreatToken from 0x%s

		transaction {

		  prepare(acct: AuthAccount) {
			let existing <- acct.storage[GreatToken.GreatNFTMinter] <- GreatToken.createGreatNFTMinter(firstID: %d, specialMod: %d)
			assert(existing == nil, message: "existed")
			destroy existing

			acct.storage[&GreatToken.GreatNFTMinter] = &acct.storage[GreatToken.GreatNFTMinter] as &GreatToken.GreatNFTMinter
		  }
		}
	`

	return []byte(fmt.Sprintf(template, nftAddr, initialID, specialMod))
}

// GenerateMintScript Creates a script that mints an NFT and put it into storage.
// The minter must have been instantiated already.
func GenerateMintScript(nftCodeAddr flow.Address) []byte {
	template := `
		import GreatToken from 0x%s

		transaction {

		  prepare(acct: AuthAccount) {
			let minter = acct.storage[&GreatToken.GreatNFTMinter] ?? panic("missing minter")
			let existing <- acct.storage[GreatToken.GreatNFT] <- minter.mint()
			destroy existing
            acct.published[&GreatToken.GreatNFT] = &acct.storage[GreatToken.GreatNFT] as &GreatToken.GreatNFT
		  }
		}
	`

	return []byte(fmt.Sprintf(template, nftCodeAddr.String()))
}

// GenerateGetNFTIDScript creates a script that retrieves an NFT from storage and returns its ID.
func GenerateGetNFTIDScript(nftCodeAddr, userAddr flow.Address) []byte {
	template := `
		import GreatToken from 0x%s

		pub fun main(): Int {
		  let acct = getAccount(0x%s)
		  let nft = acct.published[&GreatToken.GreatNFT] ?? panic("missing nft")
		  return nft.id()
		}
	`

	return []byte(fmt.Sprintf(template, nftCodeAddr, userAddr))
}

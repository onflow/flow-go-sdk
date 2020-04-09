package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/dapperlabs/cadence"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/client"
	"github.com/dapperlabs/flow-go-sdk/examples"
	"github.com/dapperlabs/flow-go-sdk/keys"
	"github.com/dapperlabs/flow-go-sdk/templates"
)

const GreatTokenContractFile = "./great-token.cdc"

func main() {
	DeployContractDemo()
}

func DeployContractDemo() {
	// Connect to an emulator running locally
	ctx := context.Background()
	flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
	examples.Handle(err)

	rootAcctAddr, rootAcctKey, rootPrivateKey := examples.RootAccount(flowClient)

	myPrivateKey := examples.RandomPrivateKey()
	myAcctKey := myPrivateKey.ToAccountKey()
	myAcctKey.Weight = keys.PublicKeyWeightThreshold

	// Generate an account creation script
	createAccountScript, err := templates.CreateAccount([]flow.AccountKey{myAcctKey}, nil)
	examples.Handle(err)

	createAccountTx := flow.NewTransaction().
		SetScript(createAccountScript).
		SetProposalKey(rootAcctAddr, rootAcctKey.ID, rootAcctKey.SequenceNumber).
		SetPayer(rootAcctAddr, rootAcctKey.ID)

	err = createAccountTx.SignContainer(
		rootAcctAddr,
		rootAcctKey.ID,
		rootPrivateKey.Signer(),
	)
	examples.Handle(err)

	err = flowClient.SendTransaction(ctx, *createAccountTx)
	examples.Handle(err)

	accountCreationTxRes := examples.WaitForSeal(ctx, flowClient, createAccountTx.ID())
	examples.Handle(accountCreationTxRes.Error)

	// Successful Tx, increment sequence number
	rootAcctKey.SequenceNumber++

	var myAddress flow.Address

	for _, event := range accountCreationTxRes.Events {
		if event.Type == flow.EventAccountCreated {
			accountCreatedEvent := flow.AccountCreatedEvent(event)
			myAddress = accountCreatedEvent.Address()
		}
	}

	fmt.Println("My Address:", myAddress.Hex())

	// Deploy the Great NFT contract
	nftCode := examples.ReadFile(GreatTokenContractFile)
	deployScript, err := templates.CreateAccount(nil, nftCode)

	deployContractTx := flow.NewTransaction().
		SetScript(deployScript).
		SetProposalKey(myAddress, myAcctKey.ID, myAcctKey.SequenceNumber).
		SetPayer(myAddress, myAcctKey.ID)

	err = deployContractTx.SignContainer(
		myAddress,
		myAcctKey.ID,
		myPrivateKey.Signer(),
	)
	examples.Handle(err)

	err = flowClient.SendTransaction(ctx, *deployContractTx)
	examples.Handle(err)

	deployContractTxResp := examples.WaitForSeal(ctx, flowClient, deployContractTx.ID())
	examples.Handle(deployContractTxResp.Error)

	// Successful Tx, increment sequence number
	myAcctKey.SequenceNumber++

	var nftAddress flow.Address

	for _, event := range deployContractTxResp.Events {
		if event.Type == flow.EventAccountCreated {
			accountCreatedEvent := flow.AccountCreatedEvent(event)
			nftAddress = accountCreatedEvent.Address()
		}
	}

	fmt.Println("My Address:", nftAddress.Hex())

	// Next, instantiate the minter
	createMinterScript := GenerateCreateMinterScript(nftAddress, 1, 2)

	createMinterTx := flow.NewTransaction().
		SetScript(createMinterScript).
		SetProposalKey(myAddress, myAcctKey.ID, myAcctKey.SequenceNumber).
		SetPayer(myAddress, myAcctKey.ID).
		AddAuthorizer(myAddress, myAcctKey.ID)

	err = createMinterTx.SignContainer(
		myAddress,
		myAcctKey.ID,
		myPrivateKey.Signer(),
	)
	examples.Handle(err)

	err = flowClient.SendTransaction(ctx, *createMinterTx)
	examples.Handle(err)

	createMinterTxResp := examples.WaitForSeal(ctx, flowClient, deployContractTx.ID())
	examples.Handle(createMinterTxResp.Error)

	// Successful Tx, increment sequence number
	myAcctKey.SequenceNumber++

	mintScript := GenerateMintScript(nftAddress)

	// Mint the NFT
	mintTx := flow.NewTransaction().
		SetScript(mintScript).
		SetProposalKey(myAddress, myAcctKey.ID, myAcctKey.SequenceNumber).
		SetPayer(myAddress, myAcctKey.ID).
		AddAuthorizer(myAddress, myAcctKey.ID)

	err = mintTx.SignContainer(
		myAddress,
		myAcctKey.ID,
		myPrivateKey.Signer(),
	)
	examples.Handle(err)

	err = flowClient.SendTransaction(ctx, *mintTx)
	examples.Handle(err)

	mintTxResp := examples.WaitForSeal(ctx, flowClient, mintTx.ID())
	examples.Handle(mintTxResp.Error)

	// Successful Tx, increment sequence number
	myAcctKey.SequenceNumber++

	fmt.Println("NFT minted!")

	result, err := flowClient.ExecuteScriptAtLatestBlock(ctx, GenerateGetNFTIDScript(nftAddress, myAddress))
	examples.Handle(err)

	myTokenID := result.(cadence.Int)

	fmt.Printf("You now own the Great NFT with ID: %d\n", myTokenID.Int())
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

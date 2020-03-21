package contracts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/utils/examples"
)

const (
	greatTokenContractFile = "./contracts/great-token.cdc"
)

func TestDeployment(t *testing.T) {
	b := examples.NewEmulator()

	// Should be able to deploy a contract as a new account with no keys.
	nftCode := examples.ReadFile(greatTokenContractFile)
	_, err := b.CreateAccount(nil, nftCode, examples.GetNonce())
	assert.NoError(t, err)
	_, err = b.CommitBlock()
	assert.NoError(t, err)
}

func TestCreateMinter(t *testing.T) {
	b := examples.NewEmulator()

	// First, deploy the contract
	nftCode := examples.ReadFile(greatTokenContractFile)
	contractAddr, err := b.CreateAccount(nil, nftCode, examples.GetNonce())
	assert.NoError(t, err)

	// GreatNFTMinter must be instantiated with initialID > 0 and
	// specialMod > 1
	t.Run("Cannot create minter with negative initial ID", func(t *testing.T) {
		tx := flow.Transaction{
			Script:         GenerateCreateMinterScript(contractAddr, -1, 2),
			Nonce:          examples.GetNonce(),
			ComputeLimit:   10,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{b.RootAccountAddress()},
		}

		examples.SignAndSubmit(t, b, tx, []flow.AccountPrivateKey{b.RootKey()}, []flow.Address{b.RootAccountAddress()}, true)
	})

	t.Run("Cannot create minter with special mod < 2", func(t *testing.T) {
		tx := flow.Transaction{
			Script:         GenerateCreateMinterScript(contractAddr, 1, 1),
			Nonce:          examples.GetNonce(),
			ComputeLimit:   10,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{b.RootAccountAddress()},
		}

		examples.SignAndSubmit(t, b, tx, []flow.AccountPrivateKey{b.RootKey()}, []flow.Address{b.RootAccountAddress()}, true)
	})

	t.Run("Should be able to create minter", func(t *testing.T) {
		tx := flow.Transaction{
			Script:         GenerateCreateMinterScript(contractAddr, 1, 2),
			Nonce:          examples.GetNonce(),
			ComputeLimit:   10,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{b.RootAccountAddress()},
		}

		examples.SignAndSubmit(t, b, tx, []flow.AccountPrivateKey{b.RootKey()}, []flow.Address{b.RootAccountAddress()}, false)
	})
}

func TestMinting(t *testing.T) {
	b := examples.NewEmulator()

	// First, deploy the contract
	nftCode := examples.ReadFile(greatTokenContractFile)
	contractAddr, err := b.CreateAccount(nil, nftCode, examples.GetNonce())
	assert.NoError(t, err)

	// Next, instantiate the minter
	createMinterTx := flow.Transaction{
		Script:         GenerateCreateMinterScript(contractAddr, 1, 2),
		Nonce:          examples.GetNonce(),
		ComputeLimit:   10,
		PayerAccount:   b.RootAccountAddress(),
		ScriptAccounts: []flow.Address{b.RootAccountAddress()},
	}

	examples.SignAndSubmit(t, b, createMinterTx, []flow.AccountPrivateKey{b.RootKey()}, []flow.Address{b.RootAccountAddress()}, false)

	// Mint the first NFT
	mintTx := flow.Transaction{
		Script:         GenerateMintScript(contractAddr),
		Nonce:          examples.GetNonce(),
		ComputeLimit:   10,
		PayerAccount:   b.RootAccountAddress(),
		ScriptAccounts: []flow.Address{b.RootAccountAddress()},
	}

	examples.SignAndSubmit(t, b, mintTx, []flow.AccountPrivateKey{b.RootKey()}, []flow.Address{b.RootAccountAddress()}, false)

	// Assert that ID/specialness are correct
	result, err := b.ExecuteScript(GenerateInspectNFTScript(contractAddr, b.RootAccountAddress(), 1, false))
	require.NoError(t, err)
	if !assert.True(t, result.Succeeded()) {
		t.Log(result.Error.Error())
	}

	// Mint a second NF
	mintTx2 := flow.Transaction{
		Script:         GenerateMintScript(contractAddr),
		Nonce:          examples.GetNonce(),
		ComputeLimit:   10,
		PayerAccount:   b.RootAccountAddress(),
		ScriptAccounts: []flow.Address{b.RootAccountAddress()},
	}

	examples.SignAndSubmit(t, b, mintTx2, []flow.AccountPrivateKey{b.RootKey()}, []flow.Address{b.RootAccountAddress()}, false)

	// Assert that ID/specialness are correct
	result, err = b.ExecuteScript(GenerateInspectNFTScript(contractAddr, b.RootAccountAddress(), 2, true))
	require.NoError(t, err)
	if !assert.True(t, result.Succeeded()) {
		t.Log(result.Error.Error())
	}
}

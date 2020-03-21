package contracts

import (
	"testing"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/keys"
	"github.com/dapperlabs/flow-go-sdk/utils/examples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	MarketContractFile = "./contracts/market.cdc"
)

func TestMarketDeployment(t *testing.T) {
	b := examples.NewEmulator()

	// Should be able to deploy a contract as a new account with no keys.
	tokenCode := examples.ReadFile(resourceTokenContractFile)
	_, err := b.CreateAccount(nil, tokenCode, examples.GetNonce())
	if !assert.Nil(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	require.NoError(t, err)

	nftCode := examples.ReadFile(NFTContractFile)
	_, err = b.CreateAccount(nil, nftCode, examples.GetNonce())
	if !assert.Nil(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	require.NoError(t, err)

	// Should be able to deploy a contract as a new account with no keys.
	marketCode := examples.ReadFile(MarketContractFile)
	_, err = b.CreateAccount(nil, marketCode, examples.GetNonce())
	if !assert.Nil(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	require.NoError(t, err)
}

func TestCreateSale(t *testing.T) {
	b := examples.NewEmulator()

	// first deploy the FT, NFT, and market code
	tokenCode := examples.ReadFile(resourceTokenContractFile)
	tokenAddr, err := b.CreateAccount(nil, tokenCode, examples.GetNonce())
	assert.Nil(t, err)
	_, err = b.CommitBlock()
	require.NoError(t, err)

	nftCode := examples.ReadFile(NFTContractFile)
	nftAddr, err := b.CreateAccount(nil, nftCode, examples.GetNonce())
	assert.Nil(t, err)
	_, err = b.CommitBlock()
	require.NoError(t, err)

	marketCode := examples.ReadFile(MarketContractFile)
	marketAddr, err := b.CreateAccount(nil, marketCode, examples.GetNonce())
	assert.Nil(t, err)
	_, err = b.CommitBlock()
	require.NoError(t, err)

	// create two new accounts
	bastianPrivateKey := examples.RandomPrivateKey()
	bastianPublicKey := bastianPrivateKey.PublicKey(keys.PublicKeyWeightThreshold)
	bastianAddress, err := b.CreateAccount([]flow.AccountPublicKey{bastianPublicKey}, nil, examples.GetNonce())
	joshPrivateKey := examples.RandomPrivateKey()
	joshPublicKey := joshPrivateKey.PublicKey(keys.PublicKeyWeightThreshold)
	joshAddress, err := b.CreateAccount([]flow.AccountPublicKey{joshPublicKey}, nil, examples.GetNonce())

	t.Run("Should be able to create FTs and NFT collections in each accounts storage", func(t *testing.T) {
		// create Fungible tokens and NFTs in each accounts storage and store references
		setupUsersTokens(t, b, tokenAddr, nftAddr, []flow.AccountPrivateKey{bastianPrivateKey, joshPrivateKey}, []flow.Address{bastianAddress, joshAddress})
	})

	t.Run("Can create sale collection", func(t *testing.T) {
		tx := flow.Transaction{
			Script:         GenerateCreateSaleScript(tokenAddr, marketAddr),
			Nonce:          examples.GetNonce(),
			ComputeLimit:   10,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{bastianAddress},
		}

		examples.SignAndSubmit(t, b, tx, []flow.AccountPrivateKey{b.RootKey(), bastianPrivateKey}, []flow.Address{b.RootAccountAddress(), bastianAddress}, false)
	})

	t.Run("Can put an NFT up for sale", func(t *testing.T) {
		tx := flow.Transaction{
			Script:         GenerateStartSaleScript(nftAddr, marketAddr, 1, 10),
			Nonce:          examples.GetNonce(),
			ComputeLimit:   10,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{bastianAddress},
		}

		examples.SignAndSubmit(t, b, tx, []flow.AccountPrivateKey{b.RootKey(), bastianPrivateKey}, []flow.Address{b.RootAccountAddress(), bastianAddress}, false)
	})

	t.Run("Cannot buy an NFT for less than the sale price", func(t *testing.T) {
		tx := flow.Transaction{
			Script:         GenerateBuySaleScript(tokenAddr, nftAddr, marketAddr, bastianAddress, 1, 9),
			Nonce:          examples.GetNonce(),
			ComputeLimit:   10,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{joshAddress},
		}

		examples.SignAndSubmit(t, b, tx, []flow.AccountPrivateKey{b.RootKey(), joshPrivateKey}, []flow.Address{b.RootAccountAddress(), joshAddress}, true)
	})

	t.Run("Cannot buy an NFT that is not for sale", func(t *testing.T) {
		tx := flow.Transaction{
			Script:         GenerateBuySaleScript(tokenAddr, nftAddr, marketAddr, bastianAddress, 2, 10),
			Nonce:          examples.GetNonce(),
			ComputeLimit:   10,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{joshAddress},
		}

		examples.SignAndSubmit(t, b, tx, []flow.AccountPrivateKey{b.RootKey(), joshPrivateKey}, []flow.Address{b.RootAccountAddress(), joshAddress}, true)
	})

	t.Run("Can buy an NFT that is for sale", func(t *testing.T) {
		tx := flow.Transaction{
			Script:         GenerateBuySaleScript(tokenAddr, nftAddr, marketAddr, bastianAddress, 1, 10),
			Nonce:          examples.GetNonce(),
			ComputeLimit:   10,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{joshAddress},
		}

		examples.SignAndSubmit(t, b, tx, []flow.AccountPrivateKey{b.RootKey(), joshPrivateKey}, []flow.Address{b.RootAccountAddress(), joshAddress}, false)

		result, err := b.ExecuteScript(GenerateInspectVaultScript(tokenAddr, bastianAddress, 40))
		require.NoError(t, err)
		if !assert.True(t, result.Succeeded()) {
			t.Log(result.Error.Error())
		}

		result, err = b.ExecuteScript(GenerateInspectVaultScript(tokenAddr, joshAddress, 20))
		require.NoError(t, err)
		if !assert.True(t, result.Succeeded()) {
			t.Log(result.Error.Error())
		}

		// Assert that the accounts' collections are correct
		result, err = b.ExecuteScript(GenerateInspectCollectionScript(nftAddr, bastianAddress, 1, false))
		require.NoError(t, err)
		if !assert.True(t, result.Succeeded()) {
			t.Log(result.Error.Error())
		}

		result, err = b.ExecuteScript(GenerateInspectCollectionScript(nftAddr, bastianAddress, 2, false))
		require.NoError(t, err)
		if !assert.True(t, result.Succeeded()) {
			t.Log(result.Error.Error())
		}

		result, err = b.ExecuteScript(GenerateInspectCollectionScript(nftAddr, joshAddress, 1, true))
		require.NoError(t, err)
		if !assert.True(t, result.Succeeded()) {
			t.Log(result.Error.Error())
		}

		result, err = b.ExecuteScript(GenerateInspectCollectionScript(nftAddr, joshAddress, 2, true))
		require.NoError(t, err)
		if !assert.True(t, result.Succeeded()) {
			t.Log(result.Error.Error())
		}
	})

}

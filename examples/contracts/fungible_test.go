package contracts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/keys"
	"github.com/dapperlabs/flow-go-sdk/utils/examples"
)

const (
	resourceTokenContractFile = "./contracts/fungible-token.cdc"
)

func TestTokenDeployment(t *testing.T) {
	b := examples.NewEmulator()

	// Should be able to deploy a contract as a new account with no keys.
	tokenCode := examples.ReadFile(resourceTokenContractFile)
	_, err := b.CreateAccount(nil, tokenCode, examples.GetNonce())
	assert.NoError(t, err)
	_, err = b.CommitBlock()
	assert.NoError(t, err)
}

func TestCreateToken(t *testing.T) {
	b := examples.NewEmulator()

	// First, deploy the contract
	tokenCode := examples.ReadFile(resourceTokenContractFile)
	contractAddr, err := b.CreateAccount(nil, tokenCode, examples.GetNonce())
	assert.NoError(t, err)

	// Vault must be instantiated with a positive balance
	t.Run("Cannot create token with negative initial balance", func(t *testing.T) {
		tx := flow.Transaction{
			Script:         GenerateCreateTokenScript(contractAddr, -7),
			Nonce:          examples.GetNonce(),
			ComputeLimit:   10,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{b.RootAccountAddress()},
		}

		examples.SignAndSubmit(t, b, tx, []flow.AccountPrivateKey{b.RootKey()}, []flow.Address{b.RootAccountAddress()}, true)
	})

	t.Run("Should be able to create token", func(t *testing.T) {
		tx := flow.Transaction{
			Script:         GenerateCreateTokenScript(contractAddr, 10),
			Nonce:          examples.GetNonce(),
			ComputeLimit:   20,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{b.RootAccountAddress()},
		}

		examples.SignAndSubmit(t, b, tx, []flow.AccountPrivateKey{b.RootKey()}, []flow.Address{b.RootAccountAddress()}, false)

		result, err := b.ExecuteScript(GenerateInspectVaultScript(contractAddr, b.RootAccountAddress(), 10))
		require.NoError(t, err)
		if !assert.True(t, result.Succeeded()) {
			t.Log(result.Error.Error())
		}
	})

	t.Run("Should be able to create multiple tokens and store them in an array", func(t *testing.T) {
		tx := flow.Transaction{
			Script:         GenerateCreateThreeTokensArrayScript(contractAddr, 10, 20, 5),
			Nonce:          examples.GetNonce(),
			ComputeLimit:   20,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{b.RootAccountAddress()},
		}

		examples.SignAndSubmit(t, b, tx, []flow.AccountPrivateKey{b.RootKey()}, []flow.Address{b.RootAccountAddress()}, false)
	})
}

func TestInAccountTransfers(t *testing.T) {
	b := examples.NewEmulator()

	// First, deploy the contract
	tokenCode := examples.ReadFile(resourceTokenContractFile)
	contractAddr, err := b.CreateAccount(nil, tokenCode, examples.GetNonce())
	assert.NoError(t, err)

	// then deploy the three tokens to an account
	tx := flow.Transaction{
		Script:         GenerateCreateThreeTokensArrayScript(contractAddr, 10, 20, 5),
		Nonce:          examples.GetNonce(),
		ComputeLimit:   20,
		PayerAccount:   b.RootAccountAddress(),
		ScriptAccounts: []flow.Address{b.RootAccountAddress()},
	}

	examples.SignAndSubmit(t, b, tx, []flow.AccountPrivateKey{b.RootKey()}, []flow.Address{b.RootAccountAddress()}, false)

	t.Run("Should be able to withdraw tokens from a vault", func(t *testing.T) {
		tx := flow.Transaction{
			Script:         GenerateWithdrawScript(contractAddr, 0, 3),
			Nonce:          examples.GetNonce(),
			ComputeLimit:   20,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{b.RootAccountAddress()},
		}

		examples.SignAndSubmit(t, b, tx, []flow.AccountPrivateKey{b.RootKey()}, []flow.Address{b.RootAccountAddress()}, false)

		// Assert that the vaults balance is correct
		result, err := b.ExecuteScript(GenerateInspectVaultArrayScript(contractAddr, b.RootAccountAddress(), 0, 7))
		require.NoError(t, err)
		if !assert.True(t, result.Succeeded()) {
			t.Log(result.Error.Error())
		}
	})

	t.Run("Should be able to withdraw and deposit tokens from one vault to another in an account", func(t *testing.T) {

		tx = flow.Transaction{
			Script:         GenerateWithdrawDepositScript(contractAddr, 1, 2, 8),
			Nonce:          examples.GetNonce(),
			ComputeLimit:   20,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{b.RootAccountAddress()},
		}

		examples.SignAndSubmit(t, b, tx, []flow.AccountPrivateKey{b.RootKey()}, []flow.Address{b.RootAccountAddress()}, false)

		// Assert that the vault's balance is correct
		result, err := b.ExecuteScript(GenerateInspectVaultArrayScript(contractAddr, b.RootAccountAddress(), 1, 12))
		require.NoError(t, err)
		if !assert.True(t, result.Succeeded()) {
			t.Log(result.Error.Error())
		}

		// Assert that the vault's balance is correct
		result, err = b.ExecuteScript(GenerateInspectVaultArrayScript(contractAddr, b.RootAccountAddress(), 2, 13))
		require.NoError(t, err)
		if !assert.True(t, result.Succeeded()) {
			t.Log(result.Error.Error())
		}
	})
}

func TestExternalTransfers(t *testing.T) {
	b := examples.NewEmulator()

	// First, deploy the token contract
	tokenCode := examples.ReadFile(resourceTokenContractFile)
	contractAddr, err := b.CreateAccount(nil, tokenCode, examples.GetNonce())
	assert.NoError(t, err)

	// then deploy the tokens to an account
	tx := flow.Transaction{
		Script:         GenerateCreateTokenScript(contractAddr, 10),
		Nonce:          examples.GetNonce(),
		ComputeLimit:   20,
		PayerAccount:   b.RootAccountAddress(),
		ScriptAccounts: []flow.Address{b.RootAccountAddress()},
	}

	examples.SignAndSubmit(t, b, tx, []flow.AccountPrivateKey{b.RootKey()}, []flow.Address{b.RootAccountAddress()}, false)

	// create a new account
	bastianPrivateKey := examples.RandomPrivateKey()
	bastianPublicKey := bastianPrivateKey.PublicKey(keys.PublicKeyWeightThreshold)
	bastianAddress, err := b.CreateAccount([]flow.AccountPublicKey{bastianPublicKey}, nil, examples.GetNonce())

	// then deploy the tokens to the new account
	tx = flow.Transaction{
		Script:         GenerateCreateTokenScript(contractAddr, 10),
		Nonce:          examples.GetNonce(),
		ComputeLimit:   20,
		PayerAccount:   b.RootAccountAddress(),
		ScriptAccounts: []flow.Address{bastianAddress},
	}

	examples.SignAndSubmit(t, b, tx, []flow.AccountPrivateKey{b.RootKey(), bastianPrivateKey}, []flow.Address{b.RootAccountAddress(), bastianAddress}, false)

	t.Run("Should be able to withdraw and deposit tokens from a vault", func(t *testing.T) {
		tx := flow.Transaction{
			Script:         GenerateDepositVaultScript(contractAddr, bastianAddress, 3),
			Nonce:          examples.GetNonce(),
			ComputeLimit:   20,
			PayerAccount:   b.RootAccountAddress(),
			ScriptAccounts: []flow.Address{b.RootAccountAddress()},
		}

		examples.SignAndSubmit(t, b, tx, []flow.AccountPrivateKey{b.RootKey()}, []flow.Address{b.RootAccountAddress()}, false)

		// Assert that the vaults' balances are correct
		result, err := b.ExecuteScript(GenerateInspectVaultScript(contractAddr, b.RootAccountAddress(), 7))
		require.NoError(t, err)
		if !assert.True(t, result.Succeeded()) {
			t.Log(result.Error.Error())
		}

		result, err = b.ExecuteScript(GenerateInspectVaultScript(contractAddr, bastianAddress, 13))
		require.NoError(t, err)
		if !assert.True(t, result.Succeeded()) {
			t.Log(result.Error.Error())
		}
	})
}

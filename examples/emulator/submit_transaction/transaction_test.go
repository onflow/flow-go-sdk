package submit_transaction_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/keys"

	utils "github.com/dapperlabs/flow-go-sdk/utils/examples"
)

func TestTransactionPayerAccount(t *testing.T) {
	b := utils.NewEmulator()

	rootAddress := b.RootAccountAddress()
	rootPrivateKey := b.RootKey()

	tx := flow.Transaction{
		Script: []byte(`
            transaction {
                execute {
                    log("Hello world!")
                }
            }
        `),
		Nonce:        1,
		ComputeLimit: 10,
		PayerAccount: rootAddress,
	}

	sig, err := keys.SignTransaction(tx, rootPrivateKey)
	require.Nil(t, err)

	tx.AddSignature(rootAddress, sig)

	err = b.AddTransaction(tx)
	require.Nil(t, err)

	result, err := b.ExecuteNextTransaction()
	require.Nil(t, err)

	assert.True(t, result.Succeeded())
}

func TestTransactionPayerAndScriptAccount(t *testing.T) {
	b := utils.NewEmulator()

	rootAddress := b.RootAccountAddress()
	rootPrivateKey := b.RootKey()

	bastianPrivateKey := utils.RandomPrivateKey()
	bastianPublicKey := bastianPrivateKey.PublicKey(keys.PublicKeyWeightThreshold)
	bastianAddress, err := b.CreateAccount([]flow.AccountPublicKey{bastianPublicKey}, nil, utils.GetNonce())

	tx := flow.Transaction{
		Script: []byte(`
            transaction {
                prepare(bastian: Account) {
                    log("Sending transaction with authorized account:")
                    log(bastian.address)
                }
            }
        `),
		Nonce:          1,
		ComputeLimit:   10,
		PayerAccount:   rootAddress,
		ScriptAccounts: []flow.Address{bastianAddress},
	}

	rootSig, err := keys.SignTransaction(tx, rootPrivateKey)
	require.Nil(t, err)

	tx.AddSignature(rootAddress, rootSig)

	bastianSig, err := keys.SignTransaction(tx, bastianPrivateKey)
	require.Nil(t, err)

	tx.AddSignature(bastianAddress, bastianSig)

	err = b.AddTransaction(tx)
	require.Nil(t, err)

	result, err := b.ExecuteNextTransaction()
	require.Nil(t, err)

	assert.True(t, result.Succeeded())
}

func TestTransactionPayerAndScriptAccounts(t *testing.T) {
	b := utils.NewEmulator()

	rootAddress := b.RootAccountAddress()
	rootPrivateKey := b.RootKey()

	bastianPrivateKey := utils.RandomPrivateKey()
	bastianPublicKey := bastianPrivateKey.PublicKey(keys.PublicKeyWeightThreshold)
	bastianAddress, err := b.CreateAccount([]flow.AccountPublicKey{bastianPublicKey}, nil, utils.GetNonce())

	jordanPrivateKey := utils.RandomPrivateKey()
	jordanPublicKey := jordanPrivateKey.PublicKey(keys.PublicKeyWeightThreshold)
	jordanAddress, err := b.CreateAccount([]flow.AccountPublicKey{jordanPublicKey}, nil, utils.GetNonce())

	tx := flow.Transaction{
		Script: []byte(`
            transaction {
                prepare(bastian: Account, jordan: Account) {
                    log("Sending transaction with authorized accounts:")
                    log(bastian.address)
                    log(jordan.address)
                }
            }
        `),
		Nonce:          1,
		ComputeLimit:   10,
		PayerAccount:   rootAddress,
		ScriptAccounts: []flow.Address{bastianAddress, jordanAddress},
	}

	rootSig, err := keys.SignTransaction(tx, rootPrivateKey)
	require.Nil(t, err)

	tx.AddSignature(rootAddress, rootSig)

	bastianSig, err := keys.SignTransaction(tx, bastianPrivateKey)
	require.Nil(t, err)

	tx.AddSignature(bastianAddress, bastianSig)

	jordanSig, err := keys.SignTransaction(tx, jordanPrivateKey)
	require.Nil(t, err)

	tx.AddSignature(jordanAddress, jordanSig)

	err = b.AddTransaction(tx)
	require.Nil(t, err)

	result, err := b.ExecuteNextTransaction()
	require.Nil(t, err)

	assert.True(t, result.Succeeded())
}

func TestTransactionMultipleKeys(t *testing.T) {
	b := utils.NewEmulator()

	laynePrivateKeyA := utils.RandomPrivateKey()
	laynePublicKeyA := laynePrivateKeyA.PublicKey(keys.PublicKeyWeightThreshold / 2)

	laynePrivateKeyB := utils.RandomPrivateKey()
	laynePublicKeyB := laynePrivateKeyB.PublicKey(keys.PublicKeyWeightThreshold / 2)

	layneAddress, err := b.CreateAccount([]flow.AccountPublicKey{laynePublicKeyA, laynePublicKeyB}, nil, utils.GetNonce())

	tx := flow.Transaction{
		Script: []byte(`
            transaction {
                prepare(layne: Account) {
                    log("Sending transaction with authorized account:")
                    log(layne.address)
                }
            }
        `),
		Nonce:          1,
		ComputeLimit:   10,
		PayerAccount:   layneAddress,
		ScriptAccounts: []flow.Address{layneAddress},
	}

	layneSigA, err := keys.SignTransaction(tx, laynePrivateKeyA)
	require.Nil(t, err)

	tx.AddSignature(layneAddress, layneSigA)

	layneSigB, err := keys.SignTransaction(tx, laynePrivateKeyB)
	require.Nil(t, err)

	tx.AddSignature(layneAddress, layneSigB)

	err = b.AddTransaction(tx)
	require.Nil(t, err)

	result, err := b.ExecuteNextTransaction()
	require.Nil(t, err)

	assert.True(t, result.Succeeded())
}

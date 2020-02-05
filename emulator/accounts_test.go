package emulator_test

import (
	"testing"

	"github.com/dapperlabs/flow-go/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/emulator"
	"github.com/dapperlabs/flow-go-sdk/keys"
	"github.com/dapperlabs/flow-go-sdk/templates"
	"github.com/dapperlabs/flow-go-sdk/utils/unittest"
)

const testContract = "pub contract Test {}"

func TestCreateAccount(t *testing.T) {
	publicKeys := unittest.PublicKeyFixtures()

	t.Run("SingleKey", func(t *testing.T) {
		b, err := emulator.NewBlockchain()
		require.NoError(t, err)

		publicKey := flow.AccountPublicKey{
			PublicKey: publicKeys[0],
			SignAlgo:  crypto.ECDSA_P256,
			HashAlgo:  crypto.SHA3_256,
			Weight:    keys.PublicKeyWeightThreshold,
		}

		createAccountScript, err := templates.CreateAccount([]flow.AccountPublicKey{publicKey}, nil)
		require.NoError(t, err)

		tx := flow.Transaction{
			Script:             createAccountScript,
			ReferenceBlockHash: nil,
			Nonce:              getNonce(),
			ComputeLimit:       10,
			PayerAccount:       b.RootAccountAddress(),
		}

		sig, err := keys.SignTransaction(tx, b.RootKey())
		assert.NoError(t, err)

		tx.AddSignature(b.RootAccountAddress(), sig)

		err = b.AddTransaction(tx)
		assert.NoError(t, err)

		result, err := b.ExecuteNextTransaction()
		assert.NoError(t, err)
		assert.True(t, result.Succeeded())

		_, err = b.CommitBlock()
		assert.NoError(t, err)

		account := b.LastCreatedAccount()

		assert.Equal(t, uint64(0), account.Balance)
		require.Len(t, account.Keys, 1)
		assert.Equal(t, publicKey, account.Keys[0])
		assert.Empty(t, account.Code)
	})

	t.Run("MultipleKeys", func(t *testing.T) {
		b, err := emulator.NewBlockchain()
		require.NoError(t, err)

		publicKeyA := flow.AccountPublicKey{
			PublicKey: publicKeys[0],
			SignAlgo:  crypto.ECDSA_P256,
			HashAlgo:  crypto.SHA3_256,
			Weight:    keys.PublicKeyWeightThreshold,
		}

		publicKeyB := flow.AccountPublicKey{
			PublicKey: publicKeys[1],
			SignAlgo:  crypto.ECDSA_P256,
			HashAlgo:  crypto.SHA3_256,
			Weight:    keys.PublicKeyWeightThreshold,
		}

		createAccountScript, err := templates.CreateAccount([]flow.AccountPublicKey{publicKeyA, publicKeyB}, nil)
		assert.NoError(t, err)

		tx := flow.Transaction{
			Script:             createAccountScript,
			ReferenceBlockHash: nil,
			Nonce:              getNonce(),
			ComputeLimit:       10,
			PayerAccount:       b.RootAccountAddress(),
		}

		sig, err := keys.SignTransaction(tx, b.RootKey())
		assert.NoError(t, err)

		tx.AddSignature(b.RootAccountAddress(), sig)

		err = b.AddTransaction(tx)
		assert.NoError(t, err)

		result, err := b.ExecuteNextTransaction()
		assert.NoError(t, err)
		assert.True(t, result.Succeeded())

		_, err = b.CommitBlock()
		assert.NoError(t, err)

		account := b.LastCreatedAccount()

		assert.Equal(t, uint64(0), account.Balance)
		require.Len(t, account.Keys, 2)
		assert.Equal(t, publicKeyA, account.Keys[0])
		assert.Equal(t, publicKeyB, account.Keys[1])
		assert.Empty(t, account.Code)
	})

	t.Run("KeysAndCode", func(t *testing.T) {
		b, err := emulator.NewBlockchain()
		require.NoError(t, err)

		publicKeyA := flow.AccountPublicKey{
			PublicKey: publicKeys[0],
			SignAlgo:  crypto.ECDSA_P256,
			HashAlgo:  crypto.SHA3_256,
			Weight:    keys.PublicKeyWeightThreshold,
		}

		publicKeyB := flow.AccountPublicKey{
			PublicKey: publicKeys[1],
			SignAlgo:  crypto.ECDSA_P256,
			HashAlgo:  crypto.SHA3_256,
			Weight:    keys.PublicKeyWeightThreshold,
		}

		code := []byte(testContract)

		createAccountScript, err := templates.CreateAccount([]flow.AccountPublicKey{publicKeyA, publicKeyB}, code)
		assert.NoError(t, err)

		tx := flow.Transaction{
			Script:             createAccountScript,
			ReferenceBlockHash: nil,
			Nonce:              getNonce(),
			ComputeLimit:       10,
			PayerAccount:       b.RootAccountAddress(),
		}

		sig, err := keys.SignTransaction(tx, b.RootKey())
		assert.NoError(t, err)

		tx.AddSignature(b.RootAccountAddress(), sig)

		err = b.AddTransaction(tx)
		assert.NoError(t, err)

		result, err := b.ExecuteNextTransaction()
		assert.NoError(t, err)
		assert.True(t, result.Succeeded())

		_, err = b.CommitBlock()
		assert.NoError(t, err)

		account := b.LastCreatedAccount()

		assert.Equal(t, uint64(0), account.Balance)
		require.Len(t, account.Keys, 2)
		assert.Equal(t, publicKeyA, account.Keys[0])
		assert.Equal(t, publicKeyB, account.Keys[1])
		assert.Equal(t, code, account.Code)
	})

	t.Run("CodeAndNoKeys", func(t *testing.T) {
		b, err := emulator.NewBlockchain()
		require.NoError(t, err)

		code := []byte(testContract)

		createAccountScript, err := templates.CreateAccount(nil, code)
		assert.NoError(t, err)

		tx := flow.Transaction{
			Script:             createAccountScript,
			ReferenceBlockHash: nil,
			Nonce:              getNonce(),
			ComputeLimit:       10,
			PayerAccount:       b.RootAccountAddress(),
		}

		sig, err := keys.SignTransaction(tx, b.RootKey())
		assert.NoError(t, err)

		tx.AddSignature(b.RootAccountAddress(), sig)

		err = b.AddTransaction(tx)
		assert.NoError(t, err)

		result, err := b.ExecuteNextTransaction()
		assert.NoError(t, err)
		assert.True(t, result.Succeeded())

		_, err = b.CommitBlock()
		assert.NoError(t, err)

		account := b.LastCreatedAccount()

		assert.Equal(t, uint64(0), account.Balance)
		assert.Empty(t, account.Keys)
		assert.Equal(t, code, account.Code)
	})

	t.Run("EventEmitted", func(t *testing.T) {
		b, err := emulator.NewBlockchain()
		require.NoError(t, err)

		publicKey := flow.AccountPublicKey{
			PublicKey: publicKeys[0],
			SignAlgo:  crypto.ECDSA_P256,
			HashAlgo:  crypto.SHA3_256,
			Weight:    keys.PublicKeyWeightThreshold,
		}

		code := []byte(testContract)

		createAccountScript, err := templates.CreateAccount([]flow.AccountPublicKey{publicKey}, code)
		assert.NoError(t, err)

		tx := flow.Transaction{
			Script:             createAccountScript,
			ReferenceBlockHash: nil,
			Nonce:              getNonce(),
			ComputeLimit:       10,
			PayerAccount:       b.RootAccountAddress(),
		}

		sig, err := keys.SignTransaction(tx, b.RootKey())
		assert.NoError(t, err)

		tx.AddSignature(b.RootAccountAddress(), sig)

		err = b.AddTransaction(tx)
		assert.NoError(t, err)

		result, err := b.ExecuteNextTransaction()
		assert.NoError(t, err)
		assert.True(t, result.Succeeded())

		block, err := b.CommitBlock()
		require.NoError(t, err)

		events, err := b.GetEvents(flow.EventAccountCreated, block.Number, block.Number)
		require.NoError(t, err)
		require.Len(t, events, 1)

		accountCreatedEvent, err := flow.DecodeAccountCreatedEvent(events[0].Payload)
		assert.Nil(t, err)
		accountAddress := accountCreatedEvent.Address()

		account, err := b.GetAccount(accountAddress)
		assert.NoError(t, err)

		assert.Equal(t, uint64(0), account.Balance)
		require.Len(t, account.Keys, 1)
		assert.Equal(t, publicKey, account.Keys[0])
		assert.Equal(t, code, account.Code)
	})

	t.Run("InvalidKeyHashingAlgorithm", func(t *testing.T) {
		b, err := emulator.NewBlockchain()
		require.NoError(t, err)

		lastAccount := b.LastCreatedAccount()

		publicKey := flow.AccountPublicKey{
			PublicKey: unittest.PublicKeyFixtures()[0],
			SignAlgo:  crypto.ECDSA_P256,
			// SHA2_384 is not compatible with ECDSA_P256
			HashAlgo: crypto.SHA2_384,
			Weight:   keys.PublicKeyWeightThreshold,
		}

		createAccountScript, err := templates.CreateAccount([]flow.AccountPublicKey{publicKey}, nil)
		require.NoError(t, err)

		tx := flow.Transaction{
			Script:             createAccountScript,
			ReferenceBlockHash: nil,
			Nonce:              getNonce(),
			ComputeLimit:       10,
			PayerAccount:       b.RootAccountAddress(),
		}

		sig, err := keys.SignTransaction(tx, b.RootKey())
		assert.NoError(t, err)

		tx.AddSignature(b.RootAccountAddress(), sig)

		err = b.AddTransaction(tx)
		assert.NoError(t, err)

		result, err := b.ExecuteNextTransaction()
		assert.NoError(t, err)
		assert.True(t, result.Reverted())

		newAccount := b.LastCreatedAccount()

		assert.Equal(t, lastAccount, newAccount)
	})

	t.Run("InvalidCode", func(t *testing.T) {
		b, err := emulator.NewBlockchain()
		require.NoError(t, err)

		lastAccount := b.LastCreatedAccount()

		code := []byte("not a valid script")

		createAccountScript, err := templates.CreateAccount(nil, code)
		assert.NoError(t, err)

		tx := flow.Transaction{
			Script:             createAccountScript,
			ReferenceBlockHash: nil,
			Nonce:              getNonce(),
			ComputeLimit:       10,
			PayerAccount:       b.RootAccountAddress(),
		}

		sig, err := keys.SignTransaction(tx, b.RootKey())
		assert.NoError(t, err)

		tx.AddSignature(b.RootAccountAddress(), sig)

		err = b.AddTransaction(tx)
		assert.NoError(t, err)

		result, err := b.ExecuteNextTransaction()
		assert.NoError(t, err)
		assert.True(t, result.Reverted())

		newAccount := b.LastCreatedAccount()

		assert.Equal(t, lastAccount, newAccount)
	})
}

func TestAddAccountKey(t *testing.T) {
	t.Run("ValidKey", func(t *testing.T) {
		b, err := emulator.NewBlockchain()
		require.NoError(t, err)

		privateKey, _ := keys.GeneratePrivateKey(keys.ECDSA_P256_SHA3_256,
			[]byte("elephant ears space cowboy octopus rodeo potato cannon pineapple"))
		publicKey := privateKey.PublicKey(keys.PublicKeyWeightThreshold)

		addKeyScript, err := templates.AddAccountKey(publicKey)
		assert.NoError(t, err)

		tx1 := flow.Transaction{
			Script:             addKeyScript,
			ReferenceBlockHash: nil,
			Nonce:              getNonce(),
			ComputeLimit:       10,
			PayerAccount:       b.RootAccountAddress(),
			ScriptAccounts:     []flow.Address{b.RootAccountAddress()},
		}

		sig, err := keys.SignTransaction(tx1, b.RootKey())
		assert.NoError(t, err)

		tx1.AddSignature(b.RootAccountAddress(), sig)

		err = b.AddTransaction(tx1)
		assert.NoError(t, err)

		result, err := b.ExecuteNextTransaction()
		assert.NoError(t, err)
		assert.True(t, result.Succeeded())

		_, err = b.CommitBlock()
		assert.NoError(t, err)

		script := []byte("transaction { execute {} }")

		tx2 := flow.Transaction{
			Script:             script,
			ReferenceBlockHash: nil,
			Nonce:              getNonce(),
			ComputeLimit:       10,
			PayerAccount:       b.RootAccountAddress(),
		}

		sig, err = keys.SignTransaction(tx2, privateKey)
		assert.NoError(t, err)

		tx2.AddSignature(b.RootAccountAddress(), sig)

		err = b.AddTransaction(tx2)
		assert.NoError(t, err)

		result, err = b.ExecuteNextTransaction()
		assert.NoError(t, err)
		assert.True(t, result.Succeeded())

		_, err = b.CommitBlock()
		assert.NoError(t, err)
	})

	t.Run("InvalidKeyHashingAlgorithm", func(t *testing.T) {
		b, err := emulator.NewBlockchain()
		require.NoError(t, err)

		publicKey := flow.AccountPublicKey{
			PublicKey: unittest.PublicKeyFixtures()[0],
			SignAlgo:  crypto.ECDSA_P256,
			// SHA2_384 is not compatible with ECDSA_P256
			HashAlgo: crypto.SHA2_384,
			Weight:   keys.PublicKeyWeightThreshold,
		}

		addKeyScript, err := templates.AddAccountKey(publicKey)
		assert.NoError(t, err)

		tx := flow.Transaction{
			Script:             addKeyScript,
			ReferenceBlockHash: nil,
			Nonce:              getNonce(),
			ComputeLimit:       10,
			PayerAccount:       b.RootAccountAddress(),
			ScriptAccounts:     []flow.Address{b.RootAccountAddress()},
		}

		sig, err := keys.SignTransaction(tx, b.RootKey())
		assert.NoError(t, err)

		tx.AddSignature(b.RootAccountAddress(), sig)

		err = b.AddTransaction(tx)
		assert.NoError(t, err)

		result, err := b.ExecuteNextTransaction()
		assert.NoError(t, err)
		assert.True(t, result.Reverted())
	})
}

func TestRemoveAccountKey(t *testing.T) {
	b, err := emulator.NewBlockchain()
	require.NoError(t, err)

	privateKey, _ := keys.GeneratePrivateKey(keys.ECDSA_P256_SHA3_256,
		[]byte("pineapple elephant ears space cowboy octopus rodeo potato cannon"))
	publicKey := privateKey.PublicKey(keys.PublicKeyWeightThreshold)

	addKeyScript, err := templates.AddAccountKey(publicKey)
	assert.NoError(t, err)

	// create transaction that adds publicKey to account keys
	tx1 := flow.Transaction{
		Script:             addKeyScript,
		ReferenceBlockHash: nil,
		Nonce:              getNonce(),
		ComputeLimit:       10,
		PayerAccount:       b.RootAccountAddress(),
		ScriptAccounts:     []flow.Address{b.RootAccountAddress()},
	}

	// sign with root key
	sig, err := keys.SignTransaction(tx1, b.RootKey())
	assert.NoError(t, err)

	tx1.AddSignature(b.RootAccountAddress(), sig)

	// submit tx1 (should succeed)
	err = b.AddTransaction(tx1)
	assert.NoError(t, err)

	result, err := b.ExecuteNextTransaction()
	assert.NoError(t, err)
	assert.True(t, result.Succeeded())

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	account, err := b.GetAccount(b.RootAccountAddress())
	assert.NoError(t, err)

	assert.Len(t, account.Keys, 2)

	// create transaction that removes root key
	tx2 := flow.Transaction{
		Script:             templates.RemoveAccountKey(0),
		ReferenceBlockHash: nil,
		Nonce:              getNonce(),
		ComputeLimit:       10,
		PayerAccount:       b.RootAccountAddress(),
		ScriptAccounts:     []flow.Address{b.RootAccountAddress()},
	}

	// sign with root key
	sig, err = keys.SignTransaction(tx2, b.RootKey())
	assert.NoError(t, err)

	tx2.AddSignature(b.RootAccountAddress(), sig)

	// submit tx2 (should succeed)
	err = b.AddTransaction(tx2)
	assert.NoError(t, err)

	result, err = b.ExecuteNextTransaction()
	assert.NoError(t, err)
	assert.True(t, result.Succeeded())

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	account, err = b.GetAccount(b.RootAccountAddress())
	assert.NoError(t, err)

	assert.Len(t, account.Keys, 1)

	// create transaction that removes remaining account key
	tx3 := flow.Transaction{
		Script:             templates.RemoveAccountKey(0),
		ReferenceBlockHash: nil,
		Nonce:              getNonce(),
		ComputeLimit:       10,
		PayerAccount:       b.RootAccountAddress(),
		ScriptAccounts:     []flow.Address{b.RootAccountAddress()},
	}

	// sign with root key (that has been removed)
	sig, err = keys.SignTransaction(tx3, b.RootKey())
	assert.NoError(t, err)

	tx3.AddSignature(b.RootAccountAddress(), sig)

	// submit tx3 (should fail)
	err = b.AddTransaction(tx3)
	assert.IsType(t, &emulator.ErrInvalidSignaturePublicKey{}, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	account, err = b.GetAccount(b.RootAccountAddress())
	assert.NoError(t, err)

	assert.Len(t, account.Keys, 1)

	// create transaction that removes remaining account key
	tx4 := flow.Transaction{
		Script:             templates.RemoveAccountKey(0),
		ReferenceBlockHash: nil,
		Nonce:              getNonce(),
		ComputeLimit:       10,
		PayerAccount:       b.RootAccountAddress(),
		ScriptAccounts:     []flow.Address{b.RootAccountAddress()},
	}

	// sign with remaining account key
	sig, err = keys.SignTransaction(tx4, privateKey)
	assert.NoError(t, err)

	tx4.AddSignature(b.RootAccountAddress(), sig)

	// submit tx4 (should succeed)
	err = b.AddTransaction(tx4)
	assert.NoError(t, err)

	result, err = b.ExecuteNextTransaction()
	assert.NoError(t, err)
	assert.True(t, result.Succeeded())

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	account, err = b.GetAccount(b.RootAccountAddress())
	assert.NoError(t, err)

	// no more keys left on account
	assert.Empty(t, account.Keys)
}

func TestUpdateAccountCode(t *testing.T) {
	codeA := []byte(`
      pub contract Test {
          pub fun a(): Int {
              return 1
          }
      }
    `)
	codeB := []byte(`
      pub contract Test {
          pub fun b(): Int {
              return 2
          }
      }
    `)

	privateKeyB, _ := keys.GeneratePrivateKey(keys.ECDSA_P256_SHA3_256,
		[]byte("elephant ears space cowboy octopus rodeo potato cannon pineapple"))
	publicKeyB := privateKeyB.PublicKey(keys.PublicKeyWeightThreshold)

	t.Run("ValidSignature", func(t *testing.T) {
		b, err := emulator.NewBlockchain()
		require.NoError(t, err)

		privateKeyA := b.RootKey()

		accountAddressA := b.RootAccountAddress()
		accountAddressB, err := b.CreateAccount([]flow.AccountPublicKey{publicKeyB}, codeA, getNonce())
		require.NoError(t, err)

		account, err := b.GetAccount(accountAddressB)
		require.NoError(t, err)

		assert.Equal(t, codeA, account.Code)

		tx := flow.Transaction{
			Script:             templates.UpdateAccountCode(codeB),
			ReferenceBlockHash: nil,
			Nonce:              getNonce(),
			ComputeLimit:       10,
			PayerAccount:       accountAddressA,
			ScriptAccounts:     []flow.Address{accountAddressB},
		}

		sigA, err := keys.SignTransaction(tx, privateKeyA)
		assert.NoError(t, err)

		sigB, err := keys.SignTransaction(tx, privateKeyB)
		assert.NoError(t, err)

		tx.AddSignature(accountAddressA, sigA)
		tx.AddSignature(accountAddressB, sigB)

		err = b.AddTransaction(tx)
		assert.NoError(t, err)

		result, err := b.ExecuteNextTransaction()
		assert.NoError(t, err)
		assert.True(t, result.Succeeded())

		_, err = b.CommitBlock()
		assert.NoError(t, err)

		account, err = b.GetAccount(accountAddressB)
		assert.NoError(t, err)

		assert.Equal(t, codeB, account.Code)
	})

	t.Run("InvalidSignature", func(t *testing.T) {
		b, err := emulator.NewBlockchain()
		require.NoError(t, err)

		privateKeyA := b.RootKey()

		accountAddressA := b.RootAccountAddress()
		accountAddressB, err := b.CreateAccount([]flow.AccountPublicKey{publicKeyB}, codeA, getNonce())
		require.NoError(t, err)

		account, err := b.GetAccount(accountAddressB)
		require.NoError(t, err)

		assert.Equal(t, codeA, account.Code)

		tx := flow.Transaction{
			Script:             templates.UpdateAccountCode(codeB),
			ReferenceBlockHash: nil,
			Nonce:              getNonce(),
			ComputeLimit:       10,
			PayerAccount:       accountAddressA,
			ScriptAccounts:     []flow.Address{accountAddressB},
		}

		sig, err := keys.SignTransaction(tx, privateKeyA)
		assert.NoError(t, err)

		tx.AddSignature(accountAddressA, sig)

		err = b.AddTransaction(tx)
		assert.IsType(t, &emulator.ErrMissingSignature{}, err)

		_, err = b.CommitBlock()
		assert.NoError(t, err)

		account, err = b.GetAccount(accountAddressB)
		assert.NoError(t, err)

		// code should not be updated
		assert.Equal(t, codeA, account.Code)
	})
}

func TestImportAccountCode(t *testing.T) {
	b, err := emulator.NewBlockchain()
	require.NoError(t, err)

	accountScript := []byte(`
      pub contract Computer {
          pub fun answer(): Int {
              return 42
          }
      }
	`)

	publicKey := b.RootKey().PublicKey(keys.PublicKeyWeightThreshold)

	address, err := b.CreateAccount([]flow.AccountPublicKey{publicKey}, accountScript, getNonce())
	assert.NoError(t, err)

	assert.Equal(t, flow.HexToAddress("02"), address)

	script := []byte(`
		// address imports can omit leading zeros
		import 0x02

		transaction {
		  execute {
			let answer = Computer.answer()
			if answer != 42 {
				panic("?!")
			}
		  }
		}
	`)

	tx := flow.Transaction{
		Script:             script,
		ReferenceBlockHash: nil,
		Nonce:              getNonce(),
		ComputeLimit:       10,
		PayerAccount:       b.RootAccountAddress(),
	}

	sig, err := keys.SignTransaction(tx, b.RootKey())
	assert.NoError(t, err)

	tx.AddSignature(b.RootAccountAddress(), sig)

	err = b.AddTransaction(tx)
	assert.NoError(t, err)

	result, err := b.ExecuteNextTransaction()
	assert.NoError(t, err)
	assert.True(t, result.Succeeded())

}

package examples

import (
	"context"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/dapperlabs/cadence/runtime/cmd"
	emulator "github.com/dapperlabs/flow-emulator"
	"github.com/dapperlabs/flow-go/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/client"
	"github.com/dapperlabs/flow-go-sdk/keys"
	"github.com/dapperlabs/flow-go-sdk/templates"
)

// ReadFile reads a file from the file system
func ReadFile(path string) []byte {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return contents
}

// GetNonce returns a nonce value that is guaranteed to be unique.
var GetNonce = func() func() uint64 {
	var nonce uint64
	return func() uint64 {
		nonce++
		return nonce
	}
}()

// RandomPrivateKey returns a randomly generated private key
func RandomPrivateKey() flow.AccountPrivateKey {
	seed := make([]byte, 40)
	rand.Read(seed)

	privateKey, err := keys.GeneratePrivateKey(keys.ECDSA_P256_SHA3_256, seed)
	if err != nil {
		panic(err)
	}

	return privateKey
}

// NewEmulator returns a emulator object for testing
func NewEmulator() *emulator.Blockchain {
	b, err := emulator.NewBlockchain()
	if err != nil {
		panic(err)
	}
	return b
}

func RootAccount() (flow.Address, flow.AccountPrivateKey) {
	privateKeyHex := "f87db87930770201010420c2e6c8cb9e8c9b9a7afe1df8ae431e68317ff7a9f42f8982b7877a9da76b28a7a00a06082a8648ce3d030107a14403420004c2c482bf01344a085af036f9413dd17a0d98a5b6fb4915c3ad4c3cb574e03ea5e2d47608093a26081c165722621bf9d8ff4b880cac0e7c586af3d86c0818a4af0203"
	privateKey := keys.MustDecodePrivateKeyHex(privateKeyHex)

	// root account always has address 0x01
	addr := flow.HexToAddress("01")

	return addr, privateKey
}

// SignAndSubmit signs a transaction with an array of signers and adds their signatures to the transaction
// Then submits the transaction to the emulator. If the private keys don't match up with the addresses,
// the transaction will not succeed.
// shouldRevert parameter indicates whether the transaction should fail or not
// This function asserts the correct result and commits the block if it passed
func SignAndSubmit(
	t *testing.T,
	b *emulator.Blockchain,
	tx flow.Transaction,
	signingKeys []flow.AccountPrivateKey,
	signingAddresses []flow.Address,
	shouldRevert bool,
) {
	// sign transaction with each signer
	for i := 0; i < len(signingAddresses); i++ {
		sig, err := keys.SignTransaction(tx, signingKeys[i])
		assert.NoError(t, err)

		tx.AddSignature(signingAddresses[i], sig)
	}

	// submit the signed transaction
	err := b.AddTransaction(tx)
	require.NoError(t, err)

	result, err := b.ExecuteNextTransaction()
	require.NoError(t, err)

	if shouldRevert {
		assert.True(t, result.Reverted())
	} else {
		if !assert.True(t, result.Succeeded()) {
			t.Log(result.Error.Error())
			cmd.PrettyPrintError(result.Error, "", map[string]string{"": ""})
		}
	}

	_, err = b.CommitBlock()
	assert.NoError(t, err)
}

func CreateAccount() (flow.AccountPrivateKey, flow.Address) {
	privateKey := RandomPrivateKey()

	addr := createAccount(
		[]flow.AccountPublicKey{privateKey.PublicKey(keys.PublicKeyWeightThreshold)},
		nil,
	)

	return privateKey, addr
}

func DeployContract(code []byte) flow.Address {
	return createAccount(nil, code)
}

func createAccount(publicKeys []flow.AccountPublicKey, code []byte) flow.Address {
	ctx := context.Background()
	flowClient, err := client.New("127.0.0.1:3569")
	Handle(err)

	rootAcctAddr, rootAcctKey := RootAccount()

	createAccountScript, err := templates.CreateAccount(publicKeys, code)
	Handle(err)

	createAccountTx := flow.Transaction{
		Script:       createAccountScript,
		Nonce:        GetNonce(),
		ComputeLimit: 10,
		PayerAccount: rootAcctAddr,
	}

	sig, err := keys.SignTransaction(createAccountTx, rootAcctKey)
	Handle(err)

	createAccountTx.AddSignature(rootAcctAddr, sig)

	err = flowClient.SendTransaction(ctx, createAccountTx)
	Handle(err)

	tx := WaitForSeal(ctx, flowClient, createAccountTx.Hash())

	accountCreatedEvent, err := flow.DecodeAccountCreatedEvent(tx.Events[0].Payload)
	Handle(err)

	return accountCreatedEvent.Address()
}

func Handle(err error) {
	if err != nil {
		fmt.Println("err:", err.Error())
		panic(err)
	}
}

func WaitForSeal(ctx context.Context, c *client.Client, hash crypto.Hash) *flow.Transaction {
	tx, err := c.GetTransaction(ctx, hash)
	Handle(err)

	fmt.Printf("Waiting for transaction %x to be sealed...\n", hash)

	for tx.Status != flow.TransactionSealed {
		time.Sleep(time.Second)
		fmt.Print(".")
		tx, err = c.GetTransaction(ctx, hash)
		Handle(err)
	}

	fmt.Println()
	fmt.Printf("Transaction %x sealed\n", hash)

	return tx
}

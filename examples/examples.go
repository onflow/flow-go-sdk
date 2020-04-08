package examples

import (
	"context"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/client"
	"github.com/dapperlabs/flow-go-sdk/crypto"
	"github.com/dapperlabs/flow-go-sdk/keys"
	"github.com/dapperlabs/flow-go-sdk/templates"
)

// ReadFile reads a file from the file system.
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

// RandomPrivateKey returns a randomly generated private key.
func RandomPrivateKey() crypto.PrivateKey {
	seed := make([]byte, 40)
	rand.Read(seed)

	privateKey, err := keys.GeneratePrivateKey(keys.ECDSA_P256_SHA3_256, seed)
	if err != nil {
		panic(err)
	}

	return privateKey
}

func RootAccount() (flow.Address, flow.AccountKey, crypto.PrivateKey) {
	privateKeyHex := "f87db87930770201010420c2e6c8cb9e8c9b9a7afe1df8ae431e68317ff7a9f42f8982b7877a9da76b28a7a00a06082a8648ce3d030107a14403420004c2c482bf01344a085af036f9413dd17a0d98a5b6fb4915c3ad4c3cb574e03ea5e2d47608093a26081c165722621bf9d8ff4b880cac0e7c586af3d86c0818a4af0203"
	privateKey := keys.MustDecodePrivateKeyHex(keys.ECDSA_P256_SHA3_256.SigningAlgorithm(), privateKeyHex)

	// root account always has address 0x01
	addr := flow.HexToAddress("01")

	accountKey := flow.AccountKey{
		PublicKey:      privateKey.PublicKey(),
		ID:             0,
		SignAlgo:       keys.ECDSA_P256_SHA2_256.SigningAlgorithm(),
		HashAlgo:       keys.ECDSA_P256_SHA3_256.HashingAlgorithm(),
		Weight:         keys.PublicKeyWeightThreshold,
		SequenceNumber: 0,
	}

	return addr, accountKey, privateKey
}

func CreateAccount() (flow.Address, flow.AccountKey, crypto.PrivateKey) {
	privateKey := RandomPrivateKey()

	accountKey := flow.AccountKey{
		PublicKey:      privateKey.PublicKey(),
		ID:             0,
		SignAlgo:       keys.ECDSA_P256_SHA3_256.SigningAlgorithm(),
		HashAlgo:       keys.ECDSA_P256_SHA3_256.HashingAlgorithm(),
		Weight:         keys.PublicKeyWeightThreshold,
		SequenceNumber: 0,
	}

	addr := createAccount(
		[]flow.AccountKey{accountKey},
		nil,
	)

	return addr, accountKey, privateKey
}

func DeployContract(code []byte) flow.Address {
	return createAccount(nil, code)
}

func createAccount(publicKeys []flow.AccountKey, code []byte) flow.Address {
	ctx := context.Background()
	flowClient, err := client.New("127.0.0.1:3569")
	Handle(err)

	rootAcctAddr, rootAcctKey, rootPrivateKey := RootAccount()

	createAccountScript, err := templates.CreateAccount(publicKeys, code)
	Handle(err)

	createAccountTx := flow.NewTransaction().
		SetScript(createAccountScript).
		SetProposalKey(rootAcctAddr, rootAcctKey.ID, rootAcctKey.SequenceNumber).
		SetPayer(rootAcctAddr, rootAcctKey.ID)

	rootKeySigner := crypto.NewNaiveSigner(rootPrivateKey, rootAcctKey.HashAlgo)

	err = createAccountTx.SignContainer(
		rootAcctAddr,
		rootAcctKey.ID,
		rootKeySigner,
	)
	Handle(err)

	err = flowClient.SendTransaction(ctx, *createAccountTx)
	Handle(err)

	result := WaitForSeal(ctx, flowClient, createAccountTx.ID())

	accountCreatedEvent := flow.AccountCreatedEvent(result.Events[0])
	Handle(err)

	return accountCreatedEvent.Address()
}

func Handle(err error) {
	if err != nil {
		fmt.Println("err:", err.Error())
		panic(err)
	}
}

func WaitForSeal(ctx context.Context, c *client.Client, id flow.Identifier) *flow.TransactionResult {
	result, err := c.GetTransactionResult(ctx, id)
	Handle(err)

	fmt.Printf("Waiting for transaction %x to be sealed...\n", id)

	for result.Status != flow.TransactionSealed {
		time.Sleep(time.Second)
		fmt.Print(".")
		result, err = c.GetTransactionResult(ctx, id)
		Handle(err)
	}

	fmt.Println()
	fmt.Printf("Transaction %x sealed\n", id)

	return result
}

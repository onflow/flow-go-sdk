package examples

import (
	"context"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"time"

	"google.golang.org/grpc"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/templates"
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

	privateKey, err := crypto.GeneratePrivateKey(crypto.ECDSA_P256, seed)
	if err != nil {
		panic(err)
	}

	return privateKey
}

func RootAccount(flowClient *client.Client) (flow.Address, flow.AccountKey, crypto.PrivateKey) {
	privateKeyHex := "f87db87930770201010420c2e6c8cb9e8c9b9a7afe1df8ae431e68317ff7a9f42f8982b7877a9da76b28a7a00a06082a8648ce3d030107a14403420004c2c482bf01344a085af036f9413dd17a0d98a5b6fb4915c3ad4c3cb574e03ea5e2d47608093a26081c165722621bf9d8ff4b880cac0e7c586af3d86c0818a4af0203"
	privateKey, sigAlgo, hashAlgo := crypto.MustDecodeWrappedPrivateKeyHex(privateKeyHex)

	// root account always has address 0x01
	addr := flow.HexToAddress("01")

	acc, err := flowClient.GetAccount(context.Background(), addr)
	Handle(err)

	accountKey := flow.AccountKey{
		PublicKey:      privateKey.PublicKey(),
		ID:             0,
		SignAlgo:       sigAlgo,
		HashAlgo:       hashAlgo,
		Weight:         flow.AccountKeyWeightThreshold,
		SequenceNumber: acc.Keys[0].SequenceNumber,
	}

	return addr, accountKey, privateKey
}

func CreateAccount() (flow.Address, flow.AccountKey, crypto.PrivateKey) {
	privateKey := RandomPrivateKey()

	accountKey := flow.AccountKey{
		PublicKey: privateKey.PublicKey(),
		SignAlgo:  crypto.ECDSA_P256,
		HashAlgo:  crypto.SHA3_256,
		Weight:    flow.AccountKeyWeightThreshold,
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
	flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
	Handle(err)

	rootAcctAddr, rootAcctKey, rootPrivateKey := RootAccount(flowClient)

	createAccountScript, err := templates.CreateAccount(publicKeys, code)
	Handle(err)

	createAccountTx := flow.NewTransaction().
		SetScript(createAccountScript).
		SetProposalKey(rootAcctAddr, rootAcctKey.ID, rootAcctKey.SequenceNumber).
		SetPayer(rootAcctAddr)

	err = createAccountTx.SignEnvelope(
		rootAcctAddr,
		rootAcctKey.ID,
		crypto.NewNaiveSigner(rootPrivateKey, rootAcctKey.HashAlgo),
	)
	Handle(err)

	err = flowClient.SendTransaction(ctx, *createAccountTx)
	Handle(err)

	result := WaitForSeal(ctx, flowClient, createAccountTx.ID())
	Handle(result.Error)

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

	for result.Status != flow.TransactionStatusSealed {
		time.Sleep(time.Second)
		fmt.Print(".")
		result, err = c.GetTransactionResult(ctx, id)
		Handle(err)
	}

	fmt.Println()
	fmt.Printf("Transaction %x sealed\n", id)

	return result
}

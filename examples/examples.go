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
	seed := make([]byte, crypto.MinSeedLengthECDSA_P256)
	rand.Read(seed)

	privateKey, err := crypto.GeneratePrivateKey(crypto.ECDSA_P256, seed)
	if err != nil {
		panic(err)
	}

	return privateKey
}

const defaultRootKeySeed = "elephant ears space cowboy octopus rodeo potato cannon pineapple"

func RootAccount(flowClient *client.Client) (flow.Address, *flow.AccountKey, crypto.Signer) {
	privateKey, _ := crypto.GeneratePrivateKey(crypto.ECDSA_P256, []byte(defaultRootKeySeed))

	// root account always has address 0x01
	addr := flow.HexToAddress("01")

	acc, err := flowClient.GetAccount(context.Background(), addr)
	Handle(err)

	accountKey := acc.Keys[0]

	signer := crypto.NewInMemorySigner(privateKey, accountKey.HashAlgo)

	return addr, accountKey, signer
}

func CreateAccount() (flow.Address, *flow.AccountKey, crypto.Signer) {
	privateKey := RandomPrivateKey()

	accountKey := flow.NewAccountKey().
		FromPrivateKey(privateKey).
		SetHashAlgo(crypto.SHA3_256).
		SetWeight(flow.AccountKeyWeightThreshold)

	addr := createAccount(
		[]*flow.AccountKey{accountKey},
		nil,
	)

	signer := crypto.NewInMemorySigner(privateKey, accountKey.HashAlgo)

	return addr, accountKey, signer
}

func DeployContract(code []byte) flow.Address {
	return createAccount(nil, code)
}

func createAccount(publicKeys []*flow.AccountKey, code []byte) flow.Address {
	ctx := context.Background()
	flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
	Handle(err)

	rootAcctAddr, rootAcctKey, rootSigner := RootAccount(flowClient)

	createAccountScript, err := templates.CreateAccount(publicKeys, code)
	Handle(err)

	createAccountTx := flow.NewTransaction().
		SetScript(createAccountScript).
		SetProposalKey(rootAcctAddr, rootAcctKey.ID, rootAcctKey.SequenceNumber).
		SetPayer(rootAcctAddr)

	err = createAccountTx.SignEnvelope(rootAcctAddr, rootAcctKey.ID, rootSigner)
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

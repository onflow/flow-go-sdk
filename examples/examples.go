/*
 * Flow Go SDK
 *
 * Copyright Flow Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package examples

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-cli/flowkit/config"
	"github.com/spf13/afero"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access"
	"github.com/onflow/flow-go-sdk/access/grpc"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/templates"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-cli/flowkit/config/json"

	"github.com/onflow/cadence/sema"
)

const configPath = "./flow.json"

var Config *config.Config

func First[T any](arr []T, f func(e T) bool) (T, error) {
	for _, e := range arr {
		if f(e) {
			return e, nil
		}
	}
	var nullT T
	return nullT, fmt.Errorf("not found")
}

func init() {
	var mockFS = afero.NewOsFs()

	var af = afero.Afero{Fs: mockFS}

	l := config.NewLoader(af)

	l.AddConfigParser(json.NewParser())

	var err error
	Config, err = l.Load([]string{configPath})
	if err != nil {
		log.Fatal(err)
	}
}

func ServiceAccount(flowClient access.Client) (flow.Address, *flow.AccountKey, crypto.Signer) {
	acc := Config.Accounts[0]
	privateKey, err := crypto.DecodePrivateKeyHex(crypto.ECDSA_P256, strings.Replace(acc.Key.PrivateKey.String(), "0x", "", 1))
	if err != nil {
		log.Fatalf("failed to decode private key %s", err)
	}

	account, err := flowClient.GetAccount(context.Background(), acc.Address)
	if err != nil {
		log.Fatalf("failed to get account %s", err)
	}
	accountKey := account.Keys[0]
	signer, err := crypto.NewInMemorySigner(privateKey, accountKey.HashAlgo)
	if err != nil {
		log.Fatalf("failed to create in mem signer: %s", err)
	}
	return flow.HexToAddress(acc.Address.Hex()), accountKey, signer
}

// RandomPrivateKey returns a randomly generated ECDSA P-256 private key.
func RandomPrivateKey() crypto.PrivateKey {
	seed := make([]byte, crypto.MinSeedLength)
	_, err := rand.Read(seed)
	Handle(err)

	privateKey, err := crypto.GeneratePrivateKey(crypto.ECDSA_P256, seed)
	Handle(err)

	return privateKey
}

func RandomTransaction(flowClient access.Client) *flow.Transaction {
	serviceAcctAddr, serviceAcctKey, serviceSigner := ServiceAccount(flowClient)

	tx := flow.NewTransaction().
		SetPayer(serviceAcctAddr).
		SetProposalKey(serviceAcctAddr, serviceAcctKey.Index, serviceAcctKey.SequenceNumber).
		SetScript([]byte("transaction { prepare(auth: AuthAccount) {} }")).
		AddAuthorizer(serviceAcctAddr).
		SetReferenceBlockID(GetReferenceBlockId(flowClient))

	err := tx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, serviceSigner)
	Handle(err)

	err = flowClient.SendTransaction(context.Background(), *tx)
	Handle(err)

	return tx
}

func RandomAccount(flowClient access.Client) (flow.Address, *flow.AccountKey, crypto.Signer) {
	privateKey := RandomPrivateKey()

	accountKey := flow.NewAccountKey().
		FromPrivateKey(privateKey).
		SetHashAlgo(crypto.SHA3_256).
		SetWeight(flow.AccountKeyWeightThreshold)

	account := CreateAccount(flowClient, []*flow.AccountKey{accountKey})
	FundAccountInEmulator(flowClient, account.Address, 10.0)
	signer, err := crypto.NewInMemorySigner(privateKey, accountKey.HashAlgo)
	Handle(err)
	return account.Address, account.Keys[0], signer
}

func GetReferenceBlockId(flowClient access.Client) flow.Identifier {
	block, err := flowClient.GetLatestBlock(context.Background(), true)
	Handle(err)

	return block.ID
}

func CreateAccountWithContracts(flowClient access.Client, publicKeys []*flow.AccountKey, contracts []templates.Contract) *flow.Account {
	serviceAcctAddr, serviceAcctKey, serviceSigner := ServiceAccount(flowClient)

	referenceBlockID := GetReferenceBlockId(flowClient)

	createAccountTx, err := templates.CreateAccount(publicKeys, contracts, serviceAcctAddr)
	Handle(err)
	createAccountTx.
		SetProposalKey(serviceAcctAddr, serviceAcctKey.Index, serviceAcctKey.SequenceNumber).
		SetReferenceBlockID(referenceBlockID).
		SetPayer(serviceAcctAddr)

	err = createAccountTx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, serviceSigner)
	Handle(err)

	ctx := context.Background()
	err = flowClient.SendTransaction(ctx, *createAccountTx)
	Handle(err)

	result := WaitForSeal(ctx, flowClient, createAccountTx.ID())
	Handle(result.Error)

	for _, event := range result.Events {
		if event.Type != flow.EventAccountCreated {
			continue
		}
		accountCreatedEvent := flow.AccountCreatedEvent(event)

		addr := accountCreatedEvent.Address()
		account, err := flowClient.GetAccount(ctx, addr)
		Handle(err)

		return account
	}
	panic("could not find an AccountCreatedEvent")
}

func ReadFile(name string) string {
	body, err := os.ReadFile(name)
	if err != nil {
		log.Fatalf("unable to read file")
	}
	return string(body)
}

/**
 * mintTokensToAccountTemplate transaction mints tokens by using the service account (in the emulator)
 * and deposits them to the recipient.
 */
var mintTokensToAccountTemplate = `
import FungibleToken from 0x%s
import FlowToken from 0x%s

transaction(recipient: Address, amount: UFix64) {
	let tokenAdmin: &FlowToken.Administrator
	let tokenReceiver: &{FungibleToken.Receiver}

	prepare(signer: auth(BorrowValue) &Account) {
		self.tokenAdmin = signer.storage
			.borrow<&FlowToken.Administrator>(from: /storage/flowTokenAdmin)
			?? panic("Signer is not the token admin")

		self.tokenReceiver = getAccount(recipient)
			.capabilities.borrow<&{FungibleToken.Receiver}>(/public/flowTokenReceiver)
			?? panic("Unable to borrow receiver reference")
	}

	execute {
		let minter <- self.tokenAdmin.createNewMinter(allowedAmount: amount)
		let mintedVault <- minter.mintTokens(amount: amount)

		self.tokenReceiver.deposit(from: <-mintedVault)

		destroy minter
	}
}
`

// FundAccountInEmulator Mints FLOW to an account. Minting only works in an emulator environment.
func FundAccountInEmulator(flowClient access.Client, address flow.Address, amount float64) {
	serviceAcctAddr, serviceAcctKey, serviceSigner := ServiceAccount(flowClient)

	referenceBlockID := GetReferenceBlockId(flowClient)

	fungibleToken, err := First[config.Contract](Config.Contracts, func(e config.Contract) bool {
		return e.Name == "FungibleToken"
	})

	Handle(err)

	flowToken, err := First[config.Contract](Config.Contracts, func(e config.Contract) bool {
		return e.Name == "FlowToken"
	})
	Handle(err)

	fungibleTokenAddress := fungibleToken.Aliases[0].Address.Hex()
	flowTokenAddress := flowToken.Aliases[0].Address.Hex()

	recipient := cadence.NewAddress(address)
	uintAmount := uint64(amount * sema.Fix64Factor)
	cadenceAmount := cadence.UFix64(uintAmount)

	fundAccountTx :=
		flow.NewTransaction().
			SetScript([]byte(fmt.Sprintf(mintTokensToAccountTemplate, fungibleTokenAddress, flowTokenAddress))).
			AddAuthorizer(serviceAcctAddr).
			AddRawArgument(jsoncdc.MustEncode(recipient)).
			AddRawArgument(jsoncdc.MustEncode(cadenceAmount)).
			SetProposalKey(serviceAcctAddr, serviceAcctKey.Index, serviceAcctKey.SequenceNumber).
			SetReferenceBlockID(referenceBlockID).
			SetPayer(serviceAcctAddr)

	err = fundAccountTx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, serviceSigner)
	Handle(err)

	ctx := context.Background()
	err = flowClient.SendTransaction(ctx, *fundAccountTx)
	Handle(err)

	result := WaitForSeal(ctx, flowClient, fundAccountTx.ID())
	Handle(result.Error)
}

func CreateAccount(flowClient access.Client, publicKeys []*flow.AccountKey) *flow.Account {
	return CreateAccountWithContracts(flowClient, publicKeys, nil)
}

func Handle(err error) {
	if err != nil {
		fmt.Println("err:", err.Error())
		panic(err)
	}
}

func NewFlowGRPCClient() *grpc.Client {
	c, err := grpc.NewClient(grpc.EmulatorHost)
	Handle(err)
	return c
}

func WaitForSeal(ctx context.Context, c access.Client, id flow.Identifier) *flow.TransactionResult {
	result, err := c.GetTransactionResult(ctx, id)
	Handle(err)

	fmt.Printf("Waiting for transaction %s to be sealed...\n", id)

	for result.Status != flow.TransactionStatusSealed {
		time.Sleep(time.Second)
		fmt.Print(".")
		result, err = c.GetTransactionResult(ctx, id)
		Handle(err)
	}

	fmt.Println()
	fmt.Printf("Transaction %s sealed\n", id)
	return result
}

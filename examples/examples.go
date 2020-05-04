/*
 * Flow Go SDK
 *
 * Copyright 2019-2020 Dapper Labs, Inc.
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
	"io/ioutil"
	"time"

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

// RandomPrivateKey returns a randomly generated ECDSA P-256 private key.
func RandomPrivateKey() crypto.PrivateKey {
	seed := make([]byte, crypto.MinSeedLengthECDSA_P256)
	rand.Read(seed)

	privateKey, err := crypto.GeneratePrivateKey(crypto.ECDSA_P256, seed)
	if err != nil {
		panic(err)
	}

	return privateKey
}

func RandomAccount(flowClient *client.Client) (flow.Address, *flow.AccountKey, crypto.Signer) {
	privateKey := RandomPrivateKey()

	accountKey := flow.NewAccountKey().
		FromPrivateKey(privateKey).
		SetHashAlgo(crypto.SHA3_256).
		SetWeight(flow.AccountKeyWeightThreshold)

	account := CreateAccount(
		flowClient,
		[]*flow.AccountKey{accountKey},
		nil,
	)

	signer := crypto.NewInMemorySigner(privateKey, accountKey.HashAlgo)

	return account.Address, account.Keys[0], signer
}

func DeployContract(flowClient *client.Client, code []byte) flow.Address {
	account := CreateAccount(flowClient, nil, code)
	return account.Address
}

func CreateAccount(flowClient *client.Client, publicKeys []*flow.AccountKey, code []byte) *flow.Account {
	ctx := context.Background()

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

	addr := accountCreatedEvent.Address()

	account, err := flowClient.GetAccount(ctx, addr)
	Handle(err)

	return account
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

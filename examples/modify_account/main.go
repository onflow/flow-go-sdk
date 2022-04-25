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

package main

import (
	"context"
	"fmt"

	"github.com/onflow/flow-go-sdk/templates"

	"google.golang.org/grpc"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/examples"
)

func main() {
	ModifyAccountDemo()
}

func prepareAndSendTx(ctx context.Context, client *client.Client, key *flow.AccountKey, tx flow.Transaction) {
	serviceAcctAddr, serviceAcctKey, serviceSigner := examples.ServiceAccount(client)
	tx.SetProposalKey(
		serviceAcctAddr,
		serviceAcctKey.Index,
		serviceAcctKey.SequenceNumber,
	)

	referenceBlockID := examples.GetReferenceBlockId(client)
	tx.SetReferenceBlockID(referenceBlockID)
	tx.SetPayer(serviceAcctAddr)

	err := tx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, serviceSigner)
	examples.Handle(err)

	err = client.SendTransaction(ctx, tx)
	examples.Handle(err)
}

func keysString(keys []*flow.AccountKey) string {
	k := ""
	for _, key := range keys {
		k = fmt.Sprintf("%s\n%s", k, key.PublicKey.String())
	}
	return k
}

func contractsString(contracts map[string][]byte) string {
	k := ""
	for name, _ := range contracts {
		if k == "" {
			k = name
			continue
		}
		k = fmt.Sprintf("%s, %s", k, name)
	}
	return k
}

func ModifyAccountDemo() {
	ctx := context.Background()
	flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
	examples.Handle(err)

	serviceAcctAddr, _, _ := examples.ServiceAccount(flowClient)

	myPrivateKey := examples.RandomPrivateKey()
	myAcctKey := flow.NewAccountKey().
		FromPrivateKey(myPrivateKey).
		SetHashAlgo(crypto.SHA3_256).
		SetWeight(flow.AccountKeyWeightThreshold)

	// create a new account without any contracts
	createAccountTx := templates.CreateAccount([]*flow.AccountKey{myAcctKey}, nil, serviceAcctAddr)
	prepareAndSendTx(ctx, flowClient, myAcctKey, *createAccountTx)

	acc, err := flowClient.GetAccount(ctx, serviceAcctAddr)
	examples.Handle(err)

	fmt.Printf("\nInitial account contracts: %s and keys: %s", contractsString(acc.Contracts), keysString(acc.Keys))

	addContractTx := templates.AddAccountContract(
		serviceAcctAddr,
		templates.Contract{
			Name:   "Foo",
			Source: "pub contract Foo {}",
		},
	)
	prepareAndSendTx(ctx, flowClient, myAcctKey, *addContractTx)

	acc, _ = flowClient.GetAccount(ctx, serviceAcctAddr)

	fmt.Printf("\nDeployed contracts on the account after 'Foo' deployment: %s", contractsString(acc.Contracts))

	updateTx := templates.UpdateAccountContract(
		serviceAcctAddr,
		templates.Contract{
			Name:   "Foo",
			Source: "pub contract Foo { pub fun hello() {} }",
		},
	)
	prepareAndSendTx(ctx, flowClient, myAcctKey, *updateTx)

	acc, _ = flowClient.GetAccount(ctx, serviceAcctAddr)
	fmt.Printf("\nContract 'Foo' after update: %s", acc.Contracts["Foo"])

	removeContractTx := templates.RemoveAccountContract(serviceAcctAddr, "Foo")
	prepareAndSendTx(ctx, flowClient, myAcctKey, *removeContractTx)

	acc, _ = flowClient.GetAccount(ctx, serviceAcctAddr)
	fmt.Printf("\nDeployed contracts on the account after 'Foo' removal: %s", contractsString((acc.Contracts)))

	newPrivKey := examples.RandomPrivateKey()
	newAcctKey := flow.NewAccountKey().
		FromPrivateKey(newPrivKey).
		SetHashAlgo(crypto.SHA3_256).
		SetWeight(flow.AccountKeyWeightThreshold)

	addKeyTx := templates.AddAccountKey(serviceAcctAddr, newAcctKey)
	prepareAndSendTx(ctx, flowClient, myAcctKey, *addKeyTx)

	acc, _ = flowClient.GetAccount(ctx, serviceAcctAddr)
	fmt.Printf("\nAccount keys after adding new key: %s", keysString(acc.Keys))

	removeKeyTx := templates.RemoveAccountKey(serviceAcctAddr, 1)
	prepareAndSendTx(ctx, flowClient, myAcctKey, *removeKeyTx)

	acc, _ = flowClient.GetAccount(ctx, serviceAcctAddr)
	fmt.Printf("\nAccount keys after removing the last key: %s", keysString(acc.Keys))
}

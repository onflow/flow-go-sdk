/*
 * Flow Go SDK
 *
 * Copyright 2019-2021 Dapper Labs, Inc.
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
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/examples"
	"github.com/onflow/flow-go-sdk/templates"
	"google.golang.org/grpc"
)

func main() {
	StorageUsageDemo()
}

func StorageUsageDemo() {
	ctx := context.Background()
	flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
	examples.Handle(err)

	serviceAcctAddr, serviceAcctKey, serviceSigner := examples.ServiceAccount(flowClient)

	// A contract that defines a resource with a string so its easier to demonstrate adding resources of different sizes
	contract := `
		pub contract StorageDemo {
			pub resource StorageTestResource {
				pub let data: String
				init(data: String) {
					self.data = data
				}
			}
			pub fun createStorageTestResource(_ data: String): @StorageTestResource {
				return <- create StorageTestResource(data: data)
			}
		}
	`
	privateKey := examples.RandomPrivateKey()

	key := flow.NewAccountKey().
		SetPublicKey(privateKey.PublicKey()).
		SetSigAlgo(privateKey.Algorithm()).
		SetHashAlgo(crypto.SHA3_256).
		SetWeight(flow.AccountKeyWeightThreshold)

	keySigner := crypto.NewInMemorySigner(privateKey, key.HashAlgo)

	demoAccount := examples.CreateAccountWithContracts(flowClient,
		[]*flow.AccountKey{key}, []templates.Contract{{
			Name:   "StorageDemo",
			Source: contract,
		}})
	serviceAcctKey.SequenceNumber++

	// try to save a very large resource to the demoAccount
	txId := sendSaveLargeResourceTransaction(
		ctx,
		flowClient,
		serviceAcctAddr,
		serviceAcctKey,
		serviceSigner,
		demoAccount,
		keySigner,
	)

	result := examples.WaitForSeal(ctx, flowClient, txId)

	if result.Error == nil {
		fmt.Println("Storage limits are off")
		return
	}

	fmt.Println("Storage limit reached")
}

func sendSaveLargeResourceTransaction(
	ctx context.Context,
	flowClient *client.Client,
	serviceAcctAddr flow.Address,
	serviceAcctKey *flow.AccountKey,
	serviceSigner crypto.Signer,
	demoAccount *flow.Account,
	demoSigner crypto.InMemorySigner,
) flow.Identifier {
	// string bigger than 100kb
	longString := longString()

	// Send a tx that emits the event in the deployed contract
	script := fmt.Sprintf(`
		import StorageDemo from 0x%s

		transaction {
			prepare(acct: AuthAccount) {
				let storageUsed = acct.storageUsed
				
				// create resource and save it on the account 
				let bigResource <- StorageDemo.createStorageTestResource("%s")
				acct.save(<-bigResource, to: /storage/StorageDemo)

				let storageUsedAfter = acct.storageUsed

				if (storageUsed == storageUsedAfter) {
					panic("storage used will change")
				}
				
				if (storageUsedAfter > acct.storageCapacity) {
					// this is where we could deposit more flow to acct to increase its storaga capacity if we wanted to
					log("Storage used is over capacity. This transaction will fail if storage limits are on on this chain.")
				}
			}
		}
	`, demoAccount.Address.Hex(), longString)

	referenceBlockID := examples.GetReferenceBlockId(flowClient)
	runScriptTx := flow.NewTransaction().
		SetScript([]byte(script)).
		SetPayer(serviceAcctAddr).
		AddAuthorizer(demoAccount.Address).
		SetReferenceBlockID(referenceBlockID).
		SetProposalKey(serviceAcctAddr, serviceAcctKey.Index, serviceAcctKey.SequenceNumber)

	err := runScriptTx.SignPayload(demoAccount.Address, demoAccount.Keys[0].Index, demoSigner)
	examples.Handle(err)

	err = runScriptTx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, serviceSigner)
	examples.Handle(err)

	err = flowClient.SendTransaction(ctx, *runScriptTx)
	examples.Handle(err)

	serviceAcctKey.SequenceNumber++

	return runScriptTx.ID()
}

func longString() string {
	// 100k bytes
	b := make([]byte, 100000)
	_, err := rand.Read(b)
	examples.Handle(err)
	longString := base64.StdEncoding.EncodeToString(b) // after encoding this is ~ 130k bytes
	return longString
}

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

	"google.golang.org/grpc"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/examples"
)

func main() {
	MultiPartySingleSignatureDemo()
}

func MultiPartySingleSignatureDemo() {
	ctx := context.Background()

	flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
	examples.Handle(err)

	privateKey1 := examples.RandomPrivateKey()
	privateKey3 := examples.RandomPrivateKey()

	key1 := flow.NewAccountKey().
		SetPublicKey(privateKey1.PublicKey()).
		SetSigAlgo(privateKey1.Algorithm()).
		SetHashAlgo(crypto.SHA3_256).
		SetWeight(flow.AccountKeyWeightThreshold)

	key1Signer := crypto.NewInMemorySigner(privateKey1, key1.HashAlgo)

	key3 := flow.NewAccountKey().
		SetPublicKey(privateKey3.PublicKey()).
		SetSigAlgo(privateKey3.Algorithm()).
		SetHashAlgo(crypto.SHA3_256).
		SetWeight(flow.AccountKeyWeightThreshold)

	key3Signer := crypto.NewInMemorySigner(privateKey3, key3.HashAlgo)

	account1 := examples.CreateAccount(flowClient, []*flow.AccountKey{key1})
	account2 := examples.CreateAccount(flowClient, []*flow.AccountKey{key3})
	referenceBlockID := examples.GetReferenceBlockId(flowClient)

	tx := flow.NewTransaction().
		SetScript([]byte(`
            transaction { 
                prepare(signer: AuthAccount) { log(signer.address) }
            }
        `)).
		SetGasLimit(100).
		SetProposalKey(account1.Address, account1.Keys[0].Index, account1.Keys[0].SequenceNumber).
		SetReferenceBlockID(referenceBlockID).
		SetPayer(account2.Address).
		AddAuthorizer(account1.Address)

	// account 1 signs the payload with key 1
	err = tx.SignPayload(account1.Address, account1.Keys[0].Index, key1Signer)
	examples.Handle(err)

	// account 2 signs the envelope with key 3
	err = tx.SignEnvelope(account2.Address, account2.Keys[0].Index, key3Signer)
	examples.Handle(err)

	err = flowClient.SendTransaction(ctx, *tx)
	examples.Handle(err)

	examples.WaitForSeal(ctx, flowClient, tx.ID())

	fmt.Println("Transaction complete!")
}

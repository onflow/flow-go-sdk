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

package main

import (
	"context"
	"fmt"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access/http"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/examples"
)

func main() {
	SinglePartyMultiSignatureDemo()
}

func SinglePartyMultiSignatureDemo() {
	ctx := context.Background()
	flowClient, err := http.NewClient(http.EmulatorHost)
	examples.Handle(err)

	privateKey1 := examples.RandomPrivateKey()
	privateKey2 := examples.RandomPrivateKey()

	key1 := flow.NewAccountKey().
		SetPublicKey(privateKey1.PublicKey()).
		SetSigAlgo(privateKey1.Algorithm()).
		SetHashAlgo(crypto.SHA3_256).
		SetWeight(flow.AccountKeyWeightThreshold / 2)

	key1Signer, err := crypto.NewInMemorySigner(privateKey1, key1.HashAlgo)
	examples.Handle(err)

	key2 := flow.NewAccountKey().
		SetPublicKey(privateKey2.PublicKey()).
		SetSigAlgo(privateKey2.Algorithm()).
		SetHashAlgo(crypto.SHA3_256).
		SetWeight(flow.AccountKeyWeightThreshold / 2)

	key2Signer, err := crypto.NewInMemorySigner(privateKey2, key2.HashAlgo)
	examples.Handle(err)

	account1 := examples.CreateAccount(flowClient, []*flow.AccountKey{key1, key2})
	// Add some flow for the transaction fees
	examples.FundAccountInEmulator(flowClient, account1.Address, 1.0)

	referenceBlockID := examples.GetReferenceBlockId(flowClient)
	tx := flow.NewTransaction().
		SetScript([]byte(`
            transaction {
                prepare(signer: auth(Storage) &Account) { log(signer.address) }
            }
        `)).
		SetComputeLimit(100).
		SetProposalKey(account1.Address, account1.Keys[0].Index, account1.Keys[0].SequenceNumber).
		SetReferenceBlockID(referenceBlockID).
		SetPayer(account1.Address).
		AddAuthorizer(account1.Address)

	// account 1 signs the envelope with key 1
	err = tx.SignEnvelope(account1.Address, account1.Keys[0].Index, key1Signer)
	examples.Handle(err)

	// account 1 signs the envelope with key 2
	err = tx.SignEnvelope(account1.Address, account1.Keys[1].Index, key2Signer)
	examples.Handle(err)

	err = flowClient.SendTransaction(ctx, *tx)
	examples.Handle(err)

	examples.WaitForSeal(ctx, flowClient, tx.ID())

	fmt.Println("Transaction complete!")
}

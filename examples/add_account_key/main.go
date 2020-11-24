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
	"github.com/onflow/flow-go-sdk/templates"
)

func main() {
	AddAccountKeyDemo()
}

func AddAccountKeyDemo() {
	ctx := context.Background()

	flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
	examples.Handle(err)

	acctAddr, acctKey, acctSigner := examples.RandomAccount(flowClient)

	// Create the new key to add to your account
	myPrivateKey := examples.RandomPrivateKey()
	myAcctKey := flow.NewAccountKey().
		FromPrivateKey(myPrivateKey).
		SetHashAlgo(crypto.SHA3_256).
		SetWeight(flow.AccountKeyWeightThreshold)

	addKeyTx := templates.AddAccountKey(acctAddr, myAcctKey)

	referenceBlockID := examples.GetReferenceBlockId(flowClient)

	addKeyTx.SetProposalKey(acctAddr, acctKey.Index, acctKey.SequenceNumber)
	addKeyTx.SetReferenceBlockID(referenceBlockID)
	addKeyTx.SetPayer(acctAddr)

	// Sign the transaction with the new account.
	err = addKeyTx.SignEnvelope(acctAddr, acctKey.Index, acctSigner)
	examples.Handle(err)

	// Send the transaction to the network.
	err = flowClient.SendTransaction(ctx, *addKeyTx)
	examples.Handle(err)

	examples.WaitForSeal(ctx, flowClient, addKeyTx.ID())

	fmt.Println("Public key added to account!")
}

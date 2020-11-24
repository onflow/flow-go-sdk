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

	"github.com/onflow/cadence"
	"google.golang.org/grpc"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/crypto/cloudkms"
	"github.com/onflow/flow-go-sdk/examples"
	"github.com/onflow/flow-go-sdk/test"
)

func main() {
	GoogleCloudKMSDemo()
}

func GoogleCloudKMSDemo() {
	ctx := context.Background()

	flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
	examples.Handle(err)

	accountAddress := test.AddressGenerator().New()
	accountKeyID := 0

	accountKMSKey := cloudkms.Key{
		ProjectID:  "my-project",
		LocationID: "global",
		KeyRingID:  "flow",
		KeyID:      "my-account",
		KeyVersion: "1",
	}

	kmsClient, err := cloudkms.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	accountKMSSigner, err := kmsClient.SignerForKey(
		ctx,
		accountAddress,
		accountKMSKey,
	)
	if err != nil {
		panic(err)
	}

	serviceAccount, err := flowClient.GetAccount(ctx, accountAddress)
	if err != nil {
		panic(err)
	}

	latestBlock, err := flowClient.GetLatestBlockHeader(ctx, true)
	if err != nil {
		panic(err)
	}

	accountKey := serviceAccount.Keys[accountKeyID]

	tx := flow.NewTransaction().
		SetScript(test.GreetingScript).
		SetReferenceBlockID(latestBlock.ID).
		SetProposalKey(accountAddress, accountKey.Index, accountKey.SequenceNumber).
		SetPayer(accountAddress)

	err = tx.AddArgument(cadence.NewString(test.GreetingGenerator().Random()))
	examples.Handle(err)

	err = tx.SignEnvelope(accountAddress, accountKey.Index, accountKMSSigner)
	examples.Handle(err)

	err = flowClient.SendTransaction(ctx, *tx)
	examples.Handle(err)

	examples.WaitForSeal(ctx, flowClient, tx.ID())

	fmt.Println("Transaction complete!")
}

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
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/examples"
	"google.golang.org/grpc"
)

func main() {
	prepareDemo()
	demo()
}

func demo(exampleCollectionID flow.Identifier) {
	ctx := context.Background()
	flowClient := examples.NewFlowClient()

	// get collection by ID
	collection, err := flowClient.GetCollection(ctx, exampleCollectionID)
	printCollection(collection, err)
}

func printCollection(collection *flow.Collection, err error) {
	examples.Handle(err)

}

func prepareDemo() {
	flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
	examples.Handle(err)
	defer func() {
		err := flowClient.Close()
		if err != nil {
			panic(err)
		}
	}()

	referenceBlockID := examples.GetReferenceBlockId(flowClient)
	acctAddr, acctKey, acctSigner := examples.RandomAccount(flowClient)
	tx := flow.NewTransaction().
		SetPayer(acctAddr).
		SetReferenceBlockID(referenceBlockID).
		SetProposalKey(acctAddr, acctKey.Index, acctKey.SequenceNumber)

	err = tx.SignEnvelope(acctAddr, 0, acctSigner)
	examples.Handle(err)

	err = flowClient.SendTransaction(context.Background(), *tx)
	examples.Handle(err)

}

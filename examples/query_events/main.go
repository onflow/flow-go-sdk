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
	"github.com/onflow/flow-go-sdk/examples"
)

func main() {
	QueryEventsDemo()
}

func QueryEventsDemo() {
	ctx := context.Background()

	//flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
	flowClient, err := client.New("access.mainnet.nodes.onflow.org:9000", grpc.WithInsecure())
	examples.Handle(err)

	//acctAddr, acctKey, acctSigner := examples.RandomAccount(flowClient)
	//
	//// Deploy a contract with an event defined
	//contract := `
	//	pub contract EventDemo {
	//		pub event Add(x: Int, y: Int, sum: Int)
	//
	//		pub fun add(x: Int, y: Int) {
	//			let sum = x + y
	//			emit Add(x: x, y: y, sum: sum)
	//		}
	//	}
	//`
	//
	//contractAccount := examples.CreateAccountWithContracts(flowClient,
	//	nil, []templates.Contract{{
	//		Name:   "EventDemo",
	//		Source: contract,
	//	}})
	//
	//// Send a tx that emits the event in the deployed contract
	//script := fmt.Sprintf(`
	//	import EventDemo from 0x%s
	//
	//	transaction {
	//		execute {
	//			EventDemo.add(x: 2, y: 3)
	//		}
	//	}
	//`, contractAccount.Address.Hex())
	//
	//referenceBlockID := examples.GetReferenceBlockId(flowClient)
	//runScriptTx := flow.NewTransaction().
	//	SetScript([]byte(script)).
	//	SetPayer(acctAddr).
	//	SetReferenceBlockID(referenceBlockID).
	//	SetProposalKey(acctAddr, acctKey.Index, acctKey.SequenceNumber)
	//
	//err = runScriptTx.SignEnvelope(acctAddr, acctKey.Index, acctSigner)
	//examples.Handle(err)
	//
	//err = flowClient.SendTransaction(ctx, *runScriptTx)
	//examples.Handle(err)
	//
	//examples.WaitForSeal(ctx, flowClient, runScriptTx.ID())

	block, err := flowClient.GetBlockByHeight(ctx, 16269420)
	examples.Handle(err)

	txIDs := make([]flow.Identifier, 0)
	events := make([]flow.Event, 0)

	for _, colGuarantee := range block.CollectionGuarantees {
		collection, err := flowClient.GetCollection(ctx, colGuarantee.CollectionID)
		examples.Handle(err)

		for _, txID := range collection.TransactionIDs {
			transactionResult, err := flowClient.GetTransactionResult(ctx, txID)
			examples.Handle(err)

			for _, event := range transactionResult.Events {
				events = append(events, event)
			}

			txIDs = append(txIDs, txID)
		}

		fmt.Printf("Block height %d, ID: %s\n", block.Height, block.ID)

		for i, event := range events {
			fmt.Printf("Found event #%d in block #%d\n", i+1, block.Height)
			fmt.Printf("Transaction ID: %s\n", event.TransactionID)
			fmt.Printf("Event ID: %s\n", event.ID())
			fmt.Println(event.String())
		}
	}

}

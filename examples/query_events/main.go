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
	"github.com/onflow/flow-go-sdk/examples"
)

func main() {
	QueryEventsDemo()
}

func QueryEventsDemo() {
	ctx := context.Background()

	flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
	examples.Handle(err)

	acctAddr, acctKey, acctSigner := examples.RandomAccount(flowClient)

	// Deploy a contract with an event defined
	contract := `
		pub contract EventDemo {
			pub event Add(x: Int, y: Int, sum: Int)

			pub fun add(x: Int, y: Int) {
				let sum = x + y
				emit Add(x: x, y: y, sum: sum)
			}
		}
	`

	contractAccount := examples.CreateAccountWithContracts(flowClient,
		nil, []templates.Contract{{
			Name:   "EventDemo",
			Source: contract,
		}})

	// Send a tx that emits the event in the deployed contract
	script := fmt.Sprintf(`
		import EventDemo from 0x%s

		transaction {
			execute {
				EventDemo.add(x: 2, y: 3)
			}
		}
	`, contractAccount.Address.Hex())

	referenceBlockID := examples.GetReferenceBlockId(flowClient)
	runScriptTx := flow.NewTransaction().
		SetScript([]byte(script)).
		SetPayer(acctAddr).
		SetReferenceBlockID(referenceBlockID).
		SetProposalKey(acctAddr, acctKey.Index, acctKey.SequenceNumber)

	err = runScriptTx.SignEnvelope(acctAddr, acctKey.Index, acctSigner)
	examples.Handle(err)

	err = flowClient.SendTransaction(ctx, *runScriptTx)
	examples.Handle(err)

	examples.WaitForSeal(ctx, flowClient, runScriptTx.ID())

	// 1
	// Query for account creation events by type
	results, err := flowClient.GetEventsForHeightRange(ctx, client.EventRangeQuery{
		Type:        "flow.AccountCreated",
		StartHeight: 0,
		EndHeight:   100,
	})
	examples.Handle(err)

	fmt.Println("\nQuery for AccountCreated event:")
	for _, block := range results {
		for i, event := range block.Events {
			fmt.Printf("Found event #%d in block #%d\n", i+1, block.Height)
			fmt.Printf("Transaction ID: %s\n", event.TransactionID)
			fmt.Printf("Event ID: %s\n", event.ID())
			fmt.Println(event.String())
		}
	}

	// 2
	// Query for our custom event by type
	results, err = flowClient.GetEventsForHeightRange(ctx, client.EventRangeQuery{
		Type:        fmt.Sprintf("AC.%s.EventDemo.EventDemo.Add", contractAccount.Address.Hex()),
		StartHeight: 0,
		EndHeight:   100,
	})
	examples.Handle(err)

	fmt.Println("\nQuery for Add event:")
	for _, block := range results {
		for i, event := range block.Events {
			fmt.Printf("Found event #%d in block #%d\n", i+1, block.Height)
			fmt.Printf("Transaction ID: %s\n", event.TransactionID)
			fmt.Printf("Event ID: %s\n", event.ID())
			fmt.Println(event.String())
		}
	}

	// 3
	// Query by transaction
	result, err := flowClient.GetTransactionResult(ctx, runScriptTx.ID())
	examples.Handle(err)

	fmt.Println("\nQuery for tx by hash:")
	for i, event := range result.Events {
		fmt.Printf("Found event #%d\n", i+1)
		fmt.Printf("Transaction ID: %s\n", event.TransactionID)
		fmt.Printf("Event ID: %s\n", event.ID())
		fmt.Println(event.String())
	}
}

/*
 * Flow Go SDK
 *
 * Copyright 2019 Dapper Labs, Inc.
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

	"github.com/onflow/flow-go-sdk/access/http"

	"github.com/onflow/flow-go-sdk/templates"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/examples"
)

func main() {
	deployedContract, runScriptTx := preapreDemo()
	demo(deployedContract, runScriptTx)
}

func demo(deployedContract *flow.Account, runScriptTx *flow.Transaction) {
	ctx := context.Background()
	flowClient, err := http.NewClient(http.EmulatorHost)
	examples.Handle(err)

	// Query for account creation events by type
	result, err := flowClient.GetEventsForHeightRange(ctx, "flow.AccountCreated", 0, 30)
	printEvents(result, err)

	customType := flow.NewEventTypeFactory().
		WithEventName("Add").
		WithContractName("EventDemo").
		WithAddress(deployedContract.Address).
		String()

	result, err = flowClient.GetEventsForHeightRange(ctx, customType, 0, 10)
	printEvents(result, err)

	// Get events directly from transaction result
	txResult, err := flowClient.GetTransactionResult(ctx, runScriptTx.ID())
	examples.Handle(err)
	printEvent(txResult.Events)
}

func printEvents(result []flow.BlockEvents, err error) {
	examples.Handle(err)

	for _, block := range result {
		printEvent(block.Events)
	}
}

func printEvent(events []flow.Event) {
	for _, event := range events {
		fmt.Printf("\n\nType: %s", event.Type)
		fmt.Printf("\nValues: %v", event.Value)
		fmt.Printf("\nTransaction ID: %s", event.TransactionID)
	}
}

func preapreDemo() (*flow.Account, *flow.Transaction) {
	ctx := context.Background()
	flowClient, err := http.NewClient(http.EmulatorHost)
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
			prepare(auth: AuthAccount) {}
			execute {
				EventDemo.add(x: 2, y: 3)
			}
		}
	`, contractAccount.Address.Hex())

	referenceBlockID := examples.GetReferenceBlockId(flowClient)
	runScriptTx := flow.NewTransaction().
		SetScript([]byte(script)).
		SetPayer(acctAddr).
		AddAuthorizer(acctAddr).
		SetReferenceBlockID(referenceBlockID).
		SetProposalKey(acctAddr, acctKey.Index, acctKey.SequenceNumber)

	err = runScriptTx.SignEnvelope(acctAddr, acctKey.Index, acctSigner)
	examples.Handle(err)

	err = flowClient.SendTransaction(ctx, *runScriptTx)
	examples.Handle(err)

	examples.WaitForSeal(ctx, flowClient, runScriptTx.ID())

	return contractAccount, runScriptTx
}

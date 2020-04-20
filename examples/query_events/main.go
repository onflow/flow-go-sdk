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
	acctAddr, acctKey, acctSigner := examples.CreateAccount()

	flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
	examples.Handle(err)

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

	contractAddr := examples.DeployContract([]byte(contract))

	// Send a tx that emits the event in the deployed contract
	script := fmt.Sprintf(`
		import EventDemo from 0x%s

		transaction {
			execute {
				EventDemo.add(x: 2, y: 3)
			}
		}
	`, contractAddr.Hex())

	runScriptTx := flow.NewTransaction().
		SetScript([]byte(script)).
		SetPayer(acctAddr).
		SetProposalKey(acctAddr, acctKey.ID, acctKey.SequenceNumber)

	err = runScriptTx.SignEnvelope(acctAddr, acctKey.ID, acctSigner)
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
			fmt.Printf("Transaction ID: %x\n", event.TransactionID)
			fmt.Printf("Event ID: %x\n", event.ID())
			fmt.Println(event.String())
		}
	}

	// 2
	// Query for our custom event by type
	results, err = flowClient.GetEventsForHeightRange(ctx, client.EventRangeQuery{
		Type:        fmt.Sprintf("A.%s.EventDemo.Add", contractAddr.Hex()),
		StartHeight: 0,
		EndHeight:   100,
	})
	examples.Handle(err)

	fmt.Println("\nQuery for Add event:")
	for _, block := range results {
		for i, event := range block.Events {
			fmt.Printf("Found event #%d in block #%d\n", i+1, block.Height)
			fmt.Printf("Transaction ID: %x\n", event.TransactionID)
			fmt.Printf("Event ID: %x\n", event.ID())
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
		fmt.Printf("Transaction ID: %x\n", event.TransactionID)
		fmt.Printf("Event ID: %x\n", event.ID())
		fmt.Println(event.String())
	}
}

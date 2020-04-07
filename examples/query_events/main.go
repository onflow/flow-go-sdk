package main

import (
	"context"
	"fmt"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/client"
	"github.com/dapperlabs/flow-go-sdk/crypto"
	"github.com/dapperlabs/flow-go-sdk/examples"
)

func main() {
	QueryEventsDemo()
}

func QueryEventsDemo() {
	ctx := context.Background()
	accountAddr, accountKey, accountPrivateKey := examples.CreateAccount()

	flowClient, err := client.New("127.0.0.1:3569")
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
		SetPayer(accountAddr, accountKey.ID).
		SetProposalKey(accountAddr, accountKey.ID, accountKey.SequenceNumber)

	err = runScriptTx.SignContainer(
		accountAddr,
		accountKey.ID,
		crypto.NewNaiveSigner(accountPrivateKey, accountKey.HashAlgo),
	)
	examples.Handle(err)

	err = flowClient.SendTransaction(ctx, *runScriptTx)
	examples.Handle(err)

	examples.WaitForSeal(ctx, flowClient, runScriptTx.ID())

	// 1
	// Query for account creation events by type
	events, err := flowClient.GetEventsForHeightRange(ctx, client.EventRangeQuery{
		Type:        "flow.AccountCreated",
		StartHeight: 0,
		EndHeight:   100,
	})
	examples.Handle(err)

	fmt.Println("\nQuery for AccountCreated event:")
	for i, event := range events {
		fmt.Printf("Found event #%d\n", i+1)
		fmt.Printf("Transaction ID: %x\n", event.TransactionID)
		fmt.Printf("Event ID: %x\n", event.ID())
		fmt.Println(event.String())
	}

	// 2
	// Query for our custom event by type
	events, err = flowClient.GetEventsForHeightRange(ctx, client.EventRangeQuery{
		Type:        fmt.Sprintf("A.%s.EventDemo.Add", contractAddr.Hex()),
		StartHeight: 0,
		EndHeight:   100,
	})
	examples.Handle(err)

	fmt.Println("\nQuery for Add event:")
	for i, event := range events {
		fmt.Printf("Found event #%d\n", i+1)
		fmt.Printf("Transaction ID: %x\n", event.TransactionID)
		fmt.Printf("Event ID: %x\n", event.ID())
		fmt.Println(event.String())
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

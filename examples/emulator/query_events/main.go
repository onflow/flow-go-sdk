package main

import (
	"context"
	"fmt"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/client"
	"github.com/dapperlabs/flow-go-sdk/keys"

	utils "github.com/dapperlabs/flow-go-sdk/utils/examples"
)

func main() {
	QueryEventsDemo()
}

func QueryEventsDemo() {
	ctx := context.Background()
	accountKey, accountAddr := utils.CreateAccount()

	flowClient, err := client.New("127.0.0.1:3569")
	utils.Handle(err)

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

	contractAddr := utils.DeployContract([]byte(contract))

	// Send a tx that emits the event in the deployed contract
	script := fmt.Sprintf(`
		import EventDemo from 0x%s

		transaction {
			execute {
				EventDemo.add(x: 2, y: 3)
			}
		}
	`, contractAddr.Hex())

	runScriptTx := flow.Transaction{
		Script:       []byte(script),
		Nonce:        utils.GetNonce(),
		ComputeLimit: 10,
		PayerAccount: accountAddr,
	}

	sig, err := keys.SignTransaction(runScriptTx, accountKey)
	utils.Handle(err)
	runScriptTx.AddSignature(accountAddr, sig)

	err = flowClient.SendTransaction(ctx, runScriptTx)
	utils.Handle(err)

	utils.WaitForSeal(ctx, flowClient, runScriptTx.Hash())

	// 1
	// Query for account creation events by type
	events, err := flowClient.GetEvents(ctx, client.EventQuery{
		Type:       "flow.AccountCreated",
		StartBlock: 0,
		EndBlock:   100,
	})
	utils.Handle(err)

	fmt.Println("\nQuery for AccountCreated event:")
	for i, event := range events {
		fmt.Printf("Found event #%d\n", i+1)
		fmt.Println("Tx Hash: ", event.TxHash.Hex())
		fmt.Println("Event ID: ", event.ID())
		fmt.Println(event.String())
	}

	// 2
	// Query for our custom event by type
	events, err = flowClient.GetEvents(ctx, client.EventQuery{
		Type:       fmt.Sprintf("A.%s.EventDemo.Add", contractAddr.Hex()),
		StartBlock: 0,
		EndBlock:   100,
	})
	utils.Handle(err)

	fmt.Println("\nQuery for Add event:")
	for i, event := range events {
		fmt.Printf("Found event #%d\n", i+1)
		fmt.Println("Tx Hash: ", event.TxHash.Hex())
		fmt.Println("Event ID: ", event.ID())
		fmt.Println(event.String())
	}

	// 3
	// Query by transaction
	tx, err := flowClient.GetTransaction(ctx, runScriptTx.Hash())
	utils.Handle(err)

	fmt.Println("\nQuery for tx by hash:")
	for i, event := range tx.Events {
		fmt.Printf("Found event #%d\n", i+1)
		fmt.Println("Tx Hash: ", event.TxHash.Hex())
		fmt.Println("Event ID: ", event.ID())
		fmt.Println(event.String())
	}
}

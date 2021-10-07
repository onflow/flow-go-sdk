<br />
<div align="center">
  <a href="https://docs.onflow.org/sdks/">
    <img src="./sdk-banner.svg" alt="Logo" width="270" height="auto">
  </a>
  <p align="center"> <br />
    <a href="https://docs.onflow.org/flow-cli/install/"><strong>View on GitHub »</strong></a> <br /><br />
    <a href="">SDK Specifications</a> ·
    <a href="">Contribute</a> ·
    <a href="">Report a Bug</a>
  </p>
</div><br />

## Overview 

This reference guide documents all the methods available in the SDK, and explains in detail how these methods work.
SDKs are open source, and you can use them according to the licence.

The library client specifications can be found here:

[![GoDoc](https://godoc.org/github.com/onflow/flow-go-sdk?status.svg)](https://pkg.go.dev/github.com/onflow/flow-go-sdk/client)


## Getting Started

### Installing
The recommended way to install Go Flow SDK is by using Go modules. 

If you already initialized your Go project you can run the following command in your terminal:
```sh
go get github.com/onflow/flow-go-sdk
```

It's usually a good practice to pin your dependencies to a specific version. 
Refer to the [SDK releases](https://github.com/onflow/flow-go-sdk/tags) page to identify the latest version.

### Importing the Library
After the library has been installed you can import it.
```go
import "github.com/onflow/flow-go-sdk"
```

## Connect
[![GoDoc](https://godoc.org/github.com/onflow/flow-go-sdk?status.svg)](https://pkg.go.dev/github.com/onflow/flow-go-sdk/client#New)

The library uses gRPC to communicate with the access nodes, and it must be configured with correct access node API URL. 

**Access API URLs** can be found [here](https://docs.onflow.org/access-api/#flow-access-node-endpoints). An error will be returned if the host is unreachable.
The Access Nodes APIs hosted by DapperLabs are accessible at:
- Testnet `access.devnet.nodes.onflow.org:9000`
- Mainnet `access.mainnet.nodes.onflow.org:9000`
- Canarynet `access.canary.nodes.onflow.org:9000`
- Local Emulator `127.0.0.1:3569` 

Mainnet is the production Flow network, Canarynet is meant to be a development network where you actively deploy new features, 
whereas Canarynet is in traditional terms a staging environment.

Example:
```go
testnet := "access.devnet.nodes.onflow.org:9000"

flow, err := client.New(testnet)
if err != nil {
	panic("failed to establish connection with the Access API")
}
```

## Query Flow Network
After you have established a connection with the access node you can query the 
Flow network to retrieve data about blocks, accounts, events and transactions. We will explore 
how to retrieve each entity in the sections bellow.

### Get Blocks
[![GoDoc](https://godoc.org/github.com/onflow/flow-go-sdk?status.svg)](https://pkg.go.dev/github.com/onflow/flow-go-sdk/client#Client.GetBlockByHeight)

Query the network for block by id, height or get the latest block. Blocks can be sealed or unsealed where 
first expresses finalized block on the network and latter still pending for verification.

**Block ID** is SHA3-256 hash of the entire block payload, but you can get that value from the block response properties. 

**Block height** expresses the height of the block in the chain, think of it like a sequence number increasing by one for each new block. 

#### Examples

Example depicts ways to get the latest block, get the block by height and by ID:

**[>> Try this example]()**
```go
func demo() {
    ctx := context.Background()
    flowClient := examples.NewFlowClient()
    
    // get the latest sealed block
    isSealed := true
    latestBlock, err := flowClient.GetLatestBlock(ctx, isSealed)
    printBlock(latestBlock, err)
    
    // get the block by ID
    blockID := latestBlock.ID.String()
    blockByID, err := flowClient.GetBlockByID(ctx, flow.HexToID(blockID))
    printBlock(blockByID, err)
    
    // get block by height
    blockByHeight, err := flowClient.GetBlockByHeight(ctx, 0)
    printBlock(blockByHeight, err)
}

func printBlock(block *flow.Block, err error) {
    examples.Handle(err)
    
    fmt.Printf("\nID: %s\n", block.ID)
    fmt.Printf("height: %d\n", block.Height)
    fmt.Printf("timestamp: %s\n\n", block.Timestamp)
}
```
Result output:
```bash
ID: 835dc83939141097aa4297aa6cf69fc600863e3b5f9241a0d7feac1868adfa4f
height: 10
timestamp: 2021-10-06 15:06:07.105382 +0000 UTC


ID: 835dc83939141097aa4297aa6cf69fc600863e3b5f9241a0d7feac1868adfa4f
height: 10
timestamp: 2021-10-06 15:06:07.105382 +0000 UTC


ID: 7bc42fe85d32ca513769a74f97f7e1a7bad6c9407f0d934c2aa645ef9cf613c7
height: 0
timestamp: 2018-12-19 22:32:30.000000042 +0000 UTC
```

### Get Accounts
[![GoDoc](https://godoc.org/github.com/onflow/flow-go-sdk?status.svg)](https://pkg.go.dev/github.com/onflow/flow-go-sdk/client#Client.GetAccount)

Retrieve the account from the network latest's block or explicitly specify the block height from which you want to retrieve the data. 
Default get account method is actually an alias for the get account at latest block method. 

**Account address** is a unique account identifier, be mindful about the `0x` prefix, you should use the prefix as a default representation but be careful and safely handle user inputs without the prefix.

#### Examples
Example depicts ways to get account at latest block and specific block height:

**[>> Try this example]()**
```go
func demo() {
    ctx := context.Background()
    flowClient := examples.NewFlowClient()
    
    // get account from the latest block
    address := flow.HexToAddress("f8d6e0586b0a20c7")
    account, err := flowClient.GetAccount(ctx, address)
    printAccount(account, err)
    
    // get account from the block by height 0
    account, err = flowClient.GetAccountAtBlockHeight(ctx, address, 0)
    printAccount(account, err)
}
    
func printAccount(account *flow.Account, err error) {
    examples.Handle(err)
    
    fmt.Printf("\nAddress: %s", account.Address.String())
    fmt.Printf("\nBalance: %d", account.Balance)
    fmt.Printf("\nContracts: %d", len(account.Contracts))
    fmt.Printf("\nKeys: %d\n", len(account.Keys))
}
```
Result output:
```bash
Address: f8d6e0586b0a20c7
Balance: 999999999999600000
Contracts: 2
Keys: 1

Address: f8d6e0586b0a20c7
Balance: 999999999999600000
Contracts: 2
Keys: 1
```

### Get Events
[![GoDoc](https://godoc.org/github.com/onflow/flow-go-sdk?status.svg)](https://pkg.go.dev/github.com/onflow/flow-go-sdk/client#Client.GetEventsForBlockIDs)

Retrieve events by a given type and in a specified block height range or list of block IDs.

**Event type** is a string that follow a standard format `A.{contract address}.{contract name}.{event name}`. 
Please read more about [events in the documentation](https://docs.onflow.org/core-contracts/flow-token/). The exception to this standard are 
core events, and you should read more about them in [this document]().

**Block height range** expresses the height of the start and end block in the chain, think of it like a sequence number increasing by one for each new block.

#### Examples
Example depicts ways to get events within block range or by block IDs:

**[>> Try this example]()**
```go
func demo(deployedContract *flow.Account, runScriptTx *flow.Transaction) {
    ctx := context.Background()
    flowClient := examples.NewFlowClient()
    
    // Query for account creation events by type
    result, err := flowClient.GetEventsForHeightRange(ctx, client.EventRangeQuery{
        Type:        "flow.AccountCreated",
        StartHeight: 0,
        EndHeight:   100,
    })
    printEvents(result, err)
    
    // Query for our custom event by type
    customType := fmt.Sprintf("AC.%s.EventDemo.EventDemo.Add", deployedContract.Address.Hex())
    result, err = flowClient.GetEventsForHeightRange(ctx, client.EventRangeQuery{
        Type:        customType,
        StartHeight: 0,
        EndHeight:   10,
    })
    printEvents(result, err)
    
    // Get events directly from transaction result
    txResult, err := flowClient.GetTransactionResult(ctx, runScriptTx.ID())
    examples.Handle(err)
    printEvent(txResult.Events)
}

func printEvents(result []client.BlockEvents, err error) {
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
```
Example output:
```bash
Type: flow.AccountCreated
Values: flow.AccountCreated(address: 0xfd43f9148d4b725d)
Transaction ID: ba9d53c8dcb0f9c2f854f93da8467a22d053eab0c540bde0b9ca2f7ad95eb78e

Type: flow.AccountCreated
Values: flow.AccountCreated(address: 0xeb179c27144f783c)
Transaction ID: 8ab7bfef3de1cf8b2ffb36559446100bf4129a9aa88d6bc59f72a467acf0c801

...

Type: A.eb179c27144f783c.EventDemo.Add
Values: A.eb179c27144f783c.EventDemo.Add(x: 2, y: 3, sum: 5)
Transaction ID: f3a2e33687ad23b0e02644ebbdcd74a7cd8ea7214065410a8007811d0bcbd353
```


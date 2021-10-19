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

📖 **Access API URLs** can be found [here](https://docs.onflow.org/access-api/#flow-access-node-endpoints). An error will be returned if the host is unreachable.
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
// retrying requests

// TODO (error handling)


## Query Flow Network
After you have established a connection with the access node you can query the 
Flow network to retrieve data about blocks, accounts, events and transactions. We will explore 
how to retrieve each entity in the sections bellow.

### Get Blocks
[![GoDoc](https://godoc.org/github.com/onflow/flow-go-sdk?status.svg)](https://pkg.go.dev/github.com/onflow/flow-go-sdk/client#Client.GetBlockByHeight)

Query the network for block by id, height or get the latest block. Blocks can be sealed or unsealed where 
first expresses finalized block on the network and latter still pending for verification.

📖 **Block ID** is SHA3-256 hash of the entire block payload, but you can get that value from the block response properties. 

📖 **Block height** expresses the height of the block in the chain, think of it like a sequence number increasing by one for each new block. 

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

📖 **Account address** is a unique account identifier, be mindful about the `0x` prefix, you should use the prefix as a default representation but be careful and safely handle user inputs without the prefix.

Account includes the following data:
- Address: the account address.
- Balance: balance of the account.
- Contracts: list of contracts deployed to the account.
- Keys: list of keys associated with the account.

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


### Get Transactions
[![GoDoc](https://godoc.org/github.com/onflow/flow-go-sdk?status.svg)](https://pkg.go.dev/github.com/onflow/flow-go-sdk/client#Client.GetTransaction)

Retrieve transactions from the network by provided transaction ID. After a transaction has been submitted, you can also get the transaction result to check the status.

📖 **Transaction ID** is a hash of the encoded transaction payload and can be calculated before submitting the transaction to the network.

📖 **Transaction Status** represents the state of transaction in the blockchain. Status can change until is finalized.

| Status      | Final | Description |
| ----------- | ----------- | ----------- |
|   UNKNOWN    |    ❌   |   The transaction has not yet been seen by the network  |
|   PENDING    |    ❌   |   The transaction has not yet been included in a block   |
|   FINALIZED    |   ❌     |  The transaction has been included in a block   |
|   EXECUTED    |   ❌    |   The transaction has been executed but the result has not yet been sealed  |
|   SEALED    |    ✅    |   The transaction has been executed and the result is sealed in a block  |
|   EXPIRED    |   ✅     |  The transaction reference block is outdated before being executed    |

⚠️ The transactionID provided must be from the current spork.

**[>> Try this example]()**
```go
func demo(txID flow.Identifier) {
    ctx := context.Background()
    flowClient := examples.NewFlowClient()
    
    tx, err := flowClient.GetTransaction(ctx, txID)
    printTransaction(tx, err)
    
    txr, err := flowClient.GetTransactionResult(ctx, txID)
    printTransactionResult(txr, err)
}

func printTransaction(tx *flow.Transaction, err error) {
    examples.Handle(err)
    
    fmt.Printf("\nID: %s", tx.ID().String())
    fmt.Printf("\nPayer: %s", tx.Payer.String())
    fmt.Printf("\nProposer: %s", tx.ProposalKey.Address.String())
    fmt.Printf("\nAuthorizers: %s", tx.Authorizers)
}

func printTransactionResult(txr *flow.TransactionResult, err error) {
    examples.Handle(err)
    
    fmt.Printf("\nStatus: %s", txr.Status.String())
    fmt.Printf("\nError: %v", txr.Error)
}
```
Example output:
```bash
ID: fb1272c57cdad79acf2fcf37576d82bf760e3008de66aa32a900c8cd16174e1c
Payer: f8d6e0586b0a20c7
Proposer: f8d6e0586b0a20c7
Authorizers: []
Status: SEALED
Error: <nil>
```


### Get Events
[![GoDoc](https://godoc.org/github.com/onflow/flow-go-sdk?status.svg)](https://pkg.go.dev/github.com/onflow/flow-go-sdk/client#Client.GetEventsForBlockIDs)

Retrieve events by a given type and in a specified block height range or list of block IDs.

📖 **Event type** is a string that follow a standard format:
```
A.{contract address}.{contract name}.{event name}
```

Please read more about [events in the documentation](https://docs.onflow.org/core-contracts/flow-token/). The exception to this standard are 
core events, and you should read more about them in [this document]().

📖 **Block height range** expresses the height of the start and end block in the chain, think of it like a sequence number increasing by one for each new block.

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

### Get Collections
[![GoDoc](https://godoc.org/github.com/onflow/flow-go-sdk?status.svg)](https://pkg.go.dev/github.com/onflow/flow-go-sdk/client#Client.GetCollection)

Retrieve a collection which is a batch of transactions that have been included in the same block. 
Collections are used to improve consensus throughput by increasing the number of transactions per block and they act as a link between block and a transaction.

📖 **Collection ID** is SHA3-256 hash of the payload.

Example retrieving a collection:
```go
func demo(exampleCollectionID flow.Identifier) {
    ctx := context.Background()
    flowClient := examples.NewFlowClient()
    
    // get collection by ID
    collection, err := flowClient.GetCollection(ctx, exampleCollectionID)
    printCollection(collection, err)
}

func printCollection(collection *flow.Collection, err error) {
    examples.Handle(err)
    
    fmt.Printf("\nID: %s", collection.ID().String())
    fmt.Printf("\nTransactions: %s", collection.TransactionIDs)
}
```
Example output:
```bash
ID: 3d7b8037381f2497d83f2f9e09422c036aae2a59d01a7693fb6003b4d0bc3595
Transactions: [cf1184e3de4bd9a7232ca3d0b9dd2cfbf96c97888298b81a05c086451fa52ec1]
```

### Execute Scripts
[![GoDoc](https://godoc.org/github.com/onflow/flow-go-sdk?status.svg)](https://pkg.go.dev/github.com/onflow/flow-go-sdk/client#Client.ExecuteScriptAtLatestBlock)

Executing scripts lets you run non-permanent Cadence scripts on the Flow blockchain and return data. You can learn more about [Cadence and scripts here](https://docs.onflow.org/cadence/language/), but we are now only interested in executing a script code and getting back the data which is then deserialized.

We can execute a script using the latest state of the Flow blockchain or we can choose to execute the script at a specific time in history defined with block height or block ID. 

📖 **Block ID** is SHA3-256 hash of the entire block payload, but you can get that value from the block response properties.

📖 **Block height** expresses the height of the block in the chain, think of it like a sequence number increasing by one for each new block.

**[>> Try this example]()**
```go
func demo() {
    ctx := context.Background()
    flowClient := examples.NewFlowClient()
    
    script := []byte(`
        pub fun main(a: Int): Int {
            return a + 10
        }
    `)
    args := []cadence.Value{ cadence.NewInt(5) }
    value, err := flowClient.ExecuteScriptAtLatestBlock(ctx, script, args)
    
    examples.Handle(err)
    fmt.Printf("\nValue: %s", value.String())
    
    complexScript := []byte(`
        pub struct User {
            pub var balance: UFix64
            pub var address: Address
            pub var name: String
    
            init(name: String, address: Address, balance: UFix64) {
                self.name = name
                self.address = address
                self.balance = balance
            }
        }
    
        pub fun main(name: String): User {
            return User(
                name: name,
                address: 0x1,
                balance: 10.0
            )
        }
    `)
    args = []cadence.Value{ cadence.NewString("Dete") }
    value, err = flowClient.ExecuteScriptAtLatestBlock(ctx, complexScript, args)
    printComplexScript(value, err)
}

type User struct {
	balance uint64
	address flow.Address
	name string
}

func printComplexScript(value cadence.Value, err error) {
    examples.Handle(err)
    fmt.Printf("\nString value: %s", value.String())
    
    s := value.(cadence.Struct)
    u := User{
        balance: s.Fields[0].ToGoValue().(uint64),
        address: s.Fields[1].ToGoValue().([flow.AddressLength]byte),
        name:    s.Fields[2].ToGoValue().(string),
    }
    
    fmt.Printf("\nName: %s", u.name)
    fmt.Printf("\nAddress: %s", u.address.String())
    fmt.Printf("\nBalance: %d", u.balance)
}
```
Example output:
```bash
Value: 15
String value: s.34a17571e1505cf6770e6ef16ca387e345e9d54d71909f23a7ec0d671cd2faf5.User(balance: 10.00000000, address: 0x1, name: "Dete")
Name: Dete
Address: 0000000000000001
Balance: 1000000000
```

## Mutate Flow Network
Flow, like most blockchains, allows anybody to submit a transaction that mutates the shared global chain state. A transaction is an object that holds a payload, which describes the state mutation, and one or more authorizations that permit the transaction to mutate the state owned by specific accounts.

Transaction data is composed and signed with help of the SDK, signed payload of transaction then gets submitted to the access node API. If transaction is invalid or signatures are not sufficient it gets rejected. 

Executing a transaction requires couple of steps:
- [Building transaction](##build-transactions).
- Signing transaction.
- Sending transaction.

### Build Transactions
[![GoDoc](https://godoc.org/github.com/onflow/flow-go-sdk?status.svg)](https://pkg.go.dev/github.com/onflow/flow-go-sdk#Transaction)
Building a transaction involves setting the transaction script, passing arguments, setting payer, proposers and authorizers. We will describe each of these properties. 

📖 **Script** field is the portion of the transaction that describes the state mutation logic. On Flow, transaction logic is written in [Cadence](https://docs.onflow.org/cadence/). Here is an example transaction script:
```
transaction(greeting: String) {
  execute {
    log(greeting.concat(", World!"))
  }
}
```

📖 **Arguments**. A transaction can accept zero or more arguments that are passed into the Cadence script. The arguments on the transaction must match the number and order declared in the Cadence script. Sample script from above accepts a single `String` argument.

📖 **[Proposal Key](https://docs.onflow.org/concepts/transaction-signing/#proposal-key)** must be provided to act as a sequence number and prevent reply and other potential attacks.

Each account key maintains a separate transaction sequence counter; the key that lends its sequence number to a transaction is called the proposal key.

A proposal key contains three fields:
- Account address
- Key index
- Sequence number

A transaction is only valid if its declared sequence number matches the current on-chain sequence number for that key. The sequence number increments by one after the transaction is executed.

📖 **[Payer](https://docs.onflow.org/concepts/transaction-signing/#signer-roles)** is the account that pays the fees for the transaction. A transaction must specify exactly one payer. The payer is only responsible for paying the network and gas fees; the transaction is not authorized to access resources or code stored in the payer account.

📖 **[Authorizers](https://docs.onflow.org/concepts/transaction-signing/#signer-roles)** are accounts that authorize a transaction to read and mutate their resources. A transaction can specify zero or more authorizers, depending on how many accounts the transaction needs to access.

The number of authorizers on the transaction must match the number of AuthAccount parameters declared in the prepare statement of the Cadence script.

Example transaction with multiple authorizers:
```
transaction {
  prepare(authorizer1: AuthAccount, authorizer2: AuthAccount) { }
}
```

**[>> Try this example]()**
```
transaction(greeting: String) {

  let guest: Address

  prepare(authorizer: AuthAccount) {
    self.guest = authorizer.address
  }

  execute {
    log(greeting.concat(",").concat(guest.toString()))
  }
}
```

```go
import (
  "context"
  "ioutil"
  "github.com/onflow/flow-go-sdk"
  "github.com/onflow/flow-go-sdk/client"
)

func main() {

  greeting, err := outil.ReadFile("Greeting2.cdc")
  if err != nil {
    panic("failed to load Cadence script")
  }

  proposerAddress := flow.HexToAddress("9a0766d93b6608b7")
  proposerKeyIndex := 3

  payerAddress := flow.HexToAddress("631e88ae7f1d7c20")
  authorizerAddress := flow.HexToAddress("7aad92e5a0715d21")

  var accessAPIHost string

  // Establish a connection with an access node
  flowClient, err := client.New(accessAPIHost)
  if err != nil {
    panic("failed to establish connection with Access API")
  }

  // Get the latest sealed block to use as a reference block
  latestBlock, err := flowClient.GetLatestBlockHeader(context.Background(), true)
  if err != nil {
    panic("failed to fetch latest block")
  }

  // Get the latest account info for this address
  proposerAccount, err := flowClient.GetAccountAtLatestBlock(context.Background(), proposerAddress)
  if err != nil {
    panic("failed to fetch proposer account")
  }

  // Get the latest sequence number for this key
  sequenceNumber := proposerAccount.Keys[proposerKeyIndex].SequenceNumber

  tx := flow.NewTransaction().
    SetScript(greeting).
    SetGasLimit(100).
    SetReferenceBlockID(latestBlock.ID).
    SetProposalKey(proposerAddress, proposerKeyIndex, sequenceNumber).
    SetPayer(payerAddress).
    AddAuthorizer(authorizerAddress)

  // Add arguments last

  hello := cadence.NewString("Hello")

  err = tx.AddArgument(hello)
  if err != nil {
    panic("invalid argument")
  }
}
```

### Sign Transactions

### Send Transactions

[![GoDoc](https://godoc.org/github.com/onflow/flow-go-sdk?status.svg)](https://pkg.go.dev/github.com/onflow/flow-go-sdk/client#Client.SendTransaction)

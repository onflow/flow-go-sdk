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

[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/ref.svg" width="130">](https://pkg.go.dev/github.com/onflow/flow-go-sdk/client)


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
[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/ref.svg" width="130">](https://pkg.go.dev/github.com/onflow/flow-go-sdk/client#New)

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

## Query Flow Network
After you have established a connection with the access node you can query the 
Flow network to retrieve data about blocks, accounts, events and transactions. We will explore 
how to retrieve each entity in the sections bellow.

### Get Blocks
[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/ref.svg" width="130">](https://pkg.go.dev/github.com/onflow/flow-go-sdk/client#Client.GetBlockByHeight)

Query the network for block by id, height or get the latest block. Blocks can be sealed or unsealed where 
first expresses finalized block on the network and latter still pending for verification.

📖 **Block ID** is SHA3-256 hash of the entire block payload, but you can get that value from the block response properties. 

📖 **Block height** expresses the height of the block in the chain, think of it like a sequence number increasing by one for each new block. 

#### Examples

Example depicts ways to get the latest block, get the block by height and by ID:

**[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/try.svg" width="130">]()**
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
[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/ref.svg" width="130">](https://pkg.go.dev/github.com/onflow/flow-go-sdk/client#Client.GetAccount)

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

**[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/try.svg" width="130">]()**
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
[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/ref.svg" width="130">](https://pkg.go.dev/github.com/onflow/flow-go-sdk/client#Client.GetTransaction)

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

**[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/try.svg" width="130">]()**
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
[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/ref.svg" width="130">](https://pkg.go.dev/github.com/onflow/flow-go-sdk/client#Client.GetEventsForBlockIDs)

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

**[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/try.svg" width="130">]()**
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
[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/ref.svg" width="130">](https://pkg.go.dev/github.com/onflow/flow-go-sdk/client#Client.GetCollection)

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
[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/ref.svg" width="130">](https://pkg.go.dev/github.com/onflow/flow-go-sdk/client#Client.ExecuteScriptAtLatestBlock)

Executing scripts lets you run non-permanent Cadence scripts on the Flow blockchain and return data. You can learn more about [Cadence and scripts here](https://docs.onflow.org/cadence/language/), but we are now only interested in executing a script code and getting back the data which is then deserialized.

We can execute a script using the latest state of the Flow blockchain or we can choose to execute the script at a specific time in history defined with block height or block ID. 

📖 **Block ID** is SHA3-256 hash of the entire block payload, but you can get that value from the block response properties.

📖 **Block height** expresses the height of the block in the chain, think of it like a sequence number increasing by one for each new block.

**[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/try.svg" width="130">]()**
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

## Transactions
A transaction is nothing more than a signed set of data that includes script code which are instructions on how to mutate the network state and properties that define and limit it's execution. All these properties are explained bellow. 

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

📖 **Gas limit** is the limit on the amount of computation a transaction requires, and it will abort if it exceeds its gas limit.
Cadence uses metering to measure the number of operations per transaction. You can read more about it in the [Cadence documentation](/cadence).

The gas limit depends on the complexity of the transaction script. Until dedicated gas estimation tooling exists, it's best to use the emulator to test complex transactions and determine a safe limit.

📖 **Reference block** specifies an expiration window (measured in blocks) during which a transaction is considered valid by the network.
A transaction will be rejected if it is submitted past its expiry block. Flow calculates transaction expiry using the _reference block_ field on a transaction.
A transaction expires after `600` blocks are committed on top of the reference block, which takes about 10 minutes at average Mainnet block rates.

### Build Transactions
[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/ref.svg" width="130">](https://pkg.go.dev/github.com/onflow/flow-go-sdk#Transaction)

Building a transaction involves setting the required properties explained above and producing a transaction object. 

Here we define a simple transaction script that will be used to execute on the network and serve as a good learning example.

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

**[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/try.svg" width="130">]()**
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

After you have successfully [built a transaction](#build-transactions) the next step in the process is to sign it.

### Sign Transactions
[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/ref.svg" width="130">](https://pkg.go.dev/github.com/onflow/flow-go-sdk#Transaction.SignEnvelope)

Flow introduces new concepts that allow for more flexibility when creating and signing transactions.
Before trying the examples below, we recommend that you read through the [transaction signature documentation](https://docs.onflow.org/concepts/accounts-and-keys/).

After you have successfully [built a transaction](#build-transactions) the next step in the process is to sign it. Flow transactions have envelope and payload signatures, and you should learn about each in the [signature documentation](https://docs.onflow.org/concepts/accounts-and-keys/#anatomy-of-a-transaction).

Quick example of building a transaction:
```go
import (
    "github.com/onflow/flow-go-sdk"
    "github.com/onflow/flow-go-sdk/crypto"
)

var (
    myAddress    flow.Address
    myAccountKey flow.AccountKey
    myPrivateKey crypto.PrivateKey
)

tx := flow.NewTransaction().
    SetScript([]byte("transaction { execute { log(\"Hello, World!\") } }")).
    SetGasLimit(100).
    SetProposalKey(myAddress, myAccountKey.Index, myAccountKey.SequenceNumber).
    SetPayer(myAddress)
```

Transaction signing is done through the `crypto.Signer` interface. The simplest (and least secure) implementation of `crypto.Signer` is `crypto.InMemorySigner`.

Signatures can be generated more securely using keys stored in a hardware device such as an [HSM](https://en.wikipedia.org/wiki/Hardware_security_module). The `crypto.Signer` interface is intended to be flexible enough to support a variety of signer implementations and is not limited to in-memory implementations.

Simple signature example:
```go
// construct a signer from your private key and configured hash algorithm
mySigner := crypto.NewInMemorySigner(myPrivateKey, myAccountKey.HashAlgo)

err := tx.SignEnvelope(myAddress, myAccountKey.Index, mySigner)
if err != nil {
    panic("failed to sign transaction")
}
```

Flow supports great flexibility when it comes to transaction signing, we can define multiple authorizers (multi-sig transactions) and have different payer account than proposer. We will explore advanced signing scenarios bellow.

### [Single party, single signature](https://docs.onflow.org/concepts/transaction-signing/#single-party-single-signature)

- Proposer, payer and authorizer are the same account (`0x01`).
- Only the envelope must be signed.
- Proposal key must have full signing weight.

| Account | Key ID | Weight |
| ------- | ------ | ------ |
| `0x01`  | 1      | 1.0    |

**[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/try.svg" width="130">](https://github.com/onflow/flow-go-sdk/tree/master/examples#single-party-single-signature)**
```go
account1, _ := c.GetAccount(ctx, flow.HexToAddress("01"))

key1 := account1.Keys[0]

// create signer from securely-stored private key
key1Signer := getSignerForKey1()

referenceBlock, _ := flow.GetLatestBlock(ctx, true)
tx := flow.NewTransaction().
    SetScript([]byte(`
        transaction {
            prepare(signer: AuthAccount) { log(signer.address) }
        }
    `)).
    SetGasLimit(100).
    SetProposalKey(account1.Address, key1.Index, key1.SequenceNumber).
    SetReferenceBlockID(referenceBlock.ID).
    SetPayer(account1.Address).
    AddAuthorizer(account1.Address)

// account 1 signs the envelope with key 1
err := tx.SignEnvelope(account1.Address, key1.Index, key1Signer)
```


### [Single party, multiple signatures](https://docs.onflow.org/concepts/transaction-signing/#single-party-multiple-signatures)

- Proposer, payer and authorizer are the same account (`0x01`).
- Only the envelope must be signed.
- Each key has weight 0.5, so two signatures are required.

| Account | Key ID | Weight |
| ------- | ------ | ------ |
| `0x01`  | 1      | 0.5    |
| `0x01`  | 2      | 0.5    |

**[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/try.svg" width="130">](https://github.com/onflow/flow-go-sdk/tree/master/examples#single-party-multiple-signatures)**
```go
account1, _ := c.GetAccount(ctx, flow.HexToAddress("01"))

key1 := account1.Keys[0]
key2 := account1.Keys[1]

// create signers from securely-stored private keys
key1Signer := getSignerForKey1()
key2Signer := getSignerForKey2()

referenceBlock, _ := flow.GetLatestBlock(ctx, true)
tx := flow.NewTransaction().
    SetScript([]byte(`
        transaction {
            prepare(signer: AuthAccount) { log(signer.address) }
        }
    `)).
    SetGasLimit(100).
    SetProposalKey(account1.Address, key1.Index, key1.SequenceNumber).
    SetReferenceBlockID(referenceBlock.ID).
    SetPayer(account1.Address).
    AddAuthorizer(account1.Address)

// account 1 signs the envelope with key 1
err := tx.SignEnvelope(account1.Address, key1.Index, key1Signer)

// account 1 signs the envelope with key 2
err = tx.SignEnvelope(account1.Address, key2.Index, key2Signer)
```

### [Multiple parties](https://docs.onflow.org/concepts/transaction-signing/#multiple-parties)

- Proposer and authorizer are the same account (`0x01`).
- Payer is a separate account (`0x02`).
- Account `0x01` signs the payload.
- Account `0x02` signs the envelope.
    - Account `0x02` must sign last since it is the payer.

| Account | Key ID | Weight |
| ------- | ------ | ------ |
| `0x01`  | 1      | 1.0    |
| `0x02`  | 3      | 1.0    |

**[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/try.svg" width="130">](https://github.com/onflow/flow-go-sdk/tree/master/examples#multiple-parties)**
```go
account1, _ := c.GetAccount(ctx, flow.HexToAddress("01"))
account2, _ := c.GetAccount(ctx, flow.HexToAddress("02"))

key1 := account1.Keys[0]
key3 := account2.Keys[0]

// create signers from securely-stored private keys
key1Signer := getSignerForKey1()
key3Signer := getSignerForKey3()

referenceBlock, _ := flow.GetLatestBlock(ctx, true)
tx := flow.NewTransaction().
    SetScript([]byte(`
        transaction {
            prepare(signer: AuthAccount) { log(signer.address) }
        }
    `)).
    SetGasLimit(100).
    SetProposalKey(account1.Address, key1.Index, key1.SequenceNumber).
    SetReferenceBlockID(referenceBlock.ID).
    SetPayer(account2.Address).
    AddAuthorizer(account1.Address)

// account 1 signs the payload with key 1
err := tx.SignPayload(account1.Address, key1.Index, key1Signer)

// account 2 signs the envelope with key 3
// note: payer always signs last
err = tx.SignEnvelope(account2.Address, key3.Index, key3Signer)
```

### [Multiple parties, two authorizers](https://docs.onflow.org/concepts/transaction-signing/#multiple-parties)

- Proposer and authorizer are the same account (`0x01`).
- Payer is a separate account (`0x02`).
- Account `0x01` signs the payload.
- Account `0x02` signs the envelope.
    - Account `0x02` must sign last since it is the payer.
- Account `0x02` is also an authorizer to show how to include two AuthAccounts into an transaction

| Account | Key ID | Weight |
| ------- | ------ | ------ |
| `0x01`  | 1      | 1.0    |
| `0x02`  | 3      | 1.0    |

**[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/try.svg" width="130">](https://github.com/onflow/flow-go-sdk/tree/master/examples#multiple-parties-two-authorizers)**
```go
account1, _ := c.GetAccount(ctx, flow.HexToAddress("01"))
account2, _ := c.GetAccount(ctx, flow.HexToAddress("02"))

key1 := account1.Keys[0]
key3 := account2.Keys[0]

// create signers from securely-stored private keys
key1Signer := getSignerForKey1()
key3Signer := getSignerForKey3()

referenceBlock, _ := flow.GetLatestBlock(ctx, true)
tx := flow.NewTransaction().
    SetScript([]byte(`
        transaction {
            prepare(signer1: AuthAccount, signer2: AuthAccount) {
              log(signer.address)
              log(signer2.address)
          }
        }
    `)).
    SetGasLimit(100).
    SetProposalKey(account1.Address, key1.Index, key1.SequenceNumber).
    SetReferenceBlockID(referenceBlock.ID).
    SetPayer(account2.Address).
    AddAuthorizer(account1.Address).
    AddAuthorizer(account2.Address)

// account 1 signs the payload with key 1
err := tx.SignPayload(account1.Address, key1.Index, key1Signer)

// account 2 signs the envelope with key 3
// note: payer always signs last
err = tx.SignEnvelope(account2.Address, key3.Index, key3Signer)
```

### [Multiple parties, multiple signatures](https://docs.onflow.org/concepts/transaction-signing/#multiple-parties)

- Proposer and authorizer are the same account (`0x01`).
- Payer is a separate account (`0x02`).
- Account `0x01` signs the payload.
- Account `0x02` signs the envelope.
    - Account `0x02` must sign last since it is the payer.
- Both accounts must sign twice (once with each of their keys).

| Account | Key ID | Weight |
| ------- | ------ | ------ |
| `0x01`  | 1      | 0.5    |
| `0x01`  | 2      | 0.5    |
| `0x02`  | 3      | 0.5    |
| `0x02`  | 4      | 0.5    |

**[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/try.svg" width="130">](https://github.com/onflow/flow-go-sdk/tree/master/examples#multiple-parties-multiple-signatures)**
```go
account1, _ := c.GetAccount(ctx, flow.HexToAddress("01"))
account2, _ := c.GetAccount(ctx, flow.HexToAddress("02"))

key1 := account1.Keys[0]
key2 := account1.Keys[1]
key3 := account2.Keys[0]
key4 := account2.Keys[1]

// create signers from securely-stored private keys
key1Signer := getSignerForKey1()
key2Signer := getSignerForKey1()
key3Signer := getSignerForKey3()
key4Signer := getSignerForKey4()

referenceBlock, _ := flow.GetLatestBlock(ctx, true)
tx := flow.NewTransaction().
    SetScript([]byte(`
        transaction {
            prepare(signer: AuthAccount) { log(signer.address) }
        }
    `)).
    SetGasLimit(100).
    SetProposalKey(account1.Address, key1.Index, key1.SequenceNumber).
    SetReferenceBlockID(referenceBlock.ID).
    SetPayer(account2.Address).
    AddAuthorizer(account1.Address)

// account 1 signs the payload with key 1
err := tx.SignPayload(account1.Address, key1.Index, key1Signer)

// account 1 signs the payload with key 2
err = tx.SignPayload(account1.Address, key2.Index, key2Signer)

// account 2 signs the envelope with key 3
// note: payer always signs last
err = tx.SignEnvelope(account2.Address, key3.Index, key3Signer)

// account 2 signs the envelope with key 4
// note: payer always signs last
err = tx.SignEnvelope(account2.Address, key4.Index, key4Signer)
```


### Send Transactions
[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/ref.svg" width="130">](https://pkg.go.dev/github.com/onflow/flow-go-sdk/client#Client.SendTransaction)

After a transaction has been [built](#build-transactions) and [signed](#sign-transactions), it can be sent to the Flow blockchain where it will be executed. If sending was successful you can then [retrieve the transaction result](#get-transactions).


**[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/try.svg" width="130">]()**
```go
func demo(tx *flow.Transaction) {
    ctx := context.Background()
    flowClient := examples.NewFlowClient()
    
    err := flowClient.SendTransaction(ctx, *tx)
    if err != nil {
        panic(err)
    }
}
```


### Create Accounts
[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/ref.svg" width="130">](https://pkg.go.dev/github.com/onflow/flow-go-sdk/templates#CreateAccount)

On Flow, account creation happens inside a transaction. Because the network allows for a many-to-many relationship between public keys and accounts, it's not possible to derive a new account address from a public key offline. 

The Flow VM uses a deterministic address generation algorithm to assign account addresses on chain. You can find more details about address generation in the [accounts & keys documentation](https://docs.onflow.org/concepts/accounts-and-keys/).

#### Public Key
Flow uses ECDSA key pairs to control access to user accounts. Each key pair can be used in combination with the SHA2-256 or SHA3-256 hashing algorithms.

⚠️ You'll need to authorize at least one public key to control your new account.

Flow represents ECDSA public keys in raw form without additional metadata. Each key is a single byte slice containing a concatenation of its X and Y components in big-endian byte form.

A Flow account can contain zero (not possible to control) or more public keys, referred to as account keys. Read more about [accounts in the documentation](https://docs.onflow.org/concepts/accounts-and-keys/#accounts).

An account key contains the following data:
- Raw public key (described above)
- Signature algorithm
- Hash algorithm
- Weight (integer between 0-1000)

Account creation happens inside a transaction, which means that somebody must pay to submit that transaction to the network. We'll call this person the account creator. Make sure you have read [sending a transaction section](#send-transactions) first. 

```go
var (
  creatorAddress    flow.Address
  creatorAccountKey *flow.AccountKey
  creatorSigner     crypto.Signer
)

var accessAPIHost string

// Establish a connection with an access node
flowClient, err := client.New(accessAPIHost)
if err != nil {
    panic("failed to establish connection with Access API")
}

// Use the templates package to create a new account creation transaction
tx := templates.CreateAccount([]*flow.AccountKey{accountKey}, nil, creatorAddress)

// Set the transaction payer and proposal key
tx.SetPayer(creatorAddress)
tx.SetProposalKey(
    creatorAddress,
    creatorAccountKey.Index,
    creatorAccountKey.SequenceNumber,
)

// Get the latest sealed block to use as a reference block
latestBlock, err := flowClient.GetLatestBlockHeader(context.Background(), true)
if err != nil {
    panic("failed to fetch latest block")
}

tx.SetReferenceBlockID(latestBlock.ID)

// Sign and submit the transaction
err = tx.SignEnvelope(creatorAddress, creatorAccountKey.Index, creatorSigner)
if err != nil {
    panic("failed to sign transaction envelope")
}

err = flowClient.SendTransaction(context.Background(), *tx)
if err != nil {
    panic("failed to send transaction to network")
}
```

After the account creation transaction has been submitted you can retrieve the new account address by [getting the transaction result](#get-transactions). 

The new account address will be emitted in a system-level `flow.AccountCreated` event.

```go
result, err := flowClient.GetTransactionResult(ctx, tx.ID())
if err != nil {
    panic("failed to get transaction result")
}

var newAddress flow.Address

if result.Status != flow.TransactionStatusSealed {
    panic("address not known until transaction is sealed")
}

for _, event := range result.Events {
    if event.Type == flow.EventAccountCreated {
        newAddress = flow.AccountCreatedEvent(event).Address()
        break
    }
}
```

### Generate Keys
[<img src="https://raw.githubusercontent.com/onflow/sdks/main/templates/documentation/ref.svg" width="130">](https://docs.onflow.org/concepts/accounts-and-keys/#supported-signature--hash-algorithms)

Flow uses [ECDSA](https://en.wikipedia.org/wiki/Elliptic_Curve_Digital_Signature_Algorithm) signatures to control access to user accounts. Each key pair can be used in combination with the `SHA2-256` or `SHA3-256` hashing algorithms.

Here's how to generate an ECDSA private key for the P-256 (secp256r1) curve.

```go
import "github.com/onflow/flow-go-sdk/crypto"

// deterministic seed phrase
// note: this is only an example, please use a secure random generator for the key seed
seed := []byte("elephant ears space cowboy octopus rodeo potato cannon pineapple")

privateKey, err := crypto.GeneratePrivateKey(crypto.ECDSA_P256, seed)

// the private key can then be encoded as bytes (i.e. for storage) 
encPrivateKey := privateKey.Encode()
// the private key has an accompanying public key
publicKey := privateKey.PublicKey()
```

The example above uses an ECDSA key pair on the P-256 (secp256r1) elliptic curve. Flow also supports the secp256k1 curve used by Bitcoin and Ethereum. Read more about [supported algorithms here](https://docs.onflow.org/concepts/accounts-and-keys/#supported-signature--hash-algorithms).

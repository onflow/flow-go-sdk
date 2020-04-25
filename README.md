# Flow Go SDK [![GoDoc](https://godoc.org/github.com/onflow/flow-go-sdk?status.svg)](https://godoc.org/github.com/onflow/flow-go-sdk)

The Flow Go SDK provides a set of packages for Go developers to build applications that interact with the Flow network.

*Note: This SDK is also fully compatible with the [Flow Emulator](https://github.com/onflow/flow/blob/master/docs/emulator.md) and can be used for local development.*

## What is Flow?

Flow is a new blockchain for open worlds. Read more about it [here](https://github.com/onflow/flow).

## Table of Contents
<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Getting Started](#getting-started)
  - [Installing](#installing)
  - [Generating Keys](#generating-keys)
  - [Creating an Account](#creating-an-account)
  - [Signing a Transaction](#signing-a-transaction)
  - [Sending a Transaction](#sending-a-transaction)
  - [Querying Transaction Results](#querying-transaction-results)
  - [Querying Blocks](#querying-blocks)
  - [Executing a Script](#executing-a-script)
  - [Querying Events](#querying-events)
    - [Event Query Format](#event-query-format)
    - [Event Results](#event-results)
  - [Querying Accounts](#querying-accounts)
- [Examples](#examples)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Getting Started

### Installing

To start using the SDK, install Go 1.13 or above and run go get:

```sh
go get github.com/onflow/flow-go-sdk
```

### Generating Keys

The signature scheme supported in Flow accounts is ECDSA. It can be coupled with the hashing algorithms SHA2-256 or SHA3-256.

Here's how you can generate an `ECDSA-P256` private key:

```go
import "github.com/onflow/flow-go-sdk/crypto"

// deterministic seed phrase (this is only an example, please use a secure random generator for the key seed)
seed := []byte("elephant ears space cowboy octopus rodeo potato cannon pineapple")

privateKey, err := crypto.GeneratePrivateKey(crypto.ECDSA_P256, seed)
```

The private key can then be encoded as bytes (i.e. for storage):

```go
encPrivateKey := privateKey.Encode()
```

A private key has an accompanying public key:

```go
publicKey := privateKey.PublicKey()
```

The example above uses an ECDSA key-pair of the elliptic curve P-256. Flow also supports the curve secp256k1. Here's how you can generate an `ECDSA-SECp256k1` private key:

```go
privateKey, err := crypto.GeneratePrivateKey(crypto.ECDSA_secp256k1, seed)
```

### Creating an Account

Once you have [generated a key-pair](#generating-keys), you can create a new account using its public key.

```go
import (
    "github.com/onflow/flow-go-sdk"
    "github.com/onflow/flow-go-sdk/crypto"
    "github.com/onflow/flow-go-sdk/templates"
)

// generate a new private key for the account (this is only an example, please use a secure random generator for the key seed)
ctx := context.Background()
// generate a new private key for the account (this is only an example, please use a secure random generator for the key seed)
seed := []byte("elephant ears space cowboy octopus rodeo potato cannon pineapple")
privateKey, _ := crypto.GeneratePrivateKey(crypto.ECDSA_P256, seed)

// get the public key from the private key
publicKey := privateKey.PublicKey()

// construct an account key from the public key
accountKey := flow.NewAccountKey().
    SetPublicKey(publicKey).
    SetHashAlgo(crypto.SHA3_256).             // pair this key with the SHA3_256 hashing algorithm
    SetWeight(flow.AccountKeyWeightThreshold) // give this key full signing weight

// generate an account creation script
// this creates an account with a single public key and no code
script, _ := templates.CreateAccount([]*flow.AccountKey{accountKey}, nil)

// connect to an emulator running locally
c, err := client.New("localhost:3569")
if err != nil {
    panic("failed to connect to emulator")
}

payer, payerKey, payerSigner := examples.RootAccount(c)

tx := flow.NewTransaction().
    SetScript(script).
    SetGasLimit(100).
    SetProposalKey(payer, payerKey.ID, payerKey.SequenceNumber).
    SetPayer(payer)

err = tx.SignEnvelope(payer, payerKey.ID, payerSigner)
if err != nil {
    panic("failed to sign transaction")
}

err = c.SendTransaction(ctx, *tx)
if err != nil {
    panic("failed to send transaction")
}

result, err := c.GetTransactionResult(ctx, tx.ID())
if err != nil {
    panic("failed to get transaction result")
}

var myAddress flow.Address

if result.Status == flow.TransactionStatusSealed {
    for _, event := range result.Events {
        if event.Type == flow.EventAccountCreated {
            accountCreatedEvent := flow.AccountCreatedEvent(event)
            myAddress = accountCreatedEvent.Address()
            }
	}
}
```

### Signing a Transaction

Below is a simple example of how to sign a transaction using a `crypto.PrivateKey`:

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
    SetScript(script).
    SetGasLimit(100).
    SetProposalKey(myAddress, myAccountKey.ID, myAccountKey.SequenceNumber).
    SetPayer(myAddress)
```

Transaction signing is done using the `crypto.Signer` interface. The simplest (and least secure) implementation of
`crypto.Signer` is `crypto.InMemorySigner`.

Signatures can be generated more securely using hardware keys stored in a device such as an [HSM](https://en.wikipedia.org/wiki/Hardware_security_module). The `crypto.Signer` interface is intended to be flexible enough to support a variety of signer implementations and is not limited to in-memory implementations.

```go
// construct a signer from your private key and configured hash algorithm
mySigner := crypto.NewInMemorySigner(myPrivateKey, myAccountKey.HashAlgo)

err := tx.SignEnvelope(myAddress, myAccountKey.ID, mySigner)
if err != nil {
    panic("failed to sign transaction")
}
```

<!--
#### How Signatures Work in Flow

TODO: link to signatures doc
-->

### Sending a Transaction

You can submit a transaction to the network using the Access API Client.

```go
import "github.com/onflow/flow-go-sdk/client"

// connect to an emulator running locally
c, err := client.New("localhost:3569")
if err != nil {
    panic("failed to connect to emulator")
}

ctx := context.Background()

err = c.SendTransaction(ctx, tx)
if err != nil {
    panic("failed to send transaction")
}
```

### Querying Transaction Results

After you have submitted a transaction, you can query its status by ID:

```go
result, err := c.GetTransactionResult(ctx, tx.ID())
if err != nil {
    panic("failed to fetch transaction result")
}
```

The result includes a `Status` field that will be one of the following values:
- `UNKNOWN` - The transaction has not yet been seen by the network.
- `PENDING` - The transaction has not yet been included in a block.
- `FINALIZED` - The transaction has been included in a block.
- `EXECUTED` - The transaction has been executed but the result has not yet been sealed.
- `SEALED` - The transaction has been executed and the result is sealed in a block.

```go
if result.Status == flow.TransactionStatusSealed {
  fmt.Println("Transaction is sealed!")
}
```

The result also contains an `Error` that holds the error information for a failed transaction.

```go
if result.Error != nil {
    fmt.Printf("Transaction failed with error: %v\n", result.Error)
}
```

### Querying Blocks

You can use the `GetLatestBlock` method to fetch the latest sealed or unsealed block:

```go
// fetch the latest sealed block
isSealed := true
latestBlock, err := c.GetLatestBlock(ctx, isSealed)
if err != nil {
    panic("failed to fetch latest sealed block")
}

// fetch the latest unsealed block
isSealed := false
latestBlock, err := c.GetLatestBlock(ctx, isSealed)
if err != nil {
    panic("failed to fetch latest unsealed block")
}
```

A block contains the following fields:

- `ID` - The ID (hash) of the block.
- `ParentBlockID` - The ID of the previous block in the chain.
- `Height` - The height of the block in the chain.
- `CollectionGuarantees` - The list of collections included in the block.

### Executing a Script

You can use the `ExecuteScriptAtLatestBlock` method to execute a read-only script against the latest sealed execution state.

This functionality can be used to read state from the blockchain.

Scripts must be in the following form:

- A single `main` function with a single return value

This is an example of a valid script:

```
fun main(): Int { return 1 }
```

```go
import "github.com/onflow/cadence"

script := []byte("fun main(): Int { return 1 }")

value, err := c.ExecuteScript(ctx, script)
if err != nil {
    panic("failed to execute script")
}

ID := value.(cadence.Int)

// convert to Go int type
myID := ID.Int()
```

### Querying Events

You can query events with the `GetEventsForHeightRange` function:

```go
import "github.com/onflow/flow-go-sdk/client"

blocks, err := c.GetEventsForHeightRange(ctx, client.EventRangeQuery{
    Type:       "flow.AccountCreated",
    StartHeight: 10,
    EndHeight:   15,
})
if err != nil {
    panic("failed to query events")
}
```

#### Event Query Format

An event query includes the following fields:

**Type**

The event type to filter by. Event types are namespaced by the account and contract in which they are declared.

For example, a `Transfer` event that was defined in the `Token` contract deployed at account `0x55555555555555555555` will have a type of `A.0x55555555555555555555.Token.Transfer`.

Read the [language documentation](https://github.com/onflow/cadence/blob/master/docs/language.md#events) for more information on how to define and emit events in Cadence.

**StartHeight, EndHeight**

The blocks to filter by. Events will be returned from blocks in the range `StartHeight` to `EndHeight`, inclusive.

#### Event Results

The `GetEventsForHeightRange` function returns events grouped by block. Each block contains a list of events matching the query in order of execution.

```go
for _, block := range blocks {
    fmt.Printf("Events for block %s:\n", block.BlockID)
    for _, event := range block.Events {
        fmt.Printf(" - %s", event)
    }
}
```

<!--
#### Decoding an Event

TODO: example for event decoding
-->

### Querying Accounts

You can query the state of an account with the `GetAccount` function:

```go
import "github.com/onflow/flow-go-sdk"

address := flow.HexToAddress("01")

account, err := c.GetAccount(ctx, address)
if err != nil {
    panic("failed to fetch account")
}
```

A `flow.Account` contains the following fields:

- `Address: flow.Address` - The account address.
- `Balance: uint64` - The account balance.
- `Code: []byte` - The code deployed at this account.
- `Keys: []flow.AccountKey` - A list of the public keys associated with this account.

## Examples

The [examples](/examples) directory contains code samples that use the SDK to interact with the [Flow Emulator](https://github.com/onflow/flow/blob/master/docs/emulator.md).

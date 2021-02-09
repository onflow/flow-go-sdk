# Flow Go SDK 
Packages for Go developers to build applications that interact with the Flow network

[![GoDoc](https://godoc.org/github.com/onflow/flow-go-sdk?status.svg)](https://godoc.org/github.com/onflow/flow-go-sdk)

The Flow Go SDK provides a set of packages for Go developers to build applications that interact with the Flow network.

*Note: This SDK is also fully compatible with the [Flow Emulator](https://docs.onflow.org/devtools/emulator/) and can be used for local development.*

## [English](#) | [Chinese](/README_zh_CN.md)

## What is Flow?

Flow is a new blockchain for open worlds. Read more about it [here](https://github.com/onflow/flow).

## Table of Contents
<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Flow Go SDK ![GoDoc](https://godoc.org/github.com/onflow/flow-go-sdk)](#flow-go-sdk-)
  - [English | [Chinese](/README_zh_CN.md)](#english--chinese)
  - [What is Flow?](#what-is-flow)
  - [Table of Contents](#table-of-contents)
- [Getting Started](#getting-started)
  - [Installing](#installing)
  - [Generating Keys](#generating-keys)
    - [Supported Curves](#supported-curves)
  - [Creating an Account](#creating-an-account)
  - [Signing a Transaction](#signing-a-transaction)
    - [How Signatures Work in Flow](#how-signatures-work-in-flow)
      - [Single party, single signature](#single-party-single-signature)
      - [Single party, multiple signatures](#single-party-multiple-signatures)
      - [Multiple parties](#multiple-parties)
      - [Multiple parties, two authorizers](#multiple-parties-two-authorizers)
      - [Multiple parties, multiple signatures](#multiple-parties-multiple-signatures)
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

# Getting Started

## Installing

To start using the SDK, install Go 1.13 or above and run go get:

```sh
go get github.com/onflow/flow-go-sdk
```

## Generating Keys

Flow uses [ECDSA](https://en.wikipedia.org/wiki/Elliptic_Curve_Digital_Signature_Algorithm) 
to control access to user accounts. Each key pair can be used in combination with
the SHA2-256 or SHA3-256 hashing algorithms.

Here's how to generate an ECDSA private key for the P-256 (secp256r1) curve:

```go
import "github.com/onflow/flow-go-sdk/crypto"

// deterministic seed phrase
// note: this is only an example, please use a secure random generator for the key seed
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

### Supported Curves

The example above uses an ECDSA key pair on the P-256 (secp256r1) elliptic curve.
Flow also supports the secp256k1 curve used by Bitcoin and Ethereum.

Here's how to generate an ECDSA private key for the secp256k1 curve:

```go
privateKey, err := crypto.GeneratePrivateKey(crypto.ECDSA_secp256k1, seed)
```

Here's a full list of the supported signature and hash algorithms: [Flow Signature & Hash Algorithms](https://github.com/onflow/flow/blob/master/docs/accounts-and-keys.md#supported-signature--hash-algorithms)

## Creating an Account

Once you have [generated a key pair](#generating-keys), you can create a new account
using its public key.

```go
import (
    "github.com/onflow/flow-go-sdk"
    "github.com/onflow/flow-go-sdk/crypto"
    "github.com/onflow/flow-go-sdk/templates"
)

ctx := context.Background()

// generate a new private key for the account
// note: this is only an example, please use a secure random generator for the key seed
seed := []byte("elephant ears space cowboy octopus rodeo potato cannon pineapple")
privateKey, _ := crypto.GeneratePrivateKey(crypto.ECDSA_P256, seed)

// get the public key
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

payer, payerKey, payerSigner := examples.ServiceAccount(c)

tx := flow.NewTransaction().
    SetScript(script).
    SetGasLimit(100).
    SetProposalKey(payer, payerKey.Index, payerKey.SequenceNumber).
    SetPayer(payer)

err = tx.SignEnvelope(payer, payerKey.Index, payerSigner)
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

## Signing a Transaction

Below is a simple example of how to sign a transaction using a `crypto.PrivateKey`.

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

Transaction signing is done through the `crypto.Signer` interface. The simplest 
(and least secure) implementation of `crypto.Signer` is `crypto.InMemorySigner`.

Signatures can be generated more securely using keys stored in a hardware device such 
as an [HSM](https://en.wikipedia.org/wiki/Hardware_security_module). The `crypto.Signer` 
interface is intended to be flexible enough to support a variety of signer implementations 
and is not limited to in-memory implementations.

```go
// construct a signer from your private key and configured hash algorithm
mySigner := crypto.NewInMemorySigner(myPrivateKey, myAccountKey.HashAlgo)

err := tx.SignEnvelope(myAddress, myAccountKey.Index, mySigner)
if err != nil {
    panic("failed to sign transaction")
}
```

### How Signatures Work in Flow

Flow introduces new concepts that allow for more flexibility when creating and signing transactions. 
Before trying the examples below, we recommend that you read through the [transaction signature documentation](https://github.com/onflow/flow/blob/master/docs/accounts-and-keys.md#signing-a-transaction).

---

#### [Single party, single signature](https://github.com/onflow/flow/blob/master/docs/accounts-and-keys.md#single-party-single-signature)

- Proposer, payer and authorizer are the same account (`0x01`).
- Only the envelope must be signed.
- Proposal key must have full signing weight.

| Account   | Key ID | Weight |
|-----------|--------|--------|
| `0x01`    | 1      | 1.0    |

```go
account1, _ := c.GetAccount(ctx, flow.HexToAddress("01"))

key1 := account1.Keys[0]

// create signer from securely-stored private key
key1Signer := getSignerForKey1()

tx := flow.NewTransaction().
    SetScript([]byte(`
        transaction { 
            prepare(signer: AuthAccount) { log(signer.address) }
        }
    `)).
    SetGasLimit(100).
    SetProposalKey(account1.Address, key1.Index, key1.SequenceNumber).
    SetPayer(account1.Address).
    AddAuthorizer(account1.Address)

// account 1 signs the envelope with key 1
err := tx.SignEnvelope(account1.Address, key1.Index, key1Signer)
```

[Full Runnable Example](/examples#single-party-single-signature)

---

#### [Single party, multiple signatures](https://github.com/onflow/flow/blob/master/docs/accounts-and-keys.md#single-party-multiple-signatures)

- Proposer, payer and authorizer are the same account (`0x01`).
- Only the envelope must be signed.
- Each key has weight 0.5, so two signatures are required.

| Account   | Key ID | Weight |
|-----------|--------|--------|
| `0x01`    | 1      | 0.5    |
| `0x01`    | 2      | 0.5    |

```go
account1, _ := c.GetAccount(ctx, flow.HexToAddress("01"))

key1 := account1.Keys[0]
key2 := account1.Keys[1]

// create signers from securely-stored private keys
key1Signer := getSignerForKey1()
key2Signer := getSignerForKey2()

tx := flow.NewTransaction().
    SetScript([]byte(`
        transaction { 
            prepare(signer: AuthAccount) { log(signer.address) }
        }
    `)).
    SetGasLimit(100).
    SetProposalKey(account1.Address, key1.Index, key1.SequenceNumber).
    SetPayer(account1.Address).
    AddAuthorizer(account1.Address)

// account 1 signs the envelope with key 1
err := tx.SignEnvelope(account1.Address, key1.Index, key1Signer)

// account 1 signs the envelope with key 2
err = tx.SignEnvelope(account1.Address, key2.Index, key2Signer)
```

[Full Runnable Example](/examples#single-party-multiple-signatures)

---

#### [Multiple parties](https://github.com/onflow/flow/blob/master/docs/accounts-and-keys.md#multiple-parties)

- Proposer and authorizer are the same account (`0x01`).
- Payer is a separate account (`0x02`).
- Account `0x01` signs the payload.
- Account `0x02` signs the envelope.
  - Account `0x02` must sign last since it is the payer.

| Account   | Key ID | Weight |
|-----------|--------|--------|
| `0x01`    | 1      | 1.0    |
| `0x02`    | 3      | 1.0    |

```go
account1, _ := c.GetAccount(ctx, flow.HexToAddress("01"))
account2, _ := c.GetAccount(ctx, flow.HexToAddress("02"))

key1 := account1.Keys[0]
key3 := account2.Keys[0]

// create signers from securely-stored private keys
key1Signer := getSignerForKey1()
key3Signer := getSignerForKey3()

tx := flow.NewTransaction().
    SetScript([]byte(`
        transaction { 
            prepare(signer: AuthAccount) { log(signer.address) }
        }
    `)).
    SetGasLimit(100).
    SetProposalKey(account1.Address, key1.Index, key1.SequenceNumber).
    SetPayer(account2.Address).
    AddAuthorizer(account1.Address)

// account 1 signs the payload with key 1
err := tx.SignPayload(account1.Address, key1.Index, key1Signer)

// account 2 signs the envelope with key 3
// note: payer always signs last
err = tx.SignEnvelope(account2.Address, key3.Index, key3Signer)
```

[Full Runnable Example](/examples#multiple-parties)

---
#### [Multiple parties, two authorizers](https://github.com/onflow/flow/blob/master/docs/accounts-and-keys.md#multiple-parties)

- Proposer and authorizer are the same account (`0x01`).
- Payer is a separate account (`0x02`).
- Account `0x01` signs the payload.
- Account `0x02` signs the envelope.
  - Account `0x02` must sign last since it is the payer.
- Account `0x02` is also an authorizer to show how to include two AuthAccounts into an transaction

| Account   | Key ID | Weight |
|-----------|--------|--------|
| `0x01`    | 1      | 1.0    |
| `0x02`    | 3      | 1.0    |

```go
account1, _ := c.GetAccount(ctx, flow.HexToAddress("01"))
account2, _ := c.GetAccount(ctx, flow.HexToAddress("02"))

key1 := account1.Keys[0]
key3 := account2.Keys[0]

// create signers from securely-stored private keys
key1Signer := getSignerForKey1()
key3Signer := getSignerForKey3()

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
    SetPayer(account2.Address).
    AddAuthorizer(account1.Address).
    AddAuthorizer(account2.Address)

// account 1 signs the payload with key 1
err := tx.SignPayload(account1.Address, key1.Index, key1Signer)

// account 2 signs the envelope with key 3
// note: payer always signs last
err = tx.SignEnvelope(account2.Address, key3.Index, key3Signer)
```

[Full Runnable Example](/examples#multiple-parties-two-authorizers)

---

#### [Multiple parties, multiple signatures](https://github.com/onflow/flow/blob/master/docs/accounts-and-keys.md#multiple-parties-multiple-signatures)

- Proposer and authorizer are the same account (`0x01`).
- Payer is a separate account (`0x02`).
- Account `0x01` signs the payload.
- Account `0x02` signs the envelope.
  - Account `0x02` must sign last since it is the payer.
- Both accounts must sign twice (once with each of their keys).

| Account   | Key ID | Weight |
|-----------|--------|--------|
| `0x01`    | 1      | 0.5    |
| `0x01`    | 2      | 0.5    |
| `0x02`    | 3      | 0.5    |
| `0x02`    | 4      | 0.5    |

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

tx := flow.NewTransaction().
    SetScript([]byte(`
        transaction { 
            prepare(signer: AuthAccount) { log(signer.address) }
        }
    `)).
    SetGasLimit(100).
    SetProposalKey(account1.Address, key1.Index, key1.SequenceNumber).
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

[Full Runnable Example](/examples#multiple-parties-multiple-signatures)

## Sending a Transaction

You can submit a transaction to the network using the Access API client.

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

## Querying Transaction Results

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

## Querying Blocks

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

## Executing a Script

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

## Querying Events

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

### Event Query Format

An event query includes the following fields:

**Type**

The event type to filter by. Event types are namespaced by the account and contract in which they are declared.

For example, a `Transfer` event that was defined in the `Token` contract deployed at account `0x55555555555555555555` will have a type of `A.0x55555555555555555555.Token.Transfer`.

Read the [language documentation](https://docs.onflow.org/cadence/language/events/) for more information on how to define and emit events in Cadence.

**StartHeight, EndHeight**

The blocks to filter by. Events will be returned from blocks in the range `StartHeight` to `EndHeight`, inclusive.

### Event Results

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
### Decoding an Event

TODO: example for event decoding
-->

## Querying Accounts

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

The [examples](/examples) directory contains code samples that use the SDK to interact with the [Flow Emulator](https://docs.onflow.org/devtools/emulator/).

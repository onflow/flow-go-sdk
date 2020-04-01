# Flow Go SDK

The Flow Go SDK provides a set of packages for Go developers to build applications that interact with the Flow network.

*Note: This SDK is also fully compatible with the Flow Emulator and can be used for local development.*

## What is Flow?

Flow is a new blockchain for open worlds. Read more about it [here](https://www.onflow.org/).

## What would you like to do today?

**Create Accounts & Contracts**

Start here to create a user agent (wallet) or a dapp, the tools you need to get started include:
- [Protobuf Definitions](https://github.com/dapperlabs/flow-go-sdk/tree/master/protobuf)
- [Emulator](https://github.com/dapperlabs/flow-go-sdk/tree/master/cmd/flow/emulator)
- [Flow Developer Preview](https://www.notion.so/flowpreview/Flow-Developer-Preview-6d5d696c8d584398a2a025185945aa5b)

**Submit Transactions**

Then you're ready to move on to sending and submitting transactions using the [Observation API Client Library](https://github.com/dapperlabs/flow-go-sdk/blob/master/client/client.go).

**Read State**

Now it's time to see what you've done, run a script to check the outcome of your transactions:
- [Executing a script](#executing-a-script)

## Collaborating on this repo - bug reports

Please submit any bug reports in this repo directly. You'll find an issue template when you create a new issue. Please fill out the indicated fields, including:

* Problem/bug
* Steps to reproduce
* Acceptance criteria (if you have any)
* Context - specifically, what is this blocking for you?

## Table of Contents
<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Generating keys](#generating-keys)
- [Creating an account](#creating-an-account)
- [Signing a transaction](#signing-a-transaction)
  - [How signatures work in Flow](#how-signatures-work-in-flow)
  - [Adding a signature to a transaction](#adding-a-signature-to-a-transaction)
- [Submitting a transaction](#submitting-a-transaction)
- [Querying transactions](#querying-transactions)
- [Querying blocks](#querying-blocks)
- [Executing a script](#executing-a-script)
- [Querying events](#querying-events)
  - [Event query format](#event-query-format)
    - [Type](#type)
    - [StartBlock, EndBlock](#startblock-endblock)
  - [Event results](#event-results)
- [Querying accounts](#querying-accounts)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Generating keys

The signature scheme supported in Flow accounts is ECDSA. It can be coupled with the hashing algorithms SHA2-256 or SHA3-256.

Here's how you can generate an `ECDSA-P256` private key using `SHA3-256`:

```go
import "github.com/dapperlabs/flow-go-sdk/keys"

// deterministic seed phrase
seed := []byte("elephant ears space cowboy octopus rodeo potato cannon pineapple")

privateKey, err := keys.GeneratePrivateKey(keys.ECDSA_P256_SHA3_256, seed)
```

The private key can then be encoded as bytes (i.e. for storage):

```go
pkBytes, err := keys.EncodePrivateKey(privateKey)
```

A private key has an accompanying public key:

```go
// keys.PublicKeyWeightThreshold is the weight required for a key to authorize an account
publicKey := privateKey.PublicKey(keys.PublicKeyWeightThreshold)
``` 

The example above uses an ECDSA key pair of the elliptic curve P-256. Flow also supports the curve secp256k1. Here's how you can generate an `ECDSA-SECp256k1` private key using `SHA2-256` :
  
```
privateKey, err := keys.GeneratePrivateKey(keys.ECDSA_SECp256k1_SHA2_256, seed)
```

## Creating an account

Once you have [generated a key-pair](#generating-keys), you can create a new account using its public key.

```go
import (
    "github.com/dapperlabs/flow-go-sdk"
    "github.com/dapperlabs/flow-go-sdk/keys"
    "github.com/dapperlabs/flow-go-sdk/templates"
)

// generate a new private key for the account
seed := []byte("elephant ears space cowboy octopus rodeo potato cannon pineapple")
privateKey, _ := keys.GeneratePrivateKey(keys.ECDSA_P256_SHA3_256, seed)

// get the public key from the private key
publicKey := privateKey.PublicKey(keys.PublicKeyWeightThreshold)

// generate an account creation script
// this creates an account with a single public key and no code
script := templates.CreateAccount([]flow.AccountPublicKey{publicKey}, nil)

payerAddress := flow.HexToAddress("01")
payerPrivateKey := keys.MustDecodePrivateKeyHex("f87db8793077020101042075e2b5704089fabfe71cc2cc9702752c4084143f3c64b11c7ab96c85a2cb4c42a00a06082a8648ce3d030107a14403420004a848edd86e68664feb1be053fe358a2af578bebcfbc022c853a70bdc8fcb7bd50665325fcb98a3300d30f685603dffdb0e95a735580955ea68a8734d532031300203")

tx := flow.Transaction{
    Script:       script,
    Nonce:        42,
    ComputeLimit: 100,
    PayerAccount: payerAddress,
}

sig, err := keys.SignTransaction(tx, payerPrivateKey)

tx.AddSignature(payerAddress, sig)

// connect to an emulator running locally
c, err := client.New("localhost:3569")
if err != nil {
    panic("failed to connect to emulator")
}

err = c.SubmitTransaction(tx)
if err != nil {
    panic("failed to submit transaction")
}

tx = c.GetTransaction(tx.Hash())

var address flow.Address

if tx.Status == flow.TransactionSealed {
    for _, event := range tx.Receipt.Events {
        if event.Type == flow.EventAccountCreated {
	        accountCreatedEvent, _ := flow.DecodeAccountCreatedEvent(event)
            address = accountCreatedEvent.Address()
        }
    }
}
```

## Signing a transaction

Below is a simple example of how to sign a transaction using an `AccountPrivateKey`:

```go
import (
    "github.com/dapperlabs/flow-go-sdk"
    "github.com/dapperlabs/flow-go-sdk/keys"
)

var (
  myAddress    flow.Address
  myPrivateKey flow.AccountPrivateKey
)

tx := flow.Transaction{
    Script:         script,
    Nonce:          42,
    ComputeLimit:   100,
    PayerAccount:   myAddress,
    ScriptAccounts: []flow.Address{myAddress},
}

sig, err := keys.SignTransaction(tx, myPrivateKey)
if err != nil {
    panic("failed to sign transaction")
}
```

### How signatures work in Flow

You may have noticed the `PayerAccount` and `ScriptAccounts` fields in the example above. These fields are used to pre-declare the accounts that will be signing a transaction.

An account can sign a transaction in two ways:

- `PayerAccount` - The account that is paying for the gas and network fees of the transaction. A transaction can only have one payer account.
- `ScriptAccounts` - The accounts that have authorized the transaction to update their state. A transaction can be authorized by many accounts.

### Adding a signature to a transaction

When you add a signature to a transaction, it must be associated with the account that it is signing for. The account must be one of the pre-declared accounts in `PayerAccount` or `ScriptAccounts`.

```go
// add the signature to the transaction
tx.AddSignature(myAccountAddress, sig)
```

A transaction is valid if it contains a valid signature from each pre-declared account.

## Submitting a transaction

You can submit a transaction to the network using the Go Flow Client.

```go
import "github.com/dapperlabs/flow-go-sdk/client"

// connect to an emulator running locally
c, err := client.New("localhost:3569")
if err != nil {
    panic("failed to connect to emulator")
}

ctx := context.Background()

err = c.SendTransaction(ctx, tx)
if err != nil {
    panic("failed to submit transaction")
}
```

## Querying transactions

After you have submitted a transaction, you can query its status by hash:

```go
txRes, err := c.GetTransaction(ctx, tx.Hash())
if err != nil {
    panic("failed to fetch transaction")
}
```

The returned transaction includes a `Status` field that will be one of the following values:
- `PENDING` - The transaction has not yet been included in a block.
- `SEALED` - The transaction has been executed and the result is sealed in a block.

## Querying blocks

You can use the `GetLatestBlock` method to fetch the latest sealed or unsealed block header:

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

A block header contains the following fields:

- `Hash` - The hash of the block.
- `PreviousBlockhash` - The hash of the previous block in the chain.
- `Number` - The block number.

_Note: the Observation API does not yet support querying of full block data._

## Executing a script

You can use the `ExecuteScript` method to execute a read-only script against the latest sealed world state.

This functionality can be used to read state from the blockchain.

Scripts must be in the following form:

- A single `main` function with a single return value

This is an example of a valid script:

```
fun main(): Int { return 1 }
```

```go
import (
    "github.com/dapperlabs/cadence"
    "github.com/dapperlabs/cadence/encoding/xdr"
)

script := []byte("fun main(): Int { return 1 }")

value, err := c.ExecuteScript(ctx, script)
if err != nil {
    panic("failed to execute script")
}

// the returned value is XDR-encoded, so it must be decoded to its original value
myTokenID, err := xdr.Decode(language.IntType{}, result)
if err != nil {
    panic("failed to decode script result")
}

ID := myTokenID.(language.Int)

// convert to Go int type
myID := ID.Int()
```

## Querying events

You can query events with the `GetEvents` function:

```go
import "github.com/dapperlabs/flow-go-sdk/client"

events, err := c.GetEvents(ctx, client.EventQuery{
    Type:       "flow.AccountCreated",
    StartBlock: 10,
    EndBlock:   15,
})
if err != nil {
    panic("failed to query events")
}
```

### Event query format

An event query includes the following fields:

#### Type

The event type to filter by. Event types are namespaced by the transaction, account, or script in which they are declared.

For example, a `Transfer` event that was defined in code deployed at account `0x55555555555555555555` will have a type of `account.0x55555555555555555555.Transfer`.

Read the [language documentation](./language.md#events) for more information on how to define and emit events in Cadence.

####  StartBlock, EndBlock

The blocks to filter by. Events will be returned from blocks in `StartBlock` to `EndBlock`, inclusive.

### Event results

The `GetEvents` function returns a list of `flow.Event` values in the order in which they were executed.

A `flow.Event` contains the following fields:

- `ID: string` - The unique identifier for the event.
- `Type: string` - The type of the event.
- `Payload: []byte` - Event fields encoded as XDR bytes.

## Querying accounts

You can query the state of an account with the `GetAccount` function:

```go
import "github.com/dapperlabs/flow-go-sdk"

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
- `Keys: []flow.AccountPublicKey` - A list of the public keys associated with this account.

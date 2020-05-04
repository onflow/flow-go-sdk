# Emulator Examples

This package contains code samples that interact with the [Flow Emulator](https://github.com/onflow/flow/blob/master/docs/emulator.md).

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Running the emulator with the Flow CLI](#running-the-emulator-with-the-flow-cli)
  - [Installation](#installation)
  - [Starting the server](#starting-the-server)
- [Running the examples](#running-the-examples)
  - [Create Account](#create-account)
  - [Add Account Key](#add-account-key)
  - [Deploy Contract](#deploy-contract)
  - [Query Events](#query-events)
  - [Transaction Signing](#transaction-signing)
    - [Single Party, Single Signature](#single-party-single-signature)
    - [Single Party, Multiple Signatures](#single-party-multiple-signatures)
    - [Multiple Parties](#multiple-parties)
    - [Multiple Parties, Multiple Signatures](#multiple-parties-multiple-signatures)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Running the emulator with the Flow CLI

The emulator is bundled with the [Flow CLI](https://github.com/onflow/flow/blob/master/docs/cli.md), a command-line interface for working with Flow.

### Installation

Follow [these steps](https://github.com/onflow/flow/blob/master/docs/cli.md) to install the Flow CLI.

### Starting the server

Start the emulator by running the following command _in this directory_:	

```sh
flow emulator start -v
```

> The -v flag enables verbose log output, which is useful for testing

## Running the examples

In a separate process, run any of the example programs below.
Watch the emulator logs to see transaction output.

### Create Account

Create a new account on Flow.

```sh
go run ./create_account/main.go
```

### Add Account Key

Add a key to an existing account.

```sh
go run ./add_account_key/main.go
```

### Deploy Contract

Deploy a Cadence smart contract.

```sh
go run ./deploy_contract/main.go
```

### Query Events

Query events emitted by transactions.

```sh
go run ./query_events/main.go
```

### Transaction Signing

#### Single Party, Single Signature

Sign a transaction with a single account.

```sh
go run ./transaction_signing/single_party/main.go
```

#### Single Party, Multiple Signatures

Sign a transaction with a single account using multiple signatures.

```sh
go run ./transaction_signing/single_party_multisig/main.go
```

#### Multiple Parties

Sign a transaction with multiple accounts.

```sh
go run ./transaction_signing/multi_party/main.go
```

#### Multiple Parties, Multiple Signatures

Sign a transaction with multiple accounts using multiple signatures.

```sh
go run ./transaction_signing/multi_party_multisig/main.go
```

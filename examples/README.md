# Emulator Examples

This package contains code samples that interact with the [Flow Emulator](https://github.com/onflow/flow/blob/master/docs/emulator.md).

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

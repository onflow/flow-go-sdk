# Emulator Examples

This package contains code samples that interact with the [Flow Emulator](https://github.com/onflow/flow/blob/master/docs/content/emulator/index.md).

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Running the emulator with the Flow CLI](#running-the-emulator-with-the-flow-cli)
    - [Installation](#installation)
    - [Starting the server](#starting-the-server)
- [Running the examples](#running-the-examples)
  - [Get Blocks](#get-blocks)
  - [Get Accounts](#get-accounts)
  - [Get Events](#query-events)
  - [Create Account](#create-account)
  - [Add Account Key](#add-account-key)
  - [Deploy Contract](#deploy-contract)
  - [Transaction Arguments](#transaction-arguments)
  - [Transaction Signing](#transaction-signing)
    - [Single Party, Single Signature](#single-party-single-signature)
    - [Single Party, Multiple Signatures](#single-party-multiple-signatures)
    - [Multiple Parties](#multiple-parties)
    - [Multiple Parties, Two authorizers](#multiple-parties-two-authorizers)
    - [Multiple Parties, Multiple Signatures](#multiple-parties-multiple-signatures)
    - [Verify Signature](#verify-signature)
        - [User Signature](#user-signature)
        - [User Signature Verify All](#user-signature-verify-all)
        - [User Signature Verify Any](#user-signature-verify-any)
    - [Verify Events](#verify-events)
    - [Create Account](#create-account)
    - [Add Account Key](#add-account-key)
    - [Deploy Contract](#deploy-contract)
    - [Query Events](#query-events)
    - [Transaction Arguments](#transaction-arguments)
    - [Transaction Signing](#transaction-signing)
        - [Single Party, Single Signature](#single-party-single-signature)
        - [Single Party, Multiple Signatures](#single-party-multiple-signatures)
        - [Multiple Parties](#multiple-parties)
        - [Multiple Parties, Two authorizers](#multiple-parties-two-authorizers)
        - [Multiple Parties, Multiple Signatures](#multiple-parties-multiple-signatures)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Running the emulator with the Flow CLI

The emulator is bundled with the [Flow CLI](https://github.com/onflow/flow-cli/blob/master/docs/index.md), a command-line interface for working with Flow.

### Installation

Follow [these steps](https://github.com/onflow/flow-cli/blob/master/docs/index.md) to install the Flow CLI.

### Starting the server

Start the emulator by running the following command _in this directory_:	

```sh
flow emulator start -v
```

> The -v flag enables verbose log output, which is useful for testing

## Running the examples

In a separate process, run any of the example programs below.
Watch the emulator logs to see transaction output.

### Get Blocks

[Get blocks by ID, height or latest on Flow.](get_blocks/main.go)
```sh
make get-blocks
```

### Get Accounts
[Get accounts by address in specific block on Flow](get_accounts/main.go)

```sh
make get-accounts
```

### Get Events

[Get events emitted by transactions.](get_events/main.go)

```sh
make get-events
```

### Create Account

[Create a new account on Flow.](./create_account/main.go)

```sh
make create-account
```

### Add Account Key

[Add a key to an existing account.](./add_account_key/main.go)

```sh
make add-account-key
```

### Deploy Contract

[Deploy a Cadence smart contract.](./deploy_contract/main.go)

```sh
make deploy-contract
```

### Transaction Arguments

[Submit a transaction with Cadence arguments.](./transaction_arguments/main.go)

```sh
make transaction-arguments
```

### Transaction Signing

#### Single Party, Single Signature

[Sign a transaction with a single account.](./transaction_signing/single_party/main.go)

```sh
make single-party
```

#### Single Party, Multiple Signatures

[Sign a transaction with a single account using multiple signatures.](./transaction_signing/single_party_multisig/main.go)

```sh
make single-party-multisig
```

#### Multiple Parties

[Sign a transaction with multiple accounts.](./transaction_signing/multi_party/main.go)

```sh
make multi-party
```

#### Multiple Parties, Two authorizers

[Sign a transaction with multiple accounts and authorize for both of them.](./transaction_signing/multi_party_two_authorizers/main.go)

```sh
make multi-party-two-authorizers
```

#### Multiple Parties, Multiple Signatures

[Sign a transaction with multiple accounts using multiple signatures.](./transaction_signing/multi_party_multisig/main.go)

```sh
make multi-party-multisig
```

### Verify Signature

#### User Signature

[Sign an arbitrary user message.](verify_signature/user_signature/main.go)

```sh
make user-signature
```

#### User Signature Verify all

[Sign an arbitrary user message and verify it by using the public keys on an account respecting the weights of each key.](verify_signature/user_signature_validate_all/main.go)

```sh
make user-signature-verify-all
```

#### User Signature Verify any

[Sign an arbitrary user message and verify it by using the public keys on an account. Return success if any public key on the account can sign the message.](verify_signature/user_signature_validate_all/main.go)

```sh
make user-signature-verify-all
```

### Verify Events

[Verify events emitted in a block.](./verify_events/main.go)

```sh
make verify-events
```

/*
 * Flow Go SDK
 *
 * Copyright 2019 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package templates

import (
	"encoding/hex"
	"fmt"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-go-sdk/crypto"

	"github.com/onflow/flow-go-sdk"
)

// Contract is a Cadence contract deployed to a Flow account.
//
// Name is the identifier for the contract within the account.
//
// Source is the the Cadence code of the contract.
type Contract struct {
	Name   string
	Source string
}

// SourceBytes returns the UTF-8 encoded source code (Source) of the contract.
func (c Contract) SourceBytes() []byte {
	return []byte(c.Source)
}

// SourceHex returns the UTF-8 encoded source code (Source) of the contract as a hex string.
func (c Contract) SourceHex() string {
	return hex.EncodeToString(c.SourceBytes())
}

func createAccountTemplate(sigAlgo crypto.SignatureAlgorithm, hashAlgo crypto.HashAlgorithm) string {
	return fmt.Sprintf(`transaction(publicKeys: [String], contracts: {String: String}) {
		prepare(signer: AuthAccount) {
			let acct = AuthAccount(payer: signer)

			for key in publicKeys {
				let pk = PublicKey(
					publicKey: key.decodeHex(),
					signatureAlgorithm: SignatureAlgorithm.%s
				)
				acct.keys.add(key, hashAlgorithm: HashAlgorithm.%s, weight: 1000.0)
			}

			for contract in contracts.keys {
				acct.contracts.add(name: contract, code: contracts[contract]!.decodeHex())
			}
		}
	}
	`, sigAlgo.String(), hashAlgo.String())
}

// CreateAccount generates a transactions that creates a new account.
//
// This template accepts a list of public keys and a contracts argument, both of which are optional.
//
// The contracts argument is a dictionary of *contract name*: *contract code (in bytes)*.
// All of the contracts will be deployed to the account.
//
// The final argument is the address of the account that will pay the account creation fee.
// This account is added as a transaction authorizer and therefore must sign the resulting transaction.
func CreateAccount(accountKeys []*flow.AccountKey, contracts []Contract, payer flow.Address) *flow.Transaction {
	publicKeys := make([]cadence.Value, len(accountKeys))

	var sigAlgo crypto.SignatureAlgorithm
	var hashAlgo crypto.HashAlgorithm

	for i, accountKey := range accountKeys {
		if sigAlgo == crypto.UnknownSignatureAlgorithm {
			sigAlgo = accountKey.SigAlgo
			hashAlgo = accountKey.HashAlgo
		}
		keyHex := hex.EncodeToString(accountKey.Encode())
		publicKeys[i] = cadence.String(keyHex)
	}

	contractKeyPairs := make([]cadence.KeyValuePair, len(contracts))

	for i, contract := range contracts {
		contractKeyPairs[i] = cadence.KeyValuePair{
			Key:   cadence.String(contract.Name),
			Value: cadence.String(contract.SourceHex()),
		}
	}

	cadencePublicKeys := cadence.NewArray(publicKeys)
	cadenceContracts := cadence.NewDictionary(contractKeyPairs)

	return flow.NewTransaction().
		SetScript([]byte(createAccountTemplate(sigAlgo, hashAlgo))).
		AddAuthorizer(payer).
		AddRawArgument(jsoncdc.MustEncode(cadencePublicKeys)).
		AddRawArgument(jsoncdc.MustEncode(cadenceContracts))
}

const updateAccountContractTemplate = `
transaction(name: String, code: String) {
	prepare(signer: AuthAccount) {
		signer.contracts.update__experimental(name: name, code: code.decodeHex())
	}
}
`

// UpdateAccountContract generates a transaction that updates a contract deployed at an account.
func UpdateAccountContract(address flow.Address, contract Contract) *flow.Transaction {
	cadenceName := cadence.String(contract.Name)
	cadenceCode := cadence.String(contract.SourceHex())

	return flow.NewTransaction().
		SetScript([]byte(updateAccountContractTemplate)).
		AddRawArgument(jsoncdc.MustEncode(cadenceName)).
		AddRawArgument(jsoncdc.MustEncode(cadenceCode)).
		AddAuthorizer(address)
}

const addAccountContractTemplate = `
transaction(name: String, code: String) {
	prepare(signer: AuthAccount) {
		signer.contracts.add(name: name, code: code.decodeHex())
	}
}
`

// AddAccountContract generates a transaction that deploys a contract to an account.
func AddAccountContract(address flow.Address, contract Contract) *flow.Transaction {
	cadenceName := cadence.String(contract.Name)
	cadenceCode := cadence.String(contract.SourceHex())

	return flow.NewTransaction().
		SetScript([]byte(addAccountContractTemplate)).
		AddRawArgument(jsoncdc.MustEncode(cadenceName)).
		AddRawArgument(jsoncdc.MustEncode(cadenceCode)).
		AddAuthorizer(address)
}

func addAccountKeyTemplate(sigAlgo crypto.SignatureAlgorithm, hashAlgo crypto.HashAlgorithm) string {
	return fmt.Sprintf(`transaction(publicKey: String) {
		prepare(signer: AuthAccount) {
			let pk = PublicKey(
				publicKey: publicKey.decodeHex(),
				signatureAlgorithm: SignatureAlgorithm.%s
			)
			acct.keys.add(key, hashAlgorithm: HashAlgorithm.%s, weight: 1000.0)
		}
	}
	`, sigAlgo.String(), hashAlgo.String())
}

// AddAccountKey generates a transaction that adds a public key to an account.
func AddAccountKey(address flow.Address, accountKey *flow.AccountKey) *flow.Transaction {
	keyHex := hex.EncodeToString(accountKey.Encode())
	cadencePublicKey := cadence.String(keyHex)

	return flow.NewTransaction().
		SetScript([]byte(addAccountKeyTemplate(accountKey.SigAlgo, accountKey.HashAlgo))).
		AddRawArgument(jsoncdc.MustEncode(cadencePublicKey)).
		AddAuthorizer(address)
}

// RemoveAccountContract generates a transaction that removes a contract with the given name
func RemoveAccountContract(address flow.Address, contractName string) *flow.Transaction {
	cadenceName := cadence.String(contractName)

	return flow.NewTransaction().
		SetScript([]byte(removeAccountContractTemplate)).
		AddRawArgument(jsoncdc.MustEncode(cadenceName)).
		AddAuthorizer(address)
}

const removeAccountContractTemplate = `
transaction(name: String) {
	prepare(signer: AuthAccount) {
		signer.contracts.remove(name: name)
	}
}
`

/*
 * Flow Go SDK
 *
 * Copyright 2019-2020 Dapper Labs, Inc.
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
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"

	"github.com/onflow/flow-go-sdk"
)

const createAccountTemplate = `
transaction(publicKeys: [[UInt8]], code: [UInt8]) {
  prepare(signer: AuthAccount) {
	let acct = AuthAccount(payer: signer)

	for key in publicKeys {
		acct.addPublicKey(key)
	}

	acct.setCode(code)
  }
}
`

// CreateAccount generates a script that creates a new account.
func CreateAccount(accountKeys []*flow.AccountKey, code []byte, payer flow.Address) *flow.Transaction {
	publicKeys := make([]cadence.Value, len(accountKeys))

	for i, accountKey := range accountKeys {
		publicKeys[i] = bytesToCadenceArray(accountKey.Encode())
	}

	cadencePublicKeys := cadence.NewArray(publicKeys)
	cadenceCode := bytesToCadenceArray(code)

	return flow.NewTransaction().
		SetScript([]byte(createAccountTemplate)).
		AddAuthorizer(payer).
		AddRawArgument(jsoncdc.MustEncode(cadencePublicKeys)).
		AddRawArgument(jsoncdc.MustEncode(cadenceCode))
}

const updateAccountCodeTemplate = `
transaction(code: [UInt8]) {
  prepare(signer: AuthAccount) {
	signer.setCode(code)
  }
}
`

// UpdateAccountCode generates a script that updates the code associated with an account.
func UpdateAccountCode(address flow.Address, code []byte) *flow.Transaction {
	cadenceCode := bytesToCadenceArray(code)

	return flow.NewTransaction().
		SetScript([]byte(updateAccountCodeTemplate)).
		AddRawArgument(jsoncdc.MustEncode(cadenceCode)).
		AddAuthorizer(address)
}

const addAccountKeyTemplate = `
transaction(publicKey: [UInt8]) {
  prepare(signer: AuthAccount) {
	signer.addPublicKey(publicKey)
  }
}
`

// AddAccountKey generates a script that adds a key to an account.
func AddAccountKey(address flow.Address, accountKey *flow.AccountKey) *flow.Transaction {
	cadencePublicKey := bytesToCadenceArray(accountKey.Encode())

	return flow.NewTransaction().
		SetScript([]byte(addAccountKeyTemplate)).
		AddRawArgument(jsoncdc.MustEncode(cadencePublicKey)).
		AddAuthorizer(address)
}

func bytesToCadenceArray(b []byte) cadence.Array {
	values := make([]cadence.Value, len(b))

	for i, v := range b {
		values[i] = cadence.NewUInt8(v)
	}

	return cadence.NewArray(values)
}

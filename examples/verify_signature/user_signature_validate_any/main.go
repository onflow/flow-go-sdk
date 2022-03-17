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

package main

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/onflow/flow-go-sdk/client/http"

	"github.com/onflow/cadence"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/examples"
)

func main() {
	UserSignatureValidateAny()
}

var script = []byte(`
import Crypto

pub fun main(
  address: Address,
  signature: String,
  message: String
): Bool {
	let keyList = Crypto.KeyList()
	
	let account = getAccount(address)
	let keys = account.keys

	let signatureBytes = signature.decodeHex()
	let messageBytes = message.utf8

	var i = 0
	while true {
		if let key = keys.get(keyIndex: i) {
			if key.isRevoked {
				// do not check revoked keys
				i = i + 1
				continue
			}
			let pk = PublicKey(
					publicKey: key.publicKey.publicKey,
					signatureAlgorithm: key.publicKey.signatureAlgorithm
			)
			if pk.verify(
				signature: signatureBytes,
				signedData: messageBytes,
				domainSeparationTag: "FLOW-V0.0-user",
				hashAlgorithm: key.hashAlgorithm
			) {
				// this key is good

				return true
			}
		} else {
			// checked all the keys, none of them match
			return false
		}
		i = i + 1
	}

	return false
}
`)

func UserSignatureValidateAny() {
	ctx := context.Background()
	flowClient, err := http.NewDefaultEmulatorClient()
	examples.Handle(err)

	privateKeyAlice := examples.RandomPrivateKey()
	accountKeyAlice := flow.NewAccountKey().
		FromPrivateKey(privateKeyAlice).
		SetHashAlgo(crypto.SHA3_256).
		SetWeight(flow.AccountKeyWeightThreshold / 2)

	privateKeyBob := examples.RandomPrivateKey()
	accountKeyBob := flow.NewAccountKey().
		FromPrivateKey(privateKeyBob).
		SetHashAlgo(crypto.SHA3_256).
		SetWeight(flow.AccountKeyWeightThreshold / 2)

	// create the account with two keys
	account := examples.CreateAccount(flowClient, []*flow.AccountKey{accountKeyAlice, accountKeyBob})

	// create the message that will be signed with one key
	message := []byte("ananas")

	signerAlice := crypto.NewInMemorySigner(privateKeyAlice, crypto.SHA3_256)

	// sign the message only with Alice
	signatureAlice, err := flow.SignUserMessage(signerAlice, message)
	examples.Handle(err)

	// call the script to verify the signature on chain
	value, err := flowClient.ExecuteScriptAtLatestBlock(
		ctx,
		script,
		[]cadence.Value{
			cadence.BytesToAddress(account.Address.Bytes()),
			cadence.String(hex.EncodeToString(signatureAlice)),
			cadence.String(message),
		},
	)
	examples.Handle(err)

	// the signature should be valid
	if value == cadence.NewBool(true) {
		fmt.Println("Signature verification succeeded")
	} else {
		fmt.Println("Signature verification failed")
	}
}

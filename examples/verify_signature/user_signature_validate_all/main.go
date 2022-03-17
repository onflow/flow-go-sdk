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
	UserSignatureValidateAll()
}

var script = []byte(`
import Crypto

pub fun main(
  address: Address,
  signatures: [String],
  keyIndexes: [Int],
  message: String,
): Bool {
	let keyList = Crypto.KeyList()
	
	let account = getAccount(address)
	let keys = account.keys

	for keyIndex in keyIndexes {
		if let key = keys.get(keyIndex: keyIndex) {
			if key.isRevoked {
				// cannot verify: the key at this index is revoked
				return false
			}
			keyList.add(
				PublicKey(
					publicKey: key.publicKey.publicKey,
					signatureAlgorithm: key.publicKey.signatureAlgorithm
				),
				hashAlgorithm: key.hashAlgorithm,
				weight: key.weight / 1000.0,
			)
		} else {
			// cannot verify: they key at this index doesn't exist
			return false
		}
	}
	
	let signatureSet: [Crypto.KeyListSignature] = []
	
	var i = 0
	for signature in signatures {
		signatureSet.append(
			Crypto.KeyListSignature(
				keyIndex: i,
				signature: signature.decodeHex()
			)
		)
		i = i + 1
	}
	
	return keyList.verify(
		signatureSet: signatureSet,
		signedData: message.utf8,
	)
}
`)

func UserSignatureValidateAll() {
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

	// create the message that will be signed
	message := []byte("ananas")

	signerAlice := crypto.NewInMemorySigner(privateKeyAlice, crypto.SHA3_256)
	signerBob := crypto.NewInMemorySigner(privateKeyBob, crypto.SHA3_256)

	// sign the message with Alice and Bob
	signatureAlice, err := flow.SignUserMessage(signerAlice, message)
	examples.Handle(err)

	signatureBob, err := flow.SignUserMessage(signerBob, message)
	examples.Handle(err)

	signatures := cadence.NewArray([]cadence.Value{
		cadence.String(hex.EncodeToString(signatureBob)),
		cadence.String(hex.EncodeToString(signatureAlice)),
	})

	// the signature indexes correspond to the key indexes on the address
	signatureIndexes := cadence.NewArray([]cadence.Value{
		cadence.NewInt(1),
		cadence.NewInt(0),
	})

	// call the script to verify the signatures on chain
	value, err := flowClient.ExecuteScriptAtLatestBlock(
		ctx,
		script,
		[]cadence.Value{
			cadence.BytesToAddress(account.Address.Bytes()),
			signatures,
			signatureIndexes,
			cadence.String(message),
		},
	)
	examples.Handle(err)

	// the signatures should be valid
	if value == cadence.NewBool(true) {
		fmt.Println("Signature verification succeeded")
	} else {
		fmt.Println("Signature verification failed")
	}
}

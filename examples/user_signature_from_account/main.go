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
	"github.com/onflow/cadence"
	"google.golang.org/grpc"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/examples"
)

func main() {
	UserSignatureFromAccountDemo()
}

var script = []byte(`
import Crypto

pub fun main(
  address: Address,
  signatures: [String],
  keyIndexes: [Int]
  message: String,
): Bool {
	let keyList = Crypto.KeyList()
	
	let account = getAccount(address)
	let keys = account.keys
	
	var i = 0
	while true {
		if let key = keys.get(keyIndex: i) {
			if key.isRevoked {
				continue
			}
			keyList.add(
				PublicKey(
					publicKey: key.publicKey.publicKey,
					signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
				),
				hashAlgorithm: key.hashAlgorithm,
				weight: key.weight / 1000.0,
			)
			i = i + 1
		} else {
			break
		}
	}
	
	let signatureSet: [Crypto.KeyListSignature] = []
	
	var j = 0
	for signature in signatures {
		signatureSet.append(
			Crypto.KeyListSignature(
				keyIndex: keyIndexes[j],
				signature: signature.decodeHex()
			)
		)
		j = j + 1
	}
	
	return keyList.verify(
		signatureSet: signatureSet,
		signedData: message.utf8,
	)
}
`)

func UserSignatureFromAccountDemo() {
	ctx := context.Background()
	flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
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

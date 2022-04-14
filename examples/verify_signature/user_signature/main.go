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
	"github.com/onflow/flow-go-sdk/test"
)

func main() {
	UserSignatureDemo()
}

var script = []byte(`
import Crypto

pub fun main(
  rawPublicKeys: [String],
  weights: [UFix64],
  signatures: [String],
  toAddress: Address,
  fromAddress: Address,
  amount: UFix64,
): Bool {
  let keyList = Crypto.KeyList()

  var i = 0
  for rawPublicKey in rawPublicKeys {
    keyList.add(
      PublicKey(
        publicKey: rawPublicKey.decodeHex(),
        signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
      ),
      hashAlgorithm: HashAlgorithm.SHA3_256,
      weight: weights[i],
    )
    i = i + 1
  }

  let signatureSet: [Crypto.KeyListSignature] = []

  var j = 0
  for signature in signatures {
    signatureSet.append(
      Crypto.KeyListSignature(
        keyIndex: j,
        signature: signature.decodeHex()
      )
    )
    j = j + 1
  }

  // assemble the same message in cadence
  let message = toAddress.toBytes()
    .concat(fromAddress.toBytes())
    .concat(amount.toBigEndianBytes())

  return keyList.verify(
    signatureSet: signatureSet,
    signedData: message,
  )
}
`)

func UserSignatureDemo() {
	ctx := context.Background()
	flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
	examples.Handle(err)

	// create the keys
	privateKeyAlice := examples.RandomPrivateKey()
	publicKeyAlice := privateKeyAlice.PublicKey()

	privateKeyBob := examples.RandomPrivateKey()
	publicKeyBob := privateKeyBob.PublicKey()

	// create the message that will be signed
	addresses := test.AddressGenerator()

	toAddress := cadence.Address(addresses.New())
	fromAddress := cadence.Address(addresses.New())
	amount, err := cadence.NewUFix64("100.00")
	examples.Handle(err)

	message := append(toAddress.Bytes(), fromAddress.Bytes()...)
	message = append(message, amount.ToBigEndianBytes()...)

	signerAlice := crypto.NewInMemorySigner(privateKeyAlice, crypto.SHA3_256)
	signerBob := crypto.NewInMemorySigner(privateKeyBob, crypto.SHA3_256)

	// sign the message with Alice and Bob
	signatureAlice, err := flow.SignUserMessage(signerAlice, message)
	examples.Handle(err)

	signatureBob, err := flow.SignUserMessage(signerBob, message)
	examples.Handle(err)

	publicKeys := cadence.NewArray([]cadence.Value{
		cadence.String(hex.EncodeToString(publicKeyAlice.Encode())),
		cadence.String(hex.EncodeToString(publicKeyBob.Encode())),
	})

	// each signature has half weight
	weightAlice, err := cadence.NewUFix64("0.5")
	examples.Handle(err)

	weightBob, err := cadence.NewUFix64("0.5")
	examples.Handle(err)

	weights := cadence.NewArray([]cadence.Value{
		weightAlice,
		weightBob,
	})

	signatures := cadence.NewArray([]cadence.Value{
		cadence.String(hex.EncodeToString(signatureAlice)),
		cadence.String(hex.EncodeToString(signatureBob)),
	})

	// call the script to verify the signatures on chain
	value, err := flowClient.ExecuteScriptAtLatestBlock(
		ctx,
		script,
		[]cadence.Value{
			publicKeys,
			weights,
			signatures,
			toAddress,
			fromAddress,
			amount,
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

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
      Crypto.PublicKey(
        publicKey: rawPublicKey.decodeHex(),
        signatureAlgorithm: Crypto.ECDSA_P256
      ),
      hashAlgorithm: Crypto.SHA3_256,
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

  let message = toAddress.toBytes()
    .concat(fromAddress.toBytes())
    .concat(amount.toBigEndianBytes())

  return keyList.isValid(
    signatureSet: signatureSet,
    signedData: message,
  )
}
`)

func bytesToCadenceArray(b []byte) cadence.Array {
	values := make([]cadence.Value, len(b))
	for i, b := range b {
		values[i] = cadence.NewUInt8(b)
	}

	return cadence.NewArray(values)
}

func UserSignatureDemo() {
	ctx := context.Background()
	flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
	examples.Handle(err)

	privateKeyA := examples.RandomPrivateKey()
	publicKeyA := privateKeyA.PublicKey()

	privateKeyB := examples.RandomPrivateKey()
	publicKeyB := privateKeyB.PublicKey()

	addresses := test.AddressGenerator()

	toAddress := cadence.Address(addresses.New())
	fromAddress := cadence.Address(addresses.New())
	amount, err := cadence.NewUFix64("100.00")
	examples.Handle(err)

	message := append(toAddress.Bytes(), fromAddress.Bytes()...)
	message = append(message, amount.ToBigEndianBytes()...)

	signerA := crypto.NewInMemorySigner(privateKeyA, crypto.SHA3_256)
	signerB := crypto.NewInMemorySigner(privateKeyB, crypto.SHA3_256)

	signatureA, err := flow.SignUserMessage(signerA, message)
	examples.Handle(err)

	signatureB, err := flow.SignUserMessage(signerB, message)
	examples.Handle(err)

	publicKeys := cadence.NewArray([]cadence.Value{
		cadence.String(hex.EncodeToString(publicKeyA.Encode())),
		cadence.String(hex.EncodeToString(publicKeyB.Encode())),
	})

	weightA, err := cadence.NewUFix64("0.5")
	examples.Handle(err)

	weightB, err := cadence.NewUFix64("0.5")
	examples.Handle(err)

	weights := cadence.NewArray([]cadence.Value{
		weightA,
		weightB,
	})

	signatures := cadence.NewArray([]cadence.Value{
		cadence.String(hex.EncodeToString(signatureA)),
		cadence.String(hex.EncodeToString(signatureB)),
	})

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

	if value == cadence.NewBool(true) {
		fmt.Println("Signature verification succeeded")
	} else {
		fmt.Println("Signature verification failed")
	}
}

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
	"fmt"

	"github.com/onflow/cadence"
	"google.golang.org/grpc"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/examples"
)

func main() {
	UserSignatureDemo()
}

var script = []byte(`
import Crypto

pub fun main(
  rawPublicKeys: [[Int]],
  weights: [UFix64],
  message: [Int], 
  signatures: [[Int]],
): Bool {
  let keyList = Crypto.KeyList()

  var i = 0
  for rawPublicKey in rawPublicKeys {
    keyList.add(
      Crypto.PublicKey(
        publicKey: rawPublicKey,
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
        signature: signature
      )
    )
    j = j + 1
  }

  return keyList.isValid(
    signatureSet: signatureSet,
    signedData: message,
  )
}
`)

func bytesToCadenceArray(b []byte) cadence.Array {
	values := make([]cadence.Value, len(b))
	for i, b := range b {
		values[i] = cadence.NewInt(int(b))
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

	message := []byte("foo")

	signerA := crypto.NewInMemorySigner(privateKeyA, crypto.SHA3_256)
	signerB := crypto.NewInMemorySigner(privateKeyB, crypto.SHA3_256)

	signatureA, err := flow.SignUserMessage(signerA, message)
	examples.Handle(err)

	signatureB, err := flow.SignUserMessage(signerB, message)
	examples.Handle(err)

	cadenceMessage := bytesToCadenceArray(message)

	cadencePublicKeys := cadence.NewArray([]cadence.Value{
		bytesToCadenceArray(publicKeyA.Encode()),
		bytesToCadenceArray(publicKeyB.Encode()),
	})

	cadenceWeightA, err := cadence.NewUFix64("0.5")
	examples.Handle(err)

	cadenceWeightB, err := cadence.NewUFix64("0.5")
	examples.Handle(err)

	cadenceWeights := cadence.NewArray([]cadence.Value{
		cadenceWeightA,
		cadenceWeightB,
	})

	cadenceSignatures := cadence.NewArray([]cadence.Value{
		bytesToCadenceArray(signatureA),
		bytesToCadenceArray(signatureB),
	})

	value, err := flowClient.ExecuteScriptAtLatestBlock(
		ctx,
		script,
		[]cadence.Value{
			cadencePublicKeys,
			cadenceWeights,
			cadenceMessage,
			cadenceSignatures,
		},
	)
	examples.Handle(err)

	if value == cadence.NewBool(true) {
		fmt.Println("Signature verification succeeded")
	} else {
		fmt.Println("Signature verification failed")
	}
}

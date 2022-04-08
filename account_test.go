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

package flow

import (
	"crypto/rand"
	"testing"

	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/stretchr/testify/assert"
)

func generateKey() crypto.PrivateKey {
	seed := make([]byte, 32)
	_, _ = rand.Read(seed)

	privateKey, _ := crypto.GeneratePrivateKey(crypto.ECDSA_P256, seed)
	return privateKey
}

func TestAccountKey(t *testing.T) {

	t.Run("Valid", func(t *testing.T) {
		privateKey := generateKey()
		weight := 500
		index := 0
		seq := uint64(1)

		key := AccountKey{
			Index:          index,
			PublicKey:      privateKey.PublicKey(),
			SigAlgo:        privateKey.Algorithm(),
			HashAlgo:       crypto.SHA3_256,
			Weight:         weight,
			SequenceNumber: seq,
			Revoked:        false,
		}

		assert.Equal(t, privateKey.PublicKey(), key.PublicKey)
		assert.Equal(t, privateKey.Algorithm(), key.SigAlgo)
		assert.Equal(t, crypto.SHA3_256, key.HashAlgo)
		assert.Equal(t, key.Weight, weight)
		assert.Equal(t, key.SequenceNumber, seq)
		assert.Equal(t, key.Index, index)

		assert.NoError(t, key.Validate())
	})

	t.Run("Invalid Key Weights", func(t *testing.T) {
		privateKey := generateKey()
		key := AccountKey{
			SigAlgo:   privateKey.Algorithm(),
			PublicKey: privateKey.PublicKey(),
			HashAlgo:  crypto.SHA3_256,
		}

		key.SetWeight(5000)
		assert.EqualError(t, key.Validate(), "invalid key weight: 5000")

		key.SetWeight(-1)
		assert.EqualError(t, key.Validate(), "invalid key weight: -1")
	})

	t.Run("Invalid Key Algorithm", func(t *testing.T) {
		privateKey := generateKey()
		key := AccountKey{
			PublicKey: privateKey.PublicKey(),
		}

		key.SetSigAlgo(privateKey.Algorithm())
		assert.EqualError(t, key.Validate(), "signing algorithm (ECDSA_P256) is incompatible with hashing algorithm (UNKNOWN)")

		key.SetHashAlgo(crypto.SHA3_256)
		key.SetSigAlgo(crypto.UnknownSignatureAlgorithm)
		assert.EqualError(t, key.Validate(), "signing algorithm (UNKNOWN) is incompatible with hashing algorithm (SHA3_256)")
	})

}

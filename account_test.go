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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/onflow/flow-go-sdk/crypto"
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
		index := uint32(0)
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

		assert.True(t, privateKey.PublicKey().Equals(key.PublicKey))
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

	t.Run("Key Algorithm", func(t *testing.T) {
		hashAlgos := []crypto.HashAlgorithm{
			crypto.UnknownHashAlgorithm,
			crypto.SHA2_256,
			crypto.SHA2_384,
			crypto.SHA3_256,
			crypto.SHA3_384,
			crypto.Keccak256,
		}
		signAlgos := []crypto.SignatureAlgorithm{
			crypto.UnknownSignatureAlgorithm,
			crypto.ECDSA_P256,
			crypto.ECDSA_secp256k1,
		}

		validPairs := map[crypto.SignatureAlgorithm]map[crypto.HashAlgorithm]bool{
			crypto.ECDSA_P256: map[crypto.HashAlgorithm]bool{
				crypto.SHA2_256: true,
				crypto.SHA3_256: true,
			},
			crypto.ECDSA_secp256k1: map[crypto.HashAlgorithm]bool{
				crypto.SHA2_256: true,
				crypto.SHA3_256: true,
			},
		}

		key := AccountKey{}
		for _, s := range signAlgos {
			for _, h := range hashAlgos {
				key.SetSigAlgo(s)
				key.SetHashAlgo(h)
				if validPairs[s][h] {
					assert.NoError(t, key.Validate())
				} else {
					assert.EqualError(t, key.Validate(), fmt.Sprintf("signing algorithm (%s) and hashing algorithm (%s) are not a valid pair for a Flow account key", s, h))
				}
			}
		}
	})
}

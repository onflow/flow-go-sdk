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

package crypto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go-sdk/crypto"
)

func TestGeneratePrivateKey(t *testing.T) {
	// key algorithms currently supported by the SDK
	supportedAlgos := []crypto.SignatureAlgorithm{
		crypto.ECDSA_P256,
		crypto.ECDSA_secp256k1,
	}

	// key algorithms not currently supported by the SDK
	unsupportedAlgos := []crypto.SignatureAlgorithm{
		crypto.BLS_BLS12381,
	}

	// key algorithm that does not represent any valid algorithm
	invalidAlgo := crypto.SignatureAlgorithm(-42)

	emptySeed := makeSeed(0)
	shortSeed := makeSeed(crypto.MinSeedLength / 2)
	equalSeed := makeSeed(crypto.MinSeedLength)
	longSeed := makeSeed(crypto.MinSeedLength * 2)

	for _, sigAlgo := range supportedAlgos {

		t.Run(sigAlgo.String(), func(t *testing.T) {

			t.Run("Nil seed", func(t *testing.T) {
				sk, err := crypto.GeneratePrivateKey(sigAlgo, nil)
				assert.Error(t, err)
				assert.Equal(t, crypto.PrivateKey{}, sk)
			})

			t.Run("Empty seed", func(t *testing.T) {
				sk, err := crypto.GeneratePrivateKey(sigAlgo, emptySeed)
				assert.Error(t, err)
				assert.Equal(t, crypto.PrivateKey{}, sk)
			})

			t.Run("Seed length too short", func(t *testing.T) {
				sk, err := crypto.GeneratePrivateKey(sigAlgo, shortSeed)
				assert.Error(t, err)
				assert.Equal(t, crypto.PrivateKey{}, sk)
			})

			t.Run("Seed length exactly equal", func(t *testing.T) {
				sk, err := crypto.GeneratePrivateKey(sigAlgo, equalSeed)
				require.NoError(t, err)
				assert.NotEqual(t, crypto.PrivateKey{}, sk)
				assert.Equal(t, sigAlgo, sk.Algorithm())
			})

			t.Run("Valid signature algorithm", func(t *testing.T) {
				sk, err := crypto.GeneratePrivateKey(sigAlgo, longSeed)
				require.NoError(t, err)
				assert.NotEqual(t, crypto.PrivateKey{}, sk)
				assert.Equal(t, sigAlgo, sk.Algorithm())
			})

			t.Run("Deterministic generation", func(t *testing.T) {
				trials := 50

				var skA crypto.PrivateKey
				var err error

				skA, err = crypto.GeneratePrivateKey(sigAlgo, longSeed)
				require.NoError(t, err)

				for i := 0; i < trials; i++ {
					skB, err := crypto.GeneratePrivateKey(sigAlgo, longSeed)
					require.NoError(t, err)
					assert.Equal(t, skA, skB) // key should be same each time
					skA = skB
				}
			})
		})
	}

	t.Run("Unsupported algorithms", func(t *testing.T) {

		for _, sigAlgo := range unsupportedAlgos {

			t.Run(sigAlgo.String(), func(t *testing.T) {
				sk, err := crypto.GeneratePrivateKey(sigAlgo, longSeed)
				assert.Error(t, err)
				assert.Equal(t, crypto.PrivateKey{}, sk)
			})
		}
	})

	t.Run("Invalid signature algorithm", func(t *testing.T) {
		sk, err := crypto.GeneratePrivateKey(invalidAlgo, longSeed)
		assert.Error(t, err)
		assert.Equal(t, crypto.PrivateKey{}, sk)
	})
}

func makeSeed(l int) []byte {
	seed := make([]byte, l)
	for i, _ := range seed {
		seed[i] = byte(i)
	}
	return seed
}

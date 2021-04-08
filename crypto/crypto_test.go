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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go-sdk/crypto"
	fgcrypto "github.com/onflow/flow-go/crypto"
)

func TestGeneratePrivateKey(t *testing.T) {
	// key algorithms currently supported by the SDK
	supportedAlgos := []crypto.SignatureAlgorithm{
		crypto.ECDSA_P256,
		crypto.ECDSA_secp256k1,
	}

	// key algorithms not currently supported by the SDK
	unsupportedAlgos := []crypto.SignatureAlgorithm{
		fgcrypto.BLSBLS12381,
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

func TestDecodePublicKeyPEM(t *testing.T) {

	const pemECDSAKeySECP256K1 = `-----BEGIN -----
MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEaN+NInGJauSEx4ErF8GwtlNTjQvjXINA
wQ86xRvlkcKK2RSaGdKyS4Dy6NAOCucCQOvK09nBhARyqwh3VLooow==
-----END -----`
	const pemKeyWithLeadingZero = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAECi6YPHhCRPZWg0sUeNAi7QdpH5E8
hbOhaN5CWXjw0HQAZeXqjoswiWlVH0baBuwAPwFcdk5fG/KW60QvOYPExA==
-----END PUBLIC KEY-----`

	expected := map[string]string{
		pemECDSAKeySECP256K1:  "0x68df8d2271896ae484c7812b17c1b0b653538d0be35c8340c10f3ac51be591c28ad9149a19d2b24b80f2e8d00e0ae70240ebcad3d9c1840472ab087754ba28a3",
		pemKeyWithLeadingZero: "0x0a2e983c784244f656834b1478d022ed07691f913c85b3a168de425978f0d0740065e5ea8e8b308969551f46da06ec003f015c764e5f1bf296eb442f3983c4c4",
	}

	t.Run("ECDSA_P256", func(t *testing.T) {
		// generate a random public key
		sk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		require.NoError(t, err)
		pk, ok := sk.Public().(*ecdsa.PublicKey)
		require.True(t, ok)
		// encode the public key
		pkEncoding, err := x509.MarshalPKIXPublicKey(pk)
		require.NoError(t, err)
		block := pem.Block{
			Bytes: pkEncoding,
		}
		pemPkEncoding := pem.EncodeToMemory(&block)
		require.NotNil(t, pemPkEncoding)

		// decode the public key
		decodedPk, err := crypto.DecodePublicKeyPEM(crypto.ECDSA_P256, string(pemPkEncoding))
		require.NoError(t, err)
		// check the decoded key is the same as the initila key
		expectedPk := elliptic.Marshal(elliptic.P256(), pk.X, pk.Y)[1:]
		t.Logf("%x\n", expectedPk)
		assert.Equal(t, expectedPk, decodedPk.Encode())
	})

	t.Run("ECDSA_secp256k1", func(t *testing.T) {
		key := pemECDSAKeySECP256K1
		pk, err := crypto.DecodePublicKeyPEM(crypto.ECDSA_secp256k1, key)
		require.NoError(t, err)

		assert.Equal(t, expected[key], pk.String())
	})

	t.Run("Key with leading zeros", func(t *testing.T) {
		key := pemKeyWithLeadingZero
		pk, err := crypto.DecodePublicKeyPEM(crypto.ECDSA_P256, key)
		require.NoError(t, err)

		assert.Equal(t, expected[key], pk.String())
	})
}

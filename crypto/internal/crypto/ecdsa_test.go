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

package crypto

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go-sdk/crypto/internal/crypto/hash"
)

// ECDSA tests
func TestECDSA(t *testing.T) {
	ecdsaCurves := []SigningAlgorithm{
		ECDSAP256,
		ECDSASecp256k1,
	}
	for i, curve := range ecdsaCurves {
		t.Logf("Testing ECDSA for curve %s", curve)
		halg := hash.NewSHA3_256()
		testGenSignVerify(t, ecdsaCurves[i], halg)
	}
}

// Signing bench
func BenchmarkECDSAP256Sign(b *testing.B) {
	halg := hash.NewSHA3_256()
	benchSign(b, ECDSAP256, halg)
}

// Verifying bench
func BenchmarkECDSAP256Verify(b *testing.B) {
	halg := hash.NewSHA3_256()
	benchVerify(b, ECDSAP256, halg)
}

// Signing bench
func BenchmarkECDSASecp256k1Sign(b *testing.B) {
	halg := hash.NewSHA3_256()
	benchSign(b, ECDSASecp256k1, halg)
}

// Verifying bench
func BenchmarkECDSASecp256k1Verify(b *testing.B) {
	halg := hash.NewSHA3_256()
	benchVerify(b, ECDSASecp256k1, halg)
}

// ECDSA tests

// TestBLSEncodeDecode tests encoding and decoding of ECDSA keys
func TestECDSAEncodeDecode(t *testing.T) {
	ecdsaCurves := []SigningAlgorithm{
		ECDSAP256,
		ECDSASecp256k1,
	}

	for _, curve := range ecdsaCurves {
		testEncodeDecode(t, curve)
	}
}

// TestECDSAEquals tests equal for ECDSA keys
func TestECDSAEquals(t *testing.T) {
	ecdsaCurves := []SigningAlgorithm{
		ECDSAP256,
		ECDSASecp256k1,
	}
	for i, curve := range ecdsaCurves {
		testEquals(t, curve, ecdsaCurves[i]^1)
	}
}

// TestECDSAUtils tests some utility functions
func TestECDSAUtils(t *testing.T) {
	ecdsaCurves := []SigningAlgorithm{
		ECDSAP256,
		ECDSASecp256k1,
	}
	ecdsaSeedLen := []int{
		KeyGenSeedMinLenECDSAP256,
		KeyGenSeedMinLenECDSASecp256k1,
	}
	ecdsaPrKeyLen := []int{
		PrKeyLenECDSAP256,
		PrKeyLenECDSASecp256k1,
	}
	ecdsaPubKeyLen := []int{
		PubKeyLenECDSAP256,
		PubKeyLenECDSASecp256k1,
	}

	for i, curve := range ecdsaCurves {
		// generate a key pair
		seed := make([]byte, ecdsaSeedLen[i])
		n, err := rand.Read(seed)
		require.Equal(t, n, ecdsaSeedLen[i])
		require.NoError(t, err)
		sk, err := GeneratePrivateKey(curve, seed)
		require.NoError(t, err)
		testKeysAlgorithm(t, sk, ecdsaCurves[i])
		testKeySize(t, sk, ecdsaPrKeyLen[i], ecdsaPubKeyLen[i])
	}
}

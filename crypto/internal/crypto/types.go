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

// SigningAlgorithm is an identifier for a signing algorithm
// (and parameters if applicable)
type SigningAlgorithm int

const (
	// Supported signing algorithms

	UnknownSigningAlgorithm SigningAlgorithm = iota
	// BLSBLS12381 is BLS on BLS 12-381 curve
	BLSBLS12381
	// ECDSAP256 is ECDSA on NIST P-256 curve
	ECDSAP256
	// ECDSASecp256k1 is ECDSA on secp256k1 curve
	ECDSASecp256k1
)

// String returns the string representation of this signing algorithm.
func (f SigningAlgorithm) String() string {
	return [...]string{"UNKNOWN", "BLS_BLS12381", "ECDSA_P256", "ECDSA_secp256k1"}[f]
}

const (
	// minimum targeted bits of security
	securityBits = 128

	// MinSeedLen is a generic minimum seed length to guarantee a minimum entropy.
	// It is used when the seed source is not necessary a CSPRG and the seed
	// should be expanded before being passed to the generation process (key generation or other).
	MinSeedLen = 2 * (securityBits / 8)

	// ECDSA

	// NIST P256
	SignatureLenECDSAP256 = 64
	PrKeyLenECDSAP256     = 32
	// PubKeyLenECDSAP256 is the size of uncompressed points on P256
	PubKeyLenECDSAP256        = 64
	KeyGenSeedMinLenECDSAP256 = PrKeyLenECDSAP256 + (securityBits / 8)

	// SECG secp256k1
	SignatureLenECDSASecp256k1 = 64
	PrKeyLenECDSASecp256k1     = 32
	// PubKeyLenECDSASecp256k1 is the size of uncompressed points on P256
	PubKeyLenECDSASecp256k1        = 64
	KeyGenSeedMinLenECDSASecp256k1 = PrKeyLenECDSASecp256k1 + (securityBits / 8)
)

// Signature is a generic type, regardless of the signature scheme
type Signature []byte

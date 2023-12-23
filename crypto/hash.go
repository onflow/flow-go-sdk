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

package crypto

import (
	"errors"
	"fmt"

	"github.com/onflow/crypto"
	"github.com/onflow/crypto/hash"
)

type Hasher = hash.Hasher
type Hash = hash.Hash

// HashAlgorithm is an identifier for a hash algorithm.
type HashAlgorithm = hash.HashingAlgorithm

const (
	UnknownHashAlgorithm HashAlgorithm = hash.UnknownHashingAlgorithm
	SHA2_256                           = hash.SHA2_256
	SHA2_384                           = hash.SHA2_384
	SHA3_256                           = hash.SHA3_256
	SHA3_384                           = hash.SHA3_384
	Keccak256                          = hash.Keccak_256
	KMAC128                            = hash.KMAC128
)

// StringToHashAlgorithm converts a string to a HashAlgorithm.
func StringToHashAlgorithm(s string) HashAlgorithm {
	switch s {
	case SHA2_256.String():
		return SHA2_256
	case SHA3_256.String():
		return SHA3_256
	case SHA2_384.String():
		return SHA2_384
	case SHA3_384.String():
		return SHA3_384
	case Keccak256.String():
		return Keccak256
	case KMAC128.String():
		return KMAC128

	default:
		return UnknownHashAlgorithm
	}
}

// NewHasher initializes and returns a new hasher with the given hash algorithm.
//
// This function returns an error if the hash algorithm is invalid.
// KMAC128 cannot be instantiated with this function. Use `NewKMAC_128` instead.
func NewHasher(algo HashAlgorithm) (Hasher, error) {
	switch algo {
	case SHA2_256:
		return NewSHA2_256(), nil
	case SHA2_384:
		return NewSHA2_384(), nil
	case SHA3_256:
		return NewSHA3_256(), nil
	case SHA3_384:
		return NewSHA3_384(), nil
	case Keccak256:
		return NewKeccak_256(), nil
	case KMAC128:
		return nil, errors.New("KMAC128 can't be instantiated with this function")
	default:
		return nil, fmt.Errorf("invalid hash algorithm %s", algo)
	}
}

// NewSHA2_256 returns a new instance of SHA2-256 hasher.
func NewSHA2_256() Hasher {
	return hash.NewSHA2_256()
}

// NewSHA2_384 returns a new instance of SHA2-384 hasher.
func NewSHA2_384() Hasher {
	return hash.NewSHA2_384()
}

// NewSHA3_256 returns a new instance of SHA3-256 hasher.
func NewSHA3_256() Hasher {
	return hash.NewSHA3_256()
}

// NewSHA3_384 returns a new instance of SHA3-384 hasher.
func NewSHA3_384() Hasher {
	return hash.NewSHA3_384()
}

// NewKeccak_256 returns a new instance of Keccak256 hasher.
func NewKeccak_256() Hasher {
	return hash.NewKeccak_256()
}

// NewKMAC_128 returns a new KMAC instance
//   - `key` is the KMAC key (the key size is compared to the security level).
//   - `customizer` is the customization string. It can be left empty if no customization
//     is required.
//
// NewKeccak_256 returns a new instance of KMAC128
func NewKMAC_128(key []byte, customizer []byte, outputSize int) (Hasher, error) {
	return hash.NewKMAC_128(key, customizer, outputSize)
}

// NewBLSHasher returns a hasher that can be used for BLS signing and verifying.
// It abstracts the complexities of meeting the right conditions of a BLS
// hasher.
//
// The hasher returned is the the expand-message step in the BLS hash-to-curve.
// It uses a xof (extendable output function) based on KMAC128. It therefore has
// a 128-bytes output.
// The `tag` parameter is a domain separation string.
//
// Check https://pkg.go.dev/github.com/onflow/crypto#NewExpandMsgXOFKMAC128 for
// more info on the hasher generation underneath.
func NewBLSHasher(tag string) hash.Hasher {
	return crypto.NewExpandMsgXOFKMAC128(tag)
}

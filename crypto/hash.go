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
	"fmt"

	"github.com/onflow/flow-go-sdk/crypto/internal/crypto/hash"
)

type Hasher = hash.Hasher
type Hash = hash.Hash

// NewHasher initializes and returns a new hasher with the given hash algorithm.
//
// This function returns an error if the hash algorithm is invalid.
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

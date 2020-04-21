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

package hash

import (
	"hash"

	"golang.org/x/crypto/sha3"
)

// sha3_256Algo, embeds commonHasher
type sha3_256Algo struct {
	*commonHasher
	hash.Hash
}

// NewSHA3_256 returns a new instance of SHA3-256 hasher
func NewSHA3_256() Hasher {
	return &sha3_256Algo{
		commonHasher: &commonHasher{
			algo:       SHA3_256,
			outputSize: HashLenSha3_256},
		Hash: sha3.New256()}
}

// ComputeHash calculates and returns the SHA3-256 output of input byte array
func (s *sha3_256Algo) ComputeHash(data []byte) Hash {
	s.Reset()
	_, _ = s.Write(data)
	digest := make(Hash, 0, HashLenSha3_256)
	return s.Sum(digest)
}

// SumHash returns the SHA3-256 output and resets the hash state
func (s *sha3_256Algo) SumHash() Hash {
	digest := make(Hash, HashLenSha3_256)
	s.Sum(digest[:0])
	s.Reset()
	return digest
}

// sha3_384Algo, embeds commonHasher
type sha3_384Algo struct {
	*commonHasher
	hash.Hash
}

// NewSHA3_384 returns a new instance of SHA3-384 hasher
func NewSHA3_384() Hasher {
	return &sha3_384Algo{
		commonHasher: &commonHasher{
			algo:       SHA3_384,
			outputSize: HashLenSha3_384},
		Hash: sha3.New384()}
}

// ComputeHash calculates and returns the SHA3-256 output of input byte array
func (s *sha3_384Algo) ComputeHash(data []byte) Hash {
	s.Reset()
	_, _ = s.Write(data)
	digest := make(Hash, 0, HashLenSha3_384)
	return s.Sum(digest)
}

// SumHash returns the SHA3-256 output and resets the hash state
func (s *sha3_384Algo) SumHash() Hash {
	digest := make(Hash, HashLenSha3_384)
	s.Sum(digest[:0])
	s.Reset()
	return digest
}

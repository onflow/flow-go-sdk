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
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

// sha2_256Algo, embeds commonHasher
type sha2_256Algo struct {
	*commonHasher
	hash.Hash
}

// NewSHA2_256 returns a new instance of SHA2-256 hasher
func NewSHA2_256() Hasher {
	return &sha2_256Algo{
		commonHasher: &commonHasher{
			algo:       SHA2_256,
			outputSize: HashLenSha2_256},
		Hash: sha256.New()}
}

// ComputeHash calculates and returns the SHA2-256 output of input byte array
func (s *sha2_256Algo) ComputeHash(data []byte) Hash {
	s.Reset()
	_, _ = s.Write(data)
	digest := make(Hash, 0, HashLenSha2_256)
	return s.Sum(digest)
}

// SumHash returns the SHA2-256 output and resets the hash state
func (s *sha2_256Algo) SumHash() Hash {
	digest := make(Hash, HashLenSha2_256)
	s.Sum(digest[:0])
	s.Reset()
	return digest
}

// sha2_384Algo, embeds commonHasher
type sha2_384Algo struct {
	*commonHasher
	hash.Hash
}

// NewSHA2_384 returns a new instance of SHA2-384 hasher
func NewSHA2_384() Hasher {
	return &sha2_384Algo{
		commonHasher: &commonHasher{
			algo:       SHA2_384,
			outputSize: HashLenSha2_384},
		Hash: sha512.New384()}
}

// ComputeHash calculates and returns the SHA2-384 output of input byte array
func (s *sha2_384Algo) ComputeHash(data []byte) Hash {
	s.Reset()
	_, _ = s.Write(data)
	digest := make(Hash, 0, HashLenSha2_384)
	return s.Sum(digest)
}

// SumHash returns the SHA2-384 output and resets the hash state
func (s *sha2_384Algo) SumHash() Hash {
	digest := make(Hash, HashLenSha2_384)
	s.Sum(digest[:0])
	s.Reset()
	return digest
}

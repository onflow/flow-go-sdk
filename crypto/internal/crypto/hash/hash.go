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
	"bytes"
	"encoding/hex"
	"io"
)

// HashingAlgorithm is an identifier for a hashing algorithm.
type HashingAlgorithm int

const (
	// Supported hashing algorithms
	UnknownHashingAlgorithm HashingAlgorithm = iota
	SHA2_256
	SHA2_384
	SHA3_256
	SHA3_384
	KMAC128
)

// String returns the string representation of this hashing algorithm.
func (f HashingAlgorithm) String() string {
	return [...]string{"UNKNOWN", "SHA2_256", "SHA2_384", "SHA3_256", "SHA3_384", "KMAC128"}[f]
}

const (
	// Lengths of hash outputs in bytes
	HashLenSha2_256 = 32
	HashLenSha2_384 = 48
	HashLenSha3_256 = 32
	HashLenSha3_384 = 48
	// KMAC
	// the parameter maximum bytes-length as defined in NIST SP 800-185
	KmacMaxParamsLen = 2040 / 8
)

// Hash is the hash algorithms output types
type Hash []byte

// Equal checks if a hash is equal to a given hash
func (h Hash) Equal(input Hash) bool {
	return bytes.Equal(h, input)
}

// Hex returns the hex string representation of the hash.
func (h Hash) Hex() string {
	return hex.EncodeToString(h)
}

// String returns the hex string representation of the hash.
func (h Hash) String() string {
	return h.Hex()
}

// Hasher interface
type Hasher interface {
	// Algorithm returns the hashing algorithm for this hasher.
	Algorithm() HashingAlgorithm
	// Size returns the hash output length
	Size() int
	// ComputeHash returns the hash output regardless of the hash state
	ComputeHash([]byte) Hash
	// Write([]bytes) (using the io.Writer interface) adds more bytes to the
	// current hash state
	io.Writer
	// SumHash returns the hash output and resets the hash state
	SumHash() Hash
	// Reset resets the hash state
	Reset()
}

// commonHasher holds the common data for all hashers
type commonHasher struct {
	algo       HashingAlgorithm
	outputSize int
}

func (a *commonHasher) Algorithm() HashingAlgorithm {
	return a.algo
}

func BytesToHash(b []byte) Hash {
	h := make([]byte, len(b))
	copy(h, b)
	return h
}

// HashesToBytes converts a slice of hashes to a slice of byte slices.
func HashesToBytes(hashes []Hash) [][]byte {
	b := make([][]byte, len(hashes))
	for i, h := range hashes {
		b[i] = h
	}
	return b
}

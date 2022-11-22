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
	"encoding/hex"
	"sync"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/onflow/flow-go-sdk/crypto"
)

// An Identifier is a 32-byte unique identifier for an entity.
type Identifier [32]byte

// EmptyID is the empty identifier.
var EmptyID = Identifier{}

// Bytes returns the bytes representation of this identifier.
func (i Identifier) Bytes() []byte {
	return i[:]
}

// Hex returns the hexadecimal string representation of this identifier.
func (i Identifier) Hex() string {
	return hex.EncodeToString(i[:])
}

// String returns the string representation of this identifier.
func (i Identifier) String() string {
	return i.Hex()
}

// BytesToID constructs an identifier from a byte slice.
func BytesToID(b []byte) Identifier {
	var id Identifier
	copy(id[:], b)
	return id
}

// HexToID constructs an identifier from a hexadecimal string.
func HexToID(h string) Identifier {
	b, _ := hex.DecodeString(h)
	return BytesToID(b)
}

// HashToID constructs an identifier from a 32-byte hash.
func HashToID(hash []byte) Identifier {
	return BytesToID(hash)
}

// BytesToHash constructs a crypto hash from byte slice.
func BytesToHash(hash []byte) crypto.Hash {
	h := make(crypto.Hash, len(hash))
	copy(h, hash)
	return h
}

type StateCommitment Identifier

// BytesToStateCommitment constructs a state commitment from a byte slice.
func BytesToStateCommitment(b []byte) StateCommitment {
	return StateCommitment(BytesToID(b))
}

// HexToStateCommitment constructs a state commitment from a hexadecimal string.
func HexToStateCommitment(h string) StateCommitment {
	return StateCommitment(HexToID(h))
}

// HashToStateCommitment constructs a state commitment from a 32-byte hash.
func HashToStateCommitment(hash []byte) StateCommitment {
	return StateCommitment(HashToID(hash))
}

// A ChainID is a unique identifier for a specific Flow network instance.
//
// Chain IDs are used used to prevent replay attacks and to support network-specific address generation.
type ChainID string

const (
	// Mainnet is the chain ID for the mainnet chain.
	Mainnet ChainID = "flow-mainnet"

	// Long-lived test networks

	// Testnet is the chain ID for the testnet chain.
	Testnet ChainID = "flow-testnet"

	// Sandboxnet is the chain ID for sandboxnet chain.
	Sandboxnet ChainID = "flow-sandboxnet"

	// Transient test networks

	// Benchnet is the chain ID for the transient benchmarking chain.
	Benchnet ChainID = "flow-benchnet"
	// Localnet is the chain ID for the local development chain.
	Localnet ChainID = "flow-localnet"
	// Emulator is the chain ID for the emulated chain.
	Emulator ChainID = "flow-emulator"
	// BftTestnet is the chain ID for testing attack vector scenarios.
	BftTestnet ChainID = "flow-bft-test-net"

	// MonotonicEmulator is the chain ID for the emulated node chain with monotonic address generation.
	MonotonicEmulator ChainID = "flow-emulator-monotonic"
)

func (id ChainID) String() string {
	return string(id)
}

// entityHasher is a thread-safe hasher used to hash Flow entities.
type entityHasher struct {
	mut    sync.Mutex
	hasher crypto.Hasher
}

func (h *entityHasher) ComputeHash(b []byte) crypto.Hash {
	h.mut.Lock()
	defer h.mut.Unlock()
	return h.hasher.ComputeHash(b)
}

// defaultEntityHasher is the default hasher used to compute Flow identifiers.
var defaultEntityHasher *entityHasher

func init() {
	defaultEntityHasher = &entityHasher{
		hasher: crypto.NewSHA3_256(),
	}
}

func rlpEncode(v interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(v)
}

func rlpDecode(b []byte, v interface{}) error {
	return rlp.DecodeBytes(b, v)
}

func mustRLPEncode(v interface{}) []byte {
	// Note(sideninja): This is a temporary workaround until cadence defines canonical format addressing the issue https://github.com/onflow/flow-go-sdk/issues/286
	if tx, ok := v.(payloadCanonicalForm); ok {
		for _, arg := range tx.Arguments {
			if arg[len(arg)-1] == byte(10) {
				arg = arg[:len(tx.Arguments)-1]
			}
		}
	}

	b, err := rlpEncode(v)
	if err != nil {
		panic(err)
	}
	return b
}

func mustRLPDecode(b []byte, v interface{}) {
	err := rlpDecode(b, v)
	if err != nil {
		panic(err)
	}
	return
}

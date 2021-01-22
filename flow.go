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

// A ChainID is a unique identifier for a specific Flow network instance.
//
// Chain IDs are used used to prevent replay attacks and to support network-specific address generation.
type ChainID string

// Mainnet is the chain ID for the mainnet node chain.
const Mainnet ChainID = "flow-mainnet"

// Testnet is the chain ID for the testnet node chain.
const Testnet ChainID = "flow-testnet"

// Emulator is the chain ID for the emulated node chain.
const Emulator ChainID = "flow-emulator"

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

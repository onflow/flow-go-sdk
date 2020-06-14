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

// A Collection is a list of transactions bundled together for inclusion in a block.
type Collection struct {
	TransactionIDs []Identifier
}

// ID returns the canonical SHA3-256 hash of this collection.
func (c Collection) ID() Identifier {
	return HashToID(defaultEntityHasher.ComputeHash(c.Encode()))
}

// Encode returns the canonical RLP byte representation of this collection.
func (c Collection) Encode() []byte {
	transactionIDs := make([][]byte, len(c.TransactionIDs))
	for i, id := range c.TransactionIDs {
		transactionIDs[i] = id.Bytes()
	}

	temp := struct {
		TransactionIDs [][]byte
	}{
		TransactionIDs: transactionIDs,
	}
	return mustRLPEncode(&temp)
}

// A CollectionGuarantee is an attestation signed by the nodes that have guaranteed a collection.
type CollectionGuarantee struct {
	CollectionID Identifier
}

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

import "time"

// Block is a set of state mutations applied to the Flow blockchain.
type Block struct {
	BlockHeader
	BlockPayload
}

// BlockHeader is a summary of a full block.
type BlockHeader struct {
	ID        Identifier
	ParentID  Identifier
	Height    uint64
	Timestamp time.Time
}

// BlockPayload is the full contents of a block.
//
// A payload contains the collection guarantees and seals for a block.
type BlockPayload struct {
	CollectionGuarantees []*CollectionGuarantee
	Seals                []*BlockSeal
}

// TODO: define block seal struct
type BlockSeal struct{}

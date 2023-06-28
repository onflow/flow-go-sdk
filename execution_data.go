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

type ExecutionData struct {
	BlockID            Identifier
	ChunkExecutionData []*ChunkExecutionData
}

type ExecutionDataStreamResponse struct {
	Height        uint64
	ExecutionData *ExecutionData
}

type ChunkExecutionData struct {
	Transactions []*Transaction
	Events       []*Event
	TrieUpdate   *TrieUpdate
}

type TrieUpdate struct {
	RootHash []byte
	Paths    [][]byte
	Payloads []*Payload
}

type Payload struct {
	KeyPart []*KeyPart
	Value   []byte
}

type KeyPart struct {
	Type  uint16
	Value []byte
}

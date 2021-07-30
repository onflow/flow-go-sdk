/*
 * Flow Go SDK
 *
 * Copyright 2019-2021 Dapper Labs, Inc.
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
	"fmt"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go/crypto/hash"
)

// List of built-in account event types.
const (
	EventAccountCreated         string = "flow.AccountCreated"
	EventAccountAdded           string = "flow.AccountKeyAdded"
	EventAccountKeyRemoved      string = "flow.AccountKeyRemoved"
	EventAccountContractAdded   string = "flow.AccountContractAdded"
	EventAccountContractUpdated string = "flow.AccountContractUpdated"
	EventAccountContractRemoved string = "flow.AccountContractRemoved"
)

type Event struct {
	// Type is the qualified event type.
	Type string
	// TransactionID is the ID of the transaction this event was emitted from.
	TransactionID Identifier
	// TransactionIndex is the index of the transaction this event was emitted from, within its containing block.
	TransactionIndex int
	// EventIndex is the index of the event within the transaction it was emitted from.
	EventIndex int
	// Value contains the event data.
	Value cadence.Event
	// Bytes representing event data.
	Payload []byte
}

// String returns the string representation of this event.
func (e Event) String() string {
	return fmt.Sprintf("%s: %s", e.Type, e.ID())
}

// ID returns the canonical SHA3-256 hash of this event.
func (e Event) ID() string {
	return defaultEntityHasher.ComputeHash(e.Encode()).Hex()
}

// Encode returns the canonical RLP byte representation of this event.
func (e Event) Encode() []byte {
	temp := struct {
		TransactionID []byte
		EventIndex    uint
	}{
		TransactionID: e.TransactionID[:],
		EventIndex:    uint(e.EventIndex),
	}
	return mustRLPEncode(&temp)
}

// Fingerprint calculates a fingerprint of an event.
func (e *Event) Fingerprint() []byte {

	return mustRLPEncode(struct {
		TxID             []byte
		Index            uint32
		Type             string
		TransactionIndex uint32
		Payload          []byte
	}{
		TxID:             e.TransactionID[:],
		Index:            uint32(e.EventIndex),
		Type:             string(e.Type),
		TransactionIndex: uint32(e.TransactionIndex),
		Payload:          e.Payload[:],
	})
}

// CalculateEventsHash calculates hash of the events, in a way compatible with the events hash
// propagated in Chunk. Providing list of events as emitted by transactions in given chunk,
// hashes should match
func CalculateEventsHash(es []Event) (crypto.Hash, error) {

	hasher := hash.NewSHA3_256()

	for _, event := range es {
		_, err := hasher.Write(event.Fingerprint())
		if err != nil {
			return nil, err
		}
	}

	return hasher.SumHash(), nil
}

// An AccountCreatedEvent is emitted when a transaction creates a new Flow account.
//
// This event contains the following fields:
// - Address: Address
type AccountCreatedEvent Event

// Address returns the address of the newly-created account.
func (evt AccountCreatedEvent) Address() Address {
	return BytesToAddress(evt.Value.Fields[0].(cadence.Address).Bytes())
}

package flow

import (
	"fmt"

	"github.com/onflow/cadence"
)

// List of built-in account event types.
const (
	EventAccountCreated string = "flow.AccountCreated"
	EventAccountUpdated string = "flow.AccountUpdated"
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
}

// String returns the string representation of this event.
func (e Event) String() string {
	return fmt.Sprintf("%s: %s", e.Type, e.ID())
}

// ID returns the canonical SHA3-256 hash of this event.
func (e Event) ID() string {
	return DefaultHasher.ComputeHash(e.Encode()).Hex()
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

// An AccountCreatedEvent is emitted when a transaction creates a new Flow account.
//
// This event contains the following fields:
// - Address: Address
type AccountCreatedEvent Event

// Address returns the address of the newly-created account.
func (evt AccountCreatedEvent) Address() Address {
	return BytesToAddress(evt.Value.Fields[0].(cadence.Address).Bytes())
}

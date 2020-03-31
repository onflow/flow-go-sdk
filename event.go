package flow

import (
	"fmt"

	"github.com/dapperlabs/cadence"
	encoding "github.com/dapperlabs/cadence/encoding/xdr"
	"github.com/dapperlabs/flow-go/crypto"
	"github.com/dapperlabs/flow-go/model/hash"
)

// List of built-in account event types.
const (
	EventAccountCreated string = "flow.AccountCreated"
	EventAccountUpdated string = "flow.AccountUpdated"
)

type Event struct {
	// Type is the qualified event type.
	Type string
	// TxHash is the hash of the transaction this event was emitted from.
	TxHash crypto.Hash
	// Index defines the ordering of events in a transaction. The first event
	// emitted has index 0, the second has index 1, and so on.
	Index uint
	// Payload contains the encoded event data.
	Payload []byte
}

// String returns the string representation of this event.
func (e Event) String() string {
	return fmt.Sprintf("%s: %s", e.Type, e.ID())
}

// ID returns a canonical identifier that is guaranteed to be unique.
func (e Event) ID() string {
	return hash.DefaultHasher.ComputeHash(e.Encode()).Hex()
}

// Encode returns the canonical encoding of the event, containing only the
// fields necessary to uniquely identify it.
func (e Event) Encode() []byte {
	w := wrapEvent(e)
	return DefaultEncoder.MustEncode(w)
}

// Defines only the fields needed to uniquely identify an event.
type eventWrapper struct {
	TxHash []byte
	Index  uint
}

func wrapEvent(e Event) eventWrapper {
	return eventWrapper{
		TxHash: e.TxHash,
		Index:  e.Index,
	}
}

type AccountCreatedEvent interface {
	Address() Address
}

var AccountCreatedEventType cadence.Type = cadence.EventType{
	TypeID: EventAccountCreated,
	Fields: []cadence.Field{
		{
			Identifier: "address",
			Type:       cadence.AddressType{},
		},
	},
}

func newAccountCreatedEventFromValue(v cadence.Value) AccountCreatedEvent {
	eventValue := v.(cadence.Event)
	return accountCreatedEvent{eventValue}
}

type accountCreatedEvent struct {
	cadence.Event
}

func (a accountCreatedEvent) Address() Address {
	return Address(a.Fields[0].(cadence.Address))
}

func DecodeAccountCreatedEvent(b []byte) (AccountCreatedEvent, error) {
	value, err := encoding.Decode(AccountCreatedEventType, b)
	if err != nil {
		return nil, err
	}

	return newAccountCreatedEventFromValue(value), nil
}

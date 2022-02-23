package flow

import (
	"github.com/onflow/cadence"
	"github.com/onflow/cadence/runtime/tests/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculateEventHash(t *testing.T) {
	events := []Event{{
		Type:             EventAccountAdded,
		TransactionID:    HexToID("ae792580417379f04c09f48455be6a19815f0a52feec15f03f8470322ee9dc88"),
		TransactionIndex: 0,
		EventIndex:       0,
		Value: cadence.NewEvent([]cadence.Value{
			cadence.NewInt(1),
			cadence.String("foo"),
		}).WithType(&cadence.EventType{
			Location:            utils.TestLocation,
			QualifiedIdentifier: "FooEvent",
			Fields: []cadence.Field{
				{
					Identifier: "a",
					Type:       cadence.IntType{},
				},
				{
					Identifier: "b",
					Type:       cadence.StringType{},
				},
			},
		}),
		Payload: nil,
	}}

	h, err := CalculateEventsHash(events)
	assert.NoError(t, err)
	assert.Equal(t, "0x7b3d0915bac54b34aba06160ec518ae3be4212a4dc18ad33163601e57f9b1e15", h.Hex())
}

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
	// 0x37623364303931356261633534623334616261303631363065633531386165336265343231326134646331386164333331363336303165353766396231653135
	assert.NoError(t, err)
	assert.Equal(t, "0x0a1e6871892be16088661236ae90c3b3eb7d99e2a367faba55927425d000f43c", h.String())
}

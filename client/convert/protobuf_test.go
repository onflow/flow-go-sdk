package convert_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go-sdk/client/convert"
	"github.com/onflow/flow-go-sdk/test"
)

func TestConvert_Block(t *testing.T) {
	blockA := test.BlockGenerator().New()

	msg := convert.BlockToMessage(*blockA)

	blockB, err := convert.MessageToBlock(msg)

	assert.NoError(t, err)
	assert.Equal(t, *blockA, blockB)
}

func TestConvert_Collection(t *testing.T) {
	colA := test.CollectionGenerator().New()

	msg := convert.CollectionToMessage(*colA)

	colB, err := convert.MessageToCollection(msg)

	assert.NoError(t, err)
	assert.Equal(t, *colA, colB)
}

func TestConvert_Transaction(t *testing.T) {
	txA := test.TransactionGenerator().New()

	msg := convert.TransactionToMessage(*txA)

	txB, err := convert.MessageToTransaction(msg)

	assert.NoError(t, err)
	assert.Equal(t, txA.ID(), txB.ID())
}

func TestConvert_Event(t *testing.T) {
	eventA := test.EventGenerator().New()

	msg, err := convert.EventToMessage(eventA)
	require.NoError(t, err)

	eventB, err := convert.MessageToEvent(msg)
	require.NoError(t, err)

	assert.Equal(t, eventA, eventB)
}

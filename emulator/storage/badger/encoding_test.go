package badger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapperlabs/flow-go/crypto"
	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/emulator/types"
	"github.com/dapperlabs/flow-go-sdk/utils/unittest"
)

func TestEncodeTransaction(t *testing.T) {
	tx := unittest.TransactionFixture()
	data, err := encodeTransaction(tx)
	require.Nil(t, err)

	var decodedTx flow.Transaction
	err = decodeTransaction(&decodedTx, data)
	require.Nil(t, err)
	assert.Equal(t, tx, decodedTx)
}

func TestEncodeBlock(t *testing.T) {
	block := types.Block{
		Number:            1234,
		PreviousBlockHash: unittest.HashFixture(32),
		TransactionHashes: []crypto.Hash{unittest.HashFixture(32)},
	}
	data, err := encodeBlock(block)
	require.Nil(t, err)

	var decodedBlock types.Block
	err = decodeBlock(&decodedBlock, data)
	require.Nil(t, err)

	assert.Equal(t, block.Number, decodedBlock.Number)
	assert.Equal(t, block.PreviousBlockHash, decodedBlock.PreviousBlockHash)
	assert.Equal(t, block.TransactionHashes, decodedBlock.TransactionHashes)
}

func TestEncodeEventList(t *testing.T) {
	eventList := []flow.Event{unittest.EventFixture(func(e *flow.Event) {
		e.Payload = []byte{1, 2, 3, 4}
	})}
	data, err := encodeEvents(eventList)
	require.Nil(t, err)

	var decodedEventList []flow.Event
	err = decodeEvents(&decodedEventList, data)
	require.Nil(t, err)
	assert.Equal(t, eventList, decodedEventList)
}

func TestEncodeChangelist(t *testing.T) {
	var clist changelist
	clist.add(1)

	data, err := encodeChangelist(clist)
	require.NoError(t, err)

	var decodedClist changelist
	err = decodeChangelist(&decodedClist, data)
	require.NoError(t, err)
	assert.Equal(t, clist, decodedClist)
}

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

package convert

import (
	"testing"
	"time"

	"github.com/onflow/cadence"
	"github.com/onflow/cadence/encoding/ccf"
	"github.com/onflow/flow/protobuf/go/flow/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/test"
)

func TestConvert_Account(t *testing.T) {
	accountA := test.AccountGenerator().New()

	msg := AccountToMessage(*accountA)

	accountB, err := MessageToAccount(msg)
	require.NoError(t, err)

	assert.Equal(t, *accountA, accountB)
}

func TestConvert_AccountKey(t *testing.T) {
	keyA := test.AccountKeyGenerator().New()

	msg := AccountKeyToMessage(keyA)

	keyB, err := MessageToAccountKey(msg)
	require.NoError(t, err)

	assert.Equal(t, keyA, keyB)
}

func TestConvert_Block(t *testing.T) {
	blockA := test.BlockGenerator().New()

	msg, err := BlockToMessage(*blockA)
	require.NoError(t, err)

	blockB, err := MessageToBlock(msg)
	require.NoError(t, err)

	assert.Equal(t, *blockA, blockB)

	t.Run("Without timestamp", func(t *testing.T) {
		blockA := test.BlockGenerator().New()

		msg, err := BlockToMessage(*blockA)
		require.NoError(t, err)

		msg.Timestamp = nil

		blockB, err = MessageToBlock(msg)
		require.NoError(t, err)

		assert.Equal(t, time.Time{}, blockB.Timestamp)
	})
}

func TestConvert_BlockHeader(t *testing.T) {
	headerA := test.BlockHeaderGenerator().New()

	msg, err := BlockHeaderToMessage(headerA)
	require.NoError(t, err)

	headerB, err := MessageToBlockHeader(msg)
	require.NoError(t, err)

	assert.Equal(t, headerA, headerB)

	t.Run("Without timestamp", func(t *testing.T) {
		headerA := test.BlockHeaderGenerator().New()

		msg, err := BlockHeaderToMessage(headerA)
		require.NoError(t, err)

		msg.Timestamp = nil

		headerB, err = MessageToBlockHeader(msg)
		require.NoError(t, err)

		assert.Equal(t, time.Time{}, headerB.Timestamp)
	})
}

func TestConvert_CadenceValue(t *testing.T) {
	t.Run("Valid value", func(t *testing.T) {
		valueA := cadence.NewInt(42)

		msg, err := CadenceValueToMessage(valueA)
		require.NoError(t, err)

		valueB, err := MessageToCadenceValue(msg, nil)
		require.NoError(t, err)

		assert.Equal(t, valueA, valueB)
	})

	t.Run("Invalid message", func(t *testing.T) {
		msg := []byte("invalid JSON-CDC bytes")

		value, err := MessageToCadenceValue(msg, nil)
		assert.Error(t, err)
		assert.Nil(t, value)
	})

	t.Run("CCF encoded value", func(t *testing.T) {
		valueA := cadence.NewInt(42)

		msg, err := ccf.Encode(valueA)
		require.NoError(t, err)

		valueB, err := MessageToCadenceValue(msg, nil)
		require.NoError(t, err)

		assert.Equal(t, valueA, valueB)
	})
}

func TestConvert_Collection(t *testing.T) {
	colA := test.CollectionGenerator().New()

	msg := CollectionToMessage(*colA)

	colB, err := MessageToCollection(msg)
	require.NoError(t, err)

	assert.Equal(t, *colA, colB)
}

func TestConvert_CollectionGuarantee(t *testing.T) {
	cgA := test.CollectionGuaranteeGenerator().New()

	msg := CollectionGuaranteeToMessage(*cgA)

	cgB, err := MessageToCollectionGuarantee(msg)
	require.NoError(t, err)

	assert.Equal(t, *cgA, cgB)
}

func TestConvert_BlockSeal(t *testing.T) {
	bsA := test.BlockSealGenerator().New()

	msg := BlockSealToMessage(*bsA)

	bsB, err := MessageToBlockSeal(msg)
	require.NoError(t, err)

	assert.Equal(t, *bsA, bsB)
}

func TestConvert_CollectionGuarantees(t *testing.T) {
	cgs := test.CollectionGuaranteeGenerator()

	cgsA := []*flow.CollectionGuarantee{
		cgs.New(),
		cgs.New(),
		cgs.New(),
	}

	msg := CollectionGuaranteesToMessages(cgsA)

	cgsB, err := MessagesToCollectionGuarantees(msg)
	require.NoError(t, err)

	assert.Equal(t, cgsA, cgsB)
}

func TestConvert_BlockSeals(t *testing.T) {
	bss := test.BlockSealGenerator()

	bssA := []*flow.BlockSeal{
		bss.New(),
		bss.New(),
		bss.New(),
	}

	msg := BlockSealsToMessages(bssA)

	bssB, err := MessagesToBlockSeals(msg)
	require.NoError(t, err)

	assert.Equal(t, bssA, bssB)
}

func TestConvert_Event(t *testing.T) {

	t.Run("JSON-CDC encoded payload", func(t *testing.T) {
		eventA := test.EventGenerator().
			WithEncoding(entities.EventEncodingVersion_JSON_CDC_V0).
			New()
		msg, err := EventToMessage(eventA)
		require.NoError(t, err)

		eventB, err := MessageToEvent(msg, nil)
		require.NoError(t, err)

		// Force evaluation of type ID, which is cached in type.
		// Necessary for equality check below
		_ = eventB.Value.Type().ID()

		assert.Equal(t, eventA, eventB)
	})

	t.Run("CCF encoded payload", func(t *testing.T) {
		eventA := test.EventGenerator().
			WithEncoding(entities.EventEncodingVersion_CCF_V0).
			New()

		msg, err := EventToMessage(eventA)
		require.NoError(t, err)

		// explicitly re-encode the payload using CCF
		msg.Payload, err = ccf.Encode(eventA.Value)
		require.NoError(t, err)

		eventB, err := MessageToEvent(msg, nil)
		require.NoError(t, err)

		// Force evaluation of type ID, which is cached in type.
		// Necessary for equality check below
		_ = eventB.Value.Type().ID()

		assert.Equal(t, eventA, eventB)
	})
}

func TestConvert_Identifier(t *testing.T) {
	idA := test.IdentifierGenerator().New()

	msg := IdentifierToMessage(idA)
	idB := MessageToIdentifier(msg)

	assert.Equal(t, idA, idB)
}

func TestConvert_Identifiers(t *testing.T) {
	ids := test.IdentifierGenerator()

	idsA := []flow.Identifier{
		ids.New(),
		ids.New(),
		ids.New(),
	}

	msg := IdentifiersToMessages(idsA)
	idsB := MessagesToIdentifiers(msg)

	assert.Equal(t, idsA, idsB)
}

func TestConvert_Transaction(t *testing.T) {
	t.Run("Without arguments", func(t *testing.T) {
		txA := test.TransactionGenerator().New()
		txA.Arguments = nil

		msg, err := TransactionToMessage(*txA)
		require.NoError(t, err)

		txB, err := MessageToTransaction(msg)
		require.NoError(t, err)

		assert.Equal(t, txA.ID(), txB.ID())
	})

	t.Run("With arguments", func(t *testing.T) {
		txA := test.TransactionGenerator().New()

		msg, err := TransactionToMessage(*txA)
		require.NoError(t, err)

		txB, err := MessageToTransaction(msg)
		require.NoError(t, err)

		assert.Equal(t, txA.ID(), txB.ID())
	})
}

func TestConvert_TransactionResult(t *testing.T) {
	resultA := test.TransactionResultGenerator().New()

	msg, err := TransactionResultToMessage(resultA)

	resultB, err := MessageToTransactionResult(msg, nil)
	require.NoError(t, err)

	// Force evaluation of type ID, which is cached in type.
	// Necessary for equality check below
	for _, event := range resultB.Events {
		_ = event.Value.Type().ID()
	}

	assert.Equal(t, resultA, resultB)
}

func TestConvert_ExecutionData(t *testing.T) {
	executionDataA := test.ExecutionDataGenerator().New()

	msg, err := BlockExecutionDataToMessage(executionDataA)
	require.NoError(t, err)

	executionDataB, err := MessageToBlockExecutionData(msg)
	require.NoError(t, err)

	assert.Equal(t, executionDataA.BlockID[:], executionDataB.BlockID[:])
	require.NotEmpty(t, executionDataA.ChunkExecutionData)

	// Force evaluation of type ID, which is cached in type.
	// Necessary for equality check below, otherwise the typeID will be empty
	for _, chunk := range executionDataB.ChunkExecutionData {
		for _, event := range chunk.Events {
			_ = event.Value.Type().ID()
		}
	}

	assert.Equal(t, executionDataA, executionDataB)
}

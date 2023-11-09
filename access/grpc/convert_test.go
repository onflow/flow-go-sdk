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

package grpc

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

	msg := accountToMessage(*accountA)

	accountB, err := messageToAccount(msg)
	require.NoError(t, err)

	assert.Equal(t, *accountA, accountB)
}

func TestConvert_AccountKey(t *testing.T) {
	keyA := test.AccountKeyGenerator().New()

	msg := accountKeyToMessage(keyA)

	keyB, err := messageToAccountKey(msg)
	require.NoError(t, err)

	assert.Equal(t, keyA, keyB)
}

func TestConvert_Block(t *testing.T) {
	blockA := test.BlockGenerator().New()

	msg, err := blockToMessage(*blockA)
	require.NoError(t, err)

	blockB, err := messageToBlock(msg)
	require.NoError(t, err)

	assert.Equal(t, *blockA, blockB)

	t.Run("Without timestamp", func(t *testing.T) {
		blockA := test.BlockGenerator().New()

		msg, err := blockToMessage(*blockA)
		require.NoError(t, err)

		msg.Timestamp = nil

		blockB, err = messageToBlock(msg)
		require.NoError(t, err)

		assert.Equal(t, time.Time{}, blockB.Timestamp)
	})
}

func TestConvert_BlockHeader(t *testing.T) {
	headerA := test.BlockHeaderGenerator().New()

	msg, err := blockHeaderToMessage(headerA)
	require.NoError(t, err)

	headerB, err := messageToBlockHeader(msg)
	require.NoError(t, err)

	assert.Equal(t, headerA, headerB)

	t.Run("Without timestamp", func(t *testing.T) {
		headerA := test.BlockHeaderGenerator().New()

		msg, err := blockHeaderToMessage(headerA)
		require.NoError(t, err)

		msg.Timestamp = nil

		headerB, err = messageToBlockHeader(msg)
		require.NoError(t, err)

		assert.Equal(t, time.Time{}, headerB.Timestamp)
	})
}

func TestConvert_CadenceValue(t *testing.T) {
	t.Run("Valid value", func(t *testing.T) {
		valueA := cadence.NewInt(42)

		msg, err := cadenceValueToMessage(valueA)
		require.NoError(t, err)

		valueB, err := messageToCadenceValue(msg, nil)
		require.NoError(t, err)

		assert.Equal(t, valueA, valueB)
	})

	t.Run("Invalid message", func(t *testing.T) {
		msg := []byte("invalid JSON-CDC bytes")

		value, err := messageToCadenceValue(msg, nil)
		assert.Error(t, err)
		assert.Nil(t, value)
	})

	t.Run("CCF encoded value", func(t *testing.T) {
		valueA := cadence.NewInt(42)

		msg, err := ccf.Encode(valueA)
		require.NoError(t, err)

		valueB, err := messageToCadenceValue(msg, nil)
		require.NoError(t, err)

		assert.Equal(t, valueA, valueB)
	})
}

func TestConvert_Collection(t *testing.T) {
	colA := test.CollectionGenerator().New()

	msg := collectionToMessage(*colA)

	colB, err := messageToCollection(msg)
	require.NoError(t, err)

	assert.Equal(t, *colA, colB)
}

func TestConvert_CollectionGuarantee(t *testing.T) {
	cgA := test.CollectionGuaranteeGenerator().New()

	msg := collectionGuaranteeToMessage(*cgA)

	cgB, err := messageToCollectionGuarantee(msg)
	require.NoError(t, err)

	assert.Equal(t, *cgA, cgB)
}

func TestConvert_BlockSeal(t *testing.T) {
	bsA := test.BlockSealGenerator().New()

	msg := blockSealToMessage(*bsA)

	bsB, err := messageToBlockSeal(msg)
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

	msg := collectionGuaranteesToMessages(cgsA)

	cgsB, err := messagesToCollectionGuarantees(msg)
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

	msg := blockSealsToMessages(bssA)

	bssB, err := messagesToBlockSeals(msg)
	require.NoError(t, err)

	assert.Equal(t, bssA, bssB)
}

func TestConvert_Event(t *testing.T) {

	t.Run("JSON-CDC encoded payload", func(t *testing.T) {
		eventA := test.EventGenerator().
			WithEncoding(entities.EventEncodingVersion_JSON_CDC_V0).
			New()
		msg, err := eventToMessage(eventA)
		require.NoError(t, err)

		eventB, err := messageToEvent(msg, nil)
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

		msg, err := eventToMessage(eventA)
		require.NoError(t, err)

		// explicitly re-encode the payload using CCF
		msg.Payload, err = ccf.Encode(eventA.Value)
		require.NoError(t, err)

		eventB, err := messageToEvent(msg, nil)
		require.NoError(t, err)

		// Force evaluation of type ID, which is cached in type.
		// Necessary for equality check below
		_ = eventB.Value.Type().ID()

		assert.Equal(t, eventA, eventB)
	})
}

func TestConvert_Identifier(t *testing.T) {
	idA := test.IdentifierGenerator().New()

	msg := identifierToMessage(idA)
	idB := messageToIdentifier(msg)

	assert.Equal(t, idA, idB)
}

func TestConvert_Identifiers(t *testing.T) {
	ids := test.IdentifierGenerator()

	idsA := []flow.Identifier{
		ids.New(),
		ids.New(),
		ids.New(),
	}

	msg := identifiersToMessages(idsA)
	idsB := messagesToIdentifiers(msg)

	assert.Equal(t, idsA, idsB)
}

func TestConvert_Transaction(t *testing.T) {
	t.Run("Without arguments", func(t *testing.T) {
		txA := test.TransactionGenerator().New()
		txA.Arguments = nil

		msg, err := transactionToMessage(*txA)
		require.NoError(t, err)

		txB, err := messageToTransaction(msg)
		require.NoError(t, err)

		assert.Equal(t, txA.ID(), txB.ID())
	})

	t.Run("With arguments", func(t *testing.T) {
		txA := test.TransactionGenerator().New()

		msg, err := transactionToMessage(*txA)
		require.NoError(t, err)

		txB, err := messageToTransaction(msg)
		require.NoError(t, err)

		assert.Equal(t, txA.ID(), txB.ID())
	})
}

func TestConvert_TransactionResult(t *testing.T) {
	resultA := test.TransactionResultGenerator().New()

	msg, err := transactionResultToMessage(resultA)

	resultB, err := messageToTransactionResult(msg, nil)
	require.NoError(t, err)

	// Force evaluation of type ID, which is cached in type.
	// Necessary for equality check below
	for _, event := range resultB.Events {
		_ = event.Value.Type().ID()
	}

	assert.Equal(t, resultA, resultB)
}

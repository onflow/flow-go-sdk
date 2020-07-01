/*
 * Flow Go SDK
 *
 * Copyright 2019-2020 Dapper Labs, Inc.
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

package convert_test

import (
	"testing"
	"time"

	"github.com/onflow/cadence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client/convert"
	"github.com/onflow/flow-go-sdk/test"
)

func TestConvert_Account(t *testing.T) {
	accountA := test.AccountGenerator().New()

	msg := convert.AccountToMessage(*accountA)

	accountB, err := convert.MessageToAccount(msg)
	require.NoError(t, err)

	assert.Equal(t, *accountA, accountB)
}

func TestConvert_AccountKey(t *testing.T) {
	keyA := test.AccountKeyGenerator().New()

	msg := convert.AccountKeyToMessage(keyA)

	keyB, err := convert.MessageToAccountKey(msg)
	require.NoError(t, err)

	assert.Equal(t, keyA, keyB)
}

func TestConvert_Block(t *testing.T) {
	blockA := test.BlockGenerator().New()

	msg, err := convert.BlockToMessage(*blockA)
	require.NoError(t, err)

	blockB, err := convert.MessageToBlock(msg)
	require.NoError(t, err)

	assert.Equal(t, *blockA, blockB)

	t.Run("Without timestamp", func(t *testing.T) {
		blockA := test.BlockGenerator().New()

		msg, err := convert.BlockToMessage(*blockA)
		require.NoError(t, err)

		msg.Timestamp = nil

		blockB, err = convert.MessageToBlock(msg)
		require.NoError(t, err)

		assert.Equal(t, time.Time{}, blockB.Timestamp)
	})
}

func TestConvert_BlockHeader(t *testing.T) {
	headerA := test.BlockHeaderGenerator().New()

	msg, err := convert.BlockHeaderToMessage(headerA)
	require.NoError(t, err)

	headerB, err := convert.MessageToBlockHeader(msg)
	require.NoError(t, err)

	assert.Equal(t, headerA, headerB)

	t.Run("Without timestamp", func(t *testing.T) {
		headerA := test.BlockHeaderGenerator().New()

		msg, err := convert.BlockHeaderToMessage(headerA)
		require.NoError(t, err)

		msg.Timestamp = nil

		headerB, err = convert.MessageToBlockHeader(msg)
		require.NoError(t, err)

		assert.Equal(t, time.Time{}, headerB.Timestamp)
	})
}

func TestConvert_CadenceValue(t *testing.T) {
	t.Run("Valid value", func(t *testing.T) {
		valueA := cadence.NewInt(42)

		msg, err := convert.CadenceValueToMessage(valueA)
		require.NoError(t, err)

		valueB, err := convert.MessageToCadenceValue(msg)
		require.NoError(t, err)

		assert.Equal(t, valueA, valueB)
	})

	t.Run("Invalid message", func(t *testing.T) {
		msg := []byte("invalid JSON-CDC bytes")

		value, err := convert.MessageToCadenceValue(msg)
		assert.Error(t, err)
		assert.Nil(t, value)
	})
}

func TestConvert_Collection(t *testing.T) {
	colA := test.CollectionGenerator().New()

	msg := convert.CollectionToMessage(*colA)

	colB, err := convert.MessageToCollection(msg)
	require.NoError(t, err)

	assert.Equal(t, *colA, colB)
}

func TestConvert_CollectionGuarantee(t *testing.T) {
	cgA := test.CollectionGuaranteeGenerator().New()

	msg := convert.CollectionGuaranteeToMessage(*cgA)

	cgB, err := convert.MessageToCollectionGuarantee(msg)
	require.NoError(t, err)

	assert.Equal(t, *cgA, cgB)
}

func TestConvert_CollectionGuarantees(t *testing.T) {
	cgs := test.CollectionGuaranteeGenerator()

	cgsA := []*flow.CollectionGuarantee{
		cgs.New(),
		cgs.New(),
		cgs.New(),
	}

	msg := convert.CollectionGuaranteesToMessages(cgsA)

	cgsB, err := convert.MessagesToCollectionGuarantees(msg)
	require.NoError(t, err)

	assert.Equal(t, cgsA, cgsB)
}

func TestConvert_Event(t *testing.T) {
	eventA := test.EventGenerator().New()

	msg, err := convert.EventToMessage(eventA)
	require.NoError(t, err)

	eventB, err := convert.MessageToEvent(msg)
	require.NoError(t, err)

	assert.Equal(t, eventA, eventB)
}

func TestConvert_Identifier(t *testing.T) {
	idA := test.IdentifierGenerator().New()

	msg := convert.IdentifierToMessage(idA)
	idB := convert.MessageToIdentifier(msg)

	assert.Equal(t, idA, idB)
}

func TestConvert_Identifiers(t *testing.T) {
	ids := test.IdentifierGenerator()

	idsA := []flow.Identifier{
		ids.New(),
		ids.New(),
		ids.New(),
	}

	msg := convert.IdentifiersToMessages(idsA)
	idsB := convert.MessagesToIdentifiers(msg)

	assert.Equal(t, idsA, idsB)
}

func TestConvert_Transaction(t *testing.T) {
	t.Run("Without arguments", func(t *testing.T) {
		txA := test.TransactionGenerator().New()
		txA.Arguments = nil

		msg, err := convert.TransactionToMessage(*txA)
		require.NoError(t, err)

		txB, err := convert.MessageToTransaction(msg)
		require.NoError(t, err)

		assert.Equal(t, txA.ID(), txB.ID())
	})

	t.Run("With arguments", func(t *testing.T) {
		txA := test.TransactionGenerator().New()

		msg, err := convert.TransactionToMessage(*txA)
		require.NoError(t, err)

		txB, err := convert.MessageToTransaction(msg)
		require.NoError(t, err)

		assert.Equal(t, txA.ID(), txB.ID())
	})
}

func TestConvert_TransactionResult(t *testing.T) {
	resultA := test.TransactionResultGenerator().New()

	msg, err := convert.TransactionResultToMessage(resultA)

	resultB, err := convert.MessageToTransactionResult(msg)
	require.NoError(t, err)

	assert.Equal(t, resultA, resultB)
}

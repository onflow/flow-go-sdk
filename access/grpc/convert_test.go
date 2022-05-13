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

package grpc_test

import (
	"testing"
	"time"

	"github.com/onflow/flow-go-sdk/access/grpc"

	"github.com/onflow/cadence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/test"
)

func TestConvert_Account(t *testing.T) {
	accountA := test.AccountGenerator().New()

	msg := grpc.AccountToMessage(*accountA)

	accountB, err := grpc.MessageToAccount(msg)
	require.NoError(t, err)

	assert.Equal(t, *accountA, accountB)
}

func TestConvert_AccountKey(t *testing.T) {
	keyA := test.AccountKeyGenerator().New()

	msg := grpc.AccountKeyToMessage(keyA)

	keyB, err := grpc.MessageToAccountKey(msg)
	require.NoError(t, err)

	assert.Equal(t, keyA, keyB)
}

func TestConvert_Block(t *testing.T) {
	blockA := test.BlockGenerator().New()

	msg, err := grpc.BlockToMessage(*blockA)
	require.NoError(t, err)

	blockB, err := grpc.MessageToBlock(msg)
	require.NoError(t, err)

	assert.Equal(t, *blockA, blockB)

	t.Run("Without timestamp", func(t *testing.T) {
		blockA := test.BlockGenerator().New()

		msg, err := grpc.BlockToMessage(*blockA)
		require.NoError(t, err)

		msg.Timestamp = nil

		blockB, err = grpc.MessageToBlock(msg)
		require.NoError(t, err)

		assert.Equal(t, time.Time{}, blockB.Timestamp)
	})
}

func TestConvert_BlockHeader(t *testing.T) {
	headerA := test.BlockHeaderGenerator().New()

	msg, err := grpc.BlockHeaderToMessage(headerA)
	require.NoError(t, err)

	headerB, err := grpc.MessageToBlockHeader(msg)
	require.NoError(t, err)

	assert.Equal(t, headerA, headerB)

	t.Run("Without timestamp", func(t *testing.T) {
		headerA := test.BlockHeaderGenerator().New()

		msg, err := grpc.BlockHeaderToMessage(headerA)
		require.NoError(t, err)

		msg.Timestamp = nil

		headerB, err = grpc.MessageToBlockHeader(msg)
		require.NoError(t, err)

		assert.Equal(t, time.Time{}, headerB.Timestamp)
	})
}

func TestConvert_CadenceValue(t *testing.T) {
	t.Run("Valid value", func(t *testing.T) {
		valueA := cadence.NewInt(42)

		msg, err := grpc.CadenceValueToMessage(valueA)
		require.NoError(t, err)

		valueB, err := grpc.MessageToCadenceValue(msg)
		require.NoError(t, err)

		assert.Equal(t, valueA, valueB)
	})

	t.Run("Invalid message", func(t *testing.T) {
		msg := []byte("invalid JSON-CDC bytes")

		value, err := grpc.MessageToCadenceValue(msg)
		assert.Error(t, err)
		assert.Nil(t, value)
	})
}

func TestConvert_Collection(t *testing.T) {
	colA := test.CollectionGenerator().New()

	msg := grpc.CollectionToMessage(*colA)

	colB, err := grpc.MessageToCollection(msg)
	require.NoError(t, err)

	assert.Equal(t, *colA, colB)
}

func TestConvert_CollectionGuarantee(t *testing.T) {
	cgA := test.CollectionGuaranteeGenerator().New()

	msg := grpc.CollectionGuaranteeToMessage(*cgA)

	cgB, err := grpc.MessageToCollectionGuarantee(msg)
	require.NoError(t, err)

	assert.Equal(t, *cgA, cgB)
}

func TestConvert_BlockSeal(t *testing.T) {
	bsA := test.BlockSealGenerator().New()

	msg := grpc.BlockSealToMessage(*bsA)

	bsB, err := grpc.MessageToBlockSeal(msg)
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

	msg := grpc.CollectionGuaranteesToMessages(cgsA)

	cgsB, err := grpc.MessagesToCollectionGuarantees(msg)
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

	msg := grpc.BlockSealsToMessages(bssA)

	bssB, err := grpc.MessagesToBlockSeals(msg)
	require.NoError(t, err)

	assert.Equal(t, bssA, bssB)
}

func TestConvert_Event(t *testing.T) {
	eventA := test.EventGenerator().New()

	msg, err := grpc.EventToMessage(eventA)
	require.NoError(t, err)

	eventB, err := grpc.MessageToEvent(msg)
	require.NoError(t, err)

	assert.Equal(t, eventA, eventB)
}

func TestConvert_Identifier(t *testing.T) {
	idA := test.IdentifierGenerator().New()

	msg := grpc.IdentifierToMessage(idA)
	idB := grpc.MessageToIdentifier(msg)

	assert.Equal(t, idA, idB)
}

func TestConvert_Identifiers(t *testing.T) {
	ids := test.IdentifierGenerator()

	idsA := []flow.Identifier{
		ids.New(),
		ids.New(),
		ids.New(),
	}

	msg := grpc.IdentifiersToMessages(idsA)
	idsB := grpc.MessagesToIdentifiers(msg)

	assert.Equal(t, idsA, idsB)
}

func TestConvert_Transaction(t *testing.T) {
	t.Run("Without arguments", func(t *testing.T) {
		txA := test.TransactionGenerator().New()
		txA.Arguments = nil

		msg, err := grpc.TransactionToMessage(*txA)
		require.NoError(t, err)

		txB, err := grpc.MessageToTransaction(msg)
		require.NoError(t, err)

		assert.Equal(t, txA.ID(), txB.ID())
	})

	t.Run("With arguments", func(t *testing.T) {
		txA := test.TransactionGenerator().New()

		msg, err := grpc.TransactionToMessage(*txA)
		require.NoError(t, err)

		txB, err := grpc.MessageToTransaction(msg)
		require.NoError(t, err)

		assert.Equal(t, txA.ID(), txB.ID())
	})
}

func TestConvert_TransactionResult(t *testing.T) {
	resultA := test.TransactionResultGenerator().New()

	msg, err := grpc.TransactionResultToMessage(resultA)

	resultB, err := grpc.MessageToTransactionResult(msg)
	require.NoError(t, err)

	assert.Equal(t, resultA, resultB)
}

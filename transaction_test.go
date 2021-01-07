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

package flow_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/test"
)

func ExampleTransaction() {
	addresses := test.AddressGenerator()

	// Mock user accounts

	adrianLaptopKey := &flow.AccountKey{
		Index:          3,
		SequenceNumber: 42,
	}

	adrianPhoneKey := &flow.AccountKey{Index: 2}
	addressA := addresses.New()

	adrian := flow.Account{
		Address: addressA,
		Keys: []*flow.AccountKey{
			adrianLaptopKey,
			adrianPhoneKey,
		},
	}

	blaineHardwareKey := &flow.AccountKey{Index: 7}
	addressB := addresses.New()

	blaine := flow.Account{
		Address: addressB,
		Keys: []*flow.AccountKey{
			blaineHardwareKey,
		},
	}

	// Transaction preparation

	tx := flow.NewTransaction().
		SetScript([]byte(`transaction { execute { log("Hello, World!") } }`)).
		SetReferenceBlockID(flow.Identifier{0x01, 0x02}).
		SetGasLimit(42).
		SetProposalKey(adrian.Address, adrianLaptopKey.Index, adrianLaptopKey.SequenceNumber).
		SetPayer(blaine.Address).
		AddAuthorizer(adrian.Address)

	fmt.Printf("Transaction ID (before signing): %s\n\n", tx.ID())

	// Signing

	err := tx.SignPayload(adrian.Address, adrianLaptopKey.Index, test.MockSigner([]byte{1}))
	if err != nil {
		panic(err)
	}

	err = tx.SignPayload(adrian.Address, adrianPhoneKey.Index, test.MockSigner([]byte{2}))
	if err != nil {
		panic(err)
	}

	err = tx.SignEnvelope(blaine.Address, blaineHardwareKey.Index, test.MockSigner([]byte{3}))
	if err != nil {
		panic(err)
	}

	fmt.Println("Payload signatures:")
	for _, sig := range tx.PayloadSignatures {
		fmt.Printf(
			"Address: %s, Key Index: %d, Signature: %x\n",
			sig.Address,
			sig.KeyIndex,
			sig.Signature,
		)
	}
	fmt.Println()

	fmt.Println("Envelope signatures:")
	for _, sig := range tx.EnvelopeSignatures {
		fmt.Printf(
			"Address: %s, Key Index: %d, Signature: %x\n",
			sig.Address,
			sig.KeyIndex,
			sig.Signature,
		)
	}
	fmt.Println()

	fmt.Printf("Transaction ID (after signing): %s\n", tx.ID())

	// Output:
	// Transaction ID (before signing): 8c362dd8b7553d48284cecc94d2ab545d513b29f930555632390fff5ca9772ee
	//
	// Payload signatures:
	// Address: f8d6e0586b0a20c7, Key Index: 2, Signature: 02
	// Address: f8d6e0586b0a20c7, Key Index: 3, Signature: 01
	//
	// Envelope signatures:
	// Address: ee82856bf20e2aa6, Key Index: 7, Signature: 03
	//
	// Transaction ID (after signing): d1a2c58aebfce1050a32edf3568ec3b69cb8637ae090b5f7444ca6b2a8de8f8b
}

func TestTransaction_SetScript(t *testing.T) {
	tx := flow.NewTransaction().
		SetScript(test.GreetingScript)

	assert.Equal(t, test.GreetingScript, tx.Script)
}

func TestTransaction_AddArgument(t *testing.T) {
	t.Run("No arguments", func(t *testing.T) {
		tx := flow.NewTransaction()
		assert.Empty(t, tx.Arguments)
	})

	t.Run("Single argument", func(t *testing.T) {
		expectedArg := cadence.NewString("foo")

		tx := flow.NewTransaction()

		err := tx.AddArgument(expectedArg)
		require.NoError(t, err)

		actualArg, err := tx.Argument(0)
		assert.NoError(t, err)

		assert.Equal(t, expectedArg, actualArg)
	})

	t.Run("Multiple arguments", func(t *testing.T) {
		expectedArgA := cadence.NewString("foo")
		expectedArgB := cadence.NewInt(42)

		tx := flow.NewTransaction()

		err := tx.AddArgument(expectedArgA)
		require.NoError(t, err)
		err = tx.AddArgument(expectedArgB)
		require.NoError(t, err)

		actualArgA, err := tx.Argument(0)
		assert.NoError(t, err)

		actualArgB, err := tx.Argument(1)
		assert.NoError(t, err)

		assert.Equal(t, expectedArgA, actualArgA)
		assert.Equal(t, expectedArgB, actualArgB)
	})
}

func TestTransaction_AddRawArgument(t *testing.T) {
	t.Run("Single argument", func(t *testing.T) {
		expectedArg := cadence.NewString("foo")

		encodedArg, err := jsoncdc.Encode(expectedArg)
		require.NoError(t, err)

		tx := flow.NewTransaction().
			AddRawArgument(encodedArg)

		actualArg, err := tx.Argument(0)
		assert.NoError(t, err)

		assert.Equal(t, expectedArg, actualArg)
	})

	t.Run("Multiple arguments", func(t *testing.T) {
		expectedArgA := cadence.NewString("foo")
		expectedArgB := cadence.NewInt(42)

		encodedArgA, err := jsoncdc.Encode(expectedArgA)
		require.NoError(t, err)

		encodedArgB, err := jsoncdc.Encode(expectedArgB)
		require.NoError(t, err)

		tx := flow.NewTransaction().
			AddRawArgument(encodedArgA).
			AddRawArgument(encodedArgB)

		actualArgA, err := tx.Argument(0)
		assert.NoError(t, err)

		actualArgB, err := tx.Argument(1)
		assert.NoError(t, err)

		assert.Equal(t, expectedArgA, actualArgA)
		assert.Equal(t, expectedArgB, actualArgB)
	})

	t.Run("Invalid argument", func(t *testing.T) {
		tx := flow.NewTransaction().
			AddRawArgument([]byte{1, 2, 3})

		actualArg, err := tx.Argument(0)
		assert.Nil(t, actualArg)
		assert.Error(t, err)
	})
}

func TestTransaction_SetReferenceBlockID(t *testing.T) {
	blockID := test.IdentifierGenerator().New()

	tx := flow.NewTransaction().
		SetReferenceBlockID(blockID)

	assert.Equal(t, blockID, tx.ReferenceBlockID)
}

func TestTransaction_SetGasLimit(t *testing.T) {
	var gasLimit uint64 = 42

	tx := flow.NewTransaction().
		SetGasLimit(gasLimit)

	assert.Equal(t, gasLimit, tx.GasLimit)
}

func TestTransaction_SetProposalKey(t *testing.T) {
	address := flow.ServiceAddress(flow.Mainnet)
	keyIndex := 7
	var sequenceNumber uint64 = 42

	tx := flow.NewTransaction().
		SetProposalKey(address, keyIndex, sequenceNumber)

	assert.Equal(t, address, tx.ProposalKey.Address)
	assert.Equal(t, keyIndex, tx.ProposalKey.KeyIndex)
	assert.Equal(t, sequenceNumber, tx.ProposalKey.SequenceNumber)
}

func TestTransaction_SetPayer(t *testing.T) {
	address := flow.ServiceAddress(flow.Mainnet)

	tx := flow.NewTransaction().
		SetPayer(address)

	assert.Equal(t, address, tx.Payer)
}

func TestTransaction_AddAuthorizer(t *testing.T) {
	addresses := test.AddressGenerator()

	addressA := addresses.New()
	addressB := addresses.New()

	tx := flow.NewTransaction().
		AddAuthorizer(addressA)

	require.Len(t, tx.Authorizers, 1)
	assert.Equal(t, addressA, tx.Authorizers[0])

	tx.AddAuthorizer(addressB)

	require.Len(t, tx.Authorizers, 2)
	assert.Equal(t, addressA, tx.Authorizers[0])
	assert.Equal(t, addressB, tx.Authorizers[1])
	assert.NotEqual(t, addressB, addressA)
}

func TestTransaction_AddPayloadSignature(t *testing.T) {
	addresses := test.AddressGenerator()

	t.Run("Invalid signer", func(t *testing.T) {
		tx := flow.NewTransaction()

		address := addresses.New()

		tx.AddPayloadSignature(address, 7, []byte{42})

		require.Len(t, tx.PayloadSignatures, 1)

		// signer cannot be found, so index is -1
		assert.Equal(t, -1, tx.PayloadSignatures[0].SignerIndex)
	})

	t.Run("Valid signers", func(t *testing.T) {
		addressA := addresses.New()
		addressB := addresses.New()

		keyIndex := 7
		sig := []byte{42}

		tx := flow.NewTransaction().
			AddAuthorizer(addressA).
			AddAuthorizer(addressB)

		// add signatures in reverse order of declaration
		tx.AddPayloadSignature(addressB, keyIndex, sig)
		tx.AddPayloadSignature(addressA, keyIndex, sig)

		require.Len(t, tx.PayloadSignatures, 2)

		assert.Equal(t, 0, tx.PayloadSignatures[0].SignerIndex)
		assert.Equal(t, addressA, tx.PayloadSignatures[0].Address)
		assert.Equal(t, keyIndex, tx.PayloadSignatures[0].KeyIndex)
		assert.Equal(t, sig, tx.PayloadSignatures[0].Signature)

		assert.Equal(t, 1, tx.PayloadSignatures[1].SignerIndex)
		assert.Equal(t, addressB, tx.PayloadSignatures[1].Address)
		assert.Equal(t, keyIndex, tx.PayloadSignatures[1].KeyIndex)
		assert.Equal(t, sig, tx.PayloadSignatures[1].Signature)
	})

	t.Run("Duplicate signers", func(t *testing.T) {
		addressA := addresses.New()
		addressB := addresses.New()

		keyIndex := 7
		sig := []byte{42}

		tx := flow.NewTransaction().
			SetProposalKey(addressA, keyIndex, 42).
			AddAuthorizer(addressB).
			AddAuthorizer(addressA)

		// add signatures in reverse order of declaration
		tx.AddPayloadSignature(addressB, keyIndex, sig)
		tx.AddPayloadSignature(addressA, keyIndex, sig)

		require.Len(t, tx.PayloadSignatures, 2)

		assert.Equal(t, 0, tx.PayloadSignatures[0].SignerIndex)
		assert.Equal(t, addressA, tx.PayloadSignatures[0].Address)
		assert.Equal(t, keyIndex, tx.PayloadSignatures[0].KeyIndex)
		assert.Equal(t, sig, tx.PayloadSignatures[0].Signature)

		assert.Equal(t, 1, tx.PayloadSignatures[1].SignerIndex)
		assert.Equal(t, addressB, tx.PayloadSignatures[1].Address)
		assert.Equal(t, keyIndex, tx.PayloadSignatures[1].KeyIndex)
		assert.Equal(t, sig, tx.PayloadSignatures[1].Signature)
	})

	t.Run("Multiple signatures", func(t *testing.T) {
		address := addresses.New()

		keyIndexA := 7
		sigA := []byte{42}

		keyIndexB := 8
		sigB := []byte{43}

		tx := flow.NewTransaction().
			AddAuthorizer(address)

		// add signatures in descending order by key index
		tx.AddPayloadSignature(address, keyIndexB, sigB)
		tx.AddPayloadSignature(address, keyIndexA, sigA)

		require.Len(t, tx.PayloadSignatures, 2)

		// signatures should be sorted in ascending order by key ID
		assert.Equal(t, 0, tx.PayloadSignatures[0].SignerIndex)
		assert.Equal(t, address, tx.PayloadSignatures[0].Address)
		assert.Equal(t, keyIndexA, tx.PayloadSignatures[0].KeyIndex)
		assert.Equal(t, sigA, tx.PayloadSignatures[0].Signature)

		assert.Equal(t, 0, tx.PayloadSignatures[1].SignerIndex)
		assert.Equal(t, address, tx.PayloadSignatures[1].Address)
		assert.Equal(t, keyIndexB, tx.PayloadSignatures[1].KeyIndex)
		assert.Equal(t, sigB, tx.PayloadSignatures[1].Signature)
	})
}

func TestTransaction_AddEnvelopeSignature(t *testing.T) {
	addresses := test.AddressGenerator()

	t.Run("Invalid signer", func(t *testing.T) {
		tx := flow.NewTransaction()

		address := addresses.New()

		tx.AddEnvelopeSignature(address, 7, []byte{42})

		require.Len(t, tx.EnvelopeSignatures, 1)

		// signer cannot be found, so index is -1
		assert.Equal(t, -1, tx.EnvelopeSignatures[0].SignerIndex)
	})

	t.Run("Valid signer", func(t *testing.T) {
		address := addresses.New()

		keyIndex := 7
		sig := []byte{42}

		tx := flow.NewTransaction().
			SetPayer(address)

		tx.AddEnvelopeSignature(address, keyIndex, sig)

		require.Len(t, tx.EnvelopeSignatures, 1)

		assert.Equal(t, 0, tx.EnvelopeSignatures[0].SignerIndex)
		assert.Equal(t, address, tx.EnvelopeSignatures[0].Address)
		assert.Equal(t, keyIndex, tx.EnvelopeSignatures[0].KeyIndex)
		assert.Equal(t, sig, tx.EnvelopeSignatures[0].Signature)
	})

	t.Run("Multiple signatures", func(t *testing.T) {
		address := addresses.New()

		keyIndexA := 7
		sigA := []byte{42}

		keyIndexB := 8
		sigB := []byte{43}

		tx := flow.NewTransaction().AddAuthorizer(address)

		// add signatures in descending order by key ID
		tx.AddEnvelopeSignature(address, keyIndexB, sigB)
		tx.AddEnvelopeSignature(address, keyIndexA, sigA)

		require.Len(t, tx.EnvelopeSignatures, 2)

		// signatures should be sorted in ascending order by key ID
		assert.Equal(t, 0, tx.EnvelopeSignatures[0].SignerIndex)
		assert.Equal(t, address, tx.EnvelopeSignatures[0].Address)
		assert.Equal(t, keyIndexA, tx.EnvelopeSignatures[0].KeyIndex)
		assert.Equal(t, sigA, tx.EnvelopeSignatures[0].Signature)

		assert.Equal(t, 0, tx.EnvelopeSignatures[1].SignerIndex)
		assert.Equal(t, address, tx.EnvelopeSignatures[1].Address)
		assert.Equal(t, keyIndexB, tx.EnvelopeSignatures[1].KeyIndex)
		assert.Equal(t, sigB, tx.EnvelopeSignatures[1].Signature)
	})
}

func TestTransaction_AbleToReconstructTransaction(t *testing.T) {
	addresses := test.AddressGenerator()
	addressOne := addresses.New()
	addressTwo := addresses.New()

	keyIndex := 7
	sig := []byte{42}

	tx := flow.NewTransaction().
		AddAuthorizer(addressOne).
		SetProposalKey(addressTwo, 0, 0).
		SetPayer(addressOne)

	tx.
		AddPayloadSignature(addressTwo, 0, sig).
		AddEnvelopeSignature(addressOne, keyIndex, sig)

	t.Run("Valid signer", func(t *testing.T) {

		require.Len(t, tx.EnvelopeSignatures, 1)
		require.Len(t, tx.PayloadSignatures, 1)

		assert.Equal(t, 0, tx.PayloadSignatures[0].SignerIndex)
		assert.Equal(t, 1, tx.EnvelopeSignatures[0].SignerIndex)
		assert.Equal(t, addressOne, tx.EnvelopeSignatures[0].Address)
		assert.Equal(t, keyIndex, tx.EnvelopeSignatures[0].KeyIndex)
		assert.Equal(t, sig, tx.EnvelopeSignatures[0].Signature)
	})

	t.Run("Valid reconstructed transaction", func(t *testing.T) {

		newTx := flow.NewTransaction().
			AddPayloadSignature(addressTwo, 0, sig).
			AddEnvelopeSignature(addressOne, keyIndex, sig)

		assert.Equal(t, -1, newTx.PayloadSignatures[0].SignerIndex)
		assert.Equal(t, -1, newTx.EnvelopeSignatures[0].SignerIndex)

		newTx.
			AddAuthorizer(addressOne).
			SetProposalKey(addressTwo, 0, 0).
			SetPayer(addressOne)

		assert.Equal(t, 0, newTx.PayloadSignatures[0].SignerIndex)
		assert.Equal(t, 1, newTx.EnvelopeSignatures[0].SignerIndex)
		assert.Equal(t, addressOne, newTx.EnvelopeSignatures[0].Address)
		assert.Equal(t, keyIndex, newTx.EnvelopeSignatures[0].KeyIndex)
		assert.Equal(t, sig, newTx.EnvelopeSignatures[0].Signature)

		assert.Equal(t, tx.ID(), newTx.ID())
	})
}

var sig, _ = hex.DecodeString("f7225388c1d69d57e6251c9fda50cbbf9e05131e5adb81e5aa0422402f048162")

func baseTx() *flow.Transaction {
	return flow.NewTransaction().
		SetScript([]byte(`transaction { execute { log("Hello, World!") } }`)).
		SetReferenceBlockID(flow.HexToID("f0e4c2f76c58916ec258f246851bea091d14d4247a2fc3e18694461b1816e13b")).
		SetGasLimit(42).
		SetProposalKey(flow.HexToAddress("01"), 4, 10).
		SetPayer(flow.HexToAddress("01")).
		AddAuthorizer(flow.HexToAddress("01")).
		AddPayloadSignature(flow.HexToAddress("01"), 4, sig)
}

func copyTxPayload(tx *flow.Transaction) *flow.Transaction {
	return &flow.Transaction{
		Script:           tx.Script,
		Arguments:        tx.Arguments,
		ReferenceBlockID: tx.ReferenceBlockID,
		GasLimit:         tx.GasLimit,
		ProposalKey:      tx.ProposalKey,
		Payer:            tx.Payer,
		Authorizers:      tx.Authorizers,
	}
}

func copyTxEnvelope(tx *flow.Transaction) *flow.Transaction {
	newTx := copyTxPayload(tx)
	newTx.PayloadSignatures = tx.PayloadSignatures
	return newTx
}

// NOTE: The following tests have identical cases in the
// JavaScript SDK to ensure parity between implementations:
// https://github.com/onflow/flow-js-sdk/blob/master/packages/encode/src/encode.test.js
func TestTransaction_RLPMessages(t *testing.T) {
	var tests = []struct {
		name     string
		tx       *flow.Transaction
		payload  string
		envelope string
	}{
		{
			name:     "Complete transaction",
			tx:       baseTx(),
			payload:  "f872b07472616e73616374696f6e207b2065786563757465207b206c6f67282248656c6c6f2c20576f726c64212229207d207dc0a0f0e4c2f76c58916ec258f246851bea091d14d4247a2fc3e18694461b1816e13b2a880000000000000001040a880000000000000001c9880000000000000001",
			envelope: "f899f872b07472616e73616374696f6e207b2065786563757465207b206c6f67282248656c6c6f2c20576f726c64212229207d207dc0a0f0e4c2f76c58916ec258f246851bea091d14d4247a2fc3e18694461b1816e13b2a880000000000000001040a880000000000000001c9880000000000000001e4e38004a0f7225388c1d69d57e6251c9fda50cbbf9e05131e5adb81e5aa0422402f048162",
		},
		{
			name:     "Complete transaction with envelope sig",
			tx:       baseTx().AddEnvelopeSignature(flow.HexToAddress("01"), 4, sig),
			payload:  "f872b07472616e73616374696f6e207b2065786563757465207b206c6f67282248656c6c6f2c20576f726c64212229207d207dc0a0f0e4c2f76c58916ec258f246851bea091d14d4247a2fc3e18694461b1816e13b2a880000000000000001040a880000000000000001c9880000000000000001",
			envelope: "f899f872b07472616e73616374696f6e207b2065786563757465207b206c6f67282248656c6c6f2c20576f726c64212229207d207dc0a0f0e4c2f76c58916ec258f246851bea091d14d4247a2fc3e18694461b1816e13b2a880000000000000001040a880000000000000001c9880000000000000001e4e38004a0f7225388c1d69d57e6251c9fda50cbbf9e05131e5adb81e5aa0422402f048162",
		},
		{
			name:     "Empty script",
			tx:       baseTx().SetScript(nil),
			payload:  "f84280c0a0f0e4c2f76c58916ec258f246851bea091d14d4247a2fc3e18694461b1816e13b2a880000000000000001040a880000000000000001c9880000000000000001",
			envelope: "f869f84280c0a0f0e4c2f76c58916ec258f246851bea091d14d4247a2fc3e18694461b1816e13b2a880000000000000001040a880000000000000001c9880000000000000001e4e38004a0f7225388c1d69d57e6251c9fda50cbbf9e05131e5adb81e5aa0422402f048162",
		},
		{
			name:     "Empty reference block",
			tx:       baseTx().SetReferenceBlockID(flow.EmptyID),
			payload:  "f872b07472616e73616374696f6e207b2065786563757465207b206c6f67282248656c6c6f2c20576f726c64212229207d207dc0a000000000000000000000000000000000000000000000000000000000000000002a880000000000000001040a880000000000000001c9880000000000000001",
			envelope: "f899f872b07472616e73616374696f6e207b2065786563757465207b206c6f67282248656c6c6f2c20576f726c64212229207d207dc0a000000000000000000000000000000000000000000000000000000000000000002a880000000000000001040a880000000000000001c9880000000000000001e4e38004a0f7225388c1d69d57e6251c9fda50cbbf9e05131e5adb81e5aa0422402f048162",
		},
		{
			name:     "Zero gas limit",
			tx:       baseTx().SetGasLimit(0),
			payload:  "f872b07472616e73616374696f6e207b2065786563757465207b206c6f67282248656c6c6f2c20576f726c64212229207d207dc0a0f0e4c2f76c58916ec258f246851bea091d14d4247a2fc3e18694461b1816e13b80880000000000000001040a880000000000000001c9880000000000000001",
			envelope: "f899f872b07472616e73616374696f6e207b2065786563757465207b206c6f67282248656c6c6f2c20576f726c64212229207d207dc0a0f0e4c2f76c58916ec258f246851bea091d14d4247a2fc3e18694461b1816e13b80880000000000000001040a880000000000000001c9880000000000000001e4e38004a0f7225388c1d69d57e6251c9fda50cbbf9e05131e5adb81e5aa0422402f048162",
		},
		{
			name:     "Empty proposal key ID",
			tx:       baseTx().SetProposalKey(flow.HexToAddress("01"), 0, 10),
			payload:  "f872b07472616e73616374696f6e207b2065786563757465207b206c6f67282248656c6c6f2c20576f726c64212229207d207dc0a0f0e4c2f76c58916ec258f246851bea091d14d4247a2fc3e18694461b1816e13b2a880000000000000001800a880000000000000001c9880000000000000001",
			envelope: "f899f872b07472616e73616374696f6e207b2065786563757465207b206c6f67282248656c6c6f2c20576f726c64212229207d207dc0a0f0e4c2f76c58916ec258f246851bea091d14d4247a2fc3e18694461b1816e13b2a880000000000000001800a880000000000000001c9880000000000000001e4e38004a0f7225388c1d69d57e6251c9fda50cbbf9e05131e5adb81e5aa0422402f048162",
		},
		{
			name:     "Empty sequence number",
			tx:       baseTx().SetProposalKey(flow.HexToAddress("01"), 4, 0),
			payload:  "f872b07472616e73616374696f6e207b2065786563757465207b206c6f67282248656c6c6f2c20576f726c64212229207d207dc0a0f0e4c2f76c58916ec258f246851bea091d14d4247a2fc3e18694461b1816e13b2a8800000000000000010480880000000000000001c9880000000000000001",
			envelope: "f899f872b07472616e73616374696f6e207b2065786563757465207b206c6f67282248656c6c6f2c20576f726c64212229207d207dc0a0f0e4c2f76c58916ec258f246851bea091d14d4247a2fc3e18694461b1816e13b2a8800000000000000010480880000000000000001c9880000000000000001e4e38004a0f7225388c1d69d57e6251c9fda50cbbf9e05131e5adb81e5aa0422402f048162",
		},
		{
			name:     "Multiple authorizers",
			tx:       baseTx().AddAuthorizer(flow.HexToAddress("02")),
			payload:  "f87bb07472616e73616374696f6e207b2065786563757465207b206c6f67282248656c6c6f2c20576f726c64212229207d207dc0a0f0e4c2f76c58916ec258f246851bea091d14d4247a2fc3e18694461b1816e13b2a880000000000000001040a880000000000000001d2880000000000000001880000000000000002",
			envelope: "f8a2f87bb07472616e73616374696f6e207b2065786563757465207b206c6f67282248656c6c6f2c20576f726c64212229207d207dc0a0f0e4c2f76c58916ec258f246851bea091d14d4247a2fc3e18694461b1816e13b2a880000000000000001040a880000000000000001d2880000000000000001880000000000000002e4e38004a0f7225388c1d69d57e6251c9fda50cbbf9e05131e5adb81e5aa0422402f048162",
		},
		{
			name:     "Single argument",
			tx:       baseTx().AddRawArgument(jsoncdc.MustEncode(cadence.NewString("foo"))),
			payload:  "f893b07472616e73616374696f6e207b2065786563757465207b206c6f67282248656c6c6f2c20576f726c64212229207d207de1a07b2274797065223a22537472696e67222c2276616c7565223a22666f6f227d0aa0f0e4c2f76c58916ec258f246851bea091d14d4247a2fc3e18694461b1816e13b2a880000000000000001040a880000000000000001c9880000000000000001",
			envelope: "f8baf893b07472616e73616374696f6e207b2065786563757465207b206c6f67282248656c6c6f2c20576f726c64212229207d207de1a07b2274797065223a22537472696e67222c2276616c7565223a22666f6f227d0aa0f0e4c2f76c58916ec258f246851bea091d14d4247a2fc3e18694461b1816e13b2a880000000000000001040a880000000000000001c9880000000000000001e4e38004a0f7225388c1d69d57e6251c9fda50cbbf9e05131e5adb81e5aa0422402f048162",
		},
		{
			name: "Multiple arguments",
			tx: baseTx().
				AddRawArgument(jsoncdc.MustEncode(cadence.NewString("foo"))).
				AddRawArgument(jsoncdc.MustEncode(cadence.NewInt(42))),
			payload:  "f8b1b07472616e73616374696f6e207b2065786563757465207b206c6f67282248656c6c6f2c20576f726c64212229207d207df83ea07b2274797065223a22537472696e67222c2276616c7565223a22666f6f227d0a9c7b2274797065223a22496e74222c2276616c7565223a223432227d0aa0f0e4c2f76c58916ec258f246851bea091d14d4247a2fc3e18694461b1816e13b2a880000000000000001040a880000000000000001c9880000000000000001",
			envelope: "f8d8f8b1b07472616e73616374696f6e207b2065786563757465207b206c6f67282248656c6c6f2c20576f726c64212229207d207df83ea07b2274797065223a22537472696e67222c2276616c7565223a22666f6f227d0a9c7b2274797065223a22496e74222c2276616c7565223a223432227d0aa0f0e4c2f76c58916ec258f246851bea091d14d4247a2fc3e18694461b1816e13b2a880000000000000001040a880000000000000001c9880000000000000001e4e38004a0f7225388c1d69d57e6251c9fda50cbbf9e05131e5adb81e5aa0422402f048162",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := hex.EncodeToString(tt.tx.PayloadMessage())
			envelope := hex.EncodeToString(tt.tx.EnvelopeMessage())

			fmt.Println(envelope)

			assert.Equal(t, tt.payload, payload)
			assert.Equal(t, tt.envelope, envelope)

			// Check tx decoding
			transactionBytes := tt.tx.Encode()
			newTx, err := flow.DecodeTransaction(transactionBytes)
			require.NoError(t, err)

			assert.Equal(t, tt.tx, newTx)
			assert.Equal(t, tt.tx.ID(), newTx.ID())

			// Check envelope decoding
			envelopeBytes, err := hex.DecodeString(envelope)
			assert.NoError(t, err)

			newTxFromEnvelope, err := flow.DecodeTransaction(envelopeBytes)
			require.NoError(t, err)

			txEnvelope := copyTxEnvelope(tt.tx)
			assert.Equal(t, txEnvelope, newTxFromEnvelope)
			assert.Equal(t, txEnvelope.ID(), newTxFromEnvelope.ID())

			// Check payload decoding
			payloadBytes, err := hex.DecodeString(payload)
			assert.NoError(t, err)

			newTxFromPayload, err := flow.DecodeTransaction(payloadBytes)
			require.NoError(t, err)

			txPayload := copyTxPayload(tt.tx)
			assert.Equal(t, txPayload, newTxFromPayload)
			assert.Equal(t, txPayload.ID(), newTxFromPayload.ID())
		})
	}
}

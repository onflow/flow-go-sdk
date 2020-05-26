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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/test"
)

func ExampleTransaction() {
	addresses := flow.NewAddressGenerator(flow.Mainnet)

	// Mock user accounts

	adrianLaptopKey := &flow.AccountKey{
		ID:             3,
		SequenceNumber: 42,
	}

	adrianPhoneKey := &flow.AccountKey{ID: 2}
	addressA, _ := addresses.NextAddress()

	adrian := flow.Account{
		Address: addressA,
		Keys: []*flow.AccountKey{
			adrianLaptopKey,
			adrianPhoneKey,
		},
	}

	blaineHardwareKey := &flow.AccountKey{ID: 7}
	addressB, _ := addresses.NextAddress()

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
		SetProposalKey(adrian.Address, adrianLaptopKey.ID, adrianLaptopKey.SequenceNumber).
		SetPayer(blaine.Address).
		AddAuthorizer(adrian.Address)

	fmt.Printf("Transaction ID (before signing): %s\n\n", tx.ID())

	// Signing

	err := tx.SignPayload(adrian.Address, adrianLaptopKey.ID, test.MockSigner([]byte{1}))
	if err != nil {
		panic(err)
	}

	err = tx.SignPayload(adrian.Address, adrianPhoneKey.ID, test.MockSigner([]byte{2}))
	if err != nil {
		panic(err)
	}

	err = tx.SignEnvelope(blaine.Address, blaineHardwareKey.ID, test.MockSigner([]byte{3}))
	if err != nil {
		panic(err)
	}

	fmt.Println("Payload signatures:")
	for _, sig := range tx.PayloadSignatures {
		fmt.Printf(
			"Address: %s, Key ID: %d, Signature: %x\n",
			sig.Address,
			sig.KeyID,
			sig.Signature,
		)
	}
	fmt.Println()

	fmt.Println("Envelope signatures:")
	for _, sig := range tx.EnvelopeSignatures {
		fmt.Printf(
			"Address: %s, Key ID: %d, Signature: %x\n",
			sig.Address,
			sig.KeyID,
			sig.Signature,
		)
	}
	fmt.Println()

	fmt.Printf("Transaction ID (after signing): %s\n", tx.ID())

	// Output:
	// Transaction ID (before signing): 6ba28e395014536e28a4b0a7c122dda7544ffee33015850ed48e8ecf975b46d7
	//
	// Payload signatures:
	// Address: e467b9dd11fa00df, Key ID: 2, Signature: 02
	// Address: e467b9dd11fa00df, Key ID: 3, Signature: 01
	//
	// Envelope signatures:
	// Address: f233dcee88fe0abe, Key ID: 7, Signature: 03
	//
	// Transaction ID (after signing): 6f2220ea52cafc2e637d56f1a7b22d697ad5ad4723cc8b04bb183d607b949c20
}

func TestTransaction_SetScript(t *testing.T) {
	tx := flow.NewTransaction().
		SetScript(test.ScriptHelloWorld)

	assert.Equal(t, test.ScriptHelloWorld, tx.Script)
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
	keyID := 7
	var sequenceNumber uint64 = 42

	tx := flow.NewTransaction().
		SetProposalKey(address, keyID, sequenceNumber)

	assert.Equal(t, address, tx.ProposalKey.Address)
	assert.Equal(t, keyID, tx.ProposalKey.KeyID)
	assert.Equal(t, sequenceNumber, tx.ProposalKey.SequenceNumber)
}

func TestTransaction_SetPayer(t *testing.T) {
	address := flow.ServiceAddress(flow.Mainnet)

	tx := flow.NewTransaction().
		SetPayer(address)

	assert.Equal(t, address, tx.Payer)
}

func TestTransaction_AddAuthorizer(t *testing.T) {
	addresses := flow.NewAddressGenerator(flow.Mainnet)

	addressA, err := addresses.NextAddress()
	require.NoError(t, err)
	addressB, err := addresses.NextAddress()
	require.NoError(t, err)

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
	addresses := flow.NewAddressGenerator(flow.Mainnet)

	t.Run("Invalid signer", func(t *testing.T) {
		tx := flow.NewTransaction()

		address, err := addresses.NextAddress()
		require.NoError(t, err)

		tx.AddPayloadSignature(address, 7, []byte{42})

		require.Len(t, tx.PayloadSignatures, 1)

		// signer cannot be found, so index is -1
		assert.Equal(t, -1, tx.PayloadSignatures[0].SignerIndex)
	})

	t.Run("Valid signers", func(t *testing.T) {
		addressA, err := addresses.NextAddress()
		require.NoError(t, err)
		addressB, err := addresses.NextAddress()
		require.NoError(t, err)

		keyID := 7
		sig := []byte{42}

		tx := flow.NewTransaction().
			AddAuthorizer(addressA).
			AddAuthorizer(addressB)

		// add signatures in reverse order of declaration
		tx.AddPayloadSignature(addressB, keyID, sig)
		tx.AddPayloadSignature(addressA, keyID, sig)

		require.Len(t, tx.PayloadSignatures, 2)

		assert.Equal(t, 0, tx.PayloadSignatures[0].SignerIndex)
		assert.Equal(t, addressA, tx.PayloadSignatures[0].Address)
		assert.Equal(t, keyID, tx.PayloadSignatures[0].KeyID)
		assert.Equal(t, sig, tx.PayloadSignatures[0].Signature)

		assert.Equal(t, 1, tx.PayloadSignatures[1].SignerIndex)
		assert.Equal(t, addressB, tx.PayloadSignatures[1].Address)
		assert.Equal(t, keyID, tx.PayloadSignatures[1].KeyID)
		assert.Equal(t, sig, tx.PayloadSignatures[1].Signature)
	})

	t.Run("Duplicate signers", func(t *testing.T) {
		addressA, err := addresses.NextAddress()
		require.NoError(t, err)
		addressB, err := addresses.NextAddress()
		require.NoError(t, err)

		keyID := 7
		sig := []byte{42}

		tx := flow.NewTransaction().
			SetProposalKey(addressA, keyID, 42).
			AddAuthorizer(addressB).
			AddAuthorizer(addressA)

		// add signatures in reverse order of declaration
		tx.AddPayloadSignature(addressB, keyID, sig)
		tx.AddPayloadSignature(addressA, keyID, sig)

		require.Len(t, tx.PayloadSignatures, 2)

		assert.Equal(t, 0, tx.PayloadSignatures[0].SignerIndex)
		assert.Equal(t, addressA, tx.PayloadSignatures[0].Address)
		assert.Equal(t, keyID, tx.PayloadSignatures[0].KeyID)
		assert.Equal(t, sig, tx.PayloadSignatures[0].Signature)

		assert.Equal(t, 1, tx.PayloadSignatures[1].SignerIndex)
		assert.Equal(t, addressB, tx.PayloadSignatures[1].Address)
		assert.Equal(t, keyID, tx.PayloadSignatures[1].KeyID)
		assert.Equal(t, sig, tx.PayloadSignatures[1].Signature)
	})

	t.Run("Multiple signatures", func(t *testing.T) {
		address, err := addresses.NextAddress()
		require.NoError(t, err)

		keyIDA := 7
		sigA := []byte{42}

		keyIDB := 8
		sigB := []byte{43}

		tx := flow.NewTransaction().
			AddAuthorizer(address)

		// add signatures in descending order by key ID
		tx.AddPayloadSignature(address, keyIDB, sigB)
		tx.AddPayloadSignature(address, keyIDA, sigA)

		require.Len(t, tx.PayloadSignatures, 2)

		// signatures should be sorted in ascending order by key ID
		assert.Equal(t, 0, tx.PayloadSignatures[0].SignerIndex)
		assert.Equal(t, address, tx.PayloadSignatures[0].Address)
		assert.Equal(t, keyIDA, tx.PayloadSignatures[0].KeyID)
		assert.Equal(t, sigA, tx.PayloadSignatures[0].Signature)

		assert.Equal(t, 0, tx.PayloadSignatures[1].SignerIndex)
		assert.Equal(t, address, tx.PayloadSignatures[1].Address)
		assert.Equal(t, keyIDB, tx.PayloadSignatures[1].KeyID)
		assert.Equal(t, sigB, tx.PayloadSignatures[1].Signature)
	})
}

func TestTransaction_AddEnvelopeSignature(t *testing.T) {
	addresses := flow.NewAddressGenerator(flow.Mainnet)

	t.Run("Invalid signer", func(t *testing.T) {
		tx := flow.NewTransaction()

		address, err := addresses.NextAddress()
		require.NoError(t, err)

		tx.AddEnvelopeSignature(address, 7, []byte{42})

		require.Len(t, tx.EnvelopeSignatures, 1)

		// signer cannot be found, so index is -1
		assert.Equal(t, -1, tx.EnvelopeSignatures[0].SignerIndex)
	})

	t.Run("Valid signer", func(t *testing.T) {
		address, err := addresses.NextAddress()
		require.NoError(t, err)

		keyID := 7
		sig := []byte{42}

		tx := flow.NewTransaction().
			SetPayer(address)

		tx.AddEnvelopeSignature(address, keyID, sig)

		require.Len(t, tx.EnvelopeSignatures, 1)

		assert.Equal(t, 0, tx.EnvelopeSignatures[0].SignerIndex)
		assert.Equal(t, address, tx.EnvelopeSignatures[0].Address)
		assert.Equal(t, keyID, tx.EnvelopeSignatures[0].KeyID)
		assert.Equal(t, sig, tx.EnvelopeSignatures[0].Signature)
	})

	t.Run("Multiple signatures", func(t *testing.T) {
		address, err := addresses.NextAddress()
		require.NoError(t, err)

		keyIDA := 7
		sigA := []byte{42}

		keyIDB := 8
		sigB := []byte{43}

		tx := flow.NewTransaction().AddAuthorizer(address)

		// add signatures in descending order by key ID
		tx.AddEnvelopeSignature(address, keyIDB, sigB)
		tx.AddEnvelopeSignature(address, keyIDA, sigA)

		require.Len(t, tx.EnvelopeSignatures, 2)

		// signatures should be sorted in ascending order by key ID
		assert.Equal(t, 0, tx.EnvelopeSignatures[0].SignerIndex)
		assert.Equal(t, address, tx.EnvelopeSignatures[0].Address)
		assert.Equal(t, keyIDA, tx.EnvelopeSignatures[0].KeyID)
		assert.Equal(t, sigA, tx.EnvelopeSignatures[0].Signature)

		assert.Equal(t, 0, tx.EnvelopeSignatures[1].SignerIndex)
		assert.Equal(t, address, tx.EnvelopeSignatures[1].Address)
		assert.Equal(t, keyIDB, tx.EnvelopeSignatures[1].KeyID)
		assert.Equal(t, sigB, tx.EnvelopeSignatures[1].Signature)
	})
}

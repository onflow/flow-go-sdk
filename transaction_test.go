package flow_test

import (
	"fmt"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/crypto"
)

func ExampleTransaction() {
	// Mock user accounts

	adrianLaptopKey := flow.AccountKey{
		ID:             3,
		SequenceNumber: 42,
	}

	adrianPhoneKey := flow.AccountKey{ID: 2}

	adrian := flow.Account{
		Address: flow.HexToAddress("01"),
		Keys: []flow.AccountKey{
			adrianLaptopKey,
			adrianPhoneKey,
		},
	}

	blaineHardwareKey := flow.AccountKey{ID: 7}

	blaine := flow.Account{
		Address: flow.HexToAddress("02"),
		Keys: []flow.AccountKey{
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

	err := tx.SignPayload(adrian.Address, adrianLaptopKey.ID, crypto.MockSigner([]byte{1}))
	if err != nil {
		panic(err)
	}

	err = tx.SignPayload(adrian.Address, adrianPhoneKey.ID, crypto.MockSigner([]byte{2}))
	if err != nil {
		panic(err)
	}

	err = tx.SignEnvelope(blaine.Address, blaineHardwareKey.ID, crypto.MockSigner([]byte{3}))
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
	// Transaction ID (before signing): 142a7862e1718bae9bae109d950a45a2540f98593edf1916ad89f4e22a0960c9
	//
	// Payload signatures:
	// Address: 0000000000000000000000000000000000000001, Key ID: 2, Signature: 02
	// Address: 0000000000000000000000000000000000000001, Key ID: 3, Signature: 01
	//
	// Envelope signatures:
	// Address: 0000000000000000000000000000000000000002, Key ID: 7, Signature: 03
	//
	// Transaction ID (after signing): ae9ef40a6f74117ded103e8bd8c8b0fb97c5a2b3503c1e7ba4aa467f662e31fb
}

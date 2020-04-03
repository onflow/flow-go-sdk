package flow_test

import (
	"fmt"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/crypto"
)

type MockSigner flow.AccountKey

func (s MockSigner) Sign(crypto.Signable) ([]byte, error) {
	return []byte{uint8(s.Index)}, nil
}

func ExampleTransaction() {
	// Mock user accounts

	adrianLaptopKey := flow.AccountKey{
		Index:          3,
		SequenceNumber: 42,
	}

	adrianPhoneKey := flow.AccountKey{Index: 2}

	adrian := flow.Account{
		Address: flow.HexToAddress("01"),
		Keys: []flow.AccountKey{
			adrianLaptopKey,
			adrianPhoneKey,
		},
	}

	blaineHardwareKey := flow.AccountKey{Index: 7}

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
		SetProposalKey(adrian.Address, adrianLaptopKey.Index, adrianLaptopKey.SequenceNumber).
		SetPayer(blaine.Address, blaineHardwareKey.Index).
		AddAuthorizer(adrian.Address, adrianLaptopKey.Index, adrianPhoneKey.Index)

	fmt.Println("Signers:")
	for _, signer := range tx.Signers() {
		fmt.Printf(
			"Address: %s, Roles: %s, Key Indices: %d\n",
			signer.Address,
			signer.Roles,
			signer.Keys,
		)
	}
	fmt.Println()

	fmt.Printf("Transaction ID (before signing): %x\n\n", tx.ID())

	// Signing

	err := tx.SignPayload(adrian.Address, adrianLaptopKey.Index, MockSigner(adrianLaptopKey))
	if err != nil {
		panic(err)
	}

	err = tx.SignPayload(adrian.Address, adrianPhoneKey.Index, MockSigner(adrianPhoneKey))
	if err != nil {
		panic(err)
	}

	err = tx.SignPayer(blaine.Address, blaineHardwareKey.Index, MockSigner(blaineHardwareKey))
	if err != nil {
		panic(err)
	}

	fmt.Println("Payload Signatures:")
	for _, set := range tx.PayloadSignatures() {
		for _, sig := range set.Signatures {
			fmt.Printf(
				"Address: %s, Key Index: %d, Signature: %x\n",
				set.Address,
				sig.KeyIndex,
				sig.Signature,
			)
		}
	}
	fmt.Println()

	fmt.Println("Payer Signatures:")
	for _, sig := range tx.PayerSignatures() {
		fmt.Printf(
			"Address: %s, Key Index: %d, Signature: %x\n",
			tx.Payer().Address,
			sig.KeyIndex,
			sig.Signature,
		)
	}
	fmt.Println()

	fmt.Printf("Transaction ID (after signing): %x\n", tx.ID())

	// Output:
	// Signers:
	// Address: 0000000000000000000000000000000000000001, Roles: [PROPOSER], Key Indices: []
	// Address: 0000000000000000000000000000000000000002, Roles: [PAYER], Key Indices: [7]
	// Address: 0000000000000000000000000000000000000001, Roles: [AUTHORIZER], Key Indices: [2 3]
	//
	// Transaction ID (before signing): 349959c09421ec233b63613f7bb60e4585fbbd8a604b788a0f18cc4f97cd0471
	//
	// Payload Signatures:
	// Address: 0000000000000000000000000000000000000001, Key Index: 2, Signature: 02
	// Address: 0000000000000000000000000000000000000001, Key Index: 3, Signature: 03
	//
	// Payer Signatures:
	// Address: 0000000000000000000000000000000000000002, Key Index: 7, Signature: 07
	//
	// Transaction ID (after signing): 370a571558f5eb9f44f367cac269757acee59c394162e952788fd4b57ec1c504
}

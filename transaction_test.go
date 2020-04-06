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
			signer.KeyIndices,
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

	err = tx.SignContainer(blaine.Address, blaineHardwareKey.Index, MockSigner(blaineHardwareKey))
	if err != nil {
		panic(err)
	}

	fmt.Println("Signatures:")
	for _, sig := range tx.Signatures {
		fmt.Printf(
			"Kind: %s, Address: %s, Key Index: %d, Signature: %x\n",
			sig.Kind,
			sig.Address,
			sig.KeyIndex,
			sig.Signature,
		)
	}
	fmt.Println()

	fmt.Printf("Transaction ID (after signing): %x\n", tx.ID())

	// Output:
	// Signers:
	// Address: 0000000000000000000000000000000000000001, Roles: [PROPOSER AUTHORIZER], Key Indices: [2 3]
	// Address: 0000000000000000000000000000000000000002, Roles: [PAYER], Key Indices: [7]
	//
	// Transaction ID (before signing): 4cd86595c7dc854b371644060c1b4cbc478726b7e3c8be2176353c169e1a76d3
	//
	// Signatures:
	// Kind: PAYLOAD, Address: 0000000000000000000000000000000000000001, Key Index: 3, Signature: 03
	// Kind: PAYLOAD, Address: 0000000000000000000000000000000000000001, Key Index: 2, Signature: 02
	// Kind: CONTAINER, Address: 0000000000000000000000000000000000000002, Key Index: 7, Signature: 07
	//
	// Transaction ID (after signing): 63271c5cb5429bcabbb3fd0f174afd1d22ca4c2e5fb237cf940ce1c61e2176f3
}

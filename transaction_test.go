package flow_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
		SetPayer(blaine.Address, blaineHardwareKey.ID).
		AddAuthorizer(adrian.Address, adrianLaptopKey.ID, adrianPhoneKey.ID)

	fmt.Println("Signers:")
	for _, signer := range tx.Signers() {
		fmt.Printf(
			"Address: %s, Roles: %s, Key IDs: %d\n",
			signer.Address,
			signer.Roles,
			signer.KeyIDs,
		)
	}
	fmt.Println()

	fmt.Printf("Transaction ID (before signing): %x\n\n", tx.ID())

	// Signing

	err := tx.SignPayload(adrian.Address, adrianLaptopKey.ID, crypto.MockSigner([]byte{1}))
	if err != nil {
		panic(err)
	}

	err = tx.SignPayload(adrian.Address, adrianPhoneKey.ID, crypto.MockSigner([]byte{2}))
	if err != nil {
		panic(err)
	}

	err = tx.SignContainer(blaine.Address, blaineHardwareKey.ID, crypto.MockSigner([]byte{3}))
	if err != nil {
		panic(err)
	}

	fmt.Println("Signatures:")
	for _, sig := range tx.Signatures {
		fmt.Printf(
			"%d - Kind: %s, Address: %s, Key ID: %d, Signature: %x\n",
			sig.Index,
			sig.Kind,
			sig.Address,
			sig.KeyID,
			sig.Signature,
		)
	}
	fmt.Println()

	fmt.Printf("Transaction ID (after signing): %x\n", tx.ID())

	// Output:
	// Signers:
	// Address: 0000000000000000000000000000000000000001, Roles: [PROPOSER AUTHORIZER], Key IDs: [2 3]
	// Address: 0000000000000000000000000000000000000002, Roles: [PAYER], Key IDs: [7]
	//
	// Transaction ID (before signing): 4cd86595c7dc854b371644060c1b4cbc478726b7e3c8be2176353c169e1a76d3
	//
	// Signatures:
	// 0 - Kind: PAYLOAD, Address: 0000000000000000000000000000000000000001, Key ID: 2, Signature: 02
	// 1 - Kind: PAYLOAD, Address: 0000000000000000000000000000000000000001, Key ID: 3, Signature: 03
	// 2 - Kind: CONTAINER, Address: 0000000000000000000000000000000000000002, Key ID: 7, Signature: 07
	//
	// Transaction ID (after signing): 395cf1a841d82569c8fedd67678c67fefda7b76436be581d405ab39bfbe35263
}

var (
	AddressA flow.Address
	AddressB flow.Address
	AddressC flow.Address
	AddressD flow.Address
	AddressE flow.Address

	RolesProposerPayerAuthorizer []flow.SignerRole
	RolesProposerPayer           []flow.SignerRole
	RolesProposerAuthorizer      []flow.SignerRole
	RolesPayerAuthorizer         []flow.SignerRole
	RolesProposer                []flow.SignerRole
	RolesPayer                   []flow.SignerRole
	RolesAuthorizer              []flow.SignerRole
)

func init() {
	AddressA = flow.HexToAddress("01")
	AddressB = flow.HexToAddress("02")
	AddressC = flow.HexToAddress("03")
	AddressD = flow.HexToAddress("04")
	AddressE = flow.HexToAddress("05")

	RolesProposerPayerAuthorizer = []flow.SignerRole{flow.SignerRoleProposer, flow.SignerRolePayer, flow.SignerRoleAuthorizer}
	RolesProposerPayer = []flow.SignerRole{flow.SignerRoleProposer, flow.SignerRolePayer}
	RolesProposerAuthorizer = []flow.SignerRole{flow.SignerRoleProposer, flow.SignerRoleAuthorizer}
	RolesPayerAuthorizer = []flow.SignerRole{flow.SignerRolePayer, flow.SignerRoleAuthorizer}
	RolesProposer = []flow.SignerRole{flow.SignerRoleProposer}
	RolesPayer = []flow.SignerRole{flow.SignerRolePayer}
	RolesAuthorizer = []flow.SignerRole{flow.SignerRoleAuthorizer}
}

func TestTransaction_Signers_SeparateSigners(t *testing.T) {
	t.Run("No authorizers", func(t *testing.T) {
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			SetPayer(AddressB, 1).
			Signers()

		require.Len(t, signers, 2)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)
		assert.Equal(t, []int{1}, signers[0].KeyIDs)

		assert.Equal(t, AddressB, signers[1].Address)
		assert.Equal(t, RolesPayer, signers[1].Roles)
		assert.Equal(t, []int{1}, signers[1].KeyIDs)
	})

	t.Run("With authorizer", func(t *testing.T) {
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			AddAuthorizer(AddressB, 1).
			SetPayer(AddressC, 1).
			Signers()

		require.Len(t, signers, 3)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)
		assert.Equal(t, []int{1}, signers[0].KeyIDs)

		assert.Equal(t, AddressB, signers[1].Address)
		assert.Equal(t, RolesAuthorizer, signers[1].Roles)
		assert.Equal(t, []int{1}, signers[1].KeyIDs)

		assert.Equal(t, AddressC, signers[2].Address)
		assert.Equal(t, RolesPayer, signers[2].Roles)
		assert.Equal(t, []int{1}, signers[2].KeyIDs)
	})
}

func TestTransaction_Signers_DeclarationOrder(t *testing.T) {
	t.Run("Payer before proposer", func(t *testing.T) {
		signers := flow.NewTransaction().
			SetPayer(AddressB, 1).
			SetProposalKey(AddressA, 1, 42).
			Signers()

		require.Len(t, signers, 2)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)
		assert.Equal(t, []int{1}, signers[0].KeyIDs)

		assert.Equal(t, AddressB, signers[1].Address)
		assert.Equal(t, RolesPayer, signers[1].Roles)
		assert.Equal(t, []int{1}, signers[1].KeyIDs)
	})

	t.Run("Authorizer before proposer", func(t *testing.T) {
		signers := flow.NewTransaction().
			AddAuthorizer(AddressB, 1).
			SetProposalKey(AddressA, 1, 42).
			SetPayer(AddressC, 1).
			Signers()

		require.Len(t, signers, 3)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)
		assert.Equal(t, []int{1}, signers[0].KeyIDs)

		assert.Equal(t, AddressB, signers[1].Address)
		assert.Equal(t, RolesAuthorizer, signers[1].Roles)
		assert.Equal(t, []int{1}, signers[1].KeyIDs)

		assert.Equal(t, AddressC, signers[2].Address)
		assert.Equal(t, RolesPayer, signers[2].Roles)
		assert.Equal(t, []int{1}, signers[2].KeyIDs)
	})

	t.Run("Authorizer after payer", func(t *testing.T) {
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			SetPayer(AddressC, 1).
			AddAuthorizer(AddressB, 1).
			Signers()

		require.Len(t, signers, 3)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)
		assert.Equal(t, []int{1}, signers[0].KeyIDs)

		assert.Equal(t, AddressB, signers[1].Address)
		assert.Equal(t, RolesAuthorizer, signers[1].Roles)
		assert.Equal(t, []int{1}, signers[1].KeyIDs)

		assert.Equal(t, AddressC, signers[2].Address)
		assert.Equal(t, RolesPayer, signers[2].Roles)
		assert.Equal(t, []int{1}, signers[2].KeyIDs)
	})

	t.Run("Authorizer before and after payer", func(t *testing.T) {
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			AddAuthorizer(AddressB, 1).
			SetPayer(AddressD, 1).
			AddAuthorizer(AddressC, 1).
			Signers()

		require.Len(t, signers, 4)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)
		assert.Equal(t, []int{1}, signers[0].KeyIDs)

		assert.Equal(t, AddressB, signers[1].Address)
		assert.Equal(t, RolesAuthorizer, signers[1].Roles)
		assert.Equal(t, []int{1}, signers[1].KeyIDs)

		assert.Equal(t, AddressC, signers[2].Address)
		assert.Equal(t, RolesAuthorizer, signers[2].Roles)
		assert.Equal(t, []int{1}, signers[2].KeyIDs)

		assert.Equal(t, AddressD, signers[3].Address)
		assert.Equal(t, RolesPayer, signers[3].Roles)
		assert.Equal(t, []int{1}, signers[3].KeyIDs)
	})
}

func TestTransaction_Signers_KeysOutOfOrder(t *testing.T) {
	signers := flow.NewTransaction().
		SetProposalKey(AddressA, 1, 42).
		SetPayer(AddressA, 4, 2, 1, 3).
		Signers()

	require.Len(t, signers, 1)

	assert.Equal(t, AddressA, signers[0].Address)
	assert.Equal(t, RolesProposerPayer, signers[0].Roles)
	assert.Equal(t, []int{1, 2, 3, 4}, signers[0].KeyIDs)
}

func TestTransaction_Signers_MultipleAuthorizers(t *testing.T) {
	signers := flow.NewTransaction().
		SetProposalKey(AddressA, 1, 42).
		AddAuthorizer(AddressB, 1).
		AddAuthorizer(AddressC, 2).
		AddAuthorizer(AddressD, 3).
		SetPayer(AddressE, 1).
		Signers()

	require.Len(t, signers, 5)

	assert.Equal(t, AddressB, signers[1].Address)
	assert.Equal(t, RolesAuthorizer, signers[1].Roles)
	assert.Equal(t, []int{1}, signers[1].KeyIDs)

	assert.Equal(t, AddressC, signers[2].Address)
	assert.Equal(t, RolesAuthorizer, signers[2].Roles)
	assert.Equal(t, []int{2}, signers[2].KeyIDs)

	assert.Equal(t, AddressD, signers[3].Address)
	assert.Equal(t, RolesAuthorizer, signers[3].Roles)
	assert.Equal(t, []int{3}, signers[3].KeyIDs)
}

func TestTransaction_Signers_ProposerPayerAuthorizerSameAddress(t *testing.T) {
	t.Run("Single key", func(t *testing.T) {
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			SetPayer(AddressA, 1).
			AddAuthorizer(AddressA, 1).
			Signers()

		require.Len(t, signers, 1)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposerPayerAuthorizer, signers[0].Roles)
		assert.Equal(
			t,
			&flow.ProposalKey{
				Address:        AddressA,
				KeyID:          1,
				SequenceNumber: 42,
			},
			signers[0].ProposalKey,
		)
		assert.Equal(t, []int{1}, signers[0].KeyIDs)
	})

	t.Run("Identical key-sets", func(t *testing.T) {
		// All key-sets contain the elements [1, 2]
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			SetPayer(AddressA, 1, 2).
			AddAuthorizer(AddressA, 1, 2).
			Signers()

		require.Len(t, signers, 1)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposerPayerAuthorizer, signers[0].Roles)
		assert.Equal(t, []int{1, 2}, signers[0].KeyIDs)
	})

	t.Run("Subset of payer key-set", func(t *testing.T) {
		// Payer key-set: [1, 2]
		// Authorizer key-set: [1]
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			SetPayer(AddressA, 1, 2).
			AddAuthorizer(AddressA, 1).
			Signers()

		require.Len(t, signers, 1)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposerPayerAuthorizer, signers[0].Roles)
		assert.Equal(t, []int{1, 2}, signers[0].KeyIDs)
	})
}

func TestTransaction_Signers_ProposerPayerSameAddress(t *testing.T) {
	t.Run("No authorizers", func(t *testing.T) {
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			SetPayer(AddressA, 1, 2).
			Signers()

		require.Len(t, signers, 1)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposerPayer, signers[0].Roles)
		assert.Equal(t, []int{1, 2}, signers[0].KeyIDs)
	})

	t.Run("With authorizer", func(t *testing.T) {
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			AddAuthorizer(AddressB, 1).
			SetPayer(AddressA, 1, 2).
			Signers()

		require.Len(t, signers, 2)

		assert.Equal(t, AddressB, signers[0].Address)
		assert.Equal(t, RolesAuthorizer, signers[0].Roles)
		assert.Equal(t, []int{1}, signers[0].KeyIDs)

		assert.Equal(t, AddressA, signers[1].Address)
		assert.Equal(t, RolesProposerPayer, signers[1].Roles)
		assert.Equal(t, []int{1, 2}, signers[1].KeyIDs)
	})

	t.Run("Disjoint key-sets", func(t *testing.T) {
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			SetPayer(AddressA, 2, 3).
			Signers()

		require.Len(t, signers, 2)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)
		assert.Equal(t, []int{1}, signers[0].KeyIDs)

		assert.Equal(t, AddressA, signers[1].Address)
		assert.Equal(t, RolesPayer, signers[1].Roles)
		assert.Equal(t, []int{2, 3}, signers[1].KeyIDs)
	})
}

func TestTransaction_Signers_PayerAuthorizerSameAddress(t *testing.T) {
	t.Run("Identical key-sets", func(t *testing.T) {
		// Payer key-set: [1, 2]
		// Authorizer key-set: [1, 2]
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			AddAuthorizer(AddressB, 1, 2).
			SetPayer(AddressB, 1, 2).
			Signers()

		require.Len(t, signers, 2)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)

		assert.Equal(t, AddressB, signers[1].Address)
		assert.Equal(t, RolesPayerAuthorizer, signers[1].Roles)
		assert.Equal(t, []int{1, 2}, signers[1].KeyIDs)
	})

	t.Run("Subset of payer key-set", func(t *testing.T) {
		// Payer key-set: [1, 2]
		// Authorizer key-set: [1]
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			AddAuthorizer(AddressB, 1).
			SetPayer(AddressB, 1, 2).
			Signers()

		require.Len(t, signers, 2)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)

		assert.Equal(t, AddressB, signers[1].Address)
		assert.Equal(t, RolesPayerAuthorizer, signers[1].Roles)
		assert.Equal(t, []int{1, 2}, signers[1].KeyIDs)
	})

	t.Run("Disjoint key-sets", func(t *testing.T) {
		// Payer key-set: [1, 2]
		// Authorizer key-set: [3, 4]
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			AddAuthorizer(AddressB, 3, 4).
			SetPayer(AddressB, 1, 2).
			Signers()

		require.Len(t, signers, 3)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)
		assert.Equal(t, []int{1}, signers[0].KeyIDs)

		assert.Equal(t, AddressB, signers[1].Address)
		assert.Equal(t, RolesAuthorizer, signers[1].Roles)
		assert.Equal(t, []int{3, 4}, signers[1].KeyIDs)

		assert.Equal(t, AddressB, signers[2].Address)
		assert.Equal(t, RolesPayer, signers[2].Roles)
		assert.Equal(t, []int{1, 2}, signers[2].KeyIDs)
	})

	t.Run("Overlapping key-sets", func(t *testing.T) {
		// Payer key-set: [1, 2]
		// Authorizer key-set: [2, 3]
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			AddAuthorizer(AddressB, 2, 3).
			SetPayer(AddressB, 1, 2).
			Signers()

		require.Len(t, signers, 3)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)
		assert.Equal(t, []int{1}, signers[0].KeyIDs)

		assert.Equal(t, AddressB, signers[1].Address)
		assert.Equal(t, RolesAuthorizer, signers[1].Roles)
		assert.Equal(t, []int{2, 3}, signers[1].KeyIDs)

		assert.Equal(t, AddressB, signers[2].Address)
		assert.Equal(t, RolesPayer, signers[2].Roles)
		assert.Equal(t, []int{1, 2}, signers[2].KeyIDs)
	})
}

func TestTransaction_Signers_ProposerAuthorizerSameAddress(t *testing.T) {
	t.Run("Overlapping key-sets", func(t *testing.T) {
		// Proposal key: 1
		// Authorizer key-set: [1, 2]
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			AddAuthorizer(AddressA, 1, 2).
			SetPayer(AddressB, 1).
			Signers()

		require.Len(t, signers, 2)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposerAuthorizer, signers[0].Roles)
		assert.Equal(t, []int{1, 2}, signers[0].KeyIDs)

		assert.Equal(t, AddressB, signers[1].Address)
		assert.Equal(t, RolesPayer, signers[1].Roles)
		assert.Equal(t, []int{1}, signers[1].KeyIDs)
	})

	t.Run("Disjoint key-sets", func(t *testing.T) {
		// Proposal key: 1
		// Authorizer key-set: [2]
		signers := flow.NewTransaction().
			SetProposalKey(AddressA, 1, 42).
			AddAuthorizer(AddressA, 2).
			SetPayer(AddressB, 1).
			Signers()

		require.Len(t, signers, 3)

		assert.Equal(t, AddressA, signers[0].Address)
		assert.Equal(t, RolesProposer, signers[0].Roles)
		assert.Equal(t, []int{1}, signers[0].KeyIDs)

		assert.Equal(t, AddressA, signers[1].Address)
		assert.Equal(t, RolesAuthorizer, signers[1].Roles)
		assert.Equal(t, []int{2}, signers[1].KeyIDs)

		assert.Equal(t, AddressB, signers[2].Address)
		assert.Equal(t, RolesPayer, signers[2].Roles)
		assert.Equal(t, []int{1}, signers[2].KeyIDs)
	})
}

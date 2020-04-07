package convert_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/client/convert"
	"github.com/dapperlabs/flow-go-sdk/crypto"
)

var (
	AddressA flow.Address
	AddressB flow.Address
	AddressC flow.Address
)

func init() {
	AddressA = flow.HexToAddress("01")
	AddressB = flow.HexToAddress("02")
	AddressC = flow.HexToAddress("03")
}

func TestConvert_Transaction(t *testing.T) {
	txA := flow.NewTransaction().
		SetScript([]byte(`transaction { execute { log("Hello, World!") } }`)).
		SetReferenceBlockID(flow.Identifier{0x01, 0x02}).
		SetGasLimit(42).
		SetProposalKey(AddressA, 1, 42).
		AddAuthorizer(AddressA, 1, 2).
		SetPayer(AddressB, 1)

	err := txA.SignPayload(AddressA, 1, crypto.MockSigner([]byte{1}))
	if err != nil {
		panic(err)
	}

	err = txA.SignPayload(AddressA, 2, crypto.MockSigner([]byte{2}))
	if err != nil {
		panic(err)
	}

	err = txA.SignContainer(AddressB, 1, crypto.MockSigner([]byte{1}))
	if err != nil {
		panic(err)
	}

	msg := convert.TransactionToMessage(*txA)

	txB, err := convert.MessageToTransaction(msg)

	assert.NoError(t, err)
	assert.Equal(t, txA.ID(), txB.ID())
}

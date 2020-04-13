package convert_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/client/convert"
	"github.com/dapperlabs/flow-go-sdk/test"
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
	txA := test.TransactionGenerator().New()

	msg := convert.TransactionToMessage(*txA)

	txB, err := convert.MessageToTransaction(msg)

	assert.NoError(t, err)
	assert.Equal(t, txA.ID(), txB.ID())
}

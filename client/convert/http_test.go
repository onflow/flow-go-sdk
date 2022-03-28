package convert

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/onflow/flow-go-sdk/test"
	"github.com/stretchr/testify/assert"
)

func Test_ConvertBlock(t *testing.T) {

}

func Test_ConvertAccount(t *testing.T) {
	httpAccount := test.AccountHTTP()
	contractName, contractCode := test.ContractHTTP()

	account, err := HTTPToAccount(&httpAccount)

	assert.NoError(t, err)
	assert.Equal(t, account.Address.String(), httpAccount.Address)
	assert.Len(t, account.Keys, 1)
	assert.Equal(t, account.Keys[0].PublicKey.String(), httpAccount.Keys[0].PublicKey)
	code, _ := base64.StdEncoding.DecodeString(contractCode)
	assert.Equal(t, account.Contracts[contractName], code)
	assert.Equal(t, fmt.Sprintf("%d", account.Balance), httpAccount.Balance)
}

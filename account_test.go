package flow_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/test"
)

func TestEncodingAccountKeySequenceNumber(t *testing.T) {

	generator := test.AccountKeyGenerator()
	accountKey := generator.New()
	accountKey.SequenceNumber = 2137

	bytes := accountKey.Encode()
	decodedAccountKey, err := flow.DecodeAccountKey(bytes)
	require.NoError(t, err)

	assert.Equal(t, accountKey.SequenceNumber, decodedAccountKey.SequenceNumber)
}

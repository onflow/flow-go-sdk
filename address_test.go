package flow_test

import (
	"encoding/json"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go-sdk"
)

type addressWrapper struct {
	Address flow.Address
}

func TestAddress_JSON(t *testing.T) {
	addr := flow.RootAddress
	data, err := json.Marshal(addressWrapper{Address: addr})
	require.Nil(t, err)

	var out addressWrapper
	err = json.Unmarshal(data, &out)
	require.Nil(t, err)
	assert.Equal(t, addr, out.Address)
}

func TestAddress_Short(t *testing.T) {
	type testcase struct {
		addr     flow.Address
		expected string
	}

	cases := []testcase{
		{
			addr:     flow.RootAddress,
			expected: "01",
		},
		{
			addr:     flow.HexToAddress("0000000002"),
			expected: "02",
		},
		{
			addr:     flow.HexToAddress("1f10"),
			expected: "1f10",
		},
		{
			addr:     flow.HexToAddress("0f10"),
			expected: "0f10",
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.addr.Short(), c.expected)
	}
}

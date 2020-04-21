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

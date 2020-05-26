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

package flow

import (
	"encoding/json"
	"math/bits"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type addressWrapper struct {
	Address Address
}

func TestAddressJSON(t *testing.T) {
	addr := ServiceAddress(Mainnet)
	data, err := json.Marshal(addressWrapper{Address: addr})
	require.Nil(t, err)

	t.Log(string(data))

	var out addressWrapper
	err = json.Unmarshal(data, &out)
	require.Nil(t, err)
	assert.Equal(t, addr, out.Address)
}

func TestFlowAddressConstants(t *testing.T) {
	// check n and k fit in 8 and 6 bytes
	assert.LessOrEqual(t, linearCodeN, 8*8)
	assert.LessOrEqual(t, linearCodeK, 6*8)

	// Test addresses for all type of networks
	networks := []ChainID{
		Mainnet,
		Testnet,
		Emulator,
	}

	for _, net := range networks {

		// check the Zero and Service constants
		expected := uint64ToAddress(chainCustomizer(net))
		assert.Equal(t, zeroAddress(net), expected)
		expected = uint64ToAddress(generatorMatrixRows[0] ^ chainCustomizer(net))
		assert.Equal(t, ServiceAddress(net), expected)

		// check the transition from account zero to service
		generator := NewAddressGenerator(net)
		address, err := generator.NextAddress()
		require.NoError(t, err)
		assert.Equal(t, address, ServiceAddress(net))

		// check high state values: generation should fail for high value states
		generator = newAddressGeneratorAtState(net, maxState)
		_, err = generator.NextAddress()
		assert.NoError(t, err)
		_, err = generator.NextAddress()
		assert.Error(t, err)

		// check ZeroAddress(net) is an invalid addresse
		z := zeroAddress(net)
		check := z.IsValid(net)
		assert.False(t, check, "should be invalid")
	}
}

const invalidCodeWord = uint64(0xab2ae42382900010)

func TestAddressGeneration(t *testing.T) {
	// seed random generator
	rand.Seed(time.Now().UnixNano())

	// loops in each test
	const loop = 3

	// Test addresses for all type of networks
	networks := []ChainID{
		Mainnet,
		Testnet,
		Emulator,
	}

	for _, net := range networks {

		// sanity check of AccountAddress function consistency
		generator := NewAddressGenerator(net)
		expectedState := zeroAddressState
		for i := 0; i < loop; i++ {
			address, err := generator.NextAddress()
			require.NoError(t, err)
			expectedState++
			expectedAddress := generateAddress(net, expectedState)
			assert.Equal(t, address, expectedAddress)
		}

		// sanity check of addresses weights in Flow.
		// All addresses hamming weights must be less than d.
		// this is only a sanity check of the implementation and not an exhaustive proof
		if net == Mainnet {
			r := rand.Intn(maxState - loop)
			generator := newAddressGeneratorAtState(net, addressState(r))
			for i := 0; i < loop; i++ {
				address, err := generator.NextAddress()
				require.NoError(t, err)
				weight := bits.OnesCount64(address.uint64())
				assert.LessOrEqual(t, linearCodeD, weight)
			}
		}

		// sanity check of address distances.
		// All distances between any two addresses must be less than d.
		// this is only a sanity check of the implementation and not an exhaustive proof
		r := rand.Intn(maxState - loop - 1)
		generator = newAddressGeneratorAtState(net, addressState(r))
		refAddress, err := generator.NextAddress()
		require.NoError(t, err)
		for i := 0; i < loop; i++ {
			address, err := generator.NextAddress()
			require.NoError(t, err)
			distance := bits.OnesCount64(address.uint64() ^ refAddress.uint64())
			assert.LessOrEqual(t, linearCodeD, distance)
		}

		// sanity check of valid account addresses.
		// All valid addresses must pass IsValid.
		r = rand.Intn(maxState - loop)
		generator = newAddressGeneratorAtState(net, addressState(r))
		for i := 0; i < loop; i++ {
			address, err := generator.NextAddress()
			require.NoError(t, err)
			check := address.IsValid(net)
			assert.True(t, check, "account address format should be valid")
		}

		// sanity check of invalid account addresses.
		// All invalid addresses must fail IsValid.
		invalidAddress := uint64ToAddress(invalidCodeWord)
		check := invalidAddress.IsValid(net)
		assert.False(t, check, "account address format should be invalid")
		r = rand.Intn(maxState - loop)
		generator = newAddressGeneratorAtState(net, addressState(r))
		for i := 0; i < loop; i++ {
			address, err := generator.NextAddress()
			require.NoError(t, err)
			invalidAddress = uint64ToAddress(address.uint64() ^ invalidCodeWord)
			check := invalidAddress.IsValid(net)
			assert.False(t, check, "account address format should be invalid")
		}
	}
}

func TestAddressesIntersection(t *testing.T) {
	// seed random generator
	rand.Seed(time.Now().UnixNano())

	// loops in each test
	const loop = 50

	// Test addresses for all type of networks
	networks := []ChainID{
		Testnet,
		Emulator,
	}

	for _, net := range networks {

		// All valid test addresses must fail Flow Mainnet check
		r := rand.Intn(maxState - loop)
		generator := newAddressGeneratorAtState(net, addressState(r))
		for i := 0; i < loop; i++ {
			address, err := generator.NextAddress()
			require.NoError(t, err)
			check := address.IsValid(Mainnet)
			assert.False(t, check, "test account address format should be invalid in Flow")
		}

		// sanity check: mainnet addresses must fail the test check
		r = rand.Intn(maxState - loop)
		generator = newAddressGeneratorAtState(Mainnet, addressState(r))
		for i := 0; i < loop; i++ {
			invalidAddress, err := generator.NextAddress()
			require.NoError(t, err)
			check := invalidAddress.IsValid(net)
			assert.False(t, check, "account address format should be invalid")
		}

		// sanity check of invalid account addresses in all networks
		require.NotEqual(t, invalidCodeWord, uint64(0))
		invalidAddress := uint64ToAddress(invalidCodeWord)
		check := invalidAddress.IsValid(net)
		assert.False(t, check, "account address format should be invalid")
		r = rand.Intn(maxState - loop)
		generator = newAddressGeneratorAtState(net, addressState(r))
		for i := 0; i < loop; i++ {
			address, err := generator.NextAddress()
			require.NoError(t, err)
			invalidAddress = uint64ToAddress(address.uint64() ^ invalidCodeWord)
			// must fail test network check
			check = invalidAddress.IsValid(net)
			assert.False(t, check, "account address format should be invalid")
			// must fail mainnet check
			check := invalidAddress.IsValid(Mainnet)
			assert.False(t, check, "account address format should be invalid")
		}
	}
}

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

func TestAddressConstants(t *testing.T) {
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
		// check the zero and service constants
		expected := uint64ToAddress(chainCustomizer(net))
		assert.Equal(t, zeroAddress(net), expected)
		expected = uint64ToAddress(generatorMatrixRows[0] ^ chainCustomizer(net))
		assert.Equal(t, ServiceAddress(net), expected)

		// check the transition from account zero to service
		generator := NewAddressGenerator(net)
		address := generator.NextAddress()
		assert.Equal(t, address, ServiceAddress(net))

		// check high state values: generation should fail for high value states
		generator = newAddressGeneratorAtState(net, maxState)
		assert.NotPanics(t, func() { generator.NextAddress() })
		assert.Panics(t, func() { generator.NextAddress() })

		// check zeroAddress(net) is an invalid address
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
	const iterations = 3

	// Test addresses for all type of networks
	networks := []ChainID{
		Mainnet,
		Testnet,
		Emulator,
	}

	for _, net := range networks {

		t.Run(net.String(), func(t *testing.T) {

			t.Run("NextAddress", func(t *testing.T) {
				// sanity check of NextAddress function consistency
				generator := NewAddressGenerator(net)
				expectedState := zeroAddressState
				for i := 0; i < iterations; i++ {
					address := generator.NextAddress()
					expectedState++
					expectedAddress := generateAddress(net, expectedState)
					assert.Equal(t, address, expectedAddress)
				}
			})

			t.Run("Address", func(t *testing.T) {
				// sanity check of Address function consistency
				generator := NewAddressGenerator(net)
				expectedState := zeroAddressState
				for i := 0; i < iterations; i++ {
					address := generator.Address()
					expectedAddress := generateAddress(net, expectedState)

					assert.Equal(t, address, expectedAddress)

					generator.Next()
					expectedState++
				}
			})

			t.Run("SetIndex", func(t *testing.T) {
				const indexA = 8
				const indexB = 16

				generatorA := NewAddressGenerator(net)
				generatorB := NewAddressGenerator(net)

				// fast-forward manually (to indexA)
				for i := 0; i < indexA; i++ {
					generatorA.Next()
				}

				// fast-forward with SetIndex (to indexA)
				generatorB.SetIndex(indexA)

				addressA1 := generatorA.Address()
				addressB1 := generatorB.Address()

				assert.Equal(t, addressA1, addressB1)

				// fast-forward manually (to indexB)
				for i := indexA; i < indexB; i++ {
					generatorA.Next()
				}

				// fast-forward with SetIndex (to indexB)
				generatorB.SetIndex(indexB)

				addressA2 := generatorA.Address()
				addressB2 := generatorB.Address()

				assert.Equal(t, addressA2, addressB2)

				// rewind with SetIndex (back to indexA)
				generatorB.SetIndex(indexA)
				addressB3 := generatorB.Address()

				assert.Equal(t, addressA1, addressB3)
			})

			t.Run("Weights", func(t *testing.T) {
				// sanity check of addresses weights in Flow.
				// All addresses hamming weights must be less than d.
				// this is only a sanity check of the implementation and not an exhaustive proof
				if net == Mainnet {
					r := rand.Intn(maxState - iterations)
					generator := newAddressGeneratorAtState(net, addressState(r))
					for i := 0; i < iterations; i++ {
						address := generator.NextAddress()
						weight := bits.OnesCount64(address.uint64())
						assert.LessOrEqual(t, linearCodeD, weight)
					}
				}
			})

			t.Run("Distances", func(t *testing.T) {
				// sanity check of address distances.
				// All distances between any two addresses must be less than d.
				// this is only a sanity check of the implementation and not an exhaustive proof
				r := rand.Intn(maxState - iterations - 1)
				generator := newAddressGeneratorAtState(net, addressState(r))
				refAddress := generator.NextAddress()
				for i := 0; i < iterations; i++ {
					address := generator.NextAddress()
					distance := bits.OnesCount64(address.uint64() ^ refAddress.uint64())
					assert.LessOrEqual(t, linearCodeD, distance)
				}
			})

			t.Run("Valid", func(t *testing.T) {
				// sanity check of valid account addresses.
				// All valid addresses must pass IsValid.
				r := rand.Intn(maxState - iterations)
				generator := newAddressGeneratorAtState(net, addressState(r))
				for i := 0; i < iterations; i++ {
					address := generator.NextAddress()
					check := address.IsValid(net)
					assert.True(t, check, "account address format should be valid")
				}
			})

			t.Run("Invalid", func(t *testing.T) {
				// sanity check of invalid account addresses.
				// All invalid addresses must fail IsValid.
				invalidAddress := uint64ToAddress(invalidCodeWord)
				check := invalidAddress.IsValid(net)
				assert.False(t, check, "account address format should be invalid")

				r := rand.Intn(maxState - iterations)
				generator := newAddressGeneratorAtState(net, addressState(r))
				for i := 0; i < iterations; i++ {
					address := generator.NextAddress()
					invalidAddress = uint64ToAddress(address.uint64() ^ invalidCodeWord)
					check := invalidAddress.IsValid(net)
					assert.False(t, check, "account address format should be invalid")
				}
			})
		})
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
			address := generator.NextAddress()
			check := address.IsValid(Mainnet)
			assert.False(t, check, "test account address format should be invalid in Flow")
		}

		// sanity check: mainnet addresses must fail the test check
		r = rand.Intn(maxState - loop)
		generator = newAddressGeneratorAtState(Mainnet, addressState(r))
		for i := 0; i < loop; i++ {
			invalidAddress := generator.NextAddress()
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
			address := generator.NextAddress()
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

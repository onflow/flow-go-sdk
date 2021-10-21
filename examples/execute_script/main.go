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

package main

import (
	"context"
	"fmt"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/examples"
)

func main() {
	prepareDemo()
	demo()
}

func demo() {
	ctx := context.Background()
	flowClient := examples.NewFlowClient()

	script := []byte(`
		pub fun main(a: Int): Int {
			return a + 10
		}
	`)
	args := []cadence.Value{cadence.NewInt(5)}
	value, err := flowClient.ExecuteScriptAtLatestBlock(ctx, script, args)

	examples.Handle(err)
	fmt.Printf("\nValue: %s", value.String())

	complexScript := []byte(`
		pub struct User {
			pub var balance: UFix64
			pub var address: Address
			pub var name: String

			init(name: String, address: Address, balance: UFix64) {
				self.name = name
				self.address = address
				self.balance = balance
			}
		}

		pub fun main(name: String): User {
			return User(
				name: name,
				address: 0x1,
				balance: 10.0
			)
		}
	`)
	args = []cadence.Value{cadence.String("Dete")}
	value, err = flowClient.ExecuteScriptAtLatestBlock(ctx, complexScript, args)
	printComplexScript(value, err)
}

type User struct {
	balance uint64
	address flow.Address
	name    string
}

func printComplexScript(value cadence.Value, err error) {
	examples.Handle(err)
	fmt.Printf("\nString value: %s", value.String())

	s := value.(cadence.Struct)
	u := User{
		balance: s.Fields[0].ToGoValue().(uint64),
		address: s.Fields[1].ToGoValue().([flow.AddressLength]byte),
		name:    s.Fields[2].ToGoValue().(string),
	}

	fmt.Printf("\nName: %s", u.name)
	fmt.Printf("\nAddress: %s", u.address.String())
	fmt.Printf("\nBalance: %d", u.balance)
}

func prepareDemo() {

}

/*
 * Flow Go SDK
 *
 * Copyright Flow Foundation
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
	"errors"
	"fmt"

	"github.com/onflow/flow-go-sdk/access/http"

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
	flowClient, err := http.NewClient(http.EmulatorHost)
	examples.Handle(err)

	script := []byte(`
		access(all) fun main(a: Int): Int {
			return a + 10
		}
	`)
	args := []cadence.Value{cadence.NewInt(5)}
	value, err := flowClient.ExecuteScriptAtLatestBlock(ctx, script, args)

	examples.Handle(err)
	fmt.Printf("\nValue: %s", value.String())

	complexScript := []byte(`
		access(all) struct User {
			access(all) var balance: UFix64
			access(all) var address: Address
			access(all) var name: String

			init(name: String, address: Address, balance: UFix64) {
				self.name = name
				self.address = address
				self.balance = balance
			}
		}

		access(all) fun main(name: String): User {
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
	balance string
	address flow.Address
	name    string
}

func printComplexScript(value cadence.Value, err error) {
	examples.Handle(err)
	fmt.Printf("\nString value: %s", value.String())

	s := value.(cadence.Struct)
	balanceCdc, ok := s.FieldsMappedByName()["balance"].(cadence.UFix64)
	if !ok {
		examples.Handle(errors.New("incorrect balance"))
	}
	addressCdc, ok := s.FieldsMappedByName()["address"].(cadence.Address)
	if !ok {
		examples.Handle(errors.New("incorrect address"))
	}
	nameCdc, ok := s.FieldsMappedByName()["name"].(cadence.String)
	if !ok {
		examples.Handle(errors.New("incorrect name"))
	}

	u := User{
		balance: balanceCdc.String(),
		address: flow.BytesToAddress(addressCdc.Bytes()),
		name:    nameCdc.String(),
	}

	fmt.Printf("\nName: %s", u.name)
	fmt.Printf("\nAddress: 0x%s", u.address.String())
	fmt.Printf("\nBalance: %s", u.balance)
}

func prepareDemo() {

}

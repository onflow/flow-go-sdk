/*
 * Flow Go SDK
 *
 * Copyright 2019 Dapper Labs, Inc.
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
	flowClient, err := http.NewClient(http.TestnetHost)
	examples.Handle(err)

	script := []byte(`
		import Weekday from 0x2a37a78609bba037
		import MetadataViews from 0x631e88ae7f1d7c20
		
		pub fun main(address: Address, nftId: UInt64): AnyStruct? {
			let collectionRef = getAccount(address).getCapability<&{MetadataViews.ResolverCollection, Weekday.WeekdayCollectionPublic}>(Weekday.WeekdayCollectionPublicPath).borrow() 
		
			if let _collectionRef = collectionRef {
		
				let ids = _collectionRef.getIDs()
		
				if ids.contains(nftId) {
					let nftViewResolver = _collectionRef.borrowViewResolver(id: nftId)
				
					return nftViewResolver.resolveView(Type<MetadataViews.NFTView>())
				}
			}
		
			return nil
		}
	`)

	args := []cadence.Value{
		cadence.BytesToAddress(flow.HexToAddress("0x2a37a78609bba037").Bytes()),
		cadence.UInt64(1),
	}
	value, err := flowClient.ExecuteScriptAtLatestBlock(ctx, script, args)

	examples.Handle(err)
	fmt.Printf("\nValue: %s", value.String())
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

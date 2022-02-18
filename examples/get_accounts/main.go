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

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/examples"
	"google.golang.org/grpc"
)

func main() {
	prepareDemo()
	demo()
}

func demo() {
	ctx := context.Background()
	flowClient := examples.NewFlowClient()

	// get account from the latest block
	address := flow.HexToAddress("f8d6e0586b0a20c7")
	account, err := flowClient.GetAccount(ctx, address)
	printAccount(account, err)

	// get account from the block by height 0
	account, err = flowClient.GetAccountAtBlockHeight(ctx, address, 0)
	printAccount(account, err)
}

func printAccount(account *flow.Account, err error) {
	examples.Handle(err)

	fmt.Printf("\nAddress: %s", account.Address.String())
	fmt.Printf("\nBalance: %d", account.Balance)
	fmt.Printf("\nContracts: %d", len(account.Contracts))
	fmt.Printf("\nKeys: %d\n", len(account.Keys))
}

func prepareDemo() {
	flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
	examples.Handle(err)
	defer func() {
		err := flowClient.Close()
		if err != nil {
			panic(err)
		}
	}()

	examples.RandomAccount(flowClient)
}

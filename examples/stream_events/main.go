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

	"github.com/onflow/flow-go-sdk/access/grpc"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/examples"
)

func main() {
	demo()
}

func demo() {
	ctx := context.Background()
	flowClient, err := grpc.NewClient("access-003.devnet46.nodes.onflow.org:9000")
	examples.Handle(err)

	block, err := flowClient.GetLatestBlock(ctx, true)
	// block, err := flowClient.GetBlockByID(ctx, flow.HexToID("7582cc6e1bb5ca1784e309ca63013e9b7ecf34b74bf7fdb029aa0faa0deb7958err"))
	examples.Handle(err)
	fmt.Printf("Block Height: %d\n", block.Height)
	fmt.Printf("Block ID: %s\n", block.ID)

	data, errChan, initErr := flowClient.SubscribeEvents(ctx, block.ID, 0, flow.EventFilter{})
	examples.Handle(initErr)

	for {
		select {
		case <-ctx.Done():
			return
		case eventData, ok := <-data:
			if !ok {
				panic("data subscription closed")
			}
			fmt.Printf("~~~ Height: %d ~~~\n", eventData.Height)
			printEvents(eventData.Events)
		case err, ok := <-errChan:
			if !ok {
				panic("error channel subscription closed")
			}
			fmt.Printf("~~~ ERROR: %s ~~~\n", err.Error())
		}
	}

}

func printEvents(events []flow.Event) {
	for _, event := range events {
		fmt.Printf("\n\nType: %s", event.Type)
		fmt.Printf("\nValues: %v", event.Value)
		fmt.Printf("\nTransaction ID: %s", event.TransactionID)
	}
}

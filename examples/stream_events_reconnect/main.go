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

// This is an example of streaming events, and handling reconnect when errors are encountered on the stream.

func demo() {
	ctx := context.Background()
	flowClient, err := grpc.NewClient("access.testnet.nodes.onflow.org:9000")
	examples.Handle(err)

	header, err := flowClient.GetLatestBlockHeader(ctx, true)
	examples.Handle(err)
	fmt.Printf("Block Height: %d\n", header.Height)
	fmt.Printf("Block ID: %s\n", header.ID)

	data, errChan, initErr := flowClient.SubscribeEventsByBlockID(ctx, header.ID, flow.EventFilter{})
	examples.Handle(initErr)

	reconnect := func(height uint64) {
		fmt.Printf("Reconnecting at block %d\n", height)

		var err error
		flowClient, err = grpc.NewClient("access.testnet.nodes.onflow.org:9000")
		examples.Handle(err)

		data, errChan, err = flowClient.SubscribeEventsByBlockHeight(ctx, height, flow.EventFilter{})
		examples.Handle(err)
	}

	// track the most recently seen block height. we will use this when reconnecting
	lastHeight := header.Height
	for {
		select {
		case <-ctx.Done():
			return

		case eventData, ok := <-data:
			if !ok {
				if ctx.Err() != nil {
					return // graceful shutdown
				}
				// unexpected close
				reconnect(lastHeight + 1)
				continue
			}

			fmt.Printf("~~~ Height: %d ~~~\n", eventData.Height)
			printEvents(eventData.Events)

			lastHeight = eventData.Height

		case err, ok := <-errChan:
			if !ok {
				if ctx.Err() != nil {
					return // graceful shutdown
				}
				// unexpected close
				reconnect(lastHeight + 1)
				continue
			}

			fmt.Printf("~~~ ERROR: %s ~~~\n", err.Error())
			reconnect(lastHeight + 1)
			continue
		}
	}

}

func printEvents(events []flow.Event) {
	for _, event := range events {
		fmt.Printf("\nType: %s\n", event.Type)
		fmt.Printf("Values: %v\n", event.Value)
		fmt.Printf("Transaction ID: %s\n", event.TransactionID)
	}
}

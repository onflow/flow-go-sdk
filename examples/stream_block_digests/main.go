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
	flowClient, err := grpc.NewClient("access.testnet.nodes.onflow.org:9000")
	examples.Handle(err)

	header, err := flowClient.GetLatestBlockHeader(ctx, true)
	examples.Handle(err)
	fmt.Printf("Block Height: %d\n", header.Height)
	fmt.Printf("Block ID: %s\n", header.ID)

	blockDigestsChan, errChan, initErr := flowClient.SubscribeBlockDigestsFromStartBlockID(ctx, header.ID, flow.BlockStatusFinalized)
	examples.Handle(initErr)

	for {
		select {
		case <-ctx.Done():
			return
		case blockDigest, ok := <-blockDigestsChan:
			if !ok {
				panic("block digest channel is closed")
			}
			examples.Print(blockDigest)
		case err, ok := <-errChan:
			if !ok {
				panic("error channel is closed")
			}
			fmt.Printf("~~~ ERROR: %s ~~~\n", err.Error())
		}
	}
}

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

	"github.com/onflow/flow-go-sdk/access"

	"github.com/onflow/flow-go-sdk/access/http"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/examples"
)

func main() {
	demo()
}

func demo() {
	ctx := context.Background()
	var flowClient access.Client
	flowClient, err := http.NewClient(http.EmulatorHost)
	examples.Handle(err)

	examples.RandomAccount(flowClient)

	// get the latest sealed block
	isSealed := true
	latestBlock, err := flowClient.GetLatestBlock(ctx, isSealed)
	printBlock(latestBlock, err)

	// get the block by ID
	blockID := latestBlock.ID.String()
	blockByID, err := flowClient.GetBlockByID(ctx, flow.HexToID(blockID))
	printBlock(blockByID, err)

	// get block by height
	blockByHeight, err := flowClient.GetBlockByHeight(ctx, 0)
	printBlock(blockByHeight, err)
}

func printBlock(block *flow.Block, err error) {
	examples.Handle(err)

	fmt.Printf("\nID: %s\n", block.ID)
	fmt.Printf("height: %d\n", block.Height)
	fmt.Printf("timestamp: %s\n\n", block.Timestamp)
	fmt.Printf("Signatures: %s\n", block.Signatures)
}

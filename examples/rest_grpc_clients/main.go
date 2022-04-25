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

	"google.golang.org/grpc/credentials/insecure"

	"github.com/onflow/flow-go-sdk/access"
	"github.com/onflow/flow-go-sdk/access/grpc"
	"github.com/onflow/flow-go-sdk/access/http"
	"github.com/onflow/flow-go-sdk/examples"

	grpcOpts "google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	var flowClient access.Client

	// initialize a http emulator client
	flowClient, err := http.NewClient(http.EmulatorHost)
	examples.Handle(err)

	latestBlock, err := flowClient.GetLatestBlock(ctx, true)
	examples.Handle(err)

	fmt.Println("Block ID:", latestBlock.BlockHeader.ID.String())

	// initialize a gPRC emulator client
	flowClient, err = grpc.NewClient(grpc.EmulatorHost)
	examples.Handle(err)

	latestBlock, err = flowClient.GetLatestBlock(ctx, true)
	examples.Handle(err)

	fmt.Println("Block ID:", latestBlock.BlockHeader.ID.String())

	// initialize http specific client
	httpClient, err := http.NewBaseClient(http.EmulatorHost)
	examples.Handle(err)

	httpBlocks, err := httpClient.GetBlocksByHeights(
		ctx,
		http.HeightQuery{
			Heights: []uint64{http.SEALED},
		},
	)
	examples.Handle(err)

	fmt.Println("Block ID:", httpBlocks[0].ID.String())

	// initialize grpc specific client
	grpcClient, err := grpc.NewBaseClient(
		grpc.EmulatorHost,
		grpcOpts.WithTransportCredentials(insecure.NewCredentials()),
	)
	examples.Handle(err)

	grpcBlock, err := grpcClient.GetLatestBlockHeader(ctx, true)
	examples.Handle(err)

	fmt.Println("Block ID:", grpcBlock.ID.String())
}

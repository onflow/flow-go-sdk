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

	"github.com/onflow/flow-go-sdk/access/http"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/examples"
)

func main() {
	id := prepareDemo()
	demo(id)
}

func demo(exampleCollectionID flow.Identifier) {
	ctx := context.Background()
	flowClient := examples.NewFlowGRPCClient()

	// get collection by ID
	collection, err := flowClient.GetCollection(ctx, exampleCollectionID)
	printCollection(collection, err)
}

func printCollection(collection *flow.LightCollection, err error) {
	examples.Handle(err)

	fmt.Printf("\nID: %s", collection.ID().String())
	fmt.Printf("\nTransactions: %s", collection.TransactionIDs)
}

func prepareDemo() flow.Identifier {
	flowClient, err := http.NewClient(http.EmulatorHost)
	examples.Handle(err)

	examples.RandomAccount(flowClient)

	block, err := flowClient.GetBlockByHeight(context.Background(), 1)
	examples.Handle(err)

	return block.CollectionGuarantees[0].CollectionID
}

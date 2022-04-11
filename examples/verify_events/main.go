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

	"github.com/onflow/flow-go-sdk/client/http"

	"github.com/onflow/flow-go-sdk"

	"github.com/onflow/flow-go-sdk/examples"
)

func main() {
	VerifyEventsDemo()
}

// VerifyEventsDemo showcases how clients can verify events they received.
// In mature Flow we provision different methods to verify the subset of events queried, however
// for now the only verification available is the hash of all the events emitted in a Chunk/Collection.
// This hash is verified by Verification Nodes so it's correctness is checked.
// For users, its possible to get all the events for chunk, calculate and compare the resulting hash
func VerifyEventsDemo() {
	ctx := context.Background()
	flowClient, err := http.NewDefaultEmulatorClient(false)
	examples.Handle(err)

	latestBlockHeader, err := flowClient.GetLatestBlockHeader(ctx, true)
	examples.Handle(err)

	block, err := flowClient.GetBlockByID(ctx, latestBlockHeader.ID)
	examples.Handle(err)

	executionResult, err := flowClient.GetExecutionResultForBlockID(ctx, block.ID)
	examples.Handle(err)

	totalEvents := 0

	for i, colGuarantee := range block.CollectionGuarantees {
		collection, err := flowClient.GetCollection(ctx, colGuarantee.CollectionID)
		examples.Handle(err)

		events := make([]flow.Event, 0)

		// only way to get all the events regardless of the type is querying by transaction
		for _, txID := range collection.TransactionIDs {
			transactionResult, err := flowClient.GetTransactionResult(ctx, txID)
			examples.Handle(err)

			for _, event := range transactionResult.Events {
				events = append(events, event)
			}
		}

		calculatedEventsHash, err := flow.CalculateEventsHash(events)
		examples.Handle(err)

		// there should be more chunks (one last is system chunk), so we don't check range here
		chunk := executionResult.Chunks[i]

		if !calculatedEventsHash.Equal(chunk.EventCollection) {
			examples.Handle(fmt.Errorf("events hash mismatch, expected %s, calculated %s", chunk.EventCollection, calculatedEventsHash))
		}

		totalEvents += len(events)
	}

	fmt.Printf("Events verified for block %s, total %d events\n", block.ID, totalEvents)
}

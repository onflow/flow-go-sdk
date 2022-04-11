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
	demo()
}

func demo() {
	ctx := context.Background()
	flowClient, err := http.NewDefaultEmulatorClient(false)
	examples.Handle(err)

	txID := examples.RandomTransaction(flowClient).ID()

	tx, err := flowClient.GetTransaction(ctx, txID)
	printTransaction(tx, err)

	txr, err := flowClient.GetTransactionResult(ctx, txID)
	printTransactionResult(txr, err)
}

func printTransaction(tx *flow.Transaction, err error) {
	examples.Handle(err)

	fmt.Printf("\nID: %s", tx.ID().String())
	fmt.Printf("\nPayer: %s", tx.Payer.String())
	fmt.Printf("\nProposer: %s", tx.ProposalKey.Address.String())
	fmt.Printf("\nAuthorizers: %s", tx.Authorizers)
}

func printTransactionResult(txr *flow.TransactionResult, err error) {
	examples.Handle(err)

	fmt.Printf("\nStatus: %s", txr.Status.String())
	fmt.Printf("\nError: %v", txr.Error)
}

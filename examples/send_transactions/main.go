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

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/examples"
)

func main() {
	tx := prepareDemo()
	demo(tx)
}

func demo(tx *flow.Transaction) {
	ctx := context.Background()
	flowClient := examples.NewFlowClient()

	err := flowClient.SendTransaction(ctx, *tx)
	if err != nil {
		panic(err)
	}
}

func prepareDemo() *flow.Transaction {
	flowClient := examples.NewFlowClient()
	defer func() {
		err := flowClient.Close()
		if err != nil {
			panic(err)
		}
	}()

	serviceAcctAddr, serviceAcctKey, serviceSigner := examples.ServiceAccount(flowClient)

	tx := flow.NewTransaction().
		SetPayer(serviceAcctAddr).
		SetProposalKey(serviceAcctAddr, serviceAcctKey.Index, serviceAcctKey.SequenceNumber).
		SetScript([]byte("transaction {}")).
		SetReferenceBlockID(examples.GetReferenceBlockId(flowClient))

	err := tx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, serviceSigner)
	examples.Handle(err)

	return tx
}

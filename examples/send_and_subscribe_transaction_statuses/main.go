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
	flowClient, err := grpc.NewClient(grpc.EmulatorHost)
	examples.Handle(err)

	serviceAcctAddr, serviceAcctKey, serviceSigner := examples.ServiceAccount(flowClient)

	tx := flow.NewTransaction().
		SetPayer(serviceAcctAddr).
		SetProposalKey(serviceAcctAddr, serviceAcctKey.Index, serviceAcctKey.SequenceNumber).
		SetScript([]byte(`
			transaction {
  				prepare(acc: &Account) {}
				execute {
    				log("test")
  				}
			}
		`)).
		AddAuthorizer(serviceAcctAddr).
		SetReferenceBlockID(examples.GetReferenceBlockId(flowClient))

	err = tx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, serviceSigner)
	examples.Handle(err)

	txResultChan, errChan, initErr := flowClient.SendAndSubscribeTransactionStatuses(ctx, *tx)
	examples.Handle(initErr)

	select {
	case <-ctx.Done():
		return
	case txResult, ok := <-txResultChan:
		if !ok {
			panic("transaction result channel is closed")
		}
		examples.Print(txResult)
	case err, ok := <-errChan:
		if !ok {
			panic("error channel is closed")
		}
		fmt.Printf("~~~ ERROR: %s ~~~\n", err.Error())
	}
}

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
	"log"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access/grpc"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/examples"
)

func main() {
	demo()
}

func demo() {
	ctx := context.Background()
	flowClient, err := grpc.NewClient("access.testnet.nodes.onflow.org:9000")
	examples.Handle(err)

	signerIndex := uint32(0)

	signerPublicAddress := flow.HexToAddress("YOUR_ACCOUNT_ADDRESS")
	signerAccount, err := flowClient.GetAccount(ctx, signerPublicAddress)
	if err != nil {
		log.Fatalf("Failed to get account: %v", err)
	}
	seqNumber := signerAccount.Keys[0].SequenceNumber

	privateKeyHex := "YOUR_PRIVATE_KEY"
	privateKey, err := crypto.DecodePrivateKeyHex(crypto.ECDSA_P256, privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}

	// Create crypto signer
	signer, err := crypto.NewInMemorySigner(privateKey, crypto.SHA3_256)
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}

	tx := flow.NewTransaction().
		SetPayer(signerPublicAddress).
		SetProposalKey(signerPublicAddress, signerIndex, seqNumber).
		SetScript([]byte(`
			transaction {
  				prepare(acc: &Account) {}
				execute {
    				log("test")
  				}
			}
		`)).
		AddAuthorizer(signerPublicAddress).
		SetReferenceBlockID(examples.GetReferenceBlockId(flowClient))

	err = tx.SignEnvelope(signerPublicAddress, signerIndex, signer)
	examples.Handle(err)

	txResultChan, errChan, initErr := flowClient.SendAndSubscribeTransactionStatuses(ctx, *tx)
	examples.Handle(initErr)

	for {
		select {
		case <-ctx.Done():
			return
		case txResult, ok := <-txResultChan:
			if !ok {
				examples.Print("transaction result channel is closed")
				return
			}
			examples.Print(txResult)
		case err := <-errChan:
			if err != nil {
				examples.Print(fmt.Errorf("~~~ ERROR: %w ~~~\n", err))
			}
			return
		}
	}
}

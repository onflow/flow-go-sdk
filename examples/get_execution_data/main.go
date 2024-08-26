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
	flowClient, err := grpc.NewClient("access-003.devnet46.nodes.onflow.org:9000")
	examples.Handle(err)

	// block, err := flowClient.GetLatestBlock(ctx, true)
	block, err := flowClient.GetBlockByID(ctx, flow.HexToID("7582cc6e1bb5ca1784e309ca63013e9b7ecf34b74bf7fdb029aa0faa0deb7958err"))
	examples.Handle(err)
	fmt.Printf("Block Height: %d\n", block.Height)
	fmt.Printf("Block ID: %s\n", block.ID)

	data, err := flowClient.GetExecutionDataByBlockID(ctx, block.ID)
	examples.Handle(err)
	printExecutionData(data)
}

func printExecutionData(ed *flow.ExecutionData) {
	for chunkNo, chunk := range ed.ChunkExecutionData {
		fmt.Printf("-- Chunk %d/%d --\n", chunkNo+1, len(ed.ChunkExecutionData))
		fmt.Printf("Transactions: %d\n", len(chunk.Transactions))
		for txNo, tx := range chunk.Transactions {
			fmt.Printf("Transaction %d/%d:\n", txNo+1, len(chunk.Transactions))
			printTransaction(tx)
		}
		fmt.Printf("Events: %d\n", len(chunk.Events))
		for eventNo, event := range chunk.Events {
			fmt.Printf("Event %d/%d: %s\n", eventNo+1, len(chunk.Events), event.Type)
		}
		if chunk.TrieUpdate != nil {
			modifiedAccounts := extractModifiedAccounts(chunk.TrieUpdate)
			fmt.Printf("Modified Accounts: %d\n", len(modifiedAccounts))
			for i, acc := range modifiedAccounts {
				fmt.Printf("Account %d/%d: %s\n", i+1, len(modifiedAccounts), acc.Hex())
			}
			fmt.Printf("TrieUpdate RootHash: %s\n", flow.BytesToHash(chunk.TrieUpdate.RootHash).Hex())
			fmt.Printf("TrieUpdate Paths:\n")
			for pathNo, path := range chunk.TrieUpdate.Paths {
				fmt.Printf("Path %d/%d: %s", pathNo+1, len(chunk.TrieUpdate.Paths), flow.BytesToHash(path))
			}
			fmt.Printf("TrieUpdate Payloads:\n")
			for payloadNo, payload := range chunk.TrieUpdate.Payloads {
				fmt.Printf("Payload %d/%d:\n", payloadNo+1, len(chunk.TrieUpdate.Payloads))
				fmt.Printf("Value: %x\n", string(payload.Value))
				for kpi, kp := range payload.KeyPart {
					fmt.Printf("KeyPart[%d].Type: %d\n", kpi, kp.Type)
					fmt.Printf("KeyPart[%d].Value: %x\n", kpi, kp.Value)
				}

			}
		} else {
			fmt.Printf("\nNo TrieUpdate")
		}
	}
}

func printTransaction(tx *flow.Transaction) {
	fmt.Printf("\nID: %s", tx.ID().String())
	fmt.Printf("\nPayer: %s", tx.Payer.String())
	fmt.Printf("\nProposer: %s", tx.ProposalKey.Address.String())
	fmt.Printf("\nAuthorizers: %s", tx.Authorizers)
}

func extractModifiedAccounts(update *flow.TrieUpdate) []flow.Address {

	accounts := map[flow.Address]struct{}{}

	for _, payload := range update.Payloads {
		address := flow.BytesToAddress(payload.KeyPart[0].Value)
		accounts[address] = struct{}{}
	}

	addresses := make([]flow.Address, 0, len(accounts))
	for address := range accounts {
		addresses = append(addresses, address)
	}

	return addresses
}

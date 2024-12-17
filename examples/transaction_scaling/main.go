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
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/crypto/hash"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access/grpc"
	"github.com/onflow/flow-go-sdk/crypto"

	"github.com/onflow/flow-go-sdk/examples"
)

/*
 * This program demonstrates how to run several transactions in parallel from the same account using multiple proposal keys.
 *
 * We are using Testnet since we want to experiment with real transaction throughput and network conditions.
 * The test starts by deploying a simple counter contract and adding multiple proposal keys to the account.
 * Then, it will execute several transactions in parallel, each using a different proposal key.
 * Each transaction will increase the counter by 1.
 *
 * IMPORTANT: This demo requires a new testnet account. Use `flow keys generate` to generate a new key pair.
 * Create and fund your Testnet account here: https://testnet-faucet.onflow.org/create-account
 *
 * Finally, set the private key and account address below before running the program.
 */

const PRIVATE_KEY = "abc123"
const ACCOUNT_ADDRESS = "0x1234567890"

const numProposalKeys = 420 // number of proposal keys to use, also number of workers (1 worker per key)
const numTxs = 420          // number of total transactions to execute

const contractCode = `
	access(all) contract Counter {

		access(self) var count: Int

		init() {
			self.count = 0
		}

		access(all) fun increase() {
			self.count = self.count + 1
		}

		access(all) view fun getCount(): Int {
			return self.count
		}
	}
`

func main() {
	// set up context and flow client
	ctx := context.Background()
	flowClient, err := grpc.NewClient(grpc.TestnetHost)
	examples.Handle(err)

	// initialize account and signer
	// this will deploy the counter contract and add the required number of proposal keys to the account
	account, signer, err := InitAccount(ctx, flowClient, PRIVATE_KEY, ACCOUNT_ADDRESS, numProposalKeys)
	examples.Handle(err)

	// print the current counter value
	startCounterValue, err := GetCounter(ctx, flowClient, account)
	examples.Handle(err)
	fmt.Printf("Initial Counter: %d\n", startCounterValue)

	// populate the job channel with the number of transactions to execute
	txChan := make(chan int, numTxs)
	for i := 0; i < numTxs; i++ {
		txChan <- i
	}

	startTime := time.Now()

	var wg sync.WaitGroup
	// start the workers
	for i := 0; i < numProposalKeys; i++ {
		wg.Add(1)

		// worker code
		// this will run in parallel for each proposal key
		go func(keyIndex int) {
			defer wg.Done()

			// consume the job channel
			for range txChan {
				fmt.Printf("[Worker %d] executing transaction\n", keyIndex)

				// execute the transaction
				err := IncreaseCounter(ctx, flowClient, account, signer, keyIndex)
				if err != nil {
					fmt.Printf("[Worker %d] Error: %v\n", keyIndex, err)
					return
				}
			}
		}(i)
	}

	close(txChan)

	// wait for all workers to finish
	wg.Wait()

	finalCounter, err := GetCounter(ctx, flowClient, account)
	examples.Handle(err)
	fmt.Printf("Final Counter: %d\n", finalCounter)

	if finalCounter-startCounterValue != numTxs {
		fmt.Printf("❌ Error: %d transactions executed, expected %d\n", numTxs, finalCounter-startCounterValue)
	} else {
		fmt.Printf("✅ Done! %d transactions executed in %s\n", numTxs, time.Since(startTime))
	}
}

// Initialize the proposer keys and deploy the counter contract
func InitAccount(ctx context.Context, flowClient *grpc.Client, privateKey string, accountAddress string, numKeys int) (*flow.Account, crypto.Signer, error) {
	pk, err := crypto.DecodePrivateKeyHex(crypto.ECDSA_P256, privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode private key %s", err)
	}
	account, err := flowClient.GetAccount(ctx, flow.HexToAddress(accountAddress))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get account %s", err)
	}
	signer, err := crypto.NewInMemorySigner(pk, hash.SHA3_256)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create in mem signer: %s", err)
	}

	if _, exists := account.Contracts["Counter"]; exists {
		return nil, nil, errors.New("contract already exists")
		// uncomment below to run this test on the same account again
		// return account, signer, nil
	}

	fmt.Println("Initializing account ...")

	initScript := `
		transaction(code: String, numKeys: Int) {

			prepare(signer: auth(AddContract, AddKey) &Account) {
				// deploy the contract
				signer.contracts.add(name: "Counter", code: code.decodeHex())

				// copy the main key with 0 weight multiple times
				// to create the required number of keys
				let key = signer.keys.get(keyIndex: 0)!
				var count: Int = 0
				while count < numKeys {
					signer.keys.add(
						publicKey: key.publicKey,
						hashAlgorithm: key.hashAlgorithm,
						weight: 0.0
					)
					count = count + 1
				}
			}
		}
	`

	codeHex := hex.EncodeToString([]byte(contractCode))

	deployContractTx := flow.NewTransaction().
		SetScript([]byte(initScript)).
		AddRawArgument(jsoncdc.MustEncode(cadence.String(codeHex))).
		AddRawArgument(jsoncdc.MustEncode(cadence.NewInt(numKeys))).
		AddAuthorizer(account.Address)

	deployContractTx.SetProposalKey(
		account.Address,
		account.Keys[0].Index,
		account.Keys[0].SequenceNumber,
	)

	err = RunTransaction(ctx, flowClient, account, signer, deployContractTx)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("✅ Account initialized")

	return account, signer, nil
}

// Run a transaction and wait for it to be sealed. Note that this function does not set the proposal key.
func RunTransaction(ctx context.Context, flowClient *grpc.Client, account *flow.Account, signer crypto.Signer, tx *flow.Transaction) error {
	latestBlock, err := flowClient.GetLatestBlock(ctx, true)
	if err != nil {
		return err
	}
	tx.SetReferenceBlockID(latestBlock.ID)
	tx.SetPayer(account.Address)

	err = SignTransaction(ctx, flowClient, account, signer, tx)
	if err != nil {
		return err
	}

	err = flowClient.SendTransaction(ctx, *tx)
	if err != nil {
		return err
	}

	txRes := examples.WaitForSeal(ctx, flowClient, tx.ID())
	if txRes.Error != nil {
		return txRes.Error
	}

	return nil
}

// Local signer is not thread safe, so we need to lock the mutex during the signing operation
var signingMutex sync.Mutex

func SignTransaction(ctx context.Context, flowClient *grpc.Client, account *flow.Account, signer crypto.Signer, tx *flow.Transaction) error {
	// Lock the mutex during the signing operation
	signingMutex.Lock()
	defer signingMutex.Unlock()

	// need an extra payload signature if the key index is not 0
	if tx.ProposalKey.KeyIndex != 0 {
		err := tx.SignPayload(account.Address, tx.ProposalKey.KeyIndex, signer)
		if err != nil {
			return err
		}
	}

	// sign the envelope with the full weight key
	err := tx.SignEnvelope(account.Address, account.Keys[0].Index, signer)
	if err != nil {
		return err
	}

	return nil
}

// Get the current counter value
func GetCounter(ctx context.Context, flowClient *grpc.Client, account *flow.Account) (int, error) {
	script := []byte(fmt.Sprintf(`
		import Counter from 0x%s

		access(all) fun main(): Int {
			return Counter.getCount()
		}

	`, account.Address.String()))
	value, err := flowClient.ExecuteScriptAtLatestBlock(ctx, script, nil)
	if err != nil {
		return -1, err
	}

	num, err := strconv.Atoi(value.String())
	if err != nil {
		return -1, err
	}

	return num, nil
}

// Increase the counter by 1 by running a transaction using the given proposal key
func IncreaseCounter(ctx context.Context, flowClient *grpc.Client, account *flow.Account, signer crypto.Signer, proposalKeyIndex int) error {
	script := []byte(fmt.Sprintf(`
		import Counter from 0x%s

		transaction() {
			prepare(signer: &Account) {
				Counter.increase()
			}
		}

	`, account.Address.String()))

	tx := flow.NewTransaction().
		SetScript(script).
		AddAuthorizer(account.Address)

	// get the latest account state including the sequence number
	account, err := flowClient.GetAccount(ctx, flow.HexToAddress(account.Address.String()))
	if err != nil {
		return err
	}
	tx.SetProposalKey(
		account.Address,
		account.Keys[proposalKeyIndex].Index,
		account.Keys[proposalKeyIndex].SequenceNumber,
	)

	return RunTransaction(ctx, flowClient, account, signer, tx)
}

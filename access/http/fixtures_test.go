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

package http

import (
	"encoding/base64"
	"fmt"

	"github.com/onflow/flow-go-sdk/access/http/models"
	"github.com/onflow/flow-go-sdk/test"
)

func contractFlowFixture() (string, string) {
	return "HelloWorld", base64.StdEncoding.EncodeToString([]byte(`
		contract HelloWorld {}
	`))
}

func accountFlowFixture() models.Account {
	name, source := contractFlowFixture()
	return models.Account{
		Address: test.AddressGenerator().New().String(),
		Balance: "10",
		Keys: []models.AccountPublicKey{
			accountKeyFlowFixture(),
		},
		Contracts:  map[string]string{name: source},
		Expandable: nil,
		Links:      nil,
	}
}

func accountKeyFlowFixture() models.AccountPublicKey {
	key := test.AccountKeyGenerator().New()
	sigAlgo := models.SigningAlgorithm(key.SigAlgo.String())
	hashAlgo := models.HashingAlgorithm(key.HashAlgo.String())

	return models.AccountPublicKey{
		Index:            "0",
		PublicKey:        key.PublicKey.String(),
		SigningAlgorithm: &sigAlgo,
		HashingAlgorithm: &hashAlgo,
		SequenceNumber:   "0",
		Weight:           "1000",
		Revoked:          false,
	}
}

func blockFlowFixture() models.Block {
	block := test.BlockGenerator().New()

	return models.Block{
		Header: &models.BlockHeader{
			Id:                   block.ID.String(),
			ParentId:             block.ParentID.String(),
			Height:               fmt.Sprintf("%d", block.Height),
			Timestamp:            block.Timestamp,
			ParentVoterSignature: base64.StdEncoding.EncodeToString([]byte("test")),
		},
		Payload: &models.BlockPayload{
			CollectionGuarantees: []models.CollectionGuarantee{{
				CollectionId: block.CollectionGuarantees[0].CollectionID.String(),
			}},
			BlockSeals: []models.BlockSeal{{
				BlockId:                      block.Seals[0].BlockID.String(),
				ResultId:                     block.Seals[0].ExecutionReceiptID.String(),
				FinalState:                   "",
				AggregatedApprovalSignatures: nil,
			}},
		},
		ExecutionResult: nil,
	}
}

func collectionFlowFixture() models.Collection {
	collection := test.CollectionGenerator().New()

	return models.Collection{
		Id: collection.ID().String(),
		Transactions: []models.Transaction{{
			Id: collection.TransactionIDs[0].String(),
		}},
	}
}

func transactionFlowFixture() models.Transaction {
	tx := test.TransactionGenerator().New()

	args := make([]string, len(tx.Arguments))
	for i, a := range tx.Arguments {
		args[i] = base64.StdEncoding.EncodeToString(a)
	}

	auths := make([]string, len(tx.Authorizers))
	for i, a := range tx.Authorizers {
		auths[i] = a.String()
	}

	return models.Transaction{
		Id:               tx.ID().String(),
		Script:           base64.StdEncoding.EncodeToString(tx.Script),
		Arguments:        args,
		ReferenceBlockId: tx.ReferenceBlockID.String(),
		GasLimit:         fmt.Sprintf("%d", tx.GasLimit),
		Payer:            tx.Payer.String(),
		ProposalKey: &models.ProposalKey{
			Address:        tx.ProposalKey.Address.String(),
			KeyIndex:       fmt.Sprintf("%d", tx.ProposalKey.KeyIndex),
			SequenceNumber: fmt.Sprintf("%d", tx.ProposalKey.SequenceNumber),
		},
		Authorizers: auths,
		PayloadSignatures: []models.TransactionSignature{{
			Address:   tx.PayloadSignatures[0].Address.String(),
			KeyIndex:  fmt.Sprintf("%d", tx.PayloadSignatures[0].KeyIndex),
			Signature: base64.StdEncoding.EncodeToString(tx.PayloadSignatures[0].Signature),
		}},
		EnvelopeSignatures: []models.TransactionSignature{{
			Address:   tx.EnvelopeSignatures[0].Address.String(),
			KeyIndex:  fmt.Sprintf("%d", tx.EnvelopeSignatures[0].KeyIndex),
			Signature: base64.StdEncoding.EncodeToString(tx.EnvelopeSignatures[0].Signature),
		}},
	}
}

func transactionResultFlowFixture() models.TransactionResult {
	txr := test.TransactionResultGenerator().New()
	status := models.SEALED_TransactionStatus

	return models.TransactionResult{
		Status:       &status,
		StatusCode:   0,
		ErrorMessage: txr.Error.Error(),
		Events: []models.Event{{
			Type_:            txr.Events[0].Type,
			TransactionId:    txr.Events[0].TransactionID.String(),
			TransactionIndex: fmt.Sprintf("%d", txr.Events[0].TransactionIndex),
			EventIndex:       fmt.Sprintf("%d", txr.Events[0].EventIndex),
			Payload:          base64.StdEncoding.EncodeToString(txr.Events[0].Payload),
		}},
	}
}

func eventsFlowFixture(n int) []models.Event {
	events := make([]models.Event, n)

	for i := 0; i < n; i++ {
		e := test.EventGenerator().New()
		events[i] = models.Event{
			Type_:            e.Type,
			TransactionId:    e.TransactionID.String(),
			TransactionIndex: fmt.Sprintf("%d", e.TransactionIndex),
			EventIndex:       fmt.Sprintf("%d", e.EventIndex),
			Payload:          base64.StdEncoding.EncodeToString(e.Payload),
		}
	}

	return events
}

func blockEventsFlowFixture() models.BlockEvents {
	block := test.BlockGenerator().New()
	events := eventsFlowFixture(4)

	return models.BlockEvents{
		BlockId:        block.ID.String(),
		BlockHeight:    fmt.Sprintf("%d", block.Height),
		BlockTimestamp: block.Timestamp,
		Events:         events,
	}
}

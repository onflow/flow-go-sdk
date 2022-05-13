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
	"testing"

	"github.com/onflow/flow-go-sdk"

	"github.com/stretchr/testify/assert"
)

func Test_ConvertBlock(t *testing.T) {
	httpBlock := BlockHTTP()

	block, err := toBlock(&httpBlock)

	assert.NoError(t, err)
	assert.Equal(t, block.ID.String(), httpBlock.Header.Id)
	assert.Equal(t, fmt.Sprintf("%d", block.Height), httpBlock.Header.Height)
	assert.Equal(t, block.Timestamp, httpBlock.Header.Timestamp)
	assert.Len(t, block.BlockPayload.Seals, len(httpBlock.Payload.BlockSeals))
	assert.Equal(t, block.ParentID.String(), httpBlock.Header.ParentId)
	assert.Len(t, block.BlockPayload.CollectionGuarantees, len(httpBlock.Payload.CollectionGuarantees))
	assert.Equal(t, block.BlockPayload.CollectionGuarantees[0].CollectionID.String(), httpBlock.Payload.CollectionGuarantees[0].CollectionId)
}

func Test_ConvertAccount(t *testing.T) {
	httpAccount := AccountHTTP()
	contractName, contractCode := ContractHTTP()

	account, err := toAccount(&httpAccount)

	assert.NoError(t, err)
	assert.Equal(t, account.Address.String(), httpAccount.Address)
	assert.Len(t, account.Keys, len(httpAccount.Keys))
	assert.Equal(t, account.Keys[0].PublicKey.String(), httpAccount.Keys[0].PublicKey)
	code, _ := base64.StdEncoding.DecodeString(contractCode)
	assert.Equal(t, account.Contracts[contractName], code)
	assert.Equal(t, fmt.Sprintf("%d", account.Balance), httpAccount.Balance)
}

func Test_ConvertCollection(t *testing.T) {
	httpColl := CollectionHTTP()

	collection := toCollection(&httpColl)

	assert.Len(t, collection.TransactionIDs, len(httpColl.Transactions))
	assert.Equal(t, collection.TransactionIDs[0].String(), httpColl.Transactions[0].Id)
}

func Test_ConvertTransaction(t *testing.T) {
	httpTx := TransactionHTTP()
	script, _ := base64.StdEncoding.DecodeString(httpTx.Script)

	tx, err := toTransaction(&httpTx)

	auths := make([]string, len(tx.Authorizers))
	for i, a := range tx.Authorizers {
		auths[i] = a.String()
	}

	assert.NoError(t, err)
	assert.Equal(t, tx.ProposalKey.Address.String(), httpTx.ProposalKey.Address)
	assert.Equal(t, fmt.Sprintf("%d", tx.ProposalKey.KeyIndex), httpTx.ProposalKey.KeyIndex)
	assert.Equal(t, fmt.Sprintf("%d", tx.ProposalKey.SequenceNumber), httpTx.ProposalKey.SequenceNumber)
	assert.Equal(t, tx.Payer.String(), httpTx.Payer)
	assert.Equal(t, tx.Script, script)
	assert.Len(t, tx.Arguments, len(httpTx.Arguments))
	assert.Equal(t, fmt.Sprintf("%d", tx.GasLimit), httpTx.GasLimit)
	assert.EqualValues(t, auths, httpTx.Authorizers)
	assert.Equal(t, tx.PayloadSignatures[0].Address.String(), httpTx.PayloadSignatures[0].Address)
	assert.Equal(t, fmt.Sprintf("%d", tx.PayloadSignatures[0].KeyIndex), httpTx.PayloadSignatures[0].KeyIndex)
	sig, err := base64.StdEncoding.DecodeString(httpTx.PayloadSignatures[0].Signature)
	assert.NoError(t, err)
	assert.Equal(t, tx.PayloadSignatures[0].Signature, sig)
	assert.Equal(t, tx.EnvelopeSignatures[0].Address.String(), httpTx.EnvelopeSignatures[0].Address)
	assert.Equal(t, fmt.Sprintf("%d", tx.EnvelopeSignatures[0].KeyIndex), httpTx.EnvelopeSignatures[0].KeyIndex)
	sig, err = base64.StdEncoding.DecodeString(httpTx.EnvelopeSignatures[0].Signature)
	assert.NoError(t, err)
	assert.Equal(t, tx.EnvelopeSignatures[0].Signature, sig)
}

func Test_ConvertTransactionResult(t *testing.T) {
	httpTxr := TransactionResultHTTP()
	txr, err := toTransactionResult(&httpTxr)

	assert.NoError(t, err)
	assert.Equal(t, txr.Status, flow.TransactionStatusSealed)
	assert.Equal(t, txr.Error.Error(), httpTxr.ErrorMessage)
	assert.Len(t, txr.Events, len(httpTxr.Events))
	payload, err := base64.StdEncoding.DecodeString(httpTxr.Events[0].Payload)
	assert.NoError(t, err)
	assert.Equal(t, txr.Events[0].Payload, payload)
	assert.Equal(t, txr.Events[0].TransactionID.String(), httpTxr.Events[0].TransactionId)
	assert.Equal(t, fmt.Sprintf("%d", txr.Events[0].TransactionIndex), httpTxr.Events[0].TransactionIndex)
}

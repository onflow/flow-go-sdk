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

package grpc

import (
	"errors"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/onflow/cadence"
	"github.com/onflow/cadence/encoding/ccf"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow/protobuf/go/flow/access"
	"github.com/onflow/flow/protobuf/go/flow/entities"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
)

var errEmptyMessage = errors.New("protobuf message is empty")

func accountToMessage(a flow.Account) *entities.Account {
	accountKeys := make([]*entities.AccountKey, len(a.Keys))
	for i, key := range a.Keys {
		accountKeys[i] = accountKeyToMessage(key)
	}

	return &entities.Account{
		Address:   a.Address.Bytes(),
		Balance:   a.Balance,
		Code:      a.Code,
		Keys:      accountKeys,
		Contracts: a.Contracts,
	}
}

func messageToAccount(m *entities.Account) (flow.Account, error) {
	if m == nil {
		return flow.Account{}, errEmptyMessage
	}

	accountKeys := make([]*flow.AccountKey, len(m.GetKeys()))
	for i, key := range m.GetKeys() {
		accountKey, err := messageToAccountKey(key)
		if err != nil {
			return flow.Account{}, err
		}

		accountKeys[i] = accountKey
	}

	return flow.Account{
		Address:   flow.BytesToAddress(m.GetAddress()),
		Balance:   m.GetBalance(),
		Code:      m.GetCode(),
		Keys:      accountKeys,
		Contracts: m.GetContracts(),
	}, nil
}

func accountKeyToMessage(a *flow.AccountKey) *entities.AccountKey {
	return &entities.AccountKey{
		Index:          uint32(a.Index),
		PublicKey:      a.PublicKey.Encode(),
		SignAlgo:       uint32(a.SigAlgo),
		HashAlgo:       uint32(a.HashAlgo),
		Weight:         uint32(a.Weight),
		SequenceNumber: uint32(a.SequenceNumber),
		Revoked:        a.Revoked,
	}
}

func messageToAccountKey(m *entities.AccountKey) (*flow.AccountKey, error) {
	if m == nil {
		return nil, errEmptyMessage
	}

	sigAlgo := crypto.SignatureAlgorithm(m.GetSignAlgo())
	hashAlgo := crypto.HashAlgorithm(m.GetHashAlgo())

	publicKey, err := crypto.DecodePublicKey(sigAlgo, m.GetPublicKey())
	if err != nil {
		return nil, err
	}

	return &flow.AccountKey{
		Index:          int(m.GetIndex()),
		PublicKey:      publicKey,
		SigAlgo:        sigAlgo,
		HashAlgo:       hashAlgo,
		Weight:         int(m.GetWeight()),
		SequenceNumber: uint64(m.GetSequenceNumber()),
		Revoked:        m.GetRevoked(),
	}, nil
}

func blockToMessage(b flow.Block) (*entities.Block, error) {

	t := timestamppb.New(b.BlockHeader.Timestamp)

	return &entities.Block{
		Id:                   b.BlockHeader.ID.Bytes(),
		ParentId:             b.BlockHeader.ParentID.Bytes(),
		Height:               b.BlockHeader.Height,
		Timestamp:            t,
		CollectionGuarantees: collectionGuaranteesToMessages(b.BlockPayload.CollectionGuarantees),
		BlockSeals:           blockSealsToMessages(b.BlockPayload.Seals),
	}, nil
}

func messageToBlock(m *entities.Block) (flow.Block, error) {
	var timestamp time.Time
	var err error

	if m.GetTimestamp() != nil {
		timestamp = m.GetTimestamp().AsTime()
	}

	header := &flow.BlockHeader{
		ID:        flow.HashToID(m.GetId()),
		ParentID:  flow.HashToID(m.GetParentId()),
		Height:    m.GetHeight(),
		Timestamp: timestamp,
	}

	guarantees, err := messagesToCollectionGuarantees(m.GetCollectionGuarantees())
	if err != nil {
		return flow.Block{}, err
	}

	seals, err := messagesToBlockSeals(m.GetBlockSeals())
	if err != nil {
		return flow.Block{}, err
	}

	payload := &flow.BlockPayload{
		CollectionGuarantees: guarantees,
		Seals:                seals,
	}

	return flow.Block{
		BlockHeader:  *header,
		BlockPayload: *payload,
	}, nil
}

func blockHeaderToMessage(b flow.BlockHeader) (*entities.BlockHeader, error) {
	t := timestamppb.New(b.Timestamp)

	return &entities.BlockHeader{
		Id:        b.ID.Bytes(),
		ParentId:  b.ParentID.Bytes(),
		Height:    b.Height,
		Timestamp: t,
	}, nil
}

func messageToBlockHeader(m *entities.BlockHeader) (flow.BlockHeader, error) {
	if m == nil {
		return flow.BlockHeader{}, errEmptyMessage
	}

	var timestamp time.Time

	if m.GetTimestamp() != nil {
		timestamp = m.GetTimestamp().AsTime()
	}

	return flow.BlockHeader{
		ID:        flow.HashToID(m.GetId()),
		ParentID:  flow.HashToID(m.GetParentId()),
		Height:    m.GetHeight(),
		Timestamp: timestamp,
	}, nil
}

func cadenceValueToMessage(value cadence.Value) ([]byte, error) {
	b, err := jsoncdc.Encode(value)
	if err != nil {
		return nil, fmt.Errorf("jsoncdc convert: %w", err)
	}

	return b, nil
}

func cadenceValuesToMessages(values []cadence.Value) ([][]byte, error) {
	msgs := make([][]byte, len(values))
	for i, val := range values {
		msg, err := cadenceValueToMessage(val)
		if err != nil {
			return nil, err
		}
		msgs[i] = msg
	}
	return msgs, nil
}

func messageToCadenceValue(m []byte, options []jsoncdc.Option) (cadence.Value, error) {
	if ccf.HasMsgPrefix(m) {
		// modern Access nodes support encoding events in CCF format
		v, err := ccf.Decode(nil, m)
		if err != nil {
			return nil, fmt.Errorf("ccf convert: %w", err)
		}
		return v, nil
	}

	v, err := jsoncdc.Decode(nil, m, options...)
	if err != nil {
		return nil, fmt.Errorf("jsoncdc convert: %w", err)
	}

	return v, nil
}

func collectionToMessage(c flow.Collection) *entities.Collection {
	transactionIDMessages := make([][]byte, len(c.TransactionIDs))
	for i, transactionID := range c.TransactionIDs {
		transactionIDMessages[i] = transactionID.Bytes()
	}

	return &entities.Collection{
		TransactionIds: transactionIDMessages,
	}
}

func messageToCollection(m *entities.Collection) (flow.Collection, error) {
	if m == nil {
		return flow.Collection{}, errEmptyMessage
	}

	transactionIDMessages := m.GetTransactionIds()

	transactionIDs := make([]flow.Identifier, len(transactionIDMessages))
	for i, transactionIDMsg := range transactionIDMessages {
		transactionIDs[i] = flow.HashToID(transactionIDMsg)
	}

	return flow.Collection{
		TransactionIDs: transactionIDs,
	}, nil
}

func collectionGuaranteeToMessage(g flow.CollectionGuarantee) *entities.CollectionGuarantee {
	return &entities.CollectionGuarantee{
		CollectionId: g.CollectionID.Bytes(),
	}
}

func blockSealToMessage(g flow.BlockSeal) *entities.BlockSeal {
	return &entities.BlockSeal{
		BlockId:            g.BlockID.Bytes(),
		ExecutionReceiptId: g.ExecutionReceiptID.Bytes(),
	}
}

func messageToCollectionGuarantee(m *entities.CollectionGuarantee) (flow.CollectionGuarantee, error) {
	if m == nil {
		return flow.CollectionGuarantee{}, errEmptyMessage
	}

	return flow.CollectionGuarantee{
		CollectionID: flow.HashToID(m.CollectionId),
	}, nil
}

func messageToBlockSeal(m *entities.BlockSeal) (flow.BlockSeal, error) {
	if m == nil {
		return flow.BlockSeal{}, errEmptyMessage
	}

	return flow.BlockSeal{
		BlockID:            flow.BytesToID(m.BlockId),
		ExecutionReceiptID: flow.BytesToID(m.ExecutionReceiptId),
	}, nil
}

func collectionGuaranteesToMessages(l []*flow.CollectionGuarantee) []*entities.CollectionGuarantee {
	results := make([]*entities.CollectionGuarantee, len(l))
	for i, item := range l {
		results[i] = collectionGuaranteeToMessage(*item)
	}
	return results
}

func blockSealsToMessages(l []*flow.BlockSeal) []*entities.BlockSeal {
	results := make([]*entities.BlockSeal, len(l))
	for i, item := range l {
		results[i] = blockSealToMessage(*item)
	}
	return results
}

func messagesToCollectionGuarantees(l []*entities.CollectionGuarantee) ([]*flow.CollectionGuarantee, error) {
	results := make([]*flow.CollectionGuarantee, len(l))
	for i, item := range l {
		temp, err := messageToCollectionGuarantee(item)
		if err != nil {
			return nil, err
		}
		results[i] = &temp
	}
	return results, nil
}

func messagesToBlockSeals(l []*entities.BlockSeal) ([]*flow.BlockSeal, error) {
	results := make([]*flow.BlockSeal, len(l))
	for i, item := range l {
		temp, err := messageToBlockSeal(item)
		if err != nil {
			return nil, err
		}
		results[i] = &temp
	}
	return results, nil
}

func eventToMessage(e flow.Event) (*entities.Event, error) {
	payload, err := cadenceValueToMessage(e.Value)
	if err != nil {
		return nil, err
	}

	return &entities.Event{
		Type:             e.Type,
		TransactionId:    e.TransactionID[:],
		TransactionIndex: uint32(e.TransactionIndex),
		EventIndex:       uint32(e.EventIndex),
		Payload:          payload,
	}, nil
}

func messageToEvent(m *entities.Event, options []jsoncdc.Option) (flow.Event, error) {
	value, err := messageToCadenceValue(m.GetPayload(), options)
	if err != nil {
		return flow.Event{}, err
	}

	eventValue, isEvent := value.(cadence.Event)
	if !isEvent {
		return flow.Event{}, fmt.Errorf("convert: expected Event value, got %s", eventValue.Type().ID())
	}

	return flow.Event{
		Type:             m.GetType(),
		TransactionID:    flow.HashToID(m.GetTransactionId()),
		TransactionIndex: int(m.GetTransactionIndex()),
		EventIndex:       int(m.GetEventIndex()),
		Payload:          m.Payload,
		Value:            eventValue,
	}, nil
}

func messagesToEvents(m []*entities.Event, options []jsoncdc.Option) ([]flow.Event, error) {
	events := make([]flow.Event, 0, len(m))
	for _, ev := range m {
		res, err := messageToEvent(ev, options)
		if err != nil {
			return nil, fmt.Errorf("convert: %w", err)
		}
		events = append(events, res)
	}
	return events, nil
}

func identifierToMessage(i flow.Identifier) []byte {
	return i.Bytes()
}

func messageToIdentifier(b []byte) flow.Identifier {
	return flow.BytesToID(b)
}

func identifiersToMessages(l []flow.Identifier) [][]byte {
	results := make([][]byte, len(l))
	for i, item := range l {
		results[i] = identifierToMessage(item)
	}
	return results
}

func messagesToIdentifiers(l [][]byte) []flow.Identifier {
	results := make([]flow.Identifier, len(l))
	for i, item := range l {
		results[i] = messageToIdentifier(item)
	}
	return results
}

func transactionToMessage(t flow.Transaction) (*entities.Transaction, error) {
	proposalKeyMessage := &entities.Transaction_ProposalKey{
		Address:        t.ProposalKey.Address.Bytes(),
		KeyId:          uint32(t.ProposalKey.KeyIndex),
		SequenceNumber: t.ProposalKey.SequenceNumber,
	}

	authMessages := make([][]byte, len(t.Authorizers))
	for i, auth := range t.Authorizers {
		authMessages[i] = auth.Bytes()
	}

	payloadSigMessages := make([]*entities.Transaction_Signature, len(t.PayloadSignatures))

	for i, sig := range t.PayloadSignatures {
		payloadSigMessages[i] = &entities.Transaction_Signature{
			Address:   sig.Address.Bytes(),
			KeyId:     uint32(sig.KeyIndex),
			Signature: sig.Signature,
		}
	}

	envelopeSigMessages := make([]*entities.Transaction_Signature, len(t.EnvelopeSignatures))

	for i, sig := range t.EnvelopeSignatures {
		envelopeSigMessages[i] = &entities.Transaction_Signature{
			Address:   sig.Address.Bytes(),
			KeyId:     uint32(sig.KeyIndex),
			Signature: sig.Signature,
		}
	}

	return &entities.Transaction{
		Script:             t.Script,
		Arguments:          t.Arguments,
		ReferenceBlockId:   t.ReferenceBlockID.Bytes(),
		GasLimit:           t.GasLimit,
		ProposalKey:        proposalKeyMessage,
		Payer:              t.Payer.Bytes(),
		Authorizers:        authMessages,
		PayloadSignatures:  payloadSigMessages,
		EnvelopeSignatures: envelopeSigMessages,
	}, nil
}

func messageToTransaction(m *entities.Transaction) (flow.Transaction, error) {
	if m == nil {
		return flow.Transaction{}, errEmptyMessage
	}

	t := flow.NewTransaction()

	t.SetScript(m.GetScript())
	t.SetReferenceBlockID(flow.HashToID(m.GetReferenceBlockId()))
	t.SetComputeLimit(m.GetGasLimit())

	for _, arg := range m.GetArguments() {
		t.AddRawArgument(arg)
	}

	proposalKey := m.GetProposalKey()
	if proposalKey != nil {
		proposalAddress := flow.BytesToAddress(proposalKey.GetAddress())
		t.SetProposalKey(proposalAddress, int(proposalKey.GetKeyId()), proposalKey.GetSequenceNumber())
	}

	payer := m.GetPayer()
	if payer != nil {
		t.SetPayer(
			flow.BytesToAddress(payer),
		)
	}

	for _, authorizer := range m.GetAuthorizers() {
		t.AddAuthorizer(
			flow.BytesToAddress(authorizer),
		)
	}

	for _, sig := range m.GetPayloadSignatures() {
		addr := flow.BytesToAddress(sig.GetAddress())
		t.AddPayloadSignature(addr, int(sig.GetKeyId()), sig.GetSignature())
	}

	for _, sig := range m.GetEnvelopeSignatures() {
		addr := flow.BytesToAddress(sig.GetAddress())
		t.AddEnvelopeSignature(addr, int(sig.GetKeyId()), sig.GetSignature())
	}

	return *t, nil
}

func transactionResultToMessage(result flow.TransactionResult) (*access.TransactionResultResponse, error) {
	eventMessages := make([]*entities.Event, len(result.Events))

	for i, event := range result.Events {
		eventMsg, err := eventToMessage(event)
		if err != nil {
			return nil, err
		}

		eventMessages[i] = eventMsg
	}

	statusCode := 0
	errorMsg := ""

	if result.Error != nil {
		statusCode = 1
		errorMsg = result.Error.Error()
	}

	return &access.TransactionResultResponse{
		Status:        entities.TransactionStatus(result.Status),
		StatusCode:    uint32(statusCode),
		ErrorMessage:  errorMsg,
		Events:        eventMessages,
		BlockId:       identifierToMessage(result.BlockID),
		BlockHeight:   result.BlockHeight,
		TransactionId: identifierToMessage(result.TransactionID),
		CollectionId:  identifierToMessage(result.CollectionID),
	}, nil
}

func messageToTransactionResult(m *access.TransactionResultResponse, options []jsoncdc.Option) (flow.TransactionResult, error) {
	eventMessages := m.GetEvents()

	events := make([]flow.Event, len(eventMessages))
	for i, eventMsg := range eventMessages {
		event, err := messageToEvent(eventMsg, options)
		if err != nil {
			return flow.TransactionResult{}, err
		}

		events[i] = event
	}

	var err error

	statusCode := m.GetStatusCode()
	if statusCode != 0 {
		errorMsg := m.GetErrorMessage()
		if errorMsg != "" {
			err = errors.New(errorMsg)
		} else {
			err = errors.New("transaction execution failed")
		}
	}

	return flow.TransactionResult{
		Status:        flow.TransactionStatus(m.GetStatus()),
		Error:         err,
		Events:        events,
		BlockID:       flow.BytesToID(m.GetBlockId()),
		BlockHeight:   m.GetBlockHeight(),
		TransactionID: flow.BytesToID(m.GetTransactionId()),
		CollectionID:  flow.BytesToID(m.GetCollectionId()),
	}, nil
}

func blockExecutionDataToMessage(
	execData *flow.ExecutionData,
) (*entities.BlockExecutionData, error) {
	chunks := make([]*entities.ChunkExecutionData, len(execData.ChunkExecutionData))
	for i, chunk := range execData.ChunkExecutionData {
		convertedChunk, err := chunkExecutionDataToMessage(chunk)
		if err != nil {
			return nil, err
		}
		chunks[i] = convertedChunk
	}

	return &entities.BlockExecutionData{
		BlockId:            identifierToMessage(execData.BlockID),
		ChunkExecutionData: chunks,
	}, nil
}

func messageToBlockExecutionData(
	m *entities.BlockExecutionData,
) (*flow.ExecutionData, error) {
	if m == nil {
		return nil, errEmptyMessage
	}

	chunks := make([]*flow.ChunkExecutionData, len(m.ChunkExecutionData))
	for i, chunk := range m.GetChunkExecutionData() {
		convertedChunk, err := messageToChunkExecutionData(chunk)
		if err != nil {
			return nil, err
		}
		chunks[i] = convertedChunk
	}

	return &flow.ExecutionData{
		BlockID:            messageToIdentifier(m.GetBlockId()),
		ChunkExecutionData: chunks,
	}, nil
}

func chunkExecutionDataToMessage(
	chunk *flow.ChunkExecutionData,
) (*entities.ChunkExecutionData, error) {

	transactions, err := executionDataCollectionToMessage(chunk.Transactions)
	if err != nil {
		return nil, err
	}

	var trieUpdate *entities.TrieUpdate
	if chunk.TrieUpdate != nil {
		trieUpdate, err = trieUpdateToMessage(chunk.TrieUpdate)
		if err != nil {
			return nil, err
		}
	}

	events := make([]*entities.Event, len(chunk.Events))
	for i, ev := range chunk.Events {
		res, err := eventToMessage(*ev)
		if err != nil {
			return nil, err
		}

		// execution data uses CCF encoding
		res.Payload, err = ccf.Encode(ev.Value)
		if err != nil {
			return nil, fmt.Errorf("ccf convert: %w", err)
		}

		events[i] = res
	}

	results := make([]*entities.ExecutionDataTransactionResult, len(chunk.TransactionResults))
	for i, res := range chunk.TransactionResults {
		result := lightTransactionResultToMessage(res)
		results[i] = result
	}

	return &entities.ChunkExecutionData{
		Collection:         transactions,
		Events:             events,
		TrieUpdate:         trieUpdate,
		TransactionResults: results,
	}, nil
}

func messageToChunkExecutionData(
	m *entities.ChunkExecutionData,
) (*flow.ChunkExecutionData, error) {

	transactions, err := messageToExecutionDataCollection(m.GetCollection())
	if err != nil {
		return nil, err
	}

	var trieUpdate *flow.TrieUpdate
	if m.GetTrieUpdate() != nil {
		trieUpdate, err = messageToTrieUpdate(m.GetTrieUpdate())
		if err != nil {
			return nil, err
		}
	}

	events := make([]*flow.Event, len(m.GetEvents()))
	for i, ev := range m.GetEvents() {
		res, err := messageToEvent(ev, nil)
		if err != nil {
			return nil, err
		}
		events[i] = &res
	}

	results := make([]*flow.LightTransactionResult, len(m.GetTransactionResults()))
	for i, res := range m.GetTransactionResults() {
		result := messageToLightTransactionResult(res)
		results[i] = &result
	}

	return &flow.ChunkExecutionData{
		Transactions:       transactions,
		Events:             events,
		TrieUpdate:         trieUpdate,
		TransactionResults: results,
	}, nil
}

func executionDataCollectionToMessage(
	txs []*flow.Transaction,
) (*entities.ExecutionDataCollection, error) {
	transactions := make([]*entities.Transaction, len(txs))
	for i, tx := range txs {
		transaction, err := transactionToMessage(*tx)
		if err != nil {
			return nil, fmt.Errorf("could not convert transaction %d: %w", i, err)
		}
		transactions[i] = transaction
	}

	return &entities.ExecutionDataCollection{
		Transactions: transactions,
	}, nil
}

func messageToExecutionDataCollection(
	m *entities.ExecutionDataCollection,
) ([]*flow.Transaction, error) {
	messages := m.GetTransactions()
	transactions := make([]*flow.Transaction, len(messages))
	for i, message := range messages {
		transaction, err := messageToTransaction(message)
		if err != nil {
			return nil, fmt.Errorf("could not convert transaction %d: %w", i, err)
		}
		transactions[i] = &transaction
	}

	if len(transactions) == 0 {
		return nil, nil
	}

	return transactions, nil
}

func trieUpdateToMessage(
	update *flow.TrieUpdate,
) (*entities.TrieUpdate, error) {

	payloads := make([]*entities.Payload, len(update.Payloads))
	for i, payload := range update.Payloads {
		keyParts := make([]*entities.KeyPart, len(payload.KeyPart))
		for j, keypart := range payload.KeyPart {
			keyParts[j] = &entities.KeyPart{
				Type:  uint32(keypart.Type),
				Value: keypart.Value,
			}
		}
		payloads[i] = &entities.Payload{
			KeyPart: keyParts,
			Value:   payload.Value,
		}
	}

	return &entities.TrieUpdate{
		RootHash: update.RootHash,
		Paths:    update.Paths,
		Payloads: payloads,
	}, nil
}

func messageToTrieUpdate(
	m *entities.TrieUpdate,
) (*flow.TrieUpdate, error) {
	rootHash := m.GetRootHash()
	paths := m.GetPaths()

	payloads := make([]*flow.Payload, len(m.Payloads))
	for i, payload := range m.GetPayloads() {
		keyParts := make([]*flow.KeyPart, len(payload.GetKeyPart()))
		for j, keypart := range payload.GetKeyPart() {
			keyParts[j] = &flow.KeyPart{
				Type:  uint16(keypart.GetType()),
				Value: keypart.GetValue(),
			}
		}
		payloads[i] = &flow.Payload{
			KeyPart: keyParts,
			Value:   payload.GetValue(),
		}
	}

	return &flow.TrieUpdate{
		RootHash: rootHash,
		Paths:    paths,
		Payloads: payloads,
	}, nil
}

func lightTransactionResultToMessage(
	result *flow.LightTransactionResult,
) *entities.ExecutionDataTransactionResult {
	return &entities.ExecutionDataTransactionResult{
		TransactionId:   identifierToMessage(result.TransactionID),
		Failed:          result.Failed,
		ComputationUsed: result.ComputationUsed,
	}
}

func messageToLightTransactionResult(
	m *entities.ExecutionDataTransactionResult,
) flow.LightTransactionResult {
	return flow.LightTransactionResult{
		TransactionID:   messageToIdentifier(m.GetTransactionId()),
		Failed:          m.Failed,
		ComputationUsed: m.GetComputationUsed(),
	}
}

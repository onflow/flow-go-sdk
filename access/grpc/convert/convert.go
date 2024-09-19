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

package convert

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

var ErrEmptyMessage = errors.New("protobuf message is empty")

func AccountToMessage(a flow.Account) *entities.Account {
	accountKeys := make([]*entities.AccountKey, len(a.Keys))
	for i, key := range a.Keys {
		accountKeys[i] = AccountKeyToMessage(key)
	}

	return &entities.Account{
		Address:   a.Address.Bytes(),
		Balance:   a.Balance,
		Code:      a.Code,
		Keys:      accountKeys,
		Contracts: a.Contracts,
	}
}

func MessageToAccount(m *entities.Account) (flow.Account, error) {
	if m == nil {
		return flow.Account{}, ErrEmptyMessage
	}

	accountKeys := make([]*flow.AccountKey, len(m.GetKeys()))
	for i, key := range m.GetKeys() {
		accountKey, err := MessageToAccountKey(key)
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

func AccountKeyToMessage(a *flow.AccountKey) *entities.AccountKey {
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

func MessageToAccountKey(m *entities.AccountKey) (*flow.AccountKey, error) {
	if m == nil {
		return nil, ErrEmptyMessage
	}

	sigAlgo := crypto.SignatureAlgorithm(m.GetSignAlgo())
	hashAlgo := crypto.HashAlgorithm(m.GetHashAlgo())

	publicKey, err := crypto.DecodePublicKey(sigAlgo, m.GetPublicKey())
	if err != nil {
		return nil, err
	}

	return &flow.AccountKey{
		Index:          m.GetIndex(),
		PublicKey:      publicKey,
		SigAlgo:        sigAlgo,
		HashAlgo:       hashAlgo,
		Weight:         int(m.GetWeight()),
		SequenceNumber: uint64(m.GetSequenceNumber()),
		Revoked:        m.GetRevoked(),
	}, nil
}

func BlockToMessage(b flow.Block) (*entities.Block, error) {

	t := timestamppb.New(b.BlockHeader.Timestamp)

	return &entities.Block{
		Id:                   b.BlockHeader.ID.Bytes(),
		ParentId:             b.BlockHeader.ParentID.Bytes(),
		Height:               b.BlockHeader.Height,
		Timestamp:            t,
		CollectionGuarantees: CollectionGuaranteesToMessages(b.BlockPayload.CollectionGuarantees),
		BlockSeals:           BlockSealsToMessages(b.BlockPayload.Seals),
	}, nil
}

func MessageToBlock(m *entities.Block) (flow.Block, error) {
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

	guarantees, err := MessagesToCollectionGuarantees(m.GetCollectionGuarantees())
	if err != nil {
		return flow.Block{}, err
	}

	seals, err := MessagesToBlockSeals(m.GetBlockSeals())
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

func BlockHeaderToMessage(b flow.BlockHeader) (*entities.BlockHeader, error) {
	t := timestamppb.New(b.Timestamp)

	return &entities.BlockHeader{
		Id:        b.ID.Bytes(),
		ParentId:  b.ParentID.Bytes(),
		Height:    b.Height,
		Timestamp: t,
	}, nil
}

func MessageToBlockHeader(m *entities.BlockHeader) (flow.BlockHeader, error) {
	if m == nil {
		return flow.BlockHeader{}, ErrEmptyMessage
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

func MessageToBlockDigest(m *access.SubscribeBlockDigestsResponse) flow.BlockDigest {
	return flow.BlockDigest{
		BlockID:   flow.BytesToID(m.GetBlockId()),
		Height:    m.GetBlockHeight(),
		Timestamp: m.GetBlockTimestamp().AsTime(),
	}
}

func BlockDigestToMessage(blockDigest flow.BlockDigest) *access.SubscribeBlockDigestsResponse {
	return &access.SubscribeBlockDigestsResponse{
		BlockId:        IdentifierToMessage(blockDigest.BlockID),
		BlockHeight:    blockDigest.Height,
		BlockTimestamp: timestamppb.New(blockDigest.Timestamp),
	}
}

func BlockStatusToEntity(blockStatus flow.BlockStatus) entities.BlockStatus {
	switch blockStatus {
	case flow.BlockStatusFinalized:
		return entities.BlockStatus_BLOCK_FINALIZED
	case flow.BlockStatusSealed:
		return entities.BlockStatus_BLOCK_SEALED
	default:
		return entities.BlockStatus_BLOCK_UNKNOWN
	}
}

func CadenceValueToMessage(value cadence.Value, encodingVersion flow.EventEncodingVersion) ([]byte, error) {
	switch encodingVersion {
	case flow.EventEncodingVersionCCF:
		b, err := ccf.Encode(value)
		if err != nil {
			return nil, fmt.Errorf("ccf convert: %w", err)
		}
		return b, nil
	case flow.EventEncodingVersionJSONCDC:
		b, err := jsoncdc.Encode(value)
		if err != nil {
			return nil, fmt.Errorf("jsoncdc convert: %w", err)
		}
		return b, nil
	default:
		return nil, fmt.Errorf("unsupported cadence encoding version: %v", encodingVersion)
	}
}

func CadenceValuesToMessages(values []cadence.Value, encodingVersion flow.EventEncodingVersion) ([][]byte, error) {
	msgs := make([][]byte, len(values))
	for i, val := range values {
		msg, err := CadenceValueToMessage(val, encodingVersion)
		if err != nil {
			return nil, err
		}
		msgs[i] = msg
	}
	return msgs, nil
}

func MessageToCadenceValue(m []byte, options []jsoncdc.Option) (cadence.Value, error) {
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

func CollectionToMessage(c flow.Collection) *entities.Collection {
	transactionIDMessages := make([][]byte, len(c.TransactionIDs))
	for i, transactionID := range c.TransactionIDs {
		transactionIDMessages[i] = transactionID.Bytes()
	}

	return &entities.Collection{
		TransactionIds: transactionIDMessages,
	}
}

func MessageToCollection(m *entities.Collection) (flow.Collection, error) {
	if m == nil {
		return flow.Collection{}, ErrEmptyMessage
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

func CollectionGuaranteeToMessage(g flow.CollectionGuarantee) *entities.CollectionGuarantee {
	return &entities.CollectionGuarantee{
		CollectionId: g.CollectionID.Bytes(),
	}
}

func BlockSealToMessage(g flow.BlockSeal) *entities.BlockSeal {
	return &entities.BlockSeal{
		BlockId:            g.BlockID.Bytes(),
		ExecutionReceiptId: g.ExecutionReceiptID.Bytes(),
	}
}

func MessageToCollectionGuarantee(m *entities.CollectionGuarantee) (flow.CollectionGuarantee, error) {
	if m == nil {
		return flow.CollectionGuarantee{}, ErrEmptyMessage
	}

	return flow.CollectionGuarantee{
		CollectionID: flow.HashToID(m.CollectionId),
	}, nil
}

func MessageToBlockSeal(m *entities.BlockSeal) (flow.BlockSeal, error) {
	if m == nil {
		return flow.BlockSeal{}, ErrEmptyMessage
	}

	return flow.BlockSeal{
		BlockID:            flow.BytesToID(m.BlockId),
		ExecutionReceiptID: flow.BytesToID(m.ExecutionReceiptId),
	}, nil
}

func CollectionGuaranteesToMessages(l []*flow.CollectionGuarantee) []*entities.CollectionGuarantee {
	results := make([]*entities.CollectionGuarantee, len(l))
	for i, item := range l {
		results[i] = CollectionGuaranteeToMessage(*item)
	}
	return results
}

func BlockSealsToMessages(l []*flow.BlockSeal) []*entities.BlockSeal {
	results := make([]*entities.BlockSeal, len(l))
	for i, item := range l {
		results[i] = BlockSealToMessage(*item)
	}
	return results
}

func MessagesToCollectionGuarantees(l []*entities.CollectionGuarantee) ([]*flow.CollectionGuarantee, error) {
	results := make([]*flow.CollectionGuarantee, len(l))
	for i, item := range l {
		temp, err := MessageToCollectionGuarantee(item)
		if err != nil {
			return nil, err
		}
		results[i] = &temp
	}
	return results, nil
}

func MessagesToBlockSeals(l []*entities.BlockSeal) ([]*flow.BlockSeal, error) {
	results := make([]*flow.BlockSeal, len(l))
	for i, item := range l {
		temp, err := MessageToBlockSeal(item)
		if err != nil {
			return nil, err
		}
		results[i] = &temp
	}
	return results, nil
}

func EventToMessage(e flow.Event, encodingVersion flow.EventEncodingVersion) (*entities.Event, error) {
	payload, err := CadenceValueToMessage(e.Value, encodingVersion)
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

func MessageToEvent(m *entities.Event, options []jsoncdc.Option) (flow.Event, error) {
	value, err := MessageToCadenceValue(m.GetPayload(), options)
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

func MessagesToEvents(m []*entities.Event, options []jsoncdc.Option) ([]flow.Event, error) {
	events := make([]flow.Event, 0, len(m))
	for _, ev := range m {
		res, err := MessageToEvent(ev, options)
		if err != nil {
			return nil, fmt.Errorf("convert: %w", err)
		}
		events = append(events, res)
	}
	return events, nil
}

func IdentifierToMessage(i flow.Identifier) []byte {
	return i.Bytes()
}

func MessageToIdentifier(b []byte) flow.Identifier {
	return flow.BytesToID(b)
}

func IdentifiersToMessages(l []flow.Identifier) [][]byte {
	results := make([][]byte, len(l))
	for i, item := range l {
		results[i] = IdentifierToMessage(item)
	}
	return results
}

func MessagesToIdentifiers(l [][]byte) []flow.Identifier {
	results := make([]flow.Identifier, len(l))
	for i, item := range l {
		results[i] = MessageToIdentifier(item)
	}
	return results
}

func TransactionToMessage(t flow.Transaction) (*entities.Transaction, error) {
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

func MessageToTransaction(m *entities.Transaction) (flow.Transaction, error) {
	if m == nil {
		return flow.Transaction{}, ErrEmptyMessage
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
		t.SetProposalKey(proposalAddress, proposalKey.GetKeyId(), proposalKey.GetSequenceNumber())
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
		t.AddPayloadSignature(addr, sig.GetKeyId(), sig.GetSignature())
	}

	for _, sig := range m.GetEnvelopeSignatures() {
		addr := flow.BytesToAddress(sig.GetAddress())
		t.AddEnvelopeSignature(addr, sig.GetKeyId(), sig.GetSignature())
	}

	return *t, nil
}

func TransactionResultToMessage(result flow.TransactionResult, encodingVersion flow.EventEncodingVersion) (*access.TransactionResultResponse, error) {
	eventMessages := make([]*entities.Event, len(result.Events))

	for i, event := range result.Events {
		eventMsg, err := EventToMessage(event, encodingVersion)
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
		BlockId:       IdentifierToMessage(result.BlockID),
		BlockHeight:   result.BlockHeight,
		TransactionId: IdentifierToMessage(result.TransactionID),
		CollectionId:  IdentifierToMessage(result.CollectionID),
	}, nil
}

func MessageToTransactionResult(m *access.TransactionResultResponse, options []jsoncdc.Option) (flow.TransactionResult, error) {
	eventMessages := m.GetEvents()

	events := make([]flow.Event, len(eventMessages))
	for i, eventMsg := range eventMessages {
		event, err := MessageToEvent(eventMsg, options)
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

func MessageToExecutionResult(execResult *entities.ExecutionResult) (*flow.ExecutionResult, error) {
	chunks := make([]*flow.Chunk, len(execResult.Chunks))
	serviceEvents := make([]*flow.ServiceEvent, len(execResult.ServiceEvents))

	for i, chunk := range execResult.Chunks {
		chunks[i] = &flow.Chunk{
			CollectionIndex:      uint(chunk.CollectionIndex),
			StartState:           flow.BytesToStateCommitment(chunk.StartState),
			EventCollection:      flow.BytesToHash(chunk.EventCollection),
			BlockID:              flow.BytesToID(chunk.BlockId),
			TotalComputationUsed: chunk.TotalComputationUsed,
			NumberOfTransactions: uint16(chunk.NumberOfTransactions),
			Index:                chunk.Index,
			EndState:             flow.BytesToStateCommitment(chunk.EndState),
		}
	}

	for i, serviceEvent := range execResult.ServiceEvents {
		serviceEvents[i] = &flow.ServiceEvent{
			Type:    serviceEvent.Type,
			Payload: serviceEvent.Payload,
		}
	}

	return &flow.ExecutionResult{
		PreviousResultID: flow.BytesToID(execResult.PreviousResultId),
		BlockID:          flow.BytesToID(execResult.BlockId),
		Chunks:           chunks,
		ServiceEvents:    serviceEvents,
	}, nil
}

func BlockExecutionDataToMessage(
	execData *flow.ExecutionData,
) (*entities.BlockExecutionData, error) {
	chunks := make([]*entities.ChunkExecutionData, len(execData.ChunkExecutionData))
	for i, chunk := range execData.ChunkExecutionData {
		convertedChunk, err := ChunkExecutionDataToMessage(chunk)
		if err != nil {
			return nil, err
		}
		chunks[i] = convertedChunk
	}

	return &entities.BlockExecutionData{
		BlockId:            IdentifierToMessage(execData.BlockID),
		ChunkExecutionData: chunks,
	}, nil
}

func MessageToBlockExecutionData(
	m *entities.BlockExecutionData,
) (*flow.ExecutionData, error) {
	if m == nil {
		return nil, ErrEmptyMessage
	}

	chunks := make([]*flow.ChunkExecutionData, len(m.ChunkExecutionData))
	for i, chunk := range m.GetChunkExecutionData() {
		convertedChunk, err := MessageToChunkExecutionData(chunk)
		if err != nil {
			return nil, err
		}
		chunks[i] = convertedChunk
	}

	return &flow.ExecutionData{
		BlockID:            MessageToIdentifier(m.GetBlockId()),
		ChunkExecutionData: chunks,
	}, nil
}

func ChunkExecutionDataToMessage(
	chunk *flow.ChunkExecutionData,
) (*entities.ChunkExecutionData, error) {

	transactions, err := ExecutionDataCollectionToMessage(chunk.Transactions)
	if err != nil {
		return nil, err
	}

	var trieUpdate *entities.TrieUpdate
	if chunk.TrieUpdate != nil {
		trieUpdate, err = TrieUpdateToMessage(chunk.TrieUpdate)
		if err != nil {
			return nil, err
		}
	}

	events := make([]*entities.Event, len(chunk.Events))
	for i, ev := range chunk.Events {
		// execution data uses CCF encoding
		res, err := EventToMessage(*ev, flow.EventEncodingVersionCCF)
		if err != nil {
			return nil, err
		}

		events[i] = res
	}

	results := make([]*entities.ExecutionDataTransactionResult, len(chunk.TransactionResults))
	for i, res := range chunk.TransactionResults {
		result := LightTransactionResultToMessage(res)
		results[i] = result
	}

	return &entities.ChunkExecutionData{
		Collection:         transactions,
		Events:             events,
		TrieUpdate:         trieUpdate,
		TransactionResults: results,
	}, nil
}

func MessageToChunkExecutionData(
	m *entities.ChunkExecutionData,
) (*flow.ChunkExecutionData, error) {

	transactions, err := MessageToExecutionDataCollection(m.GetCollection())
	if err != nil {
		return nil, err
	}

	var trieUpdate *flow.TrieUpdate
	if m.GetTrieUpdate() != nil {
		trieUpdate, err = MessageToTrieUpdate(m.GetTrieUpdate())
		if err != nil {
			return nil, err
		}
	}

	events := make([]*flow.Event, len(m.GetEvents()))
	for i, ev := range m.GetEvents() {
		res, err := MessageToEvent(ev, nil)
		if err != nil {
			return nil, err
		}
		events[i] = &res
	}

	results := make([]*flow.LightTransactionResult, len(m.GetTransactionResults()))
	for i, res := range m.GetTransactionResults() {
		result := MessageToLightTransactionResult(res)
		results[i] = &result
	}

	return &flow.ChunkExecutionData{
		Transactions:       transactions,
		Events:             events,
		TrieUpdate:         trieUpdate,
		TransactionResults: results,
	}, nil
}

func ExecutionDataCollectionToMessage(
	txs []*flow.Transaction,
) (*entities.ExecutionDataCollection, error) {
	transactions := make([]*entities.Transaction, len(txs))
	for i, tx := range txs {
		transaction, err := TransactionToMessage(*tx)
		if err != nil {
			return nil, fmt.Errorf("could not convert transaction %d: %w", i, err)
		}
		transactions[i] = transaction
	}

	return &entities.ExecutionDataCollection{
		Transactions: transactions,
	}, nil
}

func MessageToExecutionDataCollection(
	m *entities.ExecutionDataCollection,
) ([]*flow.Transaction, error) {
	messages := m.GetTransactions()
	transactions := make([]*flow.Transaction, len(messages))
	for i, message := range messages {
		transaction, err := MessageToTransaction(message)
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

func TrieUpdateToMessage(
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

func MessageToTrieUpdate(
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

func LightTransactionResultToMessage(
	result *flow.LightTransactionResult,
) *entities.ExecutionDataTransactionResult {
	return &entities.ExecutionDataTransactionResult{
		TransactionId:   IdentifierToMessage(result.TransactionID),
		Failed:          result.Failed,
		ComputationUsed: result.ComputationUsed,
	}
}

func MessageToLightTransactionResult(
	m *entities.ExecutionDataTransactionResult,
) flow.LightTransactionResult {
	return flow.LightTransactionResult{
		TransactionID:   MessageToIdentifier(m.GetTransactionId()),
		Failed:          m.Failed,
		ComputationUsed: m.GetComputationUsed(),
	}
}

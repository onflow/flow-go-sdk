/*
 * Flow Go SDK
 *
 * Copyright 2019-2020 Dapper Labs, Inc.
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

	"github.com/golang/protobuf/ptypes"
	"github.com/onflow/cadence"
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
		Index:          int(m.GetIndex()),
		PublicKey:      publicKey,
		SigAlgo:        sigAlgo,
		HashAlgo:       hashAlgo,
		Weight:         int(m.GetWeight()),
		SequenceNumber: uint64(m.GetSequenceNumber()),
		Revoked:        m.GetRevoked(),
	}, nil
}

func BlockToMessage(b flow.Block) (*entities.Block, error) {
	t, err := ptypes.TimestampProto(b.BlockHeader.Timestamp)
	if err != nil {
		return nil, fmt.Errorf("convert: failed to convert block timestamp to message: %w", err)
	}

	return &entities.Block{
		Id:                   b.BlockHeader.ID.Bytes(),
		ParentId:             b.BlockHeader.ParentID.Bytes(),
		Height:               b.BlockHeader.Height,
		Timestamp:            t,
		CollectionGuarantees: CollectionGuaranteesToMessages(b.BlockPayload.CollectionGuarantees),
	}, nil
}

func MessageToBlock(m *entities.Block) (flow.Block, error) {
	var timestamp time.Time
	var err error

	if m.GetTimestamp() != nil {
		timestamp, err = ptypes.Timestamp(m.GetTimestamp())
		if err != nil {
			return flow.Block{}, err
		}
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

	payload := &flow.BlockPayload{
		CollectionGuarantees: guarantees,
	}

	return flow.Block{
		BlockHeader:  *header,
		BlockPayload: *payload,
	}, nil
}

func BlockHeaderToMessage(b flow.BlockHeader) (*entities.BlockHeader, error) {
	t, err := ptypes.TimestampProto(b.Timestamp)
	if err != nil {
		return nil, fmt.Errorf("convert: failed to convert message to block timestamp: %w", err)
	}

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
	var err error

	if m.GetTimestamp() != nil {
		timestamp, err = ptypes.Timestamp(m.GetTimestamp())
		if err != nil {
			return flow.BlockHeader{}, err
		}
	}

	return flow.BlockHeader{
		ID:        flow.HashToID(m.GetId()),
		ParentID:  flow.HashToID(m.GetParentId()),
		Height:    m.GetHeight(),
		Timestamp: timestamp,
	}, nil
}

func CadenceValueToMessage(value cadence.Value) ([]byte, error) {
	b, err := jsoncdc.Encode(value)
	if err != nil {
		return nil, fmt.Errorf("convert: %w", err)
	}

	return b, nil
}

func CadenceValuesToMessages(values []cadence.Value) ([][]byte, error) {
	msgs := make([][]byte, len(values))
	for i, val := range values {
		msg, err := CadenceValueToMessage(val)
		if err != nil {
			return nil, fmt.Errorf("convert: %w", err)
		}
		msgs[i] = msg
	}
	return msgs, nil
}

func MessageToCadenceValue(m []byte) (cadence.Value, error) {
	v, err := jsoncdc.Decode(m)
	if err != nil {
		return nil, fmt.Errorf("convert: %w", err)
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

func MessageToCollectionGuarantee(m *entities.CollectionGuarantee) (flow.CollectionGuarantee, error) {
	if m == nil {
		return flow.CollectionGuarantee{}, ErrEmptyMessage
	}

	return flow.CollectionGuarantee{
		CollectionID: flow.HashToID(m.CollectionId),
	}, nil
}

func CollectionGuaranteesToMessages(l []*flow.CollectionGuarantee) []*entities.CollectionGuarantee {
	results := make([]*entities.CollectionGuarantee, len(l))
	for i, item := range l {
		results[i] = CollectionGuaranteeToMessage(*item)
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

func EventToMessage(e flow.Event) (*entities.Event, error) {
	payload, err := CadenceValueToMessage(e.Value)
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

func MessageToEvent(m *entities.Event) (flow.Event, error) {
	value, err := MessageToCadenceValue(m.GetPayload())
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
		Value:            eventValue,
	}, nil
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
	t.SetGasLimit(m.GetGasLimit())

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

func TransactionResultToMessage(result flow.TransactionResult) (*access.TransactionResultResponse, error) {
	eventMessages := make([]*entities.Event, len(result.Events))

	for i, event := range result.Events {
		eventMsg, err := EventToMessage(event)
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
		Status:       entities.TransactionStatus(result.Status),
		StatusCode:   uint32(statusCode),
		ErrorMessage: errorMsg,
		Events:       eventMessages,
	}, nil
}

func MessageToTransactionResult(m *access.TransactionResultResponse) (flow.TransactionResult, error) {
	eventMessages := m.GetEvents()

	events := make([]flow.Event, len(eventMessages))
	for i, eventMsg := range eventMessages {
		event, err := MessageToEvent(eventMsg)
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
		Status: flow.TransactionStatus(m.GetStatus()),
		Error:  err,
		Events: events,
	}, nil
}

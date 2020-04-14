package convert

import (
	"errors"
	"fmt"

	"github.com/dapperlabs/cadence"
	jsoncdc "github.com/dapperlabs/cadence/encoding/json"
	"github.com/dapperlabs/flow-go/crypto"
	"github.com/dapperlabs/flow/protobuf/go/flow/access"
	"github.com/dapperlabs/flow/protobuf/go/flow/entities"

	"github.com/dapperlabs/flow-go-sdk"
)

var ErrEmptyMessage = errors.New("protobuf message is empty")

func MessageToBlockHeader(m *entities.BlockHeader) flow.BlockHeader {
	return flow.BlockHeader{
		ID:       flow.HashToID(m.GetId()),
		ParentID: flow.HashToID(m.GetParentId()),
		Height:   m.GetHeight(),
	}
}

func BlockHeaderToMessage(b flow.BlockHeader) *entities.BlockHeader {
	return &entities.BlockHeader{
		Id:       b.ID[:],
		ParentId: b.ParentID[:],
		Height:   b.Height,
	}
}

func MessageToTransaction(m *entities.Transaction) (flow.Transaction, error) {
	if m == nil {
		return flow.Transaction{}, ErrEmptyMessage
	}

	t := flow.NewTransaction()

	t.SetScript(m.GetScript())
	t.SetReferenceBlockID(flow.HashToID(m.GetReferenceBlockId()))
	t.SetGasLimit(m.GetGasLimit())

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

func TransactionToMessage(t flow.Transaction) *entities.Transaction {
	proposalKeyMessage := &entities.Transaction_ProposalKey{
		Address:        t.ProposalKey.Address.Bytes(),
		KeyId:          uint32(t.ProposalKey.KeyID),
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
			KeyId:     uint32(sig.KeyID),
			Signature: sig.Signature,
		}
	}

	envelopeSigMessages := make([]*entities.Transaction_Signature, len(t.EnvelopeSignatures))

	for i, sig := range t.EnvelopeSignatures {
		envelopeSigMessages[i] = &entities.Transaction_Signature{
			Address:   sig.Address.Bytes(),
			KeyId:     uint32(sig.KeyID),
			Signature: sig.Signature,
		}
	}

	return &entities.Transaction{
		Script:             t.Script,
		ReferenceBlockId:   t.ReferenceBlockID.Bytes(),
		GasLimit:           t.GasLimit,
		ProposalKey:        proposalKeyMessage,
		Payer:              t.Payer.Bytes(),
		Authorizers:        authMessages,
		PayloadSignatures:  payloadSigMessages,
		EnvelopeSignatures: envelopeSigMessages,
	}
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

func MessageToAccount(m *entities.Account) (flow.Account, error) {
	if m == nil {
		return flow.Account{}, ErrEmptyMessage
	}

	accountKeys := make([]flow.AccountKey, len(m.Keys))
	for i, key := range m.Keys {
		accountKey, err := MessageToAccountKey(key)
		if err != nil {
			return flow.Account{}, err
		}

		accountKeys[i] = accountKey
	}

	return flow.Account{
		Address: flow.BytesToAddress(m.Address),
		Balance: m.Balance,
		Code:    m.Code,
		Keys:    accountKeys,
	}, nil
}

func AccountToMessage(a flow.Account) (*entities.Account, error) {
	accountKeys := make([]*entities.AccountKey, len(a.Keys))
	for i, key := range a.Keys {
		accountKeyMsg, err := AccountKeyToMessage(key)
		if err != nil {
			return nil, err
		}
		accountKeys[i] = accountKeyMsg
	}

	return &entities.Account{
		Address: a.Address.Bytes(),
		Balance: a.Balance,
		Code:    a.Code,
		Keys:    accountKeys,
	}, nil
}

func MessageToAccountKey(m *entities.AccountKey) (flow.AccountKey, error) {
	if m == nil {
		return flow.AccountKey{}, ErrEmptyMessage
	}

	signAlgo := crypto.SigningAlgorithm(m.GetSignAlgo())
	hashAlgo := crypto.HashingAlgorithm(m.GetHashAlgo())

	publicKey, err := crypto.DecodePublicKey(signAlgo, m.GetPublicKey())
	if err != nil {
		return flow.AccountKey{}, err
	}

	return flow.AccountKey{
		PublicKey:      publicKey,
		SignAlgo:       signAlgo,
		HashAlgo:       hashAlgo,
		Weight:         int(m.GetWeight()),
		SequenceNumber: uint64(m.GetSequenceNumber()),
	}, nil
}

func AccountKeyToMessage(a flow.AccountKey) (*entities.AccountKey, error) {
	publicKey, err := a.PublicKey.Encode()
	if err != nil {
		return nil, err
	}

	return &entities.AccountKey{
		PublicKey:      publicKey,
		SignAlgo:       uint32(a.SignAlgo),
		HashAlgo:       uint32(a.HashAlgo),
		Weight:         uint32(a.Weight),
		SequenceNumber: uint32(a.SequenceNumber),
	}, nil
}

func MessageToEvent(m *entities.Event) (flow.Event, error) {
	value, err := jsoncdc.Decode(m.GetPayload())
	if err != nil {
		return flow.Event{}, nil
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

func EventToMessage(e flow.Event) (*entities.Event, error) {
	payload, err := jsoncdc.Encode(e.Value)
	if err != nil {
		return nil, fmt.Errorf("convert: %w", err)
	}

	return &entities.Event{
		Type:             e.Type,
		TransactionId:    e.TransactionID[:],
		TransactionIndex: uint32(e.TransactionIndex),
		EventIndex:       uint32(e.EventIndex),
		Payload:          payload,
	}, nil
}

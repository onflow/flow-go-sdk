package convert

import (
	"errors"
	"fmt"

	"github.com/dapperlabs/cadence"
	jsoncdc "github.com/dapperlabs/cadence/encoding/json"
	"github.com/dapperlabs/flow-go/crypto"
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

func MessageToAccountSignature(m *entities.AccountSignature) flow.AccountSignature {
	return flow.AccountSignature{
		Account:   flow.BytesToAddress(m.GetAccount()),
		Signature: m.GetSignature(),
	}
}

func AccountSignatureToMessage(a flow.AccountSignature) *entities.AccountSignature {
	return &entities.AccountSignature{
		Account:   a.Account.Bytes(),
		Signature: a.Signature,
	}
}

func MessageToTransaction(m *entities.Transaction) (flow.Transaction, error) {
	if m == nil {
		return flow.Transaction{}, ErrEmptyMessage
	}

	scriptAccounts := make([]flow.Address, len(m.ScriptAccounts))
	for i, account := range m.ScriptAccounts {
		scriptAccounts[i] = flow.BytesToAddress(account)
	}

	signatures := make([]flow.AccountSignature, len(m.Signatures))
	for i, accountSig := range m.Signatures {
		signatures[i] = MessageToAccountSignature(accountSig)
	}

	return flow.Transaction{
		Script:         m.GetScript(),
		PayerAccount:   flow.BytesToAddress(m.PayerAccount),
		ScriptAccounts: scriptAccounts,
		Signatures:     signatures,
		Status:         flow.TransactionStatus(m.GetStatus()),
	}, nil
}

func TransactionToMessage(t flow.Transaction) *entities.Transaction {
	scriptAccounts := make([][]byte, len(t.ScriptAccounts))
	for i, account := range t.ScriptAccounts {
		scriptAccounts[i] = account.Bytes()
	}

	signatures := make([]*entities.AccountSignature, len(t.Signatures))
	for i, accountSig := range t.Signatures {
		signatures[i] = AccountSignatureToMessage(accountSig)
	}

	return &entities.Transaction{
		Script:         t.Script,
		PayerAccount:   t.PayerAccount.Bytes(),
		ScriptAccounts: scriptAccounts,
		Signatures:     signatures,
		Status:         entities.TransactionStatus(t.Status),
	}
}

func MessageToAccount(m *entities.Account) (flow.Account, error) {
	if m == nil {
		return flow.Account{}, ErrEmptyMessage
	}

	accountKeys := make([]flow.AccountPublicKey, len(m.Keys))
	for i, key := range m.Keys {
		accountKey, err := MessageToAccountPublicKey(key)
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
	accountKeys := make([]*entities.AccountPublicKey, len(a.Keys))
	for i, key := range a.Keys {
		accountKeyMsg, err := AccountPublicKeyToMessage(key)
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

func MessageToAccountPublicKey(m *entities.AccountPublicKey) (flow.AccountPublicKey, error) {
	if m == nil {
		return flow.AccountPublicKey{}, ErrEmptyMessage
	}

	signAlgo := crypto.SigningAlgorithm(m.GetSignAlgo())
	hashAlgo := crypto.HashingAlgorithm(m.GetHashAlgo())

	publicKey, err := crypto.DecodePublicKey(signAlgo, m.GetPublicKey())
	if err != nil {
		return flow.AccountPublicKey{}, err
	}

	return flow.AccountPublicKey{
		PublicKey: publicKey,
		SignAlgo:  signAlgo,
		HashAlgo:  hashAlgo,
		Weight:    int(m.GetWeight()),
	}, nil
}

func AccountPublicKeyToMessage(a flow.AccountPublicKey) (*entities.AccountPublicKey, error) {
	publicKey, err := a.PublicKey.Encode()
	if err != nil {
		return nil, err
	}

	return &entities.AccountPublicKey{
		PublicKey: publicKey,
		SignAlgo:  uint32(a.SignAlgo),
		HashAlgo:  uint32(a.HashAlgo),
		Weight:    uint32(a.Weight),
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
		Type:   m.GetType(),
		TxHash: crypto.BytesToHash(m.GetTransactionId()),
		Index:  uint(m.GetIndex()),
		Value:  eventValue,
	}, nil
}

func EventToMessage(e flow.Event) (*entities.Event, error) {
	payload, err := jsoncdc.Encode(e.Value)
	if err != nil {
		return nil, fmt.Errorf("convert: %w", err)
	}

	return &entities.Event{
		Type:          e.Type,
		TransactionId: e.TxHash,
		Index:         uint32(e.Index),
		Payload:       payload,
	}, nil
}

package convert

import (
	"errors"

	"github.com/dapperlabs/flow-go/crypto"

	"github.com/dapperlabs/flow-go-sdk"
	proto "github.com/dapperlabs/flow-go-sdk/client/protobuf/flow"
)

var ErrEmptyMessage = errors.New("protobuf message is empty")

func MessageToBlockHeader(m *proto.BlockHeader) flow.Header {
	return flow.Header{
		Parent: m.GetParentId(),
		Number: m.GetHeight(),
	}
}

func BlockHeaderToMessage(b flow.Header) *proto.BlockHeader {
	return &proto.BlockHeader{
		Id:       b.Hash(),
		ParentId: b.Parent,
		Height:   b.Number,
	}
}

func MessageToTransaction(m *proto.Transaction) (flow.Transaction, error) {
	if m == nil {
		return flow.Transaction{}, ErrEmptyMessage
	}
	return flow.Transaction{}, nil
}

func TransactionToMessage(t flow.Transaction) *proto.Transaction {
	return &proto.Transaction{}
}

func MessageToAccount(m *proto.Account) (flow.Account, error) {
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

func AccountToMessage(a flow.Account) (*proto.Account, error) {
	accountKeys := make([]*proto.AccountPublicKey, len(a.Keys))
	for i, key := range a.Keys {
		accountKeyMsg, err := AccountKeyToMessage(key)
		if err != nil {
			return nil, err
		}
		accountKeys[i] = accountKeyMsg
	}

	return &proto.Account{
		Address: a.Address.Bytes(),
		Balance: a.Balance,
		Code:    a.Code,
		Keys:    accountKeys,
	}, nil
}

func MessageToAccountKey(m *proto.AccountPublicKey) (flow.AccountKey, error) {
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
		PublicKey: publicKey,
		SignAlgo:  signAlgo,
		HashAlgo:  hashAlgo,
		Weight:    int(m.GetWeight()),
	}, nil
}

func AccountKeyToMessage(a flow.AccountKey) (*proto.AccountPublicKey, error) {
	publicKey, err := a.PublicKey.Encode()
	if err != nil {
		return nil, err
	}

	return &proto.AccountPublicKey{
		PublicKey: publicKey,
		SignAlgo:  uint32(a.SignAlgo),
		HashAlgo:  uint32(a.HashAlgo),
		Weight:    uint32(a.Weight),
	}, nil
}

func MessageToEvent(m *proto.Event) flow.Event {
	return flow.Event{
		Type:    m.GetType(),
		TxHash:  crypto.BytesToHash(m.GetTransactionId()),
		Index:   uint(m.GetIndex()),
		Payload: m.GetPayload(),
	}
}

func EventToMessage(e flow.Event) *proto.Event {
	return &proto.Event{
		Type:          e.Type,
		TransactionId: e.TxHash,
		Index:         uint32(e.Index),
		Payload:       e.Payload,
	}
}

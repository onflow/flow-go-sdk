package test

import (
	"fmt"

	"github.com/dapperlabs/cadence"
	"github.com/dapperlabs/flow-go/crypto"

	"github.com/dapperlabs/flow-go-sdk"
	sdkcrypto "github.com/dapperlabs/flow-go-sdk/crypto"
	"github.com/dapperlabs/flow-go-sdk/keys"
)

var ScriptHelloWorld = []byte(`transaction { execute { log("Hello, World!") } }`)

type Identifiers struct {
	count int
}

func IdentifierGenerator() *Identifiers {
	return &Identifiers{0}
}

func (g *Identifiers) New() flow.Identifier {
	id := newIdentifier(g.count + 1)
	g.count++
	return id
}

func newIdentifier(count int) flow.Identifier {
	var id flow.Identifier
	for i := range id {
		id[i] = uint8(count)
	}

	return id
}

type Addresses struct {
	count int
}

func AddressGenerator() *Addresses {
	return &Addresses{0}
}

func (g *Addresses) New() flow.Address {
	addr := flow.BytesToAddress([]byte{uint8(g.count)})
	g.count++
	return addr
}

type AccountKeys struct {
	count int
	ids   *Identifiers
}

func AccountKeyGenerator() *AccountKeys {
	return &AccountKeys{
		count: 0,
		ids:   IdentifierGenerator(),
	}
}

func (g *AccountKeys) New() flow.AccountKey {
	seed := make([]byte, crypto.KeyGenSeedMinLenECDSA_P256)
	for i := range seed {
		seed[i] = uint8(g.count)
	}

	privateKey, err := crypto.GeneratePrivateKey(crypto.ECDSA_P256, seed)

	if err != nil {
		panic(err)
	}

	accountKey := flow.AccountKey{
		ID:             g.count,
		PublicKey:      privateKey.PublicKey(),
		SignAlgo:       crypto.ECDSA_P256,
		HashAlgo:       crypto.SHA3_256,
		Weight:         keys.PublicKeyWeightThreshold,
		SequenceNumber: 42,
	}

	g.count++

	return accountKey
}

type Accounts struct {
	addresses   *Addresses
	accountKeys *AccountKeys
}

func AccountGenerator() *Accounts {
	return &Accounts{
		addresses:   AddressGenerator(),
		accountKeys: AccountKeyGenerator(),
	}
}

func (g *Accounts) New() *flow.Account {
	keyA := g.accountKeys.New()
	keyB := g.accountKeys.New()

	return &flow.Account{
		Address: g.addresses.New(),
		Balance: 10,
		Keys:    []flow.AccountKey{keyA, keyB},
		Code:    nil,
	}
}

type Transactions struct {
	count int
}

func TransactionGenerator() *Transactions {
	return &Transactions{0}
}

func (g *Transactions) New() *flow.Transaction {
	tx := g.NewUnsigned()

	signers := tx.Signers()

	for _, signer := range signers {
		for _, keyID := range signer.KeyIDs {
			_ = tx.Sign(
				signer.SignatureKind(),
				signer.Address,
				keyID,
				sdkcrypto.MockSigner([]byte{uint8(keyID)}),
			)
		}
	}

	return tx
}

func (g *Transactions) NewUnsigned() *flow.Transaction {
	blockID := newIdentifier(g.count + 1)

	accounts := AccountGenerator()
	accountA := accounts.New()
	accountB := accounts.New()

	return flow.NewTransaction().
		SetScript(ScriptHelloWorld).
		SetReferenceBlockID(blockID).
		SetGasLimit(42).
		SetProposalKey(accountA.Address, accountA.Keys[0].ID, accountA.Keys[0].SequenceNumber).
		AddAuthorizer(accountA.Address, accountA.Keys[0].ID, accountA.Keys[1].ID).
		SetPayer(accountB.Address, accountB.Keys[0].ID)
}

type TransactionResults struct {
	events *Events
}

func TransactionResultGenerator() *TransactionResults {
	return &TransactionResults{
		events: EventGenerator(),
	}
}

func (g *TransactionResults) New() flow.TransactionResult {
	eventA := g.events.New()
	eventB := g.events.New()

	return flow.TransactionResult{
		Status: flow.TransactionStatusSealed,
		Events: []flow.Event{
			eventA,
			eventB,
		},
	}
}

type Events struct {
	count int
	ids   *Identifiers
}

func EventGenerator() *Events {
	return &Events{
		count: 0,
		ids:   IdentifierGenerator(),
	}
}

func (g *Events) New() flow.Event {
	identifier := fmt.Sprintf("FooEvent%d", g.count)
	typeID := "test." + identifier

	testEventType := cadence.EventType{
		TypeID:     typeID,
		Identifier: identifier,
		Fields: []cadence.Field{
			{
				Identifier: "a",
				Type:       cadence.IntType{},
			},
			{
				Identifier: "b",
				Type:       cadence.StringType{},
			},
		},
	}

	testEvent := cadence.NewEvent(
		[]cadence.Value{
			cadence.NewInt(g.count),
			cadence.NewString("foo"),
		}).WithType(testEventType)

	event := flow.Event{
		Type:          typeID,
		TransactionID: g.ids.New(),
		Index:         uint(g.count),
		Value:         testEvent,
	}

	g.count++

	return event
}

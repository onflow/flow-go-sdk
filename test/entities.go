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

package test

import (
	"errors"
	"fmt"
	"time"

	"github.com/onflow/cadence"
	"github.com/onflow/cadence/runtime/common"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
)

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
	return &flow.Account{
		Address: g.addresses.New(),
		Balance: 10,
		Keys: []*flow.AccountKey{
			g.accountKeys.New(),
			g.accountKeys.New(),
		},
	}
}

type AccountKeys struct {
	count int
	ids   *Identifiers
}

func AccountKeyGenerator() *AccountKeys {
	return &AccountKeys{
		count: 1,
		ids:   IdentifierGenerator(),
	}
}

func (g *AccountKeys) New() *flow.AccountKey {
	accountKey, _ := g.NewWithSigner()
	return accountKey
}

func (g *AccountKeys) NewWithSigner() (*flow.AccountKey, crypto.Signer) {
	defer func() { g.count++ }()

	seed := make([]byte, crypto.MinSeedLength)
	for i := range seed {
		seed[i] = uint8(g.count)
	}

	privateKey, err := crypto.GeneratePrivateKey(crypto.ECDSA_P256, seed)

	if err != nil {
		panic(err)
	}

	accountKey := flow.AccountKey{
		Index:          g.count,
		PublicKey:      privateKey.PublicKey(),
		SigAlgo:        crypto.ECDSA_P256,
		HashAlgo:       crypto.SHA3_256,
		Weight:         flow.AccountKeyWeightThreshold,
		SequenceNumber: 42,
	}

	return &accountKey, crypto.NewInMemorySigner(privateKey, accountKey.HashAlgo)
}

type Addresses struct {
	generator *flow.AddressGenerator
}

func AddressGenerator() *Addresses {
	return &Addresses{
		generator: flow.NewAddressGenerator(flow.Emulator),
	}
}

func (g *Addresses) New() flow.Address {
	return g.generator.NextAddress()
}

type Blocks struct {
	headers    *BlockHeaders
	guarantees *CollectionGuarantees
}

func BlockGenerator() *Blocks {
	return &Blocks{
		headers:    BlockHeaderGenerator(),
		guarantees: CollectionGuaranteeGenerator(),
	}
}

func (g *Blocks) New() *flow.Block {
	header := g.headers.New()

	guarantees := []*flow.CollectionGuarantee{
		g.guarantees.New(),
		g.guarantees.New(),
		g.guarantees.New(),
	}

	payload := flow.BlockPayload{
		CollectionGuarantees: guarantees,
	}

	return &flow.Block{
		BlockHeader:  header,
		BlockPayload: payload,
	}
}

type BlockHeaders struct {
	count     int
	ids       *Identifiers
	startTime time.Time
}

func BlockHeaderGenerator() *BlockHeaders {
	startTime, _ := time.Parse(time.RFC3339, "2020-06-04T15:43:21+00:00")

	return &BlockHeaders{
		count:     1,
		ids:       IdentifierGenerator(),
		startTime: startTime.UTC(),
	}
}

func (g *BlockHeaders) New() flow.BlockHeader {
	defer func() { g.count++ }()

	return flow.BlockHeader{
		ID:        g.ids.New(),
		ParentID:  g.ids.New(),
		Height:    uint64(g.count),
		Timestamp: g.startTime.Add(time.Hour * time.Duration(g.count)),
	}
}

type Collections struct {
	ids *Identifiers
}

func CollectionGenerator() *Collections {
	return &Collections{
		ids: IdentifierGenerator(),
	}
}

func (g *Collections) New() *flow.Collection {
	return &flow.Collection{
		TransactionIDs: []flow.Identifier{
			g.ids.New(),
			g.ids.New(),
		},
	}
}

type CollectionGuarantees struct {
	ids *Identifiers
}

func CollectionGuaranteeGenerator() *CollectionGuarantees {
	return &CollectionGuarantees{
		ids: IdentifierGenerator(),
	}
}

func (g *CollectionGuarantees) New() *flow.CollectionGuarantee {
	return &flow.CollectionGuarantee{
		CollectionID: g.ids.New(),
	}
}

type Events struct {
	count int
	ids   *Identifiers
}

func EventGenerator() *Events {
	return &Events{
		count: 1,
		ids:   IdentifierGenerator(),
	}
}

func (g *Events) New() flow.Event {
	defer func() { g.count++ }()

	identifier := fmt.Sprintf("FooEvent%d", g.count)

	location := common.StringLocation("test")

	testEventType := &cadence.EventType{
		Location:            location,
		QualifiedIdentifier: identifier,
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

	typeID := location.TypeID(identifier)

	event := flow.Event{
		Type:             string(typeID),
		TransactionID:    g.ids.New(),
		TransactionIndex: g.count,
		EventIndex:       g.count,
		Value:            testEvent,
	}

	return event
}

type Identifiers struct {
	count int
}

func IdentifierGenerator() *Identifiers {
	return &Identifiers{1}
}

func (g *Identifiers) New() flow.Identifier {
	defer func() { g.count++ }()
	return newIdentifier(g.count + 1)
}

func newIdentifier(count int) flow.Identifier {
	var id flow.Identifier
	for i := range id {
		id[i] = uint8(count)
	}

	return id
}

type Transactions struct {
	count     int
	greetings *Greetings
}

func TransactionGenerator() *Transactions {
	return &Transactions{
		count:     1,
		greetings: GreetingGenerator(),
	}
}

func (g *Transactions) New() *flow.Transaction {
	tx := g.NewUnsigned()

	// sign payload with proposal key
	err := tx.SignPayload(
		tx.ProposalKey.Address,
		tx.ProposalKey.KeyIndex,
		MockSigner([]byte{uint8(tx.ProposalKey.KeyIndex)}),
	)
	if err != nil {
		panic(err)
	}

	// sign payload as each authorizer
	for _, addr := range tx.Authorizers {
		err = tx.SignPayload(addr, 0, MockSigner(addr.Bytes()))
		if err != nil {
			panic(err)
		}
	}

	// sign envelope as payer
	err = tx.SignEnvelope(tx.Payer, 0, MockSigner(tx.Payer.Bytes()))
	if err != nil {
		panic(err)
	}

	return tx
}

func (g *Transactions) NewUnsigned() *flow.Transaction {
	blockID := newIdentifier(g.count + 1)

	accounts := AccountGenerator()
	accountA := accounts.New()
	accountB := accounts.New()

	proposalKey := accountA.Keys[0]

	tx := flow.NewTransaction().
		SetScript(GreetingScript).
		SetReferenceBlockID(blockID).
		SetGasLimit(42).
		SetProposalKey(accountA.Address, proposalKey.Index, proposalKey.SequenceNumber).
		AddAuthorizer(accountA.Address).
		SetPayer(accountB.Address)

	err := tx.AddArgument(cadence.NewString(g.greetings.New()))
	if err != nil {
		panic(err)
	}

	return tx
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
		Error:  errors.New("transaction execution error"),
		Events: []flow.Event{
			eventA,
			eventB,
		},
	}
}

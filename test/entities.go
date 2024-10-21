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

package test

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/onflow/cadence"
	"github.com/onflow/cadence/common"
	"github.com/onflow/cadence/encoding/ccf"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow/protobuf/go/flow/entities"

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
	count uint32
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

	// error here is nil since ECDSA_P256 and SHA3_256 are compatible,
	// but keeping the error check for sanity
	signer, err := crypto.NewInMemorySigner(privateKey, accountKey.HashAlgo)
	if err != nil {
		panic(err)
	}

	return &accountKey, signer
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
	seals      *BlockSeals
	signatures *Signatures
}

func BlockGenerator() *Blocks {
	return &Blocks{
		headers:    BlockHeaderGenerator(),
		guarantees: CollectionGuaranteeGenerator(),
		seals:      BlockSealGenerator(),
		signatures: SignaturesGenerator(),
	}
}

func (g *Blocks) New() *flow.Block {
	header := g.headers.New()

	guarantees := []*flow.CollectionGuarantee{
		g.guarantees.New(),
		g.guarantees.New(),
		g.guarantees.New(),
	}

	seals := []*flow.BlockSeal{
		g.seals.New(),
	}

	payload := flow.BlockPayload{
		CollectionGuarantees: guarantees,
		Seals:                seals,
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

type LightCollection struct {
	ids *Identifiers
}

func LightCollectionGenerator() *LightCollection {
	return &LightCollection{
		ids: IdentifierGenerator(),
	}
}

func (g *LightCollection) New() *flow.Collection {
	return &flow.Collection{
		TransactionIDs: []flow.Identifier{
			g.ids.New(),
			g.ids.New(),
		},
	}
}

type FullCollection struct {
	Transactions *Transactions
}

func FullCollectionGenerator() *FullCollection {
	return &FullCollection{
		Transactions: TransactionGenerator(),
	}
}

func (c *FullCollection) New() *flow.FullCollection {
	return &flow.FullCollection{
		Transactions: []*flow.Transaction{c.Transactions.New(), c.Transactions.New()},
	}
}

type CollectionGuarantees struct {
	ids *Identifiers
}

type BlockSeals struct {
	ids   *Identifiers
	sigs  *Signatures
	bytes *Bytes
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

func BlockSealGenerator() *BlockSeals {
	return &BlockSeals{
		ids:   IdentifierGenerator(),
		sigs:  SignaturesGenerator(),
		bytes: BytesGenerator(),
	}
}

func (g *BlockSeals) New() *flow.BlockSeal {
	sigs := []*flow.AggregatedSignature{{
		VerifierSignatures: g.sigs.New(),
		SignerIds:          []flow.Identifier{g.ids.New()},
	}}

	return &flow.BlockSeal{
		BlockID:                    g.ids.New(),
		ExecutionReceiptID:         g.ids.New(),
		ExecutionReceiptSignatures: g.sigs.New(),
		ResultApprovalSignatures:   g.sigs.New(),
		FinalState:                 g.bytes.New(),
		ResultId:                   g.ids.New(),
		AggregatedApprovalSigs:     sigs,
	}
}

type Events struct {
	count    int
	ids      *Identifiers
	encoding entities.EventEncodingVersion
}

func EventGenerator(encoding entities.EventEncodingVersion) *Events {
	return &Events{
		count:    1,
		ids:      IdentifierGenerator(),
		encoding: encoding,
	}
}

func (g *Events) New() flow.Event {
	defer func() { g.count++ }()

	identifier := fmt.Sprintf("FooEvent%d", g.count)

	location := common.StringLocation("test")

	testEventType := cadence.NewEventType(
		location,
		identifier,
		[]cadence.Field{
			{
				Identifier: "a",
				Type:       cadence.IntType,
			},
			{
				Identifier: "b",
				Type:       cadence.StringType,
			},
		},
		nil,
	)

	testEvent := cadence.NewEvent(
		[]cadence.Value{
			cadence.NewInt(g.count),
			cadence.String("foo"),
		}).WithType(testEventType)

	typeID := location.TypeID(nil, identifier)

	var payload []byte
	var err error
	if g.encoding == flow.EventEncodingVersionCCF {
		payload, err = ccf.Encode(testEvent)
	} else {
		payload, err = jsoncdc.Encode(testEvent)
	}

	if err != nil {
		panic(fmt.Errorf("cannot encode test event: %w", err))
	}

	event := flow.Event{
		Type:             string(typeID),
		TransactionID:    g.ids.New(),
		TransactionIndex: g.count,
		EventIndex:       g.count,
		Value:            testEvent,
		Payload:          payload,
	}

	return event
}

type Signatures struct {
	count int
}

func SignaturesGenerator() *Signatures {
	return &Signatures{1}
}

func (g *Signatures) New() [][]byte {
	defer func() { g.count++ }()
	return [][]byte{
		[]byte(strconv.Itoa(g.count + 1)),
	}
}

func newSignatures(count int) Signatures {
	return Signatures{
		count: count,
	}
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
		SetComputeLimit(42).
		SetProposalKey(accountA.Address, proposalKey.Index, proposalKey.SequenceNumber).
		AddAuthorizer(accountA.Address).
		SetPayer(accountB.Address)

	err := tx.AddArgument(cadence.String(g.greetings.New()))
	if err != nil {
		panic(err)
	}

	return tx
}

type TransactionResults struct {
	events *Events
	ids    *Identifiers
}

func TransactionResultGenerator(encoding entities.EventEncodingVersion) *TransactionResults {
	return &TransactionResults{
		events: EventGenerator(encoding),
		ids:    IdentifierGenerator(),
	}
}

func (g *TransactionResults) New() flow.TransactionResult {
	return flow.TransactionResult{
		Status: flow.TransactionStatusSealed,
		Error:  errors.New("transaction execution error"),
		Events: []flow.Event{
			g.events.New(),
			g.events.New(),
		},
		BlockID:          g.ids.New(),
		BlockHeight:      uint64(42),
		TransactionID:    g.ids.New(),
		CollectionID:     g.ids.New(),
		ComputationUsage: uint64(42),
	}
}

type ExecutionDatas struct {
	ids    *Identifiers
	chunks *ChunkExecutionDatas
}

func ExecutionDataGenerator() *ExecutionDatas {
	return &ExecutionDatas{
		ids:    IdentifierGenerator(),
		chunks: ChunkExecutionDataGenerator(),
	}
}

func (g *ExecutionDatas) New() *flow.ExecutionData {
	return &flow.ExecutionData{
		BlockID: g.ids.New(),
		ChunkExecutionData: []*flow.ChunkExecutionData{
			g.chunks.New(),
			g.chunks.New(),
		},
	}
}

type ChunkExecutionDatas struct {
	ids         *Identifiers
	txs         *Transactions
	events      *Events
	trieUpdates *TrieUpdates
	results     *LightTransactionResults
}

func ChunkExecutionDataGenerator() *ChunkExecutionDatas {
	return &ChunkExecutionDatas{
		ids:         IdentifierGenerator(),
		txs:         TransactionGenerator(),
		events:      EventGenerator(flow.EventEncodingVersionCCF),
		trieUpdates: TrieUpdateGenerator(),
		results:     LightTransactionResultGenerator(),
	}
}

func (g *ChunkExecutionDatas) New() *flow.ChunkExecutionData {
	events := make([]*flow.Event, 0, 2)
	for i := 0; i < 2; i++ {
		event := g.events.New()
		events = append(events, &event)
	}

	return &flow.ChunkExecutionData{
		Transactions: []*flow.Transaction{
			g.txs.New(),
			g.txs.New(),
		},
		Events:     events,
		TrieUpdate: g.trieUpdates.New(),
		TransactionResults: []*flow.LightTransactionResult{
			g.results.New(),
			g.results.New(),
		},
	}
}

type TrieUpdates struct {
	ids *Identifiers
}

func TrieUpdateGenerator() *TrieUpdates {
	return &TrieUpdates{
		ids: IdentifierGenerator(),
	}
}

func (g *TrieUpdates) New() *flow.TrieUpdate {
	return &flow.TrieUpdate{
		RootHash: g.ids.New().Bytes(),
		Paths: [][]byte{
			g.ids.New().Bytes(),
		},
		Payloads: []*flow.Payload{
			{
				KeyPart: []*flow.KeyPart{
					{
						Type:  0,
						Value: g.ids.New().Bytes(),
					},
				},
				Value: g.ids.New().Bytes(),
			},
		},
	}
}

type LightTransactionResults struct {
	ids *Identifiers
}

func LightTransactionResultGenerator() *LightTransactionResults {
	return &LightTransactionResults{
		ids: IdentifierGenerator(),
	}
}

func (g *LightTransactionResults) New() *flow.LightTransactionResult {
	return &flow.LightTransactionResult{
		TransactionID:   g.ids.New(),
		Failed:          false,
		ComputationUsed: uint64(42),
	}
}

type Bytes struct {
	count int
}

func BytesGenerator() *Bytes {
	return &Bytes{
		count: 64,
	}
}

func (b *Bytes) New() []byte {
	randomBytes := make([]byte, b.count)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic("failed to generate random bytes")
	}
	return randomBytes
}

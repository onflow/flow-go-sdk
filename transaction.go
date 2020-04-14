package flow

import (
	"sort"

	"github.com/dapperlabs/flow-go/model/hash"

	"github.com/dapperlabs/flow-go-sdk/crypto"
)

// A Transaction is a full transaction object containing a payload and signatures.
type Transaction struct {
	Script             []byte
	ReferenceBlockID   Identifier
	GasLimit           uint64
	ProposalKey        ProposalKey
	Payer              Address
	Authorizers        []Address
	PayloadSignatures  []TransactionSignature
	EnvelopeSignatures []TransactionSignature
}

// NewTransaction initializes and returns an empty transaction.
func NewTransaction() *Transaction {
	return &Transaction{}
}

// ID returns the canonical SHA3-256 hash of this transaction.
func (t *Transaction) ID() Identifier {
	return HashToID(hash.DefaultHasher.ComputeHash(t.Encode()))
}

// SetScript sets the Cadence script for this transaction.
func (t *Transaction) SetScript(script []byte) *Transaction {
	t.Script = script
	return t
}

// SetReferenceBlockID sets the reference block ID for this transaction.
func (t *Transaction) SetReferenceBlockID(blockID Identifier) *Transaction {
	t.ReferenceBlockID = blockID
	return t
}

// SetGasLimit sets the gas limit for this transaction.
func (t *Transaction) SetGasLimit(limit uint64) *Transaction {
	t.GasLimit = limit
	return t
}

// SetProposalKey sets the proposal key and sequence number for this transaction.
//
// The first two arguments specify the account key to be used, and the last argument is the sequence
// number being declared.
func (t *Transaction) SetProposalKey(address Address, keyID int, sequenceNum uint64) *Transaction {
	proposalKey := ProposalKey{
		Address:        address,
		KeyID:          keyID,
		SequenceNumber: sequenceNum,
	}
	t.ProposalKey = proposalKey
	return t
}

// SetPayer sets the payer account for this transaction.
func (t *Transaction) SetPayer(address Address) *Transaction {
	t.Payer = address
	return t
}

// AddAuthorizer adds an authorizer account to this transaction.
func (t *Transaction) AddAuthorizer(address Address) *Transaction {
	t.Authorizers = append(t.Authorizers, address)
	return t
}

// signerList returns a list of unique accounts required to sign this transaction.
//
// The list is returned in the following order:
// 1. PROPOSER
// 2. PAYER
// 2. AUTHORIZERS (in insertion order)
//
// The only exception to the above ordering is for deduplication; if the same account
// is used in multiple signing roles, only the first occurrence is included in the list.
func (t *Transaction) signerList() []Address {
	signers := make([]Address, 0)
	seen := make(map[Address]struct{})

	var addSigner = func(address Address) {
		_, ok := seen[address]
		if ok {
			return
		}

		signers = append(signers, address)
		seen[address] = struct{}{}
	}

	if t.ProposalKey.Address != ZeroAddress {
		addSigner(t.ProposalKey.Address)
	}

	if t.Payer != ZeroAddress {
		addSigner(t.Payer)
	}

	for _, authorizer := range t.Authorizers {
		addSigner(authorizer)
	}

	return signers
}

// signerMap returns a mapping from address to signer index.
func (t *Transaction) signerMap() map[Address]int {
	signers := make(map[Address]int)

	for i, signer := range t.signerList() {
		signers[signer] = i
	}

	return signers
}

// SignPayload signs the transaction payload with the specified account key.
//
// The resulting signature is combined with the account address and key ID before
// being added to the transaction.
//
// This function returns an error if the signature cannot be generated.
func (t *Transaction) SignPayload(address Address, keyID int, signer crypto.Signer) error {
	sig, err := signer.Sign(t.PayloadMessage())
	if err != nil {
		// TODO: wrap error
		return err
	}

	t.AddPayloadSignature(address, keyID, sig)

	return nil
}

// SignEnvelope signs the full transaction (payload + payload signatures) with the specified account key.
//
// The resulting signature is combined with the account address and key ID before
// being added to the transaction.
//
// This function returns an error if the signature cannot be generated.
func (t *Transaction) SignEnvelope(address Address, keyID int, signer crypto.Signer) error {
	sig, err := signer.Sign(t.EnvelopeMessage())
	if err != nil {
		// TODO: wrap error
		return err
	}

	t.AddEnvelopeSignature(address, keyID, sig)

	return nil
}

// AddPayloadSignature adds a payload signature to the transaction for the given address and key ID.
func (t *Transaction) AddPayloadSignature(address Address, keyID int, sig []byte) *Transaction {
	s := t.createSignature(address, keyID, sig)

	t.PayloadSignatures = append(t.PayloadSignatures, s)
	sort.Slice(t.PayloadSignatures, compareSignatures(t.PayloadSignatures))

	return t
}

// AddEnvelopeSignature adds an envelope signature to the transaction for the given address and key ID.
func (t *Transaction) AddEnvelopeSignature(address Address, keyID int, sig []byte) *Transaction {
	s := t.createSignature(address, keyID, sig)

	t.EnvelopeSignatures = append(t.EnvelopeSignatures, s)
	sort.Slice(t.EnvelopeSignatures, compareSignatures(t.EnvelopeSignatures))

	return t
}

func (t *Transaction) createSignature(address Address, keyID int, sig []byte) TransactionSignature {
	signerIndex, signerExists := t.signerMap()[address]
	if !signerExists {
		signerIndex = -1
	}

	return TransactionSignature{
		Address:     address,
		SignerIndex: signerIndex,
		KeyID:       keyID,
		Signature:   sig,
	}
}

func (t *Transaction) PayloadMessage() []byte {
	temp := t.payloadCanonicalForm()
	return DefaultEncoder.MustEncode(&temp)
}

func (t *Transaction) payloadCanonicalForm() interface{} {
	authorizers := make([][]byte, len(t.Authorizers))
	for i, auth := range t.Authorizers {
		authorizers[i] = auth.Bytes()
	}

	return struct {
		Script                    []byte
		ReferenceBlockID          []byte
		GasLimit                  uint64
		ProposalKeyAddress        []byte
		ProposalKeyID             uint64
		ProposalKeySequenceNumber uint64
		Payer                     []byte
		Authorizers               [][]byte
	}{
		t.Script,
		t.ReferenceBlockID[:],
		t.GasLimit,
		t.ProposalKey.Address.Bytes(),
		uint64(t.ProposalKey.KeyID),
		t.ProposalKey.SequenceNumber,
		t.Payer.Bytes(),
		authorizers,
	}
}

// EnvelopeMessage returns the signable message for transaction envelope.
//
// This message is only signed by the payer account.
func (t *Transaction) EnvelopeMessage() []byte {
	temp := t.envelopeCanonicalForm()
	return DefaultEncoder.MustEncode(&temp)
}

func (t *Transaction) envelopeCanonicalForm() interface{} {
	return struct {
		Payload           interface{}
		PayloadSignatures interface{}
	}{
		t.payloadCanonicalForm(),
		signaturesList(t.PayloadSignatures).canonicalForm(),
	}
}

// Encode serializes the full transaction data including the payload and all signatures.
func (t *Transaction) Encode() []byte {
	temp := struct {
		Payload            interface{}
		PayloadSignatures  interface{}
		EnvelopeSignatures interface{}
	}{
		t.payloadCanonicalForm(),
		signaturesList(t.PayloadSignatures).canonicalForm(),
		signaturesList(t.EnvelopeSignatures).canonicalForm(),
	}

	return DefaultEncoder.MustEncode(&temp)
}

// A ProposalKey is the key that specifies the proposal key and sequence number for a transaction.
type ProposalKey struct {
	Address        Address
	KeyID          int
	SequenceNumber uint64
}

// A TransactionSignature is a signature associated with a specific account key.
type TransactionSignature struct {
	Address     Address
	SignerIndex int
	KeyID       int
	Signature   []byte
}

func (s TransactionSignature) canonicalForm() interface{} {
	return struct {
		SignerIndex uint
		KeyID       uint
		Signature   []byte
	}{
		SignerIndex: uint(s.SignerIndex), // int is not RLP-serializable
		KeyID:       uint(s.KeyID),       // int is not RLP-serializable
		Signature:   s.Signature,
	}
}

func compareSignatures(signatures []TransactionSignature) func(i, j int) bool {
	return func(i, j int) bool {
		sigA := signatures[i]
		sigB := signatures[j]
		return sigA.SignerIndex < sigB.SignerIndex || sigA.KeyID < sigB.KeyID
	}
}

type signaturesList []TransactionSignature

func (s signaturesList) canonicalForm() interface{} {
	signatures := make([]interface{}, len(s))

	for i, signature := range s {
		signatures[i] = signature.canonicalForm()
	}

	return signatures
}

type TransactionResult struct {
	Status TransactionStatus
	Error  error
	Events []Event
}

// TransactionStatus represents the status of a transaction.
type TransactionStatus int

const (
	// TransactionStatusUnknown indicates that the transaction status is not known.
	TransactionStatusUnknown TransactionStatus = iota
	// TransactionStatusPending is the status of a pending transaction.
	TransactionStatusPending
	// TransactionStatusFinalized is the status of a finalized transaction.
	TransactionStatusFinalized
	// TransactionStatusExecuted is the status of an executed transaction.
	TransactionStatusExecuted
	// TransactionStatusSealed is the status of a sealed transaction.
	TransactionStatusSealed
)

// String returns the string representation of a transaction status.
func (s TransactionStatus) String() string {
	return [...]string{"UNKNOWN", "PENDING", "FINALIZED", "EXECUTED", "SEALED"}[s]
}

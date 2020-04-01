package flow

import (
	"sort"

	"github.com/dapperlabs/flow-go/model/hash"

	"github.com/dapperlabs/flow-go-sdk/crypto"
)

// A Transaction is a full transaction object containing a payload and payer signatures.
type Transaction struct {
	Payload           TransactionPayload
	PayerSignatureSet *TransactionSignatureSet
}

// NewTransaction initializes and returns an empty transaction.
func NewTransaction() *Transaction {
	return &Transaction{}
}

// ID returns the canonical SHA3-256 hash of this transaction.
func (t *Transaction) ID() Identifier {
	return HashToID(hash.DefaultHasher.ComputeHash(t.Message()))
}

// Script returns the Cadence script for this transaction.
func (t *Transaction) Script() []byte {
	return t.Payload.Script
}

// SetScript sets the Cadence script for this transaction.
func (t *Transaction) SetScript(script []byte) *Transaction {
	t.Payload.Script = script
	return t
}

// ReferenceBlockID returns the reference block ID for this transaction.
func (t *Transaction) ReferenceBlockID() Identifier {
	return t.Payload.ReferenceBlockID
}

// SetReferenceBlockID sets the reference block ID for this transaction.
func (t *Transaction) SetReferenceBlockID(blockID Identifier) *Transaction {
	t.Payload.ReferenceBlockID = blockID
	return t
}

// GasLimit returns the gas limit for this transaction.
func (t *Transaction) GasLimit() uint64 {
	return t.Payload.GasLimit
}

// SetGasLimit sets the gas limit for this transaction.
func (t *Transaction) SetGasLimit(limit uint64) *Transaction {
	t.Payload.GasLimit = limit
	return t
}

// Proposer returns the proposer declaration for this transaction, or nil if it is not set.
func (t *Transaction) Proposer() *SignerDeclaration {
	return t.Payload.proposer
}

// SetProposer sets the proposer account for this transaction.
//
// This function takes an account address and a list of key indices representing the
// account keys that must be used for signing.
func (t *Transaction) SetProposer(address Address, keyIndices ...int) *Transaction {
	t.Payload.proposer = newSignerDeclaration(address, []SignerRole{SignerRoleProposer}, keyIndices...)
	return t
}

// SetProposerSequenceNumber sets the proposal key and sequence number for this transaction.
//
// The first argument is the index of the account key to be used as the proposal key, and the second
// argument is the sequence number of the proposal key.
func (t *Transaction) SetProposerSequenceNumber(proposalKeyIndex int, sequenceNum uint64) *Transaction {
	t.Payload.proposer.SetSequenceNumber(proposalKeyIndex, sequenceNum)
	return t
}

// Payer returns the payer declaration for this transaction, or nil if it is not set.
func (t *Transaction) Payer() *SignerDeclaration {
	return t.Payload.payer
}

// SetPayer sets the payer account for this transaction.
//
// This function takes an account address and a list of key indices representing the
// account keys that must be used for signing.
func (t *Transaction) SetPayer(address Address, keyIndices ...int) *Transaction {
	t.Payload.payer = newSignerDeclaration(address, []SignerRole{SignerRolePayer}, keyIndices...)
	return t
}

// Authorizers returns a list of signer declarations for the accounts that are authorizing
// this transaction.
func (t *Transaction) Authorizers() []*SignerDeclaration {
	return t.Payload.authorizers
}

// AddAuthorizer adds an authorizer account to this transaction.
//
// This function takes an account address and a list of key indices representing the
// account keys that must be used for signing.
func (t *Transaction) AddAuthorizer(address Address, keyIndices ...int) *Transaction {
	sd := newSignerDeclaration(address, []SignerRole{SignerRoleAuthorizer}, keyIndices...)
	t.Payload.authorizers = append(t.Payload.authorizers, sd)
	return t
}

// Signers returns a list of signer declarations for all accounts that are required
// to sign this transaction.
//
// The list is returned in the following order: the proposer is always first, followed
// by payer declaration, and then the authorizer declarations in the order in which they
// were added.
//
// In addition, the resulting list is reduced as following: any two declarations that specify
// the same account and key-set but different signer roles are combined into a single declaration.
// This allows the same account to fulfill multiple (or all) signer roles.
//
// The same account can be used in multiple signer declarations if each declaration specifies
// a unique key-set. For example, you may want to use the same account for payment and authorization
// but require a stricter key-set for payment.
//
// Two key-sets are considered equal if they contain the same key indices, regardless of order.
func (t *Transaction) Signers() []*SignerDeclaration {
	return t.Payload.Signers()
}

// SignPayload signs the transaction payload within the context of specific account key using
// the provided signer.
//
// The resulting signature is combined with the account address and key index before
// being added to the transaction.
//
// This function returns an error if the signature cannot be generated.
func (t *Transaction) SignPayload(address Address, keyIndex int, signer crypto.Signer) error {
	sig, err := signer.Sign(t.Payload)
	if err != nil {
		// TODO: wrap error
		return err
	}

	t.AddPayloadSignature(address, keyIndex, sig)
	return nil
}

// PayloadSignatures returns a list containing a signature set for each account that
// has signed the transaction payload.
//
// The list is returned in the following order: the proposer signature is always first, followed
// by the signatures of the authorizers in the order in which their signer declarations are declared.
func (t *Transaction) PayloadSignatures() []*TransactionSignatureSet {
	return t.Payload.Signatures
}

// AddPayloadSignature adds a payload signature to the transaction for the given address and key index.
func (t *Transaction) AddPayloadSignature(address Address, keyIndex int, sig []byte) *Transaction {
	for _, as := range t.Payload.Signatures {
		if as.Address == address {
			as.Add(keyIndex, sig)
			return t
		}
	}

	ts := newTransactionSignatureSet(address)
	ts.Add(keyIndex, sig)

	t.Payload.Signatures = append(t.Payload.Signatures, ts)

	return t
}

// SignPayer signs the full transaction (payload + payload signatures) within the context of the
// declared payer account and provided key index.
//
// The resulting signature is combined with the payer account address and key index before
// being added to the transaction.
//
// This function returns an error if the signature cannot be generated.
func (t *Transaction) SignPayer(key int, signer crypto.Signer) error {
	sig, err := signer.Sign(t)
	if err != nil {
		// TODO: wrap error
		return err
	}

	t.AddPayerSignature(key, sig)

	return nil
}

// PayerSignatures returns a list of payer signatures for this transaction.
//
// Each signature is associated with a different account key, and the returned
// list is ordered by key index.
func (t *Transaction) PayerSignatures() []TransactionSignature {
	return t.PayerSignatureSet.Signatures
}

// AddPayerSignature adds a payer signature to the transaction for the given key index.
func (t *Transaction) AddPayerSignature(key int, sig []byte) *Transaction {
	if t.PayerSignatureSet == nil {
		if t.Payer() == nil {
			return t
		}

		t.PayerSignatureSet = newTransactionSignatureSet(t.Payer().Address)
	}

	t.PayerSignatureSet.Add(key, sig)

	return t
}

// Message returns the signable message for the full transaction.
//
// This message is only signed by the payer account.
//
// This function conforms to the crypto.Signable interface.
func (t *Transaction) Message() []byte {
	temp := t.messageForm()
	return DefaultEncoder.MustEncode(&temp)
}

func (t *Transaction) messageForm() interface{} {
	if t.PayerSignatureSet == nil {
		return t.Payload.messageForm()
	}

	return struct {
		Payload         interface{}
		PayerSignatures interface{}
	}{
		t.Payload.messageForm(),
		signaturesList(t.PayerSignatureSet.Signatures).messageForm(), // address not included
	}
}

// A TransactionPayload is the inner portion of a transaction that contains the
// script, signers and other metadata required for transaction execution.
type TransactionPayload struct {
	Script           []byte
	ReferenceBlockID Identifier
	GasLimit         uint64

	// Signers
	proposer    *SignerDeclaration
	payer       *SignerDeclaration
	authorizers []*SignerDeclaration

	Signatures []*TransactionSignatureSet
}

// Signers returns a list of signer declarations for all accounts that are required
// to sign this transaction.
//
// The list is returned in the following order: the proposer is always first, followed
// by payer declaration, and then the authorizer declarations in the order in which they
// were added.
//
// In addition, the resulting list is reduced as following: any two declarations that specify
// the same account and key-set but different signer roles are combined into a single declaration.
// This allows the same account to fulfill multiple (or all) signer roles.
//
// The same account can be used in multiple signer declarations if each declaration specifies
// a unique key-set. For example, you may want to use the same account for payment and authorization
// but require a stricter key-set for payment.
//
// Two key-sets are considered equal if they contain the same key indices, regardless of order.
func (t TransactionPayload) Signers() []*SignerDeclaration {
	// TODO: handle case when proposer and/or payer is nil

	// proposer + payer + len(authorizers)
	signerCount := 2 + len(t.authorizers)

	signers := make([]*SignerDeclaration, signerCount)

	signers[0] = t.proposer
	signers[1] = t.payer

	for i, authorizer := range t.authorizers {
		signers[2+i] = authorizer
	}

	signerList := make([]*SignerDeclaration, 0)
	signerMap := make(map[Address]*SignerDeclaration)

	for _, signer := range signers {
		seen, ok := signerMap[signer.Address]
		if ok && signer.hasSameKeysAs(seen) {
			seen.Roles = append(seen.Roles, signer.Roles...)
			continue
		}

		signerMap[signer.Address] = signer
		signerList = append(signerList, signer)
	}

	return signerList
}

// Message returns the signable message for this transaction payload.
//
// This is the portion of the transaction that is signed by the
// proposer and authorizers.
//
// This function conforms to the crypto.Signable interface.
func (t TransactionPayload) Message() []byte {
	temp := t.messageForm()
	return DefaultEncoder.MustEncode(&temp)
}

func (t TransactionPayload) messageForm() interface{} {
	return struct {
		Script           []byte
		ReferenceBlockID []byte
		GasLimit         uint64
		Signers          interface{}
	}{
		t.Script,
		t.ReferenceBlockID[:],
		t.GasLimit,
		signersList(t.Signers()).messageForm(),
	}
}

// A SignerDeclaration specifies an account that is required to sign transaction.
//
// A declaration includes the address of the signer account, the roles
// that it fulfills, and a list of required key indices.
//
// A declaration also specifies an optional proposal key that must be set if
// the signer is fulfilling the PROPOSER role.
type SignerDeclaration struct {
	Address     Address
	Roles       []SignerRole
	Keys        []int
	ProposalKey *ProposalKey
}

// A ProposalKey is the key that specifies the sequence number for a transaction.
type ProposalKey struct {
	KeyIndex       int
	SequenceNumber uint64
}

func newSignerDeclaration(address Address, roles []SignerRole, keyIndices ...int) *SignerDeclaration {
	sortedKeys := make([]int, len(keyIndices))

	for i, key := range keyIndices {
		sortedKeys[i] = key
	}

	sort.Ints(sortedKeys)

	return &SignerDeclaration{
		Address: address,
		Roles:   roles,
		Keys:    sortedKeys,
	}
}

// SetSequenceNumber sets the proposal key and sequence number for this declaration.
func (d *SignerDeclaration) SetSequenceNumber(keyIndex int, sequenceNum uint64) *SignerDeclaration {
	d.ProposalKey = &ProposalKey{
		KeyIndex:       keyIndex,
		SequenceNumber: sequenceNum,
	}
	return d
}

func (d *SignerDeclaration) hasSameKeysAs(other *SignerDeclaration) bool {
	if len(d.Keys) != len(other.Keys) {
		return false
	}

	for i, key := range d.Keys {
		if key != other.Keys[i] {
			return false
		}
	}

	return true
}

func (d SignerDeclaration) messageForm() interface{} {
	if d.ProposalKey != nil {
		return struct {
			Address                   []byte
			Roles                     interface{}
			Keys                      interface{}
			ProposalKeyIndex          uint
			ProposalKeySequenceNumber uint64
		}{
			Address:                   d.Address[:],
			Roles:                     rolesList(d.Roles).messageForm(),
			Keys:                      keysList(d.Keys).messageForm(),
			ProposalKeyIndex:          uint(d.ProposalKey.KeyIndex),
			ProposalKeySequenceNumber: d.ProposalKey.SequenceNumber,
		}
	}

	return struct {
		Address []byte
		Roles   interface{}
		Keys    interface{}
	}{
		Address: d.Address[:],
		Roles:   rolesList(d.Roles).messageForm(),
		Keys:    keysList(d.Keys).messageForm(),
	}
}

// A SignerRole is a role fulfilled by a signer.
type SignerRole int

const (
	// SignerRoleUnknown indicates that the signer role is not known.
	SignerRoleUnknown SignerRole = iota
	// SignerRoleProposer is the role of the transaction proposer.
	SignerRoleProposer
	// SignerRolePayer is the role of the transaction payer.
	SignerRolePayer
	// SignerRoleAuthorizer is the role of a transaction authorizer.
	SignerRoleAuthorizer
)

// String returns the string representation of a signer role.
func (s SignerRole) String() string {
	return [...]string{"UNKNOWN", "PROPOSER", "PAYER", "AUTHORIZER"}[s]
}

func (s SignerRole) messageForm() interface{} {
	return uint(s)
}

// An TransactionSignatureSet is a set of signatures associated with a specific account.
type TransactionSignatureSet struct {
	Address    Address
	Signatures []TransactionSignature
}

func newTransactionSignatureSet(address Address) *TransactionSignatureSet {
	return &TransactionSignatureSet{
		Address:    address,
		Signatures: make([]TransactionSignature, 0),
	}
}

// Add adds a signature to this set for the given key index.
func (s *TransactionSignatureSet) Add(keyIndex int, sig []byte) {
	s.Signatures = signaturesList(s.Signatures).Add(keyIndex, sig)
}

func (s *TransactionSignatureSet) messageForm() interface{} {
	return struct {
		Address    []byte
		Signatures interface{}
	}{
		Address:    s.Address.Bytes(),
		Signatures: signaturesList(s.Signatures).messageForm(),
	}
}

// A TransactionSignature is a signature associated with a specific account key.
type TransactionSignature struct {
	KeyIndex  int
	Signature []byte
}

func (s TransactionSignature) messageForm() interface{} {
	return struct {
		KeyIndex  uint
		Signature []byte
	}{
		KeyIndex:  uint(s.KeyIndex), // int is not RLP-serializable
		Signature: s.Signature,
	}
}

type rolesList []SignerRole

func (l rolesList) messageForm() interface{} {
	roles := make([]interface{}, len(l))

	for i, role := range l {
		roles[i] = role.messageForm()
	}

	return roles
}

type keysList []int

func (l keysList) messageForm() interface{} {
	keys := make([]uint, len(l))

	for i, key := range l {
		keys[i] = uint(key)
	}

	return keys
}

type signersList []*SignerDeclaration

func (l signersList) messageForm() interface{} {
	signers := make([]interface{}, len(l))

	for i, signer := range l {
		signers[i] = signer.messageForm()
	}

	return signers
}

type signaturesList []TransactionSignature

func (s signaturesList) Len() int           { return len(s) }
func (s signaturesList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s signaturesList) Less(i, j int) bool { return s[i].KeyIndex < s[j].KeyIndex }

func (s signaturesList) messageForm() interface{} {
	signatures := make([]interface{}, len(s))

	for i, signature := range s {
		signatures[i] = signature.messageForm()
	}

	return signatures
}

func (s signaturesList) Add(key int, sig []byte) signaturesList {
	ts := TransactionSignature{
		KeyIndex:  key,
		Signature: sig,
	}

	s = append(s, ts)

	// sort list by key index on insertion
	sort.Sort(s)

	return s
}

// TransactionStatus represents the status of a transaction.
type TransactionStatus int

const (
	// TransactionStatusUnknown indicates that the transaction status is not known.
	TransactionStatusUnknown TransactionStatus = iota
	// TransactionPending is the status of a pending transaction.
	TransactionPending
	// TransactionFinalized is the status of a finalized transaction.
	TransactionFinalized
	// TransactionReverted is the status of a reverted transaction.
	TransactionReverted
	// TransactionSealed is the status of a sealed transaction.
	TransactionSealed
)

// String returns the string representation of a transaction status.
func (s TransactionStatus) String() string {
	return [...]string{"UNKNOWN", "PENDING", "FINALIZED", "REVERTED", "SEALED"}[s]
}

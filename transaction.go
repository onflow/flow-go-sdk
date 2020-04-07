package flow

import (
	"sort"

	"github.com/dapperlabs/flow-go/model/hash"

	"github.com/dapperlabs/flow-go-sdk/crypto"
)

// A Transaction is a full transaction object containing a payload and payer signatures.
type Transaction struct {
	Payload    TransactionPayload
	Signatures []TransactionSignature
}

// NewTransaction initializes and returns an empty transaction.
func NewTransaction() *Transaction {
	return &Transaction{}
}

// ID returns the canonical SHA3-256 hash of this transaction.
func (t *Transaction) ID() Identifier {
	return HashToID(hash.DefaultHasher.ComputeHash(t.Encode()))
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

// ProposalKey returns the proposal key declaration for this transaction, or nil if it is not set.
func (t *Transaction) ProposalKey() *ProposalKey {
	return t.Payload.ProposalKey
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

	t.Payload.ProposalKey = &proposalKey
	t.Payload.signersHaveChanged = true

	return t
}

// Payer returns the payer declaration for this transaction, or nil if it is not set.
func (t *Transaction) Payer() *TransactionPayer {
	return t.Payload.Payer
}

// SetPayer sets the payer account for this transaction.
//
// This function takes an account address and a list of key indices representing the
// account keys that must be used for signing.
func (t *Transaction) SetPayer(address Address, keyIDs ...int) *Transaction {
	payer := TransactionPayer{
		Address: address,
		KeyIDs:  keyIDs,
	}

	t.Payload.Payer = &payer
	t.Payload.signersHaveChanged = true

	return t
}

// Authorizers returns a list of signer declarations for the accounts that are authorizing
// this transaction.
func (t *Transaction) Authorizers() []*TransactionAuthorizer {
	return t.Payload.Authorizers
}

// AddAuthorizer adds an authorizer account to this transaction.
//
// This function takes an account address and a list of key indices representing the
// account keys that must be used for signing.
func (t *Transaction) AddAuthorizer(address Address, keyIDs ...int) *Transaction {
	authorizer := TransactionAuthorizer{
		Address: address,
		KeyIDs:  keyIDs,
	}

	t.Payload.Authorizers = append(t.Payload.Authorizers, &authorizer)
	t.Payload.signersHaveChanged = true

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
func (t *Transaction) Signers() []*TransactionSigner {
	return t.Payload.getSigners()
}

// SignPayload signs the transaction payload within the context of an account key.
//
// The resulting signature is combined with the account address and key ID before
// being added to the transaction.
//
// This function returns an error if the signature cannot be generated.
func (t *Transaction) SignPayload(address Address, keyID int, signer crypto.Signer) error {
	sig, err := signer.Sign(t.Payload)
	if err != nil {
		// TODO: wrap error
		return err
	}

	t.AddPayloadSignature(address, keyID, sig)

	return nil
}

// SignContainer signs the full transaction (payload + payload signatures) within the context of an account key.
//
// The resulting signature is combined with the account address and key ID before
// being added to the transaction.
//
// This function returns an error if the signature cannot be generated.
func (t *Transaction) SignContainer(address Address, keyID int, signer crypto.Signer) error {
	sig, err := signer.Sign(t)
	if err != nil {
		// TODO: wrap error
		return err
	}

	t.AddContainerSignature(address, keyID, sig)

	return nil
}

// PayloadSignatures returns a list of signatures of the transaction payload.
//
// The list is returned in the following order: the proposer signature is always first, followed
// by the signatures of the authorizers in the order in which their signer declarations are declared.
func (t *Transaction) PayloadSignatures() []TransactionSignature {
	sigs := make([]TransactionSignature, 0)

	for _, sig := range t.Signatures {
		if sig.Kind == TransactionSignatureKindPayload {
			sigs = append(sigs, sig)
		}
	}

	return sigs
}

// ContainerSignatures returns a list of signatures of the full transaction container.
func (t *Transaction) ContainerSignatures() []TransactionSignature {
	sigs := make([]TransactionSignature, 0)

	for _, sig := range t.Signatures {
		if sig.Kind == TransactionSignatureKindContainer {
			sigs = append(sigs, sig)
		}
	}

	return sigs
}

// AddPayloadSignature adds a payload signature to the transaction for the given address and key ID.
func (t *Transaction) AddPayloadSignature(address Address, keyID int, sig []byte) *Transaction {
	return t.addSignature(TransactionSignatureKindPayload, address, keyID, sig)
}

// AddContainerSignature adds a container signature to the transaction for the given address and key ID.
func (t *Transaction) AddContainerSignature(address Address, keyID int, sig []byte) *Transaction {
	return t.addSignature(TransactionSignatureKindContainer, address, keyID, sig)
}

func (t *Transaction) AddSignatureAtIndex(index int, sig []byte) *Transaction {
	sr := t.Payload.getSignatureRequirementByIndex(index)
	return t.addSignature(sr.Kind, sr.Address, sr.KeyID, sig)
}

func (t *Transaction) addSignature(
	kind TransactionSignatureKind,
	address Address,
	keyID int,
	sig []byte,
) *Transaction {
	index := t.Payload.getSignatureIndex(address, keyID)

	s := TransactionSignature{
		Index:     index,
		Kind:      kind,
		Address:   address,
		KeyID:     keyID,
		Signature: sig,
	}

	t.Signatures = append(t.Signatures, s)

	sort.Slice(t.Signatures, func(i, j int) bool {
		return t.Signatures[i].Index < t.Signatures[j].Index
	})

	return t
}

// Message returns the signable message for the full transaction.
//
// This message is only signed by the payer account.
//
// This function conforms to the crypto.Signable interface.
func (t *Transaction) Message() []byte {
	temp := t.canonicalForm()
	return DefaultEncoder.MustEncode(&temp)
}

func (t *Transaction) canonicalForm() interface{} {
	return struct {
		Payload    interface{}
		Signatures interface{}
	}{
		t.Payload.canonicalForm(),
		signaturesList(t.PayloadSignatures()).canonicalForm(),
	}
}

// Encode serializes the full transaction data including the payload and all signatures.
func (t *Transaction) Encode() []byte {
	temp := struct {
		Payload    interface{}
		Signatures interface{}
	}{
		t.Payload.canonicalForm(),
		signaturesList(t.Signatures).canonicalForm(),
	}

	return DefaultEncoder.MustEncode(&temp)
}

// A TransactionPayload is the inner portion of a transaction that contains the
// script, signers and other metadata required for transaction execution.
type TransactionPayload struct {
	Script           []byte
	ReferenceBlockID Identifier
	GasLimit         uint64
	ProposalKey      *ProposalKey
	Payer            *TransactionPayer
	Authorizers      []*TransactionAuthorizer

	// fields used to cache signer list
	signers                   []*TransactionSigner
	signersHaveChanged        bool
	signatureRequirements     []*TransactionSignatureRequirement
	signatureRequirementTable sigReqLookupTable
}

// getSigners returns a list of signer declarations for all accounts that are required
// to sign this transaction.
//
// The list is returned in the following order:
// 1. PROPOSER declaration
// 2. AUTHORIZER declarations (in insertion order)
// 3. PAYER declaration
//
// In addition, the resulting list is reduced as following:
// 1. PROPOSER can be merged into any declaration D if PROPOSER.PROPOSAL_KEY exists in D.KEYS
// 2. Any declaration D can be merged into PAYER if D.KEYS is a subset of PAYER.KEYS
//
// The same account can be used in multiple signer declarations under these conditions:
// 1. An account cannot exist in two declarations that fulfill the same role
// 2. An account cannot exist in two declarations if either declaration's key-set is a subset of the other
func (t TransactionPayload) getSigners() []*TransactionSigner {
	if t.signers != nil && !t.signersHaveChanged {
		return t.signers
	}

	var (
		proposer *TransactionSigner
		payer    *TransactionSigner
	)

	if t.ProposalKey != nil {
		proposer = newTransactionSigner(SignerRoleProposer, t.ProposalKey.Address, t.ProposalKey.KeyID)
		proposer.ProposalKey = t.ProposalKey
	}

	if t.Payer != nil {
		payer = newTransactionSigner(SignerRolePayer, t.Payer.Address, t.Payer.KeyIDs...)

		if payer.canMergeWith(proposer) {
			payer.mergeWith(proposer)
			*proposer = *payer
			payer = proposer
		}
	}

	signers := make([]*TransactionSigner, 0)

	if proposer != payer {
		signers = append(signers, proposer)
	}

	for _, authorizer := range t.Authorizers {
		auth := newTransactionSigner(SignerRoleAuthorizer, authorizer.Address, authorizer.KeyIDs...)

		// If authorizer key-set is a subset of payer key-set, merge with payer.
		// If proposer key-set is a subset of authorizer key-set, merge with proposer.
		// Otherwise, append authorizer to signer list.
		if payer != nil && payer.canMergeWith(auth) {
			payer.mergeWith(auth)
		} else if proposer != nil && auth.canMergeWith(proposer) {
			auth.mergeWith(proposer)
			*proposer = *auth
		} else {
			signers = append(signers, auth)
		}
	}

	signers = append(signers, payer)

	t.signers = signers
	t.signersHaveChanged = false

	return signers
}

type sigReqKey struct {
	address Address
	keyID   int
}

type sigReqLookupTable map[sigReqKey]*TransactionSignatureRequirement

func (t TransactionPayload) getSignatureRequirements() ([]*TransactionSignatureRequirement, sigReqLookupTable) {
	if t.signatureRequirements != nil && !t.signersHaveChanged {
		return t.signatureRequirements, t.signatureRequirementTable
	}

	signers := t.getSigners()

	signatureRequirements := make([]*TransactionSignatureRequirement, 0)
	signatureRequirementTable := make(sigReqLookupTable)

	i := 0

	for _, signer := range signers {
		for _, keyID := range signer.KeyIDs {
			sr := &TransactionSignatureRequirement{
				Index:   i,
				Address: signer.Address,
				KeyID:   keyID,
				Kind:    signer.signatureKind(),
			}

			signatureRequirementTable[sigReqKey{signer.Address, keyID}] = sr
			signatureRequirements = append(signatureRequirements, sr)

			i++
		}
	}

	t.signatureRequirements = signatureRequirements
	t.signatureRequirementTable = signatureRequirementTable

	return signatureRequirements, signatureRequirementTable
}

//
func (t TransactionPayload) getSignatureIndex(address Address, keyID int) int {
	_, signatureRequirementTable := t.getSignatureRequirements()

	sr := signatureRequirementTable[sigReqKey{address, keyID}]

	if sr == nil {
		return -1
	}

	return sr.Index
}

func (t TransactionPayload) getSignatureRequirementByIndex(index int) TransactionSignatureRequirement {
	signatureRequirements, _ := t.getSignatureRequirements()

	if index >= len(signatureRequirements) {
		return TransactionSignatureRequirement{}
	}

	return *signatureRequirements[index]
}

// Message returns the signable message for this transaction payload.
//
// This is the portion of the transaction that is signed by the
// proposer and authorizers.
//
// This function conforms to the crypto.Signable interface.
func (t TransactionPayload) Message() []byte {
	temp := t.canonicalForm()
	return DefaultEncoder.MustEncode(&temp)
}

func (t TransactionPayload) canonicalForm() interface{} {
	return struct {
		Script           []byte
		ReferenceBlockID []byte
		GasLimit         uint64
		Signers          interface{}
	}{
		t.Script,
		t.ReferenceBlockID[:],
		t.GasLimit,
		signersList(t.getSigners()).canonicalForm(),
	}
}

// A ProposalKey is the key that specifies the proposal key and sequence number for a transaction.
type ProposalKey struct {
	Address        Address
	KeyID          int
	SequenceNumber uint64
}

// A TransactionPayer specifies the account that is paying for a transaction and the
// keys required to sign.
type TransactionPayer struct {
	Address Address
	KeyIDs  []int
}

// A TransactionAuthorizer specifies an account that is authorizing a transaction and the
// keys required to sign.
type TransactionAuthorizer struct {
	Address Address
	KeyIDs  []int
}

// A TransactionSigner specifies an account that is required to sign transaction.
//
// A declaration includes the address of the signer account, the roles
// that it fulfills, and a list of required key indices.
//
// A declaration also specifies an optional proposal key that must be set if
// the signer is fulfilling the PROPOSER role.
type TransactionSigner struct {
	Address     Address
	Roles       []SignerRole
	KeyIDs      []int
	ProposalKey *ProposalKey
}

func newTransactionSigner(role SignerRole, address Address, keyIDs ...int) *TransactionSigner {
	sortedKeys := make([]int, len(keyIDs))

	for i, key := range keyIDs {
		sortedKeys[i] = key
	}

	sort.Ints(sortedKeys)

	return &TransactionSigner{
		Address: address,
		Roles:   []SignerRole{role},
		KeyIDs:  sortedKeys,
	}
}

func (d *TransactionSigner) canMergeWith(other *TransactionSigner) bool {
	if other == nil {
		return false
	}

	// can only merge with same account
	if d.Address != other.Address {
		return false
	}

	// cannot merge with an empty declaration
	if len(other.KeyIDs) == 0 {
		return false
	}

	// create lookup table for keys
	keys := make(map[int]struct{})
	for _, key := range d.KeyIDs {
		keys[key] = struct{}{}
	}

	// other can be merged into this declaration if its key-set
	// is a subset of this declaration's key-set
	for _, key := range other.KeyIDs {
		_, ok := keys[key]
		if !ok {
			return false
		}
	}

	return true
}

func (d *TransactionSigner) mergeWith(other *TransactionSigner) *TransactionSigner {
	d.Roles = append(d.Roles, other.Roles...)

	// sort roles list in following order:
	// 1 - PROPOSER
	// 2 - PAYER
	// 3 - AUTHORIZER
	sort.Slice(d.Roles, func(i, j int) bool {
		return d.Roles[i] < d.Roles[j]
	})

	// when merging, incoming key-set is always a subset of the current key-set, therefore
	// the current key-set does not change

	// copy the proposal key from the incoming declaration
	if other.ProposalKey != nil {
		d.ProposalKey = other.ProposalKey
	}

	return d
}

// signatureKind returns the portion of the transaction this signer is required to sign.
//
// If one of the signer's roles is PAYER, it must sign the CONTAINER.
// Otherwise, it must sign the PAYLOAD.
func (d TransactionSigner) signatureKind() TransactionSignatureKind {
	for _, role := range d.Roles {
		if role == SignerRolePayer {
			return TransactionSignatureKindContainer
		}
	}

	return TransactionSignatureKindPayload
}

func (d TransactionSigner) canonicalForm() interface{} {
	if d.ProposalKey != nil {
		return struct {
			Address                   []byte
			ProposalKeyID             uint
			ProposalKeySequenceNumber uint64
		}{
			Address:                   d.Address[:],
			ProposalKeyID:             uint(d.ProposalKey.KeyID),
			ProposalKeySequenceNumber: d.ProposalKey.SequenceNumber,
		}
	}

	return struct {
		Address []byte
		Roles   interface{}
		KeyIDs  interface{}
	}{
		Address: d.Address[:],
		Roles:   rolesList(d.Roles).canonicalForm(),
		KeyIDs:  keysList(d.KeyIDs).canonicalForm(),
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

func (s SignerRole) canonicalForm() interface{} {
	return uint(s)
}

// A TransactionSignatureKind is a kind of transaction signature.
type TransactionSignatureKind int

const (
	// TransactionSignatureKindUnknown indicates that the signature kind is not known.
	TransactionSignatureKindUnknown TransactionSignatureKind = iota
	// TransactionSignatureKindPayload is a signature of the transaction payload.
	TransactionSignatureKindPayload
	// TransactionSignatureKindContainer is a signature of the full transaction container.
	TransactionSignatureKindContainer
)

// String returns the string representation of a signer role.
func (s TransactionSignatureKind) String() string {
	return [...]string{"UNKNOWN", "PAYLOAD", "CONTAINER"}[s]
}

func (s TransactionSignatureKind) canonicalForm() interface{} {
	return uint(s)
}

// A TransactionSignature is a signature associated with a specific account key.
type TransactionSignature struct {
	Kind      TransactionSignatureKind
	Index     int
	Address   Address
	KeyID     int
	Signature []byte
}

func (s TransactionSignature) canonicalForm() interface{} {
	return struct {
		Index     uint
		Signature []byte
	}{
		Index:     uint(s.Index), // int is not RLP-serializable
		Signature: s.Signature,
	}
}

// A TransactionSignatureRequirement is specifies a signature that is required
// to form a valid transaction.
type TransactionSignatureRequirement struct {
	Index   int
	Address Address
	KeyID   int
	Kind    TransactionSignatureKind
}

type rolesList []SignerRole

func (l rolesList) canonicalForm() interface{} {
	roles := make([]interface{}, len(l))

	for i, role := range l {
		roles[i] = role.canonicalForm()
	}

	return roles
}

type keysList []int

func (l keysList) canonicalForm() interface{} {
	keys := make([]uint, len(l))

	for i, key := range l {
		keys[i] = uint(key)
	}

	return keys
}

type signersList []*TransactionSigner

func (l signersList) canonicalForm() interface{} {
	signers := make([]interface{}, len(l))

	for i, signer := range l {
		signers[i] = signer.canonicalForm()
	}

	return signers
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
	Events []Event
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
	// TransactionExecuted is the status of an executed transaction.
	TransactionExecuted
	// TransactionSealed is the status of a sealed transaction.
	TransactionSealed
)

// String returns the string representation of a transaction status.
func (s TransactionStatus) String() string {
	return [...]string{"UNKNOWN", "PENDING", "FINALIZED", "REVERTED", "SEALED"}[s]
}

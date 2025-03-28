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

package flow

import (
	"github.com/pkg/errors"

	"github.com/onflow/flow-go-sdk/crypto"
)

// An Account is an account on the Flow network.
type Account struct {
	Address   Address
	Balance   uint64
	Code      []byte
	Keys      []*AccountKey
	Contracts map[string][]byte
}

// AccountKeyWeightThreshold is the total key weight required to authorize access to an account.
const AccountKeyWeightThreshold int = 1000

// An AccountKey is a public key associated with an account.
type AccountKey struct {
	Index          uint32
	PublicKey      crypto.PublicKey
	SigAlgo        crypto.SignatureAlgorithm
	HashAlgo       crypto.HashAlgorithm
	Weight         int
	SequenceNumber uint64
	Revoked        bool
}

// NewAccountKey returns an empty account key.
func NewAccountKey() *AccountKey {
	return &AccountKey{}
}

// FromPrivateKey sets the public key and signature algorithm based on the provided private key.
func (a *AccountKey) FromPrivateKey(privKey crypto.PrivateKey) *AccountKey {
	a.PublicKey = privKey.PublicKey()
	a.SigAlgo = privKey.Algorithm()
	return a
}

// SetPublicKey sets the public key for this account key.
func (a *AccountKey) SetPublicKey(pubKey crypto.PublicKey) *AccountKey {
	a.PublicKey = pubKey
	a.SigAlgo = pubKey.Algorithm()
	return a
}

// SetSigAlgo sets the signature algorithm for this account key.
func (a *AccountKey) SetSigAlgo(sigAlgo crypto.SignatureAlgorithm) *AccountKey {
	a.SigAlgo = sigAlgo
	return a
}

// SetHashAlgo sets the hash algorithm for this account key.
func (a *AccountKey) SetHashAlgo(hashAlgo crypto.HashAlgorithm) *AccountKey {
	a.HashAlgo = hashAlgo
	return a
}

// SetWeight sets the weight for this account key.
func (a *AccountKey) SetWeight(weight int) *AccountKey {
	a.Weight = weight
	return a
}

// Encode returns the canonical RLP byte representation of this account key.
func (a AccountKey) Encode() []byte {
	temp := accountKeyWrapper{
		EncodedPublicKey: a.PublicKey.Encode(),
		SigAlgo:          uint(a.SigAlgo),
		HashAlgo:         uint(a.HashAlgo),
		Weight:           uint(a.Weight),
	}
	return mustRLPEncode(&temp)
}

// accountCompatibleAlgorithms returns true if the signature and hash algorithms are a valid pair
// for a key on a Flow account.
func accountCompatibleAlgorithms(sigAlgo crypto.SignatureAlgorithm, hashAlgo crypto.HashAlgorithm) bool {
	return (sigAlgo == crypto.ECDSA_P256 || sigAlgo == crypto.ECDSA_secp256k1) &&
		(hashAlgo == crypto.SHA2_256 || hashAlgo == crypto.SHA3_256)
}

// Validate returns an error if this account key is invalid.
//
// An account key can be invalid for the following reasons:
// - It specifies an incompatible signature/hash algorithm pair with regards to Flow accounts
// - It specifies an invalid key weight
func (a AccountKey) Validate() error {
	if !accountCompatibleAlgorithms(a.SigAlgo, a.HashAlgo) {
		return errors.Errorf(
			"signing algorithm (%s) and hashing algorithm (%s) are not a valid pair for a Flow account key",
			a.SigAlgo,
			a.HashAlgo,
		)
	}

	if a.Weight < 0 || a.Weight > AccountKeyWeightThreshold {
		return errors.Errorf("invalid key weight: %d", a.Weight)
	}

	return nil
}

// DecodeAccountKey decodes the RLP byte representation of an account key.
func DecodeAccountKey(b []byte) (*AccountKey, error) {
	var temp accountKeyWrapper

	err := rlpDecode(b, &temp)
	if err != nil {
		return nil, err
	}

	sigAlgo := crypto.SignatureAlgorithm(temp.SigAlgo)
	hashAlgo := crypto.HashAlgorithm(temp.HashAlgo)

	publicKey, err := crypto.DecodePublicKey(sigAlgo, temp.EncodedPublicKey)
	if err != nil {
		return nil, err
	}

	return &AccountKey{
		PublicKey: publicKey,
		SigAlgo:   sigAlgo,
		HashAlgo:  hashAlgo,
		Weight:    int(temp.Weight),
	}, nil
}

type accountKeyWrapper struct {
	EncodedPublicKey []byte
	SigAlgo          uint
	HashAlgo         uint
	Weight           uint
}

type AccountStatusFilter struct {
	EventFilter
}

type AccountStatus struct {
	BlockID      Identifier
	BlockHeight  uint64
	MessageIndex uint64
	Results      []*AccountStatusResult
}

func (a AccountStatus) GetMessageIndex() uint64 {
	return a.MessageIndex
}

type AccountStatusResult struct {
	Address Address
	Events  []Event
}

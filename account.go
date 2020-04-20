package flow

import (
	"github.com/pkg/errors"

	"github.com/onflow/flow-go-sdk/crypto"
)

// An Account is an account on the Flow network.
//
// An account can be an externally owned account or a contract account with code.
type Account struct {
	Address Address
	Balance uint64
	Code    []byte
	Keys    []*AccountKey
}

const AccountKeyWeightThreshold int = 1000

// An AccountKey is a public key associated with an account.
type AccountKey struct {
	ID             int
	PublicKey      crypto.PublicKey
	SigAlgo        crypto.SignatureAlgorithm
	HashAlgo       crypto.HashAlgorithm
	Weight         int
	SequenceNumber uint64
}

func NewAccountKey() *AccountKey {
	return &AccountKey{}
}

func (a *AccountKey) FromPrivateKey(pk crypto.PrivateKey) *AccountKey {
	a.PublicKey = pk.PublicKey()
	a.SigAlgo = pk.Algorithm()
	return a
}

func (a *AccountKey) SetPublicKey(pubKey crypto.PublicKey) *AccountKey {
	a.PublicKey = pubKey
	return a
}

func (a *AccountKey) SetSigAlgo(sigAlgo crypto.SignatureAlgorithm) *AccountKey {
	a.SigAlgo = sigAlgo
	return a
}

func (a *AccountKey) SetHashAlgo(hashAlgo crypto.HashAlgorithm) *AccountKey {
	a.HashAlgo = hashAlgo
	return a
}

func (a *AccountKey) SetWeight(weight int) *AccountKey {
	a.Weight = weight
	return a
}

func (a AccountKey) Encode() []byte {
	temp := accountPublicKeyWrapper{
		EncodedPublicKey: a.PublicKey.Encode(),
		SigAlgo:          uint(a.SigAlgo),
		HashAlgo:         uint(a.HashAlgo),
		Weight:           uint(a.Weight),
	}
	return mustRLPEncode(&temp)
}

func (a AccountKey) Validate() error {
	if !crypto.CompatibleAlgorithms(a.SigAlgo, a.HashAlgo) {
		return errors.Errorf(
			"signing algorithm (%s) is incompatible with hashing algorithm (%s)",
			a.SigAlgo,
			a.HashAlgo,
		)
	}
	return nil
}

func DecodeAccountKey(b []byte) (*AccountKey, error) {
	var temp accountPublicKeyWrapper

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

type accountPublicKeyWrapper struct {
	EncodedPublicKey []byte
	SigAlgo          uint
	HashAlgo         uint
	Weight           uint
}

package flow

import (
	"github.com/onflow/flow-go-sdk/crypto"
)

// Account represents an account on the Flow network.
//
// An account can be an externally owned account or a contract account with code.
type Account struct {
	Address Address
	Balance uint64
	Code    []byte
	Keys    []AccountKey
}

const AccountKeyWeightThreshold int = 1000

// An AccountKey is a public key associated with an account.
type AccountKey struct {
	ID             int
	PublicKey      crypto.PublicKey
	SignAlgo       crypto.SignatureAlgorithm
	HashAlgo       crypto.HashAlgorithm
	Weight         int
	SequenceNumber uint64
}

type AccountPrivateKey struct {
	PrivateKey crypto.PrivateKey
	SignAlgo   crypto.SignatureAlgorithm
	HashAlgo   crypto.HashAlgorithm
}

func (pk AccountPrivateKey) PublicKey() crypto.PublicKey {
	return pk.PrivateKey.PublicKey()
}

func (pk AccountPrivateKey) ToAccountKey() AccountKey {
	return AccountKey{
		PublicKey: pk.PublicKey(),
		SignAlgo:  pk.SignAlgo,
		HashAlgo:  pk.HashAlgo,
		Weight:    0,
	}
}

func (pk AccountPrivateKey) Signer() crypto.NaiveSigner {
	return crypto.NewNaiveSigner(pk.PrivateKey, pk.HashAlgo)
}

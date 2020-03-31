package flow

import (
	"github.com/dapperlabs/flow-go/crypto"
)

// Account represents an account on the Flow network.
//
// An account can be an externally owned account or a contract account with code.
type Account struct {
	Address Address
	Balance uint64
	Code    []byte
	Keys    []AccountPublicKey
}

// AccountPublicKey is a public key associated with an account.
//
// An account public key contains the public key, signing and hashing algorithms, and a key weight.
type AccountPublicKey struct {
	PublicKey      crypto.PublicKey
	Index          int
	SignAlgo       crypto.SigningAlgorithm
	HashAlgo       crypto.HashingAlgorithm
	Weight         int
	SequenceNumber int
}

// AccountPrivateKey is a private key associated with an account.
type AccountPrivateKey struct {
	PrivateKey crypto.PrivateKey
	SignAlgo   crypto.SigningAlgorithm
	HashAlgo   crypto.HashingAlgorithm
}

// PublicKey returns a weighted public key.
func (a AccountPrivateKey) PublicKey(weight int) AccountPublicKey {
	return AccountPublicKey{
		PublicKey: a.PrivateKey.PublicKey(),
		SignAlgo:  a.SignAlgo,
		HashAlgo:  a.HashAlgo,
		Weight:    weight,
	}
}

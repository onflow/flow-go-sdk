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
	Keys    []AccountKey
}

// An AccountKey is a public key associated with an account.
type AccountKey struct {
	PublicKey      crypto.PublicKey
	Index          int
	SignAlgo       crypto.SigningAlgorithm
	HashAlgo       crypto.HashingAlgorithm
	Weight         int
	SequenceNumber uint64
}

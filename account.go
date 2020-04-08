package flow

import (
	"github.com/dapperlabs/flow-go/crypto"

	sdkcrypto "github.com/dapperlabs/flow-go-sdk/crypto"
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
	ID             int
	PublicKey      crypto.PublicKey
	SignAlgo       crypto.SigningAlgorithm
	HashAlgo       crypto.HashingAlgorithm
	Weight         int
	SequenceNumber uint64
}

type AccountPrivateKey struct {
	PrivateKey crypto.PrivateKey
	SignAlgo   crypto.SigningAlgorithm
	HashAlgo   crypto.HashingAlgorithm
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

func (pk AccountPrivateKey) Signer() sdkcrypto.NaiveSigner {
	return sdkcrypto.NewNaiveSigner(pk.PrivateKey, pk.HashAlgo)
}

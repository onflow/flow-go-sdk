package flow

import (
	"github.com/onflow/flow-go-sdk/crypto"
)

// An Account is an account on the Flow network.
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
	SigAlgo        crypto.SignatureAlgorithm
	HashAlgo       crypto.HashAlgorithm
	Weight         int
	SequenceNumber uint64
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

func DecodeAccountKey(b []byte) (AccountKey, error) {
	var temp accountPublicKeyWrapper

	err := rlpDecode(b, &temp)
	if err != nil {
		return AccountKey{}, nil
	}

	sigAlgo := crypto.SignatureAlgorithm(temp.SigAlgo)
	hashAlgo := crypto.HashAlgorithm(temp.HashAlgo)

	publicKey, err := crypto.DecodePublicKey(sigAlgo, temp.EncodedPublicKey)
	if err != nil {
		return AccountKey{}, nil
	}

	return AccountKey{
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

type AccountPrivateKey struct {
	PrivateKey crypto.PrivateKey
	SigAlgo    crypto.SignatureAlgorithm
	HashAlgo   crypto.HashAlgorithm
}

func (pk AccountPrivateKey) PublicKey() crypto.PublicKey {
	return pk.PrivateKey.PublicKey()
}

func (pk AccountPrivateKey) Signer() crypto.NaiveSigner {
	return crypto.NewNaiveSigner(pk.PrivateKey, pk.HashAlgo)
}

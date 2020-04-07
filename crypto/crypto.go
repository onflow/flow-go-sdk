package crypto

import (
	"github.com/dapperlabs/flow-go/crypto"
)

type Signable interface {
	Message() []byte
}

type Signer interface {
	Sign(obj Signable) ([]byte, error)
}

type PrivateKey crypto.PrivateKey

type NaiveSigner struct {
	PrivateKey crypto.PrivateKey
	Hasher     crypto.Hasher
}

func NewNaiveSigner(privateKey crypto.PrivateKey, hashAlgo crypto.HashingAlgorithm) NaiveSigner {
	hasher, _ := crypto.NewHasher(hashAlgo)

	return NaiveSigner{
		PrivateKey: privateKey,
		Hasher:     hasher,
	}
}

func (s NaiveSigner) Sign(obj Signable) ([]byte, error) {
	return s.PrivateKey.Sign(obj.Message(), s.Hasher)
}

type MockSigner []byte

func (s MockSigner) Sign(Signable) ([]byte, error) {
	return s, nil
}

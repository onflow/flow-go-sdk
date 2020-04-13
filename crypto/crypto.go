package crypto

import (
	"github.com/dapperlabs/flow-go/crypto"
)

type Signer interface {
	Sign(message []byte) ([]byte, error)
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

func (s NaiveSigner) Sign(message []byte) ([]byte, error) {
	return s.PrivateKey.Sign(message, s.Hasher)
}

type MockSigner []byte

func (s MockSigner) Sign(message []byte) ([]byte, error) {
	return s, nil
}

package keys

import (
	"encoding/hex"
	"errors"

	"github.com/dapperlabs/flow-go/crypto"

	"github.com/dapperlabs/flow-go-sdk"
)

// EncodePrivateKey encodes a private key as bytes.
func EncodePrivateKey(a flow.AccountPrivateKey) ([]byte, error) {
	privateKey, err := a.PrivateKey.Encode()
	if err != nil {
		return nil, err
	}

	w := accountPrivateKeyWrapper{
		EncodedPrivateKey: privateKey,
		SignAlgo:          uint(a.SignAlgo),
		HashAlgo:          uint(a.HashAlgo),
	}

	return flow.DefaultEncoder.Encode(&w)
}

// DecodePrivateKey decodes a private key.
func DecodePrivateKey(b []byte) (a flow.AccountPrivateKey, err error) {
	var w accountPrivateKeyWrapper

	err = flow.DefaultEncoder.Decode(b, &w)
	if err != nil {
		return a, err
	}

	signAlgo := crypto.SigningAlgorithm(w.SignAlgo)
	hashAlgo := crypto.HashingAlgorithm(w.HashAlgo)

	privateKey, err := crypto.DecodePrivateKey(signAlgo, w.EncodedPrivateKey)
	if err != nil {
		return a, err
	}

	return flow.AccountPrivateKey{
		PrivateKey: privateKey,
		SignAlgo:   signAlgo,
		HashAlgo:   hashAlgo,
	}, nil
}

// DecodePrivateKeyHex decodes a private key from a hexadecimal string.
func DecodePrivateKeyHex(h string) (flow.AccountPrivateKey, error) {
	b, err := hex.DecodeString(h)
	if err != nil {
		return flow.AccountPrivateKey{}, errors.New("failed to decode hex")
	}

	a, err := DecodePrivateKey(b)
	if err != nil {
		return flow.AccountPrivateKey{}, errors.New("failed to decode private key bytes")
	}

	return a, nil
}

// MustDecodePrivateKeyHex is the same as DecodePrivateKeyHex but panics if the
// input string does not represent a valid private key.
func MustDecodePrivateKeyHex(h string) flow.AccountPrivateKey {
	k, err := DecodePrivateKeyHex(h)
	if err != nil {
		panic(err)
	}
	return k
}

// EncodePublicKey encodes a public key as bytes.
func EncodePublicKey(a flow.AccountKey) ([]byte, error) {
	publicKey, err := a.PublicKey.Encode()
	if err != nil {
		return nil, err
	}

	temp := accountPublicKeyWrapper{
		EncodedPublicKey: publicKey,
		SignAlgo:         uint(a.SignAlgo),
		HashAlgo:         uint(a.HashAlgo),
		Weight:           uint(a.Weight),
	}

	return flow.DefaultEncoder.Encode(&temp)
}

// DecodePublicKey decodes a public key.
func DecodePublicKey(b []byte) (a flow.AccountKey, err error) {
	var temp accountPublicKeyWrapper

	err = flow.DefaultEncoder.Decode(b, &temp)
	if err != nil {
		return a, err
	}

	signAlgo := crypto.SigningAlgorithm(temp.SignAlgo)
	hashAlgo := crypto.HashingAlgorithm(temp.HashAlgo)

	publicKey, err := crypto.DecodePublicKey(signAlgo, temp.EncodedPublicKey)
	if err != nil {
		return a, err
	}

	return flow.AccountKey{
		PublicKey: publicKey,
		SignAlgo:  signAlgo,
		HashAlgo:  hashAlgo,
		Weight:    int(temp.Weight),
	}, nil
}

type accountPublicKeyWrapper struct {
	EncodedPublicKey []byte
	SignAlgo         uint
	HashAlgo         uint
	Weight           uint
}

type accountPrivateKeyWrapper struct {
	EncodedPrivateKey []byte
	SignAlgo          uint
	HashAlgo          uint
}

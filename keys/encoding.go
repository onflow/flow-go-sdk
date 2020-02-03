package keys

import (
	"encoding/hex"
	"errors"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go/crypto"
)

// EncodePrivateKey encodes a private key as bytes.
func EncodePrivateKey(a flow.AccountPrivateKey) ([]byte, error) {
	privateKey, err := a.PrivateKey.Encode()
	if err != nil {
		return nil, err
	}

	w := accountPrivateKeyWrapper{
		PrivateKey: privateKey,
		SignAlgo:   uint(a.SignAlgo),
		HashAlgo:   uint(a.HashAlgo),
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

	privateKey, err := crypto.DecodePrivateKey(signAlgo, w.PrivateKey)
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
func EncodePublicKey(a flow.AccountPublicKey) ([]byte, error) {
	publicKey, err := a.PublicKey.Encode()
	if err != nil {
		return nil, err
	}

	w := accountPublicKeyWrapper{
		PublicKey: publicKey,
		SignAlgo:  uint(a.SignAlgo),
		HashAlgo:  uint(a.HashAlgo),
		Weight:    uint(a.Weight),
	}

	return flow.DefaultEncoder.Encode(&w)
}

// DecodePublicKey decodes a public key.
func DecodePublicKey(b []byte) (a flow.AccountPublicKey, err error) {
	var w accountPublicKeyWrapper

	err = flow.DefaultEncoder.Decode(b, &w)
	if err != nil {
		return a, err
	}

	signAlgo := crypto.SigningAlgorithm(w.SignAlgo)
	hashAlgo := crypto.HashingAlgorithm(w.HashAlgo)

	publicKey, err := crypto.DecodePublicKey(signAlgo, w.PublicKey)
	if err != nil {
		return a, err
	}

	return flow.AccountPublicKey{
		PublicKey: publicKey,
		SignAlgo:  signAlgo,
		HashAlgo:  hashAlgo,
		Weight:    int(w.Weight),
	}, nil
}

type accountPublicKeyWrapper struct {
	PublicKey []byte
	SignAlgo  uint
	HashAlgo  uint
	Weight    uint
}

type accountPrivateKeyWrapper struct {
	PrivateKey []byte
	SignAlgo   uint
	HashAlgo   uint
}

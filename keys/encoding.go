package keys

import (
	"encoding/hex"
	"errors"

	"github.com/dapperlabs/flow-go/crypto"

	"github.com/dapperlabs/flow-go-sdk"
)

// DecodePrivateKeyHex decodes a private key from a hexadecimal string.
func DecodePrivateKeyHex(signAlgo crypto.SigningAlgorithm, h string) (crypto.PrivateKey, error) {
	b, err := hex.DecodeString(h)
	if err != nil {
		return nil, errors.New("failed to decode hex")
	}

	return crypto.DecodePrivateKey(signAlgo, b)
}

// MustDecodePrivateKeyHex is the same as DecodePrivateKeyHex but panics if the
// input string does not represent a valid private key.
func MustDecodePrivateKeyHex(signAlgo crypto.SigningAlgorithm, h string) crypto.PrivateKey {
	k, err := DecodePrivateKeyHex(signAlgo, h)
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

	temp := struct {
		PublicKey []byte
		SignAlgo  uint
		HashAlgo  uint
		Weight    uint
	}{
		PublicKey: publicKey,
		SignAlgo:  uint(a.SignAlgo),
		HashAlgo:  uint(a.HashAlgo),
		Weight:    uint(a.Weight),
	}

	return flow.DefaultEncoder.Encode(&temp)
}

// DecodePublicKey decodes a public key.
func DecodePublicKey(b []byte) (a flow.AccountKey, err error) {
	var temp struct {
		PublicKey []byte
		SignAlgo  uint
		HashAlgo  uint
		Weight    uint
	}

	err = flow.DefaultEncoder.Decode(b, &temp)
	if err != nil {
		return a, err
	}

	signAlgo := crypto.SigningAlgorithm(temp.SignAlgo)
	hashAlgo := crypto.HashingAlgorithm(temp.HashAlgo)

	publicKey, err := crypto.DecodePublicKey(signAlgo, temp.PublicKey)
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

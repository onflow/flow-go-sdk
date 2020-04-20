package crypto

import (
	"encoding/hex"
	"errors"

	"github.com/ethereum/go-ethereum/rlp"
)

// EncodeWrappedPrivateKey encodes a private key as bytes.
func EncodeWrappedPrivateKey(privateKey PrivateKey, sigAlgo SignatureAlgorithm, hashAlgo HashAlgorithm) ([]byte, error) {
	encPrivateKey := privateKey.Encode()

	temp := accountPrivateKeyWrapper{
		EncodedPrivateKey: encPrivateKey,
		SigAlgo:           uint(sigAlgo),
		HashAlgo:          uint(hashAlgo),
	}

	return rlp.EncodeToBytes(&temp)
}

// DecodeWrappedPrivateKey decodes a private key.
func DecodeWrappedPrivateKey(b []byte) (pk PrivateKey, sigAlgo SignatureAlgorithm, hashAlgo HashAlgorithm, err error) {
	var temp accountPrivateKeyWrapper

	err = rlp.DecodeBytes(b, &temp)
	if err != nil {
		return pk, sigAlgo, hashAlgo, err
	}

	sigAlgo = SignatureAlgorithm(temp.SigAlgo)
	hashAlgo = HashAlgorithm(temp.HashAlgo)

	privateKey, err := DecodePrivateKey(sigAlgo, temp.EncodedPrivateKey)
	if err != nil {
		return pk, sigAlgo, hashAlgo, err
	}

	return privateKey, sigAlgo, hashAlgo, nil
}

// DecodeWrappedPrivateKeyHex decodes a private key from a hexadecimal string.
func DecodeWrappedPrivateKeyHex(h string) (pk PrivateKey, sigAlgo SignatureAlgorithm, hashAlgo HashAlgorithm, err error) {
	b, err := hex.DecodeString(h)
	if err != nil {
		return pk, sigAlgo, hashAlgo, errors.New("failed to decode hex")
	}

	pk, sigAlgo, hashAlgo, err = DecodeWrappedPrivateKey(b)
	if err != nil {
		return pk, sigAlgo, hashAlgo, errors.New("failed to decode private key bytes")
	}

	return pk, sigAlgo, hashAlgo, nil
}

// MustDecodeWrappedPrivateKeyHex is the same as DecodePrivateKeyHex but panics if the
// input string does not represent a valid private key.
func MustDecodeWrappedPrivateKeyHex(h string) (pk PrivateKey, sigAlgo SignatureAlgorithm, hashAlgo HashAlgorithm) {
	pk, sigAlgo, hashAlgo, err := DecodeWrappedPrivateKeyHex(h)
	if err != nil {
		panic(err)
	}

	return pk, sigAlgo, hashAlgo
}

type accountPrivateKeyWrapper struct {
	EncodedPrivateKey []byte
	SigAlgo           uint
	HashAlgo          uint
}

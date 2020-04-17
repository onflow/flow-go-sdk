// Package keys provides utilities for generating, encoding, and decoding Flow account keys.
package keys

import (
	"github.com/pkg/errors"

	"github.com/dapperlabs/flow-go/crypto"

	"github.com/onflow/flow-go-sdk"
)

// KeyType is a key format supported by Flow.
type KeyType int

// revive:disable:var-naming

const (
	UnknownKeyType KeyType = iota
	ECDSA_P256_SHA2_256
	ECDSA_P256_SHA3_256
	ECDSA_SECp256k1_SHA2_256
	ECDSA_SECp256k1_SHA3_256
)

// revive:enable

// SigningAlgorithm returns the signing algorithm for this key type.
func (k KeyType) SigningAlgorithm() crypto.SigningAlgorithm {
	switch k {
	case ECDSA_P256_SHA2_256, ECDSA_P256_SHA3_256:
		return crypto.ECDSA_P256
	case ECDSA_SECp256k1_SHA2_256, ECDSA_SECp256k1_SHA3_256:
		return crypto.ECDSA_SECp256k1
	default:
		return crypto.UnknownSigningAlgorithm
	}
}

// HashingAlgorithm returns the hashing algorithm for this key type.
func (k KeyType) HashingAlgorithm() crypto.HashingAlgorithm {
	switch k {
	case ECDSA_P256_SHA2_256, ECDSA_SECp256k1_SHA2_256:
		return crypto.SHA2_256
	case ECDSA_P256_SHA3_256, ECDSA_SECp256k1_SHA3_256:
		return crypto.SHA3_256
	default:
		return crypto.UnknownHashingAlgorithm
	}
}

// PublicKeyWeightThreshold is the total weight required for a set of keys to unlock an account.
const PublicKeyWeightThreshold int = 1000

// GeneratePrivateKey generates a private key of the specified key type.
func GeneratePrivateKey(keyType KeyType, seed []byte) (flow.AccountPrivateKey, error) {
	hasher, err := crypto.NewHasher(crypto.SHA3_384)
	if err != nil {
		return flow.AccountPrivateKey{}, err
	}

	keyGenSeed := hasher.ComputeHash(seed)

	privateKey, err := crypto.GeneratePrivateKey(keyType.SigningAlgorithm(), keyGenSeed)
	if err != nil {
		return flow.AccountPrivateKey{}, err
	}

	return flow.AccountPrivateKey{
		PrivateKey: privateKey,
		SignAlgo:   keyType.SigningAlgorithm(),
		HashAlgo:   keyType.HashingAlgorithm(),
	}, nil
}

// ValidateEncodedPublicKey returns an error if the bytes do not represent a valid public key.
func ValidateEncodedPublicKey(b []byte) error {
	publicKey, err := DecodePublicKey(b)
	if err != nil {
		return errors.Wrap(err, "invalid public key encoding")
	}

	return ValidatePublicKey(publicKey)
}

// ValidatePublicKey returns an error if the public key is invalid.
func ValidatePublicKey(publicKey flow.AccountKey) error {
	if !CompatibleAlgorithms(publicKey.SignAlgo, publicKey.HashAlgo) {
		return errors.Errorf(
			"signing algorithm (%s) is incompatible with hashing algorithm (%s)",
			publicKey.SignAlgo,
			publicKey.HashAlgo,
		)
	}

	return nil
}

// CompatibleAlgorithms returns true if the given signing and hashing algorithms are compatible.
func CompatibleAlgorithms(signAlgo crypto.SigningAlgorithm, hashAlgo crypto.HashingAlgorithm) bool {
	t := map[crypto.SigningAlgorithm]map[crypto.HashingAlgorithm]bool{
		crypto.ECDSA_P256: {
			crypto.SHA2_256: true,
			crypto.SHA3_256: true,
		},
		crypto.ECDSA_SECp256k1: {
			crypto.SHA2_256: true,
			crypto.SHA3_256: true,
		},
	}

	return t[signAlgo][hashAlgo]
}

package crypto

import (
	"encoding/hex"
	"fmt"

	"github.com/onflow/flow-go-sdk/crypto/internal/crypto"
	"github.com/onflow/flow-go-sdk/crypto/internal/crypto/hash"
)

// SignatureAlgorithm is an identifier for a signature algorithm (and parameters if applicable).
type SignatureAlgorithm int

const (
	UnknownSignatureAlgorithm SignatureAlgorithm = iota
	// BLS_BLS12381 is BLS on BLS 12-381 curve
	BLS_BLS12381
	// ECDSA_P256 is ECDSA on NIST P-256 curve
	ECDSA_P256
	// ECDSA_secp256k1 is ECDSA on secp256k1 curve
	ECDSA_secp256k1
)

// String returns the string representation of this signature algorithm.
func (f SignatureAlgorithm) String() string {
	return [...]string{"UNKNOWN", "BLS_BLS12381", "ECDSA_P256", "ECDSA_secp256k1"}[f]
}

func StringToSignatureAlgorithm(s string) SignatureAlgorithm {
	switch s {
	case BLS_BLS12381.String():
		return BLS_BLS12381
	case ECDSA_P256.String():
		return ECDSA_P256
	case ECDSA_secp256k1.String():
		return ECDSA_secp256k1
	default:
		return UnknownSignatureAlgorithm
	}
}

// HashAlgorithm is an identifier for a hash algorithm.
type HashAlgorithm int

const (
	UnknownHashAlgorithm HashAlgorithm = iota
	SHA2_256
	SHA2_384
	SHA3_256
	SHA3_384
)

// String returns the string representation of this hash algorithm.
func (f HashAlgorithm) String() string {
	return [...]string{"UNKNOWN", "SHA2_256", "SHA2_384", "SHA3_256", "SHA3_384"}[f]
}

func StringToHashAlgorithm(s string) HashAlgorithm {
	switch s {
	case SHA2_256.String():
		return SHA2_256
	case SHA2_384.String():
		return SHA2_384
	case SHA3_256.String():
		return SHA3_256
	case SHA3_384.String():
		return SHA3_384
	default:
		return UnknownHashAlgorithm
	}
}

const (
	MinSeedLengthECDSA_P256      = crypto.KeyGenSeedMinLenECDSAP256
	MinSeedLengthECDSA_secp256k1 = crypto.KeyGenSeedMinLenECDSASecp256k1
)

// KeyType is a key format supported by Flow.
type KeyType int

const (
	UnknownKeyType KeyType = iota
	ECDSA_P256_SHA2_256
	ECDSA_P256_SHA3_256
	ECDSA_secp256k1_SHA2_256
	ECDSA_secp256k1_SHA3_256
)

// SignatureAlgorithm returns the signature algorithm for this key type.
func (k KeyType) SignatureAlgorithm() SignatureAlgorithm {
	switch k {
	case ECDSA_P256_SHA2_256, ECDSA_P256_SHA3_256:
		return ECDSA_P256
	case ECDSA_secp256k1_SHA2_256, ECDSA_secp256k1_SHA3_256:
		return ECDSA_secp256k1
	default:
		return UnknownSignatureAlgorithm
	}
}

// HashAlgorithm returns the hash algorithm for this key type.
func (k KeyType) HashAlgorithm() HashAlgorithm {
	switch k {
	case ECDSA_P256_SHA2_256, ECDSA_secp256k1_SHA2_256:
		return SHA2_256
	case ECDSA_P256_SHA3_256, ECDSA_secp256k1_SHA3_256:
		return SHA3_256
	default:
		return UnknownHashAlgorithm
	}
}

type PrivateKey struct {
	private crypto.PrivateKey
}

func (pk PrivateKey) Sign(message []byte, hasher Hasher) ([]byte, error) {
	return pk.private.Sign(message, hasher)
}

func (pk PrivateKey) Algorithm() SignatureAlgorithm {
	return SignatureAlgorithm(pk.private.Algorithm())
}

func (pk PrivateKey) PublicKey() PublicKey {
	return PublicKey{publicKey: pk.private.PublicKey()}
}

func (pk PrivateKey) Encode() []byte {
	return pk.private.Encode()
}

type PublicKey struct {
	publicKey crypto.PublicKey
}

func (pk PublicKey) Verify(sig, message []byte, hasher Hasher) (bool, error) {
	return pk.publicKey.Verify(sig, message, hasher)
}

func (pk PublicKey) Algorithm() SignatureAlgorithm {
	return SignatureAlgorithm(pk.publicKey.Algorithm())
}

func (pk PublicKey) Encode() []byte {
	return pk.publicKey.Encode()
}

type Signer interface {
	Sign(message []byte) ([]byte, error)
}

type InMemorySigner struct {
	PrivateKey PrivateKey
	Hasher     Hasher
}

func NewInMemorySigner(privateKey PrivateKey, hashAlgo HashAlgorithm) InMemorySigner {
	hasher, _ := NewHasher(hashAlgo)

	return InMemorySigner{
		PrivateKey: privateKey,
		Hasher:     hasher,
	}
}

func (s InMemorySigner) Sign(message []byte) ([]byte, error) {
	return s.PrivateKey.Sign(message, s.Hasher)
}

type NaiveSigner = InMemorySigner

func NewNaiveSigner(privateKey PrivateKey, hashAlgo HashAlgorithm) NaiveSigner {
	return NewInMemorySigner(privateKey, hashAlgo)
}

func GeneratePrivateKey(sigAlgo SignatureAlgorithm, seed []byte) (PrivateKey, error) {
	privKey, err := crypto.GeneratePrivateKey(crypto.SigningAlgorithm(sigAlgo), seed)
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKey{
		private: privKey,
	}, nil
}

func DecodePrivateKey(sigAlgo SignatureAlgorithm, b []byte) (PrivateKey, error) {
	privKey, err := crypto.DecodePrivateKey(crypto.SigningAlgorithm(sigAlgo), b)
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKey{
		private: privKey,
	}, nil
}

func DecodePrivateKeyHex(sigAlgo SignatureAlgorithm, s string) (PrivateKey, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return PrivateKey{}, err
	}

	return DecodePrivateKey(sigAlgo, b)
}

func DecodePublicKey(sigAlgo SignatureAlgorithm, b []byte) (PublicKey, error) {
	pubKey, err := crypto.DecodePublicKey(crypto.SigningAlgorithm(sigAlgo), b)
	if err != nil {
		return PublicKey{}, err
	}

	return PublicKey{
		publicKey: pubKey,
	}, nil
}

func DecodePublicKeyHex(sigAlgo SignatureAlgorithm, s string) (PublicKey, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return PublicKey{}, err
	}

	return DecodePublicKey(sigAlgo, b)
}

// CompatibleAlgorithms returns true if the signature and hash algorithms are compatible.
func CompatibleAlgorithms(sigAlgo SignatureAlgorithm, hashAlgo HashAlgorithm) bool {
	switch sigAlgo {
	case ECDSA_P256:
		fallthrough
	case ECDSA_secp256k1:
		switch hashAlgo {
		case SHA2_256:
			fallthrough
		case SHA3_256:
			return true
		}
	}
	return false
}

// Package crypto ...
package crypto

import (
	"crypto/elliptic"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// revive:disable:var-naming

var ECDSA_P256Instance *ECDSAalgo
var ECDSA_SECp256k1Instance *ECDSAalgo

//  Once variables to make sure each Signer is instantiated only once
var ECDSA_P256Once sync.Once
var ECDSA_SECp256k1Once sync.Once

// revive:enable

// Signer interface
type signer interface {
	// generatePrKey generates a private key
	generatePrivateKey([]byte) (PrivateKey, error)
	// decodePrKey loads a private key from a byte array
	decodePrivateKey([]byte) (PrivateKey, error)
	// decodePubKey loads a public key from a byte array
	decodePublicKey([]byte) (PublicKey, error)
}

// commonSigner holds the common data for all signers
type commonSigner struct {
	algo            SigningAlgorithm
	prKeyLength     int
	pubKeyLength    int
	signatureLength int
}

// newNonRelicSigner initializes a signer that does not depend on the Relic library.
func newNonRelicSigner(algo SigningAlgorithm) (signer, error) {
	if algo == ECDSA_P256 {
		ECDSA_P256Once.Do(func() {
			ECDSA_P256Instance = &(ECDSAalgo{
				curve: elliptic.P256(),
				commonSigner: &commonSigner{
					algo,
					PrKeyLenECDSA_P256,
					PubKeyLenECDSA_P256,
					SignatureLenECDSA_P256,
				},
			})
		})
		return ECDSA_P256Instance, nil
	}

	if algo == ECDSA_SECp256k1 {
		ECDSA_SECp256k1Once.Do(func() {
			ECDSA_SECp256k1Instance = &(ECDSAalgo{
				curve: secp256k1(),
				commonSigner: &commonSigner{
					algo,
					PrKeyLenECDSA_SECp256k1,
					PubKeyLenECDSA_SECp256k1,
					SignatureLenECDSA_SECp256k1,
				},
			})
		})
		return ECDSA_SECp256k1Instance, nil
	}
	return nil, cryptoError{fmt.Sprintf("the signature scheme %s is not supported.", algo)}
}

// GeneratePrivateKey generates a private key of the algorithm using the entropy of the given seed
func GeneratePrivateKey(algo SigningAlgorithm, seed []byte) (PrivateKey, error) {
	signer, err := NewSigner(algo)
	if err != nil {
		return nil, err
	}
	return signer.generatePrivateKey(seed)
}

// DecodePrivateKey decodes an array of bytes into a private key of the given algorithm
func DecodePrivateKey(algo SigningAlgorithm, data []byte) (PrivateKey, error) {
	signer, err := NewSigner(algo)
	if err != nil {
		return nil, err
	}
	return signer.decodePrivateKey(data)
}

// DecodePublicKey decodes an array of bytes into a public key of the given algorithm
func DecodePublicKey(algo SigningAlgorithm, data []byte) (PublicKey, error) {
	signer, err := NewSigner(algo)
	if err != nil {
		return nil, err
	}
	return signer.decodePublicKey(data)
}

// Signature type tools

// Bytes returns a byte array of the signature data
func (s Signature) Bytes() []byte {
	return s[:]
}

// String returns a String representation of the signature data
func (s Signature) String() string {
	const zero = "00"
	var sb strings.Builder
	sb.WriteString("0x")
	for _, i := range s {
		hex := strconv.FormatUint(uint64(i), 16)
		sb.WriteString(zero[:2-len(hex)])
		sb.WriteString(hex)
	}
	return sb.String()
}

// Key Pair

// PrivateKey is an unspecified signature scheme private key
type PrivateKey interface {
	// Algorithm returns the signing algorithm related to the private key.
	Algorithm() SigningAlgorithm
	// KeySize return the key size in bytes.
	KeySize() int
	// Sign generates a signature using the provided hasher.
	Sign([]byte, Hasher) (Signature, error)
	// PublicKey returns the public key.
	PublicKey() PublicKey
	// Encode returns a bytes representation of the private key
	Encode() ([]byte, error)
}

// PublicKey is an unspecified signature scheme public key.
type PublicKey interface {
	// Algorithm returns the signing algorithm related to the public key.
	Algorithm() SigningAlgorithm
	// KeySize return the key size in bytes.
	KeySize() int
	// Verify verifies a signature of an input message using the provided hasher.
	Verify(Signature, []byte, Hasher) (bool, error)
	// Encode returns a bytes representation of the public key.
	Encode() ([]byte, error)
}

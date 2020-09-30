/*
 * Flow Go SDK
 *
 * Copyright 2019-2020 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package crypto

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
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

// StringToSignatureAlgorithm converts a string to a SignatureAlgorithm.
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

// StringToHashAlgorithm converts a string to a HashAlgorithm.
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

// A PrivateKey is a cryptographic private key that can be used for in-memory signing.
type PrivateKey struct {
	privateKey crypto.PrivateKey
}

// Sign signs the given message with this private key and the provided hasher.
//
// This function returns an error if a signature cannot be generated.
func (sk PrivateKey) Sign(message []byte, hasher Hasher) ([]byte, error) {
	return sk.privateKey.Sign(message, hasher)
}

// Algorithm returns the signature algorithm for this private key.
func (sk PrivateKey) Algorithm() SignatureAlgorithm {
	return SignatureAlgorithm(sk.privateKey.Algorithm())
}

// PublicKey returns the public key for this private key.
func (sk PrivateKey) PublicKey() PublicKey {
	return PublicKey{publicKey: sk.privateKey.PublicKey()}
}

// Encode returns the raw byte encoding of this private key.
func (sk PrivateKey) Encode() []byte {
	return sk.privateKey.Encode()
}

// A PublicKey is a cryptographic public key that can be used to verify signatures.
type PublicKey struct {
	publicKey crypto.PublicKey
}

// Verify verifies the given signature against a message with this public key and the provided hasher.
//
// This function returns true if the signature is valid for the message, and false otherwise. An error
// is returned if the signature cannot be verified.
func (pk PublicKey) Verify(sig, message []byte, hasher Hasher) (bool, error) {
	return pk.publicKey.Verify(sig, message, hasher)
}

// Algorithm returns the signature algorithm for this public key.
func (pk PublicKey) Algorithm() SignatureAlgorithm {
	return SignatureAlgorithm(pk.publicKey.Algorithm())
}

// Encode returns the raw byte encoding of this public key.
func (pk PublicKey) Encode() []byte {
	return pk.publicKey.Encode()
}

// A Signer is capable of generating cryptographic signatures.
type Signer interface {
	// Sign signs the given message with this signer.
	Sign(message []byte) ([]byte, error)
}

// An InMemorySigner is a signer that generates signatures using an in-memory private key.
//
// InMemorySigner implements simple signing that does not protect the private key against
// any tampering or side channel attacks.
type InMemorySigner struct {
	PrivateKey PrivateKey
	Hasher     Hasher
}

// NewInMemorySigner initializes and returns a new in-memory signer with the provided private key
// and hasher.
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

// NaiveSigner is an alias for InMemorySigner.
type NaiveSigner = InMemorySigner

// NewNaiveSigner is an alias for NewInMemorySigner.
func NewNaiveSigner(privateKey PrivateKey, hashAlgo HashAlgorithm) NaiveSigner {
	return NewInMemorySigner(privateKey, hashAlgo)
}

// MinSeedLength is the generic minimum seed length required to guarantee sufficient
// entropy when generating keys.
//
// This minimum is used when the seed source is not necessarily a CSPRG and the seed
// should be expanded before being passed to the key generation process.
const MinSeedLength = crypto.MinSeedLen

func keyGenerationKMACTag(sigAlgo SignatureAlgorithm) []byte {
	return []byte(fmt.Sprintf("%s Key Generation", sigAlgo))
}

// GeneratePrivateKey generates a private key with the specified signature algorithm from the given seed.
func GeneratePrivateKey(sigAlgo SignatureAlgorithm, seed []byte) (PrivateKey, error) {
	// check the seed has minimum entropy
	if len(seed) < MinSeedLength {
		return PrivateKey{}, fmt.Errorf(
			"crypto: insufficient seed length %d, must be at least %d bytes for %s",
			len(seed),
			MinSeedLength,
			sigAlgo,
		)
	}

	// expand the seed and uniformize its entropy
	var seedLen int
	switch sigAlgo {
	case ECDSA_P256:
		seedLen = crypto.KeyGenSeedMinLenECDSAP256
	case ECDSA_secp256k1:
		seedLen = crypto.KeyGenSeedMinLenECDSASecp256k1
	default:
		return PrivateKey{}, fmt.Errorf(
			"crypto: Go SDK does not support key generation for %s algorithm",
			sigAlgo,
		)
	}

	generationTag := keyGenerationKMACTag(sigAlgo)
	customizer := []byte("")
	hasher, err := hash.NewKMAC_128(generationTag, customizer, seedLen)
	if err != nil {
		return PrivateKey{}, err
	}

	hashedSeed := hasher.ComputeHash(seed)

	// generate the key
	privKey, err := crypto.GeneratePrivateKey(crypto.SigningAlgorithm(sigAlgo), hashedSeed)
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKey{
		privateKey: privKey,
	}, nil
}

// DecodePrivateKey decodes a raw byte encoded private key with the given signature algorithm.
func DecodePrivateKey(sigAlgo SignatureAlgorithm, b []byte) (PrivateKey, error) {
	privKey, err := crypto.DecodePrivateKey(crypto.SigningAlgorithm(sigAlgo), b)
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKey{
		privateKey: privKey,
	}, nil
}

// DecodePrivateKeyHex decodes a raw hex encoded private key with the given signature algorithm.
func DecodePrivateKeyHex(sigAlgo SignatureAlgorithm, s string) (PrivateKey, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return PrivateKey{}, err
	}

	return DecodePrivateKey(sigAlgo, b)
}

// DecodePublicKey decodes a raw byte encoded public key with the given signature algorithm.
func DecodePublicKey(sigAlgo SignatureAlgorithm, b []byte) (PublicKey, error) {
	pubKey, err := crypto.DecodePublicKey(crypto.SigningAlgorithm(sigAlgo), b)
	if err != nil {
		return PublicKey{}, err
	}

	return PublicKey{
		publicKey: pubKey,
	}, nil
}

// DecodePublicKeyHex decodes a raw hex encoded public key with the given signature algorithm.
func DecodePublicKeyHex(sigAlgo SignatureAlgorithm, s string) (PublicKey, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return PublicKey{}, err
	}

	return DecodePublicKey(sigAlgo, b)
}

// DecodePublicKeyHex decodes a PEM public key with the given signature algorithm.
func DecodePublicKeyPEM(sigAlgo SignatureAlgorithm, s string) (PublicKey, error) {
	block, _ := pem.Decode([]byte(s))

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return PublicKey{}, fmt.Errorf("crypto: failed to parse PEM string: %w", err)
	}

	goPublicKey := publicKey.(*ecdsa.PublicKey)

	rawPublicKey := append(
		goPublicKey.X.Bytes(),
		goPublicKey.Y.Bytes()...,
	)

	return DecodePublicKey(sigAlgo, rawPublicKey)
}

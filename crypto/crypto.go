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
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/onflow/flow-go/crypto"
	"github.com/onflow/flow-go/crypto/hash"
)

// SignatureAlgorithm is an identifier for a signature algorithm (and parameters if applicable).
type SignatureAlgorithm = crypto.SigningAlgorithm

const (
	UnknownSignatureAlgorithm SignatureAlgorithm = crypto.UnknownSigningAlgorithm
	// ECDSA_P256 is ECDSA on NIST P-256 curve
	ECDSA_P256 = crypto.ECDSAP256
	// ECDSA_secp256k1 is ECDSA on secp256k1 curve
	ECDSA_secp256k1 = crypto.ECDSASecp256k1
)

// StringToSignatureAlgorithm converts a string to a SignatureAlgorithm.
func StringToSignatureAlgorithm(s string) SignatureAlgorithm {
	switch s {
	case ECDSA_P256.String():
		return ECDSA_P256
	case ECDSA_secp256k1.String():
		return ECDSA_secp256k1
	default:
		return UnknownSignatureAlgorithm
	}
}

// HashAlgorithm is an identifier for a hash algorithm.
type HashAlgorithm = hash.HashingAlgorithm

const (
	UnknownHashAlgorithm HashAlgorithm = hash.UnknownHashingAlgorithm
	SHA2_256                           = hash.SHA2_256
	SHA2_384                           = hash.SHA2_384
	SHA3_256                           = hash.SHA3_256
	SHA3_384                           = hash.SHA3_384
)

// StringToHashAlgorithm converts a string to a HashAlgorithm.
func StringToHashAlgorithm(s string) HashAlgorithm {
	switch s {
	case SHA2_256.String():
		return SHA2_256
	case SHA3_256.String():
		return SHA3_256

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
type PrivateKey = crypto.PrivateKey

// A PublicKey is a cryptographic public key that can be used to verify signatures.
type PublicKey = crypto.PublicKey

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

// MinSeedLength is the generic minimum seed length required to make sure there is
// enough entropy to generate keys targeting 128 bits of security.
// (this is not a guarantee though).
//
// This minimum is used when the seed source is not necessarily a CSPRG and the seed
// should be expanded before being passed to the key generation process.
const MinSeedLength = 32

func keyGenerationKMACTag(sigAlgo SignatureAlgorithm) []byte {
	return []byte(fmt.Sprintf("%s Key Generation", sigAlgo))
}

// GeneratePrivateKey generates a private key with the specified signature algorithm from the given seed.
func GeneratePrivateKey(sigAlgo SignatureAlgorithm, seed []byte) (PrivateKey, error) {
	// check the seed has minimum entropy
	if len(seed) < MinSeedLength {
		return nil, fmt.Errorf(
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
		return nil, fmt.Errorf(
			"crypto: Go SDK does not support key generation for %s algorithm",
			sigAlgo,
		)
	}

	generationTag := keyGenerationKMACTag(sigAlgo)
	customizer := []byte("")
	hasher, err := hash.NewKMAC_128(generationTag, customizer, seedLen)
	if err != nil {
		return nil, err
	}

	hashedSeed := hasher.ComputeHash(seed)

	// generate the key
	privKey, err := crypto.GeneratePrivateKey(sigAlgo, hashedSeed)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

// DecodePrivateKey decodes a raw byte encoded private key with the given signature algorithm.
var DecodePrivateKey = crypto.DecodePrivateKey

// DecodePrivateKeyHex decodes a raw hex encoded private key with the given signature algorithm.
func DecodePrivateKeyHex(sigAlgo SignatureAlgorithm, s string) (PrivateKey, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return DecodePrivateKey(sigAlgo, b)
}

// DecodePublicKey decodes a raw byte encoded public key with the given signature algorithm.
var DecodePublicKey = crypto.DecodePublicKey

// DecodePublicKeyHex decodes a raw hex encoded public key with the given signature algorithm.
func DecodePublicKeyHex(sigAlgo SignatureAlgorithm, s string) (PublicKey, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return DecodePublicKey(sigAlgo, b)
}

// DecodePublicKeyHex decodes a PEM ECDSA public key with the given curve, encoded in `sigAlgo`.
//
// The function only supports ECDSA with P256 and secp256k1 curves.
func DecodePublicKeyPEM(sigAlgo SignatureAlgorithm, s string) (PublicKey, error) {

	if sigAlgo != ECDSA_P256 && sigAlgo != ECDSA_secp256k1 {
		return nil, fmt.Errorf("crypto: only ECDSA algorithms are supported")
	}

	block, rest := pem.Decode([]byte(s))
	if len(rest) > 0 {
		return nil, fmt.Errorf("crypto: failed to parse PEM string, not all bytes in PEM key were decoded: %x", rest)
	}

	// parse the public key data and extract the raw public key
	pkBytes, err := x509ParseECDSAPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("crypto: failed to parse PEM string: %w", err)
	}

	// decode the point and check the resulting key is a valid point on the curve
	return DecodePublicKey(sigAlgo, pkBytes)
}

type publicKeyInfo struct {
	Raw       asn1.RawContent
	Algorithm pkix.AlgorithmIdentifier
	PublicKey asn1.BitString
}

var (
	// object IDs of ECDSA and the 2 supported curves (https://www.secg.org/sec2-v2.pdf)
	oidPublicKeyECDSA      = asn1.ObjectIdentifier{1, 2, 840, 10045, 2, 1}
	oidNamedCurveP256      = asn1.ObjectIdentifier{1, 2, 840, 10045, 3, 1, 7}
	oidNamedCurveSECP256K1 = asn1.ObjectIdentifier{1, 3, 132, 0, 10}
)

// x509ParseECDSAPublicKey parses an ECDSA public key in PKIX, ASN.1 DER form.
//
// The function only supports curves P256 and secp256k1. It doesn't check the
// encoding represents a valid point on the curve.
func x509ParseECDSAPublicKey(derBytes []byte) ([]byte, error) {

	var pki publicKeyInfo
	if rest, err := asn1.Unmarshal(derBytes, &pki); err != nil {
		return nil, err
	} else if len(rest) != 0 {
		return nil, errors.New("x509: trailing data after ASN.1 of public-key")
	}

	// Only ECDSA is supported
	if !pki.Algorithm.Algorithm.Equal(oidPublicKeyECDSA) {
		return nil, errors.New("x509: unknown public key algorithm")
	}

	asn1Data := pki.PublicKey.RightAlign()
	paramsData := pki.Algorithm.Parameters.FullBytes
	namedCurveOID := new(asn1.ObjectIdentifier)
	rest, err := asn1.Unmarshal(paramsData, namedCurveOID)
	if err != nil {
		return nil, errors.New("x509: failed to parse ECDSA parameters as named curve")
	}
	if len(rest) != 0 {
		return nil, errors.New("x509: trailing data after ECDSA parameters")
	}

	// Check the curve is supported
	if !(namedCurveOID.Equal(oidNamedCurveP256) || namedCurveOID.Equal(oidNamedCurveSECP256K1)) {
		return nil, errors.New("x509: unsupported elliptic curve")
	}

	// the candidate field length - this function doesn't check the length is valid
	if asn1Data[0] != 4 { // uncompressed form
		return nil, errors.New("x509: only uncompressed keys are supported")
	}
	return asn1Data[1:], nil
}

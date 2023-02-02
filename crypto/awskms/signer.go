/*
 * Flow Go SDK
 *
 * Copyright 2019 Dapper Labs, Inc.
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

package awskms

import (
	"context"
	"encoding/asn1"
	"fmt"
	"math/big"

	kms "github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/onflow/flow-go-sdk/crypto"
)

var _ crypto.Signer = (*Signer)(nil)

// Signer is a AWS KMS implementation of crypto.Signer.
type Signer struct {
	ctx    context.Context
	client *kms.Client
	key    Key
	// ECDSA is the only algorithm supported by this package. The signature algorithm
	// therefore represents the elliptic curve used. The curve is needed to parse the kms signature.
	curve crypto.SignatureAlgorithm
	// public key for easier access
	publicKey crypto.PublicKey
	// Hash algorithm associated to the KMS signing key
	hashAlgo crypto.HashAlgorithm
}

// SignerForKey returns a new AWS KMS signer for an asymmetric signing key version.
//
// Only ECDSA keys on P-256 and secp256k1 curves and SHA2-256 are supported.
func (c *Client) SignerForKey(
	ctx context.Context,
	key Key,
) (*Signer, error) {
	pk, hashAlgo, err := c.GetPublicKey(ctx, key)
	if err != nil {
		return nil, err
	}

	return &Signer{
		ctx:       ctx,
		client:    c.client,
		key:       key,
		curve:     pk.Algorithm(),
		publicKey: pk,
		hashAlgo:  hashAlgo,
	}, nil
}

// Sign signs the given message using the KMS signing key for this signer.
//
// Reference: https://github.com/aws/aws-sdk-go-v2/blob/main/service/kms/api_op_Sign.go
func (s *Signer) Sign(message []byte) ([]byte, error) {

	keyArn := s.key.ARN()
	// AWS KMS supports signing messages without pre-hashing
	// up to 4096 bytes. Beyond that limit, messages must be prehashed outside KMS.
	kmsPreHashLimit := 4096
	var request *kms.SignInput

	if len(message) <= kmsPreHashLimit {
		request = &kms.SignInput{
			KeyId:            &keyArn,
			Message:          message,
			SigningAlgorithm: types.SigningAlgorithmSpecEcdsaSha256,
		}
	} else {
		// this is guaranteed to only return supported hash algos by KMS
		hasher, err := crypto.NewHasher(s.hashAlgo)
		if err != nil {
			return nil, fmt.Errorf("awskms: failed to sign: %w", err)
		}
		// pre-hash outside KMS
		hash := hasher.ComputeHash(message)
		// indicate the MessageType is digest
		request = &kms.SignInput{
			KeyId:            &keyArn,
			Message:          hash,
			SigningAlgorithm: types.SigningAlgorithmSpecEcdsaSha256,
			MessageType:      types.MessageTypeDigest,
		}
	}
	result, err := s.client.Sign(s.ctx, request)
	if err != nil {
		return nil, fmt.Errorf("awskms: failed to sign: %w", err)
	}
	sig, err := parseSignature(result.Signature, s.curve)
	if err != nil {
		return nil, fmt.Errorf("awskms: failed to parse signature: %w", err)
	}
	return sig, nil
}

// parseSignature parses an asn1 stucture (R,S) into a slice of bytes as required by the `Siger.Sign` method.
func parseSignature(kmsSignature []byte, curve crypto.SignatureAlgorithm) ([]byte, error) {
	var parsedSig struct{ R, S *big.Int }
	if _, err := asn1.Unmarshal(kmsSignature, &parsedSig); err != nil {
		return nil, fmt.Errorf("asn1.Unmarshal: %w", err)
	}

	curveOrderLen := curveOrder(curve)
	signature := make([]byte, 2*curveOrderLen)

	// left pad R and S with zeroes
	rBytes := parsedSig.R.Bytes()
	sBytes := parsedSig.S.Bytes()
	copy(signature[curveOrderLen-len(rBytes):], rBytes)
	copy(signature[len(signature)-len(sBytes):], sBytes)

	return signature, nil
}

// returns the curve order size in bytes (used to padd R and S of the ECDSA signature)
// Only P-256 and secp256k1 are supported. The calling function should make sure
// the function is only called with one of the 2 curves.
func curveOrder(curve crypto.SignatureAlgorithm) int {
	switch curve {
	case crypto.ECDSA_P256:
		return 32
	case crypto.ECDSA_secp256k1:
		return 32
	default:
		panic("the curve is not supported")
	}
}

func (s *Signer) PublicKey() crypto.PublicKey {
	return s.publicKey
}

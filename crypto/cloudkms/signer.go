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

package cloudkms

import (
	"context"
	"encoding/asn1"
	"fmt"
	"math/big"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
)

// Signer is a Google Cloud KMS implementation of crypto.Signer.
type Signer struct {
	ctx      context.Context
	client   *kms.KeyManagementClient
	address  flow.Address
	key      Key
	hashAlgo crypto.HashAlgorithm
	hasher   crypto.Hasher
}

// SignerForKey returns a new Google Cloud KMS signer for an asymmetric key version.
func (c *Client) SignerForKey(
	ctx context.Context,
	address flow.Address,
	key Key,
) (*Signer, error) {
	_, hashAlgo, err := c.GetPublicKey(ctx, key)
	if err != nil {
		return nil, err
	}

	hasher, err := crypto.NewHasher(hashAlgo)
	if err != nil {
		return nil, fmt.Errorf("cloudkms: failed to instantiate hasher: %w", err)
	}

	return &Signer{
		ctx:      ctx,
		client:   c.client,
		address:  address,
		key:      key,
		hashAlgo: hashAlgo,
		hasher:   hasher,
	}, nil
}

// Sign signs the given message using the KMS signing key for this signer.
//
// Reference: https://cloud.google.com/kms/docs/create-validate-signatures
func (s *Signer) Sign(message []byte) ([]byte, error) {
	digest := s.hasher.ComputeHash(message)

	digestMessage, err := makeDigest(s.hashAlgo, digest)
	if err != nil {
		return nil, fmt.Errorf("cloudkms: failed to construct digest: %w", err)
	}

	request := &kmspb.AsymmetricSignRequest{
		Name:   s.key.ResourceID(),
		Digest: digestMessage,
	}

	result, err := s.client.AsymmetricSign(s.ctx, request)
	if err != nil {
		return nil, fmt.Errorf("cloudkms: failed to sign: %w", err)
	}

	sig, err := parseSignature(result.Signature)
	if err != nil {
		return nil, fmt.Errorf("cloudkms: failed to parse signature: %w", err)
	}

	return sig, nil
}

func makeDigest(hashAlgo crypto.HashAlgorithm, digest []byte) (*kmspb.Digest, error) {
	switch hashAlgo {
	case crypto.SHA2_256:
		return &kmspb.Digest{Digest: &kmspb.Digest_Sha256{Sha256: digest}}, nil
	case crypto.SHA2_384:
		return &kmspb.Digest{Digest: &kmspb.Digest_Sha384{Sha384: digest}}, nil
	}

	return nil, fmt.Errorf("unsupported hash algorithm %s", hashAlgo)
}

// ecCoupleComponentSize is size of a component in either (r,s) couple for an elliptical curve signature
// or (x,y) identifying a public key. Component size is needed for encoding couples comprised of variable length
// numbers to []byte encoding. They are not always the same length, so occasionally padding is required.
// Here's how one calculates the required length of each component:
// 		ECDSA_CurveBits = 256
// 		ecCoupleComponentSize := ECDSA_CurveBits / 8
// 		if ECDSA_CurveBits % 8 > 0 {
//			ecCoupleComponentSize++
// 		}
const ecCoupleComponentSize = 32

func parseSignature(signature []byte) ([]byte, error) {
	var parsedSig struct{ R, S *big.Int }
	if _, err := asn1.Unmarshal(signature, &parsedSig); err != nil {
		return nil, fmt.Errorf("asn1.Unmarshal: %w", err)
	}

	rBytes := parsedSig.R.Bytes()
	rBytesPadded := rightPad(rBytes, ecCoupleComponentSize)

	sBytes := parsedSig.S.Bytes()
	sBytesPadded := rightPad(sBytes, ecCoupleComponentSize)

	return append(rBytesPadded, sBytesPadded...), nil
}

// rightPad pads a byte slice with empty bytes (0x00) to the given length.
func rightPad(b []byte, length int) []byte {
	padded := make([]byte, length)
	copy(padded[length-len(b):], b)
	return padded
}

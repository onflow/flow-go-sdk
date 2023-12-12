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

package cloudkms

import (
	"context"
	"fmt"
	"hash/crc32"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/crypto/internal"
)

var _ crypto.Signer = (*Signer)(nil)

// Signer is a Google Cloud KMS implementation of crypto.Signer.
type Signer struct {
	ctx    context.Context
	client *kms.KeyManagementClient
	key    Key
	// ECDSA is the only algorithm supported by this package. The signature algorithm
	// therefore represents the elliptic curve used. The curve is needed to parse the kms signature.
	curve crypto.SignatureAlgorithm
	// public key for easier access
	publicKey crypto.PublicKey
	// Hash algorithm associated to the KMS signing key
	hashAlgo crypto.HashAlgorithm
}

// SignerForKey returns a new Google Cloud KMS signer for an asymmetric signing key version.
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
// Reference: https://cloud.google.com/kms/docs/create-validate-signatures
func (s *Signer) Sign(message []byte) ([]byte, error) {

	// Google KMS supports signing messages without pre-hashing
	// up to 65536 bytes. Beyond that limit, messages must be
	// prehashed outside KMS.
	kmsPreHashLimit := 65536

	var request *kmspb.AsymmetricSignRequest
	if len(message) <= kmsPreHashLimit {
		// hash within KMS
		request = &kmspb.AsymmetricSignRequest{
			Name:       s.key.ResourceID(),
			Data:       message,
			DataCrc32C: checksum(message),
		}
	} else {
		// this is guaranteed to only return supported hash algos by KMS
		// since `s.hashAlgo` is guaranteed to be supported during signer creation
		hasher, err := crypto.NewHasher(s.hashAlgo)
		if err != nil {
			return nil, fmt.Errorf("cloudkms: failed to sign: %w", err)
		}
		// pre-hash outside KMS
		hash := hasher.ComputeHash(message)
		request = &kmspb.AsymmetricSignRequest{
			Name:         s.key.ResourceID(),
			Digest:       getDigest(s.hashAlgo, hash),
			DigestCrc32C: checksum(hash),
		}
	}
	result, err := s.client.AsymmetricSign(s.ctx, request)
	if err != nil {
		return nil, fmt.Errorf("cloudkms: failed to sign: %w", err)
	}
	sig, err := internal.ParseECDSASignature(result.Signature, s.curve)
	if err != nil {
		return nil, fmt.Errorf("cloudkms: failed to parse signature: %w", err)
	}
	return sig, nil
}

func checksum(data []byte) *wrapperspb.Int64Value {
	// compute CRC32
	table := crc32.MakeTable(crc32.Castagnoli)
	checksum := crc32.Checksum(data, table)
	val := wrapperspb.Int64(int64(checksum))
	return val
}

func (s *Signer) PublicKey() crypto.PublicKey {
	return s.publicKey
}

// returns the Digest structure for the hashing algoroithm and hash value, required by the
// signing prehash request
// This function only covers algorithms supported by KMS. It should be extended
// whenever a new hashing algorithm needs to be supported (for instance SHA3-256)
func getDigest(algo crypto.HashAlgorithm, hash []byte) *kmspb.Digest {
	if algo == crypto.SHA2_256 {
		return &kmspb.Digest{
			Digest: &kmspb.Digest_Sha256{
				Sha256: hash,
			},
		}
	}
	return nil
}

/*
 * Flow Go SDK
 *
 * Copyright Flow Foundation
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
// Only ECDSA keys on P-256 and secp256k1 curves are supported. The signer hashes
// messages with the KMS key's native hash algorithm (SHA2-256).
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

// SignerForKeyWithHashAlgo returns a Google Cloud KMS signer that hashes messages
// with the specified hash algorithm instead of the KMS key's native one.
//
// Flow accounts can register a P-256 public key with either SHA2_256 or SHA3_256
// as the on-chain hash algorithm. GCP KMS only supports SHA2_256 natively for its
// EC_SIGN_P256_SHA256 key type, so SHA3_256 must be computed client-side and the
// 32-byte digest submitted via the Digest field — KMS does not re-hash a provided
// digest and only validates its length against the key's native digest size.
func (c *Client) SignerForKeyWithHashAlgo(
	ctx context.Context,
	key Key,
	hashAlgo crypto.HashAlgorithm,
) (*Signer, error) {
	pk, _, err := c.GetPublicKey(ctx, key)
	if err != nil {
		return nil, err
	}

	if _, err := crypto.NewHasher(hashAlgo); err != nil {
		return nil, fmt.Errorf("cloudkms: unsupported hash algorithm %s: %w", hashAlgo, err)
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
// The message is always pre-hashed client-side with the signer's configured hash
// algorithm and submitted via the Digest field, which lets us produce signatures
// over a SHA3_256 digest even for KMS key types that would natively hash with
// SHA2_256.
//
// Reference: https://cloud.google.com/kms/docs/create-validate-signatures
func (s *Signer) Sign(message []byte) ([]byte, error) {
	hasher, err := crypto.NewHasher(s.hashAlgo)
	if err != nil {
		return nil, fmt.Errorf("cloudkms: failed to sign: %w", err)
	}
	hash := hasher.ComputeHash(message)

	digest := getDigest(s.hashAlgo, hash)
	if digest == nil {
		return nil, fmt.Errorf("cloudkms: unsupported hash algorithm %s", s.hashAlgo)
	}

	request := &kmspb.AsymmetricSignRequest{
		Name:         s.key.ResourceID(),
		Digest:       digest,
		DigestCrc32C: checksum(hash),
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

// returns the Digest structure for the hashing algorithm and hash value, required by the
// signing prehash request.
//
// GCP KMS exposes only sha256/sha384/sha512 oneof fields and validates the digest
// length against the key's native digest size. For 32-byte digests on
// EC_SIGN_P256_SHA256 / EC_SIGN_SECP256K1_SHA256 keys, the SHA3_256 hash fits the
// sha256 oneof shape; KMS signs the 32 bytes as-is without re-hashing.
func getDigest(algo crypto.HashAlgorithm, hash []byte) *kmspb.Digest {
	switch algo {
	case crypto.SHA2_256, crypto.SHA3_256:
		return &kmspb.Digest{
			Digest: &kmspb.Digest_Sha256{
				Sha256: hash,
			},
		}
	}
	return nil
}

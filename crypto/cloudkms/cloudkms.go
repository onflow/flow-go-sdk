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

type Key struct {
	ProjectID  string
	LocationID string
	KeyRingID  string
	KeyID      string
	KeyVersion string
}

func (k Key) Name() string {
	return fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s/cryptoKeyVersions/%s",
		k.ProjectID,
		k.LocationID,
		k.KeyRingID,
		k.KeyID,
		k.KeyVersion,
	)
}

type Signer struct {
	ctx      context.Context
	client   *kms.KeyManagementClient
	address  flow.Address
	key      Key
	hashAlgo crypto.HashAlgorithm
	hasher   crypto.Hasher
}

func NewSigner(
	ctx context.Context,
	address flow.Address,
	key Key,
	hashAlgo crypto.HashAlgorithm,
) (*Signer, error) {
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("cloudkms: failed to initialize client: %w", err)
	}

	hasher, err := crypto.NewHasher(hashAlgo)
	if err != nil {
		return nil, fmt.Errorf("cloudkms: failed to initialize hasher: %w", err)
	}

	return &Signer{
		ctx:      ctx,
		client:   client,
		address:  address,
		key:      key,
		hashAlgo: hashAlgo,
		hasher:   hasher,
	}, nil
}

func (s *Signer) Sign(message []byte) ([]byte, error) {
	digest := s.hasher.ComputeHash(message)

	digestMsg, err := makeDigest(s.hashAlgo, digest)
	if err != nil {
		return nil, fmt.Errorf("cloudkms: failed to construct digest: %w", err)
	}

	req := &kmspb.AsymmetricSignRequest{
		Name:   s.key.Name(),
		Digest: digestMsg,
	}

	resp, err := s.client.AsymmetricSign(s.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cloudkms: failed to sign: %w", err)
	}

	sig, err := parseSignature(resp.Signature)
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

const (
	// ECCoupleComponentSize is size of a component in either (r,s) couple for an elliptical curve signature
	// or (x,y) identifying a public key. Component size is needed for encoding couples comprised of variable length
	// numbers to []byte encoding. They are not always the same length, so occasionally padding is required.
	// Here's how one calculates the required length of each component:
	// 		ECDSA_CurveBits = 256
	// 		ECCoupleComponentSize := ECDSA_CurveBits / 8
	// 		if ECDSA_CurveBits % 8 > 0 {
	//			ECCoupleComponentSize++
	// 		}
	ECCoupleComponentSize = 32
)

func parseSignature(signature []byte) ([]byte, error) {
	var parsedSig struct{ R, S *big.Int }
	if _, err := asn1.Unmarshal(signature, &parsedSig); err != nil {
		return nil, fmt.Errorf("asn1.Unmarshal: %w", err)
	}

	rBytes := parsedSig.R.Bytes()
	rBytesPadded := rightPad(rBytes, ECCoupleComponentSize)

	sBytes := parsedSig.S.Bytes()
	sBytesPadded := rightPad(sBytes, ECCoupleComponentSize)

	return append(rBytesPadded, sBytesPadded...), nil
}

// rightPad pads a byte slice with empty bytes (0x00) to the given length.
func rightPad(b []byte, length int) []byte {
	padded := make([]byte, length)
	copy(padded[length-len(b):], b)
	return padded
}

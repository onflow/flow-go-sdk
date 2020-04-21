package crypto

import (
	"fmt"

	"github.com/onflow/flow-go-sdk/crypto/internal/crypto/hash"
)

type Hasher = hash.Hasher
type Hash = hash.Hash

// NewHasher initializes and returns a new hasher with the given hash algorithm.
//
// This function returns an error if the hash algorithm is invalid.
func NewHasher(algo HashAlgorithm) (Hasher, error) {
	switch algo {
	case SHA2_256:
		return NewSHA2_256(), nil
	case SHA2_384:
		return NewSHA2_384(), nil
	case SHA3_256:
		return NewSHA3_256(), nil
	case SHA3_384:
		return NewSHA3_384(), nil
	default:
		return nil, fmt.Errorf("invalid hash algorithm %s", algo)
	}
}

// NewSHA2_256 returns a new instance of SHA2-256 hasher.
func NewSHA2_256() Hasher {
	return hash.NewSHA2_256()
}

// NewSHA2_384 returns a new instance of SHA2-384 hasher.
func NewSHA2_384() Hasher {
	return hash.NewSHA2_384()
}

// NewSHA3_256 returns a new instance of SHA3-256 hasher.
func NewSHA3_256() Hasher {
	return hash.NewSHA3_256()
}

// NewSHA3_384 returns a new instance of SHA3-384 hasher.
func NewSHA3_384() Hasher {
	return hash.NewSHA3_384()
}

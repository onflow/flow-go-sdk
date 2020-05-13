package crypto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go-sdk/crypto"
)

func TestGeneratePrivateKey(t *testing.T) {
	validAlgo := crypto.ECDSA_P256
	validSeed := make([]byte, validAlgo.MinSeedLength()*2) // make seed that is twice the minimum length
	for i, _ := range validSeed {
		validSeed[i] = byte(i)
	}

	t.Run("Nil seed", func(t *testing.T) {
		sk, err := crypto.GeneratePrivateKey(validAlgo, nil)
		assert.Error(t, err)
		assert.Equal(t, crypto.PrivateKey{}, sk)
	})

	t.Run("Empty seed", func(t *testing.T) {
		sk, err := crypto.GeneratePrivateKey(validAlgo, []byte{})
		assert.Error(t, err)
		assert.Equal(t, crypto.PrivateKey{}, sk)
	})

	t.Run("Seed length too short", func(t *testing.T) {
		sk, err := crypto.GeneratePrivateKey(validAlgo, []byte("foo"))
		assert.Error(t, err)
		assert.Equal(t, crypto.PrivateKey{}, sk)
	})

	t.Run("Seed length exactly equal", func(t *testing.T) {
		seed := make([]byte, validAlgo.MinSeedLength())
		for i, _ := range seed {
			seed[i] = byte(i)
		}

		sk, err := crypto.GeneratePrivateKey(validAlgo, seed)
		assert.NoError(t, err)
		assert.NotEqual(t, crypto.PrivateKey{}, sk)
		assert.Equal(t, validAlgo, sk.Algorithm())
	})

	t.Run("Invalid signature algorithm", func(t *testing.T) {
		invalidAlgo := crypto.SignatureAlgorithm(-42)

		sk, err := crypto.GeneratePrivateKey(invalidAlgo, validSeed)
		assert.Error(t, err)
		assert.Equal(t, crypto.PrivateKey{}, sk)
	})

	t.Run("Valid signature algorithm", func(t *testing.T) {
		sk, err := crypto.GeneratePrivateKey(validAlgo, validSeed)
		assert.NoError(t, err)
		assert.NotEqual(t, crypto.PrivateKey{}, sk)
		assert.Equal(t, validAlgo, sk.Algorithm())
	})

	t.Run("Deterministic generation", func(t *testing.T) {
		trials := 100

		var skA crypto.PrivateKey
		var err error

		skA, err = crypto.GeneratePrivateKey(validAlgo, validSeed)
		require.NoError(t, err)

		for i := 0; i < trials; i++ {
			skB, err := crypto.GeneratePrivateKey(validAlgo, validSeed)
			require.NoError(t, err)
			assert.Equal(t, skA, skB) // key should be same each time
			skA = skB
		}
	})
}

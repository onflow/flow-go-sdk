package cloudkms_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go-sdk/crypto/cloudkms"
)

func TestKeyFromResourceID(t *testing.T) {
	key := cloudkms.Key{
		ProjectID:  "my-project",
		LocationID: "global",
		KeyRingID:  "flow",
		KeyID:      "my-account",
		KeyVersion: "1",
	}

	resourceID := key.ResourceID()

	assert.Equal(t, resourceID, "projects/my-project/locations/global/keyRings/flow/cryptoKeys/my-account/cryptoKeyVersions/1")

	keyFromResourceID, err := cloudkms.KeyFromResourceID(resourceID)
	require.NoError(t, err)

	assert.Equal(t, key, keyFromResourceID)
}

package cloudkms_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go-sdk/crypto"
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

// gcloudApplicationSignin signs in as an application user using gcloud command line tool
// currently assumes gcloud is already installed on the machine
// will by default pop a browser window to sign in
func gcloudApplicationSignin(kms cloudkms.Key) error {
	googleAppCreds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if len(googleAppCreds) > 0 {
		return nil
	}

	proj := kms.ProjectID
	if len(proj) == 0 {
		return fmt.Errorf(
			"could not get GOOGLE_APPLICATION_CREDENTIALS, no google service account JSON provided but private key type is KMS",
		)
	}

	loginCmd := exec.Command("gcloud", "auth", "application-default", "login", fmt.Sprintf("--project=%s", proj))

	output, err := loginCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to run %q: %s\n", loginCmd.String(), err)
	}

	squareBracketRegex := regexp.MustCompile(`(?s)\[(.*)\]`)
	regexResult := squareBracketRegex.FindAllStringSubmatch(string(output), -1)
	// Should only be one value. Second index since first index contains the square brackets
	googleApplicationCreds := regexResult[0][1]

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", googleApplicationCreds)

	return nil
}

// TestManualKMSSigning tests signing using a KMS key.
// This tests requires access to KMS and cannot be run by CI. Please use this test manually
// when making any change to the KMS signing code.
func TestManualKMSSigning(t *testing.T) {
	key := cloudkms.Key{
		ProjectID:  "dl-flow",
		LocationID: "us-east1",
		KeyRingID:  "testnet_keyring",
		KeyID:      "sdk_test",
		KeyVersion: "1",
	}

	// id is copied from the kms key: https://console.cloud.google.com/security/kms/key/manage/us-east1/testnet_keyring/sdk_test?project=dl-flow
	// this test insures the key exists in the kms keyring
	id := "projects/dl-flow/locations/us-east1/keyRings/testnet_keyring/cryptoKeys/sdk_test/cryptoKeyVersions/1"
	require.Equal(t, key.ResourceID(), id)

	// get google kms permission
	err := gcloudApplicationSignin(key)
	require.NoError(t, err)

	// initialize the client
	ctx := context.Background()
	cl, err := cloudkms.NewClient(ctx)
	require.NoError(t, err)

	// Get the public key
	pk, _, err := cl.GetPublicKey(ctx, key)
	require.NoError(t, err)

	// Sign
	msg := []byte("random_message")
	signer, err := cl.SignerForKey(ctx, key)
	require.NoError(t, err)
	sig, err := signer.Sign(msg)
	require.NoError(t, err)

	// verify
	hasher := crypto.NewSHA2_256()
	valid, err := pk.Verify(sig, msg, hasher)
	require.NoError(t, err)
	assert.True(t, valid)
}

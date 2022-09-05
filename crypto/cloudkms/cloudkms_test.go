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
// This tests requires access to KMS and cannot be run by CI.
// Please use this test manually by commenting t.Skip(),
// when making any change to the KMS signing code.
// This test assumes gcloud CLI is already installed on
// your machine.
func TestManualKMSSigning(t *testing.T) {
	// to comment when testing manually
	t.Skip()

	// KMS_TEST_KEY_RESOURCE_ID is an env var containing the resource ID of a KMS key you
	// have permissions to use.
	id := os.Getenv(`KMS_TEST_KEY_RESOURCE_ID`)
	fmt.Println(id)
	key, err := cloudkms.KeyFromResourceID(id)
	require.NoError(t, err)

	// get google kms permission
	err = gcloudApplicationSignin(key)
	require.NoError(t, err)

	// initialize the client
	ctx := context.Background()
	cl, err := cloudkms.NewClient(ctx)
	require.NoError(t, err)

	// Get the public key
	pk, _, err := cl.GetPublicKey(ctx, key)
	require.NoError(t, err)

	// signer
	signer, err := cl.SignerForKey(ctx, key)
	require.NoError(t, err)

	signAndVerify := func(t *testing.T, msgLen int) {
		// Sign
		msg := make([]byte, msgLen)
		sig, err := signer.Sign(msg)
		require.NoError(t, err)

		// verify
		hasher := crypto.NewSHA2_256()
		valid, err := pk.Verify(sig, msg, hasher)
		require.NoError(t, err)
		assert.True(t, valid)
	}

	kmsPreHashLimit := 65536
	// google KMS supports signing messages without prehashing
	// up to 65536 bytes
	t.Run("short message", func(t *testing.T) {
		signAndVerify(t, kmsPreHashLimit)
	})

	// google KMS does not support signing messages longer than 65536
	// without prehashing
	t.Run("long message", func(t *testing.T) {
		signAndVerify(t, kmsPreHashLimit+1)
	})

}

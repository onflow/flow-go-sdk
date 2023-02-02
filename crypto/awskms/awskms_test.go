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

package awskms_test

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/crypto/awskms"
)

func TestKeyFromARN(t *testing.T) {
	key := awskms.Key{
		Region:  "us-west-2",
		Account: "111122223333",
		KeyID:   "1234abcd-12ab-34cd-56ef-1234567890ab",
	}

	resourceARN := key.ARN()

	assert.Equal(t, resourceARN, "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab")

	keyFromResourceARN, err := awskms.KeyFromResourceARN(resourceARN)
	require.NoError(t, err)

	assert.Equal(t, key, keyFromResourceARN)
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

	// KMS_TEST_KEY_RESOURCE_ARN is an env var containing the resource ARN of a KMS key you
	// have permissions to use.
	os.Setenv("KMS_TEST_KEY_RESOURCE_ARN", "")
	id := os.Getenv(`KMS_TEST_KEY_RESOURCE_ARN`)
	t.Log(id)
	key, err := awskms.KeyFromResourceARN(id)
	require.NoError(t, err)

	// initialize the client
	ctx := context.Background()
	// AWS SDK uses the default credential chain to find the credentials.
	// You need to export env variables, AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY and AWS_SESSION_TOKEN
	os.Setenv("AWS_ACCESS_KEY_ID", "")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "")
	os.Setenv("AWS_SESSION_TOKEN", "")

	defaultCfg, err := config.LoadDefaultConfig(ctx)
	require.NoError(t, err)

	cl := awskms.NewClient(defaultCfg)
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

	kmsPreHashLimit := 4096
	// AWS KMS supports signing messages without prehashing
	// up to 4096 bytes
	t.Run("short message", func(t *testing.T) {
		signAndVerify(t, kmsPreHashLimit)
	})

	// google KMS does not support signing messages longer than 4096
	// without prehashing
	t.Run("long message", func(t *testing.T) {
		signAndVerify(t, kmsPreHashLimit+1)
	})

}

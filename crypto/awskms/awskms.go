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

// Package awskms provides a AWS Key Management Service (KMS)
// implementation of the crypto.Signer interface.
//
// The documentation for Google Cloud KMS can be found here: https://cloud.google.com/kms/docs
package awskms

import (
	"context"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	kms "github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"

	// kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"

	"github.com/onflow/flow-go-sdk/crypto"
)

const (
	resouceArnFormat        = "arn:aws:kms:%s:%s:key/%s"
	resourceIDArgumentCount = 5
)

// Key is a reference to a AWS KMS asymmetric signing key version.
type Key struct {
	Region  string `json:"region"`
	Account string `json:"account"`
	KeyID   string `json:"keyId"`
}

// ARN returns the KMS arn for this KMS key.
// For cross account key access, you need to pass the arn instead of just the keyID.
func (k Key) ARN() string {
	return fmt.Sprintf(
		resouceArnFormat,
		k.Region,
		k.Account,
		k.KeyID,
	)
}

// Example ARN format: "arn:aws:kms:us-west-2:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab"
func KeyFromResourceARN(resourceARN string) (Key, error) {
	key := Key{}
	spiltedARN := strings.Split(resourceARN, ":")
	if len(spiltedARN) != 6 {
		return key, fmt.Errorf("awskms: wrong format for the resourceARN: %s", resourceARN)
	}

	key.Region, key.Account = spiltedARN[3], spiltedARN[4]
	key.KeyID = strings.Split(spiltedARN[5], "/")[1]

	return key, nil
}

// Client is a client for interacting with the AWS KMS API
// using types native to the Flow Go SDK.
type Client struct {
	client *kms.Client
}

// NewClient creates a new AWS KMS client.
func NewClient(cfg aws.Config) *Client {
	client := kms.NewFromConfig(cfg)
	return &Client{
		client: client,
	}
}

// GetPublicKey fetches the public key portion of a KMS asymmetric signing key version.
//
// KMS keys of the type `CryptoKeyVersion_EC_SIGN_P256_SHA256` and `CryptoKeyVersion_EC_SIGN_SECP256K1_SHA256`
// are the only keys supported by the SDK.
//
// Ref: https://cloud.google.com/kms/docs/retrieve-public-key
func (c *Client) GetPublicKey(ctx context.Context, key Key) (crypto.PublicKey, crypto.HashAlgorithm, error) {

	keyArn := key.ARN()
	request := &kms.GetPublicKeyInput{
		KeyId: &keyArn,
	}

	result, err := c.client.GetPublicKey(ctx, request)
	if err != nil {
		return nil,
			crypto.UnknownHashAlgorithm,
			fmt.Errorf("awskms: failed to fetch public key from KMS API: %v", err)
	}

	sigAlgo := ParseSignatureAlgorithm(result.KeySpec)
	if sigAlgo == crypto.UnknownSignatureAlgorithm {
		return nil,
			crypto.UnknownHashAlgorithm,
			fmt.Errorf(
				"awskms: unsupported signature algorithm %s",
				result.KeySpec,
			)
	}

	hashAlgo := ParseHashAlgorithm(result.KeySpec)
	if hashAlgo == crypto.UnknownHashAlgorithm {
		return nil,
			crypto.UnknownHashAlgorithm,
			fmt.Errorf(
				"awskms: unsupported hash algorithm %s",
				result.KeySpec,
			)
	}

	publicKeyBytes := result.PublicKey
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicKeyPEM := pem.EncodeToMemory(block)
	publicKey, err := crypto.DecodePublicKeyPEM(sigAlgo, string(publicKeyPEM))

	if err != nil {
		return nil,
			crypto.UnknownHashAlgorithm,
			fmt.Errorf("awskms: failed to parse PEM public key: %w", err)
	}

	return publicKey, hashAlgo, nil
}

// KMSClient gives access to the KeyManagementClient,
// e.g. for closing the connection to the Google KMS API
func (c *Client) KMSClient() *kms.Client {
	return c.client
}

// ParseSignatureAlgorithm returns the `SignatureAlgorithm` corresponding to the input KMS key type.
func ParseSignatureAlgorithm(keySpec types.KeySpec) crypto.SignatureAlgorithm {
	if keySpec == types.KeySpecEccNistP256 {
		return crypto.ECDSA_P256
	}

	// not sure if this is the same
	if keySpec == types.KeySpecEccSecgP256k1 {
		return crypto.ECDSA_secp256k1
	}

	return crypto.UnknownSignatureAlgorithm
}

// ParseHashAlgorithm returns the `HashAlgorithm` corresponding to the input KMS key type.
func ParseHashAlgorithm(keySpec types.KeySpec) crypto.HashAlgorithm {
	if keySpec == types.KeySpecEccNistP256 || keySpec == types.KeySpecEccSecgP256k1 {
		return crypto.SHA2_256
	}

	// the function can be extended to return SHA3-256 if it becomes supported by KMS.
	return crypto.UnknownHashAlgorithm
}

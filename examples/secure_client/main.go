/*
 * Flow Go SDK
 *
 * Copyright 2019-2020 Dapper Labs, Inc.
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

package main

import (
	"context"
	goecdsa "crypto/ecdsa"
	"crypto/elliptic"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"

	crypto_pb "github.com/libp2p/go-libp2p-core/crypto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/examples"

	lcrypto "github.com/libp2p/go-libp2p-core/crypto"
	lcrypto_pb "github.com/libp2p/go-libp2p-core/crypto/pb"
	"github.com/libp2p/go-libp2p-core/peer"
	libp2ptls "github.com/libp2p/go-libp2p-tls"
)

// ServerAuthError is an error returned when the server authentication fails
type ServerAuthError struct {
	message string
}

// newServerAuthError constructs a new ServerAuthError
func newServerAuthError(msg string, args ...interface{}) *ServerAuthError {
	return &ServerAuthError{message: fmt.Sprintf(msg, args...)}
}

func (e ServerAuthError) Error() string {
	return e.message
}

// IsServerAuthError checks if the input error is of a ServerAuthError type
func IsServerAuthError(err error) bool {
	_, ok := err.(*ServerAuthError)
	return ok
}

func main() {
	SecurePing()
}

func SecurePing() {
	ctx := context.Background()

	secureGRPCServerAddress := "access-002.canary6.nodes.onflow.org:9001"
	accessNodePublicKey := "1a361155b3dbce94e4ec5dcffba30f85adf1d38a0793aaf3f96c2256267a17581f919bf9003bf3a0cbb8fc3a28561b11b77633c6b6fb62358c48476512c120a8"

	flowClient := secureGRPCClient(accessNodePublicKey, secureGRPCServerAddress)

	err := flowClient.Ping(ctx)
	examples.Handle(err)

	fmt.Println("ping successful")
}

// secureGRPCClient creates a secure GRPC client using the given public key
func secureGRPCClient(publicKeyHex string, secureGRPCServeraddress string) *client.Client {

	bytes, err := hex.DecodeString(publicKeyHex)
	examples.Handle(err)

	publicFlowNetworkingKey, err := crypto.DecodePublicKey(crypto.ECDSA_P256, bytes)
	examples.Handle(err)

	fmt.Printf(" using public key %s for the remote node\n", publicFlowNetworkingKey.String())

	tlsConfig, err := DefaultClientTLSConfig(publicFlowNetworkingKey)
	examples.Handle(err)

	client, err := client.New(secureGRPCServeraddress, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	examples.Handle(err)

	return client
}

// DefaultClientTLSConfig returns the default TLS client config with the given public key for a secure GRPC client
// The TLSConfig verifies that the server certifcate is valid and has the correct signature
func DefaultClientTLSConfig(publicKey crypto.PublicKey) (*tls.Config, error) {

	config := &tls.Config{
		MinVersion:         tls.VersionTLS13,
		InsecureSkipVerify: true, // This is not insecure here. We will verify the cert chain ourselves.
		ClientAuth:         tls.RequireAnyClientCert,
	}

	verifyPeerCertFunc, err := verifyPeerCertificateFunc(publicKey)
	if err != nil {
		return nil, err
	}
	config.VerifyPeerCertificate = verifyPeerCertFunc

	return config, nil
}

func verifyPeerCertificateFunc(expectedPublicKey crypto.PublicKey) (func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error, error) {

	// convert the Flow.crypto key to LibP2P key for easy comparision using LibP2P TLS utils
	expectedLibP2PKey, err := PublicKey(expectedPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate a libp2p key from a Flow key: %w", err)
	}
	remotePeerLibP2PID, err := peer.IDFromPublicKey(expectedLibP2PKey)
	if err != nil {
		return nil, fmt.Errorf("failed to derive the libp2p Peer ID from the libp2p public key: %w", err)
	}

	// We're using InsecureSkipVerify, so the verifiedChains parameter will always be empty.
	// We need to parse the certificates ourselves from the raw certs.
	verifyFunc := func(rawCerts [][]byte, _ [][]*x509.Certificate) error {

		chain := make([]*x509.Certificate, len(rawCerts))
		for i := 0; i < len(rawCerts); i++ {
			cert, err := x509.ParseCertificate(rawCerts[i])
			if err != nil {
				return newServerAuthError(err.Error())
			}
			chain[i] = cert
		}

		// libp2ptls.PubKeyFromCertChain verifies the certificate, verifies that the certificate contains the special libp2p
		// extension, extract the remote's public key and finally verifies the signature included in the certificate
		actualLibP2PKey, err := libp2ptls.PubKeyFromCertChain(chain)
		if err != nil {
			return newServerAuthError(err.Error())
		}

		// verify that the public key received is the one that is expected
		if !remotePeerLibP2PID.MatchesPublicKey(actualLibP2PKey) {
			return newServerAuthError("invalid public key received: expected %s", expectedPublicKey.String())
		}
		return nil
	}

	return verifyFunc, nil
}

// PublicKey converts a Flow public key to a LibP2P public key
func PublicKey(fpk crypto.PublicKey) (lcrypto.PubKey, error) {
	keyType, err := keyType(fpk.Algorithm())
	if err != nil {
		return nil, err
	}
	um, ok := lcrypto.PubKeyUnmarshallers[keyType]
	if !ok {
		return nil, lcrypto.ErrBadKeyType
	}

	tempBytes := fpk.Encode()

	// at this point, keytype is either KeyType_ECDSA or KeyType_Secp256k1
	// and can't hold another value
	var bytes []byte
	if keyType == crypto_pb.KeyType_ECDSA {
		var x, y big.Int
		x.SetBytes(tempBytes[:len(tempBytes)/2])
		y.SetBytes(tempBytes[len(tempBytes)/2:])
		goKey := setPubKey(elliptic.P256(), &x, &y)
		bytes, err = x509.MarshalPKIXPublicKey(goKey)
		if err != nil {
			return nil, lcrypto.ErrBadKeyType
		}
	} else if keyType == lcrypto_pb.KeyType_Secp256k1 {
		bytes = make([]byte, crypto.PubKeyLenECDSASecp256k1+1) // libp2p requires an extra byte
		bytes[0] = 4                                           // magic number in libp2p to refer to an uncompressed key
		copy(bytes[1:], tempBytes)
	}

	return um(bytes)
}

// keyType translates Flow signing algorithm constants to the corresponding LibP2P constants
func keyType(sa crypto.SignatureAlgorithm) (lcrypto_pb.KeyType, error) {
	switch sa {
	case crypto.ECDSA_P256:
		return lcrypto_pb.KeyType_ECDSA, nil
	case crypto.ECDSA_secp256k1:
		return lcrypto_pb.KeyType_Secp256k1, nil
	default:
		return -1, lcrypto.ErrBadKeyType
	}
}

// assigns two big.Int inputs to a Go ecdsa public key
func setPubKey(c elliptic.Curve, x *big.Int, y *big.Int) *goecdsa.PublicKey {
	pub := new(goecdsa.PublicKey)
	pub.Curve = c
	pub.X = x
	pub.Y = y
	return pub
}

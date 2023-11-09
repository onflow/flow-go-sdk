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
package internal

import (
	"encoding/asn1"
	"fmt"
	"math/big"

	"github.com/onflow/flow-go-sdk/crypto"
)

// parseSignature parses an asn1 stucture (R,S) into a slice of bytes as required by the `Siger.Sign` method.
func ParseSignature(kmsSignature []byte, curve crypto.SignatureAlgorithm) ([]byte, error) {
	var parsedSig struct{ R, S *big.Int }
	if _, err := asn1.Unmarshal(kmsSignature, &parsedSig); err != nil {
		return nil, fmt.Errorf("asn1.Unmarshal: %w", err)
	}

	curveOrderLen := curveOrder(curve)
	signature := make([]byte, 2*curveOrderLen)

	// left pad R and S with zeroes
	rBytes := parsedSig.R.Bytes()
	sBytes := parsedSig.S.Bytes()
	copy(signature[curveOrderLen-len(rBytes):], rBytes)
	copy(signature[len(signature)-len(sBytes):], sBytes)

	return signature, nil
}

// returns the curve order size in bytes (used to padd R and S of the ECDSA signature)
// Only P-256 and secp256k1 are supported. The calling function should make sure
// the function is only called with one of the 2 curves.
func curveOrder(curve crypto.SignatureAlgorithm) int {
	switch curve {
	case crypto.ECDSA_P256:
		return 32
	case crypto.ECDSA_secp256k1:
		return 32
	default:
		panic("the curve is not supported")
	}
}

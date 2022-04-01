/*
 * Flow Go SDK
 *
 * Copyright 2022 Dapper Labs, Inc.
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

package flow

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/rlp"
)

// AccountProofNonceMinLenBytes is the minimum length of account proof nonces in bytes.
const AccountProofNonceMinLenBytes = 32

var (
	// ErrInvalidNonce is returned when the account proof nonce passed to a function is invalid.
	ErrInvalidNonce = errors.New("invalid nonce")
)

type canonicalAccountProofV2 struct {
	AppID   string
	Address []byte
	Nonce   []byte
}

// NewAccountProofMessage creates a new account proof message for singing. The appID is optional and can be left
// empty. The encoded message returned does not include the user domain tag.
func NewAccountProofMessage(address Address, appID, nonceHex string) ([]byte, error) {
	nonceBytes, err := hex.DecodeString(strings.TrimPrefix(nonceHex, "0x"))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidNonce, err)
	}

	if len(nonceBytes) < AccountProofNonceMinLenBytes {
		return nil, fmt.Errorf("%w: nonce must be at least %d bytes", ErrInvalidNonce, AccountProofNonceMinLenBytes)
	}

	msg, err := rlp.EncodeToBytes(&canonicalAccountProofV2{
		AppID:   appID,
		Address: address.Bytes(),
		Nonce:   nonceBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("error encoding account proof message: %w", err)
	}

	return msg, nil
}

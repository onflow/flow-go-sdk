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
	// ErrInvalidAppID is returned when the account proof app ID passed to a function is invalid.
	ErrInvalidAppID = errors.New("invalid app ID")
)

type canonicalAccountProof struct {
	AppID   string
	Address []byte
	Nonce   []byte
}

// EncodeAccountProofMessage creates a new account proof message for singing. The encoded message returned does not include
// the user domain tag.
func EncodeAccountProofMessage(address Address, appID, nonceHex string) ([]byte, error) {
	if appID == "" {
		return nil, fmt.Errorf("%w: appID can't be empty", ErrInvalidAppID)
	}

	nonceBytes, err := hex.DecodeString(strings.TrimPrefix(nonceHex, "0x"))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidNonce, err)
	}

	if len(nonceBytes) < AccountProofNonceMinLenBytes {
		return nil, fmt.Errorf("%w: nonce must be at least %d bytes", ErrInvalidNonce, AccountProofNonceMinLenBytes)
	}

	msg, err := rlp.EncodeToBytes(&canonicalAccountProof{
		AppID:   appID,
		Address: address.Bytes(),
		Nonce:   nonceBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("error encoding account proof message: %w", err)
	}

	return msg, nil
}

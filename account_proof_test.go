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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccountProofMessage(t *testing.T) {
	type testCase struct {
		address        Address
		nonce          string
		appID          string
		expectedResult string
		expectedErr    error
	}

	for name, tc := range map[string]testCase{
		"valid inputs": {
			address: HexToAddress("ABC123DEF456"),
			nonce:   "3037366134636339643564623330316636626239323161663465346131393662",
			appID:   "AWESOME-APP-ID",
			// nolint: lll
			expectedResult: "f8398e415745534f4d452d4150502d4944880000abc123def456a03037366134636339643564623330316636626239323161663465346131393662",
		},
		"nonce invalid hex": {
			address:     HexToAddress("ABC123DEF456"),
			nonce:       "asdf",
			appID:       "AWESOME-APP-ID",
			expectedErr: ErrInvalidNonce,
		},
		"nonce too short": {
			address:     HexToAddress("ABC123DEF456"),
			nonce:       "222222",
			appID:       "AWESOME-APP-ID",
			expectedErr: ErrInvalidNonce,
		},
	} {
		t.Run(name, func(t *testing.T) {
			// Check the output of NewAccountProofMessage against a pre-generated message from the flow-js-sdk
			msg, err := NewAccountProofMessage(tc.address, tc.appID, tc.nonce)
			if tc.expectedErr != nil {
				assert.True(t, errors.Is(err, tc.expectedErr))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, hex.EncodeToString(msg))
			}
		})
	}
}

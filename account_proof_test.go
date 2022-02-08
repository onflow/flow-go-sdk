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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccountProofMsg(t *testing.T) {
	type testCase struct {
		addr           string
		timestamp      int64
		appDomainTag   string
		expectedResult string
		expectedErr    error
	}

	for name, tc := range map[string]testCase{
		"with domain tag": {
			addr:           "1234123412341234",
			timestamp:      int64(1644259653161),
			appDomainTag:   "AWESOME-APP-V0.0-user",
			expectedResult: "f852b8403431353734353533346634643435326434313530353032643536333032653330326437353733363537323030303030303030303030303030303030303030303088123412341234123486017ed5833629",
		},
		"without domain tag": {
			addr:           "1234123412341234",
			timestamp:      int64(1644259653161),
			expectedResult: "d088123412341234123486017ed5833629",
		},
	} {
		t.Run(name, func(t *testing.T) {
			// Check the output of NewAccountProofMsg against a pre-generated message from the flow-js-sdk
			msg, err := NewAccountProofMsg(tc.addr, tc.timestamp, tc.appDomainTag)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResult, hex.EncodeToString(msg))
		})
	}
}

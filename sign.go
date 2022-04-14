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

package flow

import (
	"fmt"

	"github.com/onflow/flow-go-sdk/crypto"
)

const domainTagLength = 32

// TransactionDomainTag is the prefix of all signed transaction payloads.
//
// A domain tag is encoded as UTF-8 bytes, right padded to a total length of 32 bytes.
var TransactionDomainTag = mustPadDomainTag("FLOW-V0.0-transaction")

// UserDomainTag is the prefix of all signed user space payloads.
//
// A domain tag is encoded as UTF-8 bytes, right padded to a total length of 32 bytes.
var UserDomainTag = mustPadDomainTag("FLOW-V0.0-user")

func mustPadDomainTag(s string) [domainTagLength]byte {
	paddedTag, err := padDomainTag(s)
	if err != nil {
		panic(err)
	}

	return paddedTag
}

// padDomainTag returns a new padded domain tag from the given string. This function returns an error if the domain
// tag is too long.
func padDomainTag(tag string) (paddedTag [domainTagLength]byte, err error) {
	if len(tag) > domainTagLength {
		return paddedTag, fmt.Errorf("domain tag %s cannot be longer than %d characters", tag, domainTagLength)
	}

	copy(paddedTag[:], tag)

	return paddedTag, nil
}

// SignUserMessage signs a message in the user domain.
//
// User messages are distinct from other signed messages (i.e. transactions), and can be
// verified directly in on-chain Cadence code.
func SignUserMessage(signer crypto.Signer, message []byte) ([]byte, error) {
	message = append(UserDomainTag[:], message...)
	return signer.Sign(message)
}

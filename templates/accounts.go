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

package templates

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/onflow/flow-go-sdk"
)

// CreateAccount generates a script that creates a new account.
func CreateAccount(accountKeys []*flow.AccountKey, code []byte) ([]byte, error) {
	publicKeys := make([][]byte, len(accountKeys))

	for i, accountKey := range accountKeys {
		publicKeys[i] = accountKey.Encode()
	}

	publicKeysStr := languageEncodeBytesArray(publicKeys)

	script := fmt.Sprintf(
		`
            transaction {
              prepare(signer: AuthAccount) {
                let acct = AuthAccount(payer: signer)

                for key in %s {
                    acct.addPublicKey(key)
                }

                acct.setCode("%s".decodeHex())
              }
            }
        `,
		publicKeysStr,
		hex.EncodeToString(code),
	)

	return []byte(script), nil
}

// UpdateAccountCode generates a script that updates the code associated with an account.
func UpdateAccountCode(code []byte) []byte {
	return []byte(fmt.Sprintf(
		`
            transaction {
              prepare(signer: AuthAccount) {
                signer.setCode("%s".decodeHex())
              }
            }
        `,
		hex.EncodeToString(code),
	))
}

// AddAccountKey generates a script that adds a key to an account.
func AddAccountKey(accountKey *flow.AccountKey) ([]byte, error) {
	accountKeyBytes := accountKey.Encode()

	publicKeyStr := languageEncodeBytes(accountKeyBytes)

	script := fmt.Sprintf(`
        transaction {
          prepare(signer: AuthAccount) {
            signer.addPublicKey(%s)
          }
        }
    `, publicKeyStr)

	return []byte(script), nil
}

// RemoveAccountKey generates a script that removes a key from an account.
func RemoveAccountKey(id int) []byte {
	script := fmt.Sprintf(`
        transaction {
          prepare(signer: AuthAccount) {
            signer.removePublicKey(%d)
          }
        }
    `, id)

	return []byte(script)
}

// languageEncodeBytes converts a byte slice to a comma-separated list of uint8 integers.
func languageEncodeBytes(b []byte) string {
	if len(b) == 0 {
		return "[]"
	}

	return strings.Join(strings.Fields(fmt.Sprintf("%d", b)), ",")
}

// languageEncodeBytesArray converts a slice of byte slices to a comma-separated list of uint8 integers.
//
// Example: [][]byte{[]byte{1}, []byte{2,3}} -> "[[1],[2,3]]"
func languageEncodeBytesArray(b [][]byte) string {
	if len(b) == 0 {
		return "[]"
	}

	return strings.Join(strings.Fields(fmt.Sprintf("%d", b)), ",")
}

// languageEncodeIntArray converts a slice of integers to a comma-separated list.
func languageEncodeIntArray(i []int) string {
	if len(i) == 0 {
		return "[]"
	}

	return strings.Join(strings.Fields(fmt.Sprintf("%d", i)), ",")
}

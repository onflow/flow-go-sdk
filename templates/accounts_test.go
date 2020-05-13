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

package templates_test

import (
	"testing"

	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/assert"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/templates"
	"github.com/onflow/flow-go-sdk/test"
)

func TestCreateAccount(t *testing.T) {
	accountKey := test.AccountKeyGenerator().New()

	t.Run("without code", func(t *testing.T) {
		script, err := templates.CreateAccount([]*flow.AccountKey{accountKey}, []byte{})
		assert.NoError(t, err)

		expectedScript := `
          transaction {
            execute {
              AuthAccount(publicKeys: [[248,71,184,64,42,225,160,83,191,67,127,193,195,59,42,198,74,98,190,196,227,30,122,0,84,2,150,50,166,98,150,46,12,117,42,18,79,65,222,163,71,18,148,7,196,110,140,101,251,190,151,215,35,142,64,203,124,8,232,74,192,73,226,143,126,177,183,153,2,3,130,3,232]], code: "".decodeHex())
            }
          }
        `

		assert.Equal(t,
			dedent.Dedent(expectedScript),
			dedent.Dedent(string(script)),
		)
	})

	t.Run("with code", func(t *testing.T) {
		script, err := templates.CreateAccount([]*flow.AccountKey{accountKey}, []byte("pub fun main() {}"))
		assert.Nil(t, err)

		expectedScript := `
          transaction {
            execute {
              AuthAccount(publicKeys: [[248,71,184,64,42,225,160,83,191,67,127,193,195,59,42,198,74,98,190,196,227,30,122,0,84,2,150,50,166,98,150,46,12,117,42,18,79,65,222,163,71,18,148,7,196,110,140,101,251,190,151,215,35,142,64,203,124,8,232,74,192,73,226,143,126,177,183,153,2,3,130,3,232]], code: "7075622066756e206d61696e2829207b7d".decodeHex())
            }
          }
        `

		assert.Equal(t,
			dedent.Dedent(expectedScript),
			dedent.Dedent(string(script)),
		)
	})
}

func TestUpdateAccountCode(t *testing.T) {
	script := templates.UpdateAccountCode([]byte("pub fun main() {}"))

	expectedScript := `
      transaction {
        prepare(signer: AuthAccount) {
          signer.setCode("7075622066756e206d61696e2829207b7d".decodeHex())
        }
      }
    `

	assert.Equal(t,
		dedent.Dedent(expectedScript),
		dedent.Dedent(string(script)),
	)
}

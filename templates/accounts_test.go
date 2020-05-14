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
              AuthAccount(publicKeys: [[248,71,184,64,141,166,11,217,138,130,124,135,226,22,34,197,7,10,227,238,68,10,191,9,39,213,219,51,249,101,44,177,48,62,184,160,77,254,65,222,162,201,234,100,238,131,238,141,124,141,6,141,184,56,108,123,171,152,105,74,249,86,224,253,174,55,24,78,2,3,130,3,232]], code: "".decodeHex())
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
              AuthAccount(publicKeys: [[248,71,184,64,141,166,11,217,138,130,124,135,226,22,34,197,7,10,227,238,68,10,191,9,39,213,219,51,249,101,44,177,48,62,184,160,77,254,65,222,162,201,234,100,238,131,238,141,124,141,6,141,184,56,108,123,171,152,105,74,249,86,224,253,174,55,24,78,2,3,130,3,232]], code: "7075622066756e206d61696e2829207b7d".decodeHex())
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

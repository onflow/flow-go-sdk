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
              AuthAccount(publicKeys: [[248,72,184,64,199,209,247,158,141,105,106,46,33,152,142,7,81,171,181,156,100,170,60,92,218,125,250,195,229,235,105,192,11,150,121,14,251,225,162,132,64,20,237,172,176,86,201,233,29,187,31,229,168,190,133,254,90,11,87,239,249,83,170,123,0,38,93,140,2,3,130,3,232,42]], code: "".decodeHex())
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
              AuthAccount(publicKeys: [[248,72,184,64,199,209,247,158,141,105,106,46,33,152,142,7,81,171,181,156,100,170,60,92,218,125,250,195,229,235,105,192,11,150,121,14,251,225,162,132,64,20,237,172,176,86,201,233,29,187,31,229,168,190,133,254,90,11,87,239,249,83,170,123,0,38,93,140,2,3,130,3,232,42]], code: "7075622066756e206d61696e2829207b7d".decodeHex())
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

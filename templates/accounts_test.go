/*
 * Flow Go SDK
 *
 * Copyright Flow Foundation
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

	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/templates"
)

func TestCreateAccount(t *testing.T) {
	// Converting transaction arguments to Cadence values can increase their size.
	// If this is not taken into account the transaction can quickly grow over the maximum transaction size limit.
	t.Run("Transaction should not grow uncontrollably in size", func(t *testing.T) {
		testContractBody := "test contract"
		tx, err := templates.CreateAccount(
			[]*flow.AccountKey{},
			[]templates.Contract{{
				Name:   "contract",
				Source: testContractBody,
			}},
			flow.HexToAddress("01"))

		require.NoError(t, err)

		txSize := len(tx.Script)
		argumentsSize := 0
		for _, argument := range tx.Arguments {
			argumentsSize += len(argument)
		}
		require.Less(t, txSize, 1000, "The create account script should not grow over 1kB.")
		require.Less(t, argumentsSize, len(testContractBody)+500,
			"The create account argument size should not grow over 500 bytes of extra data.")
	})
}

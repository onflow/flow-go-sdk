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
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/templates"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	// Converting transaction arguments to Cadence values can increase their size.
	// If this is not taken into account the transaction can quickly grow over the maximum transaction size limit.
	t.Run("Transaction should not grow uncontrollably in size", func(t *testing.T) {
		contractLen := 1000
		contractCode := make([]byte, contractLen)

		tx := templates.CreateAccount(
			[]*flow.AccountKey{},
			[]templates.Contract{{
				Name:   "contract",
				Source: string(contractCode),
			}},
			flow.HexToAddress("01"))

		txSize := len(tx.Script)
		argumentsSize := 0
		for _, argument := range tx.Arguments {
			argumentsSize += len(argument)
		}
		require.Less(t, txSize, 1000, "The create account script should not grow over 1kB.")
		require.Less(t, argumentsSize, contractLen*2+500,
			"The create account argument size should not grow over "+
				"2 times the contract code (converted to hex) + 500 bytes of extra data.")
	})
}

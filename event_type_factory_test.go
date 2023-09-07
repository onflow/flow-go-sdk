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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventTypeFactory(t *testing.T) {
	assert.Equal(t, "A.7e60df042a9c0868.FlowToken.AccountCreated", NewEventTypeFactory().
		WithEventName("AccountCreated").
		WithAddressString("7e60df042a9c0868").
		WithContractName("FlowToken").
		String())
}

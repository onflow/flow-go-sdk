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
)

type eventNameFactory struct {
	address      string
	contractName string
	eventName    string
}

func (f eventNameFactory) WithAddressString(address string) eventNameFactory {
	f.address = address
	return f
}

func (f eventNameFactory) WithAddress(address Address) eventNameFactory {
	f.address = address.Hex()
	return f
}

func (f eventNameFactory) WithContractName(contract string) eventNameFactory {
	f.contractName = contract
	return f
}

func (f eventNameFactory) WithEventName(event string) eventNameFactory {
	f.eventName = event
	return f
}

func (f eventNameFactory) Build() string {
	return fmt.Sprintf("A.%s.%s.%s", f.address, f.contractName, f.eventName)
}

// NewEvent helper function for constructing event names
func NewEvent() eventNameFactory {
	return eventNameFactory{}
}

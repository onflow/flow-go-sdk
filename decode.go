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

package flow

import (
	// NOTE: always import Cadence's stdlib package,
	// as it registers the type ID decoder for the Flow types,
	// e.g. `flow.AccountCreated`
	_ "github.com/onflow/cadence/stdlib"
	"github.com/onflow/flow/protobuf/go/flow/entities"
)

type EventEncodingVersion = entities.EventEncodingVersion

const (
	EventEncodingVersionCCF     = entities.EventEncodingVersion_CCF_V0
	EventEncodingVersionJSONCDC = entities.EventEncodingVersion_JSON_CDC_V0
)

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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/onflow/cadence/runtime/common"
)

const (
	EVMLocationPrefix = "evm"
	locationDivider   = "."
)

var _ common.Location = EVMLocation{}

type EVMLocation struct{}

func (l EVMLocation) TypeID(
	memoryGauge common.MemoryGauge,
	qualifiedIdentifier string,
) common.TypeID {
	id := fmt.Sprintf(
		"%s%s%s",
		EVMLocationPrefix,
		locationDivider,
		qualifiedIdentifier,
	)
	common.UseMemory(memoryGauge, common.NewRawStringMemoryUsage(len(id)))

	return common.TypeID(id)
}

func (l EVMLocation) QualifiedIdentifier(typeID common.TypeID) string {
	pieces := strings.SplitN(string(typeID), locationDivider, 2)

	if len(pieces) < 2 {
		return ""
	}

	return pieces[1]
}

func (l EVMLocation) String() string {
	return EVMLocationPrefix
}

func (l EVMLocation) Description() string {
	return EVMLocationPrefix
}

func (l EVMLocation) ID() string {
	return EVMLocationPrefix
}

func (l EVMLocation) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type string
	}{
		Type: "EVMLocation",
	})
}

func init() {
	common.RegisterTypeIDDecoder(
		EVMLocationPrefix,
		func(_ common.MemoryGauge, typeID string) (common.Location, string, error) {
			if typeID == "" {
				return nil, "", fmt.Errorf("invalid EVM type location ID: missing type prefix")
			}

			parts := strings.SplitN(typeID, ".", 2)
			prefix := parts[0]
			if prefix != EVMLocationPrefix {
				return EVMLocation{}, "", fmt.Errorf("invalid EVM type location ID: invalid prefix")
			}

			var qualifiedIdentifier string
			pieceCount := len(parts)
			if pieceCount > 1 {
				qualifiedIdentifier = parts[1]
			}

			return EVMLocation{}, qualifiedIdentifier, nil
		},
	)
}

package abi_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapperlabs/flow-go-sdk/language/abi"
	"github.com/dapperlabs/flow-go-sdk/utils/unittest"
)

func TestExamples(t *testing.T) {
	for _, abiName := range abi.AssetNames() {
		suffix := ".abi.json"

		if strings.HasSuffix(abiName, suffix) {
			cdcName := abiName[:len(abiName)-len(suffix)]

			cdcAsset, _ := abi.Asset(cdcName)

			if cdcAsset != nil {

				t.Run(abiName, func(t *testing.T) {
					abiAsset, err := abi.Asset(abiName)
					require.NoError(t, err)

					typesFromABI, err := abi.GetTypesFromABIJSONBytes(abiAsset)

					assert.NoError(t, err)

					typesFromCadence := abi.GetTypesFromCadenceCode(string(cdcAsset), cdcName)

					unittest.AssertEqualWithDiff(t, typesFromCadence, typesFromABI)
				})
			}
		}
	}
}

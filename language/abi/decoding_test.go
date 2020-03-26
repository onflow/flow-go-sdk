package abi

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapperlabs/flow-go-sdk/utils/unittest"
)

func TestDecodeExamples(t *testing.T) {
	// TODO: skip
	t.Skip()

	for _, abiName := range AssetNames() {
		suffix := ".abi.json"

		if strings.HasSuffix(abiName, suffix) {
			cdcName := abiName[:len(abiName)-len(suffix)]

			cdcAsset, _ := Asset(cdcName)

			if cdcAsset != nil {

				t.Run(abiName, func(t *testing.T) {
					abiAsset, err := Asset(abiName)
					require.NoError(t, err)

					typesFromABI, err := GetTypesFromABIJSONBytes(abiAsset)

					assert.NoError(t, err)

					typesFromCadence := GetTypesFromCadenceCode(string(cdcAsset), cdcName)

					unittest.AssertEqualWithDiff(t, typesFromCadence, typesFromABI)
				})
			}
		}
	}
}

package abi

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/nsf/jsondiff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xeipuuv/gojsonschema"
)

func TestEncodeExamples(t *testing.T) {

	for _, assetName := range AssetNames() {
		if strings.HasSuffix(assetName, ".cdc") {

			abiAssetName := assetName + ".abi.json"

			abiAsset, _ := Asset(abiAssetName)

			t.Run(assetName, func(t *testing.T) {

				assetBytes, err := Asset(assetName)
				require.NoError(t, err)

				generatedAbi := GetABIJSONFromCadenceCode(string(assetBytes), false, assetName)

				options := jsondiff.DefaultConsoleOptions()
				diff, _ := jsondiff.Compare(generatedAbi, abiAsset, &options)

				assert.Equal(t, diff, jsondiff.FullMatch)
			})
		}
	}
}

func TestConformanceToSchema(t *testing.T) {

	abiBytes, err := ioutil.ReadFile("abi.json")
	require.NoError(t, err)

	schemaLoader := gojsonschema.NewBytesLoader(abiBytes)

	schema, err := gojsonschema.NewSchema(schemaLoader)
	require.NoError(t, err)

	for _, assetName := range AssetNames() {
		if strings.HasSuffix(assetName, ".abi.json") {

			t.Run(assetName, func(t *testing.T) {

				assetBytes, err := Asset(assetName)
				require.NoError(t, err)

				documentLoader := gojsonschema.NewBytesLoader(assetBytes)

				result, err := schema.Validate(documentLoader)
				require.NoError(t, err)

				if !assert.True(t, result.Valid()) {
					for _, desc := range result.Errors() {
						t.Logf("- %s\n", desc)
					}
				}
			})
		}
	}
}

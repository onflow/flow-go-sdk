package abi

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/dapperlabs/cadence"
	"github.com/dapperlabs/cadence/runtime/cmd"
	"github.com/dapperlabs/cadence/runtime/sema"
)

// GenerateABI generates ABIs from provided Cadence file
func GenerateABI(args []string, pretty bool) error {
	if len(args) < 1 {
		return errors.New("no input file given")
	}

	jsonData := GetABIJSONFromCadenceFile(args[0], pretty)

	_, err := os.Stdout.Write(jsonData)

	return err
}

func exportTypesFromChecker(checker *sema.Checker) map[string]cadence.Type {
	exportedTypes := map[string]cadence.Type{}

	values := checker.UserDefinedValues()
	for _, variable := range values {
		exportedTypes[variable.Identifier] = cadence.ConvertType(variable.Type)
	}

	return exportedTypes
}

func encodeTypesAsJSON(types map[string]cadence.Type, pretty bool) ([]byte, error) {
	encoder := NewEncoder()

	for name, typ := range types {
		encoder.Encode(name, typ)
	}

	if pretty {
		return json.MarshalIndent(encoder.Get(), "", "  ")
	}
	return json.Marshal(encoder.Get())
}

func GetABIJSONFromCadenceCode(code string, pretty bool, filename string) []byte {
	checker, _ := cmd.PrepareChecker(code, filename)

	exportedTypes := exportTypesFromChecker(checker)

	jsonData, err := encodeTypesAsJSON(exportedTypes, pretty)

	if err != nil {
		panic(err)
	}

	return jsonData
}

func GetABIJSONFromCadenceFile(filename string, pretty bool) []byte {
	checker, _ := cmd.PrepareCheckerFromFile(filename)

	exportedTypes := exportTypesFromChecker(checker)

	jsonData, err := encodeTypesAsJSON(exportedTypes, pretty)

	if err != nil {
		panic(err)
	}

	return jsonData
}

func GetTypesFromCadenceFile(filename string) map[string]cadence.Type {
	checker, _ := cmd.PrepareCheckerFromFile(filename)

	exportedTypes := exportTypesFromChecker(checker)

	return exportedTypes
}

func GetTypesFromCadenceCode(code string, filename string) map[string]cadence.Type {
	checker, _ := cmd.PrepareChecker(code, filename)

	exportedTypes := exportTypesFromChecker(checker)

	return exportedTypes
}

func GetTypesFromABIJSONBytes(bytes []byte) (map[string]cadence.Type, error) {
	return Decode(bytes)
}

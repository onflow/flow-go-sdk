package codegen

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dapperlabs/flow-go-sdk/language/abi"
	"github.com/dapperlabs/flow-go-sdk/language/types"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	if len(os.Args) != 4 {
		println("use package_name input_file output_file")
		os.Exit(1)
	}

	pkg := os.Args[1]
	inputFile := os.Args[2]
	outputFile := os.Args[3]

	data, err := ioutil.ReadFile(inputFile)
	check(err)

	allTypes, err := abi.Decode(data)

	compositeTypes := map[string]types.Composite{}

	for name, typ := range allTypes {

		switch composite := typ.(type) {
		case types.Resource:
			compositeTypes[name] = composite.Composite
		case types.Struct:
			compositeTypes[name] = composite.Composite
		default:
			_, err := fmt.Fprintf(os.Stderr, "Definition %s of type %T is not supported, skipping\n", name, typ)
			check(err)
		}

		if composite, ok := typ.(types.Composite); ok {
			compositeTypes[name] = composite
		}
	}

	file, err := os.Create(outputFile)
	defer file.Close()

	check(err)

	err = GenerateGo(pkg, compositeTypes, file)
	check(err)
}

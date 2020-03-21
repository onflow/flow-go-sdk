package abi

import (
	"github.com/dapperlabs/cadence"
)

func jsonStringToType(jsonString string) cadence.Type {
	switch jsonString {
	case "AnyStruct":
		return cadence.AnyStructType{}
	case "AnyResource":
		return cadence.AnyResourceType{}
	case "Bool":
		return cadence.BoolType{}
	case "Void":
		return cadence.VoidType{}
	case "String":
		return cadence.StringType{}
	case "Int":
		return cadence.IntType{}
	case "Int8":
		return cadence.Int8Type{}
	case "Int16":
		return cadence.Int16Type{}
	case "Int32":
		return cadence.Int32Type{}
	case "Int64":
		return cadence.Int64Type{}
	case "UInt8":
		return cadence.UInt8Type{}
	case "UInt16":
		return cadence.UInt16Type{}
	case "UInt32":
		return cadence.UInt32Type{}
	case "UInt64":
		return cadence.UInt64Type{}
	}

	return nil
}

func typeToJSONString(t cadence.Type) string {
	switch t.(type) {
	case cadence.AnyStructType:
		return "AnyStruct"
	case cadence.AnyResourceType:
		return "AnyResource"
	case cadence.BoolType:
		return "Bool"
	case cadence.VoidType:
		return "Void"
	case cadence.StringType:
		return "String"
	case cadence.IntType:
		return "Int"
	case cadence.Int8Type:
		return "Int8"
	case cadence.Int16Type:
		return "Int16"
	case cadence.Int32Type:
		return "Int32"
	case cadence.Int64Type:
		return "Int64"
	case cadence.UInt8Type:
		return "UInt8"
	case cadence.UInt16Type:
		return "UInt16"
	case cadence.UInt32Type:
		return "UInt32"
	case cadence.UInt64Type:
		return "UInt64"
	}

	return ""
}

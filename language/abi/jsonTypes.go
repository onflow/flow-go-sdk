package abi

import (
	"github.com/dapperlabs/flow-go-sdk/language/types"
)

func jsonStringToType(jsonString string) types.Type {
	switch jsonString {
	case "AnyStruct":
		return types.AnyStruct{}
	case "AnyResource":
		return types.AnyResource{}
	case "Bool":
		return types.Bool{}
	case "Void":
		return types.Void{}
	case "String":
		return types.String{}
	case "Int":
		return types.Int{}
	case "UInt":
		return types.UInt{}
	case "Int8":
		return types.Int8{}
	case "Int16":
		return types.Int16{}
	case "Int32":
		return types.Int32{}
	case "Int64":
		return types.Int64{}
	case "Int128":
		return types.Int128{}
	case "Int256":
		return types.Int256{}
	case "UInt8":
		return types.UInt8{}
	case "UInt16":
		return types.UInt16{}
	case "UInt32":
		return types.UInt32{}
	case "UInt64":
		return types.UInt64{}
	case "UInt128":
		return types.UInt128{}
	case "UInt256":
		return types.UInt256{}
	case "Word8":
		return types.Word8{}
	case "Word16":
		return types.Word16{}
	case "Word32":
		return types.Word32{}
	case "Word64":
		return types.Word64{}
	}

	return nil
}

func typeToJSONString(t types.Type) string {
	switch t.(type) {
	case types.AnyStruct:
		return "AnyStruct"
	case types.AnyResource:
		return "AnyResource"
	case types.Bool:
		return "Bool"
	case types.Void:
		return "Void"
	case types.String:
		return "String"
	case types.Int:
		return "Int"
	case types.UInt:
		return "UInt"
	case types.Int8:
		return "Int8"
	case types.Int16:
		return "Int16"
	case types.Int32:
		return "Int32"
	case types.Int64:
		return "Int64"
	case types.Int128:
		return "Int128"
	case types.Int256:
		return "Int256"
	case types.UInt8:
		return "UInt8"
	case types.UInt16:
		return "UInt16"
	case types.UInt32:
		return "UInt32"
	case types.UInt64:
		return "UInt64"
	case types.UInt128:
		return "UInt128"
	case types.UInt256:
		return "UInt256"
	case types.Word8:
		return "Word8"
	case types.Word16:
		return "Word16"
	case types.Word32:
		return "Word32"
	case types.Word64:
		return "Word64"
	}

	return ""
}

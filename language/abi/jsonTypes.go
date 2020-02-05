package abi

import (
	"github.com/dapperlabs/flow-go/language"
)

func jsonStringToType(jsonString string) language.Type {
	switch jsonString {
	case "AnyStruct":
		return language.AnyStructType{}
	case "AnyResource":
		return language.AnyResourceType{}
	case "Bool":
		return language.BoolType{}
	case "Void":
		return language.VoidType{}
	case "String":
		return language.StringType{}
	case "Int":
		return language.IntType{}
	case "Int8":
		return language.Int8Type{}
	case "Int16":
		return language.Int16Type{}
	case "Int32":
		return language.Int32Type{}
	case "Int64":
		return language.Int64Type{}
	case "UInt8":
		return language.UInt8Type{}
	case "UInt16":
		return language.UInt16Type{}
	case "UInt32":
		return language.UInt32Type{}
	case "UInt64":
		return language.UInt64Type{}
	}

	return nil
}

func typeToJSONString(t language.Type) string {
	switch t.(type) {
	case language.AnyStructType:
		return "AnyStruct"
	case language.AnyResourceType:
		return "AnyResource"
	case language.BoolType:
		return "Bool"
	case language.VoidType:
		return "Void"
	case language.StringType:
		return "String"
	case language.IntType:
		return "Int"
	case language.Int8Type:
		return "Int8"
	case language.Int16Type:
		return "Int16"
	case language.Int32Type:
		return "Int32"
	case language.Int64Type:
		return "Int64"
	case language.UInt8Type:
		return "UInt8"
	case language.UInt16Type:
		return "UInt16"
	case language.UInt32Type:
		return "UInt32"
	case language.UInt64Type:
		return "UInt64"
	}

	return ""
}

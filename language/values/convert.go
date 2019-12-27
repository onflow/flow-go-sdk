package values

import (
	"fmt"
	"math/big"

	"github.com/dapperlabs/flow-go/language/runtime"
	"github.com/dapperlabs/flow-go/language/runtime/interpreter"

	"github.com/dapperlabs/flow-go-sdk/utils/sort"
)

// Convert converts a runtime value to its corresponding Go representation.
func Convert(value runtime.Value) (Value, error) {
	switch v := value.(type) {
	case interpreter.VoidValue:
		return NewVoid(), nil
	case interpreter.NilValue:
		return NewOptional(nil), nil
	case *interpreter.SomeValue:
		return convertSomeValue(v)
	case interpreter.BoolValue:
		return NewBool(bool(v)), nil
	case *interpreter.StringValue:
		return NewString(v.Str), nil
	case *interpreter.ArrayValue:
		return convertArrayValue(v)
	case interpreter.IntValue:
		return NewIntFromBig(big.NewInt(0).Set(v.Int)), nil
	case interpreter.Int8Value:
		return NewInt8(int8(v)), nil
	case interpreter.Int16Value:
		return NewInt16(int16(v)), nil
	case interpreter.Int32Value:
		return NewInt32(int32(v)), nil
	case interpreter.Int64Value:
		return NewInt64(int64(v)), nil
	case interpreter.UInt8Value:
		return NewUInt8(uint8(v)), nil
	case interpreter.UInt16Value:
		return NewUInt16(uint16(v)), nil
	case interpreter.UInt32Value:
		return NewUInt32(uint32(v)), nil
	case interpreter.UInt64Value:
		return NewUInt64(uint64(v)), nil
	case *interpreter.CompositeValue:
		return convertCompositeValue(v)
	case *interpreter.DictionaryValue:
		return convertDictionaryValue(v)
	case *interpreter.EventValue:
		return convertEventValue(v)
	case interpreter.AddressValue:
		return NewAddress(v), nil
	}

	return nil, fmt.Errorf("cannot convert value of type %T", value)
}

func convertSomeValue(v *interpreter.SomeValue) (Value, error) {
	convertedValue, err := Convert(v.Value)
	if err != nil {
		return nil, err
	}

	return NewOptional(convertedValue), nil
}

func convertArrayValue(v *interpreter.ArrayValue) (Value, error) {
	vals := make([]Value, len(v.Values))

	for i, value := range v.Values {
		convertedValue, err := Convert(value)
		if err != nil {
			return nil, err
		}

		vals[i] = convertedValue
	}

	return NewVariableSizedArray(vals), nil
}

func convertCompositeValue(v *interpreter.CompositeValue) (Value, error) {
	fields := make([]Value, len(v.Fields))

	keys := make([]string, 0, len(v.Fields))
	for key := range v.Fields {
		keys = append(keys, key)
	}

	sort.LexicographicalOrder(keys)

	for i, key := range keys {
		field := v.Fields[key]

		convertedField, err := Convert(field)
		if err != nil {
			return nil, err
		}

		fields[i] = convertedField
	}

	return NewComposite(fields), nil
}

func convertDictionaryValue(v *interpreter.DictionaryValue) (Value, error) {
	pairs := make([]KeyValuePair, v.Count())

	for i, keyValue := range v.Keys.Values {
		key := keyValue.(interpreter.HasKeyString).KeyString()
		value := v.Entries[key]

		convertedKey, err := Convert(keyValue)
		if err != nil {
			return nil, err
		}

		convertedValue, err := Convert(value)
		if err != nil {
			return nil, err
		}

		pairs[i] = KeyValuePair{
			Key:   convertedKey,
			Value: convertedValue,
		}
	}

	return NewDictionary(pairs), nil
}

func convertEventValue(v *interpreter.EventValue) (Value, error) {
	fields := make([]Value, len(v.Fields))

	for i, field := range v.Fields {
		convertedField, err := Convert(field.Value)
		if err != nil {
			return nil, err
		}

		fields[i] = convertedField
	}

	return NewEvent(fields), nil
}

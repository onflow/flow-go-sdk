package abi

import (
	"encoding/json"
	"fmt"

	"github.com/dapperlabs/cadence"
)

func getOnlyEntry(m map[string]interface{}) (string, interface{}, error) {
	if len(m) > 1 {
		return "", nil, fmt.Errorf("more than one entry in %v", m)

	}
	for k, v := range m {
		return k, v, nil
	}
	return "", nil, fmt.Errorf("no entries, but one required in %v", m)
}

func getString(m map[string]interface{}, key string) (string, error) {
	value, err := getObject(m, key)
	if err != nil {
		return "", nil
	}

	if s, ok := value.(string); ok {
		return s, nil
	}

	return "", fmt.Errorf("value for key  %s it is not a string in %v", key, m)
}

func getUInt(m map[string]interface{}, key string) (uint, error) {
	value, err := getObject(m, key)
	if err != nil {
		return 0, err
	}

	if s, ok := value.(float64); ok {
		if s >= 0 {
			return uint(s), nil
		}
	}
	return 0, fmt.Errorf("value for key  %s it %t, expected uint", key, value)
}

func getArray(m map[string]interface{}, key string) ([]interface{}, error) {
	value, err := getObject(m, key)
	if err != nil {
		return nil, nil
	}

	if s, ok := value.([]interface{}); ok {
		return s, nil
	}

	return nil, fmt.Errorf("value for key  %s it is not an array in %v", key, m)
}

func getMap(m map[string]interface{}, key string) (map[string]interface{}, error) {
	value, err := getObject(m, key)
	if err != nil {
		return nil, nil
	}

	if s, ok := value.(map[string]interface{}); ok {
		return s, nil
	}

	return nil, fmt.Errorf("value for key  %s it is not a map in %v", key, m)

}

func getIndex(a []interface{}, index int) (interface{}, error) {
	if len(a) <= index || index < 0 {
		return nil, fmt.Errorf("index %d doesn't exist in array in %v", index, a)

	}
	return a[index], nil
}

func getObject(data map[string]interface{}, key string) (interface{}, error) {
	v, ok := data[key]

	if ok {
		return v, nil
	}

	return nil, fmt.Errorf("key %s doesn't exist  in %v", key, data)
}

func toField(data map[string]interface{}) (cadence.Field, error) {
	name, err := getString(data, "name")
	if err != nil {
		return cadence.Field{}, err
	}

	typRaw, err := getObject(data, "type")
	if err != nil {
		return cadence.Field{}, err
	}

	typ, err := toType(typRaw, "")
	if err != nil {
		return cadence.Field{}, err
	}

	return cadence.Field{
		Identifier: name,
		Type:       typ,
	}, nil
}

func toFields(fields []map[string]interface{}) ([]cadence.Field, error) {
	ret := make([]cadence.Field, len(fields))

	for i, raw := range fields {
		f, err := toField(raw)
		if err != nil {
			return nil, err
		}

		ret[i] = f
	}

	return ret, nil
}

func toParameter(data map[string]interface{}) (cadence.Parameter, error) {
	name, err := getString(data, "name")
	if err != nil {
		return cadence.Parameter{}, err
	}

	label, err := getString(data, "label")
	if err != nil {
		label = ""
	}

	typRaw, err := getObject(data, "type")
	if err != nil {
		return cadence.Parameter{}, err
	}

	typ, err := toType(typRaw, "")
	if err != nil {
		return cadence.Parameter{}, err
	}

	return cadence.Parameter{
		Label:      label,
		Identifier: name,
		Type:       typ,
	}, nil
}

func interfaceToListOfMaps(input interface{}) ([]map[string]interface{}, error) {
	array, ok := input.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%v is not of expected type []interface{}", input)
	}

	ret := make([]map[string]interface{}, len(array))
	for i, a := range array {
		ret[i], ok = a.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%v is not of expected type map[string]interface{}", a)
		}
	}
	return ret, nil
}

func toComposite(data map[string]interface{}, name string) (
	string,
	[]cadence.Field,
	[][]cadence.Parameter,
	error,
) {
	fieldsRaw, err := getArray(data, "fields")
	if err != nil {
		return "", nil, nil, err
	}

	fieldsMaps, err := interfaceToListOfMaps(fieldsRaw)
	if err != nil {
		return "", nil, nil, err
	}

	fields, err := toFields(fieldsMaps)
	if err != nil {
		return "", nil, nil, err
	}

	initializersRaw, err := getArray(data, "initializers")
	if err != nil {
		return "", nil, nil, err
	}

	initializerRaw, err := getIndex(initializersRaw, 0)
	if err != nil {
		return "", nil, nil, err
	}

	initializer, err := interfaceToListOfMaps(initializerRaw)
	if err != nil {
		return "", nil, nil, err
	}

	parameters, err := toParameters(initializer)
	if err != nil {
		return "", nil, nil, err
	}

	initializers := [][]cadence.Parameter{parameters}
	return name, fields, initializers, nil
}

func toStruct(data map[string]interface{}, name string) (cadence.StructType, error) {
	identifier, fields, initializers, err := toComposite(data, name)
	if err != nil {
		return cadence.StructType{}, err
	}

	return cadence.StructType{
		// TODO:
		TypeID:       "",
		Identifier:   identifier,
		Fields:       fields,
		Initializers: initializers,
	}, nil
}

func toResource(data map[string]interface{}, name string) (cadence.ResourceType, error) {
	identifier, fields, initializers, err := toComposite(data, name)
	if err != nil {
		return cadence.ResourceType{}, err
	}

	return cadence.ResourceType{
		// TODO:
		TypeID:       "",
		Identifier:   identifier,
		Fields:       fields,
		Initializers: initializers,
	}, nil
}

func toEvent(data map[string]interface{}, name string) (cadence.EventType, error) {
	identifier, fields, initializers, err := toComposite(data, name)
	if err != nil {
		return cadence.EventType{}, err
	}

	return cadence.EventType{
		// TODO:
		TypeID:       "",
		Identifier:   identifier,
		Fields:       fields,
		Initializer:  initializers[0],
	}, nil
}

func toParameters(parameters []map[string]interface{}) ([]cadence.Parameter, error) {
	ret := make([]cadence.Parameter, len(parameters))

	for i, raw := range parameters {
		p, err := toParameter(raw)
		ret[i] = p
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func toFunction(data map[string]interface{}) (cadence.Function, error) {
	returnTypeRaw, err := getObject(data, "returnType")

	var returnType cadence.Type

	if err != nil {
		returnType = cadence.VoidType{}
	} else {
		returnType, err = toType(returnTypeRaw, "")
		if err != nil {
			return cadence.Function{}, err
		}
	}

	parametersListRaw, err := getArray(data, "parameters")
	if err != nil {
		return cadence.Function{}, err
	}

	parametersRaw, err := interfaceToListOfMaps(parametersListRaw)
	if err != nil {
		return cadence.Function{}, err
	}

	parameters, err := toParameters(parametersRaw)
	if err != nil {
		return cadence.Function{}, err
	}

	return cadence.Function{
		Parameters: parameters,
		ReturnType: returnType,
	}, nil
}

func toArray(data map[string]interface{}) (cadence.Type, error) {

	ofRaw, err := getObject(data, "of")

	if err != nil {
		return nil, err
	}

	of, err := toType(ofRaw, "")
	if err != nil {
		return nil, err
	}

	hasSize := true

	size, err := getUInt(data, "size")
	if err != nil {
		hasSize = false
	}

	if hasSize {
		return cadence.ConstantSizedArrayType{
			Size:        size,
			ElementType: of,
		}, nil
	}
	return cadence.VariableSizedArrayType{
		ElementType: of,
	}, nil
}

func toDictionary(data map[string]interface{}) (cadence.Type, error) {

	keysRaw, err := getObject(data, "keys")

	if err != nil {
		return nil, err
	}

	keys, err := toType(keysRaw, "")
	if err != nil {
		return nil, err
	}

	elementsRaw, err := getObject(data, "values")

	if err != nil {
		return nil, err
	}

	elements, err := toType(elementsRaw, "")
	if err != nil {
		return nil, err
	}

	return cadence.DictionaryType{
		KeyType:     keys,
		ElementType: elements,
	}, nil
}

func toType(data interface{}, name string) (cadence.Type, error) {

	switch v := data.(type) {

	// Simple string cases - "Int"
	case string:

		if typ := jsonStringToType(v); typ != nil {
			return typ, nil
		}

		return nil, fmt.Errorf("unsupported name %s for simple string type", v)

	// If object with key as type descriptor - <{ "<function>": XX }>
	case map[string]interface{}:

		key, value, err := getOnlyEntry(v)
		if err != nil {
			return nil, err
		}

		// when type of declaration doesn't matter as we can handle both
		switch key {
		case jsonTypeVariable:
			typ, err := toType(value, name)
			if err != nil {
				return nil, err
			}
			return cadence.Variable{
				Type: typ,
			}, nil
		case "optional":
			typ, err := toType(value, name)
			if err != nil {
				return nil, err
			}
			return cadence.OptionalType{
				Type: typ,
			}, nil
		}

		// when case require more handling
		switch v := value.(type) {
		// when type inside is simple string - { "<struct>": "SimpleString" }
		case string:
			switch key {
			case "struct":
				return cadence.StructPointer{TypeName: v}, nil
			case "resource":
				return cadence.ResourcePointer{TypeName: v}, nil
			case "event":
				return cadence.EventPointer{TypeName: v}, nil
			}

		// when type inside is complex - { "<struct>" : { "complex": "object" } }
		case map[string]interface{}:
			switch key {
			case "struct":
				return toStruct(v, name)
			case "resource":
				return toResource(v, name)
			case "event":
				return toEvent(v, name)
			case "function":
				return toFunction(v)
			case "array":
				return toArray(v)
			case "dictionary":
				return toDictionary(v)

			}
		}

	}

	return nil, fmt.Errorf("unsupported data chunk %v", data)
}

type jsonContainer struct {
	Definitions map[string]map[string]interface{}
}

func Decode(bytes []byte) (map[string]cadence.Type, error) {

	jsonRoot := jsonContainer{}

	err := json.Unmarshal(bytes, &jsonRoot)

	if err != nil {
		panic(err)
	}

	definitions := map[string]cadence.Type{}

	for name, definition := range jsonRoot.Definitions {
		typ, err := toType(definition, name)
		if err != nil {
			return nil, err
		}
		definitions[name] = typ
	}

	return definitions, nil
}

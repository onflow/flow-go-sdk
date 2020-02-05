package abi

import (
	"encoding/json"
	"fmt"

	"github.com/dapperlabs/flow-go/language"
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

func toField(data map[string]interface{}) (language.Field, error) {
	name, err := getString(data, "name")
	if err != nil {
		return language.Field{}, err
	}

	typRaw, err := getObject(data, "type")
	if err != nil {
		return language.Field{}, err
	}

	typ, err := toType(typRaw, "")
	if err != nil {
		return language.Field{}, err
	}

	return language.Field{
		Identifier: name,
		Type:       typ,
	}, nil
}

func toFields(fields []map[string]interface{}) ([]language.Field, error) {
	ret := make([]language.Field, len(fields))

	for i, raw := range fields {
		f, err := toField(raw)
		if err != nil {
			return nil, err
		}

		ret[i] = f
	}

	return ret, nil
}

func toParameter(data map[string]interface{}) (language.Parameter, error) {
	name, err := getString(data, "name")
	if err != nil {
		return language.Parameter{}, err
	}

	label, err := getString(data, "label")
	if err != nil {
		label = ""
	}

	typRaw, err := getObject(data, "type")
	if err != nil {
		return language.Parameter{}, err
	}

	typ, err := toType(typRaw, "")
	if err != nil {
		return language.Parameter{}, err
	}

	return language.Parameter{
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

func toComposite(data map[string]interface{}, name string) (language.CompositeType, error) {
	fieldsRaw, err := getArray(data, "fields")
	if err != nil {
		return language.CompositeType{}, err
	}

	fieldsMaps, err := interfaceToListOfMaps(fieldsRaw)
	if err != nil {
		return language.CompositeType{}, err
	}

	fields, err := toFields(fieldsMaps)
	if err != nil {
		return language.CompositeType{}, err
	}

	initializersRaw, err := getArray(data, "initializers")
	if err != nil {
		return language.CompositeType{}, err
	}

	initializerRaw, err := getIndex(initializersRaw, 0)
	if err != nil {
		return language.CompositeType{}, err
	}

	initializers, err := interfaceToListOfMaps(initializerRaw)
	if err != nil {
		return language.CompositeType{}, err
	}

	parameters, err := toParameters(initializers)
	if err != nil {
		return language.CompositeType{}, err
	}

	return language.CompositeType{
		Identifier: name,
		Fields:     fields,
		Initializers: [][]language.Parameter{
			parameters,
		},
	}, nil
}

func toStruct(data map[string]interface{}, name string) (language.StructType, error) {
	composite, err := toComposite(data, name)
	if err != nil {
		return language.StructType{}, err
	}

	return language.StructType{
		CompositeType: composite,
	}, nil
}

func toResource(data map[string]interface{}, name string) (language.ResourceType, error) {
	composite, err := toComposite(data, name)
	if err != nil {
		return language.ResourceType{}, err
	}

	return language.ResourceType{
		CompositeType: composite,
	}, nil
}

func toEvent(data map[string]interface{}, name string) (language.EventType, error) {
	composite, err := toComposite(data, name)
	if err != nil {
		return language.EventType{}, err
	}

	return language.EventType{
		CompositeType: composite,
	}, nil
}

func toParameters(parameters []map[string]interface{}) ([]language.Parameter, error) {
	ret := make([]language.Parameter, len(parameters))

	for i, raw := range parameters {
		p, err := toParameter(raw)
		ret[i] = p
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func toFunction(data map[string]interface{}) (language.Function, error) {
	returnTypeRaw, err := getObject(data, "returnType")

	var returnType language.Type

	if err != nil {
		returnType = language.VoidType{}
	} else {
		returnType, err = toType(returnTypeRaw, "")
		if err != nil {
			return language.Function{}, err
		}
	}

	parametersListRaw, err := getArray(data, "parameters")
	if err != nil {
		return language.Function{}, err
	}

	parametersRaw, err := interfaceToListOfMaps(parametersListRaw)
	if err != nil {
		return language.Function{}, err
	}

	parameters, err := toParameters(parametersRaw)
	if err != nil {
		return language.Function{}, err
	}

	return language.Function{
		Parameters: parameters,
		ReturnType: returnType,
	}, nil
}

func toFunctionType(data map[string]interface{}) (language.FunctionType, error) {

	returnTypeRaw, err := getObject(data, "returnType")

	var returnType language.Type

	if err != nil {
		returnType = language.VoidType{}
	} else {
		returnType, err = toType(returnTypeRaw, "")
		if err != nil {
			return language.FunctionType{}, err
		}
	}

	parametersListRaw, err := getArray(data, "parameters")
	if err != nil {
		return language.FunctionType{}, err
	}

	parameterTypes := make([]language.Type, len(parametersListRaw))

	for i, parameterTypeRaw := range parametersListRaw {
		parameterTypes[i], err = toType(parameterTypeRaw, "")
		if err != nil {
			return language.FunctionType{}, err
		}
	}

	return language.FunctionType{
		ParameterTypes: parameterTypes,
		ReturnType:     returnType,
	}, nil
}

func toArray(data map[string]interface{}) (language.Type, error) {

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
		return language.ConstantSizedArrayType{
			Size:        size,
			ElementType: of,
		}, nil
	}
	return language.VariableSizedArrayType{
		ElementType: of,
	}, nil
}

func toDictionary(data map[string]interface{}) (language.Type, error) {

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

	return language.DictionaryType{
		KeyType:     keys,
		ElementType: elements,
	}, nil
}

func toType(data interface{}, name string) (language.Type, error) {

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
			return language.Variable{
				Type: typ,
			}, nil
		case "optional":
			typ, err := toType(value, name)
			if err != nil {
				return nil, err
			}
			return language.OptionalType{
				Type: typ,
			}, nil
		}

		// when case require more handling
		switch v := value.(type) {
		// when type inside is simple string - { "<struct>": "SimpleString" }
		case string:
			switch key {
			case "struct":
				return language.StructPointer{TypeName: v}, nil
			case "resource":
				return language.ResourcePointer{TypeName: v}, nil
			case "event":
				return language.EventPointer{TypeName: v}, nil
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
				if name != "" {
					return toFunction(v)
				}
				return toFunctionType(v)
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

func Decode(bytes []byte) (map[string]language.Type, error) {

	jsonRoot := jsonContainer{}

	err := json.Unmarshal(bytes, &jsonRoot)

	if err != nil {
		panic(err)
	}

	definitions := map[string]language.Type{}

	for name, definition := range jsonRoot.Definitions {
		typ, err := toType(definition, name)
		if err != nil {
			return nil, err
		}
		definitions[name] = typ
	}

	return definitions, nil
}

package abi

import (
	"fmt"

	"github.com/dapperlabs/flow-go/language"
)

func NewEncoder() Encoder {
	return Encoder{definitions: map[string]interface{}{}}
}

type Encoder struct {
	definitions map[string]interface{}
}

func (encoder *Encoder) Encode(name string, t language.Type) {
	encoder.definitions[name] = encoder.encode(t)
}

func (encoder *Encoder) Get() interface{} {
	return ABIObject{
		encoder.definitions,
		"", // Once we setup schema, probably on withflow.org
	}
}

// region JSON Structures

type ABIObject struct {
	Definitions map[string]interface{} `json:"definitions"`
	Schema      string                 `json:"schema,omitempty"`
}

type arrayObject struct {
	Array array `json:"array"`
}

type array struct {
	Of   interface{} `json:"of"`
	Size uint        `json:"size,omitempty"`
}

type structObject struct {
	Struct compositeData `json:"struct"`
}

type compositeData struct {
	Fields       []field       `json:"fields"`
	Initializers [][]parameter `json:"initializers"`
}

type field struct {
	Name string      `json:"name"`
	Type interface{} `json:"type"`
}

type parameter struct {
	field
	Label string `json:"label,omitempty"`
}

type optionalObject struct {
	Optional interface{} `json:"optional"`
}

type function struct {
	Parameters []parameter `json:"parameters,omitempty"`
	ReturnType interface{} `json:"returnType,omitempty"`
}

type functionType struct {
	Parameters []interface{} `json:"parameters,omitempty"`
	ReturnType interface{}   `json:"returnType,omitempty"`
}

type functionObject struct {
	Function interface{} `json:"function"`
}

type dictionary struct {
	Keys   interface{} `json:"keys"`
	Values interface{} `json:"values"`
}

type dictionaryObject struct {
	Dictionary dictionary `json:"dictionary"`
}

type resourceObject struct {
	Resource compositeData `json:"resource"`
}

type eventObject struct {
	Event compositeData `json:"event"`
}

type resourcePointer struct {
	Resource string `json:"resource"`
}

type structPointer struct {
	Struct string `json:"struct"`
}

type eventPointer struct {
	Event string `json:"event"`
}

type variableObject struct {
	Variable interface{} `json:"variable"`
}

// endregion

func (encoder *Encoder) mapFields(m []language.Field) []field {
	ret := make([]field, len(m))

	for i, f := range m {
		ret[i] = field{
			Name: f.Identifier,
			Type: encoder.encode(f.Type),
		}
	}

	return ret
}

func (encoder *Encoder) mapParameters(p []language.Parameter) []parameter {
	ret := make([]parameter, len(p))

	for i := range p {
		ret[i] = parameter{
			field: field{
				Name: p[i].Identifier,
				Type: encoder.encode(p[i].Type),
			},
			Label: p[i].Label,
		}
	}

	return ret
}

func (encoder *Encoder) mapNestedParameters(p [][]language.Parameter) [][]parameter {

	ret := make([][]parameter, len(p))
	for i := range ret {
		ret[i] = encoder.mapParameters(p[i])
	}

	return ret
}

func (encoder *Encoder) mapTypes(types []language.Type) []interface{} {
	ret := make([]interface{}, len(types))

	for i, t := range types {
		ret[i] = encoder.encode(t)
	}

	return ret
}

// For function return type Void is redundant, so we remove it
func (encoder *Encoder) encodeReturnType(returnType language.Type) interface{} {
	if _, ok := returnType.(language.VoidType); ok == true {
		return nil
	}
	return encoder.encode(returnType)
}

const jsonTypeVariable = "variable"

func (encoder *Encoder) encode(t language.Type) interface{} {

	if s := typeToJSONString(t); s != "" {
		return s
	}

	switch v := (t).(type) {

	case language.VariableSizedArrayType:
		return arrayObject{array{Of: encoder.encode(v.ElementType)}}
	case language.ConstantSizedArrayType:
		return arrayObject{array{Of: encoder.encode(v.ElementType), Size: v.Size}}

	case language.OptionalType:
		return optionalObject{Optional: encoder.encode(v.Type)}

	case language.StructType:
		return structObject{
			Struct: compositeData{
				Fields:       encoder.mapFields(v.Fields),
				Initializers: encoder.mapNestedParameters(v.Initializers),
			},
		}

	case language.StructPointer:
		return structPointer{
			v.TypeName,
		}

	case language.ResourceType:
		return resourceObject{
			compositeData{
				Fields:       encoder.mapFields(v.Fields),
				Initializers: encoder.mapNestedParameters(v.Initializers),
			},
		}

	case language.ResourcePointer:
		return resourcePointer{
			v.TypeName,
		}

	case language.EventType:
		return eventObject{
			compositeData{
				Fields:       encoder.mapFields(v.Fields),
				Initializers: encoder.mapNestedParameters(v.Initializers),
			},
		}

	case language.EventPointer:
		return eventPointer{
			v.TypeName,
		}

	case language.Function:
		return functionObject{
			function{
				Parameters: encoder.mapParameters(v.Parameters),
				ReturnType: encoder.encodeReturnType(v.ReturnType),
			},
		}

	case language.FunctionType:
		return functionObject{
			functionType{
				Parameters: encoder.mapTypes(v.ParameterTypes),
				ReturnType: encoder.encodeReturnType(v.ReturnType),
			},
		}

	case language.DictionaryType:
		return dictionaryObject{
			dictionary{
				Keys:   encoder.encode(v.KeyType),
				Values: encoder.encode(v.ElementType),
			},
		}

	case language.Variable:
		return variableObject{
			encoder.encode(v.Type),
		}

	}

	panic(fmt.Errorf("unknown type of %T", t))
}

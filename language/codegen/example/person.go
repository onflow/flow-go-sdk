package example

import (
	"bytes"

	"github.com/dapperlabs/flow-go-sdk/language/encoding"
	"github.com/dapperlabs/flow-go-sdk/language/types"
	"github.com/dapperlabs/flow-go-sdk/language/values"
)

type PersonView interface {
	FullName() string
}

type personView struct {
	_fullName string
	value     values.Composite
}

func (p personView) FullName() string {
	return p._fullName
}

type PersonConstructor interface {
	Encode() ([]byte, error)
}

type personConstructor struct {
	firstName string
	lastName  string
}

func (p personConstructor) Encode() ([]byte, error) {

	var w bytes.Buffer
	encoder := encoding.NewEncoder(&w)

	err := encoder.EncodeConstantSizedArray(
		values.NewConstantSizedArray([]values.Value{
			values.NewString(p.firstName),
			values.NewString(p.lastName),
		}),
	)

	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func NewPersonConstructor(firstName string, lastName string) (PersonConstructor, error) {
	return personConstructor{
		firstName: firstName,
		lastName:  lastName,
	}, nil
}

var personType = types.Composite{
	Fields: []types.Field{
		{
			Identifier: "FullName",
			Type:       types.String{},
		},
	},
	Initializers: [][]types.Parameter{
		{
			{
				Identifier: "firstName",
				Type:       types.String{},
			},
			{
				Identifier: "lastName",
				Type:       types.String{},
			},
		},
	},
}

func DecodePersonView(b []byte) (PersonView, error) {
	r := bytes.NewReader(b)
	dec := encoding.NewDecoder(r)

	v, err := dec.DecodeComposite(personType)
	if err != nil {
		return nil, err
	}

	return personView{
		_fullName: string(v.Fields[0].(values.String)),
	}, nil
}

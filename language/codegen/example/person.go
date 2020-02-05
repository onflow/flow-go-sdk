package example

import (
	"bytes"

	"github.com/dapperlabs/flow-go/language"
	"github.com/dapperlabs/flow-go/language/encoding"
)

type PersonView interface {
	FullName() string
}

type personView struct {
	_fullName string
	value     language.Composite
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
		language.NewConstantSizedArray([]language.Value{
			language.NewString(p.firstName),
			language.NewString(p.lastName),
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

var personType = language.CompositeType{
	Fields: []language.Field{
		{
			Identifier: "FullName",
			Type:       language.StringType{},
		},
	},
	Initializers: [][]language.Parameter{
		{
			{
				Identifier: "firstName",
				Type:       language.StringType{},
			},
			{
				Identifier: "lastName",
				Type:       language.StringType{},
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
		_fullName: string(v.Fields[0].(language.String)),
	}, nil
}

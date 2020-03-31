package example

import (
	"bytes"

	"github.com/dapperlabs/cadence"
	encoding "github.com/dapperlabs/cadence/encoding/xdr"
)

type PersonView interface {
	FullName() string
}

type personView struct {
	_fullName string
	value     cadence.Struct
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

	err := encoder.EncodeArray(
		cadence.NewArray([]cadence.Value{
			cadence.NewString(p.firstName),
			cadence.NewString(p.lastName),
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

var personType = cadence.StructType{
	Fields: []cadence.Field{
		{
			Identifier: "FullName",
			Type:       cadence.StringType{},
		},
	},
	Initializers: [][]cadence.Parameter{
		{
			{
				Identifier: "firstName",
				Type:       cadence.StringType{},
			},
			{
				Identifier: "lastName",
				Type:       cadence.StringType{},
			},
		},
	},
}

func DecodePersonView(b []byte) (PersonView, error) {
	r := bytes.NewReader(b)
	dec := encoding.NewDecoder(r)

	v, err := dec.DecodeStruct(personType)
	if err != nil {
		return nil, err
	}

	return personView{
		_fullName: string(v.Fields[0].(cadence.String)),
	}, nil
}

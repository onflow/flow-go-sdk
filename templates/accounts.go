/*
 * Flow Go SDK
 *
 * Copyright 2019 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package templates

import (
	"encoding/hex"
	"fmt"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/cadence/runtime"
	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/cadence/runtime/sema"
	"github.com/onflow/flow-go-sdk/crypto"
	templates "github.com/onflow/sdks"

	"github.com/onflow/flow-go-sdk"
)

// Contract is a Cadence contract deployed to a Flow account.
//
// Name is the identifier for the contract within the account.
//
// Source is the the Cadence code of the contract.
type Contract struct {
	Name   string
	Source string
}

// SourceBytes returns the UTF-8 encoded source code (Source) of the contract.
func (c Contract) SourceBytes() []byte {
	return []byte(c.Source)
}

// SourceHex returns the UTF-8 encoded source code (Source) of the contract as a hex string.
func (c Contract) SourceHex() string {
	return hex.EncodeToString(c.SourceBytes())
}

func exportType(t sema.Type) cadence.Type {
	return runtime.ExportType(t, map[sema.TypeID]cadence.Type{})
}

func newSignAlgoValue(sigAlgo crypto.SignatureAlgorithm) (cadence.Enum, error) {
	sigAlgoValue := sema.SignatureAlgorithmECDSA_P256
	switch sigAlgo {
	case crypto.ECDSA_P256:
		sigAlgoValue = sema.SignatureAlgorithmECDSA_P256
	case crypto.ECDSA_secp256k1:
		sigAlgoValue = sema.SignatureAlgorithmECDSA_secp256k1
	default:
		return cadence.Enum{}, fmt.Errorf("cannot encode signature algorithm to cadence value: unsupported signature algorithm: %v", sigAlgo)
	}

	return cadence.NewEnum([]cadence.Value{
		cadence.NewUInt8(sigAlgoValue.RawValue()),
	}).WithType(
		exportType(sema.SignatureAlgorithmType).(*cadence.EnumType),
	), nil
}

func newHashAlgoValue(hashAlgo crypto.HashAlgorithm) (cadence.Enum, error) {
	hashAlgoValue := sema.HashAlgorithmSHA2_256
	switch hashAlgo {
	case crypto.SHA2_256:
		hashAlgoValue = sema.HashAlgorithmSHA2_256
	case crypto.SHA2_384:
		hashAlgoValue = sema.HashAlgorithmSHA2_384
	case crypto.SHA3_256:
		hashAlgoValue = sema.HashAlgorithmSHA3_256
	case crypto.SHA3_384:
		hashAlgoValue = sema.HashAlgorithmSHA3_384
	default:
		return cadence.Enum{}, fmt.Errorf("cannot encode hash algorithm to cadence value: unsupported hash algorithm: %v", hashAlgo)
	}

	return cadence.NewEnum([]cadence.Value{
		cadence.NewUInt8(hashAlgoValue.RawValue()),
	}).WithType(
		exportType(sema.HashAlgorithmType).(*cadence.EnumType),
	), nil
}

func newPublicKeyValue(pubKey crypto.PublicKey) (cadence.Struct, error) {
	pubKeyCadence := make([]cadence.Value, len(pubKey.Encode()))
	for i, k := range pubKey.Encode() {
		pubKeyCadence[i] = cadence.NewUInt8(k)
	}

	sig, err := newSignAlgoValue(pubKey.Algorithm())
	if err != nil {
		return cadence.Struct{}, fmt.Errorf("cannot encode public key to cadence value: %w", err)
	}

	return cadence.NewStruct(
		[]cadence.Value{
			cadence.NewArray(pubKeyCadence),
			sig,
		},
	).WithType(
		exportType(sema.PublicKeyType).(*cadence.StructType),
	), nil
}

// AccountKeyToCadenceCryptoKey converts a `flow.AccountKey` key to the Cadence struct `Crypto.KeyListEntry`,
// so that it can more easily be used as a parameter in scripts and transactions.
//
// example:
// ```go
// 	key := AccountKeyToCadenceCryptoKey(accountKey)
//
//	return flow.NewTransaction().
//		SetScript([]byte(templates.AddAccountKey)).
//		AddRawArgument(jsoncdc.MustEncode(key))
// ```
func AccountKeyToCadenceCryptoKey(key *flow.AccountKey) (cadence.Value, error) {
	weight, _ := cadence.NewUFix64(fmt.Sprintf("%d.0", key.Weight))
	publicKey, err := newPublicKeyValue(key.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("cannot encode account key to cadence value: %w", err)
	}
	hash, err := newHashAlgoValue(key.HashAlgo)
	if err != nil {
		return nil, fmt.Errorf("cannot encode account key to cadence value: %w", err)
	}

	return cadence.NewStruct([]cadence.Value{
		cadence.NewInt(key.Weight),
		publicKey,
		hash,
		weight,
		cadence.NewBool(false),
	}).WithType(&cadence.StructType{
		Location:            common.IdentifierLocation("Crypto"),
		QualifiedIdentifier: "Crypto.KeyListEntry",
		Fields: []cadence.Field{{
			Identifier: "keyIndex",
			Type:       cadence.IntType{},
		}, {
			Identifier: "publicKey",
			Type:       exportType(sema.PublicKeyType).(*cadence.StructType),
		}, {
			Identifier: "hashAlgorithm",
			Type:       exportType(sema.HashAlgorithmType).(*cadence.EnumType),
		}, {
			Identifier: "weight",
			Type:       cadence.UFix64Type{},
		}, {
			Identifier: "isRevoked",
			Type:       cadence.BoolType{},
		}},
	}), nil
}

// CreateAccount generates a transactions that creates a new account.
//
// This template accepts a list of public keys and a contracts argument, both of which are optional.
//
// The contracts argument is a dictionary of *contract name*: *contract code (in bytes)*.
// All of the contracts will be deployed to the account.
//
// The final argument is the address of the account that will pay the account creation fee.
// This account is added as a transaction authorizer and therefore must sign the resulting transaction.
func CreateAccount(accountKeys []*flow.AccountKey, contracts []Contract, payer flow.Address) (*flow.Transaction, error) {
	keyList := make([]cadence.Value, len(accountKeys))

	contractKeyPairs := make([]cadence.KeyValuePair, len(contracts))

	var err error
	for i, key := range accountKeys {
		keyList[i], err = AccountKeyToCadenceCryptoKey(key)
		if err != nil {
			return nil, fmt.Errorf("cannot create CreateAccount transaction: %w", err)
		}
	}

	for i, contract := range contracts {
		contractKeyPairs[i] = cadence.KeyValuePair{
			Key:   cadence.String(contract.Name),
			Value: cadence.String(contract.SourceHex()),
		}
	}

	cadencePublicKeys := cadence.NewArray(keyList)
	cadenceContracts := cadence.NewDictionary(contractKeyPairs)

	return flow.NewTransaction().
		SetScript([]byte(templates.CreateAccount)).
		AddAuthorizer(payer).
		AddRawArgument(jsoncdc.MustEncode(cadencePublicKeys)).
		AddRawArgument(jsoncdc.MustEncode(cadenceContracts)), nil
}

// UpdateAccountContract generates a transaction that updates a contract deployed at an account.
func UpdateAccountContract(address flow.Address, contract Contract) *flow.Transaction {
	cadenceName := cadence.String(contract.Name)
	cadenceCode := cadence.String(contract.SourceHex())

	return flow.NewTransaction().
		SetScript([]byte(templates.UpdateContract)).
		AddRawArgument(jsoncdc.MustEncode(cadenceName)).
		AddRawArgument(jsoncdc.MustEncode(cadenceCode)).
		AddAuthorizer(address)
}

// AddAccountContract generates a transaction that deploys a contract to an account.
func AddAccountContract(address flow.Address, contract Contract) *flow.Transaction {
	cadenceName := cadence.String(contract.Name)
	cadenceCode := cadence.String(contract.SourceHex())

	return flow.NewTransaction().
		SetScript([]byte(templates.AddContract)).
		AddRawArgument(jsoncdc.MustEncode(cadenceName)).
		AddRawArgument(jsoncdc.MustEncode(cadenceCode)).
		AddAuthorizer(address)
}

// AddAccountKey generates a transaction that adds a public key to an account.
func AddAccountKey(address flow.Address, accountKey *flow.AccountKey) (*flow.Transaction, error) {
	key, err := AccountKeyToCadenceCryptoKey(accountKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create CreateAccount transaction: %w", err)
	}

	return flow.NewTransaction().
		SetScript([]byte(templates.AddAccountKey)).
		AddRawArgument(jsoncdc.MustEncode(key)).
		AddAuthorizer(address), nil
}

// RemoveAccountKey generates a transaction that removes a key from an account.
func RemoveAccountKey(address flow.Address, keyIndex int) *flow.Transaction {
	cadenceKeyIndex := cadence.NewInt(keyIndex)

	return flow.NewTransaction().
		SetScript([]byte(templates.RemoveAccountKey)).
		AddRawArgument(jsoncdc.MustEncode(cadenceKeyIndex)).
		AddAuthorizer(address)
}

// RemoveAccountContract generates a transaction that removes a contract with the given name
func RemoveAccountContract(address flow.Address, contractName string) *flow.Transaction {
	cadenceName := cadence.String(contractName)

	return flow.NewTransaction().
		SetScript([]byte(templates.RemoveContract)).
		AddRawArgument(jsoncdc.MustEncode(cadenceName)).
		AddAuthorizer(address)
}

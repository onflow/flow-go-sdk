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

func newSignAlgoValue(sigAlgo crypto.SignatureAlgorithm) cadence.Enum {
	sigAlgoCadence := sema.SignatureAlgorithmECDSA_P256
	if sigAlgo == crypto.ECDSA_secp256k1 {
		sigAlgoCadence = sema.SignatureAlgorithmECDSA_secp256k1
	}

	return cadence.NewEnum([]cadence.Value{
		cadence.NewUInt8(sigAlgoCadence.RawValue()),
	}).WithType(
		exportType(sema.SignatureAlgorithmType).(*cadence.EnumType),
	)
}

func newHashAlgoValue(hashAlgo crypto.HashAlgorithm) cadence.Enum {
	hashAlgoCadence := sema.HashAlgorithmSHA3_256
	if hashAlgo == crypto.SHA2_256 {
		hashAlgoCadence = sema.HashAlgorithmSHA2_256
	}

	return cadence.NewEnum([]cadence.Value{
		cadence.NewUInt8(hashAlgoCadence.RawValue()),
	}).WithType(
		exportType(sema.HashAlgorithmType).(*cadence.EnumType),
	)
}

func newPublicKeyValue(pubKey crypto.PublicKey) cadence.Struct {
	pubKeyCadence := make([]cadence.Value, len(pubKey.Encode()))
	for i, k := range pubKey.Encode() {
		pubKeyCadence[i] = cadence.NewUInt8(k)
	}

	return cadence.NewStruct(
		[]cadence.Value{
			cadence.NewArray(pubKeyCadence),
			newSignAlgoValue(pubKey.Algorithm()),
			cadence.NewBool(true),
		},
	).WithType(
		exportType(sema.PublicKeyType).(*cadence.StructType),
	)
}

func newAccountKeyValue(key *flow.AccountKey) cadence.Struct {
	weight, _ := cadence.NewUFix64(fmt.Sprintf("%d", key.Weight)) // ignore err as it shouldn't fail due to validation in acc key
	return cadence.Struct{
		StructType: exportType(sema.AccountKeyType).(*cadence.StructType),
		Fields: []cadence.Value{
			cadence.NewInt(key.Index),
			newPublicKeyValue(key.PublicKey),
			newHashAlgoValue(key.HashAlgo),
			weight,
			cadence.NewBool(key.Revoked),
		},
	}
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
func CreateAccount(accountKeys []*flow.AccountKey, contracts []Contract, payer flow.Address) *flow.Transaction {
	publicKeys := make([]cadence.Value, len(accountKeys))

	for i, accountKey := range accountKeys {
		publicKeys[i] = newAccountKeyValue(accountKey)
	}

	contractKeyPairs := make([]cadence.KeyValuePair, len(contracts))

	for i, contract := range contracts {
		contractKeyPairs[i] = cadence.KeyValuePair{
			Key:   cadence.String(contract.Name),
			Value: cadence.String(contract.SourceHex()),
		}
	}

	cadencePublicKeys := cadence.NewArray(publicKeys)
	cadenceContracts := cadence.NewDictionary(contractKeyPairs)

	return flow.NewTransaction().
		SetScript([]byte(templates.CreateAccount)).
		AddAuthorizer(payer).
		AddRawArgument(jsoncdc.MustEncode(cadencePublicKeys)).
		AddRawArgument(jsoncdc.MustEncode(cadenceContracts))
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
func AddAccountKey(address flow.Address, accountKey *flow.AccountKey) *flow.Transaction {
	keyHex := hex.EncodeToString(accountKey.PublicKey.Encode())
	cadencePublicKey := cadence.String(keyHex)

	return flow.NewTransaction().
		SetScript([]byte(templates.AddAccountKey)).
		AddRawArgument(jsoncdc.MustEncode(cadencePublicKey)).
		AddAuthorizer(address)
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

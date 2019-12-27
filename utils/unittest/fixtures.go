package unittest

import (
	"encoding/hex"

	"github.com/dapperlabs/flow-go/crypto"

	"github.com/dapperlabs/flow-go-sdk"
	"github.com/dapperlabs/flow-go-sdk/keys"
)

const PublicKeyFixtureCount = 2

func PublicKeyFixtures() [PublicKeyFixtureCount]crypto.PublicKey {
	encodedKeys := [PublicKeyFixtureCount]string{
		"3059301306072a8648ce3d020106082a8648ce3d0301070342000472b074a452d0a764a1da34318f44cb16740df1cfab1e6b50e5e4145dc06e5d151c9c25244f123e53c9b6fe237504a37e7779900aad53ca26e3b57c5c3d7030c4",
		"3059301306072a8648ce3d020106082a8648ce3d03010703420004d4423e4ca70ed9fb9bb9ce771e9393e0c3a1b66f3019ed89ab410cdf8f73d5a8ca06cc093766c1a46069cf83fce2a294d3322d55bb86ac9cb5aa805c7dd8d715",
	}

	keys := [PublicKeyFixtureCount]crypto.PublicKey{}

	for i, hexKey := range encodedKeys {
		bytesKey, _ := hex.DecodeString(hexKey)
		publicKey, _ := crypto.DecodePublicKey(crypto.ECDSA_P256, bytesKey)
		keys[i] = publicKey
	}

	return keys
}

func AddressFixture() flow.Address {
	return flow.RootAddress
}

func AccountSignatureFixture() flow.AccountSignature {
	return flow.AccountSignature{
		Account:   AddressFixture(),
		Signature: []byte{1, 2, 3, 4},
	}
}

func BlockHeaderFixture() flow.Header {
	return flow.Header{
		Parent: crypto.Hash("parent"),
		Number: 100,
	}
}

func TransactionFixture(n ...func(t *flow.Transaction)) flow.Transaction {
	tx := flow.Transaction{
		Script:             []byte("pub fun main() {}"),
		ReferenceBlockHash: HashFixture(32),
		Nonce:              1,
		ComputeLimit:       10,
		PayerAccount:       AddressFixture(),
		ScriptAccounts:     []flow.Address{AddressFixture()},
		Signatures:         []flow.AccountSignature{AccountSignatureFixture()},
	}
	if len(n) > 0 {
		n[0](&tx)
	}
	return tx
}

func AccountFixture() flow.Account {
	return flow.Account{
		Address: AddressFixture(),
		Balance: 10,
		Code:    []byte("pub fun main() {}"),
		Keys:    []flow.AccountPublicKey{AccountPublicKeyFixture()},
	}
}

func AccountPublicKeyFixture() flow.AccountPublicKey {
	return flow.AccountPublicKey{
		PublicKey: PublicKeyFixtures()[0],
		SignAlgo:  crypto.ECDSA_P256,
		HashAlgo:  crypto.SHA3_256,
		Weight:    keys.PublicKeyWeightThreshold,
	}
}

func EventFixture(n ...func(e *flow.Event)) flow.Event {
	event := flow.Event{
		Type: "Transfer",
		// TODO: create proper fixture
		// Values: map[string]interface{}{
		// 	"to":   flow.ZeroAddress,
		// 	"from": flow.ZeroAddress,
		// 	"id":   1,
		// },
		Payload: []byte{},
	}

	if len(n) >= 1 {
		n[0](&event)
	}

	return event
}

func HashFixture(size int) crypto.Hash {
	hash := make(crypto.Hash, size)
	for i := 0; i < size; i++ {
		hash[i] = byte(i)
	}
	return hash
}


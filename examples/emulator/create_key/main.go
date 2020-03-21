package main

import (
	"fmt"

	"github.com/dapperlabs/flow-go-sdk/keys"

	"github.com/dapperlabs/flow-developer-demo/examples/utils"
)

func main() {
	CreateKeyDemo()
}

func CreateKeyDemo() {
	// Create a key with P256 curve and SHA3 hashing
	seed := []byte("alabama shark baseball emblem computer caterpillar")
	fmt.Println("Seed: ", string(seed))
	key1, err := keys.GeneratePrivateKey(keys.ECDSA_P256_SHA3_256, seed)
	utils.Handle(err)

	encodedKey1, err := keys.EncodePrivateKey(key1)
	utils.Handle(err)

	fmt.Printf("P256 SHA3 key: %x\n", encodedKey1)

	// Create a key with SECP256K1 curve and SHA2 hashing
	key2, err := keys.GeneratePrivateKey(keys.ECDSA_SECp256k1_SHA2_256, seed)
	utils.Handle(err)

	encodedKey2, err := keys.EncodePrivateKey(key2)
	utils.Handle(err)

	fmt.Printf("SECP256K1 SHA2 key: %x\n", encodedKey2)
}

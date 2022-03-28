package test

import (
	"encoding/base64"
	"fmt"

	"github.com/onflow/flow-go/engine/access/rest/models"
)

func ContractHTTP() (string, string) {
	return "HelloWorld", base64.StdEncoding.EncodeToString([]byte(`
		contract HelloWorld {}
	`))
}

func AccountHTTP() models.Account {
	name, source := ContractHTTP()
	return models.Account{
		Address: AddressGenerator().New().String(),
		Balance: "10",
		Keys: []models.AccountPublicKey{
			AccountKeyHTTP(),
		},
		Contracts:  map[string]string{name: source},
		Expandable: nil,
		Links:      nil,
	}
}

func AccountKeyHTTP() models.AccountPublicKey {
	key := AccountKeyGenerator().New()
	sigAlgo := models.SigningAlgorithm(key.SigAlgo.String())
	hashAlgo := models.HashingAlgorithm(key.HashAlgo.String())

	return models.AccountPublicKey{
		Index:            "0",
		PublicKey:        key.PublicKey.String(),
		SigningAlgorithm: &sigAlgo,
		HashingAlgorithm: &hashAlgo,
		SequenceNumber:   "0",
		Weight:           "1000",
		Revoked:          false,
	}
}

func BlockHTTP() models.Block {
	block := BlockGenerator().New()

	return models.Block{
		Header: models.BlockHeader{
			Id:                   block.ID.String(),
			ParentId:             block.ParentID.String(),
			Height:               fmt.Sprintf("%d", block.Height),
			Timestamp:            block.Timestamp,
			ParentVoterSignature: base64.StdEncoding.EncodeToString(block.Signatures),
		},
		Payload: models.BlockPayload{
			CollectionGuarantees: nil,
			BlockSeals:           nil,
		},
		ExecutionResult: nil,
	}
}

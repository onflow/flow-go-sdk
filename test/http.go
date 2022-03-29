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
		Header: &models.BlockHeader{
			Id:                   block.ID.String(),
			ParentId:             block.ParentID.String(),
			Height:               fmt.Sprintf("%d", block.Height),
			Timestamp:            block.Timestamp,
			ParentVoterSignature: base64.StdEncoding.EncodeToString([]byte("test")),
		},
		Payload: &models.BlockPayload{
			CollectionGuarantees: models.CollectionGuarantees{{
				CollectionId: block.CollectionGuarantees[0].CollectionID.String(),
			}},
			BlockSeals: models.BlockSeals{{
				BlockId:                      block.Seals[0].BlockID.String(),
				ResultId:                     block.Seals[0].ExecutionReceiptID.String(),
				FinalState:                   "",
				AggregatedApprovalSignatures: nil,
			}},
		},
		ExecutionResult: nil,
	}
}

func CollectionHTTP() models.Collection {
	collection := CollectionGenerator().New()

	return models.Collection{
		Id: collection.ID().String(),
		Transactions: models.Transactions{
			models.Transaction{
				Id: collection.TransactionIDs[0].String(),
			},
		},
	}
}

func TransactionHTTP() models.Transaction {
	tx := TransactionGenerator().New()

	args := make([]string, len(tx.Arguments))
	for i, a := range tx.Arguments {
		args[i] = base64.StdEncoding.EncodeToString(a)
	}

	auths := make([]string, len(tx.Authorizers))
	for i, a := range tx.Authorizers {
		auths[i] = a.String()
	}

	return models.Transaction{
		Id:               tx.ID().String(),
		Script:           base64.StdEncoding.EncodeToString(tx.Script),
		Arguments:        args,
		ReferenceBlockId: tx.ReferenceBlockID.String(),
		GasLimit:         fmt.Sprintf("%d", tx.GasLimit),
		Payer:            tx.Payer.String(),
		ProposalKey: &models.ProposalKey{
			Address:        tx.ProposalKey.Address.String(),
			KeyIndex:       fmt.Sprintf("%d", tx.ProposalKey.KeyIndex),
			SequenceNumber: fmt.Sprintf("%d", tx.ProposalKey.SequenceNumber),
		},
		Authorizers: auths,
		PayloadSignatures: models.TransactionSignatures{
			models.TransactionSignature{
				Address:   tx.PayloadSignatures[0].Address.String(),
				KeyIndex:  fmt.Sprintf("%d", tx.PayloadSignatures[0].KeyIndex),
				Signature: base64.StdEncoding.EncodeToString(tx.PayloadSignatures[0].Signature),
			},
		},
		EnvelopeSignatures: models.TransactionSignatures{
			models.TransactionSignature{
				Address:   tx.EnvelopeSignatures[0].Address.String(),
				KeyIndex:  fmt.Sprintf("%d", tx.EnvelopeSignatures[0].KeyIndex),
				Signature: base64.StdEncoding.EncodeToString(tx.EnvelopeSignatures[0].Signature),
			},
		},
	}
}

func TransactionResultHTTP() models.TransactionResult {
	txr := TransactionResultGenerator().New()
	status := models.SEALED

	return models.TransactionResult{
		Status:       &status,
		StatusCode:   0,
		ErrorMessage: txr.Error.Error(),
		Events: models.Events{
			models.Event{
				Type_:            txr.Events[0].Type,
				TransactionId:    txr.Events[0].TransactionID.String(),
				TransactionIndex: fmt.Sprintf("%d", txr.Events[0].TransactionIndex),
				EventIndex:       fmt.Sprintf("%d", txr.Events[0].EventIndex),
				Payload:          base64.StdEncoding.EncodeToString(txr.Events[0].Payload),
			},
		},
	}
}

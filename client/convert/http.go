package convert

import (
	"strconv"

	"github.com/onflow/flow-go-sdk/crypto"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go/engine/access/rest/models"
)

func HTTPToAddress(address string) flow.Address {
	return flow.HexToAddress(address)
}

func HTTPToKeys(keys models.AccountPublicKeys) []*flow.AccountKey {
	accountKeys := make([]*flow.AccountKey, len(keys))

	for i, key := range keys {
		index, _ := strconv.Atoi(key.Index)
		weight, _ := strconv.Atoi(key.Weight)
		seqNumber, _ := strconv.ParseUint(key.SequenceNumber, 10, 64)

		accountKeys[i] = &flow.AccountKey{
			Index:          index,
			PublicKey:      nil,
			SigAlgo:        crypto.StringToSignatureAlgorithm(string(*key.SigningAlgorithm)),
			HashAlgo:       crypto.StringToHashAlgorithm(string(*key.HashingAlgorithm)),
			Weight:         weight,
			SequenceNumber: seqNumber,
			Revoked:        key.Revoked,
		}
	}

	return accountKeys
}

func HTTPToAccount(account *models.Account) *flow.Account {
	balance, _ := strconv.ParseUint(account.Balance, 10, 64)

	return &flow.Account{
		Address:   HTTPToAddress(account.Address),
		Balance:   balance,
		Code:      nil,
		Keys:      HTTPToKeys(account.Keys),
		Contracts: nil,
	}
}

func HTTPToBlockHeader(header *models.BlockHeader) *flow.BlockHeader {
	height, _ := strconv.ParseUint(header.Height, 10, 64)

	return &flow.BlockHeader{
		ID:        flow.HexToID(header.Id),
		ParentID:  flow.HexToID(header.ParentId),
		Height:    height,
		Timestamp: header.Timestamp,
	}
}

func HTTPToCollectionGuarantees(guarantees []models.CollectionGuarantee) []*flow.CollectionGuarantee {
	flowGuarantees := make([]*flow.CollectionGuarantee, len(guarantees))

	for i, guarantee := range guarantees {
		flowGuarantees[i] = &flow.CollectionGuarantee{
			flow.HexToID(guarantee.CollectionId),
		}
	}

	return flowGuarantees
}

func HTTPToBlockSeals(seals []models.BlockSeal) []*flow.BlockSeal {
	flowSeal := make([]*flow.BlockSeal, len(seals))

	for i, seal := range seals {
		flowSeal[i] = &flow.BlockSeal{
			BlockID:                    flow.HexToID(seal.BlockId),
			ExecutionReceiptID:         flow.Identifier{}, // todo: assign values
			ExecutionReceiptSignatures: nil,
			ResultApprovalSignatures:   nil,
		}
	}

	return flowSeal
}

func HTTPToBlockPayload(payload *models.BlockPayload) flow.BlockPayload {
	return flow.BlockPayload{
		CollectionGuarantees: HTTPToCollectionGuarantees(payload.CollectionGuarantees),
		Seals:                HTTPToBlockSeals(payload.BlockSeals),
	}
}

func HTTPToBlock(block *models.Block) *flow.Block {
	return &flow.Block{
		BlockHeader:  *HTTPToBlockHeader(block.Header),
		BlockPayload: HTTPToBlockPayload(block.Payload),
		Signatures:   nil, // todo: assign value
	}
}

func SealedToHTTP(isSealed bool) string {
	if isSealed {
		return "sealed"
	}
	return "final"
}

func HTTPToCollection(collection *models.Collection) *flow.Collection {
	IDs := make([]flow.Identifier, len(collection.Transactions))
	for i, tx := range collection.Transactions {
		IDs[i] = flow.HexToID(tx.Id)
	}
	return &flow.Collection{TransactionIDs: IDs}
}

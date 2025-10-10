/*
 * Flow Go SDK
 *
 * Copyright Flow Foundation
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

package convert

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/onflow/cadence"
	cadenceJSON "github.com/onflow/cadence/encoding/json"
	"github.com/pkg/errors"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access/http/models"
	"github.com/onflow/flow-go-sdk/crypto"
)

func ToAddress(address string) flow.Address {
	return flow.HexToAddress(address)
}

func ToKeys(keys []models.AccountPublicKey) []*flow.AccountKey {
	accountKeys := make([]*flow.AccountKey, len(keys))

	for i, key := range keys {
		sigAlgo := crypto.StringToSignatureAlgorithm(string(*key.SigningAlgorithm))
		pkey, _ := crypto.DecodePublicKeyHex(sigAlgo, strings.TrimPrefix(key.PublicKey, "0x")) // validation is done on AN

		accountKeys[i] = &flow.AccountKey{
			Index:          MustToUint32(key.Index),
			PublicKey:      pkey,
			SigAlgo:        sigAlgo,
			HashAlgo:       crypto.StringToHashAlgorithm(string(*key.HashingAlgorithm)),
			Weight:         MustToInt(key.Weight),
			SequenceNumber: MustToUint(key.SequenceNumber),
			Revoked:        key.Revoked,
		}
	}

	return accountKeys
}

func ToContracts(contracts map[string]string) (map[string][]byte, error) {
	decoded := make(map[string][]byte, len(contracts))
	for name, code := range contracts {
		dec, err := base64.StdEncoding.DecodeString(code)
		if err != nil {
			return nil, err
		}

		decoded[name] = dec
	}

	return decoded, nil
}

func ToAccount(account *models.Account) (*flow.Account, error) {
	contracts, err := ToContracts(account.Contracts)
	if err != nil {
		return nil, err
	}

	return &flow.Account{
		Address:   ToAddress(account.Address),
		Balance:   MustToUint(account.Balance),
		Keys:      ToKeys(account.Keys),
		Contracts: contracts,
	}, nil
}

func ToBlockHeader(header *models.BlockHeader, blockStatus string) *flow.BlockHeader {
	return &flow.BlockHeader{
		ID:        flow.HexToID(header.Id),
		ParentID:  flow.HexToID(header.ParentId),
		Height:    MustToUint(header.Height),
		Timestamp: header.Timestamp,
		Status:    flow.BlockStatusFromString(blockStatus),
	}
}

func ToCollectionGuarantees(guarantees []models.CollectionGuarantee) []*flow.CollectionGuarantee {
	flowGuarantees := make([]*flow.CollectionGuarantee, len(guarantees))

	for i, guarantee := range guarantees {
		flowGuarantees[i] = &flow.CollectionGuarantee{
			CollectionID: flow.HexToID(guarantee.CollectionId),
		}
	}

	return flowGuarantees
}

func ToBlockSeals(seals []models.BlockSeal) ([]*flow.BlockSeal, error) {
	flowSeal := make([]*flow.BlockSeal, len(seals))

	for i, seal := range seals {
		signatures := make([][]byte, 0)
		for _, sig := range seal.AggregatedApprovalSignatures {
			for _, ver := range sig.VerifierSignatures {
				dec, err := base64.StdEncoding.DecodeString(ver)
				if err != nil {
					return nil, err
				}
				signatures = append(signatures, dec)
			}
		}

		flowSeal[i] = &flow.BlockSeal{
			BlockID: flow.HexToID(seal.BlockId),
			// TODO: this needs to be changed to resultID
			// https://github.com/onflow/flow-go/blob/3683183977f2ea769836d8a31997701b3dbced83/model/flow/seal.go#L42
			ExecutionReceiptID: flow.HexToID(seal.ResultId),
		}
	}

	return flowSeal, nil
}

func ToBlockPayload(payload *models.BlockPayload) (*flow.BlockPayload, error) {
	seals, err := ToBlockSeals(payload.BlockSeals)
	if err != nil {
		return nil, err
	}

	return &flow.BlockPayload{
		CollectionGuarantees: ToCollectionGuarantees(payload.CollectionGuarantees),
		Seals:                seals,
	}, nil
}

func ToBlocks(blocks []*models.Block) ([]*flow.Block, error) {
	convertedBlocks := make([]*flow.Block, len(blocks))
	for i, b := range blocks {
		converted, err := ToBlock(b)
		if err != nil {
			return nil, err
		}

		convertedBlocks[i] = converted
	}
	return convertedBlocks, nil
}

func ToBlock(block *models.Block) (*flow.Block, error) {
	payload, err := ToBlockPayload(block.Payload)
	if err != nil {
		return nil, err
	}

	return &flow.Block{
		BlockHeader:  *ToBlockHeader(block.Header, block.BlockStatus),
		BlockPayload: *payload,
	}, nil
}

func ToCollection(collection *models.Collection) *flow.Collection {
	IDs := make([]flow.Identifier, len(collection.Transactions))
	for i, tx := range collection.Transactions {
		IDs[i] = flow.HexToID(tx.Id)
	}
	return &flow.Collection{
		TransactionIDs: IDs,
	}
}

func EncodeScript(script []byte) string {
	return base64.StdEncoding.EncodeToString(script)
}

func ToScript(script string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(script)
}

func EncodeArgs(args [][]byte) []string {
	encodedArgs := make([]string, len(args))
	for i, a := range args {
		encodedArgs[i] = base64.StdEncoding.EncodeToString(a)
	}
	return encodedArgs
}

func ToArgs(arguments []string) ([][]byte, error) {
	args := make([][]byte, len(arguments))
	for i, arg := range arguments {
		a, err := base64.StdEncoding.DecodeString(arg)
		if err != nil {
			return nil, err
		}
		args[i] = a
	}

	return args, nil
}

func MustToUint(value string) uint64 {
	parsed, _ := strconv.ParseUint(value, 10, 64) // we can ignore error since these values are validated before returned
	return parsed
}

func MustToUint32(value string) uint32 {
	parsed, _ := strconv.ParseUint(value, 10, 32) // we can ignore error since these values are validated before returned
	return uint32(parsed)
}

func MustToInt(value string) int {
	parsed, _ := strconv.Atoi(value) // we can ignore error since these values are validated before returned
	return parsed
}

func EncodeCadenceArgs(args []cadence.Value) ([]string, error) {
	encArgs := make([]string, len(args))

	for i, a := range args {
		jsonArg, err := cadenceJSON.Encode(a)
		if err != nil {
			return nil, err
		}

		encArgs[i] = base64.StdEncoding.EncodeToString(jsonArg)
	}

	return encArgs, nil
}

func DecodeCadenceValue(value string, options []cadenceJSON.Option) (cadence.Value, error) {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}

	return cadenceJSON.Decode(nil, decoded, options...)
}

func ToProposalKey(key *models.ProposalKey) flow.ProposalKey {
	return flow.ProposalKey{
		Address:        flow.HexToAddress(key.Address),
		KeyIndex:       MustToUint32(key.KeyIndex),
		SequenceNumber: MustToUint(key.SequenceNumber),
	}
}

func ToSignatures(signatures []models.TransactionSignature) []flow.TransactionSignature {
	sigs := make([]flow.TransactionSignature, len(signatures))
	for i, sig := range signatures {
		signature, _ := base64.StdEncoding.DecodeString(sig.Signature) // signatures are validated and must be valid
		sigs[i] = flow.TransactionSignature{
			Address:   flow.HexToAddress(sig.Address),
			KeyIndex:  MustToUint32(sig.KeyIndex),
			Signature: signature,
		}
	}
	return sigs
}

func ToTransaction(tx *models.Transaction) (*flow.Transaction, error) {
	script, err := ToScript(tx.Script)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to decode script of transaction with ID %s", tx.Id))
	}
	args, err := ToArgs(tx.Arguments)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to decode arguments of transaction with ID %s", tx.Id))
	}

	auths := make([]flow.Address, len(tx.Authorizers))
	for i, a := range tx.Authorizers {
		auths[i] = flow.HexToAddress(a)
	}

	return &flow.Transaction{
		Script:             script,
		Arguments:          args,
		ReferenceBlockID:   flow.HexToID(tx.ReferenceBlockId),
		GasLimit:           MustToUint(tx.GasLimit),
		ProposalKey:        ToProposalKey(tx.ProposalKey),
		Payer:              flow.HexToAddress(tx.Payer),
		Authorizers:        auths,
		PayloadSignatures:  ToSignatures(tx.PayloadSignatures),
		EnvelopeSignatures: ToSignatures(tx.EnvelopeSignatures),
	}, nil
}

func ToTransactionStatus(status *models.TransactionStatus) flow.TransactionStatus {
	switch *status {
	case models.PENDING_TransactionStatus:
		return flow.TransactionStatusPending
	case models.SEALED_TransactionStatus:
		return flow.TransactionStatusSealed
	case models.FINALIZED_TransactionStatus:
		return flow.TransactionStatusFinalized
	case models.EXECUTED_TransactionStatus:
		return flow.TransactionStatusExecuted
	case models.EXPIRED_TransactionStatus:
		return flow.TransactionStatusExpired
	default:
		return flow.TransactionStatusUnknown
	}
}

func ToEvents(events []models.Event, options []cadenceJSON.Option) ([]flow.Event, error) {
	flowEvents := make([]flow.Event, len(events))
	for i, e := range events {
		payload, err := base64.StdEncoding.DecodeString(e.Payload)
		if err != nil {
			return nil, err
		}

		event, err := cadenceJSON.Decode(nil, payload, options...)
		if err != nil {
			return nil, err
		}

		flowEvents[i] = flow.Event{
			Type:             e.Type_,
			TransactionID:    flow.HexToID(e.TransactionId),
			TransactionIndex: MustToInt(e.TransactionIndex),
			EventIndex:       MustToInt(e.EventIndex),
			Value:            event.(cadence.Event),
			Payload:          payload,
		}
	}
	return flowEvents, nil
}

func ToBlockEvents(blockEvents []models.BlockEvents, options []cadenceJSON.Option) ([]flow.BlockEvents, error) {
	blocks := make([]flow.BlockEvents, len(blockEvents))
	for i, block := range blockEvents {
		events, err := ToEvents(block.Events, options)
		if err != nil {
			return nil, err
		}

		blocks[i] = flow.BlockEvents{
			BlockID:        flow.HexToID(block.BlockId),
			Height:         MustToUint(block.BlockHeight),
			BlockTimestamp: block.BlockTimestamp,
			Events:         events,
		}
	}
	return blocks, nil
}

func ToTransactionResult(txr *models.TransactionResult, options []cadenceJSON.Option) (*flow.TransactionResult, error) {
	events, err := ToEvents(txr.Events, options)
	if err != nil {
		return nil, err
	}

	var txErr error
	if txr.ErrorMessage != "" {
		txErr = fmt.Errorf("%s", txr.ErrorMessage)
	}

	return &flow.TransactionResult{
		Status:       ToTransactionStatus(txr.Status),
		Error:        txErr,
		Events:       events,
		BlockID:      flow.HexToID(txr.BlockId),
		CollectionID: flow.HexToID(txr.CollectionId),
	}, nil
}

func EncodeSignatures(signatures []flow.TransactionSignature) []models.TransactionSignature {
	sigs := make([]models.TransactionSignature, len(signatures))
	for i, sig := range signatures {
		sigs[i] = models.TransactionSignature{
			Address:   sig.Address.String(),
			KeyIndex:  fmt.Sprintf("%d", sig.KeyIndex),
			Signature: base64.StdEncoding.EncodeToString(sig.Signature),
		}
	}

	return sigs
}

func TncodeTransaction(tx flow.Transaction) ([]byte, error) {
	auths := make([]string, len(tx.Authorizers))
	for i, address := range tx.Authorizers {
		auths[i] = address.String()
	}

	return json.Marshal(models.TransactionsBody{
		Script:           EncodeScript(tx.Script),
		Arguments:        EncodeArgs(tx.Arguments),
		ReferenceBlockId: tx.ReferenceBlockID.String(),
		GasLimit:         fmt.Sprintf("%d", tx.GasLimit),
		Payer:            tx.Payer.String(),
		ProposalKey: &models.ProposalKey{
			Address:        tx.ProposalKey.Address.String(),
			KeyIndex:       fmt.Sprintf("%d", tx.ProposalKey.KeyIndex),
			SequenceNumber: fmt.Sprintf("%d", tx.ProposalKey.SequenceNumber),
		},
		Authorizers:        auths,
		PayloadSignatures:  EncodeSignatures(tx.PayloadSignatures),
		EnvelopeSignatures: EncodeSignatures(tx.EnvelopeSignatures),
	})
}

func ToExecutionResults(result models.ExecutionResult) *flow.ExecutionResult {
	events := make([]*flow.ServiceEvent, len(result.Events))
	for i, e := range result.Events {
		events[i] = &flow.ServiceEvent{
			Type:    e.Type_,
			Payload: []byte(e.Payload),
		}
	}

	chunks := make([]*flow.Chunk, len(result.Chunks))

	for i, chunk := range result.Chunks {
		chunks[i] = &flow.Chunk{
			CollectionIndex:      uint(MustToUint(chunk.CollectionIndex)),
			StartState:           flow.HexToStateCommitment(chunk.StartState),
			EventCollection:      crypto.Hash(chunk.EventCollection),
			BlockID:              flow.HexToID(chunk.BlockId),
			TotalComputationUsed: MustToUint(chunk.TotalComputationUsed),
			NumberOfTransactions: uint16(MustToUint(chunk.NumberOfTransactions)),
			Index:                MustToUint(chunk.Index),
			EndState:             flow.HexToStateCommitment(chunk.EndState),
		}
	}

	return &flow.ExecutionResult{
		PreviousResultID: flow.HexToID(result.PreviousResultId),
		BlockID:          flow.HexToID(result.BlockId),
		Chunks:           chunks,
		ServiceEvents:    events,
	}
}

func ToNetworkParameters(params *models.NetworkParameters) *flow.NetworkParameters {
	return &flow.NetworkParameters{
		ChainID: flow.ChainID(params.ChainId),
	}
}

func ToNodeVersionInfo(info *models.NodeVersionInfo) (*flow.NodeVersionInfo, error) {
	version, err := strconv.ParseUint(info.ProtocolVersion, 10, 64)
	if err != nil {
		return nil, err
	}
	sporkHeight, err := strconv.ParseUint(info.SporkRootBlockHeight, 10, 64)
	if err != nil {
		return nil, err
	}
	nodeHeight, err := strconv.ParseUint(info.NodeRootBlockHeight, 10, 64)
	if err != nil {
		return nil, err
	}

	return &flow.NodeVersionInfo{
		Semver:               info.Semver,
		Commit:               info.Commit,
		SporkId:              flow.HexToID(info.SporkId),
		ProtocolVersion:      version,
		SporkRootBlockHeight: sporkHeight,
		NodeRootBlockHeight:  nodeHeight,
	}, nil
}

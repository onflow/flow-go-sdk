package convert

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/onflow/cadence"
	cadenceJSON "github.com/onflow/cadence/encoding/json"

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
		sigAlgo := crypto.StringToSignatureAlgorithm(string(*key.SigningAlgorithm))
		pkey, _ := crypto.DecodePublicKeyHex(sigAlgo, strings.TrimPrefix(key.PublicKey, "0x")) // validation is done on AN

		accountKeys[i] = &flow.AccountKey{
			Index:          MustHTTPToInt(key.Index),
			PublicKey:      pkey,
			SigAlgo:        sigAlgo,
			HashAlgo:       crypto.StringToHashAlgorithm(string(*key.HashingAlgorithm)),
			Weight:         MustHTTPToInt(key.Weight),
			SequenceNumber: MustHTTPToUint(key.SequenceNumber),
			Revoked:        key.Revoked,
		}
	}

	return accountKeys
}

func HTTPToContracts(contracts map[string]string) (map[string][]byte, error) {
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

func HTTPToAccount(account *models.Account) (*flow.Account, error) {
	contracts, err := HTTPToContracts(account.Contracts)
	if err != nil {
		return nil, err
	}

	return &flow.Account{
		Address:   HTTPToAddress(account.Address),
		Balance:   MustHTTPToUint(account.Balance),
		Keys:      HTTPToKeys(account.Keys),
		Contracts: contracts,
	}, nil
}

func HTTPToBlockHeader(header *models.BlockHeader) *flow.BlockHeader {
	return &flow.BlockHeader{
		ID:        flow.HexToID(header.Id),
		ParentID:  flow.HexToID(header.ParentId),
		Height:    MustHTTPToUint(header.Height),
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

func HTTPToBlockSeals(seals []models.BlockSeal) ([]*flow.BlockSeal, error) {
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
			BlockID:                    flow.HexToID(seal.BlockId),
			ExecutionReceiptID:         flow.HexToID(seal.ResultId), // todo this needs to be changed to resultID https://github.com/onflow/flow-go/blob/3683183977f2ea769836d8a31997701b3dbced83/model/flow/seal.go#L42
			ExecutionReceiptSignatures: nil,                         // todo this is deprecated, should be removed
			ResultApprovalSignatures:   signatures,
		}
	}

	return flowSeal, nil
}

func HTTPToBlockPayload(payload *models.BlockPayload) (*flow.BlockPayload, error) {
	seals, err := HTTPToBlockSeals(payload.BlockSeals)
	if err != nil {
		return nil, err
	}

	return &flow.BlockPayload{
		CollectionGuarantees: HTTPToCollectionGuarantees(payload.CollectionGuarantees),
		Seals:                seals,
	}, nil
}

func HTTPToBlocks(blocks []*models.Block) ([]*flow.Block, error) {
	convertedBlocks := make([]*flow.Block, len(blocks))
	for i, b := range blocks {
		converted, err := HTTPToBlock(b)
		if err != nil {
			return nil, err
		}

		convertedBlocks[i] = converted
	}
	return convertedBlocks, nil
}

func HTTPToBlock(block *models.Block) (*flow.Block, error) {
	payload, err := HTTPToBlockPayload(block.Payload)
	if err != nil {
		return nil, err
	}

	signature, err := base64.StdEncoding.DecodeString(block.Header.ParentVoterSignature)
	if err != nil {
		return nil, err
	}

	return &flow.Block{
		BlockHeader:  *HTTPToBlockHeader(block.Header),
		BlockPayload: *payload,
		Signatures:   [][]byte{signature},
	}, nil
}

func HTTPToCollection(collection *models.Collection) *flow.Collection {
	IDs := make([]flow.Identifier, len(collection.Transactions))
	for i, tx := range collection.Transactions {
		IDs[i] = flow.HexToID(tx.Id)
	}
	return &flow.Collection{
		TransactionIDs: IDs,
	}
}

func ScriptToHTTP(script []byte) string {
	return base64.StdEncoding.EncodeToString(script)
}

func HTTPToScript(script string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(script)
}

func ArgumentsToHTTP(args [][]byte) []string {
	encodedArgs := make([]string, len(args))
	for i, a := range args {
		encodedArgs[i] = base64.StdEncoding.EncodeToString(a)
	}
	return encodedArgs
}

func HTTPToArguments(arguments []string) ([][]byte, error) {
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

func MustHTTPToUint(value string) uint64 {
	parsed, _ := strconv.ParseUint(value, 10, 64) // we can ignore error since this values are validated before returned
	return parsed
}

func MustHTTPToInt(value string) int {
	parsed, _ := strconv.Atoi(value) // we can ignore error since this values are validated before returned
	return parsed
}

func CadenceArgsToHTTP(args []cadence.Value) ([]string, error) {
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

func HTTPToCadenceValue(value string) (cadence.Value, error) {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}

	return cadenceJSON.Decode(decoded)
}

func HTTPToProposalKey(key *models.ProposalKey) flow.ProposalKey {
	return flow.ProposalKey{
		Address:        flow.HexToAddress(key.Address),
		KeyIndex:       MustHTTPToInt(key.KeyIndex),
		SequenceNumber: MustHTTPToUint(key.SequenceNumber),
	}
}

func HTTPToSignatures(signatures models.TransactionSignatures) []flow.TransactionSignature {
	sigs := make([]flow.TransactionSignature, len(signatures))
	for i, sig := range signatures {
		signature, _ := base64.StdEncoding.DecodeString(sig.Signature) // signatures are validated and must be valid
		sigs[i] = flow.TransactionSignature{
			Address:   flow.HexToAddress(sig.Address),
			KeyIndex:  MustHTTPToInt(sig.KeyIndex),
			Signature: signature,
		}
	}
	return sigs
}

func HTTPToTransaction(tx *models.Transaction) (*flow.Transaction, error) {
	script, err := HTTPToScript(tx.Script)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to decode script of transaction with ID %s", tx.Id))
	}
	args, err := HTTPToArguments(tx.Arguments)
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
		GasLimit:           MustHTTPToUint(tx.GasLimit),
		ProposalKey:        HTTPToProposalKey(tx.ProposalKey),
		Payer:              flow.HexToAddress(tx.Payer),
		Authorizers:        auths,
		PayloadSignatures:  HTTPToSignatures(tx.PayloadSignatures),
		EnvelopeSignatures: HTTPToSignatures(tx.EnvelopeSignatures),
	}, nil
}

func HTTPToTransactionStatus(status *models.TransactionStatus) flow.TransactionStatus {
	switch *status {
	case models.PENDING:
		return flow.TransactionStatusPending
	case models.SEALED:
		return flow.TransactionStatusSealed
	case models.FINALIZED:
		return flow.TransactionStatusFinalized
	case models.EXECUTED:
		return flow.TransactionStatusExecuted
	case models.EXPIRED:
		return flow.TransactionStatusExpired
	default:
		return flow.TransactionStatusUnknown
	}
}

func HTTPToEvents(events []models.Event) ([]flow.Event, error) {
	flowEvents := make([]flow.Event, len(events))
	for i, e := range events {
		payload, err := base64.StdEncoding.DecodeString(e.Payload)
		if err != nil {
			return nil, err
		}

		event, err := cadenceJSON.Decode(payload)
		if err != nil {
			return nil, err
		}

		flowEvents[i] = flow.Event{
			Type:             e.Type_,
			TransactionID:    flow.HexToID(e.TransactionId),
			TransactionIndex: MustHTTPToInt(e.TransactionIndex),
			EventIndex:       MustHTTPToInt(e.EventIndex),
			Value:            event.(cadence.Event),
			Payload:          payload,
		}
	}
	return flowEvents, nil
}

func HTTPToBlockEvents(blockEvents []models.BlockEvents) ([]flow.BlockEvents, error) {
	blocks := make([]flow.BlockEvents, len(blockEvents))
	for i, block := range blockEvents {
		events, err := HTTPToEvents(block.Events)
		if err != nil {
			return nil, err
		}

		blocks[i] = flow.BlockEvents{
			BlockID:        flow.HexToID(block.BlockId),
			Height:         MustHTTPToUint(block.BlockHeight),
			BlockTimestamp: block.BlockTimestamp,
			Events:         events,
		}
	}
	return blocks, nil
}

func HTTPToTransactionResult(txr *models.TransactionResult) (*flow.TransactionResult, error) {
	events, err := HTTPToEvents(txr.Events)
	if err != nil {
		return nil, err
	}

	var txErr error
	if txr.ErrorMessage != "" {
		txErr = fmt.Errorf(txr.ErrorMessage)
	}

	return &flow.TransactionResult{
		Status: HTTPToTransactionStatus(txr.Status),
		Error:  txErr,
		Events: events,
	}, nil
}

func SignaturesToHTTP(signatures []flow.TransactionSignature) models.TransactionSignatures {
	sigs := make(models.TransactionSignatures, len(signatures))
	for i, sig := range signatures {
		sigs[i] = models.TransactionSignature{
			Address:   sig.Address.String(),
			KeyIndex:  fmt.Sprintf("%d", sig.KeyIndex),
			Signature: base64.StdEncoding.EncodeToString(sig.Signature),
		}
	}

	return sigs
}

func TransactionToHTTP(tx flow.Transaction) ([]byte, error) {
	auths := make([]string, len(tx.Authorizers))
	for i, address := range tx.Authorizers {
		auths[i] = address.String()
	}

	return json.Marshal(models.TransactionsBody{
		Script:           ScriptToHTTP(tx.Script),
		Arguments:        ArgumentsToHTTP(tx.Arguments),
		ReferenceBlockId: tx.ReferenceBlockID.String(),
		GasLimit:         fmt.Sprintf("%d", tx.GasLimit),
		Payer:            tx.Payer.String(),
		ProposalKey: &models.ProposalKey{
			Address:        tx.ProposalKey.Address.String(),
			KeyIndex:       fmt.Sprintf("%d", tx.ProposalKey.KeyIndex),
			SequenceNumber: fmt.Sprintf("%d", tx.ProposalKey.SequenceNumber),
		},
		Authorizers:        auths,
		PayloadSignatures:  SignaturesToHTTP(tx.PayloadSignatures),
		EnvelopeSignatures: SignaturesToHTTP(tx.EnvelopeSignatures),
	})
}

func HTTPToExecutionResults(result models.ExecutionResult) *flow.ExecutionResult {
	events := make([]*flow.ServiceEvent, len(result.Events))
	for i, e := range result.Events {
		events[i] = &flow.ServiceEvent{
			Type:    e.Type_,
			Payload: []byte(e.Payload),
		}
	}

	// todo there is missing data on the http api, make sure this is consistent
	return &flow.ExecutionResult{
		PreviousResultID: flow.EmptyID,
		BlockID:          flow.HexToID(result.BlockId),
		Chunks:           nil,
		ServiceEvents:    events,
	}
}

package transactions

type FlowTransactionsResp struct {
	Data Data `json:"data"`
}
type Fields struct {
	NftID                int    `json:"nftID"`
	Expiry               int64  `json:"expiry"`
	NftType              string `json:"nftType"`
	NftUUID              int    `json:"nftUUID"`
	CustomID             string `json:"customID"`
	Purchased            bool   `json:"purchased"`
	SalePrice            int    `json:"salePrice"`
	CommissionAmount     int    `json:"commissionAmount"`
	ListingResourceID    int    `json:"listingResourceID"`
	SalePaymentVaultType string `json:"salePaymentVaultType"`
	StorefrontResourceID int    `json:"storefrontResourceID"`
}
type Events struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	EventIndex int    `json:"event_index"`
	Typename   string `json:"__typename"`
	Fields     Fields `json:"fields"`
}
type ContractTransactions struct {
	ContractID string `json:"contract_id"`
	Typename   string `json:"__typename"`
}
type Transaction struct {
	ID                     string                 `json:"id"`
	Timestamp              string                 `json:"timestamp"`
	Payer                  string                 `json:"payer"`
	Authorizers            []string               `json:"authorizers"`
	GasUsed                int                    `json:"gas_used"`
	Fee                    float64                `json:"fee"`
	Status                 string                 `json:"status"`
	BlockHeight            int                    `json:"block_height"`
	Error                  string                 `json:"error"`
	TransactionBodyHash    string                 `json:"transaction_body_hash"`
	Typename               string                 `json:"__typename"`
	Events                 []Events               `json:"events"`
	ContractTransactions   []ContractTransactions `json:"contract_transactions"`
	Proposer               string                 `json:"proposer"`
	ProposerIndex          int                    `json:"proposer_index"`
	ProposerSequenceNumber int                    `json:"proposer_sequence_number"`
}
type Participations struct {
	Roles       []string    `json:"roles"`
	Transaction Transaction `json:"transaction"`
	Typename    string      `json:"__typename"`
}
type Data struct {
	Participations []Participations `json:"participations"`
}

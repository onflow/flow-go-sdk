package transaction_detail

type FlowTransactionDetailResp struct {
	Data Data `json:"data"`
}
type TransactionBody struct {
	Body     string `json:"body"`
	Typename string `json:"__typename"`
}
type Argument struct {
	Key   string      `json:"Key"`
	Value interface{} `json:"Value"`
}
type Fields struct {
	ID              int `json:"id"`
	Duration        int `json:"duration"`
	ExpiryTimestamp int `json:"expiryTimestamp"`
}
type Events struct {
	Name       string `json:"name"`
	Fields     Fields `json:"fields"`
	EventIndex int    `json:"event_index"`
	Typename   string `json:"__typename"`
}
type Contract struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Typename string `json:"__typename"`
}
type ContractTransactions struct {
	Contract Contract `json:"contract"`
	Typename string   `json:"__typename"`
}
type Transactions struct {
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
	BlockID                string                 `json:"block_id"`
	Proposer               string                 `json:"proposer"`
	ProposerIndex          int                    `json:"proposer_index"`
	ProposerSequenceNumber int                    `json:"proposer_sequence_number"`
	TransactionBody        TransactionBody        `json:"transaction_body"`
	ExecutionEffort        float64                `json:"execution_effort"`
	Argument               []Argument             `json:"argument"`
	Events                 []Events               `json:"events"`
	ContractTransactions   []ContractTransactions `json:"contract_transactions"`
}
type Data struct {
	Transactions []Transactions `json:"transactions"`
}

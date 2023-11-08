package subgraph

type GetSubGraphRequest struct {
	OperationName string      `json:"operationName"`
	Query         string      `json:"query"`
	Variables     interface{} `json:"variables"`
}

type TransactionDetailVariables struct {
	Id string `json:"id"`
}

func (f *subgraphClient) GetTransactionDetailReq(id string) *GetSubGraphRequest {
	return &GetSubGraphRequest{
		OperationName: "TransactionDetails",
		Query:         "query TransactionDetails($id: String) {\n  transactions(where: {id: {_eq: $id}}) {\n    ...BasicTransaction\n    ...ExtraDetails\n    __typename\n  }\n}\nfragment BasicTransaction on transactions {\n  id\n  timestamp\n  payer\n  authorizers\n  gas_used\n  fee\n  status\n  block_height\n  error\n  transaction_body_hash\n  __typename\n}\nfragment ExtraDetails on transactions {\n  block_id\n  proposer\n  proposer_index\n  proposer_sequence_number\n  transaction_body_hash\n  transaction_body {\n    body\n    __typename\n  }\n  execution_effort\n  argument\n  events {\n    name\n    fields\n    event_index\n    __typename\n  }\n  contract_transactions {\n    contract {\n      name\n      address\n      __typename\n    }\n    __typename\n  }\n  __typename\n}",
		Variables: TransactionDetailVariables{
			Id: id,
		},
	}
}

type TransactionListVariables struct {
	Address           string        `json:"address"`
	ActorFilter       []interface{} `json:"actorFilter" `
	AuthorizersFilter interface{}   `json:"authorizersFilter" `
	BlockHeightFilter interface{}   `json:"blockHeightFilter" `
	EventCountFilter  interface{}   `json:"eventCountFilter" `
	GasRangeFilter    interface{}   `json:"gasRangeFilter" `
	Limit             uint64        `json:"limit"`
	Offset            uint64        `json:"offset"`
	PayerFilter       interface{}   `json:"payerFilter"`
	ProposerFilter    interface{}   `json:"proposerFilter"`
	StatusFilter      interface{}   `json:"statusFilter"`
	TimeFilter        interface{}   `json:"timeFilter"`
	TypeFilter        interface{}   `json:"typeFilter"`
}

func (f *subgraphClient) GetTransactionsByAddressReq(txsVariables *TransactionListVariables) *GetSubGraphRequest {
	return &GetSubGraphRequest{
		OperationName: "AccountTransactions",
		Query:         "query AccountTransactions($limit: Int, $offset: Int, $address: String, $authorizersFilter: String_array_comparison_exp, $blockHeightFilter: bigint_comparison_exp, $eventCountFilter: Int_comparison_exp!, $gasRangeFilter: bigint_comparison_exp, $payerFilter: String_comparison_exp, $statusFilter: String_comparison_exp, $typeFilter: String_comparison_exp, $timeFilter: timestamptz_comparison_exp, $actorFilter: [transactions_bool_exp!], $proposerFilter: String_comparison_exp) @cached(ttl: 60) {\n  participations(\n    limit: $limit\n    offset: $offset\n    where: {address: {_eq: $address}, transaction: {_not: {_or: $actorFilter}, status: $statusFilter, block_height: $blockHeightFilter, transaction_body_hash: $typeFilter, authorizers: $authorizersFilter, payer: $payerFilter, proposer: $proposerFilter, gas_used: $gasRangeFilter, timestamp: $timeFilter, events_aggregate: {count: {predicate: $eventCountFilter}}}}\n    order_by: {timestamp: desc_nulls_first}\n  ) {\n    roles\n    transaction {\n      ...BasicTransaction\n      ...TransactionEvents\n      ...TransactionContracts\n      ...ProposerInfo\n      events {\n        event_index\n        fields\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n}\nfragment BasicTransaction on transactions {\n  id\n  timestamp\n  payer\n  authorizers\n  gas_used\n  fee\n  status\n  block_height\n  error\n  transaction_body_hash\n  __typename\n}\nfragment TransactionEvents on transactions {\n  events {\n    id\n    name\n    event_index\n    __typename\n  }\n  __typename\n}\nfragment TransactionContracts on transactions {\n  contract_transactions {\n    contract_id\n    __typename\n  }\n  __typename\n}\nfragment ProposerInfo on transactions {\n  proposer\n  proposer_index\n  proposer_sequence_number\n  __typename\n}",
		Variables:     txsVariables,
	}
}

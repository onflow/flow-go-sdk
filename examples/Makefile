.PHONY: all
all: get-blocks get-accounts get-events get-collection get-network-parameters get-transactions execute-script send-transactions create-account add-account-key deploy-contract storage-usage transaction-arguments single-party single-party-multisig multi-party multi-party-multisig user-signature user-signature-validate-all user-signature-validate-any http-grpc-clients modify-account get-execution-data

.PHONY: create-account
create-account:
	go run ./create_account/main.go

.PHONY: add-account-key
add-account-key:
	go run ./add_account_key/main.go

.PHONY: deploy-contract
deploy-contract:
	go run ./deploy_contract/main.go

.PHONY: verify-events
verify-events:
	go run ./verify_events/main.go

.PHONY: storage-usage
storage-usage:
	go run ./storage_usage/main.go

.PHONY: transaction-arguments
transaction-arguments:
	go run ./transaction_arguments/main.go

.PHONY: single-party
single-party:
	go run ./transaction_signing/single_party/main.go

.PHONY: single-party-multisig
single-party-multisig:
	go run ./transaction_signing/single_party_multisig/main.go

.PHONY: multi-party
multi-party:
	go run ./transaction_signing/multi_party/main.go

.PHONY: multi-party-two-authorizers
multi-party-two-authorizers:
	go run ./transaction_signing/multi_party_two_authorizers/main.go

.PHONY: multi-party-multisig
multi-party-multisig:
	go run ./transaction_signing/multi_party_multisig/main.go

.PHONY: user-signature
user-signature:
	go run ./verify_signature/user_signature/main.go

.PHONY: get-blocks
get-blocks:
	go run ./get_blocks/main.go

.PHONY: get-accounts
get-accounts:
	go run ./get_accounts/main.go

.PHONY: get-events
get-events:
	go run ./get_events/main.go

.PHONY: get-collection
get-collection:
	go run ./get_collection/main.go

.PHONY: get-network-parameters
get-network-parameters:
	go run ./get_network_parameters/main.go

.PHONY: get-transactions
get-transactions:
	go run ./get_transactions/main.go

.PHONY: execute-script
execute-script:
	go run ./execute_script/main.go

.PHONY: send-transactions
send-transactions:
	go run ./send_transactions/main.go

.PHONY: user-signature-validate-all
user-signature-validate-all:
	go run ./verify_signature/user_signature_validate_all/main.go

.PHONY: user-signature-validate-any
user-signature-validate-any:
	go run ./verify_signature/user_signature_validate_any/main.go

.PHONY: http-grpc-clients
http-grpc-clients:
	go run ./http_grpc_clients/main.go

.PHONY: modify-account
modify-account:
	go run ./modify_account/main.go

.PHONY: get-execution-data
get-execution-data:
	go run ./get_execution_data/main.go

.PHONY: stream-events
stream-events:
	go run ./stream_events/main.go

.PHONY: stream-events-reconnect
stream-events-reconnect:
	go run ./stream_events_reconnect/main.go

.PHONY: stream-execution-data
stream-execution-data:
	go run ./stream_execution_data/main.go
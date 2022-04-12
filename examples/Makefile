.PHONY: all
all: get-blocks get-accounts get-events get-collection get-transactions execute-script send-transaction create-account add-account-key deploy-contract query-events storage-usage transaction-arguments single-party single-party-multisig multi-party multi-party-multisig user-signature user-signature-validate-all user-signature-validate-any rest-grpc-clients

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
	go run ./user_signature/main.go

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

.PHONY: get-transactions
get-transactions:
	go run ./get_transactions/main.go

.PHONY: execute-script
execute-script:
	go run ./execute_script/main.go

.PHONY: send-transaction
send-transaction:
	go run ./send_transactions/main.go

.PHONY: user-signature-validate-all
user-signature-validate-all:
	go run ./verify_signature/user_signature_validate_all/main.go

.PHONY: user-signature-validate-any
user-signature-validate-any:
	go run ./verify_signature/user_signature_validate_any/main.go

.PHONY: rest-grpc-clients
rest-grpc-clients:
	go run ./rest_grpc_clients/main.go

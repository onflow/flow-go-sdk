.PHONY: all
all: create-account add-account-key deploy-contract query-events storage-usage transaction-arguments single-party single-party-multisig multi-party multi-party-multisig user-signature

.PHONY: create-account
create-account:
	go run ./create_account/main.go

.PHONY: add-account-key
add-account-key:
	go run ./add_account_key/main.go

.PHONY: deploy-contract
deploy-contract:
	go run ./deploy_contract/main.go

.PHONY: query-events
query-events:
	go run ./query_events/main.go

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

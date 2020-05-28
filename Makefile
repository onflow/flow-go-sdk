# Ensure go bin path is in path (especially for CI)
PATH := $(PATH):$(GOPATH)/bin

.PHONY: test
test:
	go test -coverprofile=cover.out ./...

.PHONY: coverage
coverage: test
	go tool cover -html=cover.out

.PHONY: generate
generate:
	go generate ./...

.PHONY: ci
ci: generate test coverage

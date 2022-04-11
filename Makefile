# Ensure go bin path is in path (especially for CI)
PATH := $(PATH):$(GOPATH)/bin

.PHONY: test
test:
	GO111MODULE=on go test -coverprofile=coverage.txt ./...

.PHONY: coverage
coverage: test
	go tool cover -html=coverage.txt

.PHONY: generate
generate:
	go get -d github.com/vektra/mockery/cmd/mockery
	go generate ./...

.PHONY: ci
ci: check-tidy test coverage

# Ensure there is no unused dependency being added by accident and all generated code is committed
.PHONY: check-tidy
check-tidy: generate
	go mod tidy
	git diff --exit-code

.PHONY: check-headers
check-headers:
	@./check-headers.sh

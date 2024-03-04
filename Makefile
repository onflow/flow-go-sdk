# Ensure go bin path is in path (especially for CI)
PATH := $(PATH):$(GOPATH)/bin

# include script to possibly set a crypto flag for older machines
include crypto_adx_flag.mk

CGO_FLAG := CGO_CFLAGS=$(CRYPTO_FLAG)

.PHONY: test
test:
	GO111MODULE=on $(CGO_FLAG) go test -coverprofile=coverage.txt ./...

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

.PHONY: generate-openapi
generate-openapi:
	swagger-codegen generate -l go -i https://raw.githubusercontent.com/onflow/flow/master/openapi/access.yaml -D packageName=models,modelDocs=false,models -o access/http/models;
	go fmt ./access/http/models

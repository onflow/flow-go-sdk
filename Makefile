# Name of the cover profile
COVER_PROFILE := cover.out
# Disable go sum database lookup for private repos
GOPRIVATE := github.com/dapperlabs/*
# Ensure go bin path is in path (Especially for CI)
PATH := $(PATH):$(GOPATH)/bin

BINARY ?= ./cmd/flow/flow

.PHONY: install-tools
install-tools: check-go-version
	cd ${GOPATH}; \
	GO111MODULE=on go get github.com/vektra/mockery/cmd/mockery@v0.0.0-20181123154057-e78b021dcbb5; \
	GO111MODULE=on go get github.com/axw/gocov/gocov; \
	GO111MODULE=on go get github.com/matm/gocov-html; \
	GO111MODULE=on go get github.com/sanderhahn/gozip/cmd/gozip;

.PHONY: test
test:
	GO111MODULE=on go test -coverprofile=$(COVER_PROFILE) $(if $(JSON_OUTPUT),-json,) ./...

.PHONY: coverage
coverage:
ifeq ($(COVER), true)
	# file has to be called index.html
	gocov convert $(COVER_PROFILE) > cover.json
	./cover-summary.sh
	gocov-html cover.json > index.html
	# coverage.zip will automatically be picked up by teamcity
	gozip -c coverage.zip index.html
endif

.PHONY: generate
generate: generate-mocks

.PHONY: generate-mocks
generate-mocks:
	GO111MODULE=on mockery -name RPCClient -dir=client -case=underscore -output="./client/mocks" -outpkg="mocks"

.PHONY: ci
ci: install-tools generate test coverage

# Check if the go version is >1.13. flow-go-sdk only supports go versions > 1.13
.PHONY: check-go-version
check-go-version:
	go version | grep '1.13\|1.14'

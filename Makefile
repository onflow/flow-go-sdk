# The short Git commit hash
SHORT_COMMIT := $(shell git rev-parse --short HEAD)
# The Git commit hash
COMMIT := $(shell git rev-parse HEAD)
# The tag of the current commit, otherwise empty
VERSION := $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
# Name of the cover profile
COVER_PROFILE := cover.out
# Disable go sum database lookup for private repos
GOPRIVATE=github.com/dapperlabs/*
# Ensure go bin path is in path (Especially for CI)
PATH := $(PATH):$(GOPATH)/bin
# OS
UNAME := $(shell uname)

BINARY ?= ./cmd/flow/flow

# Enable docker build kit
export DOCKER_BUILDKIT := 1


.PHONY: install-tools
install-tools: check-go-version
	cd ${GOPATH}; \
	GO111MODULE=on go get github.com/golang/mock/mockgen@v1.3.1; \
	GO111MODULE=on go get github.com/kevinburke/go-bindata/...@v3.11.0; \
	GO111MODULE=on go get github.com/axw/gocov/gocov; \
	GO111MODULE=on go get github.com/matm/gocov-html; \
	GO111MODULE=on go get github.com/sanderhahn/gozip/cmd/gozip;

.PHONY: test
test:
	GO111MODULE=on go test -mod vendor -coverprofile=$(COVER_PROFILE) $(if $(JSON_OUTPUT),-json,) ./...

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
	GO111MODULE=on mockgen -destination=client/mocks/mock_client.go -package=mocks github.com/dapperlabs/flow-go-sdk/client RPCClient
	GO111MODULE=on mockgen -destination=emulator/mocks/blockchain_api.go -package=mocks github.com/dapperlabs/flow-go-sdk/emulator BlockchainAPI
	GO111MODULE=on mockgen -destination=emulator/storage/mocks/store.go -package=mocks github.com/dapperlabs/flow-go-sdk/emulator/storage Store

.PHONY: ci
ci: install-tools generate test coverage

.PHONY: binary
binary: $(BINARY)

$(BINARY):
	GO111MODULE=on go build \
		-mod vendor \
		-ldflags \
		"-X github.com/dapperlabs/flow-go-sdk/utils/build.commit=$(COMMIT) -X github.com/dapperlabs/flow-go-sdk/utils/build.semver=$(VERSION)" \
		-o $(BINARY) ./cmd/flow

.PHONY: versioned-binaries
versioned-binaries:
	$(MAKE) OS=linux versioned-binary
	$(MAKE) OS=darwin versioned-binary

.PHONY: versioned-binary
versioned-binary:
	GOOS=$(OS) GOARCH=amd64 $(MAKE) BINARY=./cmd/flow/flow-x86_64-$(OS)-$(VERSION) binary

.PHONY: install-cli
install-cli: cmd/flow/flow
	cp cmd/flow/flow $$GOPATH/bin/

.PHONY: docker-build-emulator
docker-build-emulator:
	docker build --ssh default -f cmd/flow/emulator/Dockerfile -t gcr.io/dl-flow/emulator:latest -t "gcr.io/dl-flow/emulator:$(SHORT_COMMIT)" .
ifneq (${VERSION},)
	docker tag gcr.io/dl-flow/emulator:latest gcr.io/dl-flow/emulator:${VERSION}
endif

docker-push-emulator:
	docker push gcr.io/dl-flow/emulator:latest
	docker push "gcr.io/dl-flow/emulator:$(SHORT_COMMIT)"
ifneq (${VERSION},)
	docker push "gcr.io/dl-flow/emulator:${VERSION}"
endif

# Check if the go version is >1.13. flow-go-sdk only supports go versions > 1.13
.PHONY: check-go-version
check-go-version:
	go version | grep '1.13\|1.14'

.PHONY: vendor
vendor:
	GO111MODULE=on go mod vendor

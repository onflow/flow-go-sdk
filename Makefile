# Ensure go bin path is in path (especially for CI)
PATH := $(PATH):$(GOPATH)/bin

.PHONY: install-tools
install-tools:
	cd ${GOPATH}; \
	GO111MODULE=on go get github.com/vektra/mockery/cmd/mockery@v0.0.0-20181123154057-e78b021dcbb5; \
	GO111MODULE=on go get github.com/axw/gocov/gocov; \
	GO111MODULE=on go get github.com/matm/gocov-html; \
	GO111MODULE=on go get github.com/sanderhahn/gozip/cmd/gozip;

.PHONY: test
test:
	GO111MODULE=on go test -coverprofile=cover.out ./...

.PHONY: coverage
coverage: test
	go tool cover -html=cover.out

.PHONY: generate-mocks
generate-mocks:
	GO111MODULE=on mockery -name RPCClient -dir=client -case=underscore -output="./client/mocks" -outpkg="mocks"

.PHONY: ci
ci: install-tools generate-mocks test coverage

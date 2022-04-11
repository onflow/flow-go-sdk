module github.com/onflow/flow-go-sdk

go 1.13

require (
	cloud.google.com/go/kms v1.4.0
	github.com/ethereum/go-ethereum v1.9.9
	github.com/onflow/cadence v0.20.1
	github.com/onflow/flow-go/crypto v0.24.3
	github.com/onflow/flow/protobuf/go/flow v0.2.2
	github.com/onflow/sdks v0.3.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.1-0.20210824115523-ab6dc3262822
	google.golang.org/api v0.70.0
	google.golang.org/genproto v0.0.0-20220222213610-43724f9ea8cf
	google.golang.org/grpc v1.44.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/onflow/cadence => github.com/onflow/cadence v0.21.3-0.20220419065337-d5202c162010
replace github.com/onflow/sdks => /Users/dapper/Dev/sdks

module github.com/onflow/flow-go-sdk

go 1.13

require (
	cloud.google.com/go v0.65.0
	github.com/btcsuite/btcd v0.0.0-20171128150713-2e60448ffcc6
	github.com/ethereum/go-ethereum v1.9.9
	github.com/golang/protobuf v1.4.2
	github.com/onflow/cadence v0.10.2
	github.com/onflow/flow/protobuf/go/flow v0.1.8
	github.com/pkg/errors v0.8.1
	github.com/stretchr/objx v0.1.1 // indirect
	github.com/stretchr/testify v1.5.1
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	google.golang.org/api v0.31.0 // indirect
	google.golang.org/genproto v0.0.0-20200831141814-d751682dd103
	google.golang.org/grpc v1.31.1
)

replace github.com/fxamacker/cbor/v2 => github.com/turbolent/cbor/v2 v2.2.1-0.20200911003300-cac23af49154

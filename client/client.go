// Package client provides grpc API implementation.
// Deprecated: client is deprecated use access instead.
package client

import "github.com/onflow/flow-go-sdk/access/grpc"

// New creates an gRPC client exposing all the common access APIs.
//
// Deprecated: use grpc.NewClient instead.
// Read more in the migration guide:
// https://github.com/onflow/flow-go-sdk/blob/main/docs/migration-v0.25.0.md
var New = grpc.NewBaseClient

// Client is an gRPC client implementing all API access functions.
//
// Deprecated: migrate to access.Client instead or use grpc.BaseClient for grpc specific operations.
// Read more in the migration guide:
// https://github.com/onflow/flow-go-sdk/blob/main/docs/migration-v0.25.0.md
type Client = grpc.BaseClient

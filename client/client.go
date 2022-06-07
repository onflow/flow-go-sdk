/*
 * Flow Go SDK
 *
 * Copyright 2019 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package client provides grpc API implementation.
//
// Deprecated: client is deprecated use access package instead.
package client

import (
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access/grpc"
	"github.com/onflow/flow/protobuf/go/flow/access"
)

// New creates an gRPC client exposing all the common access APIs.
//
// Deprecated: use grpc.NewClient instead.
// Read more in the migration guide:
// https://github.com/onflow/flow-go-sdk/blob/main/docs/migration-v0.25.0.md
var New = grpc.NewBaseClient

// NewFromRPCClient initializes a Flow client using a pre-configured gRPC provider.
//
// Deprecated: use grpc.NewFromRPCClient instead.
// Read more in the migration guide:
// https://github.com/onflow/flow-go-sdk/blob/main/docs/migration-v0.25.0.md
var NewFromRPCClient = grpc.NewFromRPCClient

// Client is an gRPC client implementing all API access functions.
//
// Deprecated: migrate to access.Client instead or use grpc.BaseClient for grpc specific operations.
// Read more in the migration guide:
// https://github.com/onflow/flow-go-sdk/blob/main/docs/migration-v0.25.0.md
type Client = grpc.BaseClient

// An RPCClient is an RPC client for the Flow Access API.
//
// Deprecated: use access.client instead.
// Read more in the migration guide:
// https://github.com/onflow/flow-go-sdk/blob/main/docs/migration-v0.25.0.md
type RPCClient interface {
	access.AccessAPIClient
}

// BlockEvents are the events that occurred in a specific block.
//
// Deprecated: use flow.BlockEvents instead.
// Read more in the migration guide:
// https://github.com/onflow/flow-go-sdk/blob/main/docs/migration-v0.25.0.md
type BlockEvents = flow.BlockEvents

// EventRangeQuery defines a query for Flow events.
//
// Deprecated: use grpc.EventRangeQuery instead.
// Read more in the migration guide:
// https://github.com/onflow/flow-go-sdk/blob/main/docs/migration-v0.25.0.md
type EventRangeQuery = grpc.EventRangeQuery

// RPCError is an error returned by an RPC call to an Access API.
//
// Deprecated: use grpc.RPCError instead.
// Read more in the migration guide:
// https://github.com/onflow/flow-go-sdk/blob/main/docs/migration-v0.25.0.md
type RPCError = grpc.RPCError

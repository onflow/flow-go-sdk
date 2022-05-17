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

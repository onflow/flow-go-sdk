## Migration Guide v0.25.0

The Go SDK version 0.25.0 introduced breaking changes in the API and package naming. 
Changes were required to make the implementation of the new HTTP access node API available. 

We will list all the changes and provide examples on how to migrate.

- **Renamed package: client -> access:** the `client` package was renamed to `access` 
which now includes both `grpc` package containing previously only gRPC implementation and 
also `http` package containing the new HTTP API implementation.
- **Removed package: convert:** the `convert` package was removed and all its functions were moved 
to each of the corresponding `grpc` or `http` packages. The methods were also changed to not be exported, 
so you can no longer use them outside the `convert` package.
- **New clients:** new clients were added each implementing the functions from the client interface 
and exposing a factory for creating them.
- **New Client Interface**: new client interface was created which is now network agnostic, meaning it 
doesn't any more expose additional options in the API that were used to pass gRPC specific options. You can 
still pass those options but you must use the network specific client as shown in the example bellow. 
The interface also changed some functions: 
  - `GetCollectionByID` renamed to `GetCollection`
  - `Close() error` was added


### Migration

#### Creating a Client
Creating a client for communicating with the access node has changed since it's now possible 
to pick and choose between HTTP and gRPC communication protocols. 

*Previous versions:*
```go
// initialize a gRPC emulator client
flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
```

*Version 0.25.0*:
```go
// common client interface
var flowClient access.Client

// initialize an http emulator client
flowClient, err := http.NewClient(http.EmulatorHost)

// initialize a gPRC emulator client
flowClient, err = grpc.NewClient(grpc.EmulatorHost)
```

#### Using the gRPC Client with Options
Using the client is in most cases the same except for the advance case of passing additional 
options to the gRPC client which is no longer possible in the base client, you must use a 
network specific client as shown in the advanced example:

*Previous versions:*
```go
// initialize a gRPC emulator client
flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
latestBlock, err := flowClient.GetLatestBlock(ctx, true, MaxCallSendMsgSize(100))
```

*Version 0.25.0:*
```go
// initialize a grpc network specific client
flowClient, err := NewBaseClient(
	grpc.EmulatorHost, 
	grpc.WithTransportCredentials(insecure.NewCredentials()),
)
latestBlock, err := flowClient.GetLatestBlock(ctx, true, MaxCallSendMsgSize(100))
```

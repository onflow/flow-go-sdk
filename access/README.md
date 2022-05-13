## Access Package
The access package implements network communication with the access nodes APIs. 
It also defines an `access.Client` interface exposing all the common API interactions.  

### Design
Each implementation (currently `grpc` and `http`) include the following parts:
- **Base Client** is the client that implements the client interface exposing 
all generic functionality. 
- **Client** is the network specific client that exposes any additional 
options possible by the specific API implementation.
- **Handler** takes care of actual network communication implementing 
the communication protocol. 


### Usage
If you want to use the base client you should save the instance to 
the client interface, which would allow you to easily switch between network 
implementations like so:

**General Usage**
```go
// common client interface
var flowClient client.Client

// initialize an http emulator client
flowClient, err := http.NewClient(http.EmulatorHost)

// initialize a gPRC emulator client
flowClient, err = grpc.NewClient(grpc.EmulatorHost)
```

**Network Specific Usage** if you require some network specific usage you can also 
instantiate the client like so:
```go
// initialize http specific client
httpClient, err := http.NewHTTPClient(http.EMULATOR_URL)

// initialize grpc specific client
grpcClient, err := grpc.NewGRPCClient(
    grpc.EMULATOR_URL,
    grpcOpts.WithTransportCredentials(insecure.NewCredentials()),
)
```
Read more about this in the docs.

## Development

### Testing
The testing suite is using mock network handlers which can be generated 
by running the following command in the project root directory:
```
make generate
```

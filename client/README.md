
## Client Package
Client package implements network communication with the access nodes APIs. 
It also defines a `Client` interface exposing all the common API interactions.  

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
var flowClient client.Client 
// using HTTP API
flowClient, err := http.NewDefaultTestnetClient()
// using gRPC API
flowClient, err = grpc.NewDefaultTestnetClient()
```

**Network Specific Usage** if you require some network specific usage you can also 
instantiate the client like so:
```go
flowClient, err := http.HTTPClient{}
flowClient.GetBlocksByHeights(context.Background(), BlockQuery{ Heights: []uint64{100},  })
```


## Development

### Testing
The testing suite is using mock network handlers which can be generated 
running the following command in the project root directory:
```
make generate
```

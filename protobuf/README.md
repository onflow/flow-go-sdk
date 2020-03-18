# Protobuf

This folder contains the [protocol buffer](https://developers.google.com/protocol-buffers) files that define the Observation gRPC API. 

## Generating stubs

You can use [prototool](https://github.com/uber/prototool) to generate gRPC client stubs in a variety of languages. Running the command below (in the current directory) will generate stubs for Go and Java:

```shell script
prototool generate
```

_Output files are saved to [/flow-go-sdk/protobuf/out](/out)._

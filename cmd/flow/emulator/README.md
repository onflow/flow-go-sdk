# Emulator

The Flow Emulator is a lightweight tool that emulates the behaviour of the real Flow network. 

The emulator exposes a gRPC server that implements the Flow Observation API, which is designed to have near feature parity
with the real network API.

## Running the emulator with the Flow CLI

The emulator is bundled with the Flow CLI, a command-line interface for working with Flow.

### Installation

Follow [these steps](../../../cli/README.md) to install the Flow CLI.

### Starting the server

You can start the emulator with the Flow CLI:	

```shell script	
flow emulator start --init
```

This command has several useful flags:	

- `-v` Enable verbose logging (this is useful for debugging)	
- `-p` Port to listen on (default: `3569`)	
- `-i` Time interval between blocks (default: `5s`)
- `--persist` Enable persistent storage (uses [Badger](https://github.com/dgraph-io/badger) key-value DB)
- `--dbpath` Path to store database (default: `"./flowdb"`)

### Using the emulator in a project	

You can start the emulator in your project context by running the above command
in the same directory as `flow.json`. This will configure the emulator with your 
project's root account, meaning you can use it to sign and submit transactions.	

## Running the emulator with Docker

Docker builds for the emulator are automatically built and pushed to 
`gcr.io/dl-flow/emulator`, tagged by commit and semantic version. You can also [build the image locally](#building).

### Configuration

In addition to using command-line flags, the emulator can also be configured with environment 
variables, which can be passed into the Docker image.

Here's a sample configuration:

```
FLOW_VERBOSE=true
FLOW_PORT=3596
FLOW_INTERVAL=5s
FLOW_PERSIST=true
FLOW_DBPATH=./path/to/db
```

Here's how to run the emulator Docker image on port 9001 in verbose mode:

```bash
docker run -e FLOW_PORT=9001 -e FLOW_VERBOSE=true gcr.io/dl-flow/emulator
```

The full list of configurable variables is specified [by the config of the `start` command](https://github.com/dapperlabs/flow-go-sdk/blob/master/cli/emulator/start/start.go#L20-L27).
The environment variable names are the upper-case struct fields names, prefixed with `FLOW_`.

#### Accounts

The emulator uses a `flow.json` configuration file to store persistent
configuration, including account keys. In order to start, at least one
key (called the root key) must be configured. This key is used by default
when creating other accounts.

Because Docker does not persist files by default, this file will be 
re-generated each time the emulator image restarts. For situations
where it is important that the emulator always uses the same root key (ie.
unit tests) you can specify a hex-encoded key as an environment variable.

```bash
docker run -e FLOW_ROOTKEY=<hex-encoded key> gcr.io/dl-flow/emulator
```

To generate a root key, use the `keys generate` command.
```bash
docker run gcr.io/dl-flow/emulator keys generate
```

## Building

To build the container locally, use `make docker-build-emulator` from the root of this repository.

Images are automatically built and pushed to `gcr.io/dl-flow/emulator` via [Team City](https://ci.eng.dapperlabs.com/project/Flow_FlowGo_FlowEmulator) on any push to master, i.e. Pull Request merge

## Deployment

The emulator is currently being deployed via [Team City](https://ci.eng.dapperlabs.com/project/Flow_FlowGo_FlowEmulator)
All commands relating to the deployment process live in the `Makefile` in this directory

The deployment has a persistent volume and should keep state persistent. 
The deployment file for the Ingress/Service/Deployment combo, and the
Persistent Volume Claim are located in the k8s folder.

### Accessing the deployment

#### Web Address

*This is only available from the office currently*

It should be available publicly via `https://grpc.emulator.staging.withflow.org/`

#### Port Forwarding

If you have direct access to the cluster, you can use `kubectl` to access the emulator
```bash
export EMULATOR_POD_NAME=$(kubectl get pod -n flow -l app=flow-emulator -o jsonpath="{.items[0].metadata.name}")
kubectl port-forward -n flow $EMULATOR_POD_NAME 3569:3569
```

#### From a Kubernetes Pod

If you are in the same namespace (flow), you can simply access the service via it's service name, e.g. `flow-emulator-v1`.
If you are in a different namespace, you will need to use a [Service without selectors](https://kubernetes.io/docs/concepts/services-networking/service/#services-without-selectors)
e.g.
```yaml
kind: Service
apiVersion: v1
metadata:
  name: flow-emulator-v1
  namespace: YOUR_NAMESPACE
spec:
  type: ExternalName
  externalName: flow-emulator-v1.flow.svc.cluster.local
  ports:
      protocol: TCP
      port: 3569
      targetPort: 3569
```

### Creating your own deployment

Our current deployment settings are available in the [k8s](./k8s) sub directory, if you'd like to setup the emulator in your own Kubernetes cluster. We are using the `Traefik` Ingress controller, but if that is not needed for your purposes, that can be removed, along with any corresponding annotations in the deployment file.

If not using Kubernetes, you can run the Docker container independently. Make sure to run the Docker container with the gRPC port exposed (default is `3569`). Metrics are also available on port `8080` on the `/metrics` endpoint.

To gain persistence for data on the emulator, you will have to provision a volume for the docker container. We've done this through `Persistent Volumes` on Kubernetes, but a mounted volume would suffice. The mount point can be set with the `FLOW_DBPATH` environment variable. We suggest a volume of at least 10GB (100GB for a long-term deployment).

Make sure the emulator also has access to the same `flow.json` file, or always launch it with the same root key, as mentioned above.

```bash
docker run -e FLOW_ROOTKEY=<hex-encoded key> -e FLOW_DBPATH="/flowdb" -v "$(pwd)/flowdb":"/flowdb"  -p 3569:3569 gcr.io/dl-flow/emulator
```

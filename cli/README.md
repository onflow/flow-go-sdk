# Flow CLI

The Flow CLI is a command-line interface that provides useful utilities for building Flow applications.

## Installation

The Flow CLI can be installed in one of three ways:

### Homebrew

MacOS users can install the Flow CLI with Homebrew:

```shell script
brew tap dapperlabs/tap
brew install flow-cli
```

### From a pre-built binary

_This installation method only works on MacOS/x86_64 and Linux/x86_64 architectures._

This script downloads the appropriate binary for your system and moves it to `/usr/local/bin`:

```shell script
sh -ci "$(curl -fsSL https://storage.googleapis.com/flow-cli/install.sh)"
```

### From source

_This installation method works on any system with Go >1.13 installed._

```shell script
GO111MODULE=on go get github.com/dapperlabs/flow-go-sdk/cmd/flow
```

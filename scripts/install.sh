#!/bin/sh

BASE_URL="https://storage.googleapis.com/flow-developer-preview-v1"
# The version to download, set by get_version (defaults to args[1])
VERSION="$1"
# The architecture string, set by get_architecture
ARCH=""

# Get the architecture (CPU, OS) of the current system as a string.
# Only MacOS/x86_64 and Linux/x86_64 architectures are supported.
get_architecture() {
    _ostype="$(uname -s)"
    _cputype="$(uname -m)"
    if [ "$_ostype" = Darwin ] && [ "$_cputype" = i386 ]; then
        if sysctl hw.optional.x86_64 | grep -q ': 1'; then
            _cputype=x86_64
        fi
    fi
    case "$_ostype" in
        Linux)
            _ostype=linux
            ;;
        Darwin)
            _ostype=darwin
            ;;
        *)
            echo "unrecognized OS type: $_ostype"
            return 1
            ;;
    esac
    case "$_cputype" in
        x86_64 | x86-64 | x64 | amd64)
            _cputype=x86_64
            ;;
        *)
            echo "unknown CPU type: $_cputype"
            return 1
            ;;
    esac
    _arch="${_cputype}-${_ostype}"
    ARCH="$_arch"
}

# Get the latest version from remote if none specified in args.
get_version() {
  if [ -z "$VERSION" ]
  then
    VERSION=$(curl -s "$BASE_URL/version.txt")
  fi
}

# Determine the system architecure, download the appropriate binary, and
# install it in `/usr/local/bin` with executable permission.
main() {

  get_architecture || exit 1
  get_version || exit 1

  url="$BASE_URL/flow-$ARCH-$VERSION"
  curl -s "$url" -o ./flow

  # Ensure we don't receive a not found error as response.
  if grep -q "The specified key does not exist" ./flow
  then
    echo "Version $VERSION could not be found"
    exit 1
  fi

  chmod +x ./flow
  mv ./flow /usr/local/bin
  echo "Successfully installed version $VERSION with architecture $ARCH"
}

main

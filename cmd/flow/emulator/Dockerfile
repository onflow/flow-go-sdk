# syntax=docker/dockerfile:1.0.0-experimental

# NOTE: Must be run in the context of the repo's root directory

## Build the app binary
FROM golang:1.13 AS build-app

# Build the app binary in /app
RUN mkdir /app
WORKDIR /app

# add the pubkey of github.com to knownhosts, so ssh-agent doesn't bark
RUN mkdir -p /root/.ssh && ssh-keyscan -t rsa github.com >> /root/.ssh/known_hosts
RUN git config --global 'url.ssh://git@github.com/.insteadOf' https://github.com/

COPY . .

RUN --mount=type=ssh \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GO111MODULE=on CGO_ENABLED=0 GOOOS=linux GOARCH=amd64 \
    go build -ldflags "-extldflags -static" -o ./app ./cmd/flow

## Add the binary to a fresh distroless image
FROM gcr.io/distroless/static

COPY --from=build-app /app/app /bin/app

# Expose GRPC and HTTP ports
EXPOSE 8080
EXPOSE 3569

# Run the CLI binary as the entrypoint
ENTRYPOINT ["/bin/app"]
# These arguments are separated from the entrypoint to simplify running other
# commands with this image.
CMD ["emulator", "start", "--init", "--persist"]

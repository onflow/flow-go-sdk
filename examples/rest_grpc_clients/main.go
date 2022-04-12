package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/client/grpc"
	"github.com/onflow/flow-go-sdk/client/http"
	"github.com/onflow/flow-go-sdk/examples"

	grpcOpts "google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	var flowClient client.Client

	// initialize a http emulator client
	flowClient, err := http.NewDefaultEmulatorClient(false)
	examples.Handle(err)

	latestBlock, err := flowClient.GetLatestBlock(ctx, true)
	examples.Handle(err)

	fmt.Println("Block ID:", latestBlock.BlockHeader.ID.String())

	// initialize a gPRC emulator client
	flowClient, err = grpc.NewDefaultEmulatorClient()
	examples.Handle(err)

	latestBlock, err = flowClient.GetLatestBlock(ctx, true)
	examples.Handle(err)

	fmt.Println("Block ID:", latestBlock.BlockHeader.ID.String())

	// initialize http specific client
	httpClient, err := http.NewHTTPClient(http.EMULATOR_URL)
	examples.Handle(err)

	httpBlocks, err := httpClient.GetBlocksByHeights(
		ctx,
		http.HeightQuery{
			Heights: []uint64{http.SEALED},
		},
	)
	examples.Handle(err)

	fmt.Println("Block ID:", httpBlocks[0].ID.String())

	// initialize grpc specific client
	grpcClient, err := grpc.NewGRPCClient(
		grpc.EMULATOR_URL,
		grpcOpts.WithTransportCredentials(insecure.NewCredentials()),
	)
	examples.Handle(err)

	grpcBlock, err := grpcClient.GetLatestBlockHeader(ctx, true)
	examples.Handle(err)

	fmt.Println("Block ID:", grpcBlock.ID.String())
}

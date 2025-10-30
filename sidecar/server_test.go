package main

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGetSidecarState(t *testing.T) {
	// Connect to the sidecar gRPC server
	conn, err := grpc.Dial("localhost:9191",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		t.Fatalf("Failed to connect to sidecar: %v", err)
	}
	defer conn.Close()

	// Create the client
	client := api.NewSidecarServiceClient(conn)

	// Call GetSidecarState
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetSidecarState(ctx, &api.SidecarStateRequest{})
	if err != nil {
		t.Fatalf("GetSidecarState failed: %v", err)
	}

	// Print the output
	fmt.Println("\n=== Sidecar State ===")
	fmt.Printf("Ethereum Block Height: %d\n", resp.EthBlockHeight)
	fmt.Printf("Ethereum Gas Limit: %d\n", resp.EthGasLimit)
	fmt.Printf("Ethereum Base Fee: %d\n", resp.EthBaseFee)
	fmt.Printf("Ethereum Tip Cap: %d\n", resp.EthTipCap)
	fmt.Printf("ROCK/USD Price: %s\n", resp.ROCKUSDPrice)
	fmt.Printf("BTC/USD Price: %s\n", resp.BTCUSDPrice)
	fmt.Printf("ETH/USD Price: %s\n", resp.ETHUSDPrice)
	fmt.Printf("ZEC/USD Price: %s\n", resp.ZECUSDPrice)
	fmt.Printf("Sidecar Version: %s\n", resp.SidecarVersionName)
	fmt.Printf("Ethereum Burn Events: %d\n", len(resp.EthBurnEvents))
	fmt.Printf("Redemptions: %d\n", len(resp.Redemptions))
	fmt.Printf("Solana Burn Events: %d\n", len(resp.SolanaBurnEvents))
	fmt.Printf("Solana Mint Events: %d\n", len(resp.SolanaMintEvents))

	log.Printf("Successfully retrieved sidecar state")
}

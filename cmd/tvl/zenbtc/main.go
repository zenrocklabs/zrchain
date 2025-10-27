package main

import (
	"context"
	"flag"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	
	bindings "github.com/Zenrock-Foundation/zrchain/v6/zenbtc/bindings"
)

const (
	// Default Solana mainnet RPC endpoint
	defaultSolanaRPC = "https://api.mainnet-beta.solana.com"

	// Default zenBTC program ID on Solana mainnet
	defaultProgramID = "9t9RfpterTs95eXbKQWeAriZqET13TbjwDa6VW6LJHFb"
	
	// Default Ethereum mainnet RPC endpoint
	defaultEthereumRPC = "https://eth.llamarpc.com"
	
	// zenBTC token contract address on Ethereum mainnet
	ethereumZenBTCAddress = "0x2fE9754d5D28bac0ea8971C0Ca59428b8644C776"
	
	// zenBTC decimals (same as Bitcoin)
	zenBTCDecimals = 8
)

// GetMintAddress derives the zenBTC mint address from the program ID
// using the "wrapped_mint" seed
func GetMintAddress(programID solana.PublicKey) (solana.PublicKey, uint8, error) {
	seeds := [][]byte{
		[]byte("wrapped_mint"),
	}
	addr, bump, err := solana.FindProgramAddress(seeds, programID)
	if err != nil {
		return solana.PublicKey{}, 0, err
	}
	return addr, bump, nil
}

// QueryEthereumSupply queries the zenBTC total supply on Ethereum mainnet
func QueryEthereumSupply(ctx context.Context, rpcURL string) (*big.Int, error) {
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum RPC: %w", err)
	}
	defer client.Close()

	contractAddr := common.HexToAddress(ethereumZenBTCAddress)
	token, err := bindings.NewZenBTC(contractAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate zenBTC contract: %w", err)
	}

	totalSupply, err := token.TotalSupply(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query total supply: %w", err)
	}

	return totalSupply, nil
}

// FormatZenBTC formats a raw amount (in smallest units) to zenBTC with 8 decimals
func FormatZenBTC(amount *big.Int) string {
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(zenBTCDecimals), nil)
	whole := new(big.Int).Div(amount, divisor)
	remainder := new(big.Int).Mod(amount, divisor)
	
	// Format with 8 decimal places
	return fmt.Sprintf("%s.%08d", whole.String(), remainder.Int64())
}

func main() {
	// Define command-line flags
	solanaRPC := flag.String("solana-rpc", defaultSolanaRPC, "Solana RPC endpoint URL")
	ethereumRPC := flag.String("ethereum-rpc", defaultEthereumRPC, "Ethereum RPC endpoint URL")
	programIDStr := flag.String("program-id", defaultProgramID, "zenBTC program ID on Solana")
	flag.Parse()

	// Parse program ID
	programID, err := solana.PublicKeyFromBase58(*programIDStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid program ID: %v\n", err)
		os.Exit(1)
	}

	// Derive mint address
	mintAddress, bump, err := GetMintAddress(programID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to derive mint address: %v\n", err)
		os.Exit(1)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Query Solana supply
	solanaClient := rpc.New(*solanaRPC)
	solanaSupply, err := solanaClient.GetTokenSupply(ctx, mintAddress, rpc.CommitmentFinalized)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to query Solana token supply: %v\n", err)
		os.Exit(1)
	}

	// Query Ethereum supply
	ethSupply, err := QueryEthereumSupply(ctx, *ethereumRPC)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to query Ethereum token supply: %v\n", err)
		os.Exit(1)
	}

	// Parse Solana supply as big.Int for total calculation
	solanaSupplyRaw := new(big.Int)
	solanaSupplyRaw.SetString(solanaSupply.Value.Amount, 10)

	// Calculate total supply across both chains
	totalSupply := new(big.Int).Add(solanaSupplyRaw, ethSupply)

	// Print results
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("          zenBTC Cross-Chain Total Value Locked")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()
	
	// Solana section
	fmt.Println("┌─────────────────────────────────────────────────────────┐")
	fmt.Println("│                    SOLANA MAINNET                       │")
	fmt.Println("└─────────────────────────────────────────────────────────┘")
	fmt.Println()
	fmt.Printf("RPC Endpoint:      %s\n", *solanaRPC)
	fmt.Printf("Program ID:        %s\n", programID.String())
	fmt.Printf("Mint Address:      %s\n", mintAddress.String())
	fmt.Printf("Mint Bump Seed:    %d\n", bump)
	fmt.Println()
	fmt.Printf("Total Supply:      %s zenBTC\n", solanaSupply.Value.UiAmountString)
	fmt.Printf("Raw Amount:        %s\n", solanaSupply.Value.Amount)
	fmt.Printf("Slot:              %d\n", solanaSupply.Context.Slot)
	fmt.Println()
	
	// Ethereum section
	fmt.Println("┌─────────────────────────────────────────────────────────┐")
	fmt.Println("│                   ETHEREUM MAINNET                      │")
	fmt.Println("└─────────────────────────────────────────────────────────┘")
	fmt.Println()
	fmt.Printf("RPC Endpoint:      %s\n", *ethereumRPC)
	fmt.Printf("Contract Address:  %s\n", ethereumZenBTCAddress)
	fmt.Println()
	fmt.Printf("Total Supply:      %s zenBTC\n", FormatZenBTC(ethSupply))
	fmt.Printf("Raw Amount:        %s\n", ethSupply.String())
	fmt.Println()
	
	// Total section
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("                   TOTAL ACROSS CHAINS")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Printf("Combined TVL:      %s zenBTC\n", FormatZenBTC(totalSupply))
	fmt.Printf("Raw Amount:        %s\n", totalSupply.String())
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════")
}

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/big"
	"net/http"
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

	// Zenrock API endpoint for zenBTC supply info
	zenrockAPIURL = "https://api.diamond.zenrocklabs.io/zenbtc/supply"
)

// ZenrockSupplyResponse represents the response from Zenrock API
type ZenrockSupplyResponse struct {
	CustodiedBTC  string `json:"custodiedBTC"`
	TotalZenBTC   string `json:"totalZenBTC"`
	MintedZenBTC  string `json:"mintedZenBTC"`
	PendingZenBTC string `json:"pendingZenBTC"`
	ExchangeRate  string `json:"exchangeRate"`
}

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

// QueryCustodiedBTC queries the Zenrock API for the amount of BTC custodied
func QueryCustodiedBTC(ctx context.Context) (*big.Int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", zenrockAPIURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to query Zenrock API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var supplyResp ZenrockSupplyResponse
	if err := json.Unmarshal(body, &supplyResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	custodiedBTC := new(big.Int)
	if _, ok := custodiedBTC.SetString(supplyResp.CustodiedBTC, 10); !ok {
		return nil, fmt.Errorf("failed to parse custodiedBTC value: %s", supplyResp.CustodiedBTC)
	}

	return custodiedBTC, nil
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

	// Query custodied BTC from Zenrock API
	custodiedBTC, err := QueryCustodiedBTC(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to query custodied BTC: %v\n", err)
		os.Exit(1)
	}

	// Parse Solana supply as big.Int for total calculation
	solanaSupplyRaw := new(big.Int)
	solanaSupplyRaw.SetString(solanaSupply.Value.Amount, 10)

	// Calculate total supply across both chains
	totalSupply := new(big.Int).Add(solanaSupplyRaw, ethSupply)

	// Calculate percentages
	solanaPercent := 0.0
	ethPercent := 0.0
	if totalSupply.Sign() > 0 {
		totalFloat := new(big.Float).SetInt(totalSupply)
		solanaFloat := new(big.Float).SetInt(solanaSupplyRaw)
		ethFloat := new(big.Float).SetInt(ethSupply)

		solanaRatio := new(big.Float).Quo(solanaFloat, totalFloat)
		ethRatio := new(big.Float).Quo(ethFloat, totalFloat)

		solanaPercent, _ = solanaRatio.Float64()
		ethPercent, _ = ethRatio.Float64()
		solanaPercent *= 100
		ethPercent *= 100
	}

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
	fmt.Printf("Percentage:        %.2f%%\n", solanaPercent)
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
	fmt.Printf("Percentage:        %.2f%%\n", ethPercent)
	fmt.Println()

	// Total section
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("                   TOTAL ACROSS CHAINS")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Printf("Combined TVL:      %s zenBTC\n", FormatZenBTC(totalSupply))
	fmt.Printf("Raw Amount:        %s\n", totalSupply.String())
	fmt.Println()
	fmt.Println("Distribution:")
	fmt.Printf("  Solana:          %.2f%%\n", solanaPercent)
	fmt.Printf("  Ethereum:        %.2f%%\n", ethPercent)
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("                     BACKING ASSETS")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Printf("Locked BTC:        %s BTC\n", FormatZenBTC(custodiedBTC))
	fmt.Printf("Raw Satoshis:      %s\n", custodiedBTC.String())
	fmt.Println()

	// Calculate collateralization ratio
	collateralRatio := 0.0
	if totalSupply.Sign() > 0 {
		custodiedFloat := new(big.Float).SetInt(custodiedBTC)
		totalFloat := new(big.Float).SetInt(totalSupply)
		ratio := new(big.Float).Quo(custodiedFloat, totalFloat)
		collateralRatio, _ = ratio.Float64()
	}

	fmt.Printf("Backing Ratio:     %.6f (%.4f%%)\n", collateralRatio, collateralRatio*100)
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════")
}

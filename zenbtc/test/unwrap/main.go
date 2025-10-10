package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zenrocklabs/goem/ethereum"
	bindings "github.com/zenrocklabs/zenbtc/bindings"
)

func main() {
	// Define command line flags
	privateKeyHex := flag.String("private-key", "", "Private key in hex format (with or without 0x prefix)")
	btcAddress := flag.String("btc-address", "tb1qypwjx7yj5jz0gw0vh76348ypa2ns7tfwsnhlh9", "Bitcoin testnet address")
	amount := flag.String("amount", "100000", "Amount to unwrap (in satoshis)")
	rpcURL := flag.String("rpc-url", "https://rpc.ankr.com/eth_hoodi", "Ethereum RPC URL")
	contractAddr := flag.String("contract", "0xEe6dd71ccf66E3F920a4D49a57020e0F89659407", "ZenBTC contract address")

	flag.Parse()

	// Validate required flags
	if *privateKeyHex == "" || *btcAddress == "" || *amount == "" || *contractAddr == "" {
		fmt.Println("Missing required flags. Use --help for usage information.")
		os.Exit(1)
	}

	// Connect to Ethereum client
	client, err := ethclient.Dial(*rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}

	// Verify we're on Hoodi
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}
	if chainID.Cmp(ethereum.HoodiChainId) != 0 {
		log.Fatalf("Wrong network: expected Hoodi (chain ID %s), got chain ID %s", ethereum.HoodiChainId.String(), chainID.String())
	}

	// Process private key
	privateKeyStr := strings.TrimPrefix(*privateKeyHex, "0x")
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// Get the address from the private key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Create auth
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, ethereum.HoodiChainId)
	if err != nil {
		log.Fatalf("Failed to create auth: %v", err)
	}

	// Convert BTC address to bytes
	btcBytes, err := btcAddressToBytes(*btcAddress)
	if err != nil {
		log.Fatalf("Failed to convert BTC address: %v", err)
	}

	// Parse amount
	value := new(big.Int)
	value.SetString(*amount, 10)

	// Create contract instance
	contract, err := bindings.NewZenBTC(common.HexToAddress(*contractAddr), client)
	if err != nil {
		log.Fatalf("Failed to instantiate contract: %v", err)
	}

	// Get fee estimate
	fee, err := contract.EstimateFee(&bind.CallOpts{}, value)
	if err != nil {
		log.Fatalf("Failed to estimate fee: %v", err)
	}
	fmt.Printf("Estimated fee: %s\n", fee.String())

	// Get balance
	balance, err := contract.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}
	fmt.Printf("Current balance: %s\n", balance.String())

	if balance.Cmp(value) < 0 {
		log.Fatalf("Insufficient balance: have %s, need %s", balance.String(), value.String())
	}

	// Submit transaction
	auth.GasLimit = uint64(300000) // Set appropriate gas limit

	fmt.Printf("\nSubmitting unwrap transaction:\n")
	fmt.Printf("From: %s\n", address.Hex())
	fmt.Printf("Amount: %s\n", value.String())
	fmt.Printf("BTC Address: %s\n", *btcAddress)

	tx, err := contract.Unwrap(auth, value, btcBytes)
	if err != nil {
		log.Fatalf("Failed to submit unwrap transaction: %v", err)
	}

	fmt.Printf("Transaction submitted: %s\n", tx.Hash().Hex())

	// Wait for transaction receipt
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatalf("Failed to get transaction receipt: %v", err)
	}

	fmt.Printf("Transaction mined in block: %d\n", receipt.BlockNumber)
	if receipt.Status == 0 {
		log.Fatalf("Transaction failed!")
	}
}

func btcAddressToBytes(addr string) ([]byte, error) {
	// Check if it's a bech32 address
	if strings.HasPrefix(addr, "tb1") {
		hrp, _, err := bech32.Decode(addr)
		if err != nil {
			return nil, fmt.Errorf("invalid bech32 address: %v", err)
		}
		if hrp != "tb" {
			return nil, fmt.Errorf("invalid network prefix: expected 'tb', got '%s'", hrp)
		}
		// Convert the address to bytes maintaining the bech32 format
		return []byte(addr), nil
	}

	// For legacy or nested SegWit addresses
	address, err := btcutil.DecodeAddress(addr, &chaincfg.TestNet3Params)
	if err != nil {
		return nil, fmt.Errorf("invalid BTC address: %v", err)
	}

	// Verify it's a testnet address
	if !address.IsForNet(&chaincfg.TestNet3Params) {
		return nil, fmt.Errorf("address is not for Bitcoin testnet")
	}

	// Convert to bytes
	return []byte(address.EncodeAddress()), nil
}

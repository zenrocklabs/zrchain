package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/Zenrock-Foundation/zrchain/v6/contracts/solrock"
	rock_spl_token "github.com/Zenrock-Foundation/zrchain/v6/contracts/solrock/generated/rock_spl_token"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func main() {
	var (
		signerPath  = flag.String("signer", "", "Path to the signer keypair file")
		mintAddress = flag.String("mint", "", "Mint address of the token to unwrap")
		programID   = flag.String("program", "", "Program ID for the unwrap instruction")
		solanaRPC   = flag.String("rpc", "https://api.devnet.solana.com", "Solana RPC URL")
		amount      = flag.Uint64("amount", 0, "Amount of tokens to unwrap")
		destAddr    = flag.String("dest", "", "Destination address for the unwrapped tokens")
		feeWallet   = flag.String("fee-wallet", "", "Fee wallet address for the unwrap transaction")
	)

	// Parse flags
	flag.Parse()

	// Validate required flags
	if *signerPath == "" {
		log.Fatal("Error: --signer flag is required")
	}
	if *mintAddress == "" {
		log.Fatal("Error: --mint flag is required")
	}
	if *programID == "" {
		log.Fatal("Error: --program flag is required")
	}
	if *amount == 0 {
		log.Fatal("Error: --amount flag is required and must be greater than 0")
	}
	if *destAddr == "" {
		log.Fatal("Error: --dest flag is required")
	}
	if *solanaRPC == "" {
		log.Fatal("Error: --rpc flag is required")
	}
	if *feeWallet == "" {
		log.Fatal("Error: --fee-wallet flag is required")
	}

	// Convert string addresses to PublicKey
	programPubkey := solana.MustPublicKeyFromBase58(*programID)
	mintPubkey := solana.MustPublicKeyFromBase58(*mintAddress)
	feeWalletPubkey := solana.MustPublicKeyFromBase58(*feeWallet)

	// Call the unwrap function
	if err := SubmitUnwrapTx(*solanaRPC, *signerPath, *destAddr, programPubkey, mintPubkey, feeWalletPubkey, *amount); err != nil {
		log.Fatalf("Error executing unwrap transaction: %v", err)
	}

	fmt.Println("Unwrap transaction completed successfully!")
}

func SubmitUnwrapTx(solanaRPC, signerPath, destAddr string, programID, mintAddress, feeWallet solana.PublicKey, amount uint64) error {
	signer, err := solana.PrivateKeyFromSolanaKeygenFile(signerPath)
	if err != nil {
		return fmt.Errorf("failed to load signer: %v", err)
	}

	// Decode the bech32 Cosmos address with "zen" prefix
	accAddr, err := types.AccAddressFromBech32(destAddr)
	if err != nil {
		return fmt.Errorf("failed to decode bech32 address: %v", err)
	}

	// Convert the Cosmos account address to bytes
	decodedBytes := accAddr.Bytes()
	fmt.Printf("Decoded bech32 address - Length: %d bytes\n", len(decodedBytes))

	// Convert decoded bytes to [25]uint8 - truncate if longer, pad with zeros if shorter
	var destAddrArray [25]uint8
	if len(decodedBytes) >= 25 {
		copy(destAddrArray[:], decodedBytes[:25])
	} else {
		copy(destAddrArray[:], decodedBytes)
	}

	args := rock_spl_token.UnwrapArgs{
		Value:    amount,
		DestAddr: destAddrArray,
	}

	// Create unwrap instruction using the existing solrock.Unwrap function
	instruction := solrock.Unwrap(programID, args, signer.PublicKey(), mintAddress, feeWallet)

	// Get recent blockhash
	client := rpc.New(solanaRPC)
	ctx := context.Background()

	recentBlockhash, err := client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return fmt.Errorf("failed to get recent blockhash: %v", err)
	}

	// Create transaction with the correct program ID
	tx, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		recentBlockhash.Value.Blockhash,
		solana.TransactionPayer(signer.PublicKey()),
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %v", err)
	}

	// Sign transaction
	_, err = tx.PartialSign(func(key solana.PublicKey) *solana.PrivateKey {
		if signer.PublicKey().Equals(key) {
			return &signer
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Send transaction
	signature, err := client.SendTransaction(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}

	fmt.Printf("Transaction sent successfully! https://explorer.solana.com/tx/%s?cluster=devnet\n", signature.String())
	return nil
}

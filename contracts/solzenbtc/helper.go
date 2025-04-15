package solzenbtc

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

var NONCE_ACCOUNT_LENGTH = uint64(80)
var MULTISIG_SIZE = uint64(355)

func WaitForTransactionConfirmation(ctx context.Context, client *rpc.Client, signature solana.Signature, commitment rpc.CommitmentType) (*rpc.GetTransactionResult, error) {
	timeout := time.After(5 * time.Minute) // Set a 5-minute timeout - adjust as needed

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timeout:
			return nil, fmt.Errorf("timeout waiting for transaction confirmation")
		default:
			res, err := client.GetTransaction(ctx, signature, &rpc.GetTransactionOpts{
				Commitment: commitment,
			})
			if err == nil {
				return res, nil
			}
			time.Sleep(1 * time.Second) // Add 1 second delay between attempts
		}
	}
}

func SendTransaction(
	client *rpc.Client,
	ctx context.Context,
	transaction *solana.Transaction,
) (signature solana.Signature, err error) {
	opts := rpc.TransactionOpts{
		SkipPreflight:       false,
		PreflightCommitment: rpc.CommitmentConfirmed,
	}

	return client.SendTransactionWithOpts(
		ctx,
		transaction,
		opts,
	)
}

func SignTransaction(tx *solana.Transaction, signer solana.PrivateKey) (out []solana.Signature, err error) {
	return tx.PartialSign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if signer.PublicKey().Equals(key) {
				return &signer
			}
			return nil
		},
	)
}

func GetAccountInfo(ctx context.Context, client *rpc.Client, account solana.PublicKey) (out *rpc.GetAccountInfoResult, err error) {
	return client.GetAccountInfoWithOpts(
		ctx,
		account,
		&rpc.GetAccountInfoOpts{
			Commitment: rpc.CommitmentConfirmed,
			DataSlice:  nil,
		},
	)
}

func LoadLocalWallet() (solana.PrivateKey, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return solana.PrivateKey{}, err
	}

	return solana.PrivateKeyFromSolanaKeygenFile(filepath.Join(homeDir, ".config/solana/id.json"))
}

func RequestAirdrop(ctx context.Context, client *rpc.Client, account solana.PublicKey, amount uint64) error {
	signature, err := client.RequestAirdrop(ctx, account, amount, rpc.CommitmentConfirmed)
	if err != nil {
		return err
	}

	_, err = WaitForTransactionConfirmation(
		ctx,
		client,
		signature,
		rpc.CommitmentConfirmed,
	)
	return err
}

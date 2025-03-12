package solrock

//
//import (
//	"context"
//	"encoding/hex"
//	"fmt"
//	"math/big"
//	"testing"
//
//	"github.com/Zenrock-Foundation/zrchain/v5/contracts/solrock/generated/zenbtc_spl_token"
//	"github.com/gagliardetto/solana-go"
//	ata "github.com/gagliardetto/solana-go/programs/associated-token-account"
//	"github.com/gagliardetto/solana-go/programs/system"
//	"github.com/gagliardetto/solana-go/rpc"
//	"github.com/stretchr/testify/require"
//)
//
//func TestDurableNonces2(t *testing.T) {
//	nonceAuth, err := solana.WalletFromPrivateKeyBase58("3E2YL7Sf5EZS6kXrHEd34byyF5Kzh7PCPH9Q4tRyRaTJfcHmbzr9nPUWgJ9t5Xy9LBg9BXJFdbuJwTtpdNoJS1Yh") //solana.NewWallet()
//	fmt.Printf("pubkey: %s\n", nonceAuth.PublicKey().String())
//	require.NoError(t, err)
//	nonceWallet, err := solana.WalletFromPrivateKeyBase58("7hLFPdmGFZkzPjdi75qsYfFg9gpaZttwa1bqvJmzKydGodV5PQqQ1oi7HFaT54F8UmcZb4vjuBa4RWMwFJuKy7k") //solana.NewWallet()
//	require.NoError(t, err)
//	fmt.Printf("nonceWallet: %s\n", nonceWallet.PublicKey().String())
//	sender := solana.NewWallet()
//	receiver := solana.NewWallet().PublicKey()
//
//	//err := RequestAirdrop(context.Background(), client, nonceAuth.PublicKey(), 1000*solana.LAMPORTS_PER_SOL)
//	//require.NoError(t, err)
//	//
//	//err = RequestAirdrop(context.Background(), client, sender.PublicKey(), 1000*solana.LAMPORTS_PER_SOL)
//	//require.NoError(t, err)
//	recentHash, err := solana.HashFromBase58("HmhsXVRocBP2NuTkVoSgqfpEhuygfVTW6nYVKmshXPCA")
//
//	t.Run("Creates a nonce account", func(t *testing.T) {
//		//recent, err := client.GetLatestBlockhash(context.Background(), rpc.CommitmentConfirmed)
//		require.NoError(t, err)
//
//		tx, err := solana.NewTransaction(
//			[]solana.Instruction{
//				system.NewCreateAccountInstruction(
//					uint64(0.0015*float64(solana.LAMPORTS_PER_SOL)),
//					NONCE_ACCOUNT_LENGTH,
//					solana.SystemProgramID,
//					nonceAuth.PublicKey(),
//					nonceWallet.PublicKey(),
//				).Build(),
//
//				system.NewInitializeNonceAccountInstruction(
//					nonceAuth.PublicKey(),
//					nonceWallet.PublicKey(),
//					solana.SysVarRecentBlockHashesPubkey,
//					solana.SysVarRentPubkey,
//				).Build(),
//			},
//			recentHash,
//			solana.TransactionPayer(nonceAuth.PublicKey()),
//		)
//		require.NoError(t, err)
//
//		_, err = SignTransaction(tx, nonceAuth.PrivateKey)
//		require.NoError(t, err)
//
//		_, err = SignTransaction(tx, nonceWallet.PrivateKey)
//
//		require.NoError(t, err)
//		bin, err := tx.MarshalBinary()
//		require.NoError(t, err)
//		fmt.Printf("tx: %s\n", string(bin))
//		signature, err := SendTransaction(client, context.Background(), tx)
//		require.NoError(t, err)
//
//		confirmedTx, err := WaitForTransactionConfirmation(
//			context.Background(),
//			client,
//			signature,
//			rpc.CommitmentConfirmed,
//		)
//		require.NoError(t, err)
//
//		require.Empty(t, confirmedTx.Meta.Err) // The transaction succeeded if there's no error
//
//		nonceAccount, err := GetNonceAccount(context.Background(), client, nonceWallet.PublicKey())
//		require.NoError(t, err)
//
//		require.Equal(t, nonceAccount.AuthorizedPubkey, nonceAuth.PublicKey())
//	})
//
//	t.Run("Performs a transfer using a durable nonce", func(t *testing.T) {
//		noncePubkey := nonceWallet.PublicKey()
//
//		nonceAccountBefore, err := GetNonceAccount(context.Background(), client, noncePubkey)
//		require.NoError(t, err)
//
//		tx, err := solana.NewTransaction(
//			[]solana.Instruction{
//				system.NewAdvanceNonceAccountInstruction(
//					noncePubkey,
//					solana.SysVarRecentBlockHashesPubkey,
//					nonceAuth.PublicKey(),
//				).Build(),
//				system.NewTransferInstruction(
//					uint64(0.01*float64(solana.LAMPORTS_PER_SOL)),
//					sender.PublicKey(),
//					receiver,
//				).Build(),
//			},
//			solana.Hash(nonceAccountBefore.Nonce),
//			solana.TransactionPayer(sender.PublicKey()),
//		)
//		require.NoError(t, err)
//
//		_, err = SignTransaction(tx, nonceAuth.PrivateKey)
//		require.NoError(t, err)
//
//		_, err = SignTransaction(tx, sender.PrivateKey)
//		require.NoError(t, err)
//
//		signature, err := SendTransaction(client, context.Background(), tx)
//		require.NoError(t, err)
//
//		confirmedTx, err := WaitForTransactionConfirmation(
//			context.Background(),
//			client,
//			signature,
//			rpc.CommitmentConfirmed,
//		)
//		require.NoError(t, err)
//
//		require.Empty(t, confirmedTx.Meta.Err) // The transaction succeeded if there's no error
//
//		nonceAccountAfter, err := GetNonceAccount(context.Background(), client, noncePubkey)
//		require.NoError(t, err)
//
//		require.NotEqual(t, nonceAccountBefore.Nonce, nonceAccountAfter.Nonce)
//	})
//}

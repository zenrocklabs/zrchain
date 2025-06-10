package solrock

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"
)

func TestCreateDurableNonceAccount(t *testing.T) {
	var client = rpc.New("https://api.devnet.solana.com")

	nonceAuthPubKey := solana.MustPublicKeyFromBase58("9WB3KQZ5d2jqXuufWX6vidYzjgEgnRLbxYBWYCRsfzmw")
	nonceAccPubKey := solana.MustPublicKeyFromBase58("GqojxUhvnNdHKztS2ZZ9HMAZ8VvXv4wmaB6d5MQxz9se")
	recent, err := client.GetLatestBlockhash(context.Background(), rpc.CommitmentConfirmed)
	require.NoError(t, err)

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewCreateAccountInstruction(
				uint64(0.0015*float64(solana.LAMPORTS_PER_SOL)),
				NONCE_ACCOUNT_LENGTH,
				solana.SystemProgramID,
				nonceAuthPubKey,
				nonceAccPubKey,
			).Build(),

			system.NewInitializeNonceAccountInstruction(
				nonceAuthPubKey,
				nonceAccPubKey,
				solana.SysVarRecentBlockHashesPubkey,
				solana.SysVarRentPubkey,
			).Build(),
		},
		recent.Value.Blockhash,
		solana.TransactionPayer(nonceAuthPubKey),
	)
	require.NoError(t, err)
	bin, err := tx.Message.MarshalBinary()
	require.NoError(t, err)

	fmt.Printf("unsigned transaction: %s", hex.EncodeToString(bin))
}

func TestParseDurableNonceTransaction(t *testing.T) {
	txB64 := "AAIAAwXm3dWSnU+s1iGpvDycQw5mTg4smpU65kusWo4VRYRamexoOnfU4JeVvQ7ggopVP8es3TyNU3+VlFV5tngMQp03BqfVFxksVo7gioRfc9KXiM8DXDFFshqzRNgGLqlAAAAGp9UXGSxcUSGMyUw9SvF/WNruCJuh/UTj29mKAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAifT0s57LkkWQSpOS/6DVc9hstj1+edfJ0NXqI9eVk7UCBAIAATQAAAAAYOMWAAAAAABQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAMBAgMkBgAAAObd1ZKdT6zWIam8PJxDDmZODiyalTrmS6xajhVFhFqZ"

	//rawTx, err := base64.StdEncoding.DecodeString(txB64)
	//require.NoError(t, err)

	tx := &solana.Transaction{
		Message: solana.Message{},
	}

	err := tx.UnmarshalBase64(txB64)
	for _, i := range tx.Message.AccountKeys {
		fmt.Println(i.String())
	}
	require.NoError(t, err)
	//err = tx.Message.UnmarshalWithDecoder(bin.NewBinDecoder(rawTx))
	//require.NoError(t, err)
}

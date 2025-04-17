package solrock

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/contracts/solrock/generated/rock_spl_token"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"
)

func TestBurn(t *testing.T) {
	signer, err := solana.PublicKeyFromBase58("5pVWFKJtfA52zJZ9zeHVSn1CzSZJAMGG82iS8VYwfyV5")
	require.NoError(t, err)
	feeWallet, err := solana.PublicKeyFromBase58("FzqGcRG98v1KhKxatX2Abb2z1aJ2rViQwBK5GHByKCAd")
	require.NoError(t, err)
	value := uint64(50)
	alice := "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
	aliceBytes, err := sdk.GetFromBech32(alice, "zen")
	require.NoError(t, err)
	aliceAddress := [25]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	copy(aliceAddress[:], aliceBytes)

	var client = rpc.New("https://api.devnet.solana.com")

	recent, err := client.GetLatestBlockhash(context.Background(), rpc.CommitmentConfirmed)
	require.NoError(t, err)

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			Unwrap(
				programID,
				rock_spl_token.UnwrapArgs{
					Value:    value,
					DestAddr: aliceAddress,
				},
				signer,
				mintAddress,
				feeWallet,
			),
		},
		recent.Value.Blockhash,
		solana.TransactionPayer(signer),
	)
	require.NoError(t, err)
	require.NoError(t, err)
	bin, err := tx.Message.MarshalBinary()
	require.NoError(t, err)
	fmt.Printf("transaction : %s", hex.EncodeToString(bin))
}

package keeper

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v5/contracts/solrock"
	"github.com/Zenrock-Foundation/zrchain/v5/contracts/solrock/generated/zenbtc_spl_token"
	treasuryTypes "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	ata "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
)

const durableNonceKey = "solanaDurableNonce"

func (k Keeper) PrepareSolRockMintTx(goCtx context.Context, amount uint64, signer, recipient *treasuryTypes.Key) (string, error) {
	params := k.GetParams(goCtx).Solana

	programID, err := solana.PublicKeyFromBase58(params.ProgramId)
	if err != nil {
		return "", err
	}

	nonceAccPubKey, err := solana.PublicKeyFromBase58(params.NonceAccountPubKey)
	if err != nil {
		return "", err
	}
	nonceAuthPubKey, err := solana.PublicKeyFromBase58(params.NonceAuthorityPubKey)
	if err != nil {
		return "", err
	}

	signerAddress, err := treasuryTypes.SolanaAddress(signer)
	if err != nil {
		return "", err
	}
	signerKey, err := solana.PublicKeyFromBase58(signerAddress)
	if err != nil {
		return "", err
	}

	mintKey, err := solana.PublicKeyFromBase58(params.MintAddress)
	if err != nil {
		return "", err
	}

	feeKey, err := solana.PublicKeyFromBase58(params.FeeWallet)
	if err != nil {
		return "", err
	}

	recipientAddress, err := treasuryTypes.SolanaAddress(recipient)
	if err != nil {
		return "", err
	}

	recipientKey, err := solana.PublicKeyFromBase58(recipientAddress)
	if err != nil {
		return "", err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	nonce, err := k.getSolanaDurableNonce(ctx)
	if err != nil {
		return "", err
	}

	var instructions []solana.Instruction

	feeWalletAta, _, err := solana.FindAssociatedTokenAddress(feeKey, mintKey)
	if err != nil {
		return "", err
	}

	client := rpc.New(params.RpcUrl)
	_, err = solrock.GetTokenAccount(context.Background(), client, feeWalletAta)

	if err.Error() == "not found" {
		instructions = append(
			instructions,
			ata.NewCreateInstruction(
				signerKey,
				feeKey,
				mintKey,
			).Build(),
		)
	} else {
		return "", err
	}

	receiverAta, _, err := solana.FindAssociatedTokenAddress(recipientKey, mintKey)
	if err != nil {
		return "", err
	}

	_, err = solrock.GetTokenAccount(context.Background(), client, receiverAta)

	if err.Error() == "not found" {
		instructions = append(
			instructions,
			ata.NewCreateInstruction(
				signerKey,
				recipientKey,
				mintKey,
			).Build(),
		)
	}

	instructions = append(instructions, system.NewAdvanceNonceAccountInstruction(
		nonceAccPubKey,
		solana.SysVarRecentBlockHashesPubkey,
		nonceAuthPubKey,
	).Build())
	instructions = append(instructions, solrock.Wrap(
		zenbtc_spl_token.WrapArgs{
			Value: amount,
			Fee:   params.Fee,
		},
		programID,
		mintKey,
		signerKey,
		feeKey,
		feeWalletAta,
		recipientKey,
		receiverAta,
	))

	tx, err := solana.NewTransaction(
		instructions,
		solana.Hash(nonce.Nonce),
		solana.TransactionPayer(signerKey),
	)
	if err != nil {
		return "", err
	}

	return tx.String(), nil
}

func (k Keeper) getSolanaDurableNonce(ctx sdk.Context) (system.NonceAccount, error) {
	var data []byte
	memStore := k.memStoreService.OpenMemoryStore(ctx)
	data, err := memStore.Get([]byte(durableNonceKey))
	if err != nil {
		return system.NonceAccount{}, err
	}
	params := k.GetParams(ctx).Solana
	acc, err := solana.PublicKeyFromBase58(params.NonceAccountPubKey)
	if err != nil {
		return system.NonceAccount{}, err
	}

	client := rpc.New(params.RpcUrl)
	accountInfo, err := client.GetAccountInfoWithOpts(
		ctx,
		acc,
		&rpc.GetAccountInfoOpts{
			Commitment: rpc.CommitmentConfirmed,
			DataSlice:  nil,
		},
	)
	if err != nil {
		return system.NonceAccount{}, err
	}

	data = accountInfo.Value.Data.GetBinary()

	nonceAccount := new(system.NonceAccount)
	decoder := bin.NewBorshDecoder(data)

	err = nonceAccount.UnmarshalWithDecoder(decoder)
	if err != nil {
		return system.NonceAccount{}, err
	}

	return *nonceAccount, nil

}

package keeper

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v6/contracts/solrock"
	"github.com/Zenrock-Foundation/zrchain/v6/contracts/solrock/generated/rock_spl_token"
	treasuryTypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gagliardetto/solana-go"
	ata "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/system"
)

func (k Keeper) PrepareSolRockMintTx(goCtx context.Context, amount uint64, recipient string, nonce *system.NonceAccount, fundReceiver bool) ([]byte, error) {
	params := k.GetParams(goCtx).Solana

	ctx := sdk.UnwrapSDKContext(goCtx)
	programID, err := solana.PublicKeyFromBase58(params.ProgramId)
	if err != nil {
		return nil, err
	}

	nonceAccKey, err := k.treasuryKeeper.GetKey(ctx, params.NonceAccountKey)
	if err != nil {
		return nil, err
	}

	nonceAccPubKey, err := treasuryTypes.SolanaPubkey(nonceAccKey)
	if err != nil {
		return nil, err
	}

	nonceAuthKey, err := k.treasuryKeeper.GetKey(ctx, params.NonceAuthorityKey)
	if err != nil {
		return nil, err
	}
	nonceAuthPubKey, err := treasuryTypes.SolanaPubkey(nonceAuthKey)
	if err != nil {
		return nil, err
	}

	signerKey, err := k.treasuryKeeper.GetKey(ctx, params.SignerKeyId)
	if err != nil {
		return nil, err
	}
	signerPubKey, err := treasuryTypes.SolanaPubkey(signerKey)
	if err != nil {
		return nil, err
	}

	mintKey, err := solana.PublicKeyFromBase58(params.MintAddress)
	if err != nil {
		return nil, err
	}

	feeKey, err := solana.PublicKeyFromBase58(params.FeeWallet)
	if err != nil {
		return nil, err
	}

	recipientPubKey, err := solana.PublicKeyFromBase58(recipient)
	if err != nil {
		return nil, err
	}

	var instructions []solana.Instruction

	instructions = append(instructions, system.NewAdvanceNonceAccountInstruction(
		*nonceAccPubKey,
		solana.SysVarRecentBlockHashesPubkey,
		*nonceAuthPubKey,
	).Build())

	feeWalletAta, _, err := solana.FindAssociatedTokenAddress(feeKey, mintKey)
	if err != nil {
		return nil, err
	}

	receiverAta, _, err := solana.FindAssociatedTokenAddress(recipientPubKey, mintKey)
	if err != nil {
		return nil, err
	}

	if fundReceiver {
		instructions = append(
			instructions,
			ata.NewCreateInstruction(
				*signerPubKey,
				recipientPubKey,
				mintKey,
			).Build(),
		)
	}

	instructions = append(instructions, solrock.Wrap(
		programID,
		rock_spl_token.WrapArgs{
			Value: amount,
			Fee:   params.Fee,
		},
		*signerPubKey,
		mintKey,
		feeKey,
		feeWalletAta,
		recipientPubKey,
		receiverAta,
	))

	tx, err := solana.NewTransaction(
		instructions,
		solana.Hash(nonce.Nonce),
		solana.TransactionPayer(*signerPubKey),
	)
	if err != nil {
		return nil, err
	}
	txBytes, err := tx.Message.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return txBytes, nil
}

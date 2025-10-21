package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/shared"
	dcttypes "github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
)

func (k msgServer) AdvanceSolanaNonce(ctx context.Context, msg *types.MsgAdvanceSolanaNonce) (*types.MsgAdvanceSolanaNonceResponse, error) {
	if msg.Authority != shared.AdminAuthAddr && msg.Authority != k.authority {
		return nil, fmt.Errorf("invalid authority; expected %s or %s, got %s", shared.AdminAuthAddr, k.authority, msg.Authority)
	}

	if msg.RecentBlockhash == "" {
		return nil, fmt.Errorf("recent blockhash must be provided")
	}
	if msg.Caip2ChainId == "" {
		return nil, fmt.Errorf("caip2_chain_id must be provided")
	}

	hash, err := solana.HashFromBase58(msg.RecentBlockhash)
	if err != nil {
		return nil, fmt.Errorf("invalid recent blockhash: %w", err)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	var (
		nonceAccountKeyID   uint64
		nonceAuthorityKeyID uint64
	)

	switch {
	case msg.Zenbtc:
		solParams := k.zenBTCKeeper.GetSolanaParams(sdkCtx)
		nonceAccountKeyID = solParams.NonceAccountKey
		nonceAuthorityKeyID = solParams.NonceAuthorityKey
	case msg.Asset != dcttypes.Asset_ASSET_UNSPECIFIED:
		if k.dctKeeper == nil {
			return nil, fmt.Errorf("dct keeper is not configured")
		}
		solParams, err := k.dctKeeper.GetSolanaParams(sdkCtx, msg.Asset)
		if err != nil {
			if errors.Is(err, collections.ErrNotFound) {
				return nil, fmt.Errorf("solana params not found for asset %s", msg.Asset.String())
			}
			return nil, fmt.Errorf("failed to get solana params for asset %s: %w", msg.Asset.String(), err)
		}
		if solParams == nil {
			return nil, fmt.Errorf("solana params not configured for asset %s", msg.Asset.String())
		}
		nonceAccountKeyID = solParams.NonceAccountKey
		nonceAuthorityKeyID = solParams.NonceAuthorityKey
	default:
		return nil, fmt.Errorf("either zenbtc must be true or a DCT asset must be provided")
	}

	nonceAccountKey, err := k.treasuryKeeper.GetKey(sdkCtx, nonceAccountKeyID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve nonce account key %d: %w", nonceAccountKeyID, err)
	}
	nonceAuthorityKey, err := k.treasuryKeeper.GetKey(sdkCtx, nonceAuthorityKeyID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve nonce authority key %d: %w", nonceAuthorityKeyID, err)
	}

	nonceAccountPubKey, err := treasurytypes.SolanaPubkey(nonceAccountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode nonce account pubkey: %w", err)
	}
	nonceAuthorityPubKey, err := treasurytypes.SolanaPubkey(nonceAuthorityKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode nonce authority pubkey: %w", err)
	}

	instruction := system.NewAdvanceNonceAccountInstruction(
		*nonceAccountPubKey,
		solana.SysVarRecentBlockHashesPubkey,
		*nonceAuthorityPubKey,
	).Build()

	tx, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		hash,
		solana.TransactionPayer(*nonceAuthorityPubKey),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build advance nonce transaction: %w", err)
	}

	unsignedTx, err := tx.Message.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal advance nonce transaction message: %w", err)
	}

	if _, err := k.submitSolanaTransaction(
		sdkCtx,
		msg.Authority,
		[]uint64{nonceAuthorityKeyID},
		treasurytypes.WalletType_WALLET_TYPE_SOLANA,
		msg.Caip2ChainId,
		unsignedTx,
	); err != nil {
		return nil, fmt.Errorf("failed to submit advance nonce transaction: %w", err)
	}

	return &types.MsgAdvanceSolanaNonceResponse{}, nil
}

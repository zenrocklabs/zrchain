package keeper

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/app/params"
	sidecarapitypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	zentptypes "github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	solana "github.com/gagliardetto/solana-go"
	solToken "github.com/gagliardetto/solana-go/programs/token"
)

// =========================
// ZenTP flow logic (ROCK on Solana)
// =========================

// processSolanaROCKMints processes pending mint transactions for ROCK on Solana.
func (k *Keeper) processSolanaROCKMints(ctx sdk.Context, oracleData OracleData) {
	processTransaction(
		k,
		ctx,
		k.zentpKeeper.GetSolanaParams(ctx).NonceAccountKey,
		nil,
		oracleData.SolanaMintNonces[k.zentpKeeper.GetSolanaParams(ctx).NonceAccountKey],
		func(ctx sdk.Context) ([]*zentptypes.Bridge, error) {
			return k.zentpKeeper.GetMintsWithStatusPending(ctx)
		},
		func(tx *zentptypes.Bridge) error {
			if tx.BlockHeight > 0 {
				return nil
			}
			requiredFields := []VoteExtensionField{VEFieldSolanaAccountsHash, VEFieldSolanaMintNoncesHash}
			if err := k.validateConsensusForTxFields(ctx, oracleData, requiredFields,
				"solROCK mint", fmt.Sprintf("tx_id: %d, recipient: %s, amount: %d", tx.Id, tx.RecipientAddress, tx.Amount)); err != nil {
				return fmt.Errorf("validateConsensusForTxFields: %w", err)
			}
			solParams := k.zentpKeeper.GetSolanaParams(ctx)
			// Derive the ATA for the recipient and check its state
			recipientPubKey, err := solana.PublicKeyFromBase58(tx.RecipientAddress)
			if err != nil {
				return fmt.Errorf("invalid recipient address %s for ZenTP mint: %w", tx.RecipientAddress, err)
			}
			mintPubKey, err := solana.PublicKeyFromBase58(solParams.MintAddress)
			if err != nil {
				return fmt.Errorf("invalid ZenTP mint address %s: %w", solParams.MintAddress, err)
			}
			expectedATA, _, err := solana.FindAssociatedTokenAddress(recipientPubKey, mintPubKey)
			if err != nil {
				return fmt.Errorf("failed to derive ATA for ZenTP recipient %s, mint %s: %w", tx.RecipientAddress, solParams.MintAddress, err)
			}
			fundReceiver := false
			if ata, ok := oracleData.SolanaAccounts[expectedATA.String()]; !ok || ata.State == solToken.Uninitialized {
				fundReceiver = true
			}
			nonce, ok := oracleData.SolanaMintNonces[solParams.NonceAccountKey]
			if !ok {
				return fmt.Errorf("nonce not found in oracleData.SolanaMintNonces for solParams.NonceAccountKey: %d", solParams.NonceAccountKey)
			}
			transaction, err := k.PrepareSolanaMintTx(ctx, &solanaMintTxRequest{
				amount:            tx.Amount,
				fee:               min(solParams.Fee, tx.Amount),
				recipient:         tx.RecipientAddress,
				nonce:             nonce,
				fundReceiver:      fundReceiver,
				programID:         solParams.ProgramId,
				mintAddress:       solParams.MintAddress,
				feeWallet:         solParams.FeeWallet,
				nonceAccountKey:   solParams.NonceAccountKey,
				nonceAuthorityKey: solParams.NonceAuthorityKey,
				signerKey:         solParams.SignerKeyId,
				eventID:           tx.Id,
				rock:              true,
			})
			if err != nil {
				return fmt.Errorf("PrepareSolRockMintTx: %w", err)
			}
			id, err := k.submitSolanaTransaction(ctx, tx.Creator, []uint64{solParams.SignerKeyId, solParams.NonceAuthorityKey}, treasurytypes.WalletType_WALLET_TYPE_SOLANA, tx.DestinationChain, transaction)
			if err != nil {
				return fmt.Errorf("submitSolanaTransaction: %w", err)
			}
			tx.State = zentptypes.BridgeStatus_BRIDGE_STATUS_PENDING
			tx.TxId = id
			tx.BlockHeight = ctx.BlockHeight()
			if err := k.LastUsedSolanaNonce.Set(ctx, solParams.NonceAccountKey, types.SolanaNonce{Nonce: nonce.Nonce[:]}); err != nil {
				return fmt.Errorf("LastUsedSolanaNonce.Set: %w", err)
			}
			return k.zentpKeeper.UpdateMint(ctx, tx.Id, tx)
		},
		func(tx *zentptypes.Bridge) error {
			if !fieldHasConsensus(oracleData.FieldVotePowers, VEFieldSolanaMintEventsHash) {
				return nil
			}
			if tx.BlockHeight == 0 {
				return nil
			}
			solParams := k.zentpKeeper.GetSolanaParams(ctx)
			if ctx.BlockHeight() > tx.BlockHeight+solParams.Btl {
				*tx = k.processBtlSolanaROCKMint(ctx, *tx, oracleData, *solParams)
			}
			if tx.AwaitingEventSince > 0 {
				*tx = k.processSecondaryTimeoutSolanaROCKMint(ctx, *tx, oracleData, *solParams)
			}
			return k.zentpKeeper.UpdateMint(ctx, tx.Id, tx)
		},
	)
}

// processSolanaROCKMintEvents finalizes pending ROCK mints based on oracle events and supply invariants.
func (k *Keeper) processSolanaROCKMintEvents(ctx sdk.Context, oracleData OracleData) {
	pendingMints, err := k.zentpKeeper.GetMintsWithStatusPending(ctx)
	if err != nil || len(pendingMints) == 0 {
		return
	}
	for _, pendingMint := range pendingMints {
		tx, err := k.treasuryKeeper.GetSignTransactionRequest(ctx, pendingMint.TxId)
		if err != nil {
			return
		}
		sigReq, err := k.treasuryKeeper.GetSignRequest(ctx, tx.SignRequestId)
		if err != nil {
			continue
		}
		var signatures []byte
		for _, id := range sigReq.ChildReqIds {
			childReq, err := k.treasuryKeeper.GetSignRequest(ctx, id)
			if err != nil || len(childReq.SignedData) != 1 {
				continue
			}
			signatures = append(signatures, childReq.SignedData[0].SignedData...)
		}
		sigHash := sha256.Sum256(signatures)
		for _, event := range oracleData.SolanaMintEvents {
			if bytes.Equal(event.SigHash, sigHash[:]) {
				if err := k.zentpKeeper.CheckROCKSupplyCap(ctx, sdkmath.ZeroInt()); err != nil {
					pendingMint.State = zentptypes.BridgeStatus_BRIDGE_STATUS_FAILED
					_ = k.zentpKeeper.UpdateMint(ctx, pendingMint.Id, pendingMint)
					continue
				}
				totalSupplyBefore, err := k.zentpKeeper.GetTotalROCKSupply(ctx)
				if err != nil {
					continue
				}
				if err := k.bankKeeper.BurnCoins(ctx, zentptypes.ModuleName, sdk.NewCoins(sdk.NewCoin(pendingMint.Denom, sdkmath.NewIntFromUint64(pendingMint.Amount)))); err != nil {
					continue
				}
				solanaSupply, err := k.zentpKeeper.GetSolanaROCKSupply(ctx)
				if err != nil {
					continue
				}
				if err := k.zentpKeeper.SetSolanaROCKSupply(ctx, solanaSupply.Add(sdkmath.NewIntFromUint64(pendingMint.Amount))); err != nil {
					continue
				}
				_ = k.LastCompletedZentpMintID.Set(ctx, pendingMint.Id)
				totalSupplyAfter, err := k.zentpKeeper.GetTotalROCKSupply(ctx)
				if err != nil || !totalSupplyBefore.Equal(totalSupplyAfter) {
					pendingMint.State = zentptypes.BridgeStatus_BRIDGE_STATUS_FAILED
					_ = k.zentpKeeper.UpdateMint(ctx, pendingMint.Id, pendingMint)
					continue
				}
				pendingMint.State = zentptypes.BridgeStatus_BRIDGE_STATUS_COMPLETED
				_ = k.zentpKeeper.UpdateMint(ctx, pendingMint.Id, pendingMint)
			}
		}
	}
}

// processSolanaROCKBurnEvents handles ROCK burns on Solana by minting ROCK on zrchain and crediting recipients.
func (k Keeper) processSolanaROCKBurnEvents(ctx sdk.Context, oracleData OracleData) {
	var toProcess []*sidecarapitypes.BurnEvent
	processedInThisRun := make(map[string]bool)
	for _, e := range oracleData.SolanaBurnEvents {
		if e.IsZenBTC {
			continue
		}
		if _, ok := processedInThisRun[e.TxID]; ok {
			continue
		}
		addr, err := sdk.Bech32ifyAddressBytes("zen", e.DestinationAddr[:20])
		if err != nil {
			continue
		}
		burns, err := k.zentpKeeper.GetBurns(ctx, addr, e.ChainID, e.TxID)
		if err != nil || len(burns) > 0 {
			continue
		}
		toProcess = append(toProcess, &e)
		processedInThisRun[e.TxID] = true
	}
	if len(toProcess) == 0 {
		return
	}
	processedTxHashes := make(map[string]bool)
	for _, burn := range toProcess {
		addr, err := sdk.Bech32ifyAddressBytes("zen", burn.DestinationAddr[:20])
		if err != nil {
			continue
		}
		accAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			continue
		}
		if err := k.zentpKeeper.CheckCanBurnFromSolana(ctx, sdkmath.NewIntFromUint64(burn.Amount)); err != nil {
			continue
		}
		_, bridgeFee, err := k.zentpKeeper.GetBridgeFeeParams(ctx)
		if err != nil {
			continue
		}
		bridgeFeeCoins, err := k.zentpKeeper.GetBridgeFeeAmount(ctx, burn.Amount, bridgeFee)
		if err != nil {
			continue
		}
		bridgeAmount := sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(burn.Amount).Sub(bridgeFeeCoins.AmountOf(params.BondDenom))))
		coins := sdk.NewCoins(sdk.NewCoin(params.BondDenom, sdkmath.NewIntFromUint64(burn.Amount)))
		if err := k.bankKeeper.MintCoins(ctx, zentptypes.ModuleName, coins); err != nil {
			continue
		}
		solanaSupply, err := k.zentpKeeper.GetSolanaROCKSupply(ctx)
		if err != nil {
			continue
		}
		newSolanaSupply := solanaSupply.Sub(sdkmath.NewIntFromUint64(burn.Amount))
		if newSolanaSupply.IsNegative() || k.zentpKeeper.SetSolanaROCKSupply(ctx, newSolanaSupply) != nil {
			continue
		}
		_ = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, zentptypes.ModuleName, accAddr, bridgeAmount)
		if bridgeFeeCoins.AmountOf(params.BondDenom).IsPositive() {
			_ = k.bankKeeper.SendCoinsFromModuleToModule(ctx, zentptypes.ModuleName, zentptypes.ZentpCollectorName, bridgeFeeCoins)
		}
		_ = k.zentpKeeper.AddBurn(ctx, &zentptypes.Bridge{Denom: params.BondDenom, Amount: burn.Amount, RecipientAddress: accAddr.String(), SourceChain: burn.ChainID, TxHash: burn.TxID, State: zentptypes.BridgeStatus_BRIDGE_STATUS_COMPLETED, BlockHeight: ctx.BlockHeight()})
		processedTxHashes[burn.TxID] = true
		// Emit event
		sdk.UnwrapSDKContext(ctx).EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeValidation,
				sdk.NewAttribute(types.AttributeKeyBridgeAmount, fmt.Sprintf("%d", burn.Amount)),
				sdk.NewAttribute(types.AttributeKeyBridgeFee, bridgeFeeCoins.AmountOf(params.BondDenom).String()),
				sdk.NewAttribute(types.AttributeKeyBurnDestination, addr),
			),
		)
	}
	k.ClearProcessedBackfillRequests(ctx, types.EventType_EVENT_TYPE_ZENTP_BURN, processedTxHashes)
}

package keeper

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"

	"cosmossdk.io/collections"
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
	processSolanaTxQueue(k, ctx, SolanaTxQueueArgs[*zentptypes.Bridge]{
		NonceAccountKey:          k.zentpKeeper.GetSolanaParams(ctx).NonceAccountKey,
		NonceAccount:             oracleData.SolanaMintNonces[k.zentpKeeper.GetSolanaParams(ctx).NonceAccountKey],
		DispatchRequestedChecker: TxDispatchRequestHandler[uint64]{Store: k.SolanaNonceRequested},
		GetPendingTxs: func(ctx sdk.Context) ([]*zentptypes.Bridge, error) {
			return k.zentpKeeper.GetMintsWithStatusPending(ctx)
		},
		DispatchTx: func(tx *zentptypes.Bridge) error {
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

			// Get current Solana counters for ROCK from chain state
			asset := types.Asset_ROCK
			assetKey := asset.String()
			counters, err := k.SolanaCounters.Get(ctx, assetKey)
			if err != nil {
				if errors.Is(err, collections.ErrNotFound) {
					// Initialize counters if not found
					counters = types.SolanaCounters{MintCounter: 0, RedemptionCounter: 0}
				} else {
					return fmt.Errorf("failed to get Solana counters for ROCK: %w", err)
				}
			}
			nextMintCounter := counters.MintCounter + 1
			k.Logger(ctx).Info("Read Solana mint counter from chain state",
				"asset", asset.String(),
				"current_mint_counter", counters.MintCounter,
				"next_mint_counter", nextMintCounter,
			)

			transaction, err := k.PrepareSolanaMintTx(ctx, &solanaMintTxRequest{
				amount:              tx.Amount,
				fee:                 min(solParams.Fee, tx.Amount),
				recipient:           tx.RecipientAddress,
				nonce:               nonce,
				fundReceiver:        fundReceiver,
				programID:           solParams.ProgramId,
				mintAddress:         solParams.MintAddress,
				feeWallet:           solParams.FeeWallet,
				nonceAccountKey:     solParams.NonceAccountKey,
				nonceAuthorityKey:   solParams.NonceAuthorityKey,
				signerKey:           solParams.SignerKeyId,
				multisigKey:         solParams.MultisigKeyAddress,
				eventStoreProgramID: solParams.EventStoreProgramId,
				mintCounter:         nextMintCounter,
				assetName:           asset.String(),
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

			k.Logger(ctx).Info("Dispatched Solana ROCK mint transaction",
				"tx_id", tx.Id,
				"zrchain_tx_id", id,
				"mint_counter_used", nextMintCounter,
			)

			// Counter will be incremented in processSolanaROCKMintEvents when mint is confirmed successful
			if err := k.LastUsedSolanaNonce.Set(ctx, solParams.NonceAccountKey, types.SolanaNonce{Nonce: nonce.Nonce[:]}); err != nil {
				return fmt.Errorf("LastUsedSolanaNonce.Set: %w", err)
			}
			return k.zentpKeeper.UpdateMint(ctx, tx.Id, tx)
		},
		UpdatePendingTxStatus: func(tx *zentptypes.Bridge) error {
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
	})
}

// processSolanaROCKMintEvents finalizes pending ROCK mints based on oracle events and supply invariants.
func (k *Keeper) processSolanaROCKMintEvents(ctx sdk.Context, oracleData OracleData) {
	k.Logger(ctx).Info("processSolanaROCKMintEvents: Started", "oracle_event_count", len(oracleData.SolanaMintEvents))

	pendingMints, err := k.zentpKeeper.GetMintsWithStatusPending(ctx)
	if err != nil || len(pendingMints) == 0 {
		return
	}

	// Process only the first pending mint (FIFO order)
	pendingMint := pendingMints[0]
	if pendingMint.TxId == 0 {
		k.Logger(ctx).Info("processSolanaROCKMintEvents: pending mint missing zrchain tx id", "tx_id", pendingMint.Id)
		return
	}

	signTxReq, err := k.treasuryKeeper.GetSignTransactionRequest(ctx, pendingMint.TxId)
	if err != nil {
		k.Logger(ctx).Error("processSolanaROCKMintEvents: failed to fetch sign transaction request", "tx_id", pendingMint.TxId, "error", err)
		return
	}

	mainSignReq, err := k.treasuryKeeper.GetSignRequest(ctx, signTxReq.SignRequestId)
	if err != nil {
		k.Logger(ctx).Error("processSolanaROCKMintEvents: failed to fetch sign request", "sign_req_id", signTxReq.SignRequestId, "error", err)
		return
	}

	var signatures [][]byte
	for _, childReqID := range mainSignReq.ChildReqIds {
		childReq, err := k.treasuryKeeper.GetSignRequest(ctx, childReqID)
		if err != nil || len(childReq.SignedData) == 0 || len(childReq.SignedData[0].SignedData) == 0 {
			k.Logger(ctx).Warn("processSolanaROCKMintEvents: missing signatures for sign request", "child_req", childReqID)
			return
		}
		signatures = append(signatures, childReq.SignedData[0].SignedData)
	}
	if len(signatures) == 0 {
		return
	}

	// Get expected event ID from chain state
	asset := types.Asset_ROCK
	assetKey := asset.String()
	counters, err := k.SolanaCounters.Get(ctx, assetKey)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			counters = types.SolanaCounters{MintCounter: 0, RedemptionCounter: 0}
		} else {
			k.Logger(ctx).Error("processSolanaROCKMintEvents: failed to get Solana counters", "asset", assetKey, "error", err)
			return
		}
	}
	expectedEventID := new(big.Int).SetUint64(counters.MintCounter + 1)

	// Find matching event with correct event ID
	var matchedEvent *sidecarapitypes.SolanaMintEvent
	for _, event := range oracleData.SolanaMintEvents {
		// Check if this is a ROCK event (no coin field means ROCK - legacy behavior)
		if event.Coint != sidecarapitypes.Coin_UNSPECIFIED {
			continue
		}

		eventHash := base64.StdEncoding.EncodeToString(event.SigHash)
		eventKey := collections.Join(assetKey, eventHash)

		// Check if already processed
		if alreadyProcessed, err := k.ProcessedSolanaMintEvents.Get(ctx, eventKey); err == nil && alreadyProcessed {
			k.Logger(ctx).Warn("processSolanaROCKMintEvents: Solana event already processed", "tx_id", pendingMint.Id, "event_hash", eventHash)
			continue
		} else if err != nil && !errors.Is(err, collections.ErrNotFound) {
			k.Logger(ctx).Error("processSolanaROCKMintEvents: failed to read processed event map", "tx_id", pendingMint.Id, "event_hash", eventHash, "error", err)
			continue
		}

		// Extract and validate event ID
		eventID, err := eventIDFromSolanaMintEvent(event)
		if err != nil {
			k.Logger(ctx).Error("processSolanaROCKMintEvents: failed to parse event ID", "tx_sig", event.TxSig, "error", err)
			continue
		}
		if eventID.Cmp(expectedEventID) != 0 {
			k.Logger(ctx).Info("processSolanaROCKMintEvents: skipping event with unexpected ID",
				"tx_id", pendingMint.Id,
				"expected_event_id", expectedEventID.String(),
				"event_id", eventID.String(),
				"tx_sig", event.TxSig,
			)
			continue
		}

		evtCopy := event
		matchedEvent = &evtCopy
		break
	}

	if matchedEvent == nil {
		k.Logger(ctx).Info("processSolanaROCKMintEvents: no matching Solana mint event yet", "tx_id", pendingMint.Id)
		return
	}

	eventHash := base64.StdEncoding.EncodeToString(matchedEvent.SigHash)
	eventKey := collections.Join(assetKey, eventHash)

	// Perform supply cap check
	if err := k.zentpKeeper.CheckROCKSupplyCap(ctx, sdkmath.ZeroInt()); err != nil {
		pendingMint.State = zentptypes.BridgeStatus_BRIDGE_STATUS_FAILED
		if err := k.zentpKeeper.UpdateMint(ctx, pendingMint.Id, pendingMint); err != nil {
			k.Logger(ctx).Error("error marking solROCK mint as failed (cap)", "id", pendingMint.Id, "error", err)
		}
		return
	}

	totalSupplyBefore, err := k.zentpKeeper.GetTotalROCKSupply(ctx)
	if err != nil {
		k.Logger(ctx).Error("processSolanaROCKMintEvents: failed to get total supply before", "error", err)
		return
	}

	// Burn coins from zentp module
	if err := k.bankKeeper.BurnCoins(ctx, zentptypes.ModuleName, sdk.NewCoins(sdk.NewCoin(pendingMint.Denom, sdkmath.NewIntFromUint64(pendingMint.Amount)))); err != nil {
		k.Logger(ctx).Error("processSolanaROCKMintEvents: failed to burn coins", "error", err)
		return
	}

	// Update Solana supply
	solanaSupply, err := k.zentpKeeper.GetSolanaROCKSupply(ctx)
	if err != nil {
		k.Logger(ctx).Error("processSolanaROCKMintEvents: failed to get Solana supply", "error", err)
		return
	}
	if err := k.zentpKeeper.SetSolanaROCKSupply(ctx, solanaSupply.Add(sdkmath.NewIntFromUint64(pendingMint.Amount))); err != nil {
		k.Logger(ctx).Error("processSolanaROCKMintEvents: failed to set Solana supply", "error", err)
		return
	}

	if err := k.LastCompletedZentpMintID.Set(ctx, pendingMint.Id); err != nil {
		k.Logger(ctx).Error("error setting last completed zentp mint id", "id", pendingMint.Id, "error", err)
	}

	// Verify supply invariant
	totalSupplyAfter, err := k.zentpKeeper.GetTotalROCKSupply(ctx)
	if err != nil || !totalSupplyBefore.Equal(totalSupplyAfter) {
		pendingMint.State = zentptypes.BridgeStatus_BRIDGE_STATUS_FAILED
		if err := k.zentpKeeper.UpdateMint(ctx, pendingMint.Id, pendingMint); err != nil {
			k.Logger(ctx).Error("error marking solROCK mint failed (supply invariant)", "id", pendingMint.Id, "error", err)
		}
		return
	}

	// Mark as completed
	pendingMint.State = zentptypes.BridgeStatus_BRIDGE_STATUS_COMPLETED
	if err := k.zentpKeeper.UpdateMint(ctx, pendingMint.Id, pendingMint); err != nil {
		k.Logger(ctx).Error("error marking solROCK mint completed", "id", pendingMint.Id, "error", err)
		return
	}

	// Increment mint counter after confirmed successful mint
	counters.MintCounter++
	if err := k.SolanaCounters.Set(ctx, assetKey, counters); err != nil {
		k.Logger(ctx).Error("failed to increment Solana mint counter after successful mint", "asset", assetKey, "error", err)
	} else {
		k.Logger(ctx).Info("Incremented Solana mint counter after confirmed successful mint",
			"asset", assetKey,
			"new_mint_counter", counters.MintCounter,
			"tx_id", pendingMint.Id,
		)
	}

	// Mark event as processed
	if err := k.ProcessedSolanaMintEvents.Set(ctx, eventKey, true); err != nil {
		k.Logger(ctx).Error("failed to record processed Solana mint event", "asset", assetKey, "event_hash", eventHash, "error", err)
	}

	k.Logger(ctx).Info("completed ROCK Solana mint", "tx_id", pendingMint.Id, "amount", pendingMint.Amount)
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
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, zentptypes.ModuleName, accAddr, bridgeAmount); err != nil {
			k.Logger(ctx).Error("error sending bridged coins to account", "addr", accAddr.String(), "error", err)
			continue
		}
		if bridgeFeeCoins.AmountOf(params.BondDenom).IsPositive() {
			if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, zentptypes.ModuleName, zentptypes.ZentpCollectorName, bridgeFeeCoins); err != nil {
				k.Logger(ctx).Error("error sending bridge fee to collector", "error", err)
			}
		}
		if err := k.zentpKeeper.AddBurn(ctx, &zentptypes.Bridge{Denom: params.BondDenom, Amount: burn.Amount, RecipientAddress: accAddr.String(), SourceChain: burn.ChainID, TxHash: burn.TxID, State: zentptypes.BridgeStatus_BRIDGE_STATUS_COMPLETED, BlockHeight: ctx.BlockHeight()}); err != nil {
			k.Logger(ctx).Error("error adding burn record", "tx", burn.TxID, "error", err)
		}
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

package keeper

import (
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	sdkmath "cosmossdk.io/math"
	sidecarapitypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	dcttypes "github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gagliardetto/solana-go"
	solToken "github.com/gagliardetto/solana-go/programs/token"
)

func coinToDCTAsset(coin sidecarapitypes.Coin) (dcttypes.Asset, bool) {
	switch coin {
	case sidecarapitypes.Coin_ZENZEC:
		return dcttypes.Asset_ASSET_ZENZEC, true
	default:
		return dcttypes.Asset_ASSET_UNSPECIFIED, false
	}
}

func dctAssetToCoin(asset dcttypes.Asset) (sidecarapitypes.Coin, bool) {
	switch asset {
	case dcttypes.Asset_ASSET_ZENZEC:
		return sidecarapitypes.Coin_ZENZEC, true
	default:
		return sidecarapitypes.Coin_UNSPECIFIED, false
	}
}

func (k *Keeper) processDCTMintsSolana(ctx sdk.Context, oracleData OracleData) {
	if k.dctKeeper == nil {
		return
	}

	assets, err := k.dctKeeper.ListSupportedAssets(ctx)
	if err != nil {
		k.Logger(ctx).Error("failed to list DCT assets", "error", err)
		return
	}

	k.Logger(ctx).Info("processDCTMintsSolana: starting DCT Solana mint processing", "assets", fmt.Sprint(assets))

	for _, asset := range assets {
		k.Logger(ctx).Info("processDCTMintsSolana: fetching DCT Solana params", "asset", asset.String())

		solParams, err := k.dctKeeper.GetSolanaParams(ctx, asset)
		if err != nil {
			k.Logger(ctx).Error("failed to fetch DCT Solana params", "asset", asset.String(), "error", err)
			continue
		}
		if solParams == nil {
			k.Logger(ctx).Info("no Solana params for DCT asset, skipping", "asset", asset.String())
			continue
		}

		nonceAccount := oracleData.SolanaMintNonces[solParams.NonceAccountKey]
		if nonceAccount == nil {
			k.Logger(ctx).Error("missing nonce account for asset", "asset", asset.String(), "nonce_account_key", solParams.NonceAccountKey)
			continue
		}

		processSolanaTxQueue(k, ctx, SolanaTxQueueArgs[dcttypes.PendingMintTransaction]{
			NonceAccountKey: solParams.NonceAccountKey,
			NonceAccount:    nonceAccount,
			DispatchRequestedChecker: TxDispatchRequestHandler[uint64]{
				Store: k.SolanaNonceRequested,
			},
			GetPendingTxs: func(ctx sdk.Context) ([]dcttypes.PendingMintTransaction, error) {
				pending, err := k.getPendingDCTMintTransactions(
					ctx,
					asset,
					dcttypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED,
					dcttypes.WalletType_WALLET_TYPE_SOLANA,
				)
				k.Logger(ctx).Info("processDCTMintsSolana: fetched pending DCT mint transactions", "asset", asset.String(), "pending", fmt.Sprint(pending))
				return pending, err
			},
			DispatchTx: func(tx dcttypes.PendingMintTransaction) error {
				if tx.BlockHeight > 0 {
					k.Logger(ctx).Info("waiting for pending DCT solana mint tx", "asset", asset.String(), "tx_id", tx.Id, "block_height", tx.BlockHeight)
					return nil
				}

				if err := k.dctKeeper.SetFirstPendingSolMintTransaction(ctx, asset, tx.Id); err != nil {
					return err
				}
				k.Logger(ctx).Info("processDCTMintsSolana: set first pending Sol mint transaction",
					"asset", asset.String(),
					"tx_id", tx.Id,
				)

				if len(oracleData.SolanaMintNonces) == 0 {
					k.Logger(ctx).Info("processDCTMintsSolana: no Solana nonces available for asset", "asset", asset.String())
					return fmt.Errorf("no nonce available for DCT solana mint for asset %s", asset.String())
				}

				_, ok := dctAssetToCoin(asset)
				if !ok {
					k.Logger(ctx).Info("processDCTMintsSolana: unsupported DCT asset for Solana dispatch", "asset", asset.String())
					return fmt.Errorf("unsupported DCT asset %s for Solana dispatch", asset.String())
				}

				requiredFields := []VoteExtensionField{VEFieldSolanaMintNoncesHash, VEFieldBTCUSDPrice, VEFieldSolanaAccountsHash}
				if err := k.validateConsensusForTxFields(ctx, oracleData, requiredFields,
					fmt.Sprintf("%s mint", asset.String()),
					fmt.Sprintf("tx_id: %d, recipient: %s, amount: %d", tx.Id, tx.RecipientAddress, tx.Amount)); err != nil {
					return err
				}

				// btcUSDPrice, err := sdkmath.LegacyNewDecFromStr(oracleData.BTCUSDPrice)
				// if err != nil || btcUSDPrice.IsNil() || btcUSDPrice.IsZero() {
				// 	k.Logger(ctx).Error("invalid BTC/USD price for DCT mint", "asset", asset.String(), "error", err)
				// 	return nil
				// }

				// exchangeRate, err := k.dctKeeper.GetExchangeRate(ctx, asset)
				// if err != nil {
				// 	return err
				// }

				// fee := k.CalculateFlatZenBTCMintFee(btcUSDPrice, exchangeRate)
				// fee = min(fee, tx.Amount)

				recipientPubKey, err := solana.PublicKeyFromBase58(tx.RecipientAddress)
				if err != nil {
					return fmt.Errorf("invalid recipient address %s for %s mint: %w", tx.RecipientAddress, asset.String(), err)
				}
				mintPubKey, err := solana.PublicKeyFromBase58(solParams.MintAddress)
				if err != nil {
					return fmt.Errorf("invalid %s mint address %s: %w", asset.String(), solParams.MintAddress, err)
				}
				expectedATA, _, err := solana.FindAssociatedTokenAddress(recipientPubKey, mintPubKey)
				if err != nil {
					return fmt.Errorf("failed to derive ATA for recipient %s, mint %s: %w", tx.RecipientAddress, solParams.MintAddress, err)
				}

				fundReceiver := false
				if ataAccount, ok := oracleData.SolanaAccounts[expectedATA.String()]; !ok || ataAccount.State == solToken.Uninitialized {
					fundReceiver = true
				}

				nonce, ok := oracleData.SolanaMintNonces[solParams.NonceAccountKey]
				if !ok {
					return fmt.Errorf("nonce not found in oracleData.SolanaMintNonces for key: %d", solParams.NonceAccountKey)
				}

				txPrepReq := &solanaMintTxRequest{
					amount:             tx.Amount,
					// fee:               fee,
					fee:                0,
					recipient:          tx.RecipientAddress,
					nonce:              nonce,
					fundReceiver:       fundReceiver,
					programID:          solParams.ProgramId,
					mintAddress:        solParams.MintAddress,
					feeWallet:          solParams.FeeWallet,
					nonceAccountKey:    solParams.NonceAccountKey,
					nonceAuthorityKey:  solParams.NonceAuthorityKey,
					signerKey:          solParams.SignerKeyId,
					multisigKey:        solParams.MultisigKeyAddress,
					zenbtc:             true,
					eventStoreProgramID: solParams.EventStoreProgramId,
				}

				transaction, err := k.PrepareSolanaMintTx(ctx, txPrepReq)
				if err != nil {
					k.Logger(ctx).Info("processDCTMintsSolana: PrepareSolanaMintTx failed",
						"asset", asset.String(),
						"tx_id", tx.Id,
						"error", err,
					)
					return fmt.Errorf("prepareSolanaMintTx (%s): %w", asset.String(), err)
				}
				k.Logger(ctx).Info("processDCTMintsSolana: prepared Solana mint transaction",
					"asset", asset.String(),
					"tx_id", tx.Id,
					"recipient", tx.RecipientAddress,
					"amount", tx.Amount,
					// "fee", fee,
				)

				txID, err := k.submitSolanaTransaction(
					ctx,
					tx.Creator,
					[]uint64{solParams.SignerKeyId, solParams.NonceAuthorityKey},
					treasurytypes.WalletType(tx.ChainType),
					tx.Caip2ChainId,
					transaction,
				)
				if err != nil {
					k.Logger(ctx).Info("processDCTMintsSolana: submitSolanaTransaction failed",
						"asset", asset.String(),
						"tx_id", tx.Id,
						"error", err,
					)
					return err
				}

				tx.ZrchainTxId = txID
				tx.BlockHeight = ctx.BlockHeight()
				if err := k.dctKeeper.SetPendingMintTransaction(ctx, tx); err != nil {
					return err
				}
				k.Logger(ctx).Info("processDCTMintsSolana: dispatched Solana mint transaction",
					"asset", asset.String(),
					"tx_id", tx.Id,
					"zrchain_tx_id", txID,
					"block_height", tx.BlockHeight,
				)

				solNonce := types.SolanaNonce{Nonce: nonce.Nonce[:]}
				return k.LastUsedSolanaNonce.Set(ctx, solParams.NonceAccountKey, solNonce)
			},
			UpdatePendingTxStatus: func(tx dcttypes.PendingMintTransaction) error {
				if !fieldHasConsensus(oracleData.FieldVotePowers, VEFieldSolanaMintEventsHash) {
					k.Logger(ctx).Info("Skipping Solana DCT mint retry/timeout checks – no consensus on SolanaMintEventsHash", "asset", asset.String(), "tx_id", tx.Id)
					return nil
				}
				if tx.BlockHeight == 0 {
					return k.dctKeeper.SetPendingMintTransaction(ctx, tx)
				}
				if ctx.BlockHeight() > tx.BlockHeight+solParams.Btl {
					tx = k.processBtlSolanaDCTMint(ctx, tx, oracleData, *solParams, asset)
				}
				if tx.AwaitingEventSince > 0 {
					tx = k.processSecondaryTimeoutSolanaDCTMint(ctx, tx, oracleData, *solParams, asset)
				}
				return k.dctKeeper.SetPendingMintTransaction(ctx, tx)
			},
		})
	}
}

func (k *Keeper) processSolanaDCTMintEvents(ctx sdk.Context, oracleData OracleData) {
	if k.dctKeeper == nil {
		return
	}

	assets, err := k.dctKeeper.ListSupportedAssets(ctx)
	if err != nil {
		k.Logger(ctx).Error("failed to list DCT assets for mint event processing", "error", err)
		return
	}

	for _, asset := range assets {
		coin, ok := dctAssetToCoin(asset)
		if !ok {
			continue
		}

		solParams, err := k.dctKeeper.GetSolanaParams(ctx, asset)
		if err != nil {
			k.Logger(ctx).Error("failed to fetch DCT Solana params for event processing", "asset", asset.String(), "error", err)
			continue
		}
		if solParams == nil {
			k.Logger(ctx).Info("no Solana params for DCT asset during event processing, skipping", "asset", asset.String())
			continue
		}

		firstPendingID, err := k.dctKeeper.GetFirstPendingSolMintTransaction(ctx, asset)
		if err != nil {
			if !errors.Is(err, collections.ErrNotFound) {
				k.Logger(ctx).Error("failed to fetch first pending DCT mint", "asset", asset.String(), "error", err)
			}
			continue
		}
		if firstPendingID == 0 {
			k.Logger(ctx).Info("processSolanaDCTMintEvents: no pending Solana mint id set", "asset", asset.String())
			continue
		}

		pendingMint, err := k.dctKeeper.GetPendingMintTransaction(ctx, asset, firstPendingID)
		if err != nil {
			k.Logger(ctx).Error("failed to retrieve pending DCT mint transaction", "asset", asset.String(), "id", firstPendingID, "error", err)
			continue
		}
		if pendingMint.ZrchainTxId == 0 {
			k.Logger(ctx).Info("processSolanaDCTMintEvents: pending mint missing zrchain tx id", "asset", asset.String(), "tx_id", pendingMint.Id)
			continue
		}

		signTxReq, err := k.treasuryKeeper.GetSignTransactionRequest(ctx, pendingMint.ZrchainTxId)
		if err != nil {
			k.Logger(ctx).Error("failed to fetch sign transaction request for DCT mint", "asset", asset.String(), "tx_id", pendingMint.ZrchainTxId, "error", err)
			continue
		}

		mainSignReq, err := k.treasuryKeeper.GetSignRequest(ctx, signTxReq.SignRequestId)
		if err != nil {
			k.Logger(ctx).Error("failed to fetch sign request for DCT mint", "asset", asset.String(), "sign_req_id", signTxReq.SignRequestId, "error", err)
			continue
		}

		var signatures [][]byte
		for _, childReqID := range mainSignReq.ChildReqIds {
			childReq, err := k.treasuryKeeper.GetSignRequest(ctx, childReqID)
			if err != nil || len(childReq.SignedData) == 0 || len(childReq.SignedData[0].SignedData) == 0 {
				k.Logger(ctx).Warn("missing signatures for DCT mint sign request", "asset", asset.String(), "child_req", childReqID)
				signatures = nil
				break
			}
			signatures = append(signatures, childReq.SignedData[0].SignedData)
		}
		if len(signatures) == 0 {
			continue
		}

		var matchedEvent *sidecarapitypes.SolanaMintEvent
		for _, event := range oracleData.SolanaMintEvents {
			if event.Coint != coin {
				continue
			}
			if len(event.Recipient) == 0 || len(event.Mint) == 0 {
				continue
			}

			recipientPub := solana.PublicKeyFromBytes(event.Recipient)
			if recipientPub.String() != pendingMint.RecipientAddress {
				continue
			}

			mintPub := solana.PublicKeyFromBytes(event.Mint)
			if mintPub.String() != solParams.MintAddress {
				continue
			}

			if event.Value != pendingMint.Amount {
				continue
			}

			evtCopy := event
			matchedEvent = &evtCopy
			break
		}
		if matchedEvent == nil {
			k.Logger(ctx).Info("processSolanaDCTMintEvents: no matching Solana mint event yet", "asset", asset.String(), "tx_id", pendingMint.Id)
			continue
		}

		supply, err := k.dctKeeper.GetSupply(ctx, asset)
		if err != nil {
			if !errors.Is(err, collections.ErrNotFound) {
				k.Logger(ctx).Error("failed to fetch DCT supply", "asset", asset.String(), "error", err)
				continue
			}
			supply = dcttypes.Supply{Asset: asset}
		}

		if supply.PendingAmount >= pendingMint.Amount {
			supply.PendingAmount -= pendingMint.Amount
		} else {
			supply.PendingAmount = 0
		}
		supply.MintedAmount += pendingMint.Amount
		supply.Asset = asset

		if err := k.dctKeeper.SetSupply(ctx, supply); err != nil {
			k.Logger(ctx).Error("failed to update DCT supply after mint event", "asset", asset.String(), "error", err)
			continue
		}

		pendingMint.Status = dcttypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_MINTED
		pendingMint.TxHash = matchedEvent.TxSig
		if err := k.dctKeeper.SetPendingMintTransaction(ctx, pendingMint); err != nil {
			k.Logger(ctx).Error("failed to update DCT pending mint transaction", "asset", asset.String(), "tx_id", pendingMint.Id, "error", err)
			continue
		}

		k.advanceDCTFirstPendingSolMintTransaction(ctx, asset, pendingMint.Id)

		k.Logger(ctx).Info("completed DCT Solana mint", "asset", asset.String(), "tx_id", pendingMint.Id, "recipient", pendingMint.RecipientAddress, "amount", pendingMint.Amount)
	}
}

func (k *Keeper) storeNewDCTBurnEvents(ctx sdk.Context, oracleData OracleData) {
	if k.dctKeeper == nil {
		return
	}

	processedKeys := make(map[string]bool)

	for _, burn := range oracleData.SolanaBurnEvents {
		if burn.IsZenBTC {
			continue // Skip zenBTC burns, they're handled separately
		}
		asset, ok := coinToDCTAsset(burn.Coin)
		if !ok {
			continue
		}

		key := fmt.Sprintf("%s-%d-%s", burn.TxID, burn.LogIndex, burn.ChainID)
		if processedKeys[key] {
			continue
		}
		processedKeys[key] = true

		if k.dctBurnEventExists(ctx, asset, burn.TxID, burn.LogIndex) {
			continue
		}

		burnEvent := dcttypes.BurnEvent{
			TxID:            burn.TxID,
			LogIndex:        burn.LogIndex,
			ChainID:         burn.ChainID,
			DestinationAddr: burn.DestinationAddr,
			Amount:          burn.Amount,
			Asset:           asset,
			Status:          dcttypes.BurnStatus_BURN_STATUS_UNSTAKING,
			MaturityHeight:  k.calculateRedemptionMaturityHeight(ctx),
		}

		if _, createErr := k.dctKeeper.CreateBurnEvent(ctx, asset, &burnEvent); createErr != nil {
			k.Logger(ctx).Error("failed to create DCT burn event", "asset", asset.String(), "txID", burn.TxID, "logIndex", burn.LogIndex, "error", createErr)
			continue
		}

		k.Logger(ctx).Info("recorded DCT burn awaiting maturity", "asset", asset.String(), "burn_id", burnEvent.Id, "burn_tx", burn.TxID, "amount", burn.Amount, "maturity_height", burnEvent.MaturityHeight)
	}
}

func (k *Keeper) dctBurnEventExists(ctx sdk.Context, asset dcttypes.Asset, txID string, logIndex uint64) bool {
	found := false
	_ = k.dctKeeper.WalkBurnEvents(ctx, asset, func(_ uint64, event dcttypes.BurnEvent) (bool, error) {
		if event.TxID == txID && event.LogIndex == logIndex {
			found = true
			return true, nil // stop iteration
		}
		return false, nil
	})
	return found
}

func (k *Keeper) advanceDCTFirstPendingSolMintTransaction(ctx sdk.Context, asset dcttypes.Asset, currentID uint64) {
	nextID := currentID + 1
	for {
		tx, err := k.dctKeeper.GetPendingMintTransaction(ctx, asset, nextID)
		if err != nil {
			if errors.Is(err, collections.ErrNotFound) {
				_ = k.dctKeeper.SetFirstPendingSolMintTransaction(ctx, asset, 0)
			}
			return
		}
		if tx.Status != dcttypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_MINTED {
			_ = k.dctKeeper.SetFirstPendingSolMintTransaction(ctx, asset, nextID)
			return
		}
		nextID++
	}
}

func (k Keeper) processBtlSolanaDCTMint(ctx sdk.Context, tx dcttypes.PendingMintTransaction, oracleData OracleData, solParams dcttypes.Solana, asset dcttypes.Asset) dcttypes.PendingMintTransaction {
	newBlockHeight, newAwaitingEventSince := k.handleSolanaTransactionBTLTimeout(
		ctx,
		tx.Id,
		tx.BlockHeight,
		tx.AwaitingEventSince,
		solParams.NonceAccountKey,
		solParams.Btl,
		oracleData,
	)
	tx.BlockHeight = newBlockHeight
	tx.AwaitingEventSince = newAwaitingEventSince
	return tx
}

func (k Keeper) processSecondaryTimeoutSolanaDCTMint(ctx sdk.Context, tx dcttypes.PendingMintTransaction, oracleData OracleData, solParams dcttypes.Solana, asset dcttypes.Asset) dcttypes.PendingMintTransaction {
	newBlockHeight, newAwaitingEventSince := k.handleSolanaEventArrivalTimeout(
		ctx,
		tx.Id,
		tx.RecipientAddress,
		tx.Amount,
		tx.AwaitingEventSince,
		tx.BlockHeight,
		solParams.NonceAccountKey,
		oracleData,
	)
	tx.BlockHeight = newBlockHeight
	tx.AwaitingEventSince = newAwaitingEventSince
	return tx
}

// checkForDCTRedemptionFulfilment updates supplies when treasury sign requests for DCT redemptions are fulfilled.
func (k *Keeper) checkForDCTRedemptionFulfilment(ctx sdk.Context) {
	if k.dctKeeper == nil {
		return
	}

	assets, err := k.dctKeeper.ListSupportedAssets(ctx)
	if err != nil {
		k.Logger(ctx).Error("failed to list DCT assets for redemption fulfilment", "error", err)
		return
	}

	for _, asset := range assets {
		startingIndex, _ := k.dctKeeper.GetFirstRedemptionAwaitingSign(ctx, asset)
		redemptions, err := k.GetDCTRedemptionsByStatus(ctx, asset, dcttypes.RedemptionStatus_AWAITING_SIGN, 0, startingIndex)
		if err != nil || len(redemptions) == 0 {
			continue
		}
		if err := k.dctKeeper.SetFirstRedemptionAwaitingSign(ctx, asset, redemptions[0].Data.Id); err != nil {
			k.Logger(ctx).Error("error setting first DCT redemption awaiting sign", "asset", asset.String(), "id", redemptions[0].Data.Id, "error", err)
		}
		supply, err := k.dctKeeper.GetSupply(ctx, asset)
		if err != nil {
			if !errors.Is(err, collections.ErrNotFound) {
				k.Logger(ctx).Error("failed to get DCT supply for redemption fulfilment", "asset", asset.String(), "error", err)
			}
			continue
		}
		for _, redemption := range redemptions {
			signReq, err := k.treasuryKeeper.GetSignRequest(ctx, redemption.Data.SignReqId)
			if err != nil {
				continue
			}
			if signReq.Status == treasurytypes.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING {
				continue
			}
			if signReq.Status == treasurytypes.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED {
				exchangeRate, err := k.dctKeeper.GetExchangeRate(ctx, asset)
				if err != nil {
					continue
				}
				nativeToRelease := uint64(sdkmath.LegacyNewDecFromInt(sdkmath.NewIntFromUint64(redemption.Data.Amount)).Quo(exchangeRate).TruncateInt64())
				if supply.MintedAmount < redemption.Data.Amount || supply.CustodiedAmount < nativeToRelease {
					continue
				}
				supply.MintedAmount -= redemption.Data.Amount
				supply.CustodiedAmount -= nativeToRelease
				redemption.Status = dcttypes.RedemptionStatus_COMPLETED
			}
			if signReq.Status == treasurytypes.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED {
				redemption.Data.SignReqId = 0
				redemption.Status = dcttypes.RedemptionStatus_INITIATED
			}
			if err := k.dctKeeper.SetRedemption(ctx, asset, redemption.Data.Id, redemption); err != nil {
				k.Logger(ctx).Error("error updating DCT redemption after fulfilment", "asset", asset.String(), "id", redemption.Data.Id, "error", err)
			}
		}
		if err := k.dctKeeper.SetSupply(ctx, supply); err != nil {
			k.Logger(ctx).Error("error updating DCT supply after redemption fulfilment", "asset", asset.String(), "error", err)
		}
	}
}

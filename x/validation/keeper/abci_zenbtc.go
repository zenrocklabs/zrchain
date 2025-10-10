package keeper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sidecarapitypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	solana "github.com/gagliardetto/solana-go"
	solToken "github.com/gagliardetto/solana-go/programs/token"
	zenbtctypes "github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

// =========================
// zenBTC flow logic
// =========================

// processZenBTCMintsEthereum processes pending mint transactions on EVM chains.
func (k *Keeper) processZenBTCMintsEthereum(ctx sdk.Context, oracleData OracleData) {
	processEthereumTxQueue(k, ctx, EthereumTxQueueArgs[zenbtctypes.PendingMintTransaction]{
		KeyID:                    k.zenBTCKeeper.GetEthMinterKeyID(ctx),
		RequestedNonce:           oracleData.RequestedEthMinterNonce,
		DispatchRequestedChecker: TxDispatchRequestHandler[uint64]{Store: k.EthereumNonceRequested},
		GetPendingTxs: func(ctx sdk.Context) ([]zenbtctypes.PendingMintTransaction, error) {
			return k.getPendingMintTransactions(
				ctx,
				zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED,
				zenbtctypes.WalletType_WALLET_TYPE_EVM,
			)
		},
		DispatchTx: func(tx zenbtctypes.PendingMintTransaction) error {
			if err := k.zenBTCKeeper.SetFirstPendingEthMintTransaction(ctx, tx.Id); err != nil {
				return err
			}
			requiredFields := []VoteExtensionField{VEFieldRequestedEthMinterNonce, VEFieldBTCUSDPrice}
			if err := k.validateConsensusForTxFields(ctx, oracleData, requiredFields,
				"zenBTC mint", fmt.Sprintf("tx_id: %d, recipient: %s, amount: %d", tx.Id, tx.RecipientAddress, tx.Amount)); err != nil {
				return err
			}
			exchangeRate, err := k.zenBTCKeeper.GetExchangeRate(ctx)
			if err != nil {
				return err
			}
			btcUSDPrice, err := sdkmath.LegacyNewDecFromStr(oracleData.BTCUSDPrice)
			if err != nil || btcUSDPrice.IsNil() || btcUSDPrice.IsZero() {
				k.Logger(ctx).Error("invalid BTC/USD price", "error", err)
				return nil
			}
			feeZenBTC := k.CalculateFlatZenBTCMintFee(btcUSDPrice, exchangeRate)
			feeZenBTC = min(feeZenBTC, tx.Amount)
			chainID, err := types.ValidateEVMChainID(ctx, tx.Caip2ChainId)
			if err != nil {
				return fmt.Errorf("unsupported chain ID: %w", err)
			}
			unsignedMintTxHash, unsignedMintTx, err := k.constructMintTx(
				ctx,
				tx.RecipientAddress,
				chainID,
				tx.Amount,
				feeZenBTC,
				oracleData.RequestedEthMinterNonce,
				oracleData.EthGasLimit,
				oracleData.EthBaseFee,
				oracleData.EthTipCap,
			)
			if err != nil {
				return err
			}
			k.Logger(ctx).Warn("processing zenBTC mint",
				"recipient", tx.RecipientAddress,
				"amount", tx.Amount,
				"nonce", oracleData.RequestedEthMinterNonce,
				"gas_limit", oracleData.EthGasLimit,
				"base_fee", oracleData.EthBaseFee,
				"tip_cap", oracleData.EthTipCap,
			)
			return k.submitEthereumTransaction(
				ctx,
				tx.Creator,
				k.zenBTCKeeper.GetEthMinterKeyID(ctx),
				treasurytypes.WalletType(tx.ChainType),
				chainID,
				unsignedMintTx,
				unsignedMintTxHash,
			)
		},
		OnTxConfirmed: func(tx zenbtctypes.PendingMintTransaction) error {
			supply, err := k.zenBTCKeeper.GetSupply(ctx)
			if err != nil {
				return err
			}
			supply.PendingZenBTC -= tx.Amount
			supply.MintedZenBTC += tx.Amount
			if err := k.zenBTCKeeper.SetSupply(ctx, supply); err != nil {
				return err
			}
			k.Logger(ctx).Warn("pending mint supply updated",
				"pending_mint_old", supply.PendingZenBTC+tx.Amount,
				"pending_mint_new", supply.PendingZenBTC,
			)
			k.Logger(ctx).Warn("minted supply updated",
				"minted_old", supply.MintedZenBTC-tx.Amount,
				"minted_new", supply.MintedZenBTC,
			)
			tx.Status = zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_MINTED
			if err := k.adjustDefaultValidatorBedrockBTC(ctx, sdkmath.NewIntFromUint64(tx.Amount)); err != nil {
				k.Logger(ctx).Error("error adjusting bedrock BTC on mint", "error", err)
			}
			return k.zenBTCKeeper.SetPendingMintTransaction(ctx, tx)
		},
	})
}

// processZenBTCMintsSolana processes pending zenBTC mints on Solana.
func (k *Keeper) processZenBTCMintsSolana(ctx sdk.Context, oracleData OracleData) {
	processSolanaTxQueue(k, ctx, SolanaTxQueueArgs[zenbtctypes.PendingMintTransaction]{
		NonceAccountKey:          k.zenBTCKeeper.GetSolanaParams(ctx).NonceAccountKey,
		NonceAccount:             oracleData.SolanaMintNonces[k.zenBTCKeeper.GetSolanaParams(ctx).NonceAccountKey],
		DispatchRequestedChecker: TxDispatchRequestHandler[uint64]{Store: k.SolanaNonceRequested},
		GetPendingTxs: func(ctx sdk.Context) ([]zenbtctypes.PendingMintTransaction, error) {
			pendingMints, err := k.getPendingMintTransactions(
				ctx,
				zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_DEPOSITED,
				zenbtctypes.WalletType_WALLET_TYPE_SOLANA,
			)
			k.Logger(ctx).Warn("pending zenbtc solana mints", "count", len(pendingMints))
			return pendingMints, err
		},
		DispatchTx: func(tx zenbtctypes.PendingMintTransaction) error {
			if tx.BlockHeight > 0 {
				k.Logger(ctx).Info("waiting for pending zenbtc solana mint tx", "tx_id", tx.Id, "block_height", tx.BlockHeight)
				return nil
			}
			if err := k.zenBTCKeeper.SetFirstPendingSolMintTransaction(ctx, tx.Id); err != nil {
				return err
			}
			if len(oracleData.SolanaMintNonces) == 0 {
				return fmt.Errorf("no nonce available for zenbtc solana mint")
			}
			requiredFields := []VoteExtensionField{VEFieldSolanaMintNoncesHash, VEFieldBTCUSDPrice, VEFieldSolanaAccountsHash}
			if err := k.validateConsensusForTxFields(ctx, oracleData, requiredFields,
				"zenBTC mint", fmt.Sprintf("tx_id: %d, recipient: %s, amount: %d", tx.Id, tx.RecipientAddress, tx.Amount)); err != nil {
				return err
			}
			btcUSDPrice, err := sdkmath.LegacyNewDecFromStr(oracleData.BTCUSDPrice)
			if err != nil || btcUSDPrice.IsNil() || btcUSDPrice.IsZero() {
				k.Logger(ctx).Error("invalid BTC/USD price", "error", err)
				return nil
			}
			exchangeRate, err := k.zenBTCKeeper.GetExchangeRate(ctx)
			if err != nil {
				return err
			}
			feeZenBTC := k.CalculateFlatZenBTCMintFee(btcUSDPrice, exchangeRate)
			feeZenBTC = min(feeZenBTC, tx.Amount)
			solParams := k.zenBTCKeeper.GetSolanaParams(ctx)

			recipientPubKey, err := solana.PublicKeyFromBase58(tx.RecipientAddress)
			if err != nil {
				return fmt.Errorf("invalid recipient address %s for ZenBTC mint: %w", tx.RecipientAddress, err)
			}
			mintPubKey, err := solana.PublicKeyFromBase58(solParams.MintAddress)
			if err != nil {
				return fmt.Errorf("invalid ZenBTC mint address %s: %w", solParams.MintAddress, err)
			}
			expectedATA, _, err := solana.FindAssociatedTokenAddress(recipientPubKey, mintPubKey)
			if err != nil {
				return fmt.Errorf("failed to derive ATA for recipient %s, mint %s: %w", tx.RecipientAddress, solParams.MintAddress, err)
			}
			fundReceiver := false
			if ata, ok := oracleData.SolanaAccounts[expectedATA.String()]; !ok || ata.State == solToken.Uninitialized {
				fundReceiver = true
			}
			nonce, ok := oracleData.SolanaMintNonces[solParams.NonceAccountKey]
			if !ok {
				return fmt.Errorf("nonce not found in oracleData.SolanaMintNonces for nonce account key: %d", solParams.NonceAccountKey)
			}
			txPrepReq := &solanaMintTxRequest{
				amount:            tx.Amount,
				fee:               feeZenBTC,
				recipient:         tx.RecipientAddress,
				nonce:             nonce,
				fundReceiver:      fundReceiver,
				programID:         solParams.ProgramId,
				mintAddress:       solParams.MintAddress,
				feeWallet:         solParams.FeeWallet,
				nonceAccountKey:   solParams.NonceAccountKey,
				nonceAuthorityKey: solParams.NonceAuthorityKey,
				signerKey:         solParams.SignerKeyId,
				zenbtc:            true,
			}
			transaction, err := k.PrepareSolanaMintTx(ctx, txPrepReq)
			if err != nil {
				return fmt.Errorf("PrepareSolRockMintTx: %w", err)
			}
			txID, err := k.submitSolanaTransaction(
				ctx,
				tx.Creator,
				[]uint64{solParams.SignerKeyId, solParams.NonceAuthorityKey},
				treasurytypes.WalletType(tx.ChainType),
				tx.Caip2ChainId,
				transaction,
			)
			if err != nil {
				return err
			}
			tx.ZrchainTxId = txID
			tx.BlockHeight = ctx.BlockHeight()
			if err = k.zenBTCKeeper.SetPendingMintTransaction(ctx, tx); err != nil {
				return err
			}
			solNonce := types.SolanaNonce{Nonce: nonce.Nonce[:]}
			return k.LastUsedSolanaNonce.Set(ctx, solParams.NonceAccountKey, solNonce)
		},
		UpdatePendingTxStatus: func(tx zenbtctypes.PendingMintTransaction) error {
			if !fieldHasConsensus(oracleData.FieldVotePowers, VEFieldSolanaMintEventsHash) {
				k.Logger(ctx).Debug("Skipping Solana mint retry/timeout checks – no consensus on SolanaMintEventsHash", "tx_id", tx.Id)
				return nil
			}
			if tx.BlockHeight == 0 {
				return k.zenBTCKeeper.SetPendingMintTransaction(ctx, tx)
			}
			solParams := k.zenBTCKeeper.GetSolanaParams(ctx)
			if ctx.BlockHeight() > tx.BlockHeight+solParams.Btl {
				tx = k.processBtlSolanaMint(ctx, tx, oracleData, *solParams)
			}
			if tx.AwaitingEventSince > 0 {
				tx = k.processSecondaryTimeoutSolanaMint(ctx, tx, oracleData, *solParams)
			}
			return k.zenBTCKeeper.SetPendingMintTransaction(ctx, tx)
		},
	})
}

// processSolanaZenBTCMintEvents finalizes pending Solana zenBTC mints once oracle events match signatures.
func (k *Keeper) processSolanaZenBTCMintEvents(ctx sdk.Context, oracleData OracleData) {
	k.Logger(ctx).Info("ProcessSolanaZenBTCMintEvents: Started.", "oracle_event_count", len(oracleData.SolanaMintEvents))
	firstPendingID, err := k.zenBTCKeeper.GetFirstPendingSolMintTransaction(ctx)
	if err != nil {
		return
	}
	if firstPendingID == 0 {
		return
	}
	pendingMint, err := k.zenBTCKeeper.GetPendingMintTransaction(ctx, firstPendingID)
	if err != nil {
		return
	}
	if pendingMint.ZrchainTxId == 0 {
		return
	}
	signTxReq, err := k.treasuryKeeper.GetSignTransactionRequest(ctx, pendingMint.ZrchainTxId)
	if err != nil {
		return
	}
	mainSignReq, err := k.treasuryKeeper.GetSignRequest(ctx, signTxReq.SignRequestId)
	if err != nil {
		return
	}
	var signatures [][]byte
	for _, childReqID := range mainSignReq.ChildReqIds {
		childReq, err := k.treasuryKeeper.GetSignRequest(ctx, childReqID)
		if err != nil || len(childReq.SignedData) == 0 || len(childReq.SignedData[0].SignedData) == 0 {
			return
		}
		signatures = append(signatures, childReq.SignedData[0].SignedData)
	}
	if len(signatures) == 0 {
		return
	}
	concatenated := make([]byte, 0)
	for _, s := range signatures {
		concatenated = append(concatenated, s...)
	}
	sigHash := sha256.Sum256(concatenated)
	for _, event := range oracleData.SolanaMintEvents {
		if hex.EncodeToString(event.SigHash) == hex.EncodeToString(sigHash[:]) {
			supply, err := k.zenBTCKeeper.GetSupply(ctx)
			if err != nil {
				return
			}
			supply.PendingZenBTC -= pendingMint.Amount
			supply.MintedZenBTC += pendingMint.Amount
			if err := k.zenBTCKeeper.SetSupply(ctx, supply); err != nil {
				return
			}
			pendingMint.TxHash = event.TxSig
			pendingMint.Status = zenbtctypes.MintTransactionStatus_MINT_TRANSACTION_STATUS_MINTED
			if err := k.zenBTCKeeper.SetPendingMintTransaction(ctx, pendingMint); err != nil {
				k.Logger(ctx).Error("error updating pending mint transaction", "tx_id", pendingMint.Id, "error", err)
			}
			if err := k.adjustDefaultValidatorBedrockBTC(ctx, sdkmath.NewIntFromUint64(pendingMint.Amount)); err != nil {
				k.Logger(ctx).Error("error adjusting bedrock BTC on solana event mint", "amount", pendingMint.Amount, "error", err)
			}
			break
		}
	}
}

// storeNewZenBTCBurnEventsEthereum stores new zenBTC burn events coming from Ethereum.
func (k *Keeper) storeNewZenBTCBurnEventsEthereum(ctx sdk.Context, oracleData OracleData) {
	k.storeNewZenBTCBurnEvents(ctx, oracleData.EthBurnEvents, "ethereum", "error setting EthereumNonceRequested state")
}

// storeNewZenBTCBurnEventsSolana stores new zenBTC burn events coming from Solana.
func (k *Keeper) storeNewZenBTCBurnEventsSolana(ctx sdk.Context, oracleData OracleData) {
	k.storeNewZenBTCBurnEvents(ctx, oracleData.SolanaBurnEvents, "solana", "error setting EthereumNonceRequested state for unstaker")
}

// storeNewZenBTCBurnEvents stores new burn events from a given source. Filters to zenBTC burns only.
func (k *Keeper) storeNewZenBTCBurnEvents(ctx sdk.Context, burnEvents []sidecarapitypes.BurnEvent, source string, _ string) {
	processedInThisRun := make(map[string]bool)
	processedTxHashes := make(map[string]bool)
	for _, burn := range burnEvents {
		eventKey := fmt.Sprintf("%s-%d-%s", burn.TxID, burn.LogIndex, burn.ChainID)
		if processedInThisRun[eventKey] {
			continue
		}
		processedInThisRun[eventKey] = true
		if !burn.IsZenBTC {
			continue
		}
		exists := false
		if err := k.zenBTCKeeper.WalkBurnEvents(ctx, func(id uint64, existingBurn zenbtctypes.BurnEvent) (bool, error) {
			if existingBurn.TxID == burn.TxID && existingBurn.LogIndex == burn.LogIndex && existingBurn.ChainID == burn.ChainID {
				exists = true
				return true, nil
			}
			return false, nil
		}); err != nil {
			continue
		}
		if exists {
			continue
		}

		newBurn := zenbtctypes.BurnEvent{
			TxID:            burn.TxID,
			LogIndex:        burn.LogIndex,
			ChainID:         burn.ChainID,
			DestinationAddr: burn.DestinationAddr,
			Amount:          burn.Amount,
			Status:          zenbtctypes.BurnStatus_BURN_STATUS_BURNED,
		}
		_, createErr := k.zenBTCKeeper.CreateBurnEvent(ctx, &newBurn)
		if createErr != nil {
			k.Logger(ctx).Error("error creating burn event", "burn_tx", burn.TxID, "chain_id", burn.ChainID, "error", createErr)
			continue
		}

		// Get next redemption ID by finding the highest existing ID + 1
		// Walk in descending order and stop after first item (highest ID)
		nextRedemptionID := uint64(0)
		if err := k.zenBTCKeeper.WalkRedemptionsDescending(ctx, func(id uint64, _ zenbtctypes.Redemption) (bool, error) {
			nextRedemptionID = id + 1
			return true, nil // stop after first (highest) ID
		}); err != nil {
			k.Logger(ctx).Error("error walking redemptions to find next ID", "burn_tx", burn.TxID, "error", err)
			continue
		}

		if err := k.zenBTCKeeper.SetRedemption(ctx, nextRedemptionID, zenbtctypes.Redemption{
			Data: zenbtctypes.RedemptionData{
				Id:                 nextRedemptionID,
				DestinationAddress: burn.DestinationAddr,
				Amount:             burn.Amount,
			},
			Status: zenbtctypes.RedemptionStatus_UNSTAKED,
		}); err != nil {
			k.Logger(ctx).Error("error creating redemption from burn event", "burn_tx", burn.TxID, "error", err)
			continue
		}

		k.Logger(ctx).Info("created redemption for burn event", "redemption_id", nextRedemptionID, "burn_tx", burn.TxID, "amount", burn.Amount, "destination", hex.EncodeToString(burn.DestinationAddr))
		processedTxHashes[burn.TxID] = true
	}
	k.ClearProcessedBackfillRequests(ctx, types.EventType_EVENT_TYPE_ZENBTC_BURN, processedTxHashes)
}

// processZenBTCBurnEvents constructs unstake transactions for BURNED events.
// func (k *Keeper) processZenBTCBurnEvents(ctx sdk.Context, oracleData OracleData) {
// 	processEthereumTxQueue(k, ctx, EthereumTxQueueArgs[zenbtctypes.BurnEvent]{
// 		KeyID:                    k.zenBTCKeeper.GetUnstakerKeyID(ctx),
// 		RequestedNonce:           oracleData.RequestedUnstakerNonce,
// 		DispatchRequestedChecker: TxDispatchRequestHandler[uint64]{Store: k.EthereumNonceRequested},
// 		GetPendingTxs: func(ctx sdk.Context) ([]zenbtctypes.BurnEvent, error) {
// 			return k.getPendingBurnEvents(ctx)
// 		},
// 		DispatchTx: func(be zenbtctypes.BurnEvent) error {
// 			if err := k.zenBTCKeeper.SetFirstPendingBurnEvent(ctx, be.Id); err != nil {
// 				return err
// 			}
// 			requiredFields := []VoteExtensionField{VEFieldRequestedUnstakerNonce, VEFieldBTCUSDPrice, VEFieldETHUSDPrice}
// 			if err := k.validateConsensusForTxFields(ctx, oracleData, requiredFields, "zenBTC burn unstake", fmt.Sprintf("burn_id: %d", be.Id)); err != nil {
// 				return nil
// 			}
// 			if len(be.DestinationAddr) == 0 {
// 				return fmt.Errorf("burn event %d has empty DestinationAddr", be.Id)
// 			}
// 			unsignedTxHash, unsignedTx, err := k.constructUnstakeTx(
// 				ctx,
// 				getChainIDForEigen(ctx),
// 				be.DestinationAddr,
// 				be.Amount,
// 				oracleData.RequestedUnstakerNonce,
// 				oracleData.EthBaseFee,
// 				oracleData.EthTipCap,
// 			)
// 			if err != nil {
// 				return err
// 			}
// 			creator, err := k.getAddressByKeyID(ctx, k.zenBTCKeeper.GetUnstakerKeyID(ctx), treasurytypes.WalletType_WALLET_TYPE_NATIVE)
// 			if err != nil {
// 				return err
// 			}
// 			return k.submitEthereumTransaction(ctx, creator, k.zenBTCKeeper.GetUnstakerKeyID(ctx), treasurytypes.WalletType_WALLET_TYPE_EVM, getChainIDForEigen(ctx), unsignedTx, unsignedTxHash)
// 		},
// 		OnTxConfirmed: func(be zenbtctypes.BurnEvent) error {
// 			be.Status = zenbtctypes.BurnStatus_BURN_STATUS_UNSTAKING
// 			return k.zenBTCKeeper.SetBurnEvent(ctx, be.Id, be)
// 		},
// 	})
// }

// // storeNewZenBTCRedemptions ingests new redemptions and requests completer nonce if needed.
// func (k *Keeper) storeNewZenBTCRedemptions(ctx sdk.Context, oracleData OracleData) {
// 	exchangeRate, err := k.zenBTCKeeper.GetExchangeRate(ctx)
// 	if err != nil {
// 		return
// 	}
// 	foundNew := false
// 	for _, redemption := range oracleData.Redemptions {
// 		if exists, err := k.zenBTCKeeper.HasRedemption(ctx, redemption.Id); err != nil || exists {
// 			if err != nil {
// 				k.Logger(ctx).Error("error checking redemption existence", "id", redemption.Id, "error", err)
// 			}
// 			continue
// 		}
// 		foundNew = true
// 		btcAmount := sdkmath.LegacyNewDecFromInt(sdkmath.NewIntFromUint64(redemption.Amount)).Mul(exchangeRate).TruncateInt64()
// 		if err := k.zenBTCKeeper.SetRedemption(ctx, redemption.Id, zenbtctypes.Redemption{
// 			Data:   zenbtctypes.RedemptionData{Id: redemption.Id, DestinationAddress: redemption.DestinationAddress, Amount: uint64(btcAmount)},
// 			Status: zenbtctypes.RedemptionStatus_INITIATED,
// 		}); err != nil {
// 			k.Logger(ctx).Error("error setting redemption", "id", redemption.Id, "error", err)
// 		}
// 	}
// 	if foundNew {
// 		if err := k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetCompleterKeyID(ctx), true); err != nil {
// 			k.Logger(ctx).Error("error setting completer nonce requested flag", "error", err)
// 		}
// 	}
// }

// checkForRedemptionFulfilment updates supplies when treasury sign requests for redemptions are fulfilled.
func (k *Keeper) checkForRedemptionFulfilment(ctx sdk.Context) {
	startingIndex, _ := k.zenBTCKeeper.GetFirstRedemptionAwaitingSign(ctx)
	redemptions, err := k.GetRedemptionsByStatus(ctx, zenbtctypes.RedemptionStatus_AWAITING_SIGN, 0, startingIndex)
	if err != nil || len(redemptions) == 0 {
		return
	}
	if err := k.zenBTCKeeper.SetFirstRedemptionAwaitingSign(ctx, redemptions[0].Data.Id); err != nil {
		k.Logger(ctx).Error("error setting first redemption awaiting sign", "id", redemptions[0].Data.Id, "error", err)
	}
	supply, err := k.zenBTCKeeper.GetSupply(ctx)
	if err != nil {
		return
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
			exchangeRate, err := k.zenBTCKeeper.GetExchangeRate(ctx)
			if err != nil {
				continue
			}
			btcToRelease := uint64(sdkmath.LegacyNewDecFromInt(sdkmath.NewIntFromUint64(redemption.Data.Amount)).Quo(exchangeRate).TruncateInt64())
			if supply.MintedZenBTC < redemption.Data.Amount || supply.CustodiedBTC < btcToRelease {
				continue
			}
			supply.MintedZenBTC -= redemption.Data.Amount
			supply.CustodiedBTC -= btcToRelease
			if err := k.adjustDefaultValidatorBedrockBTC(ctx, sdkmath.NewIntFromUint64(redemption.Data.Amount).Neg()); err != nil {
				k.Logger(ctx).Error("error adjusting bedrock BTC on redemption fulfilment", "id", redemption.Data.Id, "error", err)
			}
			redemption.Status = zenbtctypes.RedemptionStatus_COMPLETED
		}
		if signReq.Status == treasurytypes.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED {
			redemption.Data.SignReqId = 0
			redemption.Status = zenbtctypes.RedemptionStatus_UNSTAKED
		}
		if err := k.zenBTCKeeper.SetRedemption(ctx, redemption.Data.Id, redemption); err != nil {
			k.Logger(ctx).Error("error updating redemption after fulfilment", "id", redemption.Data.Id, "error", err)
		}
	}
	if err := k.zenBTCKeeper.SetSupply(ctx, supply); err != nil {
		k.Logger(ctx).Error("error updating supply after fulfilment", "error", err)
	}
}

// adjustDefaultValidatorBedrockBTC adds (positive) or subtracts (negative) BTC sats to the default validator's TokensBedrock (Asset_BTC)
// DEPRECATED: This function is kept for backward compatibility. Use distributeBedrockBTC instead.
func (k *Keeper) adjustDefaultValidatorBedrockBTC(ctx sdk.Context, delta sdkmath.Int) error {
	return k.distributeBedrockBTC(ctx, delta)
}

// distributeBedrockBTC adds (positive) or subtracts (negative) BTC sats across the bedrock validator set
// proportional to each validator's native stake (TokensNative)
func (k *Keeper) distributeBedrockBTC(ctx sdk.Context, delta sdkmath.Int) error {
	if delta.IsZero() {
		return nil
	}

	// Get all validators in the bedrock set
	bedrockValidators := make([]types.ValidatorHV, 0)
	totalNativeStake := sdkmath.ZeroInt()

	err := k.BedrockValidatorSet.Walk(ctx, nil, func(valAddr string, inSet bool) (bool, error) {
		if !inSet {
			return false, nil
		}

		validator, err := k.GetZenrockValidatorFromBech32(ctx, valAddr)
		if err != nil {
			k.Logger(ctx).Error("validator in bedrock set not found", "address", valAddr, "error", err)
			return false, nil
		}

		bedrockValidators = append(bedrockValidators, validator)
		totalNativeStake = totalNativeStake.Add(validator.TokensNative)
		return false, nil
	})

	if err != nil {
		return err
	}

	// If no validators in bedrock set, log warning and skip
	if len(bedrockValidators) == 0 {
		k.Logger(ctx).Warn("no validators in bedrock set, skipping bedrock BTC distribution", "delta", delta.String())
		return nil
	}

	// If total native stake is zero, distribute equally
	if totalNativeStake.IsZero() {
		k.Logger(ctx).Warn("total native stake is zero, distributing bedrock BTC equally")
		amountPerValidator := delta.Quo(sdkmath.NewInt(int64(len(bedrockValidators))))
		for _, validator := range bedrockValidators {
			if err := k.adjustValidatorBedrockBTC(ctx, validator, amountPerValidator); err != nil {
				k.Logger(ctx).Error("failed to adjust bedrock BTC for validator", "validator", validator.OperatorAddress, "error", err)
			}
		}
		return nil
	}

	// Distribute proportionally based on native stake
	remainingDelta := delta
	for i, validator := range bedrockValidators {
		var allocation sdkmath.Int
		if i == len(bedrockValidators)-1 {
			// Last validator gets the remaining amount to handle rounding
			allocation = remainingDelta
		} else {
			// Calculate proportional allocation: (validator.TokensNative / totalNativeStake) * delta
			allocation = validator.TokensNative.Mul(delta).Quo(totalNativeStake)
			remainingDelta = remainingDelta.Sub(allocation)
		}

		if err := k.adjustValidatorBedrockBTC(ctx, validator, allocation); err != nil {
			k.Logger(ctx).Error("failed to adjust bedrock BTC for validator", "validator", validator.OperatorAddress, "error", err)
		}
	}

	return nil
}

// adjustValidatorBedrockBTC adds (positive) or subtracts (negative) BTC sats to a validator's TokensBedrock (Asset_BTC)
func (k *Keeper) adjustValidatorBedrockBTC(ctx sdk.Context, validator types.ValidatorHV, delta sdkmath.Int) error {
	idx := -1
	for i, td := range validator.TokensBedrock {
		if td != nil && td.Asset == types.Asset_BTC {
			idx = i
			break
		}
	}

	if idx >= 0 {
		newAmt := validator.TokensBedrock[idx].Amount.Add(delta)
		if newAmt.IsNegative() {
			newAmt = sdkmath.ZeroInt()
		}
		validator.TokensBedrock[idx].Amount = newAmt
	} else {
		amt := delta
		if amt.IsNegative() {
			amt = sdkmath.ZeroInt()
		}
		validator.TokensBedrock = append(validator.TokensBedrock, &types.TokenData{Asset: types.Asset_BTC, Amount: amt})
	}

	return k.SetValidator(ctx, validator)
}

// rebalanceBedrockBTC rebalances bedrock BTC across the bedrock validator set
// to keep proportions aligned with native stake. This is called every block.
func (k *Keeper) rebalanceBedrockBTC(ctx sdk.Context) error {

	// Get all validators in the bedrock set
	bedrockValidators := make([]types.ValidatorHV, 0)
	totalNativeStake := sdkmath.ZeroInt()
	totalBedrockBTC := sdkmath.ZeroInt()

	err := k.BedrockValidatorSet.Walk(ctx, nil, func(valAddr string, inSet bool) (bool, error) {
		if !inSet {
			return false, nil
		}

		validator, err := k.GetZenrockValidatorFromBech32(ctx, valAddr)
		if err != nil {
			k.Logger(ctx).Error("validator in bedrock set not found during rebalance", "address", valAddr, "error", err)
			return false, nil
		}

		// Get current bedrock BTC for this validator
		currentBTC := k.getBedrockTokenAmount(validator, types.Asset_BTC)

		bedrockValidators = append(bedrockValidators, validator)
		totalNativeStake = totalNativeStake.Add(validator.TokensNative)
		totalBedrockBTC = totalBedrockBTC.Add(currentBTC)
		return false, nil
	})

	if err != nil {
		return err
	}

	// If no validators in bedrock set or no BTC to rebalance, skip
	if len(bedrockValidators) == 0 || totalBedrockBTC.IsZero() {
		return nil
	}

	// If total native stake is zero, distribute equally
	if totalNativeStake.IsZero() {
		amountPerValidator := totalBedrockBTC.Quo(sdkmath.NewInt(int64(len(bedrockValidators))))
		for _, validator := range bedrockValidators {
			currentBTC := k.getBedrockTokenAmount(validator, types.Asset_BTC)
			adjustment := amountPerValidator.Sub(currentBTC)
			if !adjustment.IsZero() {
				if err := k.adjustValidatorBedrockBTC(ctx, validator, adjustment); err != nil {
					k.Logger(ctx).Error("failed to rebalance bedrock BTC for validator", "validator", validator.OperatorAddress, "error", err)
				}
			}
		}
		return nil
	}

	// Calculate target amounts and adjust
	remainingBTC := totalBedrockBTC
	for i, validator := range bedrockValidators {
		var targetAmount sdkmath.Int
		if i == len(bedrockValidators)-1 {
			// Last validator gets the remaining amount to handle rounding
			targetAmount = remainingBTC
		} else {
			// Calculate proportional target: (validator.TokensNative / totalNativeStake) * totalBedrockBTC
			targetAmount = validator.TokensNative.Mul(totalBedrockBTC).Quo(totalNativeStake)
			remainingBTC = remainingBTC.Sub(targetAmount)
		}

		currentBTC := k.getBedrockTokenAmount(validator, types.Asset_BTC)
		adjustment := targetAmount.Sub(currentBTC)

		if !adjustment.IsZero() {
			if err := k.adjustValidatorBedrockBTC(ctx, validator, adjustment); err != nil {
				k.Logger(ctx).Error("failed to rebalance bedrock BTC for validator", "validator", validator.OperatorAddress, "error", err)
			}
		}
	}

	return nil
}

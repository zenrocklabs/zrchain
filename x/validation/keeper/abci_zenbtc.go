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
processEthereumTxQueue(k, ctx, EVMQueueArgs[zenbtctypes.PendingMintTransaction]{
		KeyID:               k.zenBTCKeeper.GetEthMinterKeyID(ctx),
		RequestedNonce:      oracleData.RequestedEthMinterNonce,
DispatchRequestedChecker: TxDispatchRequestChecker[uint64]{M: k.EthereumNonceRequested},
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
	processSolanaTxQueue(k, ctx, SolanaQueueArgs[zenbtctypes.PendingMintTransaction]{
		NonceAccountKey:     k.zenBTCKeeper.GetSolanaParams(ctx).NonceAccountKey,
		NonceAccount:        oracleData.SolanaMintNonces[k.zenBTCKeeper.GetSolanaParams(ctx).NonceAccountKey],
DispatchRequestedChecker: TxDispatchRequestChecker[uint64]{M: k.SolanaNonceRequested},
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
				eventID:           tx.Id,
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
			_ = k.zenBTCKeeper.SetPendingMintTransaction(ctx, pendingMint)
			_ = k.adjustDefaultValidatorBedrockBTC(ctx, sdkmath.NewIntFromUint64(pendingMint.Amount))
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
		if !exists {
			newBurn := zenbtctypes.BurnEvent{
				TxID:            burn.TxID,
				LogIndex:        burn.LogIndex,
				ChainID:         burn.ChainID,
				DestinationAddr: burn.DestinationAddr,
				Amount:          burn.Amount,
				Status:          zenbtctypes.BurnStatus_BURN_STATUS_BURNED,
			}
			createdID, createErr := k.zenBTCKeeper.CreateBurnEvent(ctx, &newBurn)
			if createErr == nil {
				if has, err := k.zenBTCKeeper.HasRedemption(ctx, createdID); err == nil && !has {
					_ = k.zenBTCKeeper.SetRedemption(ctx, createdID, zenbtctypes.Redemption{
						Data: zenbtctypes.RedemptionData{
							Id:                 createdID,
							DestinationAddress: burn.DestinationAddr,
							Amount:             burn.Amount,
						},
						Status: zenbtctypes.RedemptionStatus_UNSTAKED,
					})
				}
			}
			processedTxHashes[burn.TxID] = true
		}
	}
	k.ClearProcessedBackfillRequests(ctx, types.EventType_EVENT_TYPE_ZENBTC_BURN, processedTxHashes)
}

// processZenBTCBurnEvents constructs unstake transactions for BURNED events.
func (k *Keeper) processZenBTCBurnEvents(ctx sdk.Context, oracleData OracleData) {
processEthereumTxQueue(k, ctx, EVMQueueArgs[zenbtctypes.BurnEvent]{
		KeyID:               k.zenBTCKeeper.GetUnstakerKeyID(ctx),
		RequestedNonce:      oracleData.RequestedUnstakerNonce,
DispatchRequestedChecker: TxDispatchRequestChecker[uint64]{M: k.EthereumNonceRequested},
		GetPendingTxs: func(ctx sdk.Context) ([]zenbtctypes.BurnEvent, error) {
			return k.getPendingBurnEvents(ctx)
		},
		DispatchTx: func(be zenbtctypes.BurnEvent) error {
			if err := k.zenBTCKeeper.SetFirstPendingBurnEvent(ctx, be.Id); err != nil {
				return err
			}
			requiredFields := []VoteExtensionField{VEFieldRequestedUnstakerNonce, VEFieldBTCUSDPrice, VEFieldETHUSDPrice}
			if err := k.validateConsensusForTxFields(ctx, oracleData, requiredFields, "zenBTC burn unstake", fmt.Sprintf("burn_id: %d", be.Id)); err != nil {
				return nil
			}
			if len(be.DestinationAddr) == 0 {
				return fmt.Errorf("burn event %d has empty DestinationAddr", be.Id)
			}
			unsignedTxHash, unsignedTx, err := k.constructUnstakeTx(
				ctx,
				getChainIDForEigen(ctx),
				be.DestinationAddr,
				be.Amount,
				oracleData.RequestedUnstakerNonce,
				oracleData.EthBaseFee,
				oracleData.EthTipCap,
			)
			if err != nil {
				return err
			}
			creator, err := k.getAddressByKeyID(ctx, k.zenBTCKeeper.GetUnstakerKeyID(ctx), treasurytypes.WalletType_WALLET_TYPE_NATIVE)
			if err != nil {
				return err
			}
			return k.submitEthereumTransaction(ctx, creator, k.zenBTCKeeper.GetUnstakerKeyID(ctx), treasurytypes.WalletType_WALLET_TYPE_EVM, getChainIDForEigen(ctx), unsignedTx, unsignedTxHash)
		},
		OnTxConfirmed: func(be zenbtctypes.BurnEvent) error {
			be.Status = zenbtctypes.BurnStatus_BURN_STATUS_UNSTAKING
			return k.zenBTCKeeper.SetBurnEvent(ctx, be.Id, be)
		},
	})
}

// storeNewZenBTCRedemptions ingests new redemptions and requests completer nonce if needed.
func (k *Keeper) storeNewZenBTCRedemptions(ctx sdk.Context, oracleData OracleData) {
	var firstInitiated zenbtctypes.Redemption
	_ = k.zenBTCKeeper.WalkRedemptions(ctx, func(id uint64, r zenbtctypes.Redemption) (bool, error) {
		if r.Status == zenbtctypes.RedemptionStatus_INITIATED {
			firstInitiated = r
			return true, nil
		}
		return false, nil
	})
	if len(oracleData.Redemptions) == 0 {
		return
	}
	exchangeRate, err := k.zenBTCKeeper.GetExchangeRate(ctx)
	if err != nil {
		return
	}
	foundNew := false
	for _, redemption := range oracleData.Redemptions {
		if exists, _ := k.zenBTCKeeper.HasRedemption(ctx, redemption.Id); exists {
			continue
		}
		foundNew = true
		btcAmount := sdkmath.LegacyNewDecFromInt(sdkmath.NewIntFromUint64(redemption.Amount)).Mul(exchangeRate).TruncateInt64()
		_ = k.zenBTCKeeper.SetRedemption(ctx, redemption.Id, zenbtctypes.Redemption{
			Data:   zenbtctypes.RedemptionData{Id: redemption.Id, DestinationAddress: redemption.DestinationAddress, Amount: uint64(btcAmount)},
			Status: zenbtctypes.RedemptionStatus_INITIATED,
		})
	}
	if foundNew {
		_ = k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetCompleterKeyID(ctx), true)
	}
	_ = firstInitiated // keep behavior if needed by future logic
}

// processZenBTCRedemptions completes INITIATED redemptions by calling 'complete' on EigenLayer.
func (k *Keeper) processZenBTCRedemptions(ctx sdk.Context, oracleData OracleData) {
processEthereumTxQueue(k, ctx, EVMQueueArgs[zenbtctypes.Redemption]{
		KeyID:               k.zenBTCKeeper.GetCompleterKeyID(ctx),
		RequestedNonce:      oracleData.RequestedCompleterNonce,
DispatchRequestedChecker: TxDispatchRequestChecker[uint64]{M: k.EthereumNonceRequested},
		GetPendingTxs: func(ctx sdk.Context) ([]zenbtctypes.Redemption, error) {
			firstPendingID, _ := k.zenBTCKeeper.GetFirstPendingRedemption(ctx)
			return k.GetRedemptionsByStatus(ctx, zenbtctypes.RedemptionStatus_INITIATED, 2, firstPendingID)
		},
		DispatchTx: func(r zenbtctypes.Redemption) error {
			if err := k.zenBTCKeeper.SetFirstPendingRedemption(ctx, r.Data.Id); err != nil {
				return err
			}
			if err := k.validateConsensusForTxFields(ctx, oracleData, []VoteExtensionField{VEFieldRequestedCompleterNonce, VEFieldBTCUSDPrice, VEFieldETHUSDPrice},
				"zenBTC redemption", fmt.Sprintf("redemption_id: %d, amount: %d", r.Data.Id, r.Data.Amount)); err != nil {
				return err
			}
			unsignedTxHash, unsignedTx, err := k.constructCompleteTx(ctx, getChainIDForEigen(ctx), r.Data.Id, oracleData.RequestedCompleterNonce, oracleData.EthBaseFee, oracleData.EthTipCap)
			if err != nil {
				return err
			}
			creator, err := k.getAddressByKeyID(ctx, k.zenBTCKeeper.GetCompleterKeyID(ctx), treasurytypes.WalletType_WALLET_TYPE_NATIVE)
			if err != nil {
				return err
			}
			return k.submitEthereumTransaction(ctx, creator, k.zenBTCKeeper.GetCompleterKeyID(ctx), treasurytypes.WalletType_WALLET_TYPE_EVM, getChainIDForEigen(ctx), unsignedTx, unsignedTxHash)
		},
		OnTxConfirmed: func(r zenbtctypes.Redemption) error {
			r.Status = zenbtctypes.RedemptionStatus_UNSTAKED
			if err := k.zenBTCKeeper.SetRedemption(ctx, r.Data.Id, r); err != nil {
				return err
			}
			return k.EthereumNonceRequested.Set(ctx, k.zenBTCKeeper.GetStakerKeyID(ctx), true)
		},
	})
}

// checkForRedemptionFulfilment updates supplies when treasury sign requests for redemptions are fulfilled.
func (k *Keeper) checkForRedemptionFulfilment(ctx sdk.Context) {
	startingIndex, _ := k.zenBTCKeeper.GetFirstRedemptionAwaitingSign(ctx)
	redemptions, err := k.GetRedemptionsByStatus(ctx, zenbtctypes.RedemptionStatus_AWAITING_SIGN, 0, startingIndex)
	if err != nil || len(redemptions) == 0 {
		return
	}
	_ = k.zenBTCKeeper.SetFirstRedemptionAwaitingSign(ctx, redemptions[0].Data.Id)
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
			_ = k.adjustDefaultValidatorBedrockBTC(ctx, sdkmath.NewIntFromUint64(redemption.Data.Amount).Neg())
			redemption.Status = zenbtctypes.RedemptionStatus_COMPLETED
		}
		if signReq.Status == treasurytypes.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED {
			redemption.Data.SignReqId = 0
			redemption.Status = zenbtctypes.RedemptionStatus_UNSTAKED
		}
		_ = k.zenBTCKeeper.SetRedemption(ctx, redemption.Data.Id, redemption)
	}
	_ = k.zenBTCKeeper.SetSupply(ctx, supply)
}

// adjustDefaultValidatorBedrockBTC adds (positive) or subtracts (negative) BTC sats to the default validator's TokensBedrock (Asset_BTC)
func (k *Keeper) adjustDefaultValidatorBedrockBTC(ctx sdk.Context, delta sdkmath.Int) error {
	oper := k.GetBedrockDefaultValOperAddr(ctx)
	v, err := k.GetZenrockValidatorFromBech32(ctx, oper)
	if err != nil {
		return err
	}
	idx := -1
	for i, td := range v.TokensBedrock {
		if td != nil && td.Asset == types.Asset_BTC {
			idx = i
			break
		}
	}
	if idx >= 0 {
		newAmt := v.TokensBedrock[idx].Amount.Add(delta)
		if newAmt.IsNegative() {
			newAmt = sdkmath.ZeroInt()
		}
		v.TokensBedrock[idx].Amount = newAmt
	} else {
		amt := delta
		if amt.IsNegative() {
			amt = sdkmath.ZeroInt()
		}
		v.TokensBedrock = append(v.TokensBedrock, &types.TokenData{Asset: types.Asset_BTC, Amount: amt})
	}
	return k.SetValidator(ctx, v)
}

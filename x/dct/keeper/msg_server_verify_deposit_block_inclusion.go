package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/rpc"
	"strings"

	"github.com/Zenrock-Foundation/zrchain/v6/zenbtc/utils"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/bitcoin"
	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"

	"github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
)

func (k msgServer) VerifyDepositBlockInclusion(goCtx context.Context, msg *types.MsgVerifyDepositBlockInclusion) (*types.MsgVerifyDepositBlockInclusionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	asset := msg.Asset
	if asset == types.Asset_ASSET_UNSPECIFIED {
		return nil, fmt.Errorf("asset must be specified")
	}
	// IMPORTANT: zenBTC deposits are handled by the zenBTC module, not DCT
	// DCT module is for v1+ assets only (zenZEC and future wrapped assets)
	if asset == types.Asset_ASSET_ZENBTC {
		return nil, fmt.Errorf("zenBTC deposits must use the zenBTC module's VerifyDepositBlockInclusion endpoint, not DCT")
	}

	// Verify that the message creator is the BitcoinProxyAddress from module params
	params, err := k.GetAssetParams(ctx, asset)
	if err != nil {
		return nil, fmt.Errorf("failed to get module params: %w", err)
	}
	if msg.Creator != params.ProxyAddress {
		return nil, fmt.Errorf("unauthorized: only the BitcoinProxyAddress (%s) can verify deposit block inclusion, got %s", params.ProxyAddress, msg.Creator)
	}

	// Fetch the appropriate block header based on asset type
	var blockHeader api.BTCBlockHeader
	if asset == types.Asset_ASSET_ZENZEC {
		// Use ZCash block headers for ZCash deposits
		blockHeader, err = k.validationKeeper.ZcashBlockHeaders.Get(ctx, msg.BlockHeight)
		if err != nil {
			return nil, fmt.Errorf("ZCash block header not found at height %d: %w", msg.BlockHeight, err)
		}
	} else {
		// For future assets, determine the correct block header collection
		return nil, fmt.Errorf("unsupported asset type %s - block header retrieval not implemented", asset.String())
	}

	ignoreAddresses, err := k.changeAddressesForAsset(ctx, asset, msg.ChainName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving change addresses: %w", err)
	}

	outputs, _, err := bitcoin.VerifyBTCLockTransaction(msg.RawTx, msg.ChainName, int(msg.Index), msg.Proof, &blockHeader, ignoreAddresses)
	if err != nil {
		return nil, err
	}

	var matchedOutput bool
	for _, output := range outputs {
		if output.Address == msg.DepositAddr && output.Amount == msg.Amount && uint64(output.OutputIndex) == msg.Vout {
			matchedOutput = true
			break
		}
	}
	if !matchedOutput {
		return nil, fmt.Errorf("%s deposit not found in outputs from provided lock tx", asset.String())
	}

	depositKeyringAddr, err := k.GetDepositKeyringAddr(ctx, asset)
	if err != nil {
		return nil, fmt.Errorf("error retrieving deposit keyring address: %w", err)
	}

	queryResp, err := k.treasuryKeeper.KeyByAddress(ctx, &treasurytypes.QueryKeyByAddressRequest{
		Address:     msg.DepositAddr,
		KeyringAddr: depositKeyringAddr,
		KeyType:     treasurytypes.KeyType_KEY_TYPE_BITCOIN_SECP256K1,
		WalletType:  WalletTypeFromChainName(msg),
	})
	if err != nil {
		return nil, err
	}
	if queryResp.Response == nil || queryResp.Response.Wallets == nil {
		return nil, fmt.Errorf("%s deposit address does not correspond to correct key (no wallets returned)", asset.String())
	}

	metaData := queryResp.Response.Key.ZenbtcMetadata
	if metaData == nil || metaData.RecipientAddr == "" || metaData.Caip2ChainId == "" {
		return nil, errors.New("lock tx - key metadata is invalid")
	}

	// Verify asset type matches the expected asset
	// ASSET_UNSPECIFIED is allowed for backwards compatibility
	if metaData.Asset != types.Asset_ASSET_UNSPECIFIED && metaData.Asset != asset {
		return nil, fmt.Errorf("key metadata asset type mismatch: expected %s or ASSET_UNSPECIFIED, got %s", asset.String(), metaData.Asset.String())
	}

	var walletFound bool
	for _, wallet := range queryResp.Response.Wallets {
		if wallet.Address == msg.DepositAddr {
			walletFound = true
			break
		}
	}
	if !walletFound {
		return nil, fmt.Errorf("%s deposit address does not correspond to correct key (no matching wallet)", asset.String())
	}

	toBeHashed := fmt.Sprintf("%s:%d", msg.RawTx, msg.Vout)
	hash := sha256.Sum256([]byte(toBeHashed))
	lockTxKey := hex.EncodeToString(hash[:])
	lockStoreKey, err := k.lockKey(asset, lockTxKey)
	if err != nil {
		return nil, err
	}
	txExists, err := k.LockTransactions.Has(ctx, lockStoreKey)
	if err != nil {
		return nil, err
	}
	if txExists {
		return nil, errors.New("lock tx was already previously used to mint wrapped tokens")
	}

	exchangeRate, err := k.GetExchangeRate(ctx, asset)
	if err != nil {
		return nil, err
	}

	wrappedAmountDec := math.LegacyNewDecFromInt(math.NewIntFromUint64(msg.Amount)).Quo(exchangeRate)
	wrappedAmount := uint64(wrappedAmountDec.TruncateInt64())

	supply, err := k.GetSupply(ctx, asset)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return nil, err
		}
		supply = types.Supply{
			Asset:           asset,
			CustodiedAmount: 0,
			MintedAmount:    0,
			PendingAmount:   0,
		}
	}

	supply.CustodiedAmount += msg.Amount

	rewardsKeyID, err := k.GetRewardsDepositKeyID(ctx, asset)
	if err != nil {
		return nil, err
	}
	isYieldDeposit := queryResp.Response.Key.Id == rewardsKeyID

	if !isYieldDeposit {
		supply.PendingAmount += wrappedAmount
	}

	if err := k.SetSupply(ctx, supply); err != nil {
		return nil, err
	}

	k.Logger().Info("custodied supply updated",
		"asset", asset.String(),
		"previous", supply.CustodiedAmount-msg.Amount,
		"current", supply.CustodiedAmount,
	)
	if !isYieldDeposit {
		k.Logger().Info("pending mint supply updated",
			"asset", asset.String(),
			"previous", supply.PendingAmount-wrappedAmount,
			"current", supply.PendingAmount,
		)
	}

	if err := k.LockTransactions.Set(ctx, lockStoreKey, types.LockTransaction{
		RawTx:         msg.RawTx,
		Vout:          msg.Vout,
		Sender:        msg.DepositAddr,
		MintRecipient: metaData.RecipientAddr,
		Amount:        msg.Amount,
		BlockHeight:   msg.BlockHeight,
		Asset:         asset,
	}); err != nil {
		return nil, err
	}

	if isYieldDeposit {
		return &types.MsgVerifyDepositBlockInclusionResponse{}, nil
	}

	chainType := types.WalletType(metaData.ChainType)
	if chainType != types.WalletType_WALLET_TYPE_SOLANA {
		return nil, fmt.Errorf("asset %s does not support destination chain type %s", asset.String(), chainType.String())
	}

	pendingMint := &types.PendingMintTransaction{
		Caip2ChainId:     metaData.Caip2ChainId,
		ChainType:        chainType,
		RecipientAddress: metaData.RecipientAddr,
		Amount:           wrappedAmount,
		Creator:          msg.Creator,
		Asset:            asset,
	}
	if _, err := k.CreatePendingMintTransaction(ctx, pendingMint); err != nil {
		return nil, err
	}
	k.validationKeeper.Logger(ctx).Info("added pending mint transaction", "asset", asset.String(), "tx_id", pendingMint.Id, "recipient", pendingMint.RecipientAddress, "amount", pendingMint.Amount)

	solanaParams, err := k.GetSolanaParams(ctx, asset)
	if err != nil {
		return nil, err
	}
	if err := k.validationKeeper.SetSolanaRequestedNonce(ctx, solanaParams.NonceAccountKey, true); err != nil {
		return nil, err
	}
	if err := k.validationKeeper.SetSolanaDCTRequestedAccount(ctx, asset, pendingMint.RecipientAddress, true); err != nil {
		return nil, err
	}

	return &types.MsgVerifyDepositBlockInclusionResponse{}, nil
}

func debugRetrieveBlockHeaderViaRPC(chainName string, blockHeight int64) (*api.BTCBlockHeader, error) {
	if chainName == "mainnet" {
		return nil, errors.New("cannot retrieve block header from mainnet")
	}
	type GetBlockHeaderByHeightArgs struct {
		ChainName string
		Height    int64
	}

	type GetBlockHeaderByHeightReply struct {
		BlockHeader *api.BTCBlockHeader
		BlockHash   *chainhash.Hash
		Height      int32
	}

	args := GetBlockHeaderByHeightArgs{
		ChainName: chainName,
		Height:    blockHeight,
	}
	var resp GetBlockHeaderByHeightReply
	client, err := rpc.DialHTTP("tcp", "localhost"+":12345")
	if err != nil {
		return nil, err
	}

	err = client.Call("NeutrinoServer.BlockHeaderByHeight", args, &resp)
	if err != nil {
		return nil, err
	}
	return resp.BlockHeader, nil
}

// changeAddressesForAsset derives the list of change addresses for the provided asset and chain.
func (k msgServer) changeAddressesForAsset(ctx context.Context, asset types.Asset, chainName string) ([]string, error) {
	keyIDs, err := k.GetChangeAddressKeyIDs(ctx, asset)
	if err != nil {
		return nil, err
	}
	if len(keyIDs) == 0 {
		return nil, fmt.Errorf("no change key IDs configured for asset %s", asset.String())
	}

	result := make([]string, 0, len(keyIDs))

	// Zcash assets require different address generation than Bitcoin
	if asset == types.Asset_ASSET_ZENZEC {
		// Extract network from chainName (e.g., "zcash-mainnet" -> "mainnet")
		network := chainName
		if strings.HasPrefix(chainName, "zcash-") {
			network = strings.TrimPrefix(chainName, "zcash-")
		}

		for _, keyID := range keyIDs {
			key, err := k.Keeper.treasuryKeeper.KeyStore.Get(ctx, keyID)
			if err != nil {
				return nil, err
			}
			address, err := treasurytypes.ZcashAddress(&key, network)
			if err != nil {
				return nil, err
			}
			result = append(result, address)
		}
	} else {
		// Bitcoin and other Bitcoin-based assets use P2WPKH
		chaincfg := utils.ChainFromString(chainName)
		for _, keyID := range keyIDs {
			key, err := k.Keeper.treasuryKeeper.KeyStore.Get(ctx, keyID)
			if err != nil {
				return nil, err
			}
			address, err := treasurytypes.BitcoinP2WPKH(&key, chaincfg)
			if err != nil {
				return nil, err
			}
			result = append(result, address)
		}
	}

	return result, nil
}

func WalletTypeFromChainName(msg *types.MsgVerifyDepositBlockInclusion) treasurytypes.WalletType {
	switch msg.ChainName {
	// Bitcoin chains
	case "mainnet":
		return treasurytypes.WalletType_WALLET_TYPE_BTC_MAINNET
	case "regtest", "regnet":
		return treasurytypes.WalletType_WALLET_TYPE_BTC_REGNET
	case "testnet", "testnet3", "testnet4":
		return treasurytypes.WalletType_WALLET_TYPE_BTC_TESTNET
	// Zcash chains - support both hyphenated and non-hyphenated forms
	case "zcash-mainnet", "zcashmainnet":
		return treasurytypes.WalletType_WALLET_TYPE_ZCASH_MAINNET
	case "zcash-testnet", "zcashtestnet", "zcash-regtest", "zcash-regnet", "zcashregtest", "zcashregnet":
		return treasurytypes.WalletType_WALLET_TYPE_ZCASH_TESTNET
	default:
		return treasurytypes.WalletType_WALLET_TYPE_UNSPECIFIED
	}
}

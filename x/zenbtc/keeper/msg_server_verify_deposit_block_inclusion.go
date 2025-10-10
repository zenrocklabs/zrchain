package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/rpc"

	"github.com/zenrocklabs/zenbtc/utils"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/bitcoin"
	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"

	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

func (k msgServer) VerifyDepositBlockInclusion(goCtx context.Context, msg *types.MsgVerifyDepositBlockInclusion) (*types.MsgVerifyDepositBlockInclusionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	blockHeader, err := k.validationKeeper.BtcBlockHeaders.Get(ctx, msg.BlockHeight)

	//CSM For Debugging Only
	//try and get missing Blockheader over RPC - WARNING for debugging only!!!!
	// if err != nil {
	// 	bh, _ := debugRetrieveBlockHeaderViaRPC(msg.ChainName, msg.BlockHeight)
	// 	if bh != nil {
	// 		err = nil
	// 		blockHeader = *bh
	// 	}
	// }
	// END of debug code

	if err != nil {
		return nil, err
	}

	found := false

	// List of addresses to ignore - we don't want to cause a mint for change
	ignoreAddresses, err := k.ZenBTCChangeAddresses(ctx, msg.ChainName)
	if err != nil {
		return nil, errors.New("error retrieving the change addresses")
	}

	// Verify the blockheader is valid and the proof, return a list of outputs in the transaction
	outputs, _, err := bitcoin.VerifyBTCLockTransaction(msg.RawTx, msg.ChainName, int(msg.Index), msg.Proof, &blockHeader, ignoreAddresses)
	if err != nil {
		return nil, err
	}

	// Check the address & amount specified is actually in the supplied (proven) BTC Transaction
	for _, output := range outputs {
		if output.Address == msg.DepositAddr && output.Amount == msg.Amount && uint64(output.OutputIndex) == msg.Vout {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("zenBTC deposit not found in outputs from provided lock tx")
	}

	q, err := k.treasuryKeeper.KeyByAddress(ctx, &treasurytypes.QueryKeyByAddressRequest{
		Address:     msg.DepositAddr,
		KeyringAddr: k.GetDepositKeyringAddr(ctx),
		KeyType:     treasurytypes.KeyType_KEY_TYPE_BITCOIN_SECP256K1,
		WalletType:  WalletTypeFromChainName(msg),
	})
	if err != nil {
		return nil, err
	}
	if q.Response == nil || q.Response.Wallets == nil {
		return nil, errors.New("zenBTC deposit address does not correspond to correct key (no wallets returned)")
	}

	metaData := q.Response.Key.ZenbtcMetadata
	if metaData == nil || metaData.RecipientAddr == "" || metaData.Caip2ChainId == "" {
		return nil, errors.New("lock tx - Key Metadata is invalid")
	}

	found = false
	for _, wallet := range q.Response.Wallets {
		if wallet.Address == msg.DepositAddr {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("zenBTC deposit address does not correspond to correct key (no matching wallet)")
	}

	// Deposit/lock txs are stored in zenBTC module so they can't be used to mint zenBTC tokens more than once
	toBeHashed := fmt.Sprintf("%s:%d", msg.RawTx, msg.Vout)
	hash := sha256.Sum256([]byte(toBeHashed))
	lockTxKey := hex.EncodeToString(hash[:])
	txExists, err := k.LockTransactions.Has(ctx, lockTxKey)
	if err != nil {
		return nil, err
	}
	if txExists {
		return nil, errors.New("lock tx was already previously used to mint zenBTC tokens")
	}

	// Get exchange rate before updating supply
	exchangeRate, err := k.GetExchangeRate(ctx)
	if err != nil {
		return nil, err
	}

	// Calculate zenBTC amount using current exchange rate
	zenBTCAmount := uint64(math.LegacyNewDecFromInt(math.NewIntFromUint64(msg.Amount)).Quo(exchangeRate).TruncateInt64())

	supply, err := k.Supply.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return nil, err
		}
		supply = types.Supply{CustodiedBTC: 0, MintedZenBTC: 0}
	}

	supply.CustodiedBTC += msg.Amount

	// Check if it's a deposit of converted rewards from EigenLayer
	isYieldDeposit := q.Response.Key.Id == k.GetRewardsDepositKeyID(ctx)

	// Only mint zenBTC for user deposits, not yield
	if !isYieldDeposit {
		supply.PendingZenBTC += zenBTCAmount
	}

	if err := k.Supply.Set(ctx, supply); err != nil {
		return nil, err
	}

	k.Logger().Warn("custodied supply updated", "custodied_old", supply.CustodiedBTC-msg.Amount, "custodied_new", supply.CustodiedBTC)
	if !isYieldDeposit {
		k.Logger().Warn("pending mint supply updated", "pending_mint_old", supply.PendingZenBTC-zenBTCAmount, "pending_mint_new", supply.PendingZenBTC)
	}

	if err := k.LockTransactions.Set(ctx, lockTxKey, types.LockTransaction{
		RawTx:         msg.RawTx,
		Vout:          msg.Vout,
		Sender:        msg.DepositAddr,
		MintRecipient: q.Response.Key.ZenbtcMetadata.RecipientAddr,
		Amount:        msg.Amount,
		BlockHeight:   msg.BlockHeight,
	}); err != nil {
		return nil, err
	}

	// Don't mint zenBTC tokens for rewards deposits; return early
	if isYieldDeposit {
		return &types.MsgVerifyDepositBlockInclusionResponse{}, nil
	}

	tx := &types.PendingMintTransaction{
		Caip2ChainId:     q.Response.Key.ZenbtcMetadata.Caip2ChainId,
		ChainType:        types.WalletType(q.Response.Key.ZenbtcMetadata.ChainType),
		RecipientAddress: q.Response.Key.ZenbtcMetadata.RecipientAddr,
		Amount:           zenBTCAmount,
		Creator:          msg.Creator,
	}
	if _, err := k.CreatePendingMintTransaction(ctx, tx); err != nil {
		return nil, err
	}
	k.validationKeeper.Logger(ctx).Warn("added pending mint transaction", "tx", fmt.Sprintf("%+v", tx))

// Request appropriate nonces to trigger minting on destination chains (no EigenLayer staking)
chainType := types.WalletType(q.Response.Key.ZenbtcMetadata.ChainType)
if chainType == types.WalletType_WALLET_TYPE_SOLANA {
	solParams := k.GetSolanaParams(ctx)
	if err := k.validationKeeper.SetSolanaRequestedNonce(ctx, solParams.NonceAccountKey, true); err != nil {
		return nil, err
	}
	if err := k.validationKeeper.SetSolanaZenBTCRequestedAccount(ctx, q.Response.Key.ZenbtcMetadata.RecipientAddr, true); err != nil {
		return nil, err
	}
} else if chainType == types.WalletType_WALLET_TYPE_EVM {
	if err := k.validationKeeper.EthereumNonceRequested.Set(ctx, k.GetEthMinterKeyID(ctx), true); err != nil {
		return nil, err
	}
} else {
	return nil, fmt.Errorf("unsupported destination chain type: %v", chainType)
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

/*
Get the list of Change KeyID and derive the addresses for the correct Chain
*/
func (k msgServer) ZenBTCChangeAddresses(ctx context.Context, chainName string) ([]string, error) {
	keyIDs := k.GetChangeAddressKeyIDs(ctx)
	chaincfg := utils.ChainFromString(chainName)
	result := make([]string, 0)
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
	return result, nil
}

func WalletTypeFromChainName(msg *types.MsgVerifyDepositBlockInclusion) treasurytypes.WalletType {
	switch msg.ChainName {
	case "mainnet":
		return treasurytypes.WalletType_WALLET_TYPE_BTC_MAINNET
	case "regtest", "regnet":
		return treasurytypes.WalletType_WALLET_TYPE_BTC_REGNET
	case "testnet", "testnet3", "testnet4":
		return treasurytypes.WalletType_WALLET_TYPE_BTC_TESTNET
	default:
		return treasurytypes.WalletType_WALLET_TYPE_UNSPECIFIED
	}
}

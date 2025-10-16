package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
)

func (k Keeper) Keys(
	goCtx context.Context,
	req *types.QueryKeysRequest,
) (*types.QueryKeysResponse, error) {
	if req == nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "invalid arguments: request is nil")
	}

	keys, pageRes, err := query.CollectionFilteredPaginate(
		goCtx,
		k.KeyStore,
		req.Pagination,
		func(key uint64, value types.Key) (bool, error) {
			keyType := types.KeyType_KEY_TYPE_UNSPECIFIED
			switch req.WalletType {
			case types.WalletType_WALLET_TYPE_NATIVE:
				fallthrough
			case types.WalletType_WALLET_TYPE_EVM:
				keyType = types.KeyType_KEY_TYPE_ECDSA_SECP256K1
			case types.WalletType_WALLET_TYPE_SOLANA:
				keyType = types.KeyType_KEY_TYPE_EDDSA_ED25519
			case types.WalletType_WALLET_TYPE_BTC_TESTNET:
				fallthrough
			case types.WalletType_WALLET_TYPE_BTC_MAINNET:
				fallthrough
			case types.WalletType_WALLET_TYPE_BTC_REGNET:
				fallthrough
			case types.WalletType_WALLET_TYPE_ZCASH_MAINNET:
				fallthrough
			case types.WalletType_WALLET_TYPE_ZCASH_TESTNET:
				keyType = types.KeyType_KEY_TYPE_BITCOIN_SECP256K1
			}

			workspaceMatch := (req.WorkspaceAddr == "" || value.WorkspaceAddr == req.WorkspaceAddr)
			keyTypeMatch := (value.Type == keyType || keyType == types.KeyType_KEY_TYPE_UNSPECIFIED)

			return workspaceMatch && keyTypeMatch, nil
		},
		func(key uint64, value types.Key) (*types.KeyAndWalletResponse, error) {
			return &types.KeyAndWalletResponse{
				Key: &types.KeyResponse{
					Id:             value.Id,
					WorkspaceAddr:  value.WorkspaceAddr,
					KeyringAddr:    value.KeyringAddr,
					Type:           value.Type.String(),
					PublicKey:      value.PublicKey,
					Index:          value.Index,
					SignPolicyId:   value.SignPolicyId,
					ZenbtcMetadata: value.ZenbtcMetadata,
				},
				Wallets: processWallets(value, req.WalletType, req.Prefixes),
			}, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &types.QueryKeysResponse{
		Keys:       keys,
		Pagination: pageRes,
	}, nil
}

func processWallets(
	key types.Key,
	walletType types.WalletType,
	prefixes []string,
) []*types.WalletResponse {

	wallets := []*types.WalletResponse{}
	walletTypes := deriveWalletTypes(walletType)

	for _, walletType := range walletTypes {
		switch walletType {
		case types.WalletType_WALLET_TYPE_NATIVE:
			addresses := deriveNativeAddresses(key, prefixes)
			for _, address := range addresses {
				wallets = append(
					wallets,
					&types.WalletResponse{Address: address, Type: walletType.String()},
				)
			}
		case types.WalletType_WALLET_TYPE_EVM:
			if address, err := types.EthereumAddress(&key); err == nil {
				wallets = append(
					wallets,
					&types.WalletResponse{Address: address, Type: walletType.String()},
				)
			}
		case types.WalletType_WALLET_TYPE_BTC_TESTNET:
			if address, err := types.BitcoinP2WPKH(&key, &chaincfg.TestNet3Params); err == nil {
				wallets = append(
					wallets,
					&types.WalletResponse{Address: address, Type: walletType.String()},
				)
			}
		case types.WalletType_WALLET_TYPE_BTC_MAINNET:
			if address, err := types.BitcoinP2WPKH(&key, &chaincfg.MainNetParams); err == nil {
				wallets = append(
					wallets,
					&types.WalletResponse{Address: address, Type: walletType.String()},
				)
			}
		case types.WalletType_WALLET_TYPE_BTC_REGNET:
			if address, err := types.BitcoinP2WPKH(&key, &chaincfg.RegressionNetParams); err == nil {
				wallets = append(
					wallets,
					&types.WalletResponse{Address: address, Type: walletType.String()},
				)
			}
		case types.WalletType_WALLET_TYPE_SOLANA:
			if pubKey, err := types.SolanaPubkey(&key); err == nil {
				wallets = append(
					wallets,
					&types.WalletResponse{Address: pubKey.String(), Type: walletType.String()},
				)
			}
		case types.WalletType_WALLET_TYPE_ZCASH_MAINNET:
			if address, err := types.ZcashAddress(&key, "mainnet"); err == nil {
				wallets = append(
					wallets,
					&types.WalletResponse{Address: address, Type: walletType.String()},
				)
			}
		case types.WalletType_WALLET_TYPE_ZCASH_TESTNET:
			if address, err := types.ZcashAddress(&key, "testnet"); err == nil {
				wallets = append(
					wallets,
					&types.WalletResponse{Address: address, Type: walletType.String()},
				)
			}
		}
	}

	return wallets
}

func deriveWalletTypes(walletType types.WalletType) []types.WalletType {
	switch walletType {
	case types.WalletType_WALLET_TYPE_NATIVE:
		return []types.WalletType{types.WalletType_WALLET_TYPE_NATIVE}
	case types.WalletType_WALLET_TYPE_EVM:
		return []types.WalletType{types.WalletType_WALLET_TYPE_EVM}
	case types.WalletType_WALLET_TYPE_BTC_TESTNET:
		return []types.WalletType{types.WalletType_WALLET_TYPE_BTC_TESTNET}
	case types.WalletType_WALLET_TYPE_BTC_MAINNET:
		return []types.WalletType{types.WalletType_WALLET_TYPE_BTC_MAINNET}
	case types.WalletType_WALLET_TYPE_BTC_REGNET:
		return []types.WalletType{types.WalletType_WALLET_TYPE_BTC_REGNET}
	case types.WalletType_WALLET_TYPE_SOLANA:
		return []types.WalletType{types.WalletType_WALLET_TYPE_SOLANA}
	case types.WalletType_WALLET_TYPE_ZCASH_MAINNET:
		return []types.WalletType{types.WalletType_WALLET_TYPE_ZCASH_MAINNET}
	case types.WalletType_WALLET_TYPE_ZCASH_TESTNET:
		return []types.WalletType{types.WalletType_WALLET_TYPE_ZCASH_TESTNET}
	case types.WalletType_WALLET_TYPE_UNSPECIFIED:
		return []types.WalletType{
			types.WalletType_WALLET_TYPE_NATIVE,
			types.WalletType_WALLET_TYPE_EVM,
			types.WalletType_WALLET_TYPE_BTC_TESTNET,
			types.WalletType_WALLET_TYPE_BTC_MAINNET,
			types.WalletType_WALLET_TYPE_BTC_REGNET,
			types.WalletType_WALLET_TYPE_SOLANA,
			types.WalletType_WALLET_TYPE_ZCASH_MAINNET,
			types.WalletType_WALLET_TYPE_ZCASH_TESTNET,
		}
	default:
		return []types.WalletType{}
	}
}

func deriveNativeAddresses(key types.Key, prefixes []string) []string {
	var addresses []string
	if len(prefixes) == 0 {
		if address, err := types.NativeAddress(&key, ""); err == nil {
			addresses = append(addresses, address)
		}
	}
	for _, prefix := range prefixes {
		if address, err := types.NativeAddress(&key, prefix); err == nil {
			addresses = append(addresses, address)
		}
	}
	return addresses
}

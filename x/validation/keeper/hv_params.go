package keeper

import (
	"context"

	"cosmossdk.io/math"

	"github.com/Zenrock-Foundation/zrchain/v5/x/validation/types"
	"github.com/ethereum/go-ethereum/common"
)

var (
	DefaultAVSRewardsRate, _                = math.LegacyNewDecFromStr("0.03") // 0.03 == 3% APR
	DefaultBlockTime                 int64  = 1                                // seconds
	DefaultZenBTCEthBatcherAddr             = "0x17361a5050258cCeffD595957cB8fddF79cEeeEB"
	DefaultZenBTCDepositKeyringAddr         = "keyring1k6vc6vhp6e6l3rxalue9v4ux"
	DefaultZenBTCWithdrawerKeyID     uint64 = 1
	DefaultZenBTCMinterKeyID         uint64 = 2
	DefaultZenBTCChangeAddressKeyIDs        = []uint64{3}
	DefaultZenBTCUnstakerKeyID       uint64 = 4
	DefaultZenBTCBurnerKeyID         uint64 = 5
	DefaultZenBTCRewardsDepositKeyID uint64 = 6
	DefaultBitcoinProxyCreatorID            = "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
	DefaultZenBTCStrategyAddr               = "0x0000000000000000000000000000000000000000"
	DefaultStakeableAssets                  = []*types.AssetData{
		{Asset: types.Asset_ROCK, Precision: 6},
		{Asset: types.Asset_zenBTC, Precision: 8},
		{Asset: types.Asset_stETH, Precision: 18},
	}
)

// NewParams creates a new Params instance
func NewHVParams(avsRewardsRate math.LegacyDec, blockTime int64, zenBTCParams *types.ZenBTCParams) *types.HVParams {
	return &types.HVParams{
		AVSRewardsRate: avsRewardsRate,
		BlockTime:      blockTime,
		ZenBTCParams:   zenBTCParams,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultHVParams() *types.HVParams {
	return NewHVParams(
		DefaultAVSRewardsRate,
		DefaultBlockTime,
		&types.ZenBTCParams{
			ZenBTCEthBatcherAddr:      DefaultZenBTCEthBatcherAddr,
			ZenBTCDepositKeyringAddr:  DefaultZenBTCDepositKeyringAddr,
			ZenBTCMinterKeyID:         DefaultZenBTCMinterKeyID,
			ZenBTCWithdrawerKeyID:     DefaultZenBTCWithdrawerKeyID,
			ZenBTCUnstakerKeyID:       DefaultZenBTCUnstakerKeyID,
			ZenBTCBurnerKeyID:         DefaultZenBTCBurnerKeyID,
			ZenBTCRewardsDepositKeyID: DefaultZenBTCRewardsDepositKeyID,
			BitcoinProxyCreatorID:     DefaultBitcoinProxyCreatorID,
			ZenBTCChangeAddressKeyIDs: DefaultZenBTCChangeAddressKeyIDs,
			StakeableAssets:           DefaultStakeableAssets,
		},
	)
}

func (k Keeper) GetAVSRewardsRate(ctx context.Context) math.LegacyDec {
	params, err := k.HVParams.Get(ctx)
	if err != nil {
		return DefaultAVSRewardsRate
	}
	return params.AVSRewardsRate
}

func (k Keeper) GetBlockTime(ctx context.Context) int64 {
	params, err := k.HVParams.Get(ctx)
	if err != nil {
		return DefaultBlockTime
	}
	return params.BlockTime
}

func (k Keeper) GetZenBTCEthBatcherAddr(ctx context.Context) string {
	params, err := k.HVParams.Get(ctx)
	if err != nil {
		return DefaultZenBTCEthBatcherAddr
	}
	return params.ZenBTCParams.ZenBTCEthBatcherAddr
}

func (k Keeper) GetZenBTCDepositKeyringAddr(ctx context.Context) string {
	params, err := k.HVParams.Get(ctx)
	if err != nil {
		return DefaultZenBTCDepositKeyringAddr
	}
	return params.ZenBTCParams.ZenBTCDepositKeyringAddr
}

func (k Keeper) GetZenBTCMinterKeyID(ctx context.Context) uint64 {
	params, err := k.HVParams.Get(ctx)
	if err != nil {
		return DefaultZenBTCMinterKeyID
	}
	return params.ZenBTCParams.ZenBTCMinterKeyID
}

func (k Keeper) GetZenBTCUnstakerKeyID(ctx context.Context) uint64 {
	params, err := k.HVParams.Get(ctx)
	if err != nil {
		return DefaultZenBTCUnstakerKeyID
	}
	return params.ZenBTCParams.ZenBTCUnstakerKeyID
}

func (k Keeper) GetZenBTCBurnerKeyID(ctx context.Context) uint64 {
	params, err := k.HVParams.Get(ctx)
	if err != nil {
		return DefaultZenBTCBurnerKeyID
	}
	return params.ZenBTCParams.ZenBTCBurnerKeyID
}

func (k Keeper) GetZenBTCWithdrawerKeyID(ctx context.Context) uint64 {
	params, err := k.HVParams.Get(ctx)
	if err != nil {
		return DefaultZenBTCWithdrawerKeyID
	}
	return params.ZenBTCParams.ZenBTCWithdrawerKeyID
}

func (k Keeper) GetBitcoinProxyCreatorID(ctx context.Context) string {
	params, err := k.HVParams.Get(ctx)
	if err != nil {
		return DefaultBitcoinProxyCreatorID
	}
	return params.ZenBTCParams.BitcoinProxyCreatorID
}

func (k Keeper) GetZenBTCChangeAddressKeyIDs(ctx context.Context) []uint64 {
	params, err := k.HVParams.Get(ctx)
	if err != nil {
		return DefaultZenBTCChangeAddressKeyIDs
	}
	return params.ZenBTCParams.ZenBTCChangeAddressKeyIDs
}

func (k Keeper) GetZenBTCRewardsDepositKeyID(ctx context.Context) uint64 {
	params, err := k.HVParams.Get(ctx)
	if err != nil {
		return DefaultZenBTCRewardsDepositKeyID
	}
	return params.ZenBTCParams.ZenBTCRewardsDepositKeyID
}

func (k Keeper) GetStakeableAssets(ctx context.Context) []*types.AssetData {
	params, err := k.HVParams.Get(ctx)
	if err != nil {
		return DefaultStakeableAssets
	}
	return params.ZenBTCParams.StakeableAssets
}

func (k Keeper) GetZenBTCStrategyAddr(ctx context.Context) common.Address {
	params, err := k.HVParams.Get(ctx)
	if err != nil {
		return common.HexToAddress(DefaultZenBTCStrategyAddr)
	}
	return common.HexToAddress(params.ZenBTCParams.ZenBTCStrategyAddr)
}

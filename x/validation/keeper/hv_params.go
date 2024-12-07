package keeper

import (
	"context"

	"cosmossdk.io/math"

	"github.com/Zenrock-Foundation/zrchain/v5/x/validation/types"
)

var (
	DefaultAVSRewardsRate, _                = math.LegacyNewDecFromStr("0.03") // 0.03 == 3% APR
	DefaultBlockTime                 int64  = 1                                // seconds
	DefaultZenBTCEthContractAddr            = "0x0832c25DcDD7E353749F50136a191377D9bA562e"
	DefaultZenBTCDepositKeyringAddr         = "keyring1k6vc6vhp6e6l3rxalue9v4ux"
	DefaultZenBTCWithdrawerKeyID     uint64 = 1
	DefaultZenBTCMinterKeyID         uint64 = 2
	DefaultZenBTCRewardsDepositKeyID uint64 = 3
	DefaultZenBTCChangeAddressKeyIDs        = []uint64{1} //[]uint64{4, 5}
	DefaultBitcoinProxyCreatorID            = "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
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
			ZenBTCEthContractAddr:     DefaultZenBTCEthContractAddr,
			ZenBTCDepositKeyringAddr:  DefaultZenBTCDepositKeyringAddr,
			ZenBTCMinterKeyID:         DefaultZenBTCMinterKeyID,
			ZenBTCWithdrawerKeyID:     DefaultZenBTCWithdrawerKeyID,
			ZenBTCRewardsDepositKeyID: DefaultZenBTCRewardsDepositKeyID,
			BitcoinProxyCreatorID:     DefaultBitcoinProxyCreatorID,
			ZenBTCChangeAddressKeyIDs: DefaultZenBTCChangeAddressKeyIDs,
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

func (k Keeper) GetZenBTCEthContractAddr(ctx context.Context) string {
	params, err := k.HVParams.Get(ctx)
	if err != nil {
		return DefaultZenBTCEthContractAddr
	}
	return params.ZenBTCParams.ZenBTCEthContractAddr
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

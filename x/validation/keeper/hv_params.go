package keeper

import (
	"context"

	"cosmossdk.io/math"

	"github.com/Zenrock-Foundation/zrchain/v5/x/validation/types"
)

var (
	DefaultAVSRewardsRate, _                  = math.LegacyNewDecFromStr("0.03") // 0.03 == 3% APR
	DefaultBlockTime                 int64    = 1                                // seconds
	DefaultZenBTCEthContractAddr              = "0x0832c25DcDD7E353749F50136a191377D9bA562e"
	DefaultZenBTCDepositKeyringAddr           = "keyring1k6vc6vhp6e6l3rxalue9v4ux"
	DefaultZenBTCMinterKeyID         uint64   = 2
	DefaultZenBTCWithdrawerKeyID     uint64   = 1
	DefaultZenBTCChangeAddressKeyIDs []uint64 = []uint64{1, 2}
	DefaultBitcoinProxyCreatorID              = "zen13y3tm68gmu9kntcxwvmue82p6akacnpt2v7nty"
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
			ZenBTCEthContractAddr:    DefaultZenBTCEthContractAddr,
			ZenBTCDepositKeyringAddr: DefaultZenBTCDepositKeyringAddr,
			ZenBTCMinterKeyID:        DefaultZenBTCMinterKeyID,
			ZenBTCWithdrawerKeyID:    DefaultZenBTCWithdrawerKeyID,
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

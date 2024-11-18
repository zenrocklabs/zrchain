package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"

	"github.com/Zenrock-Foundation/zrchain/v5/x/mint/types"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the mint store
type Keeper struct {
	cdc            codec.BinaryCodec
	storeService   storetypes.KVStoreService
	stakingKeeper  types.StakingKeeper
	bankKeeper     types.BankKeeper
	treasuryKeeper types.TreasuryKeeper
	accountKeeper  types.AccountKeeper

	feeCollectorName string

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string

	Schema collections.Schema
	Params collections.Item[types.Params]
	Minter collections.Item[types.Minter]
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService,
	sk types.StakingKeeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	feeCollectorName string,
	authority string,
	tk types.TreasuryKeeper,
) Keeper {
	// ensure mint module account is set
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("the x/%s module account has not been set", types.ModuleName))
	}

	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:              cdc,
		storeService:     storeService,
		stakingKeeper:    sk,
		bankKeeper:       bk,
		feeCollectorName: feeCollectorName,
		authority:        authority,
		Params:           collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Minter:           collections.NewItem(sb, types.MinterKey, "minter", codec.CollValue[types.Minter](cdc)),
		treasuryKeeper:   tk,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema
	return k
}

// GetAuthority returns the x/mint module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", "x/"+types.ModuleName)
}

// StakingTokenSupply implements an alias call to the underlying staking keeper's
// StakingTokenSupply to be used in BeginBlocker.
func (k Keeper) StakingTokenSupply(ctx context.Context) (math.Int, error) {
	return k.stakingKeeper.StakingTokenSupply(ctx)
}

// BondedRatio implements an alias call to the underlying staking keeper's
// BondedRatio to be used in BeginBlocker.
func (k Keeper) BondedRatio(ctx context.Context) (math.LegacyDec, error) {
	return k.stakingKeeper.BondedRatio(ctx)
}

// MintCoins implements an alias call to the underlying supply keeper's
// MintCoins to be used in BeginBlocker.
func (k Keeper) MintCoins(ctx context.Context, newCoins sdk.Coins) error {
	if newCoins.Empty() {
		// skip as no coins need to be minted
		return nil
	}

	return k.bankKeeper.MintCoins(ctx, types.ModuleName, newCoins)
}

// AddCollectedFees implements an alias call to the underlying supply keeper's
// AddCollectedFees to be used in BeginBlocker.
func (k Keeper) AddCollectedFees(ctx context.Context, fees sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, fees)
}

func (k Keeper) TotalBondedTokens(ctx context.Context) (math.Int, error) {
	return k.stakingKeeper.TotalBondedTokens(ctx)
}

func (k Keeper) NextStakingReward(ctx context.Context, totalBondedTokens math.Int) (sdk.Coin, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}
	totalStakingReward := math.LegacyNewDecFromInt(totalBondedTokens).Mul(params.StakingYield).QuoInt(math.NewInt(int64(params.BlocksPerYear))).TruncateInt()

	return sdk.NewCoin(params.MintDenom, totalStakingReward), nil
}

func (k Keeper) ClaimKeyringFees(ctx context.Context) (sdk.Coin, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}
	bankKeeper := k.bankKeeper
	keyringRewards := bankKeeper.GetBalance(ctx, sdk.AccAddress(treasurytypes.KeyringCollectorName), params.MintDenom)
	bankKeeper.SendCoinsFromModuleToModule(ctx, treasurytypes.KeyringCollectorName, types.ModuleName, sdk.NewCoins(keyringRewards))

	return keyringRewards, nil
}

func (k Keeper) BaseDistribution(ctx context.Context, keyringRewards sdk.Coin) (sdk.Coin, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}

	burnAmount := math.LegacyNewDecFromInt(keyringRewards.Amount).Mul(params.BurnRate).TruncateInt()
	err = k.burnRewards(ctx, sdk.NewCoin(params.MintDenom, burnAmount))
	if err != nil {
		return sdk.Coin{}, err
	}

	keyringRewards.Amount = keyringRewards.Amount.Sub(burnAmount)

	protocolWalletPortion := math.LegacyNewDecFromInt(keyringRewards.Amount).Mul(params.ProtocolWalletRate).TruncateInt()
	err = k.sendProtocolWalletFees(ctx, sdk.NewCoin(params.MintDenom, protocolWalletPortion))
	if err != nil {
		return sdk.Coin{}, err
	}

	keyringRewards.Amount = keyringRewards.Amount.Sub(protocolWalletPortion)

	return keyringRewards, nil
}

func (k Keeper) burnRewards(ctx context.Context, rewards sdk.Coin) error {
	return k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(rewards))
}

func (k Keeper) sendProtocolWalletFees(ctx context.Context, protocolWalletPortion sdk.Coin) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(params.ProtocolWalletAddress), sdk.NewCoins(protocolWalletPortion))
}

func (k Keeper) CalculateTopUp(ctx context.Context, stakingRewards sdk.Coin, keyringRewardRest sdk.Coin) (sdk.Coin, error) {
	// Ceiling the amount to the nearest integer
	topUpAmount := stakingRewards.Amount.Sub(keyringRewardRest.Amount).ToLegacyDec().Ceil().TruncateInt()
	if topUpAmount.IsZero() || topUpAmount.IsNegative() {
		return sdk.Coin{}, fmt.Errorf("topUpAmount cannot be negative: %v", topUpAmount)
	}

	return sdk.NewCoin(stakingRewards.Denom, topUpAmount), nil
}

func (k Keeper) TopUpKeyringRewards(ctx context.Context, topUpAmount sdk.Coin) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}
	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.AccAddress(params.ProtocolWalletAddress), types.ModuleName, sdk.NewCoins(topUpAmount))
}

func (k Keeper) CheckModuleBalance(ctx context.Context, totalBlockStakingReward sdk.Coin) error {
	moduleBalance := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), totalBlockStakingReward.Denom)
	if moduleBalance.Amount.LT(totalBlockStakingReward.Amount) {
		return fmt.Errorf("module balance %v is less than required staking reward %v", moduleBalance, totalBlockStakingReward)
	}
	return nil
}

func (k Keeper) CalculateExcess(ctx context.Context, totalBlockStakingReward sdk.Coin, keyringRewardsRest sdk.Coin) (sdk.Coin, error) {
	excess := totalBlockStakingReward.Amount.Sub(keyringRewardsRest.Amount)
	if excess.IsZero() || excess.IsNegative() {
		return sdk.Coin{}, fmt.Errorf("excess cannot be negative: %v", excess)
	}
	return sdk.NewCoin(totalBlockStakingReward.Denom, excess), nil
}

func (k Keeper) AdditionalBurn(ctx context.Context, excess sdk.Coin) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	burnAmount := math.LegacyNewDecFromInt(excess.Amount).Mul(params.BurnRate).TruncateInt()
	return k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(params.MintDenom, burnAmount)))
}

func (k Keeper) AdditionalMpcRewards(ctx context.Context, excess sdk.Coin) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	mpcRewards := math.LegacyNewDecFromInt(excess.Amount).Mul(params.AdditionalMpcRewards).TruncateInt()
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(params.ProtocolWalletAddress), sdk.NewCoins(sdk.NewCoin(params.MintDenom, mpcRewards)))
}

func (k Keeper) AdditionalStakingRewards(ctx context.Context, excess sdk.Coin) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	stakingRewards := math.LegacyNewDecFromInt(excess.Amount).Mul(params.AdditionalStakingRewards).TruncateInt()
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, sdk.NewCoins(sdk.NewCoin(params.MintDenom, stakingRewards)))
}

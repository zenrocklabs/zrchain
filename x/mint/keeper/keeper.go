package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"

	"github.com/Zenrock-Foundation/zrchain/v6/x/mint/types"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	zenextypes "github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
	zentptypes "github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the mint store
type Keeper struct {
	cdc           codec.BinaryCodec
	storeService  storetypes.KVStoreService
	stakingKeeper types.StakingKeeper
	bankKeeper    types.BankKeeper
	accountKeeper types.AccountKeeper
	zentpKeeper   types.ZentpKeeper
	zenexKeeper   types.ZenexKeeper

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
	zk types.ZentpKeeper,
	zxk types.ZenexKeeper,
	feeCollectorName string,
	authority string,
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
		accountKeeper:    ak,
		bankKeeper:       bk,
		zentpKeeper:      zk,
		zenexKeeper:      zxk,
		feeCollectorName: feeCollectorName,
		authority:        authority,
		Params:           collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Minter:           collections.NewItem(sb, types.MinterKey, "minter", codec.CollValue[types.Minter](cdc)),
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
	keyringAddr := k.accountKeeper.GetModuleAddress(treasurytypes.KeyringCollectorName)
	keyringRewards := bankKeeper.GetBalance(ctx, keyringAddr, params.MintDenom)
	err = bankKeeper.SendCoinsFromModuleToModule(ctx, treasurytypes.KeyringCollectorName, types.ModuleName, sdk.NewCoins(keyringRewards))
	if err != nil {
		return sdk.Coin{}, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyKeyringRewards, keyringRewards.String()),
		),
	)

	return keyringRewards, nil
}

func (k Keeper) ClaimZentpFees(ctx context.Context) (sdk.Coin, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}
	bankKeeper := k.bankKeeper
	zentpAddr := k.accountKeeper.GetModuleAddress(zentptypes.ZentpCollectorName)
	zentpRewards := bankKeeper.GetBalance(ctx, zentpAddr, params.MintDenom)
	err = bankKeeper.SendCoinsFromModuleToModule(ctx, zentptypes.ZentpCollectorName, types.ModuleName, sdk.NewCoins(zentpRewards))
	if err != nil {
		return sdk.Coin{}, err
	}

	if zentpRewards.Amount.IsPositive() {
		err = k.zentpKeeper.UpdateZentpFees(ctx, zentpRewards.Amount.Uint64())
		if err != nil {
			return sdk.Coin{}, err
		}
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyZentpFees, zentpRewards.String()),
		),
	)

	return zentpRewards, nil
}

// sends urock from zentp fee collector module account
// to zenex fee collector module account
func (k Keeper) DistributeZentpFeesToZenexFeeCollector(ctx context.Context) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	zentpAddr := k.accountKeeper.GetModuleAddress(zentptypes.ZentpCollectorName)
	zentpRewards := k.bankKeeper.GetBalance(ctx, zentpAddr, params.MintDenom)

	if zentpRewards.Amount.IsPositive() {
		err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, zentptypes.ZentpCollectorName, zenextypes.ZenexFeeCollectorName, sdk.NewCoins(zentpRewards))
		if err != nil {
			return err
		}

		err = k.zentpKeeper.UpdateZentpFees(ctx, zentpRewards.Amount.Uint64())
		if err != nil {
			return err
		}
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyZentpFees, zentpRewards.String()),
		),
	)

	return nil
}

func (k Keeper) ClaimTxFees(ctx context.Context) (sdk.Coin, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}
	bankKeeper := k.bankKeeper
	feeCollectorAddr := k.accountKeeper.GetModuleAddress(k.feeCollectorName)
	feesAmount := bankKeeper.GetBalance(ctx, feeCollectorAddr, params.MintDenom)
	err = bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, types.ModuleName, sdk.NewCoins(feesAmount))
	if err != nil {
		return sdk.Coin{}, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyTxFees, feesAmount.String()),
		),
	)

	return feesAmount, nil
}

func (k Keeper) ZenexFeeProcessing(ctx context.Context) (sdk.Coin, error) {

	params, err := k.Params.Get(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}

	zenexFeeCollectorAddr := k.accountKeeper.GetModuleAddress(zenextypes.ZenexFeeCollectorName)
	zenexFeeCollectorBalance := k.bankKeeper.GetBalance(ctx, zenexFeeCollectorAddr, params.MintDenom)

	zrWalletPortion := math.LegacyNewDecFromInt(zenexFeeCollectorBalance.Amount).Mul(params.ZenbtcRewardRate).TruncateInt()
	zenbtcRewardPortion := math.LegacyNewDecFromInt(zenexFeeCollectorBalance.Amount).Mul(params.ZenbtcRewardRate).TruncateInt()

	// Convert string address to AccAddress
	zrWalletAddr, err := sdk.AccAddressFromBech32(params.ZrWalletAddress)
	if err != nil {
		return sdk.Coin{}, fmt.Errorf("invalid zr wallet address: %v", err)
	}

	err = k.sendFeesFromZenexFeeCollector(ctx, sdk.NewCoin(params.MintDenom, zrWalletPortion), zrWalletAddr)
	if err != nil {
		return sdk.Coin{}, err
	}

	zenbtcRewardAddr := k.accountKeeper.GetModuleAddress(zenextypes.ZenBtcRewardsCollectorName)

	err = k.sendFeesFromZenexFeeCollector(ctx, sdk.NewCoin(params.MintDenom, zenbtcRewardPortion), zenbtcRewardAddr)
	if err != nil {
		return sdk.Coin{}, err
	}

	leftoverBalance := zenexFeeCollectorBalance.Amount.Sub(zrWalletPortion).Sub(zenbtcRewardPortion)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyZrWalletPortion, zrWalletPortion.String()),
			sdk.NewAttribute(types.AttributeKeyZenbtcRewardPortion, zenbtcRewardPortion.String()),
		),
	)

	return sdk.NewCoin(params.MintDenom, leftoverBalance), nil
}

// func (k Keeper) burnRewards(ctx context.Context, rewards sdk.Coin) error {
// 	return k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(rewards))
// }

// sends a portion of urock from zenex fee collector module account
// to zr wallet address
func (k Keeper) sendFeesFromZenexFeeCollector(ctx context.Context, feePortion sdk.Coin, recipientAddr sdk.AccAddress) error {

	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, zenextypes.ZenexFeeCollectorName, recipientAddr, sdk.NewCoins(feePortion))
}

func (k Keeper) CalculateTopUp(ctx context.Context, stakingRewards sdk.Coin, totalRewardRest sdk.Coin) (sdk.Coin, error) {
	// Ceiling the amount to the nearest integer
	topUpAmount := stakingRewards.Amount.Sub(totalRewardRest.Amount).ToLegacyDec().Ceil().TruncateInt()
	if topUpAmount.IsNegative() {
		return sdk.Coin{}, fmt.Errorf("topUpAmount cannot be negative: %v", topUpAmount)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyTopUpAmount, topUpAmount.String()),
		),
	)

	return sdk.NewCoin(stakingRewards.Denom, topUpAmount), nil
}

func (k Keeper) CheckModuleBalance(ctx context.Context, totalBlockStakingReward sdk.Coin) error {

	// Validate input
	if totalBlockStakingReward.IsZero() {
		return fmt.Errorf("staking reward cannot be zero")
	}

	// Get the module address
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)

	if moduleAddr == nil {
		return fmt.Errorf("module address for %s not found", types.ModuleName)
	}

	// Add debug logging
	moduleBalance := k.bankKeeper.GetBalance(ctx, moduleAddr, totalBlockStakingReward.Denom)

	if moduleBalance.Amount.LT(totalBlockStakingReward.Amount) {
		return fmt.Errorf("module balance %v is less than required staking reward %v", moduleBalance, totalBlockStakingReward)
	}

	return nil
}

func (k Keeper) CalculateExcess(ctx context.Context, totalBlockStakingReward sdk.Coin, totalRewardsRest sdk.Coin) (sdk.Coin, error) {
	excess := totalRewardsRest.Amount.Sub(totalBlockStakingReward.Amount)
	if excess.IsZero() || excess.IsNegative() {
		return sdk.Coin{}, fmt.Errorf("excess cannot be negative: %v", excess)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeySurplusAmount, excess.String()),
		),
	)

	return sdk.NewCoin(totalBlockStakingReward.Denom, excess), nil
}

func (k Keeper) AdditionalBurn(ctx context.Context, excess sdk.Coin) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	burnAmount := math.LegacyNewDecFromInt(excess.Amount).Mul(params.AdditionalBurnRate).TruncateInt()
	excess.Amount = excess.Amount.Sub(burnAmount)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyAdditionalBurnAmountRewards, burnAmount.String()),
		),
	)
	return k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(params.MintDenom, burnAmount)))
}

func (k Keeper) AdditionalMpcRewards(ctx context.Context, excess sdk.Coin) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	// Convert string address to AccAddress
	protocolAddr, err := sdk.AccAddressFromBech32(params.ProtocolWalletAddress)
	if err != nil {
		return fmt.Errorf("invalid protocol wallet address: %v", err)
	}

	mpcRewards := math.LegacyNewDecFromInt(excess.Amount).Mul(params.AdditionalMpcRewards).TruncateInt()
	excess.Amount = excess.Amount.Sub(mpcRewards)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyAdditionalMpcRewards, mpcRewards.String()),
		),
	)

	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, protocolAddr, sdk.NewCoins(sdk.NewCoin(params.MintDenom, mpcRewards)))
}

func (k Keeper) AdditionalStakingRewards(ctx context.Context, excess sdk.Coin) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	stakingRewards := math.LegacyNewDecFromInt(excess.Amount).Mul(params.AdditionalStakingRewards).TruncateInt()
	excess.Amount = excess.Amount.Sub(stakingRewards)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyAdditionalStakingRewards, stakingRewards.String()),
		),
	)

	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, sdk.NewCoins(sdk.NewCoin(params.MintDenom, stakingRewards)))
}

func (k Keeper) ExcessDistribution(ctx context.Context, excessAmount sdk.Coin) error {

	err := k.AdditionalBurn(ctx, excessAmount)
	if err != nil {
		return err
	}

	err = k.AdditionalMpcRewards(ctx, excessAmount)
	if err != nil {
		return err
	}

	err = k.AdditionalStakingRewards(ctx, excessAmount)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetMintModuleBalance(ctx context.Context) (sdk.Coin, error) {

	params, err := k.Params.Get(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}

	mintModuleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	mintModuleBalance := k.bankKeeper.GetBalance(ctx, mintModuleAddr, params.MintDenom)

	return mintModuleBalance, nil
}

// claims total rewards from keyring fees, tx fees
// and sends them to the mint module account
func (k Keeper) ClaimTotalRewards(ctx context.Context) (sdk.Coin, error) {
	keyringRewards, err := k.ClaimKeyringFees(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}

	feesAmount, err := k.ClaimTxFees(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}

	return keyringRewards.Add(feesAmount), nil
}

func (k Keeper) GetModuleAccountPerms(ctx context.Context) []string {

	moduleAccount := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	return moduleAccount.GetPermissions()
}

func (k Keeper) GetParams(ctx context.Context) (types.Params, error) {
	return k.Params.Get(ctx)
}

func (k Keeper) FundMintModuleFromZenexFeeCollector(ctx context.Context, amount sdk.Coin) error {
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, zenextypes.ZenexFeeCollectorName, types.ModuleName, sdk.NewCoins(amount))
}

func (k Keeper) HandleZenTpFees(ctx context.Context) error {
	// distributes zentp fees to zenex fee collector
	err := k.DistributeZentpFeesToZenexFeeCollector(ctx)
	if err != nil {
		return err
	}

	// distributes zenex fees to zr wallet and zenbtc reward collector
	zenexRewardsRest, err := k.ZenexFeeProcessing(ctx)
	if err != nil {
		return err
	}

	// funds the mint module from leftoverzenex fee collector balance
	err = k.FundMintModuleFromZenexFeeCollector(ctx, zenexRewardsRest)
	if err != nil {
		return err
	}

	return nil
}

// swaps ROCK for BTC and funds the zenbtc reward address with BTC
// if the rock fee pool balance is greater than the required rock balance
func (k Keeper) CheckZenBtcSwapThreshold(goctx context.Context) error {

	ctx := sdk.UnwrapSDKContext(goctx)
	requiredRockBalance, err := k.zenexKeeper.GetRequiredRockBalance(ctx)
	if err != nil {
		return err
	}

	zenbtcRewardsCollectorBalance := k.zenexKeeper.GetZenBtcRewardsCollectorBalance(ctx)
	if zenbtcRewardsCollectorBalance >= requiredRockBalance {
		err = k.zenexKeeper.CreateRockBtcSwap(ctx, zenbtcRewardsCollectorBalance)
		if err != nil {
			return err
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyRockSwapForZenBtcRewardsInitiated, fmt.Sprintf("Emitted Rock Swap for ZenBtc Rewards for %d ROCK", zenbtcRewardsCollectorBalance)),
		),
	)

	return nil
}

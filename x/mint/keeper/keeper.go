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
		accountKeeper:    ak,
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
	keyringAddr := k.accountKeeper.GetModuleAddress(treasurytypes.KeyringCollectorName)
	keyringRewards := bankKeeper.GetBalance(ctx, keyringAddr, params.MintDenom)
	err = bankKeeper.SendCoinsFromModuleToModule(ctx, treasurytypes.KeyringCollectorName, types.ModuleName, sdk.NewCoins(keyringRewards))
	if err != nil {
		return sdk.Coin{}, err
	}

	return keyringRewards, nil
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

	return feesAmount, nil
}

func (k Keeper) BaseDistribution(ctx context.Context, totalRewards sdk.Coin) (sdk.Coin, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}

	burnAmount := math.LegacyNewDecFromInt(totalRewards.Amount).Mul(params.BurnRate).TruncateInt()
	err = k.burnRewards(ctx, sdk.NewCoin(params.MintDenom, burnAmount))
	if err != nil {
		return sdk.Coin{}, err
	}

	// TODO - remove
	fmt.Printf("burn amount:\t\t\t%v\n", burnAmount)

	protocolWalletPortion := math.LegacyNewDecFromInt(totalRewards.Amount).Mul(params.ProtocolWalletRate).TruncateInt()

	// TODO - remove
	fmt.Printf("for protocol wallet:\t\t%v\n", protocolWalletPortion)
	err = k.sendProtocolWalletFees(ctx, sdk.NewCoin(params.MintDenom, protocolWalletPortion))
	if err != nil {
		return sdk.Coin{}, err
	}

	totalRewards.Amount = totalRewards.Amount.Sub(protocolWalletPortion).Sub(burnAmount)

	return totalRewards, nil
}

func (k Keeper) burnRewards(ctx context.Context, rewards sdk.Coin) error {
	return k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(rewards))
}

func (k Keeper) sendProtocolWalletFees(ctx context.Context, protocolWalletPortion sdk.Coin) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	// Convert string address to AccAddress
	protocolAddr, err := sdk.AccAddressFromBech32(params.ProtocolWalletAddress)
	if err != nil {
		return fmt.Errorf("invalid protocol wallet address: %v", err)
	}
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, protocolAddr, sdk.NewCoins(protocolWalletPortion))
}

func (k Keeper) CalculateTopUp(ctx context.Context, stakingRewards sdk.Coin, totalRewardRest sdk.Coin) (sdk.Coin, error) {
	// Ceiling the amount to the nearest integer
	topUpAmount := stakingRewards.Amount.Sub(totalRewardRest.Amount).ToLegacyDec().Ceil().TruncateInt()
	if topUpAmount.IsNegative() {
		return sdk.Coin{}, fmt.Errorf("topUpAmount cannot be negative: %v", topUpAmount)
	}

	return sdk.NewCoin(stakingRewards.Denom, topUpAmount), nil
}

// func (k Keeper) TopUpTotalRewards(ctx context.Context, topUpAmount sdk.Coin) error {
// 	params, err := k.Params.Get(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	// Convert string address to AccAddress
// 	protocolAddr, err := sdk.AccAddressFromBech32(params.ProtocolWalletAddress)
// 	if err != nil {
// 		return fmt.Errorf("invalid protocol wallet address: %v", err)
// 	}

// 	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, protocolAddr, types.ModuleName, sdk.NewCoins(topUpAmount))
// }

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

	// TODO - remove
	fmt.Printf("mint module balance:\t\t%v\n", moduleBalance)

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
	return sdk.NewCoin(totalBlockStakingReward.Denom, excess), nil
}

func (k Keeper) AdditionalBurn(ctx context.Context, excess sdk.Coin) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	burnAmount := math.LegacyNewDecFromInt(excess.Amount).Mul(params.BurnRate).TruncateInt()
	// TODO - remove
	fmt.Printf("excess burn amt:\t\t%v\n", burnAmount)
	excess.Amount = excess.Amount.Sub(burnAmount)
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
	// TODO - remove
	fmt.Printf("mpc rewards:\t\t\t%v\n", mpcRewards)
	excess.Amount = excess.Amount.Sub(mpcRewards)
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, protocolAddr, sdk.NewCoins(sdk.NewCoin(params.MintDenom, mpcRewards)))
}

func (k Keeper) AdditionalStakingRewards(ctx context.Context, excess sdk.Coin) error {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return err
	}

	stakingRewards := math.LegacyNewDecFromInt(excess.Amount).Mul(params.AdditionalStakingRewards).TruncateInt()
	// TODO - remove
	fmt.Printf("staking rewards:\t\t%v\n", stakingRewards)
	excess.Amount = excess.Amount.Sub(stakingRewards)
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, sdk.NewCoins(sdk.NewCoin(params.MintDenom, stakingRewards)))
}

func (k Keeper) ExcessDistribution(ctx context.Context, excessAmount sdk.Coin) error {

	mintModuleBalance, err := k.CheckMintModuleBalance(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("mint bal before burn:\t\t%v\n", mintModuleBalance)

	err = k.AdditionalBurn(ctx, excessAmount)
	if err != nil {
		return err
	}

	mintModuleBalance, err = k.CheckMintModuleBalance(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("mint after burn:\t\t%v\n", mintModuleBalance)

	err = k.AdditionalMpcRewards(ctx, excessAmount)
	if err != nil {
		return err
	}

	mintModuleBalance, err = k.CheckMintModuleBalance(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("mint after mpc rewards:\t\t%v\n", mintModuleBalance)

	err = k.AdditionalStakingRewards(ctx, excessAmount)
	if err != nil {
		return err
	}

	mintModuleBalance, err = k.CheckMintModuleBalance(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("mint after staking rewards:\t%v\n", mintModuleBalance)

	return nil
}

func (k Keeper) CheckMintModuleBalance(ctx context.Context) (sdk.Coin, error) {

	params, err := k.Params.Get(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}

	mintModuleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	mintModuleBalance := k.bankKeeper.GetBalance(ctx, mintModuleAddr, params.MintDenom)

	return mintModuleBalance, nil
}

func (k Keeper) ClaimTotalRewards(ctx context.Context) (sdk.Coin, error) {
	keyringRewards, err := k.ClaimKeyringFees(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}

	feesAmount, err := k.ClaimTxFees(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}

	// TODO: remove
	fmt.Printf("keyring rewards:\t\t%v\n", keyringRewards)
	fmt.Printf("fee rewards:\t\t\t%v\n", feesAmount)
	return keyringRewards.Add(feesAmount), nil
}

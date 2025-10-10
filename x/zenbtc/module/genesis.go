package zenbtc

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zenrocklabs/zenbtc/x/zenbtc/keeper"
	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.Params.Set(ctx, genState.Params); err != nil {
		panic(err)
	}

	for _, lockTx := range genState.LockTransactions {
		toBeHashed := fmt.Sprintf("%s:%d", lockTx.RawTx, lockTx.Vout)
		hash := sha256.Sum256([]byte(toBeHashed))
		lockTxKey := hex.EncodeToString(hash[:])
		if err := k.LockTransactions.Set(ctx, lockTxKey, lockTx); err != nil {
			panic(err)
		}
	}

	for _, pendingMintTx := range genState.PendingMintTransactions {
		if err := k.PendingMintTransactionsMap.Set(ctx, pendingMintTx.Id, pendingMintTx); err != nil {
			panic(err)
		}
	}

	if err := k.SetFirstPendingEthMintTransaction(ctx, genState.FirstPendingEthMintTransaction); err != nil {
		panic(err)
	}

	if err := k.SetFirstPendingSolMintTransaction(ctx, genState.FirstPendingSolMintTransaction); err != nil {
		panic(err)
	}

	if err := k.PendingMintTransactionCount.Set(ctx, genState.PendingMintTransactionCount); err != nil {
		panic(err)
	}

	for _, burnEvent := range genState.BurnEvents {
		if err := k.BurnEvents.Set(ctx, burnEvent.Id, burnEvent); err != nil {
			panic(err)
		}
	}

	if err := k.SetFirstPendingBurnEvent(ctx, genState.FirstPendingBurnEvent); err != nil {
		panic(err)
	}

	if err := k.BurnEventCount.Set(ctx, genState.BurnEventCount); err != nil {
		panic(err)
	}

	for _, redemption := range genState.Redemptions {
		if err := k.Redemptions.Set(ctx, redemption.Data.Id, redemption); err != nil {
			panic(err)
		}
	}

	if err := k.SetFirstPendingRedemption(ctx, genState.FirstPendingRedemption); err != nil {
		panic(err)
	}

	if err := k.SetFirstRedemptionAwaitingSign(ctx, genState.FirstRedemptionAwaitingSign); err != nil {
		panic(err)
	}

	if err := k.Supply.Set(ctx, genState.Supply); err != nil {
		panic(err)
	}

	if err := k.SetFirstPendingStakeTransaction(ctx, genState.FirstPendingStakeTransaction); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := keeper.DefaultGenesis()
	params, err := k.Params.Get(ctx)
	if err != nil {
		panic(err)
	}
	genesis.Params = params

	// this line is used by starport scaffolding # genesis/module/export

	err = k.ExportState(ctx, genesis)
	if err != nil {
		panic(err)
	}

	return genesis
}

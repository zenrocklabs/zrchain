package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func (k *Keeper) ProcessSolanaDCTMintEventsTestHelper(ctx sdk.Context, data OracleData) {
	k.processSolanaDCTMintEvents(ctx, data)
}

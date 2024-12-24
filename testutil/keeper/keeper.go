package keeper

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"

	identitykeeper "github.com/Zenrock-Foundation/zrchain/v5/x/identity/keeper"
	policykeeper "github.com/Zenrock-Foundation/zrchain/v5/x/policy/keeper"
	treasurykeeper "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/keeper"
)

type KeeperTest struct {
	Ctx            sdk.Context
	IdentityKeeper *identitykeeper.Keeper
	TreasuryKeeper *treasurykeeper.Keeper
	PolicyKeeper   *policykeeper.Keeper
}

func NewTest(t testing.TB) *KeeperTest {
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())

	policyKeeper, ctx := PolicyKeeper(t, db, stateStore, nil)
	identityKeeper, _ := IdentityKeeper(t, &policyKeeper, db, stateStore)
	treasuryKeeper, _ := TreasuryKeeper(t, &policyKeeper, &identityKeeper, nil, nil, db, stateStore)

	return &KeeperTest{
		Ctx:            ctx,
		IdentityKeeper: &identityKeeper,
		TreasuryKeeper: &treasuryKeeper,
		PolicyKeeper:   &policyKeeper,
	}
}

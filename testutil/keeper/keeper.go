package keeper

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"

	identitykeeper "github.com/Zenrock-Foundation/zrchain/v6/x/identity/keeper"
	policykeeper "github.com/Zenrock-Foundation/zrchain/v6/x/policy/keeper"
	treasurykeeper "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/keeper"
	treasurytestutil "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/testutil"
	"github.com/golang/mock/gomock"
)

type KeeperTest struct {
	Ctx            sdk.Context
	IdentityKeeper *identitykeeper.Keeper
	TreasuryKeeper *treasurykeeper.Keeper
	PolicyKeeper   *policykeeper.Keeper
	ZentpKeeper    treasurytestutil.MockZentpKeeper
}

func NewTest(t testing.TB) *KeeperTest {
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())

	ctrl := gomock.NewController(t)
	zentpKeeper := treasurytestutil.NewMockZentpKeeper(ctrl)

	policyKeeper, ctx := PolicyKeeper(t, db, stateStore, nil)
	identityKeeper, _ := IdentityKeeper(t, &policyKeeper, db, stateStore)
	treasuryKeeper, _ := TreasuryKeeper(t, &policyKeeper, &identityKeeper, nil, db, stateStore, zentpKeeper)

	return &KeeperTest{
		Ctx:            ctx,
		IdentityKeeper: &identityKeeper,
		TreasuryKeeper: &treasuryKeeper,
		PolicyKeeper:   &policyKeeper,
		ZentpKeeper:    *zentpKeeper,
	}
}

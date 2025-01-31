package keeper

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"

	identitykeeper "github.com/Zenrock-Foundation/zrchain/v5/x/identity/keeper"
	policykeeper "github.com/Zenrock-Foundation/zrchain/v5/x/policy/keeper"
	treasurykeeper "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/keeper"
	"github.com/Zenrock-Foundation/zrchain/v5/x/treasury/testutil"
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ikmock := testutil.NewMockIdentityKeeper(ctrl)

	policyKeeper, ctx := PolicyKeeper(t, db, stateStore, nil)
	identityKeeper, _ := IdentityKeeper(t, &policyKeeper, db, stateStore)
	treasuryKeeper, _ := TreasuryKeeper(t, &policyKeeper, ikmock, nil, db, stateStore)

	return &KeeperTest{
		Ctx:            ctx,
		IdentityKeeper: &identityKeeper,
		TreasuryKeeper: &treasuryKeeper,
		PolicyKeeper:   &policyKeeper,
	}
}

package keeper_test

import (
	"fmt"
	"testing"

	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestExportGenesis(t *testing.T) {
	suite := new(ValidationKeeperTestSuite)
	suite.SetT(&testing.T{})
	keeper, _ := suite.ValidationKeeperSetupTest()
	ctx := sdk.UnwrapSDKContext(suite.ctx)

	resp := keeper.InitGenesis(ctx, testutil.TestGenesis())
	require.NotEmpty(t, resp)

	genesisState := keeper.ExportGenesis(ctx)
	require.NotNil(t, genesisState)
	fmt.Println(genesisState)
}

package simulation

import (
	"math/rand"

	"github.com/Zenrock-Foundation/zrchain/v6/testutil/sample"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zentp/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgBurn(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgBurn{
			Authority:     sample.GetAuthority(),
			ModuleAccount: "zentp",
			Denom:         "urock",
			Amount:        100,
		}

		// TODO: Handling the Burn simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "Burn simulation not implemented"), nil, nil
	}
}

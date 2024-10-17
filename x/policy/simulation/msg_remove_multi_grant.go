package simulation

import (
	"math/rand"

	"github.com/Zenrock-Foundation/zrchain/v4/x/policy/keeper"
	"github.com/Zenrock-Foundation/zrchain/v4/x/policy/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgRemoveMultiGrant(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgRemoveMultiGrant{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the MsgRemoveMultiGrant simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "MsgRemoveMultiGrant simulation not implemented"), nil, nil
	}
}

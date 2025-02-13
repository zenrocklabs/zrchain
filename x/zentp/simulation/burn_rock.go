package simulation

import (
	"math/rand"

	"github.com/Zenrock-Foundation/zrchain/v5/x/zentp/keeper"
	"github.com/Zenrock-Foundation/zrchain/v5/x/zentp/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgBurnRock(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgBurnRock{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the BurnRock simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "BurnRock simulation not implemented"), nil, nil
	}
}

package simulation

import (
	"math/rand"

	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/identity/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgAddKeyringAdmin(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgAddKeyringAdmin{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the AddKeyringAdmin simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "AddKeyringAdmin simulation not implemented"), nil, nil
	}
}

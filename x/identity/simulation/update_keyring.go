package simulation

import (
	"math/rand"

	"github.com/Zenrock-Foundation/zrchain/v4/x/identity/keeper"
	"github.com/Zenrock-Foundation/zrchain/v4/x/identity/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgUpdateKeyring(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgUpdateKeyring{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the UpdateKeyring simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "UpdateKeyring simulation not implemented"), nil, nil
	}
}

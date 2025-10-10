package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/Zenrock-Foundation/zrchain/v6/x/zenbtc/keeper"
	"github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
)

func SimulateMsgVerifyDepositBlockInclusion(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgVerifyDepositBlockInclusion{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the VerifyDepositBlockInclusion simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "VerifyDepositBlockInclusion simulation not implemented"), nil, nil
	}
}

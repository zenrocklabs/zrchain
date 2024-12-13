package keeper

import (
	"context"

	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v5/x/identity/types"
	treasury_types "github.com/Zenrock-Foundation/zrchain/v5/x/treasury/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) NewKeyring(goCtx context.Context, msg *types.MsgNewKeyring) (*types.MsgNewKeyringResponse, error) {
	keyring := &types.Keyring{
		Creator:        msg.Creator,
		Description:    msg.Description,
		Admins:         []string{msg.Creator},
		PartyThreshold: msg.PartyThreshold,
		KeyReqFee:      msg.KeyReqFee,
		SigReqFee:      msg.SigReqFee,
		IsActive:       true,
		DelegateFees:   msg.DelegateFees,
	}

	params, err := k.ParamStore.Get(goCtx)
	if err != nil {
		return nil, err
	}

	if params.KeyringCreationFee > 0 {
		if err := k.bankKeeper.SendCoinsFromAccountToModule(
			goCtx,
			sdk.MustAccAddressFromBech32(msg.Creator),
			treasury_types.KeyringCollectorName,
			sdk.NewCoins(sdk.NewCoin("urock", math.NewIntFromUint64(params.KeyringCreationFee))),
		); err != nil {
			return nil, err
		}
	}

	address, err := k.CreateKeyring(sdk.UnwrapSDKContext(goCtx), keyring)
	if err != nil {
		return nil, err
	}

	return &types.MsgNewKeyringResponse{
		Addr: address,
	}, nil
}

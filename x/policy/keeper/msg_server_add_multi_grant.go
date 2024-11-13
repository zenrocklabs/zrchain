package keeper

import (
	"context"

	"github.com/Zenrock-Foundation/zrchain/v5/x/policy/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
)

func (k msgServer) AddMultiGrant(goCtx context.Context, msg *types.MsgAddMultiGrant) (*types.MsgAddMultiGrantResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	for _, m := range msg.Msgs {
		auth := &authztypes.GenericAuthorization{
			Msg: m,
		}
		authAny, err := codectypes.NewAnyWithValue(auth)
		if err != nil {
			return nil, err
		}

		_, err = k.authzKeeper.Grant(goCtx, &authztypes.MsgGrant{
			Granter: msg.Creator,
			Grantee: msg.Grantee,
			Grant: authztypes.Grant{
				Authorization: authAny,
				Expiration:    nil,
			},
		})
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgAddMultiGrantResponse{}, nil
}

package testclient

import (
	"context"

	"cosmossdk.io/math"
	client "github.com/Zenrock-Foundation/zrchain/v6/go-client"
	"github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/gogoproto/proto"
)

func (c *TestClient) BankSend(ctx context.Context, from, to string, amount int64) (*banktypes.MsgSendResponse, error) {
	msg := &banktypes.MsgSend{
		FromAddress: from,
		ToAddress:   to,
		Amount:      types.NewCoins(types.NewCoin("urock", math.NewInt(amount))),
	}

	txres, err := c.executeTx(ctx, msg)
	if err != nil {
		return nil, err
	}

	res := &banktypes.MsgSendResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, err
}

func (c *TestClient) BankSendWithIdentity(ctx context.Context, identity client.Identity, from, to string, amount int64) (*banktypes.MsgSendResponse, error) {
	msg := &banktypes.MsgSend{
		FromAddress: from,
		ToAddress:   to,
		Amount:      types.NewCoins(types.NewCoin("urock", math.NewInt(amount))),
	}

	txres, _, err := c.executeTxWithIdentity(ctx, identity, msg)
	if err != nil {
		return nil, err
	}

	res := &banktypes.MsgSendResponse{}
	err = proto.Unmarshal(txres.MsgResponses[0].Value, res)

	return res, err
}

func (c *TestClient) GetBalance(ctx context.Context, address string, denom string) (uint64, error) {
	res, err := c.bqc.Balance(ctx, &banktypes.QueryBalanceRequest{
		Address: address,
		Denom:   denom,
	})

	if err != nil {
		return 0, err
	}

	return res.Balance.Amount.Uint64(), nil
}

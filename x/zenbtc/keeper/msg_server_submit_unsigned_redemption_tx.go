package keeper

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"

	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

func (k msgServer) SubmitUnsignedRedemptionTx(goCtx context.Context, msg *types.MsgSubmitUnsignedRedemptionTx) (*types.MsgSubmitUnsignedRedemptionTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.VerifyUnsignedRedemptionTX(ctx, msg)
	if err != nil {
		return nil, err
	}

	keyIDs := make([]uint64, len(msg.Inputs))
	hashes := make([]string, len(msg.Inputs))
	for i, input := range msg.Inputs {
		keyIDs[i] = input.Keyid
		hashes[i] = input.Hash
	}

	signReq := &treasurytypes.MsgNewSignatureRequest{
		Creator:        msg.Creator,
		KeyIds:         keyIDs,
		DataForSigning: strings.Join(hashes, ","), // hex string, each unsigned utxo is separated by comma
		CacheId:        msg.CacheId,
		ZenbtcTxBytes:  msg.Txbytes,
	}

	resp, err := k.treasuryKeeper.HandleSignatureRequest(ctx, signReq)
	if err != nil {
		return nil, err
	}

	for _, idx := range msg.RedemptionIndexes[1:] {
		redemption, err := k.Redemptions.Get(ctx, idx)
		if err != nil {
			return nil, err
		}

		redemption.Data.SignReqId = resp.SigReqId

		if err := k.Redemptions.Set(ctx, idx, redemption); err != nil {
			return nil, err
		}
	}

	return &types.MsgSubmitUnsignedRedemptionTxResponse{}, nil
}

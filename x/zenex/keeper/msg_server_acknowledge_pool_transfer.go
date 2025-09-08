package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/app/params"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/Zenrock-Foundation/zrchain/v6/x/zenex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AcknowledgePoolTransfer(goCtx context.Context, msg *types.MsgAcknowledgePoolTransfer) (*types.MsgAcknowledgePoolTransferResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Creator != k.GetParams(ctx).BtcProxyAddress {
		return nil, fmt.Errorf("message creator is not the btc proxy address: %s", msg.Creator)
	}

	if msg.Status != types.SwapStatus_SWAP_STATUS_COMPLETED && msg.Status != types.SwapStatus_SWAP_STATUS_REJECTED {
		return nil, fmt.Errorf("swap status is not completed or rejected: %s", msg.Status)
	}

	swap, err := k.SwapsStore.Get(ctx, msg.SwapId)
	if err != nil {
		return nil, err
	}

	if swap.Status != types.SwapStatus_SWAP_STATUS_SWAP_TRANSFER_REQUESTED {
		return nil, fmt.Errorf("swap status is not swap transfer requested: %s", swap.Status)
	}

	if swap.Pair != "btcrock" {

		rockKey, err := k.treasuryKeeper.GetKey(ctx, swap.RockKeyId)
		if err != nil {
			return nil, err
		}

		rockAddress, err := treasurytypes.NativeAddress(rockKey, "zen")
		if err != nil {
			return nil, err
		}

		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ZenexCollectorName, sdk.MustAccAddressFromBech32(rockAddress), sdk.NewCoins(sdk.NewCoin(params.BondDenom, math.NewIntFromUint64(swap.Data.AmountOut))))
		if err != nil {
			return nil, err
		}
	}

	swap.Status = msg.Status
	swap.SourceTxHash = msg.SourceTxHash
	err = k.SwapsStore.Set(ctx, swap.SwapId, swap)
	if err != nil {
		return nil, err
	}

	return &types.MsgAcknowledgePoolTransferResponse{}, nil
}

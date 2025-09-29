package keeper

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/exp/slices"

	"cosmossdk.io/collections"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/shared"
	sidecarapitypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	zenbtctypes "github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

func (k msgServer) ManuallyInputBitcoinHeader(ctx context.Context, msg *types.MsgManuallyInputBitcoinHeader) (*types.MsgManuallyInputBitcoinHeaderResponse, error) {
	if msg.Authority != shared.AdminAuthAddr && msg.Authority != k.authority {
		return nil, fmt.Errorf("invalid authority; expected %s or %s, got %s", shared.AdminAuthAddr, k.authority, msg.Authority)
	}

	if msg.Header.BlockHeight <= 0 {
		return nil, fmt.Errorf("block height must be greater than zero")
	}
	if msg.Header.MerkleRoot == "" {
		return nil, fmt.Errorf("merkle root must be provided")
	}
	if msg.Header.BlockHash == "" {
		return nil, fmt.Errorf("block hash must be provided")
	}

	requested, err := k.RequestedHistoricalBitcoinHeaders.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return nil, err
		}
		requested = zenbtctypes.RequestedBitcoinHeaders{}
	}

	latestHeight, err := k.LatestBtcHeaderHeight.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return nil, err
		}
		latestHeight = 0
	}

	header := sidecarapitypes.BTCBlockHeader{
		Version:     msg.Header.Version,
		PrevBlock:   msg.Header.PrevBlock,
		MerkleRoot:  msg.Header.MerkleRoot,
		TimeStamp:   msg.Header.Timestamp,
		Bits:        msg.Header.Bits,
		Nonce:       msg.Header.Nonce,
		BlockHash:   msg.Header.BlockHash,
		BlockHeight: msg.Header.BlockHeight,
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if newHeight, updated := k.processAndStoreBtcHeaderManual(sdkCtx, header.BlockHeight, &header, latestHeight); updated {
		latestHeight = newHeight
	}

	requested.Heights = slices.DeleteFunc(requested.Heights, func(h int64) bool {
		return h == header.BlockHeight
	})

	if err := k.RequestedHistoricalBitcoinHeaders.Set(ctx, requested); err != nil {
		return nil, err
	}

	if err := k.LatestBtcHeaderHeight.Set(ctx, latestHeight); err != nil {
		return nil, err
	}

	return &types.MsgManuallyInputBitcoinHeaderResponse{}, nil
}

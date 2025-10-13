package keeper

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"golang.org/x/exp/slices"

	"cosmossdk.io/collections"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/shared"
	sidecarapitypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	dcttypes "github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
)

func (k msgServer) ManuallyInputZcashHeader(ctx context.Context, msg *types.MsgManuallyInputZcashHeader) (*types.MsgManuallyInputZcashHeaderResponse, error) {
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

	requested, err := k.RequestedHistoricalZcashHeaders.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return nil, err
		}
		requested = dcttypes.RequestedZcashHeaders{}
	}

	latestHeight, err := k.LatestZcashHeaderHeight.Get(ctx)
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
		TimeStamp:   msg.Header.TimeStamp,
		Bits:        msg.Header.Bits,
		Nonce:       msg.Header.Nonce,
		BlockHash:   msg.Header.BlockHash,
		BlockHeight: msg.Header.BlockHeight,
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if newHeight, updated := k.processAndStoreZcashHeaderManual(sdkCtx, header.BlockHeight, &header, latestHeight); updated {
		latestHeight = newHeight
	}

	requested.Heights = slices.DeleteFunc(requested.Heights, func(h int64) bool {
		return h == header.BlockHeight
	})

	if err := k.RequestedHistoricalZcashHeaders.Set(ctx, requested); err != nil {
		return nil, err
	}

	if err := k.LatestZcashHeaderHeight.Set(ctx, latestHeight); err != nil {
		return nil, err
	}

	return &types.MsgManuallyInputZcashHeaderResponse{}, nil
}

// processAndStoreZcashHeaderManual processes a manually input Zcash header without triggering reorg checks.
// It returns the updated latestZcashHeaderHeight and a boolean indicating if it was updated.
func (k *Keeper) processAndStoreZcashHeaderManual(
	ctx sdk.Context,
	headerHeight int64,
	header *sidecarapitypes.BTCBlockHeader,
	latestZcashHeaderHeight int64,
) (int64, bool) {
	if headerHeight <= 0 || header == nil || header.MerkleRoot == "" {
		return latestZcashHeaderHeight, false
	}

	// Check if header already exists by comparing hashes
	existingHeader, err := k.ZcashBlockHeaders.Get(ctx, headerHeight)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		k.Logger(ctx).Error("error checking if zcash header exists", "type", "manual", "height", headerHeight, "error", err)
		return latestZcashHeaderHeight, false
	}

	// If header exists, compare hashes to see if it's different
	if err == nil {
		existingHash, err := deriveHash(existingHeader)
		if err != nil {
			k.Logger(ctx).Error("error deriving hash for existing zcash header", "type", "manual", "height", headerHeight, "error", err)
			return latestZcashHeaderHeight, false
		}

		newHash, err := deriveHash(*header)
		if err != nil {
			k.Logger(ctx).Error("error deriving hash for new zcash header", "type", "manual", "height", headerHeight, "error", err)
			return latestZcashHeaderHeight, false
		}

		if bytes.Equal(existingHash[:], newHash[:]) {
			return latestZcashHeaderHeight, false
		}
	}

	// Store the new header (either no header existed or hash is different)
	if err := k.ZcashBlockHeaders.Set(ctx, headerHeight, *header); err != nil {
		k.Logger(ctx).Error("error storing zcash header", "type", "manual", "height", headerHeight, "error", err)
		return latestZcashHeaderHeight, false
	}

	k.Logger(ctx).Info("stored new zcash header", "type", "manual", "height", headerHeight)

	// Update the latest height if this header is newer than what we had before
	if headerHeight > latestZcashHeaderHeight {
		if err := k.LatestZcashHeaderHeight.Set(ctx, headerHeight); err != nil {
			k.Logger(ctx).Error("error setting latest ZCash header height", "error", err)
		}
		return headerHeight, true
	}

	// Header was stored but it's not newer than our current latest height
	return latestZcashHeaderHeight, false
}

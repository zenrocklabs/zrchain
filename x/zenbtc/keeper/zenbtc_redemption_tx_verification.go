package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"slices"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Zenrock-Foundation/zrchain/v6/bitcoin"
	treasurytypes "github.com/Zenrock-Foundation/zrchain/v6/x/treasury/types"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/zenrocklabs/zenbtc/utils"
	"github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

func (k msgServer) VerifyUnsignedRedemptionTX(ctx sdk.Context, msg *types.MsgSubmitUnsignedRedemptionTx) error {
	// Check that the transaction creator is valid
	if err := k.checkRedemptionTXCreator(ctx, msg); err != nil {
		return fmt.Errorf("failed to verify transaction creator: %w", err)
	}

	// Check that the change address is valid and retrieve the transaction
	msgTX, err := k.checkChangeAddress(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to check change address: %w", err)
	}

	if err := k.verifyInputsInRedemption(ctx, msg.Inputs); err != nil {
		return fmt.Errorf("failed to verify input keys: %w", err)
	}

	// Verify that the outputs in the supplied BTC TX all match redemptions
	if err := k.verifyOutputsAgainstRedemptions(ctx, msg, msgTX); err != nil {
		return fmt.Errorf("failed to verify outputs against redemptions: %w", err)
	}

	// Update the redemptions to awaiting signature
	if err := k.updateRedemptions(ctx, msg.RedemptionIndexes); err != nil {
		return fmt.Errorf("failed to update redemptions to awaiting signature: %w", err)
	}

	return nil
}

func (k msgServer) updateRedemptions(ctx sdk.Context, redemptionIndices []uint64) error {
	if len(redemptionIndices) == 0 {
		return fmt.Errorf("no redemption indices provided")
	}

	// Iterate over the redemption indices, starting from the second element
	for _, redemptionIndex := range redemptionIndices[1:] {
		// Retrieve the redemption entry
		redemption, err := k.Keeper.Redemptions.Get(ctx, redemptionIndex)
		if err != nil {
			return fmt.Errorf("failed to retrieve redemption at index %d: %w", redemptionIndex, err)
		}

		// Mark the redemption as awaiting signature
		redemption.Status = types.RedemptionStatus_AWAITING_SIGN

		// Save the updated redemption entry
		if err := k.Keeper.Redemptions.Set(ctx, redemptionIndex, redemption); err != nil {
			return fmt.Errorf("failed to update redemption at index %d: %w", redemptionIndex, err)
		}
	}

	return nil
}

func (k msgServer) checkChangeAddress(ctx context.Context, msg *types.MsgSubmitUnsignedRedemptionTx) (*wire.MsgTx, error) {
	zenBTCChangeAddressKeyIDs := k.GetChangeAddressKeyIDs(ctx)
	if zenBTCChangeAddressKeyIDs == nil || len(zenBTCChangeAddressKeyIDs) == 0 {
		return nil, fmt.Errorf("failed to retrieve ZenBTCChangeAddressKeyIDs")
	}

	// Decode the transactions
	msgTX, err := bitcoin.DecodeTX(msg.Txbytes)
	if err != nil {
		return nil, fmt.Errorf("error decoding txbytes: %w", err)
	}

	// Check that the first output is the change
	if len(msgTX.TxOut) == 0 {
		return nil, fmt.Errorf("BTC transaction has zero outputs")
	}
	chaincfg := utils.ChainFromString(msg.ChainName)
	changeOutput := msgTX.TxOut[0]

	// Extract and validate change address from Output 0
	_, addrs, _, err := txscript.ExtractPkScriptAddrs(changeOutput.PkScript, chaincfg)
	if err != nil || addrs == nil || len(addrs) != 1 {
		return nil, fmt.Errorf("BTC change address invalid")
	}
	changeAddress := addrs[0].String()

	// Validate the change address against known change addresses
	validChangeAddress := false
	for _, keyID := range zenBTCChangeAddressKeyIDs {
		key, err := k.Keeper.treasuryKeeper.KeyStore.Get(ctx, keyID)
		if err != nil {
			return nil, fmt.Errorf("error retrieving change addresses for keyID %d: %w", keyID, err)
		}
		address, err := treasurytypes.BitcoinP2WPKH(&key, chaincfg)
		if err != nil {
			return nil, fmt.Errorf("error generating change address from keyID %d: %w", keyID, err)
		}
		if address == changeAddress {
			validChangeAddress = true
			break
		}
	}
	if !validChangeAddress {
		return nil, fmt.Errorf("invalid change address: %s", changeAddress)
	}
	return msgTX, nil
}

func (k msgServer) checkRedemptionTXCreator(ctx context.Context, msg *types.MsgSubmitUnsignedRedemptionTx) error {
	bitcoinProxyAddress := k.Keeper.GetBitcoinProxyAddress(ctx)
	if bitcoinProxyAddress == "" {
		return fmt.Errorf("invalid BitcoinProxyAddress")
	}

	if msg.Creator != bitcoinProxyAddress {
		return fmt.Errorf("msg.Creator (%s) must be the BitcoinProxyAddress (%s)", msg.Creator, bitcoinProxyAddress)
	}
	return nil
}

func (k msgServer) verifyOutputsAgainstRedemptions(ctx context.Context, msg *types.MsgSubmitUnsignedRedemptionTx, msgTX *wire.MsgTx) error {
	chaincfg := utils.ChainFromString(msg.ChainName)
	req := &types.QueryRedemptionsRequest{
		StartIndex: 0,
		Status:     types.RedemptionStatus_UNSTAKED,
	}

	redemptions, err := k.Keeper.GetRedemptions(ctx, req)
	if err != nil {
		return fmt.Errorf("error retrieving redemptions: %w", err)
	}
	if redemptions == nil {
		return fmt.Errorf("redemptions is nil")
	}

	// Verify that the length of RedemptionIndexes matches the number of TxOut
	if len(msg.RedemptionIndexes) != len(msgTX.TxOut) {
		return fmt.Errorf("redemptionIndexes length (%d) does not match number of outputs (%d)", len(msg.RedemptionIndexes), len(msgTX.TxOut))
	}

	// Loop through all Tx Outputs (ignoring the change output at index 0)
	for i, output := range msgTX.TxOut[1:] {
		outputIndex := uint64(i + 1) // Adjust index to correspond to TxOut slice

		// Derive output address
		_, addrs, _, err := txscript.ExtractPkScriptAddrs(output.PkScript, chaincfg)
		if err != nil || addrs == nil || len(addrs) != 1 {
			return fmt.Errorf("invalid output address at output index %d", outputIndex)
		}

		// Verify the output against redemptions
		if err = k.verifyOutputInRedemptions(redemptions, addrs[0].String(), output.Value, msg.RedemptionIndexes[outputIndex]); err != nil {
			return fmt.Errorf("invalid output %d in tx: %s, error: %w", outputIndex, hex.EncodeToString(msg.Txbytes), err)
		}
	}
	return nil
}

func (k msgServer) verifyInputsInRedemption(ctx context.Context, inputs []*types.InputHashes) error {
	changeAddressKeyIDs := k.Keeper.GetChangeAddressKeyIDs(ctx)

	for _, input := range inputs {
		// skip validation for change inputs.
		if ok := slices.Contains(changeAddressKeyIDs, input.Keyid); ok {
			continue
		}

		inputKey, err := k.treasuryKeeper.KeyStore.Get(ctx, input.Keyid)
		if err != nil {
			return err
		}

		// verify that input key is a zenBTC key
		if inputKey.ZenbtcMetadata == nil {
			return fmt.Errorf("input key is missing zenbtc_metadata: %d", input.Keyid)
		}

		// ensure that key does not contain empty "ZenbtcMetadata" structure
		if inputKey.ZenbtcMetadata.RecipientAddr == "" {
			return fmt.Errorf("input key is missing zenbtc_metadata recipient: %d", input.Keyid)
		}
	}

	return nil
}

func (k msgServer) verifyOutputInRedemptions(response *types.QueryRedemptionsResponse, outputAddress string, txAmount int64, redemptionIndex uint64) error {
	var found int

	for _, redemption := range response.Redemptions {
		// Check for matching redemption index
		if redemption.Data.Id != redemptionIndex {
			continue
		}

		// Found the expected Redemption Index
		found++

		// Validate the address
		if string(redemption.Data.DestinationAddress) != outputAddress {
			return fmt.Errorf("destination address is invalid: expected %s, got %s", outputAddress, redemption.Data.DestinationAddress)
		}

		// Validate the amount
		if int64(redemption.Data.Amount) != txAmount {
			return fmt.Errorf("output has invalid amount: expected %d, got %d", txAmount, redemption.Data.Amount)
		}
	}

	// Validate the number of found redemptions
	switch found {
	case 0:
		return fmt.Errorf("redemption not found")
	case 1:
		// Happy path: Redemption index, address, and amount are valid
		return nil
	default:
		return fmt.Errorf("redemption found more than once")
	}
}

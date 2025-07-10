//go:build test
// +build test

package keeper

import (
	"context"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Test helper methods to expose private methods for testing

// RecordMismatchedVoteExtensions is a test helper that wraps the private recordMismatchedVoteExtensions method
func (k *Keeper) RecordMismatchedVoteExtensions(ctx sdk.Context, height int64, canonicalVoteExt VoteExtension, consensusData abci.ExtendedCommitInfo) {
	k.recordMismatchedVoteExtensions(ctx, height, canonicalVoteExt, consensusData)
}

// UpdateValidatorMismatchCount is a test helper that wraps the private updateValidatorMismatchCount method
func (k *Keeper) UpdateValidatorMismatchCount(ctx sdk.Context, validatorHexAddr string, blockHeight int64) {
	k.updateValidatorMismatchCount(ctx, validatorHexAddr, blockHeight)
}

// ValidateOracleData is a test helper that wraps the private validateOracleData method
func (k *Keeper) ValidateOracleData(ctx context.Context, voteExt VoteExtension, oracleData *OracleData, fieldVotePowers map[VoteExtensionField]int64) {
	k.validateOracleData(ctx, voteExt, oracleData, fieldVotePowers)
}

// HandleValidationMismatches is a test helper that wraps the private handleValidationMismatches method
func (k *Keeper) HandleValidationMismatches(ctx context.Context, mismatches []validationMismatch, fieldVotePowers map[VoteExtensionField]int64, oracleData *OracleData) {
	k.handleValidationMismatches(ctx, mismatches, fieldVotePowers, oracleData)
}

// CheckAndJailValidatorsForMismatchedVoteExtensions is a test helper that wraps the private checkAndJailValidatorsForMismatchedVoteExtensions method
func (k *Keeper) CheckAndJailValidatorsForMismatchedVoteExtensions(ctx context.Context) error {
	return k.checkAndJailValidatorsForMismatchedVoteExtensions(ctx)
}

// ValidateVote is a test helper that wraps the private validateVote method
func (k *Keeper) ValidateVote(ctx context.Context, vote abci.ExtendedVoteInfo, currentHeight int64) (VoteExtension, error) {
	return k.validateVote(ctx, vote, currentHeight)
}

package v3

import (

	// "github.com/Zenrock-Foundation/zrchain/v6/x/mint/exported"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "mint"
)

var ParamsKey = []byte{0x01}

// UpdateParams migrates the x/mint module state from the consensus version 2 to
// version 3. Specifically, it adds several new parameters to the mint module
// and removes the legacy minter logic and replaces it with a deflationary
// model.
func UpdateParams(
	ctx sdk.Context,
	params collections.Item[types.Params],
) error {
	oldParams, err := params.Get(ctx)
	if err != nil {
		return err
	}

	currParams := oldParams
	currParams.StakingYield = math.LegacyNewDecWithPrec(7, 2)
	currParams.BurnRate = math.LegacyNewDecWithPrec(10, 2)
	currParams.ProtocolWalletRate = math.LegacyNewDecWithPrec(30, 2)
	currParams.ProtocolWalletAddress = "zen1fhln2vnudxddpymqy82vzqhnlsfh4stjd683ze"
	currParams.AdditionalStakingRewards = math.LegacyNewDecWithPrec(30, 2)
	currParams.AdditionalMpcRewards = math.LegacyNewDecWithPrec(5, 2)
	currParams.AdditionalBurnRate = math.LegacyNewDecWithPrec(25, 2)

	if err := currParams.Validate(); err != nil {
		return err
	}

	return params.Set(ctx, currParams)
}

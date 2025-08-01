package types

import (
	"encoding/json"

	"cosmossdk.io/math"
	sidecarapi "github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	zenbtc "github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

// NewGenesisState creates a new GenesisState instanc e
func NewGenesisState(
	params Params,
	validators []ValidatorHV,
	delegations []Delegation,
	lastTotalPower math.Int,
	lastValidatorPower []LastValidatorPower,
	btcBlockHeaders map[int64]sidecarapi.BTCBlockHeader,
	solanaNonce map[uint64]SolanaNonce,
	backfillRequest BackfillRequests,
	ethereumNonce map[uint64]zenbtc.NonceData,
	requestedHistoricalBitcoinHeaders zenbtc.RequestedBitcoinHeaders,
	avsRewardsPool map[string]math.Int,
	ethereumNonceRequested map[uint64]bool,
	solanaNonceRequested map[uint64]bool,
	solanaZentpAccountsRequested map[string]bool,
	solanaAccountsRequested map[string]bool,
	exported bool,
	unbondingDelegations []UnbondingDelegation,
	redelegations []Redelegation,
	hvParams HVParams,
	assetPrices map[int32]math.LegacyDec,
	lastValidVeHeight int64,
	slashEvents map[uint64]SlashEvent,
	slashEventCount uint64,
	validationInfos map[int64]ValidationInfo,
	lastUsedSolanaNonce map[uint64]SolanaNonce,
	lastUsedEthereumNonce map[uint64]zenbtc.NonceData,
) *GenesisState {
	return &GenesisState{
		Params:                            params,
		LastTotalPower:                    lastTotalPower,
		LastValidatorPowers:               lastValidatorPower,
		Validators:                        validators,
		Delegations:                       delegations,
		UnbondingDelegations:              unbondingDelegations,
		Redelegations:                     redelegations,
		Exported:                          exported,
		HVParams:                          hvParams,
		AssetPrices:                       assetPrices,
		LastValidVeHeight:                 lastValidVeHeight,
		SlashEvents:                       slashEvents,
		SlashEventCount:                   slashEventCount,
		ValidationInfos:                   validationInfos,
		BtcBlockHeaders:                   btcBlockHeaders,
		LastUsedSolanaNonce:               solanaNonce,
		BackfillRequest:                   backfillRequest,
		LastUsedEthereumNonce:             ethereumNonce,
		RequestedHistoricalBitcoinHeaders: requestedHistoricalBitcoinHeaders,
		AvsRewardsPool:                    avsRewardsPool,
		EthereumNonceRequested:            ethereumNonceRequested,
		SolanaNonceRequested:              solanaNonceRequested,
		SolanaZentpAccountsRequested:      solanaZentpAccountsRequested,
		SolanaAccountsRequested:           solanaAccountsRequested,
	}
}

// DefaultGenesisState gets the raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// GetGenesisStateFromAppState returns x/staking GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.JSONCodec, appState map[string]json.RawMessage) *GenesisState {
	var genesisState GenesisState

	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return &genesisState
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (g GenesisState) UnpackInterfaces(c codectypes.AnyUnpacker) error {
	for i := range g.Validators {
		if err := g.Validators[i].UnpackInterfaces(c); err != nil {
			return err
		}
	}
	return nil
}

package testutil

import (
	"time"

	"cosmossdk.io/math"
	validationtypes "github.com/Zenrock-Foundation/zrchain/v6/x/validation/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	zenbtc "github.com/zenrocklabs/zenbtc/x/zenbtc/types"
)

func DefaultGenesis() *validationtypes.GenesisState {
	return &validationtypes.GenesisState{
		Params: validationtypes.DefaultParams(),
	}
}

// TestGenesis creates a comprehensive test genesis state
func TestGenesis() *validationtypes.GenesisState {
	// Create validator public key with exactly 32 bytes
	pubKey := ed25519.PubKey{Key: []byte{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
		17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32,
	}}
	anyPubKey, _ := codectypes.NewAnyWithValue(&pubKey)

	// Create validator with actual data
	validator := validationtypes.ValidatorHV{
		OperatorAddress: "zenvaloper126hek6zagmp3jqf97x7pq7c0j9jqs0ndvcepy6",
		ConsensusPubkey: anyPubKey,
		Jailed:          false,
		Status:          validationtypes.Bonded,
		TokensNative:    math.NewInt(125000000000000),
		TokensAVS:       math.ZeroInt(),
		DelegatorShares: math.LegacyNewDecFromInt(math.NewInt(125000000000000)),
		Description: validationtypes.Description{
			Moniker:         "zenrock",
			Identity:        "",
			Website:         "",
			SecurityContact: "",
			Details:         "",
		},
		UnbondingHeight: 0,
		UnbondingTime:   time.Unix(0, 0).UTC(),
		Commission: validationtypes.Commission{
			CommissionRates: validationtypes.CommissionRates{
				Rate:          math.LegacyMustNewDecFromStr("0.100000000000000000"),
				MaxRate:       math.LegacyMustNewDecFromStr("0.200000000000000000"),
				MaxChangeRate: math.LegacyMustNewDecFromStr("0.010000000000000000"),
			},
			UpdateTime: time.Date(2025, 6, 26, 8, 15, 28, 251780000, time.UTC),
		},
		MinSelfDelegation:       math.OneInt(),
		UnbondingOnHoldRefCount: 0,
		UnbondingIds:            nil,
	}

	// Create delegation with actual data
	delegation := validationtypes.Delegation{
		DelegatorAddress: "zen126hek6zagmp3jqf97x7pq7c0j9jqs0ndxeaqhq",
		ValidatorAddress: "zenvaloper126hek6zagmp3jqf97x7pq7c0j9jqs0ndvcepy6",
		Shares:           math.LegacyNewDecFromInt(math.NewInt(125000000000000)),
	}

	// Create last validator power with actual data
	lastValidatorPower := validationtypes.LastValidatorPower{
		Address: "zenvaloper126hek6zagmp3jqf97x7pq7c0j9jqs0ndvcepy6",
		Power:   125000000,
	}

	// Create HV params with actual data
	hvParams := validationtypes.HVParams{
		AVSRewardsRate: math.LegacyMustNewDecFromStr("0.030000000000000000"),
		BlockTime:      5,
		StakeableAssets: []*validationtypes.AssetData{
			{
				Asset:     validationtypes.Asset_ROCK,
				PriceUSD:  math.LegacyZeroDec(),
				Precision: 6,
			},
		},
		PriceRetentionBlockRange: 100,
	}

	// Create asset prices with actual data
	assetPrices := []*validationtypes.AssetData{
		{
			Asset:     validationtypes.Asset_ROCK,
			PriceUSD:  math.LegacyZeroDec(),
			Precision: 0,
		},
		{
			Asset:     validationtypes.Asset_BTC,
			PriceUSD:  math.LegacyZeroDec(),
			Precision: 0,
		},
		{
			Asset:     validationtypes.Asset_ETH,
			PriceUSD:  math.LegacyZeroDec(),
			Precision: 0,
		},
	}

	// Create validation infos with actual data
	validationInfos := make([]validationtypes.ValidationInfo, 3)
	for i := 0; i < 3; i++ {
		validationInfos[i] = validationtypes.ValidationInfo{
			NonVotingValidators: nil,
			MismatchedVoteExtensions: []string{
				"c4848a0c008c40400d5fe4f0d546fa61f97f7d05",
			},
			BlockHeight: uint64(i),
		}
	}

	// Create backfill requests with actual data
	backfillRequest := validationtypes.BackfillRequests{
		Requests: nil,
	}

	// Create requested historical bitcoin headers with actual data
	requestedHistoricalBitcoinHeaders := []zenbtc.RequestedBitcoinHeaders{
		{
			Heights: nil,
		},
	}

	// Create AVS rewards pool with actual data
	avsRewardsPool := []string{
		"zenvaloper126hek6zagmp3jqf97x7pq7c0j9jqs0ndvcepy6",
	}

	// Create Solana nonce requested with actual data
	solanaNonceRequested := []uint64{12, 123}

	return &validationtypes.GenesisState{
		Params:         validationtypes.DefaultParams(),
		LastTotalPower: math.NewInt(125000000),
		LastValidatorPowers: []validationtypes.LastValidatorPower{
			lastValidatorPower,
		},
		Validators: []validationtypes.ValidatorHV{
			validator,
		},
		Delegations: []validationtypes.Delegation{
			delegation,
		},
		UnbondingDelegations:              nil,
		Redelegations:                     nil,
		Exported:                          true,
		HvParams:                          &hvParams,
		AssetPrices:                       assetPrices,
		LastValidVeHeight:                 0,
		SlashEvents:                       nil,
		SlashEventCount:                   0,
		ValidationInfos:                   validationInfos,
		BtcBlockHeaders:                   nil,
		LastUsedSolanaNonce:               nil,
		BackfillRequest:                   backfillRequest,
		LastUsedEthereumNonce:             nil,
		RequestedHistoricalBitcoinHeaders: requestedHistoricalBitcoinHeaders,
		AvsRewardsPool:                    avsRewardsPool,
		EthereumNonceRequested:            []uint64{},
		SolanaNonceRequested:              solanaNonceRequested,
		SolanaZentpAccountsRequested:      []string{},
		SolanaAccountsRequested:           []string{},
	}
}

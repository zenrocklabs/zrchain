package keeper

import (
	"testing"

	"cosmossdk.io/log"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/stretchr/testify/require"
)

var defaultVe = VoteExtension{
	ZRChainBlockHeight:         100,
	EigenDelegationsHash:       []byte("randomhash"),
	RequestedBtcBlockHeight:    100000,
	RequestedBtcHeaderHash:     []byte("randomhash"),
	EthBlockHeight:             100000,
	EthGasLimit:                2000000000,
	EthBaseFee:                 1000000000,
	EthTipCap:                  1000000000,
	RequestedStakerNonce:       1,
	RequestedEthMinterNonce:    1,
	RequestedUnstakerNonce:     1,
	RequestedCompleterNonce:    1,
	SolanaMintNonceHashes:      []byte("randomhash"),
	SolanaAccountsHash:         []byte("randomhash"),
	SolanaLamportsPerSignature: 1000000000,
	EthBurnEventsHash:          []byte("randomhash"),
	SolanaBurnEventsHash:       []byte("randomhash"),
	SolanaMintEventsHash:       []byte("randomhash"),
	RedemptionsHash:            []byte("randomhash"),
	ROCKUSDPrice:               "1000000000",
	BTCUSDPrice:                "1000000000",
	ETHUSDPrice:                "1000000000",
	LatestBtcBlockHeight:       100000,
	LatestBtcHeaderHash:        []byte("randomhash"),
}

func TestContainsVoteExtension(t *testing.T) {

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)
	txConfig := tx.NewTxConfig(cdc, tx.DefaultSignModes)

	voteExtBytes := []byte("vote extension data")

	tests := []struct {
		name      string
		txBytes   []byte
		expectExt bool
	}{
		{
			name:      "regular transaction",
			txBytes:   []byte("valid tx bytes"),
			expectExt: true,
		},
		{
			name:      "vote extension",
			txBytes:   voteExtBytes,
			expectExt: true,
		},
		{
			name:      "empty transaction",
			txBytes:   []byte{},
			expectExt: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasExt := ContainsVoteExtension(tt.txBytes, txConfig.TxDecoder())
			require.Equal(t, tt.expectExt, hasExt)
		})
	}
}

func TestVoteExtensionsEnabled(t *testing.T) {
	tests := []struct {
		name            string
		blockHeight     int64
		consensusParams cmtproto.ConsensusParams
		enabled         bool
	}{
		{
			name:        "vote extensions enabled",
			blockHeight: 100,
			consensusParams: cmtproto.ConsensusParams{
				Abci: &cmtproto.ABCIParams{
					VoteExtensionsEnableHeight: 50,
				},
			},
			enabled: true,
		},
		{
			name:        "vote extensions disabled",
			blockHeight: 100,
			consensusParams: cmtproto.ConsensusParams{
				Abci: &cmtproto.ABCIParams{
					VoteExtensionsEnableHeight: 101,
				},
			},
			enabled: false,
		},
		{
			name:        "block height is 1",
			blockHeight: 1,
			consensusParams: cmtproto.ConsensusParams{
				Abci: &cmtproto.ABCIParams{
					VoteExtensionsEnableHeight: 100,
				},
			},
			enabled: false,
		},
		{
			name:        "abci is nil",
			blockHeight: 100,
			consensusParams: cmtproto.ConsensusParams{
				Abci: nil,
			},
			enabled: false,
		},
		{
			name:        "vote extensions enable height is 0",
			blockHeight: 100,
			consensusParams: cmtproto.ConsensusParams{
				Abci: &cmtproto.ABCIParams{
					VoteExtensionsEnableHeight: 0,
				},
			},
			enabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := sdk.Context{}
			ctx = ctx.WithBlockHeight(tt.blockHeight)
			ctx = ctx.WithConsensusParams(tt.consensusParams)
			enabled := VoteExtensionsEnabled(ctx)
			require.Equal(t, tt.enabled, enabled)
		})
	}
}

func TestIsInvalid(t *testing.T) {

	type testCase struct {
		name          string
		fieldToModify string
		invalidValue  any
		expectInvalid bool
	}

	tests := []testCase{
		{
			name:          "valid vote extension",
			fieldToModify: "",
			expectInvalid: false,
		},
		{
			name:          "invalid zr chain block height",
			fieldToModify: "ZRChainBlockHeight",
			invalidValue:  int64(0),
			expectInvalid: true,
		},
		{
			name:          "invalid eigen delegations hash",
			fieldToModify: "EigenDelegationsHash",
			invalidValue:  []byte{},
			expectInvalid: true,
		},
		{
			name:          "invalid eth block height",
			fieldToModify: "EthBlockHeight",
			invalidValue:  uint64(0),
			expectInvalid: true,
		},
		{
			name:          "invalid eth base fee",
			fieldToModify: "EthBaseFee",
			invalidValue:  uint64(0),
			expectInvalid: true,
		},
		{
			name:          "invalid eth tip cap",
			fieldToModify: "EthTipCap",
			invalidValue:  uint64(0),
			expectInvalid: true,
		},
		{
			name:          "invalid eth gas limit",
			fieldToModify: "EthGasLimit",
			invalidValue:  uint64(0),
			expectInvalid: true,
		},
		{
			name:          "invalid eth burn events hash",
			fieldToModify: "EthBurnEventsHash",
			invalidValue:  []byte{},
			expectInvalid: true,
		},
		{
			name:          "invalid redemptions hash",
			fieldToModify: "RedemptionsHash",
			invalidValue:  []byte{},
			expectInvalid: true,
		},
		{
			name:          "invalid rock usd price",
			fieldToModify: "ROCKUSDPrice",
			invalidValue:  "",
			expectInvalid: true,
		},
		{
			name:          "invalid btc usd price",
			fieldToModify: "BTCUSDPrice",
			invalidValue:  "",
			expectInvalid: true,
		},
		{
			name:          "invalid eth usd price",
			fieldToModify: "ETHUSDPrice",
			invalidValue:  "",
			expectInvalid: true,
		},
		{
			name:          "invalid latest btc block height",
			fieldToModify: "LatestBtcBlockHeight",
			invalidValue:  int64(0),
			expectInvalid: true,
		},
		{
			name:          "invalid latest btc header hash",
			fieldToModify: "LatestBtcHeaderHash",
			invalidValue:  []byte{},
			expectInvalid: true,
		},
	}

	// Map of field names to their setters
	fieldSetters := map[string]func(*VoteExtension, any){
		"ZRChainBlockHeight":   func(ve *VoteExtension, v any) { ve.ZRChainBlockHeight = v.(int64) },
		"EigenDelegationsHash": func(ve *VoteExtension, v any) { ve.EigenDelegationsHash = v.([]byte) },
		"EthBlockHeight":       func(ve *VoteExtension, v any) { ve.EthBlockHeight = v.(uint64) },
		"EthBaseFee":           func(ve *VoteExtension, v any) { ve.EthBaseFee = v.(uint64) },
		"EthTipCap":            func(ve *VoteExtension, v any) { ve.EthTipCap = v.(uint64) },
		"EthGasLimit":          func(ve *VoteExtension, v any) { ve.EthGasLimit = v.(uint64) },
		"EthBurnEventsHash":    func(ve *VoteExtension, v any) { ve.EthBurnEventsHash = v.([]byte) },
		"RedemptionsHash":      func(ve *VoteExtension, v any) { ve.RedemptionsHash = v.([]byte) },
		"ROCKUSDPrice":         func(ve *VoteExtension, v any) { ve.ROCKUSDPrice = v.(string) },
		"BTCUSDPrice":          func(ve *VoteExtension, v any) { ve.BTCUSDPrice = v.(string) },
		"ETHUSDPrice":          func(ve *VoteExtension, v any) { ve.ETHUSDPrice = v.(string) },
		"LatestBtcBlockHeight": func(ve *VoteExtension, v any) { ve.LatestBtcBlockHeight = v.(int64) },
		"LatestBtcHeaderHash":  func(ve *VoteExtension, v any) { ve.LatestBtcHeaderHash = v.([]byte) },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.NewNopLogger()
			ve := defaultVe

			if tt.fieldToModify != "" {
				setter, ok := fieldSetters[tt.fieldToModify]
				if !ok {
					t.Fatalf("no setter found for field %s", tt.fieldToModify)
				}
				setter(&ve, tt.invalidValue)
			}

			valid := ve.IsInvalid(logger)
			require.Equal(t, tt.expectInvalid, valid)
		})
	}
}

func TestHasRequiredGasFields(t *testing.T) {
	tests := []struct {
		name            string
		fieldVotePowers map[VoteExtensionField]int64
		valid           bool
	}{
		{
			name: "valid vote extension",
			fieldVotePowers: map[VoteExtensionField]int64{
				VEFieldEthGasLimit: 100,
				VEFieldEthBaseFee:  100,
				VEFieldEthTipCap:   100,
			},
			valid: true,
		},
		{
			name: "missing ve field eth tip cap",
			fieldVotePowers: map[VoteExtensionField]int64{
				VEFieldEthGasLimit: 100,
				VEFieldEthBaseFee:  100,
			},
			valid: false,
		},
		{
			name: "missing ve field eth base fee",
			fieldVotePowers: map[VoteExtensionField]int64{
				VEFieldEthGasLimit: 100,
				VEFieldEthTipCap:   100,
			},
			valid: false,
		},
		{
			name: "missing ve field eth gas limit",
			fieldVotePowers: map[VoteExtensionField]int64{
				VEFieldEthBaseFee: 100,
				VEFieldEthTipCap:  100,
			},
			valid: false,
		},
		{
			name: "invalid ve field eth gas limit",
			fieldVotePowers: map[VoteExtensionField]int64{
				VEFieldEthGasLimit: 0,
				VEFieldEthBaseFee:  100,
				VEFieldEthTipCap:   100,
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := HasRequiredGasFields(tt.fieldVotePowers)
			require.Equal(t, tt.valid, valid)
		})
	}
}

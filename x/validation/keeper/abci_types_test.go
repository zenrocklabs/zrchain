package keeper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"math"
	"math/big"
	"strings"
	"testing"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	sidecar "github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/shared"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/stretchr/testify/require"
)

var defaultVe = VoteExtension{
	EigenDelegationsHash:    []byte("randomhash"),
	RequestedBtcBlockHeight: 100000,
	RequestedBtcHeaderHash:  []byte("randomhash"),
	EthBlockHeight:          100000,
	EthGasLimit:             2000000000,
	EthBaseFee:              1000000000,
	EthTipCap:               1000000000,
	RequestedStakerNonce:    1,
	RequestedEthMinterNonce: 1,
	RequestedUnstakerNonce:  1,
	RequestedCompleterNonce: 1,
	SolanaMintNoncesHash:    []byte("randomhash"),
	SolanaAccountsHash:      []byte("randomhash"),
	EthBurnEventsHash:       []byte("randomhash"),
	SolanaBurnEventsHash:    []byte("randomhash"),
	SolanaMintEventsHash:    []byte("randomhash"),
	RedemptionsHash:         []byte("randomhash"),
	ROCKUSDPrice:            "1000000000",
	BTCUSDPrice:             "1000000000",
	ETHUSDPrice:             "1000000000",
	ZECUSDPrice:             "1000000000",
	LatestBtcBlockHeight:    100000,
	LatestBtcHeaderHash:     []byte("randomhash"),
	SidecarVersionName:      sidecartypes.SidecarVersionName,
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
		// {
		// 	name:          "invalid eigen delegations hash",
		// 	fieldToModify: "EigenDelegationsHash",
		// 	invalidValue:  []byte{},
		// 	expectInvalid: true,
		// },
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
		{
			name:          "invalid sidecar version name",
			fieldToModify: "SidecarVersionName",
			invalidValue:  "",
			expectInvalid: true,
		},
	}

	// Map of field names to their setters
	fieldSetters := map[string]func(*VoteExtension, any){
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
		"SidecarVersionName":   func(ve *VoteExtension, v any) { ve.SidecarVersionName = v.(string) },
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

func TestIsGasField(t *testing.T) {
	tests := []struct {
		name  string
		field VoteExtensionField
		valid bool
	}{
		{
			name:  "valid gas field: eth gas limit",
			field: VEFieldEthGasLimit,
			valid: true,
		},
		{
			name:  "valid gas field: eth base fee",
			field: VEFieldEthBaseFee,
			valid: true,
		},
		{
			name:  "valid gas field: eth tip cap",
			field: VEFieldEthTipCap,
			valid: true,
		},
		{
			name:  "invalid gas field: requested btc block height",
			field: VEFieldRequestedBtcBlockHeight,
			valid: false,
		},
		{
			name:  "invalid gas field: nil",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := isGasField(tt.field)
			require.Equal(t, tt.valid, valid)
		})
	}
}

func TestFieldHasConsensus(t *testing.T) {
	tests := []struct {
		name            string
		fieldVotePowers map[VoteExtensionField]int64
		field           VoteExtensionField
		hasConsensus    bool
	}{
		{
			name: "valid field: eth gas limit",
			fieldVotePowers: map[VoteExtensionField]int64{
				VEFieldEthGasLimit: 100,
			},
			field:        VEFieldEthGasLimit,
			hasConsensus: true,
		},
		{
			name: "valid field: eth base fee",
			fieldVotePowers: map[VoteExtensionField]int64{
				VEFieldEthBaseFee: 100,
			},
			field:        VEFieldEthBaseFee,
			hasConsensus: true,
		},
		{
			name: "unvalid field: nil",
			fieldVotePowers: map[VoteExtensionField]int64{
				VEFieldEthGasLimit: 100,
			},
			hasConsensus: false,
		},
		{
			name:            "unvalid fieldvotepowers: nil",
			fieldVotePowers: nil,
			field:           VEFieldEthGasLimit,
			hasConsensus:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := fieldHasConsensus(tt.fieldVotePowers, tt.field)
			require.Equal(t, tt.hasConsensus, valid)
		})
	}
}

func TestAllFieldsHaveConsensus(t *testing.T) {
	tests := []struct {
		name            string
		fieldVotePowers map[VoteExtensionField]int64
		fields          []VoteExtensionField
		hasConsensus    []VoteExtensionField
	}{
		{
			name: "valid fields",
			fieldVotePowers: map[VoteExtensionField]int64{
				VEFieldEthGasLimit: 100,
				VEFieldEthBaseFee:  100,
				VEFieldEthTipCap:   100,
			},
			fields:       []VoteExtensionField{},
			hasConsensus: nil,
		},
		{
			name:            "empty field vote powers",
			fieldVotePowers: map[VoteExtensionField]int64{},
			fields:          []VoteExtensionField{VEFieldEthGasLimit, VEFieldEthBaseFee, VEFieldEthTipCap},
			hasConsensus:    []VoteExtensionField{VEFieldEthGasLimit, VEFieldEthBaseFee, VEFieldEthTipCap},
		},
		{
			name: "missing some fields",
			fieldVotePowers: map[VoteExtensionField]int64{
				VEFieldEthGasLimit: 100,
				VEFieldEthBaseFee:  100,
			},
			fields:       []VoteExtensionField{VEFieldEthGasLimit, VEFieldEthBaseFee, VEFieldEthTipCap},
			hasConsensus: []VoteExtensionField{VEFieldEthTipCap},
		},
		{
			name: "all fields have consensus",
			fieldVotePowers: map[VoteExtensionField]int64{
				VEFieldEthGasLimit: 100,
				VEFieldEthBaseFee:  100,
				VEFieldEthTipCap:   100,
			},
			fields:       []VoteExtensionField{VEFieldEthGasLimit, VEFieldEthBaseFee, VEFieldEthTipCap},
			hasConsensus: nil,
		},
		{
			name: "zero vote power fields",
			fieldVotePowers: map[VoteExtensionField]int64{
				VEFieldEthGasLimit: 0,
				VEFieldEthBaseFee:  0,
				VEFieldEthTipCap:   0,
			},
			fields:       []VoteExtensionField{VEFieldEthGasLimit, VEFieldEthBaseFee, VEFieldEthTipCap},
			hasConsensus: nil,
		},
		{
			name: "mixed vote powers",
			fieldVotePowers: map[VoteExtensionField]int64{
				VEFieldEthGasLimit: 100,
				VEFieldEthBaseFee:  0,
				VEFieldEthTipCap:   50,
			},
			fields:       []VoteExtensionField{VEFieldEthGasLimit, VEFieldEthBaseFee, VEFieldEthTipCap},
			hasConsensus: nil,
		},
		{
			name: "nil fields slice",
			fieldVotePowers: map[VoteExtensionField]int64{
				VEFieldEthGasLimit: 100,
				VEFieldEthBaseFee:  100,
				VEFieldEthTipCap:   100,
			},
			fields:       nil,
			hasConsensus: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := allFieldsHaveConsensus(tt.fieldVotePowers, tt.fields)
			require.Equal(t, tt.hasConsensus, valid)
		})
	}
}

func TestAnyFieldHasConsensus(t *testing.T) {
	tests := []struct {
		name            string
		fieldVotePowers map[VoteExtensionField]int64
		fields          []VoteExtensionField
		hasConsensus    bool
	}{
		{
			name: "valid fields",
			fieldVotePowers: map[VoteExtensionField]int64{
				VEFieldEthGasLimit: 100,
			},
			fields: []VoteExtensionField{
				VEFieldEthGasLimit,
				VEFieldEthBaseFee,
				VEFieldEthTipCap,
			},
			hasConsensus: true,
		},
		{
			name: "empty fields slice",
			fieldVotePowers: map[VoteExtensionField]int64{
				VEFieldEthGasLimit: 100,
			},
		},
		{
			name: "nil field vote powers",
			fields: []VoteExtensionField{
				VEFieldEthGasLimit,
			},
			hasConsensus: false,
		},
		{
			name: "zero vote power",
			fieldVotePowers: map[VoteExtensionField]int64{
				VEFieldEthGasLimit: 0,
			},
			fields: []VoteExtensionField{
				VEFieldEthGasLimit,
			},
			hasConsensus: true,
		},
		{
			name: "negative vote power",
			fieldVotePowers: map[VoteExtensionField]int64{
				VEFieldEthGasLimit: -100,
			},
			fields: []VoteExtensionField{
				VEFieldEthGasLimit,
			},
			hasConsensus: true,
		},
		{
			name: "max int64 vote power",
			fieldVotePowers: map[VoteExtensionField]int64{
				VEFieldEthGasLimit: math.MaxInt64,
			},
			fields: []VoteExtensionField{
				VEFieldEthGasLimit,
			},
			hasConsensus: true,
		},
		{
			name: "min int64 vote power",
			fieldVotePowers: map[VoteExtensionField]int64{
				VEFieldEthGasLimit: math.MinInt64,
			},
			fields: []VoteExtensionField{
				VEFieldEthGasLimit,
			},
			hasConsensus: true,
		},
		{
			name: "duplicate fields",
			fieldVotePowers: map[VoteExtensionField]int64{
				VEFieldEthGasLimit: 100,
			},
			fields: []VoteExtensionField{
				VEFieldEthGasLimit,
				VEFieldEthGasLimit,
			},
			hasConsensus: true,
		},
		{
			name:            "empty map with fields",
			fieldVotePowers: map[VoteExtensionField]int64{},
			fields: []VoteExtensionField{
				VEFieldEthGasLimit,
			},
			hasConsensus: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := anyFieldHasConsensus(tt.fieldVotePowers, tt.fields)
			require.Equal(t, tt.hasConsensus, valid)
		})
	}
}

func TestGenericGetKey(t *testing.T) {
	tests := []struct {
		name  string
		value any
		key   string
	}{
		{
			name:  "valid value",
			value: []byte("randomhash"),
			key:   "72616e646f6d68617368",
		},
		{
			name:  "valid value: int64",
			value: int64(100),
			key:   "100",
		},
		{
			name:  "valid value: uint64",
			value: uint64(100),
			key:   "100",
		},
		{
			name:  "empty value",
			value: nil,
			key:   "",
		},
		{
			name:  "empty byte slice",
			value: []byte{},
			key:   "",
		},
		{
			name:  "max int64",
			value: int64(math.MaxInt64),
			key:   "9223372036854775807",
		},
		{
			name:  "min int64",
			value: int64(math.MinInt64),
			key:   "-9223372036854775808",
		},
		{
			name:  "max uint64",
			value: uint64(math.MaxUint64),
			key:   "18446744073709551615",
		},
		{
			name:  "zero values",
			value: int64(0),
			key:   "0",
		},
		{
			name:  "negative value",
			value: int64(-100),
			key:   "-100",
		},
		{
			name:  "special characters in byte slice",
			value: []byte{0, 1, 2, 3, 255},
			key:   "00010203ff",
		},
		{
			name:  "unicode characters in byte slice",
			value: []byte("世界"),
			key:   "e4b896e7958c",
		},
		{
			name:  "very long byte slice",
			value: bytes.Repeat([]byte{1}, 1000),
			key:   strings.Repeat("01", 1000),
		},
		{
			name:  "struct value",
			value: struct{}{},
			key:   "{}",
		},
		{
			name:  "pointer value",
			value: new(int),
			key:   "0",
		},
		{
			name:  "slice value",
			value: []int{1, 2, 3},
			key:   "[1,2,3]",
		},
		{
			name:  "map value",
			value: map[string]int{"a": 1},
			key:   "{\"a\":1}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := genericGetKey(tt.value)
			require.Equal(t, tt.key, key)
		})
	}
}

func TestVoteExtensionFieldString(t *testing.T) {
	tests := []struct {
		name  string
		field VoteExtensionField
		str   string
	}{
		{
			name:  "valid field: eigen delegations hash",
			field: VEFieldEigenDelegationsHash,
			str:   "EigenDelegationsHash",
		},
		{
			name:  "valid field: eth gas limit",
			field: VEFieldEthGasLimit,
			str:   "EthGasLimit",
		},
		{
			name:  "valid field: eth base fee",
			field: VEFieldEthBaseFee,
			str:   "EthBaseFee",
		},
		{
			name:  "valid field: eth tip cap",
			field: VEFieldEthTipCap,
			str:   "EthTipCap",
		},
		{
			name:  "invalid field",
			field: 1111,
			str:   "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := tt.field.String()
			require.Equal(t, tt.str, str)
		})
	}
}

func TestFieldHandlerSetValue(t *testing.T) {
	handlers := initializeFieldHandlers()
	ve := defaultVe

	for _, handler := range handlers {
		t.Run(handler.Field.String(), func(t *testing.T) {
			// Get the original value
			originalValue := handler.GetValue(ve)

			// Create a new VoteExtension to set the value in
			newVe := defaultVe

			// Set the value
			handler.SetValue(originalValue, &newVe)

			// Get the value back and verify it matches
			gotValue := handler.GetValue(newVe)
			require.Equal(t, originalValue, gotValue, "value mismatch for field %s", handler.Field.String())
		})
	}
}

func TestNewSidecarClient(t *testing.T) {
	tests := []struct {
		name       string
		serverAddr string
		expectErr  bool
	}{
		{
			name:       "valid server address",
			serverAddr: "localhost:50051",
			expectErr:  false,
		},
		{
			name:       "invalid server address",
			serverAddr: "localhost:50052",
			expectErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sidecarClient, err := NewSidecarClient(tt.serverAddr)
			if tt.expectErr {
				require.Error(t, err)
				require.Nil(t, sidecarClient)
			} else {
				require.NoError(t, err)
				require.NotNil(t, sidecarClient)
			}
		})
	}
}

func TestProcessOracleResponse(t *testing.T) {
	keeper := Keeper{}

	tests := []struct {
		name        string
		resp        *sidecar.SidecarStateResponse
		delegations map[string]map[string]*big.Int
		expected    *OracleData
		err         error
	}{
		{
			name: "processes oracle response with zero BTC price and delegations",
			resp: &sidecar.SidecarStateResponse{
				EthBlockHeight:   100,
				EthGasLimit:      100,
				EthBaseFee:       100,
				EthTipCap:        0,
				EigenDelegations: []byte(`{"validator1":{"delegator1":100}}`),
			},
			expected: &OracleData{
				EigenDelegationsMap: map[string]map[string]*big.Int{
					"validator1": {
						"delegator1": big.NewInt(100),
					},
				},
				ValidatorDelegations: []ValidatorDelegations{
					{
						Validator: "validator1",
						Stake:     sdkmath.NewInt(100),
					},
				},
				EthBlockHeight: 100,
				EthGasLimit:    100,
				EthBaseFee:     100,
			},
			err: nil,
		},
		{
			name: "processes oracle response with typical ETH mainnet values and delegations",
			resp: &sidecar.SidecarStateResponse{
				EthBlockHeight:   100,
				EthGasLimit:      100,
				EthBaseFee:       100,
				EthTipCap:        0,
				EigenDelegations: []byte(`{"validator1":{"delegator1":100}}`),
			},
			expected: &OracleData{
				EigenDelegationsMap: map[string]map[string]*big.Int{
					"validator1": {
						"delegator1": big.NewInt(100),
					},
				},
				ValidatorDelegations: []ValidatorDelegations{
					{
						Validator: "validator1",
						Stake:     sdkmath.NewInt(100),
					},
				},
				EthBlockHeight: 100,
				EthGasLimit:    100,
				EthBaseFee:     100,
			},
		},
		{
			name: "returns error when processing delegations fails",
			resp: &sidecar.SidecarStateResponse{
				EthBlockHeight:   100,
				EthGasLimit:      100,
				EthBaseFee:       100,
				EthTipCap:        0,
				EigenDelegations: []byte(`invalid json`),
			},
			expected: nil,
			err:      &json.SyntaxError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := keeper.processOracleResponse(context.Background(), tt.resp)
			if tt.err != nil {
				require.Error(t, err)
				var syntaxErr *json.SyntaxError
				require.True(t, errors.As(err, &syntaxErr))
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
			}
		})
	}
}

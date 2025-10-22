package testutil

import (
	"encoding/json"
	"math/big"

	"cosmossdk.io/api/tendermint/abci"
	"cosmossdk.io/math"
	sidecar "github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	"github.com/Zenrock-Foundation/zrchain/v6/x/validation/keeper"
	cmtabci "github.com/cometbft/cometbft/abci/types"
)

var SampleSidecarState = &sidecar.SidecarStateResponse{
	EigenDelegations:   []byte(`{"validator1":{"delegator1":1000000000,"delegator2":2000000000},"validator2":{"delegator3":1500000000}}`),
	EthBlockHeight:     18500000,
	EthGasLimit:        30000000,
	EthBaseFee:         20000000000, // 20 gwei
	EthTipCap:          1000000000,  // 1 gwei
	EthBurnEvents:      []sidecar.BurnEvent{},
	Redemptions:        []sidecar.Redemption{},
	ROCKUSDPrice:       "1.25",
	BTCUSDPrice:        "45000.00",
	ETHUSDPrice:        "2800.00",
	ZECUSDPrice:        "35.00",
	SolanaBurnEvents:   []sidecar.BurnEvent{},
	SolanaMintEvents:   []sidecar.SolanaMintEvent{},
	SidecarVersionName: "test",
}

var SampleOracleData = keeper.OracleData{
	EigenDelegationsMap: map[string]map[string]*big.Int{
		"zenvaloper138a4gyfjyghrd4pvuhuezxa6cl0wd5cde3s8rd": {
			"zenvaloper138a4gyfjyghrd4pvuhuezxa6cl0wd5cde3s8rd": big.NewInt(1000000000000000000),
		},
	},
	ValidatorDelegations: []keeper.ValidatorDelegations{
		{
			Validator: "zenvaloper138a4gyfjyghrd4pvuhuezxa6cl0wd5cde3s8rd",
			Stake:     math.NewInt(1000000000000000000),
		},
	},
	RequestedBtcBlockHeight: 0,
	EthBlockHeight:          18500000,
	EthGasLimit:             30000000,
	EthBaseFee:              20000000000,
	EthTipCap:               1000000000,
	RequestedStakerNonce:    0,
	RequestedEthMinterNonce: 0,
	RequestedUnstakerNonce:  0,
	RequestedCompleterNonce: 0,
	ROCKUSDPrice:            "1.25",
	BTCUSDPrice:             "45000.00",
	ETHUSDPrice:             "2800.00",
	ZECUSDPrice:             "35.00",
	LatestBtcBlockHeight:    750000,
	LatestBtcBlockHeader: sidecar.BTCBlockHeader{
		Version:     0x20000000,
		PrevBlock:   "0000000000000000000000000000000000000000000000000000000000000000",
		MerkleRoot:  "0000000000000000000000000000000000000000000000000000000000000000",
		TimeStamp:   1640995200, // 2022-01-01 00:00:00 UTC
		Bits:        0x1d00ffff,
		Nonce:       0,
		BlockHash:   "0000000000000000000000000000000000000000000000000000000000000000",
		BlockHeight: 750000,
	},
	SidecarVersionName: "test",
	ConsensusData: cmtabci.ExtendedCommitInfo{
		Round: 1,
		Votes: []cmtabci.ExtendedVoteInfo{
			{
				Validator: cmtabci.Validator{
					Address: []byte("QDagxuKQqu3HMpWLmNIgCEhR9b0="), // Base64 encoded validator address
					Power:   125000000,                              // Higher power for consensus
				},
				BlockIdFlag:        2,                                                                                                  // BlockIDFlagCommit for consensus
				VoteExtension:      []byte{},                                                                                           // Empty vote extension
				ExtensionSignature: []byte("QB/lPpqzBJAW+iNF37X5PVrHpuHJ/ZmKWcFX6JdwTxYPAjomEHI9BqzF9EOSpp3CQ1/OikFMlITSR+eqIhgaCg=="), // Base64 encoded signature
			},
		},
	},
}

// Bitcoin header responses
var SampleBtcHeader = &sidecar.BitcoinBlockHeaderResponse{
	BlockHeader: &sidecar.BTCBlockHeader{
		Version:     0x20000000,
		PrevBlock:   "0000000000000000000000000000000000000000000000000000000000000000",
		MerkleRoot:  "0000000000000000000000000000000000000000000000000000000000000000",
		TimeStamp:   1640995200, // 2022-01-01 00:00:00 UTC
		Bits:        0x1d00ffff,
		Nonce:       0,
		BlockHash:   "0000000000000000000000000000000000000000000000000000000000000000",
		BlockHeight: 750000,
	},
	BlockHeight: 750000,
	TipHeight:   750001,
}

var SampleNonceResponse = &sidecar.LatestEthereumNonceForAccountResponse{
	Nonce: 42,
}

var SampleSolanaAccount = &sidecar.SolanaAccountInfoResponse{
	Account: []byte{0x01, 0x02, 0x03, 0x04}, // Mock account data
}

var (
	SampleVoteExtension        []byte
	SampleVoteExtension2       []byte
	SampleDecodedVoteExtension keeper.VoteExtension
)

func init() {
	// Initialize SampleVoteExtension
	ve, err := keeper.ConstructVoteExtension(&SampleOracleData)
	if err != nil {
		panic(err)
	}
	SampleDecodedVoteExtension = ve
	SampleVoteExtension = mustMarshal(ve)

	// Initialize SampleVoteExtensionHeight2
	oracleData2 := SampleOracleData
	oracleData2.EthBlockHeight = 19500000
	oracleData2.EthGasLimit = 32000000
	oracleData2.EthBaseFee = 22000000000
	oracleData2.EthTipCap = 1200000000
	ve2, err := keeper.ConstructVoteExtension(&oracleData2)
	if err != nil {
		panic(err)
	}
	SampleVoteExtension2 = mustMarshal(ve2)
}

func mustMarshal(v interface{}) []byte {
	bz, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bz
}

var SampleLastCommit = abci.CommitInfo{
	Round: 1,
	Votes: []*abci.VoteInfo{
		{
			Validator: &abci.Validator{
				Address: []byte("test-validator"),
				Power:   1000,
			},
			BlockIdFlag: 1,
		},
	},
}

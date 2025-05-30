package main_test

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"testing"

	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/go-client"
	sidecar "github.com/Zenrock-Foundation/zrchain/v6/sidecar"
	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/shared"
	"github.com/ethereum/go-ethereum/ethclient"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"
	// "github.com/ethereum/go-ethereum/accounts/abi/bind"
	// "github.com/ethereum/go-ethereum/common"
	// taskmanager "github.com/zenrocklabs/zenrock-avs/contracts/bindings/TaskManagerZR"
	// servicemanager "github.com/zenrocklabs/zenrock-avs/contracts/bindings/ZRServiceManager"
)

func initTestOracle() *sidecar.Oracle {
	cfg := sidecar.LoadConfig()

	var rpcAddress string
	if endpoint, ok := cfg.EthRPC[cfg.Network]; ok {
		rpcAddress = endpoint
	} else {
		log.Fatalf("No RPC endpoint found for network: %s", cfg.Network)
	}

	ethClient, err := ethclient.Dial(rpcAddress)
	if err != nil {
		log.Fatalf("failed to connect to the Ethereum client: %v", err)
	}

	solanaClient := solanarpc.New(cfg.SolanaRPC[cfg.Network])

	zrChainQueryClient, err := client.NewQueryClient(cfg.ZRChainRPC, true)
	if err != nil {
		log.Fatalf("Refresh Address Client: failed to get new client: %v", err)
	}

	return sidecar.NewOracle(cfg, ethClient, nil, solanaClient, zrChainQueryClient, true)
}

func TestGetSidecarStateByEthHeight(t *testing.T) {
	oracle := initTestOracle()

	// Sample states
	price1, _ := math.LegacyNewDecFromStr("123.45")
	delegations1 := map[string]map[string]*big.Int{
		"validator1": {"operator1": big.NewInt(1000)},
	}
	state1 := sidecartypes.OracleState{
		EthBlockHeight:   100,
		ROCKUSDPrice:     price1,
		EthBaseFee:       50,
		EigenDelegations: delegations1,
		// Populate other fields as necessary for thorough testing
	}

	price2, _ := math.LegacyNewDecFromStr("678.90")
	delegations2 := map[string]map[string]*big.Int{
		"validator2": {"operator2": big.NewInt(2000)},
	}
	state2 := sidecartypes.OracleState{
		EthBlockHeight:   200,
		ROCKUSDPrice:     price2,
		EthBaseFee:       75,
		EigenDelegations: delegations2,
	}

	oracle.SetStateCacheForTesting([]sidecartypes.OracleState{state1, state2})

	service := sidecar.NewOracleService(oracle)
	require.NotNil(t, service)

	// Test Case 1: State found
	t.Run("StateFound", func(t *testing.T) {
		req := &api.SidecarStateByEthHeightRequest{EthBlockHeight: 100}
		resp, err := service.GetSidecarStateByEthHeight(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, uint64(100), resp.EthBlockHeight)
		require.Equal(t, state1.ROCKUSDPrice.String(), resp.ROCKUSDPrice)
		require.Equal(t, state1.EthBaseFee, resp.EthBaseFee)

		expectedDelegations1JSON, _ := json.Marshal(state1.EigenDelegations)
		require.JSONEq(t, string(expectedDelegations1JSON), string(resp.EigenDelegations))
	})

	// Test Case 2: State not found
	t.Run("StateNotFound", func(t *testing.T) {
		req := &api.SidecarStateByEthHeightRequest{EthBlockHeight: 300} // This height does not exist in cache
		resp, err := service.GetSidecarStateByEthHeight(context.Background(), req)

		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "state with Ethereum block height 300 not found")
	})

	// Test Case 3: Requesting height from second state
	t.Run("SecondStateFound", func(t *testing.T) {
		req := &api.SidecarStateByEthHeightRequest{EthBlockHeight: 200}
		resp, err := service.GetSidecarStateByEthHeight(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, uint64(200), resp.EthBlockHeight)
		require.Equal(t, state2.ROCKUSDPrice.String(), resp.ROCKUSDPrice)
		require.Equal(t, state2.EthBaseFee, resp.EthBaseFee)

		expectedDelegations2JSON, _ := json.Marshal(state2.EigenDelegations)
		require.JSONEq(t, string(expectedDelegations2JSON), string(resp.EigenDelegations))
	})
}

// func TestGetTaskManagerAndStakeRegistryAddrs(t *testing.T) {
// 	o := initTestOracle()

// 	contractServiceManager, err := servicemanager.NewContractZRServiceManager(common.HexToAddress(o.Config.EthOracle.ContractAddrs.ServiceManager), o.EthClient)
// 	require.NoError(t, err)

// 	taskManagerAddr, err := contractServiceManager.TaskManagerZR(&bind.CallOpts{})
// 	require.NoError(t, err)
// 	fmt.Println("Task Manager Address:", taskManagerAddr)

// 	contractInstance, err := taskmanager.NewContractTaskManagerZR(taskManagerAddr, o.EthClient)
// 	require.NoError(t, err)

// 	stakeRegistryAddr, err := contractInstance.StakeRegistry(&bind.CallOpts{})
// 	require.NoError(t, err)
// 	fmt.Println("Stake Registry Address:", stakeRegistryAddr)
// }

func TestGetSolanaAccountInfo(t *testing.T) {
	oracle := initTestOracle()
	service := sidecar.NewOracleService(oracle)
	require.NotNil(t, service)

	tests := []struct {
		recipient     string
		pubKey        string
		expectError   bool
		expectedError string
	}{
		{
			recipient:   "HvD2fkmc19X2HMMzxpf2hN7HBHLyeftWof2hwVcbHiWK (current)",
			pubKey:      "GeHPYFhNxKdcoP1Kbf2Twc3RKEsgVh2x1aYrLb7hd63d",
			expectError: false,
		},
		{
			recipient:   "HvD2fkmc19X2HMMzxpf2hN7HBHLyeftWof2hwVcbHiWK (correct)",
			pubKey:      "B3ZjSXnczNYXQyhQft6wiNk99sW4NfW3Fe2TwhF7wTH8",
			expectError: false,
		},
		{
			recipient:   "2RoRSPwFmMcFDd3DLgoomMRFkhUSaVNSDo7XyH2tw6de",
			pubKey:      "FYsom9kuL835ZcRtHCiNtKeb7qr6h8NJ5fwfZA8Rsa3y",
			expectError: false,
		},
		{
			recipient:     "2RoRSPwFmMcFDd3DLgoomMRFkhUSaVNSDo7XyH2tw6de",
			pubKey:        "InvalidPubKey",
			expectError:   true,
			expectedError: "decode: invalid base58 digit ('I')",
		},
	}

	for _, tt := range tests {
		t.Run(tt.recipient, func(t *testing.T) {
			req := &api.SolanaAccountInfoRequest{PubKey: tt.pubKey}
			resp, err := service.GetSolanaAccountInfo(context.Background(), req)

			if tt.expectError {
				require.Error(t, err)
				if tt.expectedError != "" {
					require.Contains(t, err.Error(), tt.expectedError)
				}
				log.Printf("\nTest %s (expected error): Received error: %v\n", tt.recipient, err)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				log.Printf("\nTest %s: Received account info: %v\n", tt.recipient, resp.Account)
				// If the account is not found, resp.Account will be nil or empty
				// The original function returns an empty response if rpc.ErrNotFound, so resp itself is not nil
				if resp.Account == nil || len(resp.Account) == 0 {
					log.Printf("Test %s: Account not found or has no data.\n", tt.recipient)
				}
			}
		})
	}
}

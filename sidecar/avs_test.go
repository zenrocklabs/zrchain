package main_test

import (
	// "fmt"
	"log"
	// "testing"
	"time"

	"github.com/Zenrock-Foundation/zrchain/v5/go-client"
	sidecar "github.com/Zenrock-Foundation/zrchain/v5/sidecar"
	// "github.com/ethereum/go-ethereum/accounts/abi/bind"
	// "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
	// "github.com/stretchr/testify/require"
	// taskmanager "github.com/zenrocklabs/zenrock-avs/contracts/bindings/TaskManagerZR"
	// servicemanager "github.com/zenrocklabs/zenrock-avs/contracts/bindings/ZRServiceManager"
)

func initTestOracle() *sidecar.Oracle {
	cfg := sidecar.LoadConfig()

	var rpcAddress string
	if endpoint, ok := cfg.EthOracle.RPC[cfg.Network]; ok {
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

	return sidecar.NewOracle(cfg, ethClient, nil, solanaClient, zrChainQueryClient, time.NewTicker(sidecar.MainLoopTickerInterval))
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

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	neutrino "github.com/Zenrock-Foundation/zrchain/v5/sidecar/neutrino"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v5/sidecar/shared"

	"github.com/ethereum/go-ethereum/ethclient"

	solana "github.com/gagliardetto/solana-go/rpc"
)

func main() {
	port := flag.Int("port", 0, "Override GRPC port from config")
	flag.Parse()

	cfg := LoadConfig()

	// Override GRPC port if --port flag is provided
	if *port != 0 {
		cfg.GRPCPort = *port
	}

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	neutrinoServer := neutrino.NeutrinoServer{}
	neutrinoServer.Initialize(cfg.ProxyRPC.URL, cfg.ProxyRPC.User, cfg.ProxyRPC.Password, cfg.Neutrino.Path)

	solanaClient := solana.New(cfg.SolanaRPC[cfg.Network])

	// Align the start time to the nearest MainLoopTickerInterval
	now := time.Now()
	alignedStart := now.Truncate(MainLoopTickerInterval).Add(MainLoopTickerInterval)
	time.Sleep(alignedStart.Sub(now))

	mainLoopTicker := time.NewTicker(MainLoopTickerInterval)
	defer mainLoopTicker.Stop()

	oracle := NewOracle(cfg, ethClient, &neutrinoServer, solanaClient, mainLoopTicker)

	go startGRPCServer(oracle, cfg.GRPCPort)

	log.Printf("gRPC server listening on port %d", cfg.GRPCPort)
	log.Printf("Please wait ~%ds before launching the zrChain node for the first Ethereum state and price updates\n", MainLoopTickerInterval/time.Second)

	go func() {
		if err := oracle.runAVSContractOracleLoop(ctx); err != nil {
			log.Printf("Error in Ethereum oracle loop: %v", err)
		}
	}()

	go oracle.processUpdates()

	go func() {
		if err := oracle.runEigenOperator(); err != nil {
			log.Fatalf("Error starting EigenLayer Operator: %v", err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down gracefully...")
	cancel()
}

func (o *Oracle) processUpdates() {
	for update := range o.updateChan {
		log.Printf("Received AVS contract state for %s block %d", o.Config.EthOracle.NetworkName[o.Config.Network], update.EthBlockHeight)
		currentState := o.currentState.Load().(*sidecartypes.OracleState)
		newState := *currentState

		newState.EigenDelegations = update.EigenDelegations
		newState.EthBlockHeight = update.EthBlockHeight
		newState.EthGasLimit = update.EthGasLimit
		newState.EthBaseFee = update.EthBaseFee
		newState.EthTipCap = update.EthTipCap
		newState.SolanaLamportsPerSignature = update.SolanaLamportsPerSignature
		newState.RedemptionsEthereum = update.RedemptionsEthereum
		newState.RedemptionsSolana = update.RedemptionsSolana

		log.Printf("Received prices: ROCK/USD %f, BTC/USD %f", update.ROCKUSDPrice, update.BTCUSDPrice)
		newState.ROCKUSDPrice = update.ROCKUSDPrice
		newState.BTCUSDPrice = update.BTCUSDPrice
		newState.ETHUSDPrice = update.ETHUSDPrice
		o.currentState.Store(&newState)

		o.CacheState()
	}
}

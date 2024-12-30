package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
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
	cfg := LoadConfig()

	if !cfg.Enabled {
		for {
			slog.Info("Sidecar is disabled in config; sleeping...")
			time.Sleep(time.Hour)
		}
	}

	port := flag.Int("port", 0, "Override GRPC port from config")
	flag.Parse()

	// Override defeault GRPC port if --port flag is provided
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

	slog.Info("gRPC server listening on port", "port", cfg.GRPCPort)
	slog.Info("Please wait ~%ds before launching the zrChain node for the first Ethereum state and price updates", "seconds", MainLoopTickerInterval/time.Second)

	go func() {
		if err := oracle.runAVSContractOracleLoop(ctx); err != nil {
			slog.Error("Error in Ethereum oracle loop", "error", err)
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

	slog.Info("Shutting down gracefully...")
	cancel()
}

func (o *Oracle) processUpdates() {
	for update := range o.updateChan {
		slog.Info("Received AVS contract state for", "network", o.Config.EthOracle.NetworkName[o.Config.Network], "block", update.EthBlockHeight)
		currentState := o.currentState.Load().(*sidecartypes.OracleState)
		newState := *currentState

		newState.EigenDelegations = update.EigenDelegations
		newState.EthBlockHeight = update.EthBlockHeight
		newState.EthBaseFee = update.EthBaseFee
		newState.EthTipCap = update.EthTipCap
		newState.EthWrapGasLimit = update.EthWrapGasLimit
		newState.EthUnstakeGasLimit = update.EthUnstakeGasLimit
		newState.SolanaLamportsPerSignature = update.SolanaLamportsPerSignature
		newState.RedemptionsEthereum = update.RedemptionsEthereum
		newState.RedemptionsSolana = update.RedemptionsSolana

		slog.Info("Received prices", "ROCK/USD", update.ROCKUSDPrice, "BTC/USD", update.BTCUSDPrice, "ETH/USD", update.ETHUSDPrice)
		newState.ROCKUSDPrice = update.ROCKUSDPrice
		newState.BTCUSDPrice = update.BTCUSDPrice
		newState.ETHUSDPrice = update.ETHUSDPrice
		o.currentState.Store(&newState)

		o.CacheState()
	}
}

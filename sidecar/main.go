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

	"github.com/Zenrock-Foundation/zrchain/v5/go-client"
	neutrino "github.com/Zenrock-Foundation/zrchain/v5/sidecar/neutrino"

	"github.com/ethereum/go-ethereum/ethclient"
	solana "github.com/gagliardetto/solana-go/rpc"

	"github.com/beevik/ntp"
)

func main() {
	port := flag.Int("port", 0, "Override GRPC port from config")
	cacheFile := flag.String("cache-file", "", "Override cache file path from config")

	if !flag.Parsed() {
		flag.Parse()
	}

	cfg := LoadConfig()

	if !cfg.Enabled {
		for {
			slog.Info("Sidecar is disabled in config; sleeping...")
			time.Sleep(time.Hour)
		}
	}

	// Override default GRPC port if --port flag is provided
	if *port != 0 {
		cfg.GRPCPort = *port
	}

	// Override default state file path if --cache-file flag is provided
	if *cacheFile != "" {
		cfg.StateFile = *cacheFile
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

	zrChainQueryClient, err := client.NewQueryClient(cfg.ZRChainRPC, true)
	if err != nil {
		log.Fatalf("Refresh Address Client: failed to get new client: %v", err)
	}

	ntpTime, err := ntp.Time("time.google.com")
	if err != nil {
		log.Fatalf("Error fetching NTP time: %v", err)
	}
	// Align the start time to the nearest MainLoopTickerInterval
	alignedStart := ntpTime.Truncate(MainLoopTickerInterval).Add(MainLoopTickerInterval)
	time.Sleep(alignedStart.Sub(ntpTime))

	mainLoopTicker := time.NewTicker(MainLoopTickerInterval)
	defer mainLoopTicker.Stop()

	oracle := NewOracle(cfg, ethClient, &neutrinoServer, solanaClient, zrChainQueryClient, mainLoopTicker)

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
		slog.Info("Received prices", "ROCK/USD", update.ROCKUSDPrice, "BTC/USD", update.BTCUSDPrice, "ETH/USD", update.ETHUSDPrice)

		newState := update
		o.currentState.Store(&newState)
		o.CacheState()
	}
}

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

	"github.com/Zenrock-Foundation/zrchain/v6/go-client"
	neutrino "github.com/Zenrock-Foundation/zrchain/v6/sidecar/neutrino"

	"github.com/ethereum/go-ethereum/ethclient"
	solana "github.com/gagliardetto/solana-go/rpc"
)

func main() {
	port := flag.Int("port", 9191, "Override GRPC port from config")
	cacheFile := flag.String("cache-file", "cache.json", "Override cache file path from config")
	neutrinoPort := flag.Int("neutrino-port", 12345, "Override Neutrino RPC port (default: 12345)")
	ethRPC := flag.String("eth-rpc", "", "Override Ethereum RPC endpoint from config")
	neutrinoPath := flag.String("neutrino-path", "/neutrino_", "Path prefix for neutrino directory")
	noAVS := flag.Bool("no-avs", false, "Disable EigenLayer Operator (AVS)")
	debug := flag.Bool("debug", false, "Enable debug mode for verbose logging")

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

	// Set Neutrino port from flag or config
	neutrinoRPCPort := *neutrinoPort
	if neutrinoRPCPort == 0 && cfg.Neutrino.Port > 0 {
		neutrinoRPCPort = cfg.Neutrino.Port
	}

	var rpcAddress string
	// Use the override RPC endpoint if provided via flag
	if *ethRPC != "" {
		rpcAddress = *ethRPC
		slog.Info("Using override Ethereum RPC endpoint", "endpoint", rpcAddress)
	} else if endpoint, ok := cfg.EthRPC[cfg.Network]; ok {
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
	neutrinoServer.Initialize(cfg.ProxyRPC.URL, cfg.ProxyRPC.User, cfg.ProxyRPC.Password, cfg.Neutrino.Path, neutrinoRPCPort, *neutrinoPath)

	solanaClient := solana.New(cfg.SolanaRPC[cfg.Network])

	zrChainQueryClient, err := client.NewQueryClient(cfg.ZRChainRPC, true)
	if err != nil {
		log.Fatalf("Refresh Address Client: failed to get new client: %v", err)
	}

	oracle := NewOracle(cfg, ethClient, &neutrinoServer, solanaClient, zrChainQueryClient, *debug)

	go startGRPCServer(oracle, cfg.GRPCPort)

	slog.Info("gRPC server listening on port", "port", cfg.GRPCPort)

	go func() {
		if err := oracle.runOracleMainLoop(ctx); err != nil {
			slog.Error("Error in Ethereum oracle loop", "error", err)
		}
	}()

	if !*noAVS {
		go func() {
			if err := oracle.runEigenOperator(); err != nil {
				log.Fatalf("Error starting EigenLayer Operator: %v", err)
			}
		}()
	}

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	slog.Info("Shutting down gracefully...")
	cancel()
}

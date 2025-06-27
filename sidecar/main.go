package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Zenrock-Foundation/zrchain/v6/go-client"
	neutrino "github.com/Zenrock-Foundation/zrchain/v6/sidecar/neutrino"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/shared"

	"github.com/ethereum/go-ethereum/ethclient"
	solana "github.com/gagliardetto/solana-go/rpc"
)

func main() {
	port := flag.Int("port", 9191, "Override GRPC port from config")
	cacheFile := flag.String("cache-file", "cache.json", "Override cache file path from config")
	neutrinoPort := flag.Int("neutrino-port", 12345, "Override Neutrino RPC port (default: 12345)")
	ethRPC := flag.String("eth-rpc", "", "Override Ethereum RPC endpoint from config")
	neutrinoPath := flag.String("neutrino-path", "/neutrino_", "Path prefix for neutrino directory")
	debug := flag.Bool("debug", false, "Enable debug mode for verbose logging")
	// DEBUGGING ONLY - RISK OF SLASHING IF USED IN PRODUCTION
	noAVS := flag.Bool("no-avs", false, "Disable EigenLayer Operator (AVS)")
	skipInitialWait := flag.Bool("skip-initial-wait", false, "Skip initial NTP alignment wait and fire tick immediately")
	version := flag.Bool("version", false, "Display version information and exit")

	if !flag.Parsed() {
		flag.Parse()
	}

	// Handle version command
	if *version {
		slog.Info("zrChain Validator Sidecar", "version", sidecartypes.SidecarVersionName)
		os.Exit(0)
	}

	cfg := LoadConfig()

	slog.Info("Starting zrChain Validator Sidecar", "version", sidecartypes.SidecarVersionName)

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

	// Reset state if version requires it – firstBoot will be true only once per version
	firstBoot := resetStateForVersion(cfg.StateFile)
	if firstBoot {
		slog.Info("Completed first-boot cache reset for zrChain Validator Sidecar", "version", sidecartypes.SidecarVersionName)
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
		slog.Error("No RPC endpoint found for network", "network", cfg.Network)
		os.Exit(1)
	}

	ethClient, err := ethclient.Dial(rpcAddress)
	if err != nil {
		slog.Error("Failed to connect to the Ethereum client", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	neutrinoServer := neutrino.NeutrinoServer{}
	neutrinoServer.Initialize(cfg.ProxyRPC.URL, cfg.ProxyRPC.User, cfg.ProxyRPC.Password, cfg.Neutrino.Path, neutrinoRPCPort, *neutrinoPath)

	solanaClient := solana.New(cfg.SolanaRPC[cfg.Network])

	zrChainQueryClient, err := client.NewQueryClient(cfg.ZRChainRPC, true)
	if err != nil {
		slog.Error("Refresh Address Client: failed to get new client", "error", err)
		os.Exit(1)
	}

	oracle := NewOracle(cfg, ethClient, &neutrinoServer, solanaClient, zrChainQueryClient, *debug, *skipInitialWait)

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
				slog.Error("Error starting EigenLayer Operator", "error", err)
				os.Exit(1)
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

package main

import (
	"context"
	"flag"
	"fmt"
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
	// Parse flags first to determine debug level
	flags := parseFlags()

	// Set up coloured structured logging
	initLogger(*flags.debug)

	// Handle version command
	if *flags.version {
		slog.Info("zrChain Validator Sidecar", "version", sidecartypes.SidecarVersionName)
		os.Exit(0)
	}

	cfg := LoadConfig(*flags.configFile, *flags.configDir)

	slog.Info("Starting zrChain Validator Sidecar", "version", sidecartypes.SidecarVersionName)

	// Validate basic configuration
	if err := validateConfiguration(cfg); err != nil {
		slog.Error("Configuration validation failed", "error", err)
		os.Exit(1)
	}

	if !cfg.Enabled {
		for {
			slog.Info("Sidecar is disabled in config; sleeping...")
			time.Sleep(time.Hour)
		}
	}

	// Override default GRPC port if --port flag is provided
	if *flags.port != 0 {
		cfg.GRPCPort = *flags.port
	}

	// Override default state file path if --cache-file flag is provided
	if *flags.cacheFile != "" {
		cfg.StateFile = *flags.cacheFile
	}

	// Reset state if version requires it – firstBoot will be true only once per version
	firstBoot := resetStateForVersion(cfg.StateFile)
	if firstBoot {
		slog.Info("Completed first-boot cache reset for zrChain Validator Sidecar", "version", sidecartypes.SidecarVersionName)
	}

	// Set Neutrino port from flag or config
	neutrinoRPCPort := *flags.neutrinoPort
	if neutrinoRPCPort == 0 && cfg.Neutrino.Port > 0 {
		neutrinoRPCPort = cfg.Neutrino.Port
	}

	// Validate and create Ethereum client
	ethClient, err := validateEthereumClient(cfg, *flags.ethRPC)
	if err != nil {
		slog.Error("Ethereum client validation failed", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	neutrinoServer := neutrino.NeutrinoServer{}
	neutrinoServer.Initialize(cfg.Network, cfg.ProxyRPC.URL, cfg.ProxyRPC.User, cfg.ProxyRPC.Password, cfg.Neutrino.Path, neutrinoRPCPort, *flags.neutrinoPath)

	// Validate and create Solana client
	solanaClient, err := validateSolanaClient(cfg, *flags.noSolana)
	if err != nil {
		slog.Error("Solana client validation failed", "error", err)
		os.Exit(1)
	}

	// Create zrChain client without validation for better UX
	zrChainQueryClient, err := client.NewQueryClient(cfg.ZRChainRPC, true)
	if err != nil {
		slog.Error("Failed to create zrChain client", "error", err)
		os.Exit(1)
	}

	oracle := NewOracle(cfg, ethClient, &neutrinoServer, solanaClient, zrChainQueryClient, *flags.debug, *flags.skipInitialWait, *flags.testReset)

	go startGRPCServer(oracle, cfg.GRPCPort)

	slog.Info("gRPC server listening on port", "port", cfg.GRPCPort)

	go func() {
		if err := oracle.runOracleMainLoop(ctx); err != nil {
			slog.Error("Error in Ethereum oracle loop", "error", err)
		}
	}()

	if !*flags.noAVS {
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

// flagConfig holds all the parsed command-line flags
type flagConfig struct {
	port            *int
	cacheFile       *string
	neutrinoPort    *int
	ethRPC          *string
	neutrinoPath    *string
	debug           *bool
	configFile      *string
	configDir       *string
	version         *bool
	noAVS           *bool
	skipInitialWait *bool
	noSolana        *bool
	testReset       *bool
}

// parseFlags sets up and parses all command-line flags
func parseFlags() *flagConfig {
	flags := &flagConfig{
		port:            flag.Int("port", 9191, "Override GRPC port from config"),
		cacheFile:       flag.String("cache-file", "cache.json", "Override cache file path from config"),
		neutrinoPort:    flag.Int("neutrino-port", 12345, "Override Neutrino RPC port (default: 12345)"),
		ethRPC:          flag.String("eth-rpc", "", "Override Ethereum RPC endpoint from config"),
		neutrinoPath:    flag.String("neutrino-path", "/neutrino_", "Path prefix for neutrino directory"),
		debug:           flag.Bool("debug", false, "Enable debug mode for verbose logging"),
		configFile:      flag.String("config", "", "Override config file path (default: config.yaml)"),
		configDir:       flag.String("config-dir", "", "Directory to search for config.yaml"),
		version:         flag.Bool("version", false, "Display version information and exit"),
		noAVS:           flag.Bool("no-avs", false, "Disable EigenLayer Operator (AVS)"),
		skipInitialWait: flag.Bool("skip-initial-wait", false, "Skip initial NTP alignment wait and fire tick immediately"),
		noSolana:        flag.Bool("no-solana", false, "Disable Solana functionality for testing"),
		testReset:       flag.Bool("test-reset", false, "Force periodic oracle state reset every 2 minutes (testing only)"),
	}

	if !flag.Parsed() {
		flag.Parse()
	}

	return flags
}

// validateEthereumClient creates and validates an Ethereum client connection
func validateEthereumClient(cfg sidecartypes.Config, ethRPCOverride string) (*ethclient.Client, error) {
	var rpcAddress string
	// Use the override RPC endpoint if provided via flag
	if ethRPCOverride != "" {
		rpcAddress = ethRPCOverride
		slog.Info("Using override Ethereum RPC endpoint", "endpoint", rpcAddress)
	} else if endpoint, ok := cfg.EthRPC[cfg.Network]; ok {
		rpcAddress = endpoint
	} else {
		return nil, fmt.Errorf("no Ethereum RPC endpoint found for network: %s", cfg.Network)
	}

	ethClient, err := connectEthereumWithRetry(rpcAddress, sidecartypes.RPCConnectionMaxRetries, sidecartypes.RPCConnectionRetryDelay)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	// Additional validation: Test basic Ethereum connectivity
	ethHealthCtx, ethHealthCancel := context.WithTimeout(context.Background(), 10*time.Second)
	chainID, err := ethClient.ChainID(ethHealthCtx)
	ethHealthCancel()
	if err != nil {
		ethClient.Close()
		return nil, fmt.Errorf("failed to verify Ethereum client connectivity - ChainID call failed: %w", err)
	}

	slog.Info("Successfully verified Ethereum client connectivity", "endpoint", rpcAddress, "chainID", chainID)
	return ethClient, nil
}

// validateSolanaClient creates and validates a Solana client connection if not disabled
func validateSolanaClient(cfg sidecartypes.Config, noSolana bool) (*solana.Client, error) {
	if noSolana {
		slog.Info("Solana functionality disabled for testing")
		return nil, nil
	}

	// Validate Solana RPC endpoint is configured
	solanaEndpoint, ok := cfg.SolanaRPC[cfg.Network]
	if !ok || solanaEndpoint == "" {
		return nil, fmt.Errorf("no Solana RPC endpoint configured for network: %s. Consider using --no-solana flag if Solana functionality is not needed", cfg.Network)
	}

	solanaClient, err := connectSolanaWithRetry(solanaEndpoint, sidecartypes.RPCConnectionMaxRetries, sidecartypes.RPCConnectionRetryDelay)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Solana client: %w. Consider using --no-solana flag if Solana functionality is not needed", err)
	}

	return solanaClient, nil
}

// validateZrChainClient creates and validates a zrChain query client connection
// func validateZrChainClient(cfg sidecartypes.Config) (*client.QueryClient, error) {
// 	zrChainQueryClient, err := connectZrChainWithRetry(cfg.ZRChainRPC, sidecartypes.RPCConnectionMaxRetries, sidecartypes.RPCConnectionRetryDelay)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect to zrChain client: %w", err)
// 	}

// 	slog.Info("Successfully verified zrChain client connectivity", "endpoint", cfg.ZRChainRPC)
// 	return zrChainQueryClient, nil
// }

// validateConfiguration performs basic configuration validation
func validateConfiguration(cfg sidecartypes.Config) error {
	if cfg.ZRChainRPC == "" {
		return fmt.Errorf("zrChain RPC endpoint is required but not configured")
	}
	return nil
}

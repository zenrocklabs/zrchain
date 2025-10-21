package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Zenrock-Foundation/zrchain/v6/go-client"
	sidecarshared "github.com/Zenrock-Foundation/zrchain/v6/sidecar/shared"
	dcttypes "github.com/Zenrock-Foundation/zrchain/v6/x/dct/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptokit "github.com/cosmos/cosmos-sdk/crypto"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
	"golang.org/x/term"
	gorpc "google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

func main() {
	var (
		configFile        string
		configDir         string
		networkOverride   string
		solanaRPCOverride string
		nodeOverride      string
		cosmosChainID     string
		solanaCAIP2       string
		assetFlag         string
		zenbtc            bool
		mnemonic          string
		derivationPath    string
		fromName          string
		keyringBackend    string
		keyringHome       string
		exportPass        string
		insecureGRPC      bool
		timeout           time.Duration
	)

	flag.StringVar(&configFile, "config", "", "Path to sidecar config file (overrides autodetect)")
	flag.StringVar(&configDir, "config-dir", "../..", "Directory containing sidecar config file")
	flag.StringVar(&networkOverride, "network", "", "Override network value from config (optional)")
	flag.StringVar(&nodeOverride, "node", "grpc.dev.zenrock.tech:443", "Override ZRChain gRPC endpoint")
	flag.StringVar(&solanaRPCOverride, "solana-rpc", "", "Override Solana RPC endpoint (optional)")
	flag.StringVar(&cosmosChainID, "chain-id", "amber-1", "Cosmos SDK chain ID")
	flag.StringVar(&solanaCAIP2, "solana-caip2", "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1", "Solana CAIP-2 identifier (defaults per network)")
	flag.StringVar(&assetFlag, "asset", "", "DCT asset to advance (e.g. zenzec)")
	flag.BoolVar(&zenbtc, "zenbtc", false, "Advance zenBTC nonce account")
	flag.StringVar(&mnemonic, "mnemonic", "", "Mnemonic for the authority account (falls back to ADVANCE_NONCE_MNEMONIC env)")
	flag.StringVar(&derivationPath, "derivation-path", "m/44'/118'/0'/0/0", "Derivation path for authority account")
	flag.StringVar(&fromName, "from", "", "Key name in the Cosmos SDK keyring to sign the request with")
	flag.StringVar(&keyringBackend, "keyring-backend", "file", "Keyring backend to use (os|file|test)")
	flag.StringVar(&keyringHome, "home", filepath.Join(os.Getenv("HOME"), ".zenrockd"), "Keyring home directory")
	flag.StringVar(&exportPass, "key-passphrase", "", "Passphrase to encrypt the exported key (optional; prompt if empty)")
	flag.BoolVar(&insecureGRPC, "grpc-insecure", false, "Use insecure gRPC connection")
	flag.DurationVar(&timeout, "timeout", 30*time.Second, "Overall timeout for the operation")
	flag.Parse()

	if cosmosChainID == "" {
		log.Fatalf("chain-id must be provided")
	}

	cfg, err := loadConfig(configFile, configDir)
	if err != nil {
		log.Fatalf("failed to load sidecar config: %v", err)
	}

	network := cfg.Network
	if networkOverride != "" {
		network = networkOverride
	}

	node := nodeOverride
	if node == "" {
		node = cfg.ZRChainRPC
	}

	solanaRPC := solanaRPCOverride
	if solanaRPC == "" {
		if endpoint, ok := cfg.SolanaRPC[network]; ok && endpoint != "" {
			solanaRPC = endpoint
		} else {
			log.Fatalf("solana RPC endpoint not configured for network %s", network)
		}
	}

	if solanaCAIP2 == "" {
		if caip, ok := sidecarshared.SolanaCAIP2[network]; ok && caip != "" {
			solanaCAIP2 = caip
		} else {
			log.Fatalf("solana CAIP-2 identifier not available for network %s; provide --solana-caip2", network)
		}
	}

	if !zenbtc && assetFlag == "" {
		log.Fatalf("either --zenbtc must be true or --asset must be provided")
	}
	if zenbtc && assetFlag != "" {
		log.Printf("warning: --asset is ignored when --zenbtc is true")
	}

	var assetEnum dcttypes.Asset = dcttypes.Asset_ASSET_UNSPECIFIED
	if !zenbtc {
		var err error
		assetEnum, err = parseAsset(assetFlag)
		if err != nil {
			log.Fatalf("invalid asset %q: %v", assetFlag, err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	blockhash, err := fetchRecentBlockhash(ctx, solanaRPC)
	if err != nil {
		log.Fatalf("failed to fetch recent blockhash: %v", err)
	}

	var identity client.Identity
	switch {
	case fromName != "":
		id, loadErr := loadIdentityFromKeyring(fromName, keyringBackend, keyringHome, exportPass)
		if loadErr != nil {
			log.Fatalf("failed to load key %q from keyring: %v", fromName, loadErr)
		}
		identity = id
	default:
		if mnemonic == "" {
			mnemonic = os.Getenv("ADVANCE_NONCE_MNEMONIC")
		}
		if mnemonic == "" {
			log.Fatalf("authority mnemonic must be provided via --mnemonic, ADVANCE_NONCE_MNEMONIC, or --from")
		}
		id, deriveErr := client.NewIdentityFromSeed(derivationPath, mnemonic)
		if deriveErr != nil {
			log.Fatalf("failed to derive identity: %v", deriveErr)
		}
		identity = id
	}

	grpcConn, err := client.NewClientConn(node, insecureGRPC)
	if err != nil {
		log.Fatalf("failed to create gRPC connection: %v", err)
	}
	defer safeClose(grpcConn)

	queryClient := client.NewQueryClientWithConn(grpcConn)

	txClient, err := client.NewTxClient(identity, cosmosChainID, grpcConn, queryClient)
	if err != nil {
		log.Fatalf("failed to create tx client: %v", err)
	}

	validationClient := client.NewValidationTxClient(txClient.RawTxClient)

	log.Printf("Submitting advance nonce transaction (blockhash=%s, network=%s, caip2=%s, zenbtc=%t, asset=%s)", blockhash, network, solanaCAIP2, zenbtc, assetEnum.String())

	txHash, err := validationClient.AdvanceSolanaNonce(ctx, blockhash, solanaCAIP2, zenbtc, assetEnum)
	if err != nil {
		log.Fatalf("failed to submit advance nonce transaction: %v", err)
	}

	log.Printf("Advance nonce transaction accepted: %s", txHash)
}

func parseAsset(input string) (dcttypes.Asset, error) {
	key := strings.TrimSpace(strings.ToUpper(input))
	if key == "" {
		return dcttypes.Asset_ASSET_UNSPECIFIED, fmt.Errorf("asset string is empty")
	}
	if !strings.HasPrefix(key, "ASSET_") {
		key = "ASSET_" + key
	}
	value, ok := dcttypes.Asset_value[key]
	if !ok {
		return dcttypes.Asset_ASSET_UNSPECIFIED, fmt.Errorf("unknown asset key %s", key)
	}
	return dcttypes.Asset(value), nil
}

func fetchRecentBlockhash(ctx context.Context, endpoint string) (string, error) {
	client := solanarpc.New(endpoint)
	result, err := client.GetLatestBlockhash(ctx, solanarpc.CommitmentConfirmed)
	if err != nil {
		return "", err
	}
	if result == nil || result.Value == nil || result.Value.Blockhash.IsZero() {
		return "", fmt.Errorf("received empty blockhash from endpoint %s", endpoint)
	}
	return result.Value.Blockhash.String(), nil
}

func safeClose(conn *gorpc.ClientConn) {
	if conn == nil {
		return
	}
	if err := conn.Close(); err != nil {
		log.Printf("warning: failed to close gRPC connection: %v", err)
	}
}

func loadConfig(configFileFlag, configDirFlag string) (sidecarshared.Config, error) {
	configPath, err := resolveConfigPath(configFileFlag, configDirFlag)
	if err != nil {
		return sidecarshared.Config{}, err
	}
	return readConfig(configPath)
}

func resolveConfigPath(configFileFlag, configDirFlag string) (string, error) {
	if configFileFlag != "" {
		slog.Info("Using config file specified by flag", "path", configFileFlag)
		return configFileFlag, nil
	}

	if envPath := os.Getenv("SIDECAR_CONFIG_FILE"); envPath != "" {
		slog.Info("Using config file specified by SIDECAR_CONFIG_FILE", "path", envPath)
		return envPath, nil
	}

	if configDirFlag != "" {
		path := filepath.Join(configDirFlag, "config.yaml")
		slog.Info("Using config file from --config-dir", "path", path)
		return path, nil
	}

	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		localConfig := filepath.Join(exeDir, "config.yaml")
		if _, err := os.Stat(localConfig); err == nil {
			slog.Info("Using config file co-located with executable", "path", localConfig)
			return localConfig, nil
		}

		parentConfig := filepath.Join(filepath.Dir(exeDir), "config.yaml")
		if _, err := os.Stat(parentConfig); err == nil {
			slog.Info("Using config file from parent directory", "path", parentConfig)
			return parentConfig, nil
		}
	}

	slog.Info("Falling back to config.yaml in current working directory")
	return "config.yaml", nil
}

func readConfig(path string) (sidecarshared.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return sidecarshared.Config{}, fmt.Errorf("unable to read config from %s: %w", path, err)
	}

	var cfg sidecarshared.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return sidecarshared.Config{}, fmt.Errorf("failed to decode config %s: %w", path, err)
	}

	if cfg.ZRChainRPC == "" {
		return sidecarshared.Config{}, fmt.Errorf("zrchain_rpc must be set in %s", path)
	}
	if cfg.Network == "" {
		return sidecarshared.Config{}, fmt.Errorf("network must be set in %s", path)
	}

	return cfg, nil
}

func loadIdentityFromKeyring(name, backend, home, exportPass string) (client.Identity, error) {
	backend = strings.ToLower(strings.TrimSpace(backend))
	switch backend {
	case "":
		backend = keyring.BackendOS
	case keyring.BackendOS, keyring.BackendFile, keyring.BackendTest:
	default:
		return client.Identity{}, fmt.Errorf("unsupported keyring backend %q", backend)
	}

	interfaceRegistry := codecTypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)
	kr, err := keyring.New("zenrockd", backend, home, bufio.NewReader(os.Stdin), cdc)
	if err != nil {
		return client.Identity{}, fmt.Errorf("failed to open keyring at %s: %w", home, err)
	}

	info, err := kr.Key(name)
	if err != nil {
		return client.Identity{}, fmt.Errorf("failed to locate key %q: %w", name, err)
	}

	if backend == keyring.BackendOS {
		return client.Identity{}, errors.New("key exporting is not supported for --keyring-backend os; use --keyring-backend file/test or supply --mnemonic")
	}

	passphrase := exportPass
	if passphrase == "" {
		if env := os.Getenv("ADVANCE_NONCE_KEY_PASSPHRASE"); env != "" {
			passphrase = env
		}
	}
	if passphrase == "" {
		fmt.Printf("Enter passphrase to encrypt exported key for %q: ", name)
		bytes, readErr := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if readErr != nil {
			return client.Identity{}, fmt.Errorf("failed to read passphrase: %w", readErr)
		}
		passphrase = string(bytes)
	}

	armor, err := kr.ExportPrivKeyArmor(name, passphrase)
	if err != nil {
		return client.Identity{}, fmt.Errorf("failed to export private key: %w", err)
	}

	privKey, _, err := cryptokit.UnarmorDecryptPrivKey(armor, passphrase)
	if err != nil {
		return client.Identity{}, fmt.Errorf("failed to decrypt exported key: %w", err)
	}
	secpKey, ok := privKey.(*secp256k1.PrivKey)
	if !ok {
		return client.Identity{}, errors.New("decrypted key is not secp256k1")
	}

	addr, err := info.GetAddress()
	if err != nil {
		return client.Identity{}, fmt.Errorf("failed to derive address: %w", err)
	}

	return client.Identity{
		Address: sdktypes.AccAddress(addr),
		PrivKey: secpKey,
	}, nil
}

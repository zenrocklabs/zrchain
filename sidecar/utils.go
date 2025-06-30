package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"math/big"
	"os"
	"slices"
	"time"

	"github.com/Zenrock-Foundation/zrchain/v6/go-client"
	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/shared"
	"github.com/ethereum/go-ethereum/ethclient"
	solana "github.com/gagliardetto/solana-go"
	"github.com/gookit/color"
	"github.com/lmittmann/tint"
	"gopkg.in/yaml.v3"
)

// loadStateDataFromFile reads the state file and returns the latest state,
// all historical states, and any error.
// If the file does not exist or is empty/invalid, it returns (nil, nil, nil)
// to indicate a fresh start, unless a critical error occurs.
func loadStateDataFromFile(filename string) (latestState *sidecartypes.OracleState, historicalStates []sidecartypes.OracleState, err error) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// File not found, signal fresh start, not an error for the caller
			return nil, nil, nil
		}
		// Other error opening file
		return nil, nil, fmt.Errorf("failed to open state file %s: %w", filename, err)
	}
	defer file.Close()

	var states []sidecartypes.OracleState
	if err := json.NewDecoder(file).Decode(&states); err != nil {
		// Check for EOF which might mean an empty JSON array `[]` or just empty file
		// For an empty array or truly empty file, we treat it as a fresh start.
		// For other decode errors, it's a problem.
		// Common empty/malformed cases for json.Decoder include "EOF" for completely empty or non-JSON file,
		// and "unexpected end of JSON input" for incomplete JSON.
		if err.Error() == "EOF" || err.Error() == "unexpected end of JSON input" {
			log.Printf("State file %s is empty or contains invalid JSON, treating as fresh start.", filename)
			return nil, nil, nil
		}
		return nil, nil, fmt.Errorf("failed to decode state file %s: %w", filename, err)
	}

	if len(states) == 0 {
		// File contained an empty list of states
		log.Printf("State file %s contained no states, treating as fresh start.", filename)
		return nil, nil, nil
	}

	return &states[len(states)-1], states, nil
}

func (o *Oracle) SaveToFile(filename string) error {
	// Write to a temporary file first for atomicity
	tempFile := filename + ".tmp"
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(o.stateCache); err != nil {
		os.Remove(tempFile) // Clean up on error
		return fmt.Errorf("failed to encode state: %w", err)
	}

	if err := file.Sync(); err != nil {
		os.Remove(tempFile) // Clean up on error
		return fmt.Errorf("failed to sync temp file: %w", err)
	}

	file.Close() // Close before rename

	// Atomically replace the original file
	if err := os.Rename(tempFile, filename); err != nil {
		os.Remove(tempFile) // Clean up on error
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

func (o *Oracle) CacheState() {
	currentState := o.currentState.Load().(*sidecartypes.OracleState)
	newState := *currentState // Create a copy of the current state

	// o.currentState is already updated by processUpdates before CacheState is called.
	// The line o.currentState.Store(&newState) was redundant here.

	// Cache the new state
	o.stateCache = append(o.stateCache, newState)
	if len(o.stateCache) > sidecartypes.OracleCacheSize {
		o.stateCache = o.stateCache[1:]
	}

	if err := o.SaveToFile(o.Config.StateFile); err != nil {
		log.Printf("Error saving state to file: %v", err)
	}
}

func (o *Oracle) getStateByEthHeight(height uint64) (*sidecartypes.OracleState, error) {
	// Search in reverse order to efficiently find the most recent state with matching height
	for i := len(o.stateCache) - 1; i >= 0; i-- {
		if o.stateCache[i].EthBlockHeight == height {
			return &o.stateCache[i], nil
		}
	}
	return nil, fmt.Errorf("state with Ethereum block height %d not found", height)
}

func LoadConfig() sidecartypes.Config {
	configFile := getConfigFile()
	cfg, err := readConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
	return cfg
}

func getConfigFile() string {
	configFile := os.Getenv("SIDECAR_CONFIG_FILE")
	if configFile == "" {
		configFile = "config.yaml"
	}
	return configFile
}

func readConfig(configFile string) (sidecartypes.Config, error) {
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return sidecartypes.Config{}, fmt.Errorf("unable to read config from %s: %v", configFile, err)
	}

	rootConfig := sidecartypes.Config{}
	if err = yaml.Unmarshal(yamlFile, &rootConfig); err != nil {
		return sidecartypes.Config{}, fmt.Errorf("error unmarshalling config from %s: %v", configFile, err)
	}

	return rootConfig, nil
}

func (o *Oracle) GetSidecarState() *sidecartypes.OracleState {
	return o.currentState.Load().(*sidecartypes.OracleState)
}

func (o *Oracle) GetZrChainQueryClient() *client.QueryClient {
	return o.zrChainQueryClient
}

// SetStateCacheForTesting allows setting the oracle's state cache and current state for testing.
// If states is not empty, the last state in the slice becomes the current state.
// If states is empty or nil, it initializes with an empty state.
func (o *Oracle) SetStateCacheForTesting(states []sidecartypes.OracleState) {
	if len(states) > 0 {
		o.stateCache = make([]sidecartypes.OracleState, len(states))
		copy(o.stateCache, states)
		o.currentState.Store(&o.stateCache[len(o.stateCache)-1])
	} else {
		o.stateCache = []sidecartypes.OracleState{EmptyOracleState}
		o.currentState.Store(&EmptyOracleState)
	}
}

func (o *Oracle) initializeStateUpdate() *oracleStateUpdate {
	return &oracleStateUpdate{
		latestSolanaSigs: make(map[sidecartypes.SolanaEventType]solana.Signature),
		SolanaMintEvents: []api.SolanaMintEvent{},
		solanaBurnEvents: []api.BurnEvent{},
		eigenDelegations: make(map[string]map[string]*big.Int),
		redemptions:      []api.Redemption{},
		ethBurnEvents:    []api.BurnEvent{},
	}
}

// resetStateForVersion ensures the state cache is wiped exactly once after upgrading to a
// brand-new SidecarVersionName. It keeps a companion meta file (stateFile + ".meta") that
// stores the last version the cache was written with. If the meta file is missing or the
// version differs from the current one, the function deletes the cache file, writes the
// updated meta, and returns true (indicating first boot for this version). Subsequent boots
// for the same version leave the cache intact and return false.
func resetStateForVersion(stateFile string) bool {
	currentVersion := sidecartypes.SidecarVersionName
	metaFile := stateFile + ".meta"

	type meta struct {
		Version string `json:"version"`
	}

	// Check if current version requires cache reset
	requiresReset := slices.Contains(sidecartypes.VersionsRequiringCacheReset, currentVersion)

	if !requiresReset {
		// Current version doesn't require cache reset, just update meta file if needed
		if f, err := os.Open(metaFile); err == nil {
			defer f.Close()
			var m meta
			if err := json.NewDecoder(f).Decode(&m); err == nil && m.Version == currentVersion {
				// Meta already matches current version
				return false
			}
		}

		// Update meta file to current version without resetting cache
		if f, err := os.Create(metaFile); err == nil {
			json.NewEncoder(f).Encode(meta{Version: currentVersion})
			f.Close()
		}
		return false
	}

	// Attempt to read existing meta file
	if f, err := os.Open(metaFile); err == nil {
		defer f.Close()
		var m meta
		if err := json.NewDecoder(f).Decode(&m); err == nil && m.Version == currentVersion {
			// Cache already corresponds to current version – no reset needed.
			slog.Info("Cache is already aligned", "version", currentVersion)
			return false
		}
	}

	// Either meta file missing or version mismatch → first boot for this version.
	slog.Info("First boot detected for sidecar version requiring cache reset", "version", currentVersion)

	// Remove state file if it exists.
	if err := os.Remove(stateFile); err != nil && !os.IsNotExist(err) {
		slog.Error("Failed to delete cache file during version reset", "file", stateFile, "error", err)
	}

	// Write new meta file with current version.
	if f, err := os.Create(metaFile); err != nil {
		slog.Error("Failed to create cache meta file", "file", metaFile, "error", err)
	} else {
		if err := json.NewEncoder(f).Encode(meta{Version: currentVersion}); err != nil {
			slog.Error("Failed to write cache meta file", "file", metaFile, "error", err)
		}
		f.Close()
	}

	return true
}

var colorMap = map[string]func(string) string{
	// Core categories
	"error": func(s string) string { return color.HEX("F07178").Sprint(s) }, // Red

	// Info & identifiers
	"version": func(s string) string { return color.HEX("59C2FF").Sprint(s) }, // Cyan
	"network": func(s string) string { return color.HEX("59C2FF").Sprint(s) }, // Cyan
	"chain":   func(s string) string { return color.HEX("59C2FF").Sprint(s) }, // Cyan
	"chainID": func(s string) string { return color.HEX("59C2FF").Sprint(s) }, // Cyan
	"state":   func(s string) string { return color.HEX("59C2FF").Sprint(s) }, // Cyan

	// Timing
	"time":             func(s string) string { return color.HEX("FFD580").Sprint(s) }, // Pale Yellow
	"interval":         func(s string) string { return color.HEX("FF8F40").Sprint(s) }, // Orange
	"sleepDuration":    func(s string) string { return color.HEX("FF8F40").Sprint(s) }, // Orange
	"nextIntervalMark": func(s string) string { return color.HEX("FFD580").Sprint(s) }, // Pale Yellow

	// Transactions & signatures
	"tx":                     func(s string) string { return color.HEX("B3B1AD").Sprint(s) }, // Foreground
	"txID":                   func(s string) string { return color.HEX("B3B1AD").Sprint(s) }, // Foreground
	"txSig":                  func(s string) string { return color.HEX("B3B1AD").Sprint(s) }, // Foreground
	"txHash":                 func(s string) string { return color.HEX("B3B1AD").Sprint(s) }, // Foreground
	"lastSig":                func(s string) string { return color.HEX("539AFC").Sprint(s) }, // Blue
	"newestLastProcessedSig": func(s string) string { return color.HEX("539AFC").Sprint(s) }, // Blue
	"sigHash":                func(s string) string { return color.HEX("539AFC").Sprint(s) }, // Blue

	// Addresses
	"address":         func(s string) string { return color.HEX("B3B1AD").Sprint(s) }, // Foreground
	"destinationAddr": func(s string) string { return color.HEX("B3B1AD").Sprint(s) }, // Foreground
	"recipient":       func(s string) string { return color.HEX("B3B1AD").Sprint(s) }, // Foreground

	// Events
	"eventType":  func(s string) string { return color.HEX("95E6CB").Sprint(s) }, // Green
	"eventName":  func(s string) string { return color.HEX("95E6CB").Sprint(s) }, // Green
	"eventIndex": func(s string) string { return color.HEX("95E6CB").Sprint(s) }, // Green

	// Values / amounts
	"amount":   func(s string) string { return color.HEX("E6B450").Sprint(s) }, // Yellow
	"value":    func(s string) string { return color.HEX("E6B450").Sprint(s) }, // Yellow
	"fee":      func(s string) string { return color.HEX("F07178").Sprint(s) }, // Red
	"block":    func(s string) string { return color.HEX("FFFFFF").Sprint(s) }, // White
	"ROCK/USD": func(s string) string { return color.HEX("95E6CB").Sprint(s) }, // Green
	"BTC/USD":  func(s string) string { return color.HEX("FF8F40").Sprint(s) }, // Orange
	"ETH/USD":  func(s string) string { return color.HEX("539AFC").Sprint(s) }, // Blue

	// Metrics
	"count":        func(s string) string { return color.HEX("D2A6FF").Sprint(s) }, // Magenta
	"batchSize":    func(s string) string { return color.HEX("D2A6FF").Sprint(s) }, // Magenta
	"total":        func(s string) string { return color.HEX("D2A6FF").Sprint(s) }, // Magenta
	"inspected":    func(s string) string { return color.HEX("D2A6FF").Sprint(s) }, // Magenta
	"newTxCount":   func(s string) string { return color.HEX("D2A6FF").Sprint(s) }, // Magenta
	"requestIndex": func(s string) string { return color.HEX("D2A6FF").Sprint(s) }, // Magenta

	// Burn / Mint signatures
	"rockMintSig":   func(s string) string { return color.HEX("95E6CB").Sprint(s) }, // Green
	"zenBTCMintSig": func(s string) string { return color.HEX("95E6CB").Sprint(s) }, // Green
	"rockBurnSig":   func(s string) string { return color.HEX("E6B450").Sprint(s) }, // Yellow
	"zenBTCBurnSig": func(s string) string { return color.HEX("E6B450").Sprint(s) }, // Yellow
}

// initLogger sets up coloured structured logging
func initLogger(debug bool) {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:      level,
		TimeFormat: time.DateTime,
		AddSource:  debug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Handle timestamp coloring
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(color.HEX("FFFACD").Sprint(a.Value.String()))
				return a
			}
			// Handle other custom fields
			if f, ok := colorMap[a.Key]; ok {
				a.Value = slog.StringValue(f(a.Value.String()))
			}
			return a
		},
	})))
}

func connectWithRetry(rpcAddress string, maxRetries int, delay time.Duration) (*ethclient.Client, error) {
	var client *ethclient.Client
	var err error

	for i := 0; i < maxRetries || maxRetries == 0; i++ {
		client, err = ethclient.Dial(rpcAddress)
		if err == nil {
			// Check if client can respond to eth_blockNumber
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			_, err = client.BlockNumber(ctx)
			if err == nil {
				slog.Info("Successfully connected to Ethereum client", "rpc", rpcAddress)
				return client, err
			}
			client.Close()
		}

		slog.Warn("Retrying connection to Ethereum client", "attempt", i+1, "error", err)
		time.Sleep(delay)
	}

	return nil, err
}

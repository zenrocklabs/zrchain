package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"math/big"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"cosmossdk.io/math"
	"github.com/Zenrock-Foundation/zrchain/v6/go-client"
	"github.com/Zenrock-Foundation/zrchain/v6/sidecar/proto/api"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/shared"
	"github.com/ethereum/go-ethereum/ethclient"
	solana "github.com/gagliardetto/solana-go"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
	jsonrpc "github.com/gagliardetto/solana-go/rpc/jsonrpc"
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

var fileIOMutex sync.Mutex

// saveStatesToFile writes the oracle states to a file.
func saveStatesToFile(filename string, states []sidecartypes.OracleState) error {
	// Serialize file I/O operations to prevent concurrent writes to the same file
	fileIOMutex.Lock()
	defer fileIOMutex.Unlock()

	// Write to a temporary file first for atomicity
	tempFile := filename + ".tmp"
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(states); err != nil {
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

// appendStateToCache appends the current state to the state cache and saves it to the file.
func (o *Oracle) appendStateToCache() {
	currentState := o.currentState.Load().(*sidecartypes.OracleState)
	newState := *currentState // Create a copy of the current state

	// o.currentState is already updated by processUpdates before appendStateToCache is called.
	// The line o.currentState.Store(&newState) was redundant here.

	// Acquire exclusive lock to safely modify stateCache (write operation)
	o.stateCacheMutex.Lock()
	o.stateCache = append(o.stateCache, newState)
	if len(o.stateCache) > sidecartypes.OracleCacheSize {
		o.stateCache = o.stateCache[1:]
	}
	// Make a copy while holding the lock to pass to SaveToFile
	stateCacheCopy := slices.Clone(o.stateCache)
	o.stateCacheMutex.Unlock()

	if err := saveStatesToFile(o.Config.StateFile, stateCacheCopy); err != nil {
		log.Printf("Error saving state to file: %v", err)
	}
}

func (o *Oracle) getStateByEthHeight(height uint64) (*sidecartypes.OracleState, error) {
	// Acquire read lock to safely read stateCache (allows concurrent reads)
	o.stateCacheMutex.RLock()
	defer o.stateCacheMutex.RUnlock()

	// Search in reverse order to efficiently find the most recent state with matching height
	for i := len(o.stateCache) - 1; i >= 0; i-- {
		if o.stateCache[i].EthBlockHeight == height {
			return &o.stateCache[i], nil
		}
	}
	return nil, fmt.Errorf("state with Ethereum block height %d not found", height)
}

func LoadConfig(configFileFlag, configDirFlag string) sidecartypes.Config {
	configFile := getConfigFile(configFileFlag, configDirFlag)
	cfg, err := readConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
	return cfg
}

func getConfigFile(configFileFlag, configDirFlag string) string {
	// 1. Command-line flag has the highest precedence.
	if configFileFlag != "" {
		slog.Info("Using config file specified by flag", "path", configFileFlag)
		return configFileFlag
	}

	// 2. Environment variable is next.
	if envConfigFile := os.Getenv("SIDECAR_CONFIG_FILE"); envConfigFile != "" {
		slog.Info("Using config file specified by SIDECAR_CONFIG_FILE environment variable", "path", envConfigFile)
		return envConfigFile
	}

	// 3. --config-dir flag.
	if configDirFlag != "" {
		configPath := filepath.Join(configDirFlag, "config.yaml")
		slog.Info("Using config file from directory specified by --config-dir flag", "path", configPath)
		return configPath
	}

	// 4. Autodetection based on executable path.
	exePath, err := os.Executable()
	if err != nil {
		slog.Warn("Could not determine executable path; falling back to relative path for config", "error", err)
		return "config.yaml" // Fallback
	}

	exeDir := filepath.Dir(exePath)

	// Check for config.yaml in the same directory as the executable (local setup).
	localConfigPath := filepath.Join(exeDir, "config.yaml")
	if _, err := os.Stat(localConfigPath); err == nil {
		slog.Info("Found config file in the executable's directory", "path", localConfigPath)
		return localConfigPath
	}

	// Check for config.yaml in the parent directory (container setup).
	parentDir := filepath.Dir(exeDir)
	parentConfigPath := filepath.Join(parentDir, "config.yaml")
	if _, err := os.Stat(parentConfigPath); err == nil {
		slog.Info("Found config file in the parent of the executable's directory", "path", parentConfigPath)
		return parentConfigPath
	}

	// 5. Fallback if not found in any of the preferred locations.
	slog.Warn("Could not find config.yaml automatically; falling back to default relative path", "path", "config.yaml")
	return "config.yaml"
}

func readConfig(configFile string) (sidecartypes.Config, error) {
	slog.Info("Reading configuration file", "path", configFile)
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
	// Acquire exclusive lock to safely modify stateCache (write operation)
	o.stateCacheMutex.Lock()
	defer o.stateCacheMutex.Unlock()

	if len(states) > 0 {
		o.stateCache = make([]sidecartypes.OracleState, len(states))
		copy(o.stateCache, states)
		o.currentState.Store(&o.stateCache[len(o.stateCache)-1])
	} else {
		o.stateCache = []sidecartypes.OracleState{EmptyOracleState}
		o.currentState.Store(&EmptyOracleState)
	}
}

// Helper function to parse RPC response ID into request index
func parseRPCResponseID(resp *jsonrpc.RPCResponse, eventType string) (int, bool) {
	if resp == nil {
		slog.Error("Nil RPCResponse object in sub-batch response", "eventType", eventType)
		return 0, false
	}

	var requestIndex int
	switch id := resp.ID.(type) {
	case float64: // JSON numbers often decode to float64
		requestIndex = int(id)
	case int:
		requestIndex = id
	case uint64: // Match the type we put in the request
		requestIndex = int(id)
	case json.Number:
		idInt64, err := id.Int64()
		if err != nil {
			slog.Error("Failed to convert json.Number ID to int64", "eventType", eventType, "error", err)
			return 0, false
		}
		requestIndex = int(idInt64)
	default:
		slog.Error("Invalid response ID type received", "idType", fmt.Sprintf("%T", resp.ID), "eventType", eventType)
		return 0, false
	}
	return requestIndex, true
}

// Helper function to validate request index bounds
func validateRequestIndex(requestIndex int, batchSize int, eventType string) bool {
	if requestIndex < 0 || requestIndex >= batchSize {
		slog.Error("Invalid response ID received for sub-batch", "requestIndex", requestIndex, "eventType", eventType)
		return false
	}
	return true
}

// Helper function to generate event keys for deduplication
func generateBurnEventKey(event api.BurnEvent) string {
	return fmt.Sprintf("%s-%s-%d", event.ChainID, event.TxID, event.LogIndex)
}

func generateMintEventKey(event api.SolanaMintEvent) string {
	return base64.StdEncoding.EncodeToString(event.SigHash)
}

// Helper function to check if events already exist and merge new ones
func mergeNewBurnEvents(existingEvents []api.BurnEvent, cleanedEvents map[string]bool, newEvents []api.BurnEvent, eventTypeName string) []api.BurnEvent {
	// Create a map of existing events for quick lookup
	existingEventKeys := make(map[string]bool)
	for _, event := range existingEvents {
		key := generateBurnEventKey(event)
		existingEventKeys[key] = true
	}

	// Also check against already cleaned events
	for key := range cleanedEvents {
		existingEventKeys[key] = true
	}

	// Start with existing events
	mergedEvents := make([]api.BurnEvent, len(existingEvents))
	copy(mergedEvents, existingEvents)

	// Add new events if they don't already exist
	addedCount := 0
	skippedCount := 0
	for _, event := range newEvents {
		key := generateBurnEventKey(event)
		if !existingEventKeys[key] {
			mergedEvents = append(mergedEvents, event)
			addedCount++
			slog.Debug("Added burn event to state", "type", eventTypeName, "txID", event.TxID)
		} else {
			skippedCount++
			slog.Debug("Skipping already present burn event", "type", eventTypeName, "txID", event.TxID)
		}
	}

	// Only log if there were actual changes
	if addedCount > 0 || skippedCount > 0 {
		slog.Info("Burn event merge summary",
			"type", eventTypeName,
			"added", addedCount,
			"skipped", skippedCount,
			"totalAfterMerge", len(mergedEvents))
	}

	return mergedEvents
}

// Helper function to merge new mint events
func mergeNewMintEvents(existingEvents []api.SolanaMintEvent, cleanedEvents map[string]bool, newEvents []api.SolanaMintEvent, eventTypeName string) []api.SolanaMintEvent {
	// Create a map of existing events for quick lookup
	existingEventKeys := make(map[string]bool)
	for _, event := range existingEvents {
		key := generateMintEventKey(event)
		existingEventKeys[key] = true
	}

	// Also check against already cleaned events
	for key := range cleanedEvents {
		existingEventKeys[key] = true
	}

	// Start with existing events
	mergedEvents := make([]api.SolanaMintEvent, len(existingEvents))
	copy(mergedEvents, existingEvents)

	// Add new events if they don't already exist
	addedCount := 0
	skippedCount := 0
	for _, event := range newEvents {
		key := generateMintEventKey(event)
		if !existingEventKeys[key] {
			mergedEvents = append(mergedEvents, event)
			addedCount++
			slog.Debug("Added mint event to state", "type", eventTypeName, "txSig", event.TxSig, "key", key[:16]+"...")
		} else {
			skippedCount++
			slog.Debug("Skipping already present mint event", "type", eventTypeName, "txSig", event.TxSig, "key", key[:16]+"...")
		}
	}

	// Only log if there were actual changes or new events to process
	if addedCount > 0 || skippedCount > 0 || len(newEvents) > 0 {
		// Add detailed logging for debugging - but only when there's activity
		if addedCount > 0 || len(newEvents) > 0 {
			slog.Info("Mint event merge debug info",
				"type", eventTypeName,
				"existingCount", len(existingEvents),
				"cleanedCount", len(cleanedEvents),
				"newCount", len(newEvents),
				"existingKeys", len(existingEventKeys))
		}

		slog.Info("Mint event merge summary",
			"type", eventTypeName,
			"added", addedCount,
			"skipped", skippedCount,
			"totalAfterMerge", len(mergedEvents))
	}

	return mergedEvents
}

func (o *Oracle) initializeStateUpdate() *oracleStateUpdate {
	// Copy pending transactions from current state to preserve them across ticks
	currentState := o.currentState.Load().(*sidecartypes.OracleState)
	pendingTransactions := make(map[string]sidecartypes.PendingTxInfo)
	if currentState.PendingSolanaTxs != nil {
		for k, v := range currentState.PendingSolanaTxs {
			pendingTransactions[k] = v
		}
	}

	return &oracleStateUpdate{
		suggestedTip:            big.NewInt(0),
		estimatedGas:            0,
		ROCKUSDPrice:            math.LegacyZeroDec(),
		BTCUSDPrice:             math.LegacyZeroDec(),
		ETHUSDPrice:             math.LegacyZeroDec(),
		ZECUSDPrice:             math.LegacyZeroDec(),
		ethBurnEvents:           make([]api.BurnEvent, 0),
		cleanedEthBurnEvents:    make(map[string]bool),
		solanaBurnEvents:        make([]api.BurnEvent, 0),
		cleanedSolanaBurnEvents: make(map[string]bool),
		redemptions:             make([]api.Redemption, 0),
		SolanaMintEvents:        make([]api.SolanaMintEvent, 0),
		cleanedSolanaMintEvents: make(map[string]bool),
		latestSolanaSigs:        make(map[sidecartypes.SolanaEventType]solana.Signature),
		latestEventStoreCursors: make(map[sidecartypes.SolanaEventType]string),
		pendingTransactions:     pendingTransactions,
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

	// Special message highlighting
	"State event counts for this tick": func(s string) string { return color.HEX("FFD700").Sprint(s) }, // Bright saturated yellow
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
				a.Value = slog.StringValue(color.HEX("808080").Sprint(a.Value.String()))
				return a
			}
			// Handle message-based coloring
			if a.Key == slog.MessageKey {
				if f, ok := colorMap[a.Value.String()]; ok {
					a.Value = slog.StringValue(f(a.Value.String()))
				}
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

func connectWithRetry[T any](
	serviceName string,
	rpcAddress string,
	maxRetries int,
	delay time.Duration,
	createClient func(string) (T, error),
	healthCheck func(T, context.Context) error,
) (T, error) {
	var client T
	var err error

	for i := 0; i < maxRetries || maxRetries == 0; i++ {
		client, err = createClient(rpcAddress)
		if err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			if err = healthCheck(client, ctx); err == nil {
				slog.Info("Successfully connected to client", "service", serviceName, "rpc", rpcAddress)
				return client, nil
			}
		}

		slog.Warn("Retrying connection to client", "service", serviceName, "attempt", i+1, "error", err)
		time.Sleep(delay)
	}

	return client, err
}

func connectEthereumWithRetry(rpcAddress string, maxRetries int, delay time.Duration) (*ethclient.Client, error) {
	return connectWithRetry(
		"Ethereum",
		rpcAddress,
		maxRetries,
		delay,
		func(addr string) (*ethclient.Client, error) {
			return ethclient.Dial(addr)
		},
		func(client *ethclient.Client, ctx context.Context) error {
			_, err := client.BlockNumber(ctx)
			if err != nil {
				client.Close()
			}
			return err
		},
	)
}

func connectSolanaWithRetry(rpcAddress string, maxRetries int, delay time.Duration) (*solanarpc.Client, error) {
	return connectWithRetry(
		"Solana",
		rpcAddress,
		maxRetries,
		delay,
		func(addr string) (*solanarpc.Client, error) {
			client := solanarpc.New(addr)
			if client == nil {
				return nil, fmt.Errorf("failed to create Solana client")
			}
			return client, nil
		},
		func(client *solanarpc.Client, ctx context.Context) error {
			_, err := client.GetHealth(ctx)
			return err
		},
	)
}

// func connectZrChainWithRetry(rpcAddress string, maxRetries int, delay time.Duration) (*client.QueryClient, error) {
// 	return connectWithRetry(
// 		"zrChain",
// 		rpcAddress,
// 		maxRetries,
// 		delay,
// 		func(addr string) (*client.QueryClient, error) {
// 			client, err := client.NewQueryClient(addr, true)
// 			if err != nil {
// 				return nil, err
// 			}
// 			if client == nil {
// 				return nil, fmt.Errorf("zrChain query client is nil after creation")
// 			}
// 			return client, nil
// 		},
// 		func(client *client.QueryClient, ctx context.Context) error {
// 			_, err := client.BondedValidators(ctx, nil)
// 			return err
// 		},
// 	)
// }

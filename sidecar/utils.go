package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Zenrock-Foundation/zrchain/v6/go-client"
	sidecartypes "github.com/Zenrock-Foundation/zrchain/v6/sidecar/shared"
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
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(o.stateCache)
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

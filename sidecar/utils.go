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

func (o *Oracle) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// Initialize with empty state if file doesn't exist
			o.updateChan <- EmptyOracleState
			o.stateCache = []sidecartypes.OracleState{EmptyOracleState}
			return nil
		}
		return err
	}
	defer file.Close()

	var states []sidecartypes.OracleState
	if err := json.NewDecoder(file).Decode(&states); err != nil {
		return fmt.Errorf("failed to decode state file: %w", err)
	}

	if len(states) > 0 {
		latestState := &states[len(states)-1]
		o.updateChan <- sidecartypes.OracleState{
			EigenDelegations:           latestState.EigenDelegations,
			EthBlockHeight:             latestState.EthBlockHeight,
			EthGasLimit:                latestState.EthGasLimit,
			EthBaseFee:                 latestState.EthBaseFee,
			EthTipCap:                  latestState.EthTipCap,
			SolanaLamportsPerSignature: latestState.SolanaLamportsPerSignature,
			EthBurnEvents:              latestState.EthBurnEvents,
			CleanedEthBurnEvents:       latestState.CleanedEthBurnEvents,
			Redemptions:                latestState.Redemptions,
			ROCKUSDPrice:               latestState.ROCKUSDPrice,
			BTCUSDPrice:                latestState.BTCUSDPrice,
			ETHUSDPrice:                latestState.ETHUSDPrice,
			SolanaMintEvents:           latestState.SolanaMintEvents,
		}
		o.stateCache = states
	} else {
		// Initialize with empty state if the file is empty
		o.updateChan <- EmptyOracleState
		o.stateCache = []sidecartypes.OracleState{EmptyOracleState}
	}

	return nil
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

	// Update the ID and store the new state
	o.currentState.Store(&newState)

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

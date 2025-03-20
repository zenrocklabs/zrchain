package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v2"
)

func main() {
	config := parseFlags()

	// Display selected network info
	fmt.Printf("Using %s network (%s)\n", config.Network, config.RPCNode)

	// Get validator information
	fmt.Println("Fetching validator information...")
	addrToMoniker, err := buildValidatorMappings(config.RPCNode)
	if err != nil {
		fmt.Printf("Warning: Failed to get validator information: %v\n", err)
		fmt.Println("Proceeding without validator names.")
		addrToMoniker = make(map[string]string)
	}

	// Get block data
	blockData, err := getBlockData(config)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Process block data
	processBlockData(blockData, addrToMoniker, config.MissingOnly)
}

// parseFlags parses command line flags and returns a Config
func parseFlags() Config {
	useFileFlag := flag.String("file", "", "Use file instead of executing command (optional)")
	networkFlag := flag.String("network", "mainnet", "Network to use: devnet, testnet, or mainnet (default: mainnet)")
	rpcNodeFlag := flag.String("node", "", "RPC node URL (overrides network selection)")
	blockHeightFlag := flag.String("height", "", "Block height (default: latest)")
	missingOnlyFlag := flag.Bool("missing-only", false, "Only show validators missing vote extensions")
	flag.Parse()

	// Map network to RPC URL if node is not explicitly provided
	rpcURL := *rpcNodeFlag
	if rpcURL == "" {
		switch *networkFlag {
		case "localnet", "local", "localhost":
			rpcURL = "http://localhost:26657"
		case "devnet", "dev", "amber":
			rpcURL = "https://rpc.dev.zenrock.tech:443"
		case "testnet", "test", "gardia":
			rpcURL = "https://rpc.gardia.zenrocklabs.io:443"
		case "mainnet", "main", "diamond":
			rpcURL = "https://rpc.diamond.zenrocklabs.io:443"
		default:
			// Default to mainnet if unrecognized network
			fmt.Printf("Warning: unrecognized network '%s', defaulting to mainnet\n", *networkFlag)
			*networkFlag = "mainnet"
			rpcURL = "https://rpc.diamond.zenrocklabs.io:443"
		}
	}

	return Config{
		UseFile:     *useFileFlag,
		RPCNode:     rpcURL,
		Network:     *networkFlag,
		BlockHeight: *blockHeightFlag,
		MissingOnly: *missingOnlyFlag,
	}
}

// getBlockData retrieves block data from file or by executing command
func getBlockData(config Config) (*BlockData, error) {
	var input []byte
	var err error

	if config.UseFile != "" {
		input, err = os.ReadFile(config.UseFile)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %v", err)
		}
	} else {
		fmt.Println("Querying latest block data...")
		input, err = executeZenrockdCommand(config.RPCNode, config.BlockHeight)
		if err != nil {
			return nil, fmt.Errorf("error executing command: %v", err)
		}
	}

	// Extract the JSON object from the input
	re := regexp.MustCompile(`\{[\s\S]*\}`)
	match := re.FindString(string(input))
	if match == "" {
		return nil, fmt.Errorf("could not find a JSON object in the input")
	}

	// Parse the JSON data
	var blockData BlockData
	if err := json.Unmarshal([]byte(match), &blockData); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return &blockData, nil
}

// processBlockData analyzes the block data and displays results
func processBlockData(blockData *BlockData, addrToMoniker map[string]string, missingOnly bool) {
	if len(blockData.ConsensusData.Votes) == 0 {
		fmt.Println("No votes found in block data")
		return
	}

	// Calculate statistics
	validators, allExtensions := processVotes(blockData.ConsensusData.Votes, addrToMoniker, missingOnly)
	stats := calculateStats(validators)

	// Display participation stats
	printParticipationStats(stats)

	// Handle missing validators if requested
	if missingOnly {
		printMissingValidators(validators, stats.TotalVotingPower)
		return
	}

	// Compare extensions to find differences
	if len(allExtensions) > 0 {
		fmt.Println("\n\n=== DIFFERENCES BETWEEN VOTE EXTENSIONS ===")
		findDifferences(allExtensions, validators, stats.TotalVotingPower)
	}
}

// processVotes processes all votes and returns validator information
func processVotes(votes []ConsensusVote, addrToMoniker map[string]string, missingOnly bool) ([]ValidatorInfo, []map[string]any) {
	validators := make([]ValidatorInfo, 0, len(votes))
	allExtensions := make([]map[string]any, 0, len(votes))
	validatorsWithExtensions := 0
	totalVotingPower := 0
	totalVotedPower := 0

	if !missingOnly {
		fmt.Printf("Found %d validators in the block\n", len(votes))
	}

	for i, vote := range votes {
		// Decode the validator address
		decodedAddr := decodeValidatorAddress(vote.Validator.Address)

		// Get moniker for this validator
		moniker := addrToMoniker[decodedAddr]
		if moniker == "" {
			moniker = "Unknown"
		}

		validatorInfo := ValidatorInfo{
			Address:     vote.Validator.Address,
			Power:       vote.Validator.Power,
			Index:       i + 1,
			HasVote:     false,
			DecodedAddr: decodedAddr,
			Moniker:     moniker,
		}

		totalVotingPower += vote.Validator.Power

		// Print validator info if not in missing-only mode
		if !missingOnly {
			fmt.Printf("\n=== Validator %s (%s) (Power: %d) ===\n", decodedAddr, moniker, vote.Validator.Power)
		}

		// Process vote extension
		if vote.VoteExtension == "" {
			if !missingOnly {
				fmt.Println("No vote extension found")
			}
			allExtensions = append(allExtensions, nil)
		} else {
			validatorInfo.HasVote = true
			validatorsWithExtensions++
			totalVotedPower += vote.Validator.Power

			extension := processVoteExtension(vote.VoteExtension, missingOnly)
			allExtensions = append(allExtensions, extension)
		}

		validators = append(validators, validatorInfo)
	}

	return validators, allExtensions
}

// processVoteExtension decodes and processes a vote extension
func processVoteExtension(extensionBase64 string, missingOnly bool) map[string]any {
	// Decode base64 data
	decodedBytes, err := base64.StdEncoding.DecodeString(extensionBase64)
	if err != nil {
		if !missingOnly {
			fmt.Printf("Error decoding base64: %v\n", err)
		}
		return nil
	}

	// Parse JSON
	var extension map[string]any
	if err := json.Unmarshal(decodedBytes, &extension); err != nil {
		if !missingOnly {
			fmt.Printf("Error parsing decoded JSON: %v\n", err)
			fmt.Println("Raw decoded data:", string(decodedBytes))
		}
		return nil
	}

	// Print pretty JSON if not in missing-only mode
	if !missingOnly {
		prettyJSONBytes, err := json.MarshalIndent(extension, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting JSON: %v\n", err)
		} else {
			fmt.Println("Decoded vote extension:")
			fmt.Println(string(prettyJSONBytes))
		}
	}

	return extension
}

// calculateStats calculates statistics from validator information
func calculateStats(validators []ValidatorInfo) Stats {
	validatorsWithExtensions := 0
	totalVotedPower := 0
	totalVotingPower := 0

	for _, v := range validators {
		totalVotingPower += v.Power
		if v.HasVote {
			validatorsWithExtensions++
			totalVotedPower += v.Power
		}
	}

	return Stats{
		ValidatorsWithExtensions: validatorsWithExtensions,
		TotalValidators:          len(validators),
		TotalVotedPower:          totalVotedPower,
		TotalVotingPower:         totalVotingPower,
	}
}

// printParticipationStats prints statistics about vote extension participation
func printParticipationStats(stats Stats) {
	fmt.Printf("\n=== VOTE EXTENSION PARTICIPATION ===\n")
	fmt.Printf("Validators with vote extensions: %d of %d (%.2f%%)\n",
		stats.ValidatorsWithExtensions,
		stats.TotalValidators,
		float64(stats.ValidatorsWithExtensions)/float64(stats.TotalValidators)*100)
	fmt.Printf("Voting power with extensions: %d of %d (%.2f%%)\n",
		stats.TotalVotedPower,
		stats.TotalVotingPower,
		float64(stats.TotalVotedPower)/float64(stats.TotalVotingPower)*100)
}

// printMissingValidators prints information about validators missing vote extensions
func printMissingValidators(validators []ValidatorInfo, totalVotingPower int) {
	missingValidators := []ValidatorInfo{}
	for _, v := range validators {
		if !v.HasVote {
			missingValidators = append(missingValidators, v)
		}
	}

	fmt.Printf("\n=== VALIDATORS WITHOUT VOTE EXTENSIONS (%d validators) ===\n", len(missingValidators))
	if len(missingValidators) == 0 {
		fmt.Println("All validators submitted vote extensions!")
		return
	}

	// Sort missing validators by power (high to low)
	sort.Slice(missingValidators, func(i, j int) bool {
		return missingValidators[i].Power > missingValidators[j].Power
	})

	totalMissingPower := 0
	for _, v := range missingValidators {
		totalMissingPower += v.Power
	}

	fmt.Printf("Total missing power: %d (%.2f%% of total)\n\n",
		totalMissingPower,
		float64(totalMissingPower)/float64(totalVotingPower)*100)

	for i, v := range missingValidators {
		fmt.Printf("%3d. %s (%s) (Power: %d)\n", i+1, v.DecodedAddr, v.Moniker, v.Power)
	}
}

// buildValidatorMappings builds a mapping from bech32 consensus address to moniker
func buildValidatorMappings(rpcNode string) (map[string]string, error) {
	addressToMoniker := make(map[string]string)

	// Get validator set data
	valSetOutput, err := execCommand("bash", "-c",
		fmt.Sprintf("zenrockd --node=%s q consensus comet validator-set", rpcNode))
	if err != nil {
		return nil, fmt.Errorf("failed to get validator set: %v", err)
	}

	// Get validators with monikers
	validatorsOutput, err := execCommand("bash", "-c",
		fmt.Sprintf("zenrockd --node=%s q validation validators", rpcNode))
	if err != nil {
		return nil, fmt.Errorf("failed to get validators: %v", err)
	}

	// Parse responses
	var valSetResp ValidatorSetResponse
	if err := yaml.Unmarshal(valSetOutput, &valSetResp); err != nil {
		return nil, fmt.Errorf("failed to parse validator set: %v", err)
	}

	var validatorsResp ValidatorsResponse
	if err := yaml.Unmarshal(validatorsOutput, &validatorsResp); err != nil {
		return nil, fmt.Errorf("failed to parse validators: %v", err)
	}

	// Build pubkey to moniker mapping
	pubkeyToMoniker := make(map[string]string)
	for _, val := range validatorsResp.Validators {
		pubkeyToMoniker[val.ConsensusPublicKey.Value] = val.Description.Moniker
	}

	// Build address to moniker mapping
	for _, val := range valSetResp.Validators {
		if moniker, ok := pubkeyToMoniker[val.PubKey.Value]; ok {
			addressToMoniker[val.Address] = moniker
		}
	}

	fmt.Printf("Found %d validators with monikers\n", len(addressToMoniker))
	return addressToMoniker, nil
}

// execCommand is a helper to execute commands and capture output
func execCommand(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("command execution failed: %s\nStderr: %s", err, stderr.String())
	}

	return stdout.Bytes(), nil
}

// executeZenrockdCommand runs the zenrockd command and returns the output
func executeZenrockdCommand(rpcNode, blockHeight string) ([]byte, error) {
	var fullCmd string

	if blockHeight == "" {
		fullCmd = fmt.Sprintf(`zenrockd --node=%s query block --type=height $(zenrockd --node=%s status | jq -r '.sync_info.latest_block_height') --output=json | jq -r '.data.txs[0]' | base64 --decode`, rpcNode, rpcNode)
	} else {
		fullCmd = fmt.Sprintf(`zenrockd --node=%s query block --type=height %s --output=json | jq -r '.data.txs[0]' | base64 --decode`, rpcNode, blockHeight)
	}

	return execCommand("bash", "-c", fullCmd)
}

// decodeValidatorAddress decodes validator address from base64
func decodeValidatorAddress(base64Addr string) string {
	decoded, err := base64.StdEncoding.DecodeString(base64Addr)
	if err != nil {
		return fmt.Sprintf("Invalid(%s)", base64Addr)
	}

	return sdk.MustBech32ifyAddressBytes("zenvalcons", decoded)
}

// findDifferences finds and displays differences between vote extensions
func findDifferences(extensions []map[string]any, validators []ValidatorInfo, totalPower int) {
	// Map to track differences by field
	differences := make(map[string]map[string][]ValidatorInfo)

	// Get all possible keys from all extensions
	allKeys := getAllKeys(extensions)

	// Check each key for differences
	for key := range allKeys {
		valueMap := make(map[string][]ValidatorInfo)

		for i, ext := range extensions {
			// Skip if index is out of range
			if i >= len(validators) {
				continue
			}

			valueStr := getValueString(ext, key)
			valueMap[valueStr] = append(valueMap[valueStr], validators[i])
		}

		// If there's more than one unique value, we have a difference
		if len(valueMap) > 1 {
			differences[key] = valueMap
		}
	}

	printDifferences(differences, validators, totalPower)
}

// getAllKeys gets all keys from all extensions
func getAllKeys(extensions []map[string]any) map[string]bool {
	allKeys := make(map[string]bool)
	for _, ext := range extensions {
		if ext == nil {
			continue
		}
		for k := range ext {
			allKeys[k] = true
		}
	}
	return allKeys
}

// getValueString converts an extension value to a string representation
func getValueString(ext map[string]any, key string) string {
	if ext == nil {
		return noVoteMsg
	}

	if val, ok := ext[key]; ok {
		bytes, err := json.Marshal(val)
		if err == nil {
			return string(bytes)
		}
		return fmt.Sprintf("%v", val)
	}

	return missingFieldMsg
}

// printDifferences displays the differences between vote extensions
func printDifferences(differences map[string]map[string][]ValidatorInfo, validators []ValidatorInfo, totalPower int) {
	if len(differences) == 0 {
		fmt.Println("No differences found between vote extensions")
		return
	}

	// Create a slice of keys with their consensus percentages
	type fieldConsensus struct {
		key              string
		consensusPercent float64
	}

	fieldsByConsensus := make([]fieldConsensus, 0, len(differences))

	for key, valueMap := range differences {
		// Find highest consensus for this field
		highestConsensus := 0.0
		for _, vals := range valueMap {
			thisPower := 0
			for _, val := range vals {
				thisPower += val.Power
			}
			consensusPercent := float64(thisPower) / float64(totalPower) * 100
			if consensusPercent > highestConsensus {
				highestConsensus = consensusPercent
			}
		}

		fieldsByConsensus = append(fieldsByConsensus, fieldConsensus{
			key:              key,
			consensusPercent: highestConsensus,
		})
	}

	// Sort by consensus percentage (ascending - lowest consensus last)
	sort.Slice(fieldsByConsensus, func(i, j int) bool {
		return fieldsByConsensus[i].consensusPercent > fieldsByConsensus[j].consensusPercent
	})

	for _, field := range fieldsByConsensus {
		fmt.Printf("\nDifferences in field: %s (highest consensus: %.2f%%)\n", field.key, field.consensusPercent)
		fmt.Printf("-------------------------\n")

		valueMap := differences[field.key]
		printValueDifferences(valueMap, validators, totalPower)
	}
}

// printValueDifferences prints the differences for a specific field
func printValueDifferences(valueMap map[string][]ValidatorInfo, allValidators []ValidatorInfo, totalPower int) {
	// Convert map to slice for sorting
	valueCounts := make([]ValueCount, 0, len(valueMap))
	for v, validators := range valueMap {
		thisPower := 0
		for _, val := range validators {
			thisPower += val.Power
		}

		valueCounts = append(valueCounts, ValueCount{
			Value:      v,
			Validators: validators,
			TotalPower: thisPower,
		})
	}

	// Sort by power (descending)
	sort.Slice(valueCounts, func(i, j int) bool {
		return valueCounts[i].TotalPower > valueCounts[j].TotalPower
	})

	// Format for display
	for i, vc := range valueCounts {
		displayValue := vc.Value
		if displayValue == noVoteMsg {
			displayValue = noVoteMsg
		} else if len(displayValue) > 100 {
			displayValue = displayValue[:97] + "..."
		}

		powerPercentage := float64(vc.TotalPower) / float64(totalPower) * 100

		// Determine appropriate label based on power percentage
		var valueLabel string
		if i == 0 && powerPercentage >= 67.0 {
			valueLabel = "SUPERMAJORITY"
		} else if i == 0 && powerPercentage >= 50.0 {
			valueLabel = "MAJORITY"
		} else if i == 0 && (len(valueCounts) == 1 || valueCounts[0].TotalPower > valueCounts[1].TotalPower) {
			valueLabel = "PLURALITY"
		} else {
			valueLabel = "MINORITY"
		}

		fmt.Printf("%s VALUE (Power: %d, %.2f%% of total)\n", valueLabel, vc.TotalPower, powerPercentage)
		fmt.Printf("  %s\n", displayValue)
		fmt.Printf("  Supported by %d validators (%.2f%% of validators)\n",
			len(vc.Validators),
			float64(len(vc.Validators))/float64(len(allValidators))*100)

		// Count validators with and without votes
		votedCount := 0
		for _, v := range vc.Validators {
			if v.HasVote {
				votedCount++
			}
		}

		noVoteCount := len(vc.Validators) - votedCount
		if displayValue != noVoteMsg && noVoteCount > 0 {
			fmt.Printf("  Note: %d validators in this group did not submit a vote extension\n", noVoteCount)
		}

		fmt.Printf("  Validators:\n")
		for j, v := range vc.Validators {
			tag := ""
			if !v.HasVote {
				tag = " (no vote)"
			}
			fmt.Printf("    %3d. %s (%s) (power: %d)%s\n", j+1, v.DecodedAddr, v.Moniker, v.Power, tag)
		}
		fmt.Println()
	}
}

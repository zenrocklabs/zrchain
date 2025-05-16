package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v2"
)

const (
	// BlockIDFlagCommit = "BLOCK_ID_FLAG_COMMIT" // String representation for commit (Removed)
	// BlockIDFlagCommitOld = "2" // Numeric string representation for commit in some versions (Removed)
	ProtoBlockIDFlagCommit = 2 // Standard Tendermint proto value for a commit signature
)

// fetchRPCData performs an HTTP GET request and unmarshals the JSON response
func fetchRPCData(url string, target any) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch data from %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body) // Changed from ioutil.ReadAll
		return fmt.Errorf("failed to fetch data from %s: status %s, body: %s", url, resp.Status, string(bodyBytes))
	}

	body, err := io.ReadAll(resp.Body) // Changed from ioutil.ReadAll
	if err != nil {
		return fmt.Errorf("failed to read response body from %s: %v", url, err)
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("failed to unmarshal JSON from %s: %v. Body: %s", url, err, string(body))
	}
	return nil
}

func main() {
	config := parseFlags()

	// Display selected network info
	fmt.Printf("Using %s network (%s)\n", config.Network, config.RPCNode)

	// Get validator information
	addrToMoniker, err := buildValidatorMappings(config.RPCNode)
	if err != nil {
		fmt.Printf("Warning: Failed to get validator information: %v\n", err)
		fmt.Println("Proceeding without validator names.")
		addrToMoniker = make(map[string]string)
	}

	if config.ConsensusReportMode {
		err := processConsensusReport(config, addrToMoniker)
		if err != nil {
			fmt.Printf("Error generating consensus report: %v\n", err)
			return
		}
	} else {
		// Get block data for vote extension mode
		blockData, err := getBlockData(config)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Process block data for vote extension mode
		processBlockData(blockData, addrToMoniker, config.MissingOnly)
	}
}

// parseFlags parses command line flags and returns a Config
func parseFlags() Config {
	useFileFlag := flag.String("file", "", "Use file instead of executing command (optional)")
	networkFlag := flag.String("network", "mainnet", "Network to use: devnet, testnet, or mainnet (default: mainnet)")
	rpcNodeFlag := flag.String("node", "", "RPC node URL (overrides network selection)")
	blockHeightFlag := flag.String("height", "", "Block height (default: latest)")
	missingOnlyFlag := flag.Bool("missing-only", false, "Only show validators missing vote extensions")
	consensusReportFlag := flag.Bool("consensus-report", false, "Generate a block consensus report (alternative to vote extension analysis)")
	debugFlag := flag.Bool("debug", false, "Enable detailed RPC signature debugging output")
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
		UseFile:             *useFileFlag,
		RPCNode:             rpcURL,
		Network:             *networkFlag,
		BlockHeight:         *blockHeightFlag,
		MissingOnly:         *missingOnlyFlag,
		ConsensusReportMode: *consensusReportFlag,
		DebugMode:           *debugFlag,
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
			if strings.HasPrefix(val.Address, "zenvalcons") {
				addressToMoniker[val.Address] = moniker
			} else {
				fmt.Printf("Warning (buildValidatorMappings): validator address %s from validator-set does not have expected prefix. Moniker will not be mapped.\n", val.Address)
			}
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

// hexAddressToBech32ConsensusAddress converts a Tendermint hex address to a bech32 zenvalcons address
func hexAddressToBech32ConsensusAddress(hexAddr string) (string, error) {
	// Tendermint addresses are typically uppercase hex. Remove "0x" prefix if present.
	hexAddr = strings.TrimPrefix(hexAddr, "0x")
	rawAddrBytes, err := hex.DecodeString(hexAddr)
	if err != nil {
		return "", fmt.Errorf("failed to decode hex address '%s': %v", hexAddr, err)
	}
	// Assuming "zenvalcons" is the correct prefix for your chain's consensus addresses
	return sdk.MustBech32ifyAddressBytes("zenvalcons", rawAddrBytes), nil
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

func processConsensusReport(config Config, addrToMoniker map[string]string) error {
	var targetHeight int64
	var err error
	originalRequestedHeight := config.BlockHeight // Keep track of what user asked for

	// Determine block height
	if config.BlockHeight == "" || config.BlockHeight == "latest" {
		statusURL := fmt.Sprintf("%s/status", config.RPCNode)
		var statusResp struct {
			Result struct {
				SyncInfo struct {
					LatestBlockHeight string `json:"latest_block_height"`
				} `json:"sync_info"`
			} `json:"result"`
		}
		if err := fetchRPCData(statusURL, &statusResp); err != nil {
			return fmt.Errorf("failed to get latest block height from status: %v", err)
		}
		targetHeight, err = strconv.ParseInt(statusResp.Result.SyncInfo.LatestBlockHeight, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse latest block height: %v", err)
		}
		// fmt.Printf("Latest block height: %d\n", targetHeight) // Made quieter
	} else {
		targetHeight, err = strconv.ParseInt(config.BlockHeight, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid block height provided '%s': %v", config.BlockHeight, err)
		}
	}

	if targetHeight <= 0 {
		return fmt.Errorf("block height must be greater than 0")
	}

	// Fetch block data for targetHeight
	// fmt.Printf("Fetching block %d...\n", targetHeight) // Made quieter
	blockURL := fmt.Sprintf("%s/block?height=%d", config.RPCNode, targetHeight)
	var blockResp RPCBlockResponse
	if err := fetchRPCData(blockURL, &blockResp); err != nil {
		return fmt.Errorf("failed to get block data for height %d: %v", targetHeight, err)
	}

	// Fetch block data for targetHeight + 1 to get signatures for targetHeight
	nextHeight := targetHeight + 1
	// fmt.Printf("Fetching block %d for signatures...\n", nextHeight) // Made quieter
	nextBlockURL := fmt.Sprintf("%s/block?height=%d", config.RPCNode, nextHeight)
	var nextBlockResp RPCBlockResponse
	if err := fetchRPCData(nextBlockURL, &nextBlockResp); err != nil {
		isHeightTooHighError := strings.Contains(err.Error(), "must be less than or equal to the current blockchain height")
		requestedLatest := originalRequestedHeight == "" || originalRequestedHeight == "latest"
		// Check if the error confirms targetHeight+1 is indeed too high relative to an identified targetHeight.
		isConfirmedTip := strings.Contains(err.Error(), fmt.Sprintf("height %d must be less than or equal to the current blockchain height %d", nextHeight, targetHeight))

		if isHeightTooHighError && targetHeight > 1 && (requestedLatest || isConfirmedTip) {
			fmt.Printf("Warning: Chain tip reached. Signatures for block %d require block %d, which is not yet available.\n", targetHeight, nextHeight)
			fmt.Printf("Adjusting to report on block %d (signatures from %d).\n", targetHeight-1, targetHeight)
			targetHeight-- // Decrement targetHeight to get consensus for the previous block

			// Re-fetch block data for the new targetHeight (L-1)
			// fmt.Printf("Re-fetching block %d...\n", targetHeight) // Made quieter
			blockURL = fmt.Sprintf("%s/block?height=%d", config.RPCNode, targetHeight)
			if errBlock := fetchRPCData(blockURL, &blockResp); errBlock != nil {
				return fmt.Errorf("failed to get block data for fallback height %d: %v", targetHeight, errBlock)
			}

			// Re-fetch block data for new nextHeight (L-1)+1 = L
			nextHeight = targetHeight + 1
			// fmt.Printf("Re-fetching block %d for signatures...\n", nextHeight) // Made quieter
			nextBlockURL = fmt.Sprintf("%s/block?height=%d", config.RPCNode, nextHeight)
			if errNextBlock := fetchRPCData(nextBlockURL, &nextBlockResp); errNextBlock != nil {
				return fmt.Errorf("failed to get next block data for fallback height %d (expected %d): %v", targetHeight, nextHeight, errNextBlock)
			}
		} else {
			fmt.Printf("Warning: Failed to get block data for height %d (needed for signatures, may be latest block or other issue): %v\n", nextHeight, err)
			nextBlockResp.Result.Block.LastCommit.Signatures = []RPCSignature{}
		}
	}

	fmt.Printf("Processing consensus for H=%d (signatures from H+1=%d, validators for H=%d)...\n", targetHeight, nextHeight, targetHeight)

	// Fetch validator set for targetHeight (which might have been adjusted)
	// Tendermint RPC /validators endpoint can be paginated. Assuming <=100 active validators for now.
	// For more robust solution, handle pagination (max per_page is 100).
	// fmt.Printf("Fetching validators for block %d...\n", targetHeight) // Made quieter
	validatorsURL := fmt.Sprintf("%s/validators?height=%d&per_page=100", config.RPCNode, targetHeight)
	var validatorsResp RPCValidatorsResponse
	if err := fetchRPCData(validatorsURL, &validatorsResp); err != nil {
		return fmt.Errorf("failed to get validators for height %d: %v", targetHeight, err)
	}

	// Process data
	reportData, err := analyzeConsensusData(targetHeight, blockResp, nextBlockResp, validatorsResp, addrToMoniker, config)
	if err != nil {
		return fmt.Errorf("failed to analyze consensus data: %v", err)
	}

	// Display report
	printConsensusReport(reportData)

	return nil
}

// analyzeConsensusData processes fetched data and builds the report structure
func analyzeConsensusData(height int64, currentBlock RPCBlockResponse, nextBlock RPCBlockResponse, rpcValidators RPCValidatorsResponse, addrToMoniker map[string]string, config Config) (*ConsensusReportData, error) {
	proposerHexAddr := currentBlock.Result.Block.Header.ProposerAddress
	proposerBech32AddrForLookup, err := hexAddressToBech32ConsensusAddress(proposerHexAddr)
	var proposerMoniker string
	if err != nil {
		fmt.Printf("Warning (analyzeConsensusData): could not convert proposer hex address %s to bech32 for moniker lookup: %v\n", proposerHexAddr, err)
		proposerMoniker = "Unknown (addr conv err)" // Or try direct hex lookup if that was a fallback: addrToMoniker[proposerHexAddr]
	} else {
		proposerMoniker = addrToMoniker[proposerBech32AddrForLookup]
	}
	if proposerMoniker == "" {
		proposerMoniker = "Unknown"
	}

	report := ConsensusReportData{
		Height:          height,
		AppHash:         currentBlock.Result.Block.Header.AppHash,
		ProposerAddress: proposerHexAddr, // Keep original HEX for this field in the report
		ProposerMoniker: proposerMoniker,
	}

	activeValidators := make(map[string]ValidatorVoteInfo)
	var totalOverallVotingPower int64 = 0

	for _, val := range rpcValidators.Result.Validators {
		votingPower, err := strconv.ParseInt(val.VotingPower, 10, 64)
		if err != nil {
			fmt.Printf("Warning: could not parse voting power '%s' for validator %s: %v\n", val.VotingPower, val.Address, err)
			// Skip validator or assign 0 power if parsing fails
			continue
		}
		valHexAddr := val.Address // Original hex address from /validators RPC
		valBech32Addr, err := hexAddressToBech32ConsensusAddress(valHexAddr)
		if err != nil {
			fmt.Printf("Warning (analyzeConsensusData): could not convert validator hex address %s to bech32: %v\n", valHexAddr, err)
			valBech32Addr = valHexAddr // Use original hex for display if conversion fails
		}

		moniker := addrToMoniker[valBech32Addr] // Look up with Bech32 address
		if moniker == "" {
			moniker = "Unknown"
		}
		activeValidators[valHexAddr] = ValidatorVoteInfo{ // Key activeValidators map with original HEX for direct mapping from signatures
			Address:     valBech32Addr, // Store and display Bech32 address
			Moniker:     moniker,
			VotingPower: votingPower,
		}
		totalOverallVotingPower += votingPower
	}
	report.TotalValidators = len(activeValidators)
	report.TotalVotingPower = totalOverallVotingPower

	signatures := nextBlock.Result.Block.LastCommit.Signatures

	var currentBlockHeightForLog int64 = height
	var nextBlockInfoForLog string = "N/A (next block data not available or error during fetch)"

	if nextBlock.Result.Block.Header.Height != "" { // If Header.Height is populated
		parsedNextHeight, err := strconv.ParseInt(nextBlock.Result.Block.Header.Height, 10, 64)
		if err == nil {
			nextBlockInfoForLog = fmt.Sprintf("%d", parsedNextHeight)
		} else {
			nextBlockInfoForLog = fmt.Sprintf("Error parsing height '%s' (defaulting display)", nextBlock.Result.Block.Header.Height)
		}
	}
	// No specific 'else if' needed here because if Header.Height is empty, the default string is used.

	fmt.Printf("Found %d signatures in commit data (from block: %s) for previous block %d.\n", len(signatures), nextBlockInfoForLog, currentBlockHeightForLog)

	// Store HEX addresses of validators whose signatures are found in the commit
	signaturesPresent := make(map[string]bool)
	rawSignaturesFromRPC := nextBlock.Result.Block.LastCommit.Signatures // Keep original for reference if needed

	if config.DebugMode {
		fmt.Println("--- BEGIN DEBUG: Raw Signatures from RPC Commit Data ---")
		for i, sig := range rawSignaturesFromRPC {
			fmt.Printf("DEBUG_RPC_SIG[%d]: Address=\"%s\", BlockIDFlag=%d, Signature=\"%s\"\n", i, sig.ValidatorAddress, sig.BlockIDFlag, sig.Signature)
		}
		fmt.Println("--- END DEBUG: Raw Signatures from RPC Commit Data ---")
	}

	validSignatures := make([]RPCSignature, 0, len(rawSignaturesFromRPC))
	anomalousEmptyAddrCounts := make(map[int]int)
	totalAnomalousSignatures := 0

	for _, sig := range rawSignaturesFromRPC {
		if sig.ValidatorAddress == "" {
			anomalousEmptyAddrCounts[sig.BlockIDFlag]++
			totalAnomalousSignatures++
		} else {
			validSignatures = append(validSignatures, sig)
		}
	}

	if totalAnomalousSignatures > 0 {
		fmt.Printf("Warning: Found %d anomalous signatures with EMPTY validator addresses in commit data for block %s:\n", totalAnomalousSignatures, nextBlockInfoForLog)
		for flag, count := range anomalousEmptyAddrCounts {
			var flagDesc string
			switch flag {
			case ProtoBlockIDFlagCommit:
				flagDesc = "COMMIT"
			case 1: // Absent
				flagDesc = "ABSENT"
			case 3: // Nil
				flagDesc = "NIL"
			default:
				flagDesc = fmt.Sprintf("UNKNOWN_FLAG_%d", flag)
			}
			fmt.Printf("  - %d signatures marked as %s (BlockIDFlag=%d)\n", count, flagDesc, flag)
		}
		fmt.Println("  This could indicate an issue with data from the RPC node.")
	}

	// Process validSignatures for the report
	for _, sig := range validSignatures {
		signaturesPresent[sig.ValidatorAddress] = true // Mark signature as present using the non-empty address

		if sig.BlockIDFlag == ProtoBlockIDFlagCommit { // Voted COMMIT
			if valInfo, ok := activeValidators[sig.ValidatorAddress]; ok {
				report.AgreedValidators = append(report.AgreedValidators, valInfo)
				report.AgreedVotingPower += valInfo.VotingPower
			} else {
				fmt.Printf("Warning: Validator %s (COMMIT vote) signed but was not found in the active validator set for height %d.\n", sig.ValidatorAddress, height)
			}
		} else if sig.BlockIDFlag == 3 { // Voted NIL (3)
			if valInfo, ok := activeValidators[sig.ValidatorAddress]; ok {
				report.VotedNilValidators = append(report.VotedNilValidators, valInfo)
			} else {
				fmt.Printf("Warning: Validator %s (NIL vote) signed but was not found in the active validator set for height %d.\n", sig.ValidatorAddress, height)
			}
		} else if sig.BlockIDFlag == 1 { // Voted ABSENT (1)
			if valInfo, ok := activeValidators[sig.ValidatorAddress]; ok {
				report.AbsentValidators = append(report.AbsentValidators, valInfo)
			} else {
				fmt.Printf("Warning: Validator %s (ABSENT vote) signed but was not found in the active validator set for height %d.\n", sig.ValidatorAddress, height)
			}
		}
	}

	// Populate MissingSignatureValidators
	for valHexAddr, valInfo := range activeValidators {
		if !signaturesPresent[valHexAddr] { // If no signature of any type was found for this active validator
			report.MissingSignatureValidators = append(report.MissingSignatureValidators, valInfo)
		}
	}

	// Sort validators by voting power (descending)
	sort.Slice(report.AgreedValidators, func(i, j int) bool {
		return report.AgreedValidators[i].VotingPower > report.AgreedValidators[j].VotingPower
	})
	sort.Slice(report.VotedNilValidators, func(i, j int) bool {
		return report.VotedNilValidators[i].VotingPower > report.VotedNilValidators[j].VotingPower
	})
	sort.Slice(report.AbsentValidators, func(i, j int) bool {
		return report.AbsentValidators[i].VotingPower > report.AbsentValidators[j].VotingPower
	})
	sort.Slice(report.MissingSignatureValidators, func(i, j int) bool {
		return report.MissingSignatureValidators[i].VotingPower > report.MissingSignatureValidators[j].VotingPower
	})

	return &report, nil
}

// printConsensusReport displays the formatted consensus report
func printConsensusReport(report *ConsensusReportData) {
	fmt.Printf("\n===== BLOCK CONSENSUS REPORT =====\n")
	fmt.Printf("Block Height: %d\n", report.Height)
	fmt.Printf("App Hash: %s\n", report.AppHash)
	fmt.Printf("Proposer: %s (%s)\n", report.ProposerAddress, report.ProposerMoniker)

	// AGREED/COMMIT section
	fmt.Printf("\n----- VALIDATORS WHO VOTED <COMMIT> (for this block hash) (%d) -----\n", len(report.AgreedValidators))
	agreementValidatorPercentage := 0.0
	if report.TotalValidators > 0 {
		agreementValidatorPercentage = (float64(len(report.AgreedValidators)) / float64(report.TotalValidators)) * 100
	}
	agreementPowerPercentage := 0.0
	if report.TotalVotingPower > 0 {
		agreementPowerPercentage = (float64(report.AgreedVotingPower) / float64(report.TotalVotingPower)) * 100
	}
	fmt.Printf("  Validator Count: %d/%d (%.2f%% of active set)\n", len(report.AgreedValidators), report.TotalValidators, agreementValidatorPercentage)
	fmt.Printf("  Voting Power:    %d/%d (%.2f%% of total)\n", report.AgreedVotingPower, report.TotalVotingPower, agreementPowerPercentage)
	if len(report.AgreedValidators) > 0 {
		for i, v := range report.AgreedValidators {
			fmt.Printf("%3d. %s (%s) (Voting Power: %d)\n", i+1, v.Address, v.Moniker, v.VotingPower)
		}
	} else {
		fmt.Println("No validators explicitly cast a COMMIT vote for this block hash.")
	}

	// Section for NIL votes
	fmt.Printf("\n--- VALIDATORS WHO VOTED <NIL> (abstained/disagreed with proposed block) (%d) ---\n", len(report.VotedNilValidators))
	var nilVotingPower int64
	for _, v := range report.VotedNilValidators {
		nilVotingPower += v.VotingPower
	}
	validatorNilPercentage := 0.0
	if report.TotalValidators > 0 {
		validatorNilPercentage = (float64(len(report.VotedNilValidators)) / float64(report.TotalValidators)) * 100
	}
	powerNilPercentage := 0.0
	if report.TotalVotingPower > 0 {
		powerNilPercentage = (float64(nilVotingPower) / float64(report.TotalVotingPower)) * 100
	}
	fmt.Printf("  Validator Count: %d/%d (%.2f%% of active set)\n", len(report.VotedNilValidators), report.TotalValidators, validatorNilPercentage)
	fmt.Printf("  Voting Power:    %d/%d (%.2f%% of total)\n", nilVotingPower, report.TotalVotingPower, powerNilPercentage)
	if len(report.VotedNilValidators) > 0 {
		for i, v := range report.VotedNilValidators {
			fmt.Printf("%3d. %s (%s) (Voting Power: %d)\n", i+1, v.Address, v.Moniker, v.VotingPower)
		}
	} else {
		fmt.Println("No validators explicitly voted NIL (based on found signatures).")
	}

	// Section for ABSENT votes
	fmt.Printf("\n--- VALIDATORS RECORDED WITH <ABSENT> STATUS (vote not received) (%d) ---\n", len(report.AbsentValidators))
	var absentVotingPower int64
	for _, v := range report.AbsentValidators {
		absentVotingPower += v.VotingPower
	}
	validatorAbsentPercentage := 0.0
	if report.TotalValidators > 0 {
		validatorAbsentPercentage = (float64(len(report.AbsentValidators)) / float64(report.TotalValidators)) * 100
	}
	powerAbsentPercentage := 0.0
	if report.TotalVotingPower > 0 {
		powerAbsentPercentage = (float64(absentVotingPower) / float64(report.TotalVotingPower)) * 100
	}
	fmt.Printf("  Validator Count: %d/%d (%.2f%% of active set)\n", len(report.AbsentValidators), report.TotalValidators, validatorAbsentPercentage)
	fmt.Printf("  Voting Power:    %d/%d (%.2f%% of total)\n", absentVotingPower, report.TotalVotingPower, powerAbsentPercentage)
	if len(report.AbsentValidators) > 0 {
		for i, v := range report.AbsentValidators {
			fmt.Printf("%3d. %s (%s) (Voting Power: %d)\n", i+1, v.Address, v.Moniker, v.VotingPower)
		}
	} else {
		fmt.Println("No validators were recorded with an <ABSENT> status (vote not received).")
	}

	// Section for MISSING SIGNATURES
	fmt.Printf("\n----- VALIDATORS WITH NO SIGNATURE FOUND (in active set) (%d) -----\n", len(report.MissingSignatureValidators))
	var missingSignatureVotingPower int64
	for _, v := range report.MissingSignatureValidators {
		missingSignatureVotingPower += v.VotingPower
	}
	validatorMissingPercentage := 0.0
	if report.TotalValidators > 0 {
		validatorMissingPercentage = (float64(len(report.MissingSignatureValidators)) / float64(report.TotalValidators)) * 100
	}
	powerMissingPercentage := 0.0
	if report.TotalVotingPower > 0 {
		powerMissingPercentage = (float64(missingSignatureVotingPower) / float64(report.TotalVotingPower)) * 100
	}
	fmt.Printf("  Validator Count: %d/%d (%.2f%% of active set)\n", len(report.MissingSignatureValidators), report.TotalValidators, validatorMissingPercentage)
	fmt.Printf("  Voting Power:    %d/%d (%.2f%% of total)\n", missingSignatureVotingPower, report.TotalVotingPower, powerMissingPercentage)
	if len(report.MissingSignatureValidators) > 0 {
		for i, v := range report.MissingSignatureValidators {
			fmt.Printf("%3d. %s (%s) (Voting Power: %d)\n", i+1, v.Address, v.Moniker, v.VotingPower)
		}
	} else {
		fmt.Println("All validators in the active set had a signature in the commit data (either COMMIT, NIL, or ABSENT).")
	}
	fmt.Println()
}

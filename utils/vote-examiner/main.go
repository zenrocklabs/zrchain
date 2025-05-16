package main

import (
	"bytes"
	"encoding/base64"
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

	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v2"
)

const (
	// BlockIDFlagCommit = "BLOCK_ID_FLAG_COMMIT" // String representation for commit (Removed)
	// BlockIDFlagCommitOld = "2" // Numeric string representation for commit in some versions (Removed)
	ProtoBlockIDFlagCommit = 2 // Standard Tendermint proto value for a commit signature
)

// fetchRPCData performs an HTTP GET request and unmarshals the JSON response
func fetchRPCData(url string, target interface{}) error {
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
	fmt.Println("Fetching validator information...")
	addrToMoniker, err := buildValidatorMappings(config.RPCNode)
	if err != nil {
		fmt.Printf("Warning: Failed to get validator information: %v\n", err)
		fmt.Println("Proceeding without validator names.")
		addrToMoniker = make(map[string]string)
	}

	if config.ConsensusReportMode {
		fmt.Println("Generating block consensus report...")
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

func processConsensusReport(config Config, addrToMoniker map[string]string) error {
	var targetHeight int64
	var err error

	// Determine block height
	if config.BlockHeight == "" || config.BlockHeight == "latest" {
		// Fetch status to get the latest block height
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
		fmt.Printf("Latest block height: %d\n", targetHeight)
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
	fmt.Printf("Fetching block %d...\n", targetHeight)
	blockURL := fmt.Sprintf("%s/block?height=%d", config.RPCNode, targetHeight)
	var blockResp RPCBlockResponse
	if err := fetchRPCData(blockURL, &blockResp); err != nil {
		return fmt.Errorf("failed to get block data for height %d: %v", targetHeight, err)
	}

	// Fetch block data for targetHeight + 1 to get signatures for targetHeight
	nextHeight := targetHeight + 1
	fmt.Printf("Fetching block %d for signatures...\n", nextHeight)
	nextBlockURL := fmt.Sprintf("%s/block?height=%d", config.RPCNode, nextHeight)
	var nextBlockResp RPCBlockResponse
	if err := fetchRPCData(nextBlockURL, &nextBlockResp); err != nil {
		// This can fail if targetHeight is the absolute latest. Allow to proceed but warn.
		fmt.Printf("Warning: Failed to get block data for height %d (may be latest block): %v\n", nextHeight, err)
		// Set empty signatures if next block fetch fails
		nextBlockResp.Result.Block.LastCommit.Signatures = []RPCSignature{}
	}

	// Fetch validator set for targetHeight
	// Tendermint RPC /validators endpoint can be paginated. Assuming <=100 active validators for now.
	// For more robust solution, handle pagination (max per_page is 100).
	fmt.Printf("Fetching validators for block %d...\n", targetHeight)
	validatorsURL := fmt.Sprintf("%s/validators?height=%d&per_page=100", config.RPCNode, targetHeight)
	var validatorsResp RPCValidatorsResponse
	if err := fetchRPCData(validatorsURL, &validatorsResp); err != nil {
		return fmt.Errorf("failed to get validators for height %d: %v", targetHeight, err)
	}

	// Process data
	reportData, err := analyzeConsensusData(targetHeight, blockResp, nextBlockResp, validatorsResp, addrToMoniker)
	if err != nil {
		return fmt.Errorf("failed to analyze consensus data: %v", err)
	}

	// Display report
	printConsensusReport(reportData)

	return nil
}

// analyzeConsensusData processes fetched data and builds the report structure
func analyzeConsensusData(height int64, currentBlock RPCBlockResponse, nextBlock RPCBlockResponse, rpcValidators RPCValidatorsResponse, addrToMoniker map[string]string) (*ConsensusReportData, error) {
	report := ConsensusReportData{
		Height:          height,
		AppHash:         currentBlock.Result.Block.Header.AppHash,
		ProposerAddress: currentBlock.Result.Block.Header.ProposerAddress,
		ProposerMoniker: addrToMoniker[currentBlock.Result.Block.Header.ProposerAddress], // Moniker for proposer
	}
	if report.ProposerMoniker == "" {
		report.ProposerMoniker = "Unknown"
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
		moniker := addrToMoniker[val.Address]
		if moniker == "" {
			moniker = "Unknown"
		}
		activeValidators[val.Address] = ValidatorVoteInfo{
			Address:     val.Address,
			Moniker:     moniker,
			VotingPower: votingPower,
		}
		totalOverallVotingPower += votingPower
	}
	report.TotalValidators = len(activeValidators)
	report.TotalVotingPower = totalOverallVotingPower

	signatures := nextBlock.Result.Block.LastCommit.Signatures
	agreedValidatorMap := make(map[string]bool)
	var agreedVotingPower int64 = 0

	for _, sig := range signatures {
		// block_id_flag == 2 (BLOCK_ID_FLAG_COMMIT) means the validator signed the commit
		if sig.BlockIDFlag == ProtoBlockIDFlagCommit {
			agreedValidatorMap[sig.ValidatorAddress] = true
			if valInfo, ok := activeValidators[sig.ValidatorAddress]; ok {
				agreedVotingPower += valInfo.VotingPower
			} else {
				// This case is like the Python script's "WARNING: Validator ... signed but is not in the validators list!"
				// For now, we just note they agreed but don't add to agreed list if not in active set from /validators
				fmt.Printf("Warning: Validator %s signed but was not found in the active validator set for height %d.\n", sig.ValidatorAddress, height)
			}
		}
	}
	report.AgreedVotingPower = agreedVotingPower

	for addr, valInfo := range activeValidators {
		if agreedValidatorMap[addr] {
			report.AgreedValidators = append(report.AgreedValidators, valInfo)
		} else {
			report.DisagreedValidators = append(report.DisagreedValidators, valInfo)
		}
	}

	// Sort validators by voting power (descending)
	sort.Slice(report.AgreedValidators, func(i, j int) bool {
		return report.AgreedValidators[i].VotingPower > report.AgreedValidators[j].VotingPower
	})
	sort.Slice(report.DisagreedValidators, func(i, j int) bool {
		return report.DisagreedValidators[i].VotingPower > report.DisagreedValidators[j].VotingPower
	})

	return &report, nil
}

// printConsensusReport displays the formatted consensus report
func printConsensusReport(report *ConsensusReportData) {
	fmt.Printf("\n===== BLOCK CONSENSUS REPORT =====\n")
	fmt.Printf("Block Height: %d\n", report.Height)
	fmt.Printf("App Hash: %s\n", report.AppHash)
	fmt.Printf("Proposer: %s (%s)\n", report.ProposerAddress, report.ProposerMoniker)

	agreementPercentage := 0.0
	if report.TotalValidators > 0 {
		agreementPercentage = (float64(len(report.AgreedValidators)) / float64(report.TotalValidators)) * 100
	}
	fmt.Printf("Consensus: %d/%d validators agreed (%.2f%%)\n", len(report.AgreedValidators), report.TotalValidators, agreementPercentage)

	votingPowerPercentage := 0.0
	if report.TotalVotingPower > 0 {
		votingPowerPercentage = (float64(report.AgreedVotingPower) / float64(report.TotalVotingPower)) * 100
	}
	fmt.Printf("Consensus by Voting Power: %d/%d (%.2f%%)\n", report.AgreedVotingPower, report.TotalVotingPower, votingPowerPercentage)

	fmt.Printf("\n----- VALIDATORS WHO AGREED (%d) -----\n", len(report.AgreedValidators))
	if len(report.AgreedValidators) > 0 {
		for i, v := range report.AgreedValidators {
			fmt.Printf("%3d. %s (%s) (Voting Power: %d)\n", i+1, v.Address, v.Moniker, v.VotingPower)
		}
	} else {
		fmt.Println("No validators agreed on this block.")
	}

	fmt.Printf("\n----- VALIDATORS WHO DID NOT AGREE (%d) -----\n", len(report.DisagreedValidators))
	if len(report.DisagreedValidators) > 0 {
		for i, v := range report.DisagreedValidators {
			fmt.Printf("%3d. %s (%s) (Voting Power: %d)\n", i+1, v.Address, v.Moniker, v.VotingPower)
		}
	} else {
		fmt.Println("All active validators agreed on this block.")
	}
	fmt.Println()
}

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

// ConsensusVote represents a vote from a validator
type ConsensusVote struct {
	Validator struct {
		Address string `json:"address"`
		Power   int    `json:"power"`
	} `json:"validator"`
	VoteExtension      string `json:"vote_extension"`
	ExtensionSignature string `json:"extension_signature"`
	BlockIDFlag        int    `json:"block_id_flag"`
}

// BlockData represents the block data structure
type BlockData struct {
	ConsensusData struct {
		Votes []ConsensusVote `json:"votes"`
	} `json:"ConsensusData"`
}

// ValidatorInfo stores information about a validator
type ValidatorInfo struct {
	Address     string
	Power       int
	Index       int
	HasVote     bool
	DecodedAddr string
	Moniker     string // Added moniker field
}

// ValidatorSetEntry represents a validator in the validator set
type ValidatorSetEntry struct {
	Address          string `yaml:"address"`
	ProposerPriority string `yaml:"proposer_priority"`
	PubKey           struct {
		Type  string `yaml:"type"`
		Value string `yaml:"value"`
	} `yaml:"pub_key"`
	VotingPower string `yaml:"voting_power"`
}

// ValidatorSetResponse is the response structure for the validator set query
type ValidatorSetResponse struct {
	BlockHeight string              `yaml:"block_height"`
	Validators  []ValidatorSetEntry `yaml:"validators"`
	Pagination  struct {
		Total string `yaml:"total"`
	} `yaml:"pagination"`
}

// ValidatorEntry represents a validator in the validators query
type ValidatorEntry struct {
	ConsensusPublicKey struct {
		Type  string `yaml:"type"`
		Value string `yaml:"value"`
	} `yaml:"consensus_pubkey"`
	Description struct {
		Moniker string `yaml:"moniker"`
		Details string `yaml:"details,omitempty"`
		Website string `yaml:"website,omitempty"`
	} `yaml:"description"`
}

// ValidatorsResponse is the response structure for the validators query
type ValidatorsResponse struct {
	Validators []ValidatorEntry `yaml:"validators"`
	Pagination struct {
		Total string `yaml:"total"`
	} `yaml:"pagination"`
}

const NO_VOTE = "<no vote extension>"

func main() {
	// Command line flags
	useFilePtr := flag.String("file", "", "Use file instead of executing command (optional)")
	rpcNodePtr := flag.String("node", "https://rpc.diamond.zenrocklabs.io:443", "RPC node URL")
	blockHeightPtr := flag.String("height", "", "Block height (default: latest)")
	missingOnlyPtr := flag.Bool("missing-only", false, "Only show validators missing vote extensions")
	flag.Parse()

	var input []byte
	var err error

	// Get validator information first
	fmt.Println("Fetching validator information...")

	// Build a mapping from consensus address to moniker
	addrToMoniker, err := buildValidatorMappings(*rpcNodePtr)
	if err != nil {
		fmt.Printf("Warning: Failed to get validator information: %v\n", err)
		fmt.Println("Proceeding without validator names.")
		addrToMoniker = make(map[string]string)
	}

	// If a file is provided, read from it
	if *useFilePtr != "" {
		input, err = os.ReadFile(*useFilePtr)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}
	} else {
		// Otherwise execute the command
		fmt.Println("Querying latest block data...")

		input, err = executeZenrockdCommand(*rpcNodePtr, *blockHeightPtr)
		if err != nil {
			fmt.Printf("Error executing command: %v\n", err)
			return
		}
	}

	// Convert to string for processing
	inputStr := string(input)

	// Extract the JSON object - find text between { and the last }
	re := regexp.MustCompile(`\{[\s\S]*\}`)
	match := re.FindString(inputStr)
	if match == "" {
		fmt.Println("Could not find a JSON object in the input")
		return
	}

	// Parse the JSON data
	var blockData BlockData
	err = json.Unmarshal([]byte(match), &blockData)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	// Calculate total validator power
	totalPower := 0
	for _, vote := range blockData.ConsensusData.Votes {
		totalPower += vote.Validator.Power
	}

	// Store all decoded extensions for comparison
	allExtensions := make([]map[string]interface{}, 0)
	allValidators := make([]ValidatorInfo, 0)
	validatorsWithExtensions := 0
	totalVotingPower := 0
	totalVotedPower := 0
	missingValidators := make([]ValidatorInfo, 0)

	// Process each vote
	if !*missingOnlyPtr {
		fmt.Printf("Found %d validators in the block\n", len(blockData.ConsensusData.Votes))
	}

	for i, vote := range blockData.ConsensusData.Votes {
		validatorHasVote := false

		// Decode the base64 address
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

		if !*missingOnlyPtr {
			fmt.Printf("\n=== Validator %s (%s) (Power: %d) ===\n", decodedAddr, moniker, vote.Validator.Power)
		}

		// Skip if vote extension is empty
		if vote.VoteExtension == "" {
			if !*missingOnlyPtr {
				fmt.Println("No vote extension found")
			}
			missingValidators = append(missingValidators, validatorInfo)
		} else {
			validatorHasVote = true
			validatorsWithExtensions++
			totalVotedPower += vote.Validator.Power

			if !*missingOnlyPtr {
				// Decode base64 data
				decodedBytes, err := base64.StdEncoding.DecodeString(vote.VoteExtension)
				if err != nil {
					fmt.Printf("Error decoding base64: %v\n", err)
				} else {
					// Parse and pretty print the JSON
					var prettyJSON map[string]interface{}
					err = json.Unmarshal(decodedBytes, &prettyJSON)
					if err != nil {
						fmt.Printf("Error parsing decoded JSON: %v\n", err)
						fmt.Println("Raw decoded data:", string(decodedBytes))
					} else {
						// Format the JSON with indentation
						prettyJSONBytes, err := json.MarshalIndent(prettyJSON, "", "  ")
						if err != nil {
							fmt.Printf("Error formatting JSON: %v\n", err)
						} else {
							fmt.Println("Decoded vote extension:")
							fmt.Println(string(prettyJSONBytes))

							// Store extension for comparison
							allExtensions = append(allExtensions, prettyJSON)
						}
					}
				}
			} else {
				// Still need to decode for comparison in missing-only mode
				decodedBytes, err := base64.StdEncoding.DecodeString(vote.VoteExtension)
				if err == nil {
					var prettyJSON map[string]interface{}
					if json.Unmarshal(decodedBytes, &prettyJSON) == nil {
						allExtensions = append(allExtensions, prettyJSON)
					}
				}
			}
		}

		validatorInfo.HasVote = validatorHasVote
		allValidators = append(allValidators, validatorInfo)

		// If no extension, add a nil map to maintain index alignment
		if !validatorHasVote {
			allExtensions = append(allExtensions, nil)
		}
	}

	// Show vote extension participation stats
	fmt.Printf("\n=== VOTE EXTENSION PARTICIPATION ===\n")
	fmt.Printf("Validators with vote extensions: %d of %d (%.2f%%)\n",
		validatorsWithExtensions,
		len(blockData.ConsensusData.Votes),
		float64(validatorsWithExtensions)/float64(len(blockData.ConsensusData.Votes))*100)
	fmt.Printf("Voting power with extensions: %d of %d (%.2f%%)\n",
		totalVotedPower,
		totalVotingPower,
		float64(totalVotedPower)/float64(totalVotingPower)*100)

	// If missing-only flag is set, just list the missing validators
	if *missingOnlyPtr {
		fmt.Printf("\n=== VALIDATORS WITHOUT VOTE EXTENSIONS (%d validators) ===\n", len(missingValidators))
		if len(missingValidators) == 0 {
			fmt.Println("All validators submitted vote extensions!")
		} else {
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
		return
	}

	// Compare extensions to find differences
	if len(allExtensions) > 0 {
		fmt.Println("\n\n=== DIFFERENCES BETWEEN VOTE EXTENSIONS ===")
		findDifferences(allExtensions, allValidators, totalVotingPower)
	}
}

// buildValidatorMappings builds a mapping from bech32 consensus address to moniker
func buildValidatorMappings(rpcNode string) (map[string]string, error) {
	// Maps to store the results
	addressToMoniker := make(map[string]string)
	pubkeyToMoniker := make(map[string]string)

	// Execute command to get validator set
	fmt.Println("Fetching validator set...")
	valSetCmd := exec.Command("bash", "-c", fmt.Sprintf("zenrockd --node=%s q consensus comet validator-set", rpcNode))
	valSetOutput, err := valSetCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get validator set: %v", err)
	}

	// Execute command to get validators with monikers
	fmt.Println("Fetching validator details...")
	validatorsCmd := exec.Command("bash", "-c", fmt.Sprintf("zenrockd --node=%s q validation validators", rpcNode))
	validatorsOutput, err := validatorsCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get validators: %v", err)
	}

	// Parse validator set response
	var valSetResp ValidatorSetResponse
	err = yaml.Unmarshal(valSetOutput, &valSetResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse validator set: %v", err)
	}

	// Parse validators response
	var validatorsResp ValidatorsResponse
	err = yaml.Unmarshal(validatorsOutput, &validatorsResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse validators: %v", err)
	}

	// Build mapping from pubkey to moniker
	for _, val := range validatorsResp.Validators {
		pubkeyToMoniker[val.ConsensusPublicKey.Value] = val.Description.Moniker
	}

	// Build mapping from address to moniker
	for _, val := range valSetResp.Validators {
		if moniker, ok := pubkeyToMoniker[val.PubKey.Value]; ok {
			addressToMoniker[val.Address] = moniker
		}
	}

	fmt.Printf("Found %d validators with monikers\n", len(addressToMoniker))
	return addressToMoniker, nil
}

// executeZenrockdCommand runs the zenrockd command and returns the output
func executeZenrockdCommand(rpcNode, blockHeight string) ([]byte, error) {
	var cmd *exec.Cmd

	if blockHeight == "" {
		// Use bash to execute the full command pipeline
		fullCmd := fmt.Sprintf(`zenrockd --node=%s query block --type=height $(zenrockd --node=%s status | jq -r '.sync_info.latest_block_height') --output=json | jq -r '.data.txs[0]' | base64 --decode`, rpcNode, rpcNode)
		cmd = exec.Command("bash", "-c", fullCmd)
	} else {
		// Use specific block height
		fullCmd := fmt.Sprintf(`zenrockd --node=%s query block --type=height %s --output=json | jq -r '.data.txs[0]' | base64 --decode`, rpcNode, blockHeight)
		cmd = exec.Command("bash", "-c", fullCmd)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("command execution failed: %s\nStderr: %s", err, stderr.String())
	}

	return stdout.Bytes(), nil
}

// Decode validator address from base64
func decodeValidatorAddress(base64Addr string) string {
	// Decode from base64
	decoded, err := base64.StdEncoding.DecodeString(base64Addr)
	if err != nil {
		return fmt.Sprintf("Invalid(%s)", base64Addr)
	}

	// Format as hex
	return sdk.MustBech32ifyAddressBytes("zenvalcons", decoded)
}

func findDifferences(extensions []map[string]interface{}, validators []ValidatorInfo, totalPower int) {
	// Map to track differences by field
	differences := make(map[string]map[string][]ValidatorInfo)

	// Get all possible keys from all extensions
	allKeys := make(map[string]bool)
	for _, ext := range extensions {
		if ext == nil {
			continue
		}
		for k := range ext {
			allKeys[k] = true
		}
	}

	// Check each key for differences
	for key := range allKeys {
		valueMap := make(map[string][]ValidatorInfo)

		for i, ext := range extensions {
			// Skip if index is out of range
			if i >= len(validators) {
				continue
			}

			var valueStr string
			if ext == nil {
				valueStr = NO_VOTE
			} else if val, ok := ext[key]; ok {
				// Convert value to string representation
				bytes, err := json.Marshal(val)
				if err == nil {
					valueStr = string(bytes)
				} else {
					valueStr = fmt.Sprintf("%v", val)
				}
			} else {
				valueStr = "<missing field>"
			}

			// Add validator to the list for this value
			valueMap[valueStr] = append(valueMap[valueStr], validators[i])
		}

		// If there's more than one unique value, we have a difference
		if len(valueMap) > 1 {
			differences[key] = valueMap
		}
	}

	// Display differences in a nice format
	if len(differences) == 0 {
		fmt.Println("No differences found between vote extensions")
		return
	}

	// Sort keys for consistent output
	sortedKeys := make([]string, 0, len(differences))
	for k := range differences {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		fmt.Printf("\nDifferences in field: %s\n", key)
		fmt.Printf("-------------------------\n")

		valueMap := differences[key]

		// Sort values by number of validators (descending) and power
		type ValueCount struct {
			Value      string
			Validators []ValidatorInfo
			TotalPower int
		}

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
			if displayValue == NO_VOTE {
				displayValue = "NO VOTE EXTENSION SUBMITTED"
			} else if len(displayValue) > 100 {
				displayValue = displayValue[:97] + "..."
			}

			valueLabel := "MAJORITY"
			if i > 0 {
				valueLabel = "MINORITY"
			}

			powerPercentage := float64(vc.TotalPower) / float64(totalPower) * 100

			fmt.Printf("%s VALUE (Power: %d, %.2f%% of total)\n", valueLabel, vc.TotalPower, powerPercentage)
			fmt.Printf("  %s\n", displayValue)
			fmt.Printf("  Supported by %d validators (%.2f%% of validators)\n",
				len(vc.Validators),
				float64(len(vc.Validators))/float64(len(validators))*100)

			// Count validators with and without votes
			votedCount := 0
			for _, v := range vc.Validators {
				if v.HasVote {
					votedCount++
				}
			}

			noVoteCount := len(vc.Validators) - votedCount
			if displayValue != "NO VOTE EXTENSION SUBMITTED" && noVoteCount > 0 {
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
}
